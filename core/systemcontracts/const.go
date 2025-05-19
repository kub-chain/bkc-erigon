package systemcontracts

import (
	libcommon "github.com/ledgerwatch/erigon-lib/common"
)

var (
	// genesis contracts
	ValidatorContract          = libcommon.HexToAddress("0x0000000000000000000000000000000000001000")
	SlashContract              = libcommon.HexToAddress("0x0000000000000000000000000000000000001001")
	SystemRewardContract       = libcommon.HexToAddress("0x0000000000000000000000000000000000001002")
	LightClientContract        = libcommon.HexToAddress("0x0000000000000000000000000000000000001003")
	TokenHubContract           = libcommon.HexToAddress("0x0000000000000000000000000000000000001004")
	RelayerIncentivizeContract = libcommon.HexToAddress("0x0000000000000000000000000000000000001005")
	RelayerHubContract         = libcommon.HexToAddress("0x0000000000000000000000000000000000001006")
	GovHubContract             = libcommon.HexToAddress("0x0000000000000000000000000000000000001007")
	TokenManagerContract       = libcommon.HexToAddress("0x0000000000000000000000000000000000001008")
	MaticTokenContract         = libcommon.HexToAddress("0x0000000000000000000000000000000000001010")
	CrossChainContract         = libcommon.HexToAddress("0x0000000000000000000000000000000000002000")
	StakingContract            = libcommon.HexToAddress("0x0000000000000000000000000000000000002001")
)

var (
	// kub pos contracts
	KubPosStakeManagerMainnet        = libcommon.HexToAddress("0x443502b3F7C0934576F49CDa084f78640f56A80F")
	KubPosSlashManagerMainnet        = libcommon.HexToAddress("0xEF6F6c6fdaEAc0326FFE1413D7d7CCAA7B56b753")
	KubPosStakeManagerStorageMainnet = libcommon.HexToAddress("0xFd98aac1Fbc57e6BC16A167452DBA7af2B6a4c0d")
	KubPosStakeManagerTestnet        = libcommon.HexToAddress("0xd4a1478020092e624db300714E41b3D39fc6313f")
	KubPosSlashManagerTestnet        = libcommon.HexToAddress("0x0F717cAC6655D95Fbc6Ee60D06f03bD8ae867e4F")
	KubPosStakeManagerStorageTestnet = libcommon.HexToAddress("0xC339c3abc4dB5a730C6610058535dBEC52763dD6")
)
