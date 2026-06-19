package models

type Instance struct {
	ID            string
	Address       string
	Healthy       bool
	Weight        int
	CurrentWeight int
	ActiveRequest int32 // Kiểu dữ liệu chính xác 32 bit khác với int có thể là 32 hoặc 64
}