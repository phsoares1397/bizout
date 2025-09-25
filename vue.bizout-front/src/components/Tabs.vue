<template>
    <nav ref="tabsContainer" class="tabs-container flex overflow-x-auto no-scrollbar pl-3 mt-[68px]">
        <!-- Aba tipo de questão -->
        <!-- <div
            class="flex items-center px-4 text-sm font-medium text-gray-600 border-b-[3px] border-b-gray-600 flex-shrink-0">
            <span class="mr-1">Módulo:</span>
            <select v-model="tipoQuestao"
                class="cursor-pointer text-sm text-gray-700 bg-transparent border-none focus:outline-none">
                <option value="Concursos">Concursos</option>
                <option value="Vestibular">Vestibular</option>
                <option value="ENEM">ENEM</option>
            </select>
        </div> -->

        <!-- Outras abas -->
        <button v-for="tab in tabs" :key="tab.name" @click="handleClick(tab.name, $event)" type="button" class="cursor-pointer flex items-center px-6 text-center py-3 text-sm font-medium
               text-gray-600 hover:text-gray-800 border-b-[3px] border-b-[#ddd] flex-shrink-0 transition-colors"
            :class="{
                'text-gray-800 border-b-[#f2a81d]': ativo === tab.name,
                'border-b-[#ddd]': ativo !== tab.name,
            }">
            <span class="mr-2" v-html="(filterLocked && tab.name === 'Filtros') ? tab.iconLocked : tab.icon"></span>
            <span>{{ tab.name }}</span>
        </button>

        <span
            class="flex-1 px-6 text-center py-3 text-sm font-medium text-gray-600 hover:text-gray-800 border-b-[3px] border-b-[#ddd]"></span>
    </nav>

    <!-- Conteúdos -->
    <section v-show="ativo === 'Filtros'">
        <FilterBar />
    </section>
    <section v-if="mkrsOn" v-show="ativo === 'Marcadores' && !loadingGraphic">
        <Markers :banco="'conc'" @aplicar="aplicarMarcador" />
    </section>
    <section v-show="ativo === 'Listas'">
        <Lists :filtrosAtuais="filtros" @aplicar="aplicarListas" />
    </section>
    <section v-if="ativo === 'Desempenho' && !loadingGraphic" class="flex justify-center">
        <Desempenho :historico="historicoTeste" />
    </section>
    <section v-if="loadingGraphic" class="flex justify-center items-center py-20">
        <div class="flex flex-col items-center">
            <!-- Gráfico animado colorido -->
            <div class="flex space-x-1 mb-2">
                <span class="w-2 h-6 animate-bounce" style="background-color:#f2a81d; animation-delay:0s;"></span>
                <span class="w-2 h-8 animate-bounce" style="background-color:#2c89a0; animation-delay:0.1s;"></span>
                <span class="w-2 h-5 animate-bounce" style="background-color:#f26c4f; animation-delay:0.2s;"></span>
                <span class="w-2 h-7 animate-bounce" style="background-color:#6a5bdc; animation-delay:0.3s;"></span>
                <span class="w-2 h-6 animate-bounce" style="background-color:#4caf50; animation-delay:0.4s;"></span>
            </div>
            <span class="text-gray-700 font-medium">Buscando dados...</span>
        </div>
    </section>

</template>

<script setup>
import { reactive, ref, onMounted, provide, onBeforeMount } from 'vue'
import FilterBar from './FilterBar.vue'
import Markers from './Markers.vue'
import Lists from './Lists.vue'
import Desempenho from './Graphic.vue'

let historicoTeste = []

const filterLocked = ref(false)
const loadingGraphic = ref(false)
const mkrsOn = ref(false)

const isLoged = ref(false) // valor inicial
provide('isLoged', isLoged)

const tabs = [
    {
        name: 'Filtros',
        icon: `
        <svg xmlns="http://www.w3.org/2000/svg" fill="none"
             viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
             class="w-5 h-5">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M3 4.5h18M6 10.5h12M9 16.5h6" />
        </svg>
      `,
        iconLocked: `
        <svg xmlns="http://www.w3.org/2000/svg" fill="none"
            viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
            class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round"
                d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75M6.75 10.5h10.5a1.5 1.5 0 011.5 1.5v7.5a1.5 1.5 0 01-1.5 1.5H6.75a1.5 1.5 0 01-1.5-1.5v-7.5a1.5 1.5 0 011.5-1.5z" />
        </svg>
        `,
    },
    {
        name: 'Marcadores',
        icon: `
        <svg xmlns="http://www.w3.org/2000/svg" fill="none"
             viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
             class="w-5 h-5">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M6 4.5h12a1.5 1.5 0 011.5 1.5v14l-7.5-3-7.5 3V6a1.5 1.5 0 011.5-1.5z" />
        </svg>
      `,
    },
    {
        name: 'Desempenho',
        icon: `
        <svg xmlns="http://www.w3.org/2000/svg" fill="none"
             viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
             class="w-5 h-5">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M3 3v18h18M9 17V9m4 8v-5m4 5V5" />
        </svg>
      `,
    },
]

// implementação futura
// {
//     name: 'Listas',
//         icon: `
//         <svg xmlns="http://www.w3.org/2000/svg" fill="none"
//              viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
//              class="w-5 h-5">
//           <path stroke-linecap="round" stroke-linejoin="round"
//                 d="M11.48 3.499c.27-.823 1.47-.823 1.74 0l1.286 3.935a1
//                  1 0 00.95.69h4.138c.864 0 1.223 1.104.526 1.618l-3.35
//                  2.435a1 1 0 00-.364 1.118l1.286 3.935c.27.823-.755
//                  1.506-1.54 1.004l-3.35-2.435a1 1 0 00-1.175 0l-3.35
//                  2.435c-.785.502-1.81-.181-1.54-1.004l1.286-3.935a1
//                  1 0 00-.364-1.118l-3.35-2.435c-.697-.514-.338-1.618.526-1.618h4.138a1
//                  1 0 00.95-.69l1.286-3.935z" />
//         </svg>
//       `,
//     },

const ativo = ref('Filtros')
const tipoQuestao = ref('Concursos') // valor padrão

async function fetchGraphicData(mkrId) {
    try {
        const res = await fetch(`https://tools.bizout.com.br/mkrs-read-graphic?mkrId=${mkrId}`, {
            method: 'GET',
            credentials: 'include', // envia cookies (session_id)
        });

        if (!res.ok) {
            console.warn('Erro HTTP:', res.status, res.statusText);
            return {};
        }

        const data = await res.json();
        return data; 
    } catch (err) {
        console.error('Erro API:', err);
        return {};
    }
}

function formatGraphicData(rawData) {
    // rawData é o objeto retornado pelo backend, ex: { "2025-09-09": { nc: 3, ne: 1 }, ... }
    const formatted = [];

    for (const [date, counts] of Object.entries(rawData)) {
        formatted.push({
            data: date,
            resolvidas: counts.nc + counts.ne,
            acertos: counts.nc,
            erros: counts.ne
        });
    }

    // opcional: ordenar por data crescente
    formatted.sort((a, b) => new Date(a.data) - new Date(b.data));

    return formatted;
}

const svgMarcador = `
<span class="inline-flex items-center space-x-1">
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
    <path stroke-linecap="round" stroke-linejoin="round" d="M6 4.5h12a1.5 1.5 0 011.5 1.5v14l-7.5-3-7.5 3V6a1.5 1.5 0 011.5-1.5z"></path>
  </svg>
  <span class="mt-[1px]">Marcador ativado</span>
</span>
`;

const svgMarcadorAtivar = `
<span class="inline-flex items-center space-x-1">
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
    <path stroke-linecap="round" stroke-linejoin="round" d="M6 4.5h12a1.5 1.5 0 011.5 1.5v14l-7.5-3-7.5 3V6a1.5 1.5 0 011.5-1.5z"></path>
  </svg>
  <span class="mt-[1px]">Marcador desativado</span>
</span>
`;

async function handleClick(tabName, event) {
    const btn = event.currentTarget

    if (tabName === "Filtros" && window.filterLocked.value === true) {
        window.mostrarPopup(
            `${svgMarcador}`,
            'Desative para poder acessar os filtros.',
            'warning',
            0
        );
        return;
    } else if (tabName == "Desempenho" && window.marcadorSelecionado.value == null) {
        window.mostrarPopup(
            `${svgMarcadorAtivar}`,
            'Ative algum marcador para acompanhar o desempenho.',
            'warning',
            0
        )
        return;
    } else if (tabName == "Desempenho" && window.marcadorSelecionado.value != null) {
        loadingGraphic.value = true
        const graphicData = await fetchGraphicData(window.marcadorSelecionado.value);
        const graphicDataPrepared = formatGraphicData(graphicData);
        historicoTeste = graphicDataPrepared
        setTimeout(() => {
            loadingGraphic.value = false
        }, 2500)
    }

    mkrsOn.value = tabName === "Marcadores" || window.marcadorSelecionado != null

    ativo.value = tabName

    btn.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'start' })
}

async function checkSession() {
    try {
        const res = await fetch("https://auth.bizout.com.br/me", {
            credentials: "include" // envia automaticamente session_id
        })
        if (!res.ok) return false

        const user = await res.json()
        return user
    } catch {
        return false
    }
}

async function validateSession() {
    const user = await checkSession()
    isLoged.value = !!user
}

onBeforeMount(async () => {
    validateSession()
})

onMounted(() => {
    window.filterLocked = filterLocked
    window.checkSession = checkSession
    window.validateSession = validateSession
})

</script>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
    display: none;
}

.no-scrollbar {
    -ms-overflow-style: none;
    scrollbar-width: none;
}

@media (max-width: 640px) {
    nav.tabs-container button {
        padding-left: 1.5rem;
        padding-right: 1.5rem;
        font-size: 0.875rem;
    }
}
</style>