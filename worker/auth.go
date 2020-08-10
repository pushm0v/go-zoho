package worker

import (
	"time"

	"github.com/pushm0v/go-zoho/oauth"
)

type AuthWorker interface {
	Start()
	Stop()
	OnError(f func(err error))
	OnSuccess(f func())
}

type authWorker struct {
	client        oauth.ZohoAuthClient
	isStart       bool
	onErrorFunc   func(error)
	onSuccessFunc func()
	params        AuthWorkerParams
}

type AuthWorkerParams struct {
	SecondsBeforeRefreshToken int
}

func NewAuthWorker(client oauth.ZohoAuthClient, params AuthWorkerParams) AuthWorker {
	if params.SecondsBeforeRefreshToken <= 0 {
		params.SecondsBeforeRefreshToken = 30 // 30 seconds before expire to generate new token
	}
	return &authWorker{
		client: client,
		params: params,
	}
}

func (aw *authWorker) Start() {
	aw.loopMain()
}

func (aw *authWorker) Stop() {
	aw.isStart = false
}

func (aw *authWorker) OnError(f func(err error)) {
	aw.onErrorFunc = f
}
func (aw *authWorker) OnSuccess(f func()) {
	aw.onSuccessFunc = f
}

func (aw *authWorker) fireError(err error) {
	if aw.onErrorFunc != nil {
		aw.onErrorFunc(err)
	}
}

func (aw *authWorker) fireSuccess() {
	if aw.onSuccessFunc != nil {
		aw.onSuccessFunc()
	}
}

func (aw *authWorker) loopMain() {
	aw.isStart = true
	err := aw.client.GenerateToken()
	if err != nil {
		aw.fireError(err)
		return
	}
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			remainingSeconds := aw.client.TokenExpireTime()
			if remainingSeconds <= aw.params.SecondsBeforeRefreshToken {
				err = aw.client.RefreshToken()
				if err != nil {
					aw.fireError(err)
					return
				}
				aw.fireSuccess()
			}
		}
	}
}
