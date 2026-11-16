# Histórias de Usuário e Backlog Priorizado

## Histórias de Usuário (2 por integrante)

### H01 — Cadastro de Rota Mascarada (Pessoa 1)
- **Narrativa:** **Como** administrador do sistema, **quero** cadastrar uma API externa informando nome, descrição, slug mascarado e URL original, **para** centralizar o acesso sem expor o endereço real aos clientes.
- **Fluxo Principal:** O usuário envia os dados via endpoint administrativo; o sistema valida a URL e a unicidade do slug e salva o registro.
- **Fluxo Alternativo:** Se o slug já existir ou a URL for inválida, o sistema rejeita com HTTP `400 Bad Request` ou `409 Conflict`.
- **Requisitos Ligados:** RF-01, RNF-03.

### H02 — Proxy Reverso Mascarado (Pessoa 1)
- **Narrativa:** **Como** desenvolvedor consumidor, **quero** chamar endpoints no formato `/apihub/{slug}`, **para** consumir APIs externas usando um único domínio base.
- **Fluxo Principal:** O gateway localiza o slug, reescreve o host e o path, encaminha a requisição e devolve a resposta da API externa.
- **Fluxo Alternativo:** Se o slug não existir, retorna HTTP `404 Not Found`.
- **Requisitos Ligados:** RF-03, RNF-01.

### H03 — Catálogo de APIs Disponíveis (Pessoa 2)
- **Narrativa:** **Como** desenvolvedor consumidor, **quero** consultar uma lista com todas as APIs cadastradas, suas descrições, URLs originais e URLs mascaradas, **para** saber para que serve cada rota e como integrá-la.
- **Fluxo Principal:** O usuário acessa o endpoint/tela de catálogo e visualiza a tabela estruturada de todas as APIs registradas.
- **Fluxo Alternativo:** Caso não haja APIs cadastradas, retorna uma lista vazia com status `200 OK`.
- **Requisitos Ligados:** RF-02.

### H04 — Interface Web Simples de Visualização e Cadastro (Pessoa 2)
- **Narrativa:** **Como** usuário do APIHub, **quero** visualizar o catálogo e cadastrar novas rotas por uma página HTML simples servida pelo próprio sistema, **para** gerenciar as APIs sem depender apenas de ferramentas de linha de comando como curl.
- **Fluxo Principal:** O usuário acessa `/`, preenche o formulário de nova API e vê a lista atualizada imediatamente.
- **Fluxo Alternativo:** Mensagem de erro visível caso o cadastro falhe por slug duplicado.
- **Requisitos Ligados:** RF-01, RF-02, RNF-03.

### H05 — Tratamento de Timeout e Falha Externa (Pessoa 3)
- **Narrativa:** **Como** desenvolvedor consumidor, **quero** receber um erro claro e rápido quando uma API externa cair ou travar, **para** que minha aplicação cliente não fique travada aguardando resposta.
- **Fluxo Principal:** O proxy aguarda até o limite de 5 segundos; ao estourar, encerra a conexão e retorna HTTP `504 Gateway Timeout`.
- **Fluxo Alternativo:** Se a conexão for recusada imediatamente, retorna HTTP `502 Bad Gateway`.
- **Requisitos Ligados:** RF-05, RNF-02.

### H06 — Remoção e Desativação de APIs (Pessoa 3)
- **Narrativa:** **Como** administrador do sistema, **quero** remover ou desativar uma API cadastrada, **para** interromper imediatamente o roteamento de serviços descontinuados.
- **Fluxo Principal:** O administrador solicita exclusão pelo ID/slug; o sistema remove o mapeamento e confirma com HTTP `204 No Content`.
- **Fluxo Alternativo:** Se o ID não existir, retorna HTTP `404 Not Found`.
- **Requisitos Ligados:** RF-04.

### H07 — Logs Estruturados de Requisições Mascaradas (Pessoa 4)
- **Narrativa:** **Como** operador do sistema, **quero** visualizar no terminal logs estruturados contendo rota mascarada acessada, status HTTP retornado e tempo de resposta, **para** diagnosticar lentidão ou falhas nas APIs integradas.
- **Fluxo Principal:** A cada chamada no `/apihub/*`, um middleware registra método, slug, status e latência em ms.
- **Requisitos Ligados:** RNF-01, RNF-02.

### H08 — Edição de URL Original sem Mudar a Rota Mascarada (Pessoa 4)
- **Narrativa:** **Como** administrador do sistema, **quero** alterar a URL original de uma rota já existente mantendo o mesmo slug mascarado, **para** migrar provedores externos sem quebrar os sistemas clientes.
- **Fluxo Principal:** O administrador envia a nova URL de destino; o sistema atualiza o registro e passa a rotear as próximas chamadas para o novo endereço.
- **Requisitos Ligados:** RF-04.

---

## Backlog Priorizado (Product Owner)

| Ordem | ID | Descrição do Item | Tipo | Sprint 1? |
| :---: | :--- | :--- | :--- | :---: |
| **1** | **ARQ-01** | Estruturar projeto Go em camadas (`handler`, `service`, `repository`, `proxy`) e configurar pipeline de CI (`go test` + `go vet`) | Arquitetura | **Sim** |
| **2** | **H01** | Cadastro e persistência de APIs externas (nome, descrição, slug mascarado e URL original) | História | **Sim** |
| **3** | **H02** | Motor de proxy reverso para mascaramento de rotas (`/apihub/{slug}`) | História | **Sim** |
| **4** | **H03** | Endpoint de listagem do catálogo de APIs (URL original, mascarada e descrição) | História | **Sim** |
| **5** | **H05** | Tratamento de timeout (5s) e respostas padronizadas de erro (`502` / `504`) | História | **Sim** |
| **6** | **DOC-01** | Documentação de uso da API (exemplos `curl` e guia de execução local no `README.md`) | Documentação | **Sim** |
| **7** | **H04** | Interface web leve (HTML/JS nativo) para visualizar o catálogo e cadastrar APIs | História | Não (Sprint 2) |
| **8** | **H06** | Remoção e desativação de rotas cadastradas | História | Não (Sprint 2) |
| **9** | **H08** | Edição da URL original mantendo o slug mascarado intacto | História | Não (Sprint 2) |
| **10** | **H07** | Middleware de logs estruturados e métricas de latência por rota | História | Não (Sprint 2) |

> **Nota de discussão (Cobertura dos requisitos):** Sabemos que as histórias da Sprint 1 cobrem os requisitos essenciais porque há rastreabilidade direta 1-para-1 entre os itens priorizados e os requisitos centrais do núcleo funcional: `H01` cobre `RF-01`, `H02` cobre `RF-03`, `H03` cobre `RF-02` e `H05` cobre `RF-05` e `RNF-02`.