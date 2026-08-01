use clap::Subcommand;
use serde_json::json;

use crate::api::ApiClient;

/// 考核子命令（承接 `qtcloud-course` cli 规划的 assessment 能力）。
#[derive(Subcommand)]
pub enum AssessmentCmd {
    /// 创建考核
    Create {
        /// 班级 ID
        #[arg(long)]
        class_id: String,
        /// 标题
        #[arg(long)]
        title: String,
        /// 类型：homework / exam
        #[arg(long = "type", default_value = "homework")]
        kind: String,
        /// 满分
        #[arg(long, default_value_t = 100)]
        full_score: i64,
        /// 及格线
        #[arg(long, default_value_t = 60)]
        pass_score: i64,
        /// 截止日期（ISO）
        #[arg(long)]
        deadline: Option<String>,
    },
    /// 学员提交
    Submit {
        /// 考核 ID
        #[arg(long)]
        assessment_id: String,
        /// 学员 ID
        #[arg(long)]
        student_id: String,
    },
    /// 评分
    Grade {
        /// 提交 ID
        #[arg(long)]
        submission_id: String,
        /// 得分
        #[arg(long)]
        score: f64,
        /// 评语
        #[arg(long)]
        comment: Option<String>,
    },
    /// 考核统计
    Stats {
        /// 班级 ID（可选过滤）
        #[arg(long)]
        class_id: Option<String>,
    },
}

pub fn run(api: &ApiClient, cmd: AssessmentCmd) -> Result<String, String> {
    match cmd {
        AssessmentCmd::Create {
            class_id,
            title,
            kind,
            full_score,
            pass_score,
            deadline,
        } => {
            let mut body = json!({
                "classId": class_id,
                "title": title,
                "type": kind,
                "fullScore": full_score,
                "passScore": pass_score,
            });
            if let Some(d) = deadline {
                body["deadline"] = json!(d);
            }
            let v = api.post("assessments", &body)?;
            Ok(format!(
                "已创建考核 {}（{}，{}分）",
                v["id"].as_str().unwrap_or(""),
                v["title"].as_str().unwrap_or(""),
                v["fullScore"].as_i64().unwrap_or(0)
            ))
        }
        AssessmentCmd::Submit {
            assessment_id,
            student_id,
        } => {
            let body = json!({
                "assessmentId": assessment_id,
                "studentId": student_id,
                "submittedAt": now_iso(),
            });
            let v = api.post("submissions", &body)?;
            Ok(format!(
                "已提交 {}（考核 {}）",
                v["id"].as_str().unwrap_or(""),
                v["assessmentId"].as_str().unwrap_or("")
            ))
        }
        AssessmentCmd::Grade {
            submission_id,
            score,
            comment,
        } => {
            let mut body = json!({ "score": score });
            if let Some(c) = comment {
                body["comment"] = json!(c);
            }
            let v = api.put(&format!("submissions/{submission_id}"), &body)?;
            Ok(format!(
                "已评分 {}：{} 分{}",
                v["id"].as_str().unwrap_or(""),
                v["score"].as_f64().unwrap_or(0.0),
                v["comment"]
                    .as_str()
                    .map(|c| format!("（{c}）"))
                    .unwrap_or_default()
            ))
        }
        AssessmentCmd::Stats { class_id } => {
            let assessments = api.get("assessments")?;
            let arr = assessments.as_array().ok_or("响应不是数组")?;
            let filtered: Vec<_> = arr
                .iter()
                .filter(|v| match &class_id {
                    Some(cid) => v["classId"].as_str() == Some(cid.as_str()),
                    None => true,
                })
                .collect();
            if filtered.is_empty() {
                return Ok("暂无考核".to_string());
            }
            let submissions = api.get("submissions")?;
            let subs = submissions.as_array().ok_or("响应不是数组")?;
            let mut out = String::from("考核\t类型\t满分\t提交数\t已评分\t平均分\n");
            for v in filtered {
                let aid = v["id"].as_str().unwrap_or("");
                let mine: Vec<_> = subs
                    .iter()
                    .filter(|s| s["assessmentId"].as_str() == Some(aid))
                    .collect();
                let graded: Vec<_> = mine.iter().filter(|s| s["score"].is_number()).collect();
                let avg = if graded.is_empty() {
                    0.0
                } else {
                    graded
                        .iter()
                        .map(|s| s["score"].as_f64().unwrap_or(0.0))
                        .sum::<f64>()
                        / graded.len() as f64
                };
                out.push_str(&format!(
                    "{}\t{}\t{}\t{}\t{}\t{:.1}\n",
                    v["title"].as_str().unwrap_or(""),
                    v["type"].as_str().unwrap_or("homework"),
                    v["fullScore"].as_i64().unwrap_or(0),
                    mine.len(),
                    graded.len(),
                    avg
                ));
            }
            Ok(out.trim_end().to_string())
        }
    }
}

fn now_iso() -> String {
    use std::time::{SystemTime, UNIX_EPOCH};
    let secs = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map(|d| d.as_secs())
        .unwrap_or(0);
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
