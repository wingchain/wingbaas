
package fabric

import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils" 
	"github.com/wingbaas/platformsrv/k8s" 
	"github.com/wingbaas/platformsrv/certgenerate/fabric"
	"github.com/wingbaas/platformsrv/sdk/sdkfabric"
	"github.com/wingbaas/platformsrv/settings/fabric/txgenerate" 
	"github.com/wingbaas/platformsrv/settings/fabric/cfgtxlator" 
	"github.com/wingbaas/platformsrv/certgenerate/fabric/gopkg.in/yaml.v2"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
	"github.com/wingbaas/platformsrv/k8s/deployfabric"
	genesisconfig "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
)

func ChainAddOrg(chainId string,p public.AddOrgConfig) error { 
	var msg string
	//var srcCfg public.DeployPara
	cfg,err := sdkfabric.LoadChainCfg(chainId)
	if err != nil {
		logger.Errorf("ChainAddOrg: load chain deploy para error")
		return fmt.Errorf("ChainAddOrg: load chain deploy para error")
	}
	srcCfg,_ := sdkfabric.LoadChainCfg(chainId)
	// srcCfg.DeployType = cfg.DeployType
	// srcCfg.Version = cfg.Version
	// srcCfg.CryptoType = cfg.CryptoType
	// srcCfg.ClusterId = cfg.ClusterId
	// srcCfg.DeployNetCfg.KafkaDeployNode = cfg.DeployNetCfg.KafkaDeployNode
	// srcCfg.DeployNetCfg.ZookeeperDeployNode = cfg.DeployNetCfg.ZookeeperDeployNode
	// srcCfg.DeployNetCfg.ToolsDeployNode = cfg.DeployNetCfg.ToolsDeployNode
	// srcCfg.DeployNetCfg.PeerOrgs = cfg.DeployNetCfg.PeerOrgs[:len(cfg.DeployNetCfg.PeerOrgs)-1]
	// //srcCfg.DeployNetCfg.PeerOrgs = append(srcCfg.DeployNetCfg.PeerOrgs,cfg.DeployNetCfg.PeerOrgs...) 
	// srcCfg.DeployNetCfg.OrdererOrgs = append(srcCfg.DeployNetCfg.OrdererOrgs,cfg.DeployNetCfg.OrdererOrgs...)
	ch,_ := k8s.GetChain(chainId,cfg.ClusterId)
	if ch == nil {
		msg = "ChainAddOrg: get chain info failed"
		logger.Errorf(msg)
		return fmt.Errorf(msg)
	}
	if ch.Status >k8s.CHAIN_STATUS_FREE {
		msg = "ChainAddOrg: chain status not support operation"
		logger.Errorf(msg)
		return fmt.Errorf(msg) 
	}
	ch.Status = k8s.CHAIN_STATUS_ORGADD
	err = k8s.UpdateChainStatus(*ch)
	if err != nil {
		logger.Errorf("ChainAddOrg: set chain status error")
		return fmt.Errorf("ChainAddOrg: set chain status error")
	} 
	ch.Status = k8s.CHAIN_STATUS_FREE
	defer k8s.UpdateChainStatus(*ch) 
	certBasePath := utils.BAAS_CFG.BlockNetCfgBasePath + chainId 
	certPath := certBasePath + "/crypto-config"
	bl, _ := utils.PathExists(certPath)
	if !bl {
		msg = "ChainAddOrg: src cert path not exsist"
		logger.Errorf(msg)
		return fmt.Errorf(msg)
	}
	cfg.DeployNetCfg.PeerOrgs = append(cfg.DeployNetCfg.PeerOrgs,p.PeerOrgs...)
	bytes, err := json.Marshal(cfg.DeployNetCfg) 
	if err != nil {
		logger.Errorf("ChainAddOrg: Marshal deploy net config error")
		return fmt.Errorf("ChainAddOrg: Marshal deploy net config error")
	}
	bl = fabric.Extend(string(bytes),certPath,cfg.CryptoType)
	if !bl {
		msg = "ChainAddOrg: extend cert failed,chainId=" + chainId
		logger.Errorf(msg)
		return fmt.Errorf(msg)
	}
	root,_ := utils.GetProcessRunRoot()
	coreLocalPath := root + "/conf/core"
	coreNfsPath := utils.BAAS_CFG.NfsLocalRootDir + chainId + "/core/" 
	bl,_ = utils.CopyDir(coreLocalPath,coreNfsPath)
	if !bl {
		msg = "ChainAddOrg: copy cert to  nfs dir failed"
		logger.Errorf(msg)
		return fmt.Errorf(msg)
	} 
	var oName string
	for _,org := range srcCfg.DeployNetCfg.PeerOrgs {
		oName = org.Name
		break
	}
	chObj,err := OrgQueryChannel(chainId,oName)
	if err != nil {
		msg = "ChainAddOrg: get channel list failed"
		logger.Errorf(msg)
		return fmt.Errorf(msg)
	}
	var chList public.ChannelList
	chBytes,_ := json.Marshal(chObj)
	err = json.Unmarshal(chBytes,&chList)
	if err != nil {
		msg = "ChainAddOrg: unmarshal channel list failed"
		logger.Errorf(msg)
		return fmt.Errorf(msg)
	}
	if len(chList.Channels) < 1 {
		msg = "ChainAddOrg: at least exsist one channel,add org failed"
		logger.Errorf(msg)
		return fmt.Errorf(msg)
	}
	for _,org := range p.PeerOrgs {
		var txCnf genesisconfig.TopLevel
		var tmpOrg genesisconfig.Organization
		tmpOrg.ID = org.Name + "MSP"
		tmpOrg.Name = org.Name + "MSP"
		tmpOrg.MSPDir = certBasePath + "/crypto-config/peerOrganizations/" + org.Domain + "/msp"
		for _,h := range org.Specs {
			var anchor genesisconfig.AnchorPeer
			anchor.Host = h.Hostname
			anchor.Port = 7051
			tmpOrg.AnchorPeers = append(tmpOrg.AnchorPeers, &anchor)
			break
		}
		txCnf.Organizations = append(txCnf.Organizations,&tmpOrg)
		for _,cha := range chList.Channels {
			_,err = AddOrgV1(ch,srcCfg,cha.ChannelID,org.Name,txCnf)
			if err != nil {
				msg = "ChainAddOrg failed,name=" + org.Name
				logger.Errorf(msg)
				return fmt.Errorf(msg)
			}
		}
	}
	newBytes,err := json.Marshal(cfg) 
	if err != nil {
		logger.Errorf("ChainAddOrg: Marshal new config error")
		return fmt.Errorf("ChainAddOrg: Marshal new config error")
	} 
	blockCfgFile := utils.BAAS_CFG.BlockNetCfgBasePath + "/" + chainId + ".json" 
	err = utils.WriteFile(blockCfgFile,string(newBytes))
	if err != nil {
		logger.Errorf("ChainAddOrg: write chain new config error")
		return fmt.Errorf("ChainAddOrg: write chain new config error") 
	}
	certNfsPath := utils.BAAS_CFG.NfsLocalRootDir + chainId + "/crypto-config/"
	bl,_ = utils.CopyDir(certPath,certNfsPath)
	if !bl {
		msg = "ChainAddOrg: copy cert to  nfs dir failed"
		logger.Errorf(msg)
		return fmt.Errorf(msg)
	} 
	//deploy org 
	_,err = DeployOrg(ch.ClusterId,ch.BlockChainName,ch.BlockChainId,ch.BlockChainType,ch.Version,p.PeerOrgs)
	if err != nil {
		msg = "ChainAddOrg: DeployOrg failed"
		logger.Errorf(msg)
		return fmt.Errorf(msg)
	}
	time.Sleep(10*time.Second)
	var sdkCfg sdkfabric.GenerateParaSt
	sdkCfg.ClusterId = cfg.ClusterId
	sdkCfg.NamespaceId = ch.BlockChainName  
	sdkCfg.BlockId = chainId
	sdkfabric.UpdateOrgCfg(cfg.DeployNetCfg,sdkCfg,chList)
	logger.Debug("ChainAddOrg success")
	return nil 
}


func AddOrgV1(ch *k8s.Chain,srcCfg public.DeployPara,channelName string, orgName string, txCnf genesisconfig.TopLevel) (string,error) { 
	// logger.Debug("AddOrgV1 src cfg org=")
	// logger.Debug(srcCfg)
	var msg string
	var orderAddr string
	var orderCaFile string
	var orderMspId string
	var orderMspPath string
	coreFilePath := "/var/data/core/coreV1/"
	opCliBasePath := "/var/data/" + orgName + "/" 
	opLocalBasePath := utils.BAAS_CFG.NfsLocalRootDir + ch.BlockChainId + "/" + orgName + "/"
	bl,_ := utils.CreateDir(opLocalBasePath)
	if !bl {
		msg = "AddOrgV1: create local base dir failed"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	for _,org := range srcCfg.DeployNetCfg.OrdererOrgs {  
		orderMspId = org.Name + "MSP"
		for _,spec :=  range org.Specs {
			orderAddr = spec.Hostname + ":7050"
			orderCaFile = "/cert/crypto-config/ordererOrganizations/" + org.Domain + "/orderers/" + spec.Hostname + "." + org.Domain + "/msp/"
			orderCaFile = orderCaFile + "tlscacerts/tlsca." + org.Domain + "-cert.pem"
			orderMspPath = "/cert/crypto-config/ordererOrganizations/" + org.Domain + "/users/Admin@" + org.Domain + "/msp"
			break
		} 
		break  
	}
	var updatePara deployfabric.OrgAddToolsDeploymentPara
	updatePara.CoreCfgPath = coreFilePath
	updatePara.MspId = orderMspId
	updatePara.MspPath = orderMspPath
	updatePara.CertFile = orderCaFile 
	updatePara.AddOrgName = orgName
	outFile := opCliBasePath + "config_block_" + channelName + ".pb" 
	bl = cfgtxlator.FetchChannelCfg(ch,srcCfg,updatePara,orderAddr,channelName,outFile)
	if !bl {
		msg = "AddOrgV1: FetchChannelCfg failed,channel=" + channelName
		logger.Errorf(msg)
		return "",fmt.Errorf(msg) 
	}
	inFile := opLocalBasePath + "config_block_" + channelName + ".pb" 
	outFile = opLocalBasePath + "decoded_" + channelName + ".json"

	bl,_ = utils.PathExists(inFile)
	for i:=0;i<20;i++ {
		if bl {
			break
		}
		time.Sleep(3*time.Second)
		bl,_ = utils.PathExists(inFile)
	}

	bl = cfgtxlator.DeCodeCfg("common.Block", inFile, outFile) 
	if !bl {
		msg = "AddOrgV1: failed to DecodeCfg"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	//os.Remove(inFile)
	inFile = outFile
	outFile = opLocalBasePath + channelName + ".json"
	bl = cfgtxlator.ExtractDecodedCfg(inFile, outFile)
	if !bl {
		msg = "AddOrgV1: failed to ExtractDecodedCfg"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	} 
	//os.Remove(inFile)
	addedOrgFile := opLocalBasePath + orgName + ".json"
	by, err := yaml.Marshal(txCnf)
	if err != nil {
		msg = "AddOrgV1: marshal tx config failed"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	cnfTxYaml := opLocalBasePath + "/configtx.yaml" 
	err = utils.WriteFile(cnfTxYaml, string(by))
	if err != nil {
		msg = "AddOrgV1: write tx config to file failed"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	bl = txgenerate.CreateAddOrgCfg(opLocalBasePath, addedOrgFile, orgName + "MSP")
	if !bl {
		msg = "AddOrgV1: failed to CreateAddOrgCfg" 
		logger.Errorf(msg)
		return "",fmt.Errorf(msg) 
	} 
	inFile = outFile
	outFile = opLocalBasePath + channelName + "_add.json"
	bl = cfgtxlator.AppendOrgToCfg(orgName + "MSP", inFile, addedOrgFile, outFile, opLocalBasePath)
	if !bl {
		msg = "AddOrgV1: failed to AppendOrgToCfg" 
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	outFile = opLocalBasePath + channelName + "_encoded.pb"
	bl = cfgtxlator.EnCodeCfg("common.Config", inFile, outFile) 
	if !bl {
		msg = "AddOrgV1: failed to encodeCfg"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	//os.Remove(inFile)
	inFile = opLocalBasePath + channelName + "_add.json"
	outFile = opLocalBasePath + channelName + "_add_encoded.pb"
	bl = cfgtxlator.EnCodeCfg("common.Config", inFile, outFile)
	if !bl {
		msg = "AddOrgV1: failed to encodeCfg added"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	//os.Remove(inFile)
	inFile = opLocalBasePath + channelName + "_encoded.pb"
	updatedPbFile := outFile
	outFile = opLocalBasePath + channelName + "_updated.pb"
	bl = cfgtxlator.UpdateCfg(inFile, updatedPbFile, outFile, channelName)
	if !bl {
		msg = "AddOrgV1: failed to UpdateCfg"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	//os.Remove(inFile)
	//os.Remove(updatedPbFile)
	inFile = outFile
	outFile = opLocalBasePath + channelName + "_updated_decode.json"
	bl = cfgtxlator.DeCodeCfg("common.ConfigUpdate", inFile, outFile)
	if !bl {
		msg = "AddOrgV1: failed to DecodeCfg updated"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	//os.Remove(inFile)
	inFile = outFile
	outFile = opLocalBasePath + channelName + "_update_extract.json"
	bl = cfgtxlator.ExtractUpdateDecodedCfg(inFile, outFile)
	if !bl {
		msg = "AddOrgV1: failed to ExtractUpdateDecodedCfg"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	//os.Remove(inFile)
	inFile = outFile
	outFile = opLocalBasePath + channelName + "_update_envelop.json"
	bl = cfgtxlator.CreatEnvelopFile(channelName, inFile, outFile, opLocalBasePath)
	if !bl {
		msg = "AddOrgV1: failed to CreatEnvelopFile"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	//os.Remove(inFile)
	inFile = outFile
	outFile = opLocalBasePath + channelName + "_update_envelop.pb"
	bl = cfgtxlator.EnCodeCfg("common.Envelope", inFile, outFile)
	if !bl {
		msg = "AddOrgV1: failed to encodeCfg envelop"
		logger.Errorf(msg)
		return "",fmt.Errorf(msg)
	}
	//os.Remove(inFile)
	pbFile := opCliBasePath + channelName + "_update_envelop.pb"
	for _,org := range srcCfg.DeployNetCfg.PeerOrgs { 
		var signPara deployfabric.OrgAddToolsDeploymentPara
		signPara.CoreCfgPath = coreFilePath
		signPara.MspId = org.Name + "MSP"
		signPara.MspPath = "/cert/crypto-config/peerOrganizations/" + org.Domain + "/users/Admin@" + org.Domain + "/msp"
		signPara.CertFile = "/cert/crypto-config/peerOrganizations/" + org.Domain + "/msp/tlscacerts/tlsca." + org.Domain + "-cert.pem"
		signPara.OrgName = org.Name
		signPara.PeerDomain = org.Domain
		signPara.AddOrgName = orgName
		for _,h := range org.Specs {
			signPara.PeerAddr = h.Hostname + ":7051"
			signPara.PeerName = h.Hostname
			break
		}
		logger.Debug("AddOrgV1: SignEnvelop org="+signPara.OrgName)
		bl = cfgtxlator.SignEnvelop(ch,srcCfg,signPara,pbFile)  
		if !bl {
			msg = "AddOrgV1: failed to SignEnvelop envelop,org=" + org.Name
			logger.Errorf(msg)
			return "",fmt.Errorf(msg)
		}
	}
	inFile = outFile 
	bl = cfgtxlator.SendEnvelop(ch,srcCfg,updatePara,pbFile,channelName,orderAddr)
	if !bl {
		msg = "AddOrgV1: failed to SendEnvelop" 
		logger.Errorf(msg)
		return "",fmt.Errorf(msg) 
	}
	// os.Remove(inFile)
	// os.Remove(cnfTxYaml)
	// os.Remove(addedOrgFile)
	logger.Debug("AddOrgV1: add org success")
	return "",nil 
}