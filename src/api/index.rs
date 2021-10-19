//! # 导入队伍提交的压缩文件
//!
//! tar.gz 文件，最大 100 M
//!
//! api: /upload
//! - 方法： post
//! - 数据： tar.gz 文件
//! - 返回： [Results](../../models/struct.Results.html) [Score 二维数组](../../models/struct.Score.html)
//!

use std::{
    env::{current_dir, set_current_dir},
    fs::{create_dir_all, remove_file},
};

use rocket::{
    data::ToByteUnit,
    http::Status,
    post,
    response::status::Custom,
    serde::json::{serde_json::json, Value},
    Data,
};

use crate::models::{PasteId, Results, Score, TestMap};

#[post("/upload", data = "<data>")]
pub async fn upload(data: Data<'_>) -> Custom<Value> {
    let mut result = Results::<Vec<Vec<Score>>>::default();

    // 获取当前目录
    let current_dir = current_dir().expect("当前目录获取失败");

    // 保存文件到本地
    create_dir_all("upload").expect("创建 upload 目录失败");
    let paste_id = PasteId::new(16);
    let filename = format!("upload/{}.tar.gz", &paste_id.to_str());
    data.open(100.mebibytes())
        .into_file(&filename)
        .await
        .expect("打开文件失败");

    // 解压文件
    crate::unpacking(&filename, paste_id.to_str());

    // 比赛项目名
    let project = crate::get_project();
    // 获取测试方法表
    let api_test_map = TestMap::get(format!("go/{}/api_test.json", project));
    let chaincode_test_map = TestMap::get(format!("go/{}/chaincode_test.json", project));
    // 队伍路径
    let team = crate::get_dir(format!("./temp/{}", &paste_id.to_str()).as_str());
    // 学生项目路径
    let target = format!("{}/{}", team.to_str().unwrap_or("edu-mgmt"), &project);

    // 复制测试文件到项目
    let api_test = format!("{}/application/controller/api_test.go", &target);
    if !crate::copy_file(format!("go/{}/api_test.go", project).as_str(), &api_test) {
        result.message = Some("复制 api_test.go 失败！");
        result.code = Some(500);
        return Custom(Status::InternalServerError, json!(result));
    }
    let chaincode_test = format!("{}/chaincode/{}/chaincode_test.go", &target, &project);
    if !crate::copy_file(
        format!("go/{}/chaincode_test.go", project).as_str(),
        &chaincode_test,
    ) {
        result.message = Some("复制 chaincode_test.go 失败！");
        result.code = Some(500);
        return Custom(Status::InternalServerError, json!(result));
    }

    // 替换配置里的路径
    crate::replacen(&target, &current_dir);

    // 运行测试
    println!("开始运行 api_test 测试");
    let api_test = crate::run_test("./application/controller/", &api_test_map);
    println!("开始运行 chaincode_test 测试");
    let chaincode_test = crate::run_test(
        format!("./chaincode/{}/", &project).as_str(),
        &chaincode_test_map,
    );

    println!("评分完成");
    let scores = vec![api_test, chaincode_test];

    // 操作完删除文件
    set_current_dir(&current_dir).ok();
    remove_file(filename).expect("移除文件失败");

    result.data = Some(scores);
    Custom(Status::Ok, json!(result))
}
