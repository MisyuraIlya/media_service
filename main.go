package main

import (
	"fmt"
	"mediaService/api"
	"mediaService/system"
	"github.com/joho/godotenv"
	"time"
)

func main() {
	for {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(".env file not found")
		}

		data := api.GetPaths()

		for _, entry := range data {
			isExist := system.CheckExistFile(entry.Path)
			fmt.Println("is Exist",isExist);
			if isExist {
				api.UploadImage(entry)
			}
		}

		fmt.Println("Waiting for 24 hours...")
		time.Sleep(24 * time.Hour) 
	}
}
