package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common/constant"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common/model"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Auth struct {
	AccessKey      string
	SecretKey      string
	EndPoint       string
	SignedHeaders  string // 参与计算的头部
	HttpRequestMsg *model.HttpRequestMsg
}

type SignParams struct {
	Url    string
	Method string
	Body   []byte
}

func NewAuth(accessKey string, secretKey string, endPoint string) *Auth {
	return &Auth{
		AccessKey:      accessKey,
		SecretKey:      secretKey,
		EndPoint:       endPoint,
		HttpRequestMsg: newHttpRequestMsg(endPoint),
	}
}

func newHttpRequestMsg(endpoint string) *model.HttpRequestMsg {
	var requestMsg = model.HttpRequestMsg{Params: map[string]string{}, Headers: map[string]string{}}
	if len(endpoint) == 0 || "{endPoint}" == endpoint {
		requestMsg.Host = constant.HttpRequestDomain
	} else {
		requestMsg.Host = endpoint
	}
	requestMsg.SignedHeaders = getSignedHeaders("")
	timeStamp := getCurrentTimeSeconds()
	requestMsg.Headers[constant.HeadSignTimeStamp] = timeStamp
	requestMsg.Headers[constant.Host] = requestMsg.Host
	requestMsg.Headers[constant.ContentType] = constant.ApplicationJson
	requestMsg.Headers[constant.XCncAuthMethod] = constant.AKSK
	return &requestMsg
}

func (a *Auth) WithAuth(signParams *SignParams) {
	a.HttpRequestMsg.Headers[constant.HeadSignAccessKey] = a.AccessKey
	a.HttpRequestMsg.Uri = signParams.Url
	a.HttpRequestMsg.Url = fmt.Sprintf("https://%s%s", a.HttpRequestMsg.Headers[constant.Host], signParams.Url)
	a.HttpRequestMsg.Body = string(signParams.Body)
	signature := getSignature(a.HttpRequestMsg, a.SecretKey, a.HttpRequestMsg.Headers[constant.HeadSignTimeStamp])
	a.HttpRequestMsg.Headers[constant.Authorization] = genAuthorization(a.AccessKey, a.HttpRequestMsg.SignedHeaders, signature)
}

func getCurrentTimeSeconds() string {
	timeStamp := time.Now().UTC().Unix()
	return strconv.FormatInt(timeStamp, 10)
}

/*
*
拼接最后签名
*/
func genAuthorization(accessKey string, signedHeaders string, signature string) string {
	var build strings.Builder
	build.WriteString(constant.HeadSignAlgorithm)
	build.WriteString(" ")
	build.WriteString("Credential=")
	build.WriteString(accessKey)
	build.WriteString(", ")
	build.WriteString("SignedHeaders=")
	build.WriteString(signedHeaders)
	build.WriteString(", ")
	build.WriteString("Signature=")
	build.WriteString(signature)
	return build.String()
}

func getSignature(requestMsg *model.HttpRequestMsg, secretKey string, timeStamp string) string {
	var bodyStr = requestMsg.Body
	if len(requestMsg.Body) == 0 || "GET" == requestMsg.Method {
		bodyStr = ""
	}
	hashedRequestPayload := hmacSha256(bodyStr)
	canonicalRequest := requestMsg.Method + "\n" +
		requestMsg.Uri + "\n" +
		getQueryString(requestMsg) + "\n" +
		getCanonicalHeaders(requestMsg.Headers, requestMsg.SignedHeaders) + "\n" +
		getSignedHeaders(requestMsg.SignedHeaders) + "\n" +
		hashedRequestPayload
	stringToSign := constant.HeadSignAlgorithm + "\n" + timeStamp + "\n" + hmacSha256(canonicalRequest)
	return hmac256(secretKey, stringToSign)
}

/*
*
获取uri参数
*/
func getQueryString(requestMsg *model.HttpRequestMsg) string {
	indexOfQueryStringSeparator := strings.Index(requestMsg.Uri, "?")
	if "POST" == requestMsg.Method || indexOfQueryStringSeparator == -1 {
		return ""
	}
	s, err := url.QueryUnescape(requestMsg.Uri[indexOfQueryStringSeparator:len(requestMsg.Uri)])
	if err != nil {
		fmt.Println("decode请求参数失败.")
	}
	return s
}

/*
*
获取并排序参与签名计算的头部
*/
func getSignedHeaders(signedHeaders string) string {
	if len(signedHeaders) == 0 {
		return "content-type;host"
	}
	headers := strings.Split(strings.ToLower(signedHeaders), ";")
	sort.Strings(headers)
	return strings.Join(headers, ";")
}

/*
*
获取k-v字符串
*/
func getCanonicalHeaders(headerMap map[string]string, signedHeaders string) string {
	keys := strings.Split(signedHeaders, ";")
	var headers = make(map[string]string)
	for k, v := range headerMap {
		headers[strings.ToLower(k)] = v
	}
	var build strings.Builder
	for i := 0; i < len(keys); i++ {
		build.WriteString(keys[i])
		build.WriteString(":")
		build.WriteString(strings.ToLower(headers[keys[i]]))
		build.WriteString("\n")
	}
	return build.String()
}

/*
*
加密算法
*/
func hmacSha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashCode := hash.Sum(nil)
	result := hex.EncodeToString(hashCode)
	return strings.ToLower(result)
}

func hmac256(secretKey string, stringToSign string) string {
	value := []byte(secretKey)
	key := hmac.New(sha256.New, value)
	key.Write([]byte(stringToSign))
	result := hex.EncodeToString(key.Sum(nil))
	return strings.ToLower(result)
}
