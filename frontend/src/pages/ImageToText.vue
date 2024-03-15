<script setup>
import SelectPath from "../components/Path.vue";
import { h, ref, onMounted, reactive } from "vue";
import { Replace, UploadImage, ParsePromptFile } from '../../wailsjs/go/main/App'
import { useMessage, useNotification, NInput, NImage, NButton, NSpin } from "naive-ui";

const data = reactive(
    {
        "input": "",
        "output": "",
        "data": [
            {
                key: "",
                value: ""
            }
        ],
        "file_name": []
    }
)

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
                    pre_data.value[index].prompt = v;
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
                    pre_data.value[index].history = v;
                }
            });
        }
    }
];

const pre_data = ref([

])

const message = useMessage();
const notification = useNotification()

const placeholderInput = ref("输入")

const placeholderOutput = ref("输出")

const placeholder = reactive(["输入文件名关键字", "输出文件名关键字"])

const selectInput = (path) => {
    data.input = path
    console.log(data.input)
}
const prompt_file = ref("")

const selectPromptFile = (path) => {
    prompt_file.value = path
}

const parsePromptFile = () => {
    console.log(prompt_file.value)
    ParsePromptFile(prompt_file.value).then(res => {
        if (res.code == 0) {

            // 遍历pre_data 如果id存在，则将数据追加到data中

            pre_data.value.forEach(item => {
                let prompt = res.data.find(prompt => prompt.id == item.id)

                console.log("aaaaa", prompt);

                item["prompt"] = prompt.prompt
                item["history"] = prompt.history
            })
            message.info(res.message)
            showModal.value = false
            // pre_data.value = res.data
        } else {
            message.error(res.message)
        }
    })
}

const selectOutput = (path) => {
    data.output = path
    console.log(data.output)
}

const vis = () => {
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
    console.log(pre_data.value)
}

const replace = () => {

    Replace(data.input, data.output, data.data, data.file_name).then(res => {
        if (res.code == 0) {
            message.info(res.message)
            data.data = [{
                key: "",
                value: ""
            }]
        } else {
            message.error(res.message)
        }
    })
}

const uploadImage = () => {
    UploadImage(data.input).then(res => {
        if (res.code == 0) {
            console.log(res.data)
            message.info(res.message)
            pre_data.value = res.data
        } else {
            message.error(res.message)
        }
    })
}

const handleInputInput = () => {
    console.log(data.file_name)
}

const showModal = ref(false)

const outputText = ref("输出结果")
const handling = ref(false)

</script>

<template>
    <n-modal v-model:show="showModal" preset="card" :style="{ width: '600px' }" title="prompt">
        <template #header>
            <div>请选择关联的Prompt</div>
        </template>
        <div>
            <select-path :placeholder="placeholderInput" type="file" @click-path="selectPromptFile" />
        </div>
        <template #action>
            <div class="flex flex-row justify-end">
                <n-button strong dashed round @click="parsePromptFile">解析</n-button>
            </div>
        </template>
    </n-modal>
    <n-spin :show="handling">
        <div class="m-4 text-black">
            <div class=" bg-gray-100 rounded-xl p-3 mb-4 flex flex-col gap-3">
                <select-path :placeholder="placeholderInput" type="dir" @click-path="selectInput" />
                <n-button strong dashed round @click="uploadImage">上传到Bos(文件名为ID)</n-button>
            </div>

            <div class=" bg-gray-100 rounded-xl p-3 mb-4 flex flex-col gap-3">
                <n-data-table :columns="columns" :data="pre_data" />
                <n-button strong dashed round @click="showModal = true">上传关联prompt</n-button>
            </div>
            <div class=" bg-gray-100 rounded-xl p-3 mb-4 flex flex-col gap-3">
                <select-path :placeholder="placeholderOutput" type="dir" @click-path="selectOutput" />
                <n-button strong dashed round @click="vis">{{ outputText }}</n-button>
            </div>
        </div>
    </n-spin>

</template>
<style scoped>
</style>