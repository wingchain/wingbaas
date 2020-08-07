
package db

import(
	"fmt"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	// DB_BLOCK_PATH string = "/data/block.db"
	// DB_TX_PATH string = "/data/tx.db" 
)

func PutData(dbSubPath string,key []byte,value []byte) error { 
	root,_ := utils.GetProcessRunRoot()
	dbPath := root + dbSubPath
	db,err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		logger.Errorf("PutData: open db failed")
		return fmt.Errorf("PutData: open db failed")
	}
	defer db.Close()
	err = db.Put(key,value,nil)
	if err != nil {
		logger.Errorf("PutData: put failed")
		return fmt.Errorf("PutData: put failed")
	}
	return nil
}

func GetData(dbSubPath string,key []byte) ([]byte,error) {
	root,_ := utils.GetProcessRunRoot()
	dbPath := root + dbSubPath
	db,err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		logger.Errorf("GetData: open db failed")
		return nil,fmt.Errorf("GetData: open db failed")
	}
	defer db.Close()
	bytes,err := db.Get(key,nil)
	if err != nil {
		logger.Errorf("GetData: get failed")
		return nil,fmt.Errorf("GetData: get failed")
	}
	return bytes,nil
}

func DeleteData(dbSubPath string,key []byte) error {
	root,_ := utils.GetProcessRunRoot()
	dbPath := root + dbSubPath
	db,err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		logger.Errorf("DeleteData: open db failed")
		return fmt.Errorf("DeleteData: open db failed")
	}
	defer db.Close()
	err = db.Delete(key,nil)
	if err != nil {
		logger.Errorf("DeleteData: delete failed")
		return fmt.Errorf("DeleteData: delete failed")
	}
	return nil
}

