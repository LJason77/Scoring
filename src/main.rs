#![deny(clippy::pedantic)]
#![allow(clippy::non_ascii_literal)]

#[rocket::main]
async fn main() {
    dotenv::dotenv().ok();

    let rocket = scoring::rocket();
    if let Err(err) = rocket.launch().await {
        println!("Rocket 启动错误: {}", err);
    }
}
