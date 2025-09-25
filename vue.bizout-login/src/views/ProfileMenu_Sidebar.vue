<template>
    <div class="fixed top-0 right-0 h-full w-full bg-white shadow-xl z-50 flex flex-col p-0 pt-8">

        <!-- Spinner enquanto carrega -->
        <div v-if="loading" class="flex-1 flex items-center justify-center">
            <svg class="animate-spin h-10 w-10 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none"
                viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
            </svg>
        </div>

        <!-- Se não estiver logado -->
        <div v-else-if="!isLoggedIn" class="p-6 pt-0 flex flex-col space-y-4 flex-1 justify-center items-center">
            <p class="text-gray-600 text-sm">Bem-vindo!</p>
            <button @click="login"
                class="cursor-pointer w-full px-4 py-2 rounded-sm bg-[#2c89a0] text-white hover:bg-indigo-700 transition">
                Fazer Login
            </button>
        </div>

        <!-- Se estiver logado -->
        <div v-else class="p-6 pt-0 pb-0 flex flex-col flex-1 justify-between">
            <div class="space-y-6">
                <!-- User Info -->
                <div class="flex items-center space-x-4">
                    <!-- Se tiver imagem real -->
                    <img v-if="hasImage" :src="user_image.img" class="h-12 w-12 rounded-full object-cover"
                        alt="Avatar" />

                    <!-- Se não tiver imagem, gera span com iniciais -->
                    <span v-else
                        class="h-12 w-12 rounded-full flex items-center justify-center text-white font-semibold text-lg"
                        :style="{ backgroundColor: user_image.bgColor }">
                        {{ user_image.initials }}
                    </span>
                    <div>
                        <p class="font-semibold">{{ user.name }}</p>
                        <p class="text-sm text-gray-500">{{ user.user }}</p>
                    </div>
                </div>

                <!-- Opções -->
                <div class="flex flex-col space-y-1">
                    <!-- Meus dados -->
                    <button @click="meus_dados"
                        class="cursor-pointer flex items-center space-x-2 text-left px-4 py-2 rounded-lg hover:bg-gray-100 text-gray-700 font-medium transition-colors duration-150">
                        <!-- Ícone pequeno à esquerda -->
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-500" fill="none"
                            viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M5.121 17.804A11.959 11.959 0 0112 15c2.5 0 4.79.805 6.879 2.196M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                        <span>Meus dados</span>
                    </button>

                    <!-- Configurações -->
                    <button @click="configsHandler"
                        class="cursor-pointer flex items-center space-x-2 text-left px-4 py-2 rounded-lg hover:bg-gray-100 text-gray-700 font-medium transition-colors duration-150">
                        <!-- Ícone pequeno à esquerda -->
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-500" fill="none"
                            viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 4v1m0 14v1m8-9h1M4 12H3m15.364-6.364l.707.707M6.343 17.657l-.707.707m12.728 0l.707-.707M6.343 6.343l-.707-.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                        </svg>
                        <span>Configurações</span>
                    </button>
                </div>
            </div>

            <!-- Logout -->
            <button @click="logout"
                class="cursor-pointer flex items-center space-x-2 w-full px-4 py-2 rounded-lg hover:bg-gray-100 text-gray-700 font-medium transition-colors duration-150">
                <!-- Ícone de logout -->
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-500" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h6a2 2 0 012 2v1" />
                </svg>
                <span>Sair</span>
            </button>
        </div>

        <!-- Ícones discretos de redes sociais -->
        <div class="my-3 flex justify-center space-x-4">
            <!-- Instagram -->
            <a href="#" target="_blank" class="text-gray-400 hover:text-gray-600 transition">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                    <path
                        d="M7.75 2h8.5A5.75 5.75 0 0122 7.75v8.5A5.75 5.75 0 0116.25 22h-8.5A5.75 5.75 0 012 16.25v-8.5A5.75 5.75 0 017.75 2zm0 1.5A4.25 4.25 0 003.5 7.75v8.5A4.25 4.25 0 007.75 20.5h8.5a4.25 4.25 0 004.25-4.25v-8.5A4.25 4.25 0 0016.25 3.5h-8.5zM12 7a5 5 0 110 10 5 5 0 010-10zm0 1.5a3.5 3.5 0 100 7 3.5 3.5 0 000-7zm4.75-.75a1.25 1.25 0 110 2.5 1.25 1.25 0 010-2.5z" />
                </svg>
            </a>
            <!-- YouTube -->
            <a href="#" target="_blank" class="text-gray-400 hover:text-gray-600 transition">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                    <path
                        d="M23.498 6.186a2.998 2.998 0 00-2.111-2.115C19.125 3.5 12 3.5 12 3.5s-7.125 0-9.387.571a2.998 2.998 0 00-2.111 2.115A31.643 31.643 0 000 12a31.643 31.643 0 00.502 5.814 2.998 2.998 0 002.111 2.115C4.875 20.5 12 20.5 12 20.5s7.125 0 9.387-.571a2.998 2.998 0 002.111-2.115A31.643 31.643 0 0024 12a31.643 31.643 0 00-.502-5.814zM9.75 15.02V8.98l6.25 3.02-6.25 3.02z" />
                </svg>
            </a>
            <!-- Twitter -->
            <a href="#" target="_blank" class="text-gray-400 hover:text-gray-600 transition">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                    <path
                        d="M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z" />
                </svg>
            </a>
        </div>

    </div>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"

const isLoggedIn = ref(false)
const loading = ref(true)
const user = ref({ name: "", email: "", avatar: "" })
const router = useRouter()

const user_image = ref({
    img: "",        // caminho da imagem se houver
    initials: "",   // iniciais do nome
    bgColor: "",    // cor de fundo
})

const hasImage = ref(false)

window.headerVisibleFunc(false)

function configsHandler() {
    router.push("/config")
}

function login() {
    window.top.location = "/conta/login"
}

function meus_dados() {
    router.push("/meus-dados")
}

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

function getFirstNames(fullName) {
    if (!fullName) return ""

    const parts = fullName.trim().split(/\s+/)

    if (parts.length === 1) return parts[0]

    // Se a segunda palavra tiver 2 caracteres ou menos, considera como conector
    if (parts.length >= 3 && parts[1].length <= 2) {
        return `${parts[0]} ${parts[1]} ${parts[2]}`
    }

    return `${parts[0]} ${parts[1]}`
}

async function logout() {
    loading.value = true
    try {
        const res = await fetch("https://auth.bizout.com.br/logout", {
            method: "POST",
            credentials: "include"
        })

        if (res.ok) {
            localStorage.removeItem("user")
            setTimeout(() => {
                loading.value = false
                isLoggedIn.value = false
                window.top.reloadComponents()
            }, 1500)
        } else {
            const data = await res.json()
            window.showAlert(data.message || "Erro ao fazer logout", "error")
            loading.value = false
        }
    } catch (err) {
        console.error(err)
        window.showAlert("Erro de conexão com o servidor", "error")
        loading.value = false
    }
}

function isInIframe() {
    try {
        return window.self !== window.top
    } catch (e) {
        // Caso de cross-origin (iframe de outro domínio)
        return true
    }
}

onMounted(async () => {
    window.parent.postMessage({ action: "contractIframe" }, "*")
    loading.value = true
    if (!isInIframe()) {
        window.location.href = "https://bizout.com.br/questoes"
        return
    }
    const data = await window.checkSession()
    localStorage.setItem("user", JSON.stringify(data))
    if (data) {
        isLoggedIn.value = true
        user.value = {
            name: getFirstNames(data.fname),
            user: data.user || "",
        }
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
    } else {
        isLoggedIn.value = false
    }
    setTimeout(() => {
        loading.value = false
    }, 1300)
})
</script>