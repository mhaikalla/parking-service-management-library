package UsecaseParking

import (
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"

	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
	"github.com/mhaikalla/parking-service-management-library/pkg/file"
)

type IUsecaseParking interface {
	SetParkingIn(dc contexts.BearerContext, req *request.ParkingInRequest) (*response.BaseMessageResponse, *errs.Errs)
	SetParkingOut(dc contexts.BearerContext, req *request.ParkingOutRequest) (*response.BaseMessageResponse, *errs.Errs)
	GetParkingData(dc contexts.BearerContext, req *request.GetParkingData) (*response.GetDataParkingResponse, *errs.Errs)
	GetCountParkingData(dc contexts.BearerContext, req *request.GetCountParkingData) (*response.GetCountParkingResponse, *errs.Errs)
}

type usecaseObj struct {
	FileSystem file.IFileSystem
}
