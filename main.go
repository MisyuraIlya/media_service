package main

import (
	"fmt"
	"log"
	"mediaService/api"
	"mediaService/system"
	"time"

	"github.com/kardianos/service"
)

// program implements the service.Interface
type program struct{}

// Start is called when the service is started.
func (p *program) Start(s service.Service) error {
	// Do not block. Start the actual work in a new goroutine.
	go p.run()
	return nil
}

// run contains the main logic of your application.
func (p *program) run() {
	for {
		data := api.GetPaths()

		for _, entry := range data {
			isExist := system.CheckExistFile(entry.Path)
			fmt.Println("is Exist", isExist)
			if isExist {
				api.UploadImage(entry)
			}
		}

		fmt.Println("Waiting for 24 hours...")
		time.Sleep(24 * time.Hour)
	}
}

// Stop is called when the service is stopped.
func (p *program) Stop(s service.Service) error {
	// Perform any cleanup tasks if necessary.
	return nil
}

func main() {
	// Define service configuration.
	svcConfig := &service.Config{
		Name:        "MediaService", // Service name (no spaces)
		DisplayName: "Media Service", // Name displayed in service manager
		Description: "Service that checks file paths and uploads images every 24 hours.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Optionally, set up a logger for the service.
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Run the service. This call will block.
	if err := s.Run(); err != nil {
		logger.Error(err)
	}
}
