<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Message from 'primevue/message'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const isLoading = ref(false)
const errorMessage = ref('')

async function onSubmit() {
  if (!email.value || !password.value) {
    errorMessage.value = 'Email and password are required'
    return
  }
  
  isLoading.value = true
  errorMessage.value = ''
  
  try {
    await authStore.login({ email: email.value, password: password.value })
    router.push({ name: 'root' })
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error?.message || 'Invalid credentials or server error'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-surface-50 px-4">
    <div class="w-full max-w-md p-8 bg-white rounded-lg shadow-sm border border-surface-200">
      <div class="text-center mb-8">
        <h1 class="text-2xl font-bold text-primary-600">SuperMarket POS</h1>
        <p class="text-surface-500 mt-2">Sign in to your account</p>
      </div>

      <Message v-if="errorMessage" severity="error" :closable="false" class="mb-4">
        {{ errorMessage }}
      </Message>

      <form @submit.prevent="onSubmit" class="space-y-4">
        <div class="flex flex-col space-y-2">
          <label for="email" class="text-sm font-medium text-surface-900">Email</label>
          <InputText
            id="email"
            v-model="email"
            type="email"
            placeholder="Enter your email"
            class="w-full"
            :disabled="isLoading"
          />
        </div>

        <div class="flex flex-col space-y-2">
          <label for="password" class="text-sm font-medium text-surface-900">Password</label>
          <Password
            id="password"
            v-model="password"
            placeholder="Enter your password"
            :feedback="false"
            toggleMask
            class="w-full"
            inputClass="w-full"
            :disabled="isLoading"
          />
        </div>

        <Button
          type="submit"
          label="Sign In"
          class="w-full"
          :loading="isLoading"
        />
      </form>
    </div>
  </div>
</template>
