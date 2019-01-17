package main

import (
	"encoding/json"
	"fmt"
	ovirtsdk4 "gopkg.in/imjoey/go-ovirt.v4"
	"gopkg.in/ini.v1"
	"log"
	"time"
)

const (
	ConfigFilePath = "ovirt.ini"
)

type OvirtVM struct {
	Name   string `json:"name"`
	FQDN   string `json:"fqdn"`
	Status string `json:"status"`
}

func main() {

	config, err := ini.Load(ConfigFilePath)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err.Error())
	}

	url := config.Section("ovirt").Key("ovirt_url").String()
	user := config.Section("ovirt").Key("ovirt_username").String()
	pass := config.Section("ovirt").Key("ovirt_password").String()

	// Create the connection to the api server
	conn, err := ovirtsdk4.NewConnectionBuilder().
		URL(url).
		Username(user).
		Password(pass).
		Insecure(true).
		Compress(true).
		Timeout(time.Second * 20).
		Build()
	if err != nil {
		log.Fatalf("Make connection failed, reason: %s", err.Error())
	}

	defer conn.Close()

	// Get the reference to the "vms" service:
	vmsService := conn.SystemService().VmsService()

	// Use the "list" method of the "vms" service to list all the virtual machines
	vmsResponse, err := vmsService.List().Send()

	if err != nil {
		fmt.Printf("Failed to get vm list, reason: %v\n", err)
		return
	}
	if vms, ok := vmsResponse.Vms(); ok {
		// Print the virtual machine names and identifiers:
		for _, vm := range vms.Slice() {

			var parsedVM OvirtVM
			var jsonData []byte

			if vmName, ok := vm.Name(); ok {
				parsedVM.Name = vmName
			}

			if vmFQDN, ok := vm.Fqdn(); ok {
				parsedVM.FQDN = vmFQDN
			}

			if vmStatus, ok := vm.Status(); ok {
				parsedVM.Status = string(vmStatus)
			}

			//marshaling recieved responce
			jsonData, err := json.Marshal(parsedVM)
			if err != nil {
				log.Fatal("JSON marshaling failed: %s", err)
			}
			fmt.Println(string(jsonData))
		}
	}
}
