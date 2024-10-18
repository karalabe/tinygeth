// Copyright 2024 The go-ethereum Authors
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

package state

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"slices"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/karalabe/tinygeth/core/stateless"
	"golang.org/x/exp/maps"
)

func TestWitnessDeletionOrder(t *testing.T) {
	db := NewDatabaseForTesting()

	// Generate 3 keys with a shared prefix
	var keys []common.Hash
	for len(keys) < 5 {
		var key common.Hash
		io.ReadFull(rand.Reader, key[:])

		if bytes.HasPrefix(crypto.Keccak256Hash(key[:]).Bytes()[:3], []byte{0, 0, 0}) {
			keys = append(keys, key)
		}
	}
	// Seed the database with a single small storage accout
	state, _ := New(common.Hash{}, db)

	state.SetNonce(common.Address{255}, 1)
	for i, key := range keys {
		state.SetState(common.Address{255}, key, common.Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(i + 1)})
	}
	root, err := state.Commit(1, true)
	if err != nil {
		panic(err)
	}
	// Delete the storage account and generate the witness a few times, making
	// sure it's always the same (avoid some reordering issue).
	var first []string

	for i := 0; i < 100; i++ {
		state, _ = New(root, db)

		witness, _ := stateless.NewWitness(nil, nil)
		state.StartPrefetcher("", witness)

		for _, key := range keys {
			state.SetState(common.Address{255}, key, common.Hash{})
		}
		if _, err = state.Commit(2, true); err != nil {
			panic(err)
		}
		nodes := maps.Keys(witness.State)
		sort.Strings(nodes)

		if first == nil {
			first = nodes
			for j, node := range nodes {
				fmt.Printf("%d: %x\n", j, node)
			}
		} else if slices.Compare(first, nodes) != 0 {
			for j, want := range first {
				fmt.Printf("%d:\n - want %x\n - have %x\n\n", j, want, nodes[j])
			}
			t.Fatalf("test %d: witness mismatch", i)
		}
	}
}
