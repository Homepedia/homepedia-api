package main

import (
	"fmt"
	auth "homepedia-api/auth"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Print("Hello, World!")
	echoInstance := echo.New()
	cronInstance, err := gocron.NewScheduler()
	if err != nil {
		echoInstance.Logger.Fatal(err)
	}
	test, err := cronInstance.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func(a string, b int) {
				// do things
			},
			"hello",
			1,
		),
	)
	if err != nil {
		echoInstance.Logger.Fatal(err)
	}
	fmt.Println((test.ID()))
	cronInstance.Start()

	err = cronInstance.Shutdown()

	if err != nil {
		echoInstance.Logger.Fatal(err)
	}

	// Init service
	auth.InitService(echoInstance)
	// Start server
	echoInstance.Logger.Fatal(echoInstance.Start(":1323"))
}
