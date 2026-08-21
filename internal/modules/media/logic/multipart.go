package logic

import (
	"context"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/media/service"
	"github.com/JarvanDante/my_service/internal/shared/storage"
)

const (
	defaultPartSize = 8 << 20 // 8MiB
	minPartSize     = 5 << 20 // S3: 非末片 ≥ 5MiB
	maxPartCount    = 10000
)

func (s *sMedia) MultipartInit(ctx context.Context, in service.MultipartInitInput) (*service.MultipartInitDTO, error) {
	purpose := strings.ToLower(strings.TrimSpace(in.Purpose))
	if purpose == "" {
		purpose = "video"
	}
	if _, ok := purposeDefaults[purpose]; !ok {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "不支持的用途: %s", purpose)
	}
	if in.Size <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "文件大小无效")
	}
	filename := strings.TrimSpace(in.Filename)
	if filename == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "文件名必填")
	}

	if in.Resume && in.OperatorId > 0 {
		exist, err := s.repo.MultipartFindActive(ctx, in.OperatorId, filename, in.Size)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "查询分片会话失败")
		}
		if exist != nil && exist.Id > 0 {
			return multipartInitFromSession(exist), nil
		}
	}

	contentType := detectContentType(filename, in.ContentType)
	if err := s.validateMime(ctx, purpose, contentType); err != nil {
		return nil, err
	}
	if err := s.validateSize(ctx, purpose, in.Size); err != nil {
		return nil, err
	}

	partSize := in.PartSize
	if partSize <= 0 {
		partSize = g.Cfg().MustGet(ctx, "upload.multipart.partSize", defaultPartSize).Int64()
	}
	if partSize < minPartSize {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "分片大小不能小于 %d 字节", minPartSize)
	}
	partCount := int(math.Ceil(float64(in.Size) / float64(partSize)))
	if partCount < 1 {
		partCount = 1
	}
	if partCount > maxPartCount {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"分片数 %d 超过上限 %d, 请增大 part_size", partCount, maxPartCount,
		)
	}

	client, err := storage.Get(ctx)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}
	objectKey := buildObjectKey(purpose, filename)
	minioUploadID, err := client.NewMultipart(ctx, objectKey, contentType)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "发起分片上传失败")
	}

	uploadID := uuid.NewString()
	_, err = s.repo.MultipartCreate(ctx, &entity.MediaMultipart{
		UploadId:      uploadID,
		MinioUploadId: minioUploadID,
		Bucket:        client.Bucket(),
		ObjectKey:     objectKey,
		Purpose:       purpose,
		Filename:      filename,
		ContentType:   contentType,
		Size:          in.Size,
		PartSize:      partSize,
		PartCount:     partCount,
		Status:        entity.MultipartStatusUploading,
		CreatedBy:     in.OperatorId,
	})
	if err != nil {
		_ = client.AbortMultipart(ctx, objectKey, minioUploadID)
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "保存分片会话失败")
	}

	return &service.MultipartInitDTO{
		UploadId: uploadID, ObjectKey: objectKey, Bucket: client.Bucket(),
		Purpose: purpose, ContentType: contentType, Size: in.Size,
		PartSize: partSize, PartCount: partCount,
	}, nil
}

func multipartInitFromSession(sess *entity.MediaMultipart) *service.MultipartInitDTO {
	return &service.MultipartInitDTO{
		UploadId: sess.UploadId, ObjectKey: sess.ObjectKey, Bucket: sess.Bucket,
		Purpose: sess.Purpose, ContentType: sess.ContentType, Size: sess.Size,
		PartSize: sess.PartSize, PartCount: sess.PartCount,
	}
}

func (s *sMedia) MultipartUploadPart(ctx context.Context, in service.MultipartUploadPartInput) (*service.MultipartPartDTO, error) {
	sess, err := s.loadUploadingSession(ctx, in.UploadId, in.OperatorId)
	if err != nil {
		return nil, err
	}
	if in.PartNumber < 1 || in.PartNumber > sess.PartCount {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "分片号 %d 超出范围 1~%d", in.PartNumber, sess.PartCount)
	}
	if in.File == nil || in.File.FileHeader == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "分片文件必填")
	}
	size := in.File.Size
	if size <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "分片大小无效")
	}
	if in.PartNumber < sess.PartCount && size != sess.PartSize {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "非末片大小应为 %d 字节", sess.PartSize)
	}
	if in.PartNumber == sess.PartCount && size > sess.PartSize {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "末片不能超过 %d 字节", sess.PartSize)
	}

	client, err := storage.Get(ctx)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}
	f, err := in.File.Open()
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "打开分片失败")
	}
	defer f.Close()

	etag, err := client.PutObjectPart(ctx, sess.ObjectKey, sess.MinioUploadId, in.PartNumber, f, size)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "上传分片失败")
	}
	return &service.MultipartPartDTO{PartNumber: in.PartNumber, Etag: etag, Size: size}, nil
}

func (s *sMedia) MultipartPresign(ctx context.Context, in service.MultipartPresignInput) ([]service.MultipartPresignItemDTO, error) {
	sess, err := s.loadUploadingSession(ctx, in.UploadId, in.OperatorId)
	if err != nil {
		return nil, err
	}
	if len(in.PartNumbers) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "part_numbers 必填")
	}
	client, err := storage.Get(ctx)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}

	expires := time.Duration(g.Cfg().MustGet(ctx, "upload.multipart.presignExpire", 7200).Int64()) * time.Second
	if expires <= 0 {
		expires = 2 * time.Hour
	}
	expiresIn := int64(expires.Seconds())

	out := make([]service.MultipartPresignItemDTO, 0, len(in.PartNumbers))
	seen := map[int]struct{}{}
	for _, n := range in.PartNumbers {
		if n < 1 || n > sess.PartCount {
			return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "分片号 %d 超出范围 1~%d", n, sess.PartCount)
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		u, err := client.PresignUploadPart(ctx, sess.ObjectKey, sess.MinioUploadId, n, expires)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeInternalError, err, "预签名失败")
		}
		out = append(out, service.MultipartPresignItemDTO{
			PartNumber: n, Url: u, Method: "PUT", ExpiresIn: expiresIn,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].PartNumber < out[j].PartNumber })
	return out, nil
}

func (s *sMedia) MultipartParts(ctx context.Context, uploadId string, operatorId int64) (*service.MultipartPartsDTO, error) {
	sess, err := s.loadSession(ctx, uploadId, operatorId)
	if err != nil {
		return nil, err
	}
	dto := &service.MultipartPartsDTO{
		UploadId: sess.UploadId, Status: sess.Status, PartCount: sess.PartCount,
		List: []service.MultipartPartDTO{},
	}
	if sess.Status != entity.MultipartStatusUploading {
		return dto, nil
	}
	client, err := storage.Get(ctx)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}
	parts, err := client.ListParts(ctx, sess.ObjectKey, sess.MinioUploadId)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "查询已传分片失败")
	}
	for _, p := range parts {
		dto.List = append(dto.List, service.MultipartPartDTO{
			PartNumber: p.PartNumber, Etag: p.ETag, Size: p.Size,
		})
	}
	return dto, nil
}

func (s *sMedia) MultipartComplete(ctx context.Context, in service.MultipartCompleteInput) (*service.UploadDTO, error) {
	sess, err := s.loadUploadingSession(ctx, in.UploadId, in.OperatorId)
	if err != nil {
		return nil, err
	}
	client, err := storage.Get(ctx)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}

	var parts []storage.PartInfo
	if len(in.Parts) > 0 {
		for _, p := range in.Parts {
			if p.PartNumber < 1 || strings.TrimSpace(p.Etag) == "" {
				return nil, gerror.NewCode(gcode.CodeInvalidParameter, "parts 含无效项")
			}
			parts = append(parts, storage.PartInfo{PartNumber: p.PartNumber, ETag: p.Etag})
		}
	} else {
		listed, err := client.ListParts(ctx, sess.ObjectKey, sess.MinioUploadId)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeInternalError, err, "查询已传分片失败")
		}
		parts = listed
	}
	if len(parts) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "没有可合并的分片")
	}
	sort.Slice(parts, func(i, j int) bool { return parts[i].PartNumber < parts[j].PartNumber })

	url, err := client.CompleteMultipart(ctx, sess.ObjectKey, sess.MinioUploadId, parts)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "合并分片失败")
	}

	mediaId, err := s.repo.Create(ctx, &entity.MediaObject{
		Bucket: client.Bucket(), ObjectKey: sess.ObjectKey, Purpose: sess.Purpose,
		ContentType: sess.ContentType, Size: sess.Size, CreatedBy: in.OperatorId,
	})
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "保存媒体记录失败")
	}
	if err = s.repo.MultipartUpdateStatus(ctx, sess.UploadId, entity.MultipartStatusCompleted, mediaId); err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "更新分片会话失败")
	}

	return &service.UploadDTO{
		Id: mediaId, Url: url, ObjectKey: sess.ObjectKey, Bucket: client.Bucket(),
		Purpose: sess.Purpose, ContentType: sess.ContentType, Size: sess.Size,
	}, nil
}

func (s *sMedia) MultipartAbort(ctx context.Context, in service.MultipartAbortInput) error {
	sess, err := s.loadUploadingSession(ctx, in.UploadId, in.OperatorId)
	if err != nil {
		return err
	}
	client, err := storage.Get(ctx)
	if err != nil {
		return gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}
	if err = client.AbortMultipart(ctx, sess.ObjectKey, sess.MinioUploadId); err != nil {
		return gerror.WrapCode(gcode.CodeInternalError, err, "取消分片上传失败")
	}
	if err = s.repo.MultipartUpdateStatus(ctx, sess.UploadId, entity.MultipartStatusAborted, 0); err != nil {
		return gerror.WrapCode(gcode.CodeDbOperationError, err, "更新分片会话失败")
	}
	return nil
}

func (s *sMedia) loadSession(ctx context.Context, uploadId string, operatorId int64) (*entity.MediaMultipart, error) {
	uploadId = strings.TrimSpace(uploadId)
	if uploadId == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "upload_id 必填")
	}
	sess, err := s.repo.MultipartFindByUploadId(ctx, uploadId)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "查询分片会话失败")
	}
	if sess == nil || sess.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "分片会话不存在")
	}
	if operatorId > 0 && sess.CreatedBy > 0 && sess.CreatedBy != operatorId {
		// 超管/同站运维可放宽; P0 先限制本人会话
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "无权操作该上传会话")
	}
	return sess, nil
}

func (s *sMedia) loadUploadingSession(ctx context.Context, uploadId string, operatorId int64) (*entity.MediaMultipart, error) {
	sess, err := s.loadSession(ctx, uploadId, operatorId)
	if err != nil {
		return nil, err
	}
	if sess.Status != entity.MultipartStatusUploading {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "分片会话已结束")
	}
	return sess, nil
}
