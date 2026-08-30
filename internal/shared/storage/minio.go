// Package storage 对象存储封装(MinIO / S3 兼容)。
package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	mc        *minio.Client
	presign   *minio.Client // 按浏览器 Host 签名，禁止签完再改 Host
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

	// Region 固定可跳过 GetBucketLocation；预签名 client 指向 127.0.0.1 时不能再去拨号探活
	region := g.Cfg().MustGet(ctx, "minio.region", "us-east-1").String()
	if region == "" {
		region = "us-east-1"
	}

	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 MinIO 客户端失败: %w", err)
	}

	presign := mc
	presignEndpoint, presignSSL := endpoint, useSSL
	if publicURL != "" {
		if u, err := url.Parse(publicURL); err == nil && u.Host != "" {
			presignEndpoint = u.Host
			presignSSL = strings.EqualFold(u.Scheme, "https")
		}
	}
	if presignEndpoint != endpoint || presignSSL != useSSL {
		presign, err = minio.New(presignEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: presignSSL,
			Region: region,
		})
		if err != nil {
			return nil, fmt.Errorf("创建 MinIO 预签名客户端失败: %w", err)
		}
	}

	c := &Client{mc: mc, presign: presign, bucket: bucket, publicURL: publicURL}
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

// Get 读取默认桶对象全文(预览/前台拉密文用)。
func (c *Client) Get(ctx context.Context, objectKey string) ([]byte, error) {
	return c.GetIn(ctx, c.bucket, objectKey)
}

// GetIn 按桶读取对象。
func (c *Client) GetIn(ctx context.Context, bucket, objectKey string) ([]byte, error) {
	if bucket == "" {
		bucket = c.bucket
	}
	key := strings.TrimLeft(objectKey, "/")
	if key == "" || strings.Contains(key, "..") {
		return nil, fmt.Errorf("对象路径无效")
	}
	obj, err := c.mc.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("读取对象失败: %w", err)
	}
	defer obj.Close()
	return io.ReadAll(obj)
}

// StatIn 按桶取对象元数据。
func (c *Client) StatIn(ctx context.Context, bucket, objectKey string) (minio.ObjectInfo, error) {
	if bucket == "" {
		bucket = c.bucket
	}
	key := strings.TrimLeft(objectKey, "/")
	if key == "" || strings.Contains(key, "..") {
		return minio.ObjectInfo{}, fmt.Errorf("对象路径无效")
	}
	info, err := c.mc.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, fmt.Errorf("读取对象失败: %w", err)
	}
	return info, nil
}

// OpenIn 按桶打开对象；length<=0 表示从 offset 读到结尾。
func (c *Client) OpenIn(ctx context.Context, bucket, objectKey string, offset, length int64) (io.ReadCloser, minio.ObjectInfo, error) {
	if bucket == "" {
		bucket = c.bucket
	}
	key := strings.TrimLeft(objectKey, "/")
	if key == "" || strings.Contains(key, "..") {
		return nil, minio.ObjectInfo{}, fmt.Errorf("对象路径无效")
	}
	opts := minio.GetObjectOptions{}
	if length > 0 {
		if err := opts.SetRange(offset, offset+length-1); err != nil {
			return nil, minio.ObjectInfo{}, err
		}
	} else if offset > 0 {
		opts.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}
	obj, err := c.mc.GetObject(ctx, bucket, key, opts)
	if err != nil {
		return nil, minio.ObjectInfo{}, fmt.Errorf("读取对象失败: %w", err)
	}
	info, err := obj.Stat()
	if err != nil {
		_ = obj.Close()
		return nil, minio.ObjectInfo{}, fmt.Errorf("读取对象失败: %w", err)
	}
	return obj, info, nil
}

// ParseURL 从 my-media / my-storage 地址解析桶和 key。
func ParseURL(raw, defaultBucket string) (bucket, key string) {
	u := strings.TrimSpace(raw)
	if u == "" {
		return "", ""
	}
	if i := strings.IndexAny(u, "?#"); i >= 0 {
		u = u[:i]
	}
	for _, b := range []string{"my-storage", "my-media"} {
		marker := "/" + b + "/"
		if i := strings.Index(u, marker); i >= 0 {
			key = strings.TrimLeft(u[i+len(marker):], "/")
			if key == "" || strings.Contains(key, "..") {
				return "", ""
			}
			return b, key
		}
	}
	if defaultBucket != "" {
		key = strings.TrimLeft(u, "/")
		if strings.HasPrefix(key, defaultBucket+"/") {
			key = key[len(defaultBucket)+1:]
		}
		if key == "" || strings.Contains(key, "..") || strings.Contains(key, "://") {
			return "", ""
		}
		return defaultBucket, key
	}
	return "", ""
}

// SignPlayURL 私有桶 my-storage 签发可播/可下的 GET 地址；其它地址原样返回。
func SignPlayURL(ctx context.Context, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || !strings.Contains(raw, "/my-storage/") {
		return raw
	}
	bucket, key := ParseURL(raw, "my-storage")
	if key == "" {
		return raw
	}
	c, err := Get(ctx)
	if err != nil {
		return raw
	}
	u, err := c.PresignGetIn(ctx, bucket, key, 2*time.Hour)
	if err != nil {
		g.Log().Warningf(ctx, "签发 my-storage 播放地址失败: %v", err)
		return raw
	}
	return u
}

// PresignGetIn 按桶签发预签名 GET。
func (c *Client) PresignGetIn(ctx context.Context, bucket, objectKey string, expire time.Duration) (string, error) {
	if bucket == "" {
		bucket = c.bucket
	}
	key := strings.TrimLeft(objectKey, "/")
	if key == "" || strings.Contains(key, "..") {
		return "", fmt.Errorf("对象路径无效")
	}
	if expire <= 0 {
		expire = 2 * time.Hour
	}
	signer := c.presign
	if signer == nil {
		signer = c.mc
	}
	u, err := signer.PresignedGetObject(ctx, bucket, key, expire, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// ObjectKeyFromURL 从本桶 publicURL 解析 object key; 非法外链返回空。
func (c *Client) ObjectKeyFromURL(raw string) string {
	u := strings.TrimSpace(raw)
	if u == "" {
		return ""
	}
	if i := strings.IndexAny(u, "?#"); i >= 0 {
		u = u[:i]
	}
	if !strings.Contains(u, "://") {
		key := strings.TrimLeft(u, "/")
		if strings.HasPrefix(key, c.bucket+"/") {
			key = key[len(c.bucket)+1:]
		}
		if key == "" || strings.Contains(key, "..") {
			return ""
		}
		return key
	}
	prefix := strings.TrimRight(c.publicURL, "/") + "/" + c.bucket + "/"
	if strings.HasPrefix(u, prefix) {
		key := strings.TrimLeft(u[len(prefix):], "/")
		if key == "" || strings.Contains(key, "..") {
			return ""
		}
		return key
	}
	// 兼容 host.docker.internal / minio 内网 host, 只认 /{bucket}/ 路径
	marker := "/" + c.bucket + "/"
	if i := strings.Index(u, marker); i >= 0 {
		key := strings.TrimLeft(u[i+len(marker):], "/")
		if key == "" || strings.Contains(key, "..") {
			return ""
		}
		return key
	}
	return ""
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
