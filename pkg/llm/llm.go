package llm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
	"tools/pkg/config"

	"github.com/spf13/cast"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Completion struct {
	LogID        int64         // LogID
	UID          uint64        // UID 用户pass id
	CUID         string        // CUID  用户的cuid uid和cuid 必传1个
	UA           string        // UA 用户的ua  adapt.UA()
	COOKIE       string        //COOKIE 用户的cookie adapt.CookieStr()
	IP           string        // IP 用户的IP adapt.ClientIP()
	Query        string        // Query 问题内容
	SessionID    int64         // SessionID
	AskerAnswers []AskerAnswer // AskerAnswers 上下文
	EnableVilg   bool          // EnableVilg 是否启用文新一格,画图
	System       string        // System 人设
	SystemID     string        // SystemID
	Model        string        // Model
}

// AskerAnswer 问题答案组，涉及到传递上下文的时候都可以用这个结构，也可以自由增加字段
type AskerAnswer struct {
	AskID    int64
	AnswerID int64
	Ask      string
	Answer   string
}

type LLM struct {
	ctx    context.Context
	Client http.Client
}

type UserInfo struct {
	UID  uint64 `json:"uid"`
	CUID string `json:"cuid"`
}

type Headers struct {
	Ua       string `json:"ua,omitempty"`
	Cookie   string `json:"cookie,omitempty"`
	ClientIP string `json:"client_ip,omitempty"`
}

type LlmRequest struct {
	Model    string         `json:"model,omitempty"` // 请求模型名
	Messages []*ChatMessage `json:"messages,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role,omitempty"`    // 角色: system, user, assistant, or function
	Content string `json:"content,omitempty"` // 对话内容，不能为空
}

type Option struct {
	EnableTts      bool   `json:"enable_tts,omitempty"`      // 开启tts
	EnableVilg     bool   `json:"enable_vilg,omitempty"`     // 开启文心一格
	EnableCitation bool   `json:"enable_citation,omitempty"` // 是否开启上角标返回
	EnableTrace    bool   `json:"enable_trace,omitempty"`    // 是否开启溯源信息返回
	SessionID      string `json:"session_id,omitempty"`      // ⽤于标记多轮对话信息不⼤于36位，不允许是"0"，[0-9|a-z|A-Z|-]，开头abdcdefgh, abcdefghi-, abc--dkdk--dkd 这种都⾮法
	SystemID       string `json:"system_id,omitempty"`
}

// Req 架构网关的请求
type Body struct {
	LogID      int64       `json:"logid"`
	User       *UserInfo   `json:"user"`
	From       string      `json:"from"`
	Token      string      `json:"token"`
	ModelName  string      `json:"model_name"`
	Headers    *Headers    `json:"headers"`
	LlmRequest *LlmRequest `json:"llm_request"`
	Stream     bool        `json:"stream"`
	Option     *Option     `json:"option"`
}

type Response struct {
	ErrNo       string       `json:"err_no,omitempty"`  // 下游错误码
	ErrMsg      string       `json:"err_msg,omitempty"` // 错误信息
	LlmResponse *LlmResponse `json:"llm_response,omitempty"`
}

type CompletionResponse struct {
	ChatID  string `json:"chatId"`
	Content string `json:"content"`
}

type LlmResponse struct {
	ErrNo            int64  `json:"err_no"`
	ErrMsg           string `json:"err_msg,omitempty"`
	ChatID           string `json:"chat_id,omitempty"` // 本轮对话id，可以用这个id调dispatch赞踩接口
	Result           string `json:"result,omitempty"`  // 生成内容，普通对话返回文本，tts模式返回语音链接，多个语音链接以空格分开
	ContentType      string `json:"content_type,omitempty"`
	NeedClearHistory bool   `json:"need_clear_history,omitempty"` // 安全封禁，对话不可以继续，值为true表示用户输入存在安全风险，建议关闭当前会话，清理历史会话信息
	BanRound         int32  `json:"ban_round,omitempty"`          // 当need_clear_history为true时，次字段会告知第几轮对话有敏感信息，如果是当前问题，ban_round = -1
	IsSafe           int32  `json:"is_safe,omitempty"`
}

var instance *LLM

var once sync.Once

func NewLLM(ctx context.Context) *LLM {
	once.Do(func() {
		instance = &LLM{
			ctx: ctx,
		}
	})
	return instance
}

func (l *LLM) Completions(model, replaced string) (*CompletionResponse, error) {

	// 生成随机int
	source := rand.NewSource(time.Now().UnixNano())
	logId := rand.New(source).Int()

	completion := &Completion{
		Model: model,
		Query: replaced,
		UA:    "1440_3007_android_1.10_560",
		UID:   cast.ToUint64(config.GetInt64("user.uid")),
		CUID:  config.GetString("user.cuid"),
		LogID: cast.ToInt64(logId),
	}

	strBody := l.buildBody(completion)

	url := config.GetString("llm.host") + "/baidu.mlarch.modelgateway.ModelService/Generate"

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(strBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := l.Client.Do(req)

	if err != nil {
		wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "大模型请求错误",
			"body":  strBody,
			"error": err.Error(),
		})
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
			"type":     "error",
			"msg":      "大模型请求读取body错误",
			"payload":  strBody,
			"response": string(data),
			"error":    err.Error(),
		})
		return nil, err
	}

	wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
		"type":    "info",
		"msg":     "大模型请求成功",
		"payload": strBody,
		"data":    string(data),
	})

	var res *Response

	if err = json.Unmarshal(data, res); err != nil {
		return nil, err
	}

	if res.ErrNo != "SUCC" {
		wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "大模型请响应成功，业务状态码错误(" + res.ErrNo + ")",
			"data":  string(data),
			"error": errors.New("大模型请求失败"),
		})

		return nil, errors.New("大模型请求失败")
	}

	return &CompletionResponse{
		ChatID:  res.LlmResponse.ChatID,
		Content: res.LlmResponse.Result,
	}, nil
}

func (l *LLM) buildBody(completion *Completion) string {

	messages := make([]*ChatMessage, 0, len(completion.AskerAnswers))

	for i := 0; i < len(completion.AskerAnswers); i++ {
		// 问题
		query := &ChatMessage{
			Role:    "user",
			Content: completion.AskerAnswers[i].Ask,
		}
		messages = append(messages, query)

		// 答案
		answer := &ChatMessage{
			Role:    "assistant",
			Content: completion.AskerAnswers[i].Answer,
		}

		messages = append(messages, answer)
	}

	userQuery := &ChatMessage{
		Role:    "user",
		Content: completion.Query,
	}
	messages = append(messages, userQuery)

	if completion.System != "" {
		// system
		system := &ChatMessage{
			Role:    "system",
			Content: completion.System,
		}
		messages = append(messages, system)
	}

	params := &Body{
		LogID: completion.LogID,
		Headers: &Headers{
			Ua:       completion.UA,
			Cookie:   completion.COOKIE,
			ClientIP: completion.IP,
		},
		User: &UserInfo{
			UID:  completion.UID,
			CUID: completion.CUID,
		},
		From:      config.GetString("llm.from"),
		Token:     config.GetString("llm.token"),
		ModelName: strings.ToUpper(completion.Model),
		LlmRequest: &LlmRequest{
			Model:    "",
			Messages: messages,
		},
		Stream: false,
		Option: &Option{
			EnableTts:      true,
			EnableVilg:     completion.EnableVilg,
			EnableCitation: false,
			EnableTrace:    false,
		},
	}

	// 构建请求体
	body, _ := json.Marshal(params)

	return string(body)
}
