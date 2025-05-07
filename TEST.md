# TEST for Ethereum grpc api

## 1. support chain
### 输入
```bash
grpcurl -plaintext -d '{
  "chain": "ethereum"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.getSupportChains
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "support Chains Ethereum",
  "support": true
}
```


## 2. convert address
### 输入
```bash
grpcurl -plaintext -d '{
  "network": "testnet",
  "chain": "ethereum",
  "publicKey": "02e993166ac8fb56c438a2a0e1266f33b54dfe7b79f738d9945dbbbebf6e367c55"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.convertAddress
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "convert address success",
  "address": "0x2ec57B631580dF40d1E9e027360357eb61C7B25A"
}
```



## 3. validate address
### 输入
```bash
grpcurl -plaintext -d '{
  "chain": "ethereum",
  "network": "testnet",
  "address": "0x2ec57B631580dF40d1E9e027360357eb61C7B25A"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.validAddress
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "valid address",
  "valid": true
}
```


## 4.  block header by number
### 输入
```bash
grpcurl -plaintext -d '{
  "chain": "ethereum",
  "network": "testnet",
  "height": "990"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.getBlockHeaderByNumber
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "get latest block header success",
  "block_header": {
    "hash": "0xe02a4a9c136f627465ed63136325c7e943b527318f479b926bcc639cfaebcabc",
    "parent_hash": "0x220198524b4bd48c4b200bca3338ce2576d5cf00b96984ce243671ec64e47efd",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0x14627ea0e2B27b817DbfF94c3dA383bB73F8C30b",
    "root": "0x67e356c659d024e40d20fa4771846c43b079f3d0580647ec845c0692580453be",
    "tx_hash": "0xb96c36f85f16069347e623781648dc0bc322f4d7b72cbae725155cbe922e8847",
    "receipt_hash": "0x602093c8afd1ba1ee5477e3ef60464699b40f4580aa9b4d2a44650aeff8ec334",
    "parent_beacon_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "difficulty": "0",
    "number": "990",
    "gas_limit": "30000000",
    "gas_used": "28511072",
    "time": "1695917904",
    "extra": "d883010d01846765746888676f312e32312e31856c696e7578",
    "mix_digest": "0xef67112c590f7713e52f7fe1ceba4f456388b98a089d0b39b09a174fd716ed8f",
    "nonce": "0",
    "base_fee": "8",
    "withdrawals_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "blob_gas_used": "0",
    "excess_blob_gas": "0"
  }
}
```


## 5. block header by hash
### 输入
```bash
grpcurl -plaintext -d '{
  "chain": "ethereum",
  "network": "testnet",
  "hash": "0xc539f95dde2b26756ca6ce19aec527d8fa667752e6b495d7cfc8dbbc81c48e6c"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.getBlockHeaderByHash
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "get latest block header success",
  "block_header": {
    "hash": "0xc539f95dde2b26756ca6ce19aec527d8fa667752e6b495d7cfc8dbbc81c48e6c",
    "parent_hash": "0x87d6d2473ece7239ad34b7fa35c02f905fc5ab9bf79956e7100f60ac26bfe4f2",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0x0c10000000756BD1D14F9C837F022929A97bdA45",
    "root": "0x07ebf6c63b728e5da0e60760d03f0631bf9b14b1db8229121edb1f7037f3b6f1",
    "tx_hash": "0xdfe3bc7c40479e2c76a8b0f27c076f656529b701695e91903f8fd40085082cdb",
    "receipt_hash": "0xd68f0d4eee429bc9d7e347866fe21d7d5156ae27c392726b81101326ce81cea0",
    "parent_beacon_root": "0x988a7097183945e7c3ef9e478c0f560466a0fd8ab7ed6c7b92d921ae7c45d3ad",
    "difficulty": "0",
    "number": "3797728",
    "gas_limit": "36000000",
    "gas_used": "8128388",
    "time": "1746616164",
    "extra": "2IMBDwqEZ2V0aIhnbzEuMjQuMoVsaW51eA==",
    "mix_digest": "0x90c98005c1cf49ec85ef1f1e1417565eee7faa89af0b81b34edfb94abb4a4568",
    "nonce": "0",
    "base_fee": "329578295",
    "withdrawals_hash": "0x861c6bba8a0a94d2b530066386e76a5ad9994702d6c1825e1b621e2656fb92fd",
    "blob_gas_used": "131072",
    "excess_blob_gas": "0"
  }
}
```


## 6. get block by number
### 输入
```bash
grpcurl -plaintext -d '{
  "height": "999",
  "chain": "ethereum"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.getBlockByNumber
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "block by number success",
  "height": "999",
  "hash": "0x6f8d15acf8d5f05498c53aab3bec5eee172c3b48305ed43dadb43193a5e99bba",
  "base_fee": "0xb",
  "transactions": [
    {
      "from": "0xd3994e4d3202dd23c8497d7f75bf1647d1da1bb1",
      "to": "0xb7fb99e86f93dc3047a12932052236d853065173",
      "token_address": "0xb7fb99e86f93dc3047a12932052236d853065173",
      "contract_wallet": "0xb7fb99e86f93dc3047a12932052236d853065173",
      "hash": "0x4a880b9000ad4cce2fefb6f3bcadb46eb7eaca87de60c1a0287cc754f812c4af",
      "height": "999",
      "amount": "0xa968163f0a57b400000"
    },
    ...
    ...
    ...
    ]
}
```


## 7. get account
### 输入
```bash
grpcurl -plaintext -d '{
  "network": "testnet",
  "address": "0xeD2Eb97b84386Cb74C1232e833C8ca26FcD2e1b9",
  "chain": "ethereum",
  "contractAddress": "0x00"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.getAccount
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "get account response success",
  "network": "",
  "account_number": "0",
  "sequence": "24",
  "balance": "3809948361705892319"
}
```

## 8. get tx By address
### 输入
```bash
grpcurl -plaintext -d '{
  "address": "0xeD2Eb97b84386Cb74C1232e833C8ca26FcD2e1b9",
  "network": "testnet",
  "chain": "ethereum"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.getTxByAddress
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "get tx list success",
  "tx": [
    {
      "hash": "0xd8ce4809d87fb3115b16d66d82caaa756f5f220045ff26f7159b32ab7ea647ec",
      "index": 0,
      "from": "0x52f1984cd3e46e1214db222d3ff63712e7aceedd",
      "to": "0xed2eb97b84386cb74c1232e833c8ca26fcd2e1b9",
      "fee": "175156",
      "status": "TX_SUCCESS",
      "value": "1000000000000000000",
      "type": 1,
      "height": "3524935",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    ...
    ...
    ]
}
```


## 9. get tx by hash
### 输入
```bash
grpcurl -plaintext -d '{
  "chain": "ethereum",
  "network": "testnet",
  "hash": "0x0bd5ae9d4f4f48c19d63e44764f78d95c06038ec5c25ae738948e9830426b681"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.getTxByHash
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "get transaction success",
  "tx": {
    "hash": "0x0bd5ae9d4f4f48c19d63e44764f78d95c06038ec5c25ae738948e9830426b681",
    "index": 6,
    "from": "0x0000000000000000000000000000000000000000",
    "to": "0xcd9394575762Aa327e9a9E2bfD7d2348f070e567",
    "fee": "2651533355",
    "status": "TX_SUCCESS",
    "value": "144317870756000",
    "type": 0,
    "height": "3797942",
    "contract_address": "0x0000000000000000000000000000000000000000",
    "datetime": "",
    "data": ""
  }
}
```


## 10. todo build unsign tx
### 输入
```bash
grpcurl -plaintext -d '{
  "chain": "ethereum",
  "network": "testnet",
  "base64Tx": "eyJjaGFpbl9pZCI6IjE3MDAwIiwibm9uY2UiOjEsImZyb21fYWRkcmVzcyI6IjB4ZUQyRWI5N2I4NDM4NkNiNzRDMTIzMmU4MzNDOGNhMjZGY0QyZTFiOSIsInRvX2FkZHJlc3MiOiIweDcyZkZhQTI4OTk5M2JjYURhMkUwMTYxMjk5NUU1Yzc1ZEQ4MWNkQkMiLCJnYXNfbGltaXQiOjIxMDAwMCwiZ2FzIjoyMTAwMDAsIm1heF9mZWVfcGVyX2dhcyI6IjMwMDAwMDAwMDAwIiwibWF4X3ByaW9yaXR5X2ZlZV9wZXJfZ2FzIjoiMTUwMDAwMDAwMCIsImFtb3VudCI6IjEwMDAwMDAwMDAwMDAwMDAwMDAiLCJjb250cmFjdF9hZGRyZXNzIjoiMHgwMCJ9"
}' 127.0.0.1:8189 chainsunion.ChainsUnionService.buildUnSignTransaction
```
### 输出
```bash
{
  "code": "SUCCESS",
  "msg": "create un sign tx success",
  "un_sign_tx": "0xa4a269ecbb3b3ab0a57f7aa685d0d5d04edd91539cd0258d8282a38ea8570e47"
}
```

