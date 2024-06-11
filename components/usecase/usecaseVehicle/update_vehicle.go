package usecaseVehicle

import (
	"encoding/json"
	"time"

	models "github.com/mhaikalla/parking-service-management-library/components/models/entity"
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (ctx *usecaseObj) UpdateVehicle(dc contexts.BearerContext, req request.UpdateVehicleRequest) (*response.BaseMessageResponse, *errs.Errs) {
	resp := response.BaseMessageResponse{
		Message: "failed",
	}
	dateNow := time.Now().UTC()
	vehicleData := []models.Vehicle{}
	tableName := models.VehicleTableName
	if !ctx.FileSystem.IsFileExisting(tableName) {
		_, errCreate := ctx.FileSystem.CreateFile(tableName)
		if errCreate != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errCreate.Error())
		}
		vehicleData = append(vehicleData, models.Vehicle{
			BaseEntity: models.BaseEntity{
				Id:        len(vehicleData) + 1,
				CreatedAt: dateNow,
				UpdatedAt: dateNow,
				DeletedAt: nil,
			},
			Type:                req.Type,
			Name:                req.Name,
			FirstHourPrice:      req.FirstHourPrice,
			PricePerHourPercent: req.PricePerHourPercent,
		})
	} else {
		data, errloadData := ctx.FileSystem.LoadFile(tableName)
		if errloadData != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errloadData.Error())
		}
		if err := json.Unmarshal(data, &vehicleData); err != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(err.Error())
		}
	}

	for i, pl := range vehicleData {
		if pl.Id == req.Id {
			if pl.DeletedAt != nil {
				return nil, errs.NewErrContext().
					SetCode(errs.NotFound).
					SetMessage("Data Not Found")
			}
			vehicleData[i].UpdatedAt = dateNow
			vehicleData[i].FirstHourPrice = req.FirstHourPrice
			vehicleData[i].PricePerHourPercent = req.PricePerHourPercent
			vehicleData[i].Type = req.Type
			vehicleData[i].Name = req.Name
			break
		}
	}
	stat, err := ctx.FileSystem.SaveData(tableName, vehicleData)
	if !stat && err != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}
	resp.Message = "Success, Data Updated"
	resp.Data = req
	return &resp, nil
}
