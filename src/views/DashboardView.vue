<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>

<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold text-surface-900">
      Welcome back, {{ authStore.user?.name }}
    </h1>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div class="bg-white p-6 rounded-lg shadow-sm border border-surface-200 flex flex-col items-center justify-center h-48">
        <i class="pi pi-chart-bar text-4xl text-surface-400 mb-4"></i>
        <h3 class="text-lg font-medium text-surface-900">Analytics</h3>
        <p class="text-surface-500 mt-2 text-sm">Coming in Phase 7</p>
      </div>

      <div class="bg-white p-6 rounded-lg shadow-sm border border-surface-200 flex flex-col items-center justify-center h-48">
        <i class="pi pi-box text-4xl text-surface-400 mb-4"></i>
        <h3 class="text-lg font-medium text-surface-900">Recent Inventory</h3>
        <p class="text-surface-500 mt-2 text-sm">Coming in Phase 7</p>
      </div>
      
      <div class="bg-white p-6 rounded-lg shadow-sm border border-surface-200 flex flex-col items-center justify-center h-48">
        <i class="pi pi-shopping-cart text-4xl text-surface-400 mb-4"></i>
        <h3 class="text-lg font-medium text-surface-900">Latest Sales</h3>
        <p class="text-surface-500 mt-2 text-sm">Coming in Phase 7</p>
      </div>
    </div>
  </div>
</template>
