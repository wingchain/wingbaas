
package cfgtxlator

import (
	"os"
	"os/exec"
	"strings"
	"time"
	_ "github.com/hyperledger/fabric/protos/common"
	_ "github.com/hyperledger/fabric/protos/msp"
	_ "github.com/hyperledger/fabric/protos/orderer"
	_ "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/common/tools/configtxlator/rest"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
	"github.com/wingbaas/platformsrv/k8s/deployfabric"
)

func ExtractUpdateDecodedCfg(inFile string,outFile string)bool{ 
	var cmdstr string
	cmdstr = ". " + inFile
	args := strings.Split(cmdstr," ")
	
	cmd := exec.Command("jq",args...)
	buf,err := cmd.Output()
	if err!=nil {
		logger.Errorf("ExtractUpdateDecodedCfg: exec cmdstr failed, args=%v, err:%v", args, err)
		return false
	} 
	destf, err := os.OpenFile(outFile,os.O_WRONLY|os.O_CREATE, 0644)        
	if err != nil {
		logger.Errorf("ExtractUpdateDecodedCfg: open outFile failed, %v", err)
		return false
	}
	defer destf.Close()
	_,err = destf.Write(buf)
	if err != nil {
        logger.Errorf("ExtractUpdateDecodedCfg: write outFile failed, %v", err)
        return false
	}
	return true
}

func ExtractDecodedCfg(inFile string,outFile string)bool{
	var cmdstr string
	cmdstr = ".data.data[0].payload.data.config " + inFile
	args := strings.Split(cmdstr," ")
	
	cmd := exec.Command("jq",args...)
	buf,err := cmd.Output()
	if err!=nil {
		logger.Debug("ExtractDecodedCfg: exec cmdstr failed, args=%v, err: %v", args, err)
		return false
	} 
	destf, err := os.OpenFile(outFile,os.O_WRONLY|os.O_CREATE, 0644)        
	if err != nil {
		logger.Errorf("ExtractDecodedCfg: open outFile failed, %v", err)
		return false
	}
	defer destf.Close()
	_,err = destf.Write(buf)
	if err != nil {
		logger.Errorf("ExtractDecodedCfg: write outFile failed, %v", err)
		return false
	}
	return true
}

func AppendOrgToCfg(orgName string,srcFile string,appendFile string,outFile string,basePath string)bool{
	var cmdstr string
	var shFile = basePath + "addCmd.sh"
	content := "'.[0]*{\"channel_group\":{\"groups\":{\"Application\":{\"groups\":{\"" + orgName + "\":.[1]}}}}}'"
	cmdstr = "jq -s " + content + " " + srcFile + " " + appendFile + " > " + outFile

	destf, err := os.OpenFile(shFile,os.O_WRONLY|os.O_CREATE, 0777)        
	if err != nil {
		logger.Errorf("AppendOrgToCfg: open shell file failed, %v", err)
		return false
	}
	_,err = destf.Write([]byte(cmdstr))
	if err != nil {
		logger.Errorf("AppendOrgToCfg: write shell file failed, %v", err)
		destf.Close()
		return false
	}
	destf.Close()
	cmd := exec.Command("/bin/bash","-c",shFile)

	_,err = cmd.Output()
	if err!=nil {
		logger.Errorf("AppendOrgToCfg: exec cmdstr failed, %v", err)
		return false
	}  
	os.Remove(shFile)
	return true
}

func FetchChannelCfg(ch *k8s.Chain,cfg public.DeployPara,updatePara deployfabric.OrgAddToolsDeploymentPara,orderHost string,channelName string,outFile string)bool{
	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "FetchChannelCfg: get tools image failed"
		logger.Errorf(msg)
		return false
	}
	updatePara.ToolsImage = toolsImage
	updatePara.LogLevel = "debug"
	args := []string{"sh","-c"}
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		updatePara.OrgName = org.Name
		updatePara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			updatePara.PeerName = spec.Hostname
			break
		}
		break
	}
	randomStr := utils.GenerateRandomString(8)
	var cmdstr string
	cmdstr = "peer channel fetch config "  
	cmdstr = cmdstr + outFile + " -o " + orderHost + " -c " + channelName + " --tls --cafile " + updatePara.CertFile
	cmd := "cp -a /var/data/. /cert;"
	cmd = cmd + " echo " + randomStr + ";"
	cmd = cmd + " $(" + cmdstr + ")" + ";"
	cmd = cmd + " /bin/bash"
	updatePara.Args = append(updatePara.Args,args...)
	updatePara.Args = append(updatePara.Args,cmd)
	_,err = deployfabric.OrgAddPatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,ch.BlockChainName,ch.BlockChainId,updatePara)
	if err != nil {
		msg := "FetchChannelCfg update add org cli failed"
		logger.Errorf(msg)
		return false
	}
	time.Sleep(10*time.Second)
	return true
}

func EnCodeCfg(msgName string,inputFile string,outputFile string)bool {
	inf, err := os.Open(inputFile)        
	if err != nil {
		logger.Errorf("EnCodeCfg: open inputfile failed, %v", err)
		return false
	}
	defer inf.Close()
	outf, err := os.OpenFile(outputFile,os.O_WRONLY|os.O_CREATE, 0644)        
	if err != nil {
		logger.Errorf("EnCodeCfg: open outputFile failed, %v", err)
		return false
	}
	defer outf.Close()
	err = rest.EncodeProto(msgName,inf,outf)
	if err != nil { 
		logger.Errorf("EnCodeCfg: Error encoding: %v", err)
		return false
	}
	return true
}

func DeCodeCfg(msgName string,inputFile string,outputFile string)bool {
	inf, err := os.Open(inputFile)        
	if err != nil {
		logger.Errorf("DecodeCfg: open inputfile failed, %v", err)
		return false
	}
	defer inf.Close()
	outf, err := os.OpenFile(outputFile,os.O_WRONLY|os.O_CREATE, 0644)        
	if err != nil {
		logger.Errorf("DecodeCfg: open outputFile failed, %v", err)
		return false
	}
	defer outf.Close()
	err = rest.DecodeProto(msgName,inf,outf) 
	if err != nil { 
		logger.Errorf("DecodeCfg: Error decoding: %v", err)
		return false 
	}
	return true 
}

func UpdateCfg(srcFile string,updatedFile string,updateDestFile string,channelName string)bool{
	srcf, err := os.Open(srcFile)        
	if err != nil {
		logger.Errorf("UpdateCfg: open srcFile failed, %v", err)
		return false
	}
	defer srcf.Close()
	updatedf, err := os.Open(updatedFile)        
	if err != nil {
		logger.Errorf("UpdateCfg: open updatedFile failed, %v", err)
		return false
	}
	defer updatedf.Close()
	destf, err := os.OpenFile(updateDestFile,os.O_WRONLY|os.O_CREATE, 0644)        
	if err != nil {
		logger.Errorf("UpdateCfg: open updateDestFile failed, %v", err)
		return false
	}
	defer destf.Close()
	err = rest.ComputeUpdt(srcf,updatedf,destf,channelName)
	if err != nil { 
		logger.Errorf("UpdateCfg: Error ComputeUpdt: %v", err)
		return false
	}
	return true
}

func CreatEnvelopFile(channelName string,srcFile string,outFile string,basePath string)bool{
	var cmdstr string
	var shFile = basePath + "envelopCmd.sh"
	cmdstr = "echo '{\"payload\":{\"header\":{\"channel_header\":{\"channel_id\":\""
	cmdstr = cmdstr + channelName + "\",\"type\":2}},\"data\":{\"config_update\":'$(cat "
	cmdstr = cmdstr + srcFile + ")'}}}' | jq . > " + outFile

	destf, err := os.OpenFile(shFile,os.O_WRONLY|os.O_CREATE, 0777)        
	if err != nil {
		logger.Errorf("CreatEnvelopFile: open shell file failed, %v", err)
		return false
	}
	_,err = destf.Write([]byte(cmdstr))
	if err != nil {
		logger.Errorf("CreatEnvelopFile: write shell file failed, %v", err)
		destf.Close()
		return false
	}
	destf.Close()
	cmd := exec.Command("/bin/bash","-c",shFile)

	_,err = cmd.Output()
	if err!=nil {
		logger.Errorf("CreatEnvelopFile: exec cmdstr failed, %v", err)
		return false
	}  
	os.Remove(shFile)
	return true
}

func SignEnvelop(ch *k8s.Chain,cfg public.DeployPara,signPara deployfabric.OrgAddToolsDeploymentPara,envPBFile string)bool{
	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "SignEnvelop: get tools image failed"
		logger.Errorf(msg)
		return false
	}
	signPara.ToolsImage = toolsImage
	signPara.LogLevel = "debug"
	args := []string{"sh","-c"}
	logFile := "/var/data/" + signPara.AddOrgName + "/"  + signPara.OrgName + "_signenvelop.log"
	randomStr := utils.GenerateRandomString(8)
	var cmdstr string 
	cmdstr = "peer channel signconfigtx -f "
	cmdstr = cmdstr + envPBFile
	cmdstr = cmdstr + " > " + logFile + " 2>&1"
	cmd := "cp -a /var/data/. /cert;"
	cmd = cmd + " echo " + randomStr + ";"
	cmd = cmd + " $(" + cmdstr + ")" + ";"
	cmd = cmd + " /bin/bash"
	signPara.Args = nil
	signPara.Args = append(signPara.Args,args...)
	signPara.Args = append(signPara.Args,cmd)
	_,err = deployfabric.OrgAddPatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,ch.BlockChainName,ch.BlockChainId,signPara)
	if err != nil {
		msg := "SignEnvelop update add org cli failed"
		logger.Errorf(msg)
		return false
	}
	outFile := utils.BAAS_CFG.NfsLocalRootDir + ch.BlockChainId + "/" + signPara.AddOrgName + "/" + signPara.OrgName + "_signenvelop.log"
	success := false
	for i:=0;i<30;i++ {
		bl,_ := utils.PathExists(outFile)
		if bl {
			success = true
			break
		}else{
			time.Sleep(5*time.Second)
		}
	}
	if !success {
		msg := "SignEnvelop update add org cli failed,not receive result,org=" + signPara.OrgName
		logger.Errorf(msg)
		return false
	}
	return true
}

func SendEnvelop(ch *k8s.Chain,cfg public.DeployPara,sendPara deployfabric.OrgAddToolsDeploymentPara,envPBFile string,channelName string,order string)bool{
	toolsImage,err := utils.GetBlockImage(public.BLOCK_CHAIN_TYPE_FABRIC,cfg.Version,"tools")
	if err != nil {
		msg := "SendEnvelop: get tools image failed"
		logger.Errorf(msg)
		return false
	}
	sendPara.ToolsImage = toolsImage
	sendPara.LogLevel = "debug"
	args := []string{"sh","-c"}
	for _,org := range cfg.DeployNetCfg.PeerOrgs {
		sendPara.OrgName = org.Name
		sendPara.PeerDomain = org.Domain
		for _,spec :=  range org.Specs {
			sendPara.PeerName = spec.Hostname 
			break
		}
		break
	}
	logFile := "/var/data/" + sendPara.AddOrgName + "/"  + sendPara.OrgName + "_sendenvelop.log"
	randomStr := utils.GenerateRandomString(8)
	var cmdstr string
	cmdstr = "peer channel update -f "
	cmdstr = cmdstr + envPBFile + " -o " + order  + " -c " + channelName + " --tls --cafile " + sendPara.CertFile
	cmdstr = cmdstr + " > " + logFile + " 2>&1"
	cmd := "cp -a /var/data/. /cert;"
	cmd = cmd + " echo " + randomStr + ";" 
	cmd = cmd + " $(" + cmdstr + ")" + ";"
	cmd = cmd + " /bin/bash"
	sendPara.Args = nil
	sendPara.Args = append(sendPara.Args,args...) 
	sendPara.Args = append(sendPara.Args,cmd)
	_,err = deployfabric.OrgAddPatchToolsDeployment(cfg.ClusterId,cfg.DeployNetCfg.ToolsDeployNode,ch.BlockChainName,ch.BlockChainId,sendPara)
	if err != nil {
		msg := "SendEnvelop update add org cli failed"
		logger.Errorf(msg)
		return false
	}
	outFile := utils.BAAS_CFG.NfsLocalRootDir + ch.BlockChainId + "/" + sendPara.AddOrgName + "/" + sendPara.OrgName + "_sendenvelop.log"
	success := false
	for i:=0;i<30;i++ {
		bl,_ := utils.PathExists(outFile)
		if bl {
			success = true
			break
		}else{
			time.Sleep(5*time.Second)
		}
	}
	if !success {
		msg := "SendEnvelop update add org cli failed,not receive result,org=" + sendPara.OrgName
		logger.Errorf(msg)
		return false
	}
	return true
}