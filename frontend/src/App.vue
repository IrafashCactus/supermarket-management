<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'

const authStore = useAuthStore()

onMounted(async () => {
  await authStore.initialize()
})
</script>

<template>
  <div v-if="!authStore.isInitialized" class="flex items-center justify-center min-h-screen bg-surface-50">
    <i class="pi pi-spin pi-spinner text-4xl text-primary-500"></i>
  </div>
  <template v-else>
    <router-view />
  </template>
  <Toast />
  <ConfirmDialog />
</template>
