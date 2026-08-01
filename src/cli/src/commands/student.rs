use clap::Subcommand;
use serde_json::json;

use crate::api::ApiClient;

/// 学员子命令（承接 `qtcloud-course` cli 规划的 student 能力）。
#[derive(Subcommand)]
pub enum StudentCmd {
    /// 创建学员
    Create {
        /// 姓名
        #[arg(long)]
        name: String,
        /// 邮箱
        #[arg(long)]
        email: Option<String>,
        /// 权益计划：free / paid / vip
        #[arg(long)]
        plan: Option<String>,
    },
    /// 列出学员
    List,
    /// 查看学员详情
    Get {
        /// 学员 ID
        id: String,
    },
}

pub fn run(api: &ApiClient, cmd: StudentCmd) -> Result<String, String> {
    match cmd {
        StudentCmd::Create { name, email, plan } => {
            let mut body = json!({ "name": name });
            if let Some(e) = email {
                body["email"] = json!(e);
            }
            if let Some(p) = plan {
                body["plan"] = json!(p);
            }
            let v = api.post("students", &body)?;
            Ok(format!(
                "已创建学员 {}（{}，{}）",
                v["id"].as_str().unwrap_or(""),
                v["name"].as_str().unwrap_or(""),
                v["plan"].as_str().unwrap_or("free")
            ))
        }
        StudentCmd::List => {
            let list = api.get("students")?;
            let arr = list.as_array().ok_or("响应不是数组")?;
            if arr.is_empty() {
                return Ok("暂无学员".to_string());
            }
            let mut out = String::from("ID\t姓名\t邮箱\t计划\n");
            for v in arr {
                out.push_str(&format!(
                    "{}\t{}\t{}\t{}\n",
                    v["id"].as_str().unwrap_or(""),
                    v["name"].as_str().unwrap_or(""),
                    v["email"].as_str().unwrap_or(""),
                    v["plan"].as_str().unwrap_or("free")
                ));
            }
            Ok(out.trim_end().to_string())
        }
        StudentCmd::Get { id } => {
            let v = api.get(&format!("students/{id}"))?;
            Ok(format!(
                "学员 {}：{}（{}）\n邮箱: {}\n计划: {}",
                v["id"].as_str().unwrap_or(""),
                v["name"].as_str().unwrap_or(""),
                v["email"].as_str().unwrap_or(""),
                v["email"].as_str().unwrap_or(""),
                v["plan"].as_str().unwrap_or("free")
            ))
        }
    }
}
