<script lang="ts" setup>
import {ref, computed, PropType, defineComponent, onBeforeMount, watch} from 'vue'
import {FileApi} from "@/types/file";
import FileDrawer from "./FileDrawer.vue"
import {SizeToString} from "@/utils/size";
import { ElMessageBox } from 'element-plus'
import {AppModule} from "@/store/modules/app";
import {FileExtension} from "@/utils/file/extentions";

interface Props {
  files: FileApi[]
}

type FileRow = FileApi & {
  icon?: {
    type?: string
    text: string
  }
}



const props = defineProps({
  files: Array as PropType<FileApi[]>
})

const emit = defineEmits<{
  (e: "delete", file: FileApi): void
}>()

const files = ref<FileRow[]>()

const drawer = ref(false)
const fileDrawer = ref<FileApi>()


const search = ref('')
const filterTableData = computed(() =>
    props.files?.filter(
        (data) =>
            !search.value ||
            data.name.toLowerCase().includes(search.value.toLowerCase())
    )
)
const handleEdit = (index: number, row: FileApi) => {
  console.log(index, row)
}
const handleDelete = (file: FileApi) => {
  emit("delete", file)
}

const handleRowClick = (row: FileApi, column: any, event: any) => {
  fileDrawer.value = row
  drawer.value = true
}

const sortName = (a: FileApi, b: FileApi) => {
  return a.name.localeCompare(b.name, undefined, {sensitivity: "base"})
}

const formatSize = (row: FileApi) => {
  return SizeToString(row.size)
}

const formatDate = (row: FileApi) => {
  const date = new Date(row.created_at)
  const day = date.getDate()
  const month = date.getMonth() + 1
  const year = date.getFullYear().toString().split("").slice(2).join("")
  const hour = date.getHours()
  const min = date.getMinutes()
  return `${day < 10 ? "0".concat(day.toString()): day}-${month < 10 ? "0".concat(month.toString()) : month}-${year} à ${hour < 10 ? "0".concat(hour.toString()): hour}:${min < 10 ? "0".concat(min.toString()): min}`
}

const makeFileRow = (file: FileApi): FileRow => {
  const nameSlice = file.name.split(".")
  const ext = nameSlice.pop()
  const newFile: FileRow = {
    ...file,
    icon: {
      text: ""
    }
  }
  const extObj = FileExtension.find((item) => item.name == ext)

  if (extObj) {
    if (extObj.icon) newFile.icon!.type = extObj.icon
    else
    newFile.icon!.text = ext!.toUpperCase()

  }

  return newFile
}

const makeFilesRow = (filesSlice: FileApi[]): FileRow[] => {
  return filesSlice.map((item) => makeFileRow(item))
}

onBeforeMount(() => {
  if (props.files)
    files.value = makeFilesRow(props.files)
})

watch(() => props.files, (value) => {
  if (value)
    files.value = makeFilesRow(value)
}, {deep: true})

</script>

<template>
  <div>

  <el-table class="file-list"
            :data="files"
            :default-sort="{ prop: 'name', order: 'ascending' }"
            style="width: 100%"
            @row-click="handleRowClick" >
    <el-table-column
        label="Nom"
        prop="name"
        class-name="filename-col"
        :sort-orders="['ascending', 'descending']"
        :sort-method="sortName"
        sortable
    >
      <template #default="scope">
        <span class="file-icon">
            <span v-if="scope.row.icon && scope.row.icon.type" class="iconify" :data-icon="scope.row.icon.type"></span>
            <span v-else class="iconify" data-icon="bx:bxs-file-blank"></span>
            <div v-if="scope.row.icon.text" class="icon-text">{{scope.row.icon.text}}</div>
        </span>
        <span>{{scope.row.name}}</span>
      </template>
    </el-table-column>

    <el-table-column
        label="Taille"
        prop="size"
        :width="100"
        :formatter="formatSize"
        :sort-orders="['ascending', 'descending']"
        sortable
    />

    <el-table-column
        v-if="!AppModule.IsMobile"
        label="Ajouter le"
        prop="created_at"
        width="180"
        :formatter="formatDate"
        :sort-orders="['ascending', 'descending']"
        sortable
    />
  </el-table>

  <FileDrawer
      v-model="drawer"
      :file="fileDrawer"
      @delete="emit('delete', $event)"
  />
  </div>
</template>

<style lang="scss">

.file-list {
  td.filename-col div.cell {
    display: flex;
    flex-direction: row;
    gap: 5px;
    align-items: center;
    font-weight: bold;

    .file-icon {
      position: relative;
      font-size: 40px;
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
  .el-table__row {
    cursor: pointer;
  }
}


</style>