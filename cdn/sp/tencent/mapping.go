package tencent

import (
	"fmt"
	"github.com/run-bigpig/cloud-sdk/cdn/consts"
	"net"
	"strconv"
	"strings"
	"time"
)

// GetAreaCode 获取区域代码
func getAreaCode(t int64) string {
	switch t {
	case consts.AreaCodeChinaMainland:
		return "mainland"
	case consts.AreaCodeOversea:
		return "overseas"
	case consts.AreaCodeGlobal:
		return "global"
	default:
		return "mainland"
	}
}

// MapAreaCode 获取区域代码
func mapAreaCode(t string) int64 {
	switch t {
	case "mainland":
		return consts.AreaCodeChinaMainland
	case "overseas":
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

// ServiceTypeToChannelType 获取业务类型
func serviceTypeToChannelType(st string) int64 {
	switch st {
	case "web":
		return consts.ChannelTypeWeb
	case "download":
		return consts.ChannelTypeDownload
	case "media":
		return consts.ChannelTypeMedia
	case "hybrid":
		return consts.ChannelTypeHybrid
	default:
		return consts.ChannelTypeWeb
	}
}

// ConvertTimeZone 转换时区
func convertTimeZone(timezone string) string {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		fmt.Println(err)
		return "UTC+08:00"
	}
	_, offset := time.Now().In(location).Zone()
	return fmt.Sprintf("UTC+%02d:00", offset/3600)
}

func getProductType(t int64) string {
	switch t {
	case consts.ProductTypeCdn:
		return "cdn"
	case consts.ProductTypeEcdn:
		return "ecdn"
	default:
		return "cdn"
	}
}

func getRuleType(t int64) string {
	switch t {
	case consts.RuleTypeAll:
		return "all"
	case consts.RuleTypeFileSuffix:
		return "file"
	case consts.RuleTypeDirectory:
		return "directory"
	case consts.RuleTypePath:
		return "path"
	case consts.RuleTypeIndex:
		return "index"
	case consts.RuleTypeContentType:
		return "contentType"
	default:
		return "all"
	}
}

func getCompressRuleType(t int64) string {
	switch t {
	case consts.CompressRuleTypeAll:
		return "all"
	case consts.CompressRuleTypeFileSuffix:
		return "file"
	case consts.CompressRuleTypeContentType:
		return "contentType"
	default:
		return "all"
	}
}

// GetRulePaths 获取缓存内容
func getRulePaths(t int64, data []string) []string {
	switch {
	case consts.RuleTypeIndex == t:
		return []string{"/"}
	case len(data) == 0 || data == nil || t == consts.RuleTypeAll:
		return []string{"*"}
	default:
		return data
	}
}

// GetChannelType 获取业务类型
func getChannelType(t int64) string {
	switch t {
	case consts.ChannelTypeWeb:
		return "web"
	case consts.ChannelTypeDownload:
		return "download"
	case consts.ChannelTypeMedia:
		return "media"
	case consts.ChannelTypeHybrid:
		return "hybrid"
	default:
		return "web"
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
		return "blacklist"
	case consts.WhiteList:
		return "whitelist"
	default:
		return "blacklist"
	}
}

// GetIpfilterRules 获取ip过滤规则
func getIpfilterRules(t int64) string {
	switch t {
	case consts.RuleTypeAll:
		return "all"
	case consts.RuleTypeFileSuffix:
		return "file"
	case consts.RuleTypeDirectory:
		return "directory"
	case consts.RuleTypePath:
		return "path"
	default:
		return "all"
	}
}

// GetOriginType 获取源站类型
func getOriginType(t int64) string {
	switch t {
	case consts.OriginTypeIp:
		return "ip"
	case consts.OriginTypeDomain:
		return "domain"
	case consts.OriginTypeBucket:
		return "bucket"
	default:
		return "ip"
	}
}

// GetOriginTypeByAddresses 获取源站类型
func getOriginTypeByAddresses(addresses []string) string {
	var ip, ipv6, domain = "#_", "#_", "#_"
	for _, item := range addresses {
		ipinfo := net.ParseIP(item)
		if ipinfo != nil {
			if strings.Contains(item, ":") {
				ipv6 = "ipv6_"
			} else {
				ip = "ip_"
			}
		} else {
			domain = "domain_"
		}
	}
	return strings.Trim(strings.ReplaceAll(fmt.Sprintf("%s%s%s", ip, ipv6, domain), "#_", ""), "_")
}

// GetOriginAddressList 获取源站地址列表
func getOriginAddressList(address string, port int64, weight int64, sourcesLength int) string {
	addressList := make([]string, 0)
	addressList = append(addressList, address)
	if port != 0 {
		addressList = append(addressList, strconv.FormatInt(port, 10))
	}
	if weight != 0 && sourcesLength > 1 {
		addressList = append(addressList, strconv.FormatInt(weight, 10))
	}
	if len(addressList) == 0 {
		return ""
	}
	return strings.Join(addressList, ":")
}

// GetOriginUrlMatchMode 获取源站匹配模式
func getOriginUrlMatchMode(t int64) string {
	switch t {
	case consts.OriginUrlMatchModeFile:
		return "file"
	case consts.OriginUrlMatchModeDirectory:
		return "directory"
	default:
		return "file"
	}
}

// GetOriginHeaderAction 获取回源头操作
func getOriginHeaderAction(t int64) string {
	switch t {
	case consts.OriginHeaderActionDelete:
		return "del"
	case consts.OriginHeaderActionSet:
		return "set"
	case consts.OriginHeaderActionAdd:
		return "add"
	default:
		return "add"
	}
}

// GetOriginPriority 获取回源优先级
func getOriginPriority(t int64) string {
	switch t {
	case consts.OriginPriorityBackup:
		return "backup"
	case consts.OriginPriorityPrimary:
		return "primary"
	default:
		return "primary"
	}
}

// GetOriginMateMethodRegx 获取回源路径规则正则是否开启
func getOriginMateMethodRegx(t int64) bool {
	if t == consts.OriginMateMethodAll {
		return false
	}
	return true
}

// GetOriginMateMethodAll 获取回源路径规则是否开启
func getOriginMateMethodAll(t int64) bool {
	if t == consts.OriginMateMethodAll {
		return true
	}
	return false
}

// GetAccessAuthRange 获取访问鉴权范围
func getAccessAuthRange(t int64) string {
	switch t {
	case consts.AccessAuthRangeAll, consts.AccessAuthRangeInclude:
		return "blacklist"
	case consts.AccessAuthRangeExclude:
		return "whitelist"
	default:
		return "blacklist"
	}
}

func getAccessAuthRangeFileExtensions(data []string) []string {
	if data == nil || len(data) == 0 || data[0] == "" {
		return []string{"*"}
	}
	return data
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

// IsRequestUrlRewriteFullPath 获取是否全路径重写
func isRequestUrlRewriteFullPath(t int64) bool {
	switch t {
	case consts.RequestUrlRewriteTypeFullPath:
		return true
	default:
		return false
	}
}

// GetAccessAuthTimeFormat 获取访问鉴权时间格式
func getAccessAuthTimeFormat(t int64) string {
	switch t {
	case consts.AccessAuthTimeFormatDec:
		return "dec"
	case consts.AccessAuthTimeFormatHex:
		return "hex"
	default:
		return "dec"
	}
}

// GetAccessEffectiveType 获取访问鉴权有效期类型
func getAccessEffectiveType(t int64) string {
	switch t {
	case consts.AccessEffectiveTypeAll:
		return "all"
	case consts.AccessEffectiveTypeFileSuffix:
		return "file"
	case consts.AccessEffectiveTypeDirectory:
		return "directory"
	case consts.AccessEffectiveTypePath:
		return "path"
	case consts.AccessEffectiveTypeIndex:
		return "index"
	default:
		return "all"
	}
}

// GetAccessEffectiveContent 获取访问鉴权有效期内容
func getAccessEffectiveContent(t int64, data []string) []string {
	switch {
	case len(data) == 0 || data == nil || t == consts.AccessEffectiveTypeAll:
		return []string{"*"}
	case consts.AccessEffectiveTypeIndex == t:
		return []string{"/"}
	default:
		return data
	}
}

// GetAccessRemoteAuthTimeOutAction 获取访问鉴权超时动作
func getAccessRemoteAuthTimeOutAction(t int64) string {
	switch t {
	case consts.AccessRemoteAuthTimeOutActionReturn200:
		return "RETURN_200"
	case consts.AccessRemoteAuthTimeOutActionReturn403:
		return "RETURN_403"
	default:
		return "RETURN_200"
	}
}

func getRemoteAuthRequestMethod(t int64) string {
	switch t {
	case consts.RequestMethodGet:
		return "get"
	case consts.RequestMethodPost:
		return "post"
	case consts.RequestMethodHead:
		return "head"
	default:
		return "all"
	}
}

// GetCacheStatus 获取缓存状态
func getCacheStatus(t int64) string {
	switch t {
	case consts.CacheStatusOn:
		return "on"
	case consts.CacheStatusOff:
		return "off"
	default:
		return "on"
	}
}

// GetCacheUnit 获取缓存单位
func getCacheUnit(t int64) int64 {
	switch t {
	case consts.CacheUnitSecond:
		return 1
	case consts.CacheUnitMinute:
		return 60
	case consts.CacheUnitHour:
		return 3600
	case consts.CacheUnitDay:
		return 86400
	default:
		return 1
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

// IsCacheParameterFullPath IsFullPath 获取是否全路径缓存
func isCacheParameterFullPath(t int64) string {
	switch t {
	case consts.CacheParameterStatusOff:
		return "on"
	default:
		return "off"
	}
}

// GetTlsVersion 获取TLS版本
func getTlsVersion(t int64) string {
	switch t {
	case consts.HttpsTlsVersionSSLv0:
		return "TLSv1"
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
func getTlsVersions(ts []int64) []string {
	if ts == nil || len(ts) == 0 {
		return []string{"TLSv1,TLSv1.1,TLSv1.2"}
	}
	tlsVersions := make([]string, 0, len(ts))
	for _, t := range ts {
		tlsVersions = append(tlsVersions, getTlsVersion(t))
	}
	return tlsVersions
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
		return "brotli"
	default:
		return "gzip"
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

// GetShowContentPurgeOrPushStatus 获取显示内容刷新或预热状态
func getShowContentPurgeOrPushStatus(status string) int64 {
	switch status {
	case "fail":
		return consts.ShowContentPurgeOrPushStatusFail
	case "done":
		return consts.ShowContentPurgeOrPushStatusSuccess
	case "process":
		return consts.ShowContentPurgeOrPushStatusDoing
	default:
		return consts.ShowContentPurgeOrPushStatusFail
	}
}

// SetShowContentPurgeOrPushStatus 获取显示内容刷新或预热状态
func setShowContentPurgeOrPushStatus(status int64) string {
	switch status {
	case consts.ShowContentPurgeOrPushStatusDoing:
		return "process"
	case consts.ShowContentPurgeOrPushStatusSuccess:
		return "done"
	case consts.ShowContentPurgeOrPushStatusFail:
		return "fail"
	default:
		return "fail"
	}
}

// GetShowContentPurgeType 获取显示内容刷新类型
func getShowContentPurgeType(t int64) string {
	switch t {
	case consts.ShowContentPurgeTypeUrl:
		return "url"
	case consts.ShowContentPurgeTypePath:
		return "path"
	default:
		return "url"
	}
}

// GetDomainStatus 获取域名状态
func getDomainStatus(status string) int64 {
	switch status {
	case "rejected", "processing":
		return consts.CdnDomainStatusDeploying
	case "online":
		return consts.CdnDomainStatusDeployed
	case "closing":
		return consts.CdnDomainStatusStoping
	case "offline":
		return consts.CdnDomainStatusStoped
	default:
		return consts.CdnDomainStatusDeploying
	}
}

// SetDomainStatus
func setDomainStatus(status int64) string {
	switch status {
	case consts.CdnDomainStatusDeploying:
		return "processing"
	case consts.CdnDomainStatusDeployed:
		return "online"
	case consts.CdnDomainStatusStoped:
		return "offline"
	case consts.CdnDomainStatusDeleted:
		return "deleted"
	default:
		return "processing"
	}
}

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

func getDataAccessMetricType(t int64) string {
	switch t {
	case consts.DataAccessMetricTypeFlux:
		return "flux"
	case consts.DataAccessMetricTypeBandwidth:
		return "bandwidth"
	case consts.DataAccessMetricTypeRequest:
		return "request"
	case consts.DataAccessMetricTypeHitRequest:
		return "hitRequest"
	case consts.DataAccessMetricTypeHitFlux:
		return "hitFlux"
	case consts.DataAccessMetricTypeStatusCode2xx:
		return "2xx"
	case consts.DataAccessMetricTypeStatusCode3xx:
		return "3xx"
	case consts.DataAccessMetricTypeStatusCode4xx:
		return "4xx"
	case consts.DataAccessMetricTypeStatusCode5xx:
		return "5xx"
	default:
		return "flux"
	}
}

// GetDataOriginMetricType 获取回源数据指标
func getDataOriginMetricType(t int64) string {
	switch t {
	case consts.DataOriginMetricTypeFlux:
		return "flux"
	case consts.DataOriginMetricTypeBandwidth:
		return "bandwidth"
	case consts.DataOriginMetricTypeRequest:
		return "request"
	case consts.DataOriginMetricTypeFailRequest:
		return "failRequest"
	case consts.DataOriginMetricTypeStatusCode2xx:
		return "2xx"
	case consts.DataOriginMetricTypeStatusCode3xx:
		return "3xx"
	case consts.DataOriginMetricTypeStatusCode4xx:
		return "4xx"
	case consts.DataOriginMetricTypeStatusCode5xx:
		return "5xx"
	default:
		return "flux"
	}
}

// GetDataIntervalType 获取数据间隔
func getDataIntervalType(t int64) string {
	switch t {
	case consts.DataIntervalTypeFiveMinute:
		return "5min"
	case consts.DataIntervalTypeHour:
		return "hour"
	case consts.DataIntervalTypeDay:
		return "day"
	default:
		return "5min"
	}
}

func getCountryCode(t int64) int64 {
	switch t {
	case consts.CountryCodeCn:
		return 4460
	case consts.CountryCodeAe:
		return 386
	case consts.CountryCodeAu:
		return 4450
	case consts.CountryCodeBr:
		return 2613
	case consts.CountryCodeCa:
		return 3839
	case consts.CountryCodeDe:
		return 209
	case consts.CountryCodeCh:
		return 707
	case consts.CountryCodeEs:
		return 214
	case consts.CountryCodeFr:
		return 192
	case consts.CountryCodeGb:
		return 207
	case consts.CountryCodeId:
		return 1195
	case consts.CountryCodeIl:
		return 391
	case consts.CountryCodeIn:
		return 73
	case consts.CountryCodeIt:
		return 213
	case consts.CountryCodeJp:
		return 1044
	case consts.CountryCodeKr:
		return 3379
	case consts.CountryCodeMx:
		return 2626
	case consts.CountryCodeMy:
		return 3701
	case consts.CountryCodeNl:
		return 714
	case consts.CountryCodeNo:
		return 578
	case consts.CountryCodePh:
		return 2588
	case consts.CountryCodeQa:
		return 1233
	case consts.CountryCodeSa:
		return 471
	case consts.CountryCodeSe:
		return 208
	case consts.CountryCodeSg:
		return 1176
	case consts.CountryCodeTh:
		return 57
	case consts.CountryCodeUs:
		return 669
	case consts.CountryCodeVn:
		return 144
	case consts.CountryCodeZa:
		return 1559
	default:
		return 4460
	}
}

func getProvinceCode(t int64) int64 {
	switch t {
	case consts.ProvinceCodeAnhui:
		return 121
	case consts.ProvinceCodeBeijing:
		return 22
	case consts.ProvinceCodeChongqing:
		return 1051
	case consts.ProvinceCodeFujian:
		return 2
	case consts.ProvinceCodeGansu:
		return 1208
	case consts.ProvinceCodeGuangdong:
		return 4
	case consts.ProvinceCodeGuangxi:
		return 173
	case consts.ProvinceCodeGuizhou:
		return 118
	case consts.ProvinceCodeHainan:
		return 1441
	case consts.ProvinceCodeHebei:
		return 1069
	case consts.ProvinceCodeHeilongjiang:
		return 145
	case consts.ProvinceCodeHenan:
		return 182
	case consts.ProvinceCodeHubei:
		return 1135
	case consts.ProvinceCodeHunan:
		return 1466
	case consts.ProvinceCodeJiangsu:
		return 120
	case consts.ProvinceCodeJiangxi:
		return 1465
	case consts.ProvinceCodeJilin:
		return 1445
	case consts.ProvinceCodeLiaoning:
		return 1464
	case consts.ProvinceCodeNeimenggu:
		return 86
	case consts.ProvinceCodeNingxia:
		return 119
	case consts.ProvinceCodeQinghai:
		return 1467
	case consts.ProvinceCodeShaanxi:
		return 152
	case consts.ProvinceCodeShandong:
		return 122
	case consts.ProvinceCodeShanghai:
		return 1050
	case consts.ProvinceCodeShanxi:
		return 146
	case consts.ProvinceCodeSichuan:
		return 1068
	case consts.ProvinceCodeTianjin:
		return 1177
	case consts.ProvinceCodeXinjiang:
		return 1468
	case consts.ProvinceCodeXizang:
		return 1155
	case consts.ProvinceCodeYunnan:
		return 153
	case consts.ProvinceCodeZhejiang:
		return 1442
	case consts.ProvinceCodeGangaotai:
		return 1
	case consts.ProvinceCodeOther:
		return 0
	case consts.ProvinceCodeOverSea:
		return -1
	default:
		return 22
	}
}

func getIspCode(t int64) int64 {
	switch t {
	case consts.IspCodeYidong:
		return 1046
	case consts.IspCodeDianxin:
		return 2
	case consts.IspCodeLiantong:
		return 26
	case consts.IspCodeTietong:
		return 3947
	case consts.IspCodeJiaoyuwang:
		return 38
	case consts.IspCodeOther:
		return 0
	default:
		return 2
	}
}

func getTopUrlFilter(t int64) string {
	switch t {
	case consts.ListTopFilterFlux:
		return "flux"
	case consts.ListTopFilterRequest:
		return "request"
	default:
		return "flux"
	}
}
