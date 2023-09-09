package acaer

import (
	"microservice/internal/app/looncan"
	"microservice/internal/utils"
)

type Acaer struct {
	ID      uint64
	Name    string
	Version string

	looncans []looncan.Looncan
}

func (a *Acaer) Entsian() {
	for i := 0; i < 3; i++ {
		a.looncans = append(a.looncans, looncan.Looncan{
			Name:  utils.RandomString(10),
			Value: utils.RandomString(6),
		})
	}
}

func (a *Acaer) getLooncans() []looncan.Looncan {
	return a.looncans
}
