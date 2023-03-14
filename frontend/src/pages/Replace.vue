
<script setup>
import SelectPath from "../components/Path.vue";
import { h, ref, onMounted, reactive } from "vue";
import { Replace } from '../../wailsjs/go/main/App'
import { useMessage } from "naive-ui";

const data = reactive(
    {
        "input":"",
        "output":"",
        "data":[
            {
                key: "",
                value: ""
            }
        ],
        "file_name":[]
    }
)

const message = useMessage();

const placeholderInput = ref("输入")

const placeholderOutput = ref("输出")

const placeholder = reactive(["输入文件名关键字", "输出文件名关键字"])

const selectInput = (path) => {
    data.input = path
    console.log(data.input)
}

const selectOutput = (path) => {
    data.output = path
    console.log(data.output)
}

const replace = () => {

    Replace(data.input, data.output, data.data, data.file_name).then(res => {
        if (res.code == 0) {
            message.info(res.message)
            data.data = [ {
                key: "",
                value: ""
            }]
        } else {
            message.error(res.message)
        }
    })
}

const handleInputInput = () => {
    console.log(data.file_name)
}

</script>

<template>
    <div id="form">
        <n-grid x-gap="12" :cols="2">
            <n-gi>
                <select-path :placeholder="placeholderInput" @click-path="selectInput"/>
            </n-gi>
            <n-gi>
                <select-path :placeholder="placeholderOutput" @click-path="selectOutput"/>
            </n-gi>
        </n-grid>
        <n-grid :cols="1">
            <n-gi>
                <div style="margin:10px 0;">
                    <n-input
                        pair
                        separator="替换"
                        :placeholder="placeholder"
                        clearable
                        v-model:value="data.file_name"
                        @update:value="handleInputInput"
                    />
                </div>
            </n-gi>
        </n-grid>
        <n-grid :cols="1">
            <n-gi>
                <n-dynamic-input
                    v-model:value="data.data"
                    preset="pair"
                    show-sort-button
                    key-placeholder="old"
                    value-placeholder="new"
                    create-button-default="添加"
                >
            </n-dynamic-input>
            </n-gi>
        </n-grid>
        <n-grid :cols="1">
            <n-gi>
                <div @click="replace" class="add-item">替换</div>
            </n-gi>
        </n-grid>
    </div>
</template>
<style scoped>
#form {
    margin: 20px;
}
.add-item{
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
.add-item:hover{
    border: 1px solid #ccc;
    color: #ccc;
}
</style>