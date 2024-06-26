package model

import "encoding/json"

type CreateDomainRequest struct {
	// 版本号，当前版本号1.0.0
	Version string `json:"version"`
	//需要接入CDN的域名。支持泛域名，以符号“.”开头，如：.example.com，泛域名也包含多级“a.b.example.com”。
	//如果example.com已备案，那么域名xx.example.com则不需要备案。
	DomainName string `json:"domain-name"`
	//加速域名的服务类型（一次只能提交一个服务类型）：
	//web/web-https：网页加速/网页加速-https
	//wsa/wsa-https：全站加速/全站加速-https
	//vodstream/vod-https：点播加速/点播加速-https
	//download/dl-https：下载加速/下载加速-https
	//livestream/live-https/cloudv-live：直播加速
	//v6sa/osv6：ipv6安全加速解决方案/IPv6一体化解决方案
	//注意：
	//1、service-type中的https不代表立即开启https，比如web-https中的https并不代表立刻支持https访问，需上传完证书后才可以支持https，切记！
	ServiceType string `json:"service-type"`
	//加速域名的加速区域，如果有需要根据区域限定资源覆盖时，才需要指定加速区域。未指定加速区域时，我们将按照客户开通的服务区域，以最优的资源覆盖提供加速服务。多个区域以分号分隔，支持配置的区域如下：cn（中国大陆）、am（美洲）、emea（欧洲、中东、非洲）、apac（亚太地区）
	ServiceAreas string `json:"service-areas"`
	//源站地址
	OriginConfig OriginConfig `json:"origin-config"`
	//备注信息，最大限制1000个字符
	Comment string `json:"comment"`
	//一级cname前缀，true表示使用域名名称作为cname前缀，否则，使用14位随机串（数字+字母）作为cname前缀。
	//注意：当前缀是泛域名时，则再增加wsall作为前缀。如.baidu.com.wscloudcdn.com，会生成wsall.baidu.com.wscloudcdn.com
	CnameWithCustomizedPrefix bool `json:"cname-with-customized-prefix"`
	//标识域名是否是纯海外加速的。
	//默认是否（false）
	//true ：表示客户域名纯海外加速
	//false：表示客户域名有在中国加速
	AccelerateNoChina bool `json:"accelerate-no-china"`
}

type OriginConfig struct {
	OriginIps               string `json:"origin-ips"`
	DefaultOriginHostHeader string `json:"default-origin-host-header"`
}

type CreateDomainResponse struct {
	//httpstatus=202; 表示成功调用新增域名接口，可使用header中的x-cnc-request-id查看当前新增域名的部署情况
	HttpStatusCode int `json:"http-status-code"`
	//唯一标示的id，用于查询每次请求的任务 （适用全部接口）
	XCncRequestId string `json:"x-cnc-request-id"`
	//用于访问该域名信息的URL，其中domain-id为我司云平台为该域名生成的唯一标示，其值为字符串。
	Location string `json:"location"`
	//域名的cname
	Cname string `json:"cname"`
	//错误代码，当HTTPStatus不为202时出现，表示当前请求调用的错误类型
	Code int `json:"code"`
	//响应信息，成功时为success
	Message string `json:"message"`
}

func (o *CreateDomainRequest) Marshal() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return data
}

func (o *CreateDomainResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
