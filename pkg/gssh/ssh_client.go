package gssh

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"github.com/gofrp/tiny-frpc/pkg/util"
	"github.com/gofrp/tiny-frpc/pkg/util/log"
)

type TunnelClient struct {
	localAddr string
	sshServer string
	command   string

	sshConn *ssh.Client
	ln      net.Listener

	authMethod ssh.AuthMethod
}

func getDefaultPrivateKeyPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".ssh", "id_rsa"), nil
}

func publicKeyAuthFunc(kPath string) (ssh.AuthMethod, error) {
	key, err := os.ReadFile(kPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	return ssh.PublicKeys(signer), nil
}

func NewTunnelClient(localAddr string, sshServer string, command string) (*TunnelClient, error) {
	privateKeyPath, err := getDefaultPrivateKeyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get default private key path: %v", err)
	}

	log.Infof("get ssh private key file: [%v] to communicate with frps by ssh protocol", privateKeyPath)

	authMethod, err := publicKeyAuthFunc(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth method: %v", err)
	}

	return &TunnelClient{
		localAddr:  localAddr,
		sshServer:  sshServer,
		command:    command,
		authMethod: authMethod,
	}, nil
}

func (c *TunnelClient) Start() error {
	config := &ssh.ClientConfig{
		User:            "v0",
		Auth:            []ssh.AuthMethod{c.authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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

	session, err := c.sshConn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Start(c.command)
	if err != nil {
		return err
	}

	log.Infof("session start cmd [%v] success", c.command)

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

		log.Infof("accept a new connection. remote: %v, local: %v", conn.RemoteAddr().String(), conn.LocalAddr().String())

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
