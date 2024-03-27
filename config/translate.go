package config

import "tools/pkg/config"

func init() {
	config.Add("translate", func() map[string]interface{} {
		return map[string]interface{}{
			"host":       config.Env("TRANSLATE_HOST", ""),
			"app_id":     config.Env("TRANSLATE_APP_ID", ""),
			"app_secret": config.Env("TRANSLATE_APP_SECRET", ""),
		}
	})
}
