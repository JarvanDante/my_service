package paas

import (
	"context"
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
