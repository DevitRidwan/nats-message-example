package service

import (
	"Ridwan/Queue/src/dbms"
	"Ridwan/Queue/src/logger"
	"Ridwan/Queue/src/models"
	"Ridwan/Queue/src/properties"
	"Ridwan/Queue/src/tools"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/go-nats"
	stan "github.com/nats-io/go-nats-streaming"
	"go.uber.org/zap"
)

type textReq struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Id   int    `json:"id"`
}

type textRes struct {
	Text string `json:"text"`
}

type ConsumeService struct {
	nats       *nats.Conn
	stan       stan.Conn
	logger     *zap.Logger
	properties *properties.ServiceProperties
	dbms       *dbms.DatabaseRepository
}

func (this *ConsumeService) Init(properties *properties.ServiceProperties) {
	var err error
	this.properties = properties
	exists, _ := logger.DirExists(this.properties.Logging.OutputPaths[0])
	if exists == false {
		merr := os.MkdirAll(this.properties.Logging.OutputPaths[0], os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
	exists, _ = logger.DirExists(this.properties.Logging.ErrorOutputPaths[0])
	if exists == false {
		merr := os.MkdirAll(this.properties.Logging.ErrorOutputPaths[0], os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
	this.logger_consume()
	this.nats, err = nats.Connect(this.properties.Nats.Address)
	if err != nil {
		this.logger.Error(err.Error())
		os.Exit(1)
	}
	this.stan, err = stan.Connect(this.properties.Nats.ClusterId, this.properties.Nats.ClientId+"-record", stan.NatsURL(this.properties.Nats.Address))

	if err != nil {
		this.logger.Info("init", zap.Error(err))
	}
	this.dbms = dbms.New(this.properties.Database)
	err = this.dbms.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	this.dbms.Database.LogMode(true)
	subscription, _ := this.nats.QueueSubscribe(this.properties.Topic.GetChannel("req_auth_message"), this.properties.Topic.GetQueueName("req_auth_message"), this.authConsume)

	this.logger.Info("subscribtion", zap.String("subject", subscription.Subject), zap.Bool("status", subscription.IsValid()))

}

func (this *ConsumeService) logger_consume() {
	var logoutput []string
	var logerror []string
	t := time.Now()
	s := t.Format("2006-01-02")
	logoutput = append(logoutput, this.properties.Logging.OutputPaths[0]+this.properties.NameLogging+"-"+s+".log", this.properties.Logging.OutputPaths[1])
	logerror = append(logerror, this.properties.Logging.ErrorOutputPaths[0]+this.properties.NameLogging+"-"+s+".log", this.properties.Logging.ErrorOutputPaths[1])
	this.logger, _ = tools.LoggerGenerator(logoutput, logerror)
}

func (this *ConsumeService) authConsume(msg *nats.Msg) {
	this.logger_consume()
	var response = models.Response{Status: "failed"}
	request := &models.RequestQueue{}
	this.logger.Info("consume_auth", zap.String("subject", msg.Subject), zap.String("raw_request", string(msg.Data)))
	if err := json.Unmarshal(msg.Data, &request); err == nil {
		var isValid bool
		passwd := fmt.Sprintf("%x", sha256.Sum256([]byte(request.Password)))
		err = this.dbms.Database.Raw("select true from tbl_user where username = ? and password = ?", request.Username, passwd).Row().Scan(&isValid)
		if err != nil {
			this.logger.Info("consume_auth", zap.String("error", err.Error()))
			response.Error = "unathorize"
		} else if isValid {
			response.Status = "success"
			topic := fmt.Sprintf("%x", sha256.Sum256([]byte(request.Name)))
			this.stan.Subscribe(this.properties.Topic.GetChannel("req_consume_message")+topic, this.consumeMessage, stan.StartWithLastReceived(), stan.DurableName(topic))
		}
	} else {
		response.Error = "unrecognize request"
	}
	this.logger.Info("consume_auth", zap.String("subject", msg.Subject), zap.String("raw_request", string(msg.Data)), zap.Any("response", response))
	bytes, _ := json.Marshal(response)
	this.nats.Publish(msg.Reply, bytes)
}

func (this *ConsumeService) consumeMessage(msg *stan.Msg) {
	this.logger_consume()
	var request textReq
	this.logger.Info("consume_message", zap.String("subject", msg.Subject), zap.String("raw_request", string(msg.Data)))
	if err := json.Unmarshal(msg.Data, &request); err == nil {
		var isValid bool
		var resp = textRes{Text: request.Text}
		topic := fmt.Sprintf("%x", sha256.Sum256([]byte(request.Name)))
		this.logger.Info("consume_message_send", zap.String("subject", topic), zap.String("raw_request", string(msg.Data)), zap.Any("response", resp))
		bytes, _ := json.Marshal(resp)
		this.nats.Publish(topic, bytes)
		this.nats.FlushTimeout(30 * time.Second)
		this.dbms.Database.Raw("select * from func_consume_message(?)", request.Id).Row().Scan(&isValid)
	} else {
		fmt.Println(err)
	}
}
