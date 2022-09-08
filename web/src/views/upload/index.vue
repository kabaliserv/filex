<script lang="ts" setup>
import {reactive, ref} from "vue";
import * as tus from "tus-js-client";
import UploadCard from "./components/UploadCard.vue";
import DropFileZone from "@/views/upload/components/DropFileZone.vue";
import { SendFile } from "./utils/file";
import { sleep } from "@/utils";
import LoadingSpinner from "vue-simple-spinner/src/components/Spinner.vue"

// data
const file = ref<File>()
const options = reactive({
  expire: 0,
  password: ""
})
const uploadState = reactive({
  inProgress: false,
  isDone: false,
  showProgressBar: false,
  progress: 0
})
const showLoading = ref(false)
const uploadURL = ref("")

// functions
const StartUpload = async () => {
  if (!file.value || uploadState.inProgress) return
  showLoading.value = true
  uploadState.inProgress = true
  await sleep(1000)
  showLoading.value = false
  await sleep(100)
  await SendFile(file.value, {
    duration: options.expire,
    password: options.password,
    onProgress: (bytesSent: number, bytesTotal: number) => {
      const percentage = ((bytesSent / bytesTotal) * 100).toFixed(0)
      uploadState.progress = +percentage
    },
    onSuccess: async (url) => {
      uploadURL.value = url
      showLoading.value = true
      uploadState.inProgress = false
      uploadState.isDone = true
      await sleep(1000)
      showLoading.value = false
    }
  });

}

const OnFileChange = (value: File | undefined) => {
  if (value) {
    file.value = value
  }
}

</script>

<template>
  <div class="upload-page">
    <h1 class="title-page" v-show="!showLoading">Partagez un fichier</h1>
    <section class="upload-view" v-show="!showLoading">
      <div v-if="!uploadState.inProgress && !uploadState.isDone">
        <div v-if="!file"  class="drop-file-zone" >
          <DropFileZone @change="OnFileChange"/>
        </div>
        <UploadCard v-else v-model:file="file" v-model:expire="options.expire" @start-upload="StartUpload"/>
      </div>
      <div v-else-if="uploadState.inProgress && !uploadState.isDone">
        <el-progress :text-inside="true" :stroke-width="26" :percentage="uploadState.progress" />
      </div>
      <div v-else-if="!uploadState.inProgress && uploadState.isDone">
        <el-link :href="uploadURL">{{uploadURL}}</el-link>
      </div>
    </section>
    <LoadingSpinner v-if="showLoading" size="large" message="Chargement..."/>
  </div>
</template>

<style scoped>

.upload-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
}

.upload-view {
  max-width: 400px;
  width: 100%;
}

.title-page {
  /*position: absolute;*/
  /*top: 20vh;*/
  margin-bottom: 50px;
}

.drop-file-zone {
  /*height: 400px;*/
  background-color: #eeeeee;
  padding: 10px;
  border-radius: 30px;

}

.drop-file {
  cursor: pointer;
}

.drop-border {
  border: 4px dashed rgba(198, 198, 198, 0.65);
     border-radius: 30px;
  padding: 20px;
}

.icons-content {
  margin: 40px 20px 20px;
  color: #95afc0;
  opacity: 0.70;
  font-size: 50px;

  display: flex;
  gap: 70px;
  justify-content: center;
}

.icons-content > svg:nth-child(1) {
  transform: rotate(-45deg);
}

.icons-content > svg:nth-child(2) {
  transform: translateY(-40%);
}

.icons-content > svg:nth-child(3) {
  transform: rotate(45deg);
}

</style>