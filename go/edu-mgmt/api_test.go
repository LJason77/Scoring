package controller

import (
	"bufio"
	"bytes"
	"edu-mgmt/application/blockchain"
	"edu-mgmt/application/lib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// 根据特定请求uri，发起get请求返回响应
func get(uri string, router *gin.Engine) ([]byte, int) {
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	// 提取响应
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(result.Body)
	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body, result.StatusCode
}

// 根据特定请求uri和参数param，以表单形式传递参数，发起post请求返回响应
func postForm(uri string, param []byte, router *gin.Engine) ([]byte, int) {
	// 构造post请求
	req := httptest.NewRequest("POST", uri, strings.NewReader(bytes.NewBuffer(param).String()))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应handler接口
	router.ServeHTTP(w, req)
	// 提取响应
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(result.Body)
	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body, result.StatusCode
}

func expectApi(res int, status int, testFunc string) {
	str := "******"
	file := "./api_test_result.txt"
	switch status {
	case 1:
		if res == 200 {
			str += "成功响应测试：pass"
			writeFile(file, testFunc+":PASS\n")
		} else {
			str += "成功响应测试：fail"
			writeFile(file, testFunc+":FAIL\n")
		}
	case 2:
		if res >= 400 && res <= 403 {
			str += "参数测试：pass"
			writeFile(file, testFunc+":PASS\n")
		} else {
			str += "参数测试：fail"
			writeFile(file, testFunc+":FAIL\n")
		}
	case 3:
		if res == 500 {
			str += "错误响应测试：pass"
			writeFile(file, testFunc+":PASS\n")
		} else {
			str += "错误响应测试：fail"
			writeFile(file, testFunc+":FAIL\n")
		}
	case 4:
		str += "测试通过：pass"
		writeFile(file, testFunc+":PASS\n")
	case 5:
		str += "测试失败：fail"
		writeFile(file, testFunc+":FAIL\n")
	}
	fmt.Println(str)
}

func writeFile(filePath string, data string) {
	//os.Remove(filePath)
	//os.Create(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("文件打开失败", err)
		file, _ = os.Create(filePath)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(data)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

var routers *gin.Engine

// testing 初始化内容
func init() {
	routers = SetupRouter()
	blockchain.Init("../config_e2e.yaml")
}

// Test_SDK SDK能否访问区块链网络
func Test_SDK(t *testing.T) {
	// ************** 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("test"),
		[]byte("1010"),
	})
	// **************
	resp, errString := blockchain.ChannelExecute("set", [][]byte{
		[]byte("test"),
		[]byte("1010"),
		[]byte("测试"),
	})
	t.Logf("Test_SDK: %+v\n", resp.ChaincodeStatus)
	t.Logf("Test_SDK: %+v\n", errString)
	if resp.ChaincodeStatus == 200 && errString == "" {
		time.Sleep(2 * time.Second)
		resp, _ := blockchain.ChannelQuery("get", [][]byte{
			[]byte("test"),
			[]byte("1010"),
		})
		t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
		t.Logf("get_data: %+v\n", resp.ChaincodeStatus)
		if len(resp.Payload) != 0 {
			expectApi(int(resp.ChaincodeStatus), 4, "Test_SDK")
		} else {
			expectApi(int(resp.ChaincodeStatus), 5, "Test_SDK")
			t.Errorf("err: %+v\n", len(resp.Payload))
			t.FailNow()
		}
	} else {
		expectApi(int(resp.ChaincodeStatus), 5, "Test_SDK")
		t.Errorf("err: %+v\n", len(resp.Payload))
		t.FailNow()
	}
}

// Test_Register 注册-测试创建时间是否添加
func Test_Register(t *testing.T) {
	data := lib.Person{
		Id:       "202105062136",
		Account:  "liangsheng",
		Password: "123456",
		Name:     "testStudent",
		Type:     "student",
	}
	// ************** 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
	})
	// **************
	dataByte, _ := json.Marshal(&data)
	resp, statusCode := postForm("/register", dataByte, routers)
	t.Logf("register: %+v\n", bytes.NewBuffer(resp).String())
	t.Logf("register: %+v\n", statusCode)
	if statusCode == 200 && len(resp) != 0 {
		time.Sleep(2 * time.Second)
		resp, _ := blockchain.ChannelQuery("get", [][]byte{
			[]byte(data.Type),
			[]byte(data.Account),
		})
		t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
		t.Logf("get_data: %+v\n", resp.ChaincodeStatus)
		var result lib.Person
		_ = json.Unmarshal(resp.Payload, &result)
		if !result.CreatedAt.IsZero() {
			expectApi(int(resp.ChaincodeStatus), 4, "Test_Register")
		} else {
			expectApi(int(resp.ChaincodeStatus), 5, "Test_Register")
			t.Error("err")
			t.FailNow()
		}
	} else {
		expectApi(statusCode, 1, "Test_Register")
		t.Error("err")
		t.FailNow()
	}
}

// Test_Login1 登录-验证密码
func Test_Login1(t *testing.T) {
	data := lib.Person{
		Account:  "school01",
		Password: "1234567",
		Type:     "school",
	}
	dataByte, _ := json.Marshal(&data)
	resp, statusCode := postForm("/login", dataByte, routers)
	t.Logf("login: %+v\n", bytes.NewBuffer(resp).String())
	t.Logf("login: %+v\n", statusCode)
	if statusCode == 400 {
		expectApi(statusCode, 2, "Test_Login1")
	} else {
		expectApi(statusCode, 2, "Test_Login1")
		t.Error("err")
		t.FailNow()
	}
}

// Test_Login2 登录-新建帐号能否登录
func Test_Login2(t *testing.T) {
	data := lib.Person{
		Id:       "202105062136",
		Account:  "liangsheng",
		Password: "123456",
		Name:     "testStudent",
		Type:     "student",
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	result, statusCode := postForm("/login", dataByte, routers)
	t.Logf("login: %+v\n", bytes.NewBuffer(result).String())
	t.Logf("login: %+v\n", statusCode)
	if statusCode == 400 {
		expectApi(statusCode, 2, "Test_Login2")
	} else {
		expectApi(statusCode, 2, "Test_Login2")
		t.Error("err")
		t.FailNow()
	}

}

// Test_ConfirmUser1 认证用户-是否生成更新时间
func Test_ConfirmUser1(t *testing.T) {
	data := lib.Person{
		Id:       "202105062136",
		Account:  "liangsheng",
		Password: "123456",
		Name:     "testStudent",
		Type:     "student",
		Status:   "new",
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	tempData := struct {
		Type          string    `json:"type"`
		Account       string    `json:"account"`
		TargetType    string    `json:"targetType"`
		TargetAccount string    `json:"targetAccount"`
		UpdatedAt     time.Time `json:"updated_at"` // 更新时间
	}{
		Type:          "school",
		Account:       "school01",
		TargetAccount: "liangsheng",
		TargetType:    "student",
	}
	tempByte, _ := json.Marshal(&tempData)
	result, statusCode := postForm("/confirm_user", tempByte, routers)
	t.Logf("login: %+v\n", bytes.NewBuffer(result).String())
	t.Logf("login: %+v\n", statusCode)
	if statusCode == 200 {
		time.Sleep(2 * time.Second)
		resp, _ := blockchain.ChannelQuery("get", [][]byte{
			[]byte(data.Type),
			[]byte(data.Account),
		})
		t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
		t.Logf("get_data: %+v\n", resp.ChaincodeStatus)
		var result lib.Person
		_ = json.Unmarshal(resp.Payload, &result)
		if !result.UpdatedAt.IsZero() {
			expectApi(int(resp.ChaincodeStatus), 4, "Test_ConfirmUser1")
		} else {
			expectApi(int(resp.ChaincodeStatus), 5, "Test_ConfirmUser1")
			t.Error("err")
			t.FailNow()
		}
	} else {
		expectApi(statusCode, 1, "Test_ConfirmUser1")
		t.Error("err")
		t.FailNow()
	}
}

// Test_ConfirmUser2 认证用户-是否修改状态成功
func Test_ConfirmUser2(t *testing.T) {
	data := lib.Person{
		Id:       "202105062136",
		Account:  "liangsheng",
		Password: "123456",
		Name:     "testStudent",
		Type:     "student",
		Status:   "new",
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)

	tempData := struct {
		Type          string    `json:"type"`
		Account       string    `json:"account"`
		TargetType    string    `json:"targetType"`
		TargetAccount string    `json:"targetAccount"`
		UpdatedAt     time.Time `json:"updated_at"` // 更新时间
	}{
		Type:          "school",
		Account:       "school01",
		TargetAccount: "liangsheng",
		TargetType:    "student",
	}
	tempByte, _ := json.Marshal(&tempData)
	result, statusCode := postForm("/confirm_user", tempByte, routers)
	t.Logf("login: %+v\n", bytes.NewBuffer(result).String())
	t.Logf("login: %+v\n", statusCode)
	if statusCode == 200 {
		time.Sleep(2 * time.Second)
		resp, _ := blockchain.ChannelQuery("get", [][]byte{
			[]byte(data.Type),
			[]byte(data.Account),
		})
		t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
		t.Logf("get_data: %+v\n", resp.ChaincodeStatus)
		var result lib.Person
		_ = json.Unmarshal(resp.Payload, &result)
		if result.Status == "confirm" {
			expectApi(int(resp.ChaincodeStatus), 4, "Test_ConfirmUser2")
		} else {
			expectApi(int(resp.ChaincodeStatus), 5, "Test_ConfirmUser2")
			t.Error("err")
			t.FailNow()
		}
	} else {
		expectApi(statusCode, 1, "Test_ConfirmUser2")
		t.Error("err")
		t.FailNow()
	}
}

// Test_GetPersons1 状态为空，是否返回所有数据
func Test_GetPersons1(t *testing.T) {
	tempByte, statusCode := get("/getPersons?type=teacher&status=", routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_Login1")
		t.Error("返回值不对")
		t.Error(bytes.NewBuffer(tempByte).String())
		t.FailNow()
	}
	var data []lib.Person
	err := json.Unmarshal(tempByte, &data)
	if err != nil {
		t.Errorf("err: %+v\n", err)
		t.FailNow()
	}
	t.Logf("get_data: %+v\n", data)
	if len(data) == 5 {
		expectApi(statusCode, 4, "Test_GetPersons1")
	} else {
		expectApi(statusCode, 5, "Test_GetPersons1")
		t.Error("err")
		t.FailNow()
	}
}

// Test_GetPersons2 状态为新建，是否返回新建状态的人员列表
func Test_GetPersons2(t *testing.T) {
	data := lib.Person{
		Id:       "202105062136",
		Account:  "liangsheng",
		Password: "123456",
		Name:     "testStudent",
		Type:     "teacher",
		Status:   "new",
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	tempByte, statusCode := get("/getPersons?type=teacher&status=new", routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_GetPersons2")
		t.Error("返回值不对")
		t.FailNow()
	}
	var temp []lib.Person
	err := json.Unmarshal(tempByte, &temp)
	if err != nil {
		t.Errorf("err: %+v\n", err)
		t.FailNow()
	}
	t.Logf("get_data: %+v\n", temp)
	t.Logf("get_data: %+v\n", len(temp))
	if len(temp) == 0 {
		expectApi(statusCode, 5, "Test_GetPersons2")
		t.Error("返回数据数量不对")
		t.FailNow()
	}
	for _, datum := range temp {
		if datum.Status != "new" {
			expectApi(statusCode, 5, "Test_GetPersons2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(statusCode, 4, "Test_GetPersons2")
}

// Test_GetPersons3 状态为确认，是否返回确认状态的人员列表
func Test_GetPersons3(t *testing.T) {
	data := lib.Person{
		Id:       "202105062136",
		Account:  "liangsheng",
		Password: "123456",
		Name:     "testStudent",
		Type:     "teacher",
		Status:   "new",
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte(data.Type),
		[]byte(data.Account),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	tempByte, statusCode := get("/getPersons?type=teacher&status=confirm", routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_GetPersons3")
		t.Error("返回值不对")
		t.FailNow()
	}
	var temp []lib.Person
	err := json.Unmarshal(tempByte, &temp)
	if err != nil {
		t.Errorf("err: %+v\n", err)
		t.FailNow()
	}
	t.Logf("get_data: %+v\n", temp)
	t.Logf("get_data: %+v\n", len(temp))
	if len(temp) == 0 {
		expectApi(statusCode, 5, "Test_GetPersons3")
		t.Error("返回数据数量不对")
		t.FailNow()
	}
	for _, datum := range temp {
		if datum.Status != "confirm" {
			expectApi(statusCode, 5, "Test_GetPersons3")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(statusCode, 4, "Test_GetPersons3")
}

// Test_GetPersons4 状态为其他，是否返回状态不存在报错
func Test_GetPersons4(t *testing.T) {
	_, statusCode := get("/getPersons?type=teacher&status=test", routers)
	t.Logf("get_data: %+v\n", statusCode)
	if statusCode == 400 {
		expectApi(statusCode, 2, "Test_GetPersons4")
	} else {
		expectApi(statusCode, 2, "Test_GetPersons4")
		t.Error("err")
		t.FailNow()
	}
}

// Test_NewPaper1 是否分配论文ID
func Test_NewPaper1(t *testing.T) {
	data := lib.Paper{
		Title:    "Test_NewPaper1",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
	}
	dataByte, _ := json.Marshal(&data)
	result, statusCode := postForm("/newPaper", dataByte, routers)
	if statusCode != 200 {
		t.Logf("NewPaper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	time.Sleep(2 * time.Second)
	resp, _ := blockchain.ChannelQuery("getPapers", [][]byte{
		[]byte("paper"),
	})
	var papers []lib.Paper
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &papers)
	for _, value := range papers {
		if value.Title == "Test_NewPaper1" {
			paper = value
		}
	}

	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(paper.Id),
	})
	// ******* end *******

	if paper.Id == "" {
		expectApi(statusCode, 5, "Test_NewPaper1")
		t.Error("没有分配论文ID")
		t.FailNow()
	}
	t.Logf("get_data: %+v\n", paper)
	expectApi(statusCode, 4, "Test_NewPaper1")
}

// Test_NewPaper2 是否生成创建时间
func Test_NewPaper2(t *testing.T) {
	data := lib.Paper{
		Title:    "Test_NewPaper2",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
	}
	dataByte, _ := json.Marshal(&data)
	result, statusCode := postForm("/newPaper", dataByte, routers)
	if statusCode != 200 {
		t.Errorf("NewPaper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	time.Sleep(2 * time.Second)
	resp, _ := blockchain.ChannelQuery("getPapers", [][]byte{
		[]byte("paper"),
	})
	var papers []lib.Paper
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &papers)
	for _, value := range papers {
		if value.Title == "Test_NewPaper2" {
			paper = value
		}
	}

	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(paper.Id),
	})
	// ******* end *******

	if paper.CreatedAt.IsZero() {
		expectApi(statusCode, 5, "Test_NewPaper2")
		t.Error("没有生成创建时间")
		t.FailNow()
	}
	t.Logf("get_data: %+v\n", paper)
	expectApi(statusCode, 4, "Test_NewPaper2")
}

// Test_SubmitPaper 是否生成提交时间
func Test_SubmitPaper(t *testing.T) {
	data := lib.Paper{
		Id:       "202105080001",
		Title:    "Test_SubmitPaper",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
		Status:   lib.PaperStatusNew,
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	result, statusCode := postForm("/submitPaper", dataByte, routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_SubmitPaper")
		t.Errorf("submitPaper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.FailNow()
	}
	time.Sleep(2 * time.Second)
	resp, _ = blockchain.ChannelQuery("get", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &paper)
	t.Logf("get_data: %+v\n", paper)
	if paper.SubmitTime.IsZero() {
		expectApi(statusCode, 5, "Test_SubmitPaper")
		t.Error("没有生成更新时间")
		t.FailNow()
	}
	expectApi(statusCode, 4, "Test_SubmitPaper")
}

// Test_DeletePaper 是否删除论文成功
func Test_DeletePaper(t *testing.T) {
	data := lib.Paper{
		Id:       "202105080002",
		Title:    "Test_DeletePaper",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
		Status:   lib.PaperStatusNew,
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	result, statusCode := postForm("/deletePaper", dataByte, routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_DeletePaper")
		t.Errorf("deletePaper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.FailNow()
	}
	time.Sleep(2 * time.Second)
	resp, errString := blockchain.ChannelQuery("get", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	t.Logf("get_data: %+v\n", errString)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &paper)
	if len(resp.Payload) == 0 || paper.Status == lib.PaperStatusCancel {
		expectApi(statusCode, 4, "Test_DeletePaper")
	} else {
		expectApi(statusCode, 5, "Test_DeletePaper")
		t.Error("删除论文失败")
		t.FailNow()
	}
}

// Test_RejectPaper 是否生成退回时间
func Test_RejectPaper(t *testing.T) {
	data := lib.Paper{
		Id:         "202105080003",
		Title:      "Test_RejectPaper",
		Abstract:   "testAbstract",
		Student:    "liangsheng",
		Teacher:    "teacher01",
		RefuseInfo: "test",
		Status:     lib.PaperStatusSubmit,
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	result, statusCode := postForm("/rejectPaper", dataByte, routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_RejectPaper")
		t.Errorf("rejectPaper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.FailNow()
	}
	time.Sleep(2 * time.Second)
	resp, errString := blockchain.ChannelQuery("get", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	t.Logf("get_data: %+v\n", errString)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &paper)
	if !paper.RejectTime.IsZero() {
		expectApi(statusCode, 4, "Test_RejectPaper")
	} else {
		expectApi(statusCode, 5, "Test_RejectPaper")
		t.Error("生成退回时间失败")
		t.FailNow()
	}
}

// Test_ApprovePaper 是否生成通过时间
func Test_ApprovePaper(t *testing.T) {
	data := lib.Paper{
		Id:       "202105080004",
		Title:    "Test_ApprovePaper",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
		Status:   lib.PaperStatusSubmit,
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	result, statusCode := postForm("/approvePaper", dataByte, routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_ApprovePaper")
		t.Errorf("approvePaper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.FailNow()
	}
	time.Sleep(2 * time.Second)
	resp, errString := blockchain.ChannelQuery("get", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	t.Logf("get_data: %+v\n", errString)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &paper)
	if !paper.ApproveDate.IsZero() {
		expectApi(statusCode, 4, "Test_ApprovePaper")
	} else {
		expectApi(statusCode, 5, "Test_ApprovePaper")
		t.Error("生成通过时间失败")
		t.FailNow()
	}
}

// Test_MarkPaper 评分是否成功
func Test_MarkPaper(t *testing.T) {
	data := lib.Paper{
		Id:       "202105080005",
		Title:    "Test_MarkPaper",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
		Score:    90,
		Status:   lib.PaperStatusOralDefense,
		OralDefense: map[string]int64{
			"teacher01": -1,
			"teacher02": -1,
			"teacher03": -1,
			"teacher04": -1,
			"teacher05": -1,
		},
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	result, statusCode := postForm("/markPaper", dataByte, routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_MarkPaper")
		t.Errorf("markPaper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.FailNow()
	}
	time.Sleep(2 * time.Second)
	resp, errString := blockchain.ChannelQuery("get", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	t.Logf("get_data: %+v\n", errString)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &paper)
	if value, ok := paper.OralDefense["teacher01"]; ok && value != -1 {
		expectApi(statusCode, 4, "Test_MarkPaper")
	} else {
		expectApi(statusCode, 5, "Test_MarkPaper")
		t.Error("生成通过时间失败")
		t.FailNow()
	}
}

// Test_MarkPaper 是否生成更新时间
func Test_UpdatePaper(t *testing.T) {
	data := lib.Paper{
		Id:       "202105080006",
		Title:    "Test_UpdatePaper",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
		Status:   lib.PaperStatusNew,
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	data.Abstract = "testUpdate"
	dataByte2, _ := json.Marshal(&data)
	result, statusCode := postForm("/updatePaper", dataByte2, routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_UpdatePaper")
		t.Errorf("updatePaper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.FailNow()
	}
	time.Sleep(2 * time.Second)
	resp, errString := blockchain.ChannelQuery("get", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	t.Logf("get_data: %+v\n", errString)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &paper)
	if !paper.UpdatedAt.IsZero() {
		expectApi(statusCode, 4, "Test_UpdatePaper")
	} else {
		expectApi(statusCode, 5, "Test_UpdatePaper")
		t.Error("生成更新时间失败")
		t.FailNow()
	}
}

// Test_GetPaper0 是否获取成功
func Test_GetPaper0(t *testing.T) {
	data := lib.Paper{
		Id:       "202105080007",
		Title:    "Test_GetPaper0",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
		Status:   lib.PaperStatusNew,
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	result, statusCode := get("/paper?id=202105080007&student=liangsheng", routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_GetPaper0")
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.FailNow()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var paper lib.Paper
	_ = json.Unmarshal(result, &paper)
	if paper.Id == "202105080007" {
		expectApi(statusCode, 4, "Test_GetPaper0")
	} else {
		expectApi(statusCode, 5, "Test_GetPaper0")
		t.Error("获取失败")
		t.FailNow()
	}
}

// Test_Schedule 安排答辩是否成功
func Test_Schedule(t *testing.T) {
	data := lib.Paper{
		Id:       "202105080008",
		Title:    "Test_Schedule",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
		Status:   lib.PaperStatusApprove,
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	data.OralDefense = map[string]int64{
		"teacher01": -1,
		"teacher02": -1,
		"teacher03": -1,
		"teacher04": -1,
		"teacher05": -1,
	}
	dataByte2, _ := json.Marshal(&data)
	result, statusCode := postForm("/oralDefense", dataByte2, routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_Schedule")
		t.Errorf("oralDefense 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.FailNow()
	}
	time.Sleep(2 * time.Second)
	resp, errString := blockchain.ChannelQuery("get", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	t.Logf("get_data: %+v\n", errString)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &paper)
	if len(paper.OralDefense) == 5 {
		expectApi(statusCode, 4, "Test_Schedule")
	} else {
		expectApi(statusCode, 5, "Test_Schedule")
		t.Error("安排答辩失败")
		t.FailNow()
	}
}

// Test_ArchivePaper 论文归档是否成功
func Test_ArchivePaper(t *testing.T) {
	data := lib.Paper{
		Id:       "202105080009",
		Title:    "Test_ArchivePaper",
		Abstract: "testAbstract",
		Student:  "liangsheng",
		Teacher:  "teacher01",
		Status:   lib.PaperStatusPassed,
	}
	// ******* start ******* 延后删除测试数据
	defer blockchain.ChannelExecute("delete", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	// ******* end *******
	dataByte, _ := json.Marshal(&data)
	resp, _ := blockchain.ChannelExecute("set", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
		dataByte,
	})
	t.Logf("set_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	t.Logf("set_data: %+v\n", resp.ChaincodeStatus)
	time.Sleep(2 * time.Second)
	result, statusCode := postForm("/archivePaper", dataByte, routers)
	if statusCode != 200 {
		expectApi(statusCode, 5, "Test_ArchivePaper")
		t.Errorf("archivePaper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.FailNow()
	}
	time.Sleep(2 * time.Second)
	resp, errString := blockchain.ChannelQuery("get", [][]byte{
		[]byte("paper"),
		[]byte(data.Id),
	})
	t.Logf("get_data: %+v\n", errString)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(resp.Payload).String())
	var paper lib.Paper
	_ = json.Unmarshal(resp.Payload, &paper)
	if paper.Status == lib.PaperStatusArchive {
		expectApi(statusCode, 4, "Test_ArchivePaper")
	} else {
		expectApi(statusCode, 5, "Test_ArchivePaper")
		t.Error("归档失败")
		t.FailNow()
	}
}

// Test_GetPapers01 学生是否可以查看自己所有的论文
func Test_GetPapers01(t *testing.T) {
	result, statusCode := get("/papers?type=student&student=liangsheng", routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 8 {
		expectApi(statusCode, 4, "Test_GetPapers01")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers01")
		t.Error("失败")
		t.FailNow()
	}
}

// Test_GetPapers2 答辩列表只能看到提交给该老师的提交状态的论文
func Test_GetPapers2(t *testing.T) {
	result, statusCode := get("/papers?type=teacher&teacher=teacher01&status="+lib.PaperStatusOralDefense, routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 1 {
		expectApi(statusCode, 4, "Test_GetPapers2")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers2")
		t.Error("失败")
		t.FailNow()
	}
}

// Test_GetPapers3 论文审核列表只能查看需要当前老师审核的论文
func Test_GetPapers3(t *testing.T) {
	result, statusCode := get("/papers?type=teacher&teacher=teacher01&status="+lib.PaperStatusSubmit, routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 2 {
		expectApi(statusCode, 4, "Test_GetPapers3")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers3")
		t.Error("失败")
		t.FailNow()
	}
}

// Test_GetPapers4 学校可以查看所有论文
func Test_GetPapers4(t *testing.T) {
	result, statusCode := get("/papers?type=school", routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 9 {
		expectApi(statusCode, 4, "Test_GetPapers4")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers4")
		t.Error("失败")
		t.FailNow()
	}
}

// Test_GetPapers5 学校可以查看答辩通过状态的所有论文
func Test_GetPapers5(t *testing.T) {
	result, statusCode := get("/papers?type=school&status="+lib.PaperStatusPassed, routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 1 {
		expectApi(statusCode, 4, "Test_GetPapers5")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers5")
		t.Error("失败")
		t.FailNow()
	}
}

// Test_GetPapers6 学校可以查看归档状态的所有论文
func Test_GetPapers6(t *testing.T) {
	result, statusCode := get("/papers?type=school&status="+lib.PaperStatusArchive, routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 1 {
		expectApi(statusCode, 4, "Test_GetPapers6")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers6")
		t.Error("失败")
		t.FailNow()
	}
}

// Test_GetPapers7 学校可以查看导师通过状态的所有论文
func Test_GetPapers7(t *testing.T) {
	result, statusCode := get("/papers?type=school&status="+lib.PaperStatusApprove, routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 2 {
		expectApi(statusCode, 4, "Test_GetPapers7")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers7")
		t.Error("失败")
		t.FailNow()
	}
}

// Test_GetPapers8 学校可以查看答辩状态的所有论文
func Test_GetPapers8(t *testing.T) {
	result, statusCode := get("/papers?type=school&status="+lib.PaperStatusOralDefense, routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 2 {
		expectApi(statusCode, 4, "Test_GetPapers8")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers8")
		t.Error("失败")
		t.FailNow()
	}
}

// Test_GetPapers9 学校可以查看新建状态的所有论文
func Test_GetPapers9(t *testing.T) {
	result, statusCode := get("/papers?type=school&status="+lib.PaperStatusNew, routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 1 {
		expectApi(statusCode, 4, "Test_GetPapers9")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers9")
		t.Error("失败")
		t.FailNow()
	}
}

// Test_GetPapers10 学校可以查看提交状态的所有论文
func Test_GetPapers10(t *testing.T) {
	result, statusCode := get("/papers?type=school&status="+lib.PaperStatusSubmit, routers)
	if statusCode != 200 {
		t.Errorf("paper 函数运行错误: %+v\n", bytes.NewBuffer(result).String())
		t.Fail()
	}
	t.Logf("get_data: %+v\n", statusCode)
	t.Logf("get_data: %+v\n", bytes.NewBuffer(result).String())
	var papers []lib.Paper
	_ = json.Unmarshal(result, &papers)
	fmt.Println(len(papers))
	if len(papers) == 2 {
		expectApi(statusCode, 4, "Test_GetPapers10")
	} else {
		expectApi(statusCode, 5, "Test_GetPapers10")
		t.Error("失败")
		t.FailNow()
	}
}
