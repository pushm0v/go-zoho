# go-zoho

go-zoho is a Golang library for dealing with [ZOHO CRM](zoho.com).

## Installation

```bash
go get github.com/pushm0v/go-zoho
```

## Usage

```golang
import "github.com/pushm0v/go-zoho"

params := zoho.ZohoParams{
		GrantToken:   "123456",
		ClientID:     "123456",
		ClientSecret: "123456",
		IAMURL:       "https://accounts.zoho.com",
	}

var zohoClient = zoho.NewZoho(params)
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)