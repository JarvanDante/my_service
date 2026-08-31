package aiprovider

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/shared/aiprotocol"
)

// LocalName 自建 Python 工人。app_config.ai_provider 设为 local 后换脸/去衣走 Kafka。
const LocalName = "local"

type jobPublisher func(context.Context, aiprotocol.JobMessage) error

var (
	pubMu     sync.RWMutex
	publishFn jobPublisher
)

// SetFaceSwapPublisher 由启动装配注入 Kafka 投递。未注入时 Submit 直接失败并退款。
func SetFaceSwapPublisher(fn jobPublisher) {
	pubMu.Lock()
	defer pubMu.Unlock()
	publishFn = fn
}

type localProvider struct{}

func init() { Register(LocalName, &localProvider{}) }

func (p *localProvider) Name() string { return LocalName }

func (p *localProvider) Submit(ctx context.Context, in SubmitInput) (*SubmitOutput, error) {
	var (
		job aiprotocol.JobMessage
		err error
	)
	switch in.BizType {
	case entity.AiBizFaceSwap:
		job, err = BuildFaceSwapJob(ctx, in)
	case entity.AiBizUndress:
		job, err = BuildUndressJob(ctx, in)
	default:
		return nil, gerror.New("本地工人第一期只支持图片换脸/去衣")
	}
	if err != nil {
		return nil, err
	}
	pubMu.RLock()
	fn := publishFn
	pubMu.RUnlock()
	if fn == nil {
		return nil, gerror.New("本地换脸工人未就绪(Kafka 未装配)")
	}
	if err := fn(ctx, job); err != nil {
		return nil, gerror.Wrap(err, "投递任务失败")
	}
	return &SubmitOutput{ProviderTaskId: LocalName + "-" + in.TaskNo}, nil
}

func (p *localProvider) Query(_ context.Context, _ string) (int, map[string]any, string, error) {
	return StatusRunning, nil, "", nil
}

// BuildFaceSwapJob 从 aitask 的 params / input_url 收出对象引用。
func BuildFaceSwapJob(ctx context.Context, in SubmitInput) (aiprotocol.JobMessage, error) {
	if strings.TrimSpace(in.TaskNo) == "" {
		return aiprotocol.JobMessage{}, gerror.New("任务号为空")
	}
	mediaType := strings.ToLower(strings.TrimSpace(aiprotocol.FirstString(in.Params, "media_type")))
	if mediaType == "" {
		mediaType = aiprotocol.MediaPhoto
	}
	if mediaType != aiprotocol.MediaPhoto {
		return aiprotocol.JobMessage{}, gerror.New("第一期只支持图片换脸")
	}

	bucket, site := jobBucketSite(ctx)
	source, err := pickRef(in.Params, []string{"source", "source_url", "source_key"}, in.InputURL, bucket)
	if err != nil {
		return aiprotocol.JobMessage{}, gerror.New("换脸缺少人脸图(input_url / source)")
	}
	target, err := pickRef(in.Params, []string{"target", "target_url", "target_key", "cover"}, "", bucket)
	if err != nil {
		return aiprotocol.JobMessage{}, gerror.New("换脸缺少目标图(target_url / target)")
	}

	return aiprotocol.JobMessage{
		SchemaVersion: aiprotocol.SchemaVersion,
		JobID:         in.TaskNo,
		Biz:           aiprotocol.BizFaceSwap,
		MediaType:     aiprotocol.MediaPhoto,
		Source:        source,
		Target:        target,
		Output: aiprotocol.OutputRef{
			Bucket: bucket,
			Prefix: site + "/ai/out/" + in.TaskNo + "/",
		},
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// BuildUndressJob 去衣只要一张待处理图（Go 上传到 my-storage 后把 URL 放 input_url）。
func BuildUndressJob(ctx context.Context, in SubmitInput) (aiprotocol.JobMessage, error) {
	if strings.TrimSpace(in.TaskNo) == "" {
		return aiprotocol.JobMessage{}, gerror.New("任务号为空")
	}
	mediaType := strings.ToLower(strings.TrimSpace(aiprotocol.FirstString(in.Params, "media_type")))
	if mediaType == "" {
		mediaType = aiprotocol.MediaPhoto
	}
	if mediaType != aiprotocol.MediaPhoto {
		return aiprotocol.JobMessage{}, gerror.New("第一期只支持图片去衣")
	}

	bucket, site := jobBucketSite(ctx)
	source, err := pickRef(in.Params, []string{"source", "source_url", "source_key"}, in.InputURL, bucket)
	if err != nil {
		return aiprotocol.JobMessage{}, gerror.New("去衣缺少图片(input_url / source)")
	}

	return aiprotocol.JobMessage{
		SchemaVersion: aiprotocol.SchemaVersion,
		JobID:         in.TaskNo,
		Biz:           aiprotocol.BizUndress,
		MediaType:     aiprotocol.MediaPhoto,
		Source:        source,
		Output: aiprotocol.OutputRef{
			Bucket: bucket,
			Prefix: site + "/ai/out/" + in.TaskNo + "/",
		},
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func jobBucketSite(ctx context.Context) (bucket, site string) {
	bucket = strings.TrimSpace(g.Cfg().MustGet(ctx, "faceswap.bucket", aiprotocol.DefaultBucket).String())
	if bucket == "" {
		bucket = aiprotocol.DefaultBucket
	}
	site = strings.TrimSpace(g.Cfg().MustGet(ctx, "site.code", "my").String())
	if site == "" {
		site = "my"
	}
	return bucket, strings.Trim(site, "/")
}

func pickRef(params map[string]any, keys []string, fallbackURL, bucket string) (aiprotocol.ObjectRef, error) {
	for _, k := range keys {
		if params == nil {
			break
		}
		if v, ok := params[k]; ok {
			if ref, ok := aiprotocol.ObjectRefFromAny(v, bucket); ok && ref.OK() {
				return ref, nil
			}
		}
	}
	if fallbackURL != "" {
		if ref, ok := aiprotocol.ParseObjectRef(fallbackURL, bucket); ok && ref.OK() {
			return ref, nil
		}
	}
	return aiprotocol.ObjectRef{}, gerror.New("object ref missing")
}
