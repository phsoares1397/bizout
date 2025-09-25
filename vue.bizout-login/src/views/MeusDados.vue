<template>
    <main class="flex items-center justify-center min-h-screen p-6 pt-20">
        <div class="fixed top-0 left-0 w-full z-50 bg-white shadow-md flex items-center justify-between px-4 h-[57px]">
            <!-- Botão voltar -->
            <button @click="goBack"
                class="cursor-pointer flex items-center gap-2 px-3 py-2 rounded-lg bg-gray-100 hover:bg-gray-200 text-gray-700 transition-colors focus:outline-none focus:ring-2 focus:ring-gray-300">
                <!-- Ícone de voltar -->
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
                <span class="font-medium">Voltar</span>
            </button>
        </div>
        <div class="w-full max-w-2xl bg-white p-8 space-y-8 rounded-2xl shadow p-6">
            <!-- Cabeçalho -->
            <div class="flex flex-col items-center text-center space-y-2">
                <h1 class="text-2xl font-bold text-gray-800">Minha Conta</h1>
                <p class="text-sm text-gray-500">Gerencie seus dados pessoais e credenciais de acesso</p>
            </div>

            <!-- Seção de Perfil -->
            <div class="space-y-6">
                <div class="flex flex-col items-center space-y-4">
                    <!-- Foto de perfil -->
                    <div class="relative">
                        <img v-if="hasImage" :src="user_image.img" class="h-24 w-24 rounded-full object-cover"
                            alt="Avatar" />

                        <!-- Se não tiver imagem, gera span com iniciais -->
                        <span v-else
                            class="h-24 w-24 rounded-full flex items-center justify-center text-white font-semibold text-xl"
                            :style="{ backgroundColor: user_image.bgColor }">
                            {{ user_image.initials }}
                        </span>
                    </div>

                    <button @click="goToImage"
                        class="cursor-pointer px-4 py-2 text-sm font-medium text-white bg-[#2c89a0] rounded-lg hover:bg-[#145d6f] transition-all">
                        Alterar Foto
                    </button>
                </div>

                <!-- Dados -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700">Nome</label>
                        <input :value="user.fname" disabled
                            class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100 text-gray-500" />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700">Usuário</label>
                        <input :value="user.user" disabled
                            class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100 text-gray-500" />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700">Email</label>
                        <input :value="user.mail" disabled
                            class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100 text-gray-500" />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700">Cidade / Estado</label>
                        <input :value="user.cies" disabled
                            class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100 text-gray-500" />
                    </div>
                    <div class="md:col-span-2">
                        <label class="block text-sm font-medium text-gray-700">Nascimento</label>
                        <input :value="user.dnas || 'Não informado'" disabled type="date"
                            class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100 text-gray-500" />
                    </div>
                </div>
            </div>

            <!-- Alterar senha -->
            <div class="pt-6 border-t border-gray-200">
                <h2 class="text-lg font-semibold text-gray-800 mb-4">Alterar Senha</h2>
                <form @submit.prevent="handleChangePassword" class="space-y-4">
                    <div>
                        <label for="currentPassword" class="block text-sm font-medium text-gray-700">Senha atual</label>
                        <input v-model="passwordForm.current" type="password" id="currentPassword"
                            class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#2c89a0] focus:border-[#2c89a0] text-gray-800"
                            required />
                    </div>
                    <div>
                        <label for="newPassword" class="block text-sm font-medium text-gray-700">Nova senha</label>
                        <input v-model="passwordForm.new" type="password" id="newPassword"
                            class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#2c89a0] focus:border-[#2c89a0] text-gray-800"
                            required />
                    </div>
                    <div>
                        <label for="confirmPassword" class="block text-sm font-medium text-gray-700">Confirmar nova
                            senha</label>
                        <input v-model="passwordForm.confirm" type="password" id="confirmPassword"
                            class="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#2c89a0] focus:border-[#2c89a0] text-gray-800"
                            required />
                    </div>

                    <button type="submit"
                        class="cursor-pointer w-full py-2 px-4 bg-[#2c89a0] text-white font-medium rounded-lg hover:bg-[#145d6f] focus:outline-none focus:ring-2 focus:ring-[#2c89a0] transition-all">
                        Alterar Senha
                    </button>
                </form>
            </div>
        </div>
    </main>
</template>

<script setup>
import { reactive, ref, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const user = ref({})

window.headerVisibleFunc(false)

// formulário de alteração de senha
const passwordForm = reactive({
    current: "",
    new: "",
    confirm: ""
})

const user_image = ref({
    img: "",        // caminho da imagem se houver
    initials: "",   // iniciais do nome
    bgColor: "",    // cor de fundo
})

const hasImage = ref(false)

onMounted(async () => {
    window.parent.postMessage({ action: "expandIframe" }, "*")
    try {
        const res = await fetch("https://auth.bizout.com.br/me?op=full", {
            credentials: "include",
        })
        const data = await res.json()
        localStorage.setItem("user", JSON.stringify(data))
        user.value = data
        try {
            const imgObj = typeof data.img === "string" ? JSON.parse(data.img) : null
            if (imgObj && imgObj.image) {
                user_image.value.img = "https://bizout.com.br" + imgObj.data
                hasImage.value = true
            } else {
                user_image.value.initials = getInitials(data.fname)
                user_image.value.bgColor = imgObj?.data || "#60a5fa" // cor salva ou fallback
                hasImage.value = false
            }
        } catch {
            user_image.value.initials = getInitials(data.fname)
            user_image.value.bgColor = "#60a5fa"
            hasImage.value = false
        }
    } catch (err) {
        console.error("Erro ao buscar dados do usuário:", err)
    }
})

function goBack() {
    router.push("/perfil")
}

async function handleChangePassword() {
    if (!passwordForm.current || !passwordForm.new || !passwordForm.confirm) {
        window.mostrarPopup("Perfil", "Preencha todos os campos de senha", "warning", 2500)
        return
    }

    if (passwordForm.new !== passwordForm.confirm) {
        window.mostrarPopup("Perfil", "As senhas não coincidem", "error", 2500)
        return
    }

    try {
        const res = await fetch("https://auth.bizout.com.br/changepassword", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({
                current: passwordForm.current,
                new: passwordForm.new
            })
        })

        const data = await res.json()

        if (!res.ok) {
            window.mostrarPopup("Perfil", data.message || "Erro ao alterar senha", "error", 3000)
            return
        }

        window.showAlert("Senha alterada com sucesso!", "success")
        passwordForm.current = ""
        passwordForm.new = ""
        passwordForm.confirm = ""
    } catch (err) {
        console.error("Erro ao alterar senha:", err)
        window.mostrarPopup("Perfil", "Erro de conexão com o servidor", "error", 3000)
    }
}

function goToImage() {
    router.push("/image")
}

// utilitário: gera iniciais
function getInitials(name) {
    if (!name) return ""

    const ignoreWords = ["de", "da", "do", "dos", "das", "e"]
    // Divide o nome em partes e filtra palavras irrelevantes
    const parts = name.trim().split(/\s+/).filter(p => !ignoreWords.includes(p.toLowerCase()))

    if (parts.length === 0) return ""

    const firstInitial = parts[0][0] || ""
    const secondInitial = parts.length > 1 ? parts[1][0] : ""

    return (firstInitial + secondInitial).toUpperCase()
}

// utilitário: gera cor a partir do nome
function stringToColor(str) {
    let hash = 0
    for (let i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash)
    }
    const color = `hsl(${hash % 360}, 60%, 50%)`
    return color
}
</script>