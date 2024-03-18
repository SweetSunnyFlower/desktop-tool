package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"tools/pkg/bos"
	"tools/pkg/vis"

	// "github.com/StackExchange/wmi"
	"github.com/golang-module/carbon/v2"
	"github.com/nguyenthenguyen/docx"

	// "github.com/StackExchange/wmi"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

type Storage struct {
	Name       string
	FileSystem string
	Total      uint64
	Free       uint64
}

type storageInfo struct {
	Name       string
	Size       uint64
	FreeSpace  uint64
	FileSystem string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string, age string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name) + fmt.Sprintf("Hello %s, It's show time!", age)
}

func (a *App) GetDirs(parent string) map[string]interface{} {
	var paths []map[string]interface{}
	var err error
	if parent == "" {
		paths, err = getDirAndFileByPath(fsRootDir())
	} else {
		paths, err = getDirAndFileByPath(parent)
	}
	if err != nil {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
	}
	return map[string]interface{}{"code": 0, "data": paths, "message": "success"}
}

func getDirAndFileByPath(path string) ([]map[string]interface{}, error) {

	var paths []map[string]interface{}

	if runtime.GOOS == "windows" && path == "/" {
		systems := getStorageInfo()
		for _, system := range systems {
			paths = append(paths, map[string]interface{}{
				"path":        system.Name + "/",
				"isDir":       true,
				"name":        system.Name,
				"hasChildren": true,
			})
		}
	} else {
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {

			rootPath := path + file.Name()

			if file.IsDir() {
				children, err := os.ReadDir(rootPath)
				if err == nil {
					paths = append(paths, map[string]interface{}{
						"path":        rootPath + "/",
						"isDir":       file.IsDir(),
						"name":        file.Name(),
						"hasChildren": len(children) != 0,
					})
				}
			} else {
				paths = append(paths, map[string]interface{}{
					"path":        rootPath + "/",
					"isDir":       file.IsDir(),
					"name":        file.Name(),
					"hasChildren": false,
				})
			}
		}
	}

	return paths, nil
}

func getStorageInfo() []Storage {
	var storageinfo []storageInfo
	var loaclStorages []Storage
	// err := wmi.Query("Select * from Win32_LogicalDisk", &storageinfo)
	// if err != nil {
	// 	panic(err)
	// }

	for _, storage := range storageinfo {
		info := Storage{
			Name:       storage.Name,
			FileSystem: storage.FileSystem,
			Total:      storage.Size,
			Free:       storage.FreeSpace,
		}
		loaclStorages = append(loaclStorages, info)
	}

	return loaclStorages
}

func fsRootDir() string {
	// if runtime.GOOS == "windows" {
	// 	return os.Getenv("SystemDrive")
	// }
	return "/"
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

var response string = ""

type BosImage struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Prompt struct {
	ID      string `json:"id"`
	Prompt  string `json:"prompt"`
	History string `json:"history"`
}

type ImageToText struct {
	ID         string     `json:"id"`
	URL        string     `json:"url"`
	Prompt     string     `json:"prompt"`
	History    string     `json:"history"`
	Result     string     `json:"result"`
	HistoryMsg [][]string `json:"history_msg"`
	OcrRet     string     `json:"ocr_ret"`
	FaceRet    string     `json:"face_ret"`
}

type ImageToTextDownload struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	Prompt     string `json:"prompt"`
	History    string `json:"history"`
	Result     string `json:"result"`
	HistoryMsg string `json:"history_msg"`
	OcrRet     string `json:"ocr_ret"`
	FaceRet    string `json:"face_ret"`
}

// 解析Prompt文件
func (a *App) ParsePromptFile(path string) map[string]interface{} {

	wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
		"type": "info",
		"msg":  "开始解析文件:" + path,
	})

	// Excel读取文件内容，返回返回
	file, err := os.Open(path)
	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "打开文件:" + path + "失败",
			"error": err.Error(),
		})
		return map[string]interface{}{"code": 1, "data": map[string]interface{}{}, "message": "打开文件失败"}
	}
	defer file.Close()

	// 创建CSV reader
	reader := csv.NewReader(file)

	// 读取CSV文件中的内容
	data, err := reader.ReadAll()

	wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
		"type": "info",
		"msg":  "开始读取文件:" + path + "内容",
	})
	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "读取文件:" + path + "内容失败",
			"error": err.Error(),
		})
		return map[string]interface{}{"code": 1, "data": map[string]interface{}{}, "message": "打开文件失败"}
	}

	prompts := make([]Prompt, 0)
	for i, data := range data {
		if i == 0 {
			continue
		}
		prompts = append(prompts, Prompt{
			ID:      strings.TrimSpace(data[0]),
			Prompt:  strings.TrimSpace(data[1]),
			History: strings.TrimSpace(data[2]),
		})

	}

	wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
		"type":  "info",
		"msg":   "解析文件:" + path + "完成",
		"error": err.Error(),
		"data":  prompts,
	})

	return map[string]interface{}{"code": 0, "data": prompts, "message": "解析成功"}
}

func (a *App) UploadImage(input string) map[string]interface{} {

	imagePath := input

	if imagePath == "" {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": "请选择文件夹"}
	}

	// 判断是否为文件夹
	if !IsDir(imagePath) {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": "请选择文件夹"}
	}

	// 读取当前文件夹下源文件
	files, err := os.ReadDir(imagePath)
	if err != nil {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
	}

	count := len(files)

	isFinish := make(chan bool, count)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		go a.uploadBos(imagePath, file, isFinish)
	}

	<-isFinish

	return map[string]interface{}{"code": 0, "data": map[string]interface{}{}, "message": "上传完成"}
}

func (a *App) uploadBos(imagePath string, file fs.DirEntry, isFinish chan bool) {
	client := bos.NewBos()

	filename := file.Name()
	// 获取文件后缀名
	ext := filepath.Ext(filename)

	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		wailsruntime.EventsEmit(a.ctx, "logEvent", filename+":不支持的图片格式,直接跳过")
		isFinish <- true
		return
	}

	filenameWithoutExt := strings.TrimSuffix(filename, ext)

	// 读取文件内容
	content, err := os.Open(imagePath + "/" + filename)

	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", filename+":打开文件失败,error:"+err.Error())
		isFinish <- true
		return
	}

	byteImg, err := io.ReadAll(content)

	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", filename+":读取文件内容失败,error:"+err.Error())
		isFinish <- true
		return
	}

	_, out, err := client.Upload("", string(byteImg), "image2text")

	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", filename+":上传bos失败,error:"+err.Error())
		isFinish <- true
		return
	}

	wailsruntime.EventsEmit(a.ctx, "uploadImageEvent", BosImage{
		ID:  filenameWithoutExt,
		URL: out,
	})
	wailsruntime.EventsEmit(a.ctx, "logEvent", filename+":上传成功,url:"+out)

	isFinish <- true
}

func (a *App) Replace(input string, output string, args []map[string]string, fileName []string) map[string]interface{} {

	if len(fileName) != 0 && len(fileName) != 2 {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": "文件名替换成对出现"}
	}
	replaceFileName := false

	if len(fileName) == 2 {
		replaceFileName = true
	}

	if input == "" {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": "请选择输入路径"}
	}

	if output == "" {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": "请选择输出路径"}
	}

	if len(args) == 0 {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": "请输入替换字段"}
	}

	for _, arg := range args {
		if arg["key"] == "" {
			return map[string]interface{}{"code": 1, "data": []string{}, "message": "需要替换的字段不能为空"}
		}
	}

	// 读取当前文件夹下源文件
	files, err := os.ReadDir(input)

	if err != nil {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
	}

	_, err = os.Stat(output)

	if err != nil && os.IsNotExist(err) {
		os.Mkdir(output, 0777)
	}

	for _, file := range files {

		if file.IsDir() {
			a.Replace(input+file.Name()+"/", output+file.Name()+"/", args, fileName)
		} else {
			r, err := docx.ReadDocxFile(input + file.Name())
			name := file.Name()
			if err != nil {
				response = response + name + err.Error()
				continue
			}
			docx1 := r.Editable()
			for _, fields := range args {
				docx1.Replace(fields["key"], fields["value"], -1)
				if replaceFileName {
					name = strings.Replace(name, fileName[0], fileName[1], -1)
				}
				docx1.WriteToFile(output + name)
			}

			r.Close()
		}
	}

	if response == "" {
		response = "替换成功"
	}

	return map[string]interface{}{"code": 0, "data": []string{}, "message": response}
}

func (a *App) OpenFile(t string) map[string]interface{} {
	file, err := wailsruntime.OpenFileDialog(a.ctx, wailsruntime.OpenDialogOptions{})

	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "打开文件:" + t + "失败",
			"error": err.Error(),
		})
		return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
	}

	if file == "" {
		return map[string]interface{}{"code": 2, "data": []string{}, "message": ""}
	}

	if t == "prompt" {
		return a.ParsePromptFile(file)
	}

	return map[string]interface{}{"code": 0, "data": []string{}, "message": file}
}
func (a *App) OpenFolder(t string, data string) map[string]interface{} {
	folder, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{})

	wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
		"type": "info",
		"msg":  "打开文件夹:" + folder,
	})

	if folder == "" {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type": "info",
			"msg":  "取消打开文件夹",
		})
		return map[string]interface{}{"code": 2, "data": []string{}, "message": ""}
	}

	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "打开文件夹:" + folder + "失败",
			"error": err.Error(),
		})
		return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
	}

	wailsruntime.EventsEmit(a.ctx, "handlingEvent", true)

	response := map[string]interface{}{"code": 0, "data": []string{}, "message": folder}

	if t == "images" {
		response = a.UploadImage(folder)
	}

	// 下载模版
	if t == "download-template" {
		response = a.DownloadCsvTemplate(folder)
	}

	// 下载文生图数据
	if t == "download-iamge2text" {
		response = a.DownloadImage2Text(folder, data)
	}

	wailsruntime.EventsEmit(a.ctx, "handlingEvent", false)

	return response
}

func (a *App) Image2Text(data string) {

	wailsruntime.EventsEmit(a.ctx, "handlingEvent", true)

	defer wailsruntime.EventsEmit(a.ctx, "handlingEvent", false)

	var imageToTexts []*ImageToText

	err := json.Unmarshal([]byte(data), &imageToTexts)

	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "图生文数据解析失败",
			"error": err.Error(),
			"data":  data,
		})
	}

	visInstance := vis.NewVis()

	count := len(imageToTexts)

	for i, imageToText := range imageToTexts {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type": "info",
			"msg":  "图生文请求开始",
			"data": imageToText,
		})

		result, err := visInstance.Image2Text(imageToText.URL)
		if err != nil {
			wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
				"type":  "error",
				"msg":   "图生文请求失败",
				"error": err.Error(),
				"data":  imageToText,
			})
			continue
		}
		imageToTexts[i].FaceRet = result.FaceRet
		imageToTexts[i].OcrRet = result.OcrRet
		imageToTexts[i].HistoryMsg = result.HistoryMsg
		imageToTexts[i].Result = result.Result

		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type": "info",
			"msg":  "图生文请求成功",
			"data": result,
		})

		wailsruntime.EventsEmit(a.ctx, "image2TextEvent", imageToTexts)

		if i != count-1 {
			time.Sleep(2 * time.Second)
		}
	}
}

func (a *App) DownloadCsvTemplate(folder string) map[string]interface{} {
	outputFilePath := folder + "/template_" + carbon.Now().String() + ".csv"

	// 写入CSV文件
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "下载Prompt模版失败:下载路径" + folder,
			"error": err.Error(),
		})
		return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
	}
	defer outputFile.Close()

	// 创建CSV writer
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	records := [][]string{{"ID", "Prompt", "History"}}
	// 写入记录
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
		}
	}

	return map[string]interface{}{"code": 0, "data": []string{}, "message": "下载成功"}
}

func (a *App) DownloadImage2Text(folder string, data string) map[string]interface{} {

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
		return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
	}

	outputFilePath := folder + "/image2text_" + carbon.Now().String() + ".csv"

	// 写入CSV文件
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "文生图数据下载，文件:" + outputFilePath + "创建失败",
			"error": err.Error(),
		})
		return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
	}
	defer outputFile.Close()

	// 创建CSV writer
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	var records [][]string = [][]string{{"id", "url", "prompt", "history", "result", "face_ret", "ocr_ret", "history_msg"}}

	for _, imageToText := range imageToTexts {

		var record []string
		record = append(record, imageToText.ID)
		record = append(record, imageToText.URL)
		record = append(record, imageToText.Prompt)
		record = append(record, imageToText.History)
		record = append(record, imageToText.Result)
		record = append(record, imageToText.FaceRet)
		record = append(record, imageToText.OcrRet)
		record = append(record, imageToText.HistoryMsg)
		records = append(records, record)
	}

	// 写入记录
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			return map[string]interface{}{"code": 1, "data": []string{}, "message": err.Error()}
		}
	}

	return map[string]interface{}{"code": 0, "data": []string{}, "message": "下载成功"}
}
