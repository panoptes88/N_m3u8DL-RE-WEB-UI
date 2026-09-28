package service

import (
	"log"

	"N_m3u8DL-RE-WEB-UI/internal/model"

	"golang.org/x/net/publicsuffix"
)

// RegistrableDomain 返回可注册主域名（eTLD+1）。
// 使用公共后缀列表，能正确处理 com.cn / co.uk 这类多级后缀
// （简单取「后两段」会把 example.com.cn 算成 com.cn）。
// 无法解析时原样返回。
func RegistrableDomain(host string) string {
	if host == "" {
		return ""
	}
	if d, err := publicsuffix.EffectiveTLDPlusOne(host); err == nil {
		return d
	}
	return host
}

// MigrateProfileNames 把历史数据里「方案名 = 完整域名」的名称规范化为主域名，
// 例如 hls.ted.com -> ted.com。
// 只处理看起来是自动生成的名称（name 恰好等于 domain），
// 用户手动重命名过的方案不会被改动。可重复执行。
func MigrateProfileNames() {
	var profiles []model.DownloadProfile
	if err := model.GetDB().
		Where("name = domain").
		Where("domain <> ''").
		Find(&profiles).Error; err != nil {
		log.Printf("方案名称规范化失败: %v", err)
		return
	}

	changed := 0
	for i := range profiles {
		main := RegistrableDomain(profiles[i].Domain)
		if main == "" || main == profiles[i].Name {
			continue
		}
		if err := model.GetDB().Model(&profiles[i]).Update("name", main).Error; err != nil {
			log.Printf("更新方案 %d 名称失败: %v", profiles[i].ID, err)
			continue
		}
		changed++
	}

	if changed > 0 {
		log.Printf("已将 %d 个方案的名称由完整域名规范化为主域名", changed)
	}
}
