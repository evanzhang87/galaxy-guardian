package util

import (
	"github.com/parnurzeal/gorequest"
	"time"
)

var Request = gorequest.New().Timeout(time.Second * 15)

func init() {
	Request.Header.Set("Content-Type", "application/json")
	Request.Header.Set("Accept", "application/json")
	Request.Header.Set("User-Agent", "galaxy-guardian")
}
