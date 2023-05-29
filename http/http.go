package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

//reqUrl 请求地址
//header 请求头
//resData 返回数据
//statusCode HTTP返回Code
//errCallBackFunc 返回不为 statusCode 时处理函数
func HttpGet(reqUrl string, header map[string]string, resData interface{}, errCallBackFunc func(resp *http.Response) error) error {
	client := http.Client{
		Timeout: time.Second * 120,
	}
	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err
	}
	for k, v := range header {
		request.Header.Add(k, v)
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errCallBackFunc(resp)
	}
	if resData != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		err = json.Unmarshal(body, resData)
		if err != nil {
			return err
		}
	}
	return nil
}

//reqUrl 请求地址
//header 请求头
//reqData 请求数据
//resData 返回数据
//statusCode HTTP返回Code
//errCallBackFunc 返回不为 statusCode 时处理函数
func HttpPost(reqUrl string, header map[string]string, statusCode int, reqData interface{}, respData interface{}, errCallBackFunc func(resp *http.Response) error) error {
	client := new(http.Client)
	b, _ := json.Marshal(reqData)
	request, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	for k, v := range header {
		request.Header.Add(k, v)
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != statusCode {
		err = errCallBackFunc(resp)
		return err
	}
	if respData != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		err = json.Unmarshal(body, respData)
		if err != nil {
			return err
		}
	}
	return nil
}

//reqUrl 请求地址
//header 请求头
//reqData 请求数据
//resData 返回数据
//statusCode HTTP返回Code
//errCallBackFunc 返回不为 statusCode 时处理函数
func HttpPut(reqUrl string, header map[string]string, statusCode int, reqData interface{}, respData interface{}, errCallBackFunc func(resp *http.Response) error) error {
	client := new(http.Client)
	b, _ := json.Marshal(reqData)
	request, err := http.NewRequest("PUT", reqUrl, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	for k, v := range header {
		request.Header.Add(k, v)
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != statusCode {
		return errCallBackFunc(resp)
	}
	if respData != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		err = json.Unmarshal(body, respData)
		if err != nil {
			return err
		}
	}
	return nil
}
