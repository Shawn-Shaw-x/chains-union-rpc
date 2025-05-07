# chains-union-rpc 多链聚合 rpc

一个轻量级高性能的 Go RPC 服务，通过聚合多链的 rpc 接口，提供统一的接口调用。
---
在钱包开发中，您是否对接对每一个链都要开发一个 client 去进行 rpc 调用感到厌烦？您是否对每个链的不同接口特性开发感到混乱？如果您的回答是 YES，那么，这个项目将会拯救你。
---


## ✨ 核心特性
- 通过策略模式去调用每一条链，提供统一的接口
- 模块化结构，新增链支持只需编写某个链的代码即可
- 基于 GRPC 通信协议，简单高性能
- 易于配置，只需提供必要配置即可运行

## 📦 安装要求

- Go 1.20+
- `protoc` (Protocol Buffers compiler)
- `protoc-gen-go`, `protoc-gen-go-grpc`

## 🚀 安装指南

### Clone & Build

```bash
git clone https://github.com/your-org/MyGoRPCProject.git
cd chains-union-rpc
go mod tidy
make compile
make chains-union-rpc
./chains-union-rpc
```
```bash
grpcui -plaintext 127.0.0.1:8189
```

## ⭐️ 项目架构

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

## 🐰 项目架构图



## 🍌 实现的接口
    1. 链支持
	2. 地址转换
	3. 地址校验
	4. 根据区块号获取区块
	5. 根据 hash 获取区块
	6. 根据 hash 获取区块头
	7. 根据区块号获取区块头
	8. 根据范围获取区块头
	9. 获取账号信息
	10. 获取手续费
	11. 发送交易
	12. 根据地址获取交易
	13. 根据交易 hash 获取交易
	14. 构建未签名交易
	15. 构建已签名交易
	16. 交易解码，解析成可读形式
	17. 校验已签名交易
	18. 获取额外数据
	19. 获取某个地址的 NFT 列表
	20. 获取 NFT 的集合
	21. 获取 NFT 的细节
	22. 获取 NFT 的持有者列表
	23. 获取 NFT 的交易历史
	24. 获取某个地址的 NFT 交易历史
	25 获取范围内区块
