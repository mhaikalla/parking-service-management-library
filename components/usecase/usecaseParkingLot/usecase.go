package usecaseParkingLot

import (
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"
	"github.com/mhaikalla/parking-service-management-library/pkg/file"

	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

type IUsecaseParkingLot interface {
	CreateParkingLot(dc contexts.BearerContext, req request.CreateParkingLotRequest) (*response.BaseMessageResponse, *errs.Errs)
	UpdateParkingLot(dc contexts.BearerContext, req request.UpdateParkingLotRequest) (*response.BaseMessageResponse, *errs.Errs)
	DeleteParkingLots(dc contexts.BearerContext, req *request.DeleteParkingLotRequest) (*response.BaseMessageResponse, *errs.Errs)
	GetParkingLots(dc contexts.BearerContext, req *request.GetParkingLotRequest) (*response.GetParkingLotsResponse, *errs.Errs)
	GetDetailParkingLot(dc contexts.BearerContext, req *request.GetDetailParkingLotRequest) (*response.GetDetailParkingLotResponse, *errs.Errs)
}

type usecaseObj struct {
	FileSystem file.IFileSystem
}

func NewParkingLotUsecase(ctx ...interface{}) IUsecaseParkingLot {
	handle := usecaseObj{}
	for _, c := range ctx {
		switch c.(type) {
		case file.IFileSystem:
			handle.FileSystem = c.(file.IFileSystem)
		}
	}
	return &handle
}
