package util

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"virtual-target/log"
	"virtual-target/model"
)

// GetLineSep 根据操作系统得到不同的换行符
func GetLineSep() string {
	sysType := runtime.GOOS
	var lineSep string
	if sysType == "windows" {
		lineSep = "\r\n"
	} else {
		lineSep = "\n"
	}
	return lineSep
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return false
}

func ResolveConfigFile(configFile string) []model.TestCase {
	var cases []model.TestCase
	if strings.TrimSpace(configFile) == "" {
		return cases
	}
	lineSep := GetLineSep()
	split := strings.Split(configFile, lineSep)
	if len(split) == 0 {
		return cases
	}
	for _, v := range split {
		if v == "" {
			continue
		}
		if strings.Contains(v, "=") {
			temp := strings.Split(v, "=")
			name := temp[0]
			port, _ := strconv.Atoi(strings.TrimSpace(temp[1]))
			cases = append(cases, model.TestCase{
				Name: name,
				Port: port,
			})
		} else {
			log.Error("check your config file format")
			os.Exit(-1)
		}
	}
	return cases
}
