import { defineStore } from 'pinia'

export const useImage2TextStore = defineStore('image2text', () => {

    // 表格预览数据
    const preview = ref([])

    const image2textCount = ref(0)

    const image2textfinish = ref(false)

    const template = ref("")

    const clearTemplate = () => {
        template.value = ""
    }

    const getTemplate = () => {
        return template.value
    }

    const setTemplate = (data) => {
        template.value = data
    }

    const clearPreview = () => {
        preview.value = []
    }

    const getPreview = () => {
        return preview.value
    }

    // 图片上传完成
    const appendPreview = (data) => {
        preview.value = [...preview.value, data]
    }

    // 绑定prompt数据
    const bindPrompt = (data) => {
        preview.value.forEach(item => {
            let prompt = data.find(prompt => prompt.id == item.id)
            item["prompt"] = prompt.prompt
            item["history"] = prompt.history
        })
    }

    const bindLLM = (data) => {
        preview.value.forEach(item => {
            if (item.id != data.id) {
                return
            }
            item["chat_id"] = data.chat_id
            item["content"] = data.content
        })
    }

    const clearimage2textCount = () => {
        image2textCount.value = 0
    }

    // 图生文完成
    const bindImage2Text = (data) => {
        image2textCount.value = image2textCount.value + 1
        preview.value.forEach(item => {
            if (item.id != data.id){
                return
            }
            let vis = data
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

        image2textfinish.value = image2textCount.value == preview.value.length
    }

    // 图生文进度百分比，遍历预览数据，判断是否完成
    const getIsFinish = () => {
        return image2textfinish.value
    }

    return { preview, getPreview, bindLLM, appendPreview, getIsFinish, bindPrompt, clearPreview, bindImage2Text, clearimage2textCount, template, clearTemplate, getTemplate, setTemplate }
})