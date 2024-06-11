package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	parkingHandler "github.com/mhaikalla/parking-service-management-library/components/handlers/parking"
	parkingLotHandler "github.com/mhaikalla/parking-service-management-library/components/handlers/parkinglot"
	vehicleHandler "github.com/mhaikalla/parking-service-management-library/components/handlers/vehicle"
	"github.com/mhaikalla/parking-service-management-library/pkg/condutils"
	"github.com/mhaikalla/parking-service-management-library/pkg/file"
	"github.com/mhaikalla/parking-service-management-library/pkg/router"
	validatorRequest "github.com/mhaikalla/parking-service-management-library/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// "github.com/mhaikalla/parking-service-management-library/pkg/serverpatch"
	// "github.com/mhaikalla/parking-service-management-library/pkg/staticdatas"

	"github.com/mhaikalla/parking-service-management-library/pkg/config"
	"github.com/mhaikalla/parking-service-management-library/pkg/logs"
)

var (
	configFile = "configs/dev-default.yaml"
)

func main() {

	logger := logs.NewLogrus(condutils.Or(os.Getenv("SYSTEM_NAME"), "SYSTEM").(string))
	logger.Update()

	defer func() {
		if r := recover(); r != nil {
			errPanic := fmt.Errorf("got an error on main thread: %v", r)
			logger.Upsert("stacktrace", string(debug.Stack()))
			logger.Update()
			logger.Panic(errPanic)
		}
	}()

	viper := config.NewViperLocalProvider()

	fmt.Println("configFile: ", configFile)

	if errViper := viper.GetConfig(configFile); errViper != nil {
		logger.Fatal(errViper)
	}

	config := viper.Config()
	server := router.NewEchoServerV2(config)
	fileStorage := file.NewStorageFile(config)

	validators := validatorRequest.NewValidator()

	parkingHandler, parkingErr := parkingHandler.NewParkingHandlers(config, validators, fileStorage)
	parkingLotHandler, parkingLotErr := parkingLotHandler.NewParkingLotHandlers(config, validators, fileStorage)
	vehicleHandler, VehicleErr := vehicleHandler.NewVehicleHandlers(config, validators, fileStorage)

	if e, ok := condutils.Ors(
		parkingErr,
		parkingLotErr,
		VehicleErr,
	).(error); ok && e != nil {
		logger.Fatal(e)
	}

	server.Handle("POST", "/api/v1/parking-management/parking-in", parkingHandler.SetParkingIn())
	server.Handle("POST", "/api/v1/parking-management/parking-out", parkingHandler.SetParkingOut())

	server.Handle("GET", "/api/v1/parking-management/parking-lot/:id", parkingLotHandler.GetDetailParkingLot())
	server.Handle("GET", "/api/v1/parking-management/parking-lots", parkingLotHandler.GetParkingLot())
	server.Handle("POST", "/api/v1/parking-management/parking-lot", parkingLotHandler.CreateParkingLot())
	server.Handle("PUT", "/api/v1/parking-management/parking-lot", parkingLotHandler.UpdateParkingLot())
	server.Handle("DELETE", "/api/v1/parking-management/parking-lot", parkingLotHandler.DeleteParkingLot())

	server.Handle("GET", "/api/v1/parking-management/vehicle", vehicleHandler.GetDetailVehicle())
	server.Handle("GET", "/api/v1/parking-management/vehicles", vehicleHandler.GetVehicle())
	server.Handle("POST", "/api/v1/parking-management/vehicle", vehicleHandler.CreateVehicle())
	server.Handle("PUT", "/api/v1/parking-management/vehicle", vehicleHandler.UpdateVehicle())
	server.Handle("DELETE", "/api/v1/parking-management/vehicles", vehicleHandler.DeleteVehicle())

	ecServer := server.GetServer()
	ecServer.Use(middleware.CORS())

	go func(addr string, server *echo.Echo) {
		if startErr := server.Start(addr); startErr != nil {
			logger.Fatal(startErr)
		}
	}(condutils.Or(":8080", viper.Config()["server"]).(string), ecServer)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	defer logger.Info("server get interrupt signal")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ecServer.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}

	logger.Info("wait more 5 second for async write to db")
	time.Sleep(time.Second * 5)
	logger.Info("Exiting")
}
