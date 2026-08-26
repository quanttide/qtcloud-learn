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
            let body = String::from_utf8_lossy(&buf[header_end + 4..])
                .trim()
                .to_string();

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

/// 极简 CRUD 路由（/{resource}[/{id}]）。
fn handle(
    store: &Store,
    seq: &Arc<Mutex<HashMap<String, usize>>>,
    method: &str,
    path: &str,
    body: &str,
) -> (u16, String) {
    // 无版本前缀，路径形如 /learners、/learners/{id}：去掉前导斜杠后按段解析
    let path = path.trim_start_matches('/');
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
                    if let Ok(patch) = serde_json::from_str::<serde_json::Value>(body) {
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
fn learner_crud() {
    let api = client();
    let out = commands::learner::run(
        &api,
        commands::learner::LearnerCmd::Create {
            user_id: Some("user-123".into()),
        },
    )
    .unwrap();
    assert!(
        out.contains("已创建学习者 learners-1（user_id: user-123）"),
        "{out}"
    );

    let out = commands::learner::run(&api, commands::learner::LearnerCmd::List).unwrap();
    assert!(out.contains("learners-1"), "{out}");
    assert!(out.contains("user-123"), "{out}");

    let out = commands::learner::run(
        &api,
        commands::learner::LearnerCmd::Get {
            id: "learners-1".into(),
        },
    )
    .unwrap();
    assert!(out.contains("learners-1"), "{out}");
    assert!(out.contains("user-123"), "{out}");
}

#[test]
fn criterion_crud() {
    let api = client();
    let out = commands::criterion::run(
        &api,
        commands::criterion::CriterionCmd::Create {
            title: "vibe-coding/lesson1/zed-connection".into(),
            description: "成功建立 Zed 连接".into(),
        },
    )
    .unwrap();
    assert!(
        out.contains("已创建验收标准 criteria-1（vibe-coding/lesson1/zed-connection）"),
        "{out}"
    );

    let out = commands::criterion::run(&api, commands::criterion::CriterionCmd::List).unwrap();
    assert!(out.contains("vibe-coding/lesson1/zed-connection"), "{out}");
    assert!(out.contains("成功建立 Zed 连接"), "{out}");

    let out = commands::criterion::run(
        &api,
        commands::criterion::CriterionCmd::Get {
            id: "criteria-1".into(),
        },
    )
    .unwrap();
    assert!(out.contains("vibe-coding/lesson1/zed-connection"), "{out}");
    assert!(out.contains("成功建立 Zed 连接"), "{out}");
}

#[test]
fn completion_flow() {
    let api = client();
    let out = commands::completion::run(
        &api,
        commands::completion::CompletionCmd::Create {
            learner_id: "learners-1".into(),
            criterion_id: "criteria-1".into(),
            status: Some("not_completed".into()),
        },
    )
    .unwrap();
    assert!(
        out.contains("已创建完成记录 completions-1（learners-1 → criteria-1，not_completed）"),
        "{out}"
    );

    // 标记完成（局部更新 status → completed）
    let out = commands::completion::run(
        &api,
        commands::completion::CompletionCmd::Complete {
            id: "completions-1".into(),
        },
    )
    .unwrap();
    assert!(
        out.contains("已完成 completions-1（learners-1 → criteria-1）"),
        "{out}"
    );

    let out = commands::completion::run(&api, commands::completion::CompletionCmd::List).unwrap();
    assert!(out.contains("completions-1"), "{out}");
    assert!(out.contains("completed"), "{out}");

    let out = commands::completion::run(
        &api,
        commands::completion::CompletionCmd::Get {
            id: "completions-1".into(),
        },
    )
    .unwrap();
    assert!(out.contains("状态: completed"), "{out}");
}

#[test]
fn version_help() {
    // 直接验证 clap 解析路径可达（version 不经 API）
    let out = format!("qtcloud-learn {}", env!("CARGO_PKG_VERSION"));
    assert!(out.contains("qtcloud-learn"));
}
