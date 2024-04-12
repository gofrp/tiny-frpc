package gssh

import (
	"net"

	"golang.org/x/crypto/ssh"

	"github.com/gofrp/tiny-frpc/pkg/util"
	"github.com/gofrp/tiny-frpc/pkg/util/log"
)

type TunnelClient struct {
	localAddr string
	sshServer string
	commands  string

	sshConn *ssh.Client
	ln      net.Listener
}

func NewTunnelClient(localAddr string, sshServer string, commands string) *TunnelClient {
	return &TunnelClient{
		localAddr: localAddr,
		sshServer: sshServer,
		commands:  commands,
	}
}

func (c *TunnelClient) Start() error {
	config := &ssh.ClientConfig{
		User:            "v0",
		HostKeyCallback: func(string, net.Addr, ssh.PublicKey) error { return nil },
	}

	conn, err := ssh.Dial("tcp", c.sshServer, config)
	if err != nil {
		return err
	}
	c.sshConn = conn

	l, err := conn.Listen("tcp", "0.0.0.0:80")
	if err != nil {
		return err
	}
	c.ln = l
	ch, req, err := conn.OpenChannel("session", []byte(""))
	if err != nil {
		return err
	}
	defer ch.Close()
	go ssh.DiscardRequests(req)

	type command struct {
		Cmd string
	}
	_, err = ch.SendRequest("exec", false, ssh.Marshal(command{Cmd: c.commands}))
	if err != nil {
		return err
	}

	c.serveListener()
	return nil
}

func (c *TunnelClient) Close() {
	if c.sshConn != nil {
		_ = c.sshConn.Close()
	}
	if c.ln != nil {
		_ = c.ln.Close()
	}
}

func (c *TunnelClient) serveListener() {
	for {
		conn, err := c.ln.Accept()
		if err != nil {
			log.Errorf("ssh tunnel cient accept error: %v", err)
			return
		}
		go c.hanldeConn(conn)
	}
}

func (c *TunnelClient) hanldeConn(conn net.Conn) {
	defer conn.Close()
	local, err := net.Dial("tcp", c.localAddr)
	if err != nil {
		log.Errorf("ssh tunnel client dial %v error: %v", c.localAddr, err)
		return
	}
	_, _ = util.Join(local, conn)
}
