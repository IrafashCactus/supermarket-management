import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest } from '@/types/auth'
import { authService } from '@/services/authService'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const isInitialized = ref(false)

  const isAuthenticated = computed(() => !!accessToken.value)
  const userRole = computed(() => user.value?.role)

  function hasRole(role: User['role']) {
    return userRole.value === role
  }

  function checkRole(allowedRoles: string[]) {
    if (!userRole.value) return false
    return allowedRoles.includes(userRole.value)
  }

  async function login(credentials: LoginRequest) {
    const response = await authService.login(credentials)
    accessToken.value = response.access_token
    user.value = response.user
  }

  async function logout() {
    try {
      await authService.logout()
    } catch (e) {
      // Ignore errors on logout
    } finally {
      accessToken.value = null
      user.value = null
      window.location.href = '/login'
    }
  }

  async function refreshToken() {
    const response = await authService.refreshToken()
    accessToken.value = response.access_token
  }

  async function initialize() {
    try {
      await refreshToken()
      user.value = await authService.getMe()
    } catch (e) {
      accessToken.value = null
      user.value = null
    } finally {
      isInitialized.value = true
    }
  }

  return {
    user,
    accessToken,
    isInitialized,
    isAuthenticated,
    userRole,
    hasRole,
    checkRole,
    login,
    logout,
    refreshToken,
    initialize
  }
})
