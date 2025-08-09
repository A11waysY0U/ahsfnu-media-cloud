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
      meta: { requiresGuest: true, title: '登录 - AHSFNU 媒体云平台' }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue'),
      meta: { requiresGuest: true, title: '注册 - AHSFNU 媒体云平台' }
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue'),
          meta: { title: '仪表板 - AHSFNU 媒体云平台' }
        },
        {
          path: 'materials',
          name: 'Materials',
          component: () => import('@/views/Materials.vue'),
          meta: { title: '素材管理 - AHSFNU 媒体云平台' }
        },
        {
          path: 'materials/:id',
          name: 'MaterialDetail',
          component: () => import('@/views/MaterialDetail.vue'),
          meta: { title: '素材详情 - AHSFNU 媒体云平台' }
        },
        {
          path: 'tags',
          name: 'Tags',
          component: () => import('@/views/Tags.vue'),
          meta: { title: '标签管理 - AHSFNU 媒体云平台' }
        },
        {
          path: 'workflows',
          name: 'Workflows',
          component: () => import('@/views/Workflows.vue'),
          meta: { title: '工作流管理 - AHSFNU 媒体云平台' }
        },
        {
          path: 'workflows/:id',
          name: 'WorkflowDetail',
          component: () => import('@/views/WorkflowDetail.vue'),
          meta: { title: '工作流详情 - AHSFNU 媒体云平台' }
        },
        {
          path: 'profile',
          name: 'Profile',
          component: () => import('@/views/Profile.vue'),
          meta: { title: '个人资料 - AHSFNU 媒体云平台' }
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('@/views/Users.vue'),
          meta: { requiresAdmin: true, title: '用户管理 - AHSFNU 媒体云平台' }
        },
        {
          path: 'invite-codes',
          name: 'InviteCodes',
          component: () => import('@/views/InviteCodes.vue'),
          meta: { requiresAdmin: true, title: '邀请码管理 - AHSFNU 媒体云平台' }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue'),
      meta: { title: '页面未找到 - AHSFNU 媒体云平台' }
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
  
  // 更新页面标题
  if (to.meta.title) {
    document.title = to.meta.title as string
  } else {
    document.title = 'AHSFNU 媒体云平台'
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
