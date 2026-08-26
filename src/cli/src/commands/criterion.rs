use clap::Subcommand;
use serde_json::json;

use crate::api::ApiClient;

/// 验收标准子命令（对齐《量潮学习管理标准》Criterion 实体：id / title / description）。
#[derive(Subcommand)]
pub enum CriterionCmd {
    /// 创建验收标准
    Create {
        /// 语义标识，如 vibe-coding/lesson1/zed-connection
        #[arg(long)]
        title: String,
        /// 具体规则描述
        #[arg(long)]
        description: String,
    },
    /// 列出验收标准
    List,
    /// 查看验收标准详情
    Get {
        /// 验收标准 ID
        id: String,
    },
}

pub fn run(api: &ApiClient, cmd: CriterionCmd) -> Result<String, String> {
    match cmd {
        CriterionCmd::Create { title, description } => {
            let body = json!({
                "title": title,
                "description": description,
            });
            let v = api.post("criteria", &body)?;
            Ok(format!(
                "已创建验收标准 {}（{}）",
                v["id"].as_str().unwrap_or(""),
                v["title"].as_str().unwrap_or("")
            ))
        }
        CriterionCmd::List => {
            let list = api.get("criteria")?;
            let arr = list.as_array().ok_or("响应不是数组")?;
            if arr.is_empty() {
                return Ok("暂无验收标准".to_string());
            }
            let mut out = String::from("ID\t标题\t描述\n");
            for v in arr {
                out.push_str(&format!(
                    "{}\t{}\t{}\n",
                    v["id"].as_str().unwrap_or(""),
                    v["title"].as_str().unwrap_or(""),
                    v["description"].as_str().unwrap_or("")
                ));
            }
            Ok(out.trim_end().to_string())
        }
        CriterionCmd::Get { id } => {
            let v = api.get(&format!("criteria/{id}"))?;
            Ok(format!(
                "验收标准 {}：{}\n描述: {}",
                v["id"].as_str().unwrap_or(""),
                v["title"].as_str().unwrap_or(""),
                v["description"].as_str().unwrap_or("")
            ))
        }
    }
}
