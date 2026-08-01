use clap::Subcommand;
use serde_json::json;

use crate::api::ApiClient;

/// 学习进度子命令（学习记录由 localStorage 迁移为服务端进度数据）。
#[derive(Subcommand)]
pub enum ProgressCmd {
    /// 上报进度
    Report {
        /// 班级 ID
        #[arg(long)]
        class_id: String,
        /// 学员 ID
        #[arg(long)]
        student_id: String,
        /// 进度（0.0 ~ 1.0）
        #[arg(long, default_value_t = 0.0)]
        percent: f64,
        /// 标记完成
        #[arg(long)]
        finished: bool,
    },
    /// 查看进度
    Get {
        /// 班级 ID
        #[arg(long)]
        class_id: String,
        /// 学员 ID
        #[arg(long)]
        student_id: String,
    },
}

pub fn run(api: &ApiClient, cmd: ProgressCmd) -> Result<String, String> {
    match cmd {
        ProgressCmd::Report {
            class_id,
            student_id,
            percent,
            finished,
        } => {
            let body = json!({
                "studentId": student_id,
                "classId": class_id,
                "percent": percent,
                "finished": finished,
            });
            match find(api, &class_id, &student_id)? {
                Some(existing) => {
                    let id = existing["id"].as_str().ok_or("进度记录缺 id")?.to_string();
                    api.put(&format!("progress/{id}"), &body)?;
                    Ok(format!(
                        "已更新进度 {id}：{:.0}%{}",
                        percent * 100.0,
                        if finished { "（已完成）" } else { "" }
                    ))
                }
                None => {
                    let v = api.post("progress", &body)?;
                    Ok(format!(
                        "已上报进度 {}：{:.0}%{}",
                        v["id"].as_str().unwrap_or(""),
                        percent * 100.0,
                        if finished { "（已完成）" } else { "" }
                    ))
                }
            }
        }
        ProgressCmd::Get {
            class_id,
            student_id,
        } => match find(api, &class_id, &student_id)? {
            Some(v) => Ok(format!(
                "进度 {}：{:.0}%{}\n最近更新: {}",
                v["id"].as_str().unwrap_or(""),
                v["percent"].as_f64().unwrap_or(0.0) * 100.0,
                if v["finished"].as_bool().unwrap_or(false) {
                    "（已完成）"
                } else {
                    ""
                },
                v["updatedAt"].as_str().unwrap_or("")
            )),
            None => Ok("该学员在此班级暂无进度记录".to_string()),
        },
    }
}

/// 按班级 + 学员查找进度记录。
fn find(
    api: &ApiClient,
    class_id: &str,
    student_id: &str,
) -> Result<Option<serde_json::Value>, String> {
    let list = api.get("progress")?;
    let arr = list.as_array().ok_or("响应不是数组")?;
    Ok(arr
        .iter()
        .find(|v| {
            v["classId"].as_str() == Some(class_id)
                && v["studentId"].as_str() == Some(student_id)
        })
        .cloned())
}
