use clap::Subcommand;
use serde_json::json;

use crate::api::ApiClient;

/// 班级子命令（承接 `qtcloud-course` cli 规划的 class 能力，自其 provider 移植的模型）。
#[derive(Subcommand)]
pub enum ClassCmd {
    /// 创建班级
    Create {
        /// 班级名称
        #[arg(long)]
        name: String,
        /// 引用名称（展示用）
        #[arg(long)]
        ref_name: String,
        /// 引用 ID（Program/Course）
        #[arg(long)]
        ref_id: String,
        /// 引用类型：program / course
        #[arg(long, default_value = "program")]
        ref_type: String,
        /// 开始日期（ISO）
        #[arg(long)]
        start_date: Option<String>,
        /// 结束日期（ISO）
        #[arg(long)]
        end_date: Option<String>,
    },
    /// 列出班级
    List,
    /// 查看班级详情
    Get {
        /// 班级 ID
        id: String,
    },
}

pub fn run(api: &ApiClient, cmd: ClassCmd) -> Result<String, String> {
    match cmd {
        ClassCmd::Create {
            name,
            ref_name,
            ref_id,
            ref_type,
            start_date,
            end_date,
        } => {
            let mut body = json!({
                "name": name,
                "refName": ref_name,
                "refId": ref_id,
                "refType": ref_type,
            });
            if let Some(d) = start_date {
                body["startDate"] = json!(d);
            }
            if let Some(d) = end_date {
                body["endDate"] = json!(d);
            }
            let v = api.post("classes", &body)?;
            Ok(format!(
                "已创建班级 {}（{}）",
                v["id"].as_str().unwrap_or(""),
                v["name"].as_str().unwrap_or("")
            ))
        }
        ClassCmd::List => {
            let list = api.get("classes")?;
            let arr = list.as_array().ok_or("响应不是数组")?;
            if arr.is_empty() {
                return Ok("暂无班级".to_string());
            }
            let mut out = String::from("ID\t名称\t引用\t状态\t学员数\t进度\n");
            for v in arr {
                out.push_str(&format!(
                    "{}\t{}\t{}\t{}\t{}\t{:.0}%\n",
                    v["id"].as_str().unwrap_or(""),
                    v["name"].as_str().unwrap_or(""),
                    v["refName"].as_str().unwrap_or(""),
                    v["status"].as_str().unwrap_or("preparing"),
                    v["studentCount"].as_i64().unwrap_or(0),
                    v["progress"].as_f64().unwrap_or(0.0) * 100.0
                ));
            }
            Ok(out.trim_end().to_string())
        }
        ClassCmd::Get { id } => {
            let v = api.get(&format!("classes/{id}"))?;
            Ok(format!(
                "班级 {}：{}\n引用: {}（{}，{}）\n状态: {}\n学段: {} - {}\n学员数: {}\n进度: {:.0}%",
                v["id"].as_str().unwrap_or(""),
                v["name"].as_str().unwrap_or(""),
                v["refName"].as_str().unwrap_or(""),
                v["refType"].as_str().unwrap_or("program"),
                v["refId"].as_str().unwrap_or(""),
                v["status"].as_str().unwrap_or("preparing"),
                v["startDate"].as_str().unwrap_or(""),
                v["endDate"].as_str().unwrap_or(""),
                v["studentCount"].as_i64().unwrap_or(0),
                v["progress"].as_f64().unwrap_or(0.0) * 100.0
            ))
        }
    }
}
