// Copyright 2022 Alexandre Dumas
//
// License: LGPL

//
// This file contains the extra validators to check if it conforms the regulations.
//

package misc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Check the transaction before sending any value if:
// 1. the sender address is within the list of AllowTransfer
// 2. the value is less than or equal to LimitTransfer
func VerifySendValue(sender common.Address, tx *types.Transaction, limit *big.Int, allowList []string) bool {

	// If tx.Value is 0 or below the limit, the tx is fine to go.
	if tx.Value().Cmp(big.NewInt(0)) == 0 || (limit != nil && tx.Value().Cmp(limit) == -1) {
		return true
	}

	// If AllowTransfer is not empty, check if the sender is in the (unlimited) allow list.
	if allowList != nil && len(allowList) != 0 {
		for _, addr := range allowList {
			if addr == sender.Hex() {
				return true
			}
		}
	}

	// At last, the tx should be invalidated.
	return false
}
