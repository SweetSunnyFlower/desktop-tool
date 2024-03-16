// Package config 站点配置信息
package config

import "tools/pkg/config"

func init() {
	config.Add("bos", func() map[string]interface{} {
		return map[string]interface{}{
			"ak":               config.Env("BOS_AK", ""),
			"sk":               config.Env("BOS_SK", ""),
			"endpoint":         config.Env("BOS_ENDPOINT", ""),
			"redirectDisabled": config.Env("BOS_REDIRECT_DISABLED", false),
			"defaultBucket":    config.Env("BOS_DEFAULT_BUCKET", ""),
			"host":             config.Env("BOS_HOST", ""),
			"inner_host":       config.Env("BOS_INNER_HOST", ""),
		}
	})
}
