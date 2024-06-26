package model

import "encoding/json"

type VisitControlRequest struct {
	DomainName string `path:"domain-name" json:"-"`
	//防盗链配置
	//注意：
	//1. 需要取消防盗链配置设置时，可以传入空节点。
	//2. 表示需要设置防盗链配置时，此项必填
	VisitControlRules []*VisitControlRule `json:"visit-control-rules"`
}

type RefererControlRule struct {
	//是否允许空referer：合法refer、（合法域名、合法URL）、非法refer、（非法域名、非法URL）这六项任意一项有值，则“是否允许空referer”不能为空；合法refer、（合法域名、合法URL）、非法refer、（非法域名、非法URL）这四项都为空值，则“是否允许空referer”必须为空
	AllowNullReferer string `json:"allow-null-referer"`
	//合法referer.可以输入url或域名，支持正则，可以多个，多个以空格隔开
	ValidReferer string `json:"valid-referer"`
	//合法url，输入正确的url格式，不支持正则，可以多个，多个以分号分割。
	ValidUrl string `json:"valid-url"`
	//合法域名，不支持正则，可以多个，多个以分号分割
	ValidDomain string `json:"valid-domain"`
	//非法referer，可以输入url或域名，支持正则，可以多个，多个以空格隔开
	InvalidReferer string `json:"invalid-referer"`
	//非法url，输入正确的url格式，不支持正则，可以多个，多个以分号分割
	InvalidUrl string `json:"invalid-url"`
	//非法域名，不支持正则，可以多个，多个以分号分割
	InvalidDomain string `json:"invalid-domain"`
}

type IpControlRule struct {
	//禁止的IP段
	//支持输入IP或IP段，IP段之间用分号(;)隔开，如1.1.1.0/24;2.2.2.2
	//禁止的IP和例外的IP，只能一个有值
	ForbiddenIps string `json:"forbidden-ips"`
	//例外的IP段，支持输入IP或IP段，IP段之间用分号(;)隔开，如1.1.1.0/24;2.2.2.2，某些IP例外，不做防盗链
	AllowedIps string `json:"allowed-ips"`
}

type UaControlRule struct {
	//允许的客户端，正则匹配，不允许空格，配置多个UA如：Android|iPhone
	ValidUserAgents string `json:"valid-user-agents"`
	//禁止的客户端，正则匹配，不允许空格，配置多个UA如：Android|iPhone
	InvalidUserAgents string `json:"invalid-user-agents"`
}

type AdvanceControlRule struct {
	//禁止的访客区域，多个请用英文分号分隔。注意
	//1、仅支持iso 3166-1国家二字简称
	//2、如果有特殊区域配置需求，请联系您的专属。
	//3、同一组规则里，禁止的访客区域、允许的访客区域，不能同时配
	InvalidVisitorRegion string `json:"invalid-visitor-region"`
	//允许的访客区域，多个请用英文分号分隔。注意
	//1、仅支持iso 3166-1国家二字简称
	//2、如果有特殊区域配置需求，请联系您的专属。
	//3、同一组规则里，禁止的访客区域、允许的访客区域，不能同时配
	ValidVisitorRegion string `json:"valid-visitor-region"`
}

type CookieControlRule struct {
	//允许的cookie。填写正则格式，比如(.*)(range1|range2)(.*)。
	AllowCookie string `json:"allow-cookie"`
	//是否允许空cookie。只允许填写true或false。
	AllowNullCookie bool `json:"allow-null-cookie"`
	//禁止的cookie。填写正则格式，比如(.*)(range1|range2)(.*)。
	ForbiddenCookie string `json:"forbidden-cookie"`
}

type CustomHeaderControlRule struct {
	//来源。可选择来源于客户端还是服务端。客户端填写client，服务端填写server
	HeaderDirection string `json:"header-direction"`
	//头域白名单
	HeaderWhitelist string `json:"header-whitelist"`
	//头域值白名单
	HeaderValueWhitelist string `json:"header-value-whitelist"`
	//头域黑名单
	HeaderBlacklist string `json:"header-blacklist"`
	//头域值黑名单
	HeaderValueBlacklist string `json:"header-value-blacklist"`
}

type VisitControlRule struct {
	//url匹配模式，支持正则，如果是全部匹配，入参可以配置为：.*
	PathPattern string `json:"path-pattern"`
	//例外的url匹配模式，某些URL除外：如abc.jpg，不做防盗链功能
	//客户入参参考：^https?://[^/]+/.*\.m3u8
	ExceptPathPattern string `json:"except-path-pattern"`
	//指定常用类型：选择需要防盗链的域名是全部文件还是首页。入参参考值：
	//all：全部文件
	//homepage：首页
	CustomPattern string `json:"custom-pattern"`
	//文件类型：指定文件类型进行防盗链设置。
	//文件类型包括：gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf
	//如果需要全部类型，则直接传all。多个以分号隔开，all和具体文件类型不能同时配置。
	FileType string `json:"file-type"`
	//自定义文件类型：在指定文件类型外根据自身需求，填写适当的可识别文件类型。可以搭配file-type使用。如果file-type也有配置，实际生效的文件类型是两个入参的总和
	CustomFileType string `json:"custom-file-type"`
	//指定URL缓存：根据需求指定url进行防盗链设置
	//入参不支持含http(s):// 开头的URI格式
	SpecifyUrlPattern string `json:"specify-url-pattern"`
	//目录：指定目录进行防盗链设置
	//输入合法的目录格式。多个以英文分号隔开
	Directory string `json:"directory"`
	//例外的文件类型：指定不需要进行防盗链功能的文件类型
	//文件类型包括：gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf
	//如果需要全部类型，则直接传all。多个以分号隔开，all和具体文件类型不能同时配置
	//如果file-type=all,except-file-type=all 则表示不匹配任务文件类型
	ExceptFileType string `json:"except-file-type"`
	//例外的自定义文件类型：在指定文件类型外根据自身需求，填写适当的可识别文件类型。可以搭配except-file-type使用。如果except-file-type也有配置，实际生效的文件类型是两个入参的总和
	ExceptCustomFileType string `json:"except-custom-file-type"`
	//例外的目录：指定不需要进行进行防盗链设置的目录
	//输入合法的目录格式。多个以英文分号隔开
	ExceptDirectory string `json:"except-directory"`
	//控制方向。可选值：403和302
	//1） 403表示返回特定的错误状态码来拒绝服务（默认方式，状态码可以指定，一般为403）。
	//2） 302表示返回302 Found的重定向url，重定向的url可以指定。如果传302，rewrite-to必填
	ControlAction string `json:"control-action"`
	//指定302跳转后的url。如果control-action值为302，此项必填，值需为具体调整的url，不支持正则。如果control-action值为403，此项填不需要输入，值无效
	RewriteTo string `json:"rewrite-to"`
	//表示客户多组重定向内容的优先执行顺序。数字越大，优先级越高。
	//新增配置项时，不传默认为 10
	Priority int `json:"priority"`
	//例外的请求方法。多个以;隔开
	ExceptionalRequest string `json:"exceptional-request"`
	//是否忽略大小写。只允许填写true或false
	IgnoredCase bool `json:"ignored-case"`
	//标识IP黑白名单防盗链
	//注意：
	//1. 表示一组黑白名单防盗链，一个data-id下只能一组
	//2. 当传空标签表示清楚例外的IP段配置和禁止的IP段配置。
	IpControlRule *IpControlRule `json:"ip-control-rule"`
	//标识referer防盗链
	//注意：
	//1. 表示一组referer防盗链，一个data-id下只能一组
	//2. 当传空标签表示清除referer防盗链
	//3. 合法refer、（合法域名、合法URL）、非法refer、（非法域名、非法URL）这四项，一个data-id下只能配置一个或者都为空
	//4. 匹配条件一致或者有存在交集的情况下（匹配条件包括URL匹配模式；文件类型；自定义文件类型；目录；指定常用类型；指定url），且控制动作均为禁止时，多条配置不能同时配置<合法refer>或者<合法域名>或者<合法URL>或者（<合法域名>和<合法URL>）
	RefererControlRule *RefererControlRule `json:"referer-control-rule"`
	//标识UA头防盗链，
	//注意：
	//1. 表示一组UA头防盗链，一个data-id下只能一组
	//2. 当传空标签表示清除UA头防盗链
	UaControlRule *UaControlRule `json:"ua-control-rule"`
	//配置其他访问控制策略，比如禁止的访客区域，JSON示例：
	//advance-control-rules:{invalid-visitor-region:CN;JP;KR}
	AdvanceControlRules *AdvanceControlRule `json:"advance-control-rules"`
	//配置Cookie防盗链策略。【允许的cookie】和【禁止的cookie】只允许配置一个
	CookieControlRules *CookieControlRule `json:"cookie-control-rules"`
	//配置自定义头部防盗链。【头域黑名单】和【头域白名单】只允许配置一个
	CustomHeaderControlRules *CustomHeaderControlRule `json:"custom-header-control-rules"`
}

type VisitControlResponse struct {
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

func (o *VisitControlRequest) Marshal() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return data
}

func (o *VisitControlResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
