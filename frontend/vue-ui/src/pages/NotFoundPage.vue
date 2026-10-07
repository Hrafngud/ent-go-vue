<script setup lang="ts">
import AuthLayout from '../components/layout/AuthLayout.vue'
</script>

<template>
  <AuthLayout title="Page not found." description="This address does not lead to a page in your workspace.">
    <RouterLink class="btn btn-primary" to="/">Go to workspace</RouterLink>
  </AuthLayout>
</template>
