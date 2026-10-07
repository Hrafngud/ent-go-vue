<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { PhCaretDown, PhCheck } from '@phosphor-icons/vue'

const props = defineProps<{
  id: string
  label: string
  options: { value: string; label: string }[]
  disabled?: boolean
}>()
const model = defineModel<string>({ required: true })
const root = ref<HTMLElement>()
const trigger = ref<HTMLButtonElement>()
const open = ref(false)
const activeIndex = ref(0)
const selected = computed(() => props.options.find(option => option.value === model.value))
let search = ''
let lastTyped = 0

function openMenu() {
  if (props.disabled || !props.options.length) return
  activeIndex.value = Math.max(0, props.options.findIndex(option => option.value === model.value))
  open.value = true
}

function closeMenu() {
  open.value = false
  search = ''
}

function select(index: number) {
  const option = props.options[index]
  if (!option || props.disabled) return
  model.value = option.value
  closeMenu()
  trigger.value?.focus({ preventScroll: true })
}

function move(key: string) {
  const directions = ['ArrowDown', 'ArrowUp', 'Home', 'End']
  if (!directions.includes(key)) return false
  const wasOpen = open.value
  if (!wasOpen) openMenu()
  if (key === 'Home') activeIndex.value = 0
  else if (key === 'End') activeIndex.value = props.options.length - 1
  else if (wasOpen) {
    activeIndex.value = Math.max(0, Math.min(props.options.length - 1, activeIndex.value + (key === 'ArrowDown' ? 1 : -1)))
  }
  return true
}

function typeAhead(key: string) {
  if (!open.value) openMenu()
  const now = Date.now()
  search = now - lastTyped < 600 ? search + key : key
  lastTyped = now
  const term = search.toLocaleLowerCase()
  const index = props.options.findIndex(option => option.label.toLocaleLowerCase().startsWith(term))
  if (index >= 0) activeIndex.value = index
}

function onKeydown(event: KeyboardEvent) {
  if (props.disabled) return
  if (move(event.key)) { event.preventDefault(); return }
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault()
    if (open.value) select(activeIndex.value)
    else openMenu()
  } else if (event.key === 'Escape') {
    if (open.value) { event.preventDefault(); event.stopPropagation() }
    closeMenu()
  } else if (event.key === 'Tab') closeMenu()
  else if (event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey) {
    event.preventDefault()
    typeAhead(event.key)
  }
}

function onPointerdown(event: PointerEvent) {
  if (event.target instanceof Node && !root.value?.contains(event.target)) closeMenu()
}

watch(() => props.disabled, disabled => { if (disabled) closeMenu() })
onMounted(() => document.addEventListener('pointerdown', onPointerdown))
onBeforeUnmount(() => document.removeEventListener('pointerdown', onPointerdown))
</script>

<template>
  <div ref="root" class="dropdown w-full" :class="open ? 'dropdown-open' : 'dropdown-close'">
    <span :id="`${id}-label`" class="sr-only">{{ label }}</span>
    <button :id="id" ref="trigger" type="button" class="select flex w-full cursor-pointer items-center justify-between gap-3 bg-none text-left"
      role="combobox" aria-haspopup="listbox" :aria-labelledby="`${id}-label`" :aria-expanded="open" :aria-controls="`${id}-options`"
      :aria-activedescendant="open ? `${id}-option-${activeIndex}` : undefined" :disabled="disabled"
      @click="open ? closeMenu() : openMenu()" @keydown="onKeydown" @blur="closeMenu">
      <span class="truncate">{{ selected?.label || 'Select an option' }}</span>
      <PhCaretDown :size="16" class="shrink-0 text-base-content/60" :class="{ 'rotate-180': open }" aria-hidden="true" />
    </button>
    <ul v-show="open" :id="`${id}-options`" class="dropdown-content menu z-20 mt-2 w-full rounded-field border border-base-300 bg-base-200 p-2 shadow-lg"
      role="listbox" :aria-labelledby="`${id}-label`">
      <li v-for="(option, index) in options" :id="`${id}-option-${index}`" :key="option.value" role="option" :aria-selected="option.value === model"
        @pointerdown.prevent @pointermove="activeIndex = index" @click="select(index)">
        <span class="flex cursor-pointer items-center justify-between gap-3 rounded-field" :class="{ 'menu-focus': activeIndex === index }">
          {{ option.label }}<PhCheck :size="16" class="shrink-0" :class="{ 'invisible': option.value !== model }" aria-hidden="true" />
        </span>
      </li>
    </ul>
  </div>
</template>
