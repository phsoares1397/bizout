<template>
  <!-- Alertas -->
  <div class="fixed bottom-4 right-4 z-50 space-y-2 w-72">
    <div v-for="(alert, i) in alerts" :key="alert.id" :class="[
      'px-4 py-3 rounded-lg shadow-md text-sm transition-all duration-300',
      {
        'bg-green-100 text-green-800': alert.type === 'success',
        'bg-red-100 text-red-800': alert.type === 'error',
        'bg-yellow-100 text-yellow-800': alert.type === 'warning'
      }
    ]">
      {{ alert.message }}
    </div>
  </div>

  <!-- Popup -->
  <div v-if="visible" class="fixed inset-0 z-1000 flex items-center justify-center bg-black/40 transition-opacity">
    <div class="relative p-5 rounded-lg shadow-lg max-w-sm w-full bg-white text-gray-800 transition-all" :class="{
      'border-l-4 border-green-400': popupMessage.type === 'success',
      'border-l-4 border-red-400': popupMessage.type === 'error',
      'border-l-4 border-yellow-400': popupMessage.type === 'warning'
    }">
      <div class="flex justify-between items-start mb-3 pb-2 border-b border-gray-200">
        <h3 class="text-sm font-semibold text-gray-700">
          {{ popupMessage.title }}
        </h3>
        <button @click="fecharPopup" class="cursor-pointer text-gray-500 hover:text-gray-700">✕</button>
      </div>
      <p class="text-base font-medium">{{ popupMessage.text }}</p>
    </div>
  </div>

  <!-- Header -->
  <header v-if="headerVisible"
    class="fixed inset-x-0 top-0 z-50 h-[58px] w-full bg-[#414141] shadow flex items-center justify-between px-4">

    <!-- Botão Voltar -->
    <button @click="voltar" class="cursor-pointer flex items-center text-white hover:text-gray-300 transition">
      <!-- Ícone seta -->
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"
        stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
      </svg>
    </button>

    <!-- Logo centralizada -->
    <div class="absolute left-1/2 transform -translate-x-1/2">
      <img class="h-[30px]" src="./assets/bizout_logo.svg" />
    </div>
  </header>


  <!-- Onde mudam as telas -->
  <main v-if="headerVisible" class="pt-[58px] bg-gray-100 min-h-full">
    <Transition name="slide" mode="out-in">
      <router-view :key="$route.fullPath" />
    </Transition>
  </main>
  <main v-else class="bg-gray-100 min-h-full">
    <Transition name="slide" mode="out-in">
      <router-view :key="$route.fullPath" />
    </Transition>
  </main>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue"

const alerts = reactive([])
let nextId = 0
const visible = ref(false)
const headerVisible = ref(true)
const popupMessage = ref({ title: "", text: "", type: "success" })
let timer

function headerVisibleFunc(op) {
  headerVisible.value = op
}

async function voltar() {
  const user = await checkSession()
  if (user) {
    window.location = "https://bizout.com.br/questoes"
  } else {
    window.history.back()
  }
}

async function checkSession() {
  try {
    const res = await fetch("https://auth.bizout.com.br/me", {
      credentials: "include" // envia automaticamente session_id
    })
    if (!res.ok) return null

    const user = await res.json()
    return user
  } catch {
    return null
  }
}

function mostrarPopup(title, text, type = "success", time = 0) {
  popupMessage.value = { title, text, type }
  visible.value = true
  if (time > 0) {
    timer = setTimeout(() => {
      visible.value = false
    }, time)
  }
}

function fecharPopup() {
  visible.value = false
  clearTimeout(timer)
}

function showAlert(message, type = "success", duration = 3000) {
  const id = nextId++
  alerts.push({ id, message, type })
  setTimeout(() => {
    const index = alerts.findIndex(alert => alert.id === id)
    if (index !== -1) alerts.splice(index, 1)
  }, duration)
}

onMounted(() => {
  window.mostrarPopup = mostrarPopup
  window.showAlert = showAlert
  window.headerVisibleFunc = headerVisibleFunc
  window.checkSession = checkSession
})
</script>

<style>
.slide-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.slide-enter-to {
  transform: translateX(0);
  opacity: 1;
}

.slide-leave-from {
  transform: translateX(0);
  opacity: 1;
}

.slide-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}
</style>
