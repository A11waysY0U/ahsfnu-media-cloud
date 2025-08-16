package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ahsfnu-media-cloud/internal/config"
	"ahsfnu-media-cloud/internal/models"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type UploadService struct {
	uploadPath string
}

func NewUploadService() *UploadService {
	return &UploadService{
		uploadPath: config.AppConfig.Upload.UploadPath,
	}
}

// UploadFile 上传单个文件
func (s *UploadService) UploadFile(file *multipart.FileHeader, userID uint, workflowID *uint) (*models.Material, error) {
	// 验证文件类型
	if !s.isAllowedFileType(file.Filename) {
		return nil, fmt.Errorf("不支持的文件类型: %s", filepath.Ext(file.Filename))
	}

	// 验证文件大小
	if file.Size > config.AppConfig.Upload.MaxFileSize {
		return nil, fmt.Errorf("文件大小超过限制: %d bytes", config.AppConfig.Upload.MaxFileSize)
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 创建年月日目录结构
	now := time.Now()
	datePath := fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day())
	fullPath := filepath.Join(s.uploadPath, datePath)

	// 确保目录存在
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	// 完整的文件路径
	filePath := filepath.Join(fullPath, filename)
	relativePath := filepath.Join(datePath, filename)

	// 保存文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 检测文件类型和提取元数据
	fileType, mimeType, width, height, duration, err := s.extractFileMetadata(filePath)
	if err != nil {
		return nil, fmt.Errorf("提取文件元数据失败: %v", err)
	}

	// 创建素材记录
	material := &models.Material{
		Filename:         filename,
		OriginalFilename: file.Filename,
		FilePath:         relativePath,
		FileSize:         file.Size,
		FileType:         fileType,
		MimeType:         mimeType,
		Width:            width,
		Height:           height,
		Duration:         duration,
		UploadedBy:       userID,
		WorkflowID:       workflowID,
	}

	// 如果是图片或视频，生成缩略图
	if fileType == "image" {
		thumbName := "thumb_" + filename
		thumbPath := filepath.Join(fullPath, thumbName)
		thumbRelPath := filepath.Join(datePath, thumbName)
		if err := generateThumbnail(filePath, thumbPath); err == nil {
			material.ThumbnailPath = thumbRelPath
		}
	} else if fileType == "video" {
		thumbName := "thumb_" + strings.TrimSuffix(filename, ext) + ".jpg"
		thumbPath := filepath.Join(fullPath, thumbName)
		thumbRelPath := filepath.Join(datePath, thumbName)
		if err := generateVideoThumbnail(filePath, thumbPath); err == nil {
			material.ThumbnailPath = thumbRelPath
		}
	}

	return material, nil
}

// isAllowedFileType 检查文件类型是否允许
func (s *UploadService) isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range config.AppConfig.Upload.AllowedTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// extractFileMetadata 提取文件元数据
func (s *UploadService) extractFileMetadata(filePath string) (fileType, mimeType string, width, height, duration *int, err error) {
	// 读取文件头部来检测类型
	file, err := os.Open(filePath)
	if err != nil {
		return "unknown", "application/octet-stream", nil, nil, nil, nil
	}
	defer file.Close()

	// 读取前512字节用于类型检测
	head := make([]byte, 512)
	n, err := file.Read(head)
	if err != nil || n == 0 {
		return "unknown", "application/octet-stream", nil, nil, nil, nil
	}

	// 检测文件类型
	kind, err := filetype.Match(head)
	if err != nil || kind == filetype.Unknown {
		// 如果检测失败，尝试根据文件扩展名判断
		ext := strings.ToLower(filepath.Ext(filePath))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
			return "image", "image/" + strings.TrimPrefix(ext, "."), nil, nil, nil, nil
		case ".mp4", ".avi", ".mov", ".wmv", ".flv":
			return "video", "video/" + strings.TrimPrefix(ext, "."), nil, nil, nil, nil
		default:
			return "unknown", "application/octet-stream", nil, nil, nil, nil
		}
	}

	mimeType = kind.MIME.Value

	// 根据MIME类型判断文件类型
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		fileType = "image"
		// 这里可以添加图片尺寸提取逻辑
		// width, height = extractImageDimensions(filePath)
	case strings.HasPrefix(mimeType, "video/"):
		fileType = "video"
		// 这里可以添加视频时长提取逻辑
		// duration = extractVideoDuration(filePath)
	default:
		fileType = "unknown"
	}

	return fileType, mimeType, width, height, duration, nil
}

// 生成缩略图
func generateThumbnail(srcPath, dstPath string) error {
	img, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}
	thumb := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)
	return imaging.Save(thumb, dstPath)
}

// 生成视频缩略图（依赖系统已安装 ffmpeg）
func generateVideoThumbnail(srcPath, dstPath string) error {
	// 优先尝试直接输出到文件，显式指定图像格式与编码器
	primaryErr := ffmpeg.Input(srcPath, ffmpeg.KwArgs{"ss": "0.5"}).
		Output(dstPath, ffmpeg.KwArgs{"vframes": 1, "vf": "thumbnail,scale=200:-1", "f": "image2", "vcodec": "mjpeg"}).
		OverWriteOutput().
		Run()
	if primaryErr == nil {
		return nil
	}

	// 失败则走管道方式到内存，再用 imaging 保存，增强兼容性
	buf := &bytes.Buffer{}
	pipeErr := ffmpeg.Input(srcPath, ffmpeg.KwArgs{"ss": "0.5"}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "vf": "thumbnail,scale=200:-1", "f": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()
	if pipeErr != nil {
		return pipeErr
	}

	img, decodeErr := imaging.Decode(buf)
	if decodeErr != nil {
		return decodeErr
	}
	return imaging.Save(img, dstPath)
}

// GetFileURL 获取文件访问URL
func (s *UploadService) GetFileURL(material *models.Material) string {
	// 将 Windows 路径分隔符替换为 URL 友好的正斜杠
	clean := strings.ReplaceAll(material.FilePath, "\\", "/")
	return fmt.Sprintf("/uploads/%s", clean)
}

// GetThumbnailURL 获取缩略图访问URL
func (s *UploadService) GetThumbnailURL(material *models.Material) string {
	if material.ThumbnailPath == "" {
		return ""
	}
	clean := strings.ReplaceAll(material.ThumbnailPath, "\\", "/")
	return fmt.Sprintf("/uploads/%s", clean)
}

// DeleteFile 删除文件
func (s *UploadService) DeleteFile(material *models.Material) error {
	fullPath := filepath.Join(s.uploadPath, material.FilePath)
	_ = os.Remove(fullPath)

	// 删除缩略图（如果有）
	if material.ThumbnailPath != "" {
		thumbFullPath := filepath.Join(s.uploadPath, material.ThumbnailPath)
		_ = os.Remove(thumbFullPath)
	}
	return nil
}
