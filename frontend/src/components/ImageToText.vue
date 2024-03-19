<template>

    <div class="w-card-active relative flex" :class="imageDisplay ? 'flex-row justify-between gap-8' : 'flex-col'">
        <div
            class="bg-black rounded-full absolute right-2 top-2 px-2 py-1 flex justify-center items-center text-white font-semibold text-xs">
            {{ props.preview.id }}</div>
        <div class="overflow-hidden" :class="imageDisplay ? 'w-48 h-32' : 'h-16'">
            <n-image class="w-auto overflow-hidden" object-fit="cover" :src="props.preview.url" />
        </div>
        <div class="flex flex-col items-start w-auto mb-3">
            <div class="text-gray-700 text-sm my-2 font-bold">PROMPT&HISTORY</div>
            <div class="flex gap-2" :class="imageDisplay ? 'flex-row' : 'flex-col w-full'">
                <div class="text-xs flex flex-col justify-between gap-1 w-full">
                    <span class="mr-1">prompt:</span>
                    <textarea class="form__input" v-model="props.preview.prompt" />
                </div>
                <div class="text-xs flex flex-col justify-between gap-1 w-full">
                    <span class="mr-1">history:</span>
                    <textarea class="form__input" v-model="props.preview.history" />
                </div>
            </div>
        </div>

        <div class="flex flex-col items-start w-auto">
            <div class="text-gray-700 text-sm my-2 font-bold">Image To Text</div>
            <div class="flex gap-2" :class="imageDisplay ? 'flex-row' : 'flex-col w-full'">
                <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                    <span class="mr-1">result:</span>
                    <textarea class="form__input" v-model="props.preview.result" />
                </div>
                <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                    <span class="mr-1">history msg:</span>
                    <textarea class="form__input" v-model="props.preview.history_msg" />
                </div>
                <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                    <span class="mr-1">ocr_ret:</span>
                    <textarea class="form__input" v-model="props.preview.ocr_ret" />
                </div>
                <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                    <span class="mr-1">face_ret:</span>
                    <textarea class="form__input" v-model="props.preview.face_ret" />
                </div>
            </div>
        </div>

        <div class="flex flex-col items-start w-auto mb-3">
            <div class="text-gray-700 text-sm my-2 font-bold">LLM</div>
            <div class="flex gap-2" :class="imageDisplay ? 'flex-row' : 'flex-col w-full'">
                <div class="text-xs flex flex-col justify-between gap-1 w-full">
                    <span class="mr-1">CHAT_ID:</span>
                    <textarea class="form__input" v-model="props.preview.chat_id" />
                </div>
                <div class="text-xs flex flex-col justify-between gap-1 w-full">
                    <span class="mr-1">CONTENT:</span>
                    <textarea class="form__input" v-model="props.preview.content" />
                </div>
            </div>
        </div>

    </div>
</template>

<script setup>
import { defineProps, h, ref, onMounted, defineEmits, watchEffect } from "vue";

import { useConfigStore } from "../stores/config"

const configSotre = useConfigStore()

const imageDisplay = ref(false)

watchEffect(() => {
    imageDisplay.value = configSotre.getImageDisplay()
})

const props = defineProps({
    preview: Object,
})

const emits = defineEmits(["delete"]);

emits("delete", props.preview.id)

</script>