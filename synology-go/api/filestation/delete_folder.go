package filestation

import (
	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/api"
)

type DeleteFolderRequest struct {
	baseFileStationRequest

	folderPaths []string `synology:"path"`
	recursive   bool     `synology:"recursive"`
}

type DeleteFolderResponse struct {
	baseFileStationResponse
}

var _ api.Request = (*DeleteFolderRequest)(nil)

func NewDeleteFolderRequest(version int) *DeleteFolderRequest {
	return &DeleteFolderRequest{
		baseFileStationRequest: baseFileStationRequest{
			Version:   version,
			APIName:   "SYNO.FileStation.Delete",
			APIMethod: "delete",
		},
	}
}

func (r *DeleteFolderRequest) WithPath(p string) *DeleteFolderRequest {
	r.folderPaths = append(r.folderPaths, p)
	return r
}

func (r *DeleteFolderRequest) WithRecursive(v bool) *DeleteFolderRequest {
	r.recursive = v
	return r
}

func (r DeleteFolderResponse) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{
		{
			900: "Failed to delete file(s)/folder(s). More information in <errors> object.",
		},
		commonErrors,
	}
}
