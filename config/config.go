package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var FabricClientConf = &FabricConfig{
	UserName: "Admin",
	UserOrg:  "Org1",
	// MulOrgs:     []string{"Org1"},
	ChannelID:   "mychannel",
	ConfigFile:  "first-network.yaml",
	ChaincodeID: "mycc",
}

type FabricConfig struct {
	UserName string
	UserOrg  string
	// MulOrgs     []string
	ChannelID   string
	ConfigFile  string
	ChaincodeID string
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func LoadFabircConfig(dir string) (*FabricConfig, error) {
	path := filepath.Join(dir, "configs/fabconfig.toml")
	filePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	config := new(FabricConfig)
	if CheckFileIsExist(filePath) { //文件存在
		if _, err := toml.DecodeFile(filePath, config); err != nil {
			return nil, err
		} else {
			FabricClientConf = config
		}
	} else {
		configBuf := new(bytes.Buffer)
		if err := toml.NewEncoder(configBuf).Encode(FabricClientConf); err != nil {
			return nil, err
		}
		err := ioutil.WriteFile(filePath, configBuf.Bytes(), 0666)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
