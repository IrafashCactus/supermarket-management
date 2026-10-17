<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { productService } from '@/services/productService'
import { categoryService } from '@/services/categoryService'
import type { Product, CreateProductRequest, UpdateProductRequest, ProductListParams } from '@/types/product'
import type { Category } from '@/types/category'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Dropdown from 'primevue/dropdown'
import Tag from 'primevue/tag'
import Checkbox from 'primevue/checkbox'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'

const toast = useToast()
const confirm = useConfirm()
const authStore = useAuthStore()
const canManageProducts = computed(() => authStore.checkRole(['admin', 'manager']))

const products = ref<Product[]>([])
const loading = ref(false)
const totalRecords = ref(0)
const rows = ref(10)
const first = ref(0)

// Filters
const searchQuery = ref('')
const filterCategoryId = ref<string | null>(null)
const filterStatus = ref<boolean | null>(null)

const statusOptions = [
  { label: 'All', value: null },
  { label: 'Active', value: true },
  { label: 'Inactive', value: false }
]

const unitOfMeasureOptions = [
  { label: 'Piece', value: 'piece' },
  { label: 'Kg', value: 'kg' },
  { label: 'Litre', value: 'litre' },
  { label: 'Pack', value: 'pack' }
]

// Categories for dropdown
const categories = ref<Category[]>([])
const categoryOptions = ref<{ label: string; value: string | null }[]>([])

// Dialog
const productDialog = ref(false)
const isEdit = ref(false)
const productForm = ref<any>({
  name: '',
  barcode: '',
  sku: '',
  category_id: '',
  unit_price: 0,
  cost_price: 0,
  tax_rate: 18,
  unit_of_measure: 'piece',
  reorder_level: 0,
  image_url: '',
  is_active: true
})
const submitting = ref(false)
const selectedProductId = ref('')

let searchTimeout: ReturnType<typeof setTimeout> | null = null

async function loadCategories() {
  try {
    categories.value = await categoryService.getAllCategories()
    categoryOptions.value = [
      { label: 'All Categories', value: null },
      ...categories.value.map(c => ({ label: c.name, value: c.id }))
    ]
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load categories', life: 3000 })
  }
}

async function loadProducts(event?: any) {
  loading.value = true
  try {
    const page = event ? (event.first / event.rows) + 1 : 1
    const params: ProductListParams = {
      page,
      per_page: rows.value,
      search: searchQuery.value || undefined,
      category_id: filterCategoryId.value || undefined,
      is_active: filterStatus.value
    }
    const response = await productService.getProducts(params)
    products.value = response.data
    totalRecords.value = response.meta.total
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load products', life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategories()
  loadProducts()
})

// Debounced search
watch(searchQuery, () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    first.value = 0
    loadProducts()
  }, 400)
})

watch([filterCategoryId, filterStatus], () => {
  first.value = 0
  loadProducts()
})

function formatPrice(value: number): string {
  return `${Math.round(value).toLocaleString()} RWF`
}

function openNew() {
  productForm.value = {
    name: '',
    barcode: '',
    sku: '',
    category_id: '',
    unit_price: 0,
    cost_price: 0,
    tax_rate: 18,
    unit_of_measure: 'piece',
    reorder_level: 0,
    image_url: '',
    is_active: true
  }
  isEdit.value = false
  productDialog.value = true
}

function editProduct(product: Product) {
  productForm.value = { ...product }
  selectedProductId.value = product.id
  isEdit.value = true
  productDialog.value = true
}

async function saveProduct() {
  submitting.value = true
  try {
    if (isEdit.value) {
      const payload: UpdateProductRequest = {
        name: productForm.value.name,
        barcode: productForm.value.barcode,
        sku: productForm.value.sku,
        category_id: productForm.value.category_id,
        unit_price: productForm.value.unit_price,
        cost_price: productForm.value.cost_price,
        tax_rate: productForm.value.tax_rate,
        unit_of_measure: productForm.value.unit_of_measure,
        reorder_level: productForm.value.reorder_level,
        image_url: productForm.value.image_url,
        is_active: productForm.value.is_active
      }
      await productService.updateProduct(selectedProductId.value, payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Product updated', life: 3000 })
    } else {
      const payload: CreateProductRequest = {
        name: productForm.value.name,
        barcode: productForm.value.barcode,
        sku: productForm.value.sku,
        category_id: productForm.value.category_id,
        unit_price: productForm.value.unit_price,
        cost_price: productForm.value.cost_price,
        tax_rate: productForm.value.tax_rate,
        unit_of_measure: productForm.value.unit_of_measure,
        reorder_level: productForm.value.reorder_level,
        image_url: productForm.value.image_url
      }
      await productService.createProduct(payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Product created', life: 3000 })
    }
    productDialog.value = false
    loadProducts()
  } catch (error: any) {
    const status = error.response?.status
    const message = error.response?.data?.error?.message || 'Failed to save product'
    if (status === 409) {
      toast.add({ severity: 'error', summary: 'Duplicate Barcode', detail: 'A product with this barcode already exists. Please use a unique barcode.', life: 5000 })
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 })
    }
  } finally {
    submitting.value = false
  }
}

async function toggleProductStatus(product: Product) {
  const newStatus = !product.is_active
  const action = newStatus ? 'activate' : 'deactivate'

  confirm.require({
    message: `Are you sure you want to ${action} "${product.name}"?`,
    header: 'Confirm',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        const payload: UpdateProductRequest = {
          name: product.name,
          barcode: product.barcode,
          sku: product.sku,
          category_id: product.category_id,
          unit_price: product.unit_price,
          cost_price: product.cost_price,
          tax_rate: product.tax_rate,
          unit_of_measure: product.unit_of_measure,
          reorder_level: product.reorder_level,
          image_url: product.image_url,
          is_active: newStatus
        }
        await productService.updateProduct(product.id, payload)
        toast.add({ severity: 'success', summary: 'Success', detail: `Product ${action}d`, life: 3000 })
        loadProducts()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to change status', life: 3000 })
      }
    }
  })
}
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900">Products</h1>
      <Button v-if="canManageProducts" label="Add Product" icon="pi pi-plus" @click="openNew" />
    </div>

    <!-- Filter Bar -->
    <div class="flex flex-wrap gap-4 mb-4">
      <div class="flex-1 min-w-[200px]">
        <InputText v-model="searchQuery" placeholder="Search products..." class="w-full" />
      </div>
      <Dropdown
        v-model="filterCategoryId"
        :options="categoryOptions"
        optionLabel="label"
        optionValue="value"
        placeholder="All Categories"
        class="w-48"
      />
      <Dropdown
        v-model="filterStatus"
        :options="statusOptions"
        optionLabel="label"
        optionValue="value"
        placeholder="All Status"
        class="w-40"
      />
    </div>

    <div class="bg-white rounded-lg shadow-sm border border-surface-200">
      <DataTable
        :value="products"
        lazy
        paginator
        :rows="rows"
        :first="first"
        :totalRecords="totalRecords"
        :loading="loading"
        @page="loadProducts"
        dataKey="id"
        class="p-datatable-sm"
      >
        <Column field="name" header="Name"></Column>
        <Column field="barcode" header="Barcode"></Column>
        <Column field="category_name" header="Category"></Column>
        <Column field="unit_price" header="Unit Price">
          <template #body="{ data }">
            {{ formatPrice(data.unit_price) }}
          </template>
        </Column>
        <Column field="cost_price" header="Cost Price">
          <template #body="{ data }">
            {{ formatPrice(data.cost_price) }}
          </template>
        </Column>
        <Column field="tax_rate" header="Tax Rate">
          <template #body="{ data }">
            {{ data.tax_rate }}%
          </template>
        </Column>
        <Column field="unit_of_measure" header="UOM">
          <template #body="{ data }">
            <span class="capitalize">{{ data.unit_of_measure }}</span>
          </template>
        </Column>
        <Column field="is_active" header="Status">
          <template #body="{ data }">
            <Tag :value="data.is_active ? 'Active' : 'Inactive'" :severity="data.is_active ? 'success' : 'danger'" />
          </template>
        </Column>
        <Column v-if="canManageProducts" header="Actions" :exportable="false" style="min-width: 8rem">
          <template #body="{ data }">
            <Button icon="pi pi-pencil" text rounded @click="editProduct(data)" title="Edit" />
            <Button
              :icon="data.is_active ? 'pi pi-ban' : 'pi pi-check'"
              text
              rounded
              :severity="data.is_active ? 'danger' : 'success'"
              @click="toggleProductStatus(data)"
              :title="data.is_active ? 'Deactivate' : 'Activate'"
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- Product Form Dialog -->
    <Dialog v-model:visible="productDialog" :style="{width: '550px'}" :header="isEdit ? 'Edit Product' : 'New Product'" :modal="true">
      <div class="flex flex-col space-y-4 mt-4">
        <div class="flex flex-col space-y-2">
          <label for="prod_name">Name</label>
          <InputText id="prod_name" v-model.trim="productForm.name" required autofocus />
        </div>
        <div class="flex flex-col space-y-2">
          <label for="prod_barcode">Barcode</label>
          <InputText id="prod_barcode" v-model.trim="productForm.barcode" required />
        </div>
        <div class="flex flex-col space-y-2">
          <label for="prod_sku">SKU</label>
          <InputText id="prod_sku" v-model.trim="productForm.sku" />
        </div>
        <div class="flex flex-col space-y-2">
          <label for="prod_category">Category</label>
          <Dropdown
            id="prod_category"
            v-model="productForm.category_id"
            :options="categories"
            optionLabel="name"
            optionValue="id"
            placeholder="Select a Category"
          />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="flex flex-col space-y-2">
            <label for="prod_unit_price">Unit Price</label>
            <InputNumber id="prod_unit_price" v-model="productForm.unit_price" suffix=" RWF" :min="0" />
          </div>
          <div class="flex flex-col space-y-2">
            <label for="prod_cost_price">Cost Price</label>
            <InputNumber id="prod_cost_price" v-model="productForm.cost_price" suffix=" RWF" :min="0" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="flex flex-col space-y-2">
            <label for="prod_tax_rate">Tax Rate</label>
            <InputNumber id="prod_tax_rate" v-model="productForm.tax_rate" suffix="%" :min="0" :max="100" />
          </div>
          <div class="flex flex-col space-y-2">
            <label for="prod_uom">Unit of Measure</label>
            <Dropdown
              id="prod_uom"
              v-model="productForm.unit_of_measure"
              :options="unitOfMeasureOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="Select UOM"
            />
          </div>
        </div>
        <div class="flex flex-col space-y-2">
          <label for="prod_reorder">Reorder Level</label>
          <InputNumber id="prod_reorder" v-model="productForm.reorder_level" :min="0" />
        </div>
        <div v-if="isEdit" class="flex items-center space-x-2">
          <Checkbox id="prod_active" v-model="productForm.is_active" :binary="true" />
          <label for="prod_active">Active</label>
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" icon="pi pi-times" text @click="productDialog = false" />
        <Button label="Save" icon="pi pi-check" @click="saveProduct" :loading="submitting" />
      </template>
    </Dialog>
  </div>
</template>
