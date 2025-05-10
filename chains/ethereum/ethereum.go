package ethereum

import (
	"chains-union-rpc/chains"
	"chains-union-rpc/chains/evmbase"
	"chains-union-rpc/common/util"
	"chains-union-rpc/config"
	"chains-union-rpc/proto/chainsunion"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/shopspring/decimal"
)

const ChainName = "Ethereum"

type ChainAdaptor struct {
	ethClient     evmbase.EthClient
	ethDataClient *evmbase.EthData
}

func NewChainAdaptor(conf *config.Config) (chains.IChainAdaptor, error) {
	ethClient, err := evmbase.DialEthClient(context.Background(), conf.WalletNode.Eth.RpcUrl)
	if err != nil {
		return nil, err
	}
	ethDataClient, err := evmbase.NewEthDataClient(conf.WalletNode.Eth.DataApiUrl, conf.WalletNode.Eth.DataApiKey, time.Second*15)
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		ethClient:     ethClient,
		ethDataClient: ethDataClient,
	}, nil
}

/*链支持*/
func (c *ChainAdaptor) GetSupportChains(req *chainsunion.SupportChainsRequest) (*chainsunion.SupportChainsResponse, error) {
	return &chainsunion.SupportChainsResponse{
		Code:    chainsunion.ReturnCode_SUCCESS,
		Msg:     "support Chains " + strings.ToLower(ChainName),
		Support: true,
	}, nil
}

/*转换地址*/
func (c *ChainAdaptor) ConvertAddress(req *chainsunion.ConvertAddressRequest) (*chainsunion.ConvertAddressResponse, error) {
	publicKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		log.Error("decode public key failed:", err)
		return &chainsunion.ConvertAddressResponse{
			Code:    chainsunion.ReturnCode_ERROR,
			Msg:     "convert address failed",
			Address: common.Address{}.String(),
		}, nil
	}
	/*以太坊的地址为，未压缩公钥为 32字节，去掉 0x04 hash后取后 20 字节*/
	addressCommon := common.BytesToAddress(crypto.Keccak256(publicKeyBytes[1:])[12:])
	return &chainsunion.ConvertAddressResponse{
		Code:    chainsunion.ReturnCode_SUCCESS,
		Msg:     "convert address success",
		Address: addressCommon.String(),
	}, nil
}

/*地址校验*/
func (c *ChainAdaptor) ValidAddress(req *chainsunion.ValidAddressRequest) (*chainsunion.ValidAddressResponse, error) {
	/*20字节 40 16进制字符，加上 0x 则为 42 个字符 */
	if len(req.Address) != 42 || !strings.HasPrefix(req.Address, "0x") {
		return &chainsunion.ValidAddressResponse{
			Code:  chainsunion.ReturnCode_SUCCESS,
			Msg:   "invalid address",
			Valid: false,
		}, nil
	}
	/*正则校验*/
	ok := regexp.MustCompile("^[0-9a-fA-F]{40}$").MatchString(req.Address[2:])
	if ok {
		return &chainsunion.ValidAddressResponse{
			Code:  chainsunion.ReturnCode_SUCCESS,
			Msg:   "valid address",
			Valid: true,
		}, nil
	} else {
		return &chainsunion.ValidAddressResponse{
			Code:  chainsunion.ReturnCode_SUCCESS,
			Msg:   "invalid address",
			Valid: false,
		}, nil
	}

}

/*根据区块号获取块*/
func (c *ChainAdaptor) GetBlockByNumber(req *chainsunion.BlockNumberRequest) (*chainsunion.BlockResponse, error) {
	block, err := c.ethClient.BlockByNumber(big.NewInt(req.Height))
	if err != nil {
		log.Error("GetBlockByNumber failed:", "err", err)
		return &chainsunion.BlockResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "block by number error",
		}, nil
	}
	blockNumber, _ := block.NumberUint64()
	var txListRet []*chainsunion.BlockInfoTransactionList
	for _, v := range block.Transactions {
		bitItem := &chainsunion.BlockInfoTransactionList{
			From:           v.From,
			To:             v.To,
			TokenAddress:   v.To,
			ContractWallet: v.To,
			Hash:           v.Hash,
			Height:         blockNumber,
			Amount:         v.Value,
		}
		txListRet = append(txListRet, bitItem)
	}
	return &chainsunion.BlockResponse{
		Code:         chainsunion.ReturnCode_SUCCESS,
		Msg:          "block by number success",
		Height:       int64(blockNumber),
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txListRet,
	}, nil
}

/*根据 hash 获取区块*/
func (c *ChainAdaptor) GetBlockByHash(req *chainsunion.BlockHashRequest) (*chainsunion.BlockResponse, error) {
	block, err := c.ethClient.BlockByHash(common.HexToHash(req.Hash))
	if err != nil {
		log.Error("GetBlockByHash failed:", err)
		return &chainsunion.BlockResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "block by number error",
		}, nil
	}
	var txListRet []*chainsunion.BlockInfoTransactionList
	for _, v := range block.Transactions {
		bitItem := &chainsunion.BlockInfoTransactionList{
			From:   v.From,
			To:     v.To,
			Hash:   v.Hash,
			Amount: v.Value,
		}
		txListRet = append(txListRet, bitItem)
	}
	blockNumber, _ := block.NumberUint64()
	return &chainsunion.BlockResponse{
		Code:         chainsunion.ReturnCode_SUCCESS,
		Msg:          "block by hash success",
		Height:       int64(blockNumber),
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txListRet,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *chainsunion.BlockHeaderHashRequest) (*chainsunion.BlockHeaderResponse, error) {
	blockInfo, err := c.ethClient.BlockHeaderByHash(common.HexToHash(req.Hash))
	if err != nil {
		log.Error("GetBlockByHash failed:", "err", err)
		return &chainsunion.BlockHeaderResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "get latest block header fail",
		}, nil
	}
	blockHeader := &chainsunion.BlockHeader{
		Hash:             blockInfo.Hash().String(),
		ParentHash:       blockInfo.ParentHash.String(),
		UncleHash:        blockInfo.UncleHash.String(),
		CoinBase:         blockInfo.Coinbase.String(),
		Root:             blockInfo.Root.String(),
		TxHash:           blockInfo.TxHash.String(),
		ReceiptHash:      blockInfo.ReceiptHash.String(),
		ParentBeaconRoot: blockInfo.ParentBeaconRoot.String(),
		Difficulty:       blockInfo.Difficulty.String(),
		Number:           blockInfo.Number.String(),
		GasLimit:         blockInfo.GasLimit,
		GasUsed:          blockInfo.GasUsed,
		Time:             blockInfo.Time,
		Extra:            base64.StdEncoding.EncodeToString(blockInfo.Extra),
		MixDigest:        blockInfo.MixDigest.String(),
		Nonce:            strconv.FormatUint(blockInfo.Nonce.Uint64(), 10),
		BaseFee:          blockInfo.BaseFee.String(),
		WithdrawalsHash:  blockInfo.WithdrawalsHash.String(),
		BlobGasUsed:      *blockInfo.BlobGasUsed,
		ExcessBlobGas:    *blockInfo.ExcessBlobGas,
	}
	return &chainsunion.BlockHeaderResponse{
		Msg:         "get latest block header success",
		BlockHeader: blockHeader,
	}, nil
}

/*根据区块号获取 header*/
func (c *ChainAdaptor) GetBlockHeaderByNumber(req *chainsunion.BlockHeaderNumberRequest) (*chainsunion.BlockHeaderResponse, error) {
	var blockNumber *big.Int
	if req.Height == 0 {
		blockNumber = nil
	} else {
		blockNumber = big.NewInt(req.Height)
	}
	/*根据区块号获取 header*/
	blockInfo, err := c.ethClient.BlockHeaderByNumber(blockNumber)
	if err != nil {
		log.Error("get latest block header failed:", "err", err)
		return &chainsunion.BlockHeaderResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "get latest block header fail",
		}, nil
	}
	log.Info("get block success", "blockInfo", blockInfo, "number", blockNumber.String(), "hash", blockInfo.Hash().Hex(), "parent hash", blockInfo.ParentHash.Hex())

	blockHead := &chainsunion.BlockHeader{
		Hash:             blockInfo.Hash().Hex(),
		ParentHash:       blockInfo.ParentHash.String(),
		UncleHash:        blockInfo.UncleHash.Hex(),
		CoinBase:         blockInfo.Coinbase.String(),
		Root:             blockInfo.Root.String(),
		TxHash:           blockInfo.TxHash.String(),
		ReceiptHash:      blockInfo.ReceiptHash.String(),
		ParentBeaconRoot: common.Hash{}.String(),
		Difficulty:       blockInfo.Difficulty.String(),
		Number:           blockInfo.Number.String(),
		GasLimit:         blockInfo.GasLimit,
		GasUsed:          blockInfo.GasUsed,
		Time:             blockInfo.Time,
		Extra:            hex.EncodeToString(blockInfo.Extra),
		MixDigest:        blockInfo.MixDigest.String(),
		Nonce:            strconv.FormatUint(blockInfo.Nonce.Uint64(), 10),
		BaseFee:          blockInfo.BaseFee.String(),
		WithdrawalsHash:  common.Hash{}.String(),
		BlobGasUsed:      0,
		ExcessBlobGas:    0,
	}

	return &chainsunion.BlockHeaderResponse{
		Code:        chainsunion.ReturnCode_SUCCESS,
		Msg:         "get latest block header success",
		BlockHeader: blockHead,
	}, nil
}

/*获取账号详情*/
func (c *ChainAdaptor) GetAccount(req *chainsunion.AccountRequest) (*chainsunion.AccountResponse, error) {
	nonceResult, err := c.ethClient.TxCountByAddress(common.HexToAddress(req.Address))
	if err != nil {
		log.Error("GetAccountByAddress failed:", "err", err)
		return &chainsunion.AccountResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "get nonce by address fail",
		}, nil
	}
	balanceResult, err := c.ethDataClient.GetBalanceByAddress(req.ContractAddress, req.Address)
	if err != nil {
		return &chainsunion.AccountResponse{
			Code:    chainsunion.ReturnCode_ERROR,
			Msg:     "get token balance fail",
			Balance: "0",
		}, nil
	}
	log.Info("balanceResult balance", "balance", balanceResult.Balance, "balanceStr", balanceResult.BalanceStr)

	balanceStr := "0"
	if balanceResult.Balance != nil && balanceResult.Balance.Int() != nil {
		balanceStr = balanceResult.Balance.Int().String()
	}
	sequence := strconv.FormatUint(uint64(nonceResult), 10)
	return &chainsunion.AccountResponse{
		Code:          chainsunion.ReturnCode_SUCCESS,
		Msg:           "get account response success",
		AccountNumber: "0",
		Sequence:      sequence,
		Balance:       balanceStr,
	}, nil

}

/*获取手续费*/
func (c *ChainAdaptor) GetFee(req *chainsunion.FeeRequest) (*chainsunion.FeeResponse, error) {
	gasPrice, err := c.ethClient.SuggestGasPrice()
	if err != nil {
		log.Error("get gas price failed:", err)
		return &chainsunion.FeeResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "get suggest gas price fail",
		}, nil
	}
	gasTipCap, err := c.ethClient.SuggestGasTipCap()
	if err != nil {
		log.Error("get gas tip cap failed:", err)
		return &chainsunion.FeeResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "get suggest gas price fail",
		}, nil
	}
	return &chainsunion.FeeResponse{
		Code:      chainsunion.ReturnCode_SUCCESS,
		Msg:       "get gas price success",
		SlowFee:   gasPrice.String() + "|" + gasTipCap.String(),
		NormalFee: gasPrice.String() + "|" + gasTipCap.String() + "|" + "*2",
		FastFee:   gasPrice.String() + "|" + gasTipCap.String() + "|" + "*3",
	}, nil
}

/*发送交易*/
func (c *ChainAdaptor) SendTx(req *chainsunion.SendTxRequest) (*chainsunion.SendTxResponse, error) {
	transaction, err := c.ethClient.SendRawTransaction(req.RawTx)
	if err != nil {
		return &chainsunion.SendTxResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "send tx error" + err.Error(),
		}, err
	}
	return &chainsunion.SendTxResponse{
		Code:   chainsunion.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: transaction.String(),
	}, nil
}

/*根据地址获取交易*/
func (c *ChainAdaptor) GetTxByAddress(req *chainsunion.TxAddressRequest) (*chainsunion.TxAddressResponse, error) {
	var resp *account2.TransactionResponse[account2.AccountTxResponse]
	var err error
	if req.ContractAddress != "0x00" && req.ContractAddress != "" {
		resp, err = c.ethDataClient.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "tokentx")
	} else {
		resp, err = c.ethDataClient.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "txlist")
	}
	if err != nil {
		log.Error("GetTxByAddress failed:", err)
		return &chainsunion.TxAddressResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "get tx list fail",
			Tx:   nil,
		}, err
	} else {
		txs := resp.TransactionList
		list := make([]*chainsunion.TxMessage, 0, len(txs))
		for i := 0; i < len(txs); i++ {
			list = append(list, &chainsunion.TxMessage{
				Hash:   txs[i].TxId,
				To:     txs[i].To,
				From:   txs[i].From,
				Fee:    txs[i].TxFee,
				Status: chainsunion.TxStatus_TX_SUCCESS,
				Value:  txs[i].Amount,
				Type:   1,
				Height: txs[i].Height,
			})
		}
		log.Info("resp:", resp)
		return &chainsunion.TxAddressResponse{
			Code: chainsunion.ReturnCode_SUCCESS,
			Msg:  "get tx list success",
			Tx:   list,
		}, nil
	}
}

/*根据 hash 获取 tx*/
func (c *ChainAdaptor) GetTxByHash(req *chainsunion.TxHashRequest) (*chainsunion.TxHashResponse, error) {
	tx, err := c.ethClient.TxByHash(common.HexToHash(req.Hash))
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return &chainsunion.TxHashResponse{
				Code: chainsunion.ReturnCode_ERROR,
				Msg:  "Ethereum Tx NotFound",
			}, nil
		}
		log.Error("get transaction error", "err", err)
		return &chainsunion.TxHashResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "Ethereum Tx NotFound",
		}, nil
	}
	/*获取交易凭证*/
	receipt, err := c.ethClient.TxReceiptByHash(common.HexToHash(req.Hash))
	if err != nil {
		log.Error("get transaction receipt error", "err", err)
		return &chainsunion.TxHashResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "Get transaction receipt error",
		}, nil
	}

	/*真实接受者地址地址*/
	var beforeToAddress string
	/*合约地址，非合约则空*/
	var beforeTokenAddress string
	/*金额*/
	var beforeValue *big.Int

	/*获取 to 地址的 code，看是否合约*/
	code, err := c.ethClient.EthGetCode(common.HexToAddress(tx.To().String()))
	if err != nil {
		log.Info("Get account code fail", "err", err)
		return nil, err
	}

	/*是合约*/
	if code == "contract" {
		inputData := hexutil.Encode(tx.Data()[:])
		if len(inputData) >= 138 && inputData[:10] == "0xa9059cbb" {
			beforeToAddress = "0x" + inputData[34:74]
			trimHex := strings.TrimLeft(inputData[74:138], "0")
			rawValue, _ := hexutil.DecodeBig("0x" + trimHex)
			beforeTokenAddress = tx.To().String()
			/*解析出金额*/
			beforeValue = decimal.NewFromBigInt(rawValue, 0).BigInt()
		}
	} else {
		/*原生交易*/
		beforeToAddress = tx.To().String()
		beforeTokenAddress = common.Address{}.String()
		beforeValue = tx.Value()
	}
	var txStatus chainsunion.TxStatus
	if receipt.Status == 1 {
		txStatus = chainsunion.TxStatus_TX_SUCCESS
	} else {
		txStatus = chainsunion.TxStatus_TX_FAILED
	}
	return &chainsunion.TxHashResponse{
		Code: chainsunion.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx: &chainsunion.TxMessage{
			Hash:            tx.Hash().Hex(),
			Index:           uint32(receipt.TransactionIndex),
			From:            beforeTokenAddress,
			To:              beforeToAddress,
			Value:           beforeValue.String(),
			Fee:             tx.GasFeeCap().String(),
			Status:          txStatus,
			Type:            0,
			Height:          receipt.BlockNumber.String(),
			ContractAddress: beforeTokenAddress,
			Data:            hexutils.BytesToHex(tx.Data()),
		},
	}, nil
}

/*从传入的参数中，构建未签名交易*/
func (c *ChainAdaptor) BuildUnSignTransaction(req *chainsunion.UnSignTransactionRequest) (*chainsunion.UnSignTransactionResponse, error) {
	response := &chainsunion.UnSignTransactionResponse{
		Code: chainsunion.ReturnCode_ERROR,
	}

	/*钱包层传过来的 base64 的数据构建出来结构体（eip1159）*/
	dFeeTx, _, err := c.buildDynamicFeeTx(req.Base64Tx)
	if err != nil {
		return nil, err
	}

	log.Info("ethereum BuildUnsignTransaction", "dFeeTx", util.ToJSONString(dFeeTx))
	/*创建未签名交易hash（32 字节）*/
	rawTx, err := evmbase.CreateEip1559UnSignTx(dFeeTx, dFeeTx.ChainID)
	if err != nil {
		log.Error("create un sign tx fail", "err", err)
		response.Msg = "get un sign tx fail"
		return response, nil
	}
	log.Info("ethereum BuildUnSignTransaction", "rawTx", rawTx)
	response.Code = chainsunion.ReturnCode_SUCCESS
	response.Msg = "create un sign tx success"
	response.UnSignTx = rawTx
	return response, nil
}

/*构建已签名交易*/
func (c *ChainAdaptor) BuildSignedTransaction(req *chainsunion.SignedTransactionRequest) (*chainsunion.SignedTransactionResponse, error) {
	response := &chainsunion.SignedTransactionResponse{
		Code: chainsunion.ReturnCode_ERROR,
	}
	dFeeTx, dynamicFeeTx, err := c.buildDynamicFeeTx(req.Base64Tx)
	if err != nil {
		log.Error("build dynamicFeeTx fail", "err", err)
		return nil, err
	}

	log.Info("ethereum BuildSignedTransaction", "dFeeTx", util.ToJSONString(dFeeTx))
	log.Info("ethereum BuildSignedTransaction", "dynamicFeeTx", util.ToJSONString(dynamicFeeTx))
	log.Info("ethereum BuildSignedTransaction", "req.Signature", req.Signature)

	// 解码出来签名
	inputSignatureByteList, err := hex.DecodeString(req.Signature)
	if err != nil {
		log.Error("decode signature failed", "err", err)
		return nil, fmt.Errorf("invalid signature: %w", err)
	}
	/*签名出来交易*/
	signer, signedTx, rawTx, txHash, err := evmbase.CreateEip1559SignedTx(dFeeTx, inputSignatureByteList, dFeeTx.ChainID)
	if err != nil {
		log.Error("create signed tx fail", "err", err)
		return nil, fmt.Errorf("create signed tx fail: %w", err)
	}
	log.Info("ethereum BuildSignedTransaction", "rawTx", rawTx)
	// 验证发送者（从签名和交易中恢复出来发送者地址，确保签名正确）
	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		log.Error("recover sender failed", "err", err)
		return nil, fmt.Errorf("recover sender failed: %w", err)
	}
	if sender.Hex() != dynamicFeeTx.FromAddress {
		log.Error("sender mismatch",
			"expected", dynamicFeeTx.FromAddress,
			"got", sender.Hex(),
		)
		return nil, fmt.Errorf("sender address mismatch: expected %s, got %s",
			dynamicFeeTx.FromAddress,
			sender.Hex(),
		)
	}

	log.Info("ethereum BuildSignedTransaction", "sender", sender.Hex())

	response.Code = chainsunion.ReturnCode_SUCCESS
	response.Msg = txHash
	response.SignedTx = rawTx
	return response, nil

}

func (c *ChainAdaptor) DecodeTransaction(req *chainsunion.DecodeTransactionRequest) (*chainsunion.DecodeTransactionResponse, error) {
	return &chainsunion.DecodeTransactionResponse{
		Code:     chainsunion.ReturnCode_SUCCESS,
		Msg:      "verify tx success",
		Base64Tx: "0x000000",
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *chainsunion.VerifyTransactionRequest) (*chainsunion.VerifyTransactionResponse, error) {
	return &chainsunion.VerifyTransactionResponse{
		Code:   chainsunion.ReturnCode_SUCCESS,
		Msg:    "verify tx success",
		Verify: true,
	}, nil
}

func (c *ChainAdaptor) GetExtraData(req *chainsunion.ExtraDataRequest) (*chainsunion.ExtraDataResponse, error) {
	return &chainsunion.ExtraDataResponse{
		Code:  chainsunion.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}

func (c *ChainAdaptor) GetNftListByAddress(req *chainsunion.NftAddressRequest) (*chainsunion.NftAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetNftCollection(req *chainsunion.NftCollectionRequest) (*chainsunion.NftCollectionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetNftDetail(req *chainsunion.NftDetailRequest) (*chainsunion.NftDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetNftHolderList(req *chainsunion.NftHolderListRequest) (*chainsunion.NftHolderListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetNftTradeHistory(req *chainsunion.NftTradeHistoryRequest) (*chainsunion.NftTradeHistoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetAddressNftTradeHistory(req *chainsunion.AddressNftTradeHistoryRequest) (*chainsunion.AddressNftTradeHistoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

/*根据范围获取区块*/
func (c *ChainAdaptor) GetBlockByRange(req *chainsunion.BlockByRangeRequest) (*chainsunion.BlockByRangeResponse, error) {
	startBlock := new(big.Int)
	endBlock := new(big.Int)
	startBlock.SetString(req.Start, 10)
	endBlock.SetString(req.End, 10)
	blockRange, err := c.ethClient.BlockHeadersByRange(startBlock, endBlock, 1)
	if err != nil {
		log.Error("get block range fail", "err", err)
		return &chainsunion.BlockByRangeResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  "get block range fail",
		}, err
	}
	blockHeaderList := make([]*chainsunion.BlockHeader, 0, len(blockRange))
	for _, block := range blockRange {
		blockItem := &chainsunion.BlockHeader{
			ParentHash:       block.ParentHash.String(),
			UncleHash:        block.UncleHash.String(),
			CoinBase:         block.Coinbase.String(),
			Root:             block.Root.String(),
			TxHash:           block.TxHash.String(),
			ReceiptHash:      block.ReceiptHash.String(),
			ParentBeaconRoot: safeToString(block.ParentBeaconRoot),
			Difficulty:       block.Difficulty.String(),
			Number:           block.Number.String(),
			GasLimit:         block.GasLimit,
			GasUsed:          block.GasUsed,
			Time:             block.Time,
			Extra:            string(block.Extra),
			MixDigest:        block.MixDigest.String(),
			Nonce:            strconv.FormatUint(block.Nonce.Uint64(), 10),
			BaseFee:          block.BaseFee.String(),
			WithdrawalsHash:  safeToString(block.WithdrawalsHash),
			BlobGasUsed:      safeUint64Ptr(block.BlobGasUsed),
			ExcessBlobGas:    safeUint64Ptr(block.ExcessBlobGas),
		}
		blockHeaderList = append(blockHeaderList, blockItem)
	}
	return &chainsunion.BlockByRangeResponse{
		Code:        chainsunion.ReturnCode_SUCCESS,
		Msg:         "get block range success",
		BlockHeader: blockHeaderList,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByRange(req *chainsunion.BlockByRangeRequest) (*chainsunion.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

// buildDynamicFeeTx 构建动态费用交易的公共方法
func (c *ChainAdaptor) buildDynamicFeeTx(base64Tx string) (*types.DynamicFeeTx, *Eip1559DynamicFeeTx, error) {
	// 1. Decode base64 string
	txReqJsonByte, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, nil, err
	}

	// 2. Unmarshal JSON to struct
	var dynamicFeeTx Eip1559DynamicFeeTx
	if err := json.Unmarshal(txReqJsonByte, &dynamicFeeTx); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, nil, err
	}

	// 3. Convert string values to big.Int
	chainID := new(big.Int)
	maxPriorityFeePerGas := new(big.Int)
	maxFeePerGas := new(big.Int)
	amount := new(big.Int)

	// 各种校验
	if _, ok := chainID.SetString(dynamicFeeTx.ChainId, 10); !ok {
		return nil, nil, fmt.Errorf("invalid chain ID: %s", dynamicFeeTx.ChainId)
	}
	if _, ok := maxPriorityFeePerGas.SetString(dynamicFeeTx.MaxPriorityFeePerGas, 10); !ok {
		return nil, nil, fmt.Errorf("invalid max priority fee: %s", dynamicFeeTx.MaxPriorityFeePerGas)
	}
	if _, ok := maxFeePerGas.SetString(dynamicFeeTx.MaxFeePerGas, 10); !ok {
		return nil, nil, fmt.Errorf("invalid max fee: %s", dynamicFeeTx.MaxFeePerGas)
	}
	if _, ok := amount.SetString(dynamicFeeTx.Amount, 10); !ok {
		return nil, nil, fmt.Errorf("invalid amount: %s", dynamicFeeTx.Amount)
	}

	// 4. Handle addresses and data
	toAddress := common.HexToAddress(dynamicFeeTx.ToAddress)
	var finalToAddress common.Address
	var finalAmount *big.Int
	var buildData []byte
	log.Info("contract address check",
		"contractAddress", dynamicFeeTx.ContractAddress,
		"isEthTransfer", isEthTransfer(&dynamicFeeTx),
	)

	// 5. Handle contract interaction vs direct transfer
	if isEthTransfer(&dynamicFeeTx) {
		finalToAddress = toAddress
		finalAmount = amount
	} else {
		contractAddress := common.HexToAddress(dynamicFeeTx.ContractAddress)
		buildData = evmbase.BuildErc20Data(toAddress, amount)
		finalToAddress = contractAddress
		finalAmount = big.NewInt(0)
	}

	// 6. Create dynamic fee transaction
	dFeeTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     dynamicFeeTx.Nonce,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Gas:       dynamicFeeTx.GasLimit,
		To:        &finalToAddress,
		Value:     finalAmount,
		Data:      buildData,
	}

	return dFeeTx, &dynamicFeeTx, nil
}

// 判断是否为 ETH 转账
func isEthTransfer(tx *Eip1559DynamicFeeTx) bool {
	// 检查合约地址是否为空或零地址
	if tx.ContractAddress == "" ||
		tx.ContractAddress == "0x0000000000000000000000000000000000000000" ||
		tx.ContractAddress == "0x00" {
		return true
	}
	return false
}

/*防止string 转换 nil panic */
func safeToString(h *common.Hash) string {
	if h != nil {
		return h.String()
	}
	return "<nil>"
}

/*防止 uint64 指针 nil*/
func safeUint64Ptr(p *uint64) uint64 {
	if p != nil {
		return *p
	}
	return 0
}
