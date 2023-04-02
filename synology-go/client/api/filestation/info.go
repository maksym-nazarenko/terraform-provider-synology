package filestation

import (
	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/client/api"
)

type FileStationInfoRequest struct {
	baseFileStationRequest

	version int
}

type FileStationInfoResponse struct {
	IsManager              bool
	SupportVirtualProtocol string
	Supportsharing         bool
	Hostname               string
}

var _ api.Request = (*FileStationInfoRequest)(nil)

func NewFileStationInfoRequest(version int) *FileStationInfoRequest {
	return &FileStationInfoRequest{
		version: version,
	}
}

func (r FileStationInfoRequest) APIName() string {
	return "SYNO.FileStation.Info"
}

func (r FileStationInfoRequest) APIMethod() string {
	return "get"
}

func (r FileStationInfoRequest) APIVersion() int {
	return r.version
}

func (r FileStationInfoRequest) RequestParams() api.RequestParams {
	return nil
}

func (r FileStationInfoRequest) NewResponseInstance() api.Response {
	return &FileStationInfoResponse{}
}

func (r FileStationInfoRequest) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{commonErrors}
}
