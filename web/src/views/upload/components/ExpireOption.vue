<script lang="ts" setup>

import {computed, PropType, ref} from "vue";

interface ExpirationOption {
  id: number
  short: string
  long: string
}

interface Props {
  expListOpts: Array<ExpirationOption>
  modelValue: any
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:modelValue', value: any): void
}>()

const expireValue = computed({
  get: () => props.modelValue,
  set: (v: any) => emit("update:modelValue", v)
})



</script>

<template>
  <div class="expire-option-wrapper">
    <span class="expire-option-title">Suppression après:</span>
    <fieldset id="expire-option" class="not-selected">
      <div
          class="expire-option-item"
          v-for="item in props.expListOpts"
          :key="item.id"
      >
        <input
            type="radio"
            :id="'expire-option-item-' + item.id"
            :value="item.id"
            v-model="expireValue"
            name="expire-option"
        />
        <label :for="'expire-option-item-' + item.id" :title="item.long">
          {{ item.short }}
        </label>
      </div>
    </fieldset>
  </div>
</template>

<style lang="scss" scoped>
.expire-option-wrapper {
  display: flex;
  flex-direction: column;
  align-self: center;
  width: 100%;
  box-sizing: border-box;
}

.expire-option--title {
  font-size: 1.4rem;
  padding-bottom: 1rem;
  box-sizing: border-box;
}

#expire-option {
  display: flex;
  align-self: center;
  width: 90%;
  justify-content: space-around;
  border: 0.1rem solid var(--color-border);
  min-width: 22rem;
  max-width: 29rem;
  margin-top: 1rem;
  border-radius: 5rem;
  box-sizing: border-box;
  overflow: hidden;
}

.expire-option-item {
  display: inline-block;
  width: 100%;
  position: relative;
}

.expire-option-item input {
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    position: absolute;
    opacity: 0;
  }


.expire-option-item label {
    position: relative;
    padding: 0.7rem 0.1rem;
    display: block;
    width: 100%;
    border-radius: 5rem;
    color: var(--color-main-text);
    font-size: 1.4rem;
    font-weight: 600;
    box-sizing: border-box;
    opacity: 0.4;
    cursor: pointer;
    text-align: center;
    transition: all 0.1s linear;
}

.expire-option-item label:active {
     background-color: transparent;
}

.expire-option-item input:hover + label,
.expire-option-item input:checked + label {
    opacity: 1;
}

.expire-option-item input:focus + label {
    box-shadow: 0 0 7px var(--color-primary);
}

.expire-option-item input:checked + label {
    color: var(--color-primary-text);
    background-color: var(--color-primary);
}
</style>
