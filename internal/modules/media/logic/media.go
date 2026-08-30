package logic

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/google/uuid"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/media/domain"
	"github.com/JarvanDante/my_service/internal/modules/media/service"
	"github.com/JarvanDante/my_service/internal/shared/aesbnc"
	"github.com/JarvanDante/my_service/internal/shared/imgopt"
	"github.com/JarvanDante/my_service/internal/shared/paas"
	"github.com/JarvanDante/my_service/internal/shared/storage"
)

type sMedia struct{ repo domain.Repository }

func New(repo domain.Repository) service.IMedia { return &sMedia{repo: repo} }

var imageMimes = []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}
var avatarMimes = []string{"image/jpeg", "image/jpg", "image/png", "image/webp"}
var videoMimes = []string{
	"video/mp4", "video/quicktime", "video/x-matroska", "video/webm",
	"video/3gpp", "video/3gpp2", "video/x-m4v",
}

var purposeDefaults = map[string]struct {
	mime    []string
	maxSize int64 // KB
}{
	"image":      {imageMimes, 5120},
	"cover":      {imageMimes, 5120},
	"avatar":     {avatarMimes, 2048},
	"ad":         {imageMimes, 5120},
	"post":       {imageMimes, 5120},
	"video":      {videoMimes, 614400},
	"post_video": {videoMimes, 614400},
}

func (s *sMedia) Upload(ctx context.Context, in service.UploadInput) (*service.UploadDTO, error) {
	// File 非 nil 但内嵌 FileHeader 为 nil 时, 直接读 Filename 会 panic
	if in.File == nil || in.File.FileHeader == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "文件必填")
	}
	purpose := strings.ToLower(strings.TrimSpace(in.Purpose))
	if purpose == "" {
		purpose = "image"
	}
	if _, ok := purposeDefaults[purpose]; !ok {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "不支持的用途: %s", purpose)
	}

	filename := in.File.Filename
	size := in.File.Size
	if size <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "文件大小无效")
	}

	headerType := ""
	if in.File.Header != nil {
		headerType = in.File.Header.Get("Content-Type")
	}
	contentType := detectContentType(filename, headerType)
	if err := s.validateMime(ctx, purpose, contentType); err != nil {
		return nil, err
	}
	if err := s.validateSize(ctx, purpose, size); err != nil {
		return nil, err
	}

	f, err := in.File.Open()
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "打开上传文件失败")
	}
	defer f.Close()

	var body io.Reader = f
	if aesbnc.ShouldEncryptPurpose(purpose) {
		raw, err := io.ReadAll(f)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeInternalError, err, "读取上传文件失败")
		}
		if packed, ok := imgopt.Compress(raw, purpose); ok {
			raw = packed
			filename = imgopt.JpegName(filename)
		}
		enc, err := aesbnc.Encrypt(raw)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeInternalError, err, "图片加密失败")
		}
		body = bytes.NewReader(enc)
		filename = aesbnc.ToBncKey(filename)
		contentType = "application/octet-stream"
		size = int64(len(enc))
	}

	created, confirmed, err := putViaUnifiedStorage(ctx, filename, purpose, contentType, size, body)
	if err != nil {
		return nil, err
	}
	url := confirmed.PublicUrl
	if url == "" {
		url = created.PublicUrl
	}
	objectKey := created.Key
	if objectKey == "" {
		objectKey = created.Id
	}
	bucket := created.Bucket
	if bucket == "" {
		bucket = "my-storage"
	}

	id, err := s.repo.Create(ctx, &entity.MediaObject{
		Bucket:      bucket,
		ObjectKey:   objectKey,
		Purpose:     purpose,
		ContentType: contentType,
		Size:        confirmed.SizeBytes,
		CreatedBy:   in.OperatorId,
	})
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "保存媒体记录失败")
	}
	if confirmed.SizeBytes > 0 {
		size = confirmed.SizeBytes
	}

	return &service.UploadDTO{
		Id: id, Url: url, ObjectKey: objectKey, Bucket: bucket,
		Purpose: purpose, ContentType: contentType, Size: size,
	}, nil
}

func storageBiz(purpose string) string {
	switch strings.ToLower(strings.TrimSpace(purpose)) {
	case "video":
		return "video"
	case "post_video":
		return "post_video"
	case "avatar":
		return "avatar"
	case "cover":
		return "cover"
	case "ad":
		return "ad"
	case "post":
		return "post"
	default:
		return "image"
	}
}

func putViaUnifiedStorage(ctx context.Context, filename, purpose, contentType string, size int64, body io.Reader) (*paas.StorageCreateOut, *paas.StorageConfirmOut, error) {
	created, err := paas.CreateStorageObject(ctx, paas.StorageCreateIn{
		Filename: filename, Biz: storageBiz(purpose),
		ContentType: contentType, SizeBytes: size, Remark: "site-" + purpose,
	})
	if err != nil {
		return nil, nil, err
	}
	if err = paas.PutUploadURL(ctx, created.UploadUrl, body, size); err != nil {
		return nil, nil, gerror.WrapCode(gcode.CodeInternalError, err, "上传失败")
	}
	confirmed, err := paas.ConfirmStorageObject(ctx, created.Id)
	if err != nil {
		return nil, nil, err
	}
	return created, confirmed, nil
}

func (s *sMedia) InitStorageUpload(ctx context.Context, in service.StorageInitInput) (*service.StorageInitDTO, error) {
	purpose := strings.ToLower(strings.TrimSpace(in.Purpose))
	if purpose == "" {
		purpose = "image"
	}
	if _, ok := purposeDefaults[purpose]; !ok {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "不支持的用途: %s", purpose)
	}
	filename := strings.TrimSpace(in.Filename)
	if filename == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "文件名必填")
	}
	if in.Size <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "文件大小无效")
	}
	contentType := detectContentType(filename, in.ContentType)
	if err := s.validateMime(ctx, purpose, contentType); err != nil {
		return nil, err
	}
	if err := s.validateSize(ctx, purpose, in.Size); err != nil {
		return nil, err
	}
	out, err := paas.CreateStorageObject(ctx, paas.StorageCreateIn{
		Filename: filename, Biz: storageBiz(purpose),
		ContentType: contentType, SizeBytes: in.Size, Remark: "h5-" + purpose,
	})
	if err != nil {
		return nil, err
	}
	return &service.StorageInitDTO{
		Id: out.Id, UploadUrl: out.UploadUrl, Method: out.Method, Bucket: out.Bucket,
		ObjectKey: out.Key, ExpireSec: out.ExpireSec, PublicUrl: out.PublicUrl,
		ContentType: out.ContentType,
	}, nil
}

func (s *sMedia) ConfirmStorageUpload(ctx context.Context, id string) (*service.StorageConfirmDTO, error) {
	out, err := paas.ConfirmStorageObject(ctx, id)
	if err != nil {
		return nil, err
	}
	return &service.StorageConfirmDTO{
		Id: out.Id, Url: out.PublicUrl, Size: out.SizeBytes,
		Bucket: "my-storage", ObjectKey: out.Id,
	}, nil
}

func (s *sMedia) ResolveObjectRef(ctx context.Context, rawURL, objectKey string) (string, string, string, error) {
	client, err := storage.Get(ctx)
	if err != nil {
		return "", "", "", gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}
	bucket, key := storage.ParseURL(rawURL, "")
	if key == "" {
		okey := strings.TrimSpace(objectKey)
		if okey != "" && !strings.Contains(okey, "..") && strings.Contains(okey, "/") {
			key = okey
		}
	}
	if bucket == "" {
		if strings.Contains(rawURL, "/my-storage/") || strings.Contains(objectKey, "my-storage/") {
			bucket = "my-storage"
		} else {
			bucket = client.Bucket()
		}
	}
	if key == "" {
		return "", "", "", gerror.NewCode(gcode.CodeInvalidParameter, "无法识别的图片地址")
	}
	name := key
	if rawURL != "" {
		name = rawURL
	}
	return bucket, key, name, nil
}

func (s *sMedia) ReadObject(ctx context.Context, rawURL, objectKey string) ([]byte, string, error) {
	client, err := storage.Get(ctx)
	if err != nil {
		return nil, "", gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}
	bucket, key, name, err := s.ResolveObjectRef(ctx, rawURL, objectKey)
	if err != nil {
		return nil, "", err
	}
	data, err := client.GetIn(ctx, bucket, key)
	if err != nil {
		return nil, "", gerror.WrapCode(gcode.CodeNotFound, err, "对象不存在")
	}
	return data, name, nil
}

func (s *sMedia) validateMime(ctx context.Context, purpose, contentType string) error {
	allowed := g.Cfg().MustGet(ctx, fmt.Sprintf("upload.%s.mime", purpose)).Strings()
	if len(allowed) == 0 {
		allowed = purposeDefaults[purpose].mime
	}
	ct := strings.ToLower(contentType)
	for _, a := range allowed {
		if ct == strings.ToLower(a) {
			return nil
		}
	}
	return gerror.NewCodef(gcode.CodeInvalidParameter, "不支持的文件类型: %s", contentType)
}

func (s *sMedia) validateSize(ctx context.Context, purpose string, size int64) error {
	maxKB := g.Cfg().MustGet(ctx, fmt.Sprintf("upload.%s.maxSize", purpose)).Int64()
	if maxKB <= 0 {
		maxKB = purposeDefaults[purpose].maxSize
	}
	if size > maxKB*1024 {
		if purpose == "video" || purpose == "post_video" {
			return gerror.NewCode(gcode.CodeInvalidParameter, "视频不能超过 600M，请压缩后再上传")
		}
		return gerror.NewCodef(gcode.CodeInvalidParameter, "文件大小不能超过 %dKB", maxKB)
	}
	return nil
}

func detectContentType(filename, headerType string) string {
	ct := strings.TrimSpace(headerType)
	if ct != "" && ct != "application/octet-stream" {
		return ct
	}
	if ext := filepath.Ext(filename); ext != "" {
		if byExt := mime.TypeByExtension(ext); byExt != "" {
			return byExt
		}
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			return "image/jpeg"
		case ".png":
			return "image/png"
		case ".gif":
			return "image/gif"
		case ".webp":
			return "image/webp"
		case ".mp4":
			return "video/mp4"
		case ".mov":
			return "video/quicktime"
		case ".mkv":
			return "video/x-matroska"
		case ".webm":
			return "video/webm"
		case ".3gp", ".3gpp":
			return "video/3gpp"
		case ".m4v":
			return "video/x-m4v"
		}
	}
	return "application/octet-stream"
}

func buildObjectKey(purpose, filename string) string {
	site := genv.Get("SITE_CODE", "my").String()
	if site == "" {
		site = "my"
	}
	site = strings.ToLower(site)
	ext := strings.ToLower(filepath.Ext(filename))
	day := time.Now().Format("2006/01/02")
	return fmt.Sprintf("%s/%s/%s/%s%s", site, purpose, day, uuid.NewString(), ext)
}
