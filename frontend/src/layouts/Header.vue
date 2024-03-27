<template>
    <n-layout has-sider>
        <n-layout-sider bordered collapse-mode="width" :collapsed-width="58" :width="160" :height="1000"
            :collapsed="collapsed" show-trigger @collapse="collapsed = true" @expand="collapsed = false">
            <n-menu :options="menuOptions" :default-expanded-keys="defaultExpandedKeys"
                @update:expanded-keys="handleUpdateExpandedKeys" />
        </n-layout-sider>
        <n-layout class="h-screen">
            <router-view />
            <n-layout-footer bordered position="absolute" :class="open ? 'h-auto' : 'h-0'">
                <div class="relative nm-flat-neutral-100-lg p-2">
                    <div class="absolute   w-9 h-9 flex flex-row justify-center items-center rounded-full w-button"
                        :class="open ? '-top-5 right-16' : '-top-10 right-2'" @click="switchLog(!open)">
                        <n-icon>
                            <arrow-down-outline v-if="open" class="w-8 h-8" />
                            <arrow-up-outline v-if="!open" class="w-8 h-8" />
                        </n-icon>
                    </div>
                    <div class="absolute right-2 -top-5 w-9 h-9 flex flex-row justify-center items-center rounded-full w-button"
                        :class="open ? '-top-5' : 'top-0'" @click="clearLog()">
                        <n-icon>
                            <close-outline class="w-8 h-8" />
                        </n-icon>
                    </div>

                    <n-log :rows="10" :log="logs" show-line-numbers word-wrap language="javascript" />
                </div>
            </n-layout-footer>
        </n-layout>
    </n-layout>
</template>

<script setup>
import { h, ref, onMounted, watchEffect } from "vue";
import { NIcon, useMessage } from "naive-ui";
import { RouterLink } from "vue-router";
import { LogPrint, EventsOn, EventsOff } from "../../wailsjs/runtime"

import {
    DocumentOutline,
    HomeOutline,
    ImageOutline,
    CodeOutline,
    ArrowDownOutline,
    ArrowUpOutline,
    CloseOutline,
} from "@vicons/ionicons5";
import { useLogsStore } from "../stores/logs.js";
import { useImage2TextStore } from "../stores/image2text"
const logsStore = useLogsStore();
const image2textStore = useImage2TextStore()
const renderIcon = (icon) => {
    return () => h(NIcon, null, { default: () => h(icon) });
}

const switchLog = (value) => {
    logsStore.switchLog(value)
}

// 清空日志
const clearLog = () => {
    logsStore.clear()
}

onMounted(() => {
    // 日志事件
    EventsOn("logEvent", function (data) {
        logsStore.print(data)
    })

    // 上传图片事件
    EventsOn("uploadImageEvent", function (data) {
        image2textStore.appendPreview(data)
    })
})

const logs = ref("")

const open = ref(false)

watchEffect(() => {
    logs.value = logsStore.getLogs()
    open.value = logsStore.getOpen()
})

const menuOptions = [
    {
        label: () => h(
            RouterLink,
            {
                to: {
                    name: "welcome"
                }
            },
            { default: () => "首页" }
        ),
        key: "home",
        icon: renderIcon(HomeOutline),
    },
    {
        label: () => h(
            RouterLink,
            {
                to: {
                    name: "image-to-text"
                }
            },
            { default: () => "图生文" }
        ),
        key: "image-to-text",
        icon: renderIcon(ImageOutline)
    },
    {
        label: () => h(
            RouterLink,
            {
                to: {
                    name: "translate"
                }
            },
            { default: () => "百度翻译" }
        ),
        key: "translate",
        icon: renderIcon(ImageOutline)
    },
];

const activeKey = ref(null);
const collapsed = ref(false);
const defaultExpandedKeys = ["dance-dance-dance", "food"]

const handleUpdateExpandedKeys = (keys) => {
}
</script>
<style>
.n-layout-sider.n-layout-sider--show-content .n-layout-sider-scroll-container {
    height: 100vh !important;
}

.n-menu-item-content {
    padding-left: 17px !important;
}
.n-menu .n-menu-item-content.n-menu-item-content--selected:hover,.n-menu .n-menu-item-content.n-menu-item-content--selected::before{
    cursor: pointer;
    border: none;
    outline: none;
    border-radius: 50px;
    background: #ecf0f3;
    box-shadow: inset 4px 4px 4px #d1d9e6, inset -4px -4px 4px #f9f9f9;
}
</style>