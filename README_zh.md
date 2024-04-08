# tiny-frpc

# 简介

frp 在 >= v0.53.0 版本已经支持 ssh tunnel 模式 [ssh tunnel gateway](https://github.com/fatedier/frp?tab=readme-ov-file#ssh-tunnel-gateway)。 用户可以通过标准的 ssh 协议来跟 frps 进行通信，完成反向代理。该模式可以不依赖于 frpc 二进制。

很多用户需要在一些低配计算器或者嵌入式机器里使用反向代理，这些机器内存和存储空间有限，可能无法正常使用 frpc。本项目旨在提供最简版本的反向代理，只通过 ssh 协议与 frps 完成反向代理。

我们提供了2种二进制程序，native-ssh 的版本(frpc-nssh)和 go-ssh 版本(frpc-gssh)。2种程序都是解析 frpc toml 文件(conf/frpc_full_example.toml)，与 frps 的通信完成反向代理。

* native-ssh 版本需要你本机已经有 ssh 程序，否则无法使用。 该二进制文件较小。

* go-ssh 版本不依赖 ssh 程序，standalone，该二进制文件较大。


# 使用

## 1. 准备 frpc toml 文件 (参考 conf/frpc_full_example.toml）
注意：只支持 toml 文件格式。

## 2. 下载对应的 tiny-frpc 版本
用户自己决定是用 native-ssh 版本还是 go-ssh 版本。请移步到本项目的 releases 下载。 

## 3. 运行
> ./frpc-xxx -c frpc.toml

即可在完成反向代理。