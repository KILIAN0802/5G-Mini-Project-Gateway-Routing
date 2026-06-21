package models

import "sync/atomic"

type Instance struct {
	ID            string
	Address       string
	Healthy       atomic.Bool
	Weight       int32
	CurrentWeight int32
	ActiveRequest int32 // Kiểu dữ liệu chính xác 32 bit khác với int có thể là 32 hoặc 64
}