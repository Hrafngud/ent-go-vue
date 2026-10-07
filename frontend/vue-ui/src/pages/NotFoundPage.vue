<script setup lang="ts">
import AuthLayout from '../components/layout/AuthLayout.vue'
import { PhArrowLeft } from '@phosphor-icons/vue'
</script>

<template>
  <AuthLayout title="Page not found." description="This page may have moved or no longer exists." eyebrow="404">
    <RouterLink class="btn btn-primary w-full" to="/"><PhArrowLeft :size="20" aria-hidden="true" />Back to workspace</RouterLink>
  </AuthLayout>
</template>
