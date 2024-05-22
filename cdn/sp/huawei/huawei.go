package huawei

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	huaweisdk "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"
	"github.com/run-bigpig/cloud-sdk/cdn/consts"
	"github.com/run-bigpig/cloud-sdk/cdn/entity"
	"github.com/run-bigpig/cloud-sdk/cdn/types"
	"github.com/run-bigpig/cloud-sdk/utils"
	"github.com/spf13/cast"
	"strings"
	"time"
)

type Huawei struct {
	config *Config
	client *huaweisdk.CdnClient
	ctx    context.Context
}

type Config struct {
	Region string
	Ak     string
	Sk     string
}

// NewHuaweiSdkClient creates a new Huawei client using the provided config.
func NewHuaweiSdkClient(ctx context.Context, conf *Config) *Huawei {
	auth, err := global.NewCredentialsBuilder().WithAk(conf.Ak).WithSk(conf.Sk).SafeBuild()
	if err != nil {
		return nil
	}
	rg, err := region.SafeValueOf(conf.Region)
	if err != nil {
		return nil
	}
	huaweiClient, err := huaweisdk.CdnClientBuilder().WithRegion(rg).WithCredential(auth).SafeBuild()
	if err != nil {
		return nil
	}
	client := huaweisdk.NewCdnClient(huaweiClient)
	return &Huawei{
		config: conf,
		client: client,
		ctx:    ctx,
	}
}

func (h *Huawei) GetSdkName() string {
	return types.HuaWeiSdkName
}

// CreateDomain 创建域名
func (h *Huawei) CreateDomain(req *types.CreateDomainRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	detail, err := h.ShowDomainDetail(&types.ShowDomainDetailRequest{Domain: req.Domain})
	if err == nil {
		if detail.Status == consts.CdnDomainStatusStoped {
			//启用域名
			err = h.EnableDomain(&types.EnableDomainRequest{Domain: req.Domain})
			if err != nil {
				return err
			}
		}
		//切换区域
		err = h.changeDomainArea(req.Domain, detail.AreaCode, req.AreaCode)
		if err != nil {
			return err
		}
		return nil
	}
	request := &model.CreateDomainRequest{}
	listSource := make([]model.Sources, 0)
	if req.Sources == nil {
		return errors.New("sources is nil")
	}
	for _, v := range req.Sources {
		activeStandby := int32(0)
		if v.OriginPriority == 0 {
			activeStandby = 1
		}
		listSource = append(listSource, model.Sources{
			OriginType:    getOriginType(v.OriginType),
			IpOrDomain:    v.OriginAddressList,
			ActiveStandby: activeStandby,
		})
	}
	domainBody := &model.DomainBody{
		DomainName:   req.Domain,
		BusinessType: getChannelType(req.ChannelType),
		Sources:      listSource,
		ServiceArea:  getAreaCode(req.AreaCode),
	}
	request.Body = &model.CreateDomainRequestBody{
		Domain: domainBody,
	}
	res, err := h.client.CreateDomain(request)
	if err != nil {
		return err
	}
	if res.HttpStatusCode < 200 || res.HttpStatusCode > 299 {
		return errors.New("create domain error")
	}
	return nil
}

// changeDomainArea 切换区域
func (h *Huawei) changeDomainArea(domain string, lastAreaCode, nowAreaCode int64) error {
	if lastAreaCode == nowAreaCode {
		return nil
	}
	//更新到全球
	if lastAreaCode != consts.AreaCodeGlobal {
		err := h.UpdateDomain(&types.UpdateDomainRequest{
			Domain:       domain,
			UpdateAction: types.UpdateArea,
			CdnDomain: &entity.UpdateCdnDomainBaseConf{
				AreaCode:    consts.AreaCodeGlobal,
				SupportIpv6: 0,
			},
		})
		if err != nil {
			return err
		}
	}
	if nowAreaCode == consts.AreaCodeGlobal {
		return nil
	}
	//更新到指定区域
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second * 60)
			err := h.UpdateDomain(&types.UpdateDomainRequest{
				Domain:       domain,
				UpdateAction: types.UpdateArea,
				CdnDomain: &entity.UpdateCdnDomainBaseConf{
					AreaCode: nowAreaCode,
				},
			})
			if err == nil {
				return
			}
		}
	}()
	return nil
}

func (h *Huawei) UpdateDomain(req *types.UpdateDomainRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	request := &model.UpdateDomainFullConfigRequest{}
	modityConfigBody := &model.ModifyDomainConfigRequestBody{}
	configs := &model.Configs{}
	updateDomain := h.newUpdateDomainModel(req, configs)
	switch req.UpdateAction {
	case types.UpdateBaseConf:
		updateDomain.WithBaseConf()
	case types.UpdateArea:
		updateDomain.WithArea()
	case types.UpdateOriginConf:
		updateDomain.WithOriginConf()
	case types.UpdateOriginServerConf:
		updateDomain.WithOriginServerConf()
	case types.UpdateOriginRequestHeaderConf:
		updateDomain.WithOriginRequestHeaderConf()
	case types.UpdateOriginUrlConf:
		updateDomain.WithOriginUrlConf()
	case types.UpdateOriginAdvanceServerConf:
		updateDomain.WithOriginAdvanceConf()
	case types.UpdateResponseHeaderConf:
		updateDomain.WithResponseHeaderConf()
	case types.UpdateIntelligentCompressionConf:
		updateDomain.WithIntelligentCompressionConf()
	case types.UpdateCustomErrorPageConf:
		updateDomain.WithCustomErrorPageConf()
	case types.UpdateCacheListConf:
		updateDomain.WithCacheListConf()
	case types.UpdateCacheCodeConf:
		updateDomain.WithCacheCodeConf()
	case types.UpdateBrowserCacheConf:
		updateDomain.WithCacheBrowserConf()
	case types.UpdateRequestUrlRewriteConf:
		updateDomain.WithRequestUrlRewriteConf()
	case types.UpdateIpFilterConf:
		updateDomain.WithIpFilterConf()
	case types.UpdateRefererConf:
		updateDomain.WithRefererConf()
	case types.UpdateUserAgentConf:
		updateDomain.WithUserAgentConf()
	case types.UpdateAuthConf:
		updateDomain.WithAuthConf()
	case types.UpdateRemoteAuthConf:
		updateDomain.WithRemoteAuthConf()
	case types.UpdateIpFrequencyConf:
		updateDomain.WithIpFrequencyConf()
	case types.UpdateHttpsConf:
		updateDomain.WithHttpsConf()
	case types.UpdateRecommendConf:
		updateDomain.WithCacheBrowserConf()
		updateDomain.WithCacheListConf()
		updateDomain.WithOriginConf()
		updateDomain.WithHttpsConf()
		updateDomain.WithIntelligentCompressionConf()
		updateDomain.WithIpFrequencyConf()
	case types.UpdateFullConf:
		updateDomain.WithBaseConf()
		updateDomain.WithOriginConf()
		updateDomain.WithOriginServerConf()
		updateDomain.WithOriginUrlConf()
		updateDomain.WithOriginRequestHeaderConf()
		updateDomain.WithIpFilterConf()
		updateDomain.WithIpFrequencyConf()
		updateDomain.WithRefererConf()
		updateDomain.WithUserAgentConf()
		updateDomain.WithAuthConf()
		updateDomain.WithRemoteAuthConf()
		updateDomain.WithCacheListConf()
		updateDomain.WithCacheCodeConf()
		updateDomain.WithCacheBrowserConf()
		updateDomain.WithRequestUrlRewriteConf()
		updateDomain.WithCustomErrorPageConf()
		updateDomain.WithIntelligentCompressionConf()
		updateDomain.WithResponseHeaderConf()
		updateDomain.WithHttpsConf()

	}
	modityConfigBody.Configs = configs
	request.Body = modityConfigBody
	request.DomainName = req.Domain
	response, err := h.client.UpdateDomainFullConfig(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode > 299 {
		return errors.New("update domain error")
	}
	return nil
}

// DisableDomain 禁用域名
func (h *Huawei) DisableDomain(req *types.DisableDomainRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	domain, err := h.ShowDomainDetail(&types.ShowDomainDetailRequest{Domain: req.Domain})
	if err != nil {
		return err
	}
	request := &model.DisableDomainRequest{}
	request.DomainId = domain.DomainId
	response, err := h.client.DisableDomain(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode > 299 {
		return errors.New("disable domain error")
	}
	return nil
}

// EnableDomain 启用域名
func (h *Huawei) EnableDomain(req *types.EnableDomainRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	domain, err := h.ShowDomainDetail(&types.ShowDomainDetailRequest{Domain: req.Domain})
	if err != nil {
		return err
	}
	request := &model.EnableDomainRequest{}
	request.DomainId = domain.DomainId
	response, err := h.client.EnableDomain(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode > 299 {
		return errors.New("enable domain error")
	}
	return nil
}

// DeleteDomain 删除域名
func (h *Huawei) DeleteDomain(req *types.DeleteDomainRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	domain, err := h.ShowDomainDetail(&types.ShowDomainDetailRequest{Domain: req.Domain})
	if err != nil {
		var e *sdkerr.ServiceResponseError
		if errors.As(err, &e) {
			if e.ErrorCode == "CDN.0170" {
				return nil
			}
		}
		return err
	}
	request := &model.DeleteDomainRequest{}
	request.DomainId = domain.DomainId
	response, err := h.client.DeleteDomain(request)
	if err != nil {
		return err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode > 299 {
		return errors.New("delete domain error")
	}
	return nil
}

// PurgePathCache 刷新目录缓存
func (h *Huawei) PurgePathCache(req *types.PurgePathCacheRequest) (*types.PurgeCacheResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	typeRefreshTask := model.GetRefreshTaskRequestBodyTypeEnum().DIRECTORY
	request := &model.CreateRefreshTasksRequest{}
	mode := getContentPurgePathMode(req.Mode)
	refreshTaskbody := &model.RefreshTaskRequestBody{
		Type:        &typeRefreshTask,
		Mode:        &mode,
		Urls:        req.Paths,
		ZhUrlEncode: &req.UrlEncode,
	}
	request.Body = &model.RefreshTaskRequest{
		RefreshTask: refreshTaskbody,
	}
	response, err := h.client.CreateRefreshTasks(request)
	if err != nil {
		return nil, err
	}
	return &types.PurgeCacheResponse{TaskId: *response.RefreshTask}, nil
}

// PurgeUrlsCache 刷新URL缓存
func (h *Huawei) PurgeUrlsCache(req *types.PurgeUrlsCacheRequest) (*types.PurgeCacheResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	typeRefreshTask := model.GetRefreshTaskRequestBodyTypeEnum().FILE
	request := &model.CreateRefreshTasksRequest{}
	refreshTaskbody := &model.RefreshTaskRequestBody{
		Type:        &typeRefreshTask,
		ZhUrlEncode: &req.UrlEncode,
		Urls:        req.Urls,
	}
	request.Body = &model.RefreshTaskRequest{
		RefreshTask: refreshTaskbody,
	}
	response, err := h.client.CreateRefreshTasks(request)
	if err != nil {
		return nil, err
	}
	return &types.PurgeCacheResponse{TaskId: *response.RefreshTask}, nil
}

// PushUrlsCache 推送URL预热
func (h *Huawei) PushUrlsCache(req *types.PushUrlsCacheRequest) (*types.PushUrlsCacheResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := &model.CreatePreheatingTasksRequest{}
	preheatingTaskbody := &model.PreheatingTaskRequestBody{
		Urls:        req.Urls,
		ZhUrlEncode: &req.UrlEncode,
	}
	request.Body = &model.PreheatingTaskRequest{
		PreheatingTask: preheatingTaskbody,
	}
	response, err := h.client.CreatePreheatingTasks(request)
	if err != nil {
		return nil, err
	}
	return &types.PushUrlsCacheResponse{TaskId: *response.PreheatingTask}, nil
}

// ShowPurgeTaskStatus 展示刷新任务状态
func (h *Huawei) ShowPurgeTaskStatus(req *types.ShowPurgeTaskStatusRequest) (*types.ShowPurgeTaskStatusResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := &model.ShowHistoryTaskDetailsRequest{}
	request.HistoryTasksId = req.TaskId
	request.PageSize = utils.Int32Ptr(1)
	request.PageNumber = utils.Int32Ptr(1)
	response, err := h.client.ShowHistoryTaskDetails(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode > 299 {
		return nil, errors.New("show purge task status error")
	}
	return &types.ShowPurgeTaskStatusResponse{
		TaskId: *response.Id,
		Status: getShowContentPurgeOrPushStatus(*response.Status),
	}, nil
}

// ShowPurgeTaskList 展示刷新任务列表
func (h *Huawei) ShowPurgeTaskList(req *types.ShowPurgeTaskListRequest) (*types.ShowPurgeTaskListResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := &model.ShowHistoryTasksRequest{}
	taskType := model.GetShowHistoryTasksRequestTaskTypeEnum().REFRESH
	request.TaskType = &taskType
	fileType := getShowContentPurgeType(req.PurgeType)
	request.FileType = &fileType
	request.PageSize = utils.Int32Ptr(int32(req.Limit))
	request.PageNumber = utils.Int32Ptr(int32(req.Page))
	if req.StartTime != 0 {
		request.StartDate = utils.Int64Ptr(req.StartTime * 1000)
	}
	if req.EndTime != 0 {
		request.EndDate = utils.Int64Ptr(req.EndTime * 1000)
	}
	status := setContentPurgeOrPushStatus(req.TaskStatus)
	request.Status = &status
	response, err := h.client.ShowHistoryTasks(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode > 299 {
		return nil, errors.New("show purge task list error")
	}
	if *response.Total == 0 || *response.Tasks == nil || len(*response.Tasks) == 0 {
		return &types.ShowPurgeTaskListResponse{Total: int64(*response.Total), List: []string{}}, nil
	}
	tasks := make([]string, 0, len(*response.Tasks))
	for _, v := range *response.Tasks {
		tasks = append(tasks, *v.Id)
	}
	return &types.ShowPurgeTaskListResponse{Total: int64(*response.Total), List: tasks}, nil
}

// ShowPushTaskStatus 展示预热任务状态
func (h *Huawei) ShowPushTaskStatus(req *types.ShowPushTaskStatusRequest) (*types.ShowPushTaskStatusResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := &model.ShowHistoryTaskDetailsRequest{}
	request.HistoryTasksId = req.TaskId
	request.PageSize = utils.Int32Ptr(1)
	request.PageNumber = utils.Int32Ptr(1)
	response, err := h.client.ShowHistoryTaskDetails(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode > 299 {
		return nil, errors.New("show push task status error")
	}
	return &types.ShowPushTaskStatusResponse{
		TaskId: *response.Id,
		Status: getShowContentPurgeOrPushStatus(*response.Status),
	}, nil
}

// ShowPushTaskList 展示预热任务列表
func (h *Huawei) ShowPushTaskList(req *types.ShowPushTaskListRequest) (*types.ShowPushTaskListResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := &model.ShowHistoryTasksRequest{}
	taskType := model.GetShowHistoryTasksRequestTaskTypeEnum().PREHEATING
	request.TaskType = &taskType
	request.PageSize = utils.Int32Ptr(int32(req.Limit))
	request.PageNumber = utils.Int32Ptr(int32(req.Page))
	if req.StartTime != 0 {
		request.StartDate = utils.Int64Ptr(req.StartTime * 1000)
	}
	if req.EndTime != 0 {
		request.EndDate = utils.Int64Ptr(req.EndTime * 1000)
	}
	status := setContentPurgeOrPushStatus(req.TaskStatus)
	request.Status = &status
	response, err := h.client.ShowHistoryTasks(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode > 299 {
		return nil, errors.New("show push task list error")
	}
	if *response.Total == 0 || *response.Tasks == nil || len(*response.Tasks) == 0 {
		return &types.ShowPushTaskListResponse{Total: int64(*response.Total), List: []string{}}, nil
	}
	tasks := make([]string, 0, len(*response.Tasks))
	for _, v := range *response.Tasks {
		tasks = append(tasks, *v.Id)
	}
	return &types.ShowPushTaskListResponse{Total: int64(*response.Total), List: tasks}, nil
}

// ShowDomainDetail 展示域名详情
func (h *Huawei) ShowDomainDetail(req *types.ShowDomainDetailRequest) (*types.ShowDomainDetailResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := &model.ShowDomainDetailByNameRequest{}
	request.DomainName = req.Domain
	res, err := h.client.ShowDomainDetailByName(request)
	if err != nil {
		return nil, err
	}
	if res.HttpStatusCode < 200 || res.HttpStatusCode > 299 {
		return nil, errors.New("show domain detail error")
	}
	return &types.ShowDomainDetailResponse{
		DomainId:   *res.Domain.Id,
		Domain:     *res.Domain.DomainName,
		Cname:      *res.Domain.Cname,
		Status:     getDomainStatus(*res.Domain.DomainStatus),
		CreateTime: *res.Domain.UpdateTime,
		UpdateTime: *res.Domain.UpdateTime,
		AreaCode:   mapAreaCode(res.Domain.ServiceArea.Value()),
	}, nil
}

// ShowDomainStatusList 展示域名状态列表
func (h *Huawei) ShowDomainStatusList(req *types.ShowDomainStatusListRequest) (*types.ShowDomainStatusListResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := &model.ListDomainsRequest{}
	request.DomainStatus = utils.StringPtr(setDomainStatus(req.Status))
	request.PageNumber = utils.Int32Ptr(int32(req.Page))
	request.PageSize = utils.Int32Ptr(int32(req.Limit))
	response, err := h.client.ListDomains(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode < 200 || response.HttpStatusCode > 299 {
		return nil, errors.New("show domain status list error")
	}
	if *response.Total == 0 || *response.Domains == nil || len(*response.Domains) == 0 {
		return &types.ShowDomainStatusListResponse{Total: int64(*response.Total), List: []string{}}, nil
	}
	domains := make([]string, 0, len(*response.Domains))
	for _, v := range *response.Domains {
		domains = append(domains, *v.DomainName)
	}
	return &types.ShowDomainStatusListResponse{Total: int64(*response.Total), List: domains}, nil
}

// CreateVerifyRecord 创建域名验证记录
func (h *Huawei) CreateVerifyRecord(req *types.CreateVerifyRecordRequest) (*types.CreateVerifyRecordResponse, error) {
	return nil, nil
}

// VerifyDomainRecord 验证域名记录
func (h *Huawei) VerifyDomainRecord(req *types.VerifyDomainRecordRequest) (*types.VerifyDomainRecordResponse, error) {
	return nil, nil
}

type UpdateDomainModel struct {
	req     *types.UpdateDomainRequest
	configs *model.Configs
}

func (h *Huawei) newUpdateDomainModel(req *types.UpdateDomainRequest, configs *model.Configs) *UpdateDomainModel {
	return &UpdateDomainModel{req: req, configs: configs}
}

func (u *UpdateDomainModel) WithBaseConf() {
	//u.configs.ServiceArea = utils.StringPtr(getAreaCode(u.req.CdnDomain.AreaCode).Value())
	u.configs.Ipv6Accelerate = utils.Int32Ptr(int32(u.req.CdnDomain.SupportIpv6))
}

func (u *UpdateDomainModel) WithArea() {
	u.configs.ServiceArea = utils.StringPtr(getAreaCode(u.req.CdnDomain.AreaCode).Value())
}

func (u *UpdateDomainModel) WithOriginConf() {
	u.configs.OriginProtocol = utils.StringPtr(getOriginProtocol(u.req.OriginConf.OriginProtocol))
	u.configs.OriginFollow302Status = utils.StringPtr(getSwitch(u.req.OriginConf.OriginFollow))
	u.configs.OriginRangeStatus = utils.StringPtr(getSwitch(u.req.OriginConf.OriginRange))
	if u.req.OriginConf.OriginTimeOut == 0 {
		u.req.OriginConf.OriginTimeOut = 30
	}
	if u.req.OriginConf.OriginSniSwitch != -1 {
		u.configs.Sni = &model.Sni{
			Status:     getSwitch(u.req.OriginConf.OriginSniSwitch),
			ServerName: utils.StringPtr(u.req.OriginConf.OriginSniValue),
		}
	}
	u.configs.OriginReceiveTimeout = utils.Int32Ptr(int32(u.req.OriginConf.OriginTimeOut))
}

func (u *UpdateDomainModel) WithOriginAdvanceConf() {
	filexbleOrigins := make([]model.FlexibleOrigins, 0)
	if u.req.OriginAdvanceServerConf != nil && len(u.req.OriginAdvanceServerConf) > 0 {
		backSources := make([]model.BackSources, 0)
		for _, v := range u.req.OriginAdvanceServerConf {
			backSources = append(backSources, model.BackSources{
				SourcesType:   getOriginType(v.OriginType).Value(),
				IpOrDomain:    v.OriginAddressList,
				ObsBucketType: nil,
				HttpPort:      utils.Int32Ptr(int32(v.OriginHttpPort)),
				HttpsPort:     utils.Int32Ptr(int32(v.OriginHttpsPort)),
			})
			filexbleOrigins = append(filexbleOrigins, model.FlexibleOrigins{
				MatchType:    getOriginAdvanceUrlMatchMode(v.UrlMatchMode),
				MatchPattern: getOriginAdvanceUrlMatchRule(v.UrlMatchMode, v.UrlMatchRule),
				Priority:     int32(v.OriginPriorityValue),
				BackSources:  backSources,
			})
		}
	}
	u.configs.FlexibleOrigin = &filexbleOrigins
}

func (u *UpdateDomainModel) WithOriginServerConf() {
	sources := make([]model.SourcesConfig, 0)
	u.configs.Sources = nil
	if u.req.OriginServerConf != nil && len(u.req.OriginServerConf) > 0 {
		for _, v := range u.req.OriginServerConf {
			sources = append(sources, model.SourcesConfig{
				OriginType: getOriginType(v.OriginType).Value(),
				OriginAddr: v.OriginAddressList,
				Priority:   int32(getPrimaryOrBack(v.OriginPriority)),
				Weight:     utils.Int32Ptr(int32(v.OriginWeight)),
				HttpPort:   utils.Int32Ptr(int32(v.OriginHttpPort)),
				HttpsPort:  utils.Int32Ptr(int32(v.OriginHttpsPort)),
				HostName:   utils.StringPtr(v.OriginHost),
			})
		}
		u.configs.Sources = &sources
	}
}

func (u *UpdateDomainModel) WithOriginRequestHeaderConf() {
	originRequestHeader := make([]model.OriginRequestHeader, 0)
	if u.req.OriginRequestHeaderConf != nil && len(u.req.OriginRequestHeaderConf) > 0 {
		for _, v := range u.req.OriginRequestHeaderConf {
			originRequestHeader = append(originRequestHeader, model.OriginRequestHeader{
				Name:   v.ParameterKey,
				Value:  &v.ParameterValue,
				Action: getOriginHeaderAction(v.Action),
			})
		}
	}
	u.configs.OriginRequestHeader = &originRequestHeader
}

func (u *UpdateDomainModel) WithOriginUrlConf() {
	originUrl := make([]model.OriginRequestUrlRewrite, 0)
	if u.req.OriginUrlConf != nil && len(u.req.OriginUrlConf) > 0 {
		for _, v := range u.req.OriginUrlConf {
			originUrl = append(originUrl, model.OriginRequestUrlRewrite{
				Priority:  int32(v.Priority),
				MatchType: getOriginUrlMateMethod(v.MateMethod),
				SourceUrl: utils.StringPtr(v.RewriteUrl),
				TargetUrl: v.TargetUrl,
			})
		}
	}
	u.configs.OriginRequestUrlRewrite = &originUrl
}

func (u *UpdateDomainModel) WithResponseHeaderConf() {
	httpResponseHeader := make([]model.HttpResponseHeader, 0)
	if u.req.ResponseHeaderConf != nil && len(u.req.ResponseHeaderConf) > 0 {
		for _, v := range u.req.ResponseHeaderConf {
			httpResponseHeader = append(httpResponseHeader, model.HttpResponseHeader{
				Name:   v.ParameterKey,
				Value:  &v.ParameterValue,
				Action: getOriginHeaderAction(v.Action),
			})
		}
	}
	u.configs.HttpResponseHeader = &httpResponseHeader
}

func (u *UpdateDomainModel) WithIntelligentCompressionConf() {
	compresses := make([]string, 0)
	compress := &model.Compress{
		Status: getSwitch(u.req.IntelligentCompressionConf.Status),
	}
	if u.req.IntelligentCompressionConf.IntelligentCompressionConf != nil && len(u.req.IntelligentCompressionConf.IntelligentCompressionConf) > 0 {
		for _, v := range u.req.IntelligentCompressionConf.IntelligentCompressionConf {
			compresses = append(compresses, v.CompressContent...)
		}
		compress.Type = utils.StringPtr(getIntelligentCompressionCompressMethod(u.req.IntelligentCompressionConf.IntelligentCompressionConf[0].CompressMethod))
		compress.FileType = utils.StringPtr(getFileType(compresses))
	}
	u.configs.Compress = compress
}

func (u *UpdateDomainModel) WithCustomErrorPageConf() {
	customErrorPage := make([]model.ErrorCodeRedirectRules, 0)
	if u.req.CustomErrorPageConf != nil && len(u.req.CustomErrorPageConf) > 0 {
		for _, v := range u.req.CustomErrorPageConf {
			customErrorPage = append(customErrorPage, model.ErrorCodeRedirectRules{
				ErrorCode:  int32(v.StatusCode),
				TargetCode: int32(getRedirectCode(v.RedirectCode)),
				TargetLink: v.GoalAddress,
			})
		}
	}
	u.configs.ErrorCodeRedirectRules = &customErrorPage
}

func (u *UpdateDomainModel) WithCacheListConf() {
	cacheRules := make([]model.CacheRules, 0)
	if u.req.CacheListConf != nil && len(u.req.CacheListConf) > 0 {
		for _, v := range u.req.CacheListConf {
			cacheRule := model.CacheRules{
				MatchType:         utils.StringPtr(getRuleType(v.CacheType)),
				Priority:          int32(v.Priority),
				TtlUnit:           getCacheUnit(v.CacheUnit),
				FollowOrigin:      utils.StringPtr(isFollowOrigin(v.CacheStatus)),
				UrlParameterType:  utils.StringPtr(getCacheParameterStatus(v.ParametersStatus)),
				UrlParameterValue: utils.StringPtr(getCacheParameterValues(v.ParametersStatus, v.ParametersValue)),
			}
			if v.CacheType != consts.RuleTypeAll && v.CacheType != consts.RuleTypeIndex {
				cacheRule.MatchValue = utils.StringPtr(getRulePaths(v.CacheType, v.CacheContent))
			}
			if v.CacheTTL > 0 {
				cacheRule.Ttl = utils.Int32Ptr(int32(v.CacheTTL))
			}
			cacheRules = append(cacheRules, cacheRule)
		}
	}
	u.configs.CacheRules = &cacheRules
}

// WithCacheBrowserConf 浏览器缓存配置
func (u *UpdateDomainModel) WithCacheBrowserConf() {
	cacheBrowsers := make([]model.BrowserCacheRules, 0)
	if u.req.BrowserCacheConf != nil && len(u.req.BrowserCacheConf) > 0 {
		for _, v := range u.req.BrowserCacheConf {
			browserCache := model.BrowserCacheRules{
				CacheType: getCacheBrowserCacheStatus(v.CacheStatus),
				Condition: &model.BrowserCacheRulesCondition{
					MatchType: getRuleType(v.CacheType),
					Priority:  int32(v.Priority),
				},
			}
			if v.CacheType != consts.RuleTypeAll && v.CacheType != consts.RuleTypeIndex {
				browserCache.Condition.MatchValue = utils.StringPtr(getRulePaths(v.CacheType, v.CacheContent))
			}
			if v.CacheTTL > 0 {
				browserCache.Ttl = utils.Int32Ptr(int32(v.CacheTTL))
				browserCache.TtlUnit = utils.StringPtr(getCacheUnit(v.CacheUnit))
			}
			cacheBrowsers = append(cacheBrowsers, browserCache)
		}
	}
	u.configs.BrowserCacheRules = &cacheBrowsers
}

func (u *UpdateDomainModel) WithCacheCodeConf() {
	cacheCodeRules := make([]model.ErrorCodeCache, 0)
	if u.req.CacheCodeConf != nil && len(u.req.CacheCodeConf) > 0 {
		for _, v := range u.req.CacheCodeConf {
			cacheCodeRules = append(cacheCodeRules, model.ErrorCodeCache{
				Code: utils.Int32Ptr(int32(v.HttpCode)),
				Ttl:  utils.Int32Ptr(getCacheTtl(int32(v.CacheTTL), int32(v.CacheUnit))),
			})
		}
	}
	u.configs.ErrorCodeCache = &cacheCodeRules
}

// WithRequestUrlRewriteConf 请求URL重写配置
func (u *UpdateDomainModel) WithRequestUrlRewriteConf() {
	requestUrlRewrite := make([]model.RequestUrlRewrite, 0)
	if u.req.RequestUrlRewriteConf != nil {
		for _, v := range u.req.RequestUrlRewriteConf {
			requestUrlRewrite = append(requestUrlRewrite, model.RequestUrlRewrite{
				ExecutionMode:      "redirect",
				RedirectUrl:        v.TargetUrl,
				RedirectStatusCode: utils.Int32Ptr(int32(getRedirectCode(v.RedirectCode))),
				Condition: &model.UrlRewriteCondition{
					MatchType:  getRequestUrlRewriteType(v.MateMethod),
					MatchValue: v.RewriteUrl,
					Priority:   int32(v.Priority),
				},
			})
		}
	}
	u.configs.RequestUrlRewrite = &requestUrlRewrite
}

func (u *UpdateDomainModel) WithIpFilterConf() {
	ipFilter := &model.IpFilter{}
	if u.req.IpFilterConf.Status == consts.SwitchOff {
		ipFilter.Type = consts.OFF
	} else {
		if u.req.IpFilterConf.IpFilterConf != nil && len(u.req.IpFilterConf.IpFilterConf) > 0 {
			ipFilterList := make([]string, 0)
			ipFilter.Type = getWhiteOrBlackList(u.req.IpFilterConf.IpFilterConf[0].IpType)
			for _, v := range u.req.IpFilterConf.IpFilterConf {
				ipFilterList = append(ipFilterList, v.IpList...)
			}
			ipFilter.Value = utils.StringPtr(strings.Join(ipFilterList, ","))
		}
	}
	u.configs.IpFilter = ipFilter
}

func (u *UpdateDomainModel) WithRefererConf() {
	refererConfig := &model.RefererConfig{
		IncludeEmpty: utils.BoolPtr(cast.ToBool(u.req.RefererConf.IncludeEmpty)),
	}
	if u.req.RefererConf.Status == consts.SwitchOff {
		refererConfig.Type = consts.OFF
	} else {
		refererConfig.Type = getWhiteOrBlackList(u.req.RefererConf.RefererType)
	}
	refererConfig.Value = utils.StringPtr(strings.Join(u.req.RefererConf.RefererList, ","))
	u.configs.Referer = refererConfig
}

func (u *UpdateDomainModel) WithUserAgentConf() {
	userAgentFilter := &model.UserAgentFilter{}
	if u.req.UserAgentConf.Status == consts.SwitchOff {
		userAgentFilter.Type = consts.OFF
	} else {
		if u.req.UserAgentConf.UserAgentConf != nil && len(u.req.UserAgentConf.UserAgentConf) > 0 {
			userAgentFilter.Type = getWhiteOrBlackList(u.req.UserAgentConf.UserAgentConf[0].AgentType)
		}
		userAgentList := make([]string, 0)
		for _, v := range u.req.UserAgentConf.UserAgentConf {
			userAgentList = append(userAgentList, v.AgentList...)
		}
		userAgentFilter.Value = utils.StringPtr(strings.Join(userAgentList, ","))
	}
	u.configs.UserAgentFilter = userAgentFilter
}

func (u *UpdateDomainModel) WithAuthConf() {
	InheritConfig := &model.InheritConfig{}
	InheritConfig.Status = getSwitch(0)
	if u.req.AuthConf.InheritConf != "" {
		InheritConfig.InheritType = utils.StringPtr(getAccessAuthInheritType(u.req.AuthConf.InheritConf))
		InheritConfig.InheritTimeType = utils.StringPtr(getAccessAuthInheritTimeType(u.req.AuthConf.InteritStartTime))
	}
	u.configs.UrlAuth = &model.UrlAuth{
		Status:        getSwitch(u.req.AuthConf.Status),
		Type:          utils.StringPtr(getAccessAuthMannerType(u.req.AuthConf.AuthManner)),
		ExpireTime:    utils.Int32Ptr(int32(u.req.AuthConf.TimeValue)),
		SignMethod:    utils.StringPtr(getAccessAuthEncryptManner(u.req.AuthConf.EncryptMannger)),
		MatchType:     utils.StringPtr(getAccessAuthRange(u.req.AuthConf.AuthRange)),
		InheritConfig: InheritConfig,
		Key:           &u.req.AuthConf.AuthKey,
		BackupKey:     &u.req.AuthConf.AuthKeyBackup,
		SignArg:       &u.req.AuthConf.AuthParameter,
		TimeFormat:    utils.StringPtr(getAccessAuthTimeFormat(u.req.AuthConf.AuthManner, u.req.AuthConf.TimeFormat)),
	}
}

func (u *UpdateDomainModel) WithRemoteAuthConf() {
	customArgsRules := make([]model.CustomArgs, 0)
	customHeadersRules := make([]model.CustomArgs, 0)
	u.configs.RemoteAuth = &model.CommonRemoteAuth{
		RemoteAuthentication: getSwitch(u.req.RemoteAuthConf.Status),
		RemoteAuthRules: &model.RemoteAuthRule{
			AuthServer:            u.req.RemoteAuthConf.AuthUrl,
			RequestMethod:         getRemoteAuthRequestMethod(u.req.RemoteAuthConf.ReqMethod),
			FileTypeSetting:       getAccessRemoteAuthFileType(u.req.RemoteAuthConf.FileType),
			ReserveArgsSetting:    "reserve_all_args",
			ReserveArgs:           utils.StringPtr(""),
			AddCustomArgsRules:    &customArgsRules,
			ReserveHeadersSetting: "reserve_all_headers",
			AddCustomHeadersRules: &customHeadersRules,
			AuthSuccessStatus:     "200",
			AuthFailedStatus:      "403",
			ResponseStatus:        "403",
			Timeout:               int32(u.req.RemoteAuthConf.TimeoutDuration),
			TimeoutAction:         getAccessRemoteAuthTimeOutAction(u.req.RemoteAuthConf.TimeoutAction),
			ReserveHeaders:        utils.StringPtr(""),
		},
	}
	if u.req.RemoteAuthConf.FileType == consts.FileTypeFile {
		u.configs.RemoteAuth.RemoteAuthRules.SpecifiedFileType = utils.StringPtr(strings.Join(u.req.RemoteAuthConf.FileContent, "|"))
	}
}

func (u *UpdateDomainModel) WithIpFrequencyConf() {
	u.configs.IpFrequencyLimit = &model.IpFrequencyLimit{
		Status: getSwitch(u.req.IpFrequencyConf.Status),
		Qps:    utils.Int32Ptr(int32(u.req.IpFrequencyConf.Frequency)),
	}
}

func (u *UpdateDomainModel) WithHttpsConf() {
	if u.req.HttpsConf.HttpsStatus == consts.SwitchOff {
		u.configs.Https = &model.HttpPutBody{
			HttpsStatus: utils.StringPtr(getSwitch(u.req.HttpsConf.HttpsStatus)),
		}
		return
	}
	forceRedirect := &model.ForceRedirectConfig{
		Status:       getSwitch(u.req.HttpsConf.JumpForceStatus),
		Type:         utils.StringPtr(getHttpsJumpType(u.req.HttpsConf.JumpType)),
		RedirectCode: utils.Int32Ptr(int32(getRedirectCode(u.req.HttpsConf.JumpManner))),
	}
	u.configs.ForceRedirect = forceRedirect
	u.configs.Https = &model.HttpPutBody{
		HttpsStatus:       utils.StringPtr(getSwitch(u.req.HttpsConf.HttpsStatus)),
		CertificateName:   utils.StringPtr(u.req.HttpsConf.CertName),
		CertificateValue:  utils.StringPtr(u.req.HttpsConf.CertValue),
		PrivateKey:        utils.StringPtr(u.req.HttpsConf.CertKey),
		CertificateSource: utils.Int32Ptr(0),
		//CertificateType:    utils.StringPtr(getHttpsCertificateType(u.req.HttpsConf.CertType)),
		Http2Status:        utils.StringPtr(getSwitch(u.req.HttpsConf.HttpTwo)),
		TlsVersion:         utils.StringPtr(getTlsVersions(u.req.HttpsConf.TlsVersion)),
		OcspStaplingStatus: utils.StringPtr(getSwitch(u.req.HttpsConf.OcspStatus)),
	}
	u.configs.Hsts = &model.Hsts{
		Status:            getSwitch(u.req.HttpsConf.HstsStatus),
		MaxAge:            utils.Int32Ptr(int32(u.req.HttpsConf.HstsExpirationTime)),
		IncludeSubdomains: utils.StringPtr(getSwitch(u.req.HttpsConf.HstsSubdomain)),
	}
}

// DomainAccessDataStatic 域名访问数据统计信息
func (h *Huawei) DomainAccessDataStatic(req *types.DomainAccessDataStaticRequest) (types.DomainAccessDataStaticResponse, error) {
	if err := h.verifyMutualExclusion(req); err != nil {
		return nil, err
	}
	if req.Metric == consts.DataAccessMetricTypeHitFlux || req.Metric == consts.DataAccessMetricTypeHitRequest || (req.Isp == nil && req.District == nil && req.IpProtocol == nil && req.Protocol == nil) {
		return h.DomainServiceAreaDataStatic(req)
	}
	return h.DomainAccessLocationDataStatic(req)
}

// DomainAccessLocationDataStatic 域名地区运营商访问数据统计信息
func (h *Huawei) DomainAccessLocationDataStatic(req *types.DomainAccessDataStaticRequest) (types.DomainAccessDataStaticResponse, error) {
	request := &model.ShowDomainLocationStatsRequest{
		Action:     getAccessDataStaticType(consts.DataStaticTypeDetail),
		StartTime:  req.StartTime * 1000,
		EndTime:    req.EndTime * 1000,
		DomainName: strings.Join(req.Domains, ","),
		StatType:   getDataAccessMetricType(req.Metric),
		Interval:   utils.Int64Ptr(getDataIntervalType(req.Interval)),
		GroupBy:    utils.StringPtr("domain"),
	}
	isStatusCode := strings.Contains(request.StatType, "status_code")
	if req.Area == consts.AreaCodeChinaMainland {
		request.Country = utils.StringPtr("cn")
		if req.District != nil {
			request.Province = utils.StringPtr(getProvinceCode(*req.District))
		}
	}
	if req.Area == consts.AreaCodeOversea {
		if req.District != nil {
			request.Country = utils.StringPtr(getCountryCode(*req.District))
		} else {
			request.Country = utils.StringPtr(getAllCountryCode([]string{"cn"}))
		}
	}
	if req.Isp != nil {
		request.Isp = utils.StringPtr(getIspCode(*req.Isp))
		if *req.Isp == consts.IspCodeOther {
			request.Country = utils.StringPtr(getAllCountryCode([]string{"cn"}))
		}
	}
	if req.IpProtocol != nil {
		request.IpVersion = utils.StringPtr(getIpProtocol(*req.IpProtocol))
	}
	if req.Protocol != nil {
		request.Protocol = utils.StringPtr(getHttpProtocol(*req.Protocol))
	}
	response, err := h.client.ShowDomainLocationStats(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode != 200 {
		return nil, errors.New("show domain location stats error")
	}
	responseData := make(types.DomainAccessDataStaticResponse)
	timeStamps := utils.CalcTimeStampsWithInterval(req.StartTime, req.EndTime, *request.Interval) //获取所有时间戳
	if !isStatusCode {
		for _, domain := range req.Domains {
			responseData[domain] = h.makeDefaultData(timeStamps)
		}
	}
	if response.Result != nil {
		for resultDomain, v := range response.Result {
			jsonstr, _ := json.Marshal(v)
			jsonData := make(map[string][]int64)
			err = json.Unmarshal(jsonstr, &jsonData)
			if err != nil {
				return nil, err
			}
			//状态码处理
			if isStatusCode {
				for code, value := range jsonData {
					items := make([]*types.StaticData, 0, len(value))
					for metricIndex, metric := range value {
						if metricIndex <= len(timeStamps)-1 {
							items = append(items, &types.StaticData{
								Value: float64(metric),
								Time:  timeStamps[metricIndex],
							})
						}
					}
					responseData[fmt.Sprintf("%s#%s", resultDomain, code)] = items
				}
				continue
			}
			//其他指标处理
			if metrics, ok := jsonData[request.StatType]; ok {
				items := make([]*types.StaticData, 0, len(metrics))
				for metricIndex, metric := range metrics {
					if metricIndex <= len(timeStamps)-1 {
						items = append(items, &types.StaticData{
							Value: float64(metric),
							Time:  timeStamps[metricIndex],
						})
					}
				}
				responseData[resultDomain] = items
			}
		}
	}
	return responseData, nil
}

// DomainServiceAreaDataStatic 加速域名服务区域访问数据统计信息
func (h *Huawei) DomainServiceAreaDataStatic(req *types.DomainAccessDataStaticRequest) (types.DomainAccessDataStaticResponse, error) {
	stateRequest := &model.ShowDomainStatsRequest{
		Action:      getOriginDataStaticType(consts.DataStaticTypeDetail),
		StartTime:   req.StartTime * 1000,
		EndTime:     req.EndTime * 1000,
		DomainName:  strings.Join(req.Domains, ","),
		StatType:    getDataAccessMetricType(req.Metric),
		Interval:    utils.Int64Ptr(getDataIntervalType(req.Interval)),
		GroupBy:     utils.StringPtr("domain"),
		ServiceArea: utils.StringPtr(getAreaCode(req.Area).Value()),
	}
	isStatusCode := strings.Contains(stateRequest.StatType, "status_code")
	response, err := h.client.ShowDomainStats(stateRequest)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode != 200 {
		return nil, errors.New("show domain stats error")
	}
	responseData := make(types.DomainAccessDataStaticResponse)
	timeStamps := utils.CalcTimeStampsWithInterval(req.StartTime, req.EndTime, *stateRequest.Interval) //获取所有时间戳
	if !isStatusCode {
		for _, domain := range req.Domains {
			responseData[domain] = h.makeDefaultData(timeStamps)
		}
	}
	if response.Result != nil {
		for resultDomain, v := range response.Result {
			jsonstr, _ := json.Marshal(v)
			jsonData := make(map[string][]int64)
			err = json.Unmarshal(jsonstr, &jsonData)
			if err != nil {
				return nil, err
			}
			//状态码处理
			if isStatusCode {
				for code, value := range jsonData {
					items := make([]*types.StaticData, 0, len(value))
					for metricIndex, metric := range value {
						if metricIndex <= len(timeStamps)-1 {
							items = append(items, &types.StaticData{
								Value: float64(metric),
								Time:  timeStamps[metricIndex],
							})
						}
					}
					responseData[fmt.Sprintf("%s#%s", resultDomain, code)] = items
				}
				continue
			}
			//其他指标处理
			if metrics, ok := jsonData[stateRequest.StatType]; ok {
				items := make([]*types.StaticData, 0, len(metrics))
				for metricIndex, metric := range metrics {
					if metricIndex <= len(timeStamps)-1 {
						items = append(items, &types.StaticData{
							Value: float64(metric),
							Time:  timeStamps[metricIndex],
						})
					}
				}
				responseData[resultDomain] = items
			}
		}
	}
	return responseData, nil
}

// DomainOriginDataStatic 域名回源数据统计信息
func (h *Huawei) DomainOriginDataStatic(req *types.DomainOriginDataStaticRequest) (types.DomainOriginDataStaticResponse, error) {
	request := &model.ShowDomainStatsRequest{
		Action:     getOriginDataStaticType(consts.DataStaticTypeDetail),
		StartTime:  req.StartTime * 1000,
		EndTime:    req.EndTime * 1000,
		DomainName: strings.Join(req.Domains, ","),
		StatType:   getDataOriginMetricType(req.Metric),
		Interval:   utils.Int64Ptr(getDataIntervalType(req.Interval)),
		GroupBy:    utils.StringPtr("domain"),
	}
	isStatusCode := strings.Contains(request.StatType, "bs_status_code")
	response, err := h.client.ShowDomainStats(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode != 200 {
		return nil, errors.New("show domain location stats error")
	}
	responseData := make(types.DomainOriginDataStaticResponse)
	timeStamps := utils.CalcTimeStampsWithInterval(req.StartTime, req.EndTime, *request.Interval) //获取所有时间戳
	if !isStatusCode {
		for _, domain := range req.Domains {
			responseData[domain] = h.makeDefaultData(timeStamps)
		}
	}
	if response.Result != nil {
		for resultDomain, v := range response.Result {
			jsonstr, _ := json.Marshal(v)
			jsonData := make(map[string][]int64)
			err = json.Unmarshal(jsonstr, &jsonData)
			if err != nil {
				return nil, err
			}
			//状态码处理
			if isStatusCode {
				for code, value := range jsonData {
					items := make([]*types.StaticData, 0, len(value))
					for metricIndex, metric := range value {
						if metricIndex <= len(timeStamps)-1 {
							items = append(items, &types.StaticData{
								Value: float64(metric),
								Time:  timeStamps[metricIndex],
							})
						}
					}
					responseData[fmt.Sprintf("%s#%s", resultDomain, code)] = items
				}
				continue
			}
			//其他指标处理
			if metrics, ok := jsonData[request.StatType]; ok {
				items := make([]*types.StaticData, 0, len(metrics))
				for metricIndex, metric := range metrics {
					if metricIndex <= len(timeStamps)-1 {
						items = append(items, &types.StaticData{
							Value: float64(metric),
							Time:  timeStamps[metricIndex],
						})
					}
				}
				responseData[resultDomain] = items
			}
		}
	}
	return responseData, nil
}

// verifyMutualExclusion 验证互斥条件
func (h *Huawei) verifyMutualExclusion(req *types.DomainAccessDataStaticRequest) error {
	switch {
	case req.IpProtocol == nil:
		if req.District != nil && req.Isp != nil {
			return errors.New("when the ip protocol is empty, district and isp cannot have parameter mutual exclusion at the same time")
		}
	case req.IpProtocol != nil:
		if req.District != nil || req.Isp != nil {
			return errors.New("when the ip protocol is not empty, district and isp cannot have parameter mutual exclusion at the same time")
		}
	}
	return nil
}

// 生成回源和响应数据的默认数据
func (h *Huawei) makeDefaultData(timeStamps []int64) []*types.StaticData {
	items := make([]*types.StaticData, 0, len(timeStamps))
	for _, timeStamp := range timeStamps {
		items = append(items, &types.StaticData{
			Value: 0,
			Time:  timeStamp,
		})
	}
	return items
}

// ListTopUrlDataStatic 获取域名TOP URL访问数据
func (h *Huawei) ListTopUrlDataStatic(req *types.ListTopUrlDataStaticRequest) ([]*types.ListTopUrlDataStaticResponse, error) {
	request := &model.ShowTopUrlRequest{
		StartTime:   req.StartTime * 1000,
		EndTime:     req.EndTime * 1000,
		DomainName:  req.Domain,
		StatType:    getTopUrlFilter(req.Filter),
		ServiceArea: utils.StringPtr(getAreaCode(req.Area).Value()),
	}
	response, err := h.client.ShowTopUrl(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode != 200 {
		return nil, errors.New("show domain top url  error")
	}
	responseData := make([]*types.ListTopUrlDataStaticResponse, 0)
	if response.TopUrlSummary != nil {
		for _, v := range *response.TopUrlSummary {
			responseData = append(responseData, &types.ListTopUrlDataStaticResponse{
				Url:   *v.Url,
				Value: float64(*v.Value),
			})
		}
	}
	return responseData, nil
}

// DomainAccessTotalData 访问数据总数据
func (h *Huawei) DomainAccessTotalData(req *types.DomainAccessTotalDataRequest) (types.DataTotalDataResponse, error) {
	stateRequest := &model.ShowDomainStatsRequest{
		Action:      "summary",
		StartTime:   req.StartTime * 1000,
		EndTime:     req.EndTime * 1000,
		DomainName:  strings.Join(req.Domains, ","),
		StatType:    getDataAccessMetricType(req.Metric),
		Interval:    utils.Int64Ptr(300),
		GroupBy:     utils.StringPtr("domain"),
		ServiceArea: utils.StringPtr(getAreaCode(req.Area).Value()),
	}
	if strings.Contains(stateRequest.StatType, "status_code") {
		return nil, errors.New("status_code is not support")
	}
	response, err := h.client.ShowDomainStats(stateRequest)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode != 200 {
		return nil, errors.New("show domain stats error")
	}
	responseData := make(types.DataTotalDataResponse)
	if response.Result != nil {
		for resultDomain, v := range response.Result {
			jsonstr, _ := json.Marshal(v)
			jsonData := make(map[string]int64)
			err = json.Unmarshal(jsonstr, &jsonData)
			if err != nil {
				return nil, err
			}
			//其他指标处理
			if flux, ok := jsonData[stateRequest.StatType]; ok {
				responseData[resultDomain] = flux
			}
		}
	}
	return responseData, nil
}

// DomainOriginTotalData 回源总流量
func (h *Huawei) DomainOriginTotalData(req *types.DomainOriginTotalDataRequest) (types.DataTotalDataResponse, error) {
	request := &model.ShowDomainStatsRequest{
		Action:     "summary",
		StartTime:  req.StartTime * 1000,
		EndTime:    req.EndTime * 1000,
		DomainName: strings.Join(req.Domains, ","),
		StatType:   getDataOriginMetricType(req.Metric),
		Interval:   utils.Int64Ptr(300),
		GroupBy:    utils.StringPtr("domain"),
	}
	if strings.Contains(request.StatType, "bs_status_code") {
		return nil, errors.New("bs_status_code is not support")
	}
	response, err := h.client.ShowDomainStats(request)
	if err != nil {
		return nil, err
	}
	if response.HttpStatusCode != 200 {
		return nil, errors.New("show domain location stats error")
	}
	responseData := make(types.DataTotalDataResponse)
	if response.Result != nil {
		for resultDomain, v := range response.Result {
			jsonstr, _ := json.Marshal(v)
			jsonData := make(map[string]int64)
			err = json.Unmarshal(jsonstr, &jsonData)
			if err != nil {
				return nil, err
			}
			//其他指标处理
			if flux, ok := jsonData[request.StatType]; ok {
				responseData[resultDomain] = flux
			}
		}
	}
	return responseData, nil
}

func (h *Huawei) UserAccessRegionDistribution(req *types.UserAccessRegionDistributionRequest) (types.UserAccessRegionDistributionResponse, error) {
	request := &model.ShowDomainLocationStatsRequest{
		Action:     getAccessDataStaticType(consts.DataStaticTypeSum),
		StartTime:  req.StartTime * 1000,
		EndTime:    req.EndTime * 1000,
		DomainName: strings.Join(req.Domains, ","),
		StatType:   getDataAccessMetricType(req.Metric),
		Interval:   utils.Int64Ptr(300),
		GroupBy:    utils.StringPtr("domain,country"),
		Country:    utils.StringPtr("all"),
	}
	response, err := h.client.ShowDomainLocationStats(request)
	responseData := make(types.UserAccessRegionDistributionResponse)
	if err != nil {
		return responseData, err
	}
	if response.HttpStatusCode != 200 {
		return responseData, errors.New("show domain location stats error")
	}
	if response.Result != nil {
		for domain, v := range response.Result {
			jsonstr, _ := json.Marshal(v)
			jsonData := make(map[string]map[string]int64)
			err = json.Unmarshal(jsonstr, &jsonData)
			if err != nil {
				return responseData, err
			}
			responseData[domain] = &types.RegionDistribution{}
			for resultCountry, metric := range jsonData {
				if value, ok := metric[request.StatType]; ok {
					if resultCountry == "cn" {
						responseData[domain].MainLandValue += value
						continue
					}
					responseData[domain].OverSeaValue += value
				}
			}
		}
	}
	return responseData, nil
}
