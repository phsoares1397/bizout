<template>
    <div>
        <!-- Botão do Perfil -->
        <div class="cursor-pointer relative" @click="toggleSidebar">
            <img v-if="hasImage" :src="user_image.img" class="h-10 w-10 rounded-full object-cover" alt="Avatar" />

            <!-- Se não tiver imagem, gera span com iniciais -->
            <span v-else class="h-10 w-10 rounded-full flex items-center justify-center text-white font-normal text-sm"
                :style="{ backgroundColor: user_image.bgColor }">
                {{ user_image.initials }}
            </span>
        </div>

        <!-- Overlay -->
        <transition name="fade">
            <div v-if="open" @click="toggleSidebar" class="fixed inset-0 bg-black/40 z-40"></div>
        </transition>

        <!-- Sidebar com iframe -->
        <transition name="slide">
            <div v-if="open" :class="[
                'fixed top-0 right-0 max-w-full h-full bg-white shadow-xl z-50 flex flex-col transition-all duration-300',
                expanded ? 'w-130' : 'w-78']">
                <!-- Fechar -->
                <div class="absolute top-2 right-1 h-10 flex items-center justify-end pr-4">
                    <button @click="toggleSidebar"
                        class="cursor-pointer text-gray-600 hover:text-black text-2xl">✕</button>
                </div>

                <!-- Conteúdo via iframe -->
                <iframe src="/conta/perfil" class="flex-1 w-full border-0"></iframe>
            </div>
        </transition>
    </div>
</template>

<script setup>
import { ref, watch, onMounted } from "vue"

const open = ref(false)
const hasImage = ref(false)

const expanded = ref(false)

const user_image = ref({
    img: "",        // caminho da imagem se houver
    initials: "",   // iniciais do nome
    bgColor: "",    // cor de fundo
})

watch(open, (newValue, oldValue) => {
    expanded.value = false
})

function updateImage() {
    let data = localStorage.getItem("user")
    if (data) {
        try {
            data = JSON.parse(data)
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
    }
}

onMounted(() => {
    updateImage()
})

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

window.addEventListener("message", (event) => {
    if (event.data.action === "expandIframe") {
        expanded.value = true
    } else if (event.data.action === "contractIframe") {
        expanded.value = false
    }
})

function toggleSidebar() {
    open.value = !open.value
    updateImage()
}
</script>

<style>
/* animações */
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
    transition: transform 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
    transform: translateX(100%);
}
</style>