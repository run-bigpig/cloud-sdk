package types

import "github.com/run-bigpig/cloud-sdk/cdn/entity"

const (
	UpdateBaseConf                   = "update_cdn_domain_base_conf"         // 更新基础配置
	UpdateArea                       = "update_area"                         // 更新加速区域
	UpdateOriginConf                 = "update_origin_conf"                  // 更新回源基础配置
	UpdateOriginServerConf           = "update_origin_server_conf"           // 更新回源服务器配置
	UpdateOriginAdvanceServerConf    = "update_origin_advance_server_conf"   // 更新回源高级配置
	UpdateOriginRequestHeaderConf    = "update_origin_request_header_conf"   // 更新回源请求头配置
	UpdateOriginUrlConf              = "update_origin_url_conf"              // 更新回源url配置
	UpdateIpFilterConf               = "update_ip_filter_conf"               // 更新ip过滤配置
	UpdateIpFrequencyConf            = "update_ip_frequency_conf"            // 更新ip频率配置
	UpdateRefererConf                = "update_referer_conf"                 // 更新referer配置
	UpdateUserAgentConf              = "update_user_agent_conf"              // 更新user agent配置
	UpdateSpeedConf                  = "update_speed_conf"                   // 更新加速配置
	UpdateAuthConf                   = "update_auth_conf"                    // 更新鉴权配置
	UpdateRemoteAuthConf             = "update_remote_auth_conf"             // 更新远程鉴权配置
	UpdateCacheListConf              = "update_cache_list_conf"              // 更新缓存列表配置
	UpdateCacheCodeConf              = "update_cache_code_conf"              // 更新缓存状态码配置
	UpdateBrowserCacheConf           = "update_browser_cache_conf"           // 更新浏览器缓存配置
	UpdateRequestUrlRewriteConf      = "update_request_url_rewrite_conf"     // 更新请求url重写配置
	UpdateCustomErrorPageConf        = "update_custom_error_page_conf"       // 更新自定义错误页面配置
	UpdateIntelligentCompressionConf = "update_intelligent_compression_conf" // 更新智能压缩配置
	UpdateResponseHeaderConf         = "update_response_header_conf"         // 更新响应头配置
	UpdateHttpsConf                  = "update_https_conf"                   // 更新https配置
	UpdateRecommendConf              = "update_recommend_conf"               // 更新推荐配置
	UpdateFullConf                   = "update_full_conf"                    // 更新全部配置
	OFF                              = 0
	ON                               = 1
	HuaWeiSdkName                    = "huawei"
	TencentSdkName                   = "tencent"
	KsYunSdkName                     = "ksyun"
	WangsuSdkName                    = "wangsu"
)
const (
	TypeA = iota
	TypeB
	TypeC
	TypeD
)

type (
	IcpVerifyRequest struct {
		Domain string `json:"domain"`
	}

	CreateDomainRequest struct {
		AreaCode       int64                      `json:"area_code"`       // 加速区域代码  0中国大陆 1境外 2全球
		ChannelType    int64                      `json:"channel_type"`    // 业务类型 0 网页 1 下载 2 点播 3 全站加速
		Domain         string                     `json:"domain"`          // 域名
		OriginProtocol int64                      `json:"origin_protocol"` // 回源协议 0 http 1 https
		Sources        []*entity.OriginServerConf `json:"sources"`         // 源站信息
	}

	CreateDomainResponse struct {
		DomainId string `json:"domain_id"` // 域名ID
		Cname    string `json:"cname"`     // CNAME
	}

	IpFilterConf struct {
		Status       int64 `json:"status"` // 开关 0 关闭 1 开启
		IpFilterConf []*entity.IpFilter
	}

	UserAgentConf struct {
		Status        int64 `json:"status"` // 开关 0 关闭 1 开启
		UserAgentConf []*entity.UserAgent
	}

	SpeedConf struct {
		Status    int64 `json:"status"` // 开关 0 关闭 1 开启
		SpeedConf []*entity.SpeedConf
	}

	IntelligentCompressionConf struct {
		Status                     int64 `json:"status"` // 开关 0 关闭 1 开启
		IntelligentCompressionConf []*entity.IntelligentCompressionConf
	}

	UpdateDomainRequest struct {
		UpdateAction               string                            `json:"update_action"` // 更新动作
		Domain                     string                            `json:"domain"`        // 域名`
		DomainId                   string                            `json:"domain_id"`     // 域名ID
		CdnDomain                  *entity.UpdateCdnDomainBaseConf   `json:"cdn_domain"`
		OriginConf                 *entity.OriginConf                `json:"origin_conf"`
		OriginServerConf           []*entity.OriginServerConf        `json:"origin_server_conf"`
		OriginAdvanceServerConf    []*entity.OriginAdvanceServerConf `json:"origin_advance_server_conf"`
		OriginRequestHeaderConf    []*entity.OriginRequestHeaderConf `json:"origin_request_header_conf"`
		OriginUrlConf              []*entity.OriginUrlConf           `json:"origin_url_conf"`
		IpFilterConf               *IpFilterConf                     `json:"ip_filter_conf"`
		IpFrequencyConf            *entity.IpFrequencyConf           `json:"ip_frequency_conf"`
		RefererConf                *entity.Referer                   `json:"referer_conf"`
		UserAgentConf              *UserAgentConf                    `json:"user_agent_conf"`
		SpeedConf                  *SpeedConf                        `json:"speed_conf"`
		AuthConf                   *entity.AuthConf                  `json:"auth_conf"`
		RemoteAuthConf             *entity.RemoteAuthConf            `json:"remote_auth_conf"`
		CacheListConf              []*entity.CacheListItem           `json:"cache_list_conf"`
		CacheCodeConf              []*entity.CacheCodeListItem       `json:"cache_code_conf"`
		BrowserCacheConf           []*entity.BrowserCacheListItem    `json:"browser_cache_conf"`
		RequestUrlRewriteConf      []*entity.RequestUrlRewriteConf   `json:"request_url_rewrite_conf"`
		CustomErrorPageConf        []*entity.CustomErrorPageConf     `json:"custom_error_page_conf"`
		IntelligentCompressionConf *IntelligentCompressionConf       `json:"intelligent_compression_conf"`
		ResponseHeaderConf         []*entity.ResponseHeaderConf      `json:"response_header_conf"`
		HttpsConf                  *entity.HttpsConf                 `json:"https_conf"`
	}

	ShowDomainDetailRequest struct {
		Domain   string `json:"domain"`    // 域名`
		DomainId string `json:"domain_id"` // 域名ID
	}

	ShowDomainDetailResponse struct {
		DomainId    string `json:"domain_id"`    // 域名ID
		Domain      string `json:"domain"`       // 域名
		AreaCode    int64  `json:"area_code"`    // 加速区域代码  0中国大陆 1境外 2全球
		ChannelType int64  `json:"channel_type"` // 业务类型 0 网页 1 下载 2 点播 3 全站加速
		Cname       string `json:"cname"`        // CNAME
		Status      int64  `json:"status"`       // 状态 0部署中 1部署成功 2停止中 3停止成功 4部署失败 5删除中 6已删除
		CreateTime  int64  `json:"create_time"`  // 创建时间
		UpdateTime  int64  `json:"update_time"`  // 修改时间
	}

	ShowDomainStatusListRequest struct {
		Page   int64 `json:"page"`
		Limit  int64 `json:"limit"`
		Status int64 `json:"status"` // 状态  1部署成功, 4停止成功
	}

	ShowDomainStatusListResponse struct {
		Total int64    `json:"total"`
		List  []string `json:"list"` // 域名
	}

	DeleteDomainRequest struct {
		DomainId string `json:"domain_id"` // 域名ID
		Domain   string `json:"domain"`    // 域名
	}

	DisableDomainRequest struct {
		DomainId string `json:"domain_id"` // 域名ID
		Domain   string `json:"domain"`    // 域名
	}

	EnableDomainRequest struct {
		DomainId string `json:"domain_id"` // 域名ID
		Domain   string `json:"domain"`    // 域名
	}

	PurgeUrlsCacheRequest struct {
		Urls      []string `json:"urls"`       // 刷新的URL列表，多个URL用逗号（半角）分隔，单次请求最多可提交100个URL
		UrlEncode bool     `json:"url_encode"` // 是否对URL进行编码
	}

	PurgePathCacheRequest struct {
		Paths     []string `json:"paths"`      // 刷新的目录列表，多个目录用逗号（半角）分隔，单次请求最多可提交10个目录
		Mode      int64    `json:"mode"`       // 刷新模式0 刷新全部 1刷新部分
		UrlEncode bool     `json:"url_encode"` // 是否对URL进行编码
	}

	PurgeCacheResponse struct {
		TaskId string `json:"task_id"` // 刷新任务ID
	}

	PushUrlsCacheRequest struct {
		Urls      []string `json:"urls"`       // 预热的URL列表，多个URL用逗号（半角）分隔，单次请求最多可提交1000个URL
		UrlEncode bool     `json:"url_encode"` // 是否对URL进行编码
	}

	PushUrlsCacheResponse struct {
		TaskId string `json:"task_id"` // 预热任务ID
	}

	ShowPurgeTaskStatusRequest struct {
		TaskId string `json:"task_id"` // 任务ID
	}

	ShowPurgeTaskListRequest struct {
		StartTime  int64  `json:"start_time"`  // 开始时间
		EndTime    int64  `json:"end_time"`    // 结束时间
		Page       int64  `json:"page"`        // 页数
		Limit      int64  `json:"limit"`       // 每页条数
		PurgeType  int64  `json:"purge_type"`  // 刷新类型 0 url 1 目录
		TaskStatus int64  `json:"task_status"` // 任务状态 0 失败 1成功  2进行中
		TimeZone   string `json:"time_zone"`   // 时区
	}

	ShowPurgeTaskStatusResponse struct {
		TaskId string `json:"task_id"` // 任务ID
		Status int64  `json:"status"`  // 任务状态 0 失败 1成功  2进行中
	}

	ShowPurgeTaskListResponse struct {
		Total int64    `json:"total"`
		List  []string `json:"list"` // 任务ID列表
	}

	ShowPushTaskStatusRequest struct {
		TaskId string `json:"task_id"` // 任务ID
	}

	ShowPushTaskListRequest struct {
		StartTime  int64  `json:"start_time"`  // 开始时间
		EndTime    int64  `json:"end_time"`    // 结束时间
		Page       int64  `json:"page"`        // 页数
		Limit      int64  `json:"limit"`       // 每页条数
		TaskStatus int64  `json:"task_status"` // 任务状态 0 失败 1成功  2进行中
		TimeZone   string `json:"time_zone"`   // 时区
	}

	ShowPushTaskStatusResponse struct {
		TaskId string `json:"task_id"` // 任务ID
		Status int64  `json:"status"`  // 任务状态 0 失败 1成功  2进行中
	}

	ShowPushTaskListResponse struct {
		Total int64    `json:"total"`
		List  []string `json:"list"` // 任务ID列表
	}

	CreateVerifyRecordRequest struct {
		DomainId string `json:"domain_id"` // 域名ID
		Domain   string `json:"domain"`    // 域名
	}

	CreateVerifyRecordResponse struct {
		RecordCode    string `json:"record"`
		FileVerifyUrl string `json:"file_verify_url"`
	}

	VerifyDomainRecordRequest struct {
		DomainId   string `json:"domain_id"`   // 域名ID
		Domain     string `json:"domain"`      // 域名
		VerifyType string `json:"verify_type"` // 验证类型
	}

	VerifyDomainRecordResponse struct {
		Result bool `json:"record_code"`
	}
)

type (
	DomainAccessDataStaticRequest struct {
		Domains     []string `json:"domain"`       // 域名
		Metric      int64    `json:"metric"`       // 指标 0 流量 1 带宽 2请求数 3 命中请求数 4命中流量 5 2xx 6 3xx 7 4xx 8 5xx
		StartTime   int64    `json:"start_time"`   // 开始时间戳
		EndTime     int64    `json:"end_time"`     // 结束时间戳
		Interval    int64    `json:"interval"`     // 时间间隔 0 5分钟 1 小时  2 天
		Isp         *int64   `json:"isp"`          // 运营商
		Area        int64    `json:"area"`         // 区域 0 中国大陆 1 中国境外
		AreaType    int64    `json:"area_type"`    // 区域类型 0 server 1 client
		District    *int64   `json:"district"`     // 省份/国家/地区
		Protocol    *int64   `json:"protocol"`     // 协议 0 http 1 https
		IpProtocol  *int64   `json:"ip_protocol"`  // ip协议 0 ipv4/ 1 ipv6
		Product     int64    `json:"product"`      // 产品 0 cdn/ 1 ecdn
		ChannelType int64    `json:"channel_type"` // 0 web 1 download 2 音视频 3 全站
		TimeZone    *string  `json:"time_zone"`    // 时区
	}

	StaticData struct {
		Value float64 `json:"value"`
		Time  int64   `json:"time"`
	}

	DomainAccessDataStaticResponse map[string][]*StaticData

	DomainOriginDataStaticRequest struct {
		Domains     []string `json:"domain"`       // 域名
		Metric      int64    `json:"metric"`       // 指标 0 回源流量 1 回源带宽 2 回源请求数 3 回源失败数 4 2xx 5 3xx 6 4xx 7 5xx
		StartTime   int64    `json:"start_time"`   // 开始时间
		EndTime     int64    `json:"end_time"`     // 结束时间
		Interval    int64    `json:"interval"`     // 时间间隔 0 5分钟 1 小时  2 天
		Area        int64    `json:"area"`         // 区域 0 中国大陆 1 中国境外
		ChannelType int64    `json:"channel_type"` // 0 web 1 download 2 音视频 3 全站
		TimeZone    *string  `json:"time_zone"`    // 时区
	}

	DomainOriginDataStaticResponse map[string][]*StaticData

	ListTopUrlDataStaticRequest struct {
		Filter      int64  `json:"filter"`       // 排序条件 0 流量 1请求数
		Domain      string `json:"domain"`       // 域名
		StartTime   int64  `json:"start_time"`   // 开始时间
		EndTime     int64  `json:"end_time"`     // 结束时间
		Product     int64  `json:"product"`      // 产品 0 cdn/ 1 ecdn
		ChannelType int64  `json:"channel_type"` // 0 web 1 download 2 音视频 3 全站
		Area        int64  `json:"area"`         // 区域 0 中国大陆 1 中国境外
	}

	ListTopUrlDataStaticResponse struct {
		Url   string  `json:"url"`   // url
		Value float64 `json:"value"` // 值
	}

	DataTotalDataResponse map[string]int64

	DomainAccessTotalDataRequest struct {
		Domains     []string `json:"domain"`       // 域名
		StartTime   int64    `json:"start_time"`   // 开始时间戳
		EndTime     int64    `json:"end_time"`     // 结束时间戳
		Area        int64    `json:"area"`         // 区域 0 中国大陆 1 中国境外
		Product     int64    `json:"product"`      // 产品 0 cdn/ 1 ecdn
		ChannelType int64    `json:"channel_type"` // 0 web 1 download 2 音视频 3 全站
		Metric      int64    `json:"metric"`       // 指标 0 流量 1 带宽 2请求数 3 命中请求数 4命中流量
		TimeZone    *string  `json:"time_zone"`    // 时区
	}

	DomainOriginTotalDataRequest struct {
		Domains     []string `json:"domain"`       // 域名
		StartTime   int64    `json:"start_time"`   // 开始时间戳
		EndTime     int64    `json:"end_time"`     // 结束时间戳
		Area        int64    `json:"area"`         // 区域 0 中国大陆 1 中国境外
		ChannelType int64    `json:"channel_type"` // 0 web 1 download 2 音视频 3 全站
		Metric      int64    `json:"metric"`       // 指标 0 回源流量 1 回源带宽 2 回源请求数 3 回源失败数
		TimeZone    *string  `json:"time_zone"`    // 时区
	}

	UserAccessRegionDistributionRequest struct {
		Domains     []string `json:"domain"`       // 域名
		StartTime   int64    `json:"start_time"`   // 开始时间戳
		EndTime     int64    `json:"end_time"`     // 结束时间
		Metric      int64    `json:"metric"`       // 指标 0 流量 2请求数
		Area        int64    `json:"area"`         // 区域 0 中国大陆 1 中国境外
		Product     int64    `json:"product"`      // 产品 0 cdn/ 1 ecdn
		ChannelType int64    `json:"channel_type"` // 0 web 1 download 2 音视频 3 全站
	}

	RegionDistribution struct {
		MainLandValue int64 `json:"mainland_value"` //境内值
		OverSeaValue  int64 `json:"oversea_value"`  // 境外值
	}
	UserAccessRegionDistributionResponse map[string]*RegionDistribution
)
