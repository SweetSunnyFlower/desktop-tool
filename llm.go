package main

import (
	"encoding/json"
	"strings"
	"time"
	"tools/pkg/llm"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) LLM(model, template, data string) {

	var imageToTexts []*ImageToTextDownload

	wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
		"type": "info",
		"msg":  "开始解析图生文数据",
	})

	err := json.Unmarshal([]byte(data), &imageToTexts)

	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "文生图数据解析失败",
			"data":  data,
			"error": err.Error(),
		})
	}

	// 调用llm
	llmInstance := llm.NewLLM(a.ctx)

	for _, imageToText := range imageToTexts {

		replacements := map[string]string{
			"@id":          imageToText.ID,
			"@url":         imageToText.URL,
			"@prompt":      imageToText.Prompt,
			"@history":     imageToText.History,
			"@result":      imageToText.Result,
			"@history_msg": imageToText.HistoryMsg,
			"@ocr_ret":     imageToText.OcrRet,
			"@face_ret":    imageToText.FaceRet,
		}

		replaced := template
		for key, value := range replacements {
			replaced = strings.Replace(replaced, key, value, -1)
		}

		replay, err := llmInstance.Completions(model, replaced)

		if err != nil {
			wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
				"type":  "error",
				"msg":   "大模型请求错误",
				"error": err.Error(),
			})
		}

		wailsruntime.EventsEmit(a.ctx, "llmEvent", map[string]interface{}{
			"replay": replay,
		})

		time.Sleep(time.Second * 1)
	}
}
