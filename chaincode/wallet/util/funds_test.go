package util

import (
	"fmt"
	"math/big"
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

func TestBigInt(t *testing.T) {
	gu := new(big.Int)
	gu.SetString("0x5236"[2:], 16)
	gp := new(big.Int)
	gp.SetString("0x2540be400"[2:], 16)

	m := new(big.Int)
	m.Mul(gu, gp)
	t.Logf("0x%s * 0x%s : 0x%s\n", gu.Text(16), gp.Text(16), m.Text(16))

	b := new(big.Int)
	b.SetString("0xde0b6b3a7640000"[2:], 16)

	c := b.CmpAbs(m)
	t.Logf("c: %d", c)

	b.Sub(b, m)
	t.Logf("0x%s\n", b.Text(16))
}
