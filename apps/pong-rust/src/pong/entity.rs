use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct PongRequest {
    pub message: String,
    pub limit: i32,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PongResponse {
    pub message: String,
}