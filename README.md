"""
# Tiny-FRPC

# Introduction

Starting from version v0.53.0, frp supports the ssh tunnel mode [ssh tunnel gateway](https://github.com/fatedier/frp?tab=readme-ov-file#ssh-tunnel-gateway). Users can communicate with frps through the standard ssh protocol to complete the reverse proxy. This mode can function independently from the frpc binary.

Many users need to use reverse proxy in microsystems where memory and storage space are limited, which may prevent the normal use of frpc. The purpose of this project is to provide a minimal version of the reverse proxy, which communicates with frps through the ssh protocol only.

We provide two types of binary programs - the native-ssh version and the go-ssh version. Both programs parse the [standard file format of frpc](https://github.com/fatedier/frp/blob/dev/conf/frpc_full_example.toml), and communicate with frps to complete the reverse proxy.

* The native-ssh version requires that your machine already have an ssh program, otherwise it won't work. The binary file is smaller.

* The go-ssh version does not depend on an ssh program and is standalone, thus the binary file is larger.

# Usage

## 1. Prepare the frpc configuration file according to the [standard file format of frpc](https://github.com/fatedier/frp/blob/dev/conf/frpc_full_example.toml)
Note: Only the toml file format is supported.

## 2. Download the corresponding version of tiny-frpc
Users decide whether to use the native-ssh version or the go-ssh version.

## 3. Run
> ./tiny-frpc -c frpc.toml

And that's it - you've completed your reverse proxy.
"""