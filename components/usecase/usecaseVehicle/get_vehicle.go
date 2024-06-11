package usecaseVehicle

import (
	"encoding/json"

	models "parking-service/components/models/entity"
	"parking-service/components/models/request"
	"parking-service/components/models/response"
	"parking-service/pkg/contexts"
	"parking-service/pkg/errs"
)

func (ctx *usecaseObj) GetVehicles(dc contexts.BearerContext, req *request.GetVehicleRequest) (*response.GetVehiclesResponse, *errs.Errs) {
	resp := response.GetVehiclesResponse{}

	vehicleData := []models.Vehicle{}
	resultData := []response.GetDetailVehicleResponse{}

	tableName := models.VehicleTableName
	if !ctx.FileSystem.IsFileExisting(tableName) {
		_, errCreate := ctx.FileSystem.CreateFile(tableName)
		if errCreate != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errCreate.Error())
		}
		return nil, errs.NewErrContext().
			SetCode(errs.NotFound).
			SetMessage("Data Not Found")
	} else {
		vehicle, errloadData := ctx.FileSystem.LoadFile(tableName)
		if errloadData != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errloadData.Error())
		}
		if err := json.Unmarshal(vehicle, &vehicleData); err != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(err.Error())
		}
	}
	for _, pld := range vehicleData {
		if pld.DeletedAt != nil {
			continue
		}

		resultData = append(resultData, response.GetDetailVehicleResponse{
			BaseResponse: response.BaseResponse{
				Id:        pld.Id,
				CreatedAt: pld.CreatedAt,
				UpdatedAt: pld.UpdatedAt,
			},
			Name:                pld.Name,
			Type:                pld.Type,
			FirstHourPrice:      pld.FirstHourPrice,
			PricePerHourPercent: pld.PricePerHourPercent,
		})
	}

	resp.Data = resultData
	return &resp, nil

}
