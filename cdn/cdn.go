package cdn

import (
	"context"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/huawei"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/tencent"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu"
	"github.com/run-bigpig/cloud-sdk/cdn/types"
)

type Cdn interface {
	GetSdkName() string
	CreateDomain(req *types.CreateDomainRequest) error                                                                               // 创建域名
	UpdateDomain(req *types.UpdateDomainRequest) error                                                                               // 更新域名
	DisableDomain(req *types.DisableDomainRequest) error                                                                             // 停用域名
	EnableDomain(req *types.EnableDomainRequest) error                                                                               // 启用域名
	DeleteDomain(req *types.DeleteDomainRequest) error                                                                               // 删除域名
	CreateVerifyRecord(req *types.CreateVerifyRecordRequest) (*types.CreateVerifyRecordResponse, error)                              // 创建域名验证记录
	VerifyDomainRecord(req *types.VerifyDomainRecordRequest) (*types.VerifyDomainRecordResponse, error)                              // 验证域名
	ShowDomainDetail(req *types.ShowDomainDetailRequest) (*types.ShowDomainDetailResponse, error)                                    // 获取域名详情
	ShowDomainStatusList(req *types.ShowDomainStatusListRequest) (*types.ShowDomainStatusListResponse, error)                        // 获取指定状态域名列表
	PurgePathCache(req *types.PurgePathCacheRequest) (*types.PurgeCacheResponse, error)                                              // 刷新目录缓存
	PurgeUrlsCache(req *types.PurgeUrlsCacheRequest) (*types.PurgeCacheResponse, error)                                              // 刷新URL缓存
	PushUrlsCache(req *types.PushUrlsCacheRequest) (*types.PushUrlsCacheResponse, error)                                             // 预热URL缓存
	ShowPurgeTaskStatus(req *types.ShowPurgeTaskStatusRequest) (*types.ShowPurgeTaskStatusResponse, error)                           // 获取刷新任务状态
	ShowPushTaskStatus(req *types.ShowPushTaskStatusRequest) (*types.ShowPushTaskStatusResponse, error)                              // 获取预热任务状态
	ShowPurgeTaskList(req *types.ShowPurgeTaskListRequest) (*types.ShowPurgeTaskListResponse, error)                                 // 获取刷新任务列表
	ShowPushTaskList(req *types.ShowPushTaskListRequest) (*types.ShowPushTaskListResponse, error)                                    // 获取预热任务列表
	DomainAccessDataStatic(req *types.DomainAccessDataStaticRequest) (types.DomainAccessDataStaticResponse, error)                   // 域名访问数据统计信息
	DomainOriginDataStatic(req *types.DomainOriginDataStaticRequest) (types.DomainOriginDataStaticResponse, error)                   // 域名回源数据统计信息
	ListTopUrlDataStatic(req *types.ListTopUrlDataStaticRequest) ([]*types.ListTopUrlDataStaticResponse, error)                      // 获取TOP URL访问数据
	DomainAccessTotalData(req *types.DomainAccessTotalDataRequest) (types.DataTotalDataResponse, error)                              // 域名访问总流量
	DomainOriginTotalData(req *types.DomainOriginTotalDataRequest) (types.DataTotalDataResponse, error)                              // 域名回源数据总流量
	UserAccessRegionDistribution(req *types.UserAccessRegionDistributionRequest) (types.UserAccessRegionDistributionResponse, error) // 用户访问区域分布
}

type Config struct {
	Huawei  huawei.Config
	Tencent tencent.Config
	Wangsu  wangsu.Config
}

func NewCdn(ctx context.Context, config interface{}) Cdn {
	switch t := config.(type) {
	case huawei.Config:
		return huawei.NewHuaweiSdkClient(ctx, &t)
	case tencent.Config:
		return tencent.NewTencentSdkClient(ctx, &t)
	case wangsu.Config:
		return wangsu.NewWangsuSdkClient(ctx, &t)
	default:
		return nil
	}
}
