<template>
  <div class="login-container">
    <div class="login-background">
      <div class="login-form-container">
        <div class="login-header">
          <h1 class="login-title">AHSFNU 媒体云平台</h1>
          <p class="login-subtitle">欢迎回来，请登录您的账户</p>
        </div>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
          @submit.prevent="handleLogin"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="用户名或邮箱"
              size="large"
              prefix-icon="User"
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              size="large"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item prop="captchaCode">
            <div class="captcha-container">
              <el-input
                v-model="loginForm.captchaCode"
                placeholder="验证码"
                size="large"
                prefix-icon="Key"
                style="flex: 1; margin-right: 12px;"
              />
              <div class="captcha-image" @click="refreshCaptcha">
                <img v-if="captchaImage" :src="captchaImage" alt="验证码" />
                <el-icon v-else><Refresh /></el-icon>
              </div>
            </div>
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              size="large"
              class="login-button"
              :loading="loading"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
          
          <div class="login-footer">
            <span>还没有账户？</span>
            <el-link type="primary" @click="$router.push('/register')">立即注册</el-link>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { authAPI } from '@/api'

const router = useRouter()
const authStore = useAuthStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)
const captchaImage = ref('')
const captchaId = ref('')
const authToken = ref('')

const loginForm = reactive({
  username: '',
  password: '',
  captchaCode: ''
})

const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名或邮箱', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  captchaCode: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 4, message: '验证码长度为4位', trigger: 'blur' }
  ]
}

// 获取验证码
const getCaptcha = async () => {
  try {
    const response = await authAPI.getCaptcha()
    const { captcha_id, captcha_b64, auth_token } = response.data
    captchaId.value = captcha_id
    captchaImage.value = captcha_b64
    authToken.value = auth_token
  } catch (error) {
    ElMessage.error('获取验证码失败')
  }
}

// 刷新验证码
const refreshCaptcha = () => {
  getCaptcha()
  loginForm.captchaCode = ''
}

// 验证验证码
const verifyCaptcha = async () => {
  try {
    await authAPI.verifyCaptcha(captchaId.value, loginForm.captchaCode)
    return true
  } catch (error) {
    ElMessage.error('验证码错误')
    refreshCaptcha()
    return false
  }
}

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    
    loading.value = true
    
    // 先验证验证码
    const captchaValid = await verifyCaptcha()
    if (!captchaValid) {
      loading.value = false
      return
    }
    
    // 执行登录
    const result = await authStore.login(
      loginForm.username,
      loginForm.password,
      authToken.value
    )
    
    if (result.success) {
      ElMessage.success('登录成功')
      router.push('/dashboard')
    } else {
      ElMessage.error('登录失败，请检查用户名和密码')
      refreshCaptcha()
    }
  } catch (error) {
    ElMessage.error('登录失败')
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  getCaptcha()
})
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-background {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  padding: 40px;
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: #666;
}

.login-form {
  width: 100%;
}

.captcha-container {
  display: flex;
  align-items: center;
}

.captcha-image {
  width: 120px;
  height: 40px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  background-color: #f5f7fa;
  transition: all 0.3s;
}

.captcha-image:hover {
  border-color: #409eff;
  background-color: #f0f9ff;
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 4px;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #666;
}

.login-footer span {
  margin-right: 8px;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .login-background {
    margin: 20px;
    padding: 24px;
  }
  
  .login-title {
    font-size: 24px;
  }
  
  .captcha-image {
    width: 100px;
  }
}
</style>
