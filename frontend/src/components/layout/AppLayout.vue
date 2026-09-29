<script setup lang="ts">
import { ref } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppTopbar from './AppTopbar.vue'

const isSidebarOpen = ref(false)

function toggleSidebar() {
  isSidebarOpen.value = !isSidebarOpen.value
}
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-surface-50">
    <!-- Sidebar for mobile -->
    <div
      v-if="isSidebarOpen"
      class="fixed inset-0 z-40 bg-black bg-opacity-50 lg:hidden"
      @click="toggleSidebar"
    ></div>

    <AppSidebar
      :class="[
        isSidebarOpen ? 'translate-x-0' : '-translate-x-full',
        'fixed inset-y-0 left-0 z-50 w-64 transition-transform duration-300 lg:static lg:translate-x-0'
      ]"
    />

    <div class="flex flex-col flex-1 min-w-0 overflow-hidden">
      <AppTopbar @toggle-sidebar="toggleSidebar" />
      
      <main class="flex-1 overflow-auto p-4 sm:p-6 lg:p-8">
        <router-view />
      </main>
    </div>
  </div>
</template>
