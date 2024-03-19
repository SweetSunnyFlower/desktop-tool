package config

import "tools/pkg/config"

func init() {
	config.Add("user", func() map[string]interface{} {
		return map[string]interface{}{
			"email":    config.Env("EMAIL", ""),
			"password": config.Env("PASSWORD", ""),
			"uid":      config.Env("UID", 0),
			"cuid":     config.Env("CUID", ""),
		}
	})
}
