import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue')
        },
        {
          path: 'materials',
          name: 'Materials',
          component: () => import('@/views/Materials.vue')
        },
        {
          path: 'materials/:id',
          name: 'MaterialDetail',
          component: () => import('@/views/MaterialDetail.vue')
        },
        {
          path: 'tags',
          name: 'Tags',
          component: () => import('@/views/Tags.vue')
        },
        {
          path: 'workflows',
          name: 'Workflows',
          component: () => import('@/views/Workflows.vue')
        },
        {
          path: 'workflows/:id',
          name: 'WorkflowDetail',
          component: () => import('@/views/WorkflowDetail.vue')
        },
        {
          path: 'profile',
          name: 'Profile',
          component: () => import('@/views/Profile.vue')
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('@/views/Users.vue'),
          meta: { requiresAdmin: true }
        },
        {
          path: 'invite-codes',
          name: 'InviteCodes',
          component: () => import('@/views/InviteCodes.vue'),
          meta: { requiresAdmin: true }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // 初始化认证状态
  if (!authStore.isAuthenticated) {
    authStore.initAuth()
  }
  
  // 需要认证的页面
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }
  
  // 需要管理员权限的页面
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next('/dashboard')
    return
  }
  
  // 已登录用户不能访问登录/注册页面
  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }
  
  next()
})

export default router
