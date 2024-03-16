<template>
    <n-layout has-sider>
        <n-layout-sider bordered collapse-mode="width" :collapsed-width="58" :width="160" :height="1000"
            :collapsed="collapsed" show-trigger @collapse="collapsed = true" @expand="collapsed = false">
            <n-menu :options="menuOptions" :default-expanded-keys="defaultExpandedKeys"
                @update:expanded-keys="handleUpdateExpandedKeys" />
        </n-layout-sider>
        <n-layout class="h-screen">
            <router-view />
        </n-layout>
    </n-layout>

</template>

<script>
import { defineComponent, h, ref } from "vue";
import { NIcon, useMessage } from "naive-ui";
import { RouterLink } from "vue-router";

import {
    DocumentTextSharp,
    DocumentOutline,
    HomeOutline,
    ImageOutline,
} from "@vicons/ionicons5";
function renderIcon(icon) {
    return () => h(NIcon, null, { default: () => h(icon) });
}

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
                    name: "llm"
                }
            },
            { default: () => "文心一言" }
        ),
        key: "llm",
        icon: renderIcon(DocumentOutline)
    }
];

export default defineComponent({
    setup() {
        const message = useMessage();
        return {
            menuOptions,
            activeKey: ref(null),
            collapsed: ref(false),
            defaultExpandedKeys: ["dance-dance-dance", "food"],
            handleUpdateExpandedKeys(keys) {
            }
        };
    }
});
</script>
<style>
.n-layout-sider.n-layout-sider--show-content .n-layout-sider-scroll-container {
    height: 100vh !important;
}
.n-menu-item-content{
    padding-left: 17px !important;
}
</style>