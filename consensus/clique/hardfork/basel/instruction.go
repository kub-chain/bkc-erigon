package basel

import (
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	libcommon "github.com/ledgerwatch/erigon-lib/common"

	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/consensus/clique/hardfork"
	"github.com/ledgerwatch/erigon/core/state"
	"github.com/ledgerwatch/erigon/crypto"
)

type BaselParams struct {
	StakeManagerV3             libcommon.Address
	StakeManagerStorageV3      libcommon.Address
	BKCValidatorSetV3          libcommon.Address
	StakeManagerVaultV3        libcommon.Address
	NftContractV3              libcommon.Address
	OfficialNodeValidatorShare libcommon.Address
	SlashManagerV3             libcommon.Address
	SuperNodeAddress           libcommon.Address
	SuperNodeOwnerAddress      libcommon.Address
	BaselBlock                 *big.Int
}

// StakeManager storage slots
const (
	SLOT_POOL_AMOUNT              = 9  // pool amount
	SLOT_OFFICIAL_AMOUNT          = 11 // official amount
	SLOT_VALIDATOR_ADDRESS_TO_IDS = 20 // mapping(address => uint256[])
	SLOT_VALIDATOR_LIST           = 21 // Validator[] array
	SLOT_MINIMAL_LIST             = 22 // EnumerableSet array
	SLOT_OFFICIAL_SIGNER_OLD      = 18 // Old official signer address
	SLOT_OFFICIAL_SIGNER_NEW      = 19 // New official signer address
	SLOT_SUPER_NODE_OLD           = 30 // Old super node address
	SLOT_SUPER_NODE_NEW           = 31 // New super node address
)

// NFT Contract storage slots
const (
	SLOT_HOLDER_TOKENS        = 1
	SLOT_TOKEN_OWNERS         = 2
	SLOT_TOKEN_OWNERS_INDEXES = 3
)

const (
	SLOT_IS_OFFICAL_POOL = 19
)

func New(state *state.IntraBlockState, params BaselParams) (hardfork.HardForkInstruction, error) {

	if params.StakeManagerV3 == (libcommon.Address{}) {
		return hardfork.HardForkInstruction{}, fmt.Errorf("create Lausanne hardfork requires StakeManagerV3 address")
	}
	if params.StakeManagerStorageV3 == (libcommon.Address{}) {
		return hardfork.HardForkInstruction{}, fmt.Errorf("create Lausanne hardfork requires StakeManagerStorageV3 address")
	}
	if params.SlashManagerV3 == (libcommon.Address{}) {
		return hardfork.HardForkInstruction{}, fmt.Errorf("create Lausanne hardfork requires SlashManagerV2 address")
	}
	if params.NftContractV3 == (libcommon.Address{}) {
		return hardfork.HardForkInstruction{}, fmt.Errorf("create Lausanne hardfork requires NftContract address")
	}
	if params.BKCValidatorSetV3 == (libcommon.Address{}) {
		return hardfork.HardForkInstruction{}, fmt.Errorf("create Lausanne hardfork requires BKCValidatorSetV3 address")
	}

	instruction := hardfork.HardForkInstruction{
		Name:    "Basel",
		Storage: make(map[libcommon.Address]map[libcommon.Hash]libcommon.Hash),
		Code:    make(map[libcommon.Address][]byte),
	}
	// Deploy V3 contracts with user's modifications
	instruction.Code[params.StakeManagerV3] = StakeManagerV3ByteCode
	instruction.Code[params.StakeManagerStorageV3] = StakeManagerStorageV3ByteCode
	instruction.Code[params.SlashManagerV3] = SlashManagerV3ByteCode
	instruction.Code[params.NftContractV3] = NftContractV3ByteCode
	instruction.Code[params.BKCValidatorSetV3] = BKCValidatorSetV3ByteCode

	var nextValidatorIdValue uint256.Int
	totalSupplySolt := libcommon.BigToHash(big.NewInt(SLOT_TOKEN_OWNERS))
	state.GetState(params.NftContractV3, &totalSupplySolt, &nextValidatorIdValue)
	nextValidatorId := nextValidatorIdValue.Uint64()

	nextLength := nextValidatorId + 1

	// Read current poolAmount and increment it for official node conversion
	var poolAmountValue uint256.Int
	poolAmountSlot := libcommon.BigToHash(big.NewInt(SLOT_POOL_AMOUNT))
	state.GetState(params.StakeManagerStorageV3, &poolAmountSlot, &poolAmountValue)
	poolAmount := poolAmountValue.Uint64()
	newPoolAmount := poolAmount + 1 // Convert official node (validator 0) to pool

	nextValidatorSoltInt := new(big.Int).SetBytes(getValidatorSlot(SLOT_VALIDATOR_LIST, nextValidatorId).Bytes())

	validatorIdsArrayLengthSlot := getMappingSlot(libcommon.BytesToHash(params.SuperNodeAddress.Bytes()), SLOT_VALIDATOR_ADDRESS_TO_IDS)

	validatorIdsArrayDataSlot := crypto.Keccak256Hash(validatorIdsArrayLengthSlot.Bytes())
	validatorIdsArrayElement0Slot := libcommon.BigToHash(new(big.Int).SetBytes(validatorIdsArrayDataSlot.Bytes()))

	var currentMinimalListLengthValue uint256.Int
	minimalValidatorListLengthSlot := libcommon.BigToHash(big.NewInt(SLOT_MINIMAL_LIST))
	state.GetState(params.StakeManagerStorageV3, &minimalValidatorListLengthSlot, &currentMinimalListLengthValue)
	currentMinimalListLength := currentMinimalListLengthValue.Uint64()
	minimalValidatorListDataSlotInt := crypto.Keccak256Hash(minimalValidatorListLengthSlot.Bytes()).Big()

	// EnumerableSet._indexes mapping is at base slot + 1 (slot 23 for SLOT_MINIMAL_LIST=22)
	minimalValidatorListIndexesSlot := SLOT_MINIMAL_LIST + 1
	minimalValidatorListDataIndexes := getMappingSlot(libcommon.BytesToHash(params.SuperNodeAddress.Bytes()), uint64(minimalValidatorListIndexesSlot))

	// Calculate superNodeAndValidBlock according to the formula:
	// uint256 tmpUint256 = uint256(uint160(superNode_));
	// tmpUint256 += blockNumber_ << 160;
	tmpUint256 := new(big.Int).SetBytes(params.SuperNodeAddress.Bytes())       // uint256(uint160(superNode_))
	shiftedBlockNumber := new(big.Int).Lsh(params.BaselBlock, 160)             // blockNumber_ << 160
	superNodeAndValidBlock := new(big.Int).Add(tmpUint256, shiftedBlockNumber) // tmpUint256 += blockNumber_ << 160

	// Update new offical node siginer to address 0, shift current to old
	var currentOfficialNodeSigner uint256.Int
	slotOfficialSignerNew := libcommon.BigToHash(big.NewInt(SLOT_OFFICIAL_SIGNER_NEW))
	state.GetState(params.StakeManagerStorageV3, &slotOfficialSignerNew, &currentOfficialNodeSigner)
	newOfficialNodeSigner := new(big.Int).Add(new(big.Int).SetBytes(libcommon.Address{}.Bytes()), shiftedBlockNumber)

	instruction.Storage[params.StakeManagerStorageV3] = map[libcommon.Hash]libcommon.Hash{
		// StakeManager reference
		libcommon.BigToHash(big.NewInt(1)): libcommon.BytesToHash(params.StakeManagerV3.Bytes()),

		// Super node configuration
		libcommon.BigToHash(big.NewInt(SLOT_SUPER_NODE_OLD)): libcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		libcommon.BigToHash(big.NewInt(SLOT_SUPER_NODE_NEW)): libcommon.BigToHash(superNodeAndValidBlock),

		libcommon.BigToHash(big.NewInt(SLOT_OFFICIAL_SIGNER_NEW)): libcommon.BigToHash(newOfficialNodeSigner),
		libcommon.BigToHash(big.NewInt(SLOT_OFFICIAL_SIGNER_OLD)): libcommon.BigToHash(currentOfficialNodeSigner.ToBig()),

		// New validator struct (6 slots)
		libcommon.BigToHash(new(big.Int).Add(nextValidatorSoltInt, big.NewInt(0))): libcommon.BigToHash(big.NewInt(0)),
		libcommon.BigToHash(new(big.Int).Add(nextValidatorSoltInt, big.NewInt(1))): libcommon.BigToHash(big.NewInt(0)),
		libcommon.BigToHash(new(big.Int).Add(nextValidatorSoltInt, big.NewInt(2))): libcommon.BigToHash(big.NewInt(0)),
		libcommon.BigToHash(new(big.Int).Add(nextValidatorSoltInt, big.NewInt(3))): libcommon.BigToHash(big.NewInt(0)),
		libcommon.BigToHash(new(big.Int).Add(nextValidatorSoltInt, big.NewInt(4))): libcommon.BytesToHash(params.SuperNodeAddress.Bytes()),
		libcommon.BigToHash(new(big.Int).Add(nextValidatorSoltInt, big.NewInt(5))): packValidatorShareAndRates(libcommon.Address{}, 1, 0, 0),

		// Validator list
		libcommon.BigToHash(big.NewInt(SLOT_VALIDATOR_LIST)): libcommon.BigToHash(new(big.Int).SetUint64(nextLength)),

		// Validator IDs mapping
		validatorIdsArrayLengthSlot:   libcommon.BigToHash(big.NewInt(1)),
		validatorIdsArrayElement0Slot: libcommon.BigToHash(new(big.Int).SetUint64(nextValidatorId)),

		// Minimal validator list
		libcommon.BigToHash(minimalValidatorListLengthSlot.Big()):                                                           libcommon.BigToHash(big.NewInt(int64(currentMinimalListLength + 1))),
		libcommon.BigToHash(new(big.Int).Add(minimalValidatorListDataSlotInt, big.NewInt(int64(currentMinimalListLength)))): libcommon.BytesToHash(params.SuperNodeAddress.Bytes()),
		libcommon.BigToHash(minimalValidatorListDataIndexes.Big()):                                                          libcommon.BigToHash(big.NewInt(int64(currentMinimalListLength + 1))),

		libcommon.BigToHash(big.NewInt(SLOT_OFFICIAL_AMOUNT)): libcommon.BigToHash(libcommon.Big0),

		// Pool amount
		libcommon.BigToHash(big.NewInt(SLOT_POOL_AMOUNT)): libcommon.BigToHash(big.NewInt(int64(newPoolAmount))),
	}

	tokenOwnerSlotLength := libcommon.BigToHash(big.NewInt(SLOT_TOKEN_OWNERS))
	tokenOwnerDataSlot := crypto.Keccak256Hash(tokenOwnerSlotLength.Bytes())
	tokenOwnerDataSlotInt := new(big.Int).Add(
		tokenOwnerDataSlot.Big(),
		big.NewInt(int64(nextValidatorId)*2))

	tokenIdHash := libcommon.BigToHash(big.NewInt(int64(nextValidatorId)))

	indexesSlot := getMappingSlot(tokenIdHash, SLOT_TOKEN_OWNERS_INDEXES)

	// Get existing holder token count
	var currentHolderTokenCount uint256.Int
	holderArraySlot := getMappingSlot(libcommon.BytesToHash(params.SuperNodeOwnerAddress.Bytes()), SLOT_HOLDER_TOKENS)
	state.GetState(params.NftContractV3, &holderArraySlot, &currentHolderTokenCount)
	newHolderTokenCount := currentHolderTokenCount.Uint64() + 1

	holderArraySlotInt := holderArraySlot.Big()
	holderArrayDataSlotInt := crypto.Keccak256Hash(holderArraySlot.Bytes()).Big()

	// Add new token at the current count position (0-indexed)
	newTokenPosition := currentHolderTokenCount.Uint64()
	holderTokenArrayElementSlot := new(big.Int).Add(holderArrayDataSlotInt, big.NewInt(int64(newTokenPosition)))

	holderIndexesSlot := GetHolderTokenIndexSlot(holderArraySlot, uint64(nextValidatorId))

	// Configure BKCValidatorSet storage
	instruction.Storage[params.BKCValidatorSetV3] = map[libcommon.Hash]libcommon.Hash{
		libcommon.BigToHash(big.NewInt(3)): libcommon.BytesToHash(params.StakeManagerStorageV3.Bytes()),
	}

	// Configure NFT contract storage
	instruction.Storage[params.NftContractV3] = map[libcommon.Hash]libcommon.Hash{
		// StakeManager reference
		libcommon.BigToHash(big.NewInt(12)): libcommon.BytesToHash(params.StakeManagerStorageV3.Bytes()),

		// Token supply
		libcommon.BigToHash(big.NewInt(SLOT_TOKEN_OWNERS)): libcommon.BigToHash(big.NewInt(int64(nextLength))),

		// Token ownership entries
		libcommon.BigToHash(new(big.Int).Add(tokenOwnerDataSlotInt, big.NewInt(0))): libcommon.BigToHash(big.NewInt(int64(nextValidatorId))),
		libcommon.BigToHash(new(big.Int).Add(tokenOwnerDataSlotInt, big.NewInt(1))): libcommon.BytesToHash(params.SuperNodeOwnerAddress.Bytes()),
		libcommon.BigToHash(indexesSlot.Big()):                                      libcommon.BigToHash(big.NewInt(int64(nextLength))),

		// Holder token array
		libcommon.BigToHash(holderArraySlotInt):          libcommon.BigToHash(big.NewInt(int64(newHolderTokenCount))),
		libcommon.BigToHash(holderTokenArrayElementSlot): libcommon.BigToHash(big.NewInt(int64(nextValidatorId))),
		libcommon.BigToHash(holderIndexesSlot.Big()):     libcommon.BigToHash(big.NewInt(int64(newHolderTokenCount))),
	}

	// Configure official node validator share
	var isOfficialData uint256.Int
	isOfficialSlot := libcommon.BigToHash(big.NewInt(SLOT_IS_OFFICAL_POOL))
	state.GetState(params.OfficialNodeValidatorShare, &isOfficialSlot, &isOfficialData)

	// Set isOfficalPool to false
	isOfficialDataBytes := isOfficialData.Bytes32()
	isOfficialDataBytes[11] = 0x00

	instruction.Storage[params.OfficialNodeValidatorShare] = map[libcommon.Hash]libcommon.Hash{
		libcommon.BigToHash(big.NewInt(SLOT_IS_OFFICAL_POOL)): libcommon.BytesToHash(isOfficialDataBytes[:]),
	}

	// Preserve existing StakeManagerV3 addresses
	var committeeAddress uint256.Int
	var callHelperAddress uint256.Int
	var transferRouterAddress uint256.Int
	var officialPoolStakerAddress uint256.Int
	var kkubAddress uint256.Int

	slot1 := libcommon.BigToHash(big.NewInt(1))
	slot2 := libcommon.BigToHash(big.NewInt(2))
	slot3 := libcommon.BigToHash(big.NewInt(3))
	slot4 := libcommon.BigToHash(big.NewInt(4))
	slot8 := libcommon.BigToHash(big.NewInt(8))

	state.GetState(params.StakeManagerV3, &slot1, &committeeAddress)
	state.GetState(params.StakeManagerV3, &slot2, &callHelperAddress)
	state.GetState(params.StakeManagerV3, &slot3, &transferRouterAddress)
	state.GetState(params.StakeManagerV3, &slot4, &officialPoolStakerAddress)
	state.GetState(params.StakeManagerV3, &slot8, &kkubAddress)

	// Configure StakeManagerV3 storage
	instruction.Storage[params.StakeManagerV3] = map[libcommon.Hash]libcommon.Hash{
		// Reentrancy guard
		libcommon.BigToHash(big.NewInt(0)): libcommon.BigToHash(big.NewInt(1)),

		// Preserved addresses
		libcommon.BigToHash(big.NewInt(1)): libcommon.BytesToHash(committeeAddress.Bytes()),
		libcommon.BigToHash(big.NewInt(2)): libcommon.BytesToHash(callHelperAddress.Bytes()),
		libcommon.BigToHash(big.NewInt(3)): libcommon.BytesToHash(transferRouterAddress.Bytes()),
		libcommon.BigToHash(big.NewInt(4)): libcommon.BytesToHash(officialPoolStakerAddress.Bytes()),
		libcommon.BigToHash(big.NewInt(8)): libcommon.BytesToHash(kkubAddress.Bytes()),

		// New contract references
		libcommon.BigToHash(big.NewInt(5)): libcommon.BytesToHash(params.StakeManagerStorageV3.Bytes()),
		libcommon.BigToHash(big.NewInt(6)): libcommon.BytesToHash(params.StakeManagerVaultV3.Bytes()),
		libcommon.BigToHash(big.NewInt(7)): libcommon.BytesToHash(params.NftContractV3.Bytes()),
	}

	return instruction, nil
}

func getMappingSlot(key libcommon.Hash, baseSlot uint64) libcommon.Hash {
	slotBytes := common.LeftPadBytes(new(big.Int).SetUint64(baseSlot).Bytes(), 32)
	data := append(key.Bytes(), slotBytes...)
	return crypto.Keccak256Hash(data)
}

func getValidatorSlot(baseSlot uint64, validatorId uint64) libcommon.Hash {
	slotBytes := common.LeftPadBytes(new(big.Int).SetUint64(baseSlot).Bytes(), 32)
	arrayBase := crypto.Keccak256Hash(slotBytes)

	arrayBaseInt := new(big.Int).SetBytes(arrayBase.Bytes())
	// Each Validator struct takes 6 storage slots (0-5)
	offset := new(big.Int).Mul(new(big.Int).SetUint64(validatorId), big.NewInt(6))
	validatorSlot := new(big.Int).Add(arrayBaseInt, offset)

	return libcommon.BytesToHash(validatorSlot.Bytes())
}

// GetHolderTokenIndexSlot returns the storage slot for _holderTokens[holder]._inner._indexes[tokenId]
// EnumerableSet.UintSet structure:
// - _values array at slot 0 (relative to set base)
// - _indexes mapping at slot 1 (relative to set base)
func GetHolderTokenIndexSlot(holderSetBaseSlot libcommon.Hash, tokenId uint64) libcommon.Hash {
	// _indexes is at holderSetBaseSlot + 1
	indexesSlotBase := new(big.Int).Add(holderSetBaseSlot.Big(), big.NewInt(1))

	// mapping(uint256 => uint256) lookup for tokenId
	data := make([]byte, 64)
	tokenIdBytes := common.LeftPadBytes(new(big.Int).SetUint64(tokenId).Bytes(), 32)
	copy(data[0:32], tokenIdBytes)
	indexesSlotBytes := common.LeftPadBytes(indexesSlotBase.Bytes(), 32)
	copy(data[32:64], indexesSlotBytes)

	return crypto.Keccak256Hash(data)
}

// GetTokenOwnersIndexSlot returns the storage slot for _indexes[tokenId] at SLOT_TOKEN_OWNERS_INDEXES
// This is for the separate _indexes mapping at slot 3
func GetTokenOwnersIndexSlot(tokenId uint64) libcommon.Hash {
	// mapping(uint256 => uint256) at SLOT_TOKEN_OWNERS_INDEXES (slot 3)
	data := make([]byte, 64)
	tokenIdBytes := common.LeftPadBytes(new(big.Int).SetUint64(tokenId).Bytes(), 32)
	copy(data[0:32], tokenIdBytes)
	slotBytes := common.LeftPadBytes(new(big.Int).SetUint64(SLOT_TOKEN_OWNERS_INDEXES).Bytes(), 32)
	copy(data[32:64], slotBytes)

	return crypto.Keccak256Hash(data)
}

// packValidatorShareAndRates packs validatorShareContract, status, and commission rates into one slot (DEPRECATED - now split into slot 5 and 6)
// Layout: validatorShareContract (160 bits) + status (8 bits) + infraCommissionRate (16 bits) + commissionRate (16 bits) = 200 bits
func packValidatorShareAndRates(addr libcommon.Address, status uint8, infraRate uint16, commissionRate uint16) libcommon.Hash {
	result := new(big.Int)
	// validatorShareContract address in lower 160 bits
	result.Or(result, new(big.Int).SetBytes(addr.Bytes()))
	// Status in next 8 bits (shift left by 160)
	result.Or(result, new(big.Int).Lsh(new(big.Int).SetUint64(uint64(status)), 160))
	// infraCommissionRate in next 16 bits (shift left by 168)
	result.Or(result, new(big.Int).Lsh(new(big.Int).SetUint64(uint64(infraRate)), 168))
	// commissionRate in next 16 bits (shift left by 184)
	result.Or(result, new(big.Int).Lsh(new(big.Int).SetUint64(uint64(commissionRate)), 184))
	return libcommon.BigToHash(result)
}
