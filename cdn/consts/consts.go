package consts

// 状态数字开关
const (
	SwitchOff = iota
	SwitchOn
)

const (
	OFF = "off"
	ON  = "on"
)

// 文件类型
const (
	FileTypeAll = iota
	FileTypeFile
	FileTypeDirectory
)

// 黑白名单
const (
	BlackList = iota
	WhiteList
)

const (
	IpProtocolIpv4 = iota
	IpProtocolIpv6
)

const (
	HttpProtocolHttp = iota
	HttpProtocolHttps
	HttpProtocolQuic
)

// 加速区域代码
const (
	AreaCodeChinaMainland = iota
	AreaCodeOversea
	AreaCodeGlobal
)

// CdnType cdn类型
const (
	CdnTypeCdn = iota
	CdnTypeDcdn
	CdnTypeScdn
)

// RuleType 规则类型
const (
	RuleTypeAll = iota
	RuleTypeFileSuffix
	RuleTypeDirectory
	RuleTypePath
	RuleTypeIndex
	RuleTypeContentType
)

// ChannelType 业务类型
const (
	ChannelTypeWeb = iota
	ChannelTypeDownload
	ChannelTypeMedia
	ChannelTypeHybrid
)

// OriginProtocol 回源协议
const (
	OriginProtocolHttp = iota
	OriginProtocolHttps
	OriginProtocolFollow
)

// OriginType 源站类型
const (
	OriginTypeIp = iota
	OriginTypeDomain
	OriginTypeBucket
)

// OriginUrlMatchMode 源站url匹配模式
const (
	OriginUrlMatchModeFile = iota
	OriginUrlMatchModeDirectory
)

// OriginHeaderAction 源站请求头操作
const (
	OriginHeaderActionDelete = iota
	OriginHeaderActionSet
	OriginHeaderActionAdd
)

// OriginPriority
const (
	OriginPriorityPrimary = iota
	OriginPriorityBackup
)

// OriginMateMethod 源站匹配方法
const (
	OriginMateMethodAll = iota
	OriginMateMethodUrl
	OriginMateMethodRegx
	OriginMateMethodPath
)

// AccessAuthManner 访问鉴权方式
const (
	AccessAuthMannerTypeA = iota
	AccessAuthMannerTypeB
	AccessAuthMannerTypeC
	AccessAuthMannerTypeD
)

// AccessAuthInheritTimeType 访问鉴权继承时间类型
const (
	AccessAuthInheritTimeTypeParent = iota
	AccessAuthInheritTimeTypeSystem
)

// AccessAuthRange 访问鉴权范围
const (
	AccessAuthRangeAll = iota
	AccessAuthRangeInclude
	AccessAuthRangeExclude
)

// AccessAuthEncryptManner 访问鉴权加密方式
const (
	AccessAuthEncryptMannerMd5 = iota
	AccessAuthEncryptMannerSha256
)

// AccessAuthTimeFormat 访问鉴权时间格式
const (
	AccessAuthTimeFormatDec = iota
	AccessAuthTimeFormatHex
)

// AccessEffectiveType 访问限制生效类型
const (
	AccessEffectiveTypeAll = iota
	AccessEffectiveTypeFileSuffix
	AccessEffectiveTypeDirectory
	AccessEffectiveTypePath
	AccessEffectiveTypeIndex
)

// AccessRemoteAuthTimeOutAction 访问鉴权超时动作
const (
	AccessRemoteAuthTimeOutActionReturn200 = iota
	AccessRemoteAuthTimeOutActionReturn403
)

// CacheStatus 缓存状态
const (
	CacheStatusFollow = iota
	CacheStatusOff
	CacheStatusOn
)

// CacheUnit 缓存单位
const (
	CacheUnitSecond = iota
	CacheUnitMinute
	CacheUnitHour
	CacheUnitDay
)

// CacheParameterStatus 缓存参数状态
const (
	CacheParameterStatusOff = iota
	CacheParameterStatusAll
	CacheParameterStatusInclude
	CacheParameterStatusExclude
)

// RequestUrlRewriteType 请求URL重写类型
const (
	RequestUrlRewriteTypeDirectory = iota //目录
	RequestUrlRewriteTypeFullPath         //全路径
)

// RedirectCode 自定义错误页跳转码
const (
	RedirectCode301 = iota
	RedirectCode302
)

// IntelligentCompressionCompressMethod 智能压缩压缩方法
const (
	IntelligentCompressionCompressMethodGzip = iota
	IntelligentCompressionCompressMethodBrotli
)

// HttpsTlsVersion https tls版本
const (
	HttpsTlsVersionSSLv0 = iota
	HttpsTlsVersionSSLv1
	HttpsTlsVersionSSLv2
	HttpsTlsVersionTLSv3
)

// HttpsJumpType https跳转类型
const (
	HttpsJumpTypeHttp = iota
	HttpsJumpTypeHttps
)

// RequestMethod 请求方式
const (
	RequestMethodGet = iota
	RequestMethodPost
	RequestMethodHead
)

// HttpsCertificateType https证书类型
const (
	HttpsCertificateTypeChina = iota
	HttpsCertificateTypeGlobal
)

// CompressRuleType 压缩类型
const (
	CompressRuleTypeAll = iota
	CompressRuleTypeFileSuffix
	CompressRuleTypeContentType
)

const (
	ContentPurgePathModeAll = iota
	ContentPurgePathModeFile
)

const (
	ShowContentPurgeOrPushStatusDoing = iota
	ShowContentPurgeOrPushStatusSuccess
	ShowContentPurgeOrPushStatusFail
)

// Cdn域名状态
const (
	// Cdn域名状态
	CdnDomainStatusDeploying = iota // 部署中
	CdnDomainStatusDeployed         // 已部署
	CdnDomainStatusFaild            // 部署失败
	CdnDomainStatusStoping          // 停止中
	CdnDomainStatusStoped           // 已停止
	CdnDomainStatusDeleting         // 删除中
	CdnDomainStatusDeleted          // 已删除
)

// 内容刷新类型
const (
	ShowContentPurgeTypeUrl = iota
	ShowContentPurgeTypePath
)

// cdn产品类型
const (
	ProductTypeCdn  = iota //cdn
	ProductTypeEcdn        //ecdn
)

// 访问数据指标
const (
	DataAccessMetricTypeFlux          = iota //流量
	DataAccessMetricTypeBandwidth            //带宽
	DataAccessMetricTypeRequest              //请求总数
	DataAccessMetricTypeHitRequest           //请求命中数
	DataAccessMetricTypeHitFlux              //命中流量
	DataAccessMetricTypeStatusCode2xx        //状态码2xx
	DataAccessMetricTypeStatusCode3xx        //状态码3xx
	DataAccessMetricTypeStatusCode4xx        //状态码4xx
	DataAccessMetricTypeStatusCode5xx        //状态码5xx
)

// 回源数据指标
const (
	DataOriginMetricTypeFlux          = iota //回源流量
	DataOriginMetricTypeBandwidth            //回源带宽
	DataOriginMetricTypeRequest              //回源请求数
	DataOriginMetricTypeFailRequest          //回源失败请求数
	DataOriginMetricTypeStatusCode2xx        //回源状态码2xx
	DataOriginMetricTypeStatusCode3xx        //回源状态码3xx
	DataOriginMetricTypeStatusCode4xx        //回源状态码4xx
	DataOriginMetricTypeStatusCode5xx        //回源状态码5xx
)

// 粒度
const (
	DataIntervalTypeFiveMinute = iota
	DataIntervalTypeHour
	DataIntervalTypeDay
)

// 数据类型
const (
	DataStaticTypeSum = iota
	DataStaticTypeDetail
)

// ListTopFilter
const (
	ListTopFilterFlux    = iota //流量
	ListTopFilterRequest        //请求数
)

const (
	CountryCodeCn = iota + 1000 //中国
	CountryCodeAe               //阿联酋
	CountryCodeAu               //澳大利亚
	CountryCodeBr               //巴西
	CountryCodeCa               //加拿大
	CountryCodeCh               //瑞士
	CountryCodeDe               //德国
	CountryCodeEs               //西班牙
	CountryCodeFr               //法国
	CountryCodeGb               //英国
	CountryCodeId               //印度尼西亚
	CountryCodeIl               //以色列
	CountryCodeIn               //印度
	CountryCodeIt               //意大利
	CountryCodeJp               //日本
	CountryCodeKr               //韩国
	CountryCodeMx               //墨西哥
	CountryCodeMy               //马来西亚
	CountryCodeNl               //荷兰
	CountryCodeNo               //挪威
	CountryCodePh               //菲律宾
	CountryCodeQa               //卡塔尔
	CountryCodeSa               //沙特阿拉伯
	CountryCodeSe               //瑞典
	CountryCodeSg               //新加坡
	CountryCodeTh               //泰国
	CountryCodeUs               //美国
	CountryCodeVn               //越南
	CountryCodeZa               //南非
)

const (
	ProvinceCodeAnhui = iota + 2000
	ProvinceCodeBeijing
	ProvinceCodeChongqing
	ProvinceCodeFujian
	ProvinceCodeGansu
	ProvinceCodeGuangdong
	ProvinceCodeGuangxi
	ProvinceCodeGuizhou
	ProvinceCodeHainan
	ProvinceCodeHebei
	ProvinceCodeHeilongjiang
	ProvinceCodeHenan
	ProvinceCodeHubei
	ProvinceCodeHunan
	ProvinceCodeJiangsu
	ProvinceCodeJiangxi
	ProvinceCodeJilin
	ProvinceCodeLiaoning
	ProvinceCodeNeimenggu
	ProvinceCodeNingxia
	ProvinceCodeQinghai
	ProvinceCodeShaanxi
	ProvinceCodeShandong
	ProvinceCodeShanghai
	ProvinceCodeShanxi
	ProvinceCodeSichuan
	ProvinceCodeTianjin
	ProvinceCodeXinjiang
	ProvinceCodeXizang
	ProvinceCodeYunnan
	ProvinceCodeZhejiang
	ProvinceCodeGangaotai
	ProvinceCodeOther
	ProvinceCodeOverSea
)

const (
	IspCodeDianxin    = iota + 3000 //中国电信
	IspCodeYidong                   //中国移动
	IspCodeLiantong                 //中国联通
	IspCodeTietong                  //中国铁通
	IspCodeJiaoyuwang               //中国教育网
	IspCodeOther                    //其他运营商
)

const (
	KsYunTimeFormat = "2006-01-02T15:04-0700"
)
