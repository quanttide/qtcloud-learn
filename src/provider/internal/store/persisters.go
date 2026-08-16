package store

// 快照持久化后端：文件（本地 dev / 测试）与 OSS（生产 FC 容器——实例临时盘不持久）。

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// FilePersister 本地文件后端（原子写：临时文件 + rename）。
type FilePersister struct {
	dir string
}

func NewFilePersister(dir string) *FilePersister {
	return &FilePersister{dir: dir}
}

func (p *FilePersister) Load(key string) ([]byte, error) {
	raw, err := os.ReadFile(filepath.Join(p.dir, key))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return raw, nil
}

func (p *FilePersister) Store(key string, data []byte) error {
	if err := os.MkdirAll(p.dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(p.dir, key)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// OSSPersister OSS 对象存储后端（生产：FC 容器实例盘在重建后丢失，数据存 OSS 跨实例/发版不丢）。
type OSSPersister struct {
	bucket   *oss.Bucket
	keyPrefix string
}

// NewOSSPersister 创建 OSS 后端。
// 凭证走 FC 运行时环境变量（ALIYUN_ACCESS_KEY_ID / ALIYUN_ACCESS_KEY_SECRET，
// 与 CI 同一组 RAM 凭证；生产建议独立 RAM 角色）。
func NewOSSPersister(bucketName, endpoint, keyPrefix, accessKeyID, accessKeySecret string) (*OSSPersister, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("oss client: %w", err)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("oss bucket %s: %w", bucketName, err)
	}
	return &OSSPersister{bucket: bucket, keyPrefix: keyPrefix}, nil
}

func (p *OSSPersister) Load(key string) ([]byte, error) {
	raw, err := p.bucket.GetObject(p.keyPrefix + key)
	if err != nil {
		var ossErr oss.ServiceError
		if errors.As(err, &ossErr) && ossErr.StatusCode == 404 {
			return nil, nil
		}
		return nil, err
	}
	defer raw.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(raw); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *OSSPersister) Store(key string, data []byte) error {
	return p.bucket.PutObject(p.keyPrefix+key, bytes.NewReader(data), oss.ContentType("application/json"))
}
