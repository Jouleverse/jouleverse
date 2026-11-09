package core

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

// TestCanTransferGasBlindSync covers the BlindSync activation boundary and the
// nil-activation case in canTransferGas.
//
// A nil BlindSync means no activation height is configured, i.e. the transfer
// ban never activates. This case used to panic with a nil pointer dereference
// because the field was dereferenced unconditionally.
func TestCanTransferGasBlindSync(t *testing.T) {
	var (
		sender  = common.HexToAddress("0x1111111111111111111111111111111111111111")
		zeroVal = types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1), nil)
		posVal  = types.NewTransaction(0, common.Address{}, big.NewInt(1), 21000, big.NewInt(1), nil)
	)

	tests := []struct {
		name     string
		tx       *types.Transaction
		height   int64
		config   vm.Config
		expected bool
	}{
		// Activation height boundary: below it the ban is not enforced.
		{"below activation", posVal, 2139999, vm.Config{BlindSync: big.NewInt(2140000)}, true},
		// The activation block itself is checked.
		{"at activation", posVal, 2140000, vm.Config{BlindSync: big.NewInt(2140000)}, false},
		{"above activation", posVal, 2140001, vm.Config{BlindSync: big.NewInt(2140000)}, false},
		// Zero-value transactions are never rejected by the limit check.
		{"zero value above activation", zeroVal, 3000000, vm.Config{BlindSync: big.NewInt(2140000)}, true},

		// Nil BlindSync: ban never activates, so every transaction passes.
		// This is the case that previously panicked.
		{"nil activation, positive value", posVal, 3000000, vm.Config{}, true},
		{"nil activation, low height", posVal, 1, vm.Config{}, true},
		{"nil activation, zero value", zeroVal, 3000000, vm.Config{}, true},

		// The ethconfig default effectively never activates the ban.
		{"default max height", posVal, 8613794, vm.Config{BlindSync: maxUint63()}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// A nil BlindSync must not panic.
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("canTransferGas panicked: %v", r)
				}
			}()
			if got := canTransferGas(sender, tt.tx, big.NewInt(tt.height), tt.config); got != tt.expected {
				t.Errorf("canTransferGas(height=%d) = %v, want %v", tt.height, got, tt.expected)
			}
		})
	}
}

// maxUint63 mirrors the ethconfig default (math.MaxBig63).
func maxUint63() *big.Int {
	return new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 63), big.NewInt(1))
}
