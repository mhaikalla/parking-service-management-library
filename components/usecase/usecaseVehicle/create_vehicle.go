package usecaseVehicle

import (
	"encoding/json"
	"time"

	models "parking-service/components/models/entity"
	"parking-service/components/models/request"
	"parking-service/components/models/response"
	"parking-service/pkg/contexts"
	"parking-service/pkg/errs"
)

func (ctx *usecaseObj) CreateVehicle(dc contexts.BearerContext, req request.CreateVehicleRequest) (*response.BaseMessageResponse, *errs.Errs) {
	resp := response.BaseMessageResponse{
		Message: "failed",
	}

	vehicleData := []models.Vehicle{}
	tableName := models.VehicleTableName
	if !ctx.FileSystem.IsFileExisting(tableName) {
		_, errCreate := ctx.FileSystem.CreateFile(tableName)
		if errCreate != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errCreate.Error())
		}
	} else {
		parkingLot, errloadData := ctx.FileSystem.LoadFile(tableName)
		if errloadData != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errloadData.Error())
		}
		if err := json.Unmarshal(parkingLot, &vehicleData); err != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(err.Error())
		}
	}

	dateNow := time.Now().UTC()

	vehicleData = append(vehicleData, models.Vehicle{
		BaseEntity: models.BaseEntity{
			Id:        len(vehicleData) + 1,
			CreatedAt: dateNow,
			UpdatedAt: dateNow,
			DeletedAt: nil,
		},
		Name:                req.Name,
		Type:                req.Type,
		FirstHourPrice:      req.FirstHourPrice,
		PricePerHourPercent: req.PricePerHourPercent,
	})
	stat, err := ctx.FileSystem.SaveData(tableName, vehicleData)
	if !stat && err != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}
	resp.Message = "Success"
	resp.Data = req

	return &resp, nil
}
