package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"os"
	"sync"
	"testing"
	"time"
)

var mutex sync.RWMutex

func expectApi(status int, testFunc string) {
	str := "******"
	file := "./chaincode_test_result.txt"
	switch status {
	case 1:
		str += "测试通过：pass"
		mutex.Lock()
		writeFile(file, testFunc+":PASS\n")
		mutex.Unlock()
	case 2:
		str += "测试失败：fail"
		mutex.Lock()
		writeFile(file, testFunc+":FAIL\n")
		mutex.Unlock()
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

//TODO -----------------------新增人员---------------------------

var persons = []Person{
	{
		Id:       "",
		Account:  "Jack",
		Password: "123456",
		Name:     "Jack",
		Address:  "北京",
		Type:     Student,
	},
	{
		Id:       "1001",
		Account:  "",
		Password: "123456",
		Name:     "Jack",
		Address:  "北京",
		Type:     Student,
	},
	{
		Id:       "1001",
		Account:  "Jack",
		Password: "",
		Name:     "Jack",
		Address:  "北京",
		Type:     Student,
	},
	{
		Id:       "1001",
		Account:  "Jack",
		Password: "123456",
		Name:     "",
		Address:  "北京",
		Type:     Student,
	},

	{
		Id:       "1001",
		Account:  "Jack",
		Password: "123456",
		Name:     "Jack",
		Address:  "北京",
		Type:     "",
	},
}

// 人员状态是否更改为新增
func Test_newPerson1(t *testing.T) {
	fmt.Println("正确用例")

	stub := shim.NewMockStub("test", new(EduMgmt))

	person := Person{
		Id:         "1001",
		Account:    "Jack",
		Password:   "123456",
		Phone:      "13526889965",
		Email:      "jack@126.com",
		Name:       "Jack",
		Address:    "北京",
		College:    "北京大学",
		Department: "计算机系",
		Class:      "计算机1班",
		Type:       Student,
		CreatedAt:  time.Now(),
	}

	marshal, _ := json.Marshal(person)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("register"),
		marshal,
	})

	t.Log(resp.Status)
	t.Log(resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	key, _ := stub.CreateCompositeKey(person.Type, []string{person.Account})
	state, _ := stub.GetState(key)
	personTag := Person{}
	_ = json.Unmarshal(state, &personTag)
	t.Logf("%+v\n", personTag)

	if personTag.Status != personStatusNew {
		expectApi(2, "Test_newPerson1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_newPerson1")
}

// 参数非空校验
func Test_newPersont6(t *testing.T) {

	stub := shim.NewMockStub("test", new(EduMgmt))
	for _, person := range persons {
		marshal, _ := json.Marshal(person)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("register"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log(resp.Message)
		t.Log("返回值:" + string(resp.Payload))

		if resp.Status == 200 {
			expectApi(2, "Test_newPersont6")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_newPersont6")

}

//账号类型错误或账号是否存在
func Test_newPersont7(t *testing.T) {

	fmt.Println("账号信息错误或账号已存在")

	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	account := "Jack"
	context := "{\"" +
		"id\":\"10001\",\"" +
		"account\":\"" + account + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"Jack\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + Student + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	compositeKey, _ := stub.CreateCompositeKey(Student, []string{account})
	_ = stub.PutState(compositeKey, []byte(context))

	person := Person{
		Id:         "1001",
		Account:    account,
		Password:   "123456",
		Phone:      "13526889965",
		Email:      "jack@126.com",
		Name:       "Jack",
		Address:    "北京",
		College:    "北京大学",
		Department: "计算机系",
		Class:      "计算机1班",
		Type:       Student,
		Status:     personStatusNew,
	}

	marshal, _ := json.Marshal(person)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("register"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log(resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_newPersont7")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_newPersont7")
}

//验证学号/工号是否重复
func Test_newPersont8(t *testing.T) {

	fmt.Println("学号/工号已存在, 请重新输入")

	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	account := "May"
	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"account\":\"" + account + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"May@126.com\",\"" +
		"name\":\"May\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + Student + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	compositeKey, _ := stub.CreateCompositeKey(Student, []string{account})
	_ = stub.PutState(compositeKey, []byte(context))

	person := Person{
		Id:         id,
		Account:    "Jack",
		Password:   "123456",
		Phone:      "13526889965",
		Email:      "jack@126.com",
		Name:       "Jack",
		Address:    "北京",
		College:    "北京大学",
		Department: "计算机系",
		Class:      "计算机1班",
		Type:       Student,
		Status:     personStatusNew,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}
	marshal, _ := json.Marshal(person)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("register"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log(resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_newPersont8")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_newPersont8")
}

//TODO -----------------------updateScore测试单例---------------------------

var papers = []Paper{
	{
		Id:      "",
		Student: "Jack",
		Teacher: "Jhon",
		Score:   88,
	},
	{
		Id:      "10001",
		Student: "",
		Teacher: "Jhon",
		Score:   88,
	},
}

//参数非空校验
func Test_updateScore3(t *testing.T) {

	fmt.Println("参数非空校验")

	stub := shim.NewMockStub("test", new(EduMgmt))
	for _, paper := range papers {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("markPaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("Message:" + resp.Message)
		t.Log("返回值:" + string(resp.Payload))

		if resp.Status == 200 {
			expectApi(2, "Test_updateScore3")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_updateScore3")
}

var papers2 = []Paper{
	{
		Id:      "10001",
		Student: "Jack",
		Teacher: "Jhon",
		Score:   808,
	},
	{
		Id:      "10001",
		Student: "Jack",
		Teacher: "Jhon",
		Score:   -88,
	},
}

//验证评分是否在0-100分之间
func Test_updateScore6(t *testing.T) {

	fmt.Println("评分大于100")
	var stub = shim.NewMockStub("test", new(EduMgmt))

	for _, paper := range papers2 {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("markPaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("Message:" + resp.Message)
		t.Log("返回值:" + string(resp.Payload))

		if resp.Status == 200 {
			expectApi(2, "Test_updateScore6")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_updateScore6")
}

//验证论文是否存在
func Test_updateScore7(t *testing.T) {

	fmt.Println("论文不存在, 请输入正确的论文ID")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":10,\"" +
		"Mary\":20,\"" +
		"Jhon\":-1" +
		"},\"" +
		"status\":\"" + PaperStatusOralDefense + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      "2001",
		Student: "Jack",
		Teacher: "Jhon",
		Score:   88,
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("markPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	if resp.Status == 200 {
		expectApi(2, "Test_updateScore7")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_updateScore7")
}

//只能对答辩状态的论文评分
func Test_updateScore8(t *testing.T) {

	fmt.Println("只能对答辩状态的论文评分")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":10,\"" +
		"Mary\":20,\"" +
		"Jhon\":-1" +
		"},\"" +
		"status\":\"" + PaperStatusApprove + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "Jack",
		Teacher: "Jhon",
		Score:   88,
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("markPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	if resp.Status == 200 {
		expectApi(2, "Test_updateScore8")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_updateScore8")
}

//不能重复评分
func Test_updateScore9(t *testing.T) {

	fmt.Println("不能重复评分")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":10,\"" +
		"Mary\":20,\"" +
		"Jhon\":40" +
		"},\"" +
		"status\":\"" + PaperStatusOralDefense + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "Jack",
		Teacher: "Jhon",
		Score:   88,
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("markPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	if resp.Status == 200 {
		expectApi(2, "Test_updateScore9")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_updateScore9")
}

//验证0分是否为有效分数
func Test_updateScore11(t *testing.T) {

	fmt.Println("验证0分是否为有效分数")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":10,\"" +
		"Mary\":20,\"" +
		"Jhon\":-1,\"" +
		"Jpon\":56,\"" +
		"Juon\":78" +
		"},\"" +
		"status\":\"" + PaperStatusOralDefense + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "Jack",
		Teacher: "Jhon",
		Score:   0,
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("markPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	key, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"10001"})
	state, _ := stub.GetState(key)
	paperTag := Paper{}
	e := json.Unmarshal(state, &paperTag)
	fmt.Println(e)

	if paperTag.Status == PaperStatusPassed || paperTag.Status == PaperStatusReject {
		expectApi(1, "Test_updateScore11")
	} else {
		expectApi(2, "Test_updateScore11")
		t.Error("error")
		t.FailNow()
	}
}

//当评分人数达到5人时，是否自动计算平均值，确定论文答辩分数
func Test_updateScore12(t *testing.T) {

	fmt.Println("验证0分是否为有效分数")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":10,\"" +
		"Mary\":20,\"" +
		"Jhon\":-1,\"" +
		"Jpon\":56,\"" +
		"Juon\":78" +
		"},\"" +
		"status\":\"" + PaperStatusOralDefense + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "Jack",
		Teacher: "Jhon",
		Score:   80,
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("markPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	key, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"10001"})
	state, _ := stub.GetState(key)
	paperTag := Paper{}
	e := json.Unmarshal(state, &paperTag)
	fmt.Println(e)

	if paperTag.Status == PaperStatusPassed || paperTag.Status == PaperStatusReject {
		expectApi(1, "Test_updateScore12")
	} else {
		expectApi(2, "Test_updateScore12")
		t.Error("error")
		t.FailNow()
	}
}

//当答辩平均分为大于等于90分，是否更改状态为答辩通过
func Test_updateScore13(t *testing.T) {

	fmt.Println("验证0分是否为有效分数")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":91,\"" +
		"Mary\":91,\"" +
		"Jhon\":-1,\"" +
		"Jpon\":91,\"" +
		"Juon\":91" +
		"},\"" +
		"status\":\"" + PaperStatusOralDefense + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "Jack",
		Teacher: "Jhon",
		Score:   91,
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("markPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	key, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"10001"})
	state, _ := stub.GetState(key)
	paperTag := Paper{}
	_ = json.Unmarshal(state, &paperTag)

	if paperTag.Status != PaperStatusPassed {
		expectApi(2, "Test_updateScore13")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_updateScore13")
}

//当答辩平均分为小于90分时，是否更改状态为退回
func Test_updateScore14(t *testing.T) {

	fmt.Println("验证0分是否为有效分数")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":89,\"" +
		"Mary\":85,\"" +
		"Jhon\":-1,\"" +
		"Jpon\":45,\"" +
		"Juon\":63" +
		"},\"" +
		"status\":\"" + PaperStatusOralDefense + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "Jack",
		Teacher: "Jhon",
		Score:   22,
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("markPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	key, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"10001"})
	state, _ := stub.GetState(key)
	paperTag := Paper{}
	_ = json.Unmarshal(state, &paperTag)

	if paperTag.Status != PaperStatusReject {
		expectApi(2, "Test_updateScore14")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_updateScore14")
}

//TODO -----------------------deletePaper测试单例---------------------------
var papers3 = []Paper{
	{
		Id:      "",
		Student: "Jack",
	},
	{
		Id:      "10001",
		Student: "",
	},
}

//参数非空校验
func Test_deletePaper2(t *testing.T) {

	fmt.Println("参数非空校验")
	var stub = shim.NewMockStub("test", new(EduMgmt))

	for _, paper := range papers3 {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("deletePaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("Message:" + resp.Message)
		t.Log("返回值:" + string(resp.Payload))

		if resp.Status == 200 {
			expectApi(2, "Test_deletePaper2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_deletePaper2")
}

//验证论文是否存在
func Test_deletePaper4(t *testing.T) {

	fmt.Println("论文不存在")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":10,\"" +
		"Mary\":20,\"" +
		"Jhon\":40" +
		"},\"" +
		"status\":\"" + PaperStatusReject + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      "id",
		Student: "Jack",
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("deletePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	if resp.Status == 200 {
		expectApi(2, "Test_deletePaper4")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_deletePaper4")
}

//验证是否只能删除新建或退回状态的论文
func Test_deletePaper6(t *testing.T) {

	fmt.Println("只能删除新建或退回状态的论文")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":10,\"" +
		"Mary\":20,\"" +
		"Jhon\":40" +
		"},\"" +
		"status\":\"" + PaperStatusApprove + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "Jack",
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("deletePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	if resp.Status == 200 {
		expectApi(2, "Test_deletePaper6")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_deletePaper6")
}

//是否删除成功
func Test_deletePaper1(t *testing.T) {

	fmt.Println("正确单例")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "10001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"海贼王\",\"" +
		"abstract\":\"海贼王世界观\",\"" +
		"score\":80,\"" +
		"description\":\"海贼王世界观\",\"" +
		"student\":\"Jack\",\"" +
		"teacher\":\"Mary\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":10,\"" +
		"Mary\":20,\"" +
		"Jhon\":40" +
		"},\"" +
		"status\":\"" + PaperStatusReject + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix,
		[]string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "Jack",
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("deletePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	key, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"10001"})
	state, _ := stub.GetState(key)
	p := new(Paper)
	_ = json.Unmarshal(state, p)
	if p.Status != PaperStatusCancel {
		expectApi(2, "Test_deletePaper1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_deletePaper1")
}

//TODO -----------------------getPerson测试单例---------------------------

var persons2 = []Person{
	{
		Account:  "Jack",
		Password: "123456",
		Type:     "",
	},

	{
		Account:  "",
		Password: "123456",
		Type:     School,
	},
	{
		Account:  "Jack",
		Password: "",
		Type:     School,
	},
}

//参数非空校验
func Test_getPerson1(t *testing.T) {

	fmt.Println("参数非空校验")
	var stub = shim.NewMockStub("test", new(EduMgmt))

	for _, person := range persons2 {
		marshal, _ := json.Marshal(person)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("getPerson"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("Message:" + resp.Message)
		t.Log("返回值:" + string(resp.Payload))
		if resp.Status == 200 {
			expectApi(2, "Test_getPerson1")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_getPerson1")
}

// 账户不存在
func Test_getPerson5(t *testing.T) {

	fmt.Println("账户不存在")
	stub := shim.NewMockStub("test", new(EduMgmt))
	account := "Jack"
	context := "{\"" +
		"id\":\"10001\",\"" +
		"account\":\"" + account + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"Jack\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + Student + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(School, []string{account})
	_ = stub.PutState(compositeKey, []byte(context))

	person := Person{
		Account:  "account",
		Password: "123456",
		Type:     School,
	}

	marshal, _ := json.Marshal(person)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("getPerson"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_getPerson5")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_getPerson5")
}

//TODO -----------------------modPerson测试单例---------------------------
var persons4 = []Person{
	{
		Account:  "Jack",
		Type:     "",
		Password: Student,
		Email:    "Mary",
	},

	{
		Account:  "",
		Type:     School,
		Password: Student,
		Email:    "Mary",
	},

	{
		Account:  "Jack",
		Type:     School,
		Password: "",
		Email:    "Mary",
	},

	{
		Account:  "Jack",
		Type:     School,
		Password: School,
		Email:    "",
	},
}

//参数非空校验
func Test_modPerson1(t *testing.T) {

	fmt.Println("参数非空校验")
	var stub = shim.NewMockStub("test", new(EduMgmt))

	for _, person := range persons4 {
		marshal, _ := json.Marshal(person)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("confirmUser"),
			marshal,
			[]byte(person.Password),
			[]byte(person.Email),
		})
		t.Log(resp.Status)
		t.Log("Message:" + resp.Message)
		t.Log("返回值：" + string(resp.Payload))
		if resp.Status == 200 {
			expectApi(2, "Test_modPerson1")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_modPerson1")
}

//验证当前账号是否存在
func Test_modPerson4(t *testing.T) {

	fmt.Println("验证当前账号是否存在")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	// 构建school账户
	account := "Jack"
	context := "{\"" +
		"id\":\"1001\",\"" +
		"account\":\"" + account + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"Jack\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + School + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	compositeKey, _ := stub.CreateCompositeKey(School, []string{account})
	_ = stub.PutState(compositeKey, []byte(context))

	// 构建新建的账户
	account2 := "liangsheng"
	context2 := "{\"" +
		"id\":\"1002\",\"" +
		"account\":\"" + account2 + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"liangsheng\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + Student + "\",\"" +
		"status\":\"" + personStatusNew + "\"}"

	compositeKey2, _ := stub.CreateCompositeKey(Student, []string{account2})
	_ = stub.PutState(compositeKey2, []byte(context2))

	person := Person{
		Id:      "1001",
		Account: "accouniuhut",
		Type:    School,
	}
	marshal, _ := json.Marshal(person)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("confirmUser"),
		marshal,
		[]byte(Student),
		[]byte(account2),
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值：" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_modPerson4")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_modPerson4")
}

//验证需要修改的用户是否存在
func Test_modPerson7(t *testing.T) {

	fmt.Println("当前用户没有此权限")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	// 构建school账户
	account := "Jack"
	context := "{\"" +
		"id\":\"1001\",\"" +
		"account\":\"" + account + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"Jack\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + School + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	compositeKey, _ := stub.CreateCompositeKey(School, []string{account})
	_ = stub.PutState(compositeKey, []byte(context))

	// 构建新建的账户
	account2 := "liangsheng"
	context2 := "{\"" +
		"id\":\"1002\",\"" +
		"account\":\"" + account2 + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"liangsheng\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + Student + "\",\"" +
		"status\":\"" + personStatusNew + "\"}"

	compositeKey2, _ := stub.CreateCompositeKey(Student, []string{account2})
	_ = stub.PutState(compositeKey2, []byte(context2))

	person := Person{
		Account: account,
		Type:    School,
	}
	marshal, _ := json.Marshal(person)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("confirmUser"),
		marshal,
		[]byte(Student),
		[]byte("account2"),
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值：" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_modPerson7")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_modPerson7")
}

//人员状态是否更新为已确认和修改论文更新时间
func Test_modPerson5(t *testing.T) {

	fmt.Println("targetAccount为空")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	// 构建school账户
	account := "Jack"
	context := "{\"" +
		"id\":\"1001\",\"" +
		"account\":\"" + account + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"Jack\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + School + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	compositeKey, _ := stub.CreateCompositeKey(School, []string{account})
	_ = stub.PutState(compositeKey, []byte(context))

	// 构建新建的账户
	account2 := "liangsheng"
	context2 := "{\"" +
		"id\":\"1002\",\"" +
		"account\":\"" + account2 + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"liangsheng\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + Student + "\",\"" +
		"status\":\"" + personStatusNew + "\"}"

	compositeKey2, _ := stub.CreateCompositeKey(Student, []string{account2})
	_ = stub.PutState(compositeKey2, []byte(context2))

	person := Person{
		Id:        "1001",
		Account:   account,
		Type:      School,
		UpdatedAt: time.Now(),
	}
	marshal, _ := json.Marshal(person)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("confirmUser"),
		marshal,
		[]byte(Student),
		[]byte(account2),
	})

	state, _ := stub.GetState(compositeKey2)
	var person2 Person
	_ = json.Unmarshal(state, &person2)

	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值：" + string(resp.Payload))

	if person2.Status != personStatusConfirm {
		expectApi(2, "Test_modPerson5")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_modPerson5")
}

//TODO -----------------------getPersonList测试单例---------------------------
//是否返回所有人员列表
func Test_getPersonList1(t *testing.T) {
	fmt.Println("是否返回所有人员列表")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	//构建用户
	account := "Jack"
	context := "{\"" +
		"id\":\"1001\",\"" +
		"account\":\"" + account + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"Jack\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + Student + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	compositeKey, _ := stub.CreateCompositeKey(Student, []string{account})
	_ = stub.PutState(compositeKey, []byte(context))

	account2 := "Jack"
	context2 := "{\"" +
		"id\":\"2001\",\"" +
		"account\":\"" + account2 + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"mary@126.com\",\"" +
		"name\":\"Mary\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + School + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	compositeKey2, _ := stub.CreateCompositeKey(School, []string{account2})
	_ = stub.PutState(compositeKey2, []byte(context2))

	account3 := "Kill"
	context3 := "{\"" +
		"id\":\"1002\",\"" +
		"account\":\"" + account2 + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"kill@126.com\",\"" +
		"name\":\"Kill\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + Student + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	compositeKey3, _ := stub.CreateCompositeKey(Student, []string{account3})
	_ = stub.PutState(compositeKey3, []byte(context3))

	person := Person{
		Type: Student,
	}

	marshal, _ := json.Marshal(person)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("getPersons"),
		marshal,
	})

	var personsTag []Person

	t.Log(resp.Status)
	t.Log("返回值：" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	_ = json.Unmarshal(resp.Payload, &personsTag)

	if len(personsTag) != 2 {
		expectApi(2, "Test_getPersonList1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_getPersonList1")
}

//参数非空校验
func Test_getPersonList2(t *testing.T) {
	fmt.Println("参数非空校验")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	//构建用户
	account := "Jack"
	context := "{\"" +
		"id\":\"1001\",\"" +
		"account\":\"" + account + "\",\"" +
		"password\":\"123456\",\"" +
		"phone\":\"13526889965\",\"" +
		"email\":\"jack@126.com\",\"" +
		"name\":\"Jack\",\"" +
		"address\":\"北京\",\"" +
		"college\":\"北京大学\",\"" +
		"department\":\"计算机系\",\"" +
		"class\":\"计算机1班\",\"" +
		"type\":\"" + Student + "\",\"" +
		"status\":\"" + personStatusConfirm + "\"}"

	compositeKey, _ := stub.CreateCompositeKey(Student, []string{account})
	_ = stub.PutState(compositeKey, []byte(context))

	person := Person{
		Type: "",
	}

	marshal, _ := json.Marshal(person)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("getPersons"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("返回值：" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_getPersonList2")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_getPersonList2")
}

//TODO -----------------------newPaper测试单例---------------------------

var papers5 = []Paper{
	{
		Id:          "",
		Title:       "海贼王",
		Abstract:    "拥有金钱、名声、实力和自由的男人海贼王",
		Score:       99,
		Description: "海贼的冒险故事",
		Student:     "尾田荣一郎",
		Teacher:     "宫崎骏",
		Department:  "漫画系",
	},

	{
		Id:          "p1002",
		Title:       "",
		Abstract:    "拥有金钱、名声、实力和自由的男人海贼王",
		Score:       99,
		Description: "海贼的冒险故事",
		Student:     "尾田荣一郎",
		Teacher:     "宫崎骏",
		Department:  "漫画系",
	},
	{
		Id:          "p1002",
		Title:       "海贼王",
		Abstract:    "",
		Score:       99,
		Description: "海贼的冒险故事",
		Student:     "尾田荣一郎",
		Teacher:     "宫崎骏",
		Department:  "漫画系",
	},

	{
		Id:          "p1002",
		Title:       "海贼王",
		Abstract:    "拥有金钱、名声、实力和自由的男人海贼王",
		Score:       99,
		Description: "海贼的冒险故事",
		Student:     "",
		Teacher:     "宫崎骏",
		Department:  "漫画系",
	},

	{
		Id:          "p1002",
		Title:       "海贼王",
		Abstract:    "拥有金钱、名声、实力和自由的男人海贼王",
		Score:       99,
		Description: "海贼的冒险故事",
		Student:     "尾田荣一郎",
		Teacher:     "",
		Status:      PaperStatusNew,
		Department:  "漫画系",
	},
}

//参数非空校验
func Test_newPaper2(t *testing.T) {
	fmt.Println("参数非空校验")
	var stub = shim.NewMockStub("test", new(EduMgmt))
	for _, paper := range papers5 {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("newPaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log(resp.Message)
		t.Log("返回值:" + string(resp.Payload))
		if resp.Status == 200 {
			expectApi(2, "Test_newPaper2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_newPaper2")
}

//论文状态是否为新建
func Test_newPaper1(t *testing.T) {
	fmt.Println("论文状态是否为新建")
	stub := shim.NewMockStub("test", new(EduMgmt))
	paper := Paper{
		Id:          "p1002",
		Title:       "海贼王",
		Abstract:    "拥有金钱、名声、实力和自由的男人海贼王",
		Score:       99,
		Description: "海贼的冒险故事",
		Student:     "尾田荣一郎",
		Teacher:     "宫崎骏",
		Department:  "漫画系",
		CreatedAt:   time.Now(),
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("newPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log(resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"p1002"})
	state, _ := stub.GetState(compositeKey)
	paperTag := Paper{}
	_ = json.Unmarshal(state, &paperTag)

	if paperTag.Status != PaperStatusNew {
		expectApi(2, "Test_newPaper1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_newPaper1")
}

func Test_newPaper3(t *testing.T) {
	expectApi(1, "Test_newPaper3")
}

//验证论文title是否已存在
func Test_newPaper0(t *testing.T) {
	stub := shim.NewMockStub("test", new(EduMgmt))
	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          "p1004",
		Title:       "火影忍者",
		Abstract:    "摘要",
		Description: "test",
		Student:     "真岛浩",
		Teacher:     "宫崎骏",
		Status:      PaperStatusNew,
		Department:  "漫画系",
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("newPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("返回值:" + resp.Message)
	if resp.Status == 200 {
		expectApi(2, "Test_newPaper0")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_newPaper0")
}

// TODO -----------------------modPaper测试单例---------------------------
var papers6 = []Paper{
	{
		Id:          "",
		Title:       "火影忍者",
		Abstract:    "热血的忍者故事",
		Description: "忍者成长经历",
		Student:     "岸本齐史",
	},
	{
		Id:          "p1001",
		Title:       "",
		Abstract:    "热血的忍者故事",
		Description: "忍者成长经历",
		Student:     "岸本齐史",
	},

	{
		Id:          "p1001",
		Title:       "火影忍者",
		Abstract:    "热血的忍者故事",
		Description: "忍者成长经历",
		Student:     "",
	},
}

//参数非空校验
func Test_modPaper2(t *testing.T) {
	fmt.Println("参数非空校验")
	var stub = shim.NewMockStub("test", new(EduMgmt))

	for _, paper := range papers6 {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("updatePaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("返回值:" + resp.Message)
		if resp.Status == 200 {
			expectApi(2, "Test_modPaper2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_modPaper2")
}

// 验证修改是否成功
func Test_modPaper1(t *testing.T) {
	fmt.Println("验证修改是否成功")
	stub := shim.NewMockStub("test", new(EduMgmt))
	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"标题\",\"" +
		"abstract\":\"摘要\",\"" +
		"score\":89,\"" +
		"description\":\"描述\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          "p1001",
		Title:       "火影忍者",
		Abstract:    "热血的忍者故事",
		Description: "忍者成长经历",
		Student:     "岸本齐史",
		UpdatedAt:   time.Now(),
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("updatePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("返回值:" + resp.Message)
	if resp.Status != 200 {
		expectApi(2, "Test_modPaper1")
		t.Error("error")
	}

	// 验证是否修改成功
	state, _ := stub.GetState(compositeKey)
	var paper2 Paper
	_ = json.Unmarshal(state, &paper2)
	t.Log(fmt.Sprintf("修改前的账户：%s\n", context))
	t.Log(fmt.Sprintf("修改后的标题：%s, 摘要：%s，描述：%s", paper2.Title,
		paper2.Abstract, paper2.Description))
	if paper2.Title != paper.Title || paper2.Abstract != paper.Abstract || paper2.Description != paper.Description {
		expectApi(2, "Test_modPaper1")
		t.FailNow()
	}
	expectApi(1, "Test_modPaper1")
}

//验证论文是否存在
func Test_modPaper5(t *testing.T) {
	fmt.Println("验证论文是否存在")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"标题\",\"" +
		"abstract\":\"摘要\",\"" +
		"score\":89,\"" +
		"description\":\"描述\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          "p1002",
		Title:       "火影忍者",
		Abstract:    "热血的忍者故事",
		Description: "忍者成长经历",
		Student:     "岸本齐史",
		UpdatedAt:   time.Now(),
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("updatePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("返回值:" + resp.Message)
	if resp.Status == 200 {
		expectApi(2, "Test_modPaper5")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_modPaper5")
}

//只能修改新建或退回状态的论文
func Test_modPaper6(t *testing.T) {
	fmt.Println("只能修改新建或退回状态的论文")
	stub := shim.NewMockStub("test", new(EduMgmt))

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"标题\",\"" +
		"abstract\":\"摘要\",\"" +
		"score\":89,\"" +
		"description\":\"描述\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusArchive + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")
	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          id,
		Title:       "火影忍者",
		Abstract:    "热血的忍者故事",
		Description: "忍者成长经历",
		Student:     "岸本齐史",
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("updatePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("返回值:" + resp.Message)
	if resp.Status == 200 {
		expectApi(2, "Test_modPaper6")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_modPaper6")
}

func Test_modPaper9(t *testing.T) {
	expectApi(1, "Test_modPaper9")
}

//论文Title已存在
func Test_modPaper8(t *testing.T) {
	fmt.Println("论文Title已存在")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"标题\",\"" +
		"abstract\":\"摘要\",\"" +
		"score\":89,\"" +
		"description\":\"描述\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	id2 := "p1002"
	context2 := "{\"" +
		"id\":\"" + id2 + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"摘要\",\"" +
		"score\":89,\"" +
		"description\":\"描述\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey2, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id2})
	_ = stub.PutState(compositeKey2, []byte(context2))

	paper := Paper{
		Id:          id,
		Title:       "火影忍者",
		Abstract:    "热血的忍者故事",
		Description: "忍者成长经历",
		Student:     "岸本齐史",
	}
	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("updatePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("返回值:" + resp.Message)
	if resp.Status == 200 {
		expectApi(2, "Test_modPaper8")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_modPaper8")
}

//TODO -----------------------getPaper测试单例---------------------------

var papers7 = []Paper{
	{
		Id:      "",
		Student: "岸本齐史",
	},

	{
		Id:      "p1001",
		Student: "",
	},
}

//验证获取论文是否成功
func Test_getPaper1(t *testing.T) {

	fmt.Println("验证获取论文是否成功")

	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "岸本齐史",
	}
	marshal, _ := json.Marshal(paper)

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("getPaper"),
		marshal,
	})

	t.Log(resp.Status)
	t.Log(resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status != 200 {
		expectApi(2, "Test_getPaper1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_getPaper1")
}

//参数非空校验
func Test_getPaper2(t *testing.T) {

	fmt.Println("参数非空校验")
	var stub = shim.NewMockStub("test", new(EduMgmt))

	for paper := range papers7 {
		marshal, _ := json.Marshal(paper)

		resp := stub.MockInvoke("1", [][]byte{
			[]byte("getPaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log(resp.Message)
		t.Log("返回值:" + string(resp.Payload))
		if resp.Status == 200 {
			expectApi(2, "Test_getPaper2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_getPaper2")
}

//验证论文是否存在
func Test_getPaper4(t *testing.T) {

	fmt.Println("验证论文是否存在")

	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      "256",
		Student: "岸本齐史",
	}
	marshal, _ := json.Marshal(paper)

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("getPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log(resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_getPaper4")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_getPaper4")
}

//TODO -----------------------getPapers测试单例---------------------------
//验证是否返回所有论文
func Test_getPapers(t *testing.T) {

	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	id2 := "p1002"
	context2 := "{\"" +
		"id\":\"" + id2 + "\",\"" +
		"title\":\"妖精尾巴\",\"" +
		"abstract\":\"魔法师故事\",\"" +
		"score\":89,\"" +
		"description\":\"魔法师成长经历\",\"" +
		"student\":\"真岛浩\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":12" +
		"},\"" +
		"status\":\"" + PaperStatusArchive + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	compositeKey2, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id2})
	_ = stub.PutState(compositeKey2, []byte(context2))

	resp := stub.MockInvoke("1", [][]byte{
		[]byte("getPapers"),
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status != 200 {
		expectApi(2, "Test_getPapers")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_getPapers")
}

//提交论文是否成功
func Test_submitPaper1(t *testing.T) {

	fmt.Println("正确用例")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:         id,
		Student:    "岸本齐史",
		SubmitTime: time.Now(),
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("submitPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	compositeKey2, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"p1001"})
	state2, _ := stub.GetState(compositeKey2)

	paperTag := Paper{}
	_ = json.Unmarshal(state2, &paperTag)

	if paperTag.Status != PaperStatusSubmit {
		expectApi(2, "Test_submitPaper1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_submitPaper1")
}

//参数非空校验
func Test_submitPaper2(t *testing.T) {

	fmt.Println("参数非空校验")
	stub := shim.NewMockStub("test", new(EduMgmt))

	for paper := range papers7 {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("submitPaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("返回值:" + resp.Message)
		if resp.Status == 200 {
			expectApi(2, "Test_submitPaper2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_submitPaper2")
}

//验证论文是否存在
func Test_submitPaper4(t *testing.T) {

	fmt.Println("验证论文是否存在")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      "id",
		Student: "岸本齐史",
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("submitPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_submitPaper4")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_submitPaper4")
}

//只能提交新建状态或退回状态的论文
func Test_submitPaper6(t *testing.T) {

	fmt.Println("只能提交新建状态或退回状态的论文")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusArchive + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "岸本齐史",
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("submitPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_submitPaper6")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_submitPaper6")
}

//TODO -----------------------rejectPaper测试单例---------------------------
var t time.Time
var papers9 = []Paper{
	{
		Id:      "",
		Student: "岸本齐史",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
		RejectTime: time.Now(),
	},
	{
		Id:      "p1001",
		Student: "",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
		RejectTime: time.Now(),
	},
	{
		Id:      "p1001",
		Student: "岸本齐史",
		Teacher: "",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
		RejectTime: time.Now(),
	},
	{
		Id:         "p1001",
		Student:    "岸本齐史",
		Teacher:    "宫崎骏",
		RefuseInfo: "",
		RejectTime: time.Now(),
	},
	{
		Id:      "p1001",
		Student: "岸本齐史",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
	},
}

//参数非空校验
func Test_rejectPaper2(t *testing.T) {

	fmt.Println("参数非空校验")
	var stub = shim.NewMockStub("test", new(EduMgmt))

	for _, paper := range papers9 {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("rejectPaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("Message:" + resp.Message)
		t.Log("返回值:" + string(resp.Payload))
		if resp.Status != 400 {
			expectApi(2, "Test_rejectPaper2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_rejectPaper2")
}

//论文退回是否成功
func Test_rejectPaper1(t *testing.T) {

	fmt.Println("论文退回是否成功")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusSubmit + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "岸本齐史",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
		RejectTime: time.Now(),
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("rejectPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	compositeKey2, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"p1001"})
	state2, _ := stub.GetState(compositeKey2)

	paperTag := Paper{}
	_ = json.Unmarshal(state2, &paperTag)

	if paperTag.Status != PaperStatusReject {
		expectApi(2, "Test_rejectPaper1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_rejectPaper1")
}

//验证论文是否存在
func Test_rejectPaper6(t *testing.T) {

	fmt.Println("验证论文是否存在")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusSubmit + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      "id",
		Student: "岸本齐史",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("rejectPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_rejectPaper6")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_rejectPaper6")
}

//只能退回提交状态的论文
func Test_rejectPaper9(t *testing.T) {

	fmt.Println("只能退回提交状态的论文")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + personStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "岸本齐史",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("rejectPaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_rejectPaper9")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_rejectPaper9")
}

//TODO -----------------------approvePaper测试单例---------------------------
var papers10 = []Paper{
	{
		Id:      "",
		Student: "岸本齐史",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
	},
	{
		Id:      "p1001",
		Student: "",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
	},

	{
		Id:      "p1001",
		Student: "岸本齐史",
		Teacher: "",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
	},
}

//论文通过是否成功
func Test_approvePaper1(t *testing.T) {

	fmt.Println("论文通过是否成功")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusSubmit + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "岸本齐史",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
		ApproveDate: time.Now(),
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("approvePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	compositeKey2, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"p1001"})
	state2, _ := stub.GetState(compositeKey2)

	paperTag := Paper{}
	_ = json.Unmarshal(state2, &paperTag)

	if paperTag.Status != PaperStatusApprove {
		expectApi(2, "Test_approvePaper1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_approvePaper1")
}

//参数非空校验
func Test_approvePaper2(t *testing.T) {

	fmt.Println("参数非空校验")
	stub := shim.NewMockStub("test", new(EduMgmt))

	for paper := range papers10 {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("approvePaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("Message:" + resp.Message)
		t.Log("返回值:" + string(resp.Payload))
		if resp.Status == 200 {
			expectApi(2, "Test_approvePaper2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_approvePaper2")
}

//论文不存在
func Test_approvePaper5(t *testing.T) {

	fmt.Println("论文不存在")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusSubmit + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      "id",
		Student: "岸本齐史",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("approvePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_approvePaper5")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_approvePaper5")
}

//只能通过提交状态的论文
func Test_approvePaper8(t *testing.T) {

	fmt.Println("只能通过提交状态的论文")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + personStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:      id,
		Student: "岸本齐史",
		Teacher: "宫崎骏",
		RefuseInfo: "{\"" +
			"Jack\":20,\"" +
			"Mary\":1" +
			"}",
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("approvePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_approvePaper8")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_approvePaper8")
}

//TODO -----------------------oralDefense测试单例---------------------------
var papers11 = []Paper{
	{
		Id:          "",
		Student:     "岸本齐史",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	},
	{
		Id:          "p1001",
		Student:     "",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	},
}

//安排答辩是否成功
func Test_oralDefense1(t *testing.T) {

	fmt.Println("安排答辩是否成功")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusApprove + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          id,
		Student:     "岸本齐史",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("oralDefense"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))

	compositeKey2, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"p1001"})
	state2, _ := stub.GetState(compositeKey2)

	paperTag := Paper{}
	_ = json.Unmarshal(state2, &paperTag)

	if paperTag.Status != PaperStatusOralDefense {
		expectApi(2, "Test_oralDefense1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_oralDefense1")
}

//参数非空校验
func Test_oralDefense2(t *testing.T) {

	fmt.Println("参数非空校验")
	stub := shim.NewMockStub("test", new(EduMgmt))

	for _, paper := range papers11 {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("oralDefense"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("Message:" + resp.Message)
		t.Log("返回值:" + string(resp.Payload))
		if resp.Status == 200 {
			expectApi(2, "Test_oralDefense2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_oralDefense2")
}

//至少安排5位老师
func Test_oralDefense4(t *testing.T) {

	fmt.Println("至少安排5位老师")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusApprove + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          id,
		Student:     "岸本齐史",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20},
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("oralDefense"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_oralDefense4")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_oralDefense4")
}

//论文不存在, 请输入正确的论文ID
func Test_oralDefense6(t *testing.T) {

	fmt.Println("论文不存在, 请输入正确的论文ID")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusApprove + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          "id",
		Student:     "岸本齐史",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("oralDefense"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_oralDefense6")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_oralDefense6")
}

//导师没有通过的论文无法安排答辩
func Test_oralDefense9(t *testing.T) {

	fmt.Println("导师没有通过的论文无法安排答辩")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          id,
		Student:     "岸本齐史",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("oralDefense"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_oralDefense9")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_oralDefense9")
}

//TODO -----------------------archivePaper测试单例---------------------------
var papers12 = []Paper{
	{
		Id:          "",
		Student:     "岸本齐史",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	},
	{
		Id:          "p1001",
		Student:     "",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	},
}

//论文归档是否成功
func Test_archivePaper1(t *testing.T) {

	fmt.Println("论文归档是否成功")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusPassed + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          id,
		Student:     "岸本齐史",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("archivePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	compositeKey2, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{"p1001"})
	state2, _ := stub.GetState(compositeKey2)

	paperTag := Paper{}
	_ = json.Unmarshal(state2, &paperTag)

	if paperTag.Status != PaperStatusArchive {
		expectApi(2, "Test_archivePaper1")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_archivePaper1")
}

//参数非空校验
func Test_archivePaper2(t *testing.T) {

	fmt.Println("id为空")
	var stub = shim.NewMockStub("test", new(EduMgmt))

	for paper := range papers12 {
		marshal, _ := json.Marshal(paper)
		resp := stub.MockInvoke("1", [][]byte{
			[]byte("archivePaper"),
			marshal,
		})
		t.Log(resp.Status)
		t.Log("Message:" + resp.Message)
		t.Log("返回值:" + string(resp.Payload))
		if resp.Status == 200 {
			expectApi(2, "Test_archivePaper2")
			t.Error("error")
			t.FailNow()
		}
	}
	expectApi(1, "Test_archivePaper2")
}

//论文不存在, 请输入正确的论文ID
func Test_archivePaper4(t *testing.T) {

	fmt.Println("论文不存在, 请输入正确的论文ID")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusPassed + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          "id",
		Student:     "岸本齐史",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("archivePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_archivePaper4")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_archivePaper4")
}

//答辩通过的论文才能归档
func Test_archivePaper6(t *testing.T) {

	fmt.Println("答辩通过的论文才能归档")
	stub := shim.NewMockStub("test", new(EduMgmt))

	stub.MockTransactionStart("1")
	defer stub.MockTransactionEnd("1")

	id := "p1001"
	context := "{\"" +
		"id\":\"" + id + "\",\"" +
		"title\":\"火影忍者\",\"" +
		"abstract\":\"热血的忍者故事\",\"" +
		"score\":89,\"" +
		"description\":\"忍者成长经历\",\"" +
		"student\":\"岸本齐史\",\"" +
		"teacher\":\"宫崎骏\",\"" +
		"department\":\"漫画系\",\"" +
		"oralDefense\":" +
		"{\"" +
		"Jack\":20,\"" +
		"Mary\":1" +
		"},\"" +
		"status\":\"" + PaperStatusNew + "\",\"" +
		"reject_count\":0,\"" +
		"refuse_info\":\"拒绝详情\"}"

	compositeKey, _ := stub.CreateCompositeKey(PaperKeyPrefix, []string{id})
	_ = stub.PutState(compositeKey, []byte(context))

	paper := Paper{
		Id:          id,
		Student:     "岸本齐史",
		OralDefense: map[string]int64{"Jack": 20, "Mary": 20, "Mack": 20, "Json": 20, "Key": 20},
	}

	marshal, _ := json.Marshal(paper)
	resp := stub.MockInvoke("1", [][]byte{
		[]byte("archivePaper"),
		marshal,
	})
	t.Log(resp.Status)
	t.Log("Message:" + resp.Message)
	t.Log("返回值:" + string(resp.Payload))
	if resp.Status == 200 {
		expectApi(2, "Test_archivePaper6")
		t.Error("error")
		t.FailNow()
	}
	expectApi(1, "Test_archivePaper6")
}
