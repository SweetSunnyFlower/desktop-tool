// Package config 站点配置信息
package config

import "tools/pkg/config"

func init() {
	config.Add("cdn", func() map[string]interface{} {
		return map[string]interface{}{
			"endpoint": config.Env("CDN_ENDPOINT", ""),
		}
	})
}
