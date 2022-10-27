// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/jcbowen/jcbaseGo/helper"
	"github.com/thedevsaddam/gojsonq"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type AliOpenApiStruct struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type ConfigStruct struct {
	AliOpenApi AliOpenApiStruct `json:"AliOpenApiStruct"` // 阿里云openApi配置
	SubDomain  string           `json:"subDomain"`        // 子域名
	Type       string           `json:"type"`             // ip类型 A:ipv4, AAAA:ipv6, ALL:ipv4和ipv6
}

var Config = ConfigStruct{
	AliOpenApi: AliOpenApiStruct{
		AccessKeyId:     "阿里云AccessKey ID",
		AccessKeySecret: "阿里云AccessKey Secret",
	},
	SubDomain: "www.example.com",
	Type:      "A",
}

func init() {
	log.SetPrefix("[https://github.com/jcbowen/goAliDNS] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	isLog := flag.Bool("log", false, "是否打印日志")
	flag.Parse()
	if *isLog {
		// 定义一个以时间为文件名的日志文件
		fileName := "./data/log/" + time.Now().Format("2006-01/02") + ".log"
		exists, err := helper.DirExists(fileName, true, 0755)
		if err != nil {
			panic(err)
		}
		if !exists {
			panic("创建日志目录失败")
		}
		logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalln("打开日志文件异常")
		}
		log.SetOutput(logFile)
	}

}

func main() {
	// 执行主程序
	err := _main()
	if err != nil {
		panic(err)
	}
}

func _main() (_err error) {
	log.Println("开始获取配置信息")
	configFile := "./data/conf.json"
	configFileAbs, err := helper.GetAbsPath(configFile)
	if err != nil {
		panic(err)
	}
	if helper.FileExists(configFileAbs) {
		err = helper.ReadJsonFile(configFileAbs, &Config)
		if err != nil {
			panic(err)
		}
		log.Println("配置文件读取成功")
	} else {
		// 如果配置文件不存在，则创建配置文件
		file, _ := json.MarshalIndent(Config, "", " ")
		err := helper.CreateFile(configFileAbs, file, 0755, true)
		if err != nil {
			panic(err)
		}
		return errors.New("配置文件不存在，已创建默认配置文件，请修改配置文件后再次运行！\n配置文件路径：" + configFileAbs)
	}

	log.Println("正在检查配置文件...")

	// 检查是否配置了配置信息且配置信息是否正确
	err = checkConfig()
	if err != nil {
		return err
	}

	log.Println("配置文件检查通过")
	log.Println("需修改解析的二级域名:", Config.SubDomain)
	log.Println("解析类型:", Config.Type)
	var (
		currentIpv4 string
		currentIpv6 string
	)
	log.Println("--------------------")
	if Config.Type == "A" || Config.Type == "ALL" {
		log.Println("正在获取本机ipv4地址...")
		jsonString, err := getCurrenJsonIp(&currentIpv4, "A")
		if err != nil {
			log.Println("获取本机ipv4地址失败，错误信息：", err)
			log.Println("--------------------")
			goto label1
		}
		log.Println("获取本机json格式ipv4信息成功:", jsonString)
		log.Println("本机ipv4地址为:", currentIpv4)
		log.Println("--------------------")
	}

label1:
	if Config.Type == "AAAA" || Config.Type == "ALL" {
		log.Println("正在获取本机ipv6地址...")
		jsonString, err := getCurrenJsonIp(&currentIpv6, "AAAA")
		if err != nil {
			log.Println("获取本机ipv6地址失败，错误信息：", err)
			log.Println("--------------------")
			goto label2
		}
		log.Println("获取本机json格式ipv6信息成功:", jsonString)
		log.Println("本机ipv6地址为:", currentIpv6)
		log.Println("--------------------")
	}

label2:

	log.Println("正在获取阿里云解析记录列表...")

	// 创建Client
	client, _err := createClient(tea.String(Config.AliOpenApi.AccessKeyId), tea.String(Config.AliOpenApi.AccessKeySecret))
	if _err != nil {
		return _err
	}
	// 获取子域名解析记录列表
	result, _err := client.DescribeSubDomainRecords(&alidns20150109.DescribeSubDomainRecordsRequest{
		SubDomain: tea.String(Config.SubDomain),
	})
	if _err != nil {
		return _err
	}

	strBody := result.Body.String()
	log.Println("获取阿里云解析记录列表成功\n解析记录列表:\n", strBody)
	log.Println("--------------------")

	type recordStruct struct {
		Status     string `json:"Status"`
		Type       string `json:"Type"`
		Weight     int    `json:"Weight"`
		Value      string `json:"Value"`
		TTL        int    `json:"TTL"`
		Line       string `json:"Line"`
		RecordId   string `json:"RecordId"`
		RR         string `json:"RR"`
		DomainName string `json:"DomainName"`
		Locked     bool   `json:"Locked"`
	}

	type domainRecordsStruct struct {
		Record []recordStruct `json:"Record"`
	}

	type bodyStruct struct {
		TotalCount    int                 `json:"TotalCount"`
		PageSize      int                 `json:"PageSize"`
		RequestId     string              `json:"RequestId"`
		DomainRecords domainRecordsStruct `json:"DomainRecords"`
	}

	body := bodyStruct{}
	helper.JsonString(strBody).ToStruct(&body)

	// 开始检查是否需要修改ipv4解析记录
	if currentIpv4 != "" && (Config.Type == "A" || Config.Type == "ALL") {
		hasIpv4Change := false
		for _, v := range body.DomainRecords.Record {
			if v.Type == "A" && v.Value != currentIpv4 {
				log.Println("正在修改ipv4解析记录...")
				// 修改ipv4解析记录
				updateIpv4Result, _err := client.UpdateDomainRecord(&alidns20150109.UpdateDomainRecordRequest{
					RecordId: tea.String(v.RecordId),
					RR:       tea.String(v.RR),
					Type:     tea.String(v.Type),
					Value:    tea.String(currentIpv4),
				})
				if _err != nil {
					return _err
				}
				log.Println("修改ipv4解析记录成功\n返回信息为:\n", updateIpv4Result.Body.String())
				log.Println("--------------------")
				hasIpv4Change = true
				break
			}
		}
		if !hasIpv4Change {
			log.Println("ipv4解析记录未发生变化，无需修改")
			log.Println("--------------------")
		}
	}

	// 开始检查是否需要修改ipv6解析记录
	if currentIpv6 != "" && (Config.Type == "AAAA" || Config.Type == "ALL") {
		hasIpv6Change := false
		for _, v := range body.DomainRecords.Record {
			if v.Type == "AAAA" && v.Value != currentIpv6 {
				log.Println("正在修改ipv6解析记录...")
				// 修改ipv6解析记录
				updateIpv6Result, _err := client.UpdateDomainRecord(&alidns20150109.UpdateDomainRecordRequest{
					RecordId: tea.String(v.RecordId),
					RR:       tea.String(v.RR),
					Type:     tea.String(v.Type),
					Value:    tea.String(currentIpv6),
				})
				if _err != nil {
					return _err
				}
				log.Println("修改ipv6解析记录成功\n返回信息为:\n", updateIpv6Result.Body.String())
				log.Println("--------------------")
				hasIpv6Change = true
				break
			}
		}
		if !hasIpv6Change {
			log.Println("ipv6解析记录未发生变化，无需修改")
			log.Println("--------------------")
		}
	}

	return _err
}

// 获取json格式的ip信息
func getCurrenJsonIp(ip *string, ipType string) (string, error) {
	httpClient := http.Client{Timeout: time.Second * 5}

	var res *http.Response
	var err error

	if ipType == "A" {
		res, err = httpClient.Get("http://ipv4.test.ipv6.fastweb.it/ip/?callback=_jqjsp")
	} else {
		res, err = httpClient.Get("http://ipv6.test.ipv6.fastweb.it/ip/?callback=_jqjsp")
	}
	if err != nil {
		return "", err
	}
	robots, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = res.Body.Close()
	if err != nil {
		return "", err
	}

	strRobots := string(robots)

	//------ 新接口返回的json数据是包含在_jqjsp(和)之中的 ------/
	// 去除开头的字符_jqjsp(
	strRobots = strRobots[7:]
	// 去除结尾的字符)
	strRobots = strRobots[:len(strRobots)-2]

	*ip = gojsonq.New().FromString(strRobots).Find("ip").(string)
	return strRobots, nil
}

// createClient
/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func createClient(accessKeyId *string, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
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

// checkConfig 检查配置文件
func checkConfig() error {
	if Config.AliOpenApi.AccessKeyId == "阿里云AccessKey ID" || Config.AliOpenApi.AccessKeyId == "" {
		return errors.New("请在配置文件中配置阿里云AccessKey ID")
	}
	if Config.AliOpenApi.AccessKeySecret == "阿里云AccessKey Secret" || Config.AliOpenApi.AccessKeySecret == "" {
		return errors.New("请在配置文件中配置阿里云AccessKey Secret")
	}
	if Config.SubDomain == "www.example.com" || Config.SubDomain == "" {
		return errors.New("请在配置文件中配置需要进行动态域名解析的子域名")
	}
	if Config.Type == "" {
		return errors.New("请在配置文件中配置需要进行动态域名解析的类型")
	}
	if Config.Type != "A" && Config.Type != "AAAA" && Config.Type != "ALL" {
		return errors.New("请在配置文件中配置正确的域名解析的类型，仅支持'A'、'AAAA'和'ALL'")
	}
	return nil
}
