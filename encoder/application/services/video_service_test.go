package services_test

import (
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/zemartins81/encoderVideoGolang/application/repositories"
	"github.com/zemartins81/encoderVideoGolang/application/services"
	"github.com/zemartins81/encoderVideoGolang/domain"
	"github.com/zemartins81/encoderVideoGolang/framework/database"
	"log"
	"testing"
	"time"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func prepare() (*domain.Video, *repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "emilly.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepositoryDb(db)

	return video, repo
}

func TestVideoServiceDownload(t *testing.T) {
	video, repo := prepare()
	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("encodervideotest")
	require.Nil(t, err)
}
