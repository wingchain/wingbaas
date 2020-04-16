
package utils

import (
	"fmt"
	"bytes"
	"os"
	"os/exec"
	"encoding/json"
	"path/filepath"
	"github.com/wingbaas/platformsrv/logger" 
	"github.com/goinggo/mapstructure"
)

type Config interface {
	CfgInit(cfgFile string) error
	CfgPathInit() error
}

type BaasCfg struct {
	SrvAddr        			string		`json:"SrvAddr"`
	SrvPort        			string 		`json:"SrvPort"`
	ClusterCfgPath 			string 		`json:"ClusterCfgPath"`
	ClusterPkiBasePath      string 		`json:"ClusterPkiBasePath"`
	BlockNetCfgBasePath     string 		`json:"BlockNetCfgBasePath"`
	BlockChainVersionCfg    string		`json:"BlockChainVersionCfg"`
	NfsExternalAddr         string      `json:"NfsExternalAddr"`
	NfsInternalAddr         string		`json:"NfsInternalAddr"`
	NfsBasePath				string		`json:"NfsBasePath"`
	NfsLocalRootDir         string      `json:"NfsLocalRootDir"`
	keyStorePath			string		`json:"keyStorePath"`
}

var BAAS_CFG *BaasCfg = nil
var BLOCK_CFG_MAP map[string]interface{}

func GetProcessRunRoot() (string,error) {
	root,err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Errorf("GetProcessRunRoot: get process run root dir error,%v",err)
		return "",fmt.Errorf("GetProcessRunRoot: get process run root dir error,%v",err)
	}
	return root,nil
}

func ExecShell(s string) (string, error) {
    cmd := exec.Command("/bin/bash", "-c", s)
    var out bytes.Buffer
    cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "",fmt.Errorf("execShell: exec cmd=%s   err=%v",s,err)
	}
    return out.String(),err
}

func (cfg *BaasCfg) CfgInit(cfgFile string) error {
	bytes,err := LoadFile(cfgFile)
	if err !=nil {
		logger.Errorf("CfgInit: load baas config file error")
		return fmt.Errorf("%v", err)
	}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		logger.Errorf("CfgInit: unmarshal error")
		return fmt.Errorf("%v", err)
	}
	return nil 
}

func (cfg *BaasCfg) CfgPathInit() error {
	root,err := GetProcessRunRoot()
	if err != nil {
		logger.Errorf("CfgPathInit: %v",err)
		return fmt.Errorf("CfgPathInit: %v",err)
	}
	cfg.ClusterCfgPath = root + "/" + cfg.ClusterCfgPath
	cfg.ClusterPkiBasePath = root + "/" + cfg.ClusterPkiBasePath
	cfg.BlockNetCfgBasePath = root + "/" + cfg.BlockNetCfgBasePath
	cfg.BlockChainVersionCfg = root + "/" + cfg.BlockChainVersionCfg
	cfg.NfsLocalRootDir = root + "/" + cfg.NfsLocalRootDir
	cfg.keyStorePath = root + "/" + cfg.keyStorePath

	err = DirCheck(cfg.ClusterCfgPath)
	if err != nil {
		logger.Errorf("CfgPathInit: ClusterCfgPath init error")
		return fmt.Errorf("%v", err)
	}
	err = DirCheck(cfg.ClusterPkiBasePath)
	if err != nil {
		logger.Errorf("CfgPathInit: ClusterPkiBasePath init error")
		return fmt.Errorf("%v", err)
	}
	err = DirCheck(cfg.BlockNetCfgBasePath)
	if err != nil {
		logger.Errorf("CfgPathInit: BlockNetCfgBasePath init error")
		return fmt.Errorf("%v", err)
	}
	err = DirCheck(cfg.NfsLocalRootDir)
	if err != nil {
		logger.Errorf("CfgPathInit: NfsLocalRootDir init error")
		return fmt.Errorf("%v", err)
	}
	err = DirCheck(cfg.keyStorePath)
	if err != nil {
		logger.Errorf("CfgPathInit: keyStorePath init error")
		return fmt.Errorf("%v", err)
	}

	cmd := "sudo mount -t nfs -o resvport " + cfg.NfsExternalAddr + ":" + cfg.NfsBasePath + " " + cfg.NfsLocalRootDir
	//cmd := "mount -t nfs " + cfg.NfsInternalAddr + ":" + cfg.NfsBasePath + " " + cfg.NfsLocalRootDir
	_,err = ExecShell(cmd)
	if err != nil {
		logger.Errorf("CfgPathInit: mount nfs error, %v",err)
		return fmt.Errorf("CfgPathInit: mount nfs error, %v",err)
	}

	return nil 
}

func (cfg *BaasCfg) CfgBlockCfgInit() error {
	bytes,err := LoadFile(cfg.BlockChainVersionCfg) 
	if err != nil {
		logger.Errorf("CfgBlockCfgInit: load blockchain version config error")
		return fmt.Errorf("CfgBlockCfgInit: load blockchain version config error")
	}
	BLOCK_CFG_MAP = make(map[string]interface{})
	err = json.Unmarshal(bytes, &BLOCK_CFG_MAP)
	if err != nil {
		logger.Errorf("CfgBlockCfgInit: unmarshal to map error")
		return fmt.Errorf("%v", err)
	}
	return nil 
}

func GetBlockImage(blockType string,version string,imageKey string)(string,error) {
	obj := BLOCK_CFG_MAP[blockType]
	if obj == nil {
		logger.Errorf("GetBlockImage: not find block type=%s",blockType)
		return "",fmt.Errorf("GetBlockImage: not find block type=%s",blockType)
	}
	m1 := make(map[string]map[string]string)
	err := mapstructure.Decode(obj,&m1); 
	if err != nil {
		logger.Errorf("GetBlockImage: unmarshal block type=%s",blockType)
		return "",fmt.Errorf("GetBlockImage: unmarshal block type=%s",blockType)
	}
	image := m1[version][imageKey]
	if image == ""{
		logger.Errorf("GetBlockImage: not find image type=%s version=%s imageKey=%s",blockType,version,imageKey)
		return "",fmt.Errorf("GetBlockImage: not find image type=%s version=%s imageKey=%s",blockType,version,imageKey)
	}
	return image,nil
}


