package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"gorm.io/gorm"
	"math/rand"
	"tce/sdkInit"
	"tce/utils"
	"time"
)

// Tobacco 烟草结构体
type Tobacco struct {
	ObjectType string `json:"ObjectType"`
	TraceId    string `json:"TraceId"`  //溯源码
	Type       string `json:"Type"`     // 烟草类型
	Source     string `json:"Source"`   //烟草来源
	AstaffId   string `json:"AStaffId"` //烟草物流中心的操作员工信息
	BstaffId   string `json:"BStaffId"` //售卖商户的操作员工信息
	Date       string `json:"Date"`     // 生产日期
	Quality    string `json:"Quality"`  // 质量等级

	Location string `json:"Location"`
	State    string `json:"State"` //状态  1->表示在烟草配送中心  2->表示在商户当中 3->表示已经售出
}

// QueryResult 获取所有员工结果
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Staff
}

// Staff 员工结构体
type Staff struct {
	ObjectType string   `json:"ObjectType"`
	Id         string   `json:"Id"`
	Name       string   `json:"Name"`
	Asset      []string `json:"Asset"` //处理的烟草溯源码集合
}

// StaffId Staff 员工id
type StaffId struct {
	Id string `json:"Id"`
}

// TobaccoId StaffId Staff 员工id
type TobaccoId struct {
	Id string `json:"Id"`
}

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}

func InitService(chaincodeID, channelID string, org *sdkInit.OrgInfo, sdk *fabsdk.FabricSDK) (*ServiceSetup, error) {
	handler := &ServiceSetup{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new channel client: %s", err)
	}
	handler.Client = client

	return handler, nil
}

type TobaccoTest struct {
	TraceId         string `gorm:"column:TraceId"`         //溯源码
	TobaccoId       string `gorm:"column:TobaccoId"`       //烟草ID
	MerchantID      string `gorm:"column:MerchantID"`      // 烟草类型
	TobaccoCategory string `gorm:"column:TobaccoCategory"` //烟草来源
	State           string `gorm:"column:State"`           //状态  1->表示在库,2->表示售出
}

type TobaccoTestAuc struct {
	TraceId         string `gorm:"column:TraceId"`         //溯源码
	TobaccoId       string `gorm:"column:TobaccoId"`       //烟草ID
	MerchantID      string `gorm:"column:MerchantID"`      // 烟草类型
	TobaccoCategory string `gorm:"column:TobaccoCategory"` //烟草来源
	State           string `gorm:"column:State"`           //状态  1->表示在库,2->表示售出
}

func (m TobaccoTest) TableName() string {
	return "t_test"
}

func (m TobaccoTestAuc) TableName() string {
	return "t_testauc"
}

// DataGenerator 随机生成数量的测试数据
func DataGenerator(number int) []TobaccoTest {
	TobaccolIDList := []string{"T0001", "T0002", "T0003", "t0004", "T0005", "T0006", "T0007", "T0008", "T0009", "T0010", "T0011", "T0012", "T0013", "T0014", "T0015", "T0016", "T0017", "T0018", "T0019", "T0020"}
	MerchantIDList := []string{"M0001", "M0002", "M0003", "M0004", "M0005", "M0006", "M0007", "M0008", "M0009", "M0010", "M0011", "M0012", "M0013", "M0014", "M0015", "M0016", "M0017", "M0018", "M0019", "M0020"}
	NameList := []string{"黄金叶", "红旗渠", "中华烟", "黄鹤楼", "利群香烟", "云烟", "双喜", "玉溪", "贵烟", "五叶神", "椰树", "红玫", "羊城", "帝豪", "红河", "红塔山", "红金龙", "好日子", "喜长沙", "芙蓉王"}

	var testlist []TobaccoTest
	for i := 0; i < number; i++ {
		var tobacco = TobaccoTest{}
		tobacco.TraceId = "TC-" + utils.RandomString(8)
		tobacco.TobaccoId = TobaccolIDList[rand.Intn(20)]
		tobacco.TobaccoCategory = NameList[rand.Intn(20)]
		tobacco.MerchantID = MerchantIDList[rand.Intn(20)]
		tobacco.State = "1"
		testlist = append(testlist, tobacco)
	}
	return testlist
}

// GetAllInf 获取所有在库信息
func GetAllInf(db *gorm.DB) []TobaccoTest {
	// 查询列表
	var TobaccoOnChain []TobaccoTest
	db.Debug().Find(&TobaccoOnChain)

	return TobaccoOnChain
}

// GetAllInf 获取所有在库信息
func GetAllInfauc(db *gorm.DB) []TobaccoTestAuc {
	// 查询列表
	var Tobaccoauc []TobaccoTestAuc
	db.Debug().Find(&Tobaccoauc)

	return Tobaccoauc
}
