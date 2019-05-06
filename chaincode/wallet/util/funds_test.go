package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestRecordFunds(t *testing.T) {
	ftk := BuildFundsBaseKey("ethereum", "eth", "hash_asdasdasd")
	rec := NewRecordFunds(ftk, "wallet_asdasdasd")

	t.Log("funds token key: ", ftk)
	t.Log("funds record key: ", rec.Key)

	if !strings.EqualFold(rec.Chain, "ETHEREUM") {
		t.Error(fmt.Sprintf(
			"chain name [%s] not expect: ethereum",
			rec.Chain))
	}
	if !strings.EqualFold(rec.Token, "Eth") {
		t.Error(fmt.Sprintf(
			"token name [%s] not expect: eth",
			rec.Token))
	}
	if !strings.EqualFold(rec.FundsHash, "hash_ASDasdasd") {
		t.Error(fmt.Sprintf(
			"funds hash [%s] not expect: hash_asdasdasd",
			rec.FundsHash))
	}
	// time.Sleep(1000 * time.Millisecond)
}
