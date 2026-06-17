<!-- ============================================================
  【学生自己编写的代码】登录/注册/重置密码页面
  ============================================================ -->
<template>
  <div class="login-container">
    <!-- 左侧品牌区 -->
    <div class="brand-section">
      <div class="brand-content">
        <div class="logo-area">
          <div class="logo-icon">
            <svg viewBox="0 0 48 48" width="48" height="48">
              <path d="M6 8h20v28H6z" fill="rgba(255,255,255,0.2)" stroke="#fff" stroke-width="1.5"/>
              <path d="M22 8h20v28H22z" fill="rgba(255,255,255,0.35)" stroke="#fff" stroke-width="1.5"/>
              <rect x="10" y="13" width="12" height="2" rx="1" fill="rgba(255,255,255,0.7)"/>
              <rect x="10" y="18" width="12" height="2" rx="1" fill="rgba(255,255,255,0.7)"/>
              <rect x="10" y="23" width="8" height="2" rx="1" fill="rgba(255,255,255,0.7)"/>
              <rect x="26" y="13" width="12" height="2" rx="1" fill="rgba(255,255,255,0.9)"/>
              <rect x="26" y="18" width="12" height="2" rx="1" fill="rgba(255,255,255,0.9)"/>
              <rect x="26" y="23" width="8" height="2" rx="1" fill="rgba(255,255,255,0.9)"/>
            </svg>
          </div>
          <h1 class="brand-title">云图·智慧图书馆</h1>
          <p class="brand-subtitle">CloudLib — 知识就在指尖</p>
        </div>

        <div class="feature-list">
          <div class="feature-item">
            <span class="feature-icon">📚</span>
            <div>
              <div class="feature-title">万册藏书</div>
              <div class="feature-desc">涵盖文学、科技、历史等多领域优质图书</div>
            </div>
          </div>
          <div class="feature-item">
            <span class="feature-icon">🔍</span>
            <div>
              <div class="feature-title">智能检索</div>
              <div class="feature-desc">按书名、作者、分类快速定位目标图书</div>
            </div>
          </div>
          <div class="feature-item">
            <span class="feature-icon">⏱️</span>
            <div>
              <div class="feature-title">便捷借阅</div>
              <div class="feature-desc">在线借阅、归还，借阅记录一目了然</div>
            </div>
          </div>
          <div class="feature-item">
            <span class="feature-icon">🔐</span>
            <div>
              <div class="feature-title">安全可靠</div>
              <div class="feature-desc">邮箱验证码注册，JWT 身份认证保障数据安全</div>
            </div>
          </div>
        </div>

        <div class="brand-footer">
          <div class="stat-item">
            <span class="stat-number">10,000+</span>
            <span class="stat-label">馆藏图书</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">1,200+</span>
            <span class="stat-label">注册用户</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">98%</span>
            <span class="stat-label">满意度</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧表单区 -->
    <div class="form-section">
      <div class="form-wrapper">
        <div class="form-header">
          <h2 class="form-title">
            {{ mode === 'login' ? '欢迎回来' : mode === 'register' ? '创建账号' : '重置密码' }}
          </h2>
          <p class="form-desc">
            {{ mode === 'login' ? '登录以继续使用图书馆服务' : mode === 'register' ? '注册后即可享受在线借阅服务' : '输入邮箱验证码重置您的密码' }}
          </p>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="0"
          size="large"
          @keyup.enter="handleSubmit"
        >
          <!-- 邮箱 -->
          <el-form-item prop="email">
            <el-input v-model="form.email" placeholder="邮箱地址">
              <template #prefix><el-icon><Message /></el-icon></template>
            </el-input>
          </el-form-item>

          <!-- 验证码（注册/重置模式） -->
          <el-form-item v-if="mode !== 'login'" prop="code">
            <div style="display:flex;width:100%;">
              <el-input v-model="form.code" placeholder="4位验证码" maxlength="4">
                <template #prefix><el-icon><Key /></el-icon></template>
              </el-input>
              <el-button :disabled="cooldown > 0" style="width:120px;margin-left:10px;flex-shrink:0;" @click="sendCode">
                {{ cooldown > 0 ? `${cooldown}s 后重发` : '获取验证码' }}
              </el-button>
            </div>
          </el-form-item>

          <!-- 密码 -->
          <el-form-item v-if="mode !== 'reset'" prop="password">
            <el-input v-model="form.password" type="password" show-password placeholder="密码">
              <template #prefix><el-icon><Lock /></el-icon></template>
            </el-input>
          </el-form-item>

          <!-- 新密码（重置模式） -->
          <el-form-item v-if="mode === 'reset'" prop="newPassword">
            <el-input v-model="form.newPassword" type="password" show-password placeholder="新密码（至少6位）">
              <template #prefix><el-icon><Lock /></el-icon></template>
            </el-input>
          </el-form-item>

          <!-- 手机号（注册模式） -->
          <el-form-item v-if="mode === 'register'" prop="phone">
            <el-input v-model="form.phone" placeholder="手机号（选填）">
              <template #prefix><el-icon><Iphone /></el-icon></template>
            </el-input>
          </el-form-item>

          <!-- 姓名（注册模式） -->
          <el-form-item v-if="mode === 'register'" prop="name">
            <el-input v-model="form.name" placeholder="用户姓名">
              <template #prefix><el-icon><User /></el-icon></template>
            </el-input>
          </el-form-item>

          <!-- 提交按钮 -->
          <el-form-item>
            <el-button type="primary" style="width:100%;" size="large" :loading="loading" @click="handleSubmit">
              {{ mode === 'login' ? '登 录' : mode === 'register' ? '注 册' : '重 置 密 码' }}
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 模式切换 -->
        <div class="form-footer">
          <template v-if="mode === 'login'">
            还没有账号？
            <a href="#" @click.prevent="switchMode('register')">立即注册</a>
            <span class="footer-divider">|</span>
            <a href="#" @click.prevent="switchMode('reset')">忘记密码？</a>
          </template>
          <template v-else>
            <a href="#" @click.prevent="switchMode('login')">← 返回登录</a>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { sendCodeApi, registerApi, resetPasswordApi } from '../api/auth'
import { ElMessage } from 'element-plus'
import { Message, Key, Lock, User, Iphone } from '@element-plus/icons-vue'

const router = useRouter()
const auth = useAuthStore()
const formRef = ref(null)
const loading = ref(false)
const mode = ref('login') // login | register | reset
const cooldown = ref(0)

const form = reactive({
  email: '',
  code: '',
  password: '',
  newPassword: '',
  name: '',
  phone: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 4, message: '验证码为4位数字', trigger: 'blur' },
    { pattern: /^\d{4}$/, message: '验证码必须为数字', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ]
}

function switchMode(m) {
  mode.value = m
  form.code = ''
  form.newPassword = ''
  form.name = ''
  form.phone = ''
}

async function sendCode() {
  if (!form.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }
  const type = mode.value === 'register' ? 'register' : 'reset'
  try {
    await sendCodeApi(form.email, type)
    ElMessage.success('验证码已发送（请查看服务器日志或数据库）')
    cooldown.value = 60
    const timer = setInterval(() => {
      cooldown.value--
      if (cooldown.value <= 0) clearInterval(timer)
    }, 1000)
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

async function handleSubmit() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    if (mode.value === 'login') {
      await auth.login({ email: form.email, password: form.password })
      await auth.fetchProfile()
      ElMessage.success('登录成功')
      router.push('/books')
    } else if (mode.value === 'register') {
      const registerData = {
        email: form.email,
        code: form.code,
        password: form.password,
        name: form.name
      }
      if (form.phone) registerData.phone = form.phone
      await registerApi(registerData)
      ElMessage.success('注册成功，请登录')
      switchMode('login')
    } else if (mode.value === 'reset') {
      await resetPasswordApi({
        email: form.email,
        code: form.code,
        new_password: form.newPassword
      })
      ElMessage.success('密码重置成功，请登录')
      switchMode('login')
    }
  } catch (e) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

/* ===== 左侧品牌区 ===== */
.brand-section {
  flex: 1.2;
  background: linear-gradient(135deg, #1a365d 0%, #2d5a8e 40%, #1a365d 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

/* 装饰圆 */
.brand-section::before {
  content: '';
  position: absolute;
  top: -150px;
  right: -150px;
  width: 400px;
  height: 400px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(255,255,255,0.06) 0%, transparent 70%);
}

.brand-section::after {
  content: '';
  position: absolute;
  bottom: -100px;
  left: -100px;
  width: 300px;
  height: 300px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(255,255,255,0.04) 0%, transparent 70%);
}

.brand-content {
  position: relative;
  z-index: 1;
  padding: 60px;
  max-width: 480px;
}

.logo-area {
  margin-bottom: 50px;
  text-align: center;
}

.logo-icon {
  margin-bottom: 16px;
}

.brand-title {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px 0;
  letter-spacing: 2px;
}

.brand-subtitle {
  font-size: 15px;
  color: rgba(255,255,255,0.6);
  margin: 0;
  letter-spacing: 1px;
}

.feature-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin-bottom: 50px;
}

.feature-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.feature-icon {
  font-size: 28px;
  line-height: 1;
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.feature-title {
  font-size: 16px;
  font-weight: 600;
  color: rgba(255,255,255,0.95);
  margin-bottom: 4px;
}

.feature-desc {
  font-size: 13px;
  color: rgba(255,255,255,0.55);
  line-height: 1.5;
}

.brand-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 30px;
  padding-top: 30px;
  border-top: 1px solid rgba(255,255,255,0.1);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-number {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
}

.stat-label {
  font-size: 12px;
  color: rgba(255,255,255,0.5);
}

.stat-divider {
  width: 1px;
  height: 36px;
  background: rgba(255,255,255,0.15);
}

/* ===== 右侧表单区 ===== */
.form-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f7f8fa;
}

.form-wrapper {
  width: 400px;
  padding: 20px 10px;
}

.form-header {
  margin-bottom: 32px;
}

.form-title {
  font-size: 26px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.form-desc {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  transition: box-shadow 0.2s;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #2d5a8e inset;
}

:deep(.el-button--primary) {
  height: 44px;
  border-radius: 8px;
  font-size: 16px;
  background: #2d5a8e;
  border-color: #2d5a8e;
}

:deep(.el-button--primary:hover) {
  background: #1a365d;
  border-color: #1a365d;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
}

.form-footer {
  text-align: center;
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}

.form-footer a {
  color: #2d5a8e;
  text-decoration: none;
  font-weight: 500;
}

.form-footer a:hover {
  text-decoration: underline;
}

.footer-divider {
  margin: 0 10px;
  color: #dcdfe6;
}

/* ===== 响应式 ===== */
@media (max-width: 900px) {
  .brand-section {
    display: none;
  }
  .form-section {
    flex: 1;
  }
}
</style>


