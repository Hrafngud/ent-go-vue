<script setup lang="ts">
defineProps<{ label: string }>()
</script>

<template>
  <nav class="breadcrumbs mb-6 text-sm text-base-content/70" aria-label="Breadcrumb">
    <ul><li><RouterLink to="/admin/users">Users</RouterLink></li><li aria-current="page">{{ label }}</li></ul>
  </nav>
</template>
