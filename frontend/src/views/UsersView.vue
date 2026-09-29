<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userService } from '@/services/userService'
import type { User, CreateUserRequest, UpdateUserRequest } from '@/types/user'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Password from 'primevue/password'
import Tag from 'primevue/tag'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'

const toast = useToast()
const confirm = useConfirm()

const users = ref<User[]>([])
const loading = ref(false)
const totalRecords = ref(0)
const rows = ref(10)

const userDialog = ref(false)
const isEdit = ref(false)
const userForm = ref<any>({ name: '', email: '', role: 'cashier', password: '', is_active: true })
const passwordDialog = ref(false)
const resetPasswordValue = ref('')
const selectedUserId = ref('')
const submitting = ref(false)

const roles = [
  { label: 'Admin', value: 'admin' },
  { label: 'Manager', value: 'manager' },
  { label: 'Accountant', value: 'accountant' },
  { label: 'Storekeeper', value: 'storekeeper' },
  { label: 'Cashier', value: 'cashier' }
]

async function loadUsers(event?: any) {
  loading.value = true
  try {
    const page = event ? (event.first / event.rows) + 1 : 1
    const response = await userService.getUsers(page, rows.value)
    users.value = response.data
    totalRecords.value = response.meta.total
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load users', life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadUsers()
})

function openNew() {
  userForm.value = { name: '', email: '', role: 'cashier', password: '', is_active: true }
  isEdit.value = false
  userDialog.value = true
}

function editUser(user: User) {
  userForm.value = { ...user }
  isEdit.value = true
  userDialog.value = true
}

function openResetPassword(user: User) {
  selectedUserId.value = user.id
  resetPasswordValue.value = ''
  passwordDialog.value = true
}

async function saveUser() {
  submitting.value = true
  try {
    if (isEdit.value) {
      const payload: UpdateUserRequest = {
        name: userForm.value.name,
        email: userForm.value.email,
        role: userForm.value.role,
        is_active: userForm.value.is_active
      }
      await userService.updateUser(userForm.value.id, payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'User updated', life: 3000 })
    } else {
      const payload: CreateUserRequest = {
        name: userForm.value.name,
        email: userForm.value.email,
        role: userForm.value.role,
        password: userForm.value.password
      }
      await userService.createUser(payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'User created', life: 3000 })
    }
    userDialog.value = false
    loadUsers()
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.error?.message || 'Failed to save user', life: 3000 })
  } finally {
    submitting.value = false
  }
}

async function savePassword() {
  if (!resetPasswordValue.value) return
  
  submitting.value = true
  try {
    await userService.resetPassword(selectedUserId.value, { password: resetPasswordValue.value })
    toast.add({ severity: 'success', summary: 'Success', detail: 'Password reset', life: 3000 })
    passwordDialog.value = false
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to reset password', life: 3000 })
  } finally {
    submitting.value = false
  }
}

async function toggleUserStatus(user: User) {
  const newStatus = !user.is_active
  const action = newStatus ? 'enable' : 'disable'
  
  confirm.require({
    message: `Are you sure you want to ${action} ${user.name}?`,
    header: 'Confirm',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await userService.updateUser(user.id, {
          name: user.name,
          email: user.email,
          role: user.role,
          is_active: newStatus
        })
        toast.add({ severity: 'success', summary: 'Success', detail: `User ${action}d`, life: 3000 })
        loadUsers()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to change status', life: 3000 })
      }
    }
  })
}

function getRoleSeverity(role: string) {
  switch (role) {
    case 'admin': return 'danger'
    case 'manager': return 'warning'
    case 'accountant': return 'info'
    case 'storekeeper': return 'secondary'
    case 'cashier': return 'success'
    default: return 'info'
  }
}
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900">Users</h1>
      <Button label="Add User" icon="pi pi-plus" @click="openNew" />
    </div>

    <div class="bg-white rounded-lg shadow-sm border border-surface-200">
      <DataTable
        :value="users"
        lazy
        paginator
        :rows="rows"
        :totalRecords="totalRecords"
        :loading="loading"
        @page="loadUsers"
        dataKey="id"
        class="p-datatable-sm"
      >
        <Column field="name" header="Name"></Column>
        <Column field="email" header="Email"></Column>
        <Column field="role" header="Role">
          <template #body="{ data }">
            <Tag :value="data.role" :severity="getRoleSeverity(data.role)" class="capitalize" />
          </template>
        </Column>
        <Column field="is_active" header="Status">
          <template #body="{ data }">
            <Tag :value="data.is_active ? 'Active' : 'Inactive'" :severity="data.is_active ? 'success' : 'danger'" />
          </template>
        </Column>
        <Column header="Actions" :exportable="false" style="min-width: 12rem">
          <template #body="{ data }">
            <Button icon="pi pi-pencil" text rounded @click="editUser(data)" title="Edit" />
            <Button icon="pi pi-key" text rounded severity="warning" @click="openResetPassword(data)" title="Reset Password" />
            <Button 
              :icon="data.is_active ? 'pi pi-ban' : 'pi pi-check'" 
              text 
              rounded 
              :severity="data.is_active ? 'danger' : 'success'" 
              @click="toggleUserStatus(data)" 
              :title="data.is_active ? 'Disable' : 'Enable'" 
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- User Form Dialog -->
    <Dialog v-model:visible="userDialog" :style="{width: '450px'}" :header="isEdit ? 'Edit User' : 'New User'" :modal="true">
      <div class="flex flex-col space-y-4 mt-4">
        <div class="flex flex-col space-y-2">
          <label for="name">Name</label>
          <InputText id="name" v-model.trim="userForm.name" required autofocus />
        </div>
        <div class="flex flex-col space-y-2">
          <label for="email">Email</label>
          <InputText id="email" v-model.trim="userForm.email" required type="email" />
        </div>
        <div class="flex flex-col space-y-2">
          <label for="role">Role</label>
          <Dropdown id="role" v-model="userForm.role" :options="roles" optionLabel="label" optionValue="value" placeholder="Select a Role" />
        </div>
        <div v-if="!isEdit" class="flex flex-col space-y-2">
          <label for="password">Password</label>
          <Password id="password" v-model="userForm.password" required toggleMask :feedback="true" inputClass="w-full" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" icon="pi pi-times" text @click="userDialog = false" />
        <Button label="Save" icon="pi pi-check" @click="saveUser" :loading="submitting" />
      </template>
    </Dialog>

    <!-- Reset Password Dialog -->
    <Dialog v-model:visible="passwordDialog" :style="{width: '400px'}" header="Reset Password" :modal="true">
      <div class="flex flex-col space-y-4 mt-4">
        <div class="flex flex-col space-y-2">
          <label for="new_password">New Password</label>
          <Password id="new_password" v-model="resetPasswordValue" required toggleMask :feedback="true" inputClass="w-full" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" icon="pi pi-times" text @click="passwordDialog = false" />
        <Button label="Save" icon="pi pi-check" @click="savePassword" :loading="submitting" />
      </template>
    </Dialog>
  </div>
</template>
