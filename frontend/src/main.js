import {createApp} from 'vue';

import App from './App.vue';
// 通用字体
import 'vfonts/Lato.css';
// 等宽字体
import 'vfonts/FiraCode.css';
// 导入css样式
import "./assets/css/style.css"

import routes from './routers/router'

import naive from 'naive-ui';

import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
    history: createWebHistory(),
    routes,
})

const app = createApp(App)
app.use(naive)
app.use(router)
app.mount('#app')