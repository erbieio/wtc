 ## 环境

1. 安装nodejs（带npm）
2. 安装go
 
 ## 快速开始

 ### 安装ethereum
 ```
 go install github.com/ethereum/go-ethereum/cmd/geth@v1.10.16
 ```

### 启动虚拟节点
``` 
cd node
./restart.sh
```

### 启动服务
``` 
cd service
go run .
```