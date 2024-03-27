package main

import (
	"tools/pkg/stablediffusion"
	"tools/pkg/translate"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) Translate(word string, from string, to string) string {
	wailsruntime.EventsEmit(a.ctx, "handlingEvent", true)

	translateService := translate.NewBaiduTranslate(a.ctx)
	result, err := translateService.Translate(word, from, to)
	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type": "error",
			"msg":  "翻译失败:" + err.Error(),
		})
		return ""
	}

	sdService := stablediffusion.NewStableDiffusion(a.ctx)

	imageStr, err := sdService.TextToImage(result, "")

	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type": "error",
			"msg":  "sd失败:" + err.Error(),
		})
		return ""
	}

	wailsruntime.EventsEmit(a.ctx, "handlingEvent", false)

	return imageStr
}
