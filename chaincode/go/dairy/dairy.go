ckage main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// 定义结构体, 继承ChainCode接口
type DairyFarm struct {
}
 
// 定义数据结构体
type FarmInfo struct {
	Id string
	// 原厂场名字
	Name string
	// 生产日期
	Date string
	// 质量等级
	Quality string
	// 产量, 单位t
	Yield int
}

// 方法实现
func (t *DairyFarm) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return t.init(stub)
}
					
					//DF---牛奶类
					//DB---酒水类
					//DS---饮料类
func (t *DairyFarm) init(stub shim.ChaincodeStubInterface) peer.Response {
	// 数据初始化
	infos := []FarmInfo{
		FarmInfo{Id: "DF-001", Name: "东郊农场", Date: "2021-12-11", Quality: "优", Yield: 5},
		FarmInfo{Id: "DF-002", Name: "东郊农场", Date: "2021-12-12", Quality: "劣", Yield: 5},
		FarmInfo{Id: "DF-003", Name: "东郊农场", Date: "2021-12-13", Quality: "良", Yield: 5},
		FarmInfo{Id: "DF-004", Name: "东郊农场", Date: "2021-12-14", Quality: "优", Yield: 5},
		FarmInfo{Id: "DF-005", Name: "东郊农场", Date: "2021-12-15", Quality: "良", Yield: 5},
		FarmInfo{Id: "DF-006", Name: "西郊农场", Date: "2021-12-11", Quality: "优", Yield: 6},
		FarmInfo{Id: "DF-007", Name: "西郊农场", Date: "2021-12-12", Quality: "良", Yield: 6},
		FarmInfo{Id: "DF-008", Name: "西郊农场", Date: "2021-12-13", Quality: "劣", Yield: 6},
		FarmInfo{Id: "DF-009", Name: "西郊农场", Date: "2021-12-14", Quality: "良", Yield: 6},
		FarmInfo{Id: "DF-010", Name: "西郊农场", Date: "2021-12-15", Quality: "优", Yield: 6},
		FarmInfo{Id: "DF-011", Name: "南郊农场", Date: "2021-12-11", Quality: "良", Yield: 8},
		FarmInfo{Id: "DF-012", Name: "南郊农场", Date: "2021-12-12", Quality: "优", Yield: 8},
		FarmInfo{Id: "DF-013", Name: "南郊农场", Date: "2021-12-13", Quality: "优", Yield: 8},
		FarmInfo{Id: "DF-014", Name: "南郊农场", Date: "2021-12-14", Quality: "良", Yield: 8},
		FarmInfo{Id: "DF-015", Name: "南郊农场", Date: "2021-12-15", Quality: "优", Yield: 8},
		FarmInfo{Id: "DF-016", Name: "北郊农场", Date: "2021-12-11", Quality: "劣", Yield: 3},
		FarmInfo{Id: "DF-017", Name: "北郊农场", Date: "2021-12-12", Quality: "良", Yield: 3},
		FarmInfo{Id: "DF-018", Name: "北郊农场", Date: "2021-12-13", Quality: "优", Yield: 3},
		FarmInfo{Id: "DF-019", Name: "北郊农场", Date: "2021-12-14", Quality: "良", Yield: 3},
		FarmInfo{Id: "DF-020", Name: "北郊农场", Date: "2021-12-15", Quality: "良", Yield: 3},

		FarmInfo{Id: "DB-001", Name: "谢勒斯酿酒原厂", Date: "2021-12-11", Quality: "优", Yield: 5},
		FarmInfo{Id: "DB-002", Name: "谢勒斯酿酒原厂", Date: "2021-12-12", Quality: "优", Yield: 5},
		FarmInfo{Id: "DB-003", Name: "谢勒斯酿酒原厂", Date: "2021-12-13", Quality: "良", Yield: 5},
		FarmInfo{Id: "DB-004", Name: "谢勒斯酿酒原厂", Date: "2021-12-14", Quality: "劣", Yield: 5},
		FarmInfo{Id: "DB-005", Name: "茅畬酿造原厂", Date: "2021-12-15", Quality: "良", Yield: 5},
		FarmInfo{Id: "DB-006", Name: "茅畬酿造原厂", Date: "2021-12-11", Quality: "优", Yield: 6},
		FarmInfo{Id: "DB-007", Name: "茅畬酿造原厂", Date: "2021-12-12", Quality: "良", Yield: 6},

		FarmInfo{Id: "DS-001", Name: "南昌可乐果农场", Date: "2021-12-11", Quality: "优", Yield: 5},
		FarmInfo{Id: "DS-002", Name: "东湖可乐果农场", Date: "2021-12-12", Quality: "优", Yield: 5},
		FarmInfo{Id: "DS-003", Name: "南郊可乐果农场", Date: "2021-12-13", Quality: "良", Yield: 5},
		FarmInfo{Id: "DS-004", Name: "东郊可乐果农场", Date: "2021-12-14", Quality: "优", Yield: 5},
		FarmInfo{Id: "DS-005", Name: "东郊古柯农场", Date: "2021-12-15", Quality: "良", Yield: 5},
		FarmInfo{Id: "DS-006", Name: "西郊古柯农场", Date: "2021-12-11", Quality: "优", Yield: 6},
		FarmInfo{Id: "DS-007", Name: "西郊古柯农场", Date: "2021-12-12", Quality: "劣", Yield: 6},
	}
	// 遍历, 写入账本中
	i := 0
	for i < len(infos) {
		jsontext, error := json.Marshal(infos[i])
		if error != nil {
			return shim.Error("错误,初始化失败!!!!")
		}
		// 数据写入账本中
		stub.PutState(infos[i].Id, jsontext)
		i++
	}
	return shim.Success([]byte("初始化成功!!!"))
}

func (t *DairyFarm) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, args := stub.GetFunctionAndParameters()

	if funcName == "setvalue" {
		return t.setvalue(stub, args)
	} else if funcName == "query" {
		return t.query(stub, args)
	} else if funcName == "gethistory" {
		return t.gethistory(stub, args)
	}
	return shim.Success([]byte("invoke OK"))
}
func (t *DairyFarm) setvalue(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyID := args[0]
	error := stub.PutState(keyID, []byte(args[1]))
	if error != nil {
		return shim.Error("设置失败,重新设置")
	}
	return shim.Success([]byte("设置成功!!!"))
}

func (t *DairyFarm) query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyID := args[0]
	value, error := stub.GetState(keyID)
	if error != nil {
		return shim.Error("获取数据失败,请重新操作!")
	}
	return shim.Success(value)//输出值
}

// 根据keyID查询历史记录
func (t *DairyFarm) gethistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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

func main() {
	error := shim.Start(new(DairyFarm))
	if error != nil {
		println("程序启动失败...")
		return
	}
	fmt.Println("程序启动成功...")
}

