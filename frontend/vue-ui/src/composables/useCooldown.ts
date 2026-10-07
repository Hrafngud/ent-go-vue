import { onUnmounted, ref } from 'vue'
import { ApiError } from '../api/client'

export function useCooldown() {
  const seconds = ref(0)
  let timer: ReturnType<typeof setInterval> | undefined
  function start(error: unknown) {
    if (!(error instanceof ApiError) || error.status !== 429) return
    if (timer) clearInterval(timer)
    const end = Date.now() + error.retryAfter * 1000
    seconds.value = error.retryAfter
    timer = setInterval(() => {
      seconds.value = Math.max(0, Math.ceil((end - Date.now()) / 1000))
      if (!seconds.value) { clearInterval(timer); timer = undefined }
    }, 1000)
  }
  onUnmounted(() => { if (timer) clearInterval(timer) })
  return { seconds, start }
}
