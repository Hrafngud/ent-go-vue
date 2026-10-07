<script setup lang="ts">
defineProps<{ label: string }>()
</script>

<template>
  <nav class="breadcrumbs mb-5 py-0 text-xs text-base-content/60" aria-label="Breadcrumb">
    <ul><li><RouterLink to="/admin/users">Users</RouterLink></li><li aria-current="page">{{ label }}</li></ul>
  </nav>
</template>
