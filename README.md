# go-zoho

go-zoho is a Golang library for dealing with [ZOHO CRM](https://zoho.com). Please see [ZOHO API Docs](https://www.zoho.com/crm/developer/docs/api/) for details on API credentials.

## Installation

```bash
go get github.com/pushm0v/go-zoho
```

## Basic usage
You can use custom `http.Client` & `cache` if you want to implement Redis or similar cache, by default it will use `LocalCache`. 
```golang
import "github.com/pushm0v/go-zoho"

var params = zoho.ZohoParams{
    GrantToken:    "1000.xxxx.yyyy",
    ClientID:      "1000.abcdef",
    ClientSecret:  "client-secret",
    IamURL:        "https://accounts.zoho.com",
    CrmURL:        "https://zohoapis.com/crm",
    FileUploadURL: "https://content.zohoapis.com/crm",
    ZGID:          "12345677",
    Version:       crm.VersionV2,
}

var httpClient = http2.NewHttpClient(new(http.Client))
var localCache = cache.NewCache(cache.WithLocalCache())
var zohoClient = zoho.NewZoho(params, httpClient, cache)
```

## Auth Worker with Redis cache
Auth Worker can retrieve and refresh Access Token automatically and be used later by `zohoClient`.

```golang
params := zoho2.ZohoParams{
    GrantToken:    "1000.xxxx.yyyy",
    ClientID:      "1000.abcdef",
    ClientSecret:  "client-secret",
    IamURL:        "https://accounts.zoho.com",
    CrmURL:        "https://zohoapis.com/crm",
    FileUploadURL: "https://content.zohoapis.com/crm",
    ZGID:          "12345677",
    Version:       crm.VersionV2,
}

var ctx = context.Background()
rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
})

var onSetFunc = func(key string, value interface{}, expire int) error {
    return rdb.Set(ctx, key, value, time.Duration(expire)*time.Second).Err()
}

var onGetFunc = func(key string) (value interface{}, err error) {
    val, err := rdb.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }

    return val, nil
}

cache := cache2.NewCache(cache2.WithExternalCache(cache2.Option{
    SetFunc: onSetFunc,
    GetFunc: onGetFunc,
}))
storage := storage2.NewStorage(cache)
authParams := oauth.ZohoAuthParams{
    ClientID:     params.ClientID,
    ClientSecret: params.ClientSecret,
    GrantToken:   params.GrantToken,
    IamURL:       params.IamURL,
}
authClient := oauth.NewZohoAuthClient(authParams, http2.NewHttpClient(new(http.Client)), storage)
wParams := worker.AuthWorkerParams{
    SecondsBeforeRefreshToken: 30,
}
w := worker.NewAuthWorker(authClient, wParams)

w.OnError(func(err error) {
    fmt.Println(err)
    w.Stop()
})
w.Start()
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)