package ethereum

import (
	"chains-union-rpc/chains"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"

	"chains-union-rpc/config"
	"chains-union-rpc/proto/chainsunion"
)

func setup() (chains.IChainAdaptor, error) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		log.Error("load config failed, error:", err)
		return nil, err
	}
	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		log.Error("create chain adaptor failed, error:", err)
		return nil, err
	}
	return adaptor, nil
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	// test account
	// privateKey: 861ae7df240f80e5492065dafeb6444bdbf2d55d01e1797d2abe0db0afd4f917
	// publicKey: 02410c64fcd262512683b54576440e3d3033d825ef9f753b44c51ccdd70a7e90c3
	resp, err := adaptor.ConvertAddress(&chainsunion.ConvertAddressRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		PublicKey: "048318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed753547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5",
	})
	if err != nil {
		t.Error("convert address failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Address)

	respJson, _ := json.Marshal(resp)
	t.Logf("响应: %s", respJson)
}

func TestChainAdaptor_ValidAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.ValidAddress(&chainsunion.ValidAddressRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Address: "0x8358d847Fc823097380c4996A3D3485D9D86941f",
	})
	if err != nil {
		t.Error("valid address failed:", err)
		return
	}
	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)
}

func TestChainAdaptor_GetBlockHeaderByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockHeaderByNumber(&chainsunion.BlockHeaderNumberRequest{
		Chain:   ChainName,
		Network: "testnet",
		Height:  999,
	})
	if err != nil {
		t.Error("get block header by number failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.BlockHeader)
}

func TestChainAdaptor_GetBlockHeaderByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockHeaderByHash(&chainsunion.BlockHeaderHashRequest{
		Chain:   ChainName,
		Network: "testnet",
		Hash:    "0xd1c017c3d845f3dd5ed3ea17700ca5afdd8ec106662a4bb8a3f6edb856bfe5f4",
	})
	if err != nil {
		t.Error("get block header by hash failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.BlockHeader)
}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByNumber(&chainsunion.BlockNumberRequest{
		Chain:  ChainName,
		Height: 999,
		ViewTx: true,
	})
	if err != nil {
		t.Error("get block by number failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Transactions)
}

func TestChainAdaptor_GetBlockByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByHash(&chainsunion.BlockHashRequest{
		Chain:  ChainName,
		Hash:   "0xd1c017c3d845f3dd5ed3ea17700ca5afdd8ec106662a4bb8a3f6edb856bfe5f4",
		ViewTx: true,
	})
	if err != nil {
		t.Error("get block by hash failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Transactions)
}

func TestChainAdaptor_GetAccount(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetAccount(&chainsunion.AccountRequest{
		Chain:           ChainName,
		Network:         "mainnet",
		Address:         "0xD79053a14BC465d9C1434d4A4fAbdeA7b6a2A94b",
		ContractAddress: "0x00",
	})
	if err != nil {
		t.Error("get account failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)

	respJson, _ := json.Marshal(resp)
	t.Logf("响应: %s", respJson)
}

func TestChainAdaptor_GetFee(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetFee(&chainsunion.FeeRequest{
		Chain:   ChainName,
		Network: "mainnet",
	})
	if err != nil {
		t.Error("get account failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)

	respJson, _ := json.Marshal(resp)
	t.Logf("响应: %s", respJson)
}

func TestChainAdaptor_GetTxByAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetTxByAddress(&chainsunion.TxAddressRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Address: "0x0B70B578aBd96AAb5e80D24D1f3C28DbdE14356a",
	})
	if err != nil {
		t.Error("get transaction by address failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Tx)
}

func TestChainAdaptor_GetTxByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetTxByHash(&chainsunion.TxHashRequest{
		Chain:   ChainName,
		Network: "testnet",
		Hash:    "0xb318637c6a28109e6ee2ad599681defd87645bcc08c87b6d191225f38c8ffd1c",
	})
	if err != nil {
		t.Error("get transaction by address failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Tx)
}

func TestChainAdaptor_GetBlockByRange(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByRange(&chainsunion.BlockByRangeRequest{
		Chain:   ChainName,
		Network: "testnet",
		Start:   "997",
		End:     "999",
	})
	if err != nil {
		t.Error("get block by range failed:", err)
		return
	}

	assert.Equal(t, chainsunion.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.GetBlockHeader())
}
