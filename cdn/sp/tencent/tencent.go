package tencent

import (
	"context"
	"errors"
	"fmt"
	"github.com/run-bigpig/cloud-sdk/cdn/consts"
	"github.com/run-bigpig/cloud-sdk/cdn/entity"
	"github.com/run-bigpig/cloud-sdk/utils"
	tencentsdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	errors2 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type Tencent struct {
	config *Config
	client *tencentsdk.Client
	ctx    context.Context
}

type Config struct {
	Region   string
	Endpoint string
	Ak       string
	Sk       string
}

// NewTencentSdkClient 创建腾讯云CDN客户端
func NewTencentSdkClient(ctx context.Context, conf *Config) *Tencent {
	auth := common.NewCredential(conf.Ak, conf.Sk)
	cfp := profile.NewClientProfile()
	cfp.HttpProfile.Endpoint = conf.Endpoint
	client, _ := tencentsdk.NewClient(auth, "", cfp)
	return &Tencent{
		config: conf,
		client: client,
		ctx:    ctx,
	}
}

func (t *Tencent) GetSdkName() string {
	return types.TencentSdkName
}

func (t *Tencent) IcpVerify(req *types.IcpVerifyRequest) bool {
	request := tencentsdk.NewAddCdnDomainRequest()

	request.Domain = utils.StringPtr(req.Domain)
	request.ServiceType = utils.StringPtr("web")
	request.Origin = &tencentsdk.Origin{
		OriginType: utils.StringPtr("domain"),
	}
	request.Area = utils.StringPtr("mainland")
	_, err := t.client.AddCdnDomain(request)
	if err != nil {
		var e *errors2.TencentCloudSDKError
		if errors.As(err, &e) {
			return !(e.GetCode() == tencentsdk.RESOURCEUNAVAILABLE_CDNHOSTNOICP)
		}
	}
	return true
}

// CreateDomain 创建域名
func (t *Tencent) CreateDomain(req *types.CreateDomainRequest) error {
	var originHost string
	if req == nil {
		return errors.New("request is nil")
	}
	//检测域名是否已存在
	detail, err := t.ShowDomainDetail(&types.ShowDomainDetailRequest{Domain: req.Domain})
	if err == nil {
		if detail.Status == consts.CdnDomainStatusStoped {
			//启用域名
			err = t.EnableDomain(&types.EnableDomainRequest{Domain: req.Domain})
			if err != nil {
				return err
			}
		}
		//切换区域
		err = t.changeDomainArea(req.Domain, detail.AreaCode, req.AreaCode)
		if err != nil {
			return err
		}
		return nil
	}
	request := tencentsdk.NewAddCdnDomainRequest()
	if req.Sources == nil || len(req.Sources) == 0 {
		return errors.New("sources is nil")
	}
	request.Domain = utils.StringPtr(req.Domain)
	request.ServiceType = utils.StringPtr(getChannelType(req.ChannelType))
	request.Area = utils.StringPtr(getAreaCode(req.AreaCode))
	origins, backOrigins, originsAddresses, backOriginsAddress := make([]*string, 0), make([]*string, 0), make([]string, 0), make([]string, 0)
	primaryNum, backNum := t.cacalOriginPrimaryOrBackNumber(req.Sources)
	for i, v := range req.Sources {
		if i == 0 {
			originHost = v.OriginHost
		}
		var port int64
		switch req.OriginProtocol {
		case consts.HttpProtocolHttp:
			port = v.OriginHttpPort
		case consts.HttpProtocolHttps:
			port = v.OriginHttpsPort
		case consts.OriginProtocolFollow:
			port = 0
		default:
			port = 0
		}
		if v.OriginPriority == consts.OriginPriorityPrimary {
			originAddress := getOriginAddressList(v.OriginAddressList, port, v.OriginWeight, primaryNum)
			origins = append(origins, utils.StringPtr(originAddress))
			originsAddresses = append(originsAddresses, v.OriginAddressList)
			continue
		}
		originAddress := getOriginAddressList(v.OriginAddressList, port, v.OriginWeight, backNum)
		backOrigins = append(backOrigins, utils.StringPtr(originAddress))
		backOriginsAddress = append(backOriginsAddress, v.OriginAddressList)
	}
	origin := &tencentsdk.Origin{
		Origins:            origins,
		ServerName:         utils.StringPtr(originHost),
		OriginType:         utils.StringPtr(getOriginTypeByAddresses(originsAddresses)),
		OriginPullProtocol: utils.StringPtr(getOriginProtocol(req.OriginProtocol)),
	}
	if len(backOrigins) > 0 {
		origin.BackupOrigins = backOrigins
		origin.BackupOriginType = utils.StringPtr(getOriginTypeByAddresses(backOriginsAddress))
		origin.BackupServerName = utils.StringPtr(originHost)
	}
	request.Origin = origin
	_, err = t.client.AddCdnDomain(request)
	if err != nil {
		return err
	}
	return nil
}

// changeDomainArea 切换区域
func (t *Tencent) changeDomainArea(domain string, lastAreaCode, nowAreaCode int64) error {
	if lastAreaCode == nowAreaCode {
		return nil
	}
	//更新到全球
	if lastAreaCode != consts.AreaCodeGlobal {
		err := t.UpdateDomain(&types.UpdateDomainRequest{
			Domain:       domain,
			UpdateAction: types.UpdateArea,
			CdnDomain: &entity.UpdateCdnDomainBaseConf{
				AreaCode: consts.AreaCodeGlobal,
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
			err := t.UpdateDomain(&types.UpdateDomainRequest{
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

// CacalOriginPrimaryOrBackNumber 计算主备个数
func (t *Tencent) cacalOriginPrimaryOrBackNumber(resource []*entity.OriginServerConf) (int, int) {
	var primary, back int
	for _, v := range resource {
		if v.OriginPriority == consts.OriginPriorityPrimary {
			primary++
			continue
		}
		back++
	}
	return primary, back
}

// UpdateDomain 更新域名
func (t *Tencent) UpdateDomain(req *types.UpdateDomainRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	request := tencentsdk.NewUpdateDomainConfigRequest()
	updateDomain := t.newUpdateDomainConfigModel(req, request)
	updateDomain.WithDomain()
	switch req.UpdateAction {
	case types.UpdateBaseConf:
		updateDomain.WithBaseConf()
	case types.UpdateArea:
		updateDomain.WithArea()
	case types.UpdateOriginConf:
		updateDomain.WithOrigins()
		updateDomain.WithOriginConf()
	case types.UpdateOriginServerConf:
		updateDomain.WithOrigins()
	case types.UpdateOriginUrlConf:
		updateDomain.WithOrigins()
		updateDomain.WithOriginUrlConf()
	case types.UpdateOriginAdvanceServerConf:
		updateDomain.WithOrigins()
		updateDomain.WithOriginAdvanceConf()
	case types.UpdateOriginRequestHeaderConf:
		updateDomain.WithOriginRequestHeaderConf()
	case types.UpdateBrowserCacheConf:
		updateDomain.WithBrowserCacheConf()
	case types.UpdateCacheListConf:
		updateDomain.WithCacheListConf()
	case types.UpdateCacheCodeConf:
		updateDomain.WithCacheCodeConf()
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
	case types.UpdateSpeedConf:
		updateDomain.WithSpeedConf()
	case types.UpdateIpFrequencyConf:
		updateDomain.WithIpFrequencyConf()
	case types.UpdateHttpsConf:
		updateDomain.WithHttpsConf()
	case types.UpdateIntelligentCompressionConf:
		updateDomain.WithIntelligentCompressionConf()
	case types.UpdateResponseHeaderConf:
		updateDomain.WithResponseHeaderConf()
	case types.UpdateCustomErrorPageConf:
		updateDomain.WithCustomErrorPageConf()
	case types.UpdateRecommendConf:
		updateDomain.WithBrowserCacheConf()
		updateDomain.WithOriginConf()
		updateDomain.WithCacheListConf()
		updateDomain.WithHttpsConf()
		updateDomain.WithIntelligentCompressionConf()
		updateDomain.WithIpFrequencyConf()
	case types.UpdateFullConf:
		updateDomain.WithBaseConf()
		updateDomain.WithOriginConf()
		updateDomain.WithOrigins()
		updateDomain.WithOriginUrlConf()
		updateDomain.WithOriginAdvanceConf()
		updateDomain.WithOriginRequestHeaderConf()
		updateDomain.WithIpFilterConf()
		updateDomain.WithIpFrequencyConf()
		updateDomain.WithRefererConf()
		updateDomain.WithUserAgentConf()
		updateDomain.WithSpeedConf()
		updateDomain.WithAuthConf()
		updateDomain.WithRemoteAuthConf()
		updateDomain.WithCacheListConf()
		updateDomain.WithCacheCodeConf()
		updateDomain.WithRequestUrlRewriteConf()
		updateDomain.WithBrowserCacheConf()
		updateDomain.WithCustomErrorPageConf()
		updateDomain.WithIntelligentCompressionConf()
		updateDomain.WithResponseHeaderConf()
		updateDomain.WithHttpsConf()
	}
	_, err := t.client.UpdateDomainConfig(request)
	if err != nil {
		return err
	}
	return nil
}

// DisableDomain 停用域名
func (t *Tencent) DisableDomain(req *types.DisableDomainRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	request := tencentsdk.NewStopCdnDomainRequest()

	request.Domain = utils.StringPtr(req.Domain)

	// 返回的resp是一个StopCdnDomainResponse的实例，与请求对象对应
	_, err := t.client.StopCdnDomain(request)
	if err != nil {
		return err
	}
	return nil
}

// EnableDomain 启用域名
func (t *Tencent) EnableDomain(req *types.EnableDomainRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	request := tencentsdk.NewStartCdnDomainRequest()
	request.Domain = utils.StringPtr(req.Domain)
	_, err := t.client.StartCdnDomain(request)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDomain 删除域名
func (t *Tencent) DeleteDomain(req *types.DeleteDomainRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	//检测域名是否已存在
	_, err := t.ShowDomainDetail(&types.ShowDomainDetailRequest{Domain: req.Domain})
	if err != nil {
		var e *errors2.TencentCloudSDKError
		if errors.As(err, &e) {
			if e.GetCode() == tencentsdk.RESOURCENOTFOUND {
				return nil
			}
		}
		return err
	}
	request := tencentsdk.NewDeleteCdnDomainRequest()
	request.Domain = utils.StringPtr(req.Domain)
	_, err = t.client.DeleteCdnDomain(request)
	if err != nil {
		return err
	}
	return nil
}

// PurgePathCache 刷新目录缓存
func (t *Tencent) PurgePathCache(req *types.PurgePathCacheRequest) (*types.PurgeCacheResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	flushTypeSlice := []string{"delete", "flush"}
	request := tencentsdk.NewPurgePathCacheRequest()
	request.Paths = utils.StringPtrs(req.Paths)
	request.FlushType = utils.StringPtr(flushTypeSlice[req.Mode])
	request.UrlEncode = utils.BoolPtr(req.UrlEncode)
	response, err := t.client.PurgePathCache(request)
	if err != nil {
		return nil, err
	}
	return &types.PurgeCacheResponse{TaskId: *response.Response.TaskId}, nil
}

// PurgeUrlsCache 刷新URL缓存
func (t *Tencent) PurgeUrlsCache(req *types.PurgeUrlsCacheRequest) (*types.PurgeCacheResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewPurgeUrlsCacheRequest()
	request.Urls = utils.StringPtrs(req.Urls)
	request.UrlEncode = utils.BoolPtr(req.UrlEncode)
	response, err := t.client.PurgeUrlsCache(request)
	if err != nil {
		return nil, err
	}
	return &types.PurgeCacheResponse{TaskId: *response.Response.TaskId}, nil
}

// PushUrlsCache 预热URL缓存
func (t *Tencent) PushUrlsCache(req *types.PushUrlsCacheRequest) (*types.PushUrlsCacheResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewPushUrlsCacheRequest()

	request.Urls = utils.StringPtrs(req.Urls)
	request.UrlEncode = utils.BoolPtr(req.UrlEncode)
	request.Area = utils.StringPtr("global")
	response, err := t.client.PushUrlsCache(request)
	if err != nil {
		return nil, err
	}
	return &types.PushUrlsCacheResponse{TaskId: *response.Response.TaskId}, nil
}

// ShowPurgeTaskStatus 展示刷新任务状态
func (t *Tencent) ShowPurgeTaskStatus(req *types.ShowPurgeTaskStatusRequest) (*types.ShowPurgeTaskStatusResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewDescribePurgeTasksRequest()
	request.TaskId = utils.StringPtr(req.TaskId)
	request.Offset = utils.Int64Ptr(0)
	request.Limit = utils.Int64Ptr(1)
	response, err := t.client.DescribePurgeTasks(request)
	if err != nil {
		return nil, err
	}
	if len(response.Response.PurgeLogs) == 0 {
		return nil, errors.New("task not found")
	}
	var taskStatus int64
	for _, v := range response.Response.PurgeLogs {
		taskStatus = getShowContentPurgeOrPushStatus(*v.Status)
		if taskStatus == consts.ShowContentPurgeOrPushStatusDoing {
			break
		}
	}
	return &types.ShowPurgeTaskStatusResponse{
		TaskId: req.TaskId,
		Status: taskStatus,
	}, nil
}

// ShowPurgeTaskList 展示刷新任务列表
func (t *Tencent) ShowPurgeTaskList(req *types.ShowPurgeTaskListRequest) (*types.ShowPurgeTaskListResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewDescribePurgeTasksRequest()
	//计算翻页
	offset, limit := utils.CalcOffsetAndLimit(req.Page, req.Limit)
	request.Offset = utils.Int64Ptr(offset)
	request.Limit = utils.Int64Ptr(limit)
	request.PurgeType = utils.StringPtr(getShowContentPurgeType(req.PurgeType))
	request.StartTime = utils.StringPtr(utils.FormatTimeWithTimezone(req.StartTime, req.TimeZone))
	request.EndTime = utils.StringPtr(utils.FormatTimeWithTimezone(req.EndTime, req.TimeZone))
	request.Status = utils.StringPtr(setShowContentPurgeOrPushStatus(req.TaskStatus))
	response, err := t.client.DescribePurgeTasks(request)
	if err != nil {
		return nil, err
	}
	total := response.Response.TotalCount
	if *total == 0 {
		return &types.ShowPurgeTaskListResponse{Total: *total, List: make([]string, 0)}, nil
	}
	tasks := make([]string, 0, len(response.Response.PurgeLogs))
	for _, v := range response.Response.PurgeLogs {
		tasks = append(tasks, *v.TaskId)
	}
	return &types.ShowPurgeTaskListResponse{Total: *total, List: tasks}, nil

}

// ShowPushTaskStatus 展示预热任务状态
func (t *Tencent) ShowPushTaskStatus(req *types.ShowPushTaskStatusRequest) (*types.ShowPushTaskStatusResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewDescribePushTasksRequest()
	request.TaskId = utils.StringPtr(req.TaskId)
	request.Offset = utils.Int64Ptr(0)
	request.Limit = utils.Int64Ptr(1)
	response, err := t.client.DescribePushTasks(request)
	if err != nil {
		return nil, err
	}
	if len(response.Response.PushLogs) == 0 {
		return nil, errors.New("task not found")
	}
	var taskStatus int64
	for _, v := range response.Response.PushLogs {
		taskStatus = getShowContentPurgeOrPushStatus(*v.Status)
		if taskStatus == consts.ShowContentPurgeOrPushStatusDoing {
			break
		}
	}
	return &types.ShowPushTaskStatusResponse{
		TaskId: req.TaskId,
		Status: taskStatus,
	}, nil
}

// ShowPushTaskList 展示预热任务列表
func (t *Tencent) ShowPushTaskList(req *types.ShowPushTaskListRequest) (*types.ShowPushTaskListResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewDescribePushTasksRequest()
	//计算翻页
	offset, limit := utils.CalcOffsetAndLimit(req.Page, req.Limit)
	request.Offset = utils.Int64Ptr(offset)
	request.Limit = utils.Int64Ptr(limit)
	request.StartTime = utils.StringPtr(utils.FormatTimeWithTimezone(req.StartTime, req.TimeZone))
	request.EndTime = utils.StringPtr(utils.FormatTimeWithTimezone(req.EndTime, req.TimeZone))
	request.Status = utils.StringPtr(setShowContentPurgeOrPushStatus(req.TaskStatus))
	response, err := t.client.DescribePushTasks(request)
	if err != nil {
		return nil, err
	}
	total := int64(*response.Response.TotalCount)
	if total == 0 {
		return &types.ShowPushTaskListResponse{Total: total, List: make([]string, 0)}, nil
	}
	tasks := make([]string, 0, len(response.Response.PushLogs))
	for _, v := range response.Response.PushLogs {
		tasks = append(tasks, *v.TaskId)
	}
	return &types.ShowPushTaskListResponse{Total: total, List: tasks}, nil

}

// ShowDomainDetail 获取域名详情
func (t *Tencent) ShowDomainDetail(req *types.ShowDomainDetailRequest) (*types.ShowDomainDetailResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewDescribeDomainsRequest()
	request.Offset = utils.Int64Ptr(0)
	request.Limit = utils.Int64Ptr(1)
	request.Filters = []*tencentsdk.DomainFilter{
		{
			Name:  utils.StringPtr("domain"),
			Value: utils.StringPtrs([]string{req.Domain}),
			Fuzzy: utils.BoolPtr(false),
		},
	}
	res, err := t.client.DescribeDomains(request)
	if err != nil {
		return nil, err
	}
	if len(res.Response.Domains) == 0 {
		return nil, errors2.NewTencentCloudSDKError(tencentsdk.RESOURCENOTFOUND, "domain not found", "")
	}
	domain := res.Response.Domains[0]
	return &types.ShowDomainDetailResponse{
		DomainId:   *domain.ResourceId,
		Domain:     *domain.Domain,
		Cname:      *domain.Cname,
		AreaCode:   mapAreaCode(*domain.Area),
		Status:     getDomainStatus(*domain.Status),
		CreateTime: utils.DateTimeToTimeStamp(*domain.CreateTime),
		UpdateTime: utils.DateTimeToTimeStamp(*domain.UpdateTime),
	}, nil
}

// ShowDomainStatusList 展示域名状态列表
func (t *Tencent) ShowDomainStatusList(req *types.ShowDomainStatusListRequest) (*types.ShowDomainStatusListResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewDescribeDomainsRequest()
	//计算翻页
	offset, limit := utils.CalcOffsetAndLimit(req.Page, req.Limit)
	request.Offset = utils.Int64Ptr(offset)
	request.Limit = utils.Int64Ptr(limit)
	request.Filters = []*tencentsdk.DomainFilter{
		{
			Name:  utils.StringPtr("status"),
			Value: utils.StringPtrs([]string{setDomainStatus(req.Status)}),
			Fuzzy: utils.BoolPtr(false),
		},
	}
	res, err := t.client.DescribeDomains(request)
	if err != nil {
		return nil, err
	}
	total := res.Response.TotalNumber
	if *total == 0 {
		return &types.ShowDomainStatusListResponse{Total: *total, List: make([]string, 0)}, nil
	}
	domains := make([]string, 0, len(res.Response.Domains))
	for _, v := range res.Response.Domains {
		domains = append(domains, *v.Domain)
	}
	return &types.ShowDomainStatusListResponse{Total: *total, List: domains}, nil
}

// CreateVerifyRecord 创建域名验证记录
func (t *Tencent) CreateVerifyRecord(req *types.CreateVerifyRecordRequest) (*types.CreateVerifyRecordResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewCreateVerifyRecordRequest()

	request.Domain = utils.StringPtr(req.Domain)
	response, err := t.client.CreateVerifyRecord(request)
	if err != nil {
		return nil, err
	}
	return &types.CreateVerifyRecordResponse{RecordCode: *response.Response.Record, FileVerifyUrl: *response.Response.FileVerifyUrl}, nil
}

// VerifyDomainRecord 验证域名解析记录
func (t *Tencent) VerifyDomainRecord(req *types.VerifyDomainRecordRequest) (*types.VerifyDomainRecordResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	request := tencentsdk.NewVerifyDomainRecordRequest()
	request.Domain = utils.StringPtr(req.Domain)
	request.VerifyType = utils.StringPtr(req.VerifyType)
	response, err := t.client.VerifyDomainRecord(request)
	if err != nil {
		return nil, err
	}
	return &types.VerifyDomainRecordResponse{Result: *response.Response.Result}, nil
}

//设置更新的各项信息

type UpdateDomainConfigModel struct {
	req  *types.UpdateDomainRequest
	udcr *tencentsdk.UpdateDomainConfigRequest
}

func (t *Tencent) newUpdateDomainConfigModel(req *types.UpdateDomainRequest, m *tencentsdk.UpdateDomainConfigRequest) *UpdateDomainConfigModel {
	return &UpdateDomainConfigModel{
		req:  req,
		udcr: m,
	}
}

func (u *UpdateDomainConfigModel) WithDomain() {
	u.udcr.Domain = utils.StringPtr(u.req.Domain)
}

func (u *UpdateDomainConfigModel) WithBaseConf() {
	//u.udcr.Area = utils.StringPtr(getAreaCode(u.req.CdnDomain.AreaCode))
	u.udcr.Ipv6Access = &tencentsdk.Ipv6Access{
		Switch: utils.StringPtr(getSwitch(u.req.CdnDomain.SupportIpv6)),
	}
}

func (u *UpdateDomainConfigModel) WithArea() {
	u.udcr.Area = utils.StringPtr(getAreaCode(u.req.CdnDomain.AreaCode))
}

func (u *UpdateDomainConfigModel) WithOrigins() {
	// 设置源站信息
	var (
		origin                               = &tencentsdk.Origin{}
		originHost                           string
		origins, backOrigins                 []*string
		originsAddresses, backOriginsAddress []string
		primary, back                        []*entity.OriginServerConf
	)
	if u.req.OriginServerConf != nil {
		for _, v := range u.req.OriginServerConf {
			if v.OriginPriority == consts.OriginPriorityPrimary {
				primary = append(primary, v)
				continue
			}
			back = append(back, v)
		}
	}
	originHost, origins, originsAddresses = dealOrigin(u.req.OriginConf.OriginProtocol, primary)
	origin.Origins = origins
	origin.OriginType = utils.StringPtr(getOriginTypeByAddresses(originsAddresses))
	origin.ServerName = utils.StringPtr(originHost)
	origin.OriginPullProtocol = utils.StringPtr(getOriginProtocol(u.req.OriginConf.OriginProtocol))
	if u.req.OriginConf.OriginSniSwitch != -1 {
		origin.Sni = &tencentsdk.OriginSni{
			Switch:     utils.StringPtr(getSwitch(u.req.OriginConf.OriginSniSwitch)),
			ServerName: utils.StringPtr(u.req.OriginConf.OriginSniValue),
		}
	}
	if len(back) > 0 {
		originHost, backOrigins, backOriginsAddress = dealOrigin(u.req.OriginConf.OriginProtocol, back)
		origin.BackupOrigins = backOrigins
		origin.BackupOriginType = utils.StringPtr(getOriginTypeByAddresses(backOriginsAddress))
		origin.BackupServerName = utils.StringPtr(originHost)
	}
	u.udcr.Origin = origin
}

func dealOrigin(protocol int64, originServerConfs []*entity.OriginServerConf) (originHost string, origins []*string, address []string) {
	for i, v := range originServerConfs {
		var port int64
		if i == 0 {
			originHost = v.OriginHost
		}
		switch protocol {
		case consts.HttpProtocolHttp:
			port = v.OriginHttpPort
		case consts.HttpProtocolHttps:
			port = v.OriginHttpsPort
		case consts.OriginProtocolFollow:
			port = 0
		}
		originAddress := getOriginAddressList(v.OriginAddressList, port, v.OriginWeight, len(originServerConfs))
		origins = append(origins, utils.StringPtr(originAddress))
		address = append(address, v.OriginAddressList)
	}
	return
}

func (u *UpdateDomainConfigModel) WithOriginConf() {
	connectTimeout := uint64(5)
	receiveTimeout := uint64(10)
	if u.req.OriginConf.TcpTimeout > 5 {
		connectTimeout = uint64(u.req.OriginConf.TcpTimeout)
	}
	if u.req.OriginConf.OriginTimeOut > 10 {
		receiveTimeout = uint64(u.req.OriginConf.OriginTimeOut)
	}
	u.udcr.OriginPullTimeout = &tencentsdk.OriginPullTimeout{
		ConnectTimeout: utils.Uint64Ptr(connectTimeout),
		ReceiveTimeout: utils.Uint64Ptr(receiveTimeout),
	}
	followRedirect := &tencentsdk.FollowRedirect{Switch: utils.StringPtr(getSwitch(u.req.OriginConf.OriginFollow))}
	u.udcr.FollowRedirect = followRedirect
	u.udcr.RangeOriginPull = &tencentsdk.RangeOriginPull{Switch: utils.StringPtr(getSwitch(u.req.OriginConf.OriginRange))}
}

func (u *UpdateDomainConfigModel) WithOriginUrlConf() {
	pathRules := make([]*tencentsdk.PathRule, 0, len(u.req.OriginUrlConf))
	if u.req.OriginUrlConf != nil && len(u.req.OriginUrlConf) > 0 {
		for _, v := range u.req.OriginUrlConf {
			pathRules = append(pathRules, &tencentsdk.PathRule{
				Regex:      utils.BoolPtr(getOriginMateMethodRegx(v.MateMethod)),
				Path:       utils.StringPtr(v.RewriteUrl),
				ForwardUri: utils.StringPtr(v.TargetUrl),
				FullMatch:  utils.BoolPtr(getOriginMateMethodAll(v.MateMethod))},
			)
		}
	}
	u.udcr.Origin.PathRules = pathRules
}

func (u *UpdateDomainConfigModel) WithOriginAdvanceConf() {
	pathBasedOriginRule := make([]*tencentsdk.PathBasedOriginRule, 0, len(u.req.OriginAdvanceServerConf))
	if u.req.OriginAdvanceServerConf != nil {
		for _, v := range u.req.OriginAdvanceServerConf {
			var port int64
			switch u.req.OriginConf.OriginProtocol {
			case consts.HttpProtocolHttp:
				port = v.OriginHttpPort
			case consts.HttpProtocolHttps:
				port = v.OriginHttpsPort
			case consts.OriginProtocolFollow:
				port = 0
			}
			pathBasedOriginRule = append(pathBasedOriginRule, &tencentsdk.PathBasedOriginRule{
				RuleType:  utils.StringPtr(getRuleType(v.UrlMatchMode)),
				RulePaths: utils.StringPtrs(v.UrlMatchRule),
				Origin:    utils.StringPtrs([]string{getOriginAddressList(v.OriginAddressList, port, 0, 1)}),
			})
		}
	}
	u.udcr.Origin.PathBasedOrigin = pathBasedOriginRule
}

func (u *UpdateDomainConfigModel) WithOriginRequestHeaderConf() {
	requestHeader := &tencentsdk.RequestHeader{}
	requestHeader.Switch = utils.StringPtr(consts.OFF)
	headerRules := make([]*tencentsdk.HttpHeaderPathRule, 0)
	if u.req.OriginRequestHeaderConf != nil && len(u.req.OriginRequestHeaderConf) > 0 {
		for _, v := range u.req.OriginRequestHeaderConf {
			headerRules = append(headerRules, &tencentsdk.HttpHeaderPathRule{
				HeaderMode:  utils.StringPtr(getOriginHeaderAction(v.Action)),
				HeaderName:  utils.StringPtr(v.ParameterKey),
				HeaderValue: utils.StringPtr(v.ParameterValue),
				RuleType:    utils.StringPtr("all"),
				RulePaths:   utils.StringPtrs([]string{"*"}),
			})
		}
		requestHeader.Switch = utils.StringPtr(consts.ON)
		requestHeader.HeaderRules = headerRules
	}
	u.udcr.RequestHeader = requestHeader
}

func (u *UpdateDomainConfigModel) WithBrowserCacheConf() {
	maxAge := &tencentsdk.MaxAge{}
	maxAge.Switch = utils.StringPtr(consts.OFF)
	if u.req.BrowserCacheConf != nil && len(u.req.BrowserCacheConf) > 0 {
		rules := make([]*tencentsdk.MaxAgeRule, 0)
		for _, v := range u.req.BrowserCacheConf {
			rules = append(rules, &tencentsdk.MaxAgeRule{
				MaxAgeType:     utils.StringPtr(getRuleType(v.CacheType)),
				MaxAgeContents: utils.StringPtrs(getRulePaths(v.CacheType, v.CacheContent)),
				MaxAgeTime:     utils.Int64Ptr(v.CacheTTL * getCacheUnit(v.CacheUnit)),
				FollowOrigin:   utils.StringPtr(isFollowOrigin(v.CacheStatus)),
			})
		}
		maxAge.Switch = utils.StringPtr(consts.ON)
		maxAge.MaxAgeRules = rules
	}
	u.udcr.MaxAge = maxAge
}

func (u *UpdateDomainConfigModel) WithCacheListConf() {
	ruleCacheRules := make([]*tencentsdk.RuleCache, 0)
	cache := &tencentsdk.Cache{
		RuleCache: ruleCacheRules,
	}
	if u.req.CacheListConf != nil && len(u.req.CacheListConf) > 0 {
		for _, v := range u.req.CacheListConf {
			cacheRule := &tencentsdk.RuleCache{
				RulePaths:   utils.StringPtrs(getRulePaths(v.CacheType, v.CacheContent)),
				RuleType:    utils.StringPtr(getRuleType(v.CacheType)),
				CacheConfig: &tencentsdk.RuleCacheConfig{},
			}
			switch v.CacheStatus {
			case consts.CacheStatusFollow:
				cacheRule.CacheConfig.FollowOrigin = &tencentsdk.CacheConfigFollowOrigin{
					Switch: utils.StringPtr(consts.ON),
				}
			case consts.CacheStatusOn:
				cacheRule.CacheConfig.Cache = &tencentsdk.CacheConfigCache{
					Switch:    utils.StringPtr(consts.ON),
					CacheTime: utils.Int64Ptr(v.CacheTTL * getCacheUnit(v.CacheUnit)),
				}
			case consts.CacheStatusOff:
				cacheRule.CacheConfig.NoCache = &tencentsdk.CacheConfigNoCache{
					Switch: utils.StringPtr(consts.ON),
				}
			}
			ruleCacheRules = append(ruleCacheRules, cacheRule)
		}
		cache.RuleCache = ruleCacheRules
	}
	u.udcr.Cache = cache
	keyRules := make([]*tencentsdk.KeyRule, 0, len(u.req.CacheListConf))
	cacheKey := &tencentsdk.CacheKey{
		KeyRules: keyRules,
	}
	if u.req.CacheListConf != nil && len(u.req.CacheListConf) > 0 {
		for _, v := range u.req.CacheListConf {
			var action string
			switch v.ParametersStatus {
			case consts.CacheParameterStatusInclude:
				action = "includeCustom"
			case consts.CacheParameterStatusExclude:
				action = "excludeCustom"
			}
			if v.CacheType == consts.RuleTypeAll {
				cacheKey.QueryString = &tencentsdk.QueryStringKey{
					Switch: utils.StringPtr(consts.OFF),
				}
				if v.ParametersStatus != consts.CacheParameterStatusOff && v.ParametersStatus != consts.CacheParameterStatusAll {
					cacheKey.QueryString.Switch = utils.StringPtr(consts.ON)
					cacheKey.QueryString.Action = utils.StringPtr(action)
					cacheKey.QueryString.Value = utils.StringPtr(strings.Join(v.ParametersValue, ";"))
				}
				cacheKey.FullUrlCache = utils.StringPtr(isCacheParameterFullPath(v.ParametersStatus))
				cacheKey.IgnoreCase = utils.StringPtr(getSwitch(v.Capitalization))
				continue
			}
			keyRule := &tencentsdk.KeyRule{
				RulePaths:    utils.StringPtrs(getRulePaths(v.CacheType, v.CacheContent)),
				RuleType:     utils.StringPtr(getRuleType(v.CacheType)),
				FullUrlCache: utils.StringPtr(isCacheParameterFullPath(v.ParametersStatus)),
				IgnoreCase:   utils.StringPtr(getSwitch(v.Capitalization)),
			}
			if v.ParametersStatus != consts.CacheParameterStatusOff && v.ParametersStatus != consts.CacheParameterStatusAll {
				keyRule.QueryString = &tencentsdk.RuleQueryString{
					Switch: utils.StringPtr(consts.ON),
					Action: utils.StringPtr(action),
					Value:  utils.StringPtr(strings.Join(v.ParametersValue, ";")),
				}
			}
			keyRules = append(keyRules, keyRule)
		}
		cacheKey.KeyRules = keyRules
	}
	u.udcr.CacheKey = cacheKey
}

func (u *UpdateDomainConfigModel) WithCacheCodeConf() {
	statusCodeCache := &tencentsdk.StatusCodeCache{}
	statusCodeCache.Switch = utils.StringPtr(consts.OFF)
	codeCacheRules := make([]*tencentsdk.StatusCodeCacheRule, 0)
	if u.req.CacheCodeConf != nil && len(u.req.CacheCodeConf) > 0 {
		for _, v := range u.req.CacheCodeConf {
			codeCacheRules = append(codeCacheRules, &tencentsdk.StatusCodeCacheRule{
				StatusCode: utils.StringPtr(cast.ToString(v.HttpCode)),
				CacheTime:  utils.Int64Ptr(v.CacheTTL * getCacheUnit(v.CacheUnit)),
			})
		}
		statusCodeCache.Switch = utils.StringPtr(consts.ON)
		statusCodeCache.CacheRules = codeCacheRules
	}
	u.udcr.StatusCodeCache = statusCodeCache
}

func (u *UpdateDomainConfigModel) WithRequestUrlRewriteConf() {
	requestUrlRewrite := &tencentsdk.UrlRedirect{}
	requestUrlRewrite.Switch = utils.StringPtr(consts.OFF)
	pathRules := make([]*tencentsdk.UrlRedirectRule, 0)
	if u.req.RequestUrlRewriteConf != nil && len(u.req.RequestUrlRewriteConf) > 0 {
		for _, v := range u.req.RequestUrlRewriteConf {
			pathRules = append(pathRules, &tencentsdk.UrlRedirectRule{
				RedirectStatusCode: utils.Int64Ptr(getRedirectCode(v.RedirectCode)),
				Pattern:            utils.StringPtr(v.RewriteUrl),
				RedirectUrl:        utils.StringPtr(v.TargetUrl),
				FullMatch:          utils.BoolPtr(isRequestUrlRewriteFullPath(v.MateMethod)),
			})
		}
		requestUrlRewrite.Switch = utils.StringPtr(consts.ON)
		requestUrlRewrite.PathRules = pathRules
	}
	u.udcr.UrlRedirect = requestUrlRewrite
}

func (u *UpdateDomainConfigModel) WithIpFilterConf() {
	filterRules := make([]*tencentsdk.IpFilterPathRule, 0)
	u.udcr.IpFilter = &tencentsdk.IpFilter{
		Switch: utils.StringPtr(getSwitch(u.req.IpFilterConf.Status)),
	}
	if u.req.IpFilterConf.IpFilterConf != nil && len(u.req.IpFilterConf.IpFilterConf) > 0 {
		for _, v := range u.req.IpFilterConf.IpFilterConf {
			if len(v.IpList) > 0 {
				filterRules = append(filterRules, &tencentsdk.IpFilterPathRule{
					FilterType: utils.StringPtr(getWhiteOrBlackList(v.IpType)),
					Filters:    utils.StringPtrs(v.IpList),
					RuleType:   utils.StringPtr(getAccessEffectiveType(v.EffectiveType)),
					RulePaths:  utils.StringPtrs(getAccessEffectiveContent(v.EffectiveType, v.EffectiveRules)),
				})
			}
		}
		if len(filterRules) > 0 {
			u.udcr.IpFilter.FilterRules = filterRules
		}
	}
}

func (u *UpdateDomainConfigModel) WithRefererConf() {
	refererRules := make([]*tencentsdk.RefererRule, 0)
	u.udcr.Referer = &tencentsdk.Referer{
		Switch: utils.StringPtr(getSwitch(u.req.RefererConf.Status)),
	}
	if u.req.RefererConf != nil && len(u.req.RefererConf.RefererList) > 0 {
		refererConf := u.req.RefererConf
		refererRules = append(refererRules, &tencentsdk.RefererRule{
			RuleType:    utils.StringPtr("all"),
			RulePaths:   utils.StringPtrs([]string{"*"}),
			RefererType: utils.StringPtr(getWhiteOrBlackList(refererConf.RefererType)),
			Referers:    utils.StringPtrs(refererConf.RefererList),
			AllowEmpty:  utils.BoolPtr(cast.ToBool(refererConf.IncludeEmpty)),
		})
		u.udcr.Referer.RefererRules = refererRules
	}
}

func (u *UpdateDomainConfigModel) WithUserAgentConf() {
	userAgentFilterRules := make([]*tencentsdk.UserAgentFilterRule, 0)
	u.udcr.UserAgentFilter = &tencentsdk.UserAgentFilter{
		Switch: utils.StringPtr(getSwitch(u.req.UserAgentConf.Status)),
	}
	if u.req.UserAgentConf.UserAgentConf != nil && len(u.req.UserAgentConf.UserAgentConf) > 0 {
		for _, v := range u.req.UserAgentConf.UserAgentConf {
			if len(v.AgentList) > 0 {
				userAgentFilterRules = append(userAgentFilterRules, &tencentsdk.UserAgentFilterRule{
					RuleType:   utils.StringPtr(getAccessEffectiveType(v.EffectiveType)),
					RulePaths:  utils.StringPtrs(getAccessEffectiveContent(v.EffectiveType, v.EffectiveRules)),
					UserAgents: utils.StringPtrs(v.AgentList),
					FilterType: utils.StringPtr(getWhiteOrBlackList(v.AgentType)),
				})
			}
		}
		if len(userAgentFilterRules) > 0 {
			u.udcr.UserAgentFilter.FilterRules = userAgentFilterRules
		}
	}
}

func (u *UpdateDomainConfigModel) WithAuthConf() {
	u.udcr.Authentication = &tencentsdk.Authentication{
		Switch:        utils.StringPtr(getSwitch(u.req.AuthConf.Status)),
		AuthAlgorithm: utils.StringPtr(getAccessAuthEncryptManner(u.req.AuthConf.EncryptMannger)),
	}
	if u.req.AuthConf.Status == types.ON {
		switch u.req.AuthConf.AuthManner {
		case consts.AccessAuthMannerTypeA:
			u.udcr.Authentication.TypeA = &tencentsdk.AuthenticationTypeA{
				SecretKey:       utils.StringPtr(u.req.AuthConf.AuthKey),
				SignParam:       utils.StringPtr(u.req.AuthConf.AuthParameter),
				ExpireTime:      utils.Int64Ptr(u.req.AuthConf.TimeValue),
				FileExtensions:  utils.StringPtrs(getAccessAuthRangeFileExtensions(u.req.AuthConf.FileSuffix)),
				FilterType:      utils.StringPtr(getAccessAuthRange(u.req.AuthConf.AuthRange)),
				BackupSecretKey: utils.StringPtr(u.req.AuthConf.AuthKeyBackup),
			}
		case consts.AccessAuthMannerTypeB:
			u.udcr.Authentication.TypeB = &tencentsdk.AuthenticationTypeB{
				SecretKey:       utils.StringPtr(u.req.AuthConf.AuthKey),
				ExpireTime:      utils.Int64Ptr(u.req.AuthConf.TimeValue),
				FileExtensions:  utils.StringPtrs(getAccessAuthRangeFileExtensions(u.req.AuthConf.FileSuffix)),
				FilterType:      utils.StringPtr(getAccessAuthRange(u.req.AuthConf.AuthRange)),
				BackupSecretKey: utils.StringPtr(u.req.AuthConf.AuthKeyBackup),
			}
		case consts.AccessAuthMannerTypeC:
			u.udcr.Authentication.TypeC = &tencentsdk.AuthenticationTypeC{
				SecretKey:       utils.StringPtr(u.req.AuthConf.AuthKey),
				ExpireTime:      utils.Int64Ptr(u.req.AuthConf.TimeValue),
				FileExtensions:  utils.StringPtrs(getAccessAuthRangeFileExtensions(u.req.AuthConf.FileSuffix)),
				FilterType:      utils.StringPtr(getAccessAuthRange(u.req.AuthConf.AuthRange)),
				TimeFormat:      utils.StringPtr(getAccessAuthTimeFormat(u.req.AuthConf.TimeFormat)),
				BackupSecretKey: utils.StringPtr(u.req.AuthConf.AuthKeyBackup),
			}
		case consts.AccessAuthMannerTypeD:
			u.udcr.Authentication.TypeD = &tencentsdk.AuthenticationTypeD{
				SecretKey:       utils.StringPtr(u.req.AuthConf.AuthKey),
				ExpireTime:      utils.Int64Ptr(u.req.AuthConf.TimeValue),
				FileExtensions:  utils.StringPtrs(getAccessAuthRangeFileExtensions(u.req.AuthConf.FileSuffix)),
				FilterType:      utils.StringPtr(getAccessAuthRange(u.req.AuthConf.AuthRange)),
				TimeParam:       utils.StringPtr("t"),
				SignParam:       utils.StringPtr(u.req.AuthConf.AuthParameter),
				TimeFormat:      utils.StringPtr(getAccessAuthTimeFormat(u.req.AuthConf.TimeFormat)),
				BackupSecretKey: utils.StringPtr(u.req.AuthConf.AuthKeyBackup),
			}
		}
	}
}

func (u *UpdateDomainConfigModel) WithRemoteAuthConf() {
	u.udcr.RemoteAuthentication = &tencentsdk.RemoteAuthentication{
		Switch: utils.StringPtr(getSwitch(u.req.RemoteAuthConf.Status)),
	}
	if u.req.RemoteAuthConf.Status == types.ON {
		u.udcr.RemoteAuthentication.RemoteAuthenticationRules = []*tencentsdk.RemoteAuthenticationRule{{
			Server:            utils.StringPtr(u.req.RemoteAuthConf.AuthUrl),
			AuthMethod:        utils.StringPtr(getRemoteAuthRequestMethod(u.req.RemoteAuthConf.ReqMethod)),
			RuleType:          utils.StringPtr(getRuleType(u.req.RemoteAuthConf.FileType)),
			RulePaths:         utils.StringPtrs(getRulePaths(u.req.RemoteAuthConf.FileType, u.req.RemoteAuthConf.FileContent)),
			AuthTimeout:       utils.Int64Ptr(u.req.RemoteAuthConf.TimeoutDuration),
			AuthTimeoutAction: utils.StringPtr(getAccessRemoteAuthTimeOutAction(u.req.RemoteAuthConf.TimeoutAction)),
		}}
	}
}

func (u *UpdateDomainConfigModel) WithSpeedConf() {
	dowstreamCapping := &tencentsdk.DownstreamCapping{}
	dowstreamCapping.Switch = utils.StringPtr(getSwitch(u.req.SpeedConf.Status))
	speedRules := make([]*tencentsdk.CappingRule, 0)
	if u.req.SpeedConf.SpeedConf != nil && len(u.req.SpeedConf.SpeedConf) > 0 {
		for _, v := range u.req.SpeedConf.SpeedConf {
			speedRules = append(speedRules, &tencentsdk.CappingRule{
				RuleType:      utils.StringPtr(getRuleType(v.RuleType)),
				RulePaths:     utils.StringPtrs(getRulePaths(v.RuleType, v.RuleContent)),
				KBpsThreshold: utils.Int64Ptr(v.SpeedValues),
			})
		}
		dowstreamCapping.Switch = utils.StringPtr(getSwitch(u.req.SpeedConf.Status))
	}
	u.udcr.DownstreamCapping = dowstreamCapping
}

func (u *UpdateDomainConfigModel) WithIpFrequencyConf() {
	u.udcr.IpFreqLimit = &tencentsdk.IpFreqLimit{
		Switch: utils.StringPtr(getSwitch(u.req.IpFrequencyConf.Status)),
	}
	if u.req.IpFrequencyConf.Frequency > 0 && u.req.IpFrequencyConf.Status == types.ON {
		u.udcr.IpFreqLimit.Qps = utils.Int64Ptr(u.req.IpFrequencyConf.Frequency)
	}
}

func (u *UpdateDomainConfigModel) WithHttpsConf() {
	https := &tencentsdk.Https{Switch: utils.StringPtr(getSwitch(u.req.HttpsConf.HttpsStatus))}
	u.udcr.HttpsBilling = &tencentsdk.HttpsBilling{Switch: utils.StringPtr(getSwitch(u.req.HttpsConf.HttpsStatus))}
	if u.req.HttpsConf.HttpsStatus == types.ON {
		https.Http2 = utils.StringPtr(getSwitch(u.req.HttpsConf.HttpTwo))
		https.OcspStapling = utils.StringPtr(getSwitch(u.req.HttpsConf.OcspStatus))
		https.CertInfo = &tencentsdk.ServerCert{
			Certificate: utils.StringPtr(u.req.HttpsConf.CertValue),
			PrivateKey:  utils.StringPtr(u.req.HttpsConf.CertKey),
			From:        utils.StringPtr("upload"),
		}
		https.Hsts = &tencentsdk.Hsts{
			Switch:            utils.StringPtr(getSwitch(u.req.HttpsConf.HstsStatus)),
			MaxAge:            utils.Int64Ptr(u.req.HttpsConf.HstsExpirationTime),
			IncludeSubDomains: utils.StringPtr(getSwitch(u.req.HttpsConf.HstsSubdomain)),
		}
		https.TlsVersion = utils.StringPtrs(getTlsVersions(u.req.HttpsConf.TlsVersion))
		u.udcr.ForceRedirect = &tencentsdk.ForceRedirect{
			Switch:             utils.StringPtr(getSwitch(u.req.HttpsConf.JumpForceStatus)),
			RedirectType:       utils.StringPtr(getHttpsJumpType(u.req.HttpsConf.JumpType)),
			RedirectStatusCode: utils.Int64Ptr(getRedirectCode(u.req.HttpsConf.JumpManner)),
		}
	}
	u.udcr.Https = https
}

func (u *UpdateDomainConfigModel) WithIntelligentCompressionConf() {
	compressionRules := make([]*tencentsdk.CompressionRule, 0)
	u.udcr.Compression = &tencentsdk.Compression{
		Switch: utils.StringPtr(getSwitch(u.req.IntelligentCompressionConf.Status)),
	}
	if u.req.IntelligentCompressionConf.IntelligentCompressionConf != nil && len(u.req.IntelligentCompressionConf.IntelligentCompressionConf) > 0 {
		for _, v := range u.req.IntelligentCompressionConf.IntelligentCompressionConf {
			compressionRules = append(compressionRules, &tencentsdk.CompressionRule{
				Compress:   utils.BoolPtr(true),
				MinLength:  utils.Int64Ptr(256),
				MaxLength:  utils.Int64Ptr(30 * 1024 * 1024),
				Algorithms: utils.StringPtrs([]string{getIntelligentCompressionCompressMethod(v.CompressMethod)}),
				RuleType:   utils.StringPtr(getCompressRuleType(v.CompressType)),
				RulePaths:  utils.StringPtrs(getRulePaths(v.CompressType, v.CompressContent)),
			})
			u.udcr.Compression.CompressionRules = compressionRules
		}
	}
}

func (u *UpdateDomainConfigModel) WithResponseHeaderConf() {
	responseHeader := &tencentsdk.ResponseHeader{}
	responseHeader.Switch = utils.StringPtr(consts.OFF)
	headerRules := make([]*tencentsdk.HttpHeaderPathRule, 0)
	if u.req.ResponseHeaderConf != nil && len(u.req.ResponseHeaderConf) > 0 {
		for _, v := range u.req.ResponseHeaderConf {
			headerRules = append(headerRules, &tencentsdk.HttpHeaderPathRule{
				HeaderMode:  utils.StringPtr(getOriginHeaderAction(v.Action)),
				HeaderName:  utils.StringPtr(v.ParameterKey),
				HeaderValue: utils.StringPtr(v.ParameterValue),
				RuleType:    utils.StringPtr("all"),
				RulePaths:   utils.StringPtrs([]string{"*"}),
			})
		}
		responseHeader.Switch = utils.StringPtr(consts.ON)
		responseHeader.HeaderRules = headerRules
	}
	u.udcr.ResponseHeader = responseHeader
}

func (u *UpdateDomainConfigModel) WithCustomErrorPageConf() {
	errorPage := &tencentsdk.ErrorPage{}
	errorPage.Switch = utils.StringPtr(consts.OFF)
	errorPageRules := make([]*tencentsdk.ErrorPageRule, 0)
	if u.req.CustomErrorPageConf != nil && len(u.req.CustomErrorPageConf) > 0 {
		for _, v := range u.req.CustomErrorPageConf {
			errorPageRules = append(errorPageRules, &tencentsdk.ErrorPageRule{
				StatusCode:   utils.Int64Ptr(v.StatusCode),
				RedirectCode: utils.Int64Ptr(getRedirectCode(v.RedirectCode)),
				RedirectUrl:  utils.StringPtr(v.GoalAddress),
			})
		}
		errorPage.Switch = utils.StringPtr(consts.ON)
		errorPage.PageRules = errorPageRules
	}
	u.udcr.ErrorPage = errorPage
}

// DomainAccessDataStatic 获取域名访问数据统计
func (t *Tencent) DomainAccessDataStatic(req *types.DomainAccessDataStaticRequest) (types.DomainAccessDataStaticResponse, error) {
	if err := t.verifyMutualExclusion(req); err != nil {
		return nil, err
	}
	request := tencentsdk.NewDescribeCdnDataRequest()
	if req.TimeZone == nil {
		req.TimeZone = common.StringPtr("Asia/Shanghai")
	}
	request.StartTime = common.StringPtr(utils.FormatTimeWithTimezone(req.StartTime, *req.TimeZone))
	request.EndTime = common.StringPtr(utils.FormatTimeWithTimezone(req.EndTime, *req.TimeZone))
	request.Metric = common.StringPtr(getDataAccessMetricType(req.Metric))
	request.Domains = common.StringPtrs(req.Domains)
	request.Interval = common.StringPtr(getDataIntervalType(req.Interval))
	request.Area = common.StringPtr(getAreaCode(req.Area))
	if req.Area == consts.AreaCodeOversea && req.AreaType == 1 {
		request.AreaType = common.StringPtr("client")
	}
	request.TimeZone = common.StringPtr(convertTimeZone(*req.TimeZone))
	request.Protocol = common.StringPtr("all")
	request.IpProtocol = common.StringPtr("all")
	request.Product = common.StringPtr(getProductType(req.Product))
	request.Detail = common.BoolPtr(true)
	isStatusCode := strings.Contains(*request.Metric, "xx")
	if req.Protocol != nil {
		request.Protocol = common.StringPtr(getHttpProtocol(*req.Protocol))
	}
	if req.IpProtocol != nil {
		request.IpProtocol = common.StringPtr(getIpProtocol(*req.IpProtocol))
	}
	if req.District != nil {
		if req.Area == consts.AreaCodeChinaMainland {
			request.District = common.Int64Ptr(getProvinceCode(*req.District))
		} else {
			request.District = common.Int64Ptr(getCountryCode(*req.District))
		}
	}
	if req.Isp != nil {
		request.Isp = common.Int64Ptr(getIspCode(*req.Isp))
	}
	response, err := t.client.DescribeCdnData(request)
	if err != nil {
		return nil, err
	}
	if len(response.Response.Data) == 0 {
		return nil, errors.New("data not found")
	}
	responseData := make(types.DomainAccessDataStaticResponse)
	for _, v := range response.Response.Data {
		for _, vv := range v.CdnData {
			//状态码处理
			if isStatusCode {
				if *vv.Metric == *request.Metric {
					continue
				} else {
					resource := strings.Join(req.Domains, "|")
					items := make([]*types.StaticData, 0, len(vv.DetailData))
					for _, vvv := range vv.DetailData {
						items = append(items, &types.StaticData{
							Value: *vvv.Value,
							Time:  utils.DateTimeToTimeStampWithTimezone(*vvv.Time, *req.TimeZone),
						})
					}
					responseData[fmt.Sprintf("%s#%s", resource, *vv.Metric)] = items
				}
				continue
			}
			//其他指标处理
			items := make([]*types.StaticData, 0, len(vv.DetailData))
			for _, vvv := range vv.DetailData {
				items = append(items, &types.StaticData{
					Value: *vvv.Value,
					Time:  utils.DateTimeToTimeStampWithTimezone(*vvv.Time, *req.TimeZone),
				})
			}
			responseData[*v.Resource] = items
		}
	}
	return responseData, nil
}

// DomainOriginDataStatic 获取域名回源数据统计
func (t *Tencent) DomainOriginDataStatic(req *types.DomainOriginDataStaticRequest) (types.DomainOriginDataStaticResponse, error) {
	request := tencentsdk.NewDescribeOriginDataRequest()
	if req.TimeZone == nil {
		req.TimeZone = common.StringPtr("Asia/Shanghai")
	}
	request.StartTime = common.StringPtr(utils.FormatTimeWithTimezone(req.StartTime, *req.TimeZone))
	request.EndTime = common.StringPtr(utils.FormatTimeWithTimezone(req.EndTime, *req.TimeZone))
	request.Metric = common.StringPtr(getDataOriginMetricType(req.Metric))
	request.Domains = common.StringPtrs(req.Domains)
	request.Interval = common.StringPtr(getDataIntervalType(req.Interval))
	request.Area = common.StringPtr(getAreaCode(req.Area))
	request.TimeZone = common.StringPtr(convertTimeZone(*req.TimeZone))
	request.Detail = common.BoolPtr(true)

	response, err := t.client.DescribeOriginData(request)
	if err != nil {
		return nil, err
	}
	isStatusCode := strings.Contains(*request.Metric, "xx")
	responseData := make(types.DomainOriginDataStaticResponse)
	for _, v := range response.Response.Data {
		for _, vv := range v.OriginData {
			//状态码处理
			if isStatusCode {
				if *vv.Metric == *request.Metric {
					continue
				} else {
					resource := strings.Join(req.Domains, "|")
					items := make([]*types.StaticData, 0, len(vv.DetailData))
					for _, vvv := range vv.DetailData {
						items = append(items, &types.StaticData{
							Value: *vvv.Value,
							Time:  utils.DateTimeToTimeStampWithTimezone(*vvv.Time, *req.TimeZone),
						})
					}
					responseData[fmt.Sprintf("%s#%s", resource, *vv.Metric)] = items
				}
				continue
			}
			//其他指标处理
			items := make([]*types.StaticData, 0, len(vv.DetailData))
			for _, vvv := range vv.DetailData {
				items = append(items, &types.StaticData{
					Value: *vvv.Value,
					Time:  utils.DateTimeToTimeStampWithTimezone(*vvv.Time, *req.TimeZone),
				})
			}
			responseData[*v.Resource] = items
		}
	}
	return responseData, nil
}

// verifyMutualExclusion 验证互斥条件
func (t *Tencent) verifyMutualExclusion(req *types.DomainAccessDataStaticRequest) error {
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

// ListTopUrlDataStatic 获取TOP URL数据统计
func (t *Tencent) ListTopUrlDataStatic(req *types.ListTopUrlDataStaticRequest) ([]*types.ListTopUrlDataStaticResponse, error) {
	request := tencentsdk.NewListTopDataRequest()
	request.StartTime = common.StringPtr(utils.FormatTimeWithTimezone(req.StartTime, "Asia/Shanghai"))
	request.EndTime = common.StringPtr(utils.FormatTimeWithTimezone(req.EndTime, "Asia/Shanghai"))
	request.Metric = common.StringPtr("url")
	request.Domains = common.StringPtrs([]string{req.Domain})
	request.Area = common.StringPtr(getAreaCode(req.Area))
	request.Filter = common.StringPtr(getTopUrlFilter(req.Filter))
	request.Detail = common.BoolPtr(false)
	request.Product = common.StringPtr(getProductType(req.Product))
	response, err := t.client.ListTopData(request)
	if err != nil {
		return nil, err
	}
	responseData := make([]*types.ListTopUrlDataStaticResponse, 0)
	for _, v := range response.Response.Data {
		for _, vv := range v.DetailData {
			responseData = append(responseData, &types.ListTopUrlDataStaticResponse{
				Url:   *vv.Name,
				Value: *vv.Value,
			})
		}
	}
	return responseData, nil
}

// DomainAccessTotalFlux 获取域名总流量
func (t *Tencent) DomainAccessTotalFlux(req *types.DomainAccessTotalFluxRequest) (types.DataTotalFluxResponse, error) {
	request := tencentsdk.NewDescribeCdnDataRequest()
	if req.TimeZone == nil {
		req.TimeZone = common.StringPtr("Asia/Shanghai")
	}
	request.StartTime = common.StringPtr(utils.FormatTimeWithTimezone(req.StartTime, *req.TimeZone))
	request.EndTime = common.StringPtr(utils.FormatTimeWithTimezone(req.EndTime, *req.TimeZone))
	request.Metric = common.StringPtr("flux")
	request.Domains = common.StringPtrs(req.Domains)
	request.Interval = common.StringPtr("5min")
	request.Area = common.StringPtr(getAreaCode(req.Area))
	request.TimeZone = common.StringPtr(convertTimeZone(*req.TimeZone))
	request.Protocol = common.StringPtr("all")
	request.IpProtocol = common.StringPtr("all")
	request.Product = common.StringPtr(getProductType(req.Product))
	request.Detail = common.BoolPtr(true)
	response, err := t.client.DescribeCdnData(request)
	if err != nil {
		return nil, err
	}
	if len(response.Response.Data) == 0 {
		return nil, errors.New("data not found")
	}
	responseData := make(types.DataTotalFluxResponse)
	for _, v := range response.Response.Data {
		for _, vv := range v.CdnData {
			if vv.Metric != nil && *vv.Metric == "flux" {
				responseData[*v.Resource] += int64(*vv.SummarizedData.Value)
			}
		}
	}
	return responseData, nil
}

// DomainOriginTotalFlux 获取域名回源总流量
func (t *Tencent) DomainOriginTotalFlux(req *types.DomainOriginTotalFluxRequest) (types.DataTotalFluxResponse, error) {
	request := tencentsdk.NewDescribeOriginDataRequest()
	if req.TimeZone == nil {
		req.TimeZone = common.StringPtr("Asia/Shanghai")
	}
	request.StartTime = common.StringPtr(utils.FormatTimeWithTimezone(req.StartTime, *req.TimeZone))
	request.EndTime = common.StringPtr(utils.FormatTimeWithTimezone(req.EndTime, *req.TimeZone))
	request.Metric = common.StringPtr("flux")
	request.Domains = common.StringPtrs(req.Domains)
	request.Interval = common.StringPtr("5min")
	request.Area = common.StringPtr(getAreaCode(req.Area))
	request.TimeZone = common.StringPtr(convertTimeZone(*req.TimeZone))
	request.Detail = common.BoolPtr(true)
	response, err := t.client.DescribeOriginData(request)
	if err != nil {
		return nil, err
	}
	responseData := make(types.DataTotalFluxResponse)
	for _, v := range response.Response.Data {
		for _, vv := range v.OriginData {
			if vv.Metric != nil && *vv.Metric == "flux" {
				responseData[*v.Resource] = int64(*vv.SummarizedData.Value)
			}
		}
	}
	return responseData, nil
}

// UserAccessRegionDistribution 获取用户访问区域分布
func (t *Tencent) UserAccessRegionDistribution(req *types.UserAccessRegionDistributionRequest) (types.UserAccessRegionDistributionResponse, error) {
	request := tencentsdk.NewListTopDataRequest()
	request.StartTime = common.StringPtr(utils.FormatTimeWithTimezone(req.StartTime, "Asia/Shanghai"))
	request.EndTime = common.StringPtr(utils.FormatTimeWithTimezone(req.EndTime, "Asia/Shanghai"))
	request.Metric = common.StringPtr("district")
	request.Domains = common.StringPtrs([]string{req.Domain})
	request.Area = common.StringPtr(getAreaCode(req.Area))
	request.Filter = common.StringPtr(getDataAccessMetricType(req.Metric))
	request.Detail = common.BoolPtr(false)
	request.Product = common.StringPtr(getProductType(req.Product))
	response, err := t.client.ListTopData(request)
	responseData := types.UserAccessRegionDistributionResponse{}
	if err != nil {
		return responseData, err
	}
	for _, v := range response.Response.Data {
		for _, vv := range v.DetailData {
			switch req.Area {
			case consts.AreaCodeChinaMainland:
				if *vv.Name == "-1" {
					responseData.OverSeaValue += int64(*vv.Value)
					continue
				}
				responseData.MainLandValue += int64(*vv.Value)
			case consts.AreaCodeOversea:
				if *vv.Name == "4460" {
					responseData.MainLandValue += int64(*vv.Value)
					continue
				}
				responseData.OverSeaValue += int64(*vv.Value)
			default:
				continue
			}
		}
	}
	return responseData, nil
}
