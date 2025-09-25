<template>
    <div class="px-3 py-1 mt-[10px]">
        <!-- Grid 2 linhas de filtros -->
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-2 mb-4">
            <button @click="abrirDisciplina()" class="filtro text-left">
                <span class="block truncate">
                    {{ filtros.disciplina.length ? `Disciplinas (${filtros.disciplina.length})` : 'Disciplina' }}
                </span>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="h-5 w-5 text-gray-400">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            <button @click="abrirAssunto()" class="filtro text-left">
                <span class="block truncate">
                    {{ filtros.assuntos.length ? `Assuntos (${filtros.assuntos.length})` : 'Assuntos' }}
                </span>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="h-5 w-5 text-gray-400">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            <button @click="abrirOrgaosModal = true" class="filtro text-left">
                <span class="block truncate">
                    {{ filtros.orgao.length ? `Orgãos (${filtros.orgao.length})` : 'Orgãos' }}
                </span>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="h-5 w-5 text-gray-400">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            <button @click="abrirCargosModal = true" class="filtro text-left">
                <span class="block truncate">
                    {{ filtros.cargo.length ? `Cargos (${filtros.cargo.length})` : 'Cargos' }}
                </span>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="h-5 w-5 text-gray-400">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            <button @click="abrirBancasModal = true" class="filtro text-left">
                <span class="block truncate">
                    {{ filtros.banca.length ? `Bancas (${filtros.banca.length})` : 'Bancas' }}
                </span>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="h-5 w-5 text-gray-400">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            <button @click="abrirAnoModal = true" class="filtro text-left">
                <span class="block truncate">
                    {{ filtros.ano.length == 15 ? 'Anos' : `Anos (${filtros.ano.length})` }}
                </span>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="h-5 w-5 text-gray-400">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            <!-- Modal de seleção de anos com botões tipo "pill" -->
            <div v-if="abrirAnoModal" class="fixed inset-0 z-50 flex items-center justify-center">
                <!-- Fundo escuro -->
                <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="abrirAnoModal = false"></div>

                <!-- Conteúdo -->
                <div
                    class="relative z-10 bg-white flex flex-col rounded-2xl shadow-lg w-full max-w-md min-h-auto max-h-[75vh] p-6">
                    <h2 class="text-lg font-medium text-gray-800 mb-4">Selecionar Anos</h2>

                    <div class="flex flex-wrap gap-2 overflow-y-auto mb-4">
                        <button v-for="ano in anos" :key="ano" @click="alternarAno(ano)"
                            class="px-3 py-1 text-sm rounded-full transition-colors flex items-center gap-2" :class="{
                                'bg-blue-100 text-blue-700': filtros.ano.includes(ano),
                                'bg-gray-100 text-gray-700 hover:bg-blue-50': !filtros.ano.includes(ano)
                            }">
                            {{ ano }}
                            <span v-if="filtros.ano.includes(ano)"
                                class="text-blue-700 text-xs font-bold hover:text-red-600"
                                @click.stop="removerAno(ano)">
                                ×
                            </span>
                        </button>
                    </div>

                    <!-- Botões de ação -->
                    <div class="flex justify-start items-center mt-auto pt-0 space-x-3">
                        <button @click="selecionarTodosAnos"
                            class="text-sm text-blue-700 hover:text-blue-900 transition-colors flex items-center gap-1">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                            </svg>
                            Selecionar todos
                        </button>

                        <button @click="limparAnos"
                            class="text-sm text-red-600 hover:text-red-800 transition-colors flex items-center gap-1">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                            Limpar
                        </button>

                        <button @click="abrirAnoModal = false"
                            class="text-sm ml-auto px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 transition">
                            Fechar
                        </button>
                    </div>
                </div>
            </div>

            <button @click="abrirNivelModal = true" class="filtro text-left">
                <span class="block truncate">
                    {{ filtros.nivel.length == 3 ? 'Nível' : `Nível (${filtros.nivel.length})` }}
                </span>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="h-5 w-5 text-gray-400">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            <div v-if="abrirNivelModal" class="fixed inset-0 z-50 flex items-center justify-center">
                <!-- Fundo escuro -->
                <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="abrirNivelModal = false"></div>

                <!-- Conteúdo -->
                <div
                    class="relative z-10 bg-white flex flex-col rounded-2xl shadow-lg w-full max-w-md max-h-[75vh] p-6">
                    <h2 class="text-lg font-medium text-gray-800 mb-4">Selecionar Nível</h2>

                    <div class="flex flex-wrap gap-2 overflow-y-auto mb-4">
                        <button v-for="nivel in niveis" :key="nivel" @click="alternarNivel(nivel)"
                            class="px-3 py-1 text-sm rounded-full transition-colors flex items-center gap-2" :class="{
                                'bg-blue-100 text-blue-700': filtros.nivel.includes(nivel),
                                'bg-gray-100 text-gray-700 hover:bg-blue-50': !filtros.nivel.includes(nivel)
                            }">
                            {{ nivel }}
                            <span v-if="filtros.nivel.includes(nivel)"
                                class="text-blue-700 text-xs font-bold hover:text-red-600"
                                @click.stop="removerNivel(nivel)">
                                ×
                            </span>
                        </button>
                    </div>

                    <!-- Botões de ação -->
                    <div class="flex justify-start items-center mt-auto pt-0 space-x-3">
                        <button @click="selecionarTodosNiveis"
                            class="text-sm text-blue-700 hover:text-blue-900 transition-colors flex items-center gap-1">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                            </svg>
                            Selecionar todos
                        </button>

                        <button @click="limparNiveis"
                            class="text-sm text-red-600 hover:text-red-800 transition-colors flex items-center gap-1">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                            Limpar
                        </button>

                        <button @click="abrirNivelModal = false"
                            class="text-sm ml-auto px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 transition">
                            Fechar
                        </button>
                    </div>
                </div>
            </div>

            <button @click="abrirTipoModal = true" class="filtro text-left">
                <span class="block truncate">
                    {{ filtros.tipo.length == 2 ? 'Tipo' : `Tipo
                    (${filtros.tipo.length})` }}
                </span>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="h-5 w-5 text-gray-400">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            <div v-if="abrirTipoModal" class="fixed inset-0 z-50 flex items-center justify-center">
                <!-- Fundo escuro -->
                <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="abrirTipoModal = false"></div>

                <!-- Conteúdo -->
                <div
                    class="relative z-10 bg-white flex flex-col rounded-2xl shadow-lg w-full max-w-md max-h-[75vh] p-6">
                    <h2 class="text-lg font-medium text-gray-800 mb-4">Selecionar Tipo</h2>

                    <div class="flex flex-wrap gap-2 overflow-y-auto mb-4">
                        <button v-for="tipo in tipos" :key="tipo" @click="alternarTipo(tipo)"
                            class="px-3 py-1 text-sm rounded-full transition-colors flex items-center gap-2" :class="{
                                'bg-blue-100 text-blue-700': filtros.tipo.includes(tipo),
                                'bg-gray-100 text-gray-700 hover:bg-blue-50': !filtros.tipo.includes(tipo)
                            }">
                            {{ tipo }}
                            <span v-if="filtros.tipo.includes(tipo)"
                                class="text-blue-700 text-xs font-bold hover:text-red-600"
                                @click.stop="removerTipo(tipo)">
                                ×
                            </span>
                        </button>
                    </div>

                    <!-- Botões de ação -->
                    <div class="flex justify-start items-center mt-auto pt-0 space-x-3">
                        <button @click="selecionarTodosTipos"
                            class="text-sm text-blue-700 hover:text-blue-900 transition-colors flex items-center gap-1">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                            </svg>
                            Selecionar todos
                        </button>

                        <button @click="limparTipos"
                            class="text-sm text-red-600 hover:text-red-800 transition-colors flex items-center gap-1">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                            Limpar
                        </button>

                        <button @click="abrirTipoModal = false"
                            class="text-sm ml-auto px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 transition">
                            Fechar
                        </button>
                    </div>
                </div>
            </div>

            <button @click="abrirDificuldadeModal = true" class="filtro text-left">
                <span class="block truncate">
                    {{ filtros.dificuldade.length == 5 ? 'Dificuldades' : `Dificuldade
                    (${filtros.dificuldade.length})` }}
                </span>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="h-5 w-5 text-gray-400">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
            </button>

            <div v-if="abrirDificuldadeModal" class="fixed inset-0 z-50 flex items-center justify-center">
                <div class="absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity"
                    @click="fecharDificuldade()"></div>

                <div
                    class="relative z-10 bg-white flex flex-col rounded-2xl shadow-lg w-full max-w-md min-h-auto max-h-[60vh] p-6">
                    <h2 class="text-lg font-medium text-gray-800 mb-4">Selecionar Dificuldade</h2>

                    <div class="flex flex-wrap gap-2 mb-4">
                        <button v-for="nivel in niveisDificuldade" :key="nivel" @click="toggleDificuldade(nivel)"
                            class="px-3 py-1 text-sm rounded-full transition-colors" :class="filtros.dificuldade.includes(nivel)
                                ? 'bg-blue-100 text-blue-800 hover:bg-blue-200'
                                : 'bg-gray-100 text-gray-700 hover:bg-blue-100'">
                            {{ nivel }}
                            <span v-if="filtros.dificuldade.includes(nivel)" class="ml-1 text-blue-500">&times;</span>
                        </button>
                    </div>

                    <div class="flex justify-start items-center mt-auto pt-0 space-x-3">
                        <button @click="selecionarTodasDificuldades"
                            class="text-sm text-blue-700 hover:text-blue-900 transition-colors flex items-center gap-1">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                            </svg>
                            Selecionar todos
                        </button>

                        <button @click="limparDificuldades"
                            class="text-sm text-red-600 hover:text-red-800 transition-colors flex items-center gap-1">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                            Limpar
                        </button>

                        <button @click="fecharDificuldade()"
                            class="text-sm ml-auto px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 transition">
                            Fechar
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Linha dos filtros selecionados + botões -->
        <div class="flex flex-wrap items-center justify-between gap-3">
            <!-- filtros selecionados à esquerda -->
            <div class="flex flex-wrap items-center gap-2 max-w-full">
                <p class="text-xs tracking-wide text-gray-700">Filtros ativos:</p>

                <!-- Lista dinâmica de filtros (excluindo filtros com todos os itens selecionados) -->
                <template v-for="(valores, chave) in filtros" :key="`filtro-${chave}`">
                    <template
                        v-if="valores.length && !todosSelecionados(chave, valores) && chave != 'disciplina' && chave != 'assuntos'">
                        <template v-for="(valor, index) in valores.slice(0, 3)" :key="`${chave}-${index}`">
                            <span
                                class="bg-blue-100 text-blue-800 text-xs px-3 py-1 rounded-full flex items-center gap-1 cursor-default select-none"
                                :title="`${formatarChave(chave)}: ${valor}`">
                                {{ formatarChave(chave) }}: {{ valor.nome ? (Array.isArray(valor.nome) ? (valor.nome[0]
                                    + ' - ' + valor.nome[1]) : valor.nome) : valor }}
                                <button @click="removerFiltro(chave, valor)"
                                    class="ml-1 text-blue-700 hover:text-red-600" aria-label="Remover">&times;</button>
                            </span>
                        </template>

                        <!-- Mostrar +X caso tenha mais de 3 -->
                        <button v-if="valores.length > 3" @click="abrirModal(chave)"
                            class="bg-gray-200 text-gray-600 text-xs px-2 py-1 rounded-full hover:bg-gray-300">
                            +{{ valores.length - 3 }}
                        </button>
                    </template>

                    <template v-if="valores.length && !todosSelecionados(chave, valores) && chave == 'assuntos'">
                        <template v-for="(valor, index) in valores.slice(0, 3)" :key="`${chave}-${index}`">
                            <span
                                class="bg-blue-100 text-blue-800 text-xs px-3 py-1 rounded-full flex items-center gap-1 cursor-default select-none"
                                :title="`${formatarChave(chave)}: ${valor[0]}`">
                                {{ formatarChave(chave) }}: {{ valor[0] }}
                                <button @click="removerFiltro(chave, valor)"
                                    class="ml-1 text-blue-700 hover:text-red-600" aria-label="Remover">&times;</button>
                            </span>
                        </template>

                        <!-- Mostrar +X caso tenha mais de 3 -->
                        <button v-if="valores.length > 3" @click="abrirModal(chave)"
                            class="bg-gray-200 text-gray-600 text-xs px-2 py-1 rounded-full hover:bg-gray-300">
                            +{{ valores.length - 3 }}
                        </button>
                    </template>
                </template>
            </div>

            <!-- Botões à direita -->
            <div class="flex gap-2 ml-auto">
                <button @click="limpar"
                    class="cursor-pointer px-4 py-2 text-sm rounded-md bg-gray-200 text-gray-700 hover:bg-gray-300 transition">
                    Limpar
                </button>

                <button @click="aplicar"
                    class="cursor-pointer px-4 py-2 text-sm rounded-md bg-[#2c89a0] text-white hover:bg-blue-700 transition">
                    Aplicar
                </button>
            </div>

            <!-- Modal genérico para todos os filtros -->
            <div v-if="modalAberto"
                class="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center">
                <div class="bg-white rounded-2xl shadow-lg p-6 w-full max-w-md max-h-[70vh] overflow-y-auto relative">
                    <h3 class="text-base font-semibold text-gray-800 mb-4">Todos(as) {{ formatarChave(filtroAtual) }}s
                    </h3>
                    <div class="flex flex-wrap gap-2">
                        <template v-for="(valor, index) in filtros[filtroAtual]" :key="`modal-${filtroAtual}-${index}`">
                            <span
                                class="bg-blue-100 text-blue-800 text-xs px-3 py-1 rounded-full flex items-center gap-1 cursor-default select-none">
                                {{ valor.nome ? (Array.isArray(valor.nome) ? (valor.nome[0] + ' - ' + valor.nome[1]) :
                                    valor.nome) : valor }}
                                <button @click="removerFiltro(filtroAtual, valor)"
                                    class="ml-1 text-blue-700 hover:text-red-600" aria-label="Remover">&times;</button>
                            </span>
                        </template>
                    </div>
                    <button @click="fecharModal"
                        class="absolute top-2 right-2 text-gray-500 hover:text-gray-800 text-xl">&times;</button>
                </div>
            </div>
        </div>
    </div>

    <div v-if="abrirBancasModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Fundo com blur e opacidade -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity" @click="fecharBancas()"></div>
        <div class="relative z-10 bg-white flex flex-col rounded-2xl shadow-lg 
        max-w-[90vw] md:w-[50vw] min-w-[50vw] min-h-[75vh] max-h-[75vh] p-6">
            <h2 class="text-lg font-medium text-gray-800 mb-4">Selecionar Bancas</h2>

            <BuscaBancas v-model="filtros.banca" />

            <div class="mt-auto text-right">
                <button @click="fecharBancas()"
                    class="text-sm px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 transition">
                    Fechar
                </button>
            </div>

        </div>
    </div>

    <div v-if="abrirCargosModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Fundo com blur e opacidade -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity" @click="fecharCargos()"></div>
        <div class="relative z-10 bg-white flex flex-col rounded-2xl shadow-lg 
        max-w-[90vw] md:w-[50vw] min-w-[50vw] min-h-[75vh] max-h-[75vh] p-6">
            <h2 class="text-lg font-medium text-gray-800 mb-4">Selecionar Cargos <span
                    class="text-xs text-gray-500 ml-0">(No máximo 90)</span> </h2>

            <BuscaCargos v-model="filtros.cargo" />
        </div>
    </div>

    <div v-if="abrirOrgaosModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Fundo com blur e opacidade -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity" @click="fecharOrgaos()"></div>
        <div class="relative z-10 bg-white flex flex-col rounded-2xl shadow-lg 
            max-w-[90vw] md:w-[50vw] min-w-[50vw] min-h-[75vh] max-h-[75vh] p-6">
            <h2 class="text-lg font-medium text-gray-800 mb-4">Selecionar Orgãos</h2>

            <BuscaOrgaos v-model="filtros.orgao" />

            <div class="mt-auto text-right">
                <button @click="fecharOrgaos()"
                    class="text-sm px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 transition">
                    Fechar
                </button>
            </div>

        </div>
    </div>

    <div v-if="abrirAssuntoModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Fundo com blur e opacidade -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity" @click="fecharAssunto()">
        </div>

        <div
            class="relative z-10 bg-white flex flex-col rounded-2xl shadow-lg min-w-[90vw] md:min-w-[50vw] max-w-[90vw] min-h-[75vh] max-h-[75vh] p-6">
            <h2 class="text-lg font-medium text-gray-800 mb-4">Selecionar Assuntos</h2>

            <DisciplinaTree />

            <div class="mt-auto text-right">
                <button @click="fecharAssunto()"
                    class="text-sm px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 transition">
                    Fechar
                </button>
            </div>

        </div>
    </div>

    <!-- Modal: seleção de disciplinas (estilo moderno e minimalista) -->
    <div v-if="abrirDisciplinaModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Fundo com blur e opacidade -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity" @click="fecharModalDisciplina()">
        </div>

        <!-- Conteúdo do modal -->
        <div
            class="relative z-10 bg-white flex flex-col rounded-2xl shadow-lg min-w-[90vw] md:min-w-[50vw] max-w-[90vw] min-h-[75vh] max-h-[75vh] p-6">
            <h2 class="text-lg font-medium text-gray-800 mb-4">Selecionar Disciplinas</h2>

            <div class="relative mb-2">
                <input v-model="buscaDisciplina" type="text" placeholder="Buscar..."
                    class="w-full px-4 py-2 pl-10 text-sm border border-gray-200 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-300" />
                <svg class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 transform -translate-y-1/2" fill="none"
                    stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M21 21l-4.35-4.35m0 0A7 7 0 104.5 4.5a7 7 0 0012.15 12.15z" />
                </svg>
            </div>

            <div class="mb-2 flex-grow overflow-y-auto space-y-2 p-2 bg-gray-50 rounded-md max-h-max">
                <label v-for="d in disciplinasFiltradas" :key="d"
                    class="flex items-center gap-2 text-sm text-gray-700 hover:bg-gray-50 px-2 py-1 rounded-md cursor-pointer">
                    <input type="checkbox" :value="{ id: d[1] }" v-model="filtros.disciplina"
                        class="accent-blue-600 w-4 h-4" />
                    <span>{{ d[0] }}</span>
                </label>
            </div>

            <div class="mt-auto text-right">
                <button @click="fecharModalDisciplina()"
                    class="text-sm px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 transition">
                    Fechar
                </button>
            </div>
        </div>
    </div>

</template>

<script setup>
import { reactive, ref, computed, nextTick, watch } from 'vue'
import DisciplinaTree from '../../src/components/DisciplinaTree.vue'
import BuscaOrgaos from '../../src/components/BuscaOrgaos.vue'
import BuscaCargos from '../../src/components/BuscaCargos.vue'
import BuscaBancas from '../../src/components/BuscaBancas.vue'

let disciplinas = []
const anos = Array.from({ length: 15 }, (_, i) => 2025 - i)
const niveisDificuldade = ['Muito Fácil', 'Fácil', 'Médio', 'Difícil', 'Muito Difícil']

const abrirAnoModal = ref(false)
const abrirDificuldadeModal = ref(false)

const mostrarTodosAssuntos = ref(false)

const abrirDisciplinaModal = ref(false)
const buscaDisciplina = ref('')

const abrirAssuntoModal = ref(false)

const abrirOrgaosModal = ref(false)

const abrirCargosModal = ref(false)
const cargosSelecionados = ref([])

const abrirBancasModal = ref(false)
const bancasSelecionados = ref([])

const abrirNivelModal = ref(false)
const niveis = ['Fundamental', 'Médio', 'Superior']

const abrirTipoModal = ref(false)
const tipos = ['Certo ou Errado', 'Múltipla Escolha']

const filtros = reactive({
    disciplina: [], assuntos: [], orgao: [], cargo: [],
    banca: [], ano: [2025, 2024, 2023, 2022, 2021, 2020, 2019, 2018, 2017, 2016, 2015, 2014, 2013, 2012, 2011], nivel: ['Fundamental', 'Médio', 'Superior'], tipo: ['Certo ou Errado', 'Múltipla Escolha'], dificuldade: ['Muito Fácil', 'Fácil', 'Médio', 'Difícil', 'Muito Difícil']
})

function alternarTipo(tipo) {
    const index = filtros.tipo.indexOf(tipo)
    if (index > -1) {
        filtros.tipo.splice(index, 1)
    } else {
        filtros.tipo.push(tipo)
    }
}

function removerTipo(tipo) {
    filtros.tipo = filtros.tipo.filter((t) => t !== tipo)
}

function selecionarTodosTipos() {
    filtros.tipo = [...tipos]
}

function limparTipos() {
    filtros.tipo = []
}

function alternarNivel(nivel) {
    const index = filtros.nivel.indexOf(nivel)
    if (index > -1) {
        filtros.nivel.splice(index, 1)
    } else {
        filtros.nivel.push(nivel)
    }
}

function removerNivel(nivel) {
    filtros.nivel = filtros.nivel.filter((n) => n !== nivel)
}

function selecionarTodosNiveis() {
    filtros.nivel = [...niveis]
}

function limparNiveis() {
    filtros.nivel = []
}

function toggleDificuldade(nivel) {
    const index = filtros.dificuldade.indexOf(nivel)
    if (index > -1) {
        filtros.dificuldade.splice(index, 1)
    } else {
        filtros.dificuldade.push(nivel)
    }
}

function selecionarTodasDificuldades() {
    filtros.dificuldade = [...niveisDificuldade]
}

function limparDificuldades() {
    filtros.dificuldade = []
}

function fecharDificuldade() {
    abrirDificuldadeModal.value = false
}

function alternarAno(ano) {
    const index = filtros.ano.indexOf(ano)
    if (index === -1) {
        filtros.ano.push(ano)
    } else {
        filtros.ano.splice(index, 1)
    }
}

function removerAno(ano) {
    filtros.ano = filtros.ano.filter(a => a !== ano)
}

function selecionarTodosAnos() {
    filtros.ano = [...anos]
}

function limparAnos() {
    filtros.ano = []
}

function fecharBancas() {
    abrirBancasModal.value = false
}

function fecharCargos() {
    abrirCargosModal.value = false
}

window.fecharCargos = fecharCargos

function fecharOrgaos() {
    abrirOrgaosModal.value = false
}

function abrirDisciplina() {
    fetch('../../questoes/src/data/asst_conc/tree.json')
        .then(response => {
            if (!response.ok) {
                throw new Error('Erro na requisição: ' + response.status)
            }
            return response.json()
        })
        .then(data => {
            disciplinas = data
            abrirDisciplinaModal.value = true
        })
        .catch(error => {
            console.error('Erro:', error)
        })
}

function removeAcentos(str) {
    return str.normalize("NFD").replace(/[\u0300-\u036f]/g, "")
}

const disciplinasFiltradas = computed(() => {
    if (!buscaDisciplina.value) return disciplinas
    const busca = removeAcentos(buscaDisciplina.value.toLowerCase())
    return disciplinas.filter(d =>
        removeAcentos(d[0].toLowerCase()).includes(busca)
    )
})

function abrirAssunto() {
    nextTick(() => {
        const keys = Object.keys(filtros.disciplina)
        if (keys.length == 0 || keys == undefined) {
            window.mostrarPopup('Atenção', 'Nenhuma Disciplina selecionada', "warning", 4000)
        } else {
            let temp = []
            keys.forEach((item, i) => {
                temp.push(filtros.disciplina[item].id)
                if (i == (keys.length - 1)) {
                    localStorage.setItem("disciplinasIds", JSON.stringify(temp))
                    abrirAssuntoModal.value = true
                }
            })
        }
    })
}

function fecharModalDisciplina() {
    abrirDisciplinaModal.value = false
}

const todosPossiveis = {
    ano: [2025, 2024, 2023, 2022, 2021, 2020, 2019, 2018, 2017, 2016, 2015, 2014, 2013, 2012, 2011],
    nivel: ['Fundamental', 'Médio', 'Superior'],
    tipo: ['Certo ou Errado', 'Múltipla Escolha'],
    dificuldade: ['Muito Fácil', 'Fácil', 'Médio', 'Difícil', 'Muito Difícil']
}

function todosSelecionados(chave, valores) {
    return (
        todosPossiveis[chave] &&
        valores.length === todosPossiveis[chave].length &&
        todosPossiveis[chave].every((v) => valores.includes(v))
    )
}

function getAllChildrenById(data, targetId) {
    let result = [];

    function findNodeAndCollect(children) {
        for (const node of children) {
            if (node.id === targetId) {
                collectIds(node);
                return true;
            }
            if (findNodeAndCollect(node.filhos)) {
                return true;
            }
        }
        return false;
    }

    function collectIds(node) {
        for (const filho of node.filhos) {
            result.push(filho.id);
            collectIds(filho);
        }
    }

    findNodeAndCollect(data);
    return result;
}

function removerFiltro(chave, valor) {
    filtros[chave] = filtros[chave].filter((v) => v !== valor)
    if (chave == "assuntos") {
        const childsIds = [valor[1], ...getAllChildrenById(window.arvoreFiltrada.value, valor[1])];
        window.selecionados.value = window.selecionados.value.filter(
            id => !childsIds.includes(id)
        );
        localStorage.setItem("selecionados", JSON.stringify(window.selecionados.value))
    }
}

// Utilitário para deixar o nome do filtro legível
function formatarChave(chave) {
    return {
        disciplina: 'Disciplina',
        assuntos: 'Assunto',
        orgao: 'Órgão',
        cargo: 'Cargo',
        banca: 'Banca',
        ano: 'Ano',
        nivel: 'Nível',
        tipo: 'Tipo',
        dificuldade: 'Dificuldade'
    }[chave] || chave
}

// Controle de modal
const modalAberto = ref(false)
const filtroAtual = ref(null)

function abrirModal(chave) {
    filtroAtual.value = chave
    modalAberto.value = true
}

function fecharModal() {
    modalAberto.value = false
    filtroAtual.value = null
}

function limpar() {
    Object.assign(filtros, {
        disciplina: [], assuntos: [], orgao: [], cargo: [],
        banca: [], ano: [2025, 2024, 2023, 2022, 2021, 2020, 2019, 2018, 2017, 2016, 2015, 2014, 2013, 2012, 2011], nivel: ['Fundamental', 'Médio', 'Superior'], tipo: ['Certo ou Errado', 'Múltipla Escolha'], dificuldade: ['Muito Fácil', 'Fácil', 'Médio', 'Difícil', 'Muito Difícil']
    })
}

function normalizarFiltros(filtros, todosPossiveis) {
    const novo = { ...filtros }

    for (const chave in todosPossiveis) {
        if (Array.isArray(novo[chave])) {
            const arr1 = novo[chave]
            const arr2 = todosPossiveis[chave]

            // compara independente da ordem (como conjuntos)
            const iguais =
                arr1.length === arr2.length &&
                arr1.every(v => arr2.includes(v))

            if (iguais) {
                novo[chave] = []
            }
        }
    }

    return novo
}

let filtrosNormalizados = normalizarFiltros(filtros, todosPossiveis)


watch(
    filtros,
    (novoValor, valorAntigo) => {
        filtrosNormalizados = normalizarFiltros(novoValor, todosPossiveis)
        window.filtrosNormalizados = filtrosNormalizados
    },
    { deep: true }
);

function aplicar() {
    window.pagina = 1
    window.loadPage(1, 0, filtrosNormalizados)
}

function formatarLabel(chave) {
    // Para deixar os nomes bonitos e com maiúsculas iniciais
    const map = {
        disciplina: 'Disciplina',
        assunto: 'Assunto',
        orgao: 'Orgão',
        cargo: 'Cargo',
        banca: 'Banca',
        ano: 'Ano',
        tipo: 'Tipo',
        dificuldade: 'Dificuldade',
    }
    return map[chave] || chave
}

function fecharAssunto() {
    abrirAssuntoModal.value = false
}

window.filtros = filtros
window.filtrosNormalizados = filtrosNormalizados

</script>

<style>
@reference "../style.css";

.filtro {
    @apply cursor-pointer w-full flex items-center justify-between rounded-sm bg-gray-100 px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-gray-100;
}
</style>