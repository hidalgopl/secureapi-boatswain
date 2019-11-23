package listener

import (
	"crypto/tls"
	"fmt"
	"github.com/hidalgopl/secureapi-boatswain/internal/checks"
	"github.com/hidalgopl/secureapi-boatswain/internal/config"
	"github.com/hidalgopl/secureapi-boatswain/internal/conn"
	"github.com/hidalgopl/secureapi-boatswain/internal/messages"
	"github.com/hidalgopl/secureapi-boatswain/internal/publisher"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Listener interface {
	Listen(conf *config.Config)
}

type NatsListener struct {
	*conn.NatsHandler
	Subject   string
	QueueName string
}

func NewNatsListener(conf *config.Config) *NatsListener {
	nh := conn.NewNatsHandler(conf)
	nl := &NatsListener{
		NatsHandler: nh,
		Subject:     conf.NatsCreatedSubject,
		QueueName:   conf.NatsQueueName,
	}
	return nl

}

func Listen(conf *config.Config) error {
	// Connect Options.
	nh := NewNatsListener(conf)
	ec, err := nh.Connect()
	if err != nil {
		log.Fatal(err)
		return err
	}
	nh.QueueSub(ec)
	return nil
	//defer ec.Close()

}

func (nh *NatsListener) QueueSub(ec *nats.EncodedConn) {
	ec.QueueSubscribe(
		nh.Subject, nh.QueueName, nh.HandleTestSuite)
	ec.Flush()
	defer ec.Close()

	if err := ec.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s] in queue [%s]", nh.Subject, nh.QueueName)

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	ec.Drain()
	log.Fatalf("Exiting")
}

func (nh *NatsListener) HandleTestSuite(msg *messages.StartTestSuitePub) {
	log.Printf("Got msg: %s", msg.Print())
	testSuiteUID := msg.TestSuiteID
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	conf := config.GetConf()
	client := &http.Client{Transport: tr}
	logrus.Infof("creating result chan  with len: %v", len(msg.Tests))
	resultChan := make(chan messages.TestFinishedPub, len(msg.Tests))
	logrus.Infof("sending http req to: %s", msg.Url)
	r, err := client.Get(msg.Url)
	defer r.Body.Close()
	if err != nil {
		panic(err)
		// TODO
	}
	logrus.Info("start scheduling checks:")
	for ind, testCode := range msg.Tests {
		subject := fmt.Sprintf("test_suite.%s.test.%s.started", testSuiteUID, ind)
		pub := publisher.NewNatsPublisher(conf, subject)
		go func(testCode string) {
			// handle if someone sends wrong testCode
			checks.TestCodes[testCode](testSuiteUID, r.Header, resultChan, pub)
		}(testCode)

	}
	logrus.Info("finished scheduling tests, waiting for chan to be closed")
	tests := []messages.TestFinishedPub{}
	for range msg.Tests {
		tests = append(tests, <-resultChan)
		//fmt.Println(<-resultChan)
	}
	close(resultChan)
	finishedSubject := fmt.Sprintf("test_suite.%s.completed", testSuiteUID)
	logrus.Infof("finished subjects: %s", finishedSubject)
	pub := publisher.NewNatsPublisher(conf, finishedSubject)
	logrus.Info("set up a publisher")
	finishedMsg := messages.TestSuiteFinishedPub{
		Timestamp:   time.Now(),
		Tests:       tests,
		TestSuiteID: testSuiteUID,
		Url:         msg.Url,
	}
	logrus.Infof("created finish msg: %v", finishedMsg)
	err = pub.Publish(finishedMsg, finishedSubject)

	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("published successfully")

}
