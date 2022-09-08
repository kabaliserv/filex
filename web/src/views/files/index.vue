<script lang="ts" setup>
import {onBeforeMount, ref} from "vue";
import {FileApi} from "@/types/file";
import * as api from "@/api";
import FileList from "./components/FileList.vue";
import SideBar from "./components/SideBar.vue";
import { ArrowRight, HomeFilled, Plus } from '@element-plus/icons-vue'
import { ElNotification } from 'element-plus'
import "element-plus/es/components/notification/style/css.mjs"

const files = ref<FileApi[]>()
const filesIsFetch = ref(false)

const inputFileAsRef = ref<HTMLInputElement>()

onBeforeMount(async () => {
  const {data} = await api.files.geUserFiles()
  files.value = data
  filesIsFetch.value = true
})


const OnInputFileChange = (event : Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.item(0)

  if (file) {
    const upload = api.files.addFile(file, {
      onSuccess: (file) => {
        files.value?.push(file)
      }
    })
  }

  target.files = null
}


const handleDeleteFile = async (file: FileApi) => {
  const fileIndex = files.value?.findIndex((item) => item.id == file.id)
  if (!fileIndex) return
  const filesRemoved = files.value?.splice(fileIndex, 1)
  if (!filesRemoved || filesRemoved.length == 0) return;
  try {
    await api.files.delFile(file.id)
  } catch (e) {
    ElNotification({
      title: 'Error',
      message: 'Une erreur est survenu lors de la suppression du fichier. Veuillez contacter l\'administrateur si le problème persiste',
      type: 'error',
      offset: 100,
    })
    files.value?.push(file)
  }
}

</script>

<template>
  <div class="files-page">
    <div class="file-view__content-wrapper">
      <div class="file-view__header">
        <div class="header-breadcrumb">
          <el-breadcrumb :separator-icon="ArrowRight" separator-class="nav-separator-icon">
            <el-breadcrumb-item :to="{ path: '/' }">
              <el-icon class="home-nav-icon"><home-filled /></el-icon>
            </el-breadcrumb-item>
            <el-breadcrumb-item :to="{ path: '/' }">
              <el-icon class="new-file-nav-icon" @click="inputFileAsRef?.click()"><plus /></el-icon>
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
      </div>
      <FileList class="file-view__list-file" :files="files" @delete="handleDeleteFile"></FileList>
    </div>
    <input ref="inputFileAsRef" type="file" style="display: none" @change="OnInputFileChange">
  </div>
</template>

<style lang="scss">


.header-breadcrumb {
  padding: 10px;

  .el-breadcrumb__item {
    display: flex;
    align-items: center;
  }



  .home-nav-icon, .new-file-nav-icon {
    font-size: 20px;
    border-radius: 50%;
    padding: 5px;
    cursor: pointer;
    &:hover {
      background-color: #eaeaea;
    }
  }

  .new-file-nav-icon {
    font-size: 17px;
    background-color: #e0e0e0;
  }
}
</style>

<style lang="scss" scoped>
.files-page {
  display: flex;
  height: 100%;
}

.file-view__sidebar {
  height: 100%;
}

.file-view__content-wrapper {
  flex: 1;
  height: 100%;
  max-width: 1200px;
  margin: auto;
  border-left: 1px solid #ebeef5;
  border-right: 1px solid #ebeef5;
}

</style>
