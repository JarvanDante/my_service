package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

// PartInfo 已上传分片。
type PartInfo struct {
	PartNumber int
	ETag       string
	Size       int64
}

// NewMultipart 发起分片上传, 返回 MinIO uploadId。
func (c *Client) NewMultipart(ctx context.Context, objectKey, contentType string) (string, error) {
	key := strings.TrimLeft(objectKey, "/")
	core := minio.Core{Client: c.mc}
	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}
	uploadID, err := core.NewMultipartUpload(ctx, c.bucket, key, opts)
	if err != nil {
		return "", fmt.Errorf("发起分片上传失败: %w", err)
	}
	return uploadID, nil
}

// PresignUploadPart 预签名单片 PUT URL(客户端直传 MinIO)。
func (c *Client) PresignUploadPart(ctx context.Context, objectKey, uploadID string, partNumber int, expires time.Duration) (string, error) {
	if partNumber < 1 {
		return "", fmt.Errorf("partNumber 无效")
	}
	if expires <= 0 {
		expires = 2 * time.Hour
	}
	key := strings.TrimLeft(objectKey, "/")
	params := url.Values{}
	params.Set("uploadId", uploadID)
	params.Set("partNumber", strconv.Itoa(partNumber))
	u, err := c.mc.Presign(ctx, http.MethodPut, c.bucket, key, expires, params)
	if err != nil {
		return "", fmt.Errorf("预签名分片失败: %w", err)
	}
	return u.String(), nil
}

// PutObjectPart 服务端代传一片(H5 不直连 MinIO, 避免 CORS / 内网地址问题)。
func (c *Client) PutObjectPart(ctx context.Context, objectKey, uploadID string, partNumber int, r io.Reader, size int64) (string, error) {
	if partNumber < 1 {
		return "", fmt.Errorf("partNumber 无效")
	}
	key := strings.TrimLeft(objectKey, "/")
	core := minio.Core{Client: c.mc}
	part, err := core.PutObjectPart(ctx, c.bucket, key, uploadID, partNumber, r, size, minio.PutObjectPartOptions{})
	if err != nil {
		return "", fmt.Errorf("上传分片失败: %w", err)
	}
	return strings.Trim(part.ETag, `"`), nil
}

// ListParts 列出已上传分片(用于断点续传)。
func (c *Client) ListParts(ctx context.Context, objectKey, uploadID string) ([]PartInfo, error) {
	key := strings.TrimLeft(objectKey, "/")
	core := minio.Core{Client: c.mc}
	var out []PartInfo
	marker := 0
	for {
		res, err := core.ListObjectParts(ctx, c.bucket, key, uploadID, marker, 1000)
		if err != nil {
			return nil, fmt.Errorf("列出分片失败: %w", err)
		}
		for _, p := range res.ObjectParts {
			out = append(out, PartInfo{
				PartNumber: p.PartNumber,
				ETag:       strings.Trim(p.ETag, `"`),
				Size:       p.Size,
			})
		}
		if !res.IsTruncated {
			break
		}
		marker = res.NextPartNumberMarker
	}
	return out, nil
}

// CompleteMultipart 合并分片。
func (c *Client) CompleteMultipart(ctx context.Context, objectKey, uploadID string, parts []PartInfo) (string, error) {
	key := strings.TrimLeft(objectKey, "/")
	core := minio.Core{Client: c.mc}
	complete := make([]minio.CompletePart, 0, len(parts))
	for _, p := range parts {
		complete = append(complete, minio.CompletePart{
			PartNumber: p.PartNumber,
			ETag:       strings.Trim(p.ETag, `"`),
		})
	}
	info, err := core.CompleteMultipartUpload(ctx, c.bucket, key, uploadID, complete, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("合并分片失败: %w", err)
	}
	_ = info
	return c.PublicURL(key), nil
}

// AbortMultipart 取消未完成的分片上传。
func (c *Client) AbortMultipart(ctx context.Context, objectKey, uploadID string) error {
	key := strings.TrimLeft(objectKey, "/")
	core := minio.Core{Client: c.mc}
	if err := core.AbortMultipartUpload(ctx, c.bucket, key, uploadID); err != nil {
		return fmt.Errorf("取消分片上传失败: %w", err)
	}
	return nil
}
