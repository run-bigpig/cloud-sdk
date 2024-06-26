package model

import "encoding/json"

type HttpCodeCacheRequest struct {
	DomainName         string               `path:"domain-name" json:"-"`
	HttpCodeCacheRules []*HttpCodeCacheRule `json:"http-code-cache-rules"`
}

type HttpCodeCacheRule struct {
	CacheTtl  string   `json:"cache-ttl"`
	HttpCodes []string `json:"http-codes"`
}

type HttpCodeCacheResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (o *HttpCodeCacheRequest) Marshal() []byte {
	b, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return b
}

func (o *HttpCodeCacheResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
