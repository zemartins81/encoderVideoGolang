package services

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/zemartins81/encoderVideoGolang/application/repositories"
	"github.com/zemartins81/encoderVideoGolang/domain"
	"io"
	"log"
	"os"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVIdeoService() VideoService {
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	_, err := f.Write(body)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("video %v has been saved", v.Video.ID)
	return nil

}
