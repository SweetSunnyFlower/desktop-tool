<script setup>
// import HelloWorld from './components/HelloWorld.vue';
</script>

<template>
  <n-notification-provider>
    <n-message-provider>
      <n-watermark content="" cross selectable :font-size="16" :line-height="16" :width="192" :height="128"
        :x-offset="12" :y-offset="28" :rotate="-15">
        <router-view />
      </n-watermark>
    </n-message-provider>
  </n-notification-provider>
</template>

<style></style>
