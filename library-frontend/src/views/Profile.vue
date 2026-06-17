<!-- ============================================================
  【学生自己编写的代码】个人中心页面
  ============================================================ -->
<template>
  <el-card shadow="never">
    <template #header>
      <div style="display:flex;justify-content:space-between;align-items:center;">
        <span>个人信息</span>
      </div>
    </template>

    <el-descriptions :column="1" border v-if="auth.user">
      <el-descriptions-item label="用户ID">{{ auth.user.id }}</el-descriptions-item>
      <el-descriptions-item label="姓名">
        <span style="color:#606266;">{{ auth.user.name }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="邮箱">
        <span style="color:#606266;">{{ auth.user.email }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="手机号">
        <div style="display:flex;align-items:center;gap:10px;">
          <template v-if="!editing">
            <span :style="{ color: auth.user.phone ? '#606266' : '#bbb' }">
              {{ auth.user.phone || '未绑定' }}
            </span>
            <el-button size="small" text type="primary" @click="startEdit">
              {{ auth.user.phone ? '修改' : '绑定' }}
            </el-button>
          </template>
          <template v-else>
            <el-input
              v-model="newPhone"
              placeholder="输入新手机号"
              size="small"
              style="width:200px;"
              maxlength="11"
            />
            <el-button size="small" type="primary" :loading="saving" @click="savePhone">保存</el-button>
            <el-button size="small" @click="cancelEdit">取消</el-button>
          </template>
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="角色">
        <el-tag :type="auth.user.role === 'admin' ? 'danger' : 'info'" size="small">
          {{ auth.user.role === 'admin' ? '管理员' : '普通用户' }}
        </el-tag>
      </el-descriptions-item>
    </el-descriptions>

    <div v-else style="text-align:center;padding:40px;color:#999;">
      加载中...
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { ElMessage } from 'element-plus'
import request from '../api/request'

const auth = useAuthStore()
const editing = ref(false)
const newPhone = ref('')
const saving = ref(false)

function startEdit() {
  newPhone.value = auth.user.phone || ''
  editing.value = true
}

function cancelEdit() {
  editing.value = false
  newPhone.value = ''
}

async function savePhone() {
  const phone = newPhone.value.trim()
  if (!phone) {
    ElMessage.warning('请输入手机号')
    return
  }
  if (!/^1\d{10}$/.test(phone)) {
    ElMessage.warning('请输入11位有效手机号')
    return
  }

  saving.value = true
  try {
    await request.put('/users/me/phone', { phone })
    auth.user.phone = phone
    ElMessage.success('手机号更新成功')
    editing.value = false
  } catch (e) {
    // 已在拦截器中处理
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    await auth.fetchProfile()
  } catch (e) {
    // 已在拦截器中处理
  }
})
</script>


