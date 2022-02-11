<script lang="ts" setup>

import {reactive, ref, watch} from "vue";
import {UserModule} from "@/store/modules/user";
import type { ElForm } from 'element-plus'
import {LocationQuery, RouteLocationNormalizedLoaded, useRoute, useRouter} from "vue-router";

const router = useRouter();
const route = useRoute();

const form = reactive({
  login: "",
  password: "",
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


let redirect = "/";
const otherQuery: LocationQuery = {};



const SubmitForm = async () => {
  if (!elFormRef.value) return
  elFormRef.value.validate(async (valid: any) => {
    console.log(valid)
    if (valid) {
      try{
        await UserModule.Login({login: form.login, password: form.password})
        await router
            .push({
              path: redirect,
              query: otherQuery,
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
  password: [{ required: true, validator: validatePass, trigger: 'blur' }],
})

const getOtherQuery = (query: LocationQuery) => {
  if (typeof query.q === "string") {
    try {
      return JSON.parse(atob(query.q)) as LocationQuery;
    } catch (e) {
      console.warn(e);
      return {} as LocationQuery;
    }
  }
  return {} as LocationQuery;
};

const getRedirection = (route: RouteLocationNormalizedLoaded) => {
  const query = route.query;

  if (typeof query.redirect === "string") redirect = query.redirect;

  Object.assign(otherQuery, getOtherQuery(query));
  console.log(otherQuery);
};

getRedirection(route);

watch(() => route, getRedirection);

</script>

<template>
  <div class="login-page">
    <h1>Login Page</h1>
    <el-form
        class="login-form"
        ref="elFormRef"
        :model="form"
        :rules="rules"
        label-position="top"
        label-width="120px"
        @submit.prevent="SubmitForm"
    >
      <el-form-item label="Nom d'utilisateur / Email" prop="login">
        <el-input
            v-model="form.login"
            type="text"
            autocomplete="username"
        ></el-input>
      </el-form-item>
      <el-form-item label="Password" prop="password">
        <el-input
            v-model="form.password"
            type="password"
            autocomplete="password"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" native-type="submit">Connection</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<style lang="scss" scoped>
.login-page {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 0 20px;

  & > h1 {
    position: absolute;
    top: 20vh;
    text-align: center;
  }

  .login-form {
    max-width: 500px;
    width: 100%;
  }
}
</style>

