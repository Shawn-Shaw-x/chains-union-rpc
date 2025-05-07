package chains

import (
	"chains-union-rpc/proto/chainsunion"
)

type IChainAdaptor interface {
	// 链支持
	GetSupportChains(req *chainsunion.SupportChainsRequest) (*chainsunion.SupportChainsResponse, error)
	// 地址转换
	ConvertAddress(req *chainsunion.ConvertAddressRequest) (*chainsunion.ConvertAddressResponse, error)
	// 地址校验
	ValidAddress(req *chainsunion.ValidAddressRequest) (*chainsunion.ValidAddressResponse, error)
	// 根据区块号获取区块
	GetBlockByNumber(req *chainsunion.BlockNumberRequest) (*chainsunion.BlockResponse, error)
	// 根据 hash 获取区块
	GetBlockByHash(req *chainsunion.BlockHashRequest) (*chainsunion.BlockResponse, error)
	// 根据 hash 获取区块头
	GetBlockHeaderByHash(req *chainsunion.BlockHeaderHashRequest) (*chainsunion.BlockHeaderResponse, error)
	// 根据区块号获取区块头
	GetBlockHeaderByNumber(req *chainsunion.BlockHeaderNumberRequest) (*chainsunion.BlockHeaderResponse, error)
	// 根据范围获取区块头
	GetBlockHeaderByRange(req *chainsunion.BlockByRangeRequest) (*chainsunion.BlockByRangeResponse, error)
	// 获取账号信息
	GetAccount(req *chainsunion.AccountRequest) (*chainsunion.AccountResponse, error)
	// 获取手续费
	GetFee(req *chainsunion.FeeRequest) (*chainsunion.FeeResponse, error)
	// 发送交易
	SendTx(req *chainsunion.SendTxRequest) (*chainsunion.SendTxResponse, error)
	// 根据地址获取交易
	GetTxByAddress(req *chainsunion.TxAddressRequest) (*chainsunion.TxAddressResponse, error)
	// 根据交易 hash 获取交易
	GetTxByHash(req *chainsunion.TxHashRequest) (*chainsunion.TxHashResponse, error)
	// 构建未签名交易
	BuildUnSignTransaction(req *chainsunion.UnSignTransactionRequest) (*chainsunion.UnSignTransactionResponse, error)
	// 构建已签名交易
	BuildSignedTransaction(req *chainsunion.SignedTransactionRequest) (*chainsunion.SignedTransactionResponse, error)
	// 交易解码，解析成可读形式
	DecodeTransaction(req *chainsunion.DecodeTransactionRequest) (*chainsunion.DecodeTransactionResponse, error)
	// 校验已签名交易
	VerifySignedTransaction(req *chainsunion.VerifyTransactionRequest) (*chainsunion.VerifyTransactionResponse, error)
	// 获取额外数据
	GetExtraData(req *chainsunion.ExtraDataRequest) (*chainsunion.ExtraDataResponse, error)
	// 获取某个地址的 NFT 列表
	GetNftListByAddress(req *chainsunion.NftAddressRequest) (*chainsunion.NftAddressResponse, error)
	// 获取 NFT 的集合
	GetNftCollection(req *chainsunion.NftCollectionRequest) (*chainsunion.NftCollectionResponse, error)
	// 获取 NFT 的细节
	GetNftDetail(req *chainsunion.NftDetailRequest) (*chainsunion.NftDetailResponse, error)
	// 获取 NFT 的持有者列表
	GetNftHolderList(req *chainsunion.NftHolderListRequest) (*chainsunion.NftHolderListResponse, error)
	// 获取 NFT 的交易历史
	GetNftTradeHistory(req *chainsunion.NftTradeHistoryRequest) (*chainsunion.NftTradeHistoryResponse, error)
	// 获取某个地址的 NFT 交易历史
	GetAddressNftTradeHistory(req *chainsunion.AddressNftTradeHistoryRequest) (*chainsunion.AddressNftTradeHistoryResponse, error)
	// 获取范围内区块
	GetBlockByRange(req *chainsunion.BlockByRangeRequest) (*chainsunion.BlockByRangeResponse, error)
}
