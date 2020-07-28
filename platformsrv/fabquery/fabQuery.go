
package fabquery

import (
	"time"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/k8s"
	"github.com/wingbaas/platformsrv/sdk/sdkfabric"
	"github.com/wingbaas/platformsrv/settings/fabric"
	"github.com/wingbaas/platformsrv/settings/fabric/public"
	"github.com/wingbaas/platformsrv/db"
)

type TxMainInfo struct {
	TxId 		string 		          `json:"TxId"`
	Timestamp   string				  `json:"Timestamp"`
	Signature	[]byte 		          `json:"Signature"`
	BlockId      uint64			      `json:"BlockId"`
}

type BlockMainInfo struct {
	BlockId      uint64			      `json:"BlockId"`
	BlockHash    []byte			      `json:"BlockHash"`
	PreHash      []byte			      `json:"PreHash"`
	Timestamp   string				  `json:"Timestamp"`
	Txs          []TxMainInfo		  `json:"Txs"`
}

func TimeTransfer(t int64)string {
	timeTemplate := "2006-01-02 15:04:05"
	return time.Unix(t, 0).Format(timeTemplate) 
}

func TotalTxQueryRoutine() {
	logger.Debug("start TotalTxQueryRoutine......")
	var txRecord []public.TxRecordSt
	rdFile := utils.BAAS_CFG.BlockNetCfgBasePath + "tx.json"
	bl,_ := utils.PathExists(rdFile)
	if bl {
		bytes := utils.ReadFileBytes(rdFile)
		err := json.Unmarshal(bytes,&txRecord)
		if err != nil {
			logger.Errorf("TotalTxQueryRoutine: unmarshal txRecord failed,will stop TotalTxQueryRoutine")
			return
		}
	}
	go func(){ 
		for(true) { 
			clusters,err := k8s.GetClusters()
			if err != nil {
				time.Sleep(30*time.Second)
				continue
			}
			if len(clusters) == 0 {
				time.Sleep(30*time.Second)
				continue
			}
			for _,c := range clusters {
				chains,_ := k8s.GetChains(c.ClusterId)
				if len(chains) == 0{
					time.Sleep(30*time.Second)
					continue
				}
				for _,chain := range chains {
					cfg,err := sdkfabric.LoadChainCfg(chain.BlockChainId)
					if err != nil {
						logger.Errorf("TotalTxQueryRoutine: load chain deploy para error,chainid="+chain.BlockChainId)
						continue
					}
					var orgName string
					for _,org := range cfg.DeployNetCfg.PeerOrgs {
						orgName = org.Name
						break
					}
					chObj,err := fabric.OrgQueryChannel(chain.BlockChainId,orgName)
					if err != nil {
						logger.Errorf("TotalTxQueryRoutine: get channel list failed")
						continue
					}
					var chList public.ChannelList 
					chBytes,_ := json.Marshal(chObj)
					err = json.Unmarshal(chBytes,&chList)
					if err != nil {
						logger.Errorf("TotalTxQueryRoutine: unmarshal channel list failed")
						continue
					}
					if len(chList.Channels) < 1 {
						logger.Debug("TotalTxQueryRoutine: channel list len is 0")
						continue
					}
					
					for _,ch := range chList.Channels {
						qr,err := fabric.OrgQueryBlockChain(chain.BlockChainId,orgName,ch.ChannelID)
						if err != nil { 
							logger.Errorf("TotalTxQueryRoutine: query block height error,channel="+ch.ChannelID)
						    continue
						}
						bytes,_ := json.Marshal(qr)
						var qInfo sdkfabric.ChainInfo
						err = json.Unmarshal(bytes,&qInfo)
						if err != nil {
							logger.Errorf("TotalTxQueryRoutine: query result unmarshal error,channel="+ch.ChannelID)
						    continue
						}
						find := false
						for n,r := range txRecord {
							if r.BlockChainId == chain.BlockChainId {
								for k,c := range r.ChTx {
									if c.ChainnnelId == ch.ChannelID {
										if qInfo.Height > c.Height {
											newTx := getBlockTxCount(c.Height,qInfo.Height,chain.BlockChainId,orgName,ch.ChannelID)
											txRecord[n].ChTx[k].Height = qInfo.Height
											txRecord[n].ChTx[k].TxCount = txRecord[n].ChTx[k].TxCount + newTx
										}
										find = true
										break
									}
								}
								if !find {
									var rdCh public.ChannelTxSt
									rdCh.ChainnnelId = ch.ChannelID
									rdCh.Height = qInfo.Height
									newTx := getBlockTxCount(0,qInfo.Height,chain.BlockChainId,orgName,ch.ChannelID)
									rdCh.TxCount = newTx
									txRecord[n].ChTx = append(txRecord[n].ChTx,rdCh) 
								}
								break
							}	
						}
					}
				} 
			}
			if len(txRecord) >0 {
				rdBytes,_ := json.Marshal(txRecord)
				err = utils.WriteFile(rdFile,string(rdBytes))
				if err != nil {
					logger.Errorf("TotalTxQueryRoutine: Write tx record file error")
				}
			}
			time.Sleep(10*time.Second) 
		} 
	}()
	logger.Debug("end TotalTxQueryRoutine......") 
}

func getBlockTxCount(start uint64,end uint64,blockChainId string,orgName string,channelId string) uint64 {
	if start >= end {
		return 0
	}
	if start <0 || end <0 {
		return 0
	}
	var total uint64
	total = 0
	for i:= (start + 1);i<end;i++ {
		qr,err := fabric.OrgQueryBlockById(blockChainId,orgName,channelId,i)
		if err != nil {
			logger.Errorf("getBlockTxCount: query block tx error,blockid=%d",i)
			continue
		}
		var blockInfo sdkfabric.Block 
		bytes,_ := json.Marshal(qr)
		err = json.Unmarshal(bytes,&blockInfo)
		if err != nil {
			logger.Errorf("getBlockTxCount: query block tx unmarshal error,blockid=%d",i)
			continue
		}
		total = total + uint64(len(blockInfo.Transactions))
	} 
	return total
}

func getTotalBlockTx(blockChainId string,orgName string,channelId string) {
	logger.Debug("getTotalBlockTx")
	var txRecord []public.TxRecordSt
	rdFile := utils.BAAS_CFG.BlockNetCfgBasePath + "tx.json"
	bl,_ := utils.PathExists(rdFile)
	if bl {
		bytes := utils.ReadFileBytes(rdFile)
		err := json.Unmarshal(bytes,&txRecord)
		if err != nil {
			logger.Errorf("getTotalBlockTx: unmarshal txRecord failed")
			return
		}
	}
	qr,err := fabric.OrgQueryBlockChain(blockChainId,orgName,channelId)
	if err != nil { 
		logger.Errorf("getTotalBlockTx: query block height error,channel="+channelId)
		return
	}
	bytes,_ := json.Marshal(qr)
	var qInfo sdkfabric.ChainInfo
	err = json.Unmarshal(bytes,&qInfo)
	if err != nil {
		logger.Errorf("getTotalBlockTx: query result unmarshal error,channel="+channelId)
		return
	}
	// logger.Debug("getTotalBlockTx: block height obj=")
	// logger.Debug(qInfo)
	blockFind := false
	chFind := false
	for n,r := range txRecord { 
		if r.BlockChainId == blockChainId {
			for k,c := range r.ChTx {
				if c.ChainnnelId == channelId {
					if qInfo.Height > c.Height {
						newTx := getBlockTxCount(c.Height,qInfo.Height,blockChainId,orgName,channelId)
						txRecord[n].ChTx[k].Height = qInfo.Height
						txRecord[n].ChTx[k].TxCount = txRecord[n].ChTx[k].TxCount + newTx
					}
					chFind = true
					break
				}	
			}
			if !chFind {
				var rdCh public.ChannelTxSt
				rdCh.BlockChainId = blockChainId
				rdCh.ChainnnelId = channelId
				rdCh.Height = qInfo.Height
				newTx := getBlockTxCount(0,qInfo.Height,blockChainId,orgName,channelId)
				rdCh.TxCount = newTx
				txRecord[n].ChTx = append(txRecord[n].ChTx,rdCh) 
			}
			blockFind = true
			break
		}	
	}
	if !blockFind {
		var rdBlock public.TxRecordSt
		var rdCha public.ChannelTxSt
		rdCha.BlockChainId = blockChainId
		rdCha.ChainnnelId = channelId
		rdCha.Height = qInfo.Height
		tmpTx := getBlockTxCount(0,qInfo.Height,blockChainId,orgName,channelId)
		rdCha.TxCount = tmpTx
		rdBlock.ChTx = append(rdBlock.ChTx,rdCha)
		rdBlock.BlockChainId =  blockChainId
		txRecord = append(txRecord,rdBlock)
	}
	if len(txRecord) >0 {
		rdBytes,_ := json.Marshal(txRecord) 
		err = utils.WriteFile(rdFile,string(rdBytes))
		if err != nil {
			logger.Errorf("getTotalBlockTx: Write tx record file error")
			return
		}
	}
}

func GetBlockTx(blockChainId string,orgName string,channelId string) (interface{},error) {
	var txRecord []public.TxRecordSt
	ch,_ := k8s.GetChainById(blockChainId)
	if ch == nil {
		logger.Errorf("GetBlockTx: blockchain id not exsist")
		return nil,fmt.Errorf("GetBlockTx: blockchain id not exsist")
	}
	getTotalBlockTx(blockChainId,orgName,channelId)
	rdFile := utils.BAAS_CFG.BlockNetCfgBasePath + "tx.json"
	bl,_ := utils.PathExists(rdFile)
	if bl {
		bytes := utils.ReadFileBytes(rdFile)
		err := json.Unmarshal(bytes,&txRecord)
		if err != nil {
			logger.Errorf("GetBlockTx: unmarshal txRecord failed")
			return nil,fmt.Errorf("GetBlockTx: unmarshal txRecord failed")
		}
	} 
	for n,r := range txRecord {
		if r.BlockChainId == blockChainId { 
			for k,c := range r.ChTx {
				if c.ChainnnelId == channelId {
					return txRecord[n].ChTx[k],nil
				}
			}
		}
	}
	logger.Debug("GetBlockTx: not find")
	return nil,nil
}

func GetBlockAndTxList(start uint64,end uint64,blockChainId string,orgName string,channelId string)(interface{},error){
	if start >= end {
		return nil,fmt.Errorf("GetBlockAndTxList: start >= end error")
	}
	if start <0 || end <0 {
		return nil,fmt.Errorf("GetBlockAndTxList: start or end error")
	}
	if (end - start) >10 {
		return nil,fmt.Errorf("GetBlockAndTxList: end and start gap must <= 10")
	}
	qr,err := fabric.OrgQueryBlockChain(blockChainId,orgName,channelId)
	if err != nil { 
		logger.Errorf("GetBlockAndTxList: query block height error,channel="+channelId)
		return nil,fmt.Errorf("GetBlockAndTxList: query block height error,channel="+channelId)
	}
	bytes,_ := json.Marshal(qr)
	var qInfo sdkfabric.ChainInfo
	err = json.Unmarshal(bytes,&qInfo)
	if err != nil {
		logger.Errorf("GetBlockAndTxList: query result unmarshal error,channel="+channelId)
		return nil,fmt.Errorf("GetBlockAndTxList: query result unmarshal error,channel="+channelId)
	}
	if end > qInfo.Height {
		logger.Errorf("GetBlockAndTxList: end error")
		return nil,fmt.Errorf("GetBlockAndTxList:end error")
	}
	var blocks []BlockMainInfo
	for i:= (start + 1);i<end;i++ {
		var block BlockMainInfo
		strIndex := strconv.FormatInt(int64(i),10)
		blockBytes,err := db.GetData(db.DB_BLOCK_PATH,[]byte(strIndex))
		if blockBytes != nil && len(blockBytes) > 0 {
			err = json.Unmarshal(blockBytes,&block)
		}
		if err != nil {
			qr,err := fabric.OrgQueryBlockById(blockChainId,orgName,channelId,i)
			if err != nil {
				logger.Errorf("GetBlockAndTxList: query block tx error,blockid=%d",i)
				continue
			}
			var blockInfo sdkfabric.Block 
			bytes,_ := json.Marshal(qr)
			err = json.Unmarshal(bytes,&blockInfo)
			if err != nil {
				logger.Errorf("GetBlockAndTxList: query block tx unmarshal error,blockid=%d",i)
				continue
			}	
			block.BlockId = blockInfo.Header.Number
			block.BlockHash = blockInfo.Header.DataHash
			block.PreHash = blockInfo.Header.PreviousHash
			txCount := len(blockInfo.Transactions) 
			block.Timestamp = TimeTransfer(blockInfo.Transactions[txCount-1].ChannelHeader.Timestamp.Seconds)
			for _,tx := range blockInfo.Transactions {
				var tmpTx TxMainInfo
				tmpTx.TxId = tx.ChannelHeader.TxID
				tmpTx.Signature = tx.Signature
				tmpTx.Timestamp = TimeTransfer(tx.ChannelHeader.Timestamp.Seconds)
				tmpTx.BlockId = block.BlockId
				block.Txs = append(block.Txs,tmpTx)

				bytes,_ = json.Marshal(tmpTx)
				err = db.PutData(db.DB_TX_PATH,[]byte(tmpTx.TxId),bytes)
				if err != nil{
					logger.Errorf("GetBlockAndTxList: write tx to db error,txid=")
					logger.Debug(tmpTx.TxId)
				}
			}
			bytes,_ = json.Marshal(block)
			err = db.PutData(db.DB_BLOCK_PATH,[]byte(strIndex),bytes)
			if err != nil{
				logger.Errorf("GetBlockAndTxList: write block to db error,blockid=%d",i)
			}
		}
		blocks = append(blocks,block)
	}  
	return blocks,nil
}

func GetTxInfo(txHash string)(interface{},error) {
	var tx TxMainInfo
	bytes,err := db.GetData(db.DB_TX_PATH,[]byte(txHash))
	if bytes != nil && len(bytes) > 0 {
		err = json.Unmarshal(bytes,&tx)
	}
	if err != nil {
		logger.Errorf("GetTxInfo: query tx hash error")
		return nil,fmt.Errorf("GetTxInfo: query tx hash error")
	}
	return tx,nil  
}

func GetSearchResult(key string)(interface{},error) {
	var tx TxMainInfo
	bytes,err := db.GetData(db.DB_TX_PATH,[]byte(key))
	if bytes != nil && len(bytes) > 0 {
		err = json.Unmarshal(bytes,&tx)
	}
	if err == nil {
		return tx,nil
	}
	var block BlockMainInfo
	bytes,err = db.GetData(db.DB_BLOCK_PATH,[]byte(key))
	if bytes != nil && len(bytes) > 0 {
		err = json.Unmarshal(bytes,&block)
	}
	if err == nil {
		return block,nil 
	}
	return nil,nil 
}


