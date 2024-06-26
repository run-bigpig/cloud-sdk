package common

import (
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common/auth"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common/model"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common/util"
	"net/http"
)

type Client struct {
	httpRequestMsg *model.HttpRequestMsg
	auth           *auth.Auth
}

type Unmarshaler interface {
	Unmarshal(data []byte) error
}

func NewClient(auth *auth.Auth) *Client {
	return &Client{
		auth: auth,
	}
}

func (c *Client) send(response Unmarshaler) error {
	res, err := util.Call(c.httpRequestMsg)
	if err != nil {
		return err
	}
	return response.Unmarshal(res)
}

// CreateDomain 创建域名
func (c *Client) CreateDomain(req *model.CreateDomainRequest) (*model.CreateDomainResponse, error) {
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/domain",
		Method: http.MethodPost,
		Body:   req.Marshal(),
	})
	resp := &model.CreateDomainResponse{}
	err := c.send(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// EnableDomain 启用域名
func (c *Client) EnableDomain(req *model.EnableDomainRequest) (*model.EnableDomainResponse, error) {
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/domain/enable",
		Method: http.MethodPost,
		Body:   req.Marshal(),
	})
	resp := &model.EnableDomainResponse{}
	err := c.send(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ShowDomainList 获取域名列表
func (c *Client) ShowDomainList(req *model.ShowDomainListRequest) (*model.ShowDomainListResponse, error) {
	s := &auth.SignParams{
		Url:    "/api/domain",
		Method: http.MethodGet,
	}
	if req.CnameLabel != "" {
		s.Url += "?cname_label=" + req.CnameLabel
	}
	c.auth.WithAuth(s)
	resp := &model.ShowDomainListResponse{}
	err := c.send(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DisableDomain 停用域名
func (c *Client) DisableDomain(req *model.DisableDomainRequest) (*model.DisableDomainResponse, error) {
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/domain/disable",
		Method: http.MethodPost,
		Body:   req.Marshal(),
	})
	resp := &model.DisableDomainResponse{}
	err := c.send(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetSrcConfig 设置回源策略
func (c *Client) SetSrcConfig(req *model.SrcConfigRequest) (*model.SrcConfigResponse, error) {
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/domain/setsrcconfig",
		Method: http.MethodPost,
		Body:   req.Marshal(),
	})
	resp := &model.SrcConfigResponse{}
	err := c.send(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CacheTime 设置缓存时间
func (c *Client) CacheTime(req *model.CacheTimeRequest) (*model.CacheTimeResponse, error) {
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/config/cachetime/" + req.DomainName,
		Method: http.MethodPut,
		Body:   req.Marshal(),
	})
	resp := &model.CacheTimeResponse{}
	err := c.send(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// HttpCodeCache 设置http状态码缓存
func (c *Client) HttpCodeCache(req *model.HttpCodeCacheRequest) (*model.HttpCodeCacheResponse, error) {
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/config/httpcodecache/" + req.DomainName,
		Method: http.MethodPut,
		Body:   req.Marshal(),
	})
	resp := &model.HttpCodeCacheResponse{}
	err := c.send(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AccessSpeed 设置限速
func (c *Client) AccessSpeed(req *model.AccessSpeedRequest) (model.AccessSpeedResponse, error) {
	req.DomainName = ""
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/config/accessspeed/" + req.DomainName,
		Method: http.MethodPut,
		Body:   req.Marshal(),
	})
	resp := model.AccessSpeedResponse{}
	err := c.send(&resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// HeaderModify 设置header
func (c *Client) HeaderModify(req *model.HeaderModifyRequest) (model.HeaderModifyResponse, error) {
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/config/headermodify/" + req.DomainName,
		Method: http.MethodPut,
		Body:   req.Marshal(),
	})
	resp := model.HeaderModifyResponse{}
	err := c.send(&resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// InnerRedirect 设置内网跳转
func (c *Client) InnerRedirect(req *model.InnerRedirectRequest) (model.InnerRedirectResponse, error) {
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/config/InnerRedirect/" + req.DomainName,
		Method: http.MethodPut,
		Body:   req.Marshal(),
	})
	resp := model.InnerRedirectResponse{}
	err := c.send(&resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// VisitControl 设置访问控制
func (c *Client) VisitControl(req *model.VisitControlRequest) (model.VisitControlResponse, error) {
	c.auth.WithAuth(&auth.SignParams{
		Url:    "/api/config/visitcontrol/" + req.DomainName,
		Method: http.MethodPut,
		Body:   req.Marshal(),
	})
	resp := model.VisitControlResponse{}
	err := c.send(&resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
