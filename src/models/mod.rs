use std::{fs::File, io::BufReader, path::Path};

use rand::{thread_rng, Rng};
use serde::{Deserialize, Serialize};

/// 对象结果
///
/// 返回给前端的对象结果
#[derive(Default, Serialize, Deserialize)]
pub struct Results<'a, T> {
    /// 状态码
    ///
    /// 返回正常时为空，只在错误时出现
    #[serde(skip_serializing_if = "Option::is_none")]
    pub code: Option<u16>,
    /// 返回的内容
    ///
    /// 具体内容视情况而定
    #[serde(skip_serializing_if = "Option::is_none")]
    pub data: Option<T>,
    /// 信息
    ///
    /// 一般情况下为空，错误时出现
    #[serde(skip_serializing_if = "Option::is_none")]
    pub message: Option<&'a str>,
}

/// base62 表
const BASE62: &[u8] = b"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";

/// 唯一 id
///
/// ## 示例
///
/// ```rust
/// let paste_id = PasteId::new(50);
/// println!("{}", paste_id.to_str());
/// ```
pub struct PasteId {
    #[doc(hidden)]
    id: String,
}

impl PasteId {
    /// 生成一个新的 `PasteId`
    ///
    /// ## 参数
    ///
    /// * `size`: 生成唯一 id 的长度
    #[must_use]
    pub fn new(size: usize) -> PasteId {
        let mut id = String::with_capacity(size);
        let mut rng = thread_rng();
        for _ in 0..size {
            id.push(BASE62[rng.gen::<usize>() % 62] as char);
        }
        PasteId { id }
    }

    /// 将 `PasteId` 转换为字符串
    #[must_use]
    pub fn to_str(&self) -> &str {
        &self.id
    }
}

/// 测试项
#[derive(Debug, Serialize, Deserialize)]
pub struct TestMap {
    /// 项目
    pub item: String,
    /// 描述
    pub info: String,
    /// 分数
    pub score: f32,
}

impl TestMap {
    pub fn get<P: AsRef<Path>>(path: P) -> Vec<TestMap> {
        let file = File::open(path).expect("打开 TestMap 文件失败");
        let reader = BufReader::new(file);

        serde_json::from_reader(reader).expect("TestMap 文件转换失败")
    }
}

/// 测试的结果
#[derive(Debug, Serialize, Deserialize)]
pub struct Test {
    /// 时间
    #[serde(rename = "Time")]
    pub time: String,
    #[serde(rename = "Action")]
    pub action: String,
    #[serde(rename = "Package")]
    pub package: String,
    #[serde(rename = "Output", skip_serializing_if = "Option::is_none")]
    pub output: Option<String>,
    #[serde(rename = "Elapsed", skip_serializing_if = "Option::is_none")]
    pub elapsed: Option<f32>,
    #[serde(rename = "Test", skip_serializing_if = "Option::is_none")]
    pub test: Option<String>,
}

/// 分数
///
/// 返回给前端的对象结果
#[derive(Debug, Serialize, Deserialize)]
pub struct Score {
    /// 题目分数
    pub totals: f32,
    /// 得分
    pub score: f32,
    /// 描述
    pub info: String,
}
