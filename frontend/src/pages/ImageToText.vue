<script setup>
import SelectPath from "../components/Path.vue";
import { h, ref, onMounted, reactive } from "vue";
import { OpenFile, OpenFolder } from '../../wailsjs/go/main/App'
import { LogPrint, EventsOn } from "../../wailsjs/runtime"
import { useMessage, useNotification, NInput, NImage, NButton, NSpin } from "naive-ui";
import { DownloadOutline } from "@vicons/ionicons5";
onMounted(() => {
    // 下载模版事件
    EventsOn("downloadTemplate", function (data) {
        downloadCSV(data)
    })
    // 处理事件
    EventsOn("handlingEvent", function (data) {
        handling.value = data
    })
    // 日志事件
    EventsOn("logEvent", function (data) {
        log.value = log.value + data + "\n"
    })
    // 上传图片事件
    EventsOn("uploadImageEvent", function (data) {
        console.log(data)
        preview.value = [...preview.value, data]
    })
    // 图生文事件
    EventsOn("image2TextEvent", function (data) {
        // 输入日志
        console.log(data)

        percent.value = percent.value + 1
        outputText.value = "文生图" + `${percent.value / preview.value.length * 100}%`

        preview.value.forEach(item => {
            let vis = data.find(vis => vis.id == item.id)
            item["result"] = vis.result
            item["face_ret"] = vis.face_ret
            item["ocr_ret"] = vis.ocr_ret
            item["history_msg"] = vis.history_msg.join("|")
        })

        if (preview.value.length == percent.value) {
            image2textfinish.value = true
            notification.create({
                title: '导出通知',
                content: `请查看输出目录`,
                meta: new Date().toLocaleString(),
                onClose: () => {

                }
            })
        }
    })
})
const image2textfinish = ref(false)
const percent = ref(0)
const message = useMessage();
const notification = useNotification()
const preview = ref([])
const showModal = ref(false)
const outputText = ref("文生图")
const handling = ref(false)
const columns = [
    {
        title: "ID",
        key: "id",
    },
    {
        title: "URL",
        key: "url",
        render(row, index) {
            return h(NImage, {
                src: row.url,
                width: 100,
                height: 100,
                "object-fit": "cover",
            });
        }
    },
    {
        title: "Prompt",
        key: "prompt",
        render(row, index) {
            return h(NInput, {
                value: row.prompt,
                placeholder: "请输入提示词,或者导入",
                onUpdateValue(v) {
                    preview.value[index].prompt = v;
                }
            });
        }
    },
    {
        title: "History",
        key: "history",
        render(row, index) {
            return h(NInput, {
                value: row.history,
                placeholder: "请输入history,或者导入",
                onUpdateValue(v) {
                    preview.value[index].history = v;
                }
            });
        }
    },
    {
        title: "Result",
        key: "result",
        render(row, index) {
            return h(NInput, {
                value: row.result,
                placeholder: "请输入result",
                onUpdateValue(v) {
                    preview.value[index].result = v;
                }
            });
        }
    },
    {
        title: "face_ret",
        key: "face_ret",
        render(row, index) {
            return h(NInput, {
                value: row.face_ret,
                placeholder: "请输入face_ret",
                onUpdateValue(v) {
                    preview.value[index].face_ret = v;
                }
            });
        }
    },
    {
        title: "oct_ret",
        key: "oct_ret",
        render(row, index) {
            return h(NInput, {
                value: row.oct_ret,
                placeholder: "请输入oct_ret",
                onUpdateValue(v) {
                    preview.value[index].oct_ret = v;
                }
            });
        }
    },
    {
        title: "history_msg",
        key: "history_msg",
        render(row, index) {
            return h(NInput, {
                value: row.history_msg,
                placeholder: "请输入history_msg",
                onUpdateValue(v) {
                    preview.value[index].history_msg = v;
                }
            });
        }
    }
];

// 后端返回内容，前端下载文件
const downloadCSV = (data) => {
    const csvContent = data.map(row => row.join(",")).join("\n");
    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
    const link = document.createElement("a");
    if (link.download !== undefined) {
        const url = URL.createObjectURL(blob);
        link.setAttribute("href", url);
        link.setAttribute("download", "data.csv");
        link.style.visibility = "hidden";
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }
}

// 后端解析Prompt返回内容，前端追加内容到table中
const parsePromptFile = (response) => {
    if (response.code == 0) {
        // 遍历preview 如果id存在，则将数据追加到data中
        preview.value.forEach(item => {
            let prompt = response.data.find(prompt => prompt.id == item.id)
            item["prompt"] = prompt.prompt
            item["history"] = prompt.history
        })
        message.info(response.message)
        showModal.value = false
    } else {
        message.error(response.message)
    }
}

// 上传图片回调
const uploadImage = (response) => {
    if (response.code == 0) {
        message.info(response.message)
    } else {
        message.error(response.message)
    }
}

// 打开文件
const openFile = (type) => {
    OpenFile(type).then(res => {
        if (res.code == 2) {
            return
        }
        if (res.code == 1) {
            message.error(res.message)
            return
        }
        if (type == "prompt") {
            parsePromptFile(res)
        }
    })
}

// 打开文件夹
const openFolder = (type) => {
    let body = ""
    if (type == "image2text") {
        body = JSON.stringify(preview.value)
        console.log(type, body)
    }
    OpenFolder(type, body).then(res => {
        if (res.code == 2) {
            return
        }
        if (res.code == 1) {
            message.error(res.message)
            return
        }
        if (type == "images") {
            uploadImage(res)
        }
        if (type == "download-template") {
            message.info(res.message)
        }
        if (type == "download-data") {
        }
    })
}

const more = [
    {
        label: '下载模版',
        key: 'download-template',
        disabled: false
    },
    {
        label: '初始化',
        key: 'init',
        disabled: false
    }
]
const handleSelectMore = (item) => {
    if (item == "download-template") {
        openFolder('download-template')
    }
    if (item == "init") {
        preview.value = []
        log.value = ""
        handling.value = false
        outputText.value = "文生图"
        percent.value = 0
        showModal.value = false
        image2textfinish.value = false
    }
}
const log = ref("")
const height = ref(420)
</script>

<template>
    <n-spin :show="handling">
        <div class="m-4 text-3xl flex flex-row justify-between items-center">
            <!-- background-image: linear-gradient(120deg, #84fab0 0%, #8fd3f4 100%); -->
            <n-gradient-text gradient="linear-gradient(90deg, #84fab0 0%, #8fd3f4 100%)">
                图生文批量处理工具
            </n-gradient-text>
            <n-dropdown trigger="hover" :options="more" @select="handleSelectMore">
                <span class="text-lg hover:cursor-pointer">...</span>
            </n-dropdown>
        </div>
        <div class="m-4 text-black">
            <div class=" bg-gray-100 rounded-xl p-3 mb-4 flex flex-col gap-3">
                <n-data-table size="small" :style="{ height: `${height}px` }" flex-height :columns="columns"
                    :data="preview" />
                <div class="flex flex-row justify-between gap-3">
                    <n-button strong dashed round @click="openFolder('images')">选择照片</n-button>
                    <n-button strong dashed round @click="openFile('prompt')">上传关联prompt</n-button>
                    <n-button strong dashed round icon-placement="right" @click="openFolder('image2text')">
                        <div class="flex flex-row justify-between items-center gap-1">
                            {{ outputText }}
                            <n-icon size="14" v-if="image2textfinish">
                                <download-outline />
                            </n-icon>
                        </div>
                    </n-button>
                </div>
            </div>
        </div>
        <div class="m-4">
            <div class=" bg-gray-100 rounded-xl p-3 mb-4">
                <n-log :rows="10" :log="log" show-line-numbers word-wrap language="javascript" />
            </div>
        </div>
    </n-spin>
</template>

<style scoped></style>