<template>
    <div class="p-6 max-w-4xl mx-auto pt-20 min-h-[100vh]">
        <!-- Barra superior moderna -->
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
        <!-- Card: Gerenciamento de Dispositivos -->
        <div class="bg-white rounded-2xl shadow p-6">
            <h2 class="text-lg font-medium text-gray-700 mb-4 flex items-center gap-2">
                <!-- Ícone dispositivo -->
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-gray-500" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <rect x="7" y="2" width="10" height="20" rx="2" />
                    <line x1="12" y1="18" x2="12" y2="18" stroke-width="2" stroke-linecap="round" />
                </svg>
                Gerenciamento de Dispositivos
            </h2>

            <p class="text-sm text-gray-500 mb-4">
                Veja todos os dispositivos conectados à sua conta e encerre sessões quando necessário.
            </p>

            <!-- Lista de dispositivos -->
            <div v-if="devices.length > 0" class="space-y-3">
                <div v-for="device in devices" :key="device.session_id"
                    :class="['flex items-center justify-between p-4 rounded-xl transition', device.is_current ? 'bg-green-50 border border-green-300' : 'bg-gray-50 hover:bg-gray-100']">
                    <!-- Info -->
                    <div>
                        <p class="text-sm font-medium text-gray-800">
                            {{ device.device }} <span class="text-gray-400">•</span> {{ device.ip }}
                        </p>
                        <p class="text-xs text-gray-500">
                            Ativo desde {{ formatDate(device.created_at) }}
                        </p>

                        <!-- Sessão atual -->
                        <p v-if="device.is_current" class="text-xs font-medium text-green-600 mt-1">
                            Sessão atual
                        </p>
                    </div>

                    <!-- Botão sair -->
                    <button @click="logoutDevice(device.session_id)"
                        class="cursor-pointer flex items-center gap-1 text-sm px-3 py-1.5 bg-red-500 text-white rounded-lg hover:bg-red-600 transition"
                        :disabled="device.is_current">
                        <!-- Ícone logout -->
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24"
                            stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M15 12H3m12 0l-4-4m4 4l-4 4M21 3v18" />
                        </svg>
                        Encerrar
                    </button>
                </div>
            </div>

            <!-- Nenhum dispositivo -->
            <div v-else class="text-sm text-gray-500 italic">
                Nenhum dispositivo ativo encontrado.
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";

const devices = ref([]);
const router = useRouter()

window.headerVisibleFunc(false)

// Busca dispositivos ativos (simulado por enquanto)
onMounted(async () => {
    window.parent.postMessage({ action: "expandIframe" }, "*")
    try {
        const res = await fetch("https://auth.bizout.com.br/devices", {
            credentials: "include",
        });
        if (res.ok) {
            const data = await res.json();
            devices.value = data.sessions || []; // <-- pega apenas o array de sessões
        }
    } catch (err) {
        console.error("Erro ao buscar dispositivos:", err);
    }
});

function goBack() {
    router.push("/perfil")
}

function formatDate(dateStr) {
    const d = new Date(dateStr);
    return d.toLocaleString("pt-BR", {
        dateStyle: "short",
        timeStyle: "short",
    });
}

async function logoutDevice(sessionId) {
    try {
        await fetch(`https://auth.bizout.com.br/devices/logout/${sessionId}`, {
            method: "POST",
            credentials: "include",
        });
        devices.value = devices.value.filter((d) => d.session_id !== sessionId);
    } catch (err) {
        console.error("Erro ao encerrar sessão:", err);
    }
}
</script>