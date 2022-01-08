<script lang="ts" setup>
import * as tus from "tus-js-client";
import { ref } from "vue";
const percent = ref("");

const onInputChange = async (e: Event) => {
  // Get the selected file from the input element
  let target = e.target as HTMLInputElement;

  if (!target || !target.files) return;

  let file = target.files[0];

  const headers: {[p: string]: string} = {}
  const token = await fetch("/api/auth/upload", { method: "Post" }).then(res => {
    if (res.ok) return res.text()
    return ""
  })

  console.log("token: ", token, token.length)
  if (token.length > 0) {
    headers["Authorization"] = `Bearer ${token}`
  }

  // Create a new tus uploads
  let upload = new tus.Upload(file, {
    endpoint: "/api/files/",
    retryDelays: [0, 3000, 5000, 10000, 20000],
    headers,
    metadata: {
      filename: file.name,
      filetype: file.type,
    },
    onError: function (error) {
      console.log("Failed because: " + error);
    },
    onProgress: function (bytesUploaded, bytesTotal) {
      let percentage = ((bytesUploaded / bytesTotal) * 100).toFixed(2);
      console.log(bytesUploaded, bytesTotal, percentage + "%");
      percent.value = percentage;
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

  // Check if there are any previous uploads to continue.
  // uploads.findPreviousUploads().then(function (previousUploads) {
  //   // Found previous uploads so we select the first one.
  //   if (previousUploads.length) {
  //     uploads.resumeFromPreviousUpload(previousUploads[0]);
  //   }
  //
  //   // Start the uploads
  //   uploads.start();
  // });
  upload.start();
};

const username = ref("wilson")
const password = ref("123456789")

const onSubmit = async () => {
  const body = JSON.stringify({
    username: username.value,
    password: password.value,
  })
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body,
  })

  console.log("Login Request: ", res)
}
</script>

<template>
  <div>
    <div>
      <input type="file" @change="onInputChange" />
      <span>{{ percent }}%</span>
    </div>
    <div>
      <form action="" @submit.prevent="onSubmit">
        <input type="text" :value="username">
        <input type="text" :value="password">
        <button type="submit">Login</button>
      </form>
    </div>
  </div>
</template>
