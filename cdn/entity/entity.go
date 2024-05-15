package entity

type DomainItem struct {
	Id            int64  `json:"id"`             // 主键ID
	CustomerId    int64  `json:"customer_id"`    // 用户ID
	ResourceId    string `json:"resourceId"`     // 资源ID
	BusinessGroup int64  `json:"business_group"` // 业务组id
	AreaCode      int64  `json:"area_code"`      // 加速区域代码  0中国大陆 1境外 2全球
	CdnType       int64  `json:"cdn_type"`       // cdn类型 0cdn 1dcdn 2scdn
	ChannelType   int64  `json:"channel_type"`   // 业务类型 0 网页 1 下载 2 点播 3 全站加速
	Cname         string `json:"cname"`          // CNAME
	CreateTime    string `json:"create_time"`    // 创建时间
	Domain        string `json:"domain"`         // 域名
	DomainId      string `json:"domain_id"`      // 域名ID
	SupportIpv6   int64  `json:"support_ipv6"`   // 是否支持IPv6 0/1
	Status        int64  `json:"status"`         // 状态 0未验证 1已验证 2部署中 3部署成功 4暂停中 5已暂停
	Tag           string `json:"tag"`            // 标签
	TestUrl       string `json:"test_url"`       // 测试URL
	UpdateTime    string `json:"update_time"`    // 修改时间
}

type UpdateCdnDomainBaseConf struct {
	AreaCode    int64 `json:"area_code"`    // 加速区域代码  0中国大陆 1境外 2全球
	SupportIpv6 int64 `json:"support_ipv6"` // 是否支持IPv6 0/1
}

type OriginConf struct {
	OriginProtocol  int64  `json:"origin_protocol"`   //回源协议 0:http 1:https,2:follow
	OriginSniSwitch int64  `json:"origin_sni_switch"` //回源Sni开关
	OriginSniValue  string `json:"origin_sni_value"`  //回源Sni内容 只能是域名
	OriginRange     int64  `json:"origin_range"`      //Range回源
	OriginFollow    int64  `json:"origin_follow"`     //回源跟随301/302
	OriginEtag      int64  `json:"origin_etag"`       //回源是否验证Etag
	OriginTimeOut   int64  `json:"origin_time_out"`   //回源超时时间 秒单位
	TcpTimeout      int64  `json:"tcp_timeout"`       //tcp超时时间
}

type OriginServerConf struct {
	OriginType          int64  `json:"origin_type"`
	OriginAddressList   string `json:"origin_address_list"`
	OriginHttpPort      int64  `json:"origin_http_port"`
	OriginHttpsPort     int64  `json:"origin_https_port"`
	OriginHost          string `json:"origin_host"`
	OriginWeight        int64  `json:"origin_weight"`
	OriginPriority      int64  `json:"origin_priority"`
	OriginPriorityValue int64  `json:"origin_priority_value"`
}

type OriginAdvanceServerConf struct {
	UrlMatchMode        int64    `json:"url_match_mode"`
	UrlMatchRule        []string `json:"url_match_rule"`
	OriginType          int64    `json:"origin_type"`
	OriginAddressList   string   `json:"origin_address_list"`
	OriginHttpPort      int64    `json:"origin_http_port"`
	OriginHttpsPort     int64    `json:"origin_https_port"`
	OriginHost          string   `json:"origin_host"`
	OriginPriorityValue int64    `json:"origin_priority_value"`
}

type OriginSni struct {
	OriginSniSwitch int64  `json:"origin_sni_switch"`
	OriginSniValue  string `json:"origin_sni_value"`
}

type OriginRequestHeaderConf struct {
	Action         int64  `json:"action"`
	ParameterKey   string `json:"parameter_key"`
	ParameterValue string `json:"parameter_value"`
}

type CacheListItem struct {
	CacheType        int64    `json:"cache_type"`
	CacheContent     []string `json:"cache_content"`
	CacheTTL         int64    `json:"cache_ttl"`
	CacheUnit        int64    `json:"cache_unit"`
	CacheStatus      int64    `json:"cache_status"`
	UseRegex         int64    `json:"use_regex"`
	Priority         int64    `json:"priority"`
	ParametersStatus int64    `json:"parameters_status"`
	ParametersValue  []string `json:"parameters_value"`
	Capitalization   int64    `json:"capitalization"`
	Remark           string   `json:"remark"`
}

type CacheCodeListItem struct {
	HttpCode  int64 `json:"http_code"`
	CacheTTL  int64 `json:"cache_ttl"`
	CacheUnit int64 `json:"cache_unit"`
}

type BrowserCacheListItem struct {
	CacheType    int64    `json:"cache_type"`
	CacheContent []string `json:"cache_content"`
	CacheTTL     int64    `json:"cache_ttl"`
	CacheUnit    int64    `json:"cache_unit"`
	CacheStatus  int64    `json:"cache_status"`
	Priority     int64    `json:"priority"`
}

type IPType int

type IpFilter struct {
	IpType         int64    `json:"ip_type"` // 黑白名单, 0黑1白
	IpList         []string `json:"ip_list"`
	EffectiveType  int64    `json:"effective_type"`
	EffectiveRules []string `json:"effective_rule"`
}

type Referer struct {
	RefererType  int64    `json:"referer_type"` //黑白名单, 0黑1白
	RefererList  []string `json:"referer_list"`
	IncludeEmpty int64    `json:"include_empty"` // 允许为空, 0:N 1:Y
	Status       int64    `json:"status"`        // 是否开启, 0:N 1Y
}

type UserAgent struct {
	AgentType      int64    `json:"agent_type"`
	AgentList      []string `json:"agent_list"`
	EffectiveType  int64    `json:"effective_type"`
	EffectiveRules []string `json:"effective_rule"`
}

type AuthConf struct {
	AuthManner       int64    `json:"auth_manner"`        //鉴权模式0 type_a/1 type_b 2/type_c 3/type_d
	AuthRange        int64    `json:"auth_range"`         //鉴权范围 0:全部内容 1:指定文件后缀 2:指定目录
	FileSuffix       []string `json:"file_suffix"`        //文件后缀
	InheritConf      string   `json:"inherit_conf"`       //继承配置
	InteritStartTime int64    `json:"interit_start_time"` //继承开始时间
	AuthKey          string   `json:"auth_key"`           //鉴权key
	AuthKeyBackup    string   `json:"auth_key_backup"`    //鉴权key备份
	AuthParameter    string   `json:"auth_parameter"`     //鉴权参数
	EncryptMannger   int64    `json:"encrypt_mannger"`    //加密方式 0:md5 1:sha256
	TimeFormat       int64    `json:"time_format"`        //0 十进制
	TimeValue        int64    `json:"time_value"`         //有效时间单位秒
	Status           int64    `json:"status"`             //鉴权状态 0:关闭 1:开启
}

type RemoteAuthConf struct {
	AuthUrl         string   `json:"auth_url"`         //鉴权url
	ReqMethod       int64    `json:"req_method"`       //请求方式
	FileType        int64    `json:"file_type"`        //文件类型 0:全部内容 1:指定文件后缀 2:指定目录 3:指定文件
	FileContent     []string `json:"file_content"`     //文件内容
	TimeoutDuration int64    `json:"timeout_duration"` //超时时间
	TimeoutAction   int64    `json:"timeout_action"`   //超时动作 0:继续 1:中断
	Status          int64    `json:"status"`           //鉴权状态 0:关闭 1:开启
}

type IpFrequencyConf struct {
	Frequency int64 `json:"frequency"` //频率
	Status    int64 `json:"status"`    //状态
}

type OriginUrlConf struct {
	MateMethod int64  `json:"mate_method"` //匹配方式 0:全部内容 1:指定文件后缀 2:指定目录
	RewriteUrl string `json:"rewrite_url"` //重写url
	TargetUrl  string `json:"target_url"`  //目标url
	Priority   int64  `json:"priority"`    //优先级
}

type RequestUrlRewriteConf struct {
	MateMethod   int64  `json:"mate_method"`   //匹配方式 0:全部内容 1:指定文件后缀 2:指定目录
	RewriteUrl   string `json:"rewrite_url"`   //重写url
	TargetUrl    string `json:"target_url"`    //目标url
	RedirectCode int64  `json:"redirect_code"` //跳转码
	Priority     int64  `json:"priority"`      //优先级
}

type CustomErrorPageConf struct {
	StatusCode   int64  `json:"status_code"`   //状态码
	RedirectCode int64  `json:"redirect_code"` //跳转码
	GoalAddress  string `json:"goal_address"`  //目标地址
}

type SpeedConf struct {
	RuleType    int64    `json:"rule_type"`    //规则类型 0:全部内容 1:指定文件后缀 2:指定目录
	RuleContent []string `json:"rule_content"` //规则内容
	SpeedValues int64    `json:"speed_values"` //加速值
}

type ResponseHeaderConf struct {
	ParameterKey   string `json:"parameter_key"`   //参数key
	ParameterValue string `json:"parameter_value"` //参数值
	Action         int64  `json:"action"`          // 0:删除 1:设置
}

type IntelligentCompressionConf struct {
	CompressType    int64    `json:"compress_type"`    //压缩类型
	CompressContent []string `json:"compress_content"` //压缩类型
	CompressMethod  int64    `json:"compress_method"`  //压缩方法 0:gzip 1:brotli
	Priority        int64    `json:"priority"`         //优先级
	Status          int64    `json:"status"`           //状态
}

type HttpsConf struct {
	HttpsStatus        int64   `json:"https_status"`
	TlsVersion         []int64 `json:"tls_version"`
	HttpTwo            int64   `json:"http_two"` // HTTP 2.0 开关, 0关1开
	JumpForceStatus    int64   `json:"jump_force_status"`
	JumpType           int64   `json:"jump_type"`
	JumpManner         int64   `json:"jump_manner"`
	HstsStatus         int64   `json:"hsts_status"`
	HstsExpirationTime int64   `json:"hsts_expiration_time"`
	HstsSubdomain      int64   `json:"hsts_subdomain"`
	OcspStatus         int64   `json:"ocsp_status"`
	QuicStatus         int64   `json:"quic_status"`
	CertName           string  `json:"cert_name"`
	CertValue          string  `json:"cert_value"`
	CertKey            string  `json:"cert_key"`
	CertType           int64   `json:"cert_type"`
}

type SwitchConf struct {
	IpFilterSwitch               int64 `json:"ip_filter_switch"`
	UserAgentSwitch              int64 `json:"user_agent_switch"`
	SpeedSwitch                  int64 `json:"speed_switch"`
	IntelligentCompressionSwitch int64 `json:"intelligent_compression_switch"`
}
