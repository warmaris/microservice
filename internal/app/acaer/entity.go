package acaer

import (
	"microservice/internal/app/looncan"
	"microservice/internal/errors"
	"microservice/internal/utils"
)

type Acaer struct {
	ID      uint64
	Name    string
	Version string

	looncans []looncan.Looncan
}

// Entsian is some kind of business-logic method. It can do some business stuff, and create related entities or events.
// Entities are linked by memory. Storage creates links in database with ids.
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

type Validator struct {
	allowedVersions []string
}

func NewValidator(versions []string) *Validator {
	return &Validator{allowedVersions: versions}
}

func (v *Validator) Validate(acaer Acaer) error {
	for _, allowedVersion := range v.allowedVersions {
		if allowedVersion == acaer.Version {
			return nil
		}
	}

	return errors.NewValidationError("version", "provided version is not supported")
}
