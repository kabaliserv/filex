<script lang="ts" setup>
import DropZone from "./DropFileZone.vue";
import {ref} from "vue";
import FileView from "./FileView.vue";
import ExpireOption from "./ExpireOption.vue"

interface Props {
  file: File
  expire: number
}

const props = defineProps<Props>()
const optionsIsCollapse= ref(true)
const elOptions = ref<HTMLDivElement>()

const expirationListOption = [
  { id: 0, short: "10 Min", long: "10 Minutes" },
  { id: 1, short: "1H", long: "1 Heure" },
  { id: 2, short: "1J", long: "1 Jour" },
  { id: 3, short: "3J", long: "3 Jours" },
]

const emit = defineEmits<{
  (e: "update:file", file: File): void
  (e: "update:expire", value: number): void
  (e: "start-upload"): void
}>()

const expValue = ref(0)
const passwordOption = ref("")

const OnFileChange = (v: File | undefined) => {
  if (v) emit("update:file", v)
}

const toggleExpandOptions = () => {
  if (elOptions.value) {
    const el = elOptions.value
    if (el.style.maxHeight) {
      el.style.maxHeight = ""
    } else {
      el.style.maxHeight = `${el?.scrollHeight ?? 0}px`
    }
  }
}

const OnExpireChange = (id: number) => {
  let time = 0
  switch (id) {
    case 1:
      time = 60 * 60 * 1_000; // 1 Hour
      break;
    case 2:
      time = 24 * 60 * 60 * 1_000; // 1 Day
      break;
    case 3:
      time = 3 * 24 * 60 * 60 * 1_000 // 3 Day
    default:
      time = 10 * 60 * 1_000; // 10 Min
  }
  emit('update:expire', time)
}

</script>

<template>
  <div class="upload-card">
    <div class="wrapper">
      <FileView :file="props.file" @change="OnFileChange"/>
      <div class="options">
        <span class="title-options" @click="toggleExpandOptions">Options avancées<span class="iconify option-icon" data-icon="bx:bxs-down-arrow"></span></span>
        <div ref="elOptions" class="options-list">
          <ExpireOption :expListOpts="expirationListOption" v-model="expValue" @update:modelValue="OnExpireChange"/>
        </div>
      </div>
      <el-button type="primary" class="shared-btn" round @click="emit('start-upload')">Partager</el-button>
    </div>
  </div>
</template>

<style scoped>

.wrapper {
  /*border: 2px solid black;*/
  padding: 20px;
  border-radius: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.title-options {
  font-size: 1.4rem;
  cursor: pointer;
  user-select: none;
}

.option-icon {
  font-size: 10px;
  margin-left: 4px;
}

.options-list {
  padding: 0 10px;
  padding-top: 10px;
  transition: max-height 0.5s ease-out;
  max-height: 0;
  overflow-y: hidden;

}

.shared-btn {
  width: 100%;
  margin-top: 20px;
}

/*.collapse-transition {*/
/*  transition: height 0.5s linear;*/
/*  overflow: hidden;*/
/*}*/

.collapse {
  height: 0px;
}

</style>