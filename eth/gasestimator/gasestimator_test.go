// Copyright 2023 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package gasestimator

// Corner case: plain transfer with unused access list. The estimation should
// detect the possible early exist, but should *not* blindly return 21000 and
// instead detect that the access list ruins it.

/*
func TestEstimateGas(t *testing.T) {
	// Initialize test accounts
	var (
		accounts = newAccounts(2)
		genesis  = &core.Genesis{
			Config: params.TestChainConfig,
			Alloc: core.GenesisAlloc{
				accounts[0].addr: {Balance: big.NewInt(params.Ether)},
				accounts[1].addr: {Balance: big.NewInt(params.Ether)},
			},
		}
		genBlocks      = 10
		signer         = types.HomesteadSigner{}
		randomAccounts = newAccounts(2)
	)
	api := NewBlockChainAPI(newTestBackend(t, genBlocks, genesis, ethash.NewFaker(), func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &accounts[1].addr, Value: big.NewInt(1000), Gas: params.TxGas, GasPrice: b.BaseFee(), Data: nil}), signer, accounts[0].key)
		b.AddTx(tx)
	}))
	var tests = []struct {
		blockNumber rpc.BlockNumber
		call        TransactionArgs
		expectErr   error
		want        uint64
	}{
		// simple transfer on latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: nil,
			want:      21000,
		},
		// simple transfer with insufficient funds on latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:  &randomAccounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: core.ErrInsufficientFunds,
			want:      21000,
		},
		// empty create
		{
			blockNumber: rpc.LatestBlockNumber,
			call:        TransactionArgs{},
			expectErr:   nil,
			want:        53000,
		},
		// Test for a bug where the gas price was set to zero but the basefee non-zero
		//
		// contract BasefeeChecker {
		//    constructor() {
		//        require(tx.gasprice >= block.basefee);
		//        if (tx.gasprice > 0) {
		//            require(block.basefee > 0);
		//        }
		//    }
		//}
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:     &accounts[0].addr,
				Input:    hex2Bytes("6080604052348015600f57600080fd5b50483a1015601c57600080fd5b60003a111560315760004811603057600080fd5b5b603f80603e6000396000f3fe6080604052600080fdfea264697066735822122060729c2cee02b10748fae5200f1c9da4661963354973d9154c13a8e9ce9dee1564736f6c63430008130033"),
				GasPrice: (*hexutil.Big)(big.NewInt(1_000_000_000)), // Legacy as pricing
			},
			expectErr: nil,
			want:      67617,
		},
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:         &accounts[0].addr,
				Input:        hex2Bytes("6080604052348015600f57600080fd5b50483a1015601c57600080fd5b60003a111560315760004811603057600080fd5b5b603f80603e6000396000f3fe6080604052600080fdfea264697066735822122060729c2cee02b10748fae5200f1c9da4661963354973d9154c13a8e9ce9dee1564736f6c63430008130033"),
				MaxFeePerGas: (*hexutil.Big)(big.NewInt(1_000_000_000)), // 1559 gas pricing
			},
			expectErr: nil,
			want:      67617,
		},
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:         &accounts[0].addr,
				Input:        hex2Bytes("6080604052348015600f57600080fd5b50483a1015601c57600080fd5b60003a111560315760004811603057600080fd5b5b603f80603e6000396000f3fe6080604052600080fdfea264697066735822122060729c2cee02b10748fae5200f1c9da4661963354973d9154c13a8e9ce9dee1564736f6c63430008130033"),
				GasPrice:     nil, // No legacy gas pricing
				MaxFeePerGas: nil, // No 1559 gas pricing
			},
			expectErr: nil,
			want:      67595,
		},
	}
	for i, tc := range tests {
		result, err := Estimate(context.Background(), tc.call, &rpc.BlockNumberOrHash{BlockNumber: &tc.blockNumber})
		if tc.expectErr != nil {
			if err == nil {
				t.Errorf("test %d: want error %v, have nothing", i, tc.expectErr)
				continue
			}
			if !errors.Is(err, tc.expectErr) {
				t.Errorf("test %d: error mismatch, want %v, have %v", i, tc.expectErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("test %d: want no error, have %v", i, err)
			continue
		}
		if uint64(result) != tc.want {
			t.Errorf("test %d, result mismatch, have\n%v\n, want\n%v\n", i, uint64(result), tc.want)
		}
	}
}
*/
