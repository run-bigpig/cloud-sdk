package model

import "encoding/json"

type HeaderModifyRequest struct {
	DomainName        string              `path:"domain-name" json:"-"`
	HeaderModifyRules []*HeaderModifyRule `json:"header-modify-rules"`
}

type HeaderModifyRule struct {
	//url匹配模式，支持正则，如果是全部匹配，入参可以配置为：.*
	PathPattern string `json:"path-pattern"`
	//例外的url匹配模式，支持正则。 入参参考：
	ExceptPathPattern string `json:"except-path-pattern"`
	//匹配条件：指定常用类型，可选值为all或homepage 1. all：全部文件 2. homepage：首页
	CustomPattern string `json:"custom-pattern"`
	//匹配条件：文件类型，多个请以英文;分隔，可选值：gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf m3u8 xml f4m bootstarp ts
	FileType string `json:"file-type"`
	//匹配条件：自定义文件类型，多个请以英文分号分隔。
	CustomFileType string `json:"custom-file-type"`
	//目录
	Directory string `json:"directory"`
	//匹配条件：指定URL
	//入参不支持含http(s):// 开头的URI格式
	SpecifyUrl string `json:"specify-url"`
	//匹配的请求方式，可选值为：GET、POST、PUT、HEAD、DELETE、OPTIONS，多个请以英文分号分隔
	RequestMethod string `json:"request-method"`
	//http头的控制方向，可选值为cache2visitor/cache2origin/visitor2cache/origin2cache，单选。
	//cache2origin是指回源方向---对应配置项回源请求；
	//cache2visitor是指回客户端方向—对应配置项回客户端应答；
	//visitor2cache是指接收客户端请求
	//origin2cache是指接收源应答
	//配置接收源应答方向，添加非CACHE control头，无法传递给客户端
	HeaderDirection string `json:"header-direction"`
	//http头的控制类型，支持http头部的增删改，可选值为add|set|delete，单选。对应header-name、header-value参数
	//1. add：表示新增一个头部，头部名称为header-name，头部值为header-value
	//2. set：表示修改指定头部header-name的值为header-value
	//3. delete：表示删除头部，header-name可同时配置多个
	//注意：优先级delete>set>add。当源站有对应响应头，则按源站响应的头部响应给客户端，此处新增的无效。
	Action string `json:"action"`
	//http头正则匹配，可选值：true/false。
	//true：表示对header-name的值按正则匹配方式处理
	//false:表示对header-name的值按实际入参处理，不做正则匹配。
	//不传默认是false
	AllowRegexp *string `json:"allow-regexp"`
	//http头名称，新增或修改http头，只允许输入一个；删除http头允许输入多个，以分号“;”隔开。
	//1.当action为add：表示新增这个header-name头部
	//2.当action为set：修改这个header-name头部的值
	//3.当action为delete：删除这个header-name头部
	//注意：对特殊http头的操作是受限的，允许操作的http头及操作类型请参看【概览】-【附件2： header操作】
	HeaderName string `json:"header-name"`
	//http头域对应的值，例如：mytest.example.com
	//注意：
	//1. 当action为add或set时，该入参必须传值
	//2. 当action为delete时，该入参不用传
	//支持通过关键字获取指定变量值，如客户端ip，包含如下：
	//关键字：含义
	//#timestamp：当前时间，时间戳如1559124945
	//#request-host：请求头中的HOST
	//#request-url：请求url，包含协议域名等的全路径，如http://aaa.aa.com/a.html
	//#request-uri：请求uri，相对路径格式，如/index.html
	//#origin-ip：回源IP
	//#cache-ip：边缘节点IP
	//#server-ip：对外服务IP
	//#client-ip：客户端IP，即访客IP
	//#response-header{xxx}：获取响应头中的值，如#response-header{etag}，获取response-header中的etag值
	//#header{xxx}：获取请求的http header中的值，如#header{User-Agent}，是获取header中的User-Agent值
	//#cookie{xxx}：获取cookie中的值，如#cookie{account}，是获取cookie中设置的account的值
	HeaderValue string `json:"header-value"`
	//匹配请求头，头部值支持正则，头和头部值用空格隔开，如：Range bytes=[0-9]{9,}
	RequestHeader string `json:"request-header"`
	//表示客户多组配置的优先执行顺序。数字越大，优先级越高。 不传参默认为10，不可清空
	Priority *int `json:"priority"`
	//例外的文件类型
	ExceptFileType string `json:"except-file-type"`
	//例外的目录
	ExceptDirectory string `json:"except-directory"`
	//例外的请求方式
	ExceptRequestMethod string `json:"except-request-method"`
	//例外的请求头
	ExceptRequestHeader string `json:"except-request-header"`
}

type HeaderModifyResponse struct {
	//httpstatus=202; 表示成功调用新增域名接口，可使用header中的x-cnc-request-id查看当前新增域名的部署情况
	HttpStatusCode int `json:"http-status-code"`
	//用于访问该域名信息的URL,其中domain-id 为我司云平台为该域名生成的唯一表示,其值为字符串
	Location string `json:"location"`
	//唯一标示的id，用于查询每次请求的任务 （适用全部接口）
	XCncRequestId string `json:"x-cnc-request-id"`
	//错误代码，当HTTPStatus不为202时出现，表示当前请求调用的错误类型
	Code int `json:"code"`
	//响应信息，成功时为success
	Message string `json:"message"`
}

func (o *HeaderModifyRequest) Marshal() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return data
}

func (o *HeaderModifyResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
