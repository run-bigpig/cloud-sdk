package wangsu

import (
	"github.com/run-bigpig/cloud-sdk/cdn/consts"
	"strings"
)

func getServiceType(channelType int64) string {
	switch channelType {
	case consts.ChannelTypeWeb:
		return "web-https"
	case consts.ChannelTypeDownload:
		return "dl-https"
	case consts.ChannelTypeMedia:
		return "vod-https"
	case consts.ChannelTypeHybrid:
		return "wsa-https"
	default:
		return "web-https"
	}
}

func getServiceAreas(area int64) string {
	switch area {
	case consts.AreaCodeChinaMainland:
		return "cn"
	case consts.AreaCodeOversea:
		return strings.Join([]string{"am", "emea", "apac"}, ";")
	case consts.AreaCodeGlobal:
		return strings.Join([]string{"am", "emea", "apac", "cn"}, ";")
	default:
		return "cn"
	}
}
