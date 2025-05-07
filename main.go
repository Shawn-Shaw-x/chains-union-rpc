package main

import (
	"chains-union-rpc/chaindispatcher"
	"chains-union-rpc/config"
	"chains-union-rpc/proto/chainsunion"
	"flag"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	dispatcher, err := chaindispatcher.New(conf)
	if err != nil {
		log.Error("setup dispatcher fail", "err", err)
		panic(err)
	}

	/*注册拦截器*/
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	defer grpcServer.GracefulStop()

	/*注册实现的接口路由*/
	chainsunion.RegisterChainsUnionServiceServer(grpcServer, dispatcher)

	listen, err := net.Listen("tcp", ":"+conf.Server.Port)
	if err != nil {
		log.Error("failed to listen", "port", conf.Server.Port, "err", err)
		panic(err)
	}
	reflection.Register(grpcServer)

	log.Info("chains-union-rpc starting server", "port", conf.Server.Port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Error("failed to serve", "port", conf.Server.Port, "err", err)
		panic(err)
	}

}
