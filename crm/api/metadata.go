package api

const (
	ZOHO_CRN_API_METADATA_FIELDS_URL = "/settings/fields"
)

type ApiMetadata interface {
}

type apiMetadata struct {
	option Option
}

func NewApiMetadata(option Option) ApiMetadata {
	return &apiMetadata{
		option: option,
	}
}
