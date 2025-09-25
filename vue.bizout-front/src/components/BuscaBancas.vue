<template>
    <div class="relative">
        <!-- Campo de busca -->
        <input v-model="termo" type="text" placeholder="Buscar..."
            class="w-full px-4 py-2 pl-10 text-sm border border-gray-200 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-300" />
        <svg class="w-4 h-4 text-gray-400 absolute left-3 top-0 bottom-0 my-auto" fill="none" stroke="currentColor"
            stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round"
                d="M21 21l-4.35-4.35m0 0A7 7 0 104.5 4.5a7 7 0 0012.15 12.15z" />
        </svg>
    </div>

    <!-- Mensagem quando o campo está vazio -->
    <span v-if="termo.trim() === ''" class="text-xs text-gray-400 block mt-1 pl-1 mb-4">
        Digite para buscar...
    </span>

    <!-- Quando o campo tem texto -->
    <div v-else class="mb-2 flex-grow overflow-y-auto space-y-2 p-2 bg-gray-50 rounded-md max-h-[300px] mt-4">
        <template v-if="resultadoCarregado && resultados.length">
            <!-- Resultados -->
            <label v-for="banca in resultados" :key="banca.id"
                class="flex items-center gap-2 text-sm text-gray-700 hover:bg-gray-100 px-2 py-1 rounded-md cursor-pointer">
                <input type="checkbox" :value="banca" v-model="bancasSelecionados" class="accent-blue-600 w-4 h-4" />
                <span>{{ banca.nome[0] }} - {{ banca.nome[1] }}</span>
            </label>
        </template>

        <p v-else-if="resultadoCarregado" class="text-xs text-gray-400 px-2">
            Nenhum resultado encontrado.
        </p>
        <p v-else class="text-xs text-gray-400 px-2">Buscando...</p>
    </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'

const termo = ref('')

// Props e emits para v-model
const props = defineProps({
    modelValue: {
        type: Array,
        default: () => []
    }
})
const emit = defineEmits(['update:modelValue'])

// Lista de órgãos selecionados (local)
const bancasSelecionados = ref([...props.modelValue])

// Sincroniza ida e volta com o pai
watch(bancasSelecionados, (val) => {
    emit('update:modelValue', val)
})

const bancasPopulares = [
    { id: 1, nome: 'INSS' },
    { id: 2, nome: 'Receita Federal' },
    { id: 3, nome: 'Polícia Federal' },
    { id: 4, nome: 'Banco Central' },
    { id: 5, nome: 'CGU' }
]

const resultados = ref([])
const resultadoCarregado = ref(false)
let controller = null

watch(termo, async (valor) => {
    resultadoCarregado.value = false
    if (valor.trim().length < 3) {
        resultados.value = []
        return
    }

    if (controller) controller.abort()
    controller = new AbortController()

    try {
        const res = await fetch(`https://elastic.bizout.com.br/bnc?srch=${encodeURIComponent(valor)}`, {
            signal: controller.signal
        })

        const data = await res.json()
        resultados.value = Object.entries(data).map(([id, val]) => {
            const nome = JSON.parse(val)
            return { id, nome }
        })
    } catch (e) {
        if (e.name !== 'AbortError') {
            console.error('Erro ao buscar bancas:', e)
        }
        resultados.value = []
    } finally {
        resultadoCarregado.value = true
    }
})

function selecionarbanca(banca) {
    if (!bancasSelecionados.value.some((o) => o.id === banca.id)) {
        bancasSelecionados.value.push(banca)
    }
}
</script>
