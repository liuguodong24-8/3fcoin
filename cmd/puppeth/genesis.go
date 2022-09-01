// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"errors"
	"math"
	"math/big"
	"strings"

	"github.com/fff-chain/3f-chain/core/common"
	math2 "github.com/fff-chain/3f-chain/core/common/math"
	"github.com/fff-chain/3f-chain/core/consensus/ethash"
	"github.com/fff-chain/3f-chain/core/core"
	"github.com/fff-chain/3f-chain/core/core/types"
	"github.com/fff-chain/3f-chain/core/params"
)

// alethGenesisSpec represents the genesis specification format used by the
// C++ Ethereum implementation.
type alethGenesisSpec struct {
	SealEngine string `json:"sealEngine"`
	Params     struct {
		AccountStartNonce          math2.HexOrDecimal64   `json:"accountStartNonce"`
		MaximumExtraDataSize       common.Uint64          `json:"maximumExtraDataSize"`
		HomesteadForkBlock         *common.Big            `json:"homesteadForkBlock,omitempty"`
		DaoHardforkBlock           math2.HexOrDecimal64   `json:"daoHardforkBlock"`
		EIP150ForkBlock            *common.Big            `json:"EIP150ForkBlock,omitempty"`
		EIP158ForkBlock            *common.Big            `json:"EIP158ForkBlock,omitempty"`
		ByzantiumForkBlock         *common.Big            `json:"byzantiumForkBlock,omitempty"`
		ConstantinopleForkBlock    *common.Big            `json:"constantinopleForkBlock,omitempty"`
		ConstantinopleFixForkBlock *common.Big            `json:"constantinopleFixForkBlock,omitempty"`
		IstanbulForkBlock          *common.Big            `json:"istanbulForkBlock,omitempty"`
		MinGasLimit                common.Uint64          `json:"minGasLimit"`
		MaxGasLimit                common.Uint64          `json:"maxGasLimit"`
		TieBreakingGas             bool                   `json:"tieBreakingGas"`
		GasLimitBoundDivisor       math2.HexOrDecimal64   `json:"gasLimitBoundDivisor"`
		MinimumDifficulty          *common.Big            `json:"minimumDifficulty"`
		DifficultyBoundDivisor     *math2.HexOrDecimal256 `json:"difficultyBoundDivisor"`
		DurationLimit              *math2.HexOrDecimal256 `json:"durationLimit"`
		BlockReward                *common.Big            `json:"blockReward"`
		NetworkID                  common.Uint64          `json:"networkID"`
		ChainID                    common.Uint64          `json:"chainID"`
		AllowFutureBlocks          bool                   `json:"allowFutureBlocks"`
	} `json:"params"`

	Genesis struct {
		Nonce      types.BlockNonce `json:"nonce"`
		Difficulty *common.Big      `json:"difficulty"`
		MixHash    common.Hash      `json:"mixHash"`
		Author     common.Address   `json:"author"`
		Timestamp  common.Uint64    `json:"timestamp"`
		ParentHash common.Hash      `json:"parentHash"`
		ExtraData  common.Bytes     `json:"extraData"`
		GasLimit   common.Uint64    `json:"gasLimit"`
	} `json:"genesis"`

	Accounts map[common.Address]*alethGenesisSpecAccount `json:"accounts"`
}

// alethGenesisSpecAccount is the prefunded genesis account and/or precompiled
// contract definition.
type alethGenesisSpecAccount struct {
	Balance     *math2.HexOrDecimal256   `json:"balance,omitempty"`
	Nonce       uint64                   `json:"nonce,omitempty"`
	Precompiled *alethGenesisSpecBuiltin `json:"precompiled,omitempty"`
}

// alethGenesisSpecBuiltin is the precompiled contract definition.
type alethGenesisSpecBuiltin struct {
	Name          string                         `json:"name,omitempty"`
	StartingBlock *common.Big                    `json:"startingBlock,omitempty"`
	Linear        *alethGenesisSpecLinearPricing `json:"linear,omitempty"`
}

type alethGenesisSpecLinearPricing struct {
	Base uint64 `json:"base"`
	Word uint64 `json:"word"`
}

// newAlethGenesisSpec converts a go-ethereum genesis block into a Aleth-specific
// chain specification format.
func newAlethGenesisSpec(network string, genesis *core.Genesis) (*alethGenesisSpec, error) {
	// Only ethash is currently supported between go-ethereum and aleth
	if genesis.Config.Ethash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	// Reconstruct the chain spec in Aleth format
	spec := &alethGenesisSpec{
		SealEngine: "Ethash",
	}
	// Some defaults
	spec.Params.AccountStartNonce = 0
	spec.Params.TieBreakingGas = false
	spec.Params.AllowFutureBlocks = false

	// Dao hardfork block is a special one. The fork block is listed as 0 in the
	// config but aleth will sync with ETC clients up until the actual dao hard
	// fork block.
	spec.Params.DaoHardforkBlock = 0

	if num := genesis.Config.HomesteadBlock; num != nil {
		spec.Params.HomesteadForkBlock = (*common.Big)(num)
	}
	if num := genesis.Config.EIP150Block; num != nil {
		spec.Params.EIP150ForkBlock = (*common.Big)(num)
	}
	if num := genesis.Config.EIP158Block; num != nil {
		spec.Params.EIP158ForkBlock = (*common.Big)(num)
	}
	if num := genesis.Config.ByzantiumBlock; num != nil {
		spec.Params.ByzantiumForkBlock = (*common.Big)(num)
	}
	if num := genesis.Config.ConstantinopleBlock; num != nil {
		spec.Params.ConstantinopleForkBlock = (*common.Big)(num)
	}
	if num := genesis.Config.PetersburgBlock; num != nil {
		spec.Params.ConstantinopleFixForkBlock = (*common.Big)(num)
	}
	if num := genesis.Config.IstanbulBlock; num != nil {
		spec.Params.IstanbulForkBlock = (*common.Big)(num)
	}
	spec.Params.NetworkID = (common.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.ChainID = (common.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.MaximumExtraDataSize = (common.Uint64)(params.MaximumExtraDataSize)
	spec.Params.MinGasLimit = (common.Uint64)(params.MinGasLimit)
	spec.Params.MaxGasLimit = (common.Uint64)(math.MaxInt64)
	spec.Params.MinimumDifficulty = (*common.Big)(params.MinimumDifficulty)
	spec.Params.DifficultyBoundDivisor = (*math2.HexOrDecimal256)(params.DifficultyBoundDivisor)
	spec.Params.GasLimitBoundDivisor = (math2.HexOrDecimal64)(params.GasLimitBoundDivisor)
	spec.Params.DurationLimit = (*math2.HexOrDecimal256)(params.DurationLimit)
	spec.Params.BlockReward = (*common.Big)(ethash.FrontierBlockReward)

	spec.Genesis.Nonce = types.EncodeNonce(genesis.Nonce)
	spec.Genesis.MixHash = genesis.Mixhash
	spec.Genesis.Difficulty = (*common.Big)(genesis.Difficulty)
	spec.Genesis.Author = genesis.Coinbase
	spec.Genesis.Timestamp = (common.Uint64)(genesis.Timestamp)
	spec.Genesis.ParentHash = genesis.ParentHash
	spec.Genesis.ExtraData = genesis.ExtraData
	spec.Genesis.GasLimit = (common.Uint64)(genesis.GasLimit)

	for address, account := range genesis.Alloc {
		spec.setAccount(address, account)
	}

	spec.setPrecompile(1, &alethGenesisSpecBuiltin{Name: "ecrecover",
		Linear: &alethGenesisSpecLinearPricing{Base: 3000}})
	spec.setPrecompile(2, &alethGenesisSpecBuiltin{Name: "sha256",
		Linear: &alethGenesisSpecLinearPricing{Base: 60, Word: 12}})
	spec.setPrecompile(3, &alethGenesisSpecBuiltin{Name: "ripemd160",
		Linear: &alethGenesisSpecLinearPricing{Base: 600, Word: 120}})
	spec.setPrecompile(4, &alethGenesisSpecBuiltin{Name: "identity",
		Linear: &alethGenesisSpecLinearPricing{Base: 15, Word: 3}})
	if genesis.Config.ByzantiumBlock != nil {
		spec.setPrecompile(5, &alethGenesisSpecBuiltin{Name: "modexp",
			StartingBlock: (*common.Big)(genesis.Config.ByzantiumBlock)})
		spec.setPrecompile(6, &alethGenesisSpecBuiltin{Name: "alt_bn128_G1_add",
			StartingBlock: (*common.Big)(genesis.Config.ByzantiumBlock),
			Linear:        &alethGenesisSpecLinearPricing{Base: 500}})
		spec.setPrecompile(7, &alethGenesisSpecBuiltin{Name: "alt_bn128_G1_mul",
			StartingBlock: (*common.Big)(genesis.Config.ByzantiumBlock),
			Linear:        &alethGenesisSpecLinearPricing{Base: 40000}})
		spec.setPrecompile(8, &alethGenesisSpecBuiltin{Name: "alt_bn128_pairing_product",
			StartingBlock: (*common.Big)(genesis.Config.ByzantiumBlock)})
	}
	if genesis.Config.IstanbulBlock != nil {
		if genesis.Config.ByzantiumBlock == nil {
			return nil, errors.New("invalid genesis, istanbul fork is enabled while byzantium is not")
		}
		spec.setPrecompile(6, &alethGenesisSpecBuiltin{
			Name:          "alt_bn128_G1_add",
			StartingBlock: (*common.Big)(genesis.Config.ByzantiumBlock),
		}) // Aleth hardcoded the gas policy
		spec.setPrecompile(7, &alethGenesisSpecBuiltin{
			Name:          "alt_bn128_G1_mul",
			StartingBlock: (*common.Big)(genesis.Config.ByzantiumBlock),
		}) // Aleth hardcoded the gas policy
		spec.setPrecompile(9, &alethGenesisSpecBuiltin{
			Name:          "blake2_compression",
			StartingBlock: (*common.Big)(genesis.Config.IstanbulBlock),
		})
	}
	return spec, nil
}

func (spec *alethGenesisSpec) setPrecompile(address byte, data *alethGenesisSpecBuiltin) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.Address]*alethGenesisSpecAccount)
	}
	addr := common.Address(common.BytesToAddress([]byte{address}))
	if _, exist := spec.Accounts[addr]; !exist {
		spec.Accounts[addr] = &alethGenesisSpecAccount{}
	}
	spec.Accounts[addr].Precompiled = data
}

func (spec *alethGenesisSpec) setAccount(address common.Address, account core.GenesisAccount) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.Address]*alethGenesisSpecAccount)
	}

	a, exist := spec.Accounts[common.Address(address)]
	if !exist {
		a = &alethGenesisSpecAccount{}
		spec.Accounts[common.Address(address)] = a
	}
	a.Balance = (*math2.HexOrDecimal256)(account.Balance)
	a.Nonce = account.Nonce

}

// parityChainSpec is the chain specification format used by Parity.
type parityChainSpec struct {
	Name    string `json:"name"`
	Datadir string `json:"dataDir"`
	Engine  struct {
		Ethash struct {
			Params struct {
				MinimumDifficulty      *common.Big       `json:"minimumDifficulty"`
				DifficultyBoundDivisor *common.Big       `json:"difficultyBoundDivisor"`
				DurationLimit          *common.Big       `json:"durationLimit"`
				BlockReward            map[string]string `json:"blockReward"`
				DifficultyBombDelays   map[string]string `json:"difficultyBombDelays"`
				HomesteadTransition    common.Uint64     `json:"homesteadTransition"`
				EIP100bTransition      common.Uint64     `json:"eip100bTransition"`
			} `json:"params"`
		} `json:"Ethash"`
	} `json:"engine"`

	Params struct {
		AccountStartNonce         common.Uint64        `json:"accountStartNonce"`
		MaximumExtraDataSize      common.Uint64        `json:"maximumExtraDataSize"`
		MinGasLimit               common.Uint64        `json:"minGasLimit"`
		GasLimitBoundDivisor      math2.HexOrDecimal64 `json:"gasLimitBoundDivisor"`
		NetworkID                 common.Uint64        `json:"networkID"`
		ChainID                   common.Uint64        `json:"chainID"`
		MaxCodeSize               common.Uint64        `json:"maxCodeSize"`
		MaxCodeSizeTransition     common.Uint64        `json:"maxCodeSizeTransition"`
		EIP98Transition           common.Uint64        `json:"eip98Transition"`
		EIP150Transition          common.Uint64        `json:"eip150Transition"`
		EIP160Transition          common.Uint64        `json:"eip160Transition"`
		EIP161abcTransition       common.Uint64        `json:"eip161abcTransition"`
		EIP161dTransition         common.Uint64        `json:"eip161dTransition"`
		EIP155Transition          common.Uint64        `json:"eip155Transition"`
		EIP140Transition          common.Uint64        `json:"eip140Transition"`
		EIP211Transition          common.Uint64        `json:"eip211Transition"`
		EIP214Transition          common.Uint64        `json:"eip214Transition"`
		EIP658Transition          common.Uint64        `json:"eip658Transition"`
		EIP145Transition          common.Uint64        `json:"eip145Transition"`
		EIP1014Transition         common.Uint64        `json:"eip1014Transition"`
		EIP1052Transition         common.Uint64        `json:"eip1052Transition"`
		EIP1283Transition         common.Uint64        `json:"eip1283Transition"`
		EIP1283DisableTransition  common.Uint64        `json:"eip1283DisableTransition"`
		EIP1283ReenableTransition common.Uint64        `json:"eip1283ReenableTransition"`
		EIP1344Transition         common.Uint64        `json:"eip1344Transition"`
		EIP1884Transition         common.Uint64        `json:"eip1884Transition"`
		EIP2028Transition         common.Uint64        `json:"eip2028Transition"`
	} `json:"params"`

	Genesis struct {
		Seal struct {
			Ethereum struct {
				Nonce   types.BlockNonce `json:"nonce"`
				MixHash common.Bytes     `json:"mixHash"`
			} `json:"ethereum"`
		} `json:"seal"`

		Difficulty *common.Big    `json:"difficulty"`
		Author     common.Address `json:"author"`
		Timestamp  common.Uint64  `json:"timestamp"`
		ParentHash common.Hash    `json:"parentHash"`
		ExtraData  common.Bytes   `json:"extraData"`
		GasLimit   common.Uint64  `json:"gasLimit"`
	} `json:"genesis"`

	Nodes    []string                                   `json:"nodes"`
	Accounts map[common.Address]*parityChainSpecAccount `json:"accounts"`
}

// parityChainSpecAccount is the prefunded genesis account and/or precompiled
// contract definition.
type parityChainSpecAccount struct {
	Balance math2.HexOrDecimal256   `json:"balance"`
	Nonce   math2.HexOrDecimal64    `json:"nonce,omitempty"`
	Builtin *parityChainSpecBuiltin `json:"builtin,omitempty"`
}

// parityChainSpecBuiltin is the precompiled contract definition.
type parityChainSpecBuiltin struct {
	Name       string      `json:"name"`                  // Each builtin should has it own name
	Pricing    interface{} `json:"pricing"`               // Each builtin should has it own price strategy
	ActivateAt *common.Big `json:"activate_at,omitempty"` // ActivateAt can't be omitted if empty, default means no fork
}

// parityChainSpecPricing represents the different pricing models that builtin
// contracts might advertise using.
type parityChainSpecPricing struct {
	Linear *parityChainSpecLinearPricing `json:"linear,omitempty"`
	ModExp *parityChainSpecModExpPricing `json:"modexp,omitempty"`

	// Before the https://github.com/paritytech/parity-ethereum/pull/11039,
	// Parity uses this format to config bn pairing price policy.
	AltBnPairing *parityChainSepcAltBnPairingPricing `json:"alt_bn128_pairing,omitempty"`

	// Blake2F is the price per round of Blake2 compression
	Blake2F *parityChainSpecBlakePricing `json:"blake2_f,omitempty"`
}

type parityChainSpecLinearPricing struct {
	Base uint64 `json:"base"`
	Word uint64 `json:"word"`
}

type parityChainSpecModExpPricing struct {
	Divisor uint64 `json:"divisor"`
}

// parityChainSpecAltBnConstOperationPricing defines the price
// policy for bn const operation(used after istanbul)
type parityChainSpecAltBnConstOperationPricing struct {
	Price uint64 `json:"price"`
}

// parityChainSepcAltBnPairingPricing defines the price policy
// for bn pairing.
type parityChainSepcAltBnPairingPricing struct {
	Base uint64 `json:"base"`
	Pair uint64 `json:"pair"`
}

// parityChainSpecBlakePricing defines the price policy for blake2 f
// compression.
type parityChainSpecBlakePricing struct {
	GasPerRound uint64 `json:"gas_per_round"`
}

type parityChainSpecAlternativePrice struct {
	AltBnConstOperationPrice *parityChainSpecAltBnConstOperationPricing `json:"alt_bn128_const_operations,omitempty"`
	AltBnPairingPrice        *parityChainSepcAltBnPairingPricing        `json:"alt_bn128_pairing,omitempty"`
}

// parityChainSpecVersionedPricing represents a single version price policy.
type parityChainSpecVersionedPricing struct {
	Price *parityChainSpecAlternativePrice `json:"price,omitempty"`
	Info  string                           `json:"info,omitempty"`
}

// newParityChainSpec converts a go-ethereum genesis block into a Parity specific
// chain specification format.
func newParityChainSpec(network string, genesis *core.Genesis, bootnodes []string) (*parityChainSpec, error) {
	// Only ethash is currently supported between go-ethereum and Parity
	if genesis.Config.Ethash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	// Reconstruct the chain spec in Parity's format
	spec := &parityChainSpec{
		Name:    network,
		Nodes:   bootnodes,
		Datadir: strings.ToLower(network),
	}
	spec.Engine.Ethash.Params.BlockReward = make(map[string]string)
	spec.Engine.Ethash.Params.DifficultyBombDelays = make(map[string]string)
	// Frontier
	spec.Engine.Ethash.Params.MinimumDifficulty = (*common.Big)(params.MinimumDifficulty)
	spec.Engine.Ethash.Params.DifficultyBoundDivisor = (*common.Big)(params.DifficultyBoundDivisor)
	spec.Engine.Ethash.Params.DurationLimit = (*common.Big)(params.DurationLimit)
	spec.Engine.Ethash.Params.BlockReward["0x0"] = common.EncodeBig(ethash.FrontierBlockReward)

	// Homestead
	spec.Engine.Ethash.Params.HomesteadTransition = common.Uint64(genesis.Config.HomesteadBlock.Uint64())

	// Tangerine Whistle : 150
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-608.md
	spec.Params.EIP150Transition = common.Uint64(genesis.Config.EIP150Block.Uint64())

	// Spurious Dragon: 155, 160, 161, 170
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-607.md
	spec.Params.EIP155Transition = common.Uint64(genesis.Config.EIP155Block.Uint64())
	spec.Params.EIP160Transition = common.Uint64(genesis.Config.EIP155Block.Uint64())
	spec.Params.EIP161abcTransition = common.Uint64(genesis.Config.EIP158Block.Uint64())
	spec.Params.EIP161dTransition = common.Uint64(genesis.Config.EIP158Block.Uint64())

	// Byzantium
	if num := genesis.Config.ByzantiumBlock; num != nil {
		spec.setByzantium(num)
	}
	// Constantinople
	if num := genesis.Config.ConstantinopleBlock; num != nil {
		spec.setConstantinople(num)
	}
	// ConstantinopleFix (remove eip-1283)
	if num := genesis.Config.PetersburgBlock; num != nil {
		spec.setConstantinopleFix(num)
	}
	// Istanbul
	if num := genesis.Config.IstanbulBlock; num != nil {
		spec.setIstanbul(num)
	}
	spec.Params.MaximumExtraDataSize = (common.Uint64)(params.MaximumExtraDataSize)
	spec.Params.MinGasLimit = (common.Uint64)(params.MinGasLimit)
	spec.Params.GasLimitBoundDivisor = (math2.HexOrDecimal64)(params.GasLimitBoundDivisor)
	spec.Params.NetworkID = (common.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.ChainID = (common.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.MaxCodeSize = params.MaxCodeSize
	// geth has it set from zero
	spec.Params.MaxCodeSizeTransition = 0

	// Disable this one
	spec.Params.EIP98Transition = math.MaxInt64

	spec.Genesis.Seal.Ethereum.Nonce = types.EncodeNonce(genesis.Nonce)
	spec.Genesis.Seal.Ethereum.MixHash = genesis.Mixhash[:]
	spec.Genesis.Difficulty = (*common.Big)(genesis.Difficulty)
	spec.Genesis.Author = genesis.Coinbase
	spec.Genesis.Timestamp = (common.Uint64)(genesis.Timestamp)
	spec.Genesis.ParentHash = genesis.ParentHash
	spec.Genesis.ExtraData = genesis.ExtraData
	spec.Genesis.GasLimit = (common.Uint64)(genesis.GasLimit)

	spec.Accounts = make(map[common.Address]*parityChainSpecAccount)
	for address, account := range genesis.Alloc {
		bal := math2.HexOrDecimal256(*account.Balance)

		spec.Accounts[common.Address(address)] = &parityChainSpecAccount{
			Balance: bal,
			Nonce:   math2.HexOrDecimal64(account.Nonce),
		}
	}
	spec.setPrecompile(1, &parityChainSpecBuiltin{Name: "ecrecover",
		Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 3000}}})

	spec.setPrecompile(2, &parityChainSpecBuiltin{
		Name: "sha256", Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 60, Word: 12}},
	})
	spec.setPrecompile(3, &parityChainSpecBuiltin{
		Name: "ripemd160", Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 600, Word: 120}},
	})
	spec.setPrecompile(4, &parityChainSpecBuiltin{
		Name: "identity", Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 15, Word: 3}},
	})
	if genesis.Config.ByzantiumBlock != nil {
		spec.setPrecompile(5, &parityChainSpecBuiltin{
			Name:       "modexp",
			ActivateAt: (*common.Big)(genesis.Config.ByzantiumBlock),
			Pricing: &parityChainSpecPricing{
				ModExp: &parityChainSpecModExpPricing{Divisor: 20},
			},
		})
		spec.setPrecompile(6, &parityChainSpecBuiltin{
			Name:       "alt_bn128_add",
			ActivateAt: (*common.Big)(genesis.Config.ByzantiumBlock),
			Pricing: &parityChainSpecPricing{
				Linear: &parityChainSpecLinearPricing{Base: 500, Word: 0},
			},
		})
		spec.setPrecompile(7, &parityChainSpecBuiltin{
			Name:       "alt_bn128_mul",
			ActivateAt: (*common.Big)(genesis.Config.ByzantiumBlock),
			Pricing: &parityChainSpecPricing{
				Linear: &parityChainSpecLinearPricing{Base: 40000, Word: 0},
			},
		})
		spec.setPrecompile(8, &parityChainSpecBuiltin{
			Name:       "alt_bn128_pairing",
			ActivateAt: (*common.Big)(genesis.Config.ByzantiumBlock),
			Pricing: &parityChainSpecPricing{
				AltBnPairing: &parityChainSepcAltBnPairingPricing{Base: 100000, Pair: 80000},
			},
		})
	}
	if genesis.Config.IstanbulBlock != nil {
		if genesis.Config.ByzantiumBlock == nil {
			return nil, errors.New("invalid genesis, istanbul fork is enabled while byzantium is not")
		}
		spec.setPrecompile(6, &parityChainSpecBuiltin{
			Name:       "alt_bn128_add",
			ActivateAt: (*common.Big)(genesis.Config.ByzantiumBlock),
			Pricing: map[*common.Big]*parityChainSpecVersionedPricing{
				(*common.Big)(big.NewInt(0)): {
					Price: &parityChainSpecAlternativePrice{
						AltBnConstOperationPrice: &parityChainSpecAltBnConstOperationPricing{Price: 500},
					},
				},
				(*common.Big)(genesis.Config.IstanbulBlock): {
					Price: &parityChainSpecAlternativePrice{
						AltBnConstOperationPrice: &parityChainSpecAltBnConstOperationPricing{Price: 150},
					},
				},
			},
		})
		spec.setPrecompile(7, &parityChainSpecBuiltin{
			Name:       "alt_bn128_mul",
			ActivateAt: (*common.Big)(genesis.Config.ByzantiumBlock),
			Pricing: map[*common.Big]*parityChainSpecVersionedPricing{
				(*common.Big)(big.NewInt(0)): {
					Price: &parityChainSpecAlternativePrice{
						AltBnConstOperationPrice: &parityChainSpecAltBnConstOperationPricing{Price: 40000},
					},
				},
				(*common.Big)(genesis.Config.IstanbulBlock): {
					Price: &parityChainSpecAlternativePrice{
						AltBnConstOperationPrice: &parityChainSpecAltBnConstOperationPricing{Price: 6000},
					},
				},
			},
		})
		spec.setPrecompile(8, &parityChainSpecBuiltin{
			Name:       "alt_bn128_pairing",
			ActivateAt: (*common.Big)(genesis.Config.ByzantiumBlock),
			Pricing: map[*common.Big]*parityChainSpecVersionedPricing{
				(*common.Big)(big.NewInt(0)): {
					Price: &parityChainSpecAlternativePrice{
						AltBnPairingPrice: &parityChainSepcAltBnPairingPricing{Base: 100000, Pair: 80000},
					},
				},
				(*common.Big)(genesis.Config.IstanbulBlock): {
					Price: &parityChainSpecAlternativePrice{
						AltBnPairingPrice: &parityChainSepcAltBnPairingPricing{Base: 45000, Pair: 34000},
					},
				},
			},
		})
		spec.setPrecompile(9, &parityChainSpecBuiltin{
			Name:       "blake2_f",
			ActivateAt: (*common.Big)(genesis.Config.IstanbulBlock),
			Pricing: &parityChainSpecPricing{
				Blake2F: &parityChainSpecBlakePricing{GasPerRound: 1},
			},
		})
	}
	return spec, nil
}

func (spec *parityChainSpec) setPrecompile(address byte, data *parityChainSpecBuiltin) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.Address]*parityChainSpecAccount)
	}
	a := common.Address(common.BytesToAddress([]byte{address}))
	if _, exist := spec.Accounts[a]; !exist {
		spec.Accounts[a] = &parityChainSpecAccount{}
	}
	spec.Accounts[a].Builtin = data
}

func (spec *parityChainSpec) setByzantium(num *big.Int) {
	spec.Engine.Ethash.Params.BlockReward[common.EncodeBig(num)] = common.EncodeBig(ethash.ByzantiumBlockReward)
	spec.Engine.Ethash.Params.DifficultyBombDelays[common.EncodeBig(num)] = common.EncodeUint64(3000000)
	n := common.Uint64(num.Uint64())
	spec.Engine.Ethash.Params.EIP100bTransition = n
	spec.Params.EIP140Transition = n
	spec.Params.EIP211Transition = n
	spec.Params.EIP214Transition = n
	spec.Params.EIP658Transition = n
}

func (spec *parityChainSpec) setConstantinople(num *big.Int) {
	spec.Engine.Ethash.Params.BlockReward[common.EncodeBig(num)] = common.EncodeBig(ethash.ConstantinopleBlockReward)
	spec.Engine.Ethash.Params.DifficultyBombDelays[common.EncodeBig(num)] = common.EncodeUint64(2000000)
	n := common.Uint64(num.Uint64())
	spec.Params.EIP145Transition = n
	spec.Params.EIP1014Transition = n
	spec.Params.EIP1052Transition = n
	spec.Params.EIP1283Transition = n
}

func (spec *parityChainSpec) setConstantinopleFix(num *big.Int) {
	spec.Params.EIP1283DisableTransition = common.Uint64(num.Uint64())
}

func (spec *parityChainSpec) setIstanbul(num *big.Int) {
	spec.Params.EIP1344Transition = common.Uint64(num.Uint64())
	spec.Params.EIP1884Transition = common.Uint64(num.Uint64())
	spec.Params.EIP2028Transition = common.Uint64(num.Uint64())
	spec.Params.EIP1283ReenableTransition = common.Uint64(num.Uint64())
}

// pyEthereumGenesisSpec represents the genesis specification format used by the
// Python Ethereum implementation.
type pyEthereumGenesisSpec struct {
	Nonce      types.BlockNonce  `json:"nonce"`
	Timestamp  common.Uint64     `json:"timestamp"`
	ExtraData  common.Bytes      `json:"extraData"`
	GasLimit   common.Uint64     `json:"gasLimit"`
	Difficulty *common.Big       `json:"difficulty"`
	Mixhash    common.Hash       `json:"mixhash"`
	Coinbase   common.Address    `json:"coinbase"`
	Alloc      core.GenesisAlloc `json:"alloc"`
	ParentHash common.Hash       `json:"parentHash"`
}

// newPyEthereumGenesisSpec converts a go-ethereum genesis block into a Parity specific
// chain specification format.
func newPyEthereumGenesisSpec(network string, genesis *core.Genesis) (*pyEthereumGenesisSpec, error) {
	// Only ethash is currently supported between go-ethereum and pyethereum
	if genesis.Config.Ethash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	spec := &pyEthereumGenesisSpec{
		Nonce:      types.EncodeNonce(genesis.Nonce),
		Timestamp:  (common.Uint64)(genesis.Timestamp),
		ExtraData:  genesis.ExtraData,
		GasLimit:   (common.Uint64)(genesis.GasLimit),
		Difficulty: (*common.Big)(genesis.Difficulty),
		Mixhash:    genesis.Mixhash,
		Coinbase:   genesis.Coinbase,
		Alloc:      genesis.Alloc,
		ParentHash: genesis.ParentHash,
	}
	return spec, nil
}
