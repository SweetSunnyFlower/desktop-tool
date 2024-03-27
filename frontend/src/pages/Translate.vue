<template>
    <n-spin :show="handling">
        <div class="m-10">
            <n-input v-model:value="word" placeholder="请输入内容"></n-input>
            <n-button type="info" class="mt-4" dashed @click="translate">
                文生图
            </n-button>
            <div class="mt-4">
                {{ result }}
            </div>
            <n-image class="mt-4" width="512" v-if="result" :src="result" />
        </div>
    </n-spin>
</template>

<script setup>
import { Translate } from '../../wailsjs/go/main/App'
import { onMounted, onUnmounted } from 'vue';
import { LogPrint, EventsOn, EventsOff } from "../../wailsjs/runtime"

const word = ref('')

const result = ref('')

const handling = ref(false)

onUnmounted(() => {
    // 取消事件监听
    EventsOff("handlingEvent")
})

onMounted(() => {
    // 处理事件
    EventsOn("handlingEvent", function (data) {
        handling.value = data
    })
})

const translate = () => {
    if (!word.value) {
        return
    }
    Translate(word.value, 'auto', 'en').then(res => {
        result.value = res
    })
}


</script>