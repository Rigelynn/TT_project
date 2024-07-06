package service

import (
	"encoding/hex"
	"fmt"
	"github.com/cbergoon/merkletree"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"tce/utils"
	"time"
)

func State1(testnumber int) []TobaccoTest {
	fmt.Println("----------------------开始初始化区块链网络----------------------")
	//serviceSetup := Init.Initfabric()
	fmt.Println("----------------------开始初始化数据库----------------------")
	fmt.Println("-------------------初始化测试数据-------------------")
	TestList := DataGenerator(testnumber)
	fmt.Println("-------------------共初始化", testnumber, "烟草数据-------------------")
	return TestList
}

// State2 模拟市场行为
func State2(db *gorm.DB, simulatetime int, runrate float64) (int, []string) {
	conut := 0
	var log []string
	for i := 1; i <= simulatetime; i++ {
		//fmt.Println("第", i, "次模拟")
		rand.Seed(time.Now().UnixNano())
		randNumber := utils.GetRandNumber(100)
		fmt.Println("模拟因子:", randNumber)
		if float64(randNumber) <= runrate*100 { //作恶事件发生
			conut++
			rand.Seed(time.Now().UnixNano())
			auctionnumber := utils.GetRandNumber(1000)%2 + 1
			fmt.Println("---------注意!模拟作恶行为,作恶序列为---------", auctionnumber)
			if auctionnumber == 1 {
				log = append(log, Action1(db))
			} else if auctionnumber == 2 {
				log = append(log, Action2(db))
			} else {
				log = append(log, Action3(db))
			}

		} else { //正常行为
			fmt.Println("<正常售卖行为>")
			Action(db)
		}
	}
	fmt.Println("共模拟", simulatetime, "次市场行为", conut, "次违法行为发生")
	return conut, log
}

// State3 测试所有用户
func State3(db *gorm.DB) {
	MerchantIDList := []string{"M0001", "M0002", "M0003", "M0004", "M0005", "M0006", "M0007", "M0008", "M0009", "M0010", "M0011", "M0012", "M0013", "M0014", "M0015", "M0016", "M0017", "M0018", "M0019", "M0020"}
	AnsList := GetAllInf(db)
	//fmt.Println(AnsList)
	HashAnsMap := make(map[string]*merkletree.MerkleTree) //key v
	//构建所有人的答案默克尔树
	for k, _ := range MerchantIDList {
		var list []merkletree.Content
		for _, Merchant := range AnsList {
			if Merchant.MerchantID == MerchantIDList[k] {
				list = append(list, Merchant)
			}
		}
		t, err := merkletree.NewTree(list)
		if err != nil {
			log.Fatal(err)
		}
		mr := t.MerkleRoot()
		fmt.Println(MerchantIDList[k], "的根哈希值:", hex.EncodeToString(mr))
		HashAnsMap[MerchantIDList[k]] = t
	}

	fmt.Println("下面是线下检测结果:")
	AucList := GetAllInfauc(db)
	fmt.Println(AucList)
	HashAucMap := make(map[string]*merkletree.MerkleTree)
	//构建所有人的答案默克尔树
	for k, _ := range MerchantIDList {
		var listAuc []merkletree.Content
		for _, Merchant := range AucList {

			if Merchant.MerchantID == MerchantIDList[k] {
				listAuc = append(listAuc, Merchant)
			}
		}
		t, err := merkletree.NewTree(listAuc)
		if err != nil {
			log.Fatal(err)
		}
		mr := t.MerkleRoot()
		fmt.Println(MerchantIDList[k], "的根哈希值:", hex.EncodeToString(mr))
		HashAucMap[MerchantIDList[k]] = t
	}

}

func State4(db *gorm.DB) {
	MerchantIDList := []string{"M0001", "M0002", "M0003", "M0004", "M0005", "M0006", "M0007", "M0008", "M0009", "M0010", "M0011", "M0012", "M0013", "M0014", "M0015", "M0016", "M0017", "M0018", "M0019", "M0020"}
	AnsList := GetAllInf(db)
	AucList := GetAllInfauc(db)

	var inf_ans map[string][]TobaccoTest
	var inf_auc map[string][]TobaccoTestAuc
	inf_ans = make(map[string][]TobaccoTest)
	inf_auc = make(map[string][]TobaccoTestAuc)
	for _, v := range AnsList {
		inf_ans[v.MerchantID] = append(inf_ans[v.MerchantID], v)
	}
	for _, v := range AucList {
		inf_auc[v.MerchantID] = append(inf_auc[v.MerchantID], v)
	}
	fmt.Println(inf_auc["M0015"][0])
	fmt.Println(inf_ans["M0015"][0])

	fmt.Println(len(inf_auc["M0015"]))
	fmt.Println(inf_ans["M0015"][0].State == inf_auc["M0015"][0].State)
	for _, v := range MerchantIDList {
		if len(inf_ans[v]) != len(inf_auc[v]) {
			fmt.Println("--------商户", v, "长度不匹配--------")
			fmt.Println("商户:", v, "拥有:", len(inf_auc[v]), "商户:", v, "拥有:", len(inf_auc[v]))
			continue
		} else {
			for i := 0; i < len(inf_auc[v]); i++ {
				fmt.Println(inf_auc[v][i], inf_ans[v][i])
				if inf_ans[v][i].State != inf_auc[v][i].State || inf_ans[v][i].MerchantID != inf_auc[v][i].MerchantID {
					fmt.Println("--------商户", v, "数值不匹配--------")
					break
				}
				if i == len(inf_auc[v])-1 {
					fmt.Println("--------商户", v, "数值匹配--------")
				}
			}
		}
	}
}

func GetMyAnsTbc(AnsList []TobaccoTest) {

}
