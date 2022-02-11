<script lang="ts" setup>
import {ref} from "vue";
import * as tus from "tus-js-client";

const file = ref<File>()
const filename = ref("")
const size = ref(0)

const uploadInProgress = ref(false)
const uploadPercent = ref(0)
const uploadFinish = ref(false)

const uploadURL = ref("")

// functions
const onInputFileChange = (e: Event) => {
  let target = e.target as HTMLInputElement;
  if (!target || !target.files) return;
  file.value = target.files[0];
}

const SendFile = async () => {
  if (!file.value) return
  uploadInProgress.value = true

  const headers: {[p: string]: string} = {}
  const access: {id?: string} = await fetch("/api/files/request-upload", { method: "Post" }).then(res => {
    if (res.ok) return res.json()
    return {}
  })

  if (!access.id) return

  headers["Filex-Upload-Access-Id"] = access.id

  const metadata: Record<string, string> = {
    filename: file.value.name,
    filetype: file.value.type,
  }

  let upload = new tus.Upload(file.value, {
    endpoint: "/api/files/",
    retryDelays: [0, 3000, 5000, 10000, 20000],
    headers,
    metadata,
    onError: function (error) {
      console.log("Failed because: " + error);
    },
    onProgress: function (bytesUploaded, bytesTotal) {
      let percentage = ((bytesUploaded / bytesTotal) * 100).toFixed(2);
      console.log(bytesUploaded, bytesTotal, percentage + "%");
      uploadPercent.value = +percentage;
    },
    onSuccess: function() {
      console.log("success");
      console.log(
          "Download %s from %s",
          upload.file instanceof File ? upload.file.name : "",
          upload.url
      );
    },
    onAfterResponse: function (req, res) {
      let url = req.getURL()
      let value = res.getHeader("X-My-Header")
      console.log(upload.file)
      console.log(res)
      console.log(`Request for ${url} responded with ${value}`)
    },
  });
  upload.start()
}

</script>

<template>
  <div>
    <h1>Upload Page</h1>
    <div class="upload-view">
      <div>
        <div v-if="file">
          <div>Fichier :</div>
          <div>{{ file?.name }}</div>
          <div>{{ file?.size }}</div>
        </div>
        <div class="input-file">
          <form action="">
            <input type="file" name="file" id="upload-file" @change="onInputFileChange">
            <button @click.prevent="SendFile">Envoyer</button>
          </form>
        </div>
        <div v-if="uploadInProgress">
          <span>{{uploadPercent}} %</span>
        </div>
        <div v-if="uploadFinish">
          <span>{{uploadURL}}</span>
        </div>
      </div>
    </div>
  </div>
</template>