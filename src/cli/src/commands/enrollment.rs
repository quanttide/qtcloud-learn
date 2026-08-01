use clap::Subcommand;
use serde_json::json;

use crate::api::ApiClient;

/// 选课 / 报名子命令。
#[derive(Subcommand)]
pub enum EnrollmentCmd {
    /// 选课（报名）
    Enroll {
        /// 班级 ID
        #[arg(long)]
        class_id: String,
        /// 学员 ID
        #[arg(long)]
        student_id: String,
    },
    /// 退课
    Withdraw {
        /// 班级 ID
        #[arg(long)]
        class_id: String,
        /// 学员 ID
        #[arg(long)]
        student_id: String,
    },
    /// 列出选课记录（可按学员过滤）
    List {
        /// 学员 ID（可选过滤）
        #[arg(long)]
        student_id: Option<String>,
    },
}

pub fn run(api: &ApiClient, cmd: EnrollmentCmd) -> Result<String, String> {
    match cmd {
        EnrollmentCmd::Enroll {
            class_id,
            student_id,
        } => {
            // 幂等：已选课时直接提示
            if let Some(existing) = find(api, &class_id, &student_id)? {
                if existing["status"].as_str() == Some("enrolled") {
                    return Ok(format!("已选课（{}），无需重复报名", existing["id"].as_str().unwrap_or("")));
                }
            }
            let body = json!({
                "classId": class_id,
                "studentId": student_id,
                "status": "enrolled",
                "enrolledAt": chrono_now(),
            });
            let v = api.post("enrollments", &body)?;
            Ok(format!(
                "已报名 {} 选课 {}",
                v["id"].as_str().unwrap_or(""),
                v["classId"].as_str().unwrap_or("")
            ))
        }
        EnrollmentCmd::Withdraw {
            class_id,
            student_id,
        } => {
            let existing = find(api, &class_id, &student_id)?
                .ok_or("未找到该选课记录")?;
            let id = existing["id"].as_str().ok_or("选课记录缺 id")?.to_string();
            let body = json!({
                "classId": class_id,
                "studentId": student_id,
                "status": "withdrawn",
            });
            api.put(&format!("enrollments/{id}"), &body)?;
            Ok(format!("已退课（{id}）"))
        }
        EnrollmentCmd::List { student_id } => {
            let list = api.get("enrollments")?;
            let arr = list.as_array().ok_or("响应不是数组")?;
            let filtered: Vec<_> = arr
                .iter()
                .filter(|v| match &student_id {
                    Some(sid) => v["studentId"].as_str() == Some(sid.as_str()),
                    None => true,
                })
                .collect();
            if filtered.is_empty() {
                return Ok("暂无选课记录".to_string());
            }
            let mut out = String::from("ID\t班级\t学员\t状态\n");
            for v in filtered {
                out.push_str(&format!(
                    "{}\t{}\t{}\t{}\n",
                    v["id"].as_str().unwrap_or(""),
                    v["classId"].as_str().unwrap_or(""),
                    v["studentId"].as_str().unwrap_or(""),
                    v["status"].as_str().unwrap_or("enrolled")
                ));
            }
            Ok(out.trim_end().to_string())
        }
    }
}

/// 按班级 + 学员查找选课记录。
fn find(api: &ApiClient, class_id: &str, student_id: &str) -> Result<Option<serde_json::Value>, String> {
    let list = api.get("enrollments")?;
    let arr = list.as_array().ok_or("响应不是数组")?;
    Ok(arr
        .iter()
        .find(|v| {
            v["classId"].as_str() == Some(class_id)
                && v["studentId"].as_str() == Some(student_id)
        })
        .cloned())
}

/// 当前时间 ISO8601（避免额外 chrono 依赖，使用系统时间格式化）。
fn chrono_now() -> String {
    use std::time::{SystemTime, UNIX_EPOCH};
    let secs = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map(|d| d.as_secs())
        .unwrap_or(0);
    // 格式：YYYY-MM-DDTHH:MM:SSZ（近似 UTC）
    let days = secs / 86400;
    let rem = secs % 86400;
    let (h, m, s) = (rem / 3600, (rem % 3600) / 60, rem % 60);
    let (y, mo, d) = civil_from_days(days as i64);
    format!("{y:04}-{mo:02}-{d:02}T{h:02}:{m:02}:{s:02}Z")
}

/// 天数（自 1970-01-01）转公历日期。
fn civil_from_days(z: i64) -> (i64, i64, i64) {
    let z = z + 719468;
    let era = if z >= 0 { z } else { z - 146096 } / 146097;
    let doe = z - era * 146097;
    let yoe = (doe - doe / 1460 + doe / 36524 - doe / 146096) / 365;
    let y = yoe + era * 400;
    let doy = doe - (365 * yoe + yoe / 4 - yoe / 100);
    let mp = (5 * doy + 2) / 153;
    let d = doy - (153 * mp + 2) / 5 + 1;
    let m = if mp < 10 { mp + 3 } else { mp - 9 };
    (if m <= 2 { y + 1 } else { y }, m, d)
}
