package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zemartins81/encoderVideoGolang/domain"

	uuid "github.com/satori/go.uuid"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {

	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIDIsNotAUuid(t *testing.T) {
	video := domain.NewVideo()
	video.ID = "123"
	video.ResourceID = "abc"
	video.FilePath = "def"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceID = "abc"
	video.FilePath = "def"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Nil(t, err)
}
