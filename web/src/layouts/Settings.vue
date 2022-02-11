<script lang="ts" setup>

import {computed} from "vue";
import {useRoute} from "vue-router";
import SettingsNav from "./components/SettingsNav.vue"
import {UserModule} from "@/store/modules/user";

const route = useRoute()

const menuItems = [
  {
    name: "Profile",
    index: "profile",
    path: "/settings/profile"
  },
  {
    name: "Securité",
    index: "security",
    path: "/settings/security"
  },
]

const adminMenuItems = [
  {
    name: "Général",
    index: "admin",
    path: "/admin"
  },
  {
    name: "Utilisateurs",
    index: "users",
    path: "/admin/users"
  },
]

const activeMenu = computed(() => {
  const {path} = route
  const t = path.split("/")
  const index = t.pop()
  const page = t.pop()
  if (page == "settings")
    return index
})

const adminActiveMenu = computed(() => {
  const {path} = route
  const t = path.split("/")
  const index = t.pop()
  if (index == "admin") return index
  const page = t.pop()
  if (page == "admin")
    return index
})

</script>

<template>
  <div class="settings-page">
    <div class="page-content">
      <el-row :gutter="30">
        <el-col :span="6">
          <div class="side-nav">
            <SettingsNav title="Paramètres" :items="menuItems" :activeMenu="activeMenu" />
            <SettingsNav v-if="UserModule.admin" title="Admin" :items="adminMenuItems" :activeMenu="adminActiveMenu" />
          </div>
        </el-col>
        <el-col :span="18" style="flex: 1">
          <router-view></router-view>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<style lang="scss">
.settings-page {
  display: flex;
  justify-content: center;
}

.page-content {
  max-width: 1240px;
  width: 100%;
  padding: 30px 20px 0;
}

.side-nav {
  display: flex;
  flex-direction: column;
  gap: 30px;
}
</style>