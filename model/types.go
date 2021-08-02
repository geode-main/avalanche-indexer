package model

const (
	// Asset types
	AssetTypeFixed    = "fixed_cap"
	AssetTypeVariable = "variable_cap"
	AssetTypeNFT      = "nft"

	// PVM block types
	BlockTypeProposal = "proposal"
	BlockTypeStandard = "standard"
	BlockTypeAtomic   = "atomic"
	BlockTypeCommit   = "commit"
	BlockTypeAbort    = "abort"
	BlockTypeEvm      = "evm"

	// Transaction statuses
	TxStatusAccepted = "accepted"
	TxStatusRejected = "rejected"
	TxStatusReverted = "reverted"

	// PVM transaction types
	TxTypeCreateChain        = "p_create_chain"
	TxTypeCreateSubnet       = "p_create_subnet"
	TxTypeAddSubnetValidator = "p_add_subnet_validator"
	TxTypeAdvanceTime        = "p_advance_time"
	TxTypeRewardValidator    = "p_reward_validator"
	TxTypeAddValidator       = "p_add_validator"
	TxTypeAddDelegator       = "p_add_delegator"
	TxTypePImport            = "p_import"
	TxTypePExport            = "p_export"

	// AVM transaction types
	TxTypeBase        = "x_base"
	TxTypeXImport     = "x_import"
	TxTypeXExport     = "x_export"
	TxTypeCreateAsset = "x_create_asset"
	TxTypeOperation   = "x_operation"

	// EMV transaction types
	TxTypeAtomicExport = "c_atomic_export"
	TxTypeAtomicImport = "c_atomic_import"
	TxTypeEvm          = "c_evm"

	// Output types
	OutTypeTransfer      = "transfer"
	OutTypeStakeableLock = "stakeable_lock"
	OutTypeReward        = "reward"
	OutTypeMint          = "mint"
	OutTypeNftMint       = "nft_mint"
	OutTypeNftTransfer   = "nft_transfer"

	// Reward types
	RewardTypeValidator = "validator"
	RewardTypeDelegator = "delegator"

	// Event scopes
	EventScopeStaking = "staking"
	EventScopeRewards = "rewards"
	EventScopeNetwork = "network"

	// Event Item Types
	EventItemTypeValidator = "validator"
	EventItemTypeDelegator = "delegator"

	// Event types
	EventTypeValidatorAdded             = "validator_added"
	EventTypeValidatorFinished          = "validator_finished"
	EventTypeValidatorCommissionChanged = "validator_commission_changed"
	EventTypeDelegatorAdded             = "delegator_added"
	EventTypeDelegatorFinished          = "delegator_finished"
	EventTypeSubnetValidatorAdded       = "subnet_validator_added"
)

var (
	TransactionTypes = []string{
		TxTypeCreateChain,
		TxTypeCreateSubnet,
		TxTypeAddSubnetValidator,
		TxTypeAdvanceTime,
		TxTypeRewardValidator,
		TxTypeAddValidator,
		TxTypeAddDelegator,
		TxTypePImport,
		TxTypePExport,
		TxTypeBase,
		TxTypeXImport,
		TxTypeXExport,
		TxTypeCreateAsset,
		TxTypeOperation,
		TxTypeAtomicExport,
		TxTypeAtomicImport,
		TxTypeEvm,
	}
)
