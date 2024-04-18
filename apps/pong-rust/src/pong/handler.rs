use redis::Commands;
use tide::log;

use crate::State;

use super::entity;


pub async fn pong(mut req: tide::Request<State>) -> tide::Result<tide::Body> {
    let body: entity::PongRequest = req.body_json().await?;
    let state = req.state();

    println!("Received request: {:?}", body.clone());

    let mut conn = state.client.get_connection()?;

    let message = serde_json::to_string(&body)?;

    let published: Result<(), redis::RedisError> = conn.publish("com.github.hexennacht.ping-pong-pubsub.rust.pong", message.clone());
    match published {
        Ok(_) => {
            log::info!("Message Published: {} to: {}", message, "com.github.hexennacht.ping-pong-pubsub.rust.pong");
        },
        Err(e) => {
            log::error!("Error while publishing message: {}", e);
        }
    }

    Ok(tide::Body::from_json(&entity::PongResponse {
        message: "Success publishing message".to_string(),
    })?)
}