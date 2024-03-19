package main

import (
	"os"
	"tools/pkg/config"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) Config() map[string]interface{} {

	// 当前目录下有没有 license.txt 文件
	_, err := os.Stat("license.txt")

	if err != nil {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
	}

	cfg := config.GetConfig()

	all := map[string]map[string]interface{}{}

	for name, fn := range cfg {
		all[name] = fn()
	}

	wailsruntime.EventsEmit(a.ctx, "queryConfigEvent", all)

	return map[string]interface{}{"code": 0, "data": []string{}, "message": err.Error()}
}

func (a *App) Register(value string) {
	// 创建文件
	file, err := os.Create("license.txt")
	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "创建文件失败",
			"error": err.Error(),
		})

		return
	}
	defer file.Close()

	// 将内容写入文件
	_, err = file.Write([]byte(value))
	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "写入文件失败",
			"error": err.Error(),
		})
		return
	}
}
