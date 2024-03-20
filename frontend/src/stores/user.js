import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {

    const user = ref({
        email: '',
        uid: '',
        cuid: ''
    })

    const getUser = () => {
        return user.value
    }

    const setUser = async (data) => {
        return new Promise((resolve, reject) => {
            user.value = data
            resolve()
        })
    }

    return { user, getUser, setUser }
}, {
    storage: true,
})