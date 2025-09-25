<script setup>
import { ref, reactive, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import Quill from 'quill'
import 'quill/dist/quill.snow.css'
import { QuillDeltaToHtmlConverter } from 'quill-delta-to-html'

const props = defineProps({
    questionId: Number,
    qstDb: String,
    visible: Boolean
})

const quillContainer = ref(null)
const loading = ref(false)
const error = ref(false)
const quillInstance = ref(null)
const quillInstanceEditor = reactive({})
const editingBackup = reactive({})
const toolbarOptions = [
    ["bold", "italic", "underline", "strike"],
    ["blockquote"],
    [{ list: "ordered" }, { list: "bullet" }],
    [{ script: "sub" }, { script: "super" }],
    [{ indent: "-1" }, { indent: "+1" }],
    [{ size: ["small", !1, "large", "huge"] }],
    [{ color: [] }, { background: [] }],
    [{ align: [] }],
    ["clean"]
]

const comentarios = ref([])
const userId = ref(null)

// Controle do CAPTCHA
const mostrarCaptcha = ref(false)
const confirmLoadingCaptcha = ref(false)
const captchaAction = ref(null)

// Controle de edição
const editingCommentMap = reactive({})
const tempCommentId = ref(null)
const tempQuestionId = ref(null)

function fecharCaptchaModal() {
    mostrarCaptcha.value = false
    captchaAction.value = null
    confirmLoadingCaptcha.value = false
}

// Retorna primeiros nomes (mais legível)
function getFirstNames(fullName) {
    if (!fullName) return ""
    const parts = fullName.trim().split(/\s+/)
    if (parts.length === 1) return parts[0]
    if (parts.length >= 3 && parts[1].length <= 2) return `${parts[0]} ${parts[1]} ${parts[2]}`
    return `${parts[0]} ${parts[1]}`
}

// Envia comentários com CAPTCHA
async function enviarComentarioComCaptcha() {
    if (!captchaAction.value) return

    confirmLoadingCaptcha.value = true

    try {
        switch (captchaAction.value) {
            case 'edit_comment':
                {
                    const sendDataEdit = [
                        "conc",
                        tempCommentId.value,
                        JSON.stringify(quillInstanceEditor[tempCommentId.value].getContents())
                    ]
                    const payloadEdit = { data: JSON.stringify(sendDataEdit), token: 'ddd' }
                    const response = await fetch(`https://tools.bizout.com.br/cmtse`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                        body: new URLSearchParams(payloadEdit),
                        credentials: "include"
                    })
                    const result = await response.json()
                    if (result[0] === 200) {
                        window.mostrarPopup('Pronto!', 'Comentário editado com sucesso!', 'success', 6000)
                        editingCommentMap[tempCommentId.value] = false
                        window.fecharComentario(props.questionId)
                    } else if (result[0] === 111) {
                        window.mostrarPopup('Erro', 'Captcha inválido, tente novamente.', 'error', 6000)
                    } else {
                        console.error(result)
                        window.mostrarPopup('Erro', 'Erro desconhecido.', 'error', 6000)
                    }
                }
                break

            case 'create_comment':
                {
                    const delta = quillInstance.value.getContents()
                    if (!delta.ops || (delta.ops.length === 1 && delta.ops[0].insert.trim() === "")) {
                        window.mostrarPopup('Aviso', 'Escreva algo antes de enviar.', 'warning', 4000)
                        break
                    }
                    const sendDataCreate = [
                        props.questionId + "",
                        JSON.stringify(delta),
                        "conc"
                    ]
                    const payloadCreate = { data: JSON.stringify(sendDataCreate), token: 'ddd' }
                    const response = await fetch(`https://tools.bizout.com.br/cmtsc`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                        body: new URLSearchParams(payloadCreate),
                        credentials: "include"
                    })
                    const result = await response.json()
                    if (result[0] === 200) {
                        window.mostrarPopup('Pronto!', 'Comentário publicado!', 'success', 6000)

                        Object.keys(window.posts.value).forEach(key => {
                            const post = window.posts.value[key]
                            if (post[0] == props.questionId) post[9] += 1
                        })

                        window.fecharComentario(props.questionId)
                    } else if (result[0] === 111) {
                        window.mostrarPopup('Erro', 'Captcha inválido.', 'error', 6000)
                    } else {
                        console.error(result)
                        window.mostrarPopup('Erro', 'Erro desconhecido.', 'error', 6000)
                    }
                }
                break

            case 'remove_comment':
                {
                    const sendDataDelete = [
                        tempQuestionId.value + "",
                        tempCommentId.value + "",
                        "conc"
                    ]
                    const payloadDelete = { data: JSON.stringify(sendDataDelete), token: 'ddd' }
                    const response = await fetch(`https://tools.bizout.com.br/cmtsd`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                        body: new URLSearchParams(payloadDelete),
                        credentials: "include"
                    })
                    const result = await response.json()
                    if (result[0] === 200) {
                        window.mostrarPopup('Pronto!', 'Comentário removido!', 'success', 6000)
                        comentarios.value = comentarios.value.filter(c => c[0] != tempCommentId.value)
                        Object.keys(window.posts.value).forEach(key => {
                            const post = window.posts.value[key]
                            if (post[0] == props.questionId) post[9] -= 1
                        })
                    } else {
                        console.error(result)
                        window.mostrarPopup('Erro', 'Erro desconhecido.', 'error', 6000)
                    }
                }
                break
        }
    } catch (err) {
        console.error(err)
        window.mostrarPopup('Erro', 'Erro de conexão.', 'error', 6000)
    } finally {
        mostrarCaptcha.value = false
        confirmLoadingCaptcha.value = false
        captchaAction.value = null
    }
}

// ---------------- UTILITÁRIOS ----------------

function deltaToHtml(delta) {
    const converter = new QuillDeltaToHtmlConverter(delta.ops, {})
    return converter.convert()
}

function initQuill(placeholder = 'Escreva seu comentário...') {
    if (!quillContainer.value || quillInstance.value) return
    quillInstance.value = new Quill(quillContainer.value, {
        theme: 'snow',
        placeholder,
        modules: { toolbar: toolbarOptions }
    })
}

function destroyQuill() {
    if (quillInstance.value) {
        quillInstance.value.disable()
        quillInstance.value = null
    }
}

// ---------------- AÇÕES ----------------

function publicarComentario() {
    captchaAction.value = 'create_comment'
    mostrarCaptcha.value = true
}

function editarComentario(comentarioId, event) {
    const editor = event.currentTarget.parentElement.previousElementSibling
    editingBackup[comentarioId] = editor.cloneNode(true)
    editingCommentMap[comentarioId] = true
    tempCommentId.value = comentarioId
    quillInstanceEditor[comentarioId] = new Quill(editor, {
        theme: 'snow',
        placeholder: 'Edite seu comentário...',
        modules: { toolbar: toolbarOptions }
    })
}

function salvarComentario(comentarioId) {
    tempCommentId.value = comentarioId
    captchaAction.value = 'edit_comment'
    mostrarCaptcha.value = true
}

function cancelarEditarComentario(comentarioId, event) {
    const editor = event.currentTarget.parentElement.previousElementSibling
    if (!editor || !quillInstanceEditor[comentarioId]) return
    quillInstanceEditor[comentarioId].disable()
    editor.previousElementSibling.remove()
    editor.parentNode.replaceChild(editingBackup[comentarioId], editor)
    delete editingBackup[comentarioId]
    delete quillInstanceEditor[comentarioId]
    editingCommentMap[comentarioId] = false
}

function removerComentario(comentarioId, questionId) {
    tempCommentId.value = comentarioId
    tempQuestionId.value = questionId
    captchaAction.value = 'remove_comment'
    mostrarCaptcha.value = true
}

// ---------------- FETCH COMENTÁRIOS ----------------

async function fetchComentarios() {
    try {
        const response = await fetch(`https://tools.bizout.com.br/cmtsr?op=${props.qstDb}&id=${props.questionId}`)
        const data = await response.json()
        if (data[0] === 200) {
            loading.value = false
            error.value = false
            delete data[0]
            comentarios.value = Object.keys(data).map(key => {
                const cmtUserData = JSON.parse(data[key][2])
                cmtUserData[2] = getFirstNames(cmtUserData[2])
                return [data[key][0], data[key][1], cmtUserData, data[key][3], data[key][4], data[key][5]]
            })
            await nextTick()
            initQuill()
        } else {
            loading.value = false
            error.value = true
        }
    } catch (err) {
        loading.value = false
        error.value = true
        console.error(err)
    }
}

// ---------------- WATCHERS ----------------

watch(() => props.visible, async (isVisible) => {
    if (isVisible) {
        loading.value = true
        error.value = false
        await nextTick()
        setTimeout(fetchComentarios, 1500)
    } else {
        destroyQuill()
        mostrarCaptcha.value = false
        Object.keys(editingCommentMap).forEach(key => editingCommentMap[key] = false)
    }
})

// ---------------- LIFE CYCLE ----------------

onMounted(() => {
    try {
        userId.value = window.userData().id
    } catch (e) {
        userId.value = null
    }
})

onBeforeUnmount(() => {
    destroyQuill()
})
</script>

<template>
    <div v-if="props.visible">

        <!-- CAPTCHA Modal -->
        <div v-if="mostrarCaptcha" class="fixed inset-0 z-50 flex items-center justify-center">
            <div class="bg-white rounded-xl shadow-lg max-w-sm w-full p-6 relative">
                <div class="flex justify-between items-start mb-3 pb-2 border-b border-gray-200">
                    <h3 class="text-sm font-semibold text-gray-700">Confirmação de segurança</h3>
                    <button @click="fecharCaptchaModal"
                        class="absolute top-0 right-0 text-gray-500 hover:text-gray-700 p-2 rounded">
                        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <!-- Overlay de Loading -->
                <div v-if="confirmLoadingCaptcha"
                    class="absolute inset-0 flex items-center justify-center bg-white z-50 rounded-xl">
                    <svg class="h-10 w-10 animate-spin text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none"
                        viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4">
                        </circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
                    </svg>
                </div>

                <!-- Conteúdo CAPTCHA -->
                <div v-else>
                    <div class="bg-gray-100 border border-gray-300 rounded p-4 text-center text-sm text-gray-600 mb-4">
                        [ Captcha Aqui ]
                    </div>
                    <button class="cursor-pointer w-full mt-2 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition"
                        :disabled="confirmLoadingCaptcha" @click="enviarComentarioComCaptcha">
                        Confirmar e Enviar
                    </button>
                </div>
            </div>
        </div>

        <!-- Loading Comentários -->
        <div v-if="loading"
            class="flex items-center gap-4 p-4 bg-blue-100 text-blue-800 rounded-lg border border-blue-300">
            <svg class="h-6 w-6 animate-spin text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none"
                viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
            </svg>
            <div>
                <p class="font-semibold">Carregando comentários...</p>
                <p class="text-sm">Aguarde um momento enquanto buscamos os dados.</p>
            </div>
        </div>

        <div v-else>
            <div v-if="!error" class="mt-0 pt-0 space-y-4">

                <!-- Sem Comentários -->
                <div v-if="comentarios.length === 0"
                    class="text-sm text-gray-500 bg-gray-50 p-4 rounded-md text-center border border-gray-200">
                    Ainda não temos comentários. <span class="text-blue-600 font-medium">Seja o primeiro a
                        comentar.</span>
                </div>

                <!-- Lista Comentários -->
                <div v-else>
                    <div v-for="comentario in comentarios" :key="comentario[0]"
                        class="flex items-start gap-3 p-4 rounded-xl border border-gray-200 bg-white shadow-sm mb-3 hover:bg-gray-50 transition-all duration-200">
                        <!-- <img :src="comentario.avatar || '../src/assets/perfil_icon.png'" alt="Avatar"
                            class="w-10 h-10 rounded-full object-cover" /> -->

                        <div class="flex-1">
                            <div class="flex items-center justify-between text-sm mb-1">
                                <span class="font-medium text-gray-800">{{ comentario[2][2] }}</span>
                                <span class="text-gray-500 text-xs">{{ comentario[2][1] }}</span>
                            </div>

                            <div v-html="deltaToHtml(JSON.parse(comentario[5]))"
                                class="prose prose-sm max-w-none bg-gray-50 rounded-md p-3 text-gray-800"></div>

                            <!-- Botões -->
                            <div v-if="comentario[2][3] == userId && !editingCommentMap[comentario[0]]"
                                class="mt-2 flex gap-2 text-xs text-gray-500">
                                <button @click="editarComentario(comentario[0], $event)" aria-label="Editar comentário"
                                    class="hover:text-blue-600 transition cursor-pointer">Editar</button>
                                <span>|</span>
                                <button @click="removerComentario(comentario[0], comentario[2][0])"
                                    aria-label="Excluir comentário"
                                    class="hover:text-red-600 transition cursor-pointer">Excluir</button>
                            </div>
                            <div v-else-if="comentario[2][3] == userId && editingCommentMap[comentario[0]]"
                                class="mt-2 flex gap-2 text-xs text-gray-500">
                                <button @click="salvarComentario(comentario[0])"
                                    class="hover:text-blue-600 transition cursor-pointer">Salvar</button>
                                <span>|</span>
                                <button @click="cancelarEditarComentario(comentario[0], $event)"
                                    class="hover:text-red-600 transition cursor-pointer">Cancelar</button>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Editor Quill -->
                <div v-if="userId">
                    <div ref="quillContainer" class="bg-white rounded shadow min-h-[120px]"></div>
                    <button class="cursor-pointer mt-3 px-4 py-2 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 transition"
                        @click="publicarComentario">
                        Publicar Comentário
                    </button>
                </div>

                <!-- Mensagem para não logados -->
                <div v-else
                    class="p-4 bg-yellow-100 border border-yellow-300 rounded text-yellow-800 text-sm text-center">
                    Você precisa estar logado para publicar comentários.
                </div>
            </div>

            <!-- Erro Fetch -->
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
        </div>
    </div>
</template>

<style scoped>
.ql-toolbar.ql-snow {
    border-radius: 0.375rem 0.375rem 0 0;
    border-color: #e5e7eb;
}

.ql-container.ql-snow {
    border-radius: 0 0 0.375rem 0.375rem;
    border-color: #e5e7eb;
    font-family: inherit;
    font-size: 0.875rem;
}
</style>
