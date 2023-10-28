<script setup>
import { RouterLink } from 'vue-router'

import WelcomeItem from './WelcomeItem.vue'
import ServerIcon from './icons/IconServer.vue'
import AlpineIcon from './icons/IconAlpine.vue'
</script>

<template>
  <WelcomeItem>
    <template #icon>
      <ServerIcon />
    </template>
    <template #heading>Server</template>

    This server provides an external <a href="/img/index.jpg" target="_blank" rel="noopener">image</a> store,
    a VueJS-based <RouterLink :to="{ name: 'home' }">client</RouterLink>, and a Golang-based web <a href="/api/ping"
      target="_blank" rel="noopener">server</a>,
    which managed by Nginx.
  </WelcomeItem>

  <WelcomeItem>
    <template #icon>
      <AlpineIcon />
    </template>
    <template #heading>Ecosystem</template>

    The ecosystem is based on:
    <a href="https://www.alpinelinux.org/" target="_blank" rel="noopener">Linux</a>,
    <a href="https://nginx.org/" target="_blank" rel="noopener">Nginx</a>,
    <a href="https://www.docker.com/" target="_blank" rel="noopener">Docker</a>, and
    <a href="https://www.gnu.org/software/bash/" target="_blank" rel="noopener">Bash</a>
    with its configuration and script files.
  </WelcomeItem>
</template>
