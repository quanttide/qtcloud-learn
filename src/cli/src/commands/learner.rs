use clap::Subcommand;
use serde_json::json;

use crate::api::ApiClient;

/// 学习者子命令（对齐《量潮学习管理标准》Learner 实体：id / user_id）。
#[derive(Subcommand)]
pub enum LearnerCmd {
    /// 创建学习者
    Create {
        /// 关联 auth 领域的用户 ID（预留）
        #[arg(long)]
        user_id: Option<String>,
    },
    /// 列出学习者
    List,
    /// 查看学习者详情
    Get {
        /// 学习者 ID
        id: String,
    },
}

pub fn run(api: &ApiClient, cmd: LearnerCmd) -> Result<String, String> {
    match cmd {
        LearnerCmd::Create { user_id } => {
            let mut body = json!({});
            if let Some(u) = user_id {
                body["user_id"] = json!(u);
            }
            let v = api.post("learners", &body)?;
            Ok(format!(
                "已创建学习者 {}（user_id: {}）",
                v["id"].as_str().unwrap_or(""),
                v["user_id"].as_str().unwrap_or("-")
            ))
        }
        LearnerCmd::List => {
            let list = api.get("learners")?;
            let arr = list.as_array().ok_or("响应不是数组")?;
            if arr.is_empty() {
                return Ok("暂无学习者".to_string());
            }
            let mut out = String::from("ID\tuser_id\n");
            for v in arr {
                out.push_str(&format!(
                    "{}\t{}\n",
                    v["id"].as_str().unwrap_or(""),
                    v["user_id"].as_str().unwrap_or("")
                ));
            }
            Ok(out.trim_end().to_string())
        }
        LearnerCmd::Get { id } => {
            let v = api.get(&format!("learners/{id}"))?;
            Ok(format!(
                "学习者 {}（user_id: {}）",
                v["id"].as_str().unwrap_or(""),
                v["user_id"].as_str().unwrap_or("")
            ))
        }
    }
}
