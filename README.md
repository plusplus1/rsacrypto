# RSA 多机并行解密服务

## 一、HTTPServer,

- 功能说明：
    - 中心服务，负责对外提供加/解密操作 http 接口

- 启动:
    - `bin/center`

- 主要接口说明:
    - **`/rsa/encrypt`** ：加密操作
        - 参数:
            - `data`-*string* , 待加密的明文；
            - `pub_key`-*string* , RSA公钥,PEM格式；

        - 返回：
            - `HTTP状态码` = `200` 表示执行成功；
            - `HTTP Body` 成功时代表加密后的密文，否则表示错误信息；

    - **`/rsa/decrypt`** ：解密操作
        - 参数:
            - `data`-*string* , 待解密的密文；
            - `pri_key`-*string* ,RSA私钥,PEM格式；

        - 返回：
            - `HTTP状态码` = `200` 表示执行成功；
            - `HTTP Body` 成功时代表解密后的明文，否则表示错误信息；


---

## 二、GRPCServer

- 功能职责：
    - Worker 节点具体执行加/解密任务

- 启动：
    - `sbin/worker`

- 数据交互协议:
    - pb
        - [commonLib/rpcLib/rpc.proto](commonLib/rpcLib/rpc.proto)

-

