package model

import "encoding/json"

type CacheTimeRequest struct {
	DomainName         string               `path:"domain-name" json:"-"`
	CacheTimeBehaviors []*CacheTimeBehavior `json:"cache-time-behaviors"`
}

type CacheTimeBehavior struct {
	PathPattern                string  `json:"path-pattern"`
	ExceptPathPattern          *string `json:"except-path-pattern,omitempty"`
	CustomPattern              *string `json:"custom-pattern,omitempty"`
	FileType                   *string `json:"file-type,omitempty"`
	CustomFileType             *string `json:"custom-file-type,omitempty"`
	SpecifyUrlPattern          *string `json:"specify-url-pattern,omitempty"`
	Directory                  *string `json:"directory,omitempty"`
	CacheTTL                   string  `json:"cache-ttl"`
	IgnoreCacheControl         *bool   `json:"ignore-cache-control,omitempty"`
	IsRespectServer            *bool   `json:"is-respect-server,omitempty"`
	IgnoreLetterCase           *bool   `json:"ignore-letter-case,omitempty"`
	ReloadManage               *string `json:"reload-manage,omitempty"`
	Priority                   int     `json:"priority"`
	IgnoreAuthenticationHeader *bool   `json:"ignore-authentication-header,omitempty"`
}

type CacheTimeResponse struct {
	PreDeployId string `json:"preDeployId"`
}

func (o *CacheTimeRequest) Marshal() []byte {
	data, err := json.Marshal(o.CacheTimeBehaviors)
	if err != nil {
		return nil
	}
	return data
}

func (o *CacheTimeResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
