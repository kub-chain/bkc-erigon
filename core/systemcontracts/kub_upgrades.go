package systemcontracts

import (
	libcommon "github.com/ledgerwatch/erigon-lib/common"
)

var (
	// Mainnet contracts
	kubPosMainnetStakeManagerRecords []libcommon.CodeRecord = []libcommon.CodeRecord{
		{
			BlockNumber: 14115624,
			CodeHash:    libcommon.HexToHash("2597d35349905301461edd94aad876c701cadd298a3aa7846c1b565d681312e4"), // Chanphraya
		},
		{
			BlockNumber: 25677934,
			CodeHash:    libcommon.HexToHash("b55c1bd656d20f94ce19b7642196e8394042ddaf2135f24f3f95ab59f032acef"), // Lausanne
		},
	}
	kubPosMainnetSlashManagerRecords []libcommon.CodeRecord = []libcommon.CodeRecord{
		{
			BlockNumber: 14115624,
			CodeHash:    libcommon.HexToHash("8dffdb975547abadb304be30234bdac2e7f9fcf008b290051e2edcd669f4e0df"), // Chanphraya
		},
		{
			BlockNumber: 25677934,
			CodeHash:    libcommon.HexToHash("b0e67b0e7cceb1d29d1ecad739da6248020d9b35fbfe4674c1cf46a3574e6e48"), // Lausanne
		},
	}
	kubPosMainnetStakeManagerStorageRecords []libcommon.CodeRecord = []libcommon.CodeRecord{
		{
			BlockNumber: 14115624,
			CodeHash:    libcommon.HexToHash("456e4b696e7ee8933825f6888a0c990c243dc45632286468d150d3de11ab0c52"), // Chanphraya
		},
		{
			BlockNumber: 25677934,
			CodeHash:    libcommon.HexToHash("4144eec9835f45dd46845dd1920ba433455023a799740800dc1c0adc541fd1c3"), // Lausanne
		},
	}
	// Testnet contracts
	kubPosTestnetStakeManagerRecords []libcommon.CodeRecord = []libcommon.CodeRecord{
		{
			BlockNumber: 11712666,
			CodeHash:    libcommon.HexToHash("57a9dd2e75df9f1949c2b6ef36498e4c41af66b98b871bd9937c46a5b0de6ffe"), // Chanphraya
		},
		{
			BlockNumber: 22835041,
			CodeHash:    libcommon.HexToHash("b55c1bd656d20f94ce19b7642196e8394042ddaf2135f24f3f95ab59f032acef"), // Lausanne
		},
	}
	kubPosTestnetSlashManagerRecords []libcommon.CodeRecord = []libcommon.CodeRecord{
		{
			BlockNumber: 11712666,
			CodeHash:    libcommon.HexToHash("577b68832c6778d4653a1042aa6b358d24364004958811db0664870583dd1c30"), // Chanphraya
		},
		{
			BlockNumber: 22835041,
			CodeHash:    libcommon.HexToHash("b0e67b0e7cceb1d29d1ecad739da6248020d9b35fbfe4674c1cf46a3574e6e48"), // Lausanne
		},
	}
	kubPosTestnetStakeManagerStorageRecords []libcommon.CodeRecord = []libcommon.CodeRecord{
		{
			BlockNumber: 11712666,
			CodeHash:    libcommon.HexToHash("7def64ba6b2c753932221a13131d177977d8668f802c000d00edf165b9140b0d"), // Chanphraya
		},
		{
			BlockNumber: 22835041,
			CodeHash:    libcommon.HexToHash("4144eec9835f45dd46845dd1920ba433455023a799740800dc1c0adc541fd1c3"), // Lausanne
		},
	}
)
