package xmppc

import (
	"crypto/tls"
	"ejabberd_client_go/config"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

type XmppConn struct {
	Name     string
	Password string

	client *xmpp.Client
}

// 初始化Conn
func NewXmppConn(name, password string) *XmppConn {
	return &XmppConn{
		Name:     name,
		Password: password,
	}
}

func (xc *XmppConn) Connect() bool {
	logName := fmt.Sprintf("./logs/%v-xmpp.log", xc.Name)
	file, _ := os.OpenFile(logName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	// 创建登陆用户的配置``
	config := &xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address:   fmt.Sprintf("%v:%v", config.XMPP_DOMAIN, config.XMPP_PORT),
			TLSConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jid:        fmt.Sprintf("%v@%v", xc.Name, config.XMPP_DOMAIN),
		Credential: xmpp.Password(xc.Password),
		// StreamLogger: os.Stdout,
		StreamLogger: file,
		Insecure:     true,
	}

	// 绑定路由的回掉函数
	router := xmpp.NewRouter()
	router.HandleFunc("message", handleMessage)

	// 创建客户端
	client, err := xmpp.NewClient(config, router, errorHandler)
	if err != nil {
		logrus.Error("Xmpp connect init failed: ", err.Error())
		return false
	}

	// ***连接IM
	err = client.Connect()
	if err != nil {
		logrus.Error("Xmpp connect failed: ", err.Error())
		return false
	}

	logrus.Info("Xmpp connect successed --- with name: ", xc.Name)
	xc.client = client
	return true
}

func handleMessage(s xmpp.Sender, p stanza.Packet) {
	msg, ok := p.(stanza.Message)

	if !ok {
		_, _ = fmt.Fprintf(os.Stdout, "Ignoring packet: %T\n", p)
		return
	}

	// 输出接收到的信息到到os.Stdout中
	logrus.Info(time.Now().Format("2006-01-02 15:04:05 -0700"), "Receive message from:", msg.From, " message: ", msg.Body)
}

func errorHandler(err error) {
	fmt.Println(err.Error())
}

func (xc *XmppConn) SendMsg(msg, to string) {
	packet := stanza.NewMessage(
		stanza.Attrs{
			Type: stanza.MessageTypeChat,
			Id:   strconv.FormatInt(time.Now().Unix(), 10),
			From: fmt.Sprintf("%v@%v", xc.Name, config.XMPP_DOMAIN),
			To:   fmt.Sprintf("%v@%v", to, config.XMPP_DOMAIN),
			Lang: "en",
		})
	packet.Body = msg

	err := xc.client.Send(packet)
	if err != nil {
		logrus.Error("Send failed", err.Error())
		return
	}
	logrus.Info("Send ok")
}
