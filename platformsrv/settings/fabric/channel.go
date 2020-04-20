
package fabric

import (
	"fmt"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
	"github.com/wingbaas/platformsrv/settings/fabric/txgenerate"
)

func GenerateChannelTx(serviceRootPath string,channelName string) error {
	txPath := serviceRootPath + "/channel-artifacts/"
	err := utils.DirCheck(txPath)
	if err != nil {
		logger.Errorf("GenerateChannelTx: artifacts dir create error,%v",err)
		return fmt.Errorf("GenerateChannelTx: artifacts dir create error,%v",err)
	}
	bl := txgenerate.CreateChannelTx(serviceRootPath, channelName, txPath)
	if !bl {
		logger.Errorf("GenerateChannelTx: create channel tx error")
		return fmt.Errorf("GenerateChannelTx: create channel tx error")
	}
	return nil
}