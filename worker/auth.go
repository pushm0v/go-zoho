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
	tickerChan    chan bool
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

	err := aw.client.GenerateToken()
	if err != nil {
		aw.fireError(err)
		return
	}
	done := aw.loopMain()
	aw.tickerChan = done
	<-done
}

func (aw *authWorker) Stop() {
	aw.isStart = false
	if aw.tickerChan != nil {
		aw.tickerChan <- true
	}
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

func (aw *authWorker) every(duration time.Duration, work func(time.Time) bool) chan bool {
	ticker := time.NewTicker(duration)
	stop := make(chan bool, 1)

	go func() {
		for {
			select {
			case time := <-ticker.C:
				if !work(time) {
					stop <- true
				}
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func (aw *authWorker) loopMain() chan bool {
	aw.isStart = true
	var err error
	var stop = aw.every(1*time.Second, func(time.Time) bool {
		remainingSeconds := aw.client.TokenExpireTime()
		if remainingSeconds <= aw.params.SecondsBeforeRefreshToken {
			err = aw.client.RefreshToken()
			if err != nil {
				aw.fireError(err)
			}
			aw.fireSuccess()
		}
		return true
	})

	return stop
}
