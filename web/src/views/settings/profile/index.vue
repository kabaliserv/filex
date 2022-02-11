<script lang="ts" setup>
import { UserModule } from "@/store/modules/user";

</script>

<template>
  <div>
    <h2>Profile Publique</h2>
    <el-divider></el-divider>
    <el-row class="profile-view">
      <el-col :span="6">
        <el-form label-position="top">
          <el-form-item label="Image de profile">
            <div class="avatar-view">
              <el-avatar class="avatar-img" :src="UserModule.avatar" :size="200">
                <span v-if="!UserModule.avatar" class="iconify" data-icon="fa:user" style="font-size: 60px"></span>
              </el-avatar>
              <el-button class="avatar-edit-btn">Editer</el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-col>
      <el-col style="flex: 1">
        <el-form :label-position="'top'">
          <el-form-item label="Nom">
            <el-input :model-value="UserModule.login"></el-input>
          </el-form-item>
          <el-form-item label="Bio">
            <el-input></el-input>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>

<style lang="scss">
.profile-view {
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

<style lang="scss" scoped>
.profile-view {
  flex-direction: row-reverse;
  justify-content: flex-end;
  gap: 30px
}

.avatar-view {
  position: relative;
  .avatar-img {
    cursor: pointer;
  }

  .avatar-edit-btn {
    position: absolute;
    bottom: 5%;
    left: 5%;
  }
}

@media screen and (max-width: 900px) {
  .profile-view {
    flex-direction: column;
  }
}

</style>
