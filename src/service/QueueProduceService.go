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

type textResSend struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Id   int    `json:"id"`
}

type ProduceService struct {
	nats       *nats.Conn
	stan       stan.Conn
	logger     *zap.Logger
	properties *properties.ServiceProperties
	dbms       *dbms.DatabaseRepository
}

func (this *ProduceService) Init(properties *properties.ServiceProperties) {
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
	this.logger_produce()
	this.nats, err = nats.Connect(this.properties.Nats.Address)
	if err != nil {
		this.logger.Error(err.Error())
		os.Exit(1)
	}
	this.stan, err = stan.Connect(this.properties.Nats.ClusterId, this.properties.Nats.ClientId, stan.NatsURL(this.properties.Nats.Address))
	if err != nil {
		this.logger.Info("init", zap.Error(err))
	}
	this.dbms = dbms.New(this.properties.Database)
	err = this.dbms.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	this.dbms.Database.LogMode(true)
	subscription2, _ := this.nats.QueueSubscribe(this.properties.Topic.GetChannel("req_produce_message"), this.properties.Topic.GetQueueName("req_produce_message"), this.produceMessage)

	this.logger.Info("subscribtion", zap.String("subject", subscription2.Subject), zap.Bool("status", subscription2.IsValid()))

}

func (this *ProduceService) logger_produce() {
	var logoutput []string
	var logerror []string
	t := time.Now()
	s := t.Format("2006-01-02")
	logoutput = append(logoutput, this.properties.Logging.OutputPaths[0]+this.properties.NameLogging+"-"+s+".log", this.properties.Logging.OutputPaths[1])
	logerror = append(logerror, this.properties.Logging.ErrorOutputPaths[0]+this.properties.NameLogging+"-"+s+".log", this.properties.Logging.ErrorOutputPaths[1])
	this.logger, _ = tools.LoggerGenerator(logoutput, logerror)
}

func (this *ProduceService) produceMessage(msg *nats.Msg) {
	this.logger_produce()
	var response = models.Response{Status: "failed"}
	request := &models.RequestProduce{}
	this.logger.Info("produce_message", zap.String("subject", msg.Subject), zap.String("raw_request", string(msg.Data)))
	if err := json.Unmarshal(msg.Data, &request); err == nil {
		var isValid bool
		var id int
		passwd := fmt.Sprintf("%x", sha256.Sum256([]byte(request.Password)))
		row, err := this.dbms.Database.Raw("select * from func_produce_message (?, ?, ?, ?)", request.Name, request.Message, request.Username, passwd).Rows()

		defer row.Close()
		if err != nil {
			response.Error = err.Error()
		} else {
			for row.Next() {
				row.Scan(&isValid, &id)
			}
			if isValid {
				response.Status = "success"
				var resp = textResSend{Name: request.Name, Text: request.Message, Id: id}
				bytes1, _ := json.Marshal(resp)
				fmt.Println(resp)
				fmt.Println(string(bytes1))
				topic := fmt.Sprintf("%x", sha256.Sum256([]byte(request.Name)))
				this.stan.Publish(this.properties.Topic.GetChannel("res_consume_message")+topic, bytes1)
				this.logger.Info("consume_message_send", zap.String("subject", this.properties.Topic.GetChannel("res_consume_message")+topic), zap.String("raw_request", string(msg.Data)), zap.Any("response", request.Message))
			} else {
				response.Status = "failed"
			}
		}
	} else {
		response.Error = "unrecognize request"
	}

	this.logger.Info("produce_message", zap.String("subject", msg.Subject), zap.String("raw_request", string(msg.Data)), zap.Any("response", response))
	bytes, _ := json.Marshal(response)
	this.nats.Publish(msg.Reply, bytes)
}
