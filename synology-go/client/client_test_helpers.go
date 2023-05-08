package client

import (
	"os"
	"testing"
)

func NewTestClient(t *testing.T) Client {
	c, err := New(os.Getenv(SYNOLOGY_HOST_ENV), true)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Login(os.Getenv(SYNOLOGY_USER_ENV), os.Getenv(SYNOLOGY_PASSWORD_ENV), "webui"); err != nil {
		t.Error(err)
		return nil
	}

	return c
}
