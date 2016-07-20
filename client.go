package apns

import (
	"crypto/tls"
	"errors"
	"net"
	"strings"
)

const (
	ApnsProductionGateway = "gateway.push.apple.com:2195"
	ApnsSandboxGateway    = "gateway.sandbox.push.apple.com:2195"
)

/*
从p12导出证书命令
openssl pkcs12 -in path.p12 -out newfile.crt.pem -clcerts -nokeys
openssl pkcs12 -in path.p12 -out newfile.key.pem -nocerts -nodes
*/

type ApnsClient struct {
	Gateway string
	Conn    *tls.Conn
	Conf    *tls.Config
}

func NewClient(gateway, crtFile, keyFile string) *ApnsClient {
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		panic(err)
	}
	c := new(ApnsClient)
	c.Gateway = gateway
	c.Conn = nil
	gatewayParts := strings.Split(c.Gateway, ":")
	c.Conf = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   gatewayParts[0],
	}
	return c
}

//发送队列中的所有通知
func (this *ApnsClient) Close() error {
	return this.Conn.Close()
}

//发送队列中的所有通知
func (this *ApnsClient) Write(pn *Notification) error {
	b, err := pn.ToBytes()
	if err != nil {
		return err
	}
	l := len(b)
	s := 0
	for {
		n, err := this.Conn.Write(b[s:])
		if err != nil {
			return err
		}
		s += n
		if s >= l {
			break
		}
	}
	return nil
}

//连接到apns服务器
func (this *ApnsClient) Connect() error {
	conn, err := net.Dial("tcp", this.Gateway)
	if err != nil {
		return err
	}
	this.Conn = tls.Client(conn, this.Conf)
	if this.Conn == nil {
		return errors.New("make tls client error")
	}
	return this.Conn.Handshake()
}
