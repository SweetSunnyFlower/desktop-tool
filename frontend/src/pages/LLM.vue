<script setup>

import SelectPath from "../components/Path.vue";
import { ParseVisFile } from '../../wailsjs/go/main/App'

import { useMessage, useNotification, NInput, NImage, NButton, NSpin } from "naive-ui";


const message = useMessage();
const notification = useNotification()
const handling = ref(false);

const filePath = ref('');

const selectInput = (path) => {
    filePath.value = path
    console.log(filePath.value)
}

const placeholderInput = ref("请选择文件")

const preview = ref([])

const parseFile = () => {
    ParseVisFile(filePath.value).then(res => {
        if (res.code == 0) {
            res.data.forEach(item => {
                options.value = Object.keys(item).map(item => {
                    return { label: item, value: item };
                });
            })
            preview.value = res.data
            message.info(res.message)
        } else {
            message.error(res.message)
        }
    })
}

const columns = [
    {
        title: "ID",
        key: "id",
    },
    {
        title: "result",
        key: "result",
    },
    {
        title: "history msg",
        key: "history_msg"
    },
    {
        title: "ocr ret",
        key: "ocr_ret"
    },
    {
        title: "face ret",
        key: "face_ret"
    }
];
const options = ref([])
const template = ref('')
const mention = (value) => {
    template.value = value
}
const llm = () => {
    console.log(template.value)
}
</script>

<template>
    <n-spin :show="handling">
        <div class="m-4 text-3xl">
            <!-- background-image: linear-gradient(120deg, #84fab0 0%, #8fd3f4 100%); -->
            <n-gradient-text gradient="linear-gradient(90deg, #84fab0 0%, #8fd3f4 100%)">
                LLM大模型
            </n-gradient-text>
        </div>
        <div class="m-4 text-black">
            <div class=" bg-gray-100 rounded-xl p-3 mb-4 flex flex-col gap-3">
                <select-path :placeholder="placeholderInput" type="file" @click-path="selectInput" />
                <n-button strong dashed round @click="parseFile">选择图生文生成的文件</n-button>
                <n-data-table :columns="columns" :data="preview" />
                <n-form-item label="提问模版">
                    <n-mention type="textarea" :options="options" prefix="%" :on-update:value="mention"/>
                </n-form-item>
                <n-button strong dashed round @click="llm">批量提问</n-button>
            </div>
        </div>
    </n-spin>
</template>