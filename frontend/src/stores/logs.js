import { defineStore } from 'pinia'

export const useLogsStore = defineStore('logs', () => {
    const logs = ref("")
    const open = ref(false)
    const print = (data) => {
        if (typeof data === "string") {
            logs.value = logs.value + data + "\n"
        }
        if (typeof data === "object") {
            logs.value = logs.value + JSON.stringify(data, null, 2) + "\n"
        }
    }

    const switchLog = (value) => {
        open.value = value
    }

    const getOpen = () => {
        return open.value
    }

    const clear = () => {
        logs.value = ""
    }

    const getLogs = () => {
        return logs.value
    }


    return { logs, print, clear, getLogs, switchLog, open, getOpen }
})