// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"encoding/json"
	"fmt"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/thedevsaddam/gojsonq"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

// 获取json格式的ip信息
func getCurrenJsonIp() string {
	httpClient := http.Client{Timeout: time.Second * 5}
	res, err := httpClient.Get("https://ipv6.jsonip.com")
	if err != nil {
		res, err = httpClient.Get("https://jsonip.com")
		if err != nil {
			return ""
		}
	}
	robots, _err := ioutil.ReadAll(res.Body)
	err2 := res.Body.Close()
	if err2 != nil {
		return ""
	}
	if _err != nil {
		return ""
	}
	return string(robots)
}

// ParseIP
// 0: invalid ip
// 4: ipv4
// 6: ipv6
func ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}
	return nil, 0
}

// CreateClient
/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("alidns.cn-shenzhen.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

func _main(args []*string) (_err error) {
	file, fErr := ioutil.ReadFile("./config/config.json")
	if fErr != nil {
		return fErr
	}

	fileData := string(file)
	fmt.Println("获取配置信息成功")
	configGojsonq := gojsonq.New().FromString(fileData)
	accessKeyId := configGojsonq.Find("aliOpenApi.accessKeyId").(string)
	accessKeySecret := configGojsonq.Reset().Find("aliOpenApi.accessKeySecret").(string)
	subDomain := configGojsonq.Reset().Find("subDomain").(string)
	settingType := configGojsonq.Reset().Find("setting.type").(string)

	fmt.Println("二级域名:", subDomain)
	fmt.Println("解析类型:", settingType)

	client, _err := CreateClient(tea.String(accessKeyId), tea.String(accessKeySecret))
	if _err != nil {
		return _err
	}

	describeSubDomainRecordsRequest := &alidns20150109.DescribeSubDomainRecordsRequest{
		SubDomain: tea.String(subDomain),
	}

	// 获取子域名解析记录列表
	result, _err := client.DescribeSubDomainRecords(describeSubDomainRecordsRequest)
	if _err != nil {
		return _err
	}
	recordString := result.Body.DomainRecords.Record[0].String()
	type con struct {
		Status     string
		Type       string
		Weight     int
		Value      string
		TTL        int
		Line       string
		RecordId   string
		RR         string
		DomainName string
		Locked     bool
	}
	record := &con{}
	_err = json.Unmarshal([]byte(recordString), record)
	fmt.Println("查询子域名解析记录成功：")
	fmt.Println(recordString)

	// 获取当前IP
	var ipVersion int
	jsonString := getCurrenJsonIp()
	currenIp := gojsonq.New().FromString(jsonString).Find("ip").(string)
	_, ipVersion = ParseIP(currenIp)

	fmt.Println("当前主机的IP地址成功：", currenIp)
	//fmt.Println("当前解析的IP地址为：", record.Value)

	// 如果IP发生了变化
	if currenIp != "" && currenIp != record.Value {
		var Type = "AAAA"
		if ipVersion == 4 {
			Type = "AAA"
		}
		updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
			RecordId: tea.String(record.RecordId),
			RR:       tea.String(record.RR),
			Type:     tea.String(Type),
			Value:    tea.String(currenIp),
		}

		updateResult, __err := client.UpdateDomainRecord(updateDomainRecordRequest)
		if __err != nil {
			return __err
		}

		fmt.Println("更新解析成功，返回信息如下：")
		fmt.Println(updateResult.Body.String())
	} else {
		fmt.Println("IP未发生变化，不执行更新")
	}

	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
