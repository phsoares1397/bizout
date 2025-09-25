// Função que transforma o JSON em uma estrutura de árvore recursiva
export function construirArvore(data) {
  if (!Array.isArray(data)) return []

  function parseNode(node) {
    const [nome, id, filhos] = node
    return {
      nome,
      id,
      filhos: Array.isArray(filhos) ? filhos.map(parseNode) : []
    }
  }

  return Array.isArray(data[0]) ? data.map(parseNode) : [parseNode(data)]
}