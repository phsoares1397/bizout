<template>
    <main class="flex items-center justify-center min-h-screen bg-gray-100 px-4 py-10">
        <div class="w-full max-w-md bg-white rounded-2xl shadow-lg p-8">
            <div class="flex flex-col items-center mb-6">
                <h1 class="text-xl font-semibold text-gray-800">Criar Conta</h1>
                <p class="text-sm text-gray-500">Preencha os dados para se cadastrar</p>
            </div>

            <form @submit.prevent="handleRegister" class="space-y-4">
                <!-- Nome -->
                <div>
                    <label for="name" class="block text-sm font-medium text-gray-700">Nome</label>
                    <input v-model="form.name" type="text" id="name"
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800"
                        required />
                </div>

                <!-- Usuário -->
                <div>
                    <label for="username" class="block text-sm font-medium text-gray-700">Usuário</label>
                    <input v-model="form.username" type="text" id="username"
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800"
                        required />
                </div>

                <!-- Email -->
                <div>
                    <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
                    <input v-model="form.email" type="email" id="email"
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800"
                        required />
                </div>

                <!-- Senha -->
                <div>
                    <label for="password" class="block text-sm font-medium text-gray-700">Senha</label>
                    <input v-model="form.password" type="password" id="password"
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800"
                        required />
                </div>

                <!-- Confirmar senha -->
                <div>
                    <label for="confirmPassword" class="block text-sm font-medium text-gray-700">Confirmar Senha</label>
                    <input v-model="form.confirmPassword" type="password" id="confirmPassword"
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800"
                        required />
                </div>

                <!-- Estado -->
                <div>
                    <label for="state" class="block text-sm font-medium text-gray-700">Estado</label>
                    <select v-model="form.state" id="state" @change="fetchCities" required
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg bg-white focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800">
                        <option value="">Selecione o estado</option>
                        <option v-for="estado in estadosOrdenados" :key="estado.id" :value="estado.sigla">
                            {{ estado.nome }}
                        </option>
                    </select>
                </div>

                <!-- Cidade -->
                <div v-if="form.state">
                    <label for="city" class="block text-sm font-medium text-gray-700">Cidade</label>
                    <div class="relative">
                        <select v-model="form.city" id="city" required
                            class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg bg-white focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800"
                            :disabled="loadingCities">
                            <option value="">Selecione a cidade</option>
                            <option v-for="cidade in cidadesOrdenadas" :key="cidade.id" :value="cidade.nome">
                                {{ cidade.nome }}
                            </option>
                        </select>

                        <!-- Spinner -->
                        <div v-if="loadingCities" class="absolute inset-y-0 right-3 flex items-center">
                            <svg class="animate-spin h-5 w-5 text-indigo-600" xmlns="http://www.w3.org/2000/svg"
                                fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor"
                                    stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor"
                                    d="M4 12a8 8 0 018-8v4l3-3-3-3v4a8 8 0 00-8 8h4z"></path>
                            </svg>
                        </div>
                    </div>
                </div>

                <!-- Data de nascimento -->
                <div>
                    <label for="birth" class="block text-sm font-medium text-gray-700">Data de Nascimento
                        (opcional)</label>
                    <input v-model="form.birth" type="date" id="birth"
                        class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-gray-800" />
                </div>

                <!-- Botão -->
                <button type="submit"
                    class="cursor-pointer w-full py-2 px-4 bg-[#2c89a0] text-white font-medium rounded-lg hover:bg-[#145d6f] focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all">
                    Cadastrar
                </button>
            </form>

            <div class="mt-4 text-center">
                <button @click="goToLogin" class="cursor-pointer text-sm text-[#145d6f] hover:underline">
                    Já tem conta? Faça login
                </button>
            </div>
        </div>
    </main>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

// Lista fixa de estados
const estados = [
    { id: 11, sigla: "RO", nome: "Rondônia" },
    { id: 12, sigla: "AC", nome: "Acre" },
    { id: 13, sigla: "AM", nome: "Amazonas" },
    { id: 14, sigla: "RR", nome: "Roraima" },
    { id: 15, sigla: "PA", nome: "Pará" },
    { id: 16, sigla: "AP", nome: "Amapá" },
    { id: 17, sigla: "TO", nome: "Tocantins" },
    { id: 21, sigla: "MA", nome: "Maranhão" },
    { id: 22, sigla: "PI", nome: "Piauí" },
    { id: 23, sigla: "CE", nome: "Ceará" },
    { id: 24, sigla: "RN", nome: "Rio Grande do Norte" },
    { id: 25, sigla: "PB", nome: "Paraíba" },
    { id: 26, sigla: "PE", nome: "Pernambuco" },
    { id: 27, sigla: "AL", nome: "Alagoas" },
    { id: 28, sigla: "SE", nome: "Sergipe" },
    { id: 29, sigla: "BA", nome: "Bahia" },
    { id: 31, sigla: "MG", nome: "Minas Gerais" },
    { id: 32, sigla: "ES", nome: "Espírito Santo" },
    { id: 33, sigla: "RJ", nome: "Rio de Janeiro" },
    { id: 35, sigla: "SP", nome: "São Paulo" },
    { id: 41, sigla: "PR", nome: "Paraná" },
    { id: 42, sigla: "SC", nome: "Santa Catarina" },
    { id: 43, sigla: "RS", nome: "Rio Grande do Sul" },
    { id: 50, sigla: "MS", nome: "Mato Grosso do Sul" },
    { id: 51, sigla: "MT", nome: "Mato Grosso" },
    { id: 52, sigla: "GO", nome: "Goiás" },
    { id: 53, sigla: "DF", nome: "Distrito Federal" }
]

// Estados em ordem alfabética
const estadosOrdenados = computed(() =>
    [...estados].sort((a, b) => a.nome.localeCompare(b.nome))
)

const cidades = ref([])
const loadingCities = ref(false)

const cidadesOrdenadas = computed(() =>
    [...cidades.value].sort((a, b) => a.nome.localeCompare(b.nome))
)

const form = reactive({
    name: "",
    username: "",
    email: "",
    password: "",
    confirmPassword: "",
    state: "",
    city: "",
    birth: ""
})

document.title = "Criar conta - Bizout"

async function fetchCities() {
    if (!form.state) {
        cidades.value = []
        return
    }

    try {
        loadingCities.value = true
        cidades.value = []

        const estadoSelecionado = estados.find(e => e.sigla === form.state)
        const res = await fetch(`https://servicodados.ibge.gov.br/api/v1/localidades/estados/${estadoSelecionado.id}/municipios`)
        cidades.value = await res.json()
    } catch (err) {
        console.error("Erro ao carregar cidades:", err)
        cidades.value = []
    } finally {
        loadingCities.value = false
    }
}

async function handleRegister() {
    // Validação de campos obrigatórios
    if (!form.name || !form.username || !form.email || !form.password || !form.state || !form.city) {
        window.mostrarPopup("Cadastro", "Preencha todos os campos obrigatórios", "warning", 2500)
        return
    }

    if (form.password !== form.confirmPassword) {
        window.mostrarPopup("Cadastro", "As senhas não coincidem", "error", 2500)
        return
    }

    const cidadeEstado = `${form.city} - ${form.state}`

    const payload = {
        user: form.username,
        mail: form.email,
        pass: form.password,
        fname: form.name,
        cies: cidadeEstado,
        dnas: form.birth
    }

    try {
        // 1️⃣ Cadastro
        const res = await fetch("https://auth.bizout.com.br/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        })

        const data = await res.json()

        if (!res.ok) {
            window.mostrarPopup("Cadastro", data.message || "Erro ao criar conta", "error", 3000)
            return
        }

        window.showAlert("Cadastro realizado com sucesso!", "success")

        // 2️⃣ Login automático
        const loginRes = await fetch("https://auth.bizout.com.br/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include", // importante para receber session_id
            body: JSON.stringify({
                user: form.username,
                pass: form.password
            })
        })

        const loginData = await loginRes.json()
        if (!loginRes.ok || loginData.status !== "ok") {
            window.mostrarPopup("Login automático falhou", loginData.message || "Erro no login", "error", 3000)
            return
        } else {
            // Busca informações do usuário logado
            const meRes = await fetch("https://auth.bizout.com.br/me", {
                credentials: "include"
            })
            const meData = await meRes.json()

            // Exemplo: salvar no localStorage ou store global
            localStorage.setItem("user", JSON.stringify(meData))

            // 3️⃣ Redireciona para a tela de upload de imagem
            router.push("/image")
        }

    } catch (err) {
        console.error(err)
        window.mostrarPopup("Cadastro", "Erro de conexão com o servidor", "error", 3000)
    }
}

function goToLogin() {
    router.push("/login")
}

onMounted(async () => {
    window.parent.postMessage({ action: "contractIframe" }, "*")
    if (window.checkSession) {
        const user = await window.checkSession()
        if (user) {
            window.mostrarPopup("Cadastro", "Para realizar um novo cadastro, primeiro faça logout da sua conta atual.", "error", 0)
            setTimeout(() => {
                window.location.href = "https://bizout.com.br/questoes" // redireciona se já estiver logado
            }, 3000)
            return
        }
    }
})

</script>
