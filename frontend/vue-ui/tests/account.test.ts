import assert from 'node:assert/strict'
import test from 'node:test'
import { accountError, emailError, normalizeAccount, passwordError } from '../src/validation/account.ts'
import { ApiError } from '../src/api/client.ts'

test('normalizes names and email without changing passwords or email case', () => {
  assert.deepEqual(normalizeAccount({ name: ' Jose\u0301 ', email: ' Person@example.com ', password: ' password ' }), {
    name: 'José', email: 'Person@example.com', password: ' password ',
  })
  assert.equal(accountError({ name: "O'Connor-José", email: 'person@example.com', password: 'password' }), '')
})

test('rejects malformed names and enforces byte lengths', () => {
  for (const name of [' ', '<script>alert(1)</script>', 'Name\0', 'Name\nName', 'é'.repeat(128), '\ud800']) {
    assert.ok(accountError({ name, email: 'person@example.com', password: 'password' }), name)
  }
})

test('validates email structure and lengths', () => {
  for (const email of ['invalid', 'a@localhost', 'Person <a@example.com>', 'a..b@example.com', '.a@example.com', 'a@-example.com', `${'a'.repeat(65)}@example.com`, `a@${'a'.repeat(64)}.com`]) {
    assert.ok(emailError(email), email)
  }
  assert.equal(emailError(" O'Connor+tag@example.com "), '')
})

test('validates password bytes, preserves optional edit and legacy login passwords', () => {
  assert.equal(passwordError('😀😀'), '')
  assert.equal(passwordError('a'.repeat(72)), '')
  for (const password of ['seven77', '        ', 'password\0', '😀'.repeat(19), 'a'.repeat(73), 'password\ud800']) {
    assert.ok(passwordError(password))
  }
  const input = { name: 'Name', email: 'a@example.com' }
  assert.equal(accountError(input, true), '')
  assert.ok(accountError(input))
  assert.equal(passwordError('old123', 1), '')
})

test('bounds and defaults retry delays from untrusted responses', () => {
  assert.equal(new ApiError(429, 12.2).retryAfter, 13)
  assert.equal(new ApiError(429, Number.NaN).retryAfter, 60)
  assert.equal(new ApiError(429, 0).retryAfter, 60)
  assert.equal(new ApiError(429, 100000).retryAfter, 3600)
})
