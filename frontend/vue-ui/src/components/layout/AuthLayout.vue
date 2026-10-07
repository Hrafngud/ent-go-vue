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
  position: relative;
  isolation: isolate;
  overflow: clip;
}

.auth-split .auth-header {
  position: absolute;
  z-index: 2;
  top: 0;
  left: 0;
  padding: 2rem var(--auth-gutter);
}

.auth-composition {
  position: relative;
  display: grid;
  grid-template-columns: minmax(0, 1fr) clamp(24rem, 36vw, 38rem);
  flex: 1;
  width: 100%;
  min-height: 100dvh;
}

.auth-intro {
  position: relative;
  z-index: 1;
  min-width: 0;
  padding: clamp(8rem, 16vh, 12rem) var(--auth-gutter) 3rem;
}

.auth-rule {
  display: block;
  width: 3rem;
  height: 2px;
  margin-bottom: 2rem;
  background: var(--color-primary);
}

.auth-intro .auth-title {
  font-family: var(--font-sans);
  font-size: clamp(2.25rem, 5.2vw, 6.5rem);
  font-style: normal;
  font-weight: 500;
  line-height: 1.08;
  letter-spacing: -0.055em;
  white-space: nowrap;
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
  inset-block: 0;
  left: 12%;
  width: 64%;
  height: 100%;
  color: var(--color-primary);
  pointer-events: none;
}

.auth-panel {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  width: 100%;
  min-height: 100dvh;
  padding: clamp(2rem, 4vw, 4rem);
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
  .auth-panel {
    padding: 1.75rem;
  }
}

@media (width < 48rem) {
  .auth-split .auth-header {
    padding-block: 1.5rem;
  }

  .auth-composition {
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: auto 1fr;
  }

  .auth-intro {
    padding: 7rem var(--auth-gutter) 3rem;
  }

  .auth-intro .auth-title {
    font-size: clamp(1.75rem, 7.8vw, 3.5rem);
  }

  .auth-rule {
    margin-bottom: 1.5rem;
  }

  .auth-description {
    margin-top: 1.25rem;
  }

  .hexagon-mesh {
    left: 0;
    width: 100%;
    height: 100dvh;
    opacity: 0.6;
  }

  .auth-panel {
    min-height: auto;
    padding: 3rem var(--auth-gutter);
  }
}
</style>
