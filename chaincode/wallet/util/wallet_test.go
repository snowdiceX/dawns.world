package util

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func TestWallet(t *testing.T) {
	w := NewWallet("ethereum", "eth", "address12345678")
	if w.Key != "Wallet-ETHEREUM-ETH-ADDRESS12345678" {
		t.Error(fmt.Sprintf(
			"wrong wallet key: %s", w.Key))
	}
}

func TestWalletSum(t *testing.T) {
	w := NewWallet("ethereum", "eth", "0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2")
	w.Balance = "0x0"

	txs := []string{
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","amount":"0xde0b6b3a7640000","info":{"from":"0x4d6bb4ed029b33cf25d0810b029bd8b1a6bcab7b","to":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","amount":"0xde0b6b3a7640000","gasUsed":"0x5208","gasPrice":"0x77359400","txHash":"0xe6da8b5ff3a8cfaa84aaf9176bf69f430d653e27d5ee2fe4116dfed5e9c11314","height":"0xa35b03","status":"0x1"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","gasUsed":"0x5208","gasPrice":"0x2540be400","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xead5c972fe8bbf6f725ab8a4c7e9d40e15f35241","amount":"0x6f05b59d3b20000","gasUsed":"0x5208","gasPrice":"0x2540be400","txHash":"0x9498cd360d1e295fe294a423265548af356bca1f49faf60560a27015702039b4","height":"0xa35b42","status":"0x0"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","amount":"0x29a2241af62c0000","info":{"from":"0x3d947eb8c366d2416468675cedd00fd311d70dfb","to":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","amount":"0x29a2241af62c0000","gasUsed":"0x5208","gasPrice":"0x2540be400","txHash":"0xf86d22a587004bfe52942d332b399130b1632d6b8e4b1e2fead037bb9e4968fb","height":"0xa35b4a","status":"0x1"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","gasUsed":"0x5208","gasPrice":"0x2540be400","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xead5c972fe8bbf6f725ab8a4c7e9d40e15f35241","amount":"0xde0b6b3a7640000","gasUsed":"0x5208","gasPrice":"0x2540be400","txHash":"0xd1477743101cb2af7c712a3d1cf84790756b28ea74cf238c58f208926c08f77a","height":"0xa35b6f","status":"0x0"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","gasUsed":"0x5208","gasPrice":"0x4a817c800","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xead5c972fe8bbf6f725ab8a4c7e9d40e15f35241","amount":"0xde0b6b3a7640000","gasUsed":"0x5208","gasPrice":"0x4a817c800","txHash":"0xac9a4af995c2fb996379a8a95dd41d69a7a3d04335b971339aeb8419d0bb7db0","height":"0xa35b72","status":"0x0"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","gasUsed":"0x5236","gasPrice":"0x2540be400","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xead5c972fe8bbf6f725ab8a4c7e9d40e15f35241","amount":"0xde0b6b3a7640000","gasUsed":"0x5236","gasPrice":"0x2540be400","txHash":"0xc07a33c29304e6eb92cfce2c9d95d916686648b973a86aa791efce6c35f98998","height":"0xa35b76","status":"0x0"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","gasUsed":"0x5236","gasPrice":"0x2540be400","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xead5c972fe8bbf6f725ab8a4c7e9d40e15f35241","amount":"0xde0b6b3a7640000","gasUsed":"0x5236","gasPrice":"0x2540be400","txHash":"0x88efebd7469c323df3f7c1e57e51e586643d3df9c80ff18f70a48e93b13fea89","height":"0xa35b7d","status":"0x0"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","amount":"0x-0","gasUsed":"0xc980","gasPrice":"0x12a05f200","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xc5d0ac103d253ca6fad4ec3170391ffab6fe5bb8","amount":"0x0","gasUsed":"0xc980","gasPrice":"0x12a05f200","txHash":"0x4e9255e66cd7a948d600d87ca1fb3e8dc2ec3edfb72659e12d0918177812dae4","height":"0xa35b95","status":"0x1"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","amount":"0x-0","gasUsed":"0x8ee8","gasPrice":"0x12a05f200","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xc5d0ac103d253ca6fad4ec3170391ffab6fe5bb8","amount":"0x0","gasUsed":"0x8ee8","gasPrice":"0x12a05f200","txHash":"0x7fd50e7647dbcfb4a8620932e899881ca3781818a646932c952841e0c61fb67f","height":"0xa35b96","status":"0x1"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","amount":"0x-0","gasUsed":"0x90a2","gasPrice":"0x2540be400","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xc5d0ac103d253ca6fad4ec3170391ffab6fe5bb8","amount":"0x0","gasUsed":"0x90a2","gasPrice":"0x2540be400","txHash":"0x5cf5170be7d3e6431681a3ba36e8f1abe3e2ab005004a68a5a392a8e6d2fb1d1","height":"0xa36bd7","status":"0x1"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","amount":"0x-0","gasUsed":"0x90a2","gasPrice":"0x2540be400","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xc5d0ac103d253ca6fad4ec3170391ffab6fe5bb8","amount":"0x0","gasUsed":"0x90a2","gasPrice":"0x2540be400","txHash":"0x6656a724faa43cf6e6a57a11df5d223c9052d7a43b8240809f4326dcdecc8cdd","height":"0xa36fb7","status":"0x1"}}`,
		`{"chain":"ethereum","token":"eth","address":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","amount":"0x-0","gasUsed":"0x10ae2","gasPrice":"0x2540be400","info":{"from":"0xb0d2da0f43cd2e44e4f3a38e24945f0ca0ea95e2","to":"0xead5c972fe8bbf6f725ab8a4c7e9d40e15f35241","amount":"0x0","gasUsed":"0x10ae2","gasPrice":"0x2540be400","txHash":"0x8fb60926300c25df4661e1870adc764cacfeefb3b5d128d625f70fcd5ea09cf5","height":"0xa36fbd","status":"0x1"}}`}

	for _, txJSON := range txs {
		t.Log("tx json: ", txJSON)
		tx := &TxRegister{}
		if err := json.Unmarshal([]byte(txJSON), tx); err != nil {
			t.Errorf("tx unmarshal error: %v", err)
		}
		t.Log("tx address: ", tx.Address())
		if err := w.Sum(tx); err != nil {
			t.Errorf("wallet sum error: %v", err)
		}
		h := new(big.Int)
		h.SetString(tx.Info.Height[2:], 16)
		t.Logf("height: %s; %s; TxHash: %s", tx.Info.Height, h.Text(10), tx.Info.TxHash)
	}
	t.Logf("wallet balance: %s", w.Balance)
	if strings.ToUpper("0x3777c02e70512800") !=
		strings.ToUpper(w.Balance) {
		t.Errorf("wrong result! %s", w.Balance)
	}
}
