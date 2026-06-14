package models

type Instance struct {
	ID            string
	Address       string
	Healthy       bool
	Weight        int
	CurrentWeight int
}