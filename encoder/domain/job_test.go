package domain_test

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/zemartins81/encoderVideoGolang/domain"
	"testing"
	"time"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("path", "converted", video)
	require.NotNil(t, job)
	require.Nil(t, err)
}
