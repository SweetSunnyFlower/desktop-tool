package config

import "tools/pkg/config"

func init() {
	config.Add("user", func() map[string]interface{} {
		return map[string]interface{}{
			"uid":  config.Env("UID", 0),
			"cuid": config.Env("CUID", ""),
		}
	})
}
