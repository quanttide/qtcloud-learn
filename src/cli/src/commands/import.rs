use std::fs;
use std::path::Path;

use quanttide_learn::{Schedule, Task};
use serde::Deserialize;
use serde::Serialize;
use serde_json::Value;

use crate::api::ApiClient;

#[derive(Debug, Deserialize)]
struct SeedData {
    tasks: Vec<Task>,
    schedules: Vec<Schedule>,
}

pub fn run(api: &ApiClient, path: &Path) -> Result<String, String> {
    let raw = fs::read_to_string(path).map_err(|e| format!("读取种子数据失败: {e}"))?;
    let seed: SeedData =
        serde_json::from_str(&raw).map_err(|e| format!("解析种子数据失败: {e}"))?;

    let mut created = 0usize;
    let mut updated = 0usize;

    for task in &seed.tasks {
        if upsert(api, "tasks", &task.id, task)? {
            created += 1;
        } else {
            updated += 1;
        }
    }
    for schedule in &seed.schedules {
        if upsert(api, "schedules", &schedule.id, schedule)? {
            created += 1;
        } else {
            updated += 1;
        }
    }

    Ok(format!(
        "导入完成：tasks {} 个，schedules {} 个；新增 {}，更新 {}",
        seed.tasks.len(),
        seed.schedules.len(),
        created,
        updated
    ))
}

fn upsert<T>(api: &ApiClient, resource: &str, id: &str, entity: &T) -> Result<bool, String>
where
    T: Serialize,
{
    let body = serde_json::to_value(entity).map_err(|e| format!("序列化失败: {e}"))?;
    match api.get(&format!("{resource}/{id}")) {
        Ok(_) => {
            api.put(&format!("{resource}/{id}"), &body)?;
            Ok(false)
        }
        Err(e) if e.contains("HTTP 404") || e.contains("status code 404") => {
            let mut create_body = body;
            if let Value::Object(obj) = &mut create_body {
                obj.insert("id".to_string(), Value::String(id.to_string()));
            }
            api.post(resource, &create_body)?;
            Ok(true)
        }
        Err(e) => Err(e),
    }
}
