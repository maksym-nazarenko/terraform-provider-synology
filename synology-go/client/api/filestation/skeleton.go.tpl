package filestation

import (
	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/client/api"
)

type __TEMPLATE_TYPE_PLACEHOLDER__Request struct {
	baseFileStationRequest

	version int
}

type __TEMPLATE_TYPE_PLACEHOLDER__Response struct {
}

var _ api.Request = (*__TEMPLATE_TYPE_PLACEHOLDER__Request)(nil)

func New__TEMPLATE_TYPE_PLACEHOLDER__Request(version int) *__TEMPLATE_TYPE_PLACEHOLDER__Request {
	return &__TEMPLATE_TYPE_PLACEHOLDER__Request{
		version: version,
	}
}

func (r __TEMPLATE_TYPE_PLACEHOLDER__Request) APIName() string {
	return "SYNO.FileStation.__"
}

func (r __TEMPLATE_TYPE_PLACEHOLDER__Request) APIMethod() string {
	return "__"
}

func (r __TEMPLATE_TYPE_PLACEHOLDER__Request) APIVersion() int {
	return r.version
}

func (r __TEMPLATE_TYPE_PLACEHOLDER__Request) RequestParams() api.RequestParams {
	return map[string]string{}
}

func (r __TEMPLATE_TYPE_PLACEHOLDER__Request) NewResponseInstance() api.Response {
	return &__TEMPLATE_TYPE_PLACEHOLDER__Response{}
}

func (r __TEMPLATE_TYPE_PLACEHOLDER__Request) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{
		{
		},
		commonErrors,
	}
}
