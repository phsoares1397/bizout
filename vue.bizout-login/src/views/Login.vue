<template>
    <main class="flex items-center justify-center min-h-[calc(100vh-58px)] bg-gray-100 px-4 py-10">
        <div v-if="loading" class="flex-1 flex items-center justify-center">
            <svg class="animate-spin h-10 w-10 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none"
                viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
            </svg>
        </div>

        <div v-else class="w-full max-w-md bg-white rounded-2xl shadow-lg p-8 relative">
            <form @submit.prevent="handleLogin" class="space-y-4">
                <div class="flex flex-col items-center mb-6">
                    <h1 class="text-xl font-semibold text-gray-800">Login</h1>
                    <p class="text-sm text-gray-500">Entre com seu usuário ou email</p>
                </div>
                <div>
                    <label for="email" class="block text-sm font-medium text-gray-700">Usuário ou Email</label>
                    <input v-model="email" type="text" id="email"
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800"
                        required />
                </div>

                <div>
                    <label for="password" class="block text-sm font-medium text-gray-700">Senha</label>
                    <input v-model="password" type="password" id="password"
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800"
                        required />
                </div>

                <button type="submit"
                    class="cursor-pointer w-full py-2 px-4 bg-[#2c89a0] text-white font-medium rounded-lg hover:bg-[#145d6f] focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all">
                    Entrar
                </button>
            </form>

            <div class="mt-4">
                <button @click="handleRegister"
                    class="cursor-pointer w-full py-2 px-4 border border-[#2c89a0] text-[#145d6f] font-medium rounded-lg hover:bg-indigo-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all">
                    Criar conta
                </button>
            </div>

            <div class="mt-4 text-center">
                <button @click="router.push('/forgot-password')"
                    class="cursor-pointer text-sm text-[#145d6f] hover:underline">
                    Esqueci minha senha
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
const password = ref("")
const loading = ref(true) // controla o spinner

document.title = "Login - Bizout";

onMounted(async () => {
    window.parent.postMessage({ action: "contractIframe" }, "*")
    if (window.checkSession) {
        const user = await window.checkSession()
        if (user) {
            setTimeout(() => {
                window.location.href = "https://bizout.com.br/questoes" // redireciona se já estiver logado
            }, 1300)
            return
        }
    }
    setTimeout(() => {
        loading.value = false // se não estiver logado, mostra login
    }, 1300)
})

async function handleLogin() {
    loading.value = true
    try {
        const res = await fetch("https://auth.bizout.com.br/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ user: email.value, pass: password.value })
        })

        const data = await res.json()

        if (!res.ok || data.status !== "ok") {
            setTimeout(() => {
                loading.value = false
            }, 1500)
            window.showAlert(data.message || "Usuário ou senha incorretos", "error")
            return
        }

        window.showAlert("Login realizado com sucesso!", "success")

        // Busca informações do usuário logado
        const meRes = await fetch("https://auth.bizout.com.br/me", { credentials: "include" })
        const meData = await meRes.json()
        localStorage.setItem("user", JSON.stringify(meData))

        window.location.href = "https://bizout.com.br/questoes"
    } catch (err) {
        loading.value = false
        console.error("Erro no login:", err)
        window.showAlert("Erro ao conectar com servidor", "error")
    }
}

function handleRegister() {
    router.push("/register")
}
</script>