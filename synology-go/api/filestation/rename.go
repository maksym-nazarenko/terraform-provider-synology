package filestation

import (
	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/api"
)

type RenameRequest struct {
	baseFileStationRequest

	names []string `synology:"name"`
	paths []string `synology:"path"`
}

type File struct {
	Path  string
	Name  string
	IsDir bool
}

type RenameResponse struct {
	baseFileStationResponse

	Files []File
}

var _ api.Request = (*RenameRequest)(nil)

func NewRenameRequest(version int) *RenameRequest {
	return &RenameRequest{
		baseFileStationRequest: baseFileStationRequest{
			Version:   version,
			APIName:   "SYNO.FileStation.Rename",
			APIMethod: "rename",
		},
	}
}

func (r *RenameRequest) WithName(value string) *RenameRequest {
	r.names = append(r.names, value)
	return r
}

func (r *RenameRequest) WithPath(value string) *RenameRequest {
	r.paths = append(r.paths, value)
	return r
}

func (r RenameResponse) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{
		{
			1200: "Failed to rename it.",
		},
		commonErrors,
	}
}
