import type { UserInput } from '../api/users'

const encoder = new TextEncoder()
const emailPattern = /^[A-Za-z0-9.!#$%&'*+/=?^_`{|}~-]+@[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?)+$/
const trim = (value: string) => value.replace(/^\p{White_Space}+|\p{White_Space}+$/gu, '')
const validUnicode = (value: string) => !/[\uD800-\uDFFF]/u.test(value)

export function normalizeAccount(input: UserInput): UserInput {
  return { name: trim(input.name).normalize('NFC'), email: trim(input.email), ...(input.password ? { password: input.password } : {}) }
}

export function emailError(value: string): string {
  const email = trim(value)
  const [local = '', domain = ''] = email.split('@')
  const valid = encoder.encode(email).length <= 254 && emailPattern.test(email)
    && local.length <= 64 && !local.startsWith('.') && !local.endsWith('.') && !local.includes('..')
    && domain.split('.').every((label) => label.length <= 63)
  return valid ? '' : 'Enter a valid email address, such as you@example.com (maximum 254 bytes).'
}

export function passwordError(password: string, minimum = 8): string {
  const bytes = encoder.encode(password).length
  if (!validUnicode(password) || bytes < minimum || bytes > 72 || !trim(password) || password.includes('\0')) {
    return `Use a password between ${minimum} and 72 bytes, with at least one non-space character. Some characters count as more than one byte.`
  }
  return ''
}

export function accountError(input: UserInput, editing = false): string {
  return nameError(input.name) || emailError(input.email) || ((!editing || input.password) ? passwordError(input.password || '') : '')
}

export function nameError(value: string): string {
  const name = trim(value).normalize('NFC')
  if (!validUnicode(value) || !name || encoder.encode(name).length > 255 || /[<>\p{Cc}]/u.test(name)) {
    return 'Enter a plain-text name up to 255 bytes, without angle brackets or control characters.'
  }
  return ''
}

export function normalizeLogin(email: string): string { return trim(email) }
