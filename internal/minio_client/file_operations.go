package minio_client

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

func (m *Minio) UploadFile(objectName string, file multipart.File, fileSize int64) (string, error) {
	ctx := context.Background()

	fileBytes := make([]byte, fileSize)
	if _, err := file.Read(fileBytes); err != nil {
		return "", err
	}

	reader := bytes.NewReader(fileBytes)

	contentType := "application/octet-stream"

	_, err := m.MinioClient.PutObject(ctx, m.MinioCfg.MinioBucket, objectName, reader, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}

	 fileURL := fmt.Sprintf("https://%s/%s/%s", m.MinioCfg.Endpoint, m.MinioCfg.MinioBucket, objectName)
	 return fileURL, nil
}

func (m *Minio) DownloadFile(objectName string) (multipart.File, error) {
    ctx := context.Background()

    reader, err := m.MinioClient.GetObject(ctx, m.MinioCfg.MinioBucket, objectName, minio.GetObjectOptions{})
    if err != nil {
        return nil, err
    }

    return reader, nil
}

func (m *Minio) DeleteFile(objectName string) error {
	ctx := context.Background()

	err := m.MinioClient.RemoveObject(ctx, m.MinioCfg.MinioBucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Successfully deleted %s from MinIO\n", objectName)
	return nil
}
