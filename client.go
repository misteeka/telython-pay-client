package github

import (
	"fmt"
	transport "github.com/misteeka/fasthttp"
	"github.com/valyala/fastjson"
	"strconv"
)

type Response interface{}
type Status int

var (
	SUCCESS                Status = 100
	INVALID_REQUEST        Status = 101
	INTERNAL_SERVER_ERROR  Status = 102
	AUTHORIZATION_FAILED   Status = 103
	INVALID_CURRENCY_CODE  Status = 104
	CURRENCY_CODE_MISMATCH Status = 105
	NOT_FOUND              Status = 106
	WRONG_AMOUNT           Status = 107
	INSUFFICIENT_FUNDS     Status = 108
	TOO_MANY_REQUESTS      Status = 109
)

func StatusToString(status Status) string {
	if status == SUCCESS {
		return "SUCCESS"
	}
	if status == INVALID_REQUEST {
		return "INVALID REQUEST"
	}
	if status == INTERNAL_SERVER_ERROR {
		return "INTERNAL SERVER ERROR"
	}
	if status == AUTHORIZATION_FAILED {
		return "AUTHORIZATION FAILED"
	}
	if status == INVALID_CURRENCY_CODE {
		return "INVALID CURRENCY CODE"
	}
	if status == CURRENCY_CODE_MISMATCH {
		return "CURRENCY CODE MISMATCH"
	}
	if status == WRONG_AMOUNT {
		return "WRONG AMOUNT"
	}
	if status == NOT_FOUND {
		return "NOT FOUND"
	}
	if status == INSUFFICIENT_FUNDS {
		return "INSUFFICIENT FUNDS"
	}
	if status == TOO_MANY_REQUESTS {
		return "TOO MANY REQUESTS"
	}
	return fmt.Sprintf("%v", status)
}

const ip = "127.0.0.1" // 192.168.1.237

func get(function string) ([]byte, error) {
	resp, err := transport.Get(fmt.Sprintf("http://%s:8002/payments/%s", ip, function))
	if err != nil {
		return nil, err
	}
	response := resp.Body()
	transport.ReleaseResponse(resp)
	return response, nil
}
func post(function string, data string) ([]byte, error) {
	resp, err := transport.Post(fmt.Sprintf("http://%s:8002/payments/%s", ip, function), []byte(data))
	if err != nil {
		return nil, err
	}
	response := resp.Body()
	transport.ReleaseResponse(resp)
	return response, nil
}
func put(function string, data string) ([]byte, error) {
	resp, err := transport.Put(fmt.Sprintf("http://%s:8002/payments/%s", ip, function), []byte(data))
	if err != nil {
		return nil, err
	}
	response := resp.Body()
	transport.ReleaseResponse(resp)
	return response, nil
}

func getStatus(value *fastjson.Value) Status {
	return Status(value.GetInt("status"))
}

func SendPayment(sender uint64, receiver uint64, amount uint64, password string) (Status, error) {
	body, err := post("sendPayment", fmt.Sprintf(`{"sender":%d,"receiver":%d,"amount":%d, "password":"%s"}`, sender, receiver, amount, password))
	if err != nil {
		return 0, err
	}
	var p fastjson.Parser
	json, err := p.ParseBytes(body)
	if err != nil {
		return 0, err
	}
	return getStatus(json), nil
}

func GetBalance(accountId uint64, password string) (uint64, Status, error) {
	body, err := get("getBalance?a=" + strconv.FormatUint(accountId, 10) + "&p=" + password)
	if err != nil {
		return 0, 0, err
	}
	var p fastjson.Parser
	json, err := p.ParseBytes(body)
	if err != nil {
		return 0, 0, err
	}
	return json.GetUint64("data"), getStatus(json), nil
}

func CreateAccount(username string, password string, currency int) (uint64, Status, error) {
	body, err := post("createAccount", fmt.Sprintf(`{"username":"%s","password":"%s","currency":%d}`, username, password, currency))
	if err != nil {
		return 0, 0, err
	}
	var p fastjson.Parser
	json, err := p.ParseBytes(body)
	if err != nil {
		return 0, 0, err
	}
	return json.GetUint64("data"), getStatus(json), nil
}
