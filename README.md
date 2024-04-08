# tiny-frpc

# Introduction

Starting from version v0.53.0, frp has supported the ssh tunnel mode [ssh tunnel gateway](https://github.com/fatedier/frp?tab=readme-ov-file#ssh-tunnel-gateway). Users can communicate with frps via the standard ssh protocol to complete the reverse proxy. This mode does not rely on the frpc binary.

Many users need to use reverse proxy in some low-configuration computers or embedded machines. These machines have limited memory and storage space and may not be able to use frpc normally. This project aims to provide the simplest version of the reverse proxy, which completes the reverse proxy with frps only through the ssh protocol.

We provide two binary programs, the native-ssh version (frpc-nssh) and the go-ssh version (frpc-gssh). Both programs parse the frpc toml file (conf/frpc_full_example.toml) and complete the reverse proxy with frps communication.

* The native-ssh version requires that there is an ssh program on your machine, otherwise, it cannot be used. This binary file is smaller.

* The go-ssh version does not depend on the ssh program; it is standalone. This binary file is larger.


# Usage

## 1. Prepare frpc toml file (refer to conf/frpc_full_example.toml)
Note: Only the toml file format is supported.

## 2. Download the corresponding version of tiny-frpc
Users decide whether to use the native-ssh version or the go-ssh version. Please move to the releases of this project to download.

## 3. Run
> ./frpc-xxx -c frpc.toml

And you're done, completing the reverse proxy.