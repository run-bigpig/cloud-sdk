package model

import "encoding/json"

type SrcConfigRequest struct {
	DomainName    string         `json:"domainName"`
	UseRange      *bool          `json:"useRange,omitempty"`
	Follow301     *bool          `json:"follow301,omitempty"`
	Follow302     *bool          `json:"follow302,omitempty"`
	AdvSrcSetting *AdvSrcSetting `json:"advSrcSetting"`
}

type AdvSrcSetting struct {
	UseAdvSrc     bool          `json:"useAdvSrc"`
	DetectUrl     string        `json:"detectUrl"`
	DetectPeriod  int           `json:"detectPeriod"`
	AdvSrcConfigs AdvSrcConfigs `json:"advSrcConfigs"`
}

type AdvSrcConfigs struct {
	MasterIps []string `json:"masterIps"`
	BackupIps []string `json:"backupIps"`
}

type SrcConfigResponse struct {
	//httpstatus=202; 表示成功调用新增域名接口，可使用header中的x-cnc-request-id查看当前新增域名的部署情况
	HttpStatusCode int `json:"http-status-code"`
	//唯一标示的id，用于查询每次请求的任务 （适用全部接口）
	XCncRequestId string `json:"x-cnc-request-id"`
	//错误代码，当HTTPStatus不为202时出现，表示当前请求调用的错误类型
	Code int `json:"code"`
	//响应信息，成功时为success
	Message string `json:"message"`
}

func (o *SrcConfigRequest) Marshal() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return data
}

func (o *SrcConfigResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
