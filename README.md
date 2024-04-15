# tiny-frpc

# Introduction

As of version v0.53.0, frp has supported the ssh tunnel mode [ssh tunnel gateway](https://github.com/fatedier/frp?tab=readme-ov-file#ssh-tunnel-gateway). Users can communicate with frps using the standard ssh protocol to accomplish reverse proxying. This mode operates independently from the frpc binary.

Many users need to use reverse proxying on low-end computers or embedded devices. These machines have limited memory and storage space, and as a result, may struggle to operate the original version of frpc normally. This project aims to provide the simplest possible implementation of a reverse proxy, achieving reverse proxying solely through the ssh protocol in conjunction with frps.

We offer two binary programs: the native ssh version (tiny-frpc-ssh) and the standalone version (tiny-frpc). Both of these programs parse the frpc's toml file (conf/frpc_full_example.toml) to complete reverse proxying via communication with frps.

* The native ssh version requires that your machine already has an ssh program installed; otherwise, it cannot be used. This binary file is smaller.

* The standalone version does not rely on any ssh program on your local machine. Consequently, this binary file is larger.


# Usage

## 1. Download tiny-frpc
Users decide whether to use the native ssh version or the standalone version. Please navigate to this project's [releases](https://github.com/gofrp/tiny-frpc/releases) to download.


## 2. Prepare the frpc toml file
Note: This project only supports the toml file format.
After decompressing the package, there will be a minimal usage configuration. For a full configuration, refer to this project's conf/frpc_full_example.toml.

For example, the minimal configuration of frps is:
```
# frps.toml

bindPort = 7000

vhostHTTPPort = 80

sshTunnelGateway.bindPort = 2200
```

And the configuration for tiny-frpc is:
```
# frpc.toml

serverAddr = "x.x.x.x"

# frps ssh tunnel gateway port
serverPort = 2200

[[proxies]]
name = "test-tcp-server"
type = "tcp"
localIP = "127.0.0.1"
localPort = 5000
remotePort = 6000

[[proxies]]
name = "test-http-web"
type = "http"
localIP = "127.0.0.1"
localPort = 7080
customDomains = ["test-tiny-frpc.frps.com"]
locations = ["/", "/pic"]
```

## 3. Running
> ./tiny-frpc -c frpc.toml

or

> ./tiny-frpc-ssh -c frpc.toml

For TCP service:

> nc -zv x.x.x.x 6000

or

> telnet x.x.x.x 6000

You can test whether the intranet port has been successfully proxying to the public network. If it has, you can access the TCP service from the public network.

For HTTP service:

Assuming that the domain name 'test-tiny-frpc.frps.com' is resolved to the machine where frps resides, you can access the intranet's HTTP services through:

> curl -v 'http://test-tiny-frpc.frps.com/'



# Principle
![how tiny frpc works](doc/pic/architecture.png)


# Disclaimer

**This is currently a preview version. Compatibility is not guaranteed. It is presently for testing purposes only and should not be used in production environments!**