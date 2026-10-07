<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, useId } from 'vue'
import { PhX } from '@phosphor-icons/vue'

const props = defineProps<{ title: string; description?: string; busy?: boolean }>()
const emit = defineEmits<{ close: [] }>()
const dialog = ref<HTMLDialogElement>()
const heading = ref<HTMLElement>()
const titleId = useId()
const descriptionId = useId()
let opener: HTMLElement | null = null

function close() { if (!props.busy) emit('close') }
function trapFocus(event: KeyboardEvent) {
  if (event.key !== 'Tab') return
  const controls = [...(dialog.value?.querySelectorAll<HTMLElement>('button:not(:disabled), input:not(:disabled), select:not(:disabled), textarea:not(:disabled), a[href], [tabindex="0"]') || [])]
    .filter(control => control.getClientRects().length > 0)
  const active = document.activeElement
  const boundary = event.shiftKey ? controls[0] : controls.at(-1)
  if (active !== boundary && active !== heading.value) return
  event.preventDefault()
  const target = event.shiftKey ? controls.at(-1) : controls[0]
  target?.focus()
}
onMounted(() => {
  opener = document.activeElement instanceof HTMLElement ? document.activeElement : null
  dialog.value?.showModal()
  heading.value?.focus({ preventScroll: true })
})
onBeforeUnmount(() => {
  dialog.value?.close()
  if (opener?.isConnected) opener.focus({ preventScroll: true })
})
</script>

<template>
  <Teleport to="body">
    <dialog ref="dialog" class="modal user-modal" :aria-labelledby="titleId" :aria-describedby="description ? descriptionId : undefined"
      :aria-busy="busy" @cancel.prevent="close" @keydown="trapFocus">
      <div class="modal-box w-full max-w-2xl border border-base-300 bg-base-200 p-6 sm:p-8">
        <header class="mb-8 flex items-start justify-between gap-4">
          <div class="min-w-0">
            <p class="mb-2 text-xs font-medium tracking-widest text-primary uppercase">Workspace accounts</p>
            <h2 :id="titleId" ref="heading" tabindex="-1" class="text-2xl font-semibold tracking-tight">{{ title }}</h2>
            <p v-if="description" :id="descriptionId" class="mt-2 text-sm leading-relaxed text-base-content/65">{{ description }}</p>
          </div>
          <button type="button" class="btn btn-ghost btn-circle btn-sm shrink-0" aria-label="Close dialog" :disabled="busy" @click="close">
            <PhX :size="20" aria-hidden="true" />
          </button>
        </header>
        <slot />
      </div>
      <div class="modal-backdrop" aria-hidden="true" @click="close"></div>
    </dialog>
  </Teleport>
</template>
