package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

type CodeConfig struct {
	Name   string   `json:"name"`
	Path   string   `json:"path"`
	Stub   string   `json:"stub"`
	Params []string `json:"params"`
}

var list map[string]*CodeConfig
var dirName = ".go-gen"

func init() {
	list = make(map[string]*CodeConfig)
}

func GetCodeConfig(name string) *CodeConfig {
	if len(list) == 0 {
		initList()
	}

	res, ok := list[name]
	if !ok {
		log.Fatalf("code config %s not exist", name)
	}

	return res
}

func initStubs(dir string) {
	for _, config := range list {
		bt, err := os.ReadFile(path.Join(dir, dirName, config.Stub))
		if err != nil {
			log.Fatalf("Read stub config file failed: %v", err)
		}

		config.Stub = string(bt)
	}
}

func initList() {
	currentDir, err := os.Getwd()
	if err == nil {
		// 打印当前工作目录
		bt, err := os.ReadFile(path.Join(currentDir, dirName, "config.json"))
		if err == nil {
			err = json.Unmarshal(bt, &list)
			if err != nil {
				log.Fatalf("Error parsing config.json: %v", err)
			}

			initStubs(currentDir)

			return
		}
	}

	u, err := user.Current()
	if err == nil {

	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(u.HomeDir)
}
