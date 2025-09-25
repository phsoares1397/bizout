<template>
    <div class="w-full bg-white py-4 mt-2 px-6">

        <!-- Campo para criar novo marcador -->
        <div v-if="isLoged" class="flex items-center gap-2" :class="{ 'mb-6': !marcadores.length }">
            <input v-model="novoNome" type="text" placeholder="Nome do novo marcador"
                :disabled="marcadorSelecionado != null" class="w-full px-4 py-2 text-sm border border-gray-400 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-300 
           disabled:opacity-50 disabled:cursor-not-allowed" />

            <button @click="criarMarcador" :disabled="marcadorSelecionado != null" class="cursor-pointer px-4 py-2 text-sm rounded-md bg-[#2c89a0] text-white hover:bg-blue-700 transition
           disabled:opacity-50 disabled:cursor-not-allowed">
                Salvar
            </button>
        </div>

        <div v-else class="w-full text-center text-sm text-gray-500 italic">
            ⚠️ Para usar os marcadores, é preciso estar logado.
        </div>


        <!-- Lista de marcadores -->
        <div v-if="marcadores.length" class="flex flex-col gap-3 mt-4">
            <div class="flex flex-wrap gap-2">
                <p class="text-sm tracking-wide text-gray-700 w-full">Marcadores salvos:</p>
                <div v-for="marcador in marcadores" :key="marcador.id"
                    class="flex items-center gap-2 px-3 py-1 rounded-full border border-gray-200 bg-gray-50">
                    <!-- Nome do marcador -->
                    <span class="flex-1 text-sm font-medium text-gray-700">
                        {{ marcador.name }}
                    </span>

                    <!-- Botão de aplicar -->
                    <button v-if="marcadorSelecionado != marcador.id"
                        class="cursor-pointer px-2 py-1 text-xs rounded-md bg-blue-100 hover:bg-blue-200 text-blue-700 font-medium"
                        @click.stop="aplicarMarcador(marcador)">
                        Ativar
                    </button>
                    <button v-else
                        class="cursor-pointer px-2 py-1 text-xs rounded-md bg-amber-100 hover:bg-amber-200 text-amber-700 font-medium"
                        @click.stop="aplicarMarcador(filtrosPadrao)">
                        Desativar
                    </button>

                    <!-- Botão de reset/limpar -->
                    <button
                        class="cursor-pointer px-2 py-1 text-xs rounded-md bg-gray-100 hover:bg-gray-200 text-gray-700 font-medium"
                        @click.stop="limparMarcador(marcador.id)">
                        Resetar
                    </button>

                    <!-- Botão de remover -->
                    <button
                        class="cursor-pointer px-2 py-1 text-xs rounded-md bg-red-100 hover:bg-red-200 text-red-600 font-medium"
                        @click.stop="removerMarcador(marcador.id)">
                        Excluir
                    </button>
                </div>
            </div>
        </div>
        <p v-else-if="isLoged && !marcadores.length" class="text-sm text-gray-500 mb-4">Nenhum marcador salvo ainda.</p>

        <!-- Guia explicativo -->
        <div v-if="!marcadores.length"
            class="mt-4 text-sm text-gray-700 space-y-4 bg-gray-50 rounded-lg p-4 border border-gray-200">
            <h4 class="font-semibold text-gray-800 text-base flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-600" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                O que são marcadores?
            </h4>
            <p class="text-gray-600">
                Os marcadores permitem que você salve um conjunto de filtros de pesquisa (como banca, cargo, ano, etc.)
                para reutilizar facilmente depois.
            </p>

            <ul class="space-y-2">
                <li class="flex items-start gap-2">
                    <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" stroke-width="2"
                        viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M12 4v16m8-8H4" />
                    </svg>
                    <span><strong>Crie</strong> marcadores com filtros específicos (ex: "FCC - Direito - 2023")</span>
                </li>
                <li class="flex items-start gap-2">
                    <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" stroke-width="2"
                        viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M5 13l4 4L19 7" />
                    </svg>
                    <span><strong>Aplique</strong> filtros instantaneamente ao clicar em um marcador salvo</span>
                </li>
                <li class="flex items-start gap-2">
                    <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" stroke-width="2"
                        viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M3 10h11M9 21V3m0 0L4 8m5-5l5 5" />
                    </svg>
                    <span>
                        <strong>Acompanhe seu desempenho</strong> apenas nas questões resolvidas de um marcador.
                        Ideal para focar em temas ou bancas específicas.
                    </span>
                </li>
            </ul>
        </div>

        <div class="flex flex-wrap gap-2 mt-4">
            <template v-for="(valores, chave) in filtrosMkrs" :key="`filtro-${chave}`">
                <!-- Ignora disciplina e arrays vazios -->
                <template v-if="valores.length && chave !== 'disciplina'">
                    <template v-for="(valor, index) in valores.slice(0, 3)" :key="`${chave}-${index}`">
                        <span
                            class="inline-flex items-center bg-blue-100 text-blue-800 text-xs px-3 py-1 rounded-full cursor-default select-none"
                            :title="`${formatarChave(chave)}: ${displayValue(chave, valor)}`">
                            {{ formatarChave(chave) }}: {{ displayValue(chave, valor) }}
                        </span>
                    </template>

                    <!-- Mostrar +X caso tenha mais de 3 -->
                    <button v-if="valores.length > 3" @click="abrirModal(chave)"
                        class="inline-flex items-center bg-gray-200 text-gray-600 text-xs px-2 py-1 rounded-full hover:bg-gray-300">
                        +{{ valores.length - 3 }}
                    </button>
                </template>
            </template>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, reactive, inject } from 'vue';

const isLoged = inject('isLoged')

const props = defineProps({
    banco: { type: String, required: true } // conc | vest | enem
});

const filtrosMkrs = reactive({
    disciplina: [], assuntos: [], orgao: [], cargo: [],
    banca: [], ano: [], nivel: [], tipo: [], dificuldade: []
})

function formatarChave(chave) {
    return {
        disciplina: 'Disciplina',
        assuntos: 'Assunto',
        orgao: 'Órgão',
        cargo: 'Cargo',
        banca: 'Banca',
        ano: 'Ano',
        nivel: 'Nível',
        tipo: 'Tipo',
        dificuldade: 'Dificuldade'
    }[chave] || chave
}

function displayValue(chave, valor) {
    if (chave === 'assuntos') {
        return valor[0]; // array de arrays, só mostra o primeiro elemento
    } else if (valor && valor.nome) {
        return Array.isArray(valor.nome) ? valor.nome.join(' - ') : valor.nome;
    } else {
        return valor;
    }
}

const marcadores = ref([]);
const novoNome = ref('');
const marcadorSelecionado = ref(null);

let mkrsNC = []
let mkrsNE = []

const filtrosPadrao = reactive({
    disciplina: [], assuntos: [], orgao: [], cargo: [],
    banca: [], ano: [], nivel: [], tipo: [], dificuldade: [], name: null, id: null
})

// Mapa de codificação
const keyMap = {
    disciplina: "a",
    assuntos: "b",
    orgao: "c",
    cargo: "d",
    banca: "e",
    ano: "f",
    nivel: "g",
    tipo: "h",
    dificuldade: "i",
    name: "j"
};

// Mapa inverso para decodificar
const invKeyMap = Object.fromEntries(
    Object.entries(keyMap).map(([k, v]) => [v, k])
);

// Função para codificar
function encodeObject(obj) {
    const newObj = {};
    for (let key in obj) {
        if (obj.hasOwnProperty(key)) {
            newObj[keyMap[key]] = obj[key];
        }
    }
    return newObj;
}

// Função para decodificar
function decodeObject(obj) {
    const newObj = {};
    for (let key in obj) {
        if (obj.hasOwnProperty(key)) {
            newObj[invKeyMap[key]] = obj[key];
        }
    }
    return newObj;
}

//const decoded = decodeObject(encoded);
//console.log("Decodificado:", decoded);

// 🔌 Função para requisições
async function apiPost(url, data) {
    try {
        const res = await fetch(`https://tools.bizout.com.br${url}`, {
            method: 'POST',
            credentials: 'include', // envia cookies (session_id)
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: new URLSearchParams(data)
        });
        return await res.json();
    } catch (err) {
        console.error('Erro API:', err);
        return { 0: 500, 1: 'Erro de rede' };
    }
}

// 🔹 Carrega lista de marcadores ao montar
onMounted(async () => {
    const resp = await apiPost('/mkrs-list', { banco: props.banco });
    if (resp[0] === 200) {
        const responseJSON = JSON.parse(resp[1])
        let obj = []
        Object.keys(responseJSON).forEach((key) => {
            if (key == 0) return
            let temp = decodeObject(responseJSON[key])
            temp.id = key
            obj.push(temp)
        })
        marcadores.value = obj; // backend já retorna array de marcadores
    }
    window.marcadorSelecionado = marcadorSelecionado
    window.mkrsNC = mkrsNC
    window.mkrsNE = mkrsNE
});

// 🔹 Criar novo marcador
async function criarMarcador() {
    const nome = novoNome.value.trim();
    if (!nome) return;

    let sendMkrs = window.filtrosNormalizados
    sendMkrs.name = nome
    const encoded = encodeObject(sendMkrs);
    const resp = await apiPost('/mkrs-create', {
        banco: props.banco,
        data: JSON.stringify(encoded)
    });

    if (resp[0] === 200) {
        marcadores.value.push(resp[1]);
        novoNome.value = '';
    } else {
        alert(resp[1]);
    }
}

async function lerDadosMarcador(markerId, banco) {
    const data = {
        marker_id: markerId,
        banco: banco // ex: 'conc', 'vest' ou 'enem'
    }

    try {
        const res = await fetch('https://tools.bizout.com.br/mkrs-read-data', {
            method: 'POST',
            credentials: 'include', // envia cookies (session_id)
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: new URLSearchParams(data)
        })

        return await res.json()
    } catch (err) {
        console.error('Erro API:', err)
        return { 0: 500, 1: 'Erro de rede' }
    }
}

// 🔹 Aplicar marcador
async function aplicarMarcador(marcador) {
    marcadorSelecionado.value = marcador.id;

    // Cria uma cópia do marcador sem id e name
    let marcadorTemp = { ...marcador };
    delete marcadorTemp.id;
    delete marcadorTemp.name;

    // Atualiza cada chave de filtrosMkrs
    for (let key in filtrosMkrs) {
        if (marcadorTemp[key]) {
            filtrosMkrs[key] = Array.isArray(marcadorTemp[key]) ? marcadorTemp[key] : [];
        } else {
            filtrosMkrs[key] = [];
        }
    }

    if (marcador.id != null) {
        const resposta = await lerDadosMarcador(marcador.id, 'conc')
        if (resposta[0] === 200) {
            window.mkrsNC = resposta[1].nc
            window.mkrsNE = resposta[1].ne
        } else {
            console.warn('Erro:', resposta[1])
        }
        window.filterLocked.value = true
    } else {
        window.filterLocked.value = false
        window.mkrsNC = null
        window.mkrsNE = null
    }

    window.filtros = marcadorTemp;
    window.pagina = 1
    window.loadPage(1, 0, marcadorTemp)
}

// 🔹 Remover marcador
async function removerMarcador(id) {
    if (!confirm('Remover marcador?')) return;

    const resp = await apiPost('/mkrs-delete', {
        banco: props.banco,
        marker_id: id
    });

    if (resp[0] === 200) {
        marcadores.value = marcadores.value.filter(m => m.id !== id);
    } else {
        alert(resp[1]);
    }
}

// 🔹 Limpar todos os marcadores
async function limparMarcadores() {
    if (!confirm('Tem certeza que deseja limpar todos os marcadores?')) return;

    const resp = await apiPost('/mkrs-clear', { banco: props.banco });

    if (resp[0] === 200) {
        marcadores.value = [];
        marcadorSelecionado.value = null;
    } else {
        alert(resp[1]);
    }
}
</script>
