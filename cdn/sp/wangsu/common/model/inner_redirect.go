package model

import "encoding/json"

type InnerRedirectRequest struct {
	DomainName          string                `path:"domain-name" json:"-"`
	RewriteRuleSettings []*RewriteRuleSetting `json:"rewrite-rule-settings"`
}

type RewriteRuleSetting struct {
	//url匹配模式，支持正则，客户入参参考：.*
	//对于匹配到的URL进行内容重定向
	PathPattern string `json:"path-pattern"`
	//匹配条件：指定常用类型，可选值为all或homepage 1. all：全部文件 2. homepage：首页
	CustomPattern string `json:"custom-pattern"`
	//目录
	Directory string `json:"directory"`
	//匹配条件：文件类型，多个请以英文;分隔，可选值：gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf m3u8 xml f4m bootstarp ts
	FileType string `json:"file-type"`
	//匹配条件：自定义文件类型，多个请以英文;分隔。
	CustomFileType string `json:"custom-file-type"`
	//例外的url匹配模式，某些URL除外：如abc.jpg，不做内容重定向
	//客户入参参考：^https?://[^/]+/.*\.m3u8
	ExceptPathPattern string `json:"except-path-pattern"`
	//忽略大小写，可选值为true或false，true表示忽略大小写；false表示不忽略大小写；
	//新增配置项时，不传默认为 true
	//如果客户传了空值：如，则表示清空配置
	IgnoreLetterCase bool `json:"ignore-letter-case"`
	//改写内容的生成位置。可输入值为：Cache表示节点；
	//暂不支持其他入参格式
	PublishType string `json:"publish-type"`
	//表示客户多组重定向内容的优先执行顺序。数字越大，优先级越高。
	//新增配置项时，不传默认为 10
	Priority int `json:"priority"`
	//配置项：旧url
	//表示改写前的协议方式（即需要改写的对象），如：^https://([^/]+/.*)
	//如果是回源协议改写，则表示客户请求的原始url，配套的参数after-value，表示客户请求需要转换的回源请求。
	BeforeValue string `json:"before-value"`
	//配置项：新url
	//表示改写后的协议方式，如：http://$1
	//如果请求重定向带状态码则参考入参：301:https://$1
	//注：如果url含域名，则域名需要是本身。
	AfterValue string `json:"after-value"`
	//重定向类型；支持入参：
	//before：防盗链之前
	//after：防盗链之后
	RewriteType string `json:"rewrite-type"`
	//匹配条件：请求头
	RequestHeader string `json:"request-header"`
	//匹配条件：例外的请求头
	ExceptionRequestHeader string `json:"exception-request-header"`
}

type InnerRedirectResponse struct {
	HttpStatusCode int    `json:"http-status-code"`
	XCncRequestId  string `json:"x-cnc-request-id"`
	Code           int    `json:"code"`
	Message        string `json:"message"`
}

func (o *InnerRedirectRequest) Marshal() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return data
}

func (o *InnerRedirectResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
