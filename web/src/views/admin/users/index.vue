<script lang="ts" setup>
import {onBeforeMount, ref} from "vue";
import {UserApi} from "@/types/response";
import * as api from "@/api"
import {SizeToString} from "@/utils/size";
import {ElMessageBox} from "element-plus";
import NewUserDialog from "./components/NewUserDialog.vue"


const users = ref<UserApi[]>()

const dialogNewUserVisible = ref(false)


onBeforeMount(async () => {
  const {data} = await api.admin.getUsers()
   users.value = data
})

const sortLogin = (a: UserApi, b: UserApi) => {
  return a.login.localeCompare(b.login, undefined, {sensitivity: "base"})
}

const filterActive = (value: boolean, row: UserApi) => {
  return row.active == value
}

const formatStorageSize = (row: UserApi) => {
  const size = SizeToString(row.storage.size)
  let quota = "illimité"
  if (row.storage.activeQuota) {
    quota = SizeToString(row.storage.quota)
  }
  return `${size} / ${quota}`
}

</script>

<template>
  <div class="admin-settings-view">
    <div class="users-settings-view">
      <h2>Utilisateurs</h2>
      <el-divider/>
      <div class="settings-view">
        <div class="users-settings-actions">
          <el-button @click="dialogNewUserVisible = true">Nouvel Utilisateur</el-button>
        </div>

        <el-table
            class="user-list"
            :data="users"
            :default-sort="{ prop: 'login', order: 'ascending' }"
            style="width: 100%"
        >
          <el-table-column
              label="Nom"
              prop="login"
              :sort-orders="['ascending', 'descending']"
              :sort-method="sortLogin"
              sortable
          />

          <el-table-column
              prop="active"
              label="Status"
              width="120"
              :sort-orders="['ascending', 'descending']"
              sortable
              :filters="[
            { text: 'Activer', value: true },
            { text: 'Désactiver', value: false },
          ]"
              :filter-method="filterActive"
          >
            <template #default="scope">
              <el-tag
                  :type="scope.row.active? 'success' : 'danger'"
                  disable-transitions
              >{{ scope.row.active? 'Activer' : 'Désactiver' }}</el-tag
              >
            </template>
          </el-table-column>

          <el-table-column
              label="Stockage"
              prop="storage.size"
              width="150"
              :formatter="formatStorageSize"
          />
        </el-table>

        <NewUserDialog v-model="dialogNewUserVisible" @newUser="users?.push($event)"/>
      </div>
    </div>
  </div>
</template>

<style lang="scss">


</style>