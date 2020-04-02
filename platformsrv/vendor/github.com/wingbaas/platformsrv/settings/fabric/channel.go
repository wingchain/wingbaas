
package fabric

import (
	"fmt"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/settings/fabric/txgenerate"
)

func GenerateChannelTx(serviceRootPath string,channelName string) error {
	txPath := serviceRootPath + "/channel-artifacts/"
	bl := txgenerate.CreateChannelTx(serviceRootPath, channelName, txPath)
	if !bl {
		logger.Errorf("GenerateChannelTx: create channel tx error")
		return fmt.Errorf("GenerateChannelTx: create channel tx error")
	}
	return nil
}