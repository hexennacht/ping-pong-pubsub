use super::entity;


pub async fn pong(mut req: tide::Request<()>) -> tide::Result<tide::Body> {
    let body: entity::PongRequest = req.body_json().await?;

    println!("Received request: {:?}", body);

    Ok(tide::Body::from_json(&entity::PongResponse {
        message: "pong".to_string(),
    })?)
}