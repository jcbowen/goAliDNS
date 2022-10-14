// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"encoding/json"
	"flag"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/thedevsaddam/gojsonq"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// 获取json格式的ip信息
func getCurrenJsonIp() string {
	httpClient := http.Client{Timeout: time.Second * 5}
	//res, err := httpClient.Get("https://ipv6.jsonip.com")
	res, err := httpClient.Get("http://ipv6.test.ipv6.fastweb.it/ip/?callback=_jqjsp")
	if err != nil {
		//res, err = httpClient.Get("https://jsonip.com")
		res, err = httpClient.Get("http://ipv4.test.ipv6.fastweb.it/ip/?callback=_jqjsp")
		if err != nil {
			return ""
		}
	}
	robots, _err := io.ReadAll(res.Body)
	err2 := res.Body.Close()
	if err2 != nil {
		return ""
	}
	if _err != nil {
		return ""
	}
	strRobots := string(robots)

	//------ 新接口返回的json数据是包含在_jqjsp(和)之中的 ------/
	// 去除开头的字符_jqjsp(
	strRobots = strRobots[7:]
	// 去除结尾的字符)
	strRobots = strRobots[:len(strRobots)-2]

	return strRobots
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

func _main() (_err error) {
	file, fErr := os.ReadFile("./data/config.json")
	if fErr != nil {
		return fErr
	}

	fileData := string(file)
	log.Println("获取配置信息成功")
	configGojsonq := gojsonq.New().FromString(fileData)
	accessKeyId := configGojsonq.Find("aliOpenApi.accessKeyId").(string)
	accessKeySecret := configGojsonq.Reset().Find("aliOpenApi.accessKeySecret").(string)
	subDomain := configGojsonq.Reset().Find("subDomain").(string)
	settingType := configGojsonq.Reset().Find("setting.type").(string)

	log.Println("二级域名:", subDomain)
	log.Println("解析类型:", settingType)

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
	log.Println("查询子域名解析记录成功：\n", recordString)

	// 获取当前IP
	var ipVersion int
	jsonString := getCurrenJsonIp()
	log.Println("获取json格式ip信息成功：", jsonString)
	currenIp := gojsonq.New().FromString(jsonString).Find("ip").(string)
	_, ipVersion = ParseIP(currenIp)

	log.Println("获取当前主机的IP地址成功：", currenIp)

	// 如果IP发生了变化
	if currenIp != "" && currenIp != record.Value {
		log.Println("IP发生了变化，开始更新解析记录")

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
		log.Println("更新解析记录成功：\n", updateResult.Body.String(), "\n")
	} else {
		log.Println("IP地址未发生变化，无需更新解析", "\n")
	}

	return _err
}

func init() {
	log.SetPrefix("[https://github.com/jcbowen/goAliDNS] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	isLog := flag.Bool("log", false, "是否打印日志")
	flag.Parse()
	if *isLog {
		// 定义一个以时间为文件名的日志文件
		fileName := time.Now().Format("2006-01-02 15:02:01") + ".log"

		logFile, err := os.OpenFile("./data/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalln("打开日志文件异常")
		}
		log.SetOutput(logFile)
	}

}

func main() {
	log.Println("开始获取配置信息")

	// 执行主程序
	err := _main()
	if err != nil {
		panic(err)
	}
}
