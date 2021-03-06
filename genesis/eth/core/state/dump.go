// Copyright 2014 The go-ethereum Authors
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
	"encoding/json"
	"fmt"

	ethcmn "github.com/dappledger/AnnChain/genesis/eth/common"
	"github.com/dappledger/AnnChain/genesis/eth/rlp"
)

type DumpAccount struct {
	Balance string            `json:"balance"`
	Nonce   uint64            `json:"nonce"`
	Root    string            `json:"root"`
	Storage map[string]string `json:"storage"`
	Code    string            `json:"code"`
}

type Dump struct {
	Root     string                 `json:"root"`
	Accounts map[string]DumpAccount `json:"accounts"`
}

func (self *StateDB) RawDump() Dump {
	dump := Dump{
		Root:     ethcmn.Bytes2Hex(self.trie.Root()),
		Accounts: make(map[string]DumpAccount),
	}

	it := self.trie.Iterator()
	for it.Next() {
		addr := self.trie.GetKey(it.Key)
		var data Account
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			panic(err)
		}

		obj := newObject(nil, ethcmn.BytesToAddress(addr), data, nil)
		account := DumpAccount{
			Balance: data.Balance.String(),
			Nonce:   data.Nonce,
			Root:    ethcmn.Bytes2Hex(data.Root[:]),
			Code:    ethcmn.Bytes2Hex(obj.Code(self.db)),
			Storage: make(map[string]string),
		}
		storageIt := obj.getTrie(self.db).Iterator()
		for storageIt.Next() {
			account.Storage[ethcmn.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = ethcmn.Bytes2Hex(storageIt.Value)
		}
		dump.Accounts[ethcmn.Bytes2Hex(addr)] = account
	}
	return dump
}

func (self *StateDB) Dump() []byte {
	json, err := json.MarshalIndent(self.RawDump(), "", "    ")
	if err != nil {
		fmt.Println("dump err", err)
	}

	return json
}
