package repositories_test

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/zemartins81/encoderVideoGolang/application/repositories"
	"github.com/zemartins81/encoderVideoGolang/domain"
	"github.com/zemartins81/encoderVideoGolang/framework/database"
	"testing"
	"time"
)

func TestNewVideoRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepositoryDb(db)
	_, err := repo.Insert(video)
	if err != nil {
		return
	}

	v, err := repo.Find(video.ID)

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)
}
