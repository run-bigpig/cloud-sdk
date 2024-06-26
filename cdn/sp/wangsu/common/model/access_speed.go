package model

import "encoding/json"

type AccessSpeedRequest struct {
	DomainName       string            `path:"domain-name" json:"-"`
	AccessSpeedRules []AccessSpeedRule `json:"access-speed-rules"`
}

type AccessSpeedRule struct {
	PathPattern string  `json:"path-pattern"`
	LimitMode   string  `json:"limit-mode"`
	StartSize   *string `json:"start-size,omitempty"`
	StartTime   *string `json:"start-time,omitempty"`
	StartSpeed  *string `json:"start-speed,omitempty"`
	Speed       string  `json:"speed"`
	Priority    int     `json:"priority"`
}

type AccessSpeedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (o *AccessSpeedRequest) Marshal() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return data
}

func (o *AccessSpeedResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
