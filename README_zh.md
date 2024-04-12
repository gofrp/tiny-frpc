# tiny-frpc

# 简介

frp 在 >= v0.53.0 版本已经支持 ssh tunnel 模式 [ssh tunnel gateway](https://github.com/fatedier/frp?tab=readme-ov-file#ssh-tunnel-gateway)。 用户可以通过标准的 ssh 协议来跟 frps 进行通信，完成反向代理。该模式可以不依赖于 frpc 二进制。

很多用户需要在一些低配计算机或嵌入式设备里使用反向代理，这些机器内存和存储空间有限，可能无法正常使用原版的 frpc。本项目旨在提供最简版本的反向代理，只通过 ssh 协议与 frps 完成反向代理。

我们提供了2种二进制程序，native ssh 的版本(tiny-frpc-ssh)和 standalone 版本(tiny-frpc)。2种程序都是解析 frpc toml 文件(conf/frpc_full_example.toml)，与 frps 的通信完成反向代理。

* native ssh 版本需要你本机已经有 ssh 程序，否则无法使用。 该二进制文件较小。

* standalone 版本不依赖本机的 ssh 程序。该二进制文件较大。


# 使用

## 1. 下载对应的 tiny-frpc 版本并解压
用户自己决定是用 native ssh 版本还是 standalone 版本。请移步到本项目的 [releases](https://github.com/gofrp/tiny-frpc/releases) 下载。


## 2. 准备 frpc toml 文件 (压缩包解压之后里面有最简使用配置，完整配置参考本项目 conf/frpc_full_example.toml）
注意：本项目只支持 toml 文件格式。

举个例子：
```
serverAddr = "127.0.0.1"

# frps ssh tunnel gateway port
serverPort = 2200

[[proxies]]
name = "test-tcp"
type = "tcp"
localIP = "127.0.0.1"
localPort = 22
remotePort = 6000
```

## 3. 运行
> ./tiny-frpc -c frpc.toml

or

> ./tiny-frpc-ssh -c frpc.toml

即可在完成反向代理。


# 说明

**当前是预览版本，不保证兼容性，目前仅供测试使用，不要用于生产环境!!!**
