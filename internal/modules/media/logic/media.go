package logic

import (
	"context"
	"fmt"
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
	"github.com/JarvanDante/my_service/internal/shared/storage"
)

type sMedia struct{ repo domain.Repository }

func New(repo domain.Repository) service.IMedia { return &sMedia{repo: repo} }

var purposeDefaults = map[string]struct {
	mime    []string
	maxSize int64 // KB
}{
	"image":  {[]string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}, 5120},
	"cover":  {[]string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}, 5120},
	"avatar": {[]string{"image/jpeg", "image/jpg", "image/png", "image/webp"}, 2048},
	"video":  {[]string{"video/mp4", "video/quicktime", "video/x-matroska", "video/webm"}, 2097152},
}

func (s *sMedia) Upload(ctx context.Context, in service.UploadInput) (*service.UploadDTO, error) {
	if in.File == nil {
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

	contentType := detectContentType(filename, in.File.Header.Get("Content-Type"))
	if err := s.validateMime(ctx, purpose, contentType); err != nil {
		return nil, err
	}
	if err := s.validateSize(ctx, purpose, size); err != nil {
		return nil, err
	}

	objectKey := buildObjectKey(purpose, filename)
	client, err := storage.Get(ctx)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}

	f, err := in.File.Open()
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "打开上传文件失败")
	}
	defer f.Close()

	url, err := client.Put(ctx, objectKey, f, size, contentType)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "上传失败")
	}

	id, err := s.repo.Create(ctx, &entity.MediaObject{
		Bucket:      client.Bucket(),
		ObjectKey:   objectKey,
		Purpose:     purpose,
		ContentType: contentType,
		Size:        size,
		CreatedBy:   in.OperatorId,
	})
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "保存媒体记录失败")
	}

	return &service.UploadDTO{
		Id: id, Url: url, ObjectKey: objectKey, Bucket: client.Bucket(),
		Purpose: purpose, ContentType: contentType, Size: size,
	}, nil
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
