<template>
    <div class="w-full max-w-2xl mx-auto bg-white rounded-lg flex flex-col flex-1 min-h-0">
        <!-- Campo de busca -->
        <div class="relative mb-2 flex-shrink-0">
            <input v-model="busca" type="text" placeholder="Buscar disciplina..."
                class="w-full px-4 py-2 pl-10 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-400 outline-none" />
            <svg class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" fill="none"
                stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round"
                    d="M21 21l-4.35-4.35m0 0A7 7 0 104.5 4.5a7 7 0 0012.15 12.15z" />
            </svg>
        </div>

        <!-- Lista de nós -->
        <div class="flex-grow overflow-y-auto min-h-0 bg-gray-50 rounded-md py-2 mb-2 pr-1" id="assuntos_tree">
            <div v-if="loading && !mostrar" class="flex justify-center items-center mt-3 space-x-3 text-gray-600">
                <svg class="w-6 h-6 text-blue-500 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none"
                    viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
                </svg>
                <span>Procurando assuntos...</span>
            </div>

            <TreeItem v-if="mostrar" v-for="item in arvoreFiltrada" :key="item.id" :item="item" :busca="busca"
                v-model:selecionados="selecionados" />
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import TreeItem from './TreeItem.vue'
import { construirArvore } from '../utils/construirArvore'

const mostrar = ref(false)
const loading = ref(true)
const busca = ref('')
const selecionados = ref(JSON.parse(localStorage.getItem('selecionados')) || [])

// Sempre que selecionados mudar, salva no localStorage
watch(selecionados, (novo) => {
    localStorage.setItem('selecionados', JSON.stringify(novo))
    let assuntos_atual = []
    normalizados.value.forEach(item => assuntos_atual.push([encontrarNomePorId(arvoreFiltrada.value, item), item]))
    window.filtros.assuntos = assuntos_atual
}, { deep: true })

let json = []
let arvore = []

const disciplinasIds = JSON.parse(localStorage.getItem("disciplinasIds")) || []

function encontrarNomePorId(estrutura, idBuscado) {
    for (const item of estrutura) {
        if (item.id === idBuscado) {
            return item.nome;
        }
        if (item.filhos && item.filhos.length > 0) {
            const resultado = encontrarNomePorId(item.filhos, idBuscado);
            if (resultado) return resultado;
        }
    }
    return null; // se não encontrar
}

const fetchPromises = disciplinasIds.map(id =>
    fetch('../../questoes/src/data/asst_conc/' + id + '.json').then(res => res.json())
)

Promise.all(fetchPromises)
    .then(dados => {
        arvore = construirArvore(dados)
        loading.value = false
        mostrar.value = true
    })
    .catch(err => console.error('Erro ao carregar arquivos:', err))

// 🔍 Filtragem com hierarquia
function filtrarArvore(nos, termo) {
    return nos.map(no => {
        const filhosFiltrados = filtrarArvore(no.filhos || [], termo)
        const corresponde = no.nome.toLowerCase().includes(termo.toLowerCase())
        if (corresponde || filhosFiltrados.length > 0) {
            return {
                ...no,
                filhos: filhosFiltrados
            }
        }
        return null
    }).filter(Boolean)
}

const arvoreFiltrada = computed(() => {
    if (!busca.value.trim()) return arvore
    return filtrarArvore(arvore, busca.value)
})

// 🔁 Coleta todos os IDs recursivamente (aceita nó ou array de nós)
function coletarTodosIds(entrada) {
    if (Array.isArray(entrada)) {
        return entrada.flatMap(coletarTodosIds)
    }
    const no = entrada
    let ids = [no.id]
    if (no.filhos && no.filhos.length) {
        ids = ids.concat(no.filhos.flatMap(coletarTodosIds))
    }
    return ids
}

// ✅ Seleção normalizada
function normalizarSelecao(no) {
    if (!no.filhos || no.filhos.length === 0) {
        return selecionados.value.includes(no.id) ? [no.id] : []
    }

    const todosFilhosSelecionados = coletarTodosIds(no.filhos).every(id =>
        selecionados.value.includes(id)
    )

    if (todosFilhosSelecionados) {
        return [no.id]
    }

    const filhosNormalizados = no.filhos.flatMap(normalizarSelecao)
    const ids = [...filhosNormalizados]
    if (selecionados.value.includes(no.id)) {
        ids.push(no.id)
    }
    return ids
}

const normalizados = computed(() => {
    if (!arvore || !Array.isArray(arvore)) return []
    return arvore.flatMap(normalizarSelecao)
})

// Debug manual
window.selecionados = selecionados
window.normalizados = normalizados
window.arvoreFiltrada = arvoreFiltrada
</script>
