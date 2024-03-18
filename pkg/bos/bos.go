package bos

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"sync"
	"tools/pkg/config"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cast"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Bos struct {
	BosClient     *bos.Client
	CdnClient     *cdn.Client
	DefaultBucket string
	ctx           context.Context
}

type Bucket struct {
	Bucket string
	Host   string
	InHost string
	Env    string
	BosIDC string
}

var BosInstance *Bos

var once sync.Once

func NewBos(ctx context.Context) *Bos {
	once.Do(func() {

		wailsruntime.EventsEmit(ctx, "logEvent", map[string]interface{}{
			"type": "info",
			"msg":  "初始化Bos实例",
			"data": map[string]interface{}{
				"ak":                config.Get("bos.ak"),
				"sk":                config.Get("bos.sk"),
				"endpoint":          config.Get("bos.endpoint"),
				"redirect_disabled": config.GetBool("bos.redirect_disabled"),
			},
		})

		bosConfig := &bos.BosClientConfiguration{
			Ak:               config.Get("bos.ak"),
			Sk:               config.Get("bos.sk"),
			Endpoint:         config.Get("bos.endpoint"),
			RedirectDisabled: config.GetBool("bos.redirect_disabled"),
		}

		bosClient, err := bos.NewClientWithConfig(bosConfig)
		if err != nil {
			wailsruntime.EventsEmit(ctx, "logEvent", map[string]interface{}{
				"type":  "error",
				"msg":   "初始化Bos实例失败",
				"error": err.Error(),
			})
			panic(err)
		}

		cdnClient, err := cdn.NewClient(bosConfig.Ak, bosConfig.Sk, config.GetString("cdn.endpoint"))
		if err != nil {
			wailsruntime.EventsEmit(ctx, "logEvent", map[string]interface{}{
				"type":  "error",
				"msg":   "初始化Cdn实例失败",
				"error": err.Error(),
			})
			panic(err)
		}

		BosInstance = &Bos{
			BosClient:     bosClient,
			CdnClient:     cdnClient,
			DefaultBucket: config.GetString("bos.defaultBucket"),
			ctx:           ctx,
		}
	})

	return BosInstance
}

func (b *Bos) Upload(bucket, image string, name string) (innerpath string, outpath string, err error) {
	path := config.GetString("app.env") + "/" + cast.ToString(carbon.Now().Year()) + "/" + cast.ToString(carbon.Now().Month()) + "/" + cast.ToString(carbon.Now().Day()) + "/" + name + "/"
	filename := fmt.Sprintf("%x", md5.Sum([]byte(image))) + ".png"
	fullname := path + filename

	reader := bytes.NewBufferString(image)

	if bucket == "" {
		bucket = b.DefaultBucket
	}
	ok, _ := b.BosClient.DoesBucketExist(bucket)

	if !ok {
		_, err := b.BosClient.PutBucket(bucket)
		wailsruntime.EventsEmit(b.ctx, "logEvent", map[string]interface{}{
			"type":   "error",
			"msg":    "创建Bucket失败",
			"error":  err.Error(),
			"bucket": bucket,
		})
		if err != nil {
			return "", "", err
		}
	}

	_, err = b.BosClient.PutObjectFromStream(bucket, fullname, reader, nil)

	if err != nil {
		wailsruntime.EventsEmit(b.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "上传图片失败",
			"error": err.Error(),
			"data":  fullname,
		})

		return "", "", err
	}

	return "https://" + config.GetString("bos.inner_host") + "/" + fullname, "https://" + config.GetString("bos.host") + "/" + fullname, nil
}
