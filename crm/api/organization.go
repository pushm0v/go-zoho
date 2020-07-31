package api

import (
	"encoding/json"

	"github.com/pushm0v/go-zoho/models"
)

const (
	ZOHO_CRM_API_ORGANIZATION_URL = "/org"
)

type ApiOrganization interface {
	Details() (orgs []models.Organization, err error)
}

type apiOrganization struct {
	option Option
}

type responseOrganization struct {
	Root []models.Organization `json:"org"`
}

func NewApiOrganization(option Option) ApiOrganization {
	return &apiOrganization{
		option: option,
	}
}

func (m *apiOrganization) Details() (orgs []models.Organization, err error) {
	var params = map[string]interface{}{}

	resp, err := m.option.HttpClient.Get(m.option.ApiUrl(ZOHO_CRM_API_ORGANIZATION_URL), params)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var respOrgs = new(responseOrganization)

	err = json.NewDecoder(resp.Body).Decode(&respOrgs)
	if err != nil {
		return
	}

	return respOrgs.Root, nil
}
