package UsecaseParking

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/mhaikalla/parking-service-management-library/components/constant"
	models "github.com/mhaikalla/parking-service-management-library/components/models/entity"
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
	"github.com/mhaikalla/parking-service-management-library/pkg/file"
)

func NewParkingUsecase(ctx ...interface{}) IUsecaseParking {
	handle := usecaseObj{}
	for _, c := range ctx {
		switch c.(type) {
		case file.IFileSystem:
			handle.FileSystem = c.(file.IFileSystem)
		}
	}
	return &handle
}

func (ctx *usecaseObj) SetParkingIn(dc contexts.BearerContext, req *request.ParkingInRequest) (*response.BaseMessageResponse, *errs.Errs) {
	resp := response.BaseMessageResponse{
		Message: "failed",
		Data:    nil,
	}
	parkingStatusData := []models.ParkingVehicleStatus{}
	parkingLotData := []models.ParkingLot{}
	currentParkingLotData := models.ParkingLot{}

	fileNameParkingStatus := models.ParkingVehicleStatusTableName
	if !ctx.FileSystem.IsFileExisting(fileNameParkingStatus) {
		_, errCreate := ctx.FileSystem.CreateFile(fileNameParkingStatus)
		if errCreate != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errCreate.Error())
		}
	} else {
		parking, errloadData := ctx.FileSystem.LoadFile(fileNameParkingStatus)
		if errloadData != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errloadData.Error())
		}
		if err := json.Unmarshal(parking, &parkingStatusData); err != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(err.Error())
		}
	}

	fileNameParkingLot := models.ParkingLotTableName
	if !ctx.FileSystem.IsFileExisting(fileNameParkingLot) {
		_, errCreate := ctx.FileSystem.CreateFile(fileNameParkingLot)
		if errCreate != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errCreate.Error())
		}
	} else {
		parkinglot, errloadData := ctx.FileSystem.LoadFile(fileNameParkingLot)
		if errloadData != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errloadData.Error())
		}
		if err := json.Unmarshal(parkinglot, &parkingLotData); err != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(err.Error())
		}
	}

	dateNow := time.Now().UTC()

	for _, p := range parkingStatusData {
		if p.PlateNumber == req.PlatNomor && p.Status == constant.ParkingIn {
			return nil, errs.NewErrContext().
				SetCode(errs.BadRequest).
				SetMessage("This vehicle has already been parked")
		}
	}

	for _, pld := range parkingLotData {
		if !pld.IsParked && pld.DeletedAt == nil {
			pld.UpdatedAt = dateNow
			pld.IsParked = true
			currentParkingLotData = pld
			break
		}
	}
	if reflect.ValueOf(currentParkingLotData).IsZero() {
		return nil, errs.NewErrContext().
			SetCode(errs.BadRequest).
			SetMessage("There's No Parking Area Available")
	}

	parkingStatusData = append(parkingStatusData, models.ParkingVehicleStatus{
		BaseEntity: models.BaseEntity{
			Id:        len(parkingStatusData) + 1,
			CreatedAt: dateNow,
			UpdatedAt: dateNow,
			DeletedAt: nil,
		},
		PlateNumber:    req.PlatNomor,
		Type:           req.Tipe,
		Color:          req.Warna,
		ParkingInDate:  dateNow.String(),
		ParkingOutDate: nil,
		Status:         constant.ParkingIn,
		Price:          0,
		ParkingLot:     currentParkingLotData.Name,
	})

	saveDataParkingStatus, err := ctx.FileSystem.SaveData(models.ParkingVehicleStatusTableName, parkingStatusData)
	if !saveDataParkingStatus && err != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}
	saveDataParkingLot, errparkinglot := ctx.FileSystem.SaveData(models.ParkingLotTableName, parkingLotData)
	if !saveDataParkingLot && errparkinglot != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}
	resp.Message = "Success"
	resp.Data = req

	return &resp, nil
}

func (ctx *usecaseObj) SetParkingOut(dc contexts.BearerContext, req *request.ParkingOutRequest) (*response.BaseMessageResponse, *errs.Errs) {
	resp := response.BaseMessageResponse{
		Message: "failed",
		Data:    nil,
	}
	parkingStatusData := []models.ParkingVehicleStatus{}
	vehicleData := []models.Vehicle{}
	parkingLotData := []models.ParkingLot{}

	fileNameParkingStatus := models.ParkingVehicleStatusTableName
	fileNameVehicle := models.VehicleTableName
	fileNameParkingLot := models.ParkingLotTableName
	var wg sync.WaitGroup
	errChan := make(chan errs.Errs, 10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if !ctx.FileSystem.IsFileExisting(fileNameParkingStatus) {
			_, errCreate := ctx.FileSystem.CreateFile(fileNameParkingStatus)
			if errCreate != nil {
				errRet := errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(errCreate.Error())
				errChan <- *errRet
				return
			}
		} else {
			parking, errloadData := ctx.FileSystem.LoadFile(fileNameParkingStatus)
			if errloadData != nil {
				errRet := errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(errloadData.Error())
				errChan <- *errRet
				return
			}
			if err := json.Unmarshal(parking, &parkingStatusData); err != nil {
				errRet := errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(err.Error())
				errChan <- *errRet
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if !ctx.FileSystem.IsFileExisting(fileNameVehicle) {
			_, errCreate := ctx.FileSystem.CreateFile(fileNameVehicle)
			if errCreate != nil {
				errRet := errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(errCreate.Error())
				errChan <- *errRet
				return
			}
		} else {
			parking, errloadData := ctx.FileSystem.LoadFile(fileNameVehicle)
			if errloadData != nil {
				errRet := errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(errloadData.Error())
				errChan <- *errRet
				return
			}
			if err := json.Unmarshal(parking, &vehicleData); err != nil {
				errRet := errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(err.Error())
				errChan <- *errRet
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if !ctx.FileSystem.IsFileExisting(fileNameParkingLot) {
			_, errCreate := ctx.FileSystem.CreateFile(fileNameParkingLot)
			if errCreate != nil {
				errRet := errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(errCreate.Error())
				errChan <- *errRet
				return
			}
		} else {
			parking, errloadData := ctx.FileSystem.LoadFile(fileNameParkingLot)
			if errloadData != nil {
				errRet := errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(errloadData.Error())
				errChan <- *errRet
				return
			}
			if err := json.Unmarshal(parking, &parkingLotData); err != nil {
				errRet := errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(err.Error())
				errChan <- *errRet
				return
			}

		}
	}()

	wg.Wait()

	close(errChan)

	for err := range errChan {
		return &resp, &err
	}

	currentData := models.ParkingVehicleStatus{}
	currentVehicleData := models.Vehicle{}
	currentParkingLotData := models.ParkingLot{}

	dateNow := time.Now().UTC()
	dateNowStr := dateNow.String()

	lengthData := len(parkingStatusData)
	if lengthData > 0 {
		for i := lengthData - 1; i >= 0; i-- {
			if parkingStatusData[i].PlateNumber == req.PlatNomor {
				if parkingStatusData[i].Status == constant.ParkingOut {
					return nil, errs.NewErrContext().
						SetCode(errs.BadRequest).
						SetMessage("This vehicle has left the parking lot")
				}
				currentData = parkingStatusData[i]
				break
			}
		}
	} else {
		return nil, errs.NewErrContext().
			SetCode(errs.BadRequest).
			SetMessage("There's No Vehicle Parking With These Plate Number")
	}

	if len(vehicleData) > 0 {
		for _, vd := range vehicleData {
			if vd.Type == currentData.Type {
				currentVehicleData = vd
			}
		}
		if reflect.ValueOf(currentVehicleData).IsZero() {
			return nil, errs.NewErrContext().
				SetCode(errs.BadRequest).
				SetMessage("Vehicle Data Not Found")
		}
	} else {
		return nil, errs.NewErrContext().
			SetCode(errs.BadRequest).
			SetMessage("Vehicle Data Not Found")
	}

	if len(parkingLotData) > 0 {
		for _, pld := range parkingLotData {
			if pld.Name == currentData.ParkingLot {
				currentParkingLotData = pld
			}
		}
		if reflect.ValueOf(currentParkingLotData).IsZero() {
			return nil, errs.NewErrContext().
				SetCode(errs.BadRequest).
				SetMessage("Parking Area Data Not Found")
		}
	} else {
		return nil, errs.NewErrContext().
			SetCode(errs.BadRequest).
			SetMessage("Parking Area Not Found")
	}

	t, err := time.Parse("2006-01-02 15:04:05", currentData.ParkingInDate)
	if err != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}

	totalPrice := currentVehicleData.FirstHourPrice + (int(dateNow.Sub(t).Hours()) * (currentVehicleData.FirstHourPrice * (currentVehicleData.PricePerHourPercent / 100)))
	parkingStatusData = append(parkingStatusData, models.ParkingVehicleStatus{
		BaseEntity: models.BaseEntity{
			Id:        len(parkingStatusData) + 1,
			CreatedAt: dateNow,
			UpdatedAt: dateNow,
			DeletedAt: nil,
		},
		PlateNumber:    currentData.PlateNumber,
		Type:           currentData.Type,
		Color:          currentData.Color,
		ParkingInDate:  currentData.ParkingInDate,
		ParkingOutDate: &dateNowStr,
		Status:         constant.ParkingOut,
		Price:          totalPrice,
	})
	stat, err := ctx.FileSystem.SaveData(models.ParkingVehicleStatusTableName, parkingStatusData)
	if !stat && err != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}

	// stat, err := ctx.FileSystem.SaveData(models., parkingData)
	// if !stat && err != nil {
	// 	return nil, errs.NewErrContext().
	// 		SetCode(errs.InternalServerError).
	//
	// 		SetMessage(err.Error())
	// }
	resp.Message = "Success"
	resp.Data = req

	return &resp, nil
}
