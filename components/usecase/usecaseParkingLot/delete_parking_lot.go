package usecaseParkingLot

import (
	"encoding/json"
	"strconv"
	"time"

	models "parking-service/components/models/entity"
	"parking-service/components/models/request"
	"parking-service/components/models/response"
	"parking-service/pkg/contexts"
	"parking-service/pkg/errs"
)

func (ctx *usecaseObj) DeleteParkingLots(dc contexts.BearerContext, req *request.DeleteParkingLotRequest) (*response.BaseMessageResponse, *errs.Errs) {
	resp := response.BaseMessageResponse{
		Message: "failed",
	}
	tableName := models.ParkingLotTableName
	dateNow := time.Now().UTC()
	parkingLotData := []models.ParkingLot{}
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
		id, errConv := strconv.Atoi(req.ParkingLotId)
		if errConv != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errConv.Error())
		}
		for i, pl := range parkingLotData {
			if pl.Id == id {
				if pl.DeletedAt != nil {
					return nil, errs.NewErrContext().
						SetCode(errs.NotFound).
						SetMessage("Data Not Found")
				}
				if pl.IsParked {
					return nil, errs.NewErrContext().
						SetCode(errs.NotFound).
						SetMessage("This Parking Area was Filled")
				}
				parkingLotData[i].DeletedAt = &dateNow
				break
			}
		}
	}
	stat, err := ctx.FileSystem.SaveData(tableName, parkingLotData)
	if !stat && err != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}
	resp.Message = "Success"
	resp.Data = nil
	return &resp, nil

}
