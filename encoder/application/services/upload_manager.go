package services

import (
	"context"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

// UploadObject uploads an object to a storage client.
//
// Parameters:
// - objectpath: the path of the object to upload.
// - client: the storage client where the object will be uploaded.
// - ctx: the context for the upload operation.
// Returns an error if any.
func (vu *VideoUpload) UploadObject(objectpath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectpath, os.Getenv("localStoragePath")+"/")

	f, err := os.Open(objectpath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}
	return nil

}
