package filestation

import (
	"testing"

	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRenameFolder(t *testing.T) {
	sharedFolder := "/test-folder"
	testFolder := "integration-test"
	renamedTestFolder := "integration-test-renamed"
	c := client.NewTestClient(t)
	r := NewCreateFolderRequest(2).WithFolderPath(sharedFolder).WithName(testFolder).WithForceParent(true)

	resp := CreateFolderResponse{}
	err := c.Do(r, &resp)
	require.NoError(t, err)
	if !resp.Success() {
		t.Fatal(resp.GetError().Error())
	}
	assert.Len(t, resp.Folders, 1)

	// rename it
	renameReq := NewRenameRequest(2).WithPath(sharedFolder + "/" + testFolder).WithName(renamedTestFolder)
	renameResp := RenameResponse{}
	err = c.Do(renameReq, &renameResp)
	require.NoError(t, err)
	if !resp.Success() {
		t.Fatal(resp.GetError().Error())
	}
	assert.Len(t, renameResp.Files, 1)

	defer func() {
		r := NewDeleteFolderRequest(2).
			WithRecursive(true).
			WithPath(sharedFolder)
		resp := DeleteFolderResponse{}
		err := c.Do(r, &resp)
		if err != nil {
			t.Error(err)
		}
		if !resp.Success() {
			t.Error(resp.GetError().Error())
		}
	}()
}
