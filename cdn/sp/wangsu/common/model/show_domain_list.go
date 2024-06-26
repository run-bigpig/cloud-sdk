package model

import "encoding/json"

type ShowDomainListRequest struct {
	//共用一级别名标示，可选入参，不选表示查询账号下所有域名
	//客户存在较多一级域名共用的需求，因此在接口中引入cname-label标识，即拥有相同cname-label的一组域名，共用一级cname。关于cname-label的具体使用方式和注意事项，请参看【创建加速域名】和【修改域名配置】接口
	CnameLabel string `json:"cname-label,omitempty"`
}

type ShowDomainListResponse struct {
	//域名ID
	DomainID string `json:"domain-id"`
	//加速域名的名称
	DomainName string `json:"domain-name"`
	//加速域名的服务类型，取值：
	//web/web-https：网页加速/网页加速-https
	//wsa/wsa-https：全站加速/全站加速-https
	//vodstream/vod-https：点播加速/点播加速-https
	//download/dl-https：下载加速/下载加速-https
	//livestream/live-https/cloudv-live：直播加速/直播加速-https/云直播
	//appa/s-appa：应用加速/应用安全加速解决方案
	ServiceType string `json:"service-type"`
	//加速域名对应的CNAME域名，例如：7nt6mrh7sdkslj.cdn30.com
	Cname string `json:"cname"`
	//加速域名的部署状态：Deployed表示该加速域名配置完成部署；InProgress表示该加速域名配置的部署任务还在进行中，可能处于排队、部署中或失败任意一种状态
	Status string `json:"status"`
	//加速域名的CDN服务状态：当取消加速域名CDN服务后，此项为false；当恢复加速域名CDN服务后，此项为true
	CdnServiceStatus string `json:"cdn-service-status"`
	//加速域名的启用状态：当禁用加速域名服务后，此项为false；当启用加速域名服务后，此项为true
	Enabled string `json:"enabled"`
}

func (o *ShowDomainListRequest) Marshal() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return data
}

func (o *ShowDomainListResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
