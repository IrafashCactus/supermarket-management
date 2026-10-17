<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { categoryService } from '@/services/categoryService'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '@/types/category'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'

const toast = useToast()
const confirm = useConfirm()

const categories = ref<Category[]>([])
const loading = ref(false)
const totalRecords = ref(0)
const rows = ref(10)

const categoryDialog = ref(false)
const isEdit = ref(false)
const categoryForm = ref<any>({ name: '', description: '' })
const selectedCategoryId = ref('')
const submitting = ref(false)

async function loadCategories(event?: any) {
  loading.value = true
  try {
    const page = event ? (event.first / event.rows) + 1 : 1
    const response = await categoryService.getCategories(page, rows.value)
    categories.value = response.data
    totalRecords.value = response.meta.total
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load categories', life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategories()
})

function openNew() {
  categoryForm.value = { name: '', description: '' }
  isEdit.value = false
  categoryDialog.value = true
}

function editCategory(category: Category) {
  categoryForm.value = { ...category }
  selectedCategoryId.value = category.id
  isEdit.value = true
  categoryDialog.value = true
}

async function saveCategory() {
  submitting.value = true
  try {
    if (isEdit.value) {
      const payload: UpdateCategoryRequest = {
        name: categoryForm.value.name,
        description: categoryForm.value.description
      }
      await categoryService.updateCategory(selectedCategoryId.value, payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Category updated', life: 3000 })
    } else {
      const payload: CreateCategoryRequest = {
        name: categoryForm.value.name,
        description: categoryForm.value.description
      }
      await categoryService.createCategory(payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Category created', life: 3000 })
    }
    categoryDialog.value = false
    loadCategories()
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.error?.message || 'Failed to save category', life: 3000 })
  } finally {
    submitting.value = false
  }
}

function deleteCategory(category: Category) {
  if (category.product_count > 0) {
    toast.add({ severity: 'error', summary: 'Cannot Delete', detail: `Category "${category.name}" has ${category.product_count} product(s). Remove all products first.`, life: 5000 })
    return
  }

  confirm.require({
    message: `Are you sure you want to delete "${category.name}"?`,
    header: 'Confirm Delete',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await categoryService.deleteCategory(category.id)
        toast.add({ severity: 'success', summary: 'Success', detail: 'Category deleted', life: 3000 })
        loadCategories()
      } catch (error: any) {
        toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.error?.message || 'Failed to delete category', life: 3000 })
      }
    }
  })
}
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900">Categories</h1>
      <Button label="Add Category" icon="pi pi-plus" @click="openNew" />
    </div>

    <div class="bg-white rounded-lg shadow-sm border border-surface-200">
      <DataTable
        :value="categories"
        lazy
        paginator
        :rows="rows"
        :totalRecords="totalRecords"
        :loading="loading"
        @page="loadCategories"
        dataKey="id"
        class="p-datatable-sm"
      >
        <Column field="name" header="Name"></Column>
        <Column field="description" header="Description"></Column>
        <Column field="product_count" header="Product Count"></Column>
        <Column header="Actions" :exportable="false" style="min-width: 8rem">
          <template #body="{ data }">
            <Button icon="pi pi-pencil" text rounded @click="editCategory(data)" title="Edit" />
            <Button icon="pi pi-trash" text rounded severity="danger" @click="deleteCategory(data)" title="Delete" />
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- Category Form Dialog -->
    <Dialog v-model:visible="categoryDialog" :style="{width: '450px'}" :header="isEdit ? 'Edit Category' : 'New Category'" :modal="true">
      <div class="flex flex-col space-y-4 mt-4">
        <div class="flex flex-col space-y-2">
          <label for="name">Name</label>
          <InputText id="name" v-model.trim="categoryForm.name" required autofocus />
        </div>
        <div class="flex flex-col space-y-2">
          <label for="description">Description</label>
          <Textarea id="description" v-model="categoryForm.description" rows="3" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" icon="pi pi-times" text @click="categoryDialog = false" />
        <Button label="Save" icon="pi pi-check" @click="saveCategory" :loading="submitting" />
      </template>
    </Dialog>
  </div>
</template>
