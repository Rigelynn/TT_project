package service

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cbergoon/merkletree"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) UpdateTobacco(tobacco Tobacco) (string, error) {

	eventID := "updateTobacco"
	//注册事件
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)

	//事件defer
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将staff对象序列化成为字节数组
	b, err := json.Marshal(tobacco)

	if err != nil {
		return "", fmt.Errorf("指定的tobacco对象序列化时发生错误")
	}

	//req 是执行调用链码需要的参数

	req := channel.Request{
		ChaincodeID: t.ChaincodeID,                //通道名字
		Fcn:         "updateTobacco",              //函数名
		Args:        [][]byte{b, []byte(eventID)}, //函数参数
	}

	//t.Client.Execute(req) 是在后端执行该函数
	res, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	//事件结果
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(res.TransactionID), nil
}

func (t *ServiceSetup) Query1(entityID [][]byte) (channel.Response, error) {
	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "query1",
		Args:        entityID,
	}
	resp, _ := t.Client.Query(req)

	return resp, nil
}

func (t *ServiceSetup) Query2(entityID [][]byte) (channel.Response, error) {
	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "query2",
		Args:        entityID,
	}
	resp, _ := t.Client.Query(req)

	return resp, nil

}

func (t *ServiceSetup) AddStaff(staff Staff) (string, error) {

	eventID := "addStaff"
	//注册事件
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)

	//事件defer
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将staff对象序列化成为字节数组
	b, err := json.Marshal(staff)

	if err != nil {
		return "", fmt.Errorf("指定的staff对象序列化时发生错误")
	}

	//req 是执行调用链码需要的参数

	req := channel.Request{
		ChaincodeID: t.ChaincodeID,                //通道名字
		Fcn:         "addStaff",                   //函数名
		Args:        [][]byte{b, []byte(eventID)}, //函数参数
	}

	//t.Client.Execute(req) 是在后端执行该函数
	res, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	//事件结果
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(res.TransactionID), nil
}

func (t *ServiceSetup) QueryStaff(entityID [][]byte) (channel.Response, error) {

	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "queryStaff",
		Args:        entityID,
	}

	resp, _ := t.Client.Query(req)

	return resp, nil
}

func (t *ServiceSetup) QueryAllStaff() (channel.Response, error) {

	eventID := "QueryAllStaff"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "QueryAllStaff",
		Args:        [][]byte{[]byte(eventID)},
	}
	res, _ := t.Client.Query(req)

	_ = eventResult(notifier, eventID)

	return res, nil
}

func Copy(test TobaccoTest) TobaccoTestAuc {
	var testauc TobaccoTestAuc
	testauc.MerchantID = test.MerchantID
	testauc.State = test.State
	testauc.TobaccoCategory = test.TobaccoCategory
	testauc.TraceId = test.TraceId
	testauc.TobaccoId = test.TobaccoId
	return testauc
}

func (test TobaccoTest) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"TraceId":         test.TraceId,
		"TobaccoId":       test.TobaccoId,
		"TobaccoCategory": test.TobaccoCategory,
		"State":           test.State,
	})
}

// CalculateHash hashes the values of a TestContent
func (t TobaccoTestAuc) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.State + t.TraceId + t.TobaccoCategory + t.MerchantID + t.TobaccoId)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t TobaccoTestAuc) Equals(other merkletree.Content) (bool, error) {
	otherTC, ok := other.(TobaccoTestAuc)
	if !ok {
		return false, errors.New("value is not of type TestContent")
	}
	return t.State == otherTC.State && t.TraceId == otherTC.TraceId && t.TobaccoCategory == otherTC.TobaccoCategory && t.MerchantID == otherTC.MerchantID, nil
}

// CalculateHash hashes the values of a TestContent
func (t TobaccoTest) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.State + t.TraceId + t.TobaccoCategory + t.MerchantID + t.TobaccoId)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t TobaccoTest) Equals(other merkletree.Content) (bool, error) {
	otherTC, ok := other.(TobaccoTest)
	if !ok {
		return false, errors.New("value is not of type TestContent")
	}
	return t.State == otherTC.State && t.TraceId == otherTC.TraceId && t.TobaccoCategory == otherTC.TobaccoCategory && t.MerchantID == otherTC.MerchantID, nil
}

func GetIdentity(test TobaccoTest) []byte {
	ans, _ := json.Marshal(test)
	return ans
}
func GetIdentityAuc(test TobaccoTestAuc) []byte {
	ans, _ := json.Marshal(test)
	return ans
}
