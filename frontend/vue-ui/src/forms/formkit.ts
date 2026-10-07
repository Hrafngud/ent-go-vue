import { isNode, type FormKitNode, type FormKitPlugin } from '@formkit/core'
import { emailError, nameError, normalizeAccount, normalizeLogin, passwordError } from '../validation/account'

// Normalize only declared identity fields. Passwords always retain their exact bytes.
function sanitizeValue(node: FormKitNode, value: unknown): unknown {
  if (typeof value !== 'string') return value
  if (node.props.sanitize === 'name') return normalizeAccount({ name: value, email: '' }).name
  if (node.props.sanitize === 'email') return normalizeLogin(value)
  return value
}

function sanitizeValues(node: FormKitNode, values: Record<string, unknown>): Record<string, unknown> {
  const result = { ...values }
  for (const child of node.children) {
    if (!isNode(child)) continue
    const value = result[child.name]
    result[child.name] = child.type === 'group' && value && typeof value === 'object'
      ? sanitizeValues(child, value as Record<string, unknown>)
      : sanitizeValue(child, value)
  }
  return result
}

export const sanitizePlugin: FormKitPlugin = (node) => {
  node.addProps(['sanitize'])
  if (node.type === 'input') {
    node.hook.input((value, next) => next(
      typeof value === 'string' && node.props.sanitize === 'name' ? value.normalize('NFC') : value,
    ))
  }
  if (node.props.type === 'form') {
    node.hook.submit((values, next) => next(sanitizeValues(node, values)))
  }
}

export const accountRules = {
  account_name: (node: FormKitNode) => !nameError(String(node.value ?? '')),
  account_email: (node: FormKitNode) => !emailError(String(node.value ?? '')),
  password_bytes: (node: FormKitNode, minimum = 8) => !passwordError(String(node.value ?? ''), Number(minimum)),
}

const sectionClasses: Record<string, string> = {
  outer: 'space-y-2',
  label: 'block text-sm font-medium',
  inner: 'relative',
  help: 'text-xs leading-relaxed text-base-content/60',
  messages: 'space-y-1',
  message: 'text-sm text-error',
}

export const formkitOptions = {
  plugins: [sanitizePlugin],
  rules: accountRules,
  messages: {
    en: {
      validation: {
        account_name: () => 'Enter a plain-text name up to 255 bytes, without angle brackets or control characters.',
        account_email: () => 'Enter a valid email address, such as you@example.com (maximum 254 bytes).',
        password_bytes: ({ node, args }: { node: FormKitNode; args: unknown[] }) => passwordError(String(node.value ?? ''), Number(args[0] ?? 8)),
        confirm: () => 'Your passwords do not match.',
      },
    },
  },
  config: {
    validationVisibility: 'blur',
    rootClasses(section: string, node: FormKitNode) {
      const classes = node.props.type === 'form'
        ? ({ message: 'alert alert-error text-sm', messages: 'space-y-2' } as Record<string, string>)[section] || ''
        : sectionClasses[section] || ''
      const input = section === 'input' ? (node.props.type === 'select' ? 'select w-full' : 'input w-full') : ''
      return Object.fromEntries([`formkit-${section}`, ...`${classes} ${input}`.split(' ').filter(Boolean)].map(name => [name, true]))
    },
  },
}
