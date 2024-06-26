package wangsu

import (
	"context"
	"fmt"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common/auth"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common/constant"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common/model"
	"github.com/run-bigpig/cloud-sdk/cdn/types"
)

type Wangsu struct {
	ctx    context.Context
	config *Config
	client *common.Client
}

type Config struct {
	Ak       string
	Sk       string
	Endpoint string
}

func NewWangsuSdkClient(ctx context.Context, conf *Config) *Wangsu {
	return &Wangsu{
		ctx:    ctx,
		config: conf,
		client: common.NewClient(auth.NewAuth(conf.Ak, conf.Sk, conf.Endpoint)),
	}
}

func (w *Wangsu) GetSdkName() string {
	return types.WangsuSdkName
}

func (w *Wangsu) CreateDomain(data *types.CreateDomainRequest) error {
	req := &model.CreateDomainRequest{
		Version:                   constant.SDKVersion,
		DomainName:                data.Domain,
		ServiceType:               getServiceType(data.ChannelType),
		ServiceAreas:              getServiceAreas(data.AreaCode),
		Comment:                   "",
		CnameWithCustomizedPrefix: true,
		AccelerateNoChina:         false,
	}
	res, err := w.client.CreateDomain(req)
	if err != nil {
		return err
	}
	if res.HttpStatusCode != 202 {
		return fmt.Errorf("create domain fail code %d message %s ", res.Code, res.Message)
	}
	return nil
}

func (w *Wangsu) UpdateDomain(data *types.UpdateDomainRequest) error {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) DisableDomain(data *types.DisableDomainRequest) error {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) EnableDomain(data *types.EnableDomainRequest) error {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) DeleteDomain(data *types.DeleteDomainRequest) error {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) CreateVerifyRecord(data *types.CreateVerifyRecordRequest) (*types.CreateVerifyRecordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) VerifyDomainRecord(data *types.VerifyDomainRecordRequest) (*types.VerifyDomainRecordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) ShowDomainDetail(data *types.ShowDomainDetailRequest) (*types.ShowDomainDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) ShowDomainStatusList(data *types.ShowDomainStatusListRequest) (*types.ShowDomainStatusListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) PurgePathCache(data *types.PurgePathCacheRequest) (*types.PurgeCacheResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) PurgeUrlsCache(data *types.PurgeUrlsCacheRequest) (*types.PurgeCacheResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) PushUrlsCache(data *types.PushUrlsCacheRequest) (*types.PushUrlsCacheResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) ShowPurgeTaskStatus(data *types.ShowPurgeTaskStatusRequest) (*types.ShowPurgeTaskStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) ShowPushTaskStatus(data *types.ShowPushTaskStatusRequest) (*types.ShowPushTaskStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) ShowPurgeTaskList(data *types.ShowPurgeTaskListRequest) (*types.ShowPurgeTaskListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) ShowPushTaskList(data *types.ShowPushTaskListRequest) (*types.ShowPushTaskListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) DomainAccessDataStatic(data *types.DomainAccessDataStaticRequest) (types.DomainAccessDataStaticResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) DomainOriginDataStatic(data *types.DomainOriginDataStaticRequest) (types.DomainOriginDataStaticResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) ListTopUrlDataStatic(data *types.ListTopUrlDataStaticRequest) ([]*types.ListTopUrlDataStaticResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) DomainAccessTotalData(data *types.DomainAccessTotalDataRequest) (types.DataTotalDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) DomainOriginTotalData(data *types.DomainOriginTotalDataRequest) (types.DataTotalDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Wangsu) UserAccessRegionDistribution(data *types.UserAccessRegionDistributionRequest) (types.UserAccessRegionDistributionResponse, error) {
	//TODO implement me
	panic("implement me")
}
