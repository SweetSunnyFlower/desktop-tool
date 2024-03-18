import { defineStore } from 'pinia'

export const useImage2TextStore = defineStore('image2text', () => {

    // 表格预览数据
    const preview = ref([])

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

    // 图生文进度百分比，遍历预览数据，判断是否完成
    const getIsFinish = () => {
        if (preview.value.length == 0) {
            return false
        }
        let image2TextCount = 0
        // 查找出result != "" 的数量
        preview.value.forEach(item => {
            console.log(item)
            if (item.result && item.result != "") {
                image2TextCount = image2TextCount + 1
            }
        })

        return image2TextCount == preview.value.length
    }

    return { preview, getPreview, appendPreview, getIsFinish, bindPrompt, clearPreview }
})