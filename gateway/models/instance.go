package models

import (
	"net/url"
	"sync/atomic"
)

type Instance struct {
	ID            string
	Address       string
	TargetURL     *url.URL
	TargetURLStr  string
	Healthy       atomic.Bool
	Weight        int32
	CurrentWeight int32
	ActiveRequest int32
}