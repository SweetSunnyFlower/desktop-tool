package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	// "github.com/StackExchange/wmi"
	"github.com/nguyenthenguyen/docx"
	// "github.com/StackExchange/wmi"
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
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type Prompt struct {
	ID      int    `json:"id"`
	Prompt  string `json:"prompt"`
	History string `json:"history"`
}

// 解析Prompt文件
func (a *App) ParsePromptFile(file_path string) map[string]interface{} {

	// mock upload image
	prompts := []Prompt{
		{
			ID:     1,
			Prompt: "a little boy",
		},
		{
			ID:     2,
			Prompt: "a little girl",
		},
	}

	return map[string]interface{}{"code": 0, "data": prompts, "message": response}
}

func (a *App) UploadImage(input string) map[string]interface{} {
	imagePath := input

	// 判断是否为文件夹
	if !IsDir(imagePath) {
		return map[string]interface{}{"code": 1, "data": []string{}, "message": "文件名替换成对出现"}
	}

	if response == "" {
		response = "上传成功"
	}

	// mock upload image
	bosImages := []BosImage{
		{
			ID:  1,
			URL: "https://www.baidu.com/img/bd_logo1.png",
		},
		{
			ID:  2,
			URL: "https://www.baidu.com/img/bd_logo1.png",
		},
	}

	return map[string]interface{}{"code": 0, "data": bosImages, "message": response}
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
