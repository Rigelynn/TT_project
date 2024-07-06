package service

import (
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"tce/utils"
	"time"
)

// Action 正常交易->商家卖出烟草
func Action(db *gorm.DB) {
	fmt.Println("售卖信息如下:")
	infOnChain := GetAllInf(db)
	var flag = 0
	for flag == 0 {
		//记录内随机抽取卖出
		luckindex := utils.GetRandNumber(len(infOnChain))
		lucklist := infOnChain[luckindex]
		if lucklist.State == "1" {
			fmt.Println("抽取的是", lucklist, "记录,即该烟草卖出并记录在案")
			lucklist.State = "2"
			infOnChain[luckindex] = lucklist
			traceid := lucklist.TraceId
			db.Model(TobaccoTest{}).Where("TraceId = ?", traceid).Updates(&lucklist)
			testauc := Copy(lucklist)
			db.Model(TobaccoTestAuc{}).Where("TraceId = ?", traceid).Updates(testauc)
			fmt.Println("更新完成")
			flag = 1
		}
	}
}

// Action1 作恶1->商家未经备案卖出烟草
func Action1(db *gorm.DB) string {
	var log string
	fmt.Println("<商家未经备案卖出烟草>")
	fmt.Println("问题烟草信息:")
	infauc := GetAllInfauc(db)
	var flag = 0
	for flag == 0 {
		rand.Seed(time.Now().UnixNano())
		luckindex := utils.GetRandNumber(len(infauc))
		lucklist := infauc[luckindex]
		fmt.Println(lucklist)
		if lucklist.State == "1" {
			lucklist.State = "2"
			db.Model(TobaccoTestAuc{}).Where("TraceId = ?", lucklist.TraceId).Updates(lucklist)
			flag = 1
			log = "-->" + lucklist.TraceId + lucklist.MerchantID
		}
	}
	fmt.Println("模拟状态:成功")
	return log
}

// Action2 作恶2->卖家从外部获取烟草
func Action2(db *gorm.DB) string {
	fmt.Println("<卖家从外部获取烟草>")
	var log string
	NewTobacco := DataGenerator(1)
	NewTobaccoAuc := Copy(NewTobacco[0])
	if err := db.Create(&NewTobaccoAuc).Error; err != nil {
		fmt.Println(err)
		return "nil"
	}
	fmt.Println("外部获取烟草内容是:")
	fmt.Println(NewTobaccoAuc)
	log = "<--" + NewTobaccoAuc.TraceId + NewTobaccoAuc.MerchantID
	fmt.Println("模拟状态:成功")
	return log

}

// Action3 作恶3->交换部分烟草
func Action3(db *gorm.DB) string {
	MerchantIDList := []string{"M0001", "M0002", "M0003", "M0004", "M0005", "M0006", "M0007", "M0008", "M0009", "M0010", "M0011", "M0012", "M0013", "M0014", "M0015", "M0016", "M0017", "M0018", "M0019", "M0020"}
	var log string
	fmt.Println("<交换部分烟草>")
	//随机抽取两家
	Merchant1 := utils.GetRandNumber(20) + 1
	Merchant2 := (Merchant1 + utils.GetRandNumber(4)) % 20
	auc := GetAllInfauc(db)
	for _, v1 := range auc {
		if v1.MerchantID == MerchantIDList[Merchant2] {
			for _, v2 := range auc {
				if v2.MerchantID == MerchantIDList[Merchant1] {
					//交换序列
					v1.MerchantID, v2.MerchantID = v2.MerchantID, v1.MerchantID
					//重新插入
					db.Model(TobaccoTestAuc{}).Where("TraceId = ?", v1.TraceId).Updates(v1)
					db.Model(TobaccoTestAuc{}).Where("TraceId = ?", v2.TraceId).Updates(v2)
					fmt.Println("MerchantID交换成功,内容如下:")
					fmt.Println(v1, "<——>", v2)
					log = v1.MerchantID + "<-->" + v2.TobaccoId + "溯源码" + v1.TraceId
					break
				}
			}
			break
		}
	}
	return log
}
