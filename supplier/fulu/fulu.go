package fulu

import (
	"bytes"
	"context"
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
	return App{OpenApiURL: apiUrl, Key: key, Secret: secret}
}

// LoopCodeList 当发起充值遇到如下错误代码，自动切换至下一个商品
var LoopCodeList = []int{3001, 3002, 3003, 3004, 3005, 3008}

func inLoopCodeList(c int) bool {
	for _, code := range LoopCodeList {
		if code == c {
			return true
		}
	}
	return false
}

// DirectChannel 为OrderDirectAddChannel提供便利
type DirectChannel struct {
	wait     chan struct{}
	loopWait chan struct{}
	loopStop bool
	data     *OrderDirectAddResponse
	err      error
}

var ChannelCancel = errors.New("channel cancel")

// GetWait 等待数据反馈
func (o *DirectChannel) GetWait() <-chan struct{} {
	return o.wait
}

// GetData 获取数据
func (o *DirectChannel) GetData() *OrderDirectAddResponse {
	return o.data
}

// GetError 获取错误
func (o *DirectChannel) GetError() error {
	return o.err
}

// loopNext 进入下一次循环
func (o *DirectChannel) loopNext() {
	go func() {
		o.loopWait <- struct{}{}
	}()
}

func (o *DirectChannel) sendData(data *OrderDirectAddResponse, err error, isNext bool) {
	o.data = data
	o.err = err
	if isNext {
		o.loopNext()
	}
	o.wait <- struct{}{}
}

func (o *DirectChannel) close() {
	close(o.wait)
	close(o.loopWait)
}

// OrderDirectAddChannel 直充下单接口（多商品轮询channel版）
func (a App) OrderDirectAddChannel(ctx context.Context, ProductID []int, arg *OrderDirectAddParam) *DirectChannel {
	channel := &DirectChannel{wait: make(chan struct{}), loopWait: make(chan struct{}), data: nil, err: nil}
	channel.loopNext()
	go func() {
		defer channel.close()
		var ret OrderDirectAddResponse
		var err error
		i := 0
		for {
			select {
			case <-channel.loopWait:
				if channel.loopStop {
					channel.sendData(nil, ChannelCancel, false)
					return
				}
				id := ProductID[i]
				ret, err = a.OrderDirectAdd(arg.Clone(id, EmptyOrderNo))
				i++
				ret.ProductID = id
				if inLoopCodeList(ret.ParentInfo.Code) {
					if i == len(ProductID) {
						channel.sendData(&ret, err, false)
						return
					}
					channel.sendData(&ret, err, true) // 如果错误代码位于列表，继续执行后续id
					break
				}
				channel.sendData(&ret, err, false)
				return
			case <-ctx.Done():
				channel.loopStop = true
			}
		}
	}()
	return channel
}

// OrderDirectAddLoop 直充下单接口（多商品轮询普通版）
func (a App) OrderDirectAddLoop(ProductID []int, arg *OrderDirectAddParam) (OrderDirectAddResponse, error) {
	var ret OrderDirectAddResponse
	var err error
	for _, id := range ProductID {
		ret, err = a.OrderDirectAdd(arg.Clone(id, EmptyOrderNo))
		if inLoopCodeList(ret.ParentInfo.Code) {
			continue
		}
		if err != nil {
			return ret, err
		}
	}
	return ret, err
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
