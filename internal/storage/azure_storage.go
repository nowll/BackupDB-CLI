package storage

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/nowll/BackupDB-CLI/internal/domain"
)

type azureStorage struct {
	client    *azblob.Client
	container string
}

func NewAzureStorage(config domain.StorageConfig) (domain.StorageRepository, error) {
	client, err := azblob.NewClientFromConnectionString(
		config.Options["connection_string"],
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create azure client: %w", err)
	}

	return &azureStorage{
		client:    client,
		container: config.Bucket,
	}, nil
}

func (s *azureStorage) Upload(ctx context.Context, key string, data []byte) error {
	_, err := s.client.UploadBuffer(
		ctx,
		s.container,
		key,
		data,
		nil,
	)

	if err != nil {
		return fmt.Errorf("upload to azure: %w", err)
	}

	return nil
}

func (s *azureStorage) Download(ctx context.Context, key string) ([]byte, error) {
	response, err := s.client.DownloadStream(ctx, s.container, key, nil)
	if err != nil {
		return nil, fmt.Errorf("download from azure: %w", err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read azure response: %w", err)
	}

	return data, nil
}

func (s *azureStorage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteBlob(ctx, s.container, key, nil)
	if err != nil {
		return fmt.Errorf("delete from azure: %w", err)
	}

	return nil
}

func (s *azureStorage) List(ctx context.Context, prefix string) ([]string, error) {
	pager := s.client.NewListBlobsFlatPager(s.container, &azblob.ListBlobsFlatOptions{
		Prefix: &prefix,
	})

	var keys []string
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list azure blobs: %w", err)
		}

		for _, blob := range resp.Segment.BlobItems {
			keys = append(keys, *blob.Name)
		}
	}

	return keys, nil
}

func (s *azureStorage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.ServiceClient().GetProperties(ctx, s.container)
	if err != nil {
		return false, nil
	}

	return true, nil
}
