// Package config 站点配置信息
package config

import "tools/pkg/config"

func init() {
	config.Add("app", func() map[string]interface{} {
		return map[string]interface{}{

			// 应用名称
			"name": config.Env("APP_NAME", "Image To Text"),

			"width": config.Env("APP_WIDTH", 1024),

			"height": config.Env("APP_HEIGHT", 768),

			// 当前环境，用以区分多环境，一般为 local, stage, production, test
			"env": config.Env("APP_ENV", "production"),

			// 是否进入调试模式
			"debug": config.Env("APP_DEBUG", false),
		}
	})
}
