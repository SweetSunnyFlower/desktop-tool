<template>
    <n-spin :show="handling">
        <div class="m-4 text-3xl flex flex-row justify-between items-center relative text-gray-700">
            <!-- background-image: linear-gradient(120deg, #84fab0 0%, #8fd3f4 100%); -->
            <div class="flex flex-row items-baseline gap-2">
                <!-- <n-gradient-text gradient="linear-gradient(90deg, #84fab0 0%, #8fd3f4 100%)">
                    图生文批量处理工具
                </n-gradient-text> -->
                <div class="flex flex-row font-bold text-xs gap-1">
                    <div class="flex flex-row  gap-1 w-52">
                        <div class="text-xs">横向</div>
                        <div @click="changeDisplay"
                            :class="imageDisplay ? 'h-4 w-10 nm-inset-gray-200 rounded-full flex justify-start' : 'h-4 w-10 nm-inset-gray-300 rounded-full flex justify-end'">
                            <div class="h-4 w-4 bg-gray-200 shadow-gray-200 transform scale-110 rounded-full"></div>
                        </div>
                        <div class="text-xs">纵向</div>
                    </div>
                    <n-slider v-if="!imageDisplay" v-model:value="colNumber" :step="1" :max="5" :min="1"/>
                </div>
            </div>
            <div class="flex lg:flex-row sm:flex-col md:flex-col sm:top-0  justify-between gap-3 right-2 z-50">
                <button class="w-button px-4 py-2" @click="openFolder('images')">
                    <n-gradient-text gradient="linear-gradient(90deg, #84fab0 0%, #8fd3f4 100%)">
                        选择照片
                    </n-gradient-text>
                </button>
                <button class="w-button px-4 py-2" @click="openFolder('download-template')">
                    <n-gradient-text gradient="linear-gradient(90deg, #fee140 0%, #96e6a1 100%)">
                        下载prompt模版
                    </n-gradient-text>
                </button>
                <button class="w-button px-4 py-2" @click="openFile('prompt')">
                    <n-gradient-text gradient="linear-gradient(90deg, #4facfe 0%, #00f2fe 100%)">
                        上传关联prompt
                    </n-gradient-text>
                </button>
                <button class="w-button px-4 py-2" @click="image2Text">
                    <n-gradient-text gradient="linear-gradient(90deg, #fa709a 0%, #fee140 100%)">
                        图生文
                    </n-gradient-text>
                </button>
                <button class="w-button px-4 py-2 relative overflow-visible" @click="showmetion = !showmetion">
                    <n-gradient-text gradient="linear-gradient(90deg, #a8edea 0%, #fed6e3 100%)">
                        设置提问模版
                    </n-gradient-text>
                    <div class="absolute top-12 left-0 w-64 overflow-hidden">
                        <n-select @click.stop="showmetion = showmetion" v-if="showmetion" v-model:value="model"
                            @update:value="handleModelUpdateValue" size="tiny" :options="models" />
                        <n-mention v-if="showmetion" @click.stop="showmetion = showmetion" class="text-left" type="text"
                            :value="template" :options="templateOptions" prefix="@" :on-update:value="mention" />
                    </div>
                </button>
                <button class="w-button px-4 py-2" @click="llm">
                    <n-gradient-text gradient="linear-gradient(90deg, #d299c2 0%, #96e6a1 100%)">
                        大模型对话
                    </n-gradient-text>
                </button>
                <button class="w-button px-4 py-2" v-if="image2textfinish" @click="openFolder('download-iamge2text')">
                    <n-gradient-text gradient="linear-gradient(90deg, #89f7fe 0%, #66a6ff 100%)">
                        下载
                    </n-gradient-text>
                </button>
                <button class="w-button px-4 py-2" @click="clear">
                    <n-gradient-text gradient="linear-gradient(90deg, #8fd3f4 0%, #84fab0 100%)">
                        清理数据
                    </n-gradient-text>
                </button>
            </div>
        </div>
        <div class="grid-cols-1"></div>
        <div class="grid-cols-2"></div>
        <div class="grid-cols-3"></div>
        <div class="grid-cols-4"></div>
        <div class="grid-cols-5"></div>
        <div class="m-4 text-gray-700 grid gap-5" :class="imageDisplay ? '' : `lg:grid-cols-${colNumber} sm:grid-cols-${colNumber}`">
            <image-to-text-view v-for="(item, index) in preview" :preview="item" />
            <!-- <div class="nm-flat-white-xs p-3 mb-4 flex flex-col gap-3">
                <n-data-table size="small" ref="tableRef" :bordered="false" :single-line="false" :scroll-x="1800"
                    :row-key="rowKey" @update:checked-row-keys="handleCheck" :style="{ height: `${height}px` }"
                    flex-height :columns="columns" :data="preview" />
            </div> -->
        </div>
    </n-spin>
</template>

<script setup>
import ImageToTextView from "../components/ImageToText.vue"
import { h, ref, onMounted, onUnmounted, watchEffect } from "vue";
import { OpenFile, OpenFolder, Image2Text, LLM } from '../../wailsjs/go/main/App'
import { LogPrint, EventsOn, EventsOff } from "../../wailsjs/runtime"
import { useMessage, useNotification, NInput, NImage, NSpin } from "naive-ui";
import { DownloadOutline } from "@vicons/ionicons5";
import { useImage2TextStore } from "../stores/image2text"
import { useConfigStore } from "../stores/config"

const image2textStore = useImage2TextStore()
const configSotre = useConfigStore()

const showmetion = ref(false)
const imageDisplay = ref(false)
const model = ref("")
const models = [
    {
        label: '大模型3.5',
        value: 'qianfan_chat'
    },
    {
        label: '大模型4.0',
        value: 'EB40'
    },
]

const colNumber = ref(3)

const changeDisplay = () => {
    configSotre.setImageDisplay()
}

// 预览数据
const preview = ref([])
const image2textfinish = ref(false)

const template = ref("")

const templateOptions = ref([
    {
        label: 'id',
        value: 'id',
    },
    {
        label: 'url',
        value: 'url',
    },
    {
        label: 'prompt',
        value: 'prompt',
    },
    {
        label: 'history',
        value: 'history',
    },
    {
        label: 'result',
        value: 'result',
    },
    {
        label: 'history_msg',
        label: 'history_msg',
    },
    {
        label: 'ocr_ret',
        value: 'ocr_ret',
    },
    {
        label: 'face_ret',
        value: 'face_ret',
    },
    {
        label: 'chat_id',
        value: 'chat_id',
    },
    {
        label: 'content',
        value: 'content',
    }
])

const mention = (value) => {
    image2textStore.setTemplate(value)
}

const handleModelUpdateValue = (value) => {
    configSotre.setModel(value)
}

watchEffect(() => {
    model.value = configSotre.getModel()
    image2textfinish.value = preview.value.length > 0
    preview.value = image2textStore.getPreview()
    imageDisplay.value = configSotre.getImageDisplay()
    template.value = image2textStore.getTemplate()
    // preview.value.forEach(item => {
    //     templateOptions.value = Object.keys(item).map(item => {
    //         return { label: item, value: item };
    //     });
    // })
})

onUnmounted(() => {
    // 取消事件监听
    EventsOff("handlingEvent")
    EventsOff("image2TextEvent")
})
onMounted(() => {
    // 处理事件
    EventsOn("handlingEvent", function (data) {
        handling.value = data
    })

    // 图生文事件
    EventsOn("image2TextEvent", function (data) {
        image2textStore.bindImage2Text(data)
    })

    EventsOn("llmEvent", function (data) {
        image2textStore.bindLLM(data)
    })
})

const message = useMessage();
const notification = useNotification()
const outputText = ref("图生文")
const handling = ref(false)
const columns = [
    {
        type: "selection",
    },
    {
        title: "ID",
        key: "id",
        fixed: "left",
        width: 100,
        align: "center",
    },
    {
        width: 100,
        title: "URL",
        align: "center",
        key: "url",
        fixed: "left",
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
        title: "关联Prompts",
        key: "prompts",
        align: "center",
        children: [
            {
                title: "Prompt",
                key: "prompt",
                align: "center",
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
                align: "center",
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
        ]
    },
    {
        title: "文生图",
        key: "image2text",
        align: "center",
        children: [
            {
                title: "Result",
                key: "result",
                align: "center",
                render(row, index) {
                    return h(NInput, {
                        value: row.result,
                        type: "textarea",
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
                align: "center",
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
                align: "center",
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
                align: "center",
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
        ]
    },
];

// 后端解析Prompt返回内容，前端追加内容到table中
const parsePromptFile = (response) => {
    if (response.code == 0) {
        // 遍历preview 如果id存在，则将数据追加到data中
        image2textStore.bindPrompt(response.data)
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

    if (type == "download-iamge2text") {
        body = JSON.stringify(preview.value)
    }

    OpenFolder(type, body).then(res => {
        if (res.code == 2) {
            return
        }
        if (res.code == 1) {
            message.error(res.message)
            return
        }
        if (type == "download-template") {
            message.info(res.message)
        }
        if (type == "download-iamge2text") {
            message.info(res.message)
        }
        if (type == "download-data") {
        }
    })
}

const setModel = () => {
    model.value = model.value
}

const clear = () => {
    image2textStore.clearPreview()
    image2textStore.clearimage2textCount()
}

// 图生文接口
const image2Text = () => {
    let body = JSON.stringify(preview.value)
    Image2Text(body)
}

const llm = () => {
    let body = JSON.stringify(preview.value)
    LLM(model.value, template.value, body)
}
const tableRef = ref();
const height = ref(420)
const checkedRowKeysRef = ref([]);

const rowKey = (row) => row.id
const handleCheck = (rowKeys) => {
    checkedRowKeysRef.value = rowKeys;
}
</script>


<style scoped></style>