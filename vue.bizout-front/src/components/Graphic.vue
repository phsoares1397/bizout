<template>
    <div v-if="visivel" class="container">
        <div class="grafico">
            <h2 class="titulo">
                <svg class="titulo-icone" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M3 17v-6a1 1 0 011-1h3v7H4a1 1 0 01-1-1zM10 12v5h3v-5h-3zM17 8v9h3v-9h-3z" />
                </svg>
                Desempenho Hoje
            </h2>
            <canvas v-if="chartDataHoje.labels.length && chartDataHoje.datasets.length" ref="canvasHoje"></canvas>
            <p v-else class="sem-dados">Sem dados para hoje.</p>
        </div>

        <div class="grafico" style="flex: 3">
            <h2 class="titulo">
                <svg class="titulo-icone" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M3 16l4-4 4 4 6-6 4 4" />
                </svg>
                Histórico Geral
            </h2>
            <canvas v-if="chartDataHistorico.labels.length && chartDataHistorico.datasets.length"
                ref="canvasHistorico"></canvas>
            <p v-else class="sem-dados">Sem dados históricos.</p>
        </div>
    </div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick, onBeforeUnmount } from "vue";
import Chart from "chart.js/auto";

const props = defineProps({
    historico: {
        type: Array,
        default: () => [],
    },
    visivel: {
        type: Boolean,
        default: true,
    },
});

const defaultChartData = { labels: [], datasets: [] };
const chartDataHoje = ref({ ...defaultChartData });
const chartDataHistorico = ref({ ...defaultChartData });

const canvasHoje = ref(null);
const canvasHistorico = ref(null);
let chartHojeInstance = null;
let chartHistoricoInstance = null;

const opcoesBasicas = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { position: "top" },
        title: { display: false },
    },
    scales: {
        y: {
            beginAtZero: true,
            precision: 0,
        },
    },
};

const opcoesHistorico = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { position: "top" },
        title: { display: false },
    },
    scales: {
        x: {
            ticks: {
                maxTicksLimit: 7,  // limita a quantidade de labels visíveis
                autoSkip: true,
                maxRotation: 0,
                minRotation: 0,
            },
        },
        y: {
            beginAtZero: true,
            precision: 0,
        },
    },
};

function prepararDados(historico) {
    if (
        !Array.isArray(historico) ||
        historico.length === 0 ||
        !historico.every(
            (item) =>
                item.data &&
                typeof item.resolvidas === "number" &&
                typeof item.acertos === "number" &&
                typeof item.erros === "number"
        )
    ) {
        chartDataHoje.value = { ...defaultChartData };
        chartDataHistorico.value = { ...defaultChartData };
        return;
    }

    const hoje = new Date();

    const ano = hoje.getFullYear();
    const mes = String(hoje.getMonth() + 1).padStart(2, '0');
    const dia = String(hoje.getDate()).padStart(2, '0');

    const hojeStr = `${ano}-${mes}-${dia}`;

    const hojeStats =
        historico.find((item) => item.data === hojeStr) || {
            resolvidas: 0,
            acertos: 0,
            erros: 0,
        };

    chartDataHoje.value = {
        labels: ["Resolvidas", "Acertos", "Erros"],
        datasets: [
            {
                label: "Hoje",
                backgroundColor: ["#3B82F6", "#10B981", "#EF4444"],
                data: [hojeStats.resolvidas, hojeStats.acertos, hojeStats.erros],
            },
        ],
    };

    chartDataHistorico.value = {
        labels: historico.map((item) => item.data),
        datasets: [
            {
                label: "Resolvidas",
                borderColor: "#3B82F6",
                backgroundColor: "#DBEAFE",
                data: historico.map((item) => item.resolvidas),
                fill: false,
                tension: 0.3,
            },
            {
                label: "Acertos",
                borderColor: "#10B981",
                backgroundColor: "#D1FAE5",
                data: historico.map((item) => item.acertos),
                fill: false,
                tension: 0.3,
            },
            {
                label: "Erros",
                borderColor: "#EF4444",
                backgroundColor: "#FECACA",
                data: historico.map((item) => item.erros),
                fill: false,
                tension: 0.3,
            },
        ],
    };
}

function criarOuAtualizarGraficos() {
    if (canvasHoje.value && chartDataHoje.value.labels.length) {
        if (chartHojeInstance) {
            chartHojeInstance.destroy();
            chartHojeInstance = null;
        }
        chartHojeInstance = new Chart(canvasHoje.value, {
            type: "bar",
            data: chartDataHoje.value,
            options: opcoesBasicas,
        });
    }

    if (canvasHistorico.value && chartDataHistorico.value.labels.length) {
        if (chartHistoricoInstance) {
            chartHistoricoInstance.destroy();
            chartHistoricoInstance = null;
        }
        chartHistoricoInstance = new Chart(canvasHistorico.value, {
            type: "line",
            data: chartDataHistorico.value,
            options: opcoesHistorico,
        });
    }
}

onBeforeUnmount(() => {
    if (chartHojeInstance) {
        chartHojeInstance.destroy();
        chartHojeInstance = null;
    }
    if (chartHistoricoInstance) {
        chartHistoricoInstance.destroy();
        chartHistoricoInstance = null;
    }
});

watch(
    () => props.historico,
    async (novoHistorico) => {
        if (!props.visivel) return;
        prepararDados(novoHistorico);
        await nextTick();
        criarOuAtualizarGraficos();
    }
);

watch(
    () => props.visivel,
    async (novoValor) => {
        if (novoValor) {
            prepararDados(props.historico);
            await nextTick();
            criarOuAtualizarGraficos();
        } else {
            if (chartHojeInstance) {
                chartHojeInstance.destroy();
                chartHojeInstance = null;
            }
            if (chartHistoricoInstance) {
                chartHistoricoInstance.destroy();
                chartHistoricoInstance = null;
            }
        }
    }
);

onMounted(async () => {
    if (props.visivel) {
        prepararDados(props.historico);
        await nextTick();
        criarOuAtualizarGraficos();
    }
});
</script>

<style scoped>
.container {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    width: calc(100vw - 2rem);
    max-width: calc(100vw - 2rem) !important;
    padding: 1rem;
    margin-bottom: 2rem;
}

.grafico {
    flex: 1;
}

canvas {
    width: 100%;
    height: 300px;
    max-height: 300px;
}

.titulo {
    display: flex;
    align-items: center;
    font-weight: 600;
    font-size: 1.15rem;
    color: #374151;
    padding-left: 0.75rem;
    margin-bottom: 1rem;
    font-family: 'Inter', sans-serif;
}

.titulo-icone {
    width: 24px;
    height: 24px;
    margin-right: 0.5rem;
    stroke: #f2a81d;
}

.sem-dados {
    color: #6b7280;
    /* cinza */
    font-style: italic;
}

/* A partir de 768px, os gráficos ficam lado a lado */
@media (min-width: 768px) {
    .container {
        flex-direction: row;
    }
}
</style>