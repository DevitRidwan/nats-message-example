package main

import (
	"Ridwan/Queue/src/info"
	"Ridwan/Queue/src/properties"
	service2 "Ridwan/Queue/src/service"
	"Ridwan/Queue/src/tools"
	"fmt"
	"sync"
)

func main() {
	info.PrintHeader()
	prop := properties.ServiceProperties{}
	err := tools.Ekstration("configPath", "configName", &prop)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		service := service2.QueueManagement{}
		service.Init(&prop)
		var wg sync.WaitGroup
		wg.Add(1)
		wg.Wait()
	}
}
