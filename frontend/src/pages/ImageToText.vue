<script setup>
import SelectPath from "../components/Path.vue";
import { h, ref, onMounted, reactive } from "vue";
import { Replace, UploadImage, ParsePromptFile } from '../../wailsjs/go/main/App'
import { useMessage } from "naive-ui";
import { NInput, NImage } from "naive-ui";

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

                console.log("aaaaa",prompt);

                item["prompt"] = prompt.prompt
            })
            console.log(pre_data.value)
            message.info(res.message)
            pre_data.value = res.data
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

</script>

<template>
    <n-modal v-model:show="showModal" preset="dialog" title="Dialog">
        <template #header>
            <div>标题</div>
        </template>
        <div>
            <select-path :placeholder="placeholderInput" @click-path="selectPromptFile" />
        </div>
        <template #action>
            <n-button strong success round @click="parsePromptFile">解析</n-button>
        </template>
    </n-modal>
    <div id="form">
        <n-grid x-gap="12" :cols="2">
            <n-gi>
                <select-path :placeholder="placeholderInput" @click-path="selectInput" />
            </n-gi>
            <n-gi>
                <select-path :placeholder="placeholderOutput" @click-path="selectOutput" />
            </n-gi>
        </n-grid>
        <n-grid x-gap="12" :cols="2">
            <n-gi>
                <div @click="uploadImage" class="add-item">上传Bos</div>
            </n-gi>
            <n-gi>
                <div @click="showModal = true" class="add-item">上传关联prompt</div>
            </n-gi>
        </n-grid>
        <n-data-table :columns="columns" :data="pre_data" />
        <div @click="vis" class="add-item">vis图生文</div>
    </div>
</template>
<style scoped>
#form {
    margin: 20px;
}

.add-item {
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid black;
    border-radius: 4px;
    margin: 10px 0;
    cursor: pointer;
    transition: all 0.5 ease-in;
}

.add-item:hover {
    border: 1px solid #ccc;
    color: #ccc;
}
</style>