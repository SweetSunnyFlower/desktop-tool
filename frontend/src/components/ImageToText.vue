<template>

    <div class="w-card-active relative">
        <div class="bg-black rounded-full absolute right-2 top-2 px-2 py-1 flex justify-center items-center text-white font-semibold text-xs">{{ props.preview.id }}</div>
        <div class="overflow-hidden h-16">
            <n-image class="w-auto overflow-hidden" object-fit="cover" :src="props.preview.url" />
        </div>
        <div class="flex flex-col items-start mb-3">
            <div class="text-gray-700 text-sm my-2 font-bold">PROMPT&HISTORY</div>
            <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                <span class="mr-1">prompt:</span>
                <textarea class="form__input" v-model="props.preview.prompt"/>
            </div>
            <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                <span class="mr-1">history:</span>
                <textarea class="form__input" v-model="props.preview.history"/>
            </div>
        </div>

        <div class="flex flex-col items-start">
            <div class="text-gray-700 text-sm my-2 font-bold">Image To Text</div>
            <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                <span class="mr-1">result:</span>
                <textarea class="form__input" v-model="props.preview.result"/>
            </div>
            <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                <span class="mr-1">history msg:</span>
                <textarea class="form__input" v-model="props.preview.history_msg"/>
            </div>
            <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                <span class="mr-1">ocr_ret:</span>
                <textarea class="form__input" v-model="props.preview.ocr_ret"/>
            </div>
            <div class="text-xs mb-2 flex flex-col justify-between gap-1 w-full">
                <span class="mr-1">face_ret:</span>
                <textarea class="form__input" v-model="props.preview.face_ret"/>
            </div>
        </div>
    </div>
</template>

<script setup>
import { defineProps, h, ref, onMounted, defineEmits, watch } from "vue";

const props = defineProps({
    preview: Object,
})

const emits = defineEmits(["delete"]);

emits("delete", props.preview.id)

</script>