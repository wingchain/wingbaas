
package alliance

import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/utils"
)

type Alliance struct {
	Id			string  `json:"Id"`
	Name		string  `json:"Name,omitempty"`
	Describe	string  `json:"Describe,,omitempty"` 
	Creator		string  `json:"Creator"`
	CreateTime  string  `json:"CreateTime,omitempty"`
	JoinTime  string    `json:"JoinTime,omitempty"` 
} 

const (
	ALLIANCE_FILE	string = "alliances.json" 
)

func AppendAlliance(alliance Alliance) (string,error) {
	cfgPath := utils.BAAS_CFG.AlliancePath + ALLIANCE_FILE
	exsist,_ := utils.PathExists(cfgPath)
	var alliances []Alliance
	if exsist {
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&alliances) 
			if err != nil {
				logger.Errorf("AppendAlliance: unmarshal alliances error,%v", err)
				return "",fmt.Errorf("%v", err)
			}
		}else {
			logger.Errorf("AppendAlliance: load alliance list file error,%v", err)
			return "",fmt.Errorf("%v", err)
		}
	}
	for _,a := range alliances {
		if a.Name == alliance.Name {
			logger.Errorf("%s%s","AppendAlliance: alliance name already exsist ",a.Name)
			return "",fmt.Errorf("%s%s","AppendAlliance: alliance name already exsist ",a.Name)
		}
	}
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	alliance.Id = utils.GenerateRandomString(16)
	alliance.CreateTime = nowTime 
	alliances = append(alliances,alliance)
	bytes, err := json.Marshal(alliances)
	if err != nil {
		logger.Errorf("AppendAlliance: marshal alliances error,%v", err)
		return "",fmt.Errorf("%v", err)
	}
	err = utils.WriteFile(cfgPath,string(bytes))
	if err != nil {
		logger.Errorf("AppendAlliance: Write alliance list file error,%v", err)
		return "",fmt.Errorf("%v", err)
	} 
	return alliance.Id,nil
}

// func AddAlliance(alliance Alliance) error {
// 	err := AppendAlliance(alliance)
// 	if err != nil {
// 		logger.Errorf("AddAlliance: AddAlliance error,%v", err)
// 		return fmt.Errorf("%v", err)
// 	}
// 	AllianceFile := alliance.Name + ".json"
// 	cfgPath := utils.BAAS_CFG.AlliancePath + AllianceFile
// 	exsist,_ := utils.PathExists(cfgPath)
// 	if exsist {
// 		logger.Errorf("%s%s","AddAlliance: alliance file already exsist ",alliance.Name)
// 		return fmt.Errorf("%s%s","AddAlliance: alliance file already exsist ",alliance.Name)
// 	}
// 	bytes, err := json.Marshal(alliance)
// 	if err != nil {
// 		logger.Errorf("AddAlliance: marshal alliance error,%v", err)
// 		return fmt.Errorf("%v", err)
// 	}
// 	err = utils.WriteFile(cfgPath,string(bytes))
// 	if err != nil {
// 		logger.Errorf("AddAlliance: Write alliance file error,%v", err)
// 		return fmt.Errorf("%v", err)
// 	}
// 	return nil
// }
 
func GetAlliances()([]Alliance,error) {
	cfgPath := utils.BAAS_CFG.AlliancePath + ALLIANCE_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var alliances []Alliance
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&alliances)
			if err != nil {
				logger.Errorf("GetAlliances: unmarshal alliances error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			return alliances,nil
		}else {
			logger.Errorf("GetAlliances: load alliance list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Debug("GetAlliances: not find alliance list file")
	return nil,fmt.Errorf("GetAlliances:not find alliance list file") 
}

func GetAllianceByName(name string)(*Alliance,error) {
	cfgPath := utils.BAAS_CFG.AlliancePath + ALLIANCE_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var alliances []Alliance
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&alliances)
			if err != nil {
				logger.Errorf("GetAllianceByName: unmarshal alliances error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			for _,a := range alliances {
				if a.Name == name {
					return &a,nil
				}
			}
			return nil,fmt.Errorf("GetAllianceByName: not find alliance in list")
		}else {
			logger.Errorf("GetAllianceByName: load alliance list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Debug("GetAllianceByName: not find alliance list file")
	return nil,fmt.Errorf("GetAllianceByName:not find alliance list file")
}

func GetAllianceById(id string)(*Alliance,error) {
	cfgPath := utils.BAAS_CFG.AlliancePath + ALLIANCE_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var alliances []Alliance
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&alliances)
			if err != nil {
				logger.Errorf("GetAllianceById: unmarshal alliances error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			for _,a := range alliances {
				if a.Id == id {
					return &a,nil
				}
			}
			return nil,fmt.Errorf("GetAllianceById: not find alliance in list")
		}else {
			logger.Errorf("GetAllianceById: load alliance list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Debug("GetAllianceById: not find alliance list file")
	return nil,fmt.Errorf("GetAllianceById:not find alliance list file")
}

func GetAllianceChains(allianceId string)([]k8s.Chain,error) {
	cfgPath := utils.BAAS_CFG.BlockNetCfgBasePath + k8s.CHAIN_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var chains []k8s.Chain
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&chains) 
			if err != nil {
				logger.Errorf("GetAllianceChains: unmarshal chains error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			var chs []k8s.Chain
			for _,c := range chains {
				if c.AllianceId == allianceId {
					chs = append(chs,c)
				}
			} 
			return chs,nil 
		}else {
			logger.Errorf("GetAllianceChains: load chain list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Debug("GetAllianceChains: not find chain list file")
	return nil,fmt.Errorf("GetAllianceChains: not find chain list file")
}

func GetAllianceClusters(allianceId string)([]k8s.Cluster,error) {
	cfgPath := utils.BAAS_CFG.ClusterCfgPath  + k8s.CLUSTER_CFG_FILE
	exsist,_ := utils.PathExists(cfgPath)
	var clusters []k8s.Cluster
	if exsist {
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&clusters)
			if err != nil {
				logger.Errorf("GetAllianceClusters: unmarshal clusters error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			var cs []k8s.Cluster
			for _,c := range clusters {
				if c.AllianceId == allianceId {
					cs = append(cs,c)
				}
			} 
			return cs,nil 
		}else {
			logger.Errorf("GetAllianceClusters: load cluster config error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Errorf("GetAllianceClusters: load cluster config error")
	return nil,fmt.Errorf("GetAllianceClusters: load cluster config error")
}

func DeleteAllianceById(mail string,id string)error {
	cfgPath := utils.BAAS_CFG.AlliancePath + ALLIANCE_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var alliances []Alliance
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&alliances)
			if err != nil {
				logger.Errorf("DeleteAllianceById: unmarshal alliances error,%v", err)
				return fmt.Errorf("%v", err) 
			}
			var as []Alliance
			find := false
			for _,a := range alliances {
				tmpAlliance := a
				if a.Id != id {
					as = append(as,tmpAlliance)
				}else{
					find = true
					if a.Creator != mail {
						logger.Debug("DeleteAllianceById: user is not the creator")
						return fmt.Errorf("DeleteAllianceById: user is not the creator")
					}
				}
			}
			if !find {
				logger.Debug("DeleteAllianceById: not find the alliance,id=" + id)
				return fmt.Errorf("DeleteAllianceById: not find the alliance,id=" + id)
			}
			err = DeleteUserAlliance(id) 
			if err != nil {
				logger.Errorf("DeleteAllianceById: delete user alliance error")
				return fmt.Errorf("DeleteAllianceById: delete user alliance error") 
			}
			bytes, err := json.Marshal(as)
			if err != nil {
				logger.Errorf("DeleteAllianceById: marshal alliances error,%v", err)
				return fmt.Errorf("%v", err)
			}
			err = utils.WriteFile(cfgPath,string(bytes))
			if err != nil {  
				logger.Errorf("DeleteAllianceById: Write alliances file error,%v", err)
				return fmt.Errorf("%v", err)
			} 
			return nil
		}else {
			logger.Errorf("DeleteAllianceById: load alliance list file error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}
	logger.Errorf("DeleteAllianceById: not find alliance list file")
	return fmt.Errorf("DeleteAllianceById:not find alliance list file")
}