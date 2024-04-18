use std::thread;

use redis::Commands;
use tide::log;

pub mod pong;

#[derive(Clone)]
pub struct State {
    pub client: redis::Client,
}

impl State {
    pub fn new(addr: String) -> Self {
        let client = redis::Client::open(addr).expect("Failed to connect to Redis");
        return State { client };
    }
}

#[async_std::main]
async fn main() -> tide::Result<()> {
    femme::start();

    let client = redis::Client::open("redis://127.0.0.1:6379")?;

    thread::spawn(move || {
        let mut conn = client.get_connection().unwrap();
        let mut pubsub = conn.as_pubsub();

        pubsub.subscribe("com.github.hexennacht.ping-pong-pubsub.go.ping").unwrap();

        loop {
            let msg = pubsub.get_message().expect("kenapa error disini");
            let payload: String = msg.get_payload().expect("jangan sampai error");

            let mut request = serde_json::from_str::<pong::entity::PongRequest>(&payload).expect("Kenapa error");

            log::info!("Message Received: {} from: {}", payload, "com.github.hexennacht.ping-pong-pubsub.go.ping");

            if request.limit <= 0 {
                log::info!("Message limit reached, stopping the loop for message: {}", payload);
                continue;
            }

            request.limit -= 1;

            let response = serde_json::to_string(&request).expect("apalagi ini ngga boleh error");

            let mut new_conn = client.get_connection().expect("ini juga sama jangan sampai error");
            let publish: Result<(), redis::RedisError> = new_conn.publish("com.github.hexennacht.ping-pong-pubsub.rust.pong", response.clone());
            match publish {
                Ok(_) => {
                    log::info!("Message Published: {} to: {}", response, "com.github.hexennacht.ping-pong-pubsub.rust.pong");
                },
                Err(e) => {
                    log::error!("Error while publishing message: {}", e);
                }
                
            }
            continue;
        }
    });

    let mut app = tide::with_state(State::new("redis://127.0.0.1:6379".to_string()));
    
    app.with(tide::log::LogMiddleware::new());

    app.at("/pong").post(pong::handler::pong);

    app.listen("127.0.0.1:8000").await?;
    Ok(())
}

