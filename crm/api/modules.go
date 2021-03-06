package api

import (
	"encoding/json"

	"github.com/pushm0v/go-zoho/models"
)

const (
	ZOHO_CRM_API_MODULES_URL = "/settings/modules"
)

type ApiModules interface {
	List() (modules []models.Modules, err error)
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

func (m *apiModules) List() (modules []models.Modules, err error) {
	var params = map[string]interface{}{}

	resp, err := m.option.HttpClient.Get(m.option.ApiUrl(ZOHO_CRM_API_MODULES_URL, false), params)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var respModules = new(responseModules)

	err = json.NewDecoder(resp.Body).Decode(&respModules)
	if err != nil {
		return
	}

	return respModules.Root, nil
}
