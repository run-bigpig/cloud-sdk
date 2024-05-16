module github.com/run-bigpig/cloud-sdk

go 1.21

require (
	github.com/huaweicloud/huaweicloud-sdk-go-v3 v0.1.95
	github.com/spf13/cast v1.6.0
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn v1.0.920
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.920
)

require (
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	go.mongodb.org/mongo-driver v1.12.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/huaweicloud/huaweicloud-sdk-go-v3 => github.com/run-bigpig/huaweicloud-sdk-go-v3 v0.0.0-20240515082902-ea132da5ad6f
