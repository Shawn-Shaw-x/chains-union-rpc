package chaindispatcher

import (
	"chains-union-rpc/chains"
	"chains-union-rpc/chains/ethereum"
	"chains-union-rpc/config"
	"chains-union-rpc/proto/chainsunion"
	"context"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"strings"
)

type CommonReply = chainsunion.SupportChainsResponse

type ChainType = string

type CommonRequest interface {
	GetChain() string
}

type ChainDispatcher struct {
	registry map[ChainType]chains.IChainAdaptor
}

/*新建处理器*/
func New(conf *config.Config) (*ChainDispatcher, error) {
	dispatcher := &ChainDispatcher{
		registry: make(map[ChainType]chains.IChainAdaptor),
	}
	/*适配器工厂，根据key 按策略分发给不同 adaptor 实现*/
	chainAdaptorFactoryMap := map[string]func(conf *config.Config) (chains.IChainAdaptor, error){
		strings.ToLower(ethereum.ChainName): ethereum.NewChainAdaptor,
	}
	supportedChains := []string{
		strings.ToLower(ethereum.ChainName),
	}

	/*全部 new 出来放到工厂 map 中*/
	for _, c := range conf.Chains {
		chainName := strings.ToLower(c)
		if factory, ok := chainAdaptorFactoryMap[chainName]; ok {
			adaptor, err := factory(conf)
			if err != nil {
				log.Crit("Failed to create chain adaptor", "err", err)
			}
			dispatcher.registry[chainName] = adaptor
		} else {
			log.Error("Unsupported chain type", "chain", chainName, "supportedChains", supportedChains)
		}
	}
	return dispatcher, nil
}

/*拦截器*/
func (d *ChainDispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	/*错误处理,防止出错全部程序崩溃*/
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			log.Debug("panic", "stack", string(debug.Stack()))
			err = status.Errorf(codes.Internal, "panic: %v", e)
		}
	}()
	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]

	/*获取请求参数*/
	chainName := req.(CommonRequest).GetChain()
	log.Info(method, "chain", chainName, "req", req)

	resp, err = handler(ctx, req)
	log.Debug("finish handling", "resp", resp, "err", err)
	return
}

/*接口预处理*/
func (d *ChainDispatcher) preHandler(req interface{}) (resp *CommonReply, chainName string) {
	chainName = strings.ToLower(req.(CommonRequest).GetChain())
	log.Debug("chain", "chainName", chainName, "req", req)
	if _, ok := d.registry[ChainType(req.(CommonRequest).GetChain())]; !ok {
		return &CommonReply{
			Code:    chainsunion.ReturnCode_ERROR,
			Msg:     config.UnsupportedChain,
			Support: false,
		}, chainName
	}
	return nil, chainName
}

func (d *ChainDispatcher) GetSupportChains(ctx context.Context, request *chainsunion.SupportChainsRequest) (*chainsunion.SupportChainsResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.SupportChainsResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetSupportChains(request)
}

func (d *ChainDispatcher) ConvertAddress(ctx context.Context, request *chainsunion.ConvertAddressRequest) (*chainsunion.ConvertAddressResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.ConvertAddressResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].ConvertAddress(request)
}

func (d *ChainDispatcher) ValidAddress(ctx context.Context, request *chainsunion.ValidAddressRequest) (*chainsunion.ValidAddressResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.ValidAddressResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].ValidAddress(request)
}

func (d *ChainDispatcher) GetBlockByNumber(ctx context.Context, request *chainsunion.BlockNumberRequest) (*chainsunion.BlockResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.BlockResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetBlockByNumber(request)
}

func (d *ChainDispatcher) GetBlockByHash(ctx context.Context, request *chainsunion.BlockHashRequest) (*chainsunion.BlockResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.BlockResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetBlockByHash(request)

}

func (d *ChainDispatcher) GetBlockHeaderByHash(ctx context.Context, request *chainsunion.BlockHeaderHashRequest) (*chainsunion.BlockHeaderResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.BlockHeaderResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetBlockHeaderByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByNumber(ctx context.Context, request *chainsunion.BlockHeaderNumberRequest) (*chainsunion.BlockHeaderResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.BlockHeaderResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetBlockHeaderByNumber(request)

}

func (d *ChainDispatcher) GetBlockHeaderByRange(ctx context.Context, request *chainsunion.BlockByRangeRequest) (*chainsunion.BlockByRangeResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.BlockByRangeResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetBlockHeaderByRange(request)
}

func (d *ChainDispatcher) GetAccount(ctx context.Context, request *chainsunion.AccountRequest) (*chainsunion.AccountResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.AccountResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetAccount(request)
}

func (d *ChainDispatcher) GetFee(ctx context.Context, request *chainsunion.FeeRequest) (*chainsunion.FeeResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.FeeResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetFee(request)
}

func (d *ChainDispatcher) SendTx(ctx context.Context, request *chainsunion.SendTxRequest) (*chainsunion.SendTxResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.SendTxResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].SendTx(request)
}

func (d *ChainDispatcher) GetTxByAddress(ctx context.Context, request *chainsunion.TxAddressRequest) (*chainsunion.TxAddressResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.TxAddressResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetTxByAddress(request)
}

func (d *ChainDispatcher) GetTxByHash(ctx context.Context, request *chainsunion.TxHashRequest) (*chainsunion.TxHashResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.TxHashResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetTxByHash(request)
}

func (d *ChainDispatcher) BuildUnSignTransaction(ctx context.Context, request *chainsunion.UnSignTransactionRequest) (*chainsunion.UnSignTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.UnSignTransactionResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].BuildUnSignTransaction(request)
}

func (d *ChainDispatcher) BuildSignedTransaction(ctx context.Context, request *chainsunion.SignedTransactionRequest) (*chainsunion.SignedTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.SignedTransactionResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].BuildSignedTransaction(request)
}

func (d *ChainDispatcher) DecodeTransaction(ctx context.Context, request *chainsunion.DecodeTransactionRequest) (*chainsunion.DecodeTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.DecodeTransactionResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].DecodeTransaction(request)
}

func (d *ChainDispatcher) VerifySignedTransaction(ctx context.Context, request *chainsunion.VerifyTransactionRequest) (*chainsunion.VerifyTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.VerifyTransactionResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].VerifySignedTransaction(request)
}

func (d *ChainDispatcher) GetExtraData(ctx context.Context, request *chainsunion.ExtraDataRequest) (*chainsunion.ExtraDataResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &chainsunion.ExtraDataResponse{
			Code: chainsunion.ReturnCode_ERROR,
			Msg:  config.UnsupportedChain + " " + chainName,
		}, nil
	}
	return d.registry[chainName].GetExtraData(request)
}

func (d *ChainDispatcher) GetNftListByAddress(ctx context.Context, request *chainsunion.NftAddressRequest) (*chainsunion.NftAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *ChainDispatcher) GetNftCollection(ctx context.Context, request *chainsunion.NftCollectionRequest) (*chainsunion.NftCollectionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *ChainDispatcher) GetNftDetail(ctx context.Context, request *chainsunion.NftDetailRequest) (*chainsunion.NftDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *ChainDispatcher) GetNftHolderList(ctx context.Context, request *chainsunion.NftHolderListRequest) (*chainsunion.NftHolderListResponse, error) {

	//TODO implement me
	panic("implement me")
}

func (d *ChainDispatcher) GetNftTradeHistory(ctx context.Context, request *chainsunion.NftTradeHistoryRequest) (*chainsunion.NftTradeHistoryResponse, error) {

	//TODO implement me
	panic("implement me")
}

func (d *ChainDispatcher) GetAddressNftTradeHistory(ctx context.Context, request *chainsunion.AddressNftTradeHistoryRequest) (*chainsunion.AddressNftTradeHistoryResponse, error) {

	//TODO implement me
	panic("implement me")
}
