package handler

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"N_m3u8DL-RE-WEB-UI/internal/config"
	"N_m3u8DL-RE-WEB-UI/internal/service"

	"github.com/gin-gonic/gin"
)

// 视频 MIME 类型映射（包级别变量，避免重复创建）
var videoMimeTypes = map[string]string{
	".mp4":  "video/mp4",
	".webm": "video/webm",
	".mkv":  "video/x-matroska",
	".avi":  "video/x-msvideo",
	".mov":  "video/quicktime",
	".flv":  "video/x-flv",
	".wmv":  "video/x-ms-wmv",
	".m4v":  "video/x-m4v",
	".3gp":  "video/3gpp",
	".mpg":  "video/mpeg",
	".mpeg": "video/mpeg",
}

func ListFiles(c *gin.Context) {
	cfg := config.Load()
	files, err := service.ListFiles(cfg.DownloadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

func DownloadFile(c *gin.Context) {
	filename := c.Query("name")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	// 防止路径遍历
	filename = filepath.Base(filename)

	cfg := config.Load()
	filePath := filepath.Join(cfg.DownloadDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", getVideoMimeType(filename))
	// 使用 mime.FormatMediaType 生成 RFC 5987 兼容的头：
	// 同时给出 ASCII 回退名与 filename*=UTF-8''…，
	// 此前直接用 url.PathEscape 拼 filename=，中文文件名在部分浏览器会乱码
	c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{
		"filename": filename,
	}))
	c.File(filePath)
}

func DeleteFile(c *gin.Context) {
	filename := c.Param("name")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	// 防止路径遍历
	filename = filepath.Base(filename)

	cfg := config.Load()
	filePath := filepath.Join(cfg.DownloadDir, filename)

	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

type BatchDeleteFilesRequest struct {
	Names []string `json:"names" binding:"required,min=1"`
}

// BatchDeleteFiles 批量删除文件。
// 使用 DELETE /files + JSON body，避免与 /files/:name 产生路由冲突。
// 前端此前是 N 次串行请求，50 个文件就要 50 次往返。
func BatchDeleteFiles(c *gin.Context) {
	var req BatchDeleteFilesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	cfg := config.Load()
	deleted := make([]string, 0, len(req.Names))
	failed := make([]string, 0)

	for _, name := range req.Names {
		if name == "" {
			continue
		}
		// 防止路径遍历
		safeName := filepath.Base(name)
		if err := os.Remove(filepath.Join(cfg.DownloadDir, safeName)); err != nil {
			failed = append(failed, name)
			continue
		}
		deleted = append(deleted, name)
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted": deleted,
		"failed":  failed,
	})
}

// Health 健康检查：无需鉴权，供容器 healthcheck 使用。
// 不能复用 /api/user，因为它需要登录、未登录返回 401，会让健康检查永远失败。
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func getVideoMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if mimeType, ok := videoMimeTypes[ext]; ok {
		return mimeType
	}
	return "application/octet-stream"
}
