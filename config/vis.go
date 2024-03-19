package config

import "tools/pkg/config"

func init() {
	config.Add("vis", func() map[string]interface{} {
		return map[string]interface{}{
			"limit":          config.Env("VIS_REQUEST_LIMIT", 2),
			"host":           config.Env("VIS_HOST", "http://127.0.0.1:8090"),
			"auth_key":       config.Env("VIS_AUTH_KEY", ""),
			"image2text_url": config.Env("VIS_IMAGE2TEXT_URL", ""),
			"business_name":  config.Env("VIS_BUSINESS_NAME", ""),
			"feature_name":   config.Env("VIS_FEATURE_NAME", ""),
		}
	})
}
