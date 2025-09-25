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
            <div class="flex justify-between items-start mb-3 pb-2 border-b border-gray-200">
                <h3 class="text-sm font-semibold text-gray-700" v-html="popupMessage.title"></h3>
                <div class="relative w-6 h-6">
                    <svg v-if="popupTimer > 0" class="absolute top-0 left-0 w-full h-full rotate-[-90deg]"
                        viewBox="0 0 36 36">
                        <circle class="text-gray-300" stroke="currentColor" stroke-width="3" fill="transparent" r="16"
                            cx="18" cy="18" />
                        <circle class="text-[#f2a81d]" :stroke-dasharray="circumference" :stroke-dashoffset="offset"
                            stroke-linecap="round" stroke="currentColor" stroke-width="3" fill="transparent" r="16"
                            cx="18" cy="18" />
                    </svg>
                    <button @click="fecharPopup"
                        class="absolute inset-0 flex items-center justify-center text-gray-500 hover:text-gray-700 cursor-pointer"
                        aria-label="Fechar popup">
                        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
            </div>
            <p class="text-base font-medium" v-html="popupMessage.text"></p>
        </div>
    </div>

    <div class="min-h-screen mx-auto px-0 lg:px-5">
        <header
            class="fixed inset-x-0 top-0 z-50 h-[58px] w-full bg-[#414141] shadow flex items-center justify-between px-4">
            <div class="absolute left-1/2 transform -translate-x-1/2 text-xl font-bold text-gray-800">
                <img class="h-[30px]" src="https://bizout.com.br/questoes/assets/bizout_logo-BozI29p7.svg" alt="Logo">
            </div>
            <div class="flex items-center space-x-4">
                <ProfileMenu />
            </div>
        </header>

        <main class="mt-[78px] pt-[68px] p-1 lg:p-4">
            <div v-if="pending" class="flex justify-center items-center space-x-3 text-gray-600">
                <svg class="w-6 h-6 text-blue-500 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none"
                    viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
                </svg>
                <span>Procurando questão...</span>
            </div>

            <div v-else-if="posts && Object.keys(posts).length > 0">
                <div class="flex items-center mb-8">
                    <a href="/questoes"
                        class="flex items-center space-x-2 text-gray-800 font-medium hover:text-blue-600 transition-colors duration-200">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                        </svg>
                        <span>Mostrar todas as questões</span>
                    </a>
                </div>

                <div v-if="posts.length === 0" class="text-center text-gray-500 py-8">
                    <svg class="w-12 h-12 mx-auto mb-4 text-gray-300" fill="none" viewBox="0 0 24 24"
                        stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round"
                            d="M9 17v-6h6v6m2 4H7a2 2 0 01-2-2V7a2 2 0 012-2h5l5 5v9a2 2 0 01-2 2z" />
                    </svg>
                    <p class="text-lg font-semibold">Nenhuma questão encontrada.</p>
                </div>

                <div v-else>
                    <div class="flex items-baseline gap-1">
                        <span class="black sm:hidden text-sm pb-2 uppercase tracking-wide text-gray-500">C{{ posts[0][0]
                        }}</span>
                    </div>

                    <div
                        class="flex flex-col items-start gap-2 rounded-lg border border-gray-200 mb-4 bg-gray-50 shadow-sm p-4">
                        <div class="flex flex-wrap items-center gap-2 text-sm text-gray-700">
                            <span class="hidden md:inline text-sm uppercase tracking-wide text-gray-500 mr-2">C{{
                                posts[0][0] }}:</span>
                            <span v-for="(assuntos, assuntos_index) in posts[0][3]" :key="assuntos_index"
                                class="px-2 py-1 bg-gray-200 text-gray-800 text-sm rounded-md">
                                {{ assuntos }}
                            </span>
                        </div>
                        <div class="text-sm text-gray-700">
                            <span class="font-semibold">Ano:</span> {{ posts[0][2][3] }} - {{ posts[0][2][2] }}
                            <span class="font-semibold">Banca:</span> {{ posts[0][5] }}
                            <span class="font-semibold">Orgão:</span> {{ posts[0][6][0] }}
                        </div>
                    </div>

                    <p class="text-gray-900 p-4" v-html="posts[0][7][1]"></p>
                    <p v-if="posts[0][7][0]" class="text-gray-900 my-4 mt-0 p-4" v-html="posts[0][7][2]"></p>

                    <div class="relative p-4 pt-0 pb-2">
                        <div v-if="mkrsNC.includes(posts[0][0]) || mkrsNE.includes(posts[0][0])" @click="respondida()"
                            class="absolute w-full h-full bg-gray-100 opacity-50 rounded-lg cursor-no-drop"></div>

                        <button v-for="(alt, i) in posts[0][8][1]" :key="i" :data="i + 1"
                            @click="ans_select(posts[0][0], i)"
                            :class="['flex items-start gap-3 w-full rounded-lg border border-gray-300 px-4 py-2 text-left hover:bg-gray-100 focus:bg-indigo-50 focus:outline-none transition-colors cursor-pointer select-none mb-3',
                                { 'ans_selected': selectedIndices[posts[0][0]] === i, 'ans_selected-correct': correctIndices[posts[0][0]] === i, 'ans_selected-wrong': wrongIndices[posts[0][0]] === i },
                                { 'pointer-events-none': mkrsNC.includes(posts[0][0]) || mkrsNE.includes(posts[0][0]) }]">
                            <span class="font-semibold uppercase text-gray-700">{{ String.fromCharCode(97 + i)
                            }}.</span>
                            <span class="text-gray-800" v-html="alt"></span>
                            <span
                                v-if="selectedIndices[posts[0][0]] === i && (mkrsNC.includes(posts[0][0]) || mkrsNE.includes(posts[0][0]))"
                                class="ml-2 px-2 py-1 text-xs font-semibold rounded-full text-gray-800">
                                (Alternativa correta)
                            </span>
                        </button>
                    </div>

                    <div class="flex justify-center w-full">
                        <div
                            class="mb-4 inline-flex flex-wrap gap-3 text-gray-600 rounded-lg bg-gray-100/60 backdrop-blur-sm transition-all hover:border-gray-600 hover:shadow-md">
                            <button
                                class="cursor-pointer relative group flex items-center gap-1 px-2 py-2 rounded-md hover:bg-gray-200"
                                type="button" @click="responder(0, posts[0][0])"
                                :class="{ 'pointer-events-none': mkrsNC.includes(posts[0][0]) || mkrsNE.includes(posts[0][0]) }">
                                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"
                                    stroke-width="1.8">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                                </svg>
                                <span class="sm:inline">Responder</span>
                            </button>
                        </div>
                    </div>

                    <CommentSection :questionId="posts[0][0]" :qstDb="'conc'" :visible="true" class="mb-8" />
                </div>
            </div>

            <div v-else class="flex items-center gap-4 p-4 bg-red-100 text-red-800 rounded-lg border border-red-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M18.364 5.636l-1.414 1.414M5.636 5.636l1.414 1.414M12 2v2m0 16v2m6.364-6.364l-1.414-1.414M5.636 18.364l1.414-1.414M4 12H2m20 0h-2" />
                </svg>
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
            <nav class="flex flex-wrap gap-4 text-sm">
                <a href="/" class="hover:text-white transition">Início</a>
                <a href="/concursos" class="hover:text-white transition">Questões de Concursos</a>
            </nav>
            <div class="flex gap-4">
                <a href="https://twitter.com" target="_blank" aria-label="Twitter" class="hover:text-white transition">
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                        <path
                            d="M23 3a10.9 10.9 0 01-3.14 1.53A4.48 4.48 0 0022.4 1a9.06 9.06 0 01-2.88 1.1 4.52 4.52 0 00-7.86 4.13A12.8 12.8 0 013 2.1a4.52 4.52 0 001.4 6.05 4.52 4.52 0 01-2.05-.57v.06a4.52 4.52 0 003.63 4.43 4.52 4.52 0 01-2.04.08 4.52 4.52 0 004.22 3.13A9.05 9.05 0 013 19.54a12.77 12.77 0 006.92 2.03c8.3 0 12.84-6.88 12.84-12.84 0-.2 0-.42-.02-.63A9.22 9.22 0 0023 3z">
                        </path>
                    </svg>
                </a>
            </div>
        </div>
        <div class="text-center text-xs text-gray-500 py-4 border-t border-gray-800">
            Desenvolvido por Pedro Soares. Bizout. 2025.
        </div>
    </footer>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import ProfileMenu from '../../../components/ProfileMenu.vue'
import CommentSection from '../../../components/CommentSection.vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const id = route.params.id

const selectedIndices = reactive({})
const correctIndices = reactive({})
const wrongIndices = reactive({})
const alerts = reactive([])
const mkrsNC = ref(window.mkrsNC || [])
const mkrsNE = ref(window.mkrsNE || [])

const visible = ref(false)
const popupMessage = ref({ title: '', text: '', type: 'success' })
const popupTimer = ref(0)
const offset = ref(0)
const circumference = 2 * Math.PI * 16

const { data: posts, pending, error } = useFetch(`http://127.0.0.1:3030/id?idQst=${id}`, {
    transform: (data) => {
        let atual_qsts = data
        const n = Object.keys(atual_qsts)
        for (var i = 0; i < (n.length - 1); i++) {
            if (atual_qsts[i][2][2] == "1") atual_qsts[i][2][2] = "Nível Fundamental";
            else if (atual_qsts[i][2][2] == "2") atual_qsts[i][2][2] = "Nível Médio";
            else if (atual_qsts[i][2][2] == "3") atual_qsts[i][2][2] = "Nível Superior";
            let cont = atual_qsts[i][7].replace(/\\"/g, '"').split("#x@x#")
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
        const keys = Object.keys(atual_qsts)
        if (keys.length > 0) delete atual_qsts[keys[keys.length - 1]]
        return atual_qsts
    }
})

console.log(posts)

function ans_select(postId, index) {
    selectedIndices[postId] = selectedIndices[postId] === index ? undefined : index
    wrongIndices[postId] = undefined
    correctIndices[postId] = undefined
}

function responder(qst, qst_id) {
    const correct = parseInt(posts.value[qst][8][0].split("##")[1])
    if (selectedIndices[qst_id] === undefined) {
        mostrarPopup('Q' + qst_id, 'Nenhuma alternativa selecionada', 'warning')
        return
    }
    if (correct === (selectedIndices[qst_id] + 1)) {
        correctIndices[qst_id] = selectedIndices[qst_id]
        mostrarPopup('Q' + qst_id, 'Resposta correta!', 'success')
    }
    else {
        wrongIndices[qst_id] = selectedIndices[qst_id]
        mostrarPopup('Q' + qst_id, 'Resposta errada!', 'error')
    }
}

function mostrarPopup(title, text, type = 'success', time = 0) {
    popupMessage.value = { title, text, type }
    visible.value = true
    popupTimer.value = time
}

function fecharPopup() {
    visible.value = false
    popupTimer.value = 0
}

function respondida() {
    mostrarPopup("Atenção", "Você já respondeu esta questão.", "warning", 3000)
}
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
</style>
