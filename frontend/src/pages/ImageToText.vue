<script setup>
import SelectPath from "../components/Path.vue";
import { h, ref, onMounted, reactive } from "vue";
import { Replace, UploadImage, ParsePromptFile } from '../../wailsjs/go/main/App'
import { useMessage, useNotification, NInput, NImage, NButton, NSpin } from "naive-ui";

const message = useMessage();
const notification = useNotification()
const imagesDir = ref("")
const promptFile = ref("")
const outputDir = ref("")
const preview = ref([])
const imagePlaceholder = ref("选择图片")
const promptPlaceholder = ref("选择图片提示词文件")
const placeholderOutput = ref("选择输出位置")
const showModal = ref(false)
const outputText = ref("输出结果")
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
    }
];

const selectImages = (path) => {
    imagesDir.value = path
}
const selectOutputDir = (path) => {
    outputDir.value = path
}
const selectPromptFile = (path) => {
    promptFile.value = path
}

const parsePromptFile = () => {
    ParsePromptFile(promptFile.value).then(response => {
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
    })
}

const imageToText = () => {
    handling.value = true
    let percent = 0
    let markAsRead = false
    const n = notification.create({
        title: '重要通知',
        content: `请不要操作当前界面`,
        meta: new Date().toLocaleString(),
        action: () =>
            h(
                NButton,
                {
                    text: true,
                    type: 'primary',
                    onClick: () => {
                        markAsRead = true
                        n.destroy()
                    }
                },
                {
                    default: () => '已读'
                }
            ),
        onClose: () => {
            if (!markAsRead) {
                message.warning('请设为已读')
                return false
            }
        }
    })
    outputText.value = "输出结果" + `${percent}%`
    let interval = setInterval(() => {
        if (percent >= 100) {
            handling.value = false
            clearInterval(interval)
            const n = notification.create({
                title: '导出通知',
                content: `请查看输出目录`,
                meta: new Date().toLocaleString(),
                onClose: () => {

                }
            })
            return
        }
        percent = percent + 10
        outputText.value = "输出结果" + `${percent}%`
    }, 30)
}

const uploadImage = () => {
    handling.value = true
    UploadImage(imagesDir.value).then(response => {
        console.log(response)
        handling.value = false
        if (response.code == 0) {
            message.info(response.message)
            preview.value = response.data
        } else {
            message.error(response.message)
        }
    })
}
</script>

<template>
    <n-modal v-model:show="showModal" preset="card" :style="{ width: '600px' }" title="prompt">
        <template #header>
            <div>请选择关联的Prompt</div>
        </template>
        <div>
            <select-path :placeholder="promptPlaceholder" type="file" @click-path="selectPromptFile" />
        </div>
        <template #action>
            <div class="flex flex-row justify-end">
                <n-button strong dashed round @click="parsePromptFile">解析</n-button>
            </div>
        </template>
    </n-modal>
    <n-spin :show="handling">
        <div class="m-4 text-3xl flex flex-row justify-between items-baseline">
            <!-- background-image: linear-gradient(120deg, #84fab0 0%, #8fd3f4 100%); -->
            <n-gradient-text gradient="linear-gradient(90deg, #84fab0 0%, #8fd3f4 100%)">
                图生文批量处理工具
            </n-gradient-text>
            <n-button type="error" dashed round class="text-xs">
                初始化
            </n-button>
        </div>
        <div class="m-4 text-black">
            <div class=" bg-gray-100 rounded-xl p-3 mb-4 flex flex-col gap-3">
                <select-path :placeholder="imagePlaceholder" type="dir" @click-path="selectImages" />
                <n-button strong dashed round @click="uploadImage">上传到Bos(文件名为ID)</n-button>
            </div>

            <div class=" bg-gray-100 rounded-xl p-3 mb-4 flex flex-col gap-3">
                <n-data-table :columns="columns" :data="preview" />
                <n-button strong dashed round @click="showModal = true">上传关联prompt</n-button>
            </div>
            <div class=" bg-gray-100 rounded-xl p-3 mb-4 flex flex-col gap-3">
                <select-path :placeholder="placeholderOutput" type="dir" @click-path="selectOutputDir" />
                <n-button strong dashed round @click="imageToText">{{ outputText }}</n-button>
            </div>
        </div>
    </n-spin>
</template>
<style scoped></style>