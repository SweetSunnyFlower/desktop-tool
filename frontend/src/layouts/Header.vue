<template>
<n-layout has-sider>
    <n-layout-sider
        bordered
        collapse-mode="width"
        :collapsed-width="64"
        :width="240"
        :height="1000"
        :collapsed="collapsed"
        show-trigger
        @collapse="collapsed = true"
        @expand="collapsed = false"
      >
    <n-menu
        :options="menuOptions"
        :default-expanded-keys="defaultExpandedKeys"
        @update:expanded-keys="handleUpdateExpandedKeys"
    />
    </n-layout-sider>
    <n-layout>
        <router-view/>
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
    HomeOutline
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
        label: "Word",
        key: "word",
        icon: renderIcon(DocumentTextSharp),
        children:[
            {
                label: () => h(
                RouterLink,
                {
                    to: {
                        name: "replace"
                    }
                },
                { default: () => "批量替换文本" }
                ),
                key: "word-replace",
                icon: renderIcon(DocumentOutline)
            }
        ]
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
.n-layout-sider.n-layout-sider--show-content .n-layout-sider-scroll-container{
    height: 100vh!important;
}
</style>