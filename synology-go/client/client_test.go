package client

import (
	"os"
	"testing"

	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/client/api"
	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/client/api/filestation"
)

func newClient() (*client, error) {
	c, err := New("dev-synology:5001", true)
	if err != nil {
		return nil, err
	}

	if err := c.Login("api-client", os.Getenv("SYNOLOGY_PASSWORD"), "webui"); err != nil {
		return nil, err
	}

	return c, nil
}

func TestFilestationInfo(t *testing.T) {
	c, err := newClient()
	if err != nil {
		t.Fatal(err)
	}

	req := filestation.NewFileStationInfoRequest(2)
	resp := filestation.FileStationInfoResponse{}
	err = c.Do(req, &resp)
	if err != nil {
		t.Log(err)
	}
	t.Logf("%+v\n", resp)

	t.Fail()
}

func TestFilestationRename(t *testing.T) {
	c, err := newClient()
	if err != nil {
		t.Fatal(err)
	}

	req := filestation.NewFileStationRenameRequest(2).
		WithPath("/some_folder").
		WithName("/renamed_folder")

	resp := filestation.FileStationRenameResponse{}
	err = c.Do(req, &resp)
	if err != nil {
		t.Log(err)
	}
	t.Logf("%+v\n", resp)

	t.Fail()
}

func TestHandleErrors(t *testing.T) {
	c := client{}

	req := func() []api.ErrorSummary {
		return []api.ErrorSummary{map[int]string{
			100: "error 100",
			101: "error 101",
			102: "error 102",
		}}
	}
	resp := api.GenericResponse{
		Success: false,
		Data:    nil,
		Error: api.SynologyError{
			Code: 100,
			Errors: []api.ErrorItem{
				{Code: 101},
				{Code: 102, Details: api.ErrorFields{"path": "/some/path", "code": 100, "reason": "a reason"}},
				{Code: 103},
			},
		},
	}
	synoErr := c.handleErrors(errorDescriber(req), resp)

	t.Errorf("%+#v\n", synoErr)

	t.Fail()
}

func TestCreateFolder(t *testing.T) {
	c, err := newClient()
	if err != nil {
		t.Fatal(err)
	}
	r := filestation.NewCreateFolderRequest(2).
		WithFolderPath("/test-folder").
		WithName("folder_from_tests")

	resp := filestation.CreateFolderResponse{}
	err = c.Do(r, &resp)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", resp)

	t.Fail()
}

type errorDescriber func() []api.ErrorSummary

func (d errorDescriber) ErrorSummaries() []api.ErrorSummary {
	return d()
}
