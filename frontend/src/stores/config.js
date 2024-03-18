import { defineStore } from 'pinia'

export const useConfigStore = defineStore('config', () => {

    const imageDisplay = ref(true)

    const getImageDisplay = () => {
        return imageDisplay.value
    }

    const setImageDisplay = () => {
        imageDisplay.value = !imageDisplay.value
    }

    return {imageDisplay, getImageDisplay, setImageDisplay}
})