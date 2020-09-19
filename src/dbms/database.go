package dbms

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	POSTGRESQL = 1
)

type Properties struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	MaxIdleConns int    `json:"max_idle_conns" mapstructure:"max_idle_conns"`
}

type DatabaseRepository struct {
	Database         *gorm.DB
	properties       *Properties
	connectionType   int
	connectionNumber int
	status           bool
	mutexStart       sync.Mutex
	mutexEnd         sync.Mutex
	mutexConn        sync.Mutex
	group            sync.WaitGroup
}

func New(properties *Properties) *DatabaseRepository {
	fmt.Println(properties)
	return &DatabaseRepository{properties: properties}

}

func (this *DatabaseRepository) Connect() error {
	var err error
	// if this.connectionType == POSTGRESQL {
	this.Database, err = gorm.Open("postgres", "host="+this.properties.Host+" port="+this.properties.Port+" user="+this.properties.Username+" dbname="+this.properties.Name+" password="+this.properties.Password+" sslmode=disable ")
	// } else {
	// 	err = fmt.Errorf("dbms : unknown Database connector type %q", this.connectionType)
	// }
	if err == nil {
		this.Database.DB().SetMaxIdleConns(this.properties.MaxIdleConns)
		this.status = true
		this.connectionNumber = 0
	} else {
		this.status = false
	}
	return err
}
