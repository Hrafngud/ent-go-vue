import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useFeedbackStore = defineStore('feedback', () => {
  const message = ref('')
  function clear() { message.value = '' }
  return { message, clear }
})
