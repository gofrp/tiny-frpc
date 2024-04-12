# tiny-frpc

# Introduction

Beginning with version v0.53.0, frp has supported SSH tunnel mode [SSH tunnel gateway](https://github.com/fatedier/frp?tab=readme-ov-file#ssh-tunnel-gateway). Users can communicate with frps using the standard SSH protocol, thereby accomplishing reverse proxying. This mode operates independently from the frpc binary.

Many users need to use reverse proxies on low-end computers or embedded machines. These machines may have limited memory and storage space, potentially preventing frpc from operating properly. This project aims to provide the simplest possible implementation of a reverse proxy, achieving reverse proxying via the SSH protocol in conjunction with frps.

We offer two binary programs: the native-SSH version (frpc-nssh) and the go-SSH version (frpc-gssh). Both of these programs parse frpc's toml file (conf/frpc_full_example.toml) and employ frps's communication to accomplish reverse proxying.

* The native-SSH version requires an SSH program on your machine; without one, it cannot function. This binary file is smaller.

* The go-SSH version does not rely on the SSH program; it's standalone, thus the binary file is larger.


# Usage

## 1. Download and uncompress the corresponding tiny-frpc version
Users decide whether to use the native-SSH version or the go-SSH version. Please visit the [releases](https://github.com/gofrp/tiny-frpc/releases) of this project to download.

## 2. Prepare the frpc toml file (There's a minimal usage configuration inside the decompressed package, for full configuration refer to this project's conf/frpc_full_example.toml)
Note: This project only supports the toml file format.

For example:
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

## 3. Run
> ./frpc-gssh -c frpc.toml

or

> ./frpc-nssh -c frpc.toml

This is all it takes to set up the reverse proxy.