/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package bridge

import (
	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/fabric-amcl/amcl"
	cryptolib "github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/idemix"
)

// NewRandOrPanic return a new amcl PRG or panic
func NewRandOrPanic() *amcl.RAND {
	rng, err := cryptolib.GetRand()
	if err != nil {
		panic(err)
	}
	return rng
}
