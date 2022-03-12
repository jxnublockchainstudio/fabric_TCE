package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type DairyFarm struct {
}
type FarmInfo struct {
	Id      string //奶牛ID
	Name    string // 奶牛场名字
	Date    string // 生产日期
	Quality string // 质量等级
	Yield   int    // 牛奶产量, 单位t
	Owner   string //拥有者
	Price   int    //价格
	State   int    //交易状态 1--滞留在库 2--已对下家交易结算
}
type MachinInfo struct {
	Id     string // 加工厂ID 默认为 1
	FromId string // 奶源奶牛ID

	Date     string // 生产日期
	Name     string // 加工厂名
	Validity int    // 保质期

	State int
}

type DistribuInfo struct {
	Id           string //配送厂ID
	FromMachinId string //加工厂来源ID
	FromId       string //奶源奶牛Id

	Date string //配送日期
	Name string //配送厂名

	State int
}
type Fund struct {
	Owner  string //所属组织
	CardId string //银行卡ID
	Amount int    //银行卡账户余额
}

/*
	交易结构体，对于牛奶厂和加工厂，产品是从牛奶从到加工厂的，所以交易的过程是
	1、记录要交易牛奶的ID（由调用者填入（也就是加工厂）），由该ID获得该产品的价格
	2、调用者（加工厂）输入自己的CardID 和 牛奶厂的CardID ，加工厂的Amount -=Price	牛奶厂的Amount +=Price
	3、用随机函数生成一个定长的字符串当作唯一的交易ID，并且打印出来，并记录当前系统时间 这样一个交易的结构就完成了
	4、对于调用者，只需在调用的时候给出以下参数，产品ID，对方银行账号的ID 自己银行的ID
	5、银行ID查验--为了防止调用者调用他人银行ID作为自己银行ID从而不减少自身余额，我们在账号下定义密码，由创建时填写
	这样就要修改“Qurey”函数时打印出来的东西，不能打印密码!
*/

type Transaction struct {
	TransId    string //交易id
	FromCardId string //发送方银行卡ID
	ToCardId   string //接收方银行卡ID
	Time       string //交易时间
	ProductID  string //交易产品id
	Account    int
	//ParentOrderNo string //父订单号
}

var err error

func (t *DairyFarm) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return t.init(stub)
}

//初始化数据
func (t *DairyFarm) init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Start Initing ...")

	//农场
	farminfos := []FarmInfo{
		FarmInfo{Id: "DF-101", Name: "东郊农场", Date: "2018-12-11", Quality: "优", Yield: 5, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-102", Name: "东郊农场", Date: "2018-12-12", Quality: "优", Yield: 5, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-103", Name: "东郊农场", Date: "2018-12-13", Quality: "良", Yield: 5, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-104", Name: "东郊农场", Date: "2018-12-14", Quality: "优", Yield: 5, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-105", Name: "东郊农场", Date: "2018-12-15", Quality: "良", Yield: 5, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-106", Name: "西郊农场", Date: "2018-12-11", Quality: "优", Yield: 6, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-107", Name: "西郊农场", Date: "2018-12-12", Quality: "良", Yield: 6, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-108", Name: "西郊农场", Date: "2018-12-13", Quality: "优", Yield: 6, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-109", Name: "西郊农场", Date: "2018-12-14", Quality: "良", Yield: 6, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-110", Name: "西郊农场", Date: "2018-12-15", Quality: "优", Yield: 6, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-111", Name: "南郊农场", Date: "2018-12-11", Quality: "良", Yield: 8, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-112", Name: "南郊农场", Date: "2018-12-12", Quality: "优", Yield: 8, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-113", Name: "南郊农场", Date: "2018-12-13", Quality: "优", Yield: 8, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-114", Name: "南郊农场", Date: "2018-12-14", Quality: "良", Yield: 8, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-115", Name: "南郊农场", Date: "2018-12-15", Quality: "优", Yield: 8, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-116", Name: "北郊农场", Date: "2018-12-11", Quality: "良", Yield: 3, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-117", Name: "北郊农场", Date: "2018-12-12", Quality: "良", Yield: 3, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-118", Name: "北郊农场", Date: "2018-12-13", Quality: "优", Yield: 3, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-119", Name: "北郊农场", Date: "2018-12-14", Quality: "良", Yield: 3, Owner: "组织一", Price: 2, State: 1},
		FarmInfo{Id: "DF-120", Name: "北郊农场", Date: "2018-12-15", Quality: "良", Yield: 3, Owner: "组织一", Price: 2, State: 1},
	}
	for i := 0; i < len(farminfos); i++ {
		jsonText, err := json.Marshal(farminfos[i])
		if err != nil {
			return shim.Error("init error, json marshal fail...")
		}
		//数据写入账本
		stub.PutState(farminfos[i].Id, jsonText)
	}

	return shim.Success([]byte("init ledger OK!!!"))
}
func (t *DairyFarm) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funKey, args := stub.GetFunctionAndParameters()
	if funKey == "CreateUser" {
		return t.CreateUser(stub, args)
	} else if funKey == "SetValue" {
		return t.SetValue(stub, args)
	} else if funKey == "Query" {
		return t.Query(stub, args)
	} else if funKey == "GetHistory" {
		return t.GetHistory(stub, args)
	} else if funKey == "Transaction" {
		return t.Transaction(stub, args)
	} else if funKey == "TransactionDT" {
		return t.TransactionDT(stub, args)
	} else if funKey == "Maciningforpd" {
		return t.Maciningforpd(stub, args)
	} else if funKey == "Distribuforpd" {
		return t.Distribuforpd(stub, args)
	} else {
		return shim.Success([]byte("funKey is wrong!!!"))
	}
}
func (t *DairyFarm) CreateUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		fmt.Println("number of args is wrong we need 2 strings ,CardId and your password")
	}
	var temp Fund
	temp.CardId = args[0]
	temp.Owner = args[1]
	temp.Amount, _ = strconv.Atoi("100")
	jsonText, _ := json.Marshal(temp)
	stub.PutState(args[0], jsonText)
	return shim.Success([]byte("CreateUser ok! ..."))
}
func (t *DairyFarm) SetValue(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect arguments. Expecting a key and 7 values")
	}
	fmt.Printf("len(args)=%d", len(args))
	keyId := args[0]
	var temp FarmInfo
	temp.Id = args[0]
	temp.Name = args[1]
	temp.Date = args[2]
	temp.Quality = args[3]
	temp.Yield, _ = strconv.Atoi(args[4])
	temp.Owner = args[5]
	temp.Price, _ = strconv.Atoi(args[6])
	temp.State, _ = strconv.Atoi(args[7])
	jsonText, _ := json.Marshal(temp)
	err := stub.PutState(keyId, jsonText)
	if err != nil {
		return shim.Error("PuState fail...")
	}
	return shim.Success([]byte("SetValue sucess!!!"))
}

func (t *DairyFarm) Query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyId := args[0]
	value, err := stub.GetState(keyId)
	if err != nil {
		return shim.Error("GetState fail...")
	}
	return shim.Success(value)
}

//输入产品的id 来查询获得产品，返回产品
func (t *DairyFarm) Queryproduct(stub shim.ChaincodeStubInterface, args string) (*FarmInfo, peer.Response) {

	keyId := args
	productAsbytes, err := stub.GetState(keyId)
	if err != nil {
		return nil, shim.Error("GetState fail...")
	}
	if productAsbytes == nil {
		return nil, shim.Error("GetState fail...")
	}
	prdt := new(FarmInfo)
	_ = json.Unmarshal(productAsbytes, prdt)
	return prdt, shim.Error("Unmarshal fail...")
}

//输入用户CardId 返回用户结构体
func (t *DairyFarm) QueryUsr(stub shim.ChaincodeStubInterface, args string) (peer.Response, *Fund) {

	keyId := args
	productAsbytes, err := stub.GetState(keyId)
	if err != nil {
		return shim.Error("GetState fail..."), nil
	}
	if productAsbytes == nil {
		return shim.Error("GetState fail..."), nil
	}
	user := new(Fund)
	_ = json.Unmarshal(productAsbytes, user)
	return shim.Error("GetState fail..."), user
}
func (t *DairyFarm) GetHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyiter, error := stub.GetHistoryForKey(args[0])
	if error != nil {
		return shim.Error("GetHistoryForKey fail...")
	}
	defer keyiter.Close()
	// 通过迭代器对象遍历结果
	var myList []string
	for keyiter.HasNext() {
		// 获取当前值
		result, error := keyiter.Next()
		if error != nil {
			return shim.Error("keyiter.Next() fail...")
		}
		// 获取需要的信息
		txID := result.TxId
		txValue := result.Value
		txTime := result.Timestamp
		txStatus := result.IsDelete
		tm := time.Unix(txTime.Seconds, 0)
		datastr := tm.Format("2006-01-02 15:04:05")
		all := fmt.Sprintf("%s, %s, %s, %t", txID, txValue, datastr, txStatus)
		myList = append(myList, all)
	}
	// 数据格式化为json
	jsonText, error := json.Marshal(myList)
	if error != nil {
		return shim.Error("json.Marshal(myList) fail...")
	}
	return shim.Success(jsonText)
}
func Exist(stub shim.ChaincodeStubInterface, a string, num int) int {
	existid := a
	existInfo, _ := stub.GetState(existid)
	farm := new(FarmInfo)
	json.Unmarshal(existInfo, farm)
	switch {
	case num == 1:
		if farm.Owner != "组织二" {
			return 1
		}
	case num == 2:
		if farm.Owner != "组织三" {
			return 2
		}
	}
	return 0
}

//组合KEY

func (t *DairyFarm) Maciningforpd(stub shim.ChaincodeStubInterface, args []string) peer.Response { //对于交易完成的订单先记录KEY值

	if Exist(stub, args[0], 1) == 1 {
		return shim.Error("Org2 not exist id ...")
	}
	var temp MachinInfo
	temp.Id = args[0] + "_1"
	//str := args[0] + "_1"
	temp.FromId = args[0]
	temp.Name = args[1]
	temp.Date = time.Now().Format("2006-01-02 15:04:05")
	temp.Validity, _ = strconv.Atoi(args[2])
	temp.State = 1
	jsontest, _ := json.Marshal(temp)
	stub.PutState(temp.Id, jsontest)
	return shim.Success([]byte("Machining ok! ..."))
}

func (t *DairyFarm) Transaction(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("Transaction starting ...")

	// 系统记录
	//var TransId string //交易id

	//用户输入
	var FromCardId string //发送方银行卡号
	var ToCardId string   //接收方银行卡号
	var ProductID string  //交易产品id

	//var product FarmInfo
	//var fund Fund
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting right number of information!!!")
	}
	ProductID = args[0]
	FromCardId = args[1]
	ToCardId = args[2]
	//查询出产品
	produc, _ := t.Queryproduct(stub, ProductID)
	//查询出用户
	if err != nil {
		fmt.Println("Query Error")
	}
	_, usr := t.QueryUsr(stub, FromCardId)
	_, usr1 := t.QueryUsr(stub, ToCardId)

	if usr.Amount < produc.Price {
		fmt.Println("Insufficient transaction failed...")
		return shim.Error("Balance is less than the price of traded goods...")
	}

	//用户认证成功 余额足

	produc.Owner = usr1.Owner
	produc.State = 2
	usr.Amount += produc.Price
	usr1.Amount -= produc.Price

	jsontest, _ := json.Marshal(usr)
	jsontest1, _ := json.Marshal(usr1)
	jsontest2, _ := json.Marshal(produc)

	stub.PutState(usr.CardId, jsontest)

	stub.PutState(usr1.CardId, jsontest1)

	stub.PutState(produc.Id, jsontest2)

	//系统交易完成

	//订单编写
	//订单完成后 ，对于结构体Machinifo 要对其赋值

	var tx Transaction // 1->2的交易信息

	tx.TransId = randSeq(6)
	tx.FromCardId = args[1]
	tx.ToCardId = args[2]
	tx.Time = time.Now().Format("2006-01-02 15:04:05")
	tx.ProductID = ProductID
	tx.Account = produc.Price

	jsonTx, errsss := json.Marshal(tx) //序列化

	if errsss != nil {
		return shim.Error("json.Marshal(tx) fail...")
	}
	stub.PutState(tx.TransId, jsonTx) //1->2的交易信息 存入数据库

	t.SetKeyForMa(stub, args[0])

	return shim.Success(jsonTx)
}

func (t *DairyFarm) SetKeyForMa(stub shim.ChaincodeStubInterface, a string) (string, peer.Response) { //对于交易完成的订单先记录KEY值

	var temp MachinInfo
	temp.Id = a + "_1"
	temp.FromId = a
	temp.State = 1
	temp.Date = ""
	temp.Name = ""
	temp.Validity = 0
	temp.State = 2
	jsontest, _ := json.Marshal(temp)
	stub.PutState(temp.Id, jsontest)
	return temp.Id, shim.Success([]byte("set key for machining ok! ...")) //对ID 组合ID（）  组合ID= 产品ID +"组织名"
}
func (t *DairyFarm) Distribuforpd(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if Exist(stub, args[0], 2) == 2 {
		return shim.Error("Org3 not exist id ...")
	}

	var temp DistribuInfo
	Mid, _ := t.SetKeyForMa(stub, args[0])
	temp.Id = args[0] + "_2"
	//str := args[0] + "_2"
	temp.FromMachinId = Mid
	temp.FromId = args[0]
	temp.Name = args[1]
	temp.Date = time.Now().Format("2006-01-02 15:04:05")
	temp.State = 1
	jsontest, _ := json.Marshal(temp)
	stub.PutState(temp.Id, jsontest)
	return shim.Success([]byte("Distributing ok! ..."))
}
func (t *DairyFarm) TransactionDT(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("TransactionDT starting ...")
	var FromCardId string //发送方银行卡号
	var ToCardId string   //接收方银行卡号
	var ProductID string  //交易产品id

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting right number of information!!!")
	}
	ProductID = args[0]
	FromCardId = args[1]
	ToCardId = args[2]
	produc, _ := t.Queryproduct(stub, ProductID)
	_, usr := t.QueryUsr(stub, FromCardId)
	_, usr1 := t.QueryUsr(stub, ToCardId)
	if usr.Amount < produc.Price {
		fmt.Println("Insufficient transaction failed...")
		return shim.Error("Balance is less than the price of traded goods...")
	}
	produc.Owner = usr.Owner
	produc.State = 2
	usr.Amount += produc.Price
	usr1.Amount -= produc.Price

	jsontest, _ := json.Marshal(usr)
	jsontest1, _ := json.Marshal(usr1)
	jsontest2, _ := json.Marshal(produc)

	stub.PutState(usr.CardId, jsontest)

	stub.PutState(usr1.CardId, jsontest1)

	stub.PutState(produc.Id, jsontest2)
	var txdt Transaction
	txdt.TransId = randSeq(6)
	txdt.FromCardId = args[1]
	txdt.ToCardId = args[2]
	txdt.Time = time.Now().Format("2006-01-02 15:04:05")
	txdt.ProductID = ProductID
	txdt.Account = produc.Price
	jsonTxDT, _ := json.Marshal(txdt)
	stub.PutState(txdt.TransId, jsonTxDT)

	//setkeyfordt
	Mid, _ := t.SetKeyForMa(stub, args[0])
	t.SetKeyForDt(stub, args[0], Mid)

	return shim.Success(jsonTxDT)
}
func (t *DairyFarm) SetKeyForDt(stub shim.ChaincodeStubInterface, a string, b string) peer.Response {
	var temp DistribuInfo
	temp.Id = a + "_2"
	temp.FromMachinId = b
	temp.FromId = a
	temp.State = 1
	temp.Date = ""
	temp.Name = ""
	jsontest, _ := json.Marshal(temp)
	stub.PutState(temp.Id, jsontest)
	return shim.Success([]byte("set key for distributing ok! ..."))
}

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func main() {
	err := shim.Start(new(DairyFarm))
	if err != nil {
		fmt.Println("程序启动失败...")
		return
	}
	fmt.Println("程序启动成功...")
}

/* setvalue 指令

peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n zz --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt  -c '{"Args":["SetValue","DF-121","南昌","2020-1-30","优","6","组织一","20","1"]}'

trancation 指令

*/
