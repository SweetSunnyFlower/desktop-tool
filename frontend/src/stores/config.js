import { defineStore } from 'pinia'

export const useConfigStore = defineStore('config', () => {

    const imageDisplay = ref(true)

    const getImageDisplay = () => {
        return imageDisplay.value
    }

    const setImageDisplay = () => {
        imageDisplay.value = !imageDisplay.value
    }

    const model = ref("qianfan_chat")

    const getModel = () => {
        return model.value
    }

    const setModel = (value) => {
        model.value = value
    }

    return {imageDisplay, getImageDisplay, setImageDisplay, model, getModel, setModel}
})