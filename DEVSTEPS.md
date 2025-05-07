## 1. 项目搭建
    搭建 chains-union-rpc
    项目目录如下：
```bash
├── bin               protobuf 命令管理       
├── chains            支持的链       
├── chaindispatcher   接口分发
├── common            通用工具库
├── config            配置代码
├── proto             grpc 生成的 protobuf代码
├── main.go           程序主入口
├── go.mod            依赖管理
├── config.yml        配置文件
├── Makefile          shell 命令管理
├── README.md         项目文档
├── DEVSTEPTS.md      项目开发步骤

```
## 2. 编写 proto 文件
- proto 中，编写 account.proto文件 
- 定义消息、接口
- 调用 make proto 生成相应的 go 代码

## 3. 搭建 grpc 框架
- 新建 grpcServer
- 注册 chainDispatcher
- tcp 端口监听启动
- 实现 dispatcher.go 负责接口的转发给不同链的调用，实例需继承 IChainAdaptor
- 构建 所有链调用实例的代码，如 ethereum.go, 实现 IChainAdaptor

## 4. 实现对接 rpc 节点和数据平台的 client
- 实现 ethclient.go 对接区块链 rpc 节点
- 实现 erc20data.go 对接区块链数据平台

## 5. 实现对接链的接口
```bash
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
```

