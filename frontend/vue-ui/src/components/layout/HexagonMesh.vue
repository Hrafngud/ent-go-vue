<script setup lang="ts">
// Share edges so every line in the honeycomb has the same visual weight.
const edges = new Map<string, string>()
const accents: string[] = []
const radius = 38
const rowHeight = Math.sqrt(3) * radius

function project(x: number, y: number) {
  return [x, y + 0.0009 * (x - 360) ** 2 - 70].map(value => value.toFixed(2)).join(',')
}

for (let column = 0; column < 14; column++) {
  for (let row = 0; row < 12; row++) {
    const x = column * radius * 1.5
    const y = row * rowHeight + (column % 2) * rowHeight / 2 - 50
    const vertices = Array.from({ length: 6 }, (_, corner) => {
      const angle = corner * Math.PI / 3
      return project(x + radius * Math.cos(angle), y + radius * Math.sin(angle))
    })
    vertices.forEach((start, corner) => {
      const end = vertices[(corner + 1) % 6]!
      edges.set([start, end].sort().join('|'), `M${start}L${end}`)
    })
    if ((column === 5 && row === 4) || (column === 8 && row === 6) || (column === 6 && row === 8)) {
      accents.push(`M${vertices.join('L')}Z`)
    }
  }
}

const mesh = [...edges.values()].join('')
</script>

<template>
  <svg class="hexagon-mesh" viewBox="0 0 740 720" fill="none" aria-hidden="true" focusable="false">
    <defs>
      <radialGradient id="login-mesh-fade" cx="48%" cy="48%" r="52%">
        <stop offset="0" stop-color="white" />
        <stop offset="0.5" stop-color="white" stop-opacity="0.7" />
        <stop offset="1" stop-color="white" stop-opacity="0" />
      </radialGradient>
      <mask id="login-mesh-mask">
        <rect width="740" height="720" fill="url(#login-mesh-fade)" />
      </mask>
    </defs>
    <g mask="url(#login-mesh-mask)" transform="rotate(-18 370 360)">
      <path :d="mesh" stroke="currentColor" stroke-width="0.8" opacity="0.18" />
      <path v-for="cell in accents" :key="cell" :d="cell" fill="currentColor" fill-opacity="0.035" stroke="currentColor" stroke-width="1.1" opacity="0.5" />
    </g>
  </svg>
</template>
