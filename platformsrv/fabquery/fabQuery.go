
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

func TimeReverse(t string)int64 {
	timeTemplate := "2006-01-02 15:04:05"
	stamp, _ := time.ParseInLocation(timeTemplate,t,time.Local) 
	return stamp.Unix() 
}

func TimeReverseStamp(t string)time.Time {
	timeTemplate := "2006-01-02 15:04:05"
	stamp, _ := time.ParseInLocation(timeTemplate,t,time.Local) 
	return stamp
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
	TX_DB := "/data/" + blockChainId + "/tx.db"
	BLOCK_DB := "/data/" + blockChainId + "/block.db" 
	if start >= end {
		return nil,fmt.Errorf("GetBlockAndTxList: start >= end error")
	}
	if start <0 || end <0 {
		return nil,fmt.Errorf("GetBlockAndTxList: start or end error")
	}
	if (end - start) >31 {
		return nil,fmt.Errorf("GetBlockAndTxList: end and start gap must <= 30")
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
		blockBytes,err := db.GetData(BLOCK_DB,[]byte(strIndex))
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
				err = db.PutData(TX_DB,[]byte(tmpTx.TxId),bytes)
				if err != nil{
					logger.Errorf("GetBlockAndTxList: write tx to db error,txid=")
					logger.Debug(tmpTx.TxId)
				}
			}
			bytes,_ = json.Marshal(block)
			err = db.PutData(BLOCK_DB,[]byte(strIndex),bytes)
			if err != nil{
				logger.Errorf("GetBlockAndTxList: write block to db error,blockid=%d",i)
			}
		}
		blocks = append(blocks,block)
	}  
	return blocks,nil
}

func GetTxInfo(blockChainId string,txHash string)(interface{},error) {
	var tx TxMainInfo
	TX_DB := "/data/" + blockChainId + "/tx.db"
	bytes,err := db.GetData(TX_DB,[]byte(txHash))
	if bytes != nil && len(bytes) > 0 {
		err = json.Unmarshal(bytes,&tx)
	}
	if err != nil {
		logger.Errorf("GetTxInfo: query tx hash error")
		return nil,fmt.Errorf("GetTxInfo: query tx hash error")
	}
	return tx,nil  
}

func GetSearchResult(blockChainId string,key string)(interface{},error) {
	var tx TxMainInfo
	TX_DB := "/data/" + blockChainId + "/tx.db"
	BLOCK_DB := "/data/" + blockChainId + "/block.db" 
	bytes,err := db.GetData(TX_DB,[]byte(key))
	if bytes != nil && len(bytes) > 0 {
		err = json.Unmarshal(bytes,&tx)
	}
	if err == nil {
		return tx,nil
	}
	var block BlockMainInfo
	bytes,err = db.GetData(BLOCK_DB,[]byte(key))
	if bytes != nil && len(bytes) > 0 {
		err = json.Unmarshal(bytes,&block)
	}
	if err == nil {
		return block,nil 
	}
	return nil,nil 
}

func GetPassDayTxCount(blockChainId string,orgName string,channelId string,day int)(interface{},error) { 
	TX_DB := "/data/" + blockChainId + "/tx.db"
	BLOCK_DB := "/data/" + blockChainId + "/block.db" 
	if day < -14 || day >= 0 {
		logger.Errorf("GetPassDayTxCount: only query passed 14 days tx")
		return nil,fmt.Errorf("GetPassDayTxCount: only query passed 14 days tx")
	}
	qr,err := fabric.OrgQueryBlockChain(blockChainId,orgName,channelId)
	if err != nil { 
		logger.Errorf("GetPassDayTxCount: query block height error,channel="+channelId)
		return nil,fmt.Errorf("GetPassDayTxCount: query block height error,channel="+channelId)
	}
	bytes,_ := json.Marshal(qr)
	var qInfo sdkfabric.ChainInfo
	err = json.Unmarshal(bytes,&qInfo)
	if err != nil {
		logger.Errorf("GetPassDayTxCount: query result unmarshal error,channel="+channelId)
		return nil,fmt.Errorf("GetPassDayTxCount: query result unmarshal error,channel="+channelId)
	}
	startTime := time.Now().AddDate(0, 0, -1)
	startSecond := startTime.Unix()
	logger.Debug("start time second=")
	logger.Debug(startSecond)
	var txs []TxMainInfo
	for i:= (qInfo.Height - 1);i>0;i-- {
		var block BlockMainInfo
		strIndex := strconv.FormatInt(int64(i),10)
		blockBytes,err := db.GetData(BLOCK_DB,[]byte(strIndex))
		if blockBytes != nil && len(blockBytes) > 0 {
			err = json.Unmarshal(blockBytes,&block)
		}
		if err != nil {
			qr,err := fabric.OrgQueryBlockById(blockChainId,orgName,channelId,i)
			if err != nil {
				logger.Errorf("GetPassDayTxCount: query block tx error,blockid=%d",i)
				continue
			}
			var blockInfo sdkfabric.Block 
			bytes,_ := json.Marshal(qr)
			err = json.Unmarshal(bytes,&blockInfo)
			if err != nil {
				logger.Errorf("GetPassDayTxCount: query block tx unmarshal error,blockid=%d",i)
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
				err = db.PutData(TX_DB,[]byte(tmpTx.TxId),bytes) 
				if err != nil{
					logger.Errorf("GetPassDayTxCount: write tx to db error,txid=")
					logger.Debug(tmpTx.TxId)
				}
				if tx.ChannelHeader.Timestamp.Seconds >= startSecond {
					txs = append(txs,tmpTx)
				}else {
					return txs,nil
				}
			}
			bytes,_ = json.Marshal(block)
			err = db.PutData(BLOCK_DB,[]byte(strIndex),bytes)
			if err != nil{
				logger.Errorf("GetPassDayTxCount: write block to db error,blockid=%d",i)
			}
		}else{
			for _,tx := range block.Txs {
				tmpTx := tx
				txTime := TimeReverse(tx.Timestamp)
				if txTime >= startSecond {
					txs = append(txs,tmpTx)
				}else {
					return txs,nil
				}
			}
		}
	}
	return txs,nil
}

func GetDayTxCount(start string,end string,blockChainId string,orgName string,channelId string)(interface{},error) { 
	qr,err := fabric.OrgQueryBlockChain(blockChainId,orgName,channelId)
	if err != nil { 
		logger.Errorf("GetDayTxCount: query block height error,channel="+channelId)
		return nil,fmt.Errorf("GetDayTxCount: query block height error,channel="+channelId)
	}
	bytes,_ := json.Marshal(qr)
	var qInfo sdkfabric.ChainInfo
	err = json.Unmarshal(bytes,&qInfo)
	if err != nil {
		logger.Errorf("GetDayTxCount: query result unmarshal error,channel="+channelId)
		return nil,fmt.Errorf("GetDayTxCount: query result unmarshal error,channel="+channelId)
	}
	TX_DB := "/data/" + blockChainId + "/tx.db"
	BLOCK_DB := "/data/" + blockChainId + "/block.db" 
	startTime := TimeReverse(start)
	endTime := TimeReverse(end)
	// logger.Debug("start time=")
	// logger.Debug(startTime)
	// logger.Debug("end time=")
	// logger.Debug(endTime)
	// logger.Debug("now time=")
	// logger.Debug(time.Now().Unix())
	var txs []TxMainInfo
	for i:= (qInfo.Height - 1);i>0;i-- {
		var block BlockMainInfo
		strIndex := strconv.FormatInt(int64(i),10)
		blockBytes,err := db.GetData(BLOCK_DB,[]byte(strIndex))
		if blockBytes != nil && len(blockBytes) > 0 {
			err = json.Unmarshal(blockBytes,&block)
		}
		if err != nil {
			qr,err := fabric.OrgQueryBlockById(blockChainId,orgName,channelId,i)
			if err != nil {
				logger.Errorf("GetDayTxCount: query block tx error,blockid=%d",i)
				continue
			}
			var blockInfo sdkfabric.Block 
			bytes,_ := json.Marshal(qr)
			err = json.Unmarshal(bytes,&blockInfo)
			if err != nil {
				logger.Errorf("GetDayTxCount: query block tx unmarshal error,blockid=%d",i) 
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
				err = db.PutData(TX_DB,[]byte(tmpTx.TxId),bytes)
				if err != nil{
					logger.Errorf("GetDayTxCount: write tx to db error,txid=")
					logger.Debug(tmpTx.TxId)
				}
				// logger.Debug("query tx timestamp=")
				// logger.Debug(tx.ChannelHeader.Timestamp.Seconds)
				if tx.ChannelHeader.Timestamp.Seconds >= startTime &&  tx.ChannelHeader.Timestamp.Seconds <= endTime {
					txs = append(txs,tmpTx)
				}else if tx.ChannelHeader.Timestamp.Seconds < startTime {
					return txs,nil
				}
			}
			bytes,_ = json.Marshal(block)
			err = db.PutData(BLOCK_DB,[]byte(strIndex),bytes)
			if err != nil{
				logger.Errorf("GetDayTxCount: write block to db error,blockid=%d",i)
			}
		}else{
			for _,tx := range block.Txs {
				tmpTx := tx
				txTime := TimeReverse(tx.Timestamp)
				// logger.Debug("query db tx timestamp=")
				// logger.Debug(txTime)
				if txTime >= startTime &&  txTime <= endTime {
					txs = append(txs,tmpTx)
				}else if txTime < startTime {
					return txs,nil
				}
			}
		}
	}
	return txs,nil
}

type DayTx struct{
	TimeStamp string `json:"TimeStamp"`
	TxCount int64 `json:"TxCount"`
}

type DayTxStatistic struct{
	Txs []DayTx `json:"Txs"`
	Max int64 `json:"Max"`
	Min int64 `json:"Min"`
}

func GetEveryDayTxCount(start string,days int,blockChainId string,orgName string,channelId string)(interface{},error) { 
	if TimeReverse(start) >= time.Now().Unix() {
		logger.Errorf("GetEveryDayTxCount: start time must < now time")
		return nil,fmt.Errorf("GetEveryDayTxCount: start time must < now time")
	}
	end := ((TimeReverseStamp(start)).AddDate(0, 0, days)).Unix() 
	if end > time.Now().Unix() {
		logger.Errorf("GetEveryDayTxCount: end time must < now time")
		return nil,fmt.Errorf("GetEveryDayTxCount: end time must < now time")
	}
	qr,err := fabric.OrgQueryBlockChain(blockChainId,orgName,channelId)
	if err != nil { 
		logger.Errorf("GetEveryDayTxCount: query block height error,channel="+channelId)
		return nil,fmt.Errorf("GetEveryDayTxCount: query block height error,channel="+channelId)
	}
	bytes,_ := json.Marshal(qr)
	var qInfo sdkfabric.ChainInfo
	err = json.Unmarshal(bytes,&qInfo)
	if err != nil {
		logger.Errorf("GetEveryDayTxCount: query result unmarshal error,channel="+channelId)
		return nil,fmt.Errorf("GetEveryDayTxCount: query result unmarshal error,channel="+channelId)
	}
	TX_DB := "/data/" + blockChainId + "/tx.db"
	BLOCK_DB := "/data/" + blockChainId + "/block.db" 
	var daysTx DayTxStatistic
	var index uint64
	daysTx.Min = 9000000
	index = 1
	for day := 0; day<=days; day++ {
		startTime := ((TimeReverseStamp(start)).AddDate(0, 0, day)).Unix()
		endTime := ((TimeReverseStamp(start)).AddDate(0, 0, day+1)).Unix()
		// logger.Debug("start time=")
		// logger.Debug(startTime)
		// logger.Debug("end time=")
		// logger.Debug(endTime)
		// logger.Debug("now time=")
		// logger.Debug(time.Now().Unix())
		var i uint64
		var total int64
		var dayTx DayTx
		total = 0
		for i=index; i<=(qInfo.Height - 1);i++ {
			var block BlockMainInfo
			stop := false
			strIndex := strconv.FormatInt(int64(i),10)
			blockBytes,err := db.GetData(BLOCK_DB,[]byte(strIndex))
			if blockBytes != nil && len(blockBytes) > 0 {
				err = json.Unmarshal(blockBytes,&block)
			}
			if err != nil {
				qr,err := fabric.OrgQueryBlockById(blockChainId,orgName,channelId,i)
				if err != nil {
					logger.Errorf("GetEveryDayTxCount: query block tx error,blockid=%d",i)
					continue
				}
				var blockInfo sdkfabric.Block 
				bytes,_ := json.Marshal(qr)
				err = json.Unmarshal(bytes,&blockInfo)
				if err != nil {
					logger.Errorf("GetEveryDayTxCount: query block tx unmarshal error,blockid=%d",i) 
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
					err = db.PutData(TX_DB,[]byte(tmpTx.TxId),bytes)
					if err != nil{
						logger.Errorf("GetEveryDayTxCount: write tx to db error,txid=")
						logger.Debug(tmpTx.TxId)
					}
					// logger.Debug("query tx timestamp=")
					// logger.Debug(tx.ChannelHeader.Timestamp.Seconds)
					if tx.ChannelHeader.Timestamp.Seconds >= startTime &&  tx.ChannelHeader.Timestamp.Seconds <= endTime {
						total++
					}else if tx.ChannelHeader.Timestamp.Seconds > endTime {
						//logger.Debug("query tx stop")
						stop = true
						break
					}
				}
				bytes,_ = json.Marshal(block)
				err = db.PutData(BLOCK_DB,[]byte(strIndex),bytes)
				if err != nil{
					logger.Errorf("GetEveryDayTxCount: write block to db error,blockid=%d",i)
				}
			}else{
				for _,tx := range block.Txs {
					txTime := TimeReverse(tx.Timestamp)
					// logger.Debug("query db tx timestamp=")
					// logger.Debug(txTime)
					if txTime >= startTime &&  txTime <= endTime {
						total++
					}else if txTime > endTime {
						//logger.Debug("query db tx stop")
						stop = true
						break
					}
				}
			}
			if stop {
				index = i
				break
			}
		}
		dayTx.TimeStamp = TimeTransfer(startTime)
		dayTx.TxCount = total
		daysTx.Txs = append(daysTx.Txs,dayTx)
		if total > daysTx.Max {
			daysTx.Max = total
		}
		if total < daysTx.Min {
			daysTx.Min = total
		}
		total = 0
	}
	return daysTx,nil 
}

func GetEveryDayTxCountMulRoutine(start string,days int,blockChainId string,orgName string,channelId string)(interface{},error) { 
	if days > 30 || days < 1 {
		logger.Errorf("GetEveryDayTxCountMulRoutine: days must >= 1 and <= 30")
		return nil,fmt.Errorf("GetEveryDayTxCountMulRoutine: days must >= 1 and <= 30")
	}
	if TimeReverse(start) >= time.Now().Unix() {
		logger.Errorf("GetEveryDayTxCountMulRoutine: start time must < now time")
		return nil,fmt.Errorf("GetEveryDayTxCountMulRoutine: start time must < now time")
	}
	end := ((TimeReverseStamp(start)).AddDate(0, 0, days)).Unix() 
	if end > time.Now().Unix() {
		logger.Errorf("GetEveryDayTxCountMulRoutine: end time must < now time")
		return nil,fmt.Errorf("GetEveryDayTxCountMulRoutine: end time must < now time")
	}
	qr,err := fabric.OrgQueryBlockChain(blockChainId,orgName,channelId)
	if err != nil { 
		logger.Errorf("GetEveryDayTxCountMulRoutine: query block height error,channel="+channelId)
		return nil,fmt.Errorf("GetEveryDayTxCountMulRoutine: query block height error,channel="+channelId)
	}
	bytes,_ := json.Marshal(qr)
	var qInfo sdkfabric.ChainInfo
	err = json.Unmarshal(bytes,&qInfo)
	if err != nil {
		logger.Errorf("GetEveryDayTxCountMulRoutine: query result unmarshal error,channel="+channelId)
		return nil,fmt.Errorf("GetEveryDayTxCountMulRoutine: query result unmarshal error,channel="+channelId)
	}
	var dayTxs []DayTx
	var chs chan int
	dayTxs = make([]DayTx,days+1)
	chs = make(chan int,days+1)
	for day := 0; day<=days; day++ {
		startTime := ((TimeReverseStamp(start)).AddDate(0, 0, day)).Unix()
		endTime := ((TimeReverseStamp(start)).AddDate(0, 0, day+1)).Unix()
		go TxCountRoutine(startTime,endTime,blockChainId,orgName,channelId,qInfo.Height,day,dayTxs,chs)
	}
	result := 0
	for i:=0;i<=days;i++ {
		result += <-chs
	}
	logger.Debugf("GetEveryDayTxCountMulRoutine: thread end count=%d",result) 
	var daysTx DayTxStatistic
	daysTx.Max = 0
	daysTx.Min = 9000000
	daysTx.Txs = append(daysTx.Txs,dayTxs...)

	for _,dt := range dayTxs {
		if dt.TxCount > daysTx.Max {
			daysTx.Max = dt.TxCount
		}
		if dt.TxCount < daysTx.Min {
			daysTx.Min = dt.TxCount
		}
	}	
	return daysTx,nil 
}

func TxCountRoutine(startTime int64,endTime int64,blockChainId string,orgName string,channelId string,blockHeight uint64,id int,m []DayTx,chs chan int){
	var startIndex,endIndex uint64
	var i uint64
	var total int64
	var dayTx DayTx
	total = 0
	endIndex = blockHeight - 1
    if endIndex > 2 {
		startIndex = uint64(endIndex/2)
	}else{
		startIndex = 1
	}
	TX_DB := "/data/" + blockChainId + "/tx.db"
	BLOCK_DB := "/data/" + blockChainId + "/block.db" 
	for true {
		if startIndex<1 || endIndex<1 || startIndex >= endIndex {
			logger.Debugf("TxCountRoutine circle break,startIndex=%d  endIndex=%d",startIndex,endIndex)
			break
		}
		i = startIndex
		var block BlockMainInfo
		stop := false
		strIndex := strconv.FormatInt(int64(i),10)
		blockBytes,err := db.GetData(BLOCK_DB,[]byte(strIndex))
		if blockBytes != nil && len(blockBytes) > 0 {
			err = json.Unmarshal(blockBytes,&block)
		}
		if err != nil {
			qr,err := fabric.OrgQueryBlockById(blockChainId,orgName,channelId,i)
			if err != nil {
				logger.Errorf("TxCountRoutine: query block tx error,blockid=%d",i)
				continue
			}
			var blockInfo sdkfabric.Block 
			bytes,_ := json.Marshal(qr)
			err = json.Unmarshal(bytes,&blockInfo)
			if err != nil {
				logger.Errorf("TxCountRoutine: query block tx unmarshal error,blockid=%d",i) 
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
				err = db.PutData(TX_DB,[]byte(tmpTx.TxId),bytes)
				if err != nil{
					logger.Errorf("TxCountRoutine: write tx to db error,txid=")
					logger.Debug(tmpTx.TxId)
				}
				// logger.Debug("query tx timestamp=")
				// logger.Debug(tx.ChannelHeader.Timestamp.Seconds)
				if tx.ChannelHeader.Timestamp.Seconds >= startTime &&  tx.ChannelHeader.Timestamp.Seconds <= endTime {
					total++
					stop = true
				}else if tx.ChannelHeader.Timestamp.Seconds > endTime {
					//logger.Debug("query tx stop")
					startIndex = uint64((startIndex+1+endIndex)/2)
					break
				}else if tx.ChannelHeader.Timestamp.Seconds < startTime {
					if startIndex >= 1 {
						endIndex = startIndex - 1
					}else{
						endIndex = 0
					}
					startIndex = uint64((startIndex+endIndex)/2)
					break
				}
			}
			bytes,_ = json.Marshal(block)
			err = db.PutData(BLOCK_DB,[]byte(strIndex),bytes)
			if err != nil{
				logger.Errorf("TxCountRoutine: write block to db error,blockid=%d",i) 
			}
		}else{
			for _,tx := range block.Txs {
				txTime := TimeReverse(tx.Timestamp)
				// logger.Debug("query db tx timestamp=")
				// logger.Debug(txTime)
				if txTime >= startTime &&  txTime <= endTime {
					total++
					stop = true
				}else if txTime > endTime {
					//logger.Debug("query db tx stop")
					startIndex = uint64((startIndex+1+endIndex)/2)
					break
				}else if txTime < startTime {
					if startIndex >= 1 {
						endIndex = startIndex - 1
					}else{
						endIndex = 0
					}
					startIndex = uint64((startIndex+endIndex)/2)
					break
				}
			}
		} 
		//logger.Debugf("thread id=%d  startIndex=%d  endIndex=%d",id,startIndex,endIndex)
		if stop {
			logger.Debugf("TxCountRoutine circle break,find tx in block id=%d",i)
			break
		} 
	}
	dayTx.TimeStamp = TimeTransfer(startTime)
	dayTx.TxCount = total	
	m[id] = dayTx
	chs <- 1
} 


