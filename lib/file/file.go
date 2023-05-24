package file

import (
	"context"
)

type Uploader interface {
	// Upload uploads a file to the file storage and returns the file's URL
	Upload(ctx context.Context, name string, file []byte) error
}

type Downloader interface {
	// Download downloads a file from the file storage and returns the file's URL
	Download(ctx context.Context, name string) ([]byte, error)
}

type Lister interface {
	// List lists all files in the file storage
	List(ctx context.Context, path string) ([]string, error)
}

type Deleter interface {
	// Delete deletes a file from the file storage
	Delete(ctx context.Context, name string) error
}

type FileHandler interface {
	Uploader
	Downloader
	Lister
	Deleter
}
