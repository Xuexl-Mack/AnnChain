// Copyright 2017 ZhongAn Information Technology Services Co.,Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	at "github.com/dappledger/AnnChain/angine/types"
	ethcmn "github.com/dappledger/AnnChain/genesis/eth/common"
	"github.com/dappledger/AnnChain/genesis/types"
)

type DoManageData struct {
	app *GenesisApp
	op  *types.ManageDataOp
	tx  *types.Transaction
}

//type ValueCategory struct {
//	value    string
//	category string
//}

func (ct *DoManageData) PreCheck() at.Result {
	return at.NewResultOK([]byte{}, "")
}

func (ct *DoManageData) CheckValid(stateDup *stateDup) error {
	if ct.op.Source == types.ZERO_ADDRESS {
		return at.NewError(at.CodeType_BaseUnknownAddress, at.CodeType_BaseUnknownAddress.String())
	}
	return nil
}

func (ct *DoManageData) Apply(stateDup *stateDup) error {

	db := ct.app.dataM

	toAdd := make(map[string]types.ValueCategory)

	for idx, r := range ct.op.DataName {
		if err := db.QueryOneAccData(ct.op.Source, r); err != nil {
			return at.NewError(at.CodeType_BaseInvalidInput, err.Error())
		}
		var vc types.ValueCategory
		vc.Value = ct.op.Data[idx]
		vc.Category = ct.op.Category[idx]
		toAdd[r] = vc
	}

	for k, v := range toAdd {
		_, err := db.AddAccData(ct.op.Source, k, v.Value, v.Category)
		if err != nil {
			return at.NewError(at.CodeType_SaveFailed, at.CodeType_SaveFailed.String()+err.Error())
		}
		if len(k) > types.AccDataLength || len(k) == 0 {
			return at.NewError(at.CodeType_InvalidTx, at.CodeType_InvalidTx.String())
		}

		if len(v.Value) > types.AccDataLength {
			return at.NewError(at.CodeType_AccDataLengthError, at.CodeType_AccDataLengthError.String())
		}

		if len(v.Category) > types.AccCategoryLength {
			return at.NewError(at.CodeType_AccCategoryLengthError, at.CodeType_AccCategoryLengthError.String())
		}
	}
	ct.SetEffects(ct.op.Source, toAdd)
	return nil
}

func (ct *DoManageData) SetEffects(source ethcmn.Address, toadd map[string]types.ValueCategory) {
	ct.op.SetEffects(&types.ActionManageData{
		ActionBase: types.ActionBase{
			Typei:       types.OP_S_MANAGEDATA.OpInt(),
			Type:        types.OP_S_MANAGEDATA,
			FromAccount: source,
			Nonce:       ct.tx.Nonce(),
			BaseFee:     ct.tx.BaseFee(),
			Memo:        ct.tx.Data.Memo,
		},
		KeyPair: toadd,
		Source:  source,
	}, []types.EffectObject{})
	return
}
