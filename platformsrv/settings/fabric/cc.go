
package fabric

import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
)

func RecordCC(blockChainId string,channelId string,ccName string,ccVersion string,flag string,status uint64) error {
	logger.Debug("RecordCC")
	var ccRecord []public.CCRecordSt 
	rdFile := utils.BAAS_CFG.BlockNetCfgBasePath + "cc.json"
	bl,_ := utils.PathExists(rdFile)
	if bl {
		bytes := utils.ReadFileBytes(rdFile)
		err := json.Unmarshal(bytes,&ccRecord)
		if err != nil {
			logger.Errorf("RecordCC: unmarshal ccRecord failed")
			return fmt.Errorf("RecordCC: unmarshal ccRecord failed")
		}
	}
	chain,_ := k8s.GetChainById(blockChainId)
	if chain == nil {
		logger.Errorf("RecordCC: get chain error")
		return fmt.Errorf("RecordCC: get chain error")
	}
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	var cc public.CCSt
	cc.BlockChainId = blockChainId
	cc.BlockChainName = chain.BlockChainName
	cc.CCName = ccName
	cc.CCVersion = ccVersion 
	cc.Status = status 
	if flag == "deploy" {
		cc.CreateTime = nowTime
		cc.UpdateTime = nowTime
	}else if flag == "update" {
		cc.UpdateTime = nowTime
	}
	var tmpCh public.ChannelCCSt
	tmpCh.ChainnnelId = channelId
	tmpCh.CCRecord = append(tmpCh.CCRecord,cc)

	blockFind := false
	chFind := false
	ccFind := false

	for n,r := range ccRecord { 
		if r.BlockChainId == blockChainId {
			for k,c := range r.ChCC {
				if c.ChainnnelId == channelId {
					for m,s := range c.CCRecord {
						if s.CCName == ccName {
							//ccRecord[n].ChCC[k].CCRecord[m].CCName = ccName
							ccRecord[n].ChCC[k].CCRecord[m].CCVersion = ccVersion
							ccRecord[n].ChCC[k].CCRecord[m].Status = status
							if flag == "deploy" {
								ccRecord[n].ChCC[k].CCRecord[m].CreateTime = nowTime
								ccRecord[n].ChCC[k].CCRecord[m].UpdateTime = nowTime
							}else if flag == "update" {
								ccRecord[n].ChCC[k].CCRecord[m].UpdateTime = nowTime
							}
							ccFind = true
							break
						}
					}	
					if !ccFind {
						ccRecord[n].ChCC[k].CCRecord = append(ccRecord[n].ChCC[k].CCRecord,cc)
					}
					chFind = true
					break
				}	
			}
			if !chFind {
				ccRecord[n].ChCC = append(ccRecord[n].ChCC,tmpCh)
			}
			blockFind = true
			break
		}	
	}
	if !blockFind {
		var record public.CCRecordSt
		record.BlockChainId = blockChainId
		record.BlockChainName = chain.BlockChainName
		record.ChCC = append(record.ChCC,tmpCh)
		ccRecord = append(ccRecord,record)
	}
	if len(ccRecord) >0 {
		rdBytes,_ := json.Marshal(ccRecord) 
		err := utils.WriteFile(rdFile,string(rdBytes))
		if err != nil {
			logger.Errorf("RecordCC: Write cc record file error")
			return fmt.Errorf("RecordCC: Write cc record file error")
		}
	}
	return nil
}

func GetCCRecord(blockChainId string,orgName string,channelId string) (interface{},error) {
	var ccRecord []public.CCRecordSt
	rdFile := utils.BAAS_CFG.BlockNetCfgBasePath + "cc.json"
	bl,_ := utils.PathExists(rdFile)
	if bl {
		bytes := utils.ReadFileBytes(rdFile)
		err := json.Unmarshal(bytes,&ccRecord)
		if err != nil {
			logger.Errorf("GetCCRecord: unmarshal ccRecord failed")
			return nil,fmt.Errorf("GetCCRecord: unmarshal ccRecord failed")
		}
	} 
	for n,r := range ccRecord {
		if r.BlockChainId == blockChainId { 
			for k,c := range r.ChCC {
				if c.ChainnnelId == channelId {
					return ccRecord[n].ChCC[k].CCRecord,nil 
				}
			}
		}
	}
	logger.Debug("GetCCRecord: not find")
	return nil,nil
}

func CheckCCRecord(blockChainId string,orgName string,channelId string) error {
	var ccRecord []public.CCRecordSt
	rdFile := utils.BAAS_CFG.BlockNetCfgBasePath + "cc.json"
	bl,_ := utils.PathExists(rdFile)
	if bl {
		bytes := utils.ReadFileBytes(rdFile)
		err := json.Unmarshal(bytes,&ccRecord)
		if err != nil {
			logger.Errorf("CheckCCRecord: unmarshal ccRecord failed")
			return fmt.Errorf("CheckCCRecord: unmarshal ccRecord failed")
		}
	} 
	for n,r := range ccRecord {
		if r.BlockChainId == blockChainId { 
			for k,c := range r.ChCC {
				if c.ChainnnelId == channelId {
					//return ccRecord[n].ChCC[k].CCRecord,nil 
					for _,cr := range ccRecord[n].ChCC[k].CCRecord {
						if cr.Status == public.CC_DEPLOY_OK {
							return nil
						}
					}
				}
			}
		}
	}
	logger.Debug("CheckCCRecord: not deploy successed cc")
	return fmt.Errorf("CheckCCRecord: not deploy successed cc") 
}