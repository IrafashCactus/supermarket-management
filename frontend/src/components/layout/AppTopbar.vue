<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import Button from 'primevue/button'

const emit = defineEmits(['toggle-sidebar'])
const authStore = useAuthStore()

async function handleLogout() {
  await authStore.logout()
}
</script>

<template>
  <header class="flex items-center justify-between h-16 px-4 bg-white border-b border-surface-200 sm:px-6 lg:px-8">
    <div class="flex items-center">
      <Button
        icon="pi pi-bars"
        text
        rounded
        class="lg:hidden mr-2"
        @click="emit('toggle-sidebar')"
      />
    </div>
    
    <div class="flex items-center space-x-4">
      <div class="flex flex-col items-end hidden sm:flex text-sm">
        <span class="font-medium text-surface-900">{{ authStore.user?.name }}</span>
        <span class="text-surface-500 capitalize">{{ authStore.user?.role }}</span>
      </div>
      
      <Button
        icon="pi pi-sign-out"
        severity="danger"
        text
        rounded
        aria-label="Logout"
        @click="handleLogout"
        title="Logout"
      />
    </div>
  </header>
</template>
