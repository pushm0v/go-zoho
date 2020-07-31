package api

import (
	"encoding/json"

	"github.com/pushm0v/go-zoho/models"
)

const (
	ZOHO_CRN_API_METADATA_FIELDS_URL = "/settings/fields"
)

type ApiMetadata interface {
	ListFields(string) (fields []models.Fields, err error)
}

type apiMetadata struct {
	option Option
}

type responseFields struct {
	Root []models.Fields `json:"fields"`
}

func NewApiMetadata(option Option) ApiMetadata {
	return &apiMetadata{
		option: option,
	}
}

func (m *apiMetadata) ListFields(moduleName string) (fields []models.Fields, err error) {
	var params = map[string]interface{}{
		"module": moduleName,
	}

	resp, err := m.option.HttpClient.Get(m.option.ApiUrl(ZOHO_CRN_API_METADATA_FIELDS_URL), params)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var responseFields = new(responseFields)

	err = json.NewDecoder(resp.Body).Decode(&responseFields)
	if err != nil {
		return
	}

	return responseFields.Root, nil
}
