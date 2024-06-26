package util

import (
	"fmt"
	"github.com/run-bigpig/cloud-sdk/cdn/sp/wangsu/common/model"
	"io"
	"net/http"
	"strings"
)

func Call(requestMsg *model.HttpRequestMsg) ([]byte, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(requestMsg.Method, requestMsg.Url, strings.NewReader(requestMsg.Body))
	if err != nil {
		return nil, err
	}
	for key := range requestMsg.Headers {
		req.Header.Set(key, requestMsg.Headers[key])
	}
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
