package filestation

import (
	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/client/api"
)

type FileStationRenameRequest struct {
	baseFileStationRequest

	version int
	path    string
	name    string
}

type File struct {
	Path  string
	Name  string
	IsDir bool
}

type FileStationRenameResponse struct {
	Files []File
}

var _ api.Request = (*FileStationRenameRequest)(nil)

func NewFileStationRenameRequest(version int) *FileStationRenameRequest {
	return &FileStationRenameRequest{
		version: version,
	}
}

func (r *FileStationRenameRequest) WithPath(path string) *FileStationRenameRequest {
	r.path = path
	return r
}

func (r *FileStationRenameRequest) WithName(name string) *FileStationRenameRequest {
	r.name = name
	return r
}

func (r FileStationRenameRequest) APIName() string {
	return "SYNO.FileStation.Rename"
}

func (r FileStationRenameRequest) APIMethod() string {
	return "rename"
}

func (r FileStationRenameRequest) APIVersion() int {
	return r.version
}

func (r FileStationRenameRequest) RequestParams() api.RequestParams {
	return map[string]string{
		"path": r.path,
		"name": r.name,
	}
}

func (r FileStationRenameRequest) NewResponseInstance() api.Response {
	return &FileStationRenameResponse{}
}

func (r FileStationRenameRequest) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{
		{
			1200: "Failed to rename it.",
		},
		commonErrors,
	}
}
