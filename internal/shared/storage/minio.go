// Package storage 对象存储封装(MinIO / S3 兼容)。
package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	mc        *minio.Client
	bucket    string
	publicURL string
}

var (
	once    sync.Once
	shared  *Client
	initErr error
)

// Get 懒加载单例。配置变更需重启进程(Nacos Watch 暂不热重建客户端)。
func Get(ctx context.Context) (*Client, error) {
	once.Do(func() {
		shared, initErr = newClient(ctx)
	})
	return shared, initErr
}

func newClient(ctx context.Context) (*Client, error) {
	endpoint := g.Cfg().MustGet(ctx, "minio.endpoint", "127.0.0.1:19000").String()
	accessKey := g.Cfg().MustGet(ctx, "minio.accessKey", "minioadmin").String()
	secretKey := g.Cfg().MustGet(ctx, "minio.secretKey", "minioadmin123").String()
	bucket := g.Cfg().MustGet(ctx, "minio.bucket", "my-media").String()
	useSSL := g.Cfg().MustGet(ctx, "minio.useSSL", false).Bool()
	publicURL := strings.TrimRight(g.Cfg().MustGet(ctx, "minio.publicURL", "http://127.0.0.1:19000").String(), "/")

	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 MinIO 客户端失败: %w", err)
	}

	c := &Client{mc: mc, bucket: bucket, publicURL: publicURL}
	if err = c.ensureBucket(ctx); err != nil {
		return nil, err
	}
	// 本地开发默认开匿名读, 便于 publicURL 直链预览; 生产可配 minio.publicRead=false
	if g.Cfg().MustGet(ctx, "minio.publicRead", true).Bool() {
		if err = c.ensurePublicRead(ctx); err != nil {
			g.Log().Warningf(ctx, "设置 bucket 匿名读失败(可手动在 MinIO 控制台开): %v", err)
		}
	}
	g.Log().Infof(ctx, "对象存储就绪 endpoint=%s bucket=%s", endpoint, bucket)
	return c, nil
}

func (c *Client) ensureBucket(ctx context.Context) error {
	ok, err := c.mc.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("检查 bucket 失败: %w", err)
	}
	if ok {
		return nil
	}
	if err = c.mc.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("创建 bucket 失败: %w", err)
	}
	g.Log().Infof(ctx, "已创建 MinIO bucket: %s", c.bucket)
	return nil
}

func (c *Client) ensurePublicRead(ctx context.Context) error {
	policy := fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {"AWS": ["*"]},
    "Action": ["s3:GetObject"],
    "Resource": ["arn:aws:s3:::%s/*"]
  }]
}`, c.bucket)
	return c.mc.SetBucketPolicy(ctx, c.bucket, policy)
}

// Bucket 当前桶名。
func (c *Client) Bucket() string { return c.bucket }

// PublicURL 拼可访问 URL(依赖 bucket 策略或网关代理; 本地开发通常开匿名读)。
func (c *Client) PublicURL(objectKey string) string {
	key := strings.TrimLeft(objectKey, "/")
	return fmt.Sprintf("%s/%s/%s", c.publicURL, c.bucket, key)
}

// Put 上传对象(流式, 不整文件进内存)。
func (c *Client) Put(ctx context.Context, objectKey string, r io.Reader, size int64, contentType string) (string, error) {
	key := strings.TrimLeft(objectKey, "/")
	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}
	_, err := c.mc.PutObject(ctx, c.bucket, key, r, size, opts)
	if err != nil {
		return "", fmt.Errorf("上传到 MinIO 失败: %w", err)
	}
	return c.PublicURL(key), nil
}
