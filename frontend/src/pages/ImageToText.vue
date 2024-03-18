<script setup>
import { h, ref, onMounted, onUnmounted } from "vue";
import { OpenFile, OpenFolder, Image2Text } from '../../wailsjs/go/main/App'
import { LogPrint, EventsOn, EventsOff } from "../../wailsjs/runtime"
import { useMessage, useNotification, NInput, NImage, NSpin } from "naive-ui";
import { DownloadOutline } from "@vicons/ionicons5";
onUnmounted(() => {
    // 取消事件监听
    EventsOff("handlingEvent")
    EventsOff("uploadImageEvent")
    EventsOff("image2TextEvent")
})
onMounted(() => {
    // 处理事件
    EventsOn("handlingEvent", function (data) {
        handling.value = data
    })
    // 上传图片事件
    EventsOn("uploadImageEvent", function (data) {
        console.log(data)
        preview.value = [...preview.value, data]
        console.log(preview.value)
    })
    // 图生文事件
    EventsOn("image2TextEvent", function (data) {
        // 输入日志
        percent.value = percent.value + 1
        outputText.value = "文生图" + `${(percent.value / preview.value.length * 100).toFixed(2)}%`

        preview.value.forEach(item => {
            let vis = data.find(vis => vis.id == item.id)
            // 定义一个字符串变量用于保存结果
            let result = "";

            // 遍历二维数组并连接字符串
            for (let i = 0; i < vis.history_msg.length; i++) {
                let row = vis.history_msg[i];
                for (let j = 0; j < row.length; j++) {
                    result += row[j];
                    // 在每个元素后面添加 "|"，除了最后一个元素
                    if (j < row.length - 1) {
                        result += "|";
                    }
                }
                // 在每一行的末尾添加换行符，如果不需要可以去掉
                result += "\n";
            }

            item["result"] = vis.result
            item["face_ret"] = vis.face_ret
            item["ocr_ret"] = vis.ocr_ret
            item["history_msg"] = result
        })

        if (preview.value.length == percent.value) {
            image2textfinish.value = true
            notification.create({
                title: '图生文完成',
                content: `如需下载，请点击下载按钮`,
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
        if (type == "images") {
            uploadImage(res)
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


const clear = () => {
    preview.value = []
    handling.value = false
    outputText.value = "文生图"
    percent.value = 0
    showModal.value = false
    image2textfinish.value = false
}

// 图生文接口
const image2Text = () => {
    let body = JSON.stringify(preview.value)
    Image2Text(body).then(res => {
        console.log(res)
    })
}
const tableRef = ref();

const height = ref(420)
const rowKey = (row) => row.id
const checkedRowKeysRef = ref([]);
const handleCheck = (rowKeys) => {
    checkedRowKeysRef.value = rowKeys;
}
</script>

<template>
    <n-spin :show="handling">
        <div class="m-4 text-3xl flex flex-row justify-between items-center relative text-gray-700">
            <!-- background-image: linear-gradient(120deg, #84fab0 0%, #8fd3f4 100%); -->
            <n-gradient-text gradient="linear-gradient(90deg, #84fab0 0%, #8fd3f4 100%)">
                图生文批量处理工具
            </n-gradient-text>
            <div class="flex lg:flex-row sm:flex-col sm:top-0 justify-between gap-3 absolute right-2 z-50">
                <button
                    class="w-button px-4 py-2"
                    @click="openFolder('images')">
                    选择照片
                </button>
                <button
                    class="w-button px-4 py-2"
                    @click="openFolder('download-template')">
                    下载prompt模版
                </button>
                <button
                    class="w-button px-4 py-2"
                    @click="openFile('prompt')">
                    上传关联prompt
                </button>
                <button
                    class="w-button px-4 py-2"
                    @click="image2Text">
                    <div class="flex flex-row justify-between items-center gap-2">
                        {{ outputText }}
                        <n-icon size="14" v-if="image2textfinish" @click.stop="openFolder('download-iamge2text')">
                            <download-outline />
                        </n-icon>
                    </div>
                </button>
                <button
                    class="w-button px-4 py-2"
                    @click="clear">
                    清理数据
                </button>
            </div>
        </div>
        <div class="m-4 text-black">
            <div class="nm-flat-white-xs p-3 mb-4 flex flex-col gap-3">
                <n-data-table size="small" ref="tableRef" :bordered="false" :single-line="false" :scroll-x="1800"
                    :row-key="rowKey" @update:checked-row-keys="handleCheck" :style="{ height: `${height}px` }"
                    flex-height :columns="columns" :data="preview" />
            </div>
        </div>
    </n-spin>
</template>

<style scoped></style>