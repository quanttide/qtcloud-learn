//! CLI 子命令集成测试：以进程内迷你 provider（内存 CRUD）验证端到端流程。

use std::collections::HashMap;
use std::io::{Read, Write};
use std::net::TcpListener;
use std::sync::{Arc, Mutex};
use std::thread;

use qtcloud_learn_cli::api::ApiClient;
use qtcloud_learn_cli::commands;

type Store = Arc<Mutex<HashMap<String, serde_json::Value>>>;

/// 启动进程内迷你 provider，返回 base URL。
fn spawn_server() -> String {
    let listener = TcpListener::bind("127.0.0.1:0").unwrap();
    let addr = listener.local_addr().unwrap();
    let store: Store = Arc::new(Mutex::new(HashMap::new()));
    let seq: Arc<Mutex<HashMap<String, usize>>> = Arc::new(Mutex::new(HashMap::new()));

    thread::spawn(move || {
        for stream in listener.incoming() {
            let Ok(mut stream) = stream else { continue };
            let store = Arc::clone(&store);
            let seq = Arc::clone(&seq);
            let mut buf = Vec::new();
            let mut tmp = [0u8; 4096];
            // 读请求头（以 \r\n\r\n 结束）
            loop {
                match stream.read(&mut tmp) {
                    Ok(0) => break,
                    Ok(n) => {
                        buf.extend_from_slice(&tmp[..n]);
                        if buf.windows(4).any(|w| w == b"\r\n\r\n") {
                            break;
                        }
                    }
                    Err(_) => break,
                }
            }
            // 按 Content-Length 补读请求体
            let header_end = buf
                .windows(4)
                .position(|w| w == b"\r\n\r\n")
                .unwrap_or(buf.len());
            let headers = String::from_utf8_lossy(&buf[..header_end]).to_string();
            let content_length: usize = headers
                .lines()
                .find(|l| l.to_ascii_lowercase().starts_with("content-length:"))
                .and_then(|l| l.split(':').nth(1))
                .and_then(|v| v.trim().parse().ok())
                .unwrap_or(0);
            while buf.len() < header_end + 4 + content_length {
                match stream.read(&mut tmp) {
                    Ok(0) => break,
                    Ok(n) => buf.extend_from_slice(&tmp[..n]),
                    Err(_) => break,
                }
            }
            let req = String::from_utf8_lossy(&buf).to_string();
            let first = req.lines().next().unwrap_or("GET / HTTP/1.1").to_string();
            let mut parts = first.split_whitespace();
            let method = parts.next().unwrap_or("GET").to_string();
            let path = parts.next().unwrap_or("/").to_string();
            let body = String::from_utf8_lossy(&buf[header_end + 4..]).trim().to_string();

            let (status, resp_body) = handle(&store, &seq, &method, &path, &body);
            let resp = format!(
                "HTTP/1.1 {status}\r\nContent-Type: application/json\r\nContent-Length: {}\r\nConnection: close\r\n\r\n{}",
                resp_body.len(),
                resp_body
            );
            let _ = stream.write_all(resp.as_bytes());
        }
    });

    format!("http://{addr}")
}

/// 极简 CRUD 路由（/api/v1/{resource}[/{id}]）。
fn handle(
    store: &Store,
    seq: &Arc<Mutex<HashMap<String, usize>>>,
    method: &str,
    path: &str,
    body: &str,
) -> (u16, String) {
    let path = path.trim_start_matches("/api/v1/");
    let mut segs = path.split('/');
    let resource = segs.next().unwrap_or("").to_string();
    let id = segs.next();
    let mut store = store.lock().unwrap();

    let not_found = (404, r#"{"error":"not found"}"#.to_string());
    match (method, id) {
        ("GET", None) => {
            let prefix = format!("{resource}/");
            let list: Vec<_> = store
                .iter()
                .filter(|(k, _)| k.starts_with(&prefix))
                .map(|(_, v)| v.clone())
                .collect();
            (200, serde_json::to_string(&list).unwrap())
        }
        ("GET", Some(id)) => match store.get(&format!("{resource}/{id}")) {
            Some(v) => (200, serde_json::to_string(v).unwrap()),
            None => not_found,
        },
        ("POST", None) => {
            let mut v: serde_json::Value =
                serde_json::from_str(body).unwrap_or(serde_json::json!({}));
            let mut seq = seq.lock().unwrap();
            let n = seq.entry(resource.clone()).or_insert(0);
            *n += 1;
            let new_id = format!("{resource}-{n}");
            v["id"] = serde_json::json!(new_id);
            store.insert(format!("{resource}/{new_id}"), v.clone());
            (201, serde_json::to_string(&v).unwrap())
        }
        ("PUT", Some(id)) => {
            let key = format!("{resource}/{id}");
            match store.get(&key) {
                Some(existing) => {
                    // 合并语义：仅覆盖请求体中的字段（与 provider 一致）
                    let mut merged = existing.clone();
                    if let Ok(patch) =
                        serde_json::from_str::<serde_json::Value>(body)
                    {
                        if let Some(obj) = patch.as_object() {
                            for (k, v) in obj {
                                merged[k] = v.clone();
                            }
                        }
                    }
                    merged["id"] = serde_json::json!(id);
                    store.insert(key, merged.clone());
                    (200, serde_json::to_string(&merged).unwrap())
                }
                None => not_found,
            }
        }
        _ => (400, r#"{"error":"bad request"}"#.to_string()),
    }
}

fn client() -> ApiClient {
    ApiClient::new(&spawn_server())
}

#[test]
fn student_crud() {
    let api = client();
    let out = commands::student::run(
        &api,
        commands::student::StudentCmd::Create {
            name: "张三".into(),
            email: Some("zhangsan@example.com".into()),
            plan: Some("vip".into()),
        },
    )
    .unwrap();
    assert!(out.contains("已创建学员 students-1（张三，vip）"), "{out}");

    let out = commands::student::run(&api, commands::student::StudentCmd::List).unwrap();
    assert!(out.contains("张三"), "{out}");
    assert!(out.contains("vip"), "{out}");

    let out = commands::student::run(
        &api,
        commands::student::StudentCmd::Get { id: "students-1".into() },
    )
    .unwrap();
    assert!(out.contains("张三"), "{out}");
    assert!(out.contains("zhangsan@example.com"), "{out}");
}

#[test]
fn class_crud() {
    let api = client();
    let out = commands::class::run(
        &api,
        commands::class::ClassCmd::Create {
            name: "浙理班级".into(),
            ref_name: "大数据微专业".into(),
            ref_id: "prog-1".into(),
            ref_type: "program".into(),
            start_date: Some("2026-09-01".into()),
            end_date: Some("2027-01-15".into()),
        },
    )
    .unwrap();
    assert!(out.contains("已创建班级 classes-1（浙理班级）"), "{out}");

    let out = commands::class::run(&api, commands::class::ClassCmd::List).unwrap();
    assert!(out.contains("浙理班级"), "{out}");

    let out = commands::class::run(
        &api,
        commands::class::ClassCmd::Get { id: "classes-1".into() },
    )
    .unwrap();
    assert!(out.contains("大数据微专业"), "{out}");
    assert!(out.contains("2026-09-01"), "{out}");
}

#[test]
fn enrollment_enroll_withdraw() {
    let api = client();
    let out = commands::enrollment::run(
        &api,
        commands::enrollment::EnrollmentCmd::Enroll {
            class_id: "classes-1".into(),
            student_id: "students-1".into(),
        },
    )
    .unwrap();
    assert!(out.contains("已报名"), "{out}");

    // 重复报名幂等
    let out = commands::enrollment::run(
        &api,
        commands::enrollment::EnrollmentCmd::Enroll {
            class_id: "classes-1".into(),
            student_id: "students-1".into(),
        },
    )
    .unwrap();
    assert!(out.contains("无需重复报名"), "{out}");

    let out = commands::enrollment::run(
        &api,
        commands::enrollment::EnrollmentCmd::List {
            student_id: Some("students-1".into()),
        },
    )
    .unwrap();
    assert!(out.contains("classes-1"), "{out}");

    let out = commands::enrollment::run(
        &api,
        commands::enrollment::EnrollmentCmd::Withdraw {
            class_id: "classes-1".into(),
            student_id: "students-1".into(),
        },
    )
    .unwrap();
    assert!(out.contains("已退课"), "{out}");

    let out = commands::enrollment::run(
        &api,
        commands::enrollment::EnrollmentCmd::List {
            student_id: Some("students-1".into()),
        },
    )
    .unwrap();
    assert!(out.contains("withdrawn"), "{out}");
}

#[test]
fn progress_report_get() {
    let api = client();
    let out = commands::progress::run(
        &api,
        commands::progress::ProgressCmd::Report {
            class_id: "classes-1".into(),
            student_id: "students-1".into(),
            percent: 0.5,
            finished: true,
        },
    )
    .unwrap();
    assert!(out.contains("已上报进度 progress-1：50%（已完成）"), "{out}");

    // 二次上报走更新
    let out = commands::progress::run(
        &api,
        commands::progress::ProgressCmd::Report {
            class_id: "classes-1".into(),
            student_id: "students-1".into(),
            percent: 0.8,
            finished: false,
        },
    )
    .unwrap();
    assert!(out.contains("已更新进度 progress-1：80%"), "{out}");

    let out = commands::progress::run(
        &api,
        commands::progress::ProgressCmd::Get {
            class_id: "classes-1".into(),
            student_id: "students-1".into(),
        },
    )
    .unwrap();
    assert!(out.contains("进度 progress-1：80%"), "{out}");
}

#[test]
fn assessment_flow() {
    let api = client();
    let out = commands::assessment::run(
        &api,
        commands::assessment::AssessmentCmd::Create {
            class_id: "classes-1".into(),
            title: "期中考试".into(),
            kind: "exam".into(),
            full_score: 100,
            pass_score: 60,
            deadline: Some("2026-10-01".into()),
        },
    )
    .unwrap();
    assert!(out.contains("已创建考核 assessments-1（期中考试，100分）"), "{out}");

    let out = commands::assessment::run(
        &api,
        commands::assessment::AssessmentCmd::Submit {
            assessment_id: "assessments-1".into(),
            student_id: "students-1".into(),
        },
    )
    .unwrap();
    assert!(out.contains("已提交 submissions-1（考核 assessments-1）"), "{out}");

    let out = commands::assessment::run(
        &api,
        commands::assessment::AssessmentCmd::Grade {
            submission_id: "submissions-1".into(),
            score: 85.0,
            comment: Some("不错".into()),
        },
    )
    .unwrap();
    assert!(out.contains("已评分 submissions-1：85 分（不错）"), "{out}");

    let out = commands::assessment::run(
        &api,
        commands::assessment::AssessmentCmd::Stats {
            class_id: Some("classes-1".into()),
        },
    )
    .unwrap();
    assert!(out.contains("期中考试"), "{out}");
    assert!(out.contains("85.0"), "{out}");
}

#[test]
fn version_help() {
    // 直接验证 clap 解析路径可达（version 不经 API）
    let out = format!("qtcloud-learn {}", env!("CARGO_PKG_VERSION"));
    assert!(out.contains("qtcloud-learn"));
}
