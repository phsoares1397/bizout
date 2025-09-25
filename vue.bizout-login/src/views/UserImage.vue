<template>
    <main :class="[
        'flex items-center justify-center bg-gray-100 px-4 transition-all',
        noHeader ? 'min-h-[100vh]' : 'min-h-[calc(100vh-58px)]']">

        <div class="w-full max-w-md bg-white rounded-2xl shadow-lg p-8">
            <!-- Título -->
            <div class="flex flex-col items-center mb-6">
                <h1 class="text-xl font-semibold text-gray-800">Perfil</h1>
                <p class="text-sm text-gray-500">Escolha sua imagem de perfil (opcional)</p>
            </div>

            <!-- Avatar com transição fade + scale -->
            <div class="relative flex items-center justify-center mb-6">
                <transition name="fade-scale" mode="out-in">
                    <img v-if="hasImage" :key="'img'" :src="user_image.img" class="h-24 w-24 rounded-full object-cover"
                        alt="Avatar" />
                    <span v-else :key="'initials'"
                        class="h-24 w-24 rounded-full flex items-center justify-center text-white font-semibold text-xl transition-all duration-300"
                        :style="{ backgroundColor: user_image.bgColor }">
                        {{ user_image.initials }}
                    </span>
                </transition>
            </div>

            <!-- Upload -->
            <div class="flex flex-col items-center space-y-2">
                <!-- Botão gerar aleatória -->
                <button @click="handleDefault" class="flex items-center justify-center gap-2 cursor-pointer py-2 px-4 bg-gray-500 text-white font-medium rounded-lg 
                 hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all duration-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                        stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M4 4h16M4 20h16M4 12h16" />
                    </svg>
                    Gerar aleatória
                </button>

                <!-- Input oculto + label para escolher imagem -->
                <input type="file" accept="image/*" @change="handleFileChange" id="avatarInput" class="hidden" />
                <label for="avatarInput" class="flex items-center justify-center gap-2 cursor-pointer py-2 px-4 bg-[#2c89a0] text-white rounded-lg 
                 hover:bg-[#145d6f] transition-all duration-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                        stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2M12 12V4m0 0l-4 4m4-4l4 4" />
                    </svg>
                    Escolher imagem
                </label>
            </div>

            <!-- Botão Salvar -->
            <div class="mt-6 space-y-2">
                <button @click="handleSave" :disabled="!imageChange" :class="[
                    'w-full py-2 px-4 font-medium rounded-lg transition-all duration-300 focus:outline-none focus:ring-2',
                    imageChange
                        ? 'bg-green-600 text-white hover:bg-green-700 focus:ring-green-400 cursor-pointer'
                        : 'bg-gray-200 text-gray-400 cursor-not-allowed'
                ]">
                    Salvar
                </button>
                <button @click="handleCancell" class="cursor-pointer w-full py-2 px-4 border border-gray-400 text-gray-700 font-medium rounded-lg 
                 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-300 transition-all duration-200">
                    Voltar
                </button>
            </div>
        </div>
    </main>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const storedUser = localStorage.getItem("user")
const userObj = storedUser ? JSON.parse(storedUser) : { fname: "" }
const noHeader = ref(false)

// Avatar
const user_image = ref({ img: "", initials: "", bgColor: "" })
const defaultImage = { img: "", initials: "", bgColor: "" }
const imageChange = ref(false)
const hasImage = ref(false)
const file = ref(null)
const avatarURL = ref("")

// Cores aleatórias
const colors = [
    "#ef4444", "#f87171", "#b91c1c", "#dc2626",
    "#f97316", "#fb923c", "#ea580c", "#c2410c",
    "#facc15", "#fde047", "#eab308", "#ca8a04",
    "#22c55e", "#4ade80", "#16a34a", "#15803d",
    "#06b6d4", "#67e8f9", "#0891b2", "#0e7490",
    "#3b82f6", "#60a5fa", "#2563eb", "#1d4ed8",
    "#8b5cf6", "#a78bfa", "#7c3aed", "#6d28d9",
    "#ec4899", "#f472b6", "#db2777", "#be185d",
    "#6b7280", "#9ca3af", "#d1d5db", "#374151"
]

// Computed: iniciais
const initials = computed(() => {
    const fname = userObj.fname || ""
    if (!fname) return ""
    const ignoreWords = ["de", "da", "do", "dos", "das", "e"]
    const parts = fname.split(" ").filter(p => p && !ignoreWords.includes(p.toLowerCase()))
    if (parts.length === 0) return ""
    const first = parts[0][0] || ""
    const second = parts.length > 1 ? parts[1][0] : ""
    return (first + second).toUpperCase()
})

// Watch profundo para alterações de avatar
watch(
    user_image,
    (newVal) => {
        imageChange.value =
            newVal.img !== defaultImage.img ||
            newVal.initials !== defaultImage.initials ||
            newVal.bgColor !== defaultImage.bgColor
    },
    { deep: true }
)

// Upload de arquivo
function handleFileChange(event) {
    const selectedFile = event.target.files[0]
    if (!selectedFile) return
    file.value = selectedFile
    avatarURL.value = URL.createObjectURL(selectedFile)
    user_image.value.img = avatarURL.value
    user_image.value.initials = ""
    user_image.value.bgColor = ""
    hasImage.value = true
}

// Gerar avatar aleatório
function handleDefault() {
    file.value = null
    avatarURL.value = ""
    user_image.value.img = ""
    user_image.value.initials = initials.value
    user_image.value.bgColor = colors[Math.floor(Math.random() * colors.length)]
    hasImage.value = false
}

// Botão voltar
function handleCancell() {
    router.push("/meus-dados")
}

// Salvar avatar
async function handleSave() {
    try {
        const formData = new FormData()
        if (file.value) {
            formData.append("avatar", file.value)
        } else {
            formData.append("avatar", JSON.stringify({ image: false, data: user_image.value.bgColor }))
        }

        const res = await fetch("https://auth.bizout.com.br/userimage", {
            method: "POST",
            body: formData,
            credentials: "include"
        })
        const data = await res.json()
        if (!res.ok || data.status !== "ok") {
            window.mostrarPopup("Perfil", data.message || "Erro ao salvar imagem", "error", 2500)
            return
        }

        imageChange.value = false
        defaultImage.img = user_image.value.img
        defaultImage.initials = user_image.value.initials
        defaultImage.bgColor = user_image.value.bgColor

        router.push("/meus-dados")
    } catch (err) {
        console.error(err)
        window.mostrarPopup("Perfil", "Erro ao conectar com servidor", "error", 2500)
    }
}

// onMounted: carregar avatar do backend
onMounted(async () => {
    if (!document.querySelector("header")) noHeader.value = true
    if (!storedUser) return

    try {
        const res = await fetch("https://auth.bizout.com.br/me?op=full", {
            credentials: "include",
        })
        const data = await res.json()

        let imgObj = null
        try { imgObj = typeof data.img === "string" ? JSON.parse(data.img) : null } catch { }

        if (imgObj?.image && imgObj.data) {
            user_image.value.img = "https://bizout.com.br" + imgObj.data
            hasImage.value = true
        } else {
            user_image.value.initials = initials.value
            user_image.value.bgColor = imgObj?.data || "#60a5fa"
            hasImage.value = false
        }

        defaultImage.img = user_image.value.img
        defaultImage.initials = user_image.value.initials
        defaultImage.bgColor = user_image.value.bgColor

    } catch (err) {
        console.error("Erro ao buscar dados do usuário:", err)
    }
})
</script>

<style>
/* Fade + scale para avatar */
.fade-scale-enter-active,
.fade-scale-leave-active {
    transition: all 0.3s ease;
}

.fade-scale-enter-from,
.fade-scale-leave-to {
    opacity: 0;
    transform: scale(0.8);
}

.fade-scale-enter-to,
.fade-scale-leave-from {
    opacity: 1;
    transform: scale(1);
}
</style>