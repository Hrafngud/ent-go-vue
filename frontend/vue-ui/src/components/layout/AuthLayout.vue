<script setup lang="ts">
import BrandMark from './BrandMark.vue'
import HexagonMesh from './HexagonMesh.vue'
defineProps<{ title: string; description?: string; eyebrow?: string; split?: boolean }>()
</script>

<template>
  <div class="flex min-h-dvh flex-col" :class="{ 'auth-split': split }">
    <header class="auth-header px-5 py-6 sm:px-10 sm:py-8"><BrandMark /></header>
    <main id="main-content" :class="split ? 'auth-composition' : 'flex flex-1 items-center justify-center px-5 py-6 sm:py-10'">
      <div v-if="split" class="auth-intro page-enter">
        <span class="auth-rule" aria-hidden="true"></span>
        <h1 id="auth-title" tabindex="-1" class="auth-title">{{ title }}</h1>
        <p v-if="description" class="auth-description">{{ description }}</p>
      </div>
      <HexagonMesh v-if="split" />
      <section :class="split ? 'auth-panel page-enter' : 'surface page-enter w-full max-w-md p-6 sm:p-10'" aria-labelledby="auth-title">
        <p v-if="eyebrow" class="mb-5 text-xs font-medium tracking-widest text-primary uppercase">{{ eyebrow }}</p>
        <h1 v-if="!split" id="auth-title" tabindex="-1" class="auth-title">{{ title }}</h1>
        <p v-if="description && !split" class="mt-3 text-sm leading-relaxed text-base-content/65">{{ description }}</p>
        <div :class="{ 'mt-8': !split }"><slot /></div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.auth-split {
  --auth-gutter: clamp(1.25rem, 5vw, 6rem);
  overflow: clip;
}

.auth-split .auth-header {
  padding: 2rem var(--auth-gutter);
}

.auth-composition {
  position: relative;
  isolation: isolate;
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(0, 1fr);
  align-content: start;
  gap: clamp(3rem, 8vw, 9rem);
  flex: 1;
  width: 100%;
  max-width: 100rem;
  margin-inline: auto;
  padding: clamp(3rem, 6vh, 5rem) var(--auth-gutter) 5rem;
}

.auth-intro {
  position: relative;
  z-index: 1;
}

.auth-rule {
  display: block;
  width: 3rem;
  height: 2px;
  margin-bottom: 2rem;
  background: var(--color-primary);
}

.auth-intro .auth-title {
  max-width: 7ch;
  font-size: clamp(5rem, 8.5vw, 8.5rem);
  line-height: 0.96;
  letter-spacing: -0.045em;
}

.auth-description {
  margin-top: 1.75rem;
  color: color-mix(in oklch, var(--color-base-content) 65%, var(--color-base-100));
  font-size: 1rem;
  line-height: 1.6;
}

.hexagon-mesh {
  position: absolute;
  z-index: -1;
  top: 3rem;
  left: 23%;
  width: min(58vw, 54rem);
  color: var(--color-primary);
  pointer-events: none;
}

.auth-panel {
  position: relative;
  width: 100%;
  margin-top: clamp(8rem, 12vw, 12rem);
  padding: clamp(1.75rem, 3vw, 3rem);
  border: 1px solid var(--color-base-300);
  border-radius: 1rem;
  background: var(--color-base-200);
}

.auth-panel :deep(.input),
.auth-panel :deep(.btn-primary) {
  min-height: 3.25rem;
}

.auth-panel :deep(.input) {
  border-width: 1px;
  background: var(--color-base-100);
}

@media (width < 64rem) {
  .auth-composition {
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
    gap: 2rem;
  }

  .auth-intro .auth-title {
    font-size: clamp(4.5rem, 9vw, 6rem);
  }

  .auth-panel {
    margin-top: 8rem;
    padding: 1.75rem;
  }
}

@media (width < 48rem) {
  .auth-split .auth-header {
    padding-block: 1.5rem;
  }

  .auth-composition {
    grid-template-columns: minmax(0, 1fr);
    gap: 3rem;
    padding-block: 2rem 3rem;
  }

  .auth-intro .auth-title {
    font-size: clamp(4rem, 14vw, 6rem);
  }

  .auth-rule {
    margin-bottom: 1.5rem;
  }

  .auth-description {
    margin-top: 1.25rem;
  }

  .hexagon-mesh {
    top: -2rem;
    left: 24%;
    width: 32rem;
    opacity: 0.6;
  }

  .auth-panel {
    margin-top: 0;
    padding: 1.5rem;
  }
}
</style>
