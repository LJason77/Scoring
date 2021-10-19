package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	School  = "school"
	Student = "student"
	Teacher = "teacher"

	PaperKeyPrefix         = "paper"
	PaperStatusNew         = "新建"
	PaperStatusSubmit      = "提交"
	PaperStatusCancel      = "取消"
	PaperStatusReject      = "退回"
	PaperStatusApprove     = "导师通过"
	PaperStatusOralDefense = "答辩"
	PaperStatusPassed      = "答辩通过"
	PaperStatusArchive     = "归档"

	TimeFormat = "2006-01-02 15:04:05"

	personStatusNew     = "new"
	personStatusConfirm = "confirm"
)

// EduMgmt 链码
type EduMgmt struct {
}

// Paper 论文
type Paper struct {
	Id          string           `json:"id"`           // 论文编号 主键
	Title       string           `json:"title"`        // 标题
	Abstract    string           `json:"abstract"`     // 摘要
	Score       int64            `json:"score"`        // 成绩
	Description string           `json:"description"`  // 描述(备注)
	Student     string           `json:"student"`      // 学生名称
	Teacher     string           `json:"teacher"`      // 导师名称
	Department  string           `json:"department"`   // 学院
	OralDefense map[string]int64 `json:"oralDefense"`  // 答辩信息
	Status      string           `json:"status"`       // 状态      ----请使用上面const中的常量来标记状态
	RejectCount uint64           `json:"reject_count"` // 被拒次数
	RefuseInfo  string           `json:"refuse_info"`  // 拒绝详情
	SubmitTime  time.Time        `json:"submit_time"`  // 论文提交时间
	RejectTime  time.Time        `json:"reject_time"`  // 论文退回时间
	ApproveDate time.Time        `json:"approve_date"` // 导师通过时间
	CreatedAt   time.Time        `json:"created_at"`   // 论文创建时间
	UpdatedAt   time.Time        `json:"updated_at"`   // 论文修改更新时间
}

// Person 人员信息
type Person struct {
	Id         string    `json:"id"`         // 编号（工号、学号）
	Account    string    `json:"account"`    // 登录账号 主键
	Password   string    `json:"password"`   // 登录密码
	Phone      string    `json:"phone"`      // 联系电话
	Email      string    `json:"email"`      // 电子邮箱
	Name       string    `json:"name"`       // 姓名
	Address    string    `json:"address"`    // 地址
	College    string    `json:"college"`    // 学院
	Department string    `json:"department"` // 系部
	Class      string    `json:"class"`      // 班级
	Type       string    `json:"type"`       // 人员类型，学生，老师，学校
	Status     string    `json:"status"`     // 人员状态 [new,confirm]
	CreatedAt  time.Time `json:"created_at"` // 创建时间
	UpdatedAt  time.Time `json:"updated_at"` // 更新时间
}

// Init 初始化
func (s *EduMgmt) Init(stub shim.ChaincodeStubInterface) peer.Response {

	// 格式化时间
	now, err := time.ParseInLocation(TimeFormat, time.Now().Format(TimeFormat), time.Local)
	if err != nil {
		return shim.Error(err.Error())
	}

	for i := 1; i <= 5; i++ {
		// 导师
		teacher := Person{
			Id:         "t100" + strconv.Itoa(i),
			Account:    "teacher0" + strconv.Itoa(i),
			Password:   "123456",
			Phone:      "1888888888" + strconv.Itoa(i),
			Email:      "teacher0" + strconv.Itoa(i) + "@gdzce.cn",
			Name:       strconv.Itoa(i) + "号老师",
			Address:    "广州市海珠区",
			College:    "外语学院",
			Department: "外语学院",
			Class:      "应用日语一班",
			Type:       Teacher,
			Status:     personStatusConfirm,
			CreatedAt:  now,
		}
		personBytes, err := json.Marshal(teacher)
		if err != nil {
			shim.Error(err.Error())
		}
		// 通过用户类型和用户账号存储用户
		compositeKey, err := stub.CreateCompositeKey(teacher.Type, []string{teacher.Account})
		if err != nil {
			return shim.Error("CreateCompositeKey error: " + err.Error())
		}
		err = stub.PutState(compositeKey, personBytes)
		if err != nil {
			return shim.Error("PutState error: " + err.Error())
		}
	}

	// 教务处人员
	school := Person{
		Id:         "s1001",
		Account:    "school01",
		Password:   "123456",
		Phone:      "18888888888",
		Email:      "school01@gdzce.cn",
		Name:       "学校A",
		Address:    "广州市海珠区",
		College:    "外语学院",
		Department: "",
		Class:      "",
		Type:       School,
		Status:     personStatusConfirm,
		CreatedAt:  now,
	}
	personBytes, err := json.Marshal(school)
	if err != nil {
		shim.Error(err.Error())
	}
	// 通过用户类型和用户账号存储用户
	compositeKey, err := stub.CreateCompositeKey(school.Type, []string{school.Account})
	if err != nil {
		return shim.Error("CreateCompositeKey error: " + err.Error())
	}
	err = stub.PutState(compositeKey, personBytes)
	if err != nil {
		return shim.Error("PutState error: " + err.Error())
	}

	//***************测试数据
	datas := []Paper{
		{
			Id:       "202105080001111",
			Title:    "testTitle1111",
			Abstract: "testAbstract1111",
			Student:  "liangsheng",
			Teacher:  "teacher01",
			Status:   PaperStatusNew,
		},
		{
			Id:       "202105080002222",
			Title:    "testTitle2222",
			Abstract: "testAbstract2222",
			Student:  "liangsheng",
			Teacher:  "teacher01",
			Status:   PaperStatusApprove,
		},
		{
			Id:       "202105080003333",
			Title:    "testTitle3333",
			Abstract: "testAbstract3333",
			Student:  "liangsheng",
			Teacher:  "teacher01",
			Status:   PaperStatusOralDefense,
			OralDefense: map[string]int64{
				"teacher01": -1,
				"teacher02": -1,
				"teacher03": -1,
				"teacher04": -1,
				"teacher05": -1,
			},
		},
		{
			Id:       "202105080004444",
			Title:    "testTitle4444",
			Abstract: "testAbstract4444",
			Student:  "liangsheng",
			Teacher:  "teacher02",
			Status:   PaperStatusOralDefense,
			OralDefense: map[string]int64{
				"teacher06": -1,
				"teacher02": -1,
				"teacher03": -1,
				"teacher04": -1,
				"teacher05": -1,
			},
		},
		{
			Id:       "202105080005555",
			Title:    "testTitle5555",
			Abstract: "testAbstract5555",
			Student:  "liangsheng",
			Teacher:  "teacher01",
			Status:   PaperStatusSubmit,
		},
		{
			Id:       "202105080006666",
			Title:    "testTitle6666",
			Abstract: "testAbstract6666",
			Student:  "liangsheng",
			Teacher:  "teacher01",
			Status:   PaperStatusApprove,
		},
		{
			Id:       "202105080007777",
			Title:    "testTitle7777",
			Abstract: "testAbstract7777",
			Student:  "liangsheng",
			Teacher:  "teacher01",
			Status:   PaperStatusSubmit,
		},
		{
			Id:       "202105080008888",
			Title:    "testTitle8888",
			Abstract: "testAbstract8888",
			Student:  "liangsheng",
			Teacher:  "teacher01",
			Status:   PaperStatusPassed,
		},
		{
			Id:       "202105080009999",
			Title:    "testTitle9999",
			Abstract: "testAbstract9999",
			Student:  "lilin",
			Teacher:  "teacher01",
			Status:   PaperStatusArchive,
		},
	}
	for _, data := range datas {
		personBytes, err := json.Marshal(data)
		if err != nil {
			shim.Error(err.Error())
		}
		// 通过用户类型和用户账号存储用户
		compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{data.Id})
		if err != nil {
			return shim.Error("CreateCompositeKey error: " + err.Error())
		}
		err = stub.PutState(compositeKey, personBytes)
		if err != nil {
			return shim.Error("PutState error: " + err.Error())
		}
	}
	//***************

	return shim.Success(nil)
}

// Invoke 业务处理
func (s *EduMgmt) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, args := stub.GetFunctionAndParameters()

	switch funcName {
	case "register": // 新建用户
		return register(stub, args)
	case "getPerson": // 获取单个用户
		return getPerson(stub, args)
	case "confirmUser": // 修改用户信息
		return confirmUser(stub, args)
	case "getPersons": // 获取用户列表
		return getPersons(stub, args)
	case "newPaper": // 新建论文
		return newPaper(stub, args)
	case "getPapers": // 获取论文列表
		return getPapers(stub, args)
	case "getPaper": // 获取单个论文
		return getPaper(stub, args)
	case "updatePaper": // 修改论文
		return updatePaper(stub, args)
	case "deletePaper": // 删除论文
		return deletePaper(stub, args)
	case "submitPaper": // 学生提交论文
		return submitPaper(stub, args)
	case "rejectPaper": // 导师退回论文
		return rejectPaper(stub, args)
	case "approvePaper": // 导师通过论文
		return approvePaper(stub, args)
	case "oralDefense": // 安排论文答辩
		return oralDefense(stub, args)
	case "markPaper": // 导师评分
		return markPaper(stub, args)
	case "archivePaper": // 论文归档
		return archivePaper(stub, args)
	case "get":
		return get(stub, args)
	case "getList":
		return getList(stub, args)
	case "set":
		return set(stub, args)
	case "delete":
		return deleteData(stub, args)
	default:
		return shim.Error(fmt.Sprintf("不支持的智能合约: %s", funcName))
	}
}

// register 人员注册
func register(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Person
	err := json.Unmarshal([]byte(args[0]), &requestData)
	if err != nil {
		return shim.Error("反序列化参数报错: " + err.Error())
	}

	// 必填参数的非空校验
	if requestData.CreatedAt.IsZero() || requestData.Id == "" || requestData.Account == "" || requestData.Password == "" || requestData.Name == "" || requestData.Type == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(requestData.Type, []string{requestData.Account})
	if err != nil {
		return shim.Error("创建组合键失败: " + err.Error())
	}

	// 验证账号是否存在
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) != 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "账号信息错误或账号已存在",
		}
	}

	// 验证用户id是否重复
	queryIterator, err := stub.GetStateByPartialCompositeKey(requestData.Type, []string{})
	if err != nil {
		return shim.Error("GetStateByPartialCompositeKey error: " + err.Error())
	}
	defer queryIterator.Close()
	// 遍历相应类型的所有帐号
	for queryIterator.HasNext() {
		next, err := queryIterator.Next()
		if err != nil {
			return shim.Error("queryIterator.Next() error: " + err.Error())
		}
		var chainData Person
		err = json.Unmarshal(next.GetValue(), &chainData)
		if err != nil {
			return shim.Error("Unmarshal error: " + err.Error())
		}
		// 验证
		if chainData.Id == requestData.Id {
			return peer.Response{
				Status:  shim.ERRORTHRESHOLD,
				Message: "学号/工号已存在, 请重新输入",
			}
		}
	}

	// 更改状态
	requestData.Status = personStatusNew

	txBytes, err := json.Marshal(requestData)
	if err != nil {
		return shim.Error("序列化person对象出错: " + err.Error())
	}
	err = stub.PutState(compositeKey, txBytes)
	if err != nil {
		return shim.Error("PutState error: " + err.Error())
	}
	return shim.Success(nil)
}

// markPaper 答辩评分
func markPaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 验证参数的正确性
	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error("反序列化参数报错: " + err.Error())
	}

	// 必填参数的非空校验
	if requestData.Id == "" || requestData.Student == "" || requestData.Teacher == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	if requestData.Score < 0 || requestData.Score > 100 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "评分应在0-100分之间",
		}
	}

	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 从账本中获取值
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "论文不存在, 请输入正确的论文ID",
		}
	}
	var chainData Paper
	err = json.Unmarshal(state, &chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	if chainData.Status != PaperStatusOralDefense {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "只能对答辩状态的论文评分",
		}
	}
	// 检查是否答辩列表里的老师
	if value, ok := chainData.OralDefense[requestData.Teacher]; ok {
		// 不等于初始值-1, 代表已经修改过评分
		if value != -1 {
			return peer.Response{
				Status:  shim.ERRORTHRESHOLD,
				Message: fmt.Sprintf("%s, 您已给予 %v 分, 不能重复评分", requestData.Teacher, value),
			}
		}
		chainData.OralDefense[requestData.Teacher] = requestData.Score
	} else {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: fmt.Sprintf("%s, 您无权对论文[%s]进行评分", requestData.Teacher, chainData.Title),
		}
	}
	// 统计当前论文已评分的老师数量
	var count, sum int64
	for _, value := range chainData.OralDefense {
		if value > -1 {
			count++
			sum += value
		}
	}
	// 当评分人数达到5人, 则自动进行论文评分
	if count == 5 {
		// 取评分平均值, 大于或等于90则答辩通过
		average := sum / count
		if average >= 90 {
			chainData.Status = PaperStatusPassed
		} else {
			chainData.Status = PaperStatusReject
		}
		chainData.Score = average
	}

	// 写回区块链账本
	state, err = json.Marshal(chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Marshal error: %s", err))
	}
	err = stub.PutState(compositeKey, state)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState error: %s", err))
	}
	return shim.Success(nil)
}

// 删除论文
func deletePaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error("反序列化参数报错: " + err.Error())
	}
	// 检查必要参数
	if requestData.Id == "" || requestData.Student == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 从账本中获取值
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "论文不存在, 请输入正确的论文ID",
		}
	}

	var chainData Paper
	err = json.Unmarshal(state, &chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}

	if chainData.Student != requestData.Student {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "不能删除其他人的论文",
		}
	}

	if chainData.Status != PaperStatusNew && chainData.Status != PaperStatusReject {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "只能删除新建或退回状态的论文",
		}
	}

	// 将论文状态修改为已删除再DelState
	chainData.Status = PaperStatusCancel

	txBytes, err := json.Marshal(chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Marshal error: %s", err))
	}
	// 写入区块链
	err = stub.PutState(compositeKey, txBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState error: %s", err))
	}

	return shim.Success(nil)
}

// 获取单个用户 登录
func getPerson(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Person
	err := json.Unmarshal([]byte(args[0]), &requestData)
	if err != nil {
		return shim.Error("反序列化参数报错: " + err.Error())
	}

	// 必填参数的非空校验
	if requestData.Type == "" || requestData.Account == "" || requestData.Password == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	key, err := stub.CreateCompositeKey(requestData.Type, []string{requestData.Account})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 通过主键从区块链查找相关的数据
	state, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("query persons error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "账户类型错误或账户不存在",
		}
	}
	// 将数据返回
	return shim.Success(state)
}

// 修改用户信息
// 只有学校用户才能更改用户信息
func confirmUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var requestData Person
	err := json.Unmarshal([]byte(args[0]), &requestData)
	if err != nil {
		return shim.Error("反序列化参数报错: " + err.Error())
	}

	// 必填参数的非空校验
	targetType := args[1]
	targetAccount := args[2]
	if requestData.UpdatedAt.IsZero() || requestData.Type == "" || requestData.Account == "" || targetType == "" || targetAccount == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	compositeKey, err := stub.CreateCompositeKey(requestData.Type, []string{requestData.Account})
	if err != nil {
		return shim.Error("CreateCompositeKey error: " + err.Error())
	}

	state, _ := stub.GetState(compositeKey)
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "查询账号失败",
		}
	}

	var chainData Person
	err = json.Unmarshal(state, &chainData)
	if err != nil {
		return shim.Error("Unmarshal error: " + err.Error())
	}

	if chainData.Type != School {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "当前用户没有此权限",
		}
	}

	compositeKey, err = stub.CreateCompositeKey(targetType, []string{targetAccount})
	if err != nil {
		return shim.Error("CreateCompositeKey error: " + err.Error())
	}

	state, err = stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "需要修改的用户不存在",
		}
	}

	var chainTargetData Person
	err = json.Unmarshal(state, &chainTargetData)
	if err != nil {
		return shim.Error("Unmarshal error: " + err.Error())
	}

	// 更改状态
	chainTargetData.Status = personStatusConfirm
	chainTargetData.UpdatedAt = requestData.UpdatedAt

	txBytes, err := json.Marshal(chainTargetData)
	if err != nil {
		return shim.Error("序列化person对象出错: " + err.Error())
	}
	err = stub.PutState(compositeKey, txBytes)
	if err != nil {
		return shim.Error("PutState error: " + err.Error())
	}
	return shim.Success(nil)
}

// getPersons 获取人员信息
func getPersons(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Person
	err := json.Unmarshal([]byte(args[0]), &requestData)
	if err != nil {
		return shim.Error("反序列化参数报错: " + err.Error())
	}
	// 检查必要参数
	if requestData.Type == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	var persons []Person
	// 通过主键从区块链查找相关的数据
	resultIterator, err := stub.GetStateByPartialCompositeKey(requestData.Type, []string{})
	if err != nil {
		return shim.Error(fmt.Sprintf("query persons error: %s", err))
	}
	defer resultIterator.Close()

	// 检查返回的数据是否为空，不为空则遍历数据，否则返回空数组
	for resultIterator.HasNext() {
		val, err := resultIterator.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("resultIterator error: %s", err))
		}

		var person Person
		if err := json.Unmarshal(val.GetValue(), &person); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}
		// 不返回密码
		person.Password = ""
		persons = append(persons, person)
	}

	// 序列化数据
	result, err := json.Marshal(persons)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal error: %s", err))
	}

	return shim.Success(result)
}

// newPaper 新增论文
func newPaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error("反序列化参数报错: " + err.Error())
	}
	// 检查必要参数
	if requestData.CreatedAt.IsZero() || requestData.Id == "" || requestData.Title == "" || requestData.Abstract == "" || requestData.Student == "" || requestData.Teacher == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 检测论文是否存在
	papers, err := stub.GetStateByPartialCompositeKey(PaperKeyPrefix, []string{})
	if err != nil {
		return shim.Error("获取论文失败")
	}
	verifyData := new(Paper)
	for papers.HasNext() {
		next, err := papers.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal(next.GetValue(), verifyData)
		if err != nil {
			return shim.Error(err.Error())
		}

		if verifyData.Title == requestData.Title {
			return peer.Response{
				Status:  shim.ERRORTHRESHOLD,
				Message: fmt.Sprintf("论文Title[%s]已存在", requestData.Title),
			}
		}
	}
	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}

	// 修改论文状态
	requestData.Status = PaperStatusNew
	txBytes, err := json.Marshal(requestData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Marshal error: %s", err))
	}
	// 写入区块链
	err = stub.PutState(compositeKey, txBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState error: %s", err))
	}

	return shim.Success(nil)
}

// updatePaper 修改论文
func updatePaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error("反序列化参数报错: " + err.Error())
	}
	// 检查必要参数
	if requestData.UpdatedAt.IsZero() || requestData.Id == "" || requestData.Title == "" || requestData.Student == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 从账本中获取值
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "需要修改的论文不存在",
		}
	}

	var chainData Paper
	err = json.Unmarshal(state, &chainData)
	if err != nil {
		return shim.Error("Unmarshal error: " + err.Error())
	}

	if chainData.Status != PaperStatusNew && chainData.Status != PaperStatusReject {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "只能修改新建或退回状态的论文",
		}
	}

	if chainData.Student != requestData.Student {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "不能修改其他人的论文",
		}
	}

	// 检测论文title是否重复
	papers, err := stub.GetStateByPartialCompositeKey(PaperKeyPrefix, []string{})
	if err != nil {
		return shim.Error("获取论文失败")
	}
	verifyData := new(Paper)
	for papers.HasNext() {
		next, err := papers.Next()
		if err != nil {
			return peer.Response{}
		}
		err = json.Unmarshal(next.GetValue(), verifyData)
		if err != nil {
			return peer.Response{}
		}

		if verifyData.Title == requestData.Title && verifyData.Id != requestData.Id {
			return peer.Response{
				Status:  shim.ERRORTHRESHOLD,
				Message: fmt.Sprintf("论文Title[%s]已存在", requestData.Title),
			}
		}
	}

	// 修改时间
	chainData.Title = requestData.Title
	chainData.Abstract = requestData.Abstract
	chainData.Description = requestData.Description
	chainData.UpdatedAt = requestData.UpdatedAt

	txBytes, err := json.Marshal(chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Marshal error: %s", err))
	}
	// 写入区块链
	err = stub.PutState(compositeKey, txBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState error: %s", err))
	}

	return shim.Success(nil)
}

// getPaper 获取单个论文
func getPaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error("反序列化参数报错: " + err.Error())
	}
	// 检查必要参数
	if requestData.Id == "" || requestData.Student == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 从账本中获取值
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "论文不存在, 请输入正确的论文ID",
		}
	}

	return shim.Success(state)
}

// getPapers 查询论文
func getPapers(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var papers []Paper
	// 通过主键从区块链查找相关的数据
	resultIterator, err := stub.GetStateByPartialCompositeKey(PaperKeyPrefix, []string{})
	if err != nil {
		return shim.Error(fmt.Sprintf("query papers error: %s", err))
	}
	defer resultIterator.Close()

	// 检查返回的数据是否为空，不为空则遍历数据，否则返回空数组
	for resultIterator.HasNext() {
		val, err := resultIterator.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("resultIterator error: %s", err))
		}
		var paper Paper
		if err := json.Unmarshal(val.GetValue(), &paper); err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error: %s", err))
		}
		papers = append(papers, paper)
	}
	// 序列化数据
	result, err := json.Marshal(papers)

	if err != nil {
		return shim.Error(fmt.Sprintf("marshal error: %s", err))
	}

	return shim.Success(result)
}

// submitPaper 提交论文
func submitPaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查必要参数
	if requestData.SubmitTime.IsZero() || requestData.Id == "" || requestData.Student == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 从账本中获取值
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "论文不存在, 请输入正确的论文ID",
		}
	}
	var chainData Paper
	err = json.Unmarshal(state, &chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查是否该学生的论文
	if chainData.Student != requestData.Student {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "不能提交其他人的论文",
		}
	}
	if chainData.Status != PaperStatusNew && chainData.Status != PaperStatusReject {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "只能提交新建状态或退回状态的论文",
		}
	}

	// 更改字段的值
	chainData.Status = PaperStatusSubmit
	chainData.SubmitTime = requestData.SubmitTime
	chainData.RefuseInfo = "" // 清空导师退回意见

	// 写回区块链账本
	state, err = json.Marshal(chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Marshal error: %s", err))
	}
	err = stub.PutState(compositeKey, state)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState error: %s", err))
	}

	return shim.Success(nil)
}

// submitPaper 退回论文
func rejectPaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查必要参数
	if requestData.RejectTime.IsZero() || requestData.Id == "" || requestData.Student == "" || requestData.Teacher == "" || requestData.RefuseInfo == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 从账本中获取值
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "论文不存在, 请输入正确的论文ID",
		}
	}
	var chainData Paper
	err = json.Unmarshal(state, &chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查是否该学生或该老师的论文
	if chainData.Teacher != requestData.Teacher || chainData.Student != requestData.Student {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "不能修改其他人的论文",
		}
	}
	if chainData.Status != PaperStatusSubmit {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "只能退回提交状态的论文",
		}
	}
	// 更改字段的值
	chainData.Status = PaperStatusReject          // 退回
	chainData.RejectCount += 1                    // 退回次数加 1
	chainData.RefuseInfo = requestData.RefuseInfo // 拒绝信息
	chainData.RejectTime = requestData.RejectTime

	// 写回区块链账本
	state, err = json.Marshal(chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Marshal error: %s", err))
	}
	err = stub.PutState(compositeKey, state)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState error: %s", err))
	}

	return shim.Success(nil)
}

// submitPaper 导师通过论文
func approvePaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查必要参数
	if requestData.ApproveDate.IsZero() || requestData.Id == "" || requestData.Teacher == "" || requestData.Student == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 从账本中获取值
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "需要修改的论文不存在",
		}
	}
	var chainData Paper
	err = json.Unmarshal(state, &chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查是否该学生或该老师的论文
	if chainData.Teacher != requestData.Teacher || chainData.Student != requestData.Student {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "不能修改其他人的论文",
		}
	}

	if chainData.Status != PaperStatusSubmit {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "只能通过提交状态的论文",
		}
	}

	// 更改字段的值
	chainData.Status = PaperStatusApprove // 导师通过
	chainData.ApproveDate = requestData.ApproveDate

	// 写回区块链账本
	state, err = json.Marshal(chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Marshal error: %s", err))
	}
	err = stub.PutState(compositeKey, state)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState error: %s", err))
	}

	return shim.Success(nil)
}

// oralDefense 安排论文答辩
func oralDefense(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查必要参数
	if requestData.Id == "" || requestData.Student == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 至少5个老师
	if len(requestData.OralDefense) <= 4 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "至少安排5位老师",
		}
	}

	// 从账本中取出该交易
	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 从账本中获取值
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "论文不存在, 请输入正确的论文ID",
		}
	}
	var chainData Paper
	err = json.Unmarshal(state, &chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查是否该学生的论文
	if chainData.Student != requestData.Student {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "不能修改其他人的论文",
		}
	}
	if chainData.Status != PaperStatusApprove {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "导师没有通过的论文无法安排答辩",
		}
	}

	// 更改字段的值
	chainData.Status = PaperStatusOralDefense       // 答辩
	chainData.OralDefense = requestData.OralDefense // 答辩导师信息

	// 写回区块链账本
	state, err = json.Marshal(chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Marshal error: %s", err))
	}
	err = stub.PutState(compositeKey, state)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState error: %s", err))
	}

	return shim.Success(nil)
}

// archivePaper 学校归档论文
func archivePaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var requestData Paper
	if err := json.Unmarshal([]byte(args[0]), &requestData); err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查必要参数
	if requestData.Id == "" || requestData.Student == "" {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "必填参数为空",
		}
	}

	// 从账本中取出该交易
	// 创建组合键
	compositeKey, err := stub.CreateCompositeKey(PaperKeyPrefix, []string{requestData.Id})
	if err != nil {
		return shim.Error(fmt.Sprintf("CreateCompositeKey error: %s", err))
	}
	// 从账本中获取值
	state, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetState error: %s", err))
	}
	if len(state) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "论文不存在, 请输入正确的论文ID",
		}
	}
	var chainData Paper
	err = json.Unmarshal(state, &chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unmarshal error: %s", err))
	}
	// 检查是否该学生的论文
	if chainData.Student != requestData.Student {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "不能修改其他人的论文",
		}
	}
	if chainData.Status != PaperStatusPassed {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: "答辩通过的论文才能归档",
		}
	}

	// 更改字段的值
	chainData.Status = PaperStatusArchive // 学校归档

	// 写回区块链账本
	state, err = json.Marshal(chainData)
	if err != nil {
		return shim.Error(fmt.Sprintf("Marshal error: %s", err))
	}
	err = stub.PutState(compositeKey, state)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState error: %s", err))
	}

	return shim.Success(nil)
}

// -------------------------- 以下评分测试专用，请勿修改 -----------------------------
// 评分测试专用，请勿修改
func set(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	compositeKey, _ := stub.CreateCompositeKey(args[0], []string{args[1]})
	err := stub.PutState(compositeKey, []byte(args[2]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to set key: %s with error: %s", compositeKey, err.Error()))
	}
	return shim.Success(nil)
}

// 评分测试专用，请勿修改
func get(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	compositeKey, _ := stub.CreateCompositeKey(args[0], []string{args[1]})
	value, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get key: %s with error: %s", compositeKey, err.Error()))
	}
	if len(value) == 0 {
		return peer.Response{
			Status:  shim.ERRORTHRESHOLD,
			Message: fmt.Sprintf("key not found: %s", compositeKey),
		}
	}
	return shim.Success(value)
}

// 评分测试专用，请勿修改
func getList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	iteratorData, err := stub.GetStateByPartialCompositeKey(args[0], []string{})
	if err != nil {
		return shim.Error(err.Error())
	}
	if args[0] == "paper" {
		var papers []Paper
		for iteratorData.HasNext() {
			next, err := iteratorData.Next()
			if err != nil {
				return shim.Error(err.Error())
			}
			var paper Paper
			err = json.Unmarshal(next.GetValue(), &paper)
			if err != nil {
				return shim.Error(err.Error())
			}
			papers = append(papers, paper)
		}
		marshal, err := json.Marshal(papers)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(marshal)
	} else {
		var persons []Person
		for iteratorData.HasNext() {
			next, err := iteratorData.Next()
			if err != nil {
				return shim.Error(err.Error())
			}
			var person Person
			err = json.Unmarshal(next.GetValue(), &person)
			if err != nil {
				return shim.Error(err.Error())
			}
			persons = append(persons, person)
		}
		marshal, err := json.Marshal(persons)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(marshal)
	}
}

// 评分测试专用，请勿修改
func deleteData(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	compositeKey, _ := stub.CreateCompositeKey(args[0], []string{args[1]})
	_ = stub.DelState(compositeKey)
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(EduMgmt))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
