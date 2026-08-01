//! Provider API 客户端（`/api/v1` 前缀）。

use serde_json::Value;

/// 指向 `qtcloud-learn/provider` 的 HTTP 客户端。
pub struct ApiClient {
    base_url: String,
}

impl ApiClient {
    pub fn new(base_url: &str) -> Self {
        ApiClient {
            base_url: base_url.trim_end_matches('/').to_string(),
        }
    }

    pub fn get(&self, resource: &str) -> Result<Value, String> {
        self.request("GET", resource, None)
    }

    pub fn post(&self, resource: &str, body: &Value) -> Result<Value, String> {
        self.request("POST", resource, Some(body))
    }

    pub fn put(&self, resource: &str, body: &Value) -> Result<Value, String> {
        self.request("PUT", resource, Some(body))
    }

    fn request(
        &self,
        method: &str,
        resource: &str,
        body: Option<&Value>,
    ) -> Result<Value, String> {
        let url = format!("{}/api/v1/{}", self.base_url, resource);
        let resp = match method {
            "GET" => ureq::get(&url).call(),
            "POST" => ureq::post(&url).send_json(body.unwrap_or(&Value::Null)),
            "PUT" => ureq::put(&url).send_json(body.unwrap_or(&Value::Null)),
            _ => return Err(format!("不支持的请求方法: {method}")),
        }
        .map_err(|e| format!("请求失败: {e}"))?;

        let status = resp.status();
        let text = resp
            .into_string()
            .map_err(|e| format!("读取响应失败: {e}"))?;

        if !(200..300).contains(&status) {
            return Err(format!("HTTP {status}: {text}"));
        }
        if text.trim().is_empty() {
            return Ok(Value::Null);
        }
        serde_json::from_str(&text).map_err(|e| format!("无效 JSON: {e}"))
    }
}
