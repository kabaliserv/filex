<script lang="ts" setup>

import {useRoute, useRouter} from "vue-router";
import {onBeforeMount, ref} from "vue";
import {ApiUpload} from "@/types/api";
import {uploads} from "@/api"
import {SizeToString} from "@/utils/size";

const route = useRoute()
const router = useRouter()
let uploadId = ""
let token = ""

{
  const u = route.query["u"]
  if (typeof u == "string") {
    uploadId = u
  }
  const t = route.query["t"]
  if (typeof t == "string") {
    token = t
  }
}

const uploadData = ref<ApiUpload>()
const uploadRequestAuth = ref(false)

onBeforeMount(async () => {
  const header: Record<string, string> = {}
  if (token.length > 0) {
    header["Authorization"] = `Bearer ${token}`
  }
  const res = await uploads.findById(uploadId)
  if (res.status == 401) {
    await router.push({
      path: "/login",
    })
  }
  else if (res.status == 200) uploadData.value = res.data
})

const DownloadFile = async () => {
  const res = await uploads.requestDownload(uploadId)
  if (res.status == 200) {
    const a = document.createElement("a")
    a.href = res.data.url
    a.click()
  }
}

</script>

<template>
  <div class="download-page">
    <div v-if="uploadData"  class="download-view">
      <div class="download-wrapper">
        <h1 class="title-page">Partage de fichiers</h1>
        <div class="main-content">
          <div class="content-header">Fichier:</div>
          <div class="metadata-file">
            <div class="filename">{{ uploadData?.file.name }}</div>
            <div class="filesize">{{ SizeToString(uploadData?.file.size) }}</div>
          </div>
          <div class="file-actions">
            <el-button type="primary" round style="width: 100%; padding: 20px 0; font-size: 18px" @click="DownloadFile">Télécharger</el-button>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <h1>404</h1>
      <span>Not Found</span>
    </div>
  </div>
</template>

<style scoped>

.download-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
}

.download-view {
  max-width: 400px;
  width: 100%;
}
.title-page {
  /*position: absolute;*/
  /*top: 20vh;*/
  margin-bottom: 50px;
  text-align: center;
}

.content-header {
  font-size: 20px;
  margin-bottom: 20px;
}

.metadata-file {
  padding: 0 10px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

</style>