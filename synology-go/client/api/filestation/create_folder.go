package filestation

import (
	"strconv"
	"strings"

	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/client/api"
)

type CreateFolderRequest struct {
	baseFileStationRequest

	version     int
	folderPaths []string
	names       []string
	forceParent bool
}

type CreateFolderResponse struct {
	Folders []struct {
		Path  string
		Name  string
		IsDir bool
	}
}

var _ api.Request = (*CreateFolderRequest)(nil)

func NewCreateFolderRequest(version int) *CreateFolderRequest {
	return &CreateFolderRequest{
		version: version,
	}
}

func (r CreateFolderRequest) APIName() string {
	return "SYNO.FileStation.CreateFolder"
}

func (r CreateFolderRequest) APIMethod() string {
	return "create"
}

func (r CreateFolderRequest) APIVersion() int {
	return r.version
}

func (r *CreateFolderRequest) WithFolderPath(value string) *CreateFolderRequest {
	r.folderPaths = append(r.folderPaths, value)
	return r
}

func (r *CreateFolderRequest) WithName(value string) *CreateFolderRequest {
	r.names = append(r.names, value)
	return r
}

func (r *CreateFolderRequest) WithForceParent(value bool) *CreateFolderRequest {
	r.forceParent = value
	return r
}

// todo: create generic function to handle different Go types
func (r CreateFolderRequest) RequestParams() api.RequestParams {
	return map[string]string{
		"folder_path":  "[\"" + strings.Join(r.folderPaths, "\",\"") + "\"]",
		"name":         "[\"" + strings.Join(r.names, "\",\"") + "\"]",
		"force_parent": strconv.FormatBool(r.forceParent),
	}
}

func (r CreateFolderRequest) NewResponseInstance() api.Response {
	return &CreateFolderResponse{}
}

func (r CreateFolderRequest) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{
		{
			1100: "Failed to create a folder. More information in <errors> object.",
			1101: "The number of folders to the parent folder would exceed the system limitation.",
		},
		commonErrors,
	}
}
