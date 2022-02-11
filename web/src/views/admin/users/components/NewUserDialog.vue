<script lang="ts" setup>
import {h, reactive, ref, watch} from "vue";
import {ElForm, ElMessage, ElMessageBox} from "element-plus";
import type { Action as ElAction } from 'element-plus'
import {urlAlphabet, customAlphabet} from "nanoid"
import {NewUser} from "@/types/request";
import * as api from "@/api"
import {UserApi} from "@/types/response";

const passwordAlphabet = '-$&#%azertyuiopqsdfghjklmwxcvbn_AZERTYUIOPQSDFGHJKLMWXCVBN'
const generatePassword = customAlphabet(passwordAlphabet, 30)

type Props = {
  modelValue: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false
})

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void
  (e: "newUser", value: UserApi): void
}>()

const newUserDialog = ref(props.modelValue)

const elFormRef = ref<InstanceType<typeof ElForm>>()

const newUser = reactive({
  login: "",
  email: "",
  password: "",
  autoGeneratePassword: false,
  active: true,
  admin: false,
  quota: "10000",
  activeQuota: true,
})

const resetUser = () => {
  newUser.login = ""
  newUser.email = ""
  newUser.password = ""
  newUser.autoGeneratePassword = false
  newUser.active = true
  newUser.admin = false
  newUser.quota = "10000"
  newUser.activeQuota = true
}

const validateLogin = (rule: any, value: string, callback: any) => {
  if (value === '' || value.length < 3) {
    callback(new Error("S'il vous plaît entrez un identifiant valide"))
  } else {
    callback()
  }
}

const validateEmail = (rule: any, value: string, callback: any) => {
  const regexEmail = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
  if (value === '' || value.length < 5 || !regexEmail.test(value)) {
    callback(new Error("S'il vous plaît entrez un email valide"))
  } else {
    callback()
  }
}

const elFormRules = reactive({
  login: [{ required: true, validator: validateLogin ,trigger: 'blur' }],
  email: [{required: true}, {validator: validateEmail ,trigger: 'blur'}],
  password: [{ required: true }],
})

const onlyNumber = (event: KeyboardEvent) => {
  if (
      event.ctrlKey ||
      event.shiftKey ||
      event.altKey ||
      event.metaKey ||
      event.key == "Backspace" ||
      event.key == "Tab" ||
      event.key.startsWith("Arrow")
  ) return
  const num = parseInt(event.key)
  if (isNaN(num)) {
    event.preventDefault()
    return false
  }
  return true

}

const handleConfirmNewUser = async () => {
  if (!elFormRef.value) return
  elFormRef.value.validate(async (valid: any) => {
    if (valid) {
      const userReq: NewUser = {
        login: newUser.login,
        email: newUser.email,
        password: newUser.password,
        active: newUser.active,
        admin: newUser.admin,
        quota: 0,
        activeQuota: newUser.activeQuota
      }
      if (newUser.activeQuota) {
        const num = parseInt(newUser.quota)
        if (!isNaN(num))
        userReq.quota = num * 1024 * 1024
      }
      try {
        const { data } = await api.users.addUser(userReq)
        if (newUser.autoGeneratePassword) {
          ElMessageBox({
            title: "Utilisateur",
            message: h("div",
                {style: "display: flex; flex-direction: column"},
                [
                  h("span", [
                    h("span", {style: "font-weight: bold;"}, "Login: "),
                    h("span", newUser.login)
                  ]),
                  h("span", [
                    h("span", {style: "font-weight: bold;"}, "Password: "),
                    h("span", newUser.password)
                  ]),
                ])
          })
        }
        emit("newUser", data)
        handleCloseDialog()

      } catch (e) {
        ElMessage.error({message: "Erreur lors de la creation de l'utilisateur"})
      }
    } else {
      return false
    }
  })
}

const handleBeforeCloseDialog = (done: () => void) => {
  ElMessageBox.confirm("Êtes vous sur de vouloir fermer cette fenêtre ?")
      .then(() => {
        done()
        handleCloseDialog()
      })
      .catch(() => {
        // catch error
      })
}

const handleCloseDialog = () => {
  emit("update:modelValue", false)
  newUserDialog.value = false
  resetUser()
}


watch(() => props.modelValue, (value) => {
  newUserDialog.value = value
})

watch(() => newUser.autoGeneratePassword, (value) => {
  if (value) {
    newUser.password = generatePassword()
    elFormRef.value?.validateField("password", () => {})
  }
   else
     newUser.password = ""
})
</script>

<template>
  <el-dialog
      v-model="newUserDialog"
      title="Nouvel Utilisateur"
      width="500px"
      top="15vh"
      :before-close="handleBeforeCloseDialog"
  >
    <el-form
        ref="elFormRef"
        :model="newUser"
        :rules="elFormRules"
        label-position="top"
        label-width="120px"
        @submit="handleConfirmNewUser"
    >
      <el-form-item label="Nom d'utilisateur" prop="login">
        <el-input
            v-model="newUser.login"
            type="text"
            autocomplete="username"
        ></el-input>
      </el-form-item>
      <el-form-item label="Email" prop="email">
        <el-input
            v-model="newUser.email"
            type="text"
            autocomplete="email"
        ></el-input>
      </el-form-item>
      <el-form-item label="Mot de passe" prop="password">
        <el-input
            v-model="newUser.password"
            type="password"
            autocomplete="new-password"
            :disabled="newUser.autoGeneratePassword"
        ></el-input>
        <span style="font-size: 15px">Générer Aléatoirement <el-switch v-model="newUser.autoGeneratePassword" size="large"/></span>
      </el-form-item>
      <el-form-item label="Quota" prop="quota">
        <template #label>
          Quota <el-switch v-model="newUser.activeQuota" :size="'small'"/>
        </template>
        <el-input
            v-model="newUser.quota"
            type="text"
            autocomplete="new-password"
            @keydown="onlyNumber"
            :disabled="!newUser.activeQuota"
        >
          <template #suffix>MB</template>
        </el-input>
      </el-form-item>
    </el-form>
    <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleCloseDialog">Cancel</el-button>
          <el-button type="primary" @click="handleConfirmNewUser"
          >Confirm</el-button
          >
        </span>
    </template>
  </el-dialog>
</template>

<style lang="scss"></style>
