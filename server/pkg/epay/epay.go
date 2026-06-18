package epay

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Client struct {
	ApiUrl     string
	Pid        string
	Key        string // V1 MD5 密钥
	PrivateKey string // V2 RSA 商户私钥
	PublicKey  string // V2 RSA 平台公钥
	Version    string // "v1" or "v2"
}

type CreateOrderReq struct {
	OutTradeNo string
	Type       string
	Name       string
	Money      string
	NotifyUrl  string
	ReturnUrl  string
	ClientIP   string
	Param      string
}

type CreateOrderResp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	TradeNo string `json:"trade_no"`
	PayUrl  string `json:"payurl"`
	QrCode  string `json:"qrcode"`
	PayType string `json:"pay_type"`
	PayInfo string `json:"pay_info"`
}

type NotifyParams struct {
	Pid         string
	TradeNo     string
	OutTradeNo  string
	Type        string
	Name        string
	Money       string
	TradeStatus string
	Param       string
	Sign        string
	SignType    string
	RawQuery    map[string]string
}

func (c *Client) buildSignParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || k == "sign_type" || params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	return strings.Join(parts, "&")
}

func (c *Client) signMD5(params map[string]string) string {
	str := c.buildSignParams(params) + c.Key
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func (c *Client) signRSA(params map[string]string) (string, error) {
	str := c.buildSignParams(params)
	block, _ := pem.Decode([]byte(c.PrivateKey))
	if block == nil {
		return "", errors.New("解析商户私钥失败")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		key2, err2 := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err2 != nil {
			return "", errors.New("解析商户私钥失败")
		}
		key = key2
	}
	rsaKey := key.(*rsa.PrivateKey)
	hashed := sha256.Sum256([]byte(str))
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (c *Client) verifyMD5(params map[string]string, sign string) bool {
	return c.signMD5(params) == sign
}

func (c *Client) verifyRSA(params map[string]string, sign string) bool {
	str := c.buildSignParams(params)
	block, _ := pem.Decode([]byte(c.PublicKey))
	if block == nil {
		return false
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}
	rsaPub := pub.(*rsa.PublicKey)
	hashed := sha256.Sum256([]byte(str))
	sigBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}
	return rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], sigBytes) == nil
}

func (c *Client) Sign(params map[string]string) (string, error) {
	if c.Version == "v2" {
		return c.signRSA(params)
	}
	return c.signMD5(params), nil
}

func (c *Client) Verify(params map[string]string, sign string) bool {
	if c.Version == "v2" {
		return c.verifyRSA(params, sign)
	}
	return c.verifyMD5(params, sign)
}

func (c *Client) GetPayUrl(req CreateOrderReq) (string, error) {
	params := map[string]string{
		"pid":          c.Pid,
		"type":         req.Type,
		"out_trade_no": req.OutTradeNo,
		"notify_url":   req.NotifyUrl,
		"return_url":   req.ReturnUrl,
		"name":         req.Name,
		"money":        req.Money,
	}
	if req.Param != "" {
		params["param"] = req.Param
	}
	var signType string
	var submitUrl string
	if c.Version == "v2" {
		params["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
		signType = "RSA"
		submitUrl = strings.TrimRight(c.ApiUrl, "/") + "/api/pay/submit"
	} else {
		signType = "MD5"
		submitUrl = strings.TrimRight(c.ApiUrl, "/") + "/submit.php"
	}
	sign, err := c.Sign(params)
	if err != nil {
		return "", err
	}
	params["sign"] = sign
	params["sign_type"] = signType
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return submitUrl + "?" + v.Encode(), nil
}

func (c *Client) CreateOrder(req CreateOrderReq) (*CreateOrderResp, error) {
	params := map[string]string{
		"pid":          c.Pid,
		"type":         req.Type,
		"out_trade_no": req.OutTradeNo,
		"notify_url":   req.NotifyUrl,
		"name":         req.Name,
		"money":        req.Money,
		"clientip":     req.ClientIP,
	}
	if req.ReturnUrl != "" {
		params["return_url"] = req.ReturnUrl
	}
	if req.Param != "" {
		params["param"] = req.Param
	}
	var signType string
	var apiUrl string
	if c.Version == "v2" {
		params["method"] = "web"
		params["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
		signType = "RSA"
		apiUrl = strings.TrimRight(c.ApiUrl, "/") + "/api/pay/create"
	} else {
		signType = "MD5"
		apiUrl = strings.TrimRight(c.ApiUrl, "/") + "/mapi.php"
	}
	sign, err := c.Sign(params)
	if err != nil {
		return nil, err
	}
	params["sign"] = sign
	params["sign_type"] = signType
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result CreateOrderResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("解析支付响应失败")
	}
	return &result, nil
}

func ParseNotify(r *http.Request) NotifyParams {
	q := r.URL.Query()
	p := NotifyParams{
		Pid:         q.Get("pid"),
		TradeNo:     q.Get("trade_no"),
		OutTradeNo:  q.Get("out_trade_no"),
		Type:        q.Get("type"),
		Name:        q.Get("name"),
		Money:       q.Get("money"),
		TradeStatus: q.Get("trade_status"),
		Param:       q.Get("param"),
		Sign:        q.Get("sign"),
		SignType:    q.Get("sign_type"),
		RawQuery:    make(map[string]string),
	}
	for k := range q {
		if q.Get(k) != "" {
			p.RawQuery[k] = q.Get(k)
		}
	}
	return p
}

func (n NotifyParams) ToMap() map[string]string {
	m := make(map[string]string)
	for k, v := range n.RawQuery {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if v != "" {
			m[k] = v
		}
	}
	return m
}
