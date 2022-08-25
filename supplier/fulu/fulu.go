package fulu

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const Name = "fulu"

const OpenApiURL = "http://openapi.fulu.com/api/getway"
const TestOpenApiURL = "http://pre.openapi.fulu.com/api/getway"

const TimeFormat1 = "2006-01-02 15:04:05"
const TimeFormat2 = "20060102150405"

type App struct {
	OpenApiURL string
	Key        string
	Secret     string
}

// NewApp 创建接口实例
func NewApp(apiUrl, key, secret string) App {
	return App{
		OpenApiURL: apiUrl,
		Key:        key,
		Secret:     secret,
	}
}

// OrderCardAdd 卡密下单接口
func (a App) OrderCardAdd(arg *OrderCardAddParam) (OrderCardAddResponse, error) {

	var result OrderCardAddResponse

	if arg.CustomerOrderNo == "" {
		arg.CustomerOrderNo = CreateOrderID()
	}

	parent, data, err := a.httpPost("fulu.order.card.add", arg)
	if err != nil {
		return result, err
	}
	result.ParentInfo = parent

	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil

}

// OrderDirectAdd 直充下单接口
func (a App) OrderDirectAdd(arg *OrderDirectAddParam) (OrderDirectAddResponse, error) {

	var result OrderDirectAddResponse

	if arg.CustomerOrderNo == "" {
		arg.CustomerOrderNo = CreateOrderID()
	}

	parent, data, err := a.httpPost("fulu.order.direct.add", arg)
	if err != nil {
		return result, err
	}
	result.ParentInfo = parent

	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}

// OrderInfoGet 订单查询接口
func (a App) OrderInfoGet(arg *OrderInfoGetParam) (OrderInfoGetResponse, error) {

	var result OrderInfoGetResponse

	_, data, err := a.httpPost("fulu.order.info.get", arg)
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}

// GoodsListGet 商品列表
func (a App) GoodsListGet(arg *GoodsListGetParam) ([]GoodsListGetResponse, error) {

	_, data, err := a.httpPost("fulu.goods.list.get", arg)
	if err != nil {
		return nil, err
	}

	var body []GoodsListGetResponse
	if err = json.Unmarshal(data, &body); err != nil {
		return nil, err
	}

	return body, nil
}

// GoodsStockCheck 商品库存校验接口
func (a App) GoodsStockCheck(arg *GoodsStockCheckParam) (GoodsStockCheckResponse, error) {

	var body GoodsStockCheckResponse

	_, data, err := a.httpPost("fulu.goods.stock.check", arg)
	if err != nil {
		return body, err
	}

	if err = json.Unmarshal(data, &body); err != nil {
		return body, err
	}

	return body, nil
}

// GoodsInfoGet 获取商品信息接口
func (a App) GoodsInfoGet(arg *GoodsInfoGetParam) (GoodsInfoGetResponse, error) {

	var body GoodsInfoGetResponse

	_, data, err := a.httpPost("fulu.goods.info.get", arg)
	if err != nil {
		return body, err
	}

	if err = json.Unmarshal(data, &body); err != nil {
		return body, err
	}

	return body, nil
}

// OtherAPI 其它接口，使用该方法
func (a App) OtherAPI(api string, arg interface{}) (*ParentBody, []byte, error) {
	return a.httpPost(api, arg)
}

var ErrVerifyFail = errors.New("err verify fail")

func (a App) httpPost(api string, arg interface{}) (*ParentBody, []byte, error) {

	param := Param{}
	param.Set("app_key", a.Key)
	param.Set("method", api)
	param.Set("timestamp", time.Now().Format(TimeFormat1))
	param.Set("version", "2.0")
	param.Set("format", "json")
	param.Set("charset", "utf-8")
	param.Set("sign_type", "md5")
	param.Set("app_auth_token", "")
	if arg == nil {
		param.Set("biz_content", "{}")
	} else {
		data, _ := json.Marshal(arg)
		param.Set("biz_content", string(data))
	}

	req, err := http.NewRequest(http.MethodPost, a.OpenApiURL, a.sign(param))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", "RuiJie Spider (https://www.3721hy.com/)")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, nil, err
	}

	var ret ParentBody
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, nil, err
	}

	if !a.verify(&ret) {
		return nil, nil, ErrVerifyFail
	}

	return &ret, []byte(ret.Result), nil
}

func (a App) verify(param *ParentBody) bool {
	if param.Sign == "" {
		return true
	}
	dataStr := param.Result
	dataLen := len(dataStr)
	dataArr := make([]string, dataLen, dataLen)
	for i, str := range dataStr {
		dataArr[i] = string(str)
	}
	sort.Strings(dataArr)
	dataStr = strings.Join(dataArr, "") + a.Secret
	h := md5.New()
	h.Write([]byte(dataStr))
	dataStr = hex.EncodeToString(h.Sum(nil))
	return strings.ToLower(dataStr) == param.Sign
}

func (a App) sign(param Param) io.Reader {

	if param.Has("sign") {
		param.Del("sign")
	}

	data, _ := json.Marshal(param)
	dataStr := string(data)
	dataLen := len(dataStr)
	dataArr := make([]string, dataLen, dataLen)
	for i, str := range string(data) {
		dataArr[i] = string(str)
	}

	sort.Strings(dataArr)
	dataStr = strings.Join(dataArr, "") + a.Secret
	h := md5.New()
	h.Write([]byte(dataStr))
	dataStr = hex.EncodeToString(h.Sum(nil))

	param.Set("sign", strings.ToLower(dataStr))
	data, _ = json.Marshal(param)

	return bytes.NewReader(data)
}
