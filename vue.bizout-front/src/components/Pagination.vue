<template>
    <nav v-if="totalPaginas > 1" class="flex flex-wrap items-center justify-center gap-2 select-none" role="navigation"
        aria-label="Paginação">
        <!-- ← anterior -->
        <button :disabled="pagina === 1" @click="irPara(pagina - 1)" class="botao-seta" aria-label="Anterior">
            <svg viewBox="0 0 24 24" class="w-4 h-4">
                <path d="M15 6l-6 6 6 6" fill="none" stroke="currentColor" stroke-linecap="round"
                    stroke-linejoin="round" stroke-width="2" />
            </svg>
        </button>

        <!-- páginas visíveis -->
        <button v-for="p in paginasVisiveis" :key="p" @click="irPara(p)" :class="[
            'px-3 py-1 rounded-md text-sm cursor-pointer',
            p === pagina
                ? 'bg-blue-600 text-white'
                : 'text-gray-700 hover:bg-gray-100'
        ]">
            {{ p }}
        </button>

        <!-- campo digitar página -->
        <input v-model="input" @keyup.enter="confirmarInput" @blur="confirmarInput" class="w-14 px-2 py-1 text-sm text-center border border-gray-300 rounded-md
               focus:outline-none focus:ring-2 focus:ring-blue-300" :placeholder="pagina"
            aria-label="Ir para página" />

        <!-- próximo → -->
        <button :disabled="pagina === totalPaginas" @click="irPara(pagina + 1)" class="botao-seta" aria-label="Próxima">
            <svg viewBox="0 0 24 24" class="w-4 h-4 rotate-180">
                <path d="M15 6l-6 6 6 6" fill="none" stroke="currentColor" stroke-linecap="round"
                    stroke-linejoin="round" stroke-width="2" />
            </svg>
        </button>
    </nav>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

/* ---------- props ---------- */
const props = defineProps({
    total: { type: Number, required: true },
    porPagina: { type: Number, default: 10 },
    page: { type: Number, default: 1 }
})
const emit = defineEmits(['update:page'])

/* ---------- derivado ---------- */
const totalPaginas = computed(() =>
    Math.max(1, Math.ceil(props.total / props.porPagina))
)

/* ---------- página vinculada ---------- */
const pagina = computed({
    get: () => Math.min(Math.max(1, props.page), totalPaginas.value),
    set: p => emit('update:page', p)
})

/* ---------- janela de botões ---------- */
const paginasVisiveis = computed(() => {
    const range = 5
    const half = Math.floor(range / 2)
    let start = Math.max(1, pagina.value - half)
    let end = start + range - 1
    if (end > totalPaginas.value) {
        end = totalPaginas.value
        start = Math.max(1, end - range + 1)
    }
    return Array.from({ length: end - start + 1 }, (_, i) => start + i)
})

/* ---------- navegação ---------- */
function irPara(p) {
    const alvo = Number(p)
    if (!Number.isNaN(alvo) && alvo >= 1 && alvo <= totalPaginas.value) {
        pagina.value = alvo
    }
}

/* ---------- input direto ---------- */
const input = ref('')

function confirmarInput() {
    if (input.value.trim() !== '') {
        irPara(+input.value)
        input.value = ''
    }
}

/* mantém placeholder sempre com página atual */
watch(pagina, () => { input.value = '' })
</script>

<style scoped>
@reference "../style.css";

.botao-seta {
    @apply cursor-pointer p-2 text-gray-600 enabled:hover:bg-gray-100 enabled:hover:text-gray-800 rounded-md disabled:opacity-30 disabled:cursor-default;
}
</style>