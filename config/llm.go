package config

import "tools/pkg/config"

func init() {
	config.Add("llm", func() map[string]interface{} {
		return map[string]interface{}{
			"host":  config.Env("LLM_HOST", "http://127.0.0.1:8090"),
			"from":  config.Env("LLM_FROM", ""),
			"token": config.Env("LLM_TOKEN", ""),
		}
	})
}
