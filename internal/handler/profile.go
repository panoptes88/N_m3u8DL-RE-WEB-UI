package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"N_m3u8DL-RE-WEB-UI/internal/model"

	"github.com/gin-gonic/gin"
)

// CreateProfileRequest 创建方案请求
type CreateProfileRequest struct {
	Name             string `json:"name" binding:"required"`
	Domain           string `json:"domain"`
	ThreadCount      int    `json:"thread_count"`
	RetryCount       int    `json:"retry_count"`
	Headers          string `json:"headers"`
	BaseURL          string `json:"base_url"`
	DelAfterDone     *bool  `json:"del_after_done"`
	BinaryMerge      *bool  `json:"binary_merge"`
	AutoSelect       *bool  `json:"auto_select"`
	Key              string `json:"key"`
	DecryptionEngine string `json:"decryption_engine"`
	CustomArgs       string `json:"custom_args"`
	CustomProxy      string `json:"custom_proxy"`
}

// ListProfiles 获取方案列表
func ListProfiles(c *gin.Context) {
	var profiles []model.DownloadProfile
	domain := c.Query("domain")

	query := model.GetDB().Model(&model.DownloadProfile{})
	if domain != "" {
		query = query.Where("domain = ?", domain)
	}
	query = query.Order("created_at DESC").Find(&profiles)

	c.JSON(http.StatusOK, profiles)
}

// CreateProfile 创建方案
func CreateProfile(c *gin.Context) {
	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 设置默认值
	if req.ThreadCount <= 0 {
		req.ThreadCount = 32
	}
	if req.RetryCount <= 0 {
		req.RetryCount = 15
	}
	if req.DecryptionEngine == "" {
		req.DecryptionEngine = "MP4DECRYPT"
	}

	// 如果没有指定域名，尝试从BaseURL或Headers中提取
	if req.Domain == "" && req.BaseURL != "" {
		if u, err := url.Parse(req.BaseURL); err == nil {
			req.Domain = u.Hostname()
		}
	}

	// 处理布尔指针，默认值
	delAfterDone := true
	if req.DelAfterDone != nil {
		delAfterDone = *req.DelAfterDone
	}
	binaryMerge := false
	if req.BinaryMerge != nil {
		binaryMerge = *req.BinaryMerge
	}
	autoSelect := false
	if req.AutoSelect != nil {
		autoSelect = *req.AutoSelect
	}

	profile := &model.DownloadProfile{
		Name:             req.Name,
		Domain:           req.Domain,
		ThreadCount:      req.ThreadCount,
		RetryCount:       req.RetryCount,
		Headers:          req.Headers,
		BaseURL:          req.BaseURL,
		DelAfterDone:     delAfterDone,
		BinaryMerge:      binaryMerge,
		AutoSelect:       autoSelect,
		DecryptionEngine: req.DecryptionEngine,
		CustomArgs:       req.CustomArgs,
		CustomProxy:      req.CustomProxy,
	}

	if err := model.GetDB().Create(profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建方案失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetProfile 获取方案详情
func GetProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的方案ID"})
		return
	}

	var profile model.DownloadProfile
	if err := model.GetDB().First(&profile, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "方案不存在"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile 更新方案
func UpdateProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的方案ID"})
		return
	}

	var profile model.DownloadProfile
	if err := model.GetDB().First(&profile, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "方案不存在"})
		return
	}

	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 更新字段
	if req.Name != "" {
		profile.Name = req.Name
	}
	if req.Domain != "" {
		profile.Domain = req.Domain
	}
	if req.ThreadCount > 0 {
		profile.ThreadCount = req.ThreadCount
	}
	if req.RetryCount > 0 {
		profile.RetryCount = req.RetryCount
	}
	if req.Headers != "" {
		profile.Headers = req.Headers
	}
	if req.BaseURL != "" {
		profile.BaseURL = req.BaseURL
	}
	if req.DelAfterDone != nil {
		profile.DelAfterDone = *req.DelAfterDone
	}
	if req.BinaryMerge != nil {
		profile.BinaryMerge = *req.BinaryMerge
	}
	if req.AutoSelect != nil {
		profile.AutoSelect = *req.AutoSelect
	}
	if req.DecryptionEngine != "" {
		profile.DecryptionEngine = req.DecryptionEngine
	}
	if req.CustomArgs != "" {
		profile.CustomArgs = req.CustomArgs
	}
	if req.CustomProxy != "" {
		profile.CustomProxy = req.CustomProxy
	}

	if err := model.GetDB().Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新方案失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// DeleteProfile 删除方案
func DeleteProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的方案ID"})
		return
	}

	if err := model.GetDB().Delete(&model.DownloadProfile{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除方案失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetProfileByDomain 根据域名获取匹配的方案
func GetProfileByDomain(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名参数不能为空"})
		return
	}

	var profiles []model.DownloadProfile
	model.GetDB().Where("domain = ?", domain).Order("updated_at DESC").Find(&profiles)

	if len(profiles) == 0 {
		c.JSON(http.StatusOK, nil)
		return
	}

	// 返回最新的方案
	c.JSON(http.StatusOK, profiles[0])
}

// SaveTaskAsProfile 将任务保存为方案
func SaveTaskAsProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var task model.Task
	if err := model.GetDB().First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	// 提取域名
	domain := ""
	if task.BaseURL != "" {
		if u, err := url.Parse(task.BaseURL); err == nil {
			domain = u.Hostname()
		}
	}
	if domain == "" && task.URL != "" {
		if u, err := url.Parse(task.URL); err == nil {
			domain = u.Hostname()
		}
	}

	// 生成方案名称
	name := domain
	if name == "" {
		name = "未命名方案"
	}

	profile := &model.DownloadProfile{
		Name:             name,
		Domain:           domain,
		ThreadCount:      task.ThreadCount,
		RetryCount:       task.RetryCount,
		Headers:          task.Headers,
		BaseURL:          task.BaseURL,
		DelAfterDone:     task.DelAfterDone,
		BinaryMerge:      task.BinaryMerge,
		AutoSelect:       task.AutoSelect,
		DecryptionEngine: task.DecryptionEngine,
		CustomArgs:       task.CustomArgs,
		CustomProxy:      task.CustomProxy,
	}

	if err := model.GetDB().Create(profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存方案失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
