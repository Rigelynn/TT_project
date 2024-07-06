package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"tce/Init"
	"tce/service"
	"tce/utils"
	"time"
)

// 设置常数    链码名称  版本

const (
	cc_name      = "TCE_cc"
	cc_version   = "1.0.0"
	username     = "root"       //账号
	password     = "root123..." //密码
	host         = "121.4.26.4" //数据库地址，可以是Ip或者域名
	port         = 3306         //数据库端口
	Dbname       = "tobacco"    //数据库名
	runrate      = 0.1          //违法比例
	testnumber   = 5000         //测试数据总量
	simulatetime = 50           //市场行为数量
)

func main() {

	//初始化区块链SDK
	//fmt.Println("----------------------阶段1:初始化----------------------")
	db := Init.InitGrom()
	//TestList := service.State1(testnumber)
	//fmt.Println("----------------------阶段1:完成----------------------")
	//fmt.Println("-------------------", testnumber, "条烟草数据上链-------------------")
	//for _, v := range TestList {
	//	var auc service.TobaccoTestAuc
	//	auc.State = v.State
	//	auc.TraceId = v.TraceId
	//	auc.TobaccoId = v.TobaccoId
	//	auc.TobaccoCategory = v.TobaccoCategory
	//	auc.MerchantID = v.MerchantID
	//
	//	if err := db.Create(&auc).Error; err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	if err := db.Create(&v).Error; err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//}

	fmt.Println("----------------------阶段2:模拟市场行为----------------------")
	fmt.Println("----------------------规定作恶比例为", runrate, "----------------------")
	number, _ := service.State2(db, simulatetime, runrate)
	fmt.Println("----------------------阶段2:模拟市场成功,", number, "次违法交易----------------------")
	fmt.Println("----------------------阶段3:开始验证----------------------")
	start := time.Now()
	service.State4(db)
	cost := time.Since(start)
	fmt.Println("---------------耗时：", cost, "-------")
}

// GetAllInfauc // GetAllInf 获取所有在库信息
//
//	func GetAllInf(db *gorm.DB) []service.TobaccoTest {
//		// 查询列表
//		var TobaccoOnChain []service.TobaccoTest
//		db.Debug().Find(&TobaccoOnChain)
//
//		return TobaccoOnChain
//	}
//
// // GetAllInf 获取所有在库信息
func GetAllInfauc(db *gorm.DB) []service.TobaccoTestAuc {
	// 查询列表
	var Tobaccoauc []service.TobaccoTestAuc
	db.Debug().Find(&Tobaccoauc)

	return Tobaccoauc
}

func DataGenerator(number int) []service.TobaccoTest {
	TobaccolIDList := []string{"T0001", "T0002", "T0003", "t0004", "T0005", "T0006", "T0007", "T0008", "T0009", "T0010", "T0011", "T0012", "T0013", "T0014", "T0015", "T0016", "T0017", "T0018", "T0019", "T0020"}
	MerchantIDList := []string{"M0001", "M0002", "M0003", "M0004", "M0005", "M0006", "M0007", "M0008", "M0009", "M0010", "M0011", "M0012", "M0013", "M0014", "M0015", "M0016", "M0017", "M0018", "M0019", "M0020"}
	NameList := []string{"黄金叶", "红旗渠", "中华烟", "黄鹤楼", "利群香烟", "云烟", "双喜", "玉溪", "贵烟", "五叶神", "椰树", "红玫", "羊城", "帝豪", "红河", "红塔山", "红金龙", "好日子", "喜长沙", "芙蓉王"}

	var testlist []service.TobaccoTest
	for i := 0; i < number; i++ {
		var tobacco = service.TobaccoTest{}
		tobacco.TraceId = "TC-" + utils.RandomString(6)
		tobacco.TobaccoId = TobaccolIDList[rand.Intn(20)]
		tobacco.TobaccoCategory = NameList[rand.Intn(20)]
		tobacco.MerchantID = MerchantIDList[rand.Intn(20)]
		tobacco.State = "1"
		testlist = append(testlist, tobacco)
	}
	return testlist
}

func CollectRoute(r *gin.Engine, setup *service.ServiceSetup) *gin.Engine {
	//路由
	r.POST("/UpdateTobacco", func(ctx *gin.Context) {

		var requestUser = service.Tobacco{}
		err := ctx.ShouldBind(&requestUser)
		if err != nil {
			ctx.String(http.StatusNotFound, "绑定form失败")
			return
		}

		fmt.Println("开始更新烟草信息主体......")
		fmt.Println(requestUser)

		//返回结果
		Txid, err := setup.UpdateTobacco(requestUser)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"txid": Txid,
		})

	})

	r.POST("/AddStaff", func(ctx *gin.Context) {

		var requestUser = service.Staff{}
		err := ctx.ShouldBind(&requestUser)
		if err != nil {
			ctx.String(http.StatusNotFound, "绑定form失败")
			return
		}

		fmt.Println("开始添加员工主体......")
		fmt.Println(requestUser)

		//返回结果
		Txid, err := setup.AddStaff(requestUser)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"txid": Txid,
		})

	})

	r.POST("/Query1", func(ctx *gin.Context) {

		var requestUser = service.TobaccoId{}
		err := ctx.ShouldBind(&requestUser)
		requestUser.Id = requestUser.Id + "_" + "1"
		if err != nil {
			return
		}
		var bodyBytes [][]byte
		bodyBytes = append(bodyBytes, []byte(requestUser.Id))
		resp, err := setup.Query1(bodyBytes)
		var data interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"失败":  "",
				"错误码": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"data": data,
			})
		}
	})

	r.POST("/Query2", func(ctx *gin.Context) {

		var requestUser = service.TobaccoId{}
		err := ctx.ShouldBind(&requestUser)
		requestUser.Id = requestUser.Id + "_" + "2"
		if err != nil {
			return
		}
		fmt.Println(requestUser.Id)
		var bodyBytes [][]byte
		bodyBytes = append(bodyBytes, []byte(requestUser.Id))
		resp, err := setup.Query2(bodyBytes)

		var data interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"失败":  "",
				"错误码": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"data": data,
			})
		}

	})

	r.POST("/QueryStaff", func(ctx *gin.Context) {

		var staffId = service.StaffId{}

		err := ctx.ShouldBind(&staffId)

		if err != nil {
			return
		}
		fmt.Println(staffId)
		var bodyBytes [][]byte
		bodyBytes = append(bodyBytes, []byte(staffId.Id))

		resp, _ := setup.QueryStaff(bodyBytes)
		var data interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"data": data,
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"data": data,
			})
		}

	})

	r.GET("/QueryAllStaff", func(ctx *gin.Context) {

		resp, _ := setup.QueryAllStaff()

		var data []map[string]interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{
				"code": err.Error(),
				"data": data,
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"data": data,
			})
		}

	})
	return r
}
