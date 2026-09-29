import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: () => import('@/components/layout/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'root',
          redirect: () => {
            const authStore = useAuthStore()
            if (authStore.hasRole('cashier')) {
              return { name: 'pos' } // not implemented yet, but keeping logic
            }
            return { name: 'dashboard' }
          }
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
          meta: { roles: ['admin', 'manager', 'accountant'] }
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/UsersView.vue'),
          meta: { roles: ['admin'] }
        },
        {
          path: 'pos',
          name: 'pos',
          component: () => import('@/views/DashboardView.vue'), // placeholder for POS
          meta: { roles: ['cashier'] }
        }
      ]
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth ?? true

  if (requiresAuth && !authStore.isAuthenticated) {
    return next({ name: 'login' })
  }

  if (to.name === 'login' && authStore.isAuthenticated) {
    return next({ name: 'root' })
  }

  if (to.meta.roles) {
    const roles = to.meta.roles as string[]
    if (!authStore.checkRole(roles)) {
      return next({ name: 'root' })
    }
  }

  next()
})

export default router
