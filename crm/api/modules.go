package api

import (
	"encoding/json"

	"github.com/pushm0v/go-zoho/models"
)

const (
	ZOHO_CRM_API_MODULES_URL = "/settings/modules"
)

type ApiModules interface {
	List() (err error, modules []models.Modules)
}

type apiModules struct {
	option Option
}

type responseModules struct {
	Root []models.Modules `json:"modules"`
}

func NewApiModules(option Option) ApiModules {
	return &apiModules{
		option: option,
	}
}

func (m *apiModules) List() (err error, modules []models.Modules) {
	var params = map[string]interface{}{}

	resp, err := m.option.HttpClient.Get(m.option.ApiUrl(ZOHO_CRM_API_MODULES_URL), params)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var respModules = new(responseModules)

	err = json.NewDecoder(resp.Body).Decode(&respModules)
	if err != nil {
		return
	}

	return nil, respModules.Root
}
