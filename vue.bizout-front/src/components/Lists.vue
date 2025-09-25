<template>
    <div class="w-full bg-white py-4 mt-2 px-6">
        <!-- Campo para criar nova lista -->
        <div class="flex items-center gap-2" :class="{ 'mb-6': !listas.length }">
            <input v-model="novoNome" type="text" placeholder="Nome da nova lista"
                class="w-full px-4 py-2 text-sm border border-gray-400 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-300" />
            <button @click="criarLista"
                class="cursor-pointer px-4 py-2 text-sm rounded-md bg-[#2c89a0] text-white hover:bg-blue-700 transition">
                Salvar
            </button>
        </div>

        <!-- Lista de listas salvas -->
        <div v-if="listas.length" class="flex flex-wrap gap-2 mt-4">
            <p class="text-sm tracking-wide flex items-center justify-center text-gray-700">Listas salvas:</p>
            <button v-for="(lista, index) in listas" :key="index"
                class="cursor-pointer px-3 py-1 text-sm rounded-full transition-colors flex items-center gap-2" :class="{
                    'bg-blue-100 text-blue-700': listaSelecionada === lista.nome,
                    'bg-gray-100 text-gray-700 hover:bg-blue-50': listaSelecionada !== lista.nome
                }" @click="abrirLista(lista)">
                {{ lista.nome }}
                <span class="text-red-600 text-xs font-bold hover:text-red-800 ml-1 cursor-pointer"
                    @click.stop="removerLista(index)">
                    ×
                </span>
            </button>
        </div>
        <p v-else class="text-sm text-gray-500 mb-4">Nenhuma lista salva ainda.</p>

        <!-- Guia explicativo -->
        <div v-if="!listas.length"
            class="mt-4 text-sm text-gray-700 space-y-4 bg-gray-50 rounded-lg p-4 border border-gray-200">
            <h4 class="font-semibold text-gray-800 text-base flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-600" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                O que são listas?
            </h4>
            <p class="text-gray-600">
                As listas permitem que você salve questões específicas para revisar ou praticar mais tarde.
            </p>

            <ul class="space-y-2">
                <li class="flex items-start gap-2">
                    <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" stroke-width="2"
                        viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M12 4v16m8-8H4" />
                    </svg>
                    <span><strong>Crie</strong> listas com questões específicas que deseja focar</span>
                </li>
                <li class="flex items-start gap-2">
                    <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" stroke-width="2"
                        viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M5 13l4 4L19 7" />
                    </svg>
                    <span><strong>Acesse</strong> rapidamente as questões salvas em uma lista</span>
                </li>
                <li class="flex items-start gap-2">
                    <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" stroke-width="2"
                        viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M3 10h11M9 21V3m0 0L4 8m5-5l5 5" />
                    </svg>
                    <span><strong>Reforce</strong> conteúdos com revisão seletiva das questões de cada lista</span>
                </li>
            </ul>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';

const props = defineProps({
    questoesSelecionadas: Array // array de IDs ou objetos de questões
});
const emit = defineEmits(['abrir']);

const LOCAL_STORAGE_KEY = 'minhasListas';
const listas = ref([]);
const novoNome = ref('');
const listaSelecionada = ref(null);

// Carrega listas do localStorage
onMounted(() => {
    const salvas = localStorage.getItem(LOCAL_STORAGE_KEY);
    if (salvas) {
        try {
            listas.value = JSON.parse(salvas);
        } catch (e) {
            console.error('Erro ao carregar listas:', e);
        }
    }
});

// Atualiza localStorage sempre que listas mudam
watch(listas, (novas) => {
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(novas));
}, { deep: true });

function criarLista() {
    const nome = novoNome.value.trim();
    if (!nome) return;

    const existe = listas.value.some(l => l.nome.toLowerCase() === nome.toLowerCase());
    if (existe) {
        alert('Já existe uma lista com esse nome.');
        return;
    }

    listas.value.push({
        nome,
        questoes: [...props.questoesSelecionadas] // salva uma cópia das questões selecionadas
    });
    novoNome.value = '';
}

function abrirLista(lista) {
    listaSelecionada.value = lista.nome;
    emit('abrir', lista.questoes);
}

function removerLista(index) {
    if (confirm('Remover esta lista?')) {
        listas.value.splice(index, 1);
    }
}
</script>
