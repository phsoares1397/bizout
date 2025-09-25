<template>
  <div class="fixed bottom-4 right-4 z-50 space-y-2 w-72">
    <div v-for="(alert, i) in alerts" :key="alert.id" :class="[
      'px-4 py-3 rounded-lg shadow-md text-sm transition-all duration-300',
      {
        'bg-green-100 text-green-800': alert.type === 'success',
        'bg-red-100 text-red-800': alert.type === 'error',
        'bg-yellow-100 text-yellow-800': alert.type === 'warning'
      }
    ]">
      <div v-html="alert.message"></div>
    </div>
  </div>

  <div v-if="visible" class="fixed inset-0 z-1000 flex items-center justify-center bg-black/40 transition-opacity">
    <div class="relative p-5 rounded-lg shadow-lg max-w-sm w-full bg-white text-gray-800 transition-all" :class="{
      'border-l-4 border-green-400': popupMessage.type === 'success',
      'border-l-4 border-red-400': popupMessage.type === 'error',
      'border-l-4 border-yellow-400': popupMessage.type === 'warning'
    }">
      <!-- Título com botão de fechar -->
      <div class="flex justify-between items-start mb-3 pb-2 border-b border-gray-200">
        <h3 class="text-sm font-semibold text-gray-700" v-html="popupMessage.title"></h3>
        <div class="relative w-6 h-6">
          <!-- círculo animado -->
          <svg v-if="popupTimer > 0" class="absolute top-0 left-0 w-full h-full rotate-[-90deg]" viewBox="0 0 36 36">
            <circle class="text-gray-300" stroke="currentColor" stroke-width="3" fill="transparent" r="16" cx="18"
              cy="18" />
            <circle class="text-[#f2a81d]" :stroke-dasharray="circumference" :stroke-dashoffset="offset"
              stroke-linecap="round" stroke="currentColor" stroke-width="3" fill="transparent" r="16" cx="18" cy="18" />
          </svg>

          <!-- botão de fechar -->
          <button @click="fecharPopup"
            class="absolute inset-0 flex items-center justify-center text-gray-500 hover:text-gray-700 cursor-pointer">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
      <!-- Mensagem -->
      <p class="text-base font-medium" v-html="popupMessage.text"></p>
    </div>
  </div>

  <div class="min-h-screen mx-auto px-0 lg:px-5">
    <!-- Navbar -->
    <header
      class="fixed inset-x-0 top-0 z-50 h-[58px] w-full bg-[#414141] shadow flex items-center justify-between px-4">
      <!-- Logo -->
      <div class="flex items-center space-x-2">
        <svg class="w-6 h-6 text-white cursor-pointer" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </div>

      <!-- Nome -->
      <div class="absolute left-1/2 transform -translate-x-1/2 text-xl font-bold text-gray-800">
        <img class="h-[30px]" src="../src/assets/bizout_logo.svg">
      </div>

      <!-- Ícones -->
      <div class="flex items-center space-x-4">
        <ProfileMenu />
      </div>
    </header>

    <Tabs :key="siteVersion" />

    <!-- Conteúdo -->
    <main class="mt-[8px] pt-[18px] p-1 lg:p-4"> <!-- w-full md:max-w-[90vw] mx-auto -->
      <!-- Tela de carregamento -->
      <div v-if="loading && !loaded" class="flex justify-center items-center space-x-3 text-gray-600">
        <svg class="w-6 h-6 text-blue-500 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none"
          viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
        </svg>
        <span>Procurando questões...</span>
      </div>

      <!-- Lista quando já chegou resposta -->
      <div v-else-if="loaded">

        <!-- Se não houver questões -->
        <div v-if="atual_num_qsts === 0 || posts.length === 0" class="text-center text-gray-500 py-8">
          <svg class="w-12 h-12 mx-auto mb-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor"
            stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round"
              d="M9 17v-6h6v6m2 4H7a2 2 0 01-2-2V7a2 2 0 012-2h5l5 5v9a2 2 0 01-2 2z" />
          </svg>
          <p class="text-lg font-semibold">Nenhuma questão encontrada.</p>
          <p class="text-sm">Tente alterar os filtros ou a pesquisa para encontrar questões.</p>
        </div>

        <!-- Se houver questões -->
        <div v-else>
          <!-- Quantidade de questões -->
          <div class="flex justify-between items-center gap-2 text-gray-600 text-sm pl-1">
            <!-- Lado esquerdo: número de questões -->
            <div class="flex items-center gap-2">
              <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"
                stroke-width="1.6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h7" />
              </svg>
              <span>
                <span class="font-semibold text-gray-800">{{ atual_num_qsts }}</span>
                questões disponíveis
              </span>
            </div>

            <!-- Lado direito: botão de modo de visualização -->
            <div class="flex items-center gap-3 mb-4">
              <span class="text-gray-700 font-medium">Modo de visualização:</span>

              <button @click="toggleColumns" class="flex items-center gap-2 px-3 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg 
                 hover:bg-gray-100 transition-colors shadow-sm">
                <svg v-if="isTwoColumns" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none"
                  viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <!-- Ícone de 2 colunas -->
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 4h6v16H4V4zm10 0h6v16h-6V4z" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor" stroke-width="2">
                  <!-- Ícone de 1 coluna -->
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 4h16v16H4V4z" />
                </svg>
                <span class="sr-only">Toggle columns</span>
              </button>
            </div>
          </div>

          <!-- Paginação superior -->
          <Pagination v-model:page="pagina" :total="atual_num_qsts" :por-pagina="10" class="mt-3 mb-6" />

          <!-- Lista de questões -->
          <div :class="isTwoColumns ? 'grid grid-cols-1 md:grid-cols-2 gap-6' : 'grid grid-cols-1 gap-6'">
            <div v-for="(post, index) in posts" :key="index" :qst-num="index"
              class="mb-4 p-0 bg-white relative transition-opacity">

              <div v-if="mkrsNC.includes(parseInt(post[0]))"
                class="absolute top-4 right-2 px-2 py-1 bg-green-100 text-green-800 text-xs font-semibold rounded">
                Você acertou ✅
              </div>
              <div v-else-if="mkrsNE.includes(parseInt(post[0]))"
                class="absolute top-4 right-2 px-2 py-1 bg-amber-100 text-amber-800 text-xs font-semibold rounded">
                Você errou ❌
              </div>

              <div class="flex items-baseline gap-1">
                <span class="black sm:hidden text-sm pb-2 uppercase tracking-wide text-gray-500"
                  style="margin: auto;">C{{
                    post[0] }}</span>
              </div>

              <div
                class="flex flex-col items-baseline gap-1 rounded-lg border border-gray-200 mb-4 bg-gray-50 shadow-sm p-4">
                <span class="overflow-x-auto max-w-full flex items-center text-sm pb-2 tracking-wide text-gray-700">
                  <span class="hidden md:block text-sm uppercase tracking-wide text-gray-500 mr-2">C{{ post[0] }}:
                  </span>
                  <a v-for="(assuntos, assuntos_index) in post[3]" :key="assuntos_index" href=""
                    class="whitespace-nowrap px-[10px] py-[6px] mr-2 bg-gray-200 text-gray-800 text-sm rounded-md cursor-pointer">{{
                      assuntos }} </a>
                </span>
                <span class="text-sm pb-2 tracking-wide text-gray-700">
                  <span class="font-semibold mt-15">Ano:</span> {{ post[2][3] }} - {{ post[2][2] }}
                  <span class="font-semibold">Banca: </span>{{ post[5] }}
                  <span class="font-semibold">Orgão: </span>{{ post[6][0] }}
                </span>
              </div>

              <!-- <div class="border-t border-gray-300 mt-2 mb-4"></div> -->

              <p class="text-gray-900 p-4" v-html="post[7][1]"></p>
              <p v-if="post[7][0]" class="text-gray-900 my-4 p-4" v-html="post[7][2]"></p>

              <div class="relative p-4 pb-2">
                <div v-if="mkrsNC.includes(parseInt(post[0])) || mkrsNE.includes(post[0])" @click="respondida()"
                  class="absolute w-full h-full bg-gray-100 opacity-50 rounded-lg cursor-no-drop">
                </div>
                <button v-for="(alt, i) in post[8][1]" :key="i" :data="i + 1" @click="ans_select(post[0], i)" :class="['flex items-start gap-3 w-full rounded-lg border border-gray-300 px-4 py-2 text-left hover:bg-gray-100 focus:bg-indigo-50 focus:outline-none transition-colors cursor-pointer select-none mb-3',
                  { 'ans_selected': selectedIndices[post[0]] === i, 'ans_selected-correct': correctIndices[post[0]] === i, 'ans_selected-wrong': wrongIndices[post[0]] === i },
                  { 'pointer-events-none': mkrsNC.includes(post[0]) || mkrsNE.includes(post[0]) }]">
                  <span class="font-semibold uppercase text-gray-700">{{ String.fromCharCode(97 + i) }}.</span>
                  <span class="text-gray-800" v-html="alt"></span>
                  <span v-if="selectedIndices[post[0]] === i && (mkrsNC.includes(post[0]) || mkrsNE.includes(post[0]))"
                    class="ml-2 px-2 py-1 text-xs font-semibold rounded-full text-gray-800">
                    (Alternativa correta)
                  </span>
                </button>
              </div>

              <!-- Ações da questão -->
              <div class="relative w-full flex justify-center">
                <div class="mb-4 inline-flex justify-center gap-3 text-gray-600 rounded-lg bg-white p-1">
                  <button class="cursor-pointer inline-flex items-center gap-1 px-2 py-2 rounded-md hover:bg-gray-100"
                    type="button" @click="responder(index, post[0])"
                    :class="{ 'pointer-events-none': mkrsNC.includes(post[0]) || mkrsNE.includes(post[0]) }">
                    <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                    </svg>
                    <span class="hidden sm:inline">Responder</span>
                  </button>

                  <button class="cursor-pointer inline-flex items-center gap-1 px-2 py-2 rounded-md hover:bg-gray-100"
                    type="button" @click="mostrarComentarios[post[0]] = !mostrarComentarios[post[0]]">
                    <svg v-if="!mostrarComentarios[post[0]]" class="w-5 h-5" fill="none" viewBox="0 0 24 24"
                      stroke="currentColor" stroke-width="1.8">
                      <path stroke-linecap="round" stroke-linejoin="round"
                        d="M8 10h8m-8 4h5m-5 8l-4-4H5a7 7 0 100-14h8a7 7 0 017 7v1a7 7 0 01-7 7H8z" />
                    </svg>
                    <svg v-else class="w-5 h-5 text-[#f2a81d]" fill="none" viewBox="0 0 24 24" stroke="currentColor"
                      stroke-width="2.6">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    <span class="hidden sm:inline">Comentários ({{ post[9] }})</span>
                  </button>

                  <button class="cursor-pointer inline-flex items-center gap-1 px-2 py-2 rounded-md hover:bg-gray-100"
                    type="button">
                    <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                      <path stroke-linecap="round" stroke-linejoin="round"
                        d="M5 5v14h14V7.828a2 2 0 00-.586-1.414l-2.828-2.828A2 2 0 0014.172 3H5zm4 9h6v5H9v-5z" />
                    </svg>
                    <span class="hidden sm:inline">Salvar</span>
                  </button>

                  <button class="cursor-pointer inline-flex items-center gap-1 px-2 py-2 rounded-md hover:bg-gray-100"
                    type="button">
                    <svg height="16" id="svg2" version="1.1" width="16" xmlns="http://www.w3.org/2000/svg"
                      xmlns:cc="http://creativecommons.org/ns#" xmlns:dc="http://purl.org/dc/elements/1.1/"
                      xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
                      xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
                      xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
                      xmlns:svg="http://www.w3.org/2000/svg">
                      <defs id="defs4" />
                      <g id="layer1" transform="translate(0,-1036.3622)">
                        <path
                          d="m -22.410713,-3.3303571 a 2.3660715,2.3660715 0 1 1 -4.732143,0 2.3660715,2.3660715 0 1 1 4.732143,0 z"
                          id="path2985" style="fill:currentColor;fill-opacity:1;stroke:none"
                          transform="matrix(0.84528301,0,0,0.84528301,33.943395,1042.1773)" />
                        <path
                          d="m -22.410713,-3.3303571 a 2.3660715,2.3660715 0 1 1 -4.732143,0 2.3660715,2.3660715 0 1 1 4.732143,0 z"
                          id="path2985-1" style="fill:currentColor;fill-opacity:1;stroke:none"
                          transform="matrix(0.84528301,0,0,0.84528301,33.943395,1052.1773)" />
                        <path
                          d="m -22.410713,-3.3303571 a 2.3660715,2.3660715 0 1 1 -4.732143,0 2.3660715,2.3660715 0 1 1 4.732143,0 z"
                          id="path2985-1-7" style="fill:currentColor;fill-opacity:1;stroke:none"
                          transform="matrix(0.84528301,0,0,0.84528301,23.943395,1047.1773)" />
                        <path d="M 13,3 3,8 13,13" id="path3791"
                          style="fill:none;stroke:currentColor;stroke-width:1px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1"
                          transform="translate(0,1036.3622)" />
                      </g>
                    </svg>
                    <span class="hidden sm:inline">Compartilhar</span>
                  </button>

                  <button class="cursor-pointer inline-flex items-center gap-1 px-2 py-2 rounded-md hover:bg-gray-100"
                    type="button">
                    <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M6 3v18l6-4 6 4V3H6z" />
                    </svg>
                    <span class="hidden sm:inline">Reportar Erro</span>
                  </button>
                </div>
              </div>

              <CommentSection :questionId="post[0]" :qstDb="'conc'" :visible="mostrarComentarios[post[0]]"
                class="mb-8" />
            </div>
          </div>

          <!-- Paginação inferior -->
          <Pagination v-model:page="pagina" :total="atual_num_qsts" :por-pagina="10" class="mt-3 mb-6" />
        </div>
      </div>

      <div v-else class="flex items-center gap-4 p-4 bg-red-100 text-red-800 rounded-lg border border-red-300">
        <!-- Ícone de erro -->
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0" fill="none" viewBox="0 0 24 24"
          stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M18.364 5.636l-1.414 1.414M5.636 5.636l1.414 1.414M12 2v2m0 16v2m6.364-6.364l-1.414-1.414M5.636 18.364l1.414-1.414M4 12H2m20 0h-2" />
        </svg>

        <!-- Texto -->
        <div>
          <p class="font-semibold">Problema de conexão</p>
          <p class="text-sm">Verifique sua internet ou tente novamente mais tarde.</p>
        </div>
      </div>

    </main>
  </div>

  <footer class="bg-gray-900 text-gray-200 mt-8">
    <div
      class="max-w-7xl mx-auto px-4 py-8 flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">

      <!-- Links principais -->
      <nav class="flex flex-wrap gap-4 text-sm">
        <a href="http://localhost" class="hover:text-white transition">Início</a>
        <a href="/concursos" class="hover:text-white transition">Questões de Concursos</a>
        <a href="/vestibulares" class="hover:text-white transition">Questões de Vestibulares</a>
        <a href="/enem" class="hover:text-white transition">Questões de Enem</a>
      </nav>

      <!-- Redes sociais -->
      <div class="flex gap-4">
        <a href="https://twitter.com" target="_blank" aria-label="Twitter" class="hover:text-white transition">
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path
              d="M23 3a10.9 10.9 0 01-3.14 1.53A4.48 4.48 0 0022.4 1a9.06 9.06 0 01-2.88 1.1 4.52 4.52 0 00-7.86 4.13A12.8 12.8 0 013 2.1a4.52 4.52 0 001.4 6.05 4.52 4.52 0 01-2.05-.57v.06a4.52 4.52 0 003.63 4.43 4.52 4.52 0 01-2.04.08 4.52 4.52 0 004.22 3.13A9.05 9.05 0 013 19.54a12.77 12.77 0 006.92 2.03c8.3 0 12.84-6.88 12.84-12.84 0-.2 0-.42-.02-.63A9.22 9.22 0 0023 3z">
            </path>
          </svg>
        </a>
        <a href="https://youtube.com" target="_blank" aria-label="YouTube" class="hover:text-white transition">
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path
              d="M19.615 3.184a2.99 2.99 0 00-2.113-2.113C15.544.5 12 0.5 12 0.5s-3.544 0-5.502.571a2.99 2.99 0 00-2.113 2.113C3.5 5.146 3.5 8.7 3.5 8.7s0 3.555.885 5.515a2.99 2.99 0 002.113 2.113C8.456 17.5 12 17.5 12 17.5s3.544 0 5.502-.572a2.99 2.99 0 002.113-2.113c.885-1.961.885-5.515.885-5.515s0-3.554-.885-5.501zM9.75 12.25V5.25l6 3.5-6 3.5z">
            </path>
          </svg>
        </a>
        <a href="https://instagram.com" target="_blank" aria-label="Instagram" class="hover:text-white transition">
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path
              d="M12 2.163c3.204 0 3.584.012 4.85.07 1.17.054 1.97.24 2.43.403a4.92 4.92 0 011.787 1.14 4.92 4.92 0 011.14 1.788c.163.462.35 1.26.403 2.43.058 1.265.07 1.645.07 4.849s-.012 3.584-.07 4.85c-.054 1.17-.24 1.97-.403 2.43a4.902 4.902 0 01-1.14 1.787 4.902 4.902 0 01-1.788 1.14c-.462.163-1.26.35-2.43.403-1.266.058-1.646.07-4.85.07s-3.584-.012-4.85-.07c-1.17-.054-1.97-.24-2.43-.403a4.902 4.902 0 01-1.787-1.14 4.902 4.902 0 01-1.14-1.788c-.163-.462-.35-1.26-.403-2.43C2.175 15.584 2.163 15.204 2.163 12s.012-3.584.07-4.85c.054-1.17.24-1.97.403-2.43a4.92 4.92 0 011.14-1.787 4.92 4.92 0 011.788-1.14c.462-.163 1.26-.35 2.43-.403C8.416 2.175 8.796 2.163 12 2.163zm0 1.838c-3.19 0-3.565.012-4.82.069-1.07.048-1.65.22-2.036.366a3.042 3.042 0 00-1.104.723 3.042 3.042 0 00-.723 1.104c-.147.386-.318.966-.366 2.036-.058 1.255-.069 1.63-.069 4.82s.012 3.565.069 4.82c.048 1.07.22 1.65.366 2.036.17.415.387.765.723 1.104.338.336.688.554 1.104.723.386.147.966.318 2.036.366 1.255.058 1.63.069 4.82.069s3.565-.012 4.82-.069c1.07-.048 1.65-.22 2.036-.366a3.042 3.042 0 001.104-.723 3.042 3.042 0 00.723-1.104c.147-.386.318-.966.366-2.036.058-1.255.069-1.63.069-4.82s-.012-3.565-.069-4.82c-.048-1.07-.22-1.65-.366-2.036a3.042 3.042 0 00-.723-1.104 3.042 3.042 0 00-1.104-.723c-.386-.147-.966-.318-2.036-.366-1.255-.058-1.63-.069-4.82-.069zm0 4.838a5 5 0 110 10 5 5 0 010-10zm0 1.838a3.162 3.162 0 100 6.324 3.162 3.162 0 000-6.324zm6.406-.338a1.18 1.18 0 11-2.36 0 1.18 1.18 0 012.36 0z">
            </path>
          </svg>
        </a>
      </div>
    </div>

    <!-- Rodapé de crédito -->
    <div class="text-center text-xs text-gray-500 py-4 border-t border-gray-800">
      Desenvolvido por Pedro Soares. Bizout. 2025.
    </div>
  </footer>

</template>

<script setup>
import { ref, onMounted, watch, reactive, computed, nextTick } from 'vue'
import ProfileMenu from "./components/ProfileMenu_Sidebar.vue"
import Tabs from './components/Tabs.vue'
import Pagination from './components/Pagination.vue'
import CommentSection from './components/CommentSection.vue'

const isTwoColumns = ref(false) // começa em modo 2 colunas

function toggleColumns() {
  isTwoColumns.value = !isTwoColumns.value
}

const posts = ref([])
const loading = ref(true)
const loaded = ref(false)
var atual_qsts
const pagina = ref(1)
const atual_num_qsts = ref()
const selectedIndices = reactive({})
const correctIndices = reactive({})
const wrongIndices = reactive({})
const alert = reactive({
  message: '',
  type: '',     // 'success', 'error', 'warning'
  visible: false
})

const mkrsNC = ref(window.mkrsNC || [])
const mkrsNE = ref(window.mkrsNE || [])

let nextId = 0

const alerts = reactive([])
const mostrarComentarios = reactive({})

const visible = ref(false)
const popupMessage = ref({ title: '', text: '', type: 'success' })
const popupTimer = ref(0)
const offset = ref(0)
const circumference = 2 * Math.PI * 16
const marcadorSelecionado = ref(null)
let timer

const siteVersion = ref(0);

// sempre que quiser remontar tudo:
function reloadComponents() {
  siteVersion.value++;
}

function mostrarPopup(title, text, type = 'success', time = 0) {
  popupMessage.value = { title, text, type }
  visible.value = true
  popupTimer.value = time
  offset.value = 0

  if (time > 0) {
    const interval = 50
    const totalSteps = time / interval
    let step = 0

    timer = setInterval(() => {
      step++
      offset.value = (step / totalSteps) * circumference
      if (step >= totalSteps) {
        clearInterval(timer)
        visible.value = false
        popupTimer.value = 0
      }
    }, interval)
  }
}

function fecharPopup() {
  visible.value = false
  popupTimer.value = 0
  offset.value = 0
  try {
    clearInterval(timer)
  } catch (e) { }
}

function showAlert(message, type = 'success', duration = 3000) {
  const id = nextId++

  alerts.push({ id, message, type })

  // Remove automaticamente após 'duration'
  setTimeout(() => {
    const index = alerts.findIndex(alert => alert.id === id)
    if (index !== -1) alerts.splice(index, 1)
  }, duration)
}

function userData() {
  try {
    return JSON.parse(localStorage.getItem('user'))
  } catch (e) {
    return null
  }
}

function fecharComentario(comentarioId) {
  mostrarComentarios[comentarioId] = false
}

const filtrosPadrao = reactive({
  disciplina: [], assuntos: [], orgao: [], cargo: [],
  banca: [], ano: [], nivel: [], tipo: [], dificuldade: []
})

// mapas de conversão
const mapNivel = { Fundamental: 1, Médio: 2, Superior: 3 }
const mapTipo = { "Certo ou Errado": 2, "Múltipla Escolha": 1 }
const mapDificuldade = { "Muito Fácil": 1, "Fácil": 2, "Médio": 3, "Difícil": 4, "Muito Difícil": 5 }

async function loadPage(pg = 1, set = 0, filtros) {
  mkrsNC.value = window.mkrsNC || []
  mkrsNE.value = window.mkrsNE || []
  try {
    loading.value = true
    loaded.value = false
    pagina.value = pg

    const params = new URLSearchParams()
    params.append("pg", pg)
    params.append("set", set)

    // Extrai IDs de disciplinas
    const disciplinaIds = filtros.disciplina.map(d => d.id)

    // Map de filtros: chave do objeto -> nome do parâmetro na URL
    const filtroMap = {
      orgao: "org",
      cargo: "crg",
      banca: "bnc",
      ano: "ano",
      nivel: "nv",
      tipo: "tp",
      dificuldade: "dfcd",
    }

    // Função auxiliar para adicionar array não vazio
    const addParam = (key, arr) => {
      if (arr.length > 0) params.append(key, JSON.stringify(arr))
    }

    // asst_rt: apenas disciplinas
    addParam("asst_rt", disciplinaIds)

    // asst: assuntos + disciplinas, sem duplicados
    const assuntosIds = window.filtros.assuntos.map(a => a[1]);
    const asstIds = Array.from(new Set([...assuntosIds]))
    addParam("asst", asstIds)

    // adiciona os demais filtros dinamicamente
    for (const [filtroKey, paramName] of Object.entries(filtroMap)) {
      let valores = filtros[filtroKey]

      // converte quando for nível, tipo ou dificuldade
      if (filtroKey === "nivel") {
        valores = valores.map(v => mapNivel[v] ?? v)
      } else if (filtroKey === "tipo") {
        valores = valores.map(v => mapTipo[v] ?? v)
      } else if (filtroKey === "dificuldade") {
        valores = valores.map(v => mapDificuldade[v] ?? v)
      }

      addParam(paramName, valores)
    }

    const url = `http://localhost:3030/v?${params.toString()}`
    const res = await fetch(url)
    if (!res.ok) throw new Error("Erro ao buscar posts")

    console.log(url)

    atual_qsts = await res.json()
    parseQstsData()
    posts.value = atual_qsts
    window.posts = posts

    // limpa comentários antigos
    Object.keys(mostrarComentarios).forEach(k => delete mostrarComentarios[k])

    setTimeout(() => {
      loaded.value = true
    }, 1000)
  } catch (e) {
    console.error("Erro ao carregar questões:", e)
  } finally {
    setTimeout(() => {
      loading.value = false
      window.pegarQuestoesMarcadas()
    }, 1000)
  }
}

function ans_select(postId, index) {
  if (selectedIndices[postId] == index) selectedIndices[postId] = undefined
  else selectedIndices[postId] = index

  wrongIndices[postId] = undefined
  correctIndices[postId] = undefined
}

async function marcarQuestaoCorreta(markerId, qst_id, banco) {
  const data = {
    marker_id: markerId,
    qid: qst_id,
    banco: banco // ex: 'conc', 'vest' ou 'enem'
  }

  try {
    const res = await fetch('http://localhost:4040/mkrs-insert-correct', {
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

async function marcarQuestaoErrada(markerId, qst_id, banco) {
  const data = {
    marker_id: markerId,
    qid: qst_id,
    banco: banco // ex: 'conc', 'vest' ou 'enem'
  }

  try {
    const res = await fetch('http://localhost:4040/mkrs-insert-wrong', {
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

function pegarQuestoesMarcadas() {
  for (const key in atual_qsts) {
    if (!Object.hasOwn(atual_qsts, key)) continue

    const qst_id = atual_qsts[key][0] // id da questão
    var correct = parseInt(atual_qsts[key][8][0].split("##")[1]) - 1

    if (mkrsNC.value.includes(qst_id) || mkrsNE.value.includes(qst_id)) {
      selectedIndices[qst_id] = correct
    }
  }
}

const svgWarning = `
<span class="inline-flex items-center space-x-1">
  <svg xmlns="http://www.w3.org/2000/svg" 
      class="w-5 h-5 text-yellow-500 inline-block"
      viewBox="0 0 20 20" fill="currentColor">
    <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.721-1.36 3.486 0l6.518 11.602c.75 1.336-.213 2.999-1.742 2.999H3.48c-1.53 0-2.492-1.663-1.742-2.999L8.257 3.1zM11 14a1 1 0 10-2 0 1 1 0 002 0zm-1-2a.75.75 0 01-.75-.75v-3.5a.75.75 0 011.5 0v3.5A.75.75 0 0110 12z" clip-rule="evenodd" />
  </svg>
  <span class="mt-[1px]">Atenção</span>
</span>
`;

function respondida() {
  mostrarPopup(`${svgWarning}`, "Você já respondeu essa questão neste marcador.", "warning", 3000)
}

async function responder(qst, qst_id) {
  var correct = parseInt(atual_qsts[qst][8][0].split("##")[1])
  if (selectedIndices[qst_id] == undefined) {
    window.mostrarPopup('Q' + qst_id, 'Nenhuma alternativa selecionada', 'warning', 2000)
    return
  }
  if (correct == (selectedIndices[qst_id] + 1)) {
    wrongIndices[qst_id] = undefined
    correctIndices[qst_id] = selectedIndices[qst_id]
    if (window.marcadorSelecionado.value != null) {
      const resposta = await marcarQuestaoCorreta(window.marcadorSelecionado.value, qst_id, 'conc')
      if (resposta[0] === 200) {
        window.mkrsNC.push(qst_id)
        mkrsNC.value = window.mkrsNC
        window.mostrarPopup('Q' + qst_id, 'Resposta correta!', 'success', 2000)
      } else {
        console.warn('Erro:', resposta[1])
      }
    } else {
      window.mostrarPopup('Q' + qst_id, 'Resposta correta!', 'success', 2000)
    }
  }
  else {
    correctIndices[qst_id] = undefined
    wrongIndices[qst_id] = selectedIndices[qst_id]
    if (window.marcadorSelecionado.value != null) {
      const resposta = await marcarQuestaoErrada(window.marcadorSelecionado.value, qst_id, 'conc')
      if (resposta[0] === 200) {
        window.mkrsNE.push(qst_id)
        mkrsNE.value = window.mkrsNE
        window.mostrarPopup('Q' + qst_id, 'Resposta errada!', 'error', 2000)
      } else {
        console.warn('Erro:', resposta[1])
      }
    } else {
      window.mostrarPopup('Q' + qst_id, 'Resposta errada!', 'error', 2000)
    }
  }
  window.pegarQuestoesMarcadas()
}

function parseQstsData() {
  const n = Object.keys(atual_qsts)
  let cont
  for (var i = 0; i < (n.length - 1); i++) {
    if (atual_qsts[i][2][2] == "1") atual_qsts[i][2][2] = "Nível Fundamental";
    else if (atual_qsts[i][2][2] == "2") atual_qsts[i][2][2] = "Nível Médio";
    else if (atual_qsts[i][2][2] == "3") atual_qsts[i][2][2] = "Nível Superior";
    cont = atual_qsts[i][7].replace(/\\"/g, '"').split("#x@x#")
    if (cont[0] == "") cont = ["false", cont[0], cont[1]]
    else cont = ["true", cont[0], cont[1]]
    atual_qsts[i][7] = cont
    cont = atual_qsts[i][8].replace(/\\"/g, '"').split("#x@x#");
    atual_qsts[i][8] = [cont[cont.length - 1], cont.slice(0, -1)]
    cont = JSON.parse(atual_qsts[i][5])
    atual_qsts[i][5] = cont[0] + " " + cont[1]
    cont = JSON.parse(atual_qsts[i][6])
    atual_qsts[i][6] = [cont[0] + " - " + cont[1], cont[2], cont[3]]
    cont = JSON.parse(atual_qsts[i][4])
    atual_qsts[i][4] = cont[0]
    cont = JSON.parse(atual_qsts[i][3])
    atual_qsts[i][3] = cont
  }
  atual_num_qsts.value = atual_qsts[n.length - 1]
  atual_qsts = removeLastProp(atual_qsts)
}

function removeLastProp(obj) {
  const keys = Object.keys(obj)
  if (keys.length === 0) return
  const lastKey = keys[keys.length - 1]
  delete obj[lastKey]
  return obj
}

onMounted(async () => {
  window.userData = userData
  window.mostrarPopup = mostrarPopup
  window.showAlert = showAlert
  window.fecharComentario = fecharComentario
  window.atual_qsts = atual_qsts
  window.loadPage = loadPage
  window.pagina = pagina
  localStorage.removeItem('selecionados')
  window.pegarQuestoesMarcadas = pegarQuestoesMarcadas
  window.reloadComponents = reloadComponents
  window.marcadorSelecionado = marcadorSelecionado

  watch(
    pagina, // reactive/ref que você quer observar
    (newVal, oldVal) => {
      setTimeout(() => {
        window.scrollTo({ top: 0, behavior: 'smooth' })
      }, 0)
      loadPage(newVal, 0, window.filtrosNormalizados)
    },
    { immediate: true }
  )
})
</script>

<style>
.ans_selected {
  background-color: #f2a81d59;
}

.ans_selected-correct {
  background-color: #1df22a59 !important;
}

.ans_selected-wrong {
  background-color: #f21d1d59 !important;
}

@keyframes fade-in-out {
  0% {
    opacity: 0;
    transform: scale(0.95);
  }

  10% {
    opacity: 1;
    transform: scale(1);
  }

  90% {
    opacity: 1;
  }

  100% {
    opacity: 0;
    transform: scale(0.95);
  }
}

.animate-fade-in-out {
  animation: fade-in-out 2.5s ease-in-out forwards;
}
</style>