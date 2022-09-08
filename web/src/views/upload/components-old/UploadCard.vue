
<script lang="ts" setup>
import {defineComponent, nextTick, reactive, ref} from "vue";

import anime from "animejs";

import KSFileList from "./KSFileList.vue";
import OptionsCard from "./OptionsCard.vue";
import KSButton from "@/components/common/KSButton.vue";

import {SendFile, file} from "../utils/file";
import { ProcessState, sleep } from "@/_utils";
import { UploadActionTypes } from "@/store/modules/upload/action-types";

const progressBar = ref<HTMLElement>();
const uploadWrapper = ref<HTMLElement>();

const state = reactive({
  uploadIsStart: false,
  showProgressBar: false,
  uploadProgress: 0,
  //#################
  successUpload: false,
  errorUpload: false,
  //#################
  progressBarOffsetWidth: 0,
  progressBarIntervalTimerAnime: 0,
  progressBarOffsetLeft: 10, // -(number)%
  progressBarAnimeFinish: false,
  //#################
  animeInstanceWrapper: null as anime.AnimeInstance | null,
})

let animeInstanceWrapper: anime.AnimeInstance | undefined


const animateBeforeStart = () => {
  state.uploadIsStart = true;
  const instance = anime({
    targets: uploadWrapper.value,
    height: "2rem",
    delay: 450,
    duration: 700,
    easing: "linear",
    autoplay: false
  });
  nextTick(() => {
    instance.play()
  });
  return instance
}


const StartUpload = async () => {
  const instance = animateBeforeStart()
  await instance.finished
  state.showProgressBar = true;
  nextTick(() => {
    if (uploadWrapper.value instanceof HTMLElement) {
      state.progressBarOffsetWidth = state.progressBarOffsetLeft;
    }
    if (file)
      SendFile(file, {})
  });
}

</script>

<template>
  <div class="ks-upload">
    <div
        ref="uploadWrapper"
        :class="[
        'ks-upload-wrapper',
        { running: uploadIsStart },
        { error: errorUpload },
      ]"
    >
      <transition name="scale">
        <div class="ks-upload-content" v-if="!uploadIsStart">
          <div class="ks-upload__header">
            <span class="header-title">Partage de fichiers</span>
          </div>
          <div class="ks-upload__body">
            <div class="file-view">

            </div>
            <OptionsCard />
          </div>
          <div class="ks-upload__actions">
            <el-button
                class="send-upload"
                :disable="!file"
                @click="StartUpload"
            >
              Partager
            </el-button>
          </div>
        </div>
      </transition>
      <div
          v-if="showProgressBar"
          ref="progressBar"
          :class="[
          'progress-bar',
          { finish: successUpload },
          { error: errorUpload },
        ]"
          :style="{
          left: `-${progressBarOffsetLeft}%`,
        }"
      ></div>
    </div>
  </div>
</template>


<style lang="scss" scoped>
.scale-enter-active,
.scale-leave-active {
  transition: all 0.5s linear;
}

.scale-enter-from,
.scale-leave-to {
  transform: scale(0.9);
  opacity: 0;
}
.ks-upload {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  margin: auto;
  transition: height 0.5s ease;

html[data-useragent*="Android"] &,
html[data-useragent*="iPhone"] & {
  max-height: 100%;
}
}
.ks-upload-wrapper {
  border-top: 0.1rem solid var(--color-border);
  position: relative;
  height: 100%;
  width: 100%;
  transition: width 0.7s 0.45s, border 0.5s 0.5s, border-radius 0.5s 0.7s;
  transition-timing-function: linear;
  overflow: hidden;
  margin: auto;
  padding: 1rem 2rem;

&.running {
 // padding: 0;
   width: 90%;
   border: 2px solid var(--color-border);
   border-radius: 3rem;
 }

&.error {
   border-color: var(--color-error);
 }
}

.ks-upload-content {
  display: flex;
  flex-direction: column;
  justify-content: space-between;

  height: 100%;
}
.ks-upload__header {
  padding: 1rem 0 2rem 0;
  text-align: center;

& .header-title {
    font-size: 2.7rem;
    font-weight: 600;
  }
}
.ks-upload__body {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 1.5rem 0;
  height: 100%;
}
.ks-upload__actions {
  display: flex;
  text-align: center;
  height: 20%;
  padding-top: 2.4rem;
& > .send-upload {
    width: 70%;
    margin: auto;
  }
}

.progress-bar {
  position: absolute;
  top: 0;
  height: 100%;
  background-color: var(--color-primary);
  transition: width 0.5s linear, background-color 0.2s linear;
  border-radius: 5rem;

&.finish {
   background-color: var(--color-success);
 }

&.error {
   background-color: var(--color-error);
 }
}

@media screen and (min-height: 800px) and (min-width: 600px) {
  .ks-upload {
    max-width: 40rem;
    max-height: 71rem;
  }

  .ks-upload-wrapper {
    transition: width 0.7s 0.45s, border 0.2s 0.5s;
    width: 100% !important;
    border: 2px solid var(--color-border);
    border-radius: 3rem;
  }

  .ks-upload-content {
    padding: 2.4rem;
  }
}
</style>
