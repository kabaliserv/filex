<script lang="ts" setup>

import {FileApi} from "@/types/file";
import {ref, watch} from "vue";

import {SizeToString} from "@/utils/size";
import {ElMessageBox} from "element-plus/es";
import {AppModule} from "@/store/modules/app";

type FileRow = FileApi & {
  icon?: {
    type?: string
    text: string
  }
}

interface Props {
  file?: FileRow
  modelValue: boolean
}

const props = withDefaults(defineProps<Props>(),{
  modelValue: false
})

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void
  (e: "delete", file: FileRow): void
}>()

const drawer = ref(false)

const sharedFile = () => {

}

const renameFile= () => {

}

const deleteFile = () => {
  console.log('click')
  ElMessageBox.confirm(`Êtes vous sur de vouloir supprimer ce fichier ?`)
      .then(() => {
        emit("delete", props.file!)
        drawer.value = false
      })
      .catch(() => {
        // catch error
      })
}


watch(() => props.modelValue, (value: boolean) => {
  drawer.value = value
})

</script>

<template>
  <el-drawer v-model="drawer"
             title="Fichier"
             :direction="AppModule.IsMobile ? 'btt': 'rtl'"
             :size="AppModule.IsMobile ? '70%': '30%'"
             @open="emit('update:modelValue', true)"
             @close="emit('update:modelValue', false)"
  >
    <div v-if="props.file && drawer" class="file-drawer-view">
      <div class="file-icon">
        <span class="icon-wrapper">
          <span v-if="props.file.icon && props.file.icon.type" class="iconify" :data-icon="props.file.icon.type"></span>
          <span v-else class="iconify" data-icon="bx:bxs-file-blank"></span>
          <div v-if="props.file.icon.text" class="icon-text">{{props.file.icon.text}}</div>
        </span>
      </div>
      <div class="file-info">
        <span class="filename">{{props.file?.name}}</span>
        <div class="file-props">
          <span class="filesize">Taille: {{SizeToString(props.file?.size)}}</span>
        </div>
      </div>
      <div class="file-actions">
        <div v-if="false" @click="sharedFile">
          <span class="iconify share-icon" data-icon="bi:share-fill"></span>
        </div>
        <div v-if="false" @click="renameFile">
          <span class="iconify rename-icon" data-icon="ic:sharp-drive-file-rename-outline"></span>
        </div>
        <div @click="deleteFile" title="Supprimer">
          <span class="iconify delete-icon" data-icon="bi:trash-fill"></span>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<style lang="scss" scoped>

.file-info {
  display: flex;
  flex-direction: column;
  gap: 10px;

  .filename {
    font-size: 20px;
    font-weight: bold;
  }
}

.file-actions {
  display: flex;
  justify-content: space-around;
  margin-top: 50px;
  font-size: 30px;
  
  &  svg {
    cursor: pointer;
    color: blue;
  }

  &  .delete-icon {
    color: red;
  }
}

.file-icon {
  font-size: 200px;
  display: flex;
  justify-content: center;

  .icon-wrapper {
    position: relative;
    .icon-text {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -40%);
      z-index: 1;
      font-size: 9px;
      word-break: break-word;
      overflow-x: visible;
      height: 15px;
      width: 100%;
      text-align: center;
      color: white;

    }
  }
}

</style>