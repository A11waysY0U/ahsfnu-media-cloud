import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authAPI } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isAuthenticated = computed(() => !!token.value)

  // 初始化状态
  const initAuth = () => {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')
    
    if (storedToken && storedUser) {
      token.value = storedToken
      user.value = JSON.parse(storedUser)
    }
  }

  // 登录
  const login = async (username: string, password: string, authToken: string) => {
    try {
      const response = await authAPI.login({ username, password, auth_token: authToken })
      const { token: newToken, user: userData } = response.data
      
      token.value = newToken
      user.value = userData
      
      localStorage.setItem('token', newToken)
      localStorage.setItem('user', JSON.stringify(userData))
      
      return { success: true }
    } catch (error) {
      return { success: false, error }
    }
  }

  // 注册
  const register = async (username: string, email: string, password: string, inviteCode: string, authToken: string) => {
    try {
      const response = await authAPI.register({ 
        username, 
        email, 
        password, 
        invite_code: inviteCode, 
        auth_token: authToken 
      })
      const { token: newToken, user: userData } = response.data
      
      token.value = newToken
      user.value = userData
      
      localStorage.setItem('token', newToken)
      localStorage.setItem('user', JSON.stringify(userData))
      
      return { success: true }
    } catch (error) {
      return { success: false, error }
    }
  }

  // 登出
  const logout = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // 获取个人资料
  const fetchProfile = async () => {
    try {
      const response = await authAPI.getProfile()
      user.value = response.data
      localStorage.setItem('user', JSON.stringify(response.data))
      return { success: true }
    } catch (error) {
      return { success: false, error }
    }
  }

  // 更新个人资料
  const updateProfile = async (data: Partial<User>) => {
    try {
      const response = await authAPI.updateProfile(data)
      user.value = response.data
      localStorage.setItem('user', JSON.stringify(response.data))
      return { success: true }
    } catch (error) {
      return { success: false, error }
    }
  }

  // 修改密码
  const changePassword = async (currentPassword: string, newPassword: string) => {
    try {
      await authAPI.changePassword(currentPassword, newPassword)
      return { success: true }
    } catch (error) {
      return { success: false, error }
    }
  }

  // 检查是否为管理员
  const isAdmin = computed(() => user.value?.role === 'admin')

  return {
    user,
    token,
    isAuthenticated,
    isAdmin,
    initAuth,
    login,
    register,
    logout,
    fetchProfile,
    updateProfile,
    changePassword,
  }
})
