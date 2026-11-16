# Mini-Spec: Do Ticket Vago à Especificação Executável

## Ticket Vago Original
> *"Fazer o redirecionamento mascarado das APIs cadastradas."*

---

## Mini-Spec Reescrita

### 1. Objetivo e Público
Permitir que **desenvolvedores consumidores** realizem requisições HTTP (GET, POST, PUT, DELETE) contra rotas internas padronizadas no formato `/apihub/{servico}/{recurso}` e recebam a resposta da API externa correspondente de forma transparente, ocultando a infraestrutura e o domínio original.

### 2. Critérios de Aceitação Verificáveis
1. **Resolução de Rota:** Dado um cadastro ativo com slug `pagamento/pix` apontando para `https://api.banco.exemplo/v2/pix`, uma chamada `POST /apihub/pagamento/pix` deve repassar método, headers autorizados, query strings e body para a URL original e devolver o status code e o body originais.
2. **Ocultação de Origem:** A resposta entregue ao cliente não pode conter headers que revelem a URL original (ex.: remover/sobrescrever `X-Powered-By`, `Server` ou redirecionamentos `Location` absolutos apontando para o host externo).
3. **Rota Inexistente:** Chamadas para um slug não cadastrado ou inativo devem retornar HTTP `404 Not Found` com JSON `{"erro": "rota não encontrada"}` sem realizar chamadas externas.

### 3. Requisito Não Funcional Mensurável
- **Overhead de Proxy (RNF-01):** O tempo adicionado pelo motor de proxy do APIHub (excluindo o tempo de resposta da API externa mockada) deve ser **inferior a 15 ms** em 95% das requisições em testes locais.

### 4. Contexto de Arquitetura
- A requisição entra pela **Camada de Entrega/HTTP (`handler`)**, consulta o mapeamento na **Camada de Domínio/Serviço (`service`)** que obtém os dados do **Repositório (`repository`)**, e aciona o **Adaptador de Proxy (`proxy`)** baseado em `net/http/httputil.ReverseProxy`. Nenhuma consulta ao armazenamento de rotas pode ocorrer diretamente dentro da função de transporte do proxy.

### 5. Armadilha Técnica (Edge Case Crítico)
- **Sobrescrita do Header `Host` e Concatenação de Path:** Ao usar `httputil.ReverseProxy` em Go, um erro comum é manter o header `req.Host` original da requisição recebida (`localhost:8080`). Servidores externos (especialmente atrás de CDNs como Cloudflare ou AWS API Gateway) rejeitam requisições com erro `403 Forbidden` ou `421 Misdirected Request` se o `req.Host` não corresponder ao domínio da URL de destino. Além disso, o prefixo `/apihub/{slug}` deve ser removido antes de juntar ao path da URL original para evitar URLs duplicadas como `/v2/pix/apihub/pagamento/pix`.