package storage

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/nowll/BackupDB-CLI/internal/domain"
	"google.golang.org/api/iterator"
)

type gcsStorage struct {
	client *storage.Client
	bucket string
}

func NewGCSStorage(config domain.StorageConfig) (domain.StorageRepository, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("create gcs client: %w", err)
	}

	return &gcsStorage{
		client: client,
		bucket: config.Bucket,
	}, nil
}

func (s *gcsStorage) Upload(ctx context.Context, key string, data []byte) error {
	obj := s.client.Bucket(s.bucket).Object(key)
	writer := obj.NewWriter(ctx)

	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("write to gcs: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close gcs writer: %w", err)
	}

	return nil
}

func (s *gcsStorage) Download(ctx context.Context, key string) ([]byte, error) {
	obj := s.client.Bucket(s.bucket).Object(key)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("create gcs reader: %w", err)
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read from gcs: %w", err)
	}

	return data, nil
}

func (s *gcsStorage) Delete(ctx context.Context, key string) error {
	obj := s.client.Bucket(s.bucket).Object(key)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("delete from gcs: %w", err)
	}

	return nil
}

func (s *gcsStorage) List(ctx context.Context, prefix string) ([]string, error) {
	it := s.client.Bucket(s.bucket).Objects(ctx, &storage.Query{
		Prefix: prefix,
	})

	var keys []string
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iterate gcs objects: %w", err)
		}
		keys = append(keys, attrs.Name)
	}

	return keys, nil
}

func (s *gcsStorage) Exists(ctx context.Context, key string) (bool, error) {
	obj := s.client.Bucket(s.bucket).Object(key)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
