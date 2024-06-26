package model

import "encoding/json"

type CacheKeyRequest struct {
	DomainName    string          `path:"domain-name" json:"-"`
	CacheKeyRules []*CacheKeyRule `json:"cacheKeyRules"`
}

type CacheKeyRule struct {
	//url匹配模式，支持正则，如果是全部匹配，入参可以配置为：.*
	PathPattern string `json:"pathPattern"`
	//指定具体的uri，如/test/specifyurl
	SpecifyUrl string `json:"specifyUrl"`
	//是否完全匹配specifyUrl，可选择为true和false。
	//为true则完全匹配；为false则模糊匹配，如指定/test/uri，请求/test/uri?p=1也会匹配
	FullMatch4SpecifyUrl bool `json:"fullMatch4SpecifyUrl"`
	//指定常用类型：选择缓存域名的是全部文件还是首页。入参参考值： all：全部文件 homepage：首页
	CustomPattern string `json:"customPattern"`
	//文件类型：指定需要缓存的文件类型。 文件类型包括：gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf 如果需要全部类型，则直接传all。多个以分号隔开，all和具体文件类型不能同时配置。
	FileType string `json:"fileType"`
	//自定义文件类型：在指定文件类型外根据自身需求，填写适当的可识别文件类型。可以搭配file-type使用。如果file-type也有配置，实际生效的文件类型是两个入参的总和
	CustomFileType string `json:"customFileType"`
	//目录：指定目录缓存。 输入合法的目录格式。多个以英文分号隔开
	Directory string `json:"directory"`
	//是否忽略大小写：允许值为true和false，默认为忽略
	IgnoreCase bool `json:"ignoreCase"`
	//头部名称
	//例如：指定头部lang，lang的值一致则缓存一份
	HeaderName string `json:"headerName"`
	//头部值的参数名，
	//例如：指定头部Cookie，头部值的参数名为name。则name的值一致则缓存一份。
	ParameterOfHeader string `json:"parameterOfHeader"`
	//优先级，表示客户多组配置的优先执行顺序。数字越大，优先级越高。不传默认为10，不可清空。
	Priority *int `json:"priority"`
}

type CacheKeyResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (o *CacheKeyRequest) Marshal() []byte {
	b, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return b
}

func (o *CacheKeyResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
