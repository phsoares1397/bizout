<template>
    <div class="ml-2">
        <div class="flex items-center justify-between hover:bg-gray-50 rounded px-2 py-1 group">
            <div class="flex items-center space-x-2">
                <button v-if="item.filhos?.length" @click.stop="expandido = !expandido"
                    class="w-4 h-4 text-gray-500 hover:text-gray-700 transition">
                    <ChevronDownIcon :class="[
                        'w-4 h-4 transform transition-transform duration-200',
                        { '-rotate-90': !expandido },
                    ]" />
                </button>
                <span v-else class="w-4 h-4 inline-block" />

                <input ref="checkboxRef" type="checkbox" :checked="estaMarcado" :indeterminate="indeterminate"
                    @change="alternarSelecao" class="accent-blue-600 w-4 h-4" />

                <span class="text-sm text-gray-800" v-html="realceTexto(item.nome)" />
            </div>
        </div>

        <div v-if="expandido" class="ml-4 pl-2">
            <TreeItem v-for="filho in item.filhos" :key="filho.id" :item="filho" :busca="busca"
                v-model:selecionados="selecionados" />
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from "vue";
import { ChevronDown as ChevronDownIcon } from "lucide-vue-next";

const props = defineProps({
    item: Object,
    busca: String,
});

const selecionados = defineModel("selecionados");
const expandido = ref(false);
const checkboxRef = ref(null);

// Verifica se todos os descendentes estão selecionados
function todosDescendentesSelecionados(no) {
    if (!no.filhos || no.filhos.length === 0) {
        return selecionados.value.includes(no.id);
    }
    return no.filhos.every((filho) => todosDescendentesSelecionados(filho));
}

// Verifica se algum descendente está selecionado (mas não todos)
function algumDescendenteSelecionado(no) {
    if (!no.filhos || no.filhos.length === 0) {
        return selecionados.value.includes(no.id);
    }
    return no.filhos.some((filho) => {
        return (
            selecionados.value.includes(filho.id) || algumDescendenteSelecionado(filho)
        );
    });
}

// Atualiza a seleção recursivamente na árvore para os pais
function atualizarSelecaoRecursiva(no) {
    if (!no.filhos || no.filhos.length === 0) return;

    no.filhos.forEach((filho) => atualizarSelecaoRecursiva(filho));

    if (todosDescendentesSelecionados(no)) {
        if (!selecionados.value.includes(no.id)) {
            selecionados.value.push(no.id);
        }
    } else {
        const index = selecionados.value.indexOf(no.id);
        if (index !== -1) {
            selecionados.value.splice(index, 1);
        }
    }
}

const estaMarcado = computed(() => {
    return selecionados.value.includes(props.item.id);
});

const indeterminate = computed(() => {
    // Se não tem filhos, não é indeterminate
    if (!props.item.filhos || props.item.filhos.length === 0) return false;

    // Se tem filhos, indeterminate é true se algum descendente está selecionado, mas nem todos
    const algum = algumDescendenteSelecionado(props.item);
    const todos = todosDescendentesSelecionados(props.item);
    return algum && !todos;
});

function alternarSelecao() {
    const itemId = props.item.id;

    if (selecionados.value.includes(itemId)) {
        // Remover item e todos os seus descendentes
        removerSelecionadosRecursivo(props.item);
    } else {
        // Adicionar item e todos os seus descendentes
        adicionarSelecionadosRecursivo(props.item);
    }
}

// Adiciona o nó e todos os descendentes na seleção
function adicionarSelecionadosRecursivo(no) {
    if (!selecionados.value.includes(no.id)) {
        selecionados.value.push(no.id);
    }
    if (no.filhos && no.filhos.length) {
        no.filhos.forEach((filho) => adicionarSelecionadosRecursivo(filho));
    }
}

// Remove o nó e todos os descendentes da seleção
function removerSelecionadosRecursivo(no) {
    const index = selecionados.value.indexOf(no.id);
    if (index !== -1) {
        selecionados.value.splice(index, 1);
    }
    if (no.filhos && no.filhos.length) {
        no.filhos.forEach((filho) => removerSelecionadosRecursivo(filho));
    }
}

// Atualiza seleção sempre que mudarem os selecionados
watch(
    selecionados,
    () => {
        atualizarSelecaoRecursiva(props.item);
        nextTick(() => {
            if (checkboxRef.value) {
                checkboxRef.value.indeterminate = indeterminate.value;
            }
        });
    },
    { deep: true, immediate: true }
);

onMounted(() => {
    if (checkboxRef.value) {
        checkboxRef.value.indeterminate = indeterminate.value;
    }
});

function realceTexto(texto) {
    if (!props.busca) return texto;
    const re = new RegExp(`(${props.busca})`, "gi");
    expandido.value = true;
    return texto.replace(
        re,
        `<mark class="bg-yellow-100 text-yellow-800 font-medium">$1</mark>`
    );
}
</script>