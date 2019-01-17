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

	getTags(conn)
	getVMs(conn)

	defer conn.Close()
}

func getTags(conn *ovirtsdk4.Connection) {
	// Get the reference to the "tag" service:
	tagService := conn.SystemService().TagsService()

	resp, err := tagService.List().Send()
	if err != nil {
		fmt.Printf("Failed to get tag list, reason: %v\n", err)
		return
	}
	if tagSlice, ok := resp.Tags(); ok {
		for _, tag := range tagSlice.Slice() {
			fmt.Printf("Tag: (")
			if name, ok := tag.Name(); ok {
				fmt.Printf(" name: %v", name)
			}
			fmt.Println(")")
		}
	}
}

// Builds list of VM
func getVMs(conn *ovirtsdk4.Connection) {
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

			// Marshaling recieved responce
			jsonData, err := json.Marshal(parsedVM)
			if err != nil {
				log.Fatal("JSON marshaling failed: %s", err)
			}
			fmt.Println(string(jsonData))
		}
	}
}
