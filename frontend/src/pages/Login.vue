<template>
  <div class="main">
    <div class="container b-container" id="b-container">
      <div class="form" id="b-form">
        <h2 class="form_title title">
          <n-gradient-text gradient="linear-gradient(90deg, #84fab0 0%, #8fd3f4 100%)">
            产品工具
          </n-gradient-text>
        </h2>
        <input class="form__input" type="text" v-model="user.email" placeholder="用户名">
        <input class="form__input" type="text" v-model="user.uid" placeholder="uid">
        <input class="form__input" type="text" v-model="user.cuid" placeholder="cuid">
        <button class="form__button button submit" @click="login">登录</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from "vue-router";
import { Config, SetUser } from "../../wailsjs/go/main/App"
import { LogPrint, EventsOn, EventsOff } from "../../wailsjs/runtime"
import { useMessage, useNotification, NInput, NImage, NSpin } from "naive-ui";
import { useUserStore } from "../stores/user"
import { useRoute, useRouter } from 'vue-router'

import { onMounted, onUnmounted, watchEffect } from 'vue'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const message = useMessage();

const user = ref({
  email: "",
  uid: "",
  cuid: "",
})

const login = async () => {
  userStore.setUser(user.value).then(res => {
    SetUser(user.value.email, user.value.uid, user.value.cuid).then((res) => {
      if (res.code == 0) {
        router.push('/admin/welcome')
      }
    })
  })
}

watchEffect(() => {

})

onUnmounted(() => {
  // EventsOff("queryConfigEvent")
})

onMounted(() => {
  // 从本地获取用户信息
  let user = userStore.getUser()

  user.value = user

  // Config().then(response => {
  //   console.log(response)
  //   if (response.code == 0) {
  //       // 遍历preview 如果id存在，则将数据追加到data中
  //       message.info(response.message)
  //   } else {
  //       message.error(response.message)
  //   }
  // })
  // EventsOn("queryConfigEvent", (res) => {
  //   email.value = res.user.email
  //   password.value = res.user.password
  // })
})

</script>

<style scoped>
*,
*::after,
*::before {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  user-select: none;
}

/* Generic */
body {
  width: 100%;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: "Montserrat", sans-serif;
  font-size: 12px;
  background-color: #ecf0f3;
  color: #a0a5a8;
}

/**/
.main {
  position: relative;
  width: 100%;
  height: 100vh;
  background-color: #ecf0f3;
  box-shadow: 10px 10px 10px #d1d9e6, -10px -10px 10px #f9f9f9;
  border-radius: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

/* @media (max-width: 1200px) {
  .main {
    transform: scale(0.7);
  }
}
@media (max-width: 1000px) {
  .main {
    transform: scale(0.6);
  }
}
@media (max-width: 800px) {
  .main {
    transform: scale(0.5);
  }
}
@media (max-width: 600px) {
  .main {
    transform: scale(0.4);
  }
} */

.container {
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  top: 0;
  width: 600px;
  height: 100%;
  padding: 25px;
  background-color: #ecf0f3;
  transition: 1.25s;
}

.form {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  width: 100%;
  height: 100%;
}

.form__icon {
  object-fit: contain;
  width: 30px;
  margin: 0 5px;
  opacity: 0.5;
  transition: 0.15s;
}

.form__icon:hover {
  opacity: 1;
  transition: 0.15s;
  cursor: pointer;
}

.form__input {
  width: 350px;
  height: 40px;
  margin: 4px 0;
  padding-left: 25px;
  font-size: 13px;
  letter-spacing: 0.15px;
  border: none;
  outline: none;
  font-family: "Montserrat", sans-serif;
  background-color: #ecf0f3;
  transition: 0.25s ease;
  border-radius: 8px;
  box-shadow: inset 2px 2px 4px #d1d9e6, inset -2px -2px 4px #f9f9f9;
}

.form__input:focus {
  box-shadow: inset 4px 4px 4px #d1d9e6, inset -4px -4px 4px #f9f9f9;
}

.form__span {
  margin-top: 30px;
  margin-bottom: 12px;
}

.form__link {
  color: #181818;
  font-size: 15px;
  margin-top: 25px;
  border-bottom: 1px solid #a0a5a8;
  line-height: 2;
}

.title {
  font-size: 34px;
  font-weight: 700;
  line-height: 3;
  color: #181818;
}

.description {
  font-size: 14px;
  letter-spacing: 0.25px;
  text-align: center;
  line-height: 1.6;
}

.button {
  cursor: pointer;
  width: 348px;
  height: 40px;
  border-radius: 25px;
  margin-top: 20px;
  font-weight: 700;
  font-size: 14px;
  letter-spacing: 1.15px;
  /* color: #f9f9f9; */
  border: none;
  outline: none;
  border-radius: 50px;
  background: #ecf0f3;
  box-shadow: inset 4px 4px 4px #d1d9e6, inset -4px -4px 4px #f9f9f9;
}

/**/
.a-container {
  z-index: 100;
}

.b-container {
  z-index: 0;
}

.switch {
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  width: 400px;
  padding: 50px;
  z-index: 200;
  transition: 1.25s;
  background-color: #ecf0f3;
  overflow: hidden;
  box-shadow: 4px 4px 10px #d1d9e6, -4px -4px 10px #f9f9f9;
}

.switch__circle {
  position: absolute;
  width: 500px;
  height: 500px;
  border-radius: 50%;
  background-color: #ecf0f3;
  box-shadow: inset 8px 8px 12px #d1d9e6, inset -8px -8px 12px #f9f9f9;
  bottom: -60%;
  left: -60%;
  transition: 1.25s;
}

.switch__circle--t {
  top: -30%;
  left: 60%;
  width: 300px;
  height: 300px;
}

.switch__container {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  position: absolute;
  width: 400px;
  padding: 50px 55px;
  transition: 1.25s;
}

.switch__button {
  cursor: pointer;
}

.switch__button:hover {
  box-shadow: 6px 6px 10px #d1d9e6, -6px -6px 10px #f9f9f9;
  transform: scale(0.985);
  transition: 0.25s;
}

.switch__button:active,
.switch__button:focus {
  box-shadow: 2px 2px 6px #d1d9e6, -2px -2px 6px #f9f9f9;
  transform: scale(0.97);
  transition: 0.25s;
}

/**/
.is-txr {
  left: calc(100% - 400px);
  transition: 1.25s;
  transform-origin: left;
}

.is-txl {
  left: 0;
  transition: 1.25s;
  transform-origin: right;
}

.is-z200 {
  z-index: 200;
  transition: 1.25s;
}

.is-hidden {
  visibility: hidden;
  opacity: 0;
  position: absolute;
  transition: 1.25s;
}

.is-gx {
  animation: is-gx 1.25s;
}

@keyframes is-gx {

  0%,
  10%,
  100% {
    width: 400px;
  }

  30%,
  50% {
    width: 500px;
  }
}
</style>