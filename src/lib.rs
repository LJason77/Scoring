#![deny(clippy::pedantic)]
#![allow(clippy::non_ascii_literal)]

use std::{
    env::set_current_dir,
    fs::{copy, create_dir_all, read_dir, File},
    io::{BufRead, BufReader, Read, Write},
    path::{Path, PathBuf},
    process::{exit, Command},
};

use rocket::{catchers, fs::FileServer, routes, Build, Rocket};

use api::index;

use crate::models::{Score, Test, TestMap};

pub mod api;
pub mod catchers;
pub mod models;

#[must_use]
pub fn rocket() -> Rocket<Build> {
    rocket::build()
        // 静态文件
        .mount("/", FileServer::from("static/"))
        // 评分
        .mount("/", routes![index::upload])
        // 错误处理
        .register("/", catchers![catchers::not_found])
}

/// 解压文件
pub fn unpacking(file: &str, name: &str) {
    let dir = format!("temp/{}/", name);
    create_dir_all(&dir).expect("创建临时目录失败");
    Command::new("tar")
        .arg("zxf")
        .arg(file)
        .args(["-C", &dir])
        .output()
        .expect("解压失败");
    println!("解压完成");
}

/// 复制文件
#[must_use]
pub fn copy_file(source: &str, target: &str) -> bool {
    if !Path::new(source).exists() {
        return false;
    }
    copy(source, target).is_ok()
}

/// 获得队伍路径
#[must_use]
pub fn get_dir(dir: &str) -> PathBuf {
    let read_dir = read_dir(PathBuf::from(dir)).expect("获得队伍路径失败");
    for dir in read_dir {
        let file = dir.expect("获取子路径失败").path();
        if file.is_dir() {
            return file;
        }
    }
    PathBuf::new()
}

/// 运行测试
#[must_use]
pub fn run_test(test: &str, test_map: &[TestMap]) -> Vec<Score> {
    create_dir_all("test").expect("创建 test 目录失败");

    let mut scores = Vec::new();
    for test_item in test_map {
        let api_test = Command::new("/usr/go/bin/go")
            .arg("test")
            .arg("-count=1")
            .arg("-json")
            .arg("--run")
            .arg(&test_item.item)
            .arg(test)
            .output()
            .expect("运行失败")
            .stdout;
        // 写入文件
        let mut api_test_txt =
            File::create(format!("test/{}.json", &test_item.item)).expect("创建文件失败");
        api_test_txt.write(&api_test).ok();
        // 读取文件
        let file =
            File::open(format!("test/{}.json", &test_item.item)).expect("文件打开失败");

        let mut vec = Vec::new();
        for line in BufReader::new(file).lines().flatten() {
            vec.push(serde_json::from_str::<Test>(&line).expect("解析失败"));
        }

        // 生成结果
        for test in vec {
            // 跳过 output、run
            if test.action == "output" || test.action == "run" || test.test.is_none() {
                continue;
            }
            let mut score = Score {
                totals: test_item.score,
                score: 0.0,
                info: test_item.info.clone(),
            };
            if test.test.as_ref().expect("") == &test_item.item {
                score.score = test_item.score;
            }
            scores.push(score);
        }
    }
    scores
}

/// 替换文本
pub fn replacen(path: &str, current_dir: &Path) {
    // 切换工作路径
    set_current_dir(Path::new(path)).ok();

    let mut config_e2e =
        File::open("application/config_e2e.yaml").expect("打开配置文件失败");
    let mut data = String::new();
    config_e2e.read_to_string(&mut data).ok();
    // 提前关闭文件
    drop(config_e2e);

    let new_data = data.replace(
        "/home/zcedu/go/src/gdzce.cn/edu-mgmt",
        current_dir.to_str().unwrap_or(""),
    );

    // 写入文件
    let mut dst = File::create("application/config_e2e.yaml").expect("创建文件失败");
    dst.write(new_data.as_bytes()).ok();
    println!("替换路径完成");
}

/// 返回项目名
#[must_use]
pub fn get_project() -> String {
    dotenv::var("project").unwrap_or_else(|_| {
        println!("环境变量：project 未找到");
        exit(0)
    })
}
