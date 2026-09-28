package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"N_m3u8DL-RE-WEB-UI/internal/model"
	"N_m3u8DL-RE-WEB-UI/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateProfileRequest 创建方案请求
type CreateProfileRequest struct {
	Name               string `json:"name" binding:"required"`
	Domain             string `json:"domain"`
	ThreadCount        int    `json:"thread_count"`
	RetryCount         int    `json:"retry_count"`
	Headers            string `json:"headers"`
	BaseURL            string `json:"base_url"`
	DelAfterDone       *bool  `json:"del_after_done"`
	BinaryMerge        *bool  `json:"binary_merge"`
	AutoSelect         *bool  `json:"auto_select"`
	SkipSegmentsCheck  *bool  `json:"skip_segments_check"`
	ConcurrentDownload *bool  `json:"concurrent_download"`
	Key                string `json:"key"`
	DecryptionEngine   string `json:"decryption_engine"`
	CustomArgs         string `json:"custom_args"`
	CustomProxy        string `json:"custom_proxy"`
}

// ListProfiles 获取方案列表
func ListProfiles(c *gin.Context) {
	var profiles []model.DownloadProfile
	domain := c.Query("domain")

	query := model.GetDB().Model(&model.DownloadProfile{})
	if domain != "" {
		query = query.Where("domain = ?", domain)
	}
	if err := query.Order("created_at DESC").Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取方案列表失败: " + err.Error()})
		return
	}

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
	skipSegmentsCheck := false
	if req.SkipSegmentsCheck != nil {
		skipSegmentsCheck = *req.SkipSegmentsCheck
	}
	concurrentDownload := false
	if req.ConcurrentDownload != nil {
		concurrentDownload = *req.ConcurrentDownload
	}

	profile := &model.DownloadProfile{
		Name:               req.Name,
		Domain:             req.Domain,
		ThreadCount:        req.ThreadCount,
		RetryCount:         req.RetryCount,
		Headers:            req.Headers,
		BaseURL:            req.BaseURL,
		DelAfterDone:       delAfterDone,
		BinaryMerge:        binaryMerge,
		AutoSelect:         autoSelect,
		SkipSegmentsCheck:  skipSegmentsCheck,
		ConcurrentDownload: concurrentDownload,
		DecryptionEngine:   req.DecryptionEngine,
		CustomArgs:         req.CustomArgs,
		CustomProxy:        req.CustomProxy,
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
	if req.SkipSegmentsCheck != nil {
		profile.SkipSegmentsCheck = *req.SkipSegmentsCheck
	}
	if req.ConcurrentDownload != nil {
		profile.ConcurrentDownload = *req.ConcurrentDownload
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
	if err := model.GetDB().Where("domain = ?", domain).Order("updated_at DESC").Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询方案失败: " + err.Error()})
		return
	}

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

	// 默认名称只取可注册主域名（hls.ted.com -> ted.com），避免把子域名也当名称
	baseName := service.RegistrableDomain(domain)
	if baseName == "" {
		baseName = "未命名方案"
	}

	// 待保存的方案内容：域名 + 全部下载配置。
	// 有意不含名称（按需求不比较名称）、也不含解密密钥
	// —— 密钥逐视频不同，存进方案只会在加载时注入错误的密钥。
	// 域名必须计入内容，否则不同站点但配置相同会被误判成「已存在」。
	target := &model.DownloadProfile{
		Domain:             domain,
		ThreadCount:        task.ThreadCount,
		RetryCount:         task.RetryCount,
		Headers:            task.Headers,
		BaseURL:            task.BaseURL,
		DelAfterDone:       task.DelAfterDone,
		BinaryMerge:        task.BinaryMerge,
		AutoSelect:         task.AutoSelect,
		SkipSegmentsCheck:  task.SkipSegmentsCheck,
		ConcurrentDownload: task.ConcurrentDownload,
		DecryptionEngine:   task.DecryptionEngine,
		CustomArgs:         task.CustomArgs,
		CustomProxy:        task.CustomProxy,
	}

	var existing []model.DownloadProfile
	if err := model.GetDB().Find(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询方案失败: " + err.Error()})
		return
	}

	// 1) 已存在内容完全一致的方案 → 视为重复，不新增
	targetKey := profileContentKey(target)
	for i := range existing {
		if profileContentKey(&existing[i]) == targetKey {
			c.JSON(http.StatusOK, gin.H{
				"profile": existing[i],
				"action":  "exists",
			})
			return
		}
	}

	// 2) 内容不一致但名称已被占用 → 自动往后加序号：ted.com -> ted.com2 -> ted.com3
	taken := make(map[string]bool, len(existing))
	for i := range existing {
		taken[existing[i].Name] = true
	}
	finalName := uniqueProfileName(taken, baseName)

	target.Name = finalName
	if err := model.GetDB().Create(target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存方案失败: " + err.Error()})
		return
	}

	action := "created"
	if finalName != baseName {
		action = "renamed"
	}
	c.JSON(http.StatusOK, gin.H{
		"profile":        target,
		"action":         action,
		"requested_name": baseName,
	})
}

// profileContentKey 生成用于判断「方案内容是否相同」的指纹。
// 包含域名与全部影响下载行为的配置；
// 有意排除 Name（按需求不比较名称）与 Key（解密密钥逐视频不同，不随方案保存）。
func profileContentKey(p *model.DownloadProfile) string {
	return strings.Join([]string{
		p.Domain,
		strconv.Itoa(p.ThreadCount),
		strconv.Itoa(p.RetryCount),
		p.Headers,
		p.BaseURL,
		strconv.FormatBool(p.DelAfterDone),
		strconv.FormatBool(p.BinaryMerge),
		strconv.FormatBool(p.AutoSelect),
		strconv.FormatBool(p.SkipSegmentsCheck),
		strconv.FormatBool(p.ConcurrentDownload),
		p.DecryptionEngine,
		p.CustomArgs,
		p.CustomProxy,
		// \x1f（单元分隔符）不会出现在这些字段里，避免拼接歧义
	}, "\x1f")
}

// uniqueProfileName 在名称已被占用时依次尝试 name2、name3……
func uniqueProfileName(taken map[string]bool, base string) string {
	if base == "" {
		base = "未命名方案"
	}
	if !taken[base] {
		return base
	}
	for i := 2; i <= 999; i++ {
		candidate := fmt.Sprintf("%s%d", base, i)
		if !taken[candidate] {
			return candidate
		}
	}
	return fmt.Sprintf("%s_%d", base, time.Now().Unix())
}
