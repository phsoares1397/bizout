<template>
    <main class="flex items-center justify-center min-h-[calc(100vh-58px)] bg-gray-100 px-4">
        <div class="w-full max-w-md bg-white rounded-2xl shadow-lg p-8">
            <!-- Logo / Título -->
            <div class="flex flex-col items-center mb-6">
                <h1 class="text-xl font-semibold text-gray-800">Recuperar Senha</h1>
            </div>

            <!-- Formulário -->
            <form @submit.prevent="handleSubmit" class="space-y-4">
                <div>
                    <label for="email" class="block text-sm font-medium text-gray-700">Email ou Usuário</label>
                    <input v-model="email" type="text" id="email" placeholder="Digite seu email ou usuário"
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800"
                        required />
                </div>

                <button type="submit"
                    class="cursor-pointer w-full py-2 px-4 bg-[#2c89a0] text-white font-medium rounded-lg hover:bg-[#145d6f] focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all">
                    Enviar link
                </button>
            </form>

            <div class="mt-4 text-center">
                <button @click="goToLogin" class="cursor-pointer text-sm text-[#145d6f] hover:underline">
                    Voltar ao login
                </button>
            </div>
        </div>
    </main>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const email = ref("")

document.title = "Recuperar senha - Bizout";

function handleSubmit() {
    if (!email.value) {
        window.mostrarPopup("Recuperar Senha", "Digite seu email ou usuário", "warning", 2500)
        return
    }

    // Aqui você chamaria a API de envio de link
    window.mostrarPopup("Recuperar Senha", "Link de redefinição enviado com sucesso!", "success", 3000)
    email.value = ""
}

function goToLogin() {
    router.push("/login")
}

onMounted(async () => {
    window.parent.postMessage({ action: "contractIframe" }, "*")
})
</script>