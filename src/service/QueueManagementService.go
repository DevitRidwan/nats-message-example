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
	"go.uber.org/zap"
)

type QueueManagement struct {
	nats       *nats.Conn
	logger     *zap.Logger
	properties *properties.ServiceProperties
	dbms       *dbms.DatabaseRepository
}

func (this *QueueManagement) Init(properties *properties.ServiceProperties) {
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
	this.logger_queue()
	this.nats, err = nats.Connect(this.properties.Nats.Address)
	if err != nil {
		this.logger.Error(err.Error())
		os.Exit(1)
	}
	this.dbms = dbms.New(this.properties.Database)
	err = this.dbms.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	this.dbms.Database.LogMode(true)
	subscription, _ := this.nats.QueueSubscribe(this.properties.Topic.GetChannel("req_create_queue"), this.properties.Topic.GetQueueName("req_create_queue"), this.createQueue)
	subscription1, _ := this.nats.QueueSubscribe(this.properties.Topic.GetChannel("req_delete_queue"), this.properties.Topic.GetQueueName("req_delete_queue"), this.deleteQueue)
	this.logger.Info("subscribtion", zap.String("subject", subscription.Subject), zap.Bool("status", subscription.IsValid()))
	this.logger.Info("subscribtion", zap.String("subject", subscription1.Subject), zap.Bool("status", subscription1.IsValid()))

}

func (this *QueueManagement) logger_queue() {
	var logoutput []string
	var logerror []string
	t := time.Now()
	s := t.Format("2006-01-02")
	logoutput = append(logoutput, this.properties.Logging.OutputPaths[0]+this.properties.NameLogging+"-"+s+".log", this.properties.Logging.OutputPaths[1])
	logerror = append(logerror, this.properties.Logging.ErrorOutputPaths[0]+this.properties.NameLogging+"-"+s+".log", this.properties.Logging.ErrorOutputPaths[1])
	this.logger, _ = tools.LoggerGenerator(logoutput, logerror)
}

func (this *QueueManagement) createQueue(msg *nats.Msg) {
	this.logger_queue()
	var response = models.Response{Status: "failed"}
	request := &models.RequestQueue{}
	this.logger.Info("create_queue", zap.String("subject", msg.Subject), zap.String("raw_request", string(msg.Data)))
	if err := json.Unmarshal(msg.Data, &request); err == nil {
		var isValid bool
		passwd := fmt.Sprintf("%x", sha256.Sum256([]byte(request.Password)))
		row, err := this.dbms.Database.Raw("select * from func_create_queue (?, ?, ?)", request.Name, request.Username, passwd).Rows()

		defer row.Close()
		if err != nil {
			response.Error = err.Error()
		} else {
			for row.Next() {
				row.Scan(&isValid)
				if isValid {
					response.Status = "success"
				} else {
					response.Status = "failed"
				}
			}
		}
	} else {
		response.Error = "unrecognize request"
	}

	this.logger.Info("create_queue", zap.String("subject", msg.Subject), zap.String("raw_request", string(msg.Data)), zap.Any("response", response))
	bytes, _ := json.Marshal(response)
	this.nats.Publish(msg.Reply, bytes)
}

func (this *QueueManagement) deleteQueue(msg *nats.Msg) {
	this.logger_queue()
	var response = models.Response{Status: "failed"}
	request := &models.RequestQueue{}
	this.logger.Info("delete_queue", zap.String("subject", msg.Subject), zap.String("raw_request", string(msg.Data)))
	if err := json.Unmarshal(msg.Data, &request); err == nil {
		var isValid bool
		passwd := fmt.Sprintf("%x", sha256.Sum256([]byte(request.Password)))
		row, err := this.dbms.Database.Raw("select * from func_delete_queue (?, ?, ?)", request.Name, request.Username, passwd).Rows()

		defer row.Close()
		if err != nil {
			response.Error = err.Error()
		} else {
			for row.Next() {
				row.Scan(&isValid)
				if isValid {
					response.Status = "success"
				} else {
					response.Status = "failed"
				}
			}
		}
	} else {
		response.Error = "unrecognize request"
	}

	this.logger.Info("delete_queue", zap.String("subject", msg.Subject), zap.String("raw_request", string(msg.Data)), zap.Any("response", response))
	bytes, _ := json.Marshal(response)
	this.nats.Publish(msg.Reply, bytes)
}
