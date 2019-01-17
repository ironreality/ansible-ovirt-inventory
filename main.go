package main

import (
	"fmt"
	ovirtsdk4 "gopkg.in/imjoey/go-ovirt.v4"
	"gopkg.in/ini.v1"
	"log"
	"time"
)

const (
	ConfigFilePath = "ovirt.ini"
)

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

			if vmName, ok := vm.Name(); ok {
				fmt.Printf("VM name: %v\n", vmName)
			}

			if vmFQDN, ok := vm.Fqdn(); ok {
				fmt.Printf("FQDN: %v\n", vmFQDN)
			}

			if vmType, ok := vm.Type(); ok {
				fmt.Printf("type: %v\n", vmType)
			}

			if vmStatus, ok := vm.Status(); ok {
				fmt.Printf("status: %v\n", vmStatus)
			}

			if vmOS, ok := vm.Os(); ok {
				t, _ := vmOS.Type()
				fmt.Printf("OS: %v\n", t)
			}

			if vmHost, ok := vm.Host(); ok {
				fmt.Printf("Host: %v\n", vmHost.Name)
			}

			fmt.Println("")
		}
	}
}
