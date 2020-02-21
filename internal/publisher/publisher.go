package publisher

import (
	"github.com/hidalgopl/secureapi-boatswain/internal/config"
	"github.com/hidalgopl/secureapi-boatswain/internal/conn"
	"github.com/sirupsen/logrus"
)

type Publisher interface {
	Publish(msg interface{}, subject string) error
}

type NatsPublisher struct {
	*conn.NatsHandler
	SubjectName string
}

func NewNatsPublisher(conf *config.Config, subject string) *NatsPublisher {
	nh := conn.NewNatsHandler(conf)
	nl := &NatsPublisher{
		NatsHandler: nh,
		SubjectName: subject,
	}
	return nl
}

func (np *NatsPublisher) Publish(msg interface{}, subject string) error {
	ec, err := np.Connect()
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("publishing to: %s", subject)
	err = ec.Publish(subject, msg)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
