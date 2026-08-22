package paas

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

type StorageCreateIn struct {
	Filename    string
	Biz         string
	ContentType string
	SizeBytes   int64
	Remark      string
}

type StorageCreateOut struct {
	Id          string `json:"id"`
	UploadUrl   string `json:"upload_url"`
	Method      string `json:"method"`
	Bucket      string `json:"bucket"`
	Key         string `json:"key"`
	ExpireSec   int    `json:"expire_sec"`
	PublicUrl   string `json:"public_url"`
	ContentType string `json:"content_type"`
}

type StorageConfirmOut struct {
	Id        string `json:"id"`
	Status    int    `json:"status"`
	SizeBytes int64  `json:"size_bytes"`
	PublicUrl string `json:"public_url"`
}

func storageBase(ctx context.Context) (string, error) {
	base := strings.TrimRight(g.Cfg().MustGet(ctx, "paas.storage_base").String(), "/")
	if base == "" {
		return "", gerror.New("未配置 paas.storage_base")
	}
	return base, nil
}

func storageClient(ctx context.Context) (string, *gclient.Client, error) {
	base, err := storageBase(ctx)
	if err != nil {
		return "", nil, err
	}
	key := g.Cfg().MustGet(ctx, "paas.app_key").String()
	secret := g.Cfg().MustGet(ctx, "paas.app_secret").String()
	if key == "" || secret == "" {
		return "", nil, gerror.New("未配置 paas.app_key / app_secret")
	}
	c := g.Client().ContentJson().
		SetHeader("X-App-Key", key).
		SetHeader("X-App-Secret", secret)
	return base, c, nil
}

func CreateStorageObject(ctx context.Context, in StorageCreateIn) (*StorageCreateOut, error) {
	base, c, err := storageClient(ctx)
	if err != nil {
		return nil, err
	}
	r, err := c.Post(ctx, base+"/open/objects", g.Map{
		"filename":     in.Filename,
		"biz":          in.Biz,
		"content_type": in.ContentType,
		"size_bytes":   in.SizeBytes,
		"remark":       in.Remark,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "请求统一存储失败")
	}
	defer r.Close()
	out := &StorageCreateOut{}
	if err := parseEnvelope(r.ReadAll(), out); err != nil {
		return nil, err
	}
	return out, nil
}

func ConfirmStorageObject(ctx context.Context, id string) (*StorageConfirmOut, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, gerror.New("对象ID必填")
	}
	base, c, err := storageClient(ctx)
	if err != nil {
		return nil, err
	}
	r, err := c.Post(ctx, base+"/open/objects/"+id+"/confirm", g.Map{})
	if err != nil {
		return nil, gerror.Wrap(err, "请求统一存储失败")
	}
	defer r.Close()
	out := &StorageConfirmOut{}
	if err := parseEnvelope(r.ReadAll(), out); err != nil {
		return nil, err
	}
	return out, nil
}

// PutUploadURL 按统一存储签发的预签名地址写入文件。不要带 Content-Type。
// 预签名 Host 是给浏览器的 127.0.0.1:19000；容器内改拨 minio.endpoint，Host 头保持原值以免验签失败。
func PutUploadURL(ctx context.Context, uploadURL string, body io.Reader, size int64) error {
	uploadURL = strings.TrimSpace(uploadURL)
	if uploadURL == "" {
		return gerror.New("缺少统一存储上传地址")
	}
	u, err := url.Parse(uploadURL)
	if err != nil {
		return gerror.Wrap(err, "统一存储上传地址无效")
	}
	signedHost := u.Host
	if internal := strings.TrimSpace(g.Cfg().MustGet(ctx, "minio.endpoint").String()); internal != "" {
		host := u.Hostname()
		if host == "127.0.0.1" || host == "localhost" || host == "::1" {
			u.Host = internal
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), body)
	if err != nil {
		return gerror.Wrap(err, "构造统一存储上传请求失败")
	}
	req.Host = signedHost
	if size > 0 {
		req.ContentLength = size
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gerror.Wrap(err, "写入统一存储失败")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 180))
		msg := strings.TrimSpace(string(raw))
		if msg == "" {
			msg = resp.Status
		}
		return gerror.Newf("写入统一存储失败(%d): %s", resp.StatusCode, msg)
	}
	return nil
}
