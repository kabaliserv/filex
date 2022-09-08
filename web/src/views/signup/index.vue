<script lang="ts" setup>
import {reactive, ref, watch} from "vue";
import {UserModule} from "@/store/modules/user";
import { auth } from "@/api"
import type { ElForm } from 'element-plus'
import {LocationQuery, RouteLocationNormalizedLoaded, useRoute, useRouter} from "vue-router";

const router = useRouter();
const route = useRoute();

const form = reactive({
  login: "",
  email: "",
  password: "",
  confirmPassword: "",
})

const elFormRef = ref<InstanceType<typeof ElForm>>()

const validateLogin = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error("S'il vous plaît entrez votre identifiant"))
  } else {
    callback()
  }
}

const validatePass = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error("S'il vous plaît entrez votre mot de passe"))
  } else {
    callback()
  }
}

const SubmitForm = async () => {
  if (!elFormRef.value) return
  elFormRef.value.validate(async (valid: any) => {
    if (valid) {
      try{
        await auth.signup({login: form.login, email: form.email, password: form.password})
        await UserModule.Login({login: form.login, password: form.password})
        await router
            .push({
              path: "/",
            })
            .catch(console.warn);
      } catch (e) {
        console.log(e)
      }
    } else {
      return false
    }
  })
}

const rules = reactive({
  login: [{ required: true, validator: validateLogin, trigger: 'blur' }],
  email: [{required: true}],
  password: [{ required: true, validator: validatePass, trigger: 'blur' }],
  confirmPassword: [{required: true}],
})

</script>

<template>
  <div class="signup-view">
    <div class="view-wrapper">
      <h1 class="title-page">S'inscrire</h1>
      <div class="form-wrapper">
        <el-form
            ref="elFormRef"
            :model="form"
            :rules="rules"
            label-position="top"
            label-width="120px"
            @submit.prevent="SubmitForm"
        >
          <el-form-item label="Nom d'utilisateur" prop="login">
            <el-input
                v-model="form.login"
                type="text"
                autocomplete="username"
            ></el-input>
          </el-form-item>
          <el-form-item label="Email" prop="email">
            <el-input
                v-model="form.email"
                type="text"
                autocomplete="email"
            ></el-input>
          </el-form-item>
          <el-form-item label="Mot de passe" prop="password">
            <el-input
                v-model="form.password"
                type="password"
                autocomplete="new-password"
            ></el-input>
          </el-form-item>
          <el-form-item label="Confirmer mot de passe" prop="confirmPassword">
            <el-input
                v-model="form.confirmPassword"
                type="password"
                autocomplete="new-password"
            ></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" native-type="submit">Connection</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.signup-view {
  display: flex;
  position: relative;
  flex-direction: column;
  justify-content: center;
  padding: 20px;
  box-sizing: border-box;
}

.title-page {
  position: absolute;
  width: 100%;
  top: -10vh;
  text-align: center;
}

.view-wrapper {
  position: relative;
  max-width: 500px;
  width: 100%;
  margin: 0 auto;

}

.form-wrapper {
}

</style>
