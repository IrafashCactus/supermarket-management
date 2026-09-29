<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const allMenuItems = [
  { label: 'Dashboard', icon: 'pi pi-home', routeName: 'dashboard', roles: ['admin', 'manager', 'accountant'] },
  { label: 'Products', icon: 'pi pi-box', routeName: 'products', roles: ['admin', 'manager', 'storekeeper'], disabled: true },
  { label: 'Categories', icon: 'pi pi-tags', routeName: 'categories', roles: ['admin', 'manager'], disabled: true },
  { label: 'Inventory', icon: 'pi pi-list', routeName: 'inventory', roles: ['admin', 'manager', 'storekeeper'], disabled: true },
  { label: 'Suppliers', icon: 'pi pi-truck', routeName: 'suppliers', roles: ['admin', 'manager', 'storekeeper'], disabled: true },
  { label: 'Purchases', icon: 'pi pi-shopping-bag', routeName: 'purchases', roles: ['admin', 'manager', 'storekeeper'], disabled: true },
  { label: 'POS', icon: 'pi pi-desktop', routeName: 'pos', roles: ['cashier'], disabled: true },
  { label: 'Sales', icon: 'pi pi-chart-line', routeName: 'sales', roles: ['admin', 'manager', 'accountant', 'cashier'], disabled: true },
  { label: 'Reports', icon: 'pi pi-file', routeName: 'reports', roles: ['admin', 'manager', 'accountant'], disabled: true },
  { label: 'Users', icon: 'pi pi-users', routeName: 'users', roles: ['admin'] },
]

const menuItems = computed(() => {
  return allMenuItems.filter(item => authStore.checkRole(item.roles))
})

function navigate(routeName: string, disabled?: boolean) {
  if (disabled) return
  router.push({ name: routeName })
}
</script>

<template>
  <div class="flex flex-col h-full bg-white border-r border-surface-200">
    <div class="flex items-center justify-center h-16 border-b border-surface-200">
      <span class="text-xl font-bold text-primary-600">SuperMarket POS</span>
    </div>
    
    <div class="flex-1 overflow-y-auto py-4">
      <nav class="space-y-1 px-3">
        <a
          v-for="item in menuItems"
          :key="item.routeName"
          href="#"
          @click.prevent="navigate(item.routeName, item.disabled)"
          :class="[
            route.name === item.routeName
              ? 'bg-primary-50 text-primary-700'
              : 'text-surface-700 hover:bg-surface-50',
            item.disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer',
            'flex items-center px-3 py-2 text-sm font-medium rounded-md'
          ]"
        >
          <i :class="[item.icon, 'mr-3 text-lg']"></i>
          {{ item.label }}
          <span v-if="item.disabled" class="ml-auto text-xs bg-surface-200 px-2 py-1 rounded">Soon</span>
        </a>
      </nav>
    </div>
  </div>
</template>
