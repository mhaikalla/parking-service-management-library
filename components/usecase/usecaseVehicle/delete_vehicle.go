package usecaseVehicle

import (
	"encoding/json"
	"strconv"
	"time"

	models "github.com/mhaikalla/parking-service-management-library/components/models/entity"
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (ctx *usecaseObj) DeleteVehicles(dc contexts.BearerContext, req *request.DeleteVehicleRequest) (*response.BaseMessageResponse, *errs.Errs) {
	resp := response.BaseMessageResponse{
		Message: "failed",
	}
	tableName := models.VehicleTableName
	dateNow := time.Now().UTC()
	vehicleData := []models.Vehicle{}

	if !ctx.FileSystem.IsFileExisting(tableName) {
		_, errCreate := ctx.FileSystem.CreateFile(tableName)
		if errCreate != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errCreate.Error())
		}
		return nil, errs.NewErrContext().
			SetCode(errs.NotFound).
			SetMessage("Parking Area Not Found")
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
		id, errConv := strconv.Atoi(req.VehicleId)
		if errConv != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errConv.Error())
		}
		for i, pl := range vehicleData {
			if pl.Id == id {
				if pl.DeletedAt != nil {
					return nil, errs.NewErrContext().
						SetCode(errs.NotFound).
						SetMessage("Data Not Found")
				}
				vehicleData[i].DeletedAt = &dateNow
				break
			}
		}
	}
	stat, err := ctx.FileSystem.SaveData(tableName, vehicleData)
	if !stat && err != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}
	resp.Message = "Success"
	resp.Data = nil
	return &resp, nil
}
