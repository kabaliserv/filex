<script lang="ts" setup>

import {useRoute, useRouter} from "vue-router";
import {UserModule} from "@/store/modules/user";

const router = useRouter()
const route = useRoute()

const GotoLogin = () => {
  if (route.path != "/login")
    router.push({path: "/login"})
}

const GotoSignUp = () => {
  if (route.path != "/signup")
    router.push({path: "/signup"})
}

const Logout = async () => {
  await UserModule.Logout()
  await router.push({path: "/login"})
}

</script>

<template>
  <header>
    <div class="header-wrapper">
      <div class="logo">
        <a href="/">
          <span>fileX</span>
        </a>
      </div>
      <nav>
        <el-link v-if="!UserModule.auth" :underline="false" @click="GotoLogin">Se Connecter</el-link>
        <el-link v-if="!UserModule.auth" :underline="false" @click="GotoSignUp">S'inscrire</el-link>

        <el-popover v-if="UserModule.auth" placement="bottom" trigger="click">
          <template #reference>
            <div v-if="UserModule.auth" class="user-item">
              <span>{{UserModule.login}}</span>
              <el-avatar v-if="UserModule.avatar" :src="UserModule.avatar" :size="30"></el-avatar>
              <el-avatar v-else :size="30">
                <span class="iconify" data-icon="fa:user"></span>
              </el-avatar>
            </div>
          </template>
          <router-link to="/settings">Paramètre</router-link>
          <el-link :underline="false" @click="Logout">Se Déconnecter</el-link>
        </el-popover>
      </nav>
    </div>
  </header>
</template>

<style lang="scss" scoped>
header {
  border-bottom: 1px solid #e6e6e6;
  background-color: #99c4ff;

}
.header-wrapper {
  display: flex;
  max-width: 1240px;
  margin: 0 auto;
  justify-content: space-between;
  align-items: center;
  height: 60px;
  padding: 0 30px;
}
nav {
  display: flex;
  gap: 10px;
}

.user-item {
  display: flex;
  gap: 10px;
  align-items: center;
  cursor: pointer;
  border-radius: 7px;
  padding: 5px 10px;
  &:hover {
    background-color: rgba(255, 255, 255, 0.5);
    transition: background-color 0.2s  ease-out;
  }
}

</style>
