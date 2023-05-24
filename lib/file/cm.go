package file

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/chartmuseum/storage"
)

var (
	ErrUnsupportedDriver = errors.New("unsupported driver")
)

type cloudStorage struct {
	backend storage.Backend
}

type CloudStorageConfig struct {
	// Driver is the driver to use for the file storage
	// Supported drivers: s3, gcs, azure
	Driver string
	// Bucket is the bucket to use for the file storage
	Bucket string
	// Prefix is the prefix to use for the file storage
	Prefix string
	// Region is the region to use for the file storage
	// Only used for s3
	Region string
	// Endpoint is the endpoint to use for the file storage
	// Only used for s3
	Endpoint string
	// SSE is the server-side encryption to use for the file storage
	// Only used for s3
	SSE string
}

func NewCloudStorage(config CloudStorageConfig) (FileHandler, error) {
	var backend storage.Backend
	switch config.Driver {
	case "s3":
		backend = storage.NewAmazonS3BackendWithCredentials(config.Bucket, config.Prefix, config.Region, config.Endpoint, config.SSE, credentials.NewCredentials(&credentials.EnvProvider{}))
	case "gcs":
		backend = storage.NewGoogleCSBackend(config.Bucket, config.Prefix)
	case "azure":
		backend = storage.NewMicrosoftBlobBackend(config.Bucket, config.Prefix)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedDriver, config.Driver)
	}
	return &cloudStorage{
		backend: backend,
	}, nil
}

func (c *cloudStorage) Upload(ctx context.Context, name string, bts []byte) error {
	return c.backend.PutObject(name, bts)
}

func (c *cloudStorage) Download(ctx context.Context, name string) ([]byte, error) {
	obj, err := c.backend.GetObject(name)
	if err != nil {
		return nil, err
	}
	return obj.Content, nil
}

func (c *cloudStorage) List(ctx context.Context, path string) ([]string, error) {
	objs, err := c.backend.ListObjects(path)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, obj := range objs {
		names = append(names, obj.Path)
	}
	return names, nil
}

func (c *cloudStorage) Delete(ctx context.Context, name string) error {
	return c.backend.DeleteObject(name)
}
