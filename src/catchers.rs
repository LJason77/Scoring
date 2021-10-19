use rocket::{
    catch,
    http::Status,
    response::status::Custom,
    serde::json::{json, Value},
    Request,
};

use crate::models::Results;

#[catch(404)]
pub fn not_found(status: Status, req: &Request) -> Custom<Value> {
    let mut result = Results::<Value>::default();
    println!(
        "{}：\nIP：{:?}\nurl：{:?}\n",
        status,
        &req.client_ip(),
        &req.uri().path().as_str()
    );
    result.message = Some("路由未找到");
    Custom(status, json!(result))
}
