package main

import (
	"tools/pkg/config"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) Config() map[string]interface{} {

	cfg := config.GetConfig()

	all := map[string]map[string]interface{}{}

	for name, fn := range cfg {
		all[name] = fn()
	}

	wailsruntime.EventsEmit(a.ctx, "queryConfigEvent", all)

	return map[string]interface{}{"code": 0, "data": []string{}, "message": "ok"}
}

func (a *App) SetUser(email, uid, cuid string) map[string]interface{} {

	config.Add("user", func() map[string]interface{} {
		return map[string]interface{}{
			"email": email,
			"uid":   uid,
			"cuid":  cuid,
		}
	})

	config.LoadConfig()

	return map[string]interface{}{"code": 0, "data": []string{}, "message": "ok"}
}
