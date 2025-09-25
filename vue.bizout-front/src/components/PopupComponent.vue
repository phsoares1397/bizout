<template>
    <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 transition-opacity">
        <div class="relative p-5 rounded-lg shadow-lg max-w-sm w-full bg-white text-gray-800 transition-all" :class="{
            'border-l-4 border-green-400': alertData.type === 'success',
            'border-l-4 border-red-400': alertData.type === 'error',
            'border-l-4 border-yellow-400': alertData.type === 'warning'
        }">
            <!-- Título e botão de fechar -->
            <div class="flex justify-between items-start mb-3 pb-2 border-b border-gray-200">
                <h3 class="text-sm font-semibold text-gray-700">
                    {{ alertData.title }}
                </h3>
                <div class="relative w-6 h-6">
                    <!-- Círculo animado -->
                    <svg v-if="timer > 0" class="absolute top-0 left-0 w-full h-full rotate-[-90deg]"
                        viewBox="0 0 36 36">
                        <circle class="text-gray-300" stroke="currentColor" stroke-width="3" fill="transparent" r="16"
                            cx="18" cy="18" />
                        <circle class="text-[#f2a81d]" :stroke-dasharray="circumference" :stroke-dashoffset="offset"
                            stroke-linecap="round" stroke="currentColor" stroke-width="3" fill="transparent" r="16"
                            cx="18" cy="18" />
                    </svg>

                    <!-- Botão de fechar -->
                    <button @click="fechar"
                        class="absolute inset-0 flex items-center justify-center text-gray-500 hover:text-gray-700 cursor-pointer">
                        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
            </div>

            <!-- Mensagem -->
            <p class="text-base font-medium">{{ alertData.text }}</p>
        </div>
    </div>
</template>

<script setup lang="ts">

interface AlertData {
    title: string
    text: string
    type: 'success' | 'error' | 'warning'
}

const props = defineProps<{
    visible: boolean
    alertData: AlertData
    timer: number
    circumference: number
    offset: number
}>()

let x = props 
x = x

const emit = defineEmits<{
    (e: 'close'): void
}>()

function fechar() {
    emit('close')
}
</script>