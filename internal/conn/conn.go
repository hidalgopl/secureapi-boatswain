package conn

import (
	"fmt"
	"github.com/hidalgopl/secureapi-boatswain/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"time"
)

type NatsHandler struct {
	username       string
	password       string
	url            string
	totalWait      time.Duration
	reconnectDelay time.Duration
}

func (nh *NatsHandler) setupConnOptions(opts []nats.Option) []nats.Option {

	opts = append(opts, nats.ReconnectWait(nh.reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(nh.totalWait/nh.reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		logrus.Infof("Disconnected due to: %s, will attempt reconnects for %.0fm", err, nh.totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		logrus.Infof("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		logrus.Errorf("Exiting: %v", nc.LastError())
	}))
	return opts
}

func (nh *NatsHandler) Connect() (*nats.EncodedConn, error) {
	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = nh.setupConnOptions(opts)
	natsUrl := fmt.Sprintf("%s", nh.url)
	// Connect to NATS
	logrus.Infof("trying to connect to nats on %s", natsUrl)
	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		logrus.Error(err)
		return &nats.EncodedConn{}, err
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		logrus.Error(err)
		return &nats.EncodedConn{}, err
	}
	return ec, nil
}

func NewNatsHandler(conf *config.Config) *NatsHandler {
	conn := &NatsHandler{
		username:       conf.NatsUsername,
		password:       conf.NatsPass,
		url:            conf.NatsUrl,
		totalWait:      10 * time.Minute,
		reconnectDelay: time.Second,
	}
	return conn
}