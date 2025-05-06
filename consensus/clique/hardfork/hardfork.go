package hardfork

import (
	"github.com/holiman/uint256"
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/core/state"
	"github.com/ledgerwatch/log/v3"
)

type HardForkInstruction struct {
	Name    string
	Storage map[libcommon.Address]map[libcommon.Hash]libcommon.Hash
	Code    map[libcommon.Address][]byte
}

func ApplyHardfork(state *state.IntraBlockState, instruction HardForkInstruction) {
	for address, storage := range instruction.Storage {
		for key, value := range storage {
			var k = key
			v, overflow := uint256.FromBig(value.Big())
			if !overflow {
				state.SetState(address, &k, *v)
				log.Debug("Set storage", "address", address.Hex, "key", key.Hex(), "value", value.Hex())
			} else {
				log.Error("Set storage (apply hardfork) failed to convert value to uint256", "value", value.Hex())
			}
		}
	}

	for address, bytecode := range instruction.Code {
		state.SetCode(address, bytecode)
	}
}

// IntToHash function is a utility function that allows us to convert
// slot numer to hash easily. Take a note that int is max at 64 bits.
func IntToHash(storageSlot int) libcommon.Hash {
	return libcommon.BytesToHash([]byte{byte(storageSlot)})
}
