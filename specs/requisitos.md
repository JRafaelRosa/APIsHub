# Especificação de Requisitos — APIHub

## Requisitos Funcionais (RF)

- **RF-01 — Cadastro de API Externa:** O sistema deve permitir registrar uma nova API informando nome, descrição de finalidade, rota mascarada (slug único, ex.: `/pagamento/pix`) e URL de destino original.
    - *Critério de aceitação:* Ao enviar um payload válido com slug inédito e URL bem-formada, o sistema persiste o registro e retorna status HTTP `201 Created` com os dados salvos.

- **RF-02 — Listagem e Catálogo de APIs:** O sistema deve listar todas as APIs cadastradas, exibindo nome, descrição, URL original e a URL mascarada completa exposta pelo APIHub.
    - *Critério de aceitação:* Uma requisição `GET /api/catalogo` retorna HTTP `200 OK` com a lista JSON contendo todos os campos obrigatórios de cada API ativa.

- **RF-03 — Roteamento e Mascaramento de Requisições (Proxy):** Ao receber uma requisição em `/apihub/{slug}` (e sub-rotas), o sistema deve encaminhar o método HTTP, query params e corpo da requisição para a URL original correspondente e devolver a resposta ao cliente sem expor a URL original nos headers de resposta.
    - *Critério de aceitação:* Uma chamada `GET /apihub/clima/regiao?uf=PR` retorna exatamente o status code e o payload da API original cadastrada, omitindo o domínio original dos cabeçalhos devolvidos ao cliente.

- **RF-04 — Edição e Remoção de Rotas:** O sistema deve permitir atualizar a URL de destino ou descrição de uma API cadastrada, bem como desativar/remover um mapeamento existente.
    - *Critério de aceitação:* Após remover ou desativar uma rota, chamadas subsequentes para `/apihub/{slug}` devem retornar imediatamente HTTP `404 Not Found`.

- **RF-05 — Tratamento de Indisponibilidade da API Externa:** Caso a API original demore a responder ou esteja fora do ar, o sistema deve interceptar o erro de rede e retornar uma resposta padronizada ao cliente.
    - *Critério de aceitação:* Se o servidor de destino recusar conexão ou exceder o timeout, o APIHub retorna HTTP `502 Bad Gateway` ou `504 Gateway Timeout` em formato JSON padronizado.

## Requisitos Não Funcionais (RNF)

- **RNF-01 — Baixa Latência de Repasse (Mensurável):** O overhead de processamento interno do APIHub (tempo entre receber a requisição e despachar ao destino, somado ao retorno) não deve ultrapassar **15 ms** no percentil 95 (p95) sob carga de 100 requisições simultâneas locais.
    - *Critério de aceitação:* Teste de benchmark/carga local demonstrando p95 de overhead interno ≤ 15 ms.

- **RNF-02 — Timeout de Segurança:** Nenhuma requisição encaminhada a uma API externa pode bloquear recursos do servidor indefinidamente, devendo respeitar um timeout máximo configurável de **5 segundos**.
    - *Critério de aceitação:* Requisições a um mock server com atraso superior a 5s são abortadas pelo gateway em até 5,1s com retorno `504 Gateway Timeout`.

- **RNF-03 — Portabilidade e Execução:** O sistema deve compilar em um único binário estático em Go, sem dependências externas de runtime além do sistema operacional.
    - *Critério de aceitação:* Execução bem-sucedida de `CGO_ENABLED=0 go build` gerando binário funcional.