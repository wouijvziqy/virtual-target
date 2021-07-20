package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"virtual-target/log"
	"virtual-target/model"
	"virtual-target/util"
)

// StartNoConfig 不使用配置文件的启动函数
func StartNoConfig(host string) {
	ports := generatePorts()
	responses := getAllResponses()
	cases := resolveResponses(responses)
	var index = ports[0]
	for _, v := range cases {
		for i := 0; i < len(ports); i++ {
			// 确保从上次成功端口下一个开始启动
			if i > 0 && ports[i] < index {
				continue
			}
			// 如果启动失败从下一个端口启动
			if !startListen(v, host, ports[i]) {
				continue
			} else {
				// 启动成功后设置下一个case在从下一个端口开始尝试
				index = ports[i+1]
				break
			}
		}
	}
}

func StartUseConfig(host string, inputConfig []model.TestCase) {
	responses := getAllResponses()
	modules := resolveResponses(responses)
	var index int
	for i := 0; i < len(modules); i++ {
		for j := 0; j < len(inputConfig); j++ {
			if inputConfig[j].Name == modules[i].Name {
				inputConfig[j].Response = modules[i].Response
				index++
				break
			}
		}
	}
	if index != len(modules) {
		log.Error("the number of configuration items does not match")
		log.Error("check your config file")
		os.Exit(-1)
	}
	for _, input := range inputConfig {
		if !startListen(input, host, input.Port) {
			log.Error(fmt.Sprintf("%s->port:%d error", input.Name, input.Port))
		}
	}
}

// 处理响应文件
func resolveResponses(responses []string) []model.TestCase {
	var res []model.TestCase
	for _, v := range responses {
		resp, _ := ioutil.ReadFile(v)
		t := model.TestCase{
			Name: strings.Split(
				strings.Split(v, "/")[2],
				".RESPONSE")[0],
			Response: resolveContentLength(string(resp)),
		}
		res = append(res, t)
	}
	return res
}

// 清除并重新计算响应的Content-Length头
func resolveContentLength(resp string) string {
	lineSep := util.GetLineSep()
	header := strings.Split(resp, lineSep+lineSep)[0]
	body := strings.Split(resp, lineSep+lineSep)[1]
	headers := strings.Split(header, lineSep)
	for i := 1; i < len(headers); i++ {
		if strings.TrimSpace(headers[i]) != "" {
			key := strings.Split(headers[i], ":")[0]
			if key == "Content-Length" {
				for j := i; j < len(headers)-1; j++ {
					headers[j] = headers[j+1]
				}
				headers[len(headers)-1] = ""
			}
		}
	}
	newContentLength := fmt.Sprintf("Content-Length: %d", +len(body))
	var newResp = bytes.Buffer{}
	newResp.WriteString(strings.TrimSpace(strings.Join(headers, lineSep)))
	newResp.WriteString(lineSep)
	newResp.WriteString(newContentLength)
	newResp.WriteString(lineSep + lineSep)
	newResp.WriteString(body)
	return newResp.String()
}

// 尝试监听端口
func startListen(testCase model.TestCase, host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	log.Info(fmt.Sprintf("start %s at port %d", testCase.Name, port))
	go doListen(l, testCase)
	return true
}

// 默认从1024开始生成端口
func generatePorts() []int {
	var ports []int
	for i := 1024; i < 65535; i++ {
		ports = append(ports, i)
	}
	return ports
}

// 获得modules目录下所有响应文件名
func getAllResponses() []string {
	var responses []string
	var dirname = "./modules"
	dir, _ := ioutil.ReadDir(dirname)
	for _, file := range dir {
		if !file.IsDir() {
			filename := dirname + "/" + file.Name()
			responses = append(responses, filename)
		}
	}
	return responses
}

// 监听
func doListen(l net.Listener, testCase model.TestCase) {
	for {
		conn, _ := l.Accept()
		_, _ = conn.Write([]byte(testCase.Response))
	}
}
