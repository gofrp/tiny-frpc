# tiny-frpc

# Introduction

Starting from version v0.53.0, frp has supported the ssh tunnel mode [ssh tunnel gateway](https://github.com/fatedier/frp?tab=readme-ov-file#ssh-tunnel-gateway). Users can communicate with frps using the standard ssh protocol to accomplish reverse proxying. This mode operates independently from the frpc binary.

Many users need to use reverse proxying on low-end computers or embedded devices. These pieces of hardware have limited memory and storage space and may struggle to operate the original version of frpc normally. This project aims to provide the simplest version of a reverse proxy, which communicates with frps via the ssh protocol to complete reverse proxying.

We offer two binary programs: the native ssh version (tiny-frpc-ssh) and the standalone version (tiny-frpc). Both these programs parse the frpc toml file (conf/frpc_full_example.toml) and complete reverse proxying via communication with frps.

* The native ssh version requires an ssh program on your machine, otherwise, it cannot be used. This binary file is smaller.

* The standalone version does not rely on any ssh program installed on your machine. This binary file is larger.


# Usage

## 1. Download and uncompress the tiny-frpc version that corresponds to your needs
Users have the choice between the native ssh version or the standalone version. Please navigate to this project's [releases](https://github.com/gofrp/tiny-frpc/releases) to download.

## 2. Prepare the frpc toml file (within the decompressed package you'll find a simple usage configuration, for full configuration refer to conf/frpc_full_example.toml in this project)
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
> ./tiny-frpc -c frpc.toml

or

> ./tiny-frpc-ssh -c frpc.toml

And just like that, the reverse proxy is set up.


# Disclaimer

**This is currently a preview version. Compatibility is not guaranteed. It is currently intended for testing purposes only and should not be used in production environments!!!**