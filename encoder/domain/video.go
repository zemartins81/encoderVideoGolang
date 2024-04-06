package domain

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type Video struct {
	ID         int       `valid:"uuid"`
	ResourceID string    `valid:"notnull"`
	FilePath   string    `valid:"notnull"`
	CreatedAt  time.Time `valid:"notnull"`
}

var validate *validator.Validate

func init() {
  validate = validator.New(validator.WithRequiredStructEnabled())
}
 
func NewVideo() *Video {
	return &Video{}
}

func (video *Video) Validate() error {

  err := validate.Struct(video)

	if err != nil {
		return err
	}

	return nil

}
