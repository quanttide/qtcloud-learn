use clap::Subcommand;
use serde_json::json;

use crate::api::ApiClient;

/// 完成记录子命令（对齐《量潮学习管理标准》Completion 实体：
/// learner_id / criterion_id / status / created_at / updated_at）。
#[derive(Subcommand)]
pub enum CompletionCmd {
    /// 创建完成记录
    Create {
        /// 学习者 ID
        #[arg(long)]
        learner_id: String,
        /// 验收标准 ID
        #[arg(long)]
        criterion_id: String,
        /// 通过状态：completed / not_completed（缺省 not_completed）
        #[arg(long)]
        status: Option<String>,
    },
    /// 标记完成（局部更新 status → completed）
    Complete {
        /// 完成记录 ID
        id: String,
    },
    /// 列出完成记录
    List,
    /// 查看完成记录详情
    Get {
        /// 完成记录 ID
        id: String,
    },
}

pub fn run(api: &ApiClient, cmd: CompletionCmd) -> Result<String, String> {
    match cmd {
        CompletionCmd::Create {
            learner_id,
            criterion_id,
            status,
        } => {
            let mut body = json!({
                "learner_id": learner_id,
                "criterion_id": criterion_id,
            });
            if let Some(s) = status {
                body["status"] = json!(s);
            }
            let v = api.post("completions", &body)?;
            Ok(format!(
                "已创建完成记录 {}（{} → {}，{}）",
                v["id"].as_str().unwrap_or(""),
                v["learner_id"].as_str().unwrap_or(""),
                v["criterion_id"].as_str().unwrap_or(""),
                v["status"].as_str().unwrap_or("")
            ))
        }
        CompletionCmd::Complete { id } => {
            let v = api.put(
                &format!("completions/{id}"),
                &json!({ "status": "completed" }),
            )?;
            Ok(format!(
                "已完成 {}（{} → {}）",
                v["id"].as_str().unwrap_or(""),
                v["learner_id"].as_str().unwrap_or(""),
                v["criterion_id"].as_str().unwrap_or("")
            ))
        }
        CompletionCmd::List => {
            let list = api.get("completions")?;
            let arr = list.as_array().ok_or("响应不是数组")?;
            if arr.is_empty() {
                return Ok("暂无完成记录".to_string());
            }
            let mut out = String::from("ID\t学习者\t验收标准\t状态\t创建时间\n");
            for v in arr {
                out.push_str(&format!(
                    "{}\t{}\t{}\t{}\t{}\n",
                    v["id"].as_str().unwrap_or(""),
                    v["learner_id"].as_str().unwrap_or(""),
                    v["criterion_id"].as_str().unwrap_or(""),
                    v["status"].as_str().unwrap_or(""),
                    v["created_at"].as_str().unwrap_or("")
                ));
            }
            Ok(out.trim_end().to_string())
        }
        CompletionCmd::Get { id } => {
            let v = api.get(&format!("completions/{id}"))?;
            Ok(format!(
                "完成记录 {}：{} → {}\n状态: {}\n创建: {}\n更新: {}",
                v["id"].as_str().unwrap_or(""),
                v["learner_id"].as_str().unwrap_or(""),
                v["criterion_id"].as_str().unwrap_or(""),
                v["status"].as_str().unwrap_or(""),
                v["created_at"].as_str().unwrap_or(""),
                v["updated_at"].as_str().unwrap_or("")
            ))
        }
    }
}
