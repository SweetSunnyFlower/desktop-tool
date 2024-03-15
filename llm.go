package main

type VisPrompt struct {
	ID         int    `json:"id"`
	Result     string `json:"result"`
	HistoryMsg string `json:"history_msg"`
	OcrRet     string `json:"ocr_ret"`
	FaceRet    string `json:"face_ret"`
}

// 解析Prompt文件
func (a *App) ParseVisFile(file_path string) map[string]interface{} {
	// mock upload image
	prompts := []VisPrompt{
		{
			ID:         1,
			Result:     "a little boy",
			HistoryMsg: "HistoryMsg",
			OcrRet:     "OcrRet",
			FaceRet:    "FaceRet",
		},
		{
			ID:         2,
			Result:     "a little girl",
			HistoryMsg: "HistoryMsg",
			OcrRet:     "OcrRet",
			FaceRet:    "FaceRet",
		},
	}

	return map[string]interface{}{"code": 0, "data": prompts, "message": "解析成功"}
}
