package config

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

type CodeStub struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Stub     string   `json:"stub"`
	CodeStub string   `json:"code_stub"`
	Params   []string `json:"params"`
}

type Config struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Github      string               `json:"github"`
	Stubs       map[string]*CodeStub `json:"stubs"`
}

var conf *Config
var dirName = ".go-gen"
var configName = "config.json"

//go:embed .go-gen
var defaultConfigDir embed.FS

func init() {
	conf = &Config{}
}

func CreateConfigDir(dir string, force bool) error {
	configPath := path.Join(dir, dirName, configName)
	bt, err := os.ReadFile(configPath)
	if err == nil && !force {
		return errors.New(fmt.Sprintf("%s already exists", configPath))
	}

	defaultPath := path.Join(dir, dirName, configName)
	bt, _ = os.ReadFile(defaultPath)

	err = os.WriteFile(configPath, bt, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("create config.json failed, %v", err))
	}

	initDefaultConfig()

	for _, stub := range conf.Stubs {
		codeStub := stub.CodeStub
		err = os.WriteFile(path.Join(configPath, stub.Stub), []byte(codeStub), 0644)
		if err != nil {
			return errors.New(fmt.Sprintf("create stub failed, %v", err))
		}
	}

	return nil
}

func GetCodeConfig(name string) *CodeStub {
	if len(conf.Stubs) == 0 {
		initConfigs()
	}

	res, ok := conf.Stubs[name]
	if !ok {
		log.Fatalf("code config %s not exist", name)
	}

	return res
}

func initStubs(dir string) {
	for _, config := range conf.Stubs {
		bt, err := os.ReadFile(path.Join(dir, dirName, config.Stub))
		if err != nil {
			log.Fatalf("Read stub config file failed: %v", err)
		}

		config.CodeStub = string(bt)
	}
}

func initConfig(dir string) error {
	bt, err := os.ReadFile(path.Join(dir, dirName, configName))
	if err != nil {
		return err
	}

	err = json.Unmarshal(bt, &conf)
	if err != nil {
		log.Fatalf("Error parsing config.json: %v", err)
	}

	initStubs(dir)
	return nil
}

func initDefaultConfig() {
	bt, err := defaultConfigDir.ReadFile(path.Join(dirName, configName))
	err = json.Unmarshal(bt, &conf)
	if err != nil {
		log.Fatalf("Error parsing config.json: %v", err)
	}

	for _, config := range conf.Stubs {
		bt, err := defaultConfigDir.ReadFile(path.Join(dirName, config.Stub))
		if err != nil {
			log.Fatalf("Read stub config file failed: %v", err)
		}

		config.CodeStub = string(bt)
	}
}

func initConfigs() {
	currentDir, err := os.Getwd()
	if err == nil {
		err = initConfig(currentDir)
		if err == nil {
			return
		}
	}

	u, err := user.Current()
	if err == nil {
		err = initConfig(u.HomeDir)
		if err == nil {
			return
		}
	}

	initDefaultConfig()
}
