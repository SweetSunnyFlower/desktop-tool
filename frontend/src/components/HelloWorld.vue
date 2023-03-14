<script setup>
import {reactive} from 'vue'
import { Replace} from '../../wailsjs/go/main/App'

const data = reactive({
  input: "",
  output:"",
  replaceData:[],
})

function replace() {
  Replace(data.input,data.output, data.replaceData).then(result => {
    console.log(result)
  })
}

function addReplace() {
  data.replaceData.push(["",""])
}

function delReplace(index) {
  data.replaceData.splice(index, 1)
}

</script>

<template>
  <main>
    <div>
      <div id="input" class="input-box">
        <input v-model="data.input" placeholder="输入文件夹" autocomplete="off" class="input" type="text"/>
        <input v-model="data.output" placeholder="输出文件夹" autocomplete="off" class="input" type="text"/>
        <button class="btn" @click="replace">替换</button>
      </div>
      <div class="input-box" v-for="(re,index) in data.replaceData" :key="index">
        <input v-model="re[0]" placeholder="Old" autocomplete="off" class="input" type="text"/>
        <input v-model="re[1]" placeholder="New" autocomplete="off" class="input" type="text"/>
        <button class="btn" @click="delReplace(index)">删除</button>
      </div>
    </div>
    <div>
      <button class="btn" @click="addReplace">添加替换词</button>
    </div>
  </main>
</template>

<style scoped>
main{
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}
.btn{
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}
.input-box .btn {
  /* width: 60px; */
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  margin: 10px 10px;
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
