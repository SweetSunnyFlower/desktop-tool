package config

import "tools/pkg/config"

func init() {
	config.Add("sd", func() map[string]interface{} {
		return map[string]interface{}{
			"host":   config.Env("SD_HOST", "http://127.0.0.1:8090"),
			"width":  config.Env("SD_WIDTH", 1024),
			"height": config.Env("SD_HEIGHT", 1024),
		}
	})
}
