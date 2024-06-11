package usecaseParkingLot

import (
	"encoding/json"

	models "github.com/mhaikalla/parking-service-management-library/components/models/entity"
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (ctx *usecaseObj) GetParkingLots(dc contexts.BearerContext, req *request.GetParkingLotRequest) (*response.GetParkingLotsResponse, *errs.Errs) {
	resp := response.GetParkingLotsResponse{}

	parkingLotData := []models.ParkingLot{}
	resultData := []response.GetDetailParkingLotResponse{}

	tableName := models.ParkingLotTableName
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
		parkingLot, errloadData := ctx.FileSystem.LoadFile(tableName)
		if errloadData != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errloadData.Error())
		}
		if err := json.Unmarshal(parkingLot, &parkingLotData); err != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(err.Error())
		}
	}
	for _, pld := range parkingLotData {
		if pld.DeletedAt != nil {
			continue
		}

		resultData = append(resultData, response.GetDetailParkingLotResponse{
			BaseResponse: response.BaseResponse{
				Id:        pld.Id,
				CreatedAt: pld.CreatedAt,
				UpdatedAt: pld.UpdatedAt,
			},
			Name:     pld.Name,
			Floor:    pld.Floor,
			IsParked: pld.IsParked,
		})
	}

	resp.Data = resultData
	return &resp, nil

}
