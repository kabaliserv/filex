<script lang="ts" setup>

import {ref} from "vue";

const emit = defineEmits<{
  (e: 'change', file: File | undefined): void
}>()

const dropZone = ref<HTMLDivElement>();
const fakeInputFile = ref<HTMLInputElement>()


const openBrowserFiles = () => {
  fakeInputFile.value?.click()
}

</script>

<template>
  <div ref="dropZone" class="dropzone" @click="openBrowserFiles">
    <div class="icons-wrapper">
      <div class="icons-content">
        <span class="iconify" data-icon="fa-solid:file-image"></span>
        <span class="iconify" data-icon="fa-solid:file-alt"></span>
        <span class="iconify" data-icon="fa-solid:file-pdf"></span>
      </div>
    </div>
    <div class="drop-zone-msg">
      <el-link type="primary">
        Cliquer ici
      </el-link>
      pour ajouter un fichier
    </div>
    <input
        type="file"
        class="fake-input-file"
        style="display: none"
        ref="fakeInputFile"
        @change="emit('change', $event.target.files[0])"
    />
  </div>
</template>

<style>
.dropzone {
  cursor: pointer;
  border: 4px dashed rgba(198, 198, 198, 0.65);
  border-radius: 30px;
  padding: 20px;
  box-sizing: border-box;
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
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

.drop-zone-msg {
  margin-top: 20px;

}
.drop-zone-msg span {
  font-size: 16px;
}

</style>
