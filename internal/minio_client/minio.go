package minio_client

import (
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	MinioClient *minio.Client
	MinioCfg    *config.Minio
}

func (m *Minio) New(cfg *config.Minio) error {
	var err error
	m.MinioCfg = cfg
	m.MinioClient, err = minio.New(m.MinioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.MinioCfg.MinioUser, m.MinioCfg.MinioPassword, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}
	return nil
}
