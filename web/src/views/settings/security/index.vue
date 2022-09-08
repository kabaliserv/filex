<script lang="ts" setup>

import {reactive, ref} from "vue";
import {ElForm, ElMessage} from "element-plus";
import {sleep} from "@/utils";
import * as api from "@/api"
import {changePassword} from "@/api/users";
import {UserModule} from "@/store/modules/user";
import {AxiosResponse} from "axios";

const elFormRef = ref<InstanceType<typeof ElForm>>()

const changePasswordForm = reactive({
  oldPassword: "",
  newPassword: "",
  confirmPassword: "",
})


const loadingButton = ref(false)

const oldPasswordErrorMessage = ref<string>()


const validatePass = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error("S'il vous plaît entrez votre mot de passe"))
  } else {
    callback()
  }
}

const validateConfirmPass = (rule: any, value: string, callback: any) => {
  if (value == '') {
    callback(new Error("S'il vous plaît confirmer votre mot de passe"))
  } else if (value !== changePasswordForm.newPassword) {
    callback(new Error("La confirmation du mot de passe n'est pas correct !"))
  }else {
    callback()
  }
}

const elFormRules = reactive({
  oldPassword: [{required: true}, {validator: validatePass, trigger: 'blur'}],
  newPassword: [{required: true}, {validator: validatePass, trigger: 'blur'}],
  confirmPassword: [{required: true}, {validator: validateConfirmPass, trigger: 'blur'}],
})

const HandleChangePassword = () => {
  if (!elFormRef.value) return
  elFormRef.value?.validate(async (valid: any) => {
    if (valid) {
      loadingButton.value = true

      try {
        await api.users.changePassword(
            UserModule.id,
            {
              old_password: changePasswordForm.oldPassword,
              new_password: changePasswordForm.newPassword
            })

        ElMessage({
          type: "success",
          message: "Le mot de passe à bien été mis à jour"
        })

        elFormRef.value?.resetFields()
      } catch (error: any) {
        const response: AxiosResponse = error.response
        if (response && response.status == 403) {
          oldPasswordErrorMessage.value = "Le mot de passe est incorrecte !"
        } else {
          ElMessage({
            type: "error",
            message: "Une erreur est survenu lors de la mis à jour du mot de passe, veuillez contacter l'administrateur si le problème persiste."
          })
          // Something happened in setting up the request and triggered an Error
          console.log('Error', error.message);
        }
      }

      loadingButton.value = false
    }
  })
}
</script>

<template>
  <div class="security-view">
    <div class="password-view">
      <h2>Changer de mot de passe</h2>
      <el-divider/>
      <el-form
          ref="elFormRef"
          :rules="elFormRules"
          :model="changePasswordForm"
          label-position="top"
          @submit.prevent="HandleChangePassword"
      >
        <el-form-item label="Ancien mot de passe" prop="oldPassword" :error="oldPasswordErrorMessage">
          <el-input type="password" v-model="changePasswordForm.oldPassword" :disabled="loadingButton"></el-input>
        </el-form-item>
        <el-form-item label="Nouveau mot de passe" prop="newPassword">
          <el-input type="password" v-model="changePasswordForm.newPassword" :disabled="loadingButton"></el-input>
        </el-form-item>
        <el-form-item label="Confirmer mot de passe" prop="confirmPassword">
          <el-input type="password" v-model="changePasswordForm.confirmPassword" :disabled="loadingButton"></el-input>
        </el-form-item>
        <div class="password-desc"></div>
        <div class="form-actions">
          <el-button native-type="submit" :loading="loadingButton">Changer</el-button>
        </div>
      </el-form>

    </div>
  </div>
</template>

<style lang="scss">
.el-form-item {
  max-width: 500px;
}
.password-view {
  .el-input {
    max-width: 350px;
  }

  @media screen and (max-width: 900px) {
    .el-input {
      //max-width: 300px;
      max-width: 500px;
    }
  }
}

</style>