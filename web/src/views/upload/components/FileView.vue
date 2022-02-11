<script lang="ts" setup>
import {PropType, ref} from "vue";
import {SizeToString} from "@/utils/size"

const props = defineProps({
  file: Object as PropType<File>
})

const emit = defineEmits<{
  (e: 'change', file: File | undefined): void
}>()

const fakeInputFile = ref<HTMLInputElement>()


const openBrowserFiles = () => {
  fakeInputFile.value?.click()
}

</script>
<template>
  <div class="file-view">
    <div class="view-header">
      <span class="title">Fichier:</span>
      <el-link :underline="false" @click="openBrowserFiles">changer de fichier</el-link>
    </div>
    <div class="view-body">
      <span class="filename" :title="props.file.name">{{props.file.name}}</span>
      <span class="filesize">{{ SizeToString(props.file.size)}}</span>
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

<style scoped>
.view-header {
  display: flex;
  justify-content: space-between;
  /*margin-bottom: 10px;*/
}

.view-header .title {
  font-weight: bold;
}

.view-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-start;
  padding: 10px 0;
  margin: 10px 0;

}

.view-body > * {
  opacity: 0.7;
  font-weight: bold;
}

.filename {

}
</style>
