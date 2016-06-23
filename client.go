package apns

import (
	"crypto/tls"
	"errors"
	"log"
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
	Pool    []*Notification
}

func NewClient(gateway, crtFile, keyFile string) *ApnsClient {
	c := new(ApnsClient)
	c.Gateway = gateway
	c.Conn = nil
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		return nil
	}
	gatewayParts := strings.Split(c.Gateway, ":")
	c.Conf = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   gatewayParts[0],
	}
	c.Pool = []*Notification{}
	return c
}

//发送队列中的所有通知
func (this *ApnsClient) Send() error {
	err := this.connect()
	if err != nil {
		return err
	}
	for _, d := range this.Pool {
		b, err := d.ToBytes()
		if err != nil {
			log.Println("apns pool tobytes error", err)
			continue
		}
		this.Conn.Write(b)
	}
	return this.Conn.Close()
}

//加入对列
func (this *ApnsClient) Append(pn *Notification) {
	this.Pool = append(this.Pool, pn)
}

//连接到apns服务器
func (this *ApnsClient) connect() error {
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
