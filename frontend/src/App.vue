<script setup>
// import HelloWorld from './components/HelloWorld.vue';
import hljs from 'highlight.js/lib/core'
import javascript from 'highlight.js/lib/languages/javascript'
hljs.registerLanguage('javascript', javascript)

</script>

<template>
  <n-config-provider :hljs="hljs">
    <n-notification-provider>
      <n-message-provider>
        <n-watermark content="" cross selectable :font-size="16" :line-height="16" :width="192" :height="128"
          :x-offset="12" :y-offset="28" :rotate="-15">
          <router-view />
        </n-watermark>
      </n-message-provider>
    </n-notification-provider>
  </n-config-provider>

</template>

<style></style>
