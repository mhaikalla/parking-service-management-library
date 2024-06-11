package UsecaseParking

import (
	"parking-service/components/models/request"
	"parking-service/components/models/response"

	"parking-service/pkg/contexts"
	"parking-service/pkg/errs"
	"parking-service/pkg/file"
)

type IUsecaseParking interface {
	SetParkingIn(dc contexts.BearerContext, req *request.ParkingInRequest) (*response.BaseMessageResponse, *errs.Errs)
	SetParkingOut(dc contexts.BearerContext, req *request.ParkingOutRequest) (*response.BaseMessageResponse, *errs.Errs)
}

type usecaseObj struct {
	FileSystem file.IFileSystem
}
