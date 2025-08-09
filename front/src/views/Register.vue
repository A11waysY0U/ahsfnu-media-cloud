<template>
  <div class="register-container">
    <div class="register-background">
      <div class="register-form-container">
        <div class="register-header">
          <h1 class="register-title">AHSFNU 媒体云平台</h1>
          <p class="register-subtitle">创建您的账户</p>
        </div>
        
        <el-form
          ref="registerFormRef"
          :model="registerForm"
          :rules="registerRules"
          class="register-form"
          @submit.prevent="handleRegister"
        >
          <el-form-item prop="username">
            <el-input
              v-model="registerForm.username"
              placeholder="用户名 (3-50字符)"
              size="large"
              prefix-icon="User"
            />
          </el-form-item>
          
          <el-form-item prop="email">
            <el-input
              v-model="registerForm.email"
              placeholder="邮箱地址"
              size="large"
              prefix-icon="Message"
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="密码 (最少6位)"
              size="large"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="确认密码"
              size="large"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item prop="inviteCode">
            <el-input
              v-model="registerForm.inviteCode"
              placeholder="邀请码"
              size="large"
              prefix-icon="Key"
            />
          </el-form-item>
          
          <el-form-item prop="captchaCode">
            <div class="captcha-container">
              <el-input
                v-model="registerForm.captchaCode"
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
              class="register-button"
              :loading="loading"
              @click="handleRegister"
            >
              注册
            </el-button>
          </el-form-item>
          
          <div class="register-footer">
            <span>已有账户？</span>
            <el-link type="primary" @click="$router.push('/login')">立即登录</el-link>
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

const registerFormRef = ref<FormInstance>()
const loading = ref(false)
const captchaImage = ref('')
const captchaId = ref('')
const authToken = ref('')

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  inviteCode: '',
  captchaCode: ''
})

const registerRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在3到50个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  inviteCode: [
    { required: true, message: '请输入邀请码', trigger: 'blur' }
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
  registerForm.captchaCode = ''
}

// 验证验证码
const verifyCaptcha = async () => {
  try {
    await authAPI.verifyCaptcha(captchaId.value, registerForm.captchaCode)
    return true
  } catch (error) {
    ElMessage.error('验证码错误')
    refreshCaptcha()
    return false
  }
}

// 处理注册
const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  try {
    await registerFormRef.value.validate()
    
    loading.value = true
    
    // 先验证验证码
    const captchaValid = await verifyCaptcha()
    if (!captchaValid) {
      loading.value = false
      return
    }
    
    // 执行注册
    const result = await authStore.register(
      registerForm.username,
      registerForm.email,
      registerForm.password,
      registerForm.inviteCode,
      authToken.value
    )
    
    if (result.success) {
      ElMessage.success('注册成功')
      router.push('/dashboard')
    } else {
      ElMessage.error('注册失败，请检查输入信息')
      refreshCaptcha()
    }
  } catch (error) {
    ElMessage.error('注册失败')
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
.register-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-background {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  padding: 40px;
  width: 100%;
  max-width: 400px;
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.register-subtitle {
  font-size: 14px;
  color: #666;
}

.register-form {
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

.register-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
}

.register-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #666;
}

.register-footer span {
  margin-right: 8px;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .register-background {
    margin: 20px;
    padding: 24px;
  }
  
  .register-title {
    font-size: 24px;
  }
  
  .captcha-image {
    width: 100px;
  }
}
</style>
