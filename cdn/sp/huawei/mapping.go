package huawei

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	"github.com/run-bigpig/cloud-sdk/cdn/consts"
	"strings"
)

// GetAreaCode 获取区域代码
func getAreaCode(t int64) model.DomainBodyServiceArea {
	switch t {
	case consts.AreaCodeChinaMainland:
		return model.GetDomainBodyServiceAreaEnum().MAINLAND_CHINA
	case consts.AreaCodeOversea:
		return model.GetDomainBodyServiceAreaEnum().OUTSIDE_MAINLAND_CHINA
	case consts.AreaCodeGlobal:
		return model.GetDomainBodyServiceAreaEnum().GLOBAL
	default:
		return model.GetDomainBodyServiceAreaEnum().MAINLAND_CHINA
	}
}

// MapAreaCode 获取区域代码
func mapAreaCode(t string) int64 {
	switch t {
	case "mainland_china":
		return consts.AreaCodeChinaMainland
	case "outside_mainland_china":
		return consts.AreaCodeOversea
	case "global":
		return consts.AreaCodeGlobal
	default:
		return consts.AreaCodeChinaMainland
	}

}

func getCdnType(t int64) string {
	switch t {
	case consts.CdnTypeCdn:
		return "cdn"
	case consts.CdnTypeDcdn:
		return "dcdn"
	case consts.CdnTypeScdn:
		return "scdn"
	default:
		return "cdn"
	}
}

func getRuleType(t int64) string {
	switch t {
	case consts.RuleTypeAll:
		return "all"
	case consts.RuleTypeFileSuffix:
		return "file_extension"
	case consts.RuleTypeDirectory:
		return "catalog"
	case consts.RuleTypePath:
		return "full_path"
	case consts.RuleTypeIndex:
		return "home_page"
	case consts.RuleTypeContentType:
		return "contentType"
	default:
		return "all"
	}
}

// GetRulePaths 获取缓存内容
func getRulePaths(t int64, data []string) string {
	switch {
	case len(data) == 0 || data == nil || t == consts.RuleTypeAll:
		return ""
	case t == consts.RuleTypeFileSuffix:
		return getFileType(data)
	case t == consts.RuleTypeDirectory, t == consts.RuleTypePath:
		return strings.ReplaceAll(strings.Join(data, ","), "*", "\\*")
	default:
		return "data"
	}
}

// GetChannelType 获取业务类型
func getChannelType(t int64) model.DomainBodyBusinessType {
	switch t {
	case consts.ChannelTypeWeb:
		return model.GetDomainBodyBusinessTypeEnum().WEB
	case consts.ChannelTypeDownload:
		return model.GetDomainBodyBusinessTypeEnum().DOWNLOAD
	case consts.ChannelTypeMedia:
		return model.GetDomainBodyBusinessTypeEnum().VIDEO
	case consts.ChannelTypeHybrid:
		return model.GetDomainBodyBusinessTypeEnum().WHOLE_SITE
	default:
		return model.GetDomainBodyBusinessTypeEnum().WEB
	}
}

// GetOriginProtocol 获取回源协议
func getOriginProtocol(t int64) string {
	switch t {
	case consts.OriginProtocolHttp:
		return "http"
	case consts.OriginProtocolHttps:
		return "https"
	case consts.OriginProtocolFollow:
		return "follow"
	default:
		return "http"
	}
}

// GetSwitch 获取开关状态
func getSwitch(t int64) string {
	switch t {
	case consts.SwitchOff:
		return "off"
	case consts.SwitchOn:
		return "on"
	default:
		return "off"
	}
}

// GetWhiteOrBlackList 获取黑白名单
func getWhiteOrBlackList(t int64) string {
	switch t {
	case consts.BlackList:
		return "black"
	case consts.WhiteList:
		return "white"
	default:
		return "black"
	}
}

// GetAccessAuthInheritType 获取继承类型
func getAccessAuthInheritType(inheritType string) string {
	inheritType = strings.ToLower(inheritType)
	if strings.Contains(inheritType, ",") || strings.Contains(inheritType, "，") {
		str := strings.ReplaceAll(inheritType, "，", ",")
		strSlice := strings.Split(str, ",")
		return strings.Join(strSlice, ",")
	}
	return inheritType
}

// GetAccessAuthInheritTimeType 获取继承时间类型
func getAccessAuthInheritTimeType(t int64) string {
	switch t {
	case consts.AccessAuthInheritTimeTypeParent:
		return "parent_url_time"
	case consts.AccessAuthInheritTimeTypeSystem:
		return "sys_time"
	default:
		return "parent_url_time"
	}
}

// GetOriginType 获取源站类型
func getOriginType(t int64) model.SourcesOriginType {
	switch t {
	case consts.OriginTypeIp:
		return model.GetSourcesOriginTypeEnum().IPADDR
	case consts.OriginTypeDomain:
		return model.GetSourcesOriginTypeEnum().DOMAIN
	case consts.OriginTypeBucket:
		return model.GetSourcesOriginTypeEnum().OBS_BUCKET
	default:
		return model.GetSourcesOriginTypeEnum().IPADDR
	}
}

// GetOriginAdvanceUrlMatchMode 获取源站url匹配模式
func getOriginAdvanceUrlMatchMode(t int64) string {
	switch t {
	case consts.OriginUrlMatchModeFile:
		return "file_extension"
	case consts.OriginUrlMatchModeDirectory:
		return "file_path"
	default:
		return "file_extension"
	}
}

// GetOriginAdvanceUrlMatchRule 获取源站url匹配规则
func getOriginAdvanceUrlMatchRule(t int64, data []string) string {
	newData := make([]string, 0, len(data))
	switch {
	case t == consts.OriginUrlMatchModeFile:
		for _, item := range data {
			newData = append(newData, fmt.Sprintf(".%s", item))
		}
		return strings.Join(newData, ";")
	case t == consts.OriginUrlMatchModeDirectory:
		for _, item := range data {
			newData = append(newData, fmt.Sprintf("/%s", strings.TrimLeft(item, "/")))
		}
		return strings.Join(newData, ";")
	default:
		return strings.Join(data, ";")
	}
}

// GetCacheParameterStatus 获取回源参数
func getCacheParameterStatus(t int64) string {
	switch t {
	case consts.CacheParameterStatusAll:
		return "ignore_url_params"
	case consts.CacheParameterStatusOff:
		return "full_url"
	case consts.CacheParameterStatusInclude:
		return "reserve_params"
	case consts.CacheParameterStatusExclude:
		return "del_params"
	default:
		return "gnore_url_params"
	}
}

// GetCacheParameterValues 获取参数值 最多十条
func getCacheParameterValues(t int64, data []string) string {
	if len(data) > 10 {
		data = data[:10]
	}
	switch {
	case t == consts.CacheParameterStatusInclude:
		return strings.Join(data, ",")
	case t == consts.CacheParameterStatusExclude:
		return strings.Join(data, ",")
	default:
		return ""
	}
}

// GetOriginHeaderAction 获取回源头操作
func getOriginHeaderAction(t int64) string {
	switch t {
	case consts.OriginHeaderActionDelete:
		return "delete"
	case consts.OriginHeaderActionSet:
		return "set"
	case consts.OriginHeaderActionAdd:
		return "add"
	default:
		return "add"
	}
}

// GetOriginUrlMateMethod 获取回源路径规则
func getOriginUrlMateMethod(t int64) string {
	switch t {
	case consts.OriginMateMethodAll:
		return "all"
	case consts.OriginMateMethodUrl:
		return "file_path"
	case consts.OriginMateMethodRegx:
		return "wildcard"
	case consts.OriginMateMethodPath:
		return "full_path"
	default:
		return "all"
	}
}

// GetAccessAuthRange 获取访问鉴权范围
func getAccessAuthRange(t int64) string {
	switch t {
	case consts.AccessAuthRangeAll:
		return "all"
	case consts.AccessAuthRangeInclude:
		return "all"
	case consts.AccessAuthRangeExclude:
		return "all"
	default:
		return "all"
	}
}

// GetAccessAuthMannerType 获取访问鉴权类型
func getAccessAuthMannerType(t int64) string {
	switch t {
	case consts.AccessAuthMannerTypeA:
		return "type_a"
	case consts.AccessAuthMannerTypeB:
		return "type_b"
	case consts.AccessAuthMannerTypeC:
		return "type_c1"
	case consts.AccessAuthMannerTypeD:
		return "type_c2"
	default:
		return "type_a"
	}
}

func getAccessAuthTimeFormat(authManner, timeFormat int64) string {
	switch authManner {
	case consts.AccessAuthMannerTypeA:
		return "dec"
	case consts.AccessAuthMannerTypeB:
		return "dec"
	case consts.AccessAuthMannerTypeC:
		return "hex"
	case consts.AccessAuthMannerTypeD:
		switch timeFormat {
		case consts.AccessAuthTimeFormatDec:
			return "dec"
		case consts.AccessAuthTimeFormatHex:
			return "hex"
		default:
			return "dec"
		}
	default:
		return "dec"
	}
}

// GetAccessRemoteAuthFileType 获取访问鉴权文件类型
func getAccessRemoteAuthFileType(t int64) string {
	switch t {
	case consts.FileTypeAll:
		return "all"
	case consts.FileTypeFile:
		return "specific_file"
	default:
		return "all"
	}
}

// GetAccessAuthEncryptManner 获取访问鉴权加密方式
func getAccessAuthEncryptManner(t int64) string {
	switch t {
	case consts.AccessAuthEncryptMannerMd5:
		return "md5"
	case consts.AccessAuthEncryptMannerSha256:
		return "sha256"
	default:
		return "md5"
	}
}

// GetAccessRemoteAuthTimeOutAction 获取访问鉴权超时动作
func getAccessRemoteAuthTimeOutAction(t int64) string {
	switch t {
	case consts.AccessRemoteAuthTimeOutActionReturn200:
		return "pass"
	case consts.AccessRemoteAuthTimeOutActionReturn403:
		return "forbid"
	default:
		return "pass"
	}
}

// GetCacheUnit 获取缓存单位
func getCacheUnit(t int64) string {
	switch t {
	case consts.CacheUnitSecond:
		return "s"
	case consts.CacheUnitMinute:
		return "m"
	case consts.CacheUnitHour:
		return "h"
	case consts.CacheUnitDay:
		return "d"
	default:
		return "s"
	}
}

func getCacheTtl(cacheTtl int32, cacheUnit int32) int32 {
	switch cacheUnit {
	case consts.CacheUnitSecond:
		return cacheTtl
	case consts.CacheUnitMinute:
		return cacheTtl * 60
	case consts.CacheUnitHour:
		return cacheTtl * 60 * 60
	case consts.CacheUnitDay:
		return cacheTtl * 60 * 60 * 24
	default:
		return cacheTtl
	}
}

// GetCacheBrowserCacheStatus 获取浏览器缓存状态
func getCacheBrowserCacheStatus(t int64) string {
	switch t {
	case consts.CacheStatusFollow:
		return "follow_origin"
	case consts.CacheStatusOn:
		return "ttl"
	case consts.CacheStatusOff:
		return "never"
	default:
		return "ttl"
	}
}

// GetRequestUrlRewriteType 获取请求url重写执行模式
func getRequestUrlRewriteType(t int64) string {
	switch t {
	case consts.RequestUrlRewriteTypeDirectory:
		return "catalog"
	case consts.RequestUrlRewriteTypeFullPath:
		return "full_path"
	default:
		return "catalog"
	}
}

// IsFollowOrigin 获取是否跟随回源
func isFollowOrigin(t int64) string {
	switch t {
	case consts.CacheStatusFollow:
		return "on"
	default:
		return "off"
	}
}

// GetTlsVersion 获取TLS版本
func getTlsVersion(t int64) string {
	switch t {
	case consts.HttpsTlsVersionSSLv0:
		return "TLSv1.0"
	case consts.HttpsTlsVersionSSLv1:
		return "TLSv1.1"
	case consts.HttpsTlsVersionSSLv2:
		return "TLSv1.2"
	case consts.HttpsTlsVersionTLSv3:
		return "TLSv1.3"
	default:
		return "TLSv1.0"
	}
}

// GetTlsVersions 获取TLS版本
func getTlsVersions(ts []int64) string {
	if ts == nil || len(ts) == 0 {
		return strings.Join([]string{"TLSv1.0,TLSv1.1,TLSv1.2"}, ",")
	}
	tlsVersions := make([]string, 0, len(ts))
	for _, t := range ts {
		tlsVersions = append(tlsVersions, getTlsVersion(t))
	}
	return strings.Join(tlsVersions, ",")
}

// GetHttpsCertificateType 获取https证书类型
func getHttpsCertificateType(t int64) string {
	switch t {
	case consts.HttpsCertificateTypeChina:
		return "server_sm"
	case consts.HttpsCertificateTypeGlobal:
		return "server"
	default:
		return "server_sm"
	}
}

// GetHttpsJumpType 获取https跳转类型
func getHttpsJumpType(t int64) string {
	switch t {
	case consts.HttpsJumpTypeHttp:
		return "http"
	case consts.HttpsJumpTypeHttps:
		return "https"
	default:
		return "http"
	}
}

// GetRedirectCode 获取重定向跳转码
func getRedirectCode(t int64) int64 {
	switch t {
	case consts.RedirectCode301:
		return 301
	case consts.RedirectCode302:
		return 302
	default:
		return 301
	}
}

// GetIntelligentCompressionCompressMethod 获取智能压缩压缩方法
func getIntelligentCompressionCompressMethod(t int64) string {
	switch t {
	case consts.IntelligentCompressionCompressMethodGzip:
		return "gzip"
	case consts.IntelligentCompressionCompressMethodBrotli:
		return "br"
	default:
		return "gzip"
	}
}

func getRemoteAuthRequestMethod(t int64) string {
	switch t {
	case consts.RequestMethodGet:
		return "GET"
	case consts.RequestMethodPost:
		return "POST"
	case consts.RequestMethodHead:
		return "HEAD"
	default:
		return "GET"
	}
}

// 判断切片是否是递增的
func isIncreasingSlice(slice []int64) bool {
	if len(slice) < 2 {
		return false
	}
	if len(slice) == 2 {
		return slice[1] == slice[0]+1
	}
	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[i-1]+1 {
			return false
		}
	}
	return true
}

func getFileType(data []string) string {
	newData := make([]string, len(data))
	copy(newData, data)
	for index, item := range newData {
		newData[index] = fmt.Sprintf(".%s", item)
	}
	return strings.Join(newData, ",")
}

// GetContentPurgePathMode 获取刷新模式
func getContentPurgePathMode(t int64) model.RefreshTaskRequestBodyMode {
	switch t {
	case consts.ContentPurgePathModeAll:
		return model.GetRefreshTaskRequestBodyModeEnum().ALL
	case consts.ContentPurgePathModeFile:
		return model.GetRefreshTaskRequestBodyModeEnum().DETECT_MODIFY_REFRESH
	default:
		return model.GetRefreshTaskRequestBodyModeEnum().ALL
	}
}

// GetIpProtocol 获取ip协议
func getIpProtocol(t int64) string {
	switch t {
	case consts.IpProtocolIpv4:
		return "ipv4"
	case consts.IpProtocolIpv6:
		return "ipv6"
	default:
		return "ipv4"
	}
}

// GetHttpProtocol 获取http协议
func getHttpProtocol(t int64) string {
	switch t {
	case consts.HttpProtocolHttp:
		return "http"
	case consts.HttpProtocolHttps:
		return "https"
	default:
		return "http"
	}
}

// GetShowContentPurgeOrPushStatus 获取展示刷新或预热状态
func getShowContentPurgeOrPushStatus(status string) int64 {
	switch status {
	case "task_done":
		return consts.ShowContentPurgeOrPushStatusSuccess
	case "task_inprocess":
		return consts.ShowContentPurgeOrPushStatusDoing
	default:
		return consts.ShowContentPurgeOrPushStatusFail
	}
}

// SetContentPurgeOrPushStatus 设置刷新或预热状态
func setContentPurgeOrPushStatus(status int64) model.ShowHistoryTasksRequestStatus {
	switch status {
	case consts.ShowContentPurgeOrPushStatusDoing:
		return model.GetShowHistoryTasksRequestStatusEnum().TASK_INPROCESS
	case consts.ShowContentPurgeOrPushStatusSuccess:
		return model.GetShowHistoryTasksRequestStatusEnum().TASK_DONE
	default:
		return model.GetShowHistoryTasksRequestStatusEnum().TASK_INPROCESS
	}
}

// GetShowContentPurgeType 获取刷新类型
func getShowContentPurgeType(t int64) model.ShowHistoryTasksRequestFileType {
	switch t {
	case consts.ShowContentPurgeTypeUrl:
		return model.GetShowHistoryTasksRequestFileTypeEnum().FILE
	case consts.ShowContentPurgeTypePath:
		return model.GetShowHistoryTasksRequestFileTypeEnum().DIRECTORY
	default:
		return model.GetShowHistoryTasksRequestFileTypeEnum().FILE
	}
}

// GetDomainStatus 获取域名状态
func getDomainStatus(status string) int64 {
	switch status {
	case "online":
		return consts.CdnDomainStatusDeployed
	case "offline":
		return consts.CdnDomainStatusStoped
	case "configuring":
		return consts.CdnDomainStatusDeploying
	case "configure_failed":
		return consts.CdnDomainStatusFaild
	default:
		return consts.CdnDomainStatusDeploying
	}
}

// SetDomainStatus 设置域名状态
func setDomainStatus(status int64) string {
	switch status {
	case consts.CdnDomainStatusDeployed:
		return "online"
	case consts.CdnDomainStatusStoped:
		return "offline"
	case consts.CdnDomainStatusDeploying, consts.CdnDomainStatusStoping:
		return "configuring"
	case consts.CdnDomainStatusFaild:
		return "configure_failed"
	default:
		return "configuring"
	}
}

func getPrimaryOrBack(t int64) int64 {
	switch t {
	case consts.OriginPriorityPrimary:
		return 70
	case consts.OriginPriorityBackup:
		return 30
	default:
		return 70
	}
}

// businessTypeToChannelType 获取业务类型
func businessTypeToChannelType(t string) int64 {
	switch t {
	case "web":
		return consts.ChannelTypeWeb
	case "download":
		return consts.ChannelTypeDownload
	case "video":
		return consts.ChannelTypeMedia
	case "wholeSite":
		return consts.ChannelTypeHybrid
	default:
		return consts.ChannelTypeWeb
	}
}

func getDataAccessMetricType(t int64) string {
	switch t {
	case consts.DataAccessMetricTypeFlux:
		return "flux"
	case consts.DataAccessMetricTypeBandwidth:
		return "bw"
	case consts.DataAccessMetricTypeRequest:
		return "req_num"
	case consts.DataAccessMetricTypeHitRequest:
		return "hit_num"
	case consts.DataAccessMetricTypeHitFlux:
		return "hit_flux"
	case consts.DataAccessMetricTypeStatusCode2xx:
		return "status_code_2xx"
	case consts.DataAccessMetricTypeStatusCode3xx:
		return "status_code_3xx"
	case consts.DataAccessMetricTypeStatusCode4xx:
		return "status_code_4xx"
	case consts.DataAccessMetricTypeStatusCode5xx:
		return "status_code_5xx"
	default:
		return "flux"
	}
}

// GetDataOriginMetricType 获取回源数据指标
func getDataOriginMetricType(t int64) string {
	switch t {
	case consts.DataOriginMetricTypeFlux:
		return "bs_flux"
	case consts.DataOriginMetricTypeBandwidth:
		return "bs_bw"
	case consts.DataOriginMetricTypeRequest:
		return "bs_num"
	case consts.DataOriginMetricTypeFailRequest:
		return "bs_fail_num"
	case consts.DataOriginMetricTypeStatusCode2xx:
		return "bs_status_code_2xx"
	case consts.DataOriginMetricTypeStatusCode3xx:
		return "bs_status_code_3xx"
	case consts.DataOriginMetricTypeStatusCode4xx:
		return "bs_status_code_4xx"
	case consts.DataOriginMetricTypeStatusCode5xx:
		return "bs_status_code_5xx"
	default:
		return "bs_flux"
	}
}

// GetDataStaticType 获取数据统计类型
func getAccessDataStaticType(t int64) string {
	switch t {
	case consts.DataStaticTypeSum:
		return "location_summary"
	case consts.DataStaticTypeDetail:
		return "location_detail"
	default:
		return "location_detail"
	}
}

func getOriginDataStaticType(t int64) string {
	switch t {
	case consts.DataStaticTypeSum:
		return "summary"
	case consts.DataStaticTypeDetail:
		return "detail"
	default:
		return "detail"
	}
}

// GetDataIntervalType 获取数据间隔
func getDataIntervalType(t int64) int64 {
	switch t {
	case consts.DataIntervalTypeFiveMinute:
		return 300
	case consts.DataIntervalTypeHour:
		return 3600
	case consts.DataIntervalTypeDay:
		return 86400
	default:
		return 300
	}
}

func getCountryCode(t int64) string {
	switch t {
	case consts.CountryCodeCn:
		return "cn"
	case consts.CountryCodeAe:
		return "ae"
	case consts.CountryCodeAu:
		return "au"
	case consts.CountryCodeBr:
		return "br"
	case consts.CountryCodeCa:
		return "ca"
	case consts.CountryCodeDe:
		return "de"
	case consts.CountryCodeCh:
		return "ch"
	case consts.CountryCodeEs:
		return "es"
	case consts.CountryCodeFr:
		return "fr"
	case consts.CountryCodeGb:
		return "gb"
	case consts.CountryCodeId:
		return "id"
	case consts.CountryCodeIl:
		return "il"
	case consts.CountryCodeIn:
		return "in"
	case consts.CountryCodeIt:
		return "it"
	case consts.CountryCodeJp:
		return "jp"
	case consts.CountryCodeKr:
		return "kr"
	case consts.CountryCodeMx:
		return "mx"
	case consts.CountryCodeMy:
		return "my"
	case consts.CountryCodeNl:
		return "nl"
	case consts.CountryCodeNo:
		return "no"
	case consts.CountryCodePh:
		return "ph"
	case consts.CountryCodeQa:
		return "qa"
	case consts.CountryCodeSa:
		return "sa"
	case consts.CountryCodeSe:
		return "se"
	case consts.CountryCodeSg:
		return "sg"
	case consts.CountryCodeTh:
		return "th"
	case consts.CountryCodeUs:
		return "us"
	case consts.CountryCodeVn:
		return "vn"
	case consts.CountryCodeZa:
		return "za"
	default:
		return "cn"
	}
}

// GetAllCountryCode 获取国家代码
func getAllCountryCode(igore []string) string {
	allCode := []string{"cn", "ae", "au", "br", "ca", "de", "ch", "es", "fr", "gb", "id", "il", "in", "it", "jp", "kr", "mx", "my", "nl", "no", "ph", "qa", "sa", "se", "sg", "th", "us", "vn", "za"}
	if igore == nil || len(igore) == 0 {
		return strings.Join(allCode, ",")
	}
	for _, v := range igore {
		for i, vv := range allCode {
			if v == vv {
				allCode = append(allCode[:i], allCode[i+1:]...)
			}
		}
	}
	return strings.Join(allCode, ",")
}

func getProvinceCode(t int64) string {
	switch t {
	case consts.ProvinceCodeAnhui:
		return "anhui"
	case consts.ProvinceCodeBeijing:
		return "beijing"
	case consts.ProvinceCodeChongqing:
		return "chongqing"
	case consts.ProvinceCodeFujian:
		return "fujian"
	case consts.ProvinceCodeGansu:
		return "gansu"
	case consts.ProvinceCodeGuangdong:
		return "guangdong"
	case consts.ProvinceCodeGuangxi:
		return "guangxi"
	case consts.ProvinceCodeGuizhou:
		return "guizhou"
	case consts.ProvinceCodeHainan:
		return "hainan"
	case consts.ProvinceCodeHebei:
		return "hebei"
	case consts.ProvinceCodeHeilongjiang:
		return "heilongjiang"
	case consts.ProvinceCodeHenan:
		return "henan"
	case consts.ProvinceCodeHubei:
		return "hubei"
	case consts.ProvinceCodeHunan:
		return "hunan"
	case consts.ProvinceCodeJiangsu:
		return "jiangsu"
	case consts.ProvinceCodeJiangxi:
		return "jiangxi"
	case consts.ProvinceCodeJilin:
		return "jilin"
	case consts.ProvinceCodeLiaoning:
		return "liaoling"
	case consts.ProvinceCodeNeimenggu:
		return "neimenggu"
	case consts.ProvinceCodeNingxia:
		return "ningxia"
	case consts.ProvinceCodeQinghai:
		return "qinghai"
	case consts.ProvinceCodeShaanxi:
		return "shaanxi"
	case consts.ProvinceCodeShandong:
		return "shandong"
	case consts.ProvinceCodeShanghai:
		return "shanghai"
	case consts.ProvinceCodeShanxi:
		return "shanxi"
	case consts.ProvinceCodeSichuan:
		return "sichuan"
	case consts.ProvinceCodeTianjin:
		return "tianjin"
	case consts.ProvinceCodeXinjiang:
		return "xinjiang"
	case consts.ProvinceCodeXizang:
		return "xizang"
	case consts.ProvinceCodeYunnan:
		return "yunnan"
	case consts.ProvinceCodeZhejiang:
		return "zhejiang"
	case consts.ProvinceCodeGangaotai:
		return "gangaotai"
	case consts.ProvinceCodeOther:
		return "qita"
	default:
		return "beijing"
	}
}

func getIspCode(t int64) string {
	switch t {
	case consts.IspCodeYidong:
		return "yidong"
	case consts.IspCodeDianxin:
		return "dianxin"
	case consts.IspCodeLiantong:
		return "liantong"
	case consts.IspCodeTietong:
		return "tietong"
	case consts.IspCodeJiaoyuwang:
		return "jiaoyuwang"
	case consts.IspCodeOther:
		return "qita"
	default:
		return "dianxin"
	}
}

func getTopUrlFilter(t int64) string {
	switch t {
	case consts.ListTopFilterFlux:
		return "flux"
	case consts.ListTopFilterRequest:
		return "req_num"
	default:
		return "flux"
	}
}
