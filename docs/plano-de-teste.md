---

### 8. `docs/plano-de-teste.md` (Entregável E7)
```markdown
# Plano de Teste — APIHub

## 1. O que será testado
Serão testados o cadastro e validação de rotas mascaradas, a listagem do catálogo de APIs, a reescrita de URLs/headers no motor de proxy reverso, o tratamento de rotas inexistentes (`404`) e o comportamento diante de falhas/timeouts de APIs externas (`502`/`504`).

## 2. Níveis de Teste
- **Teste de Unidade (`*_test.go`):** Validação de regras de negócio no pacote `service` (formatação de slugs, rejeição de URLs malformadas e duplicidade).
- **Teste de Componente / Integração (`httptest`):** Testes subindo servidores HTTP simulados (`httptest.NewServer`) para validar o fluxo completo `Handler -> Service -> Proxy -> Mock API Externa`.
- **Teste de Sistema / Carga:** Execução do binário compilado com disparo concorrente de requisições para validar o limite de latência do `RNF-01`.
- **Teste de Aceitação:** Verificação manual e automatizada dos critérios de aceitação definidos para cada história da Sprint.

## 3. Item do Backlog com TDD (Teste Antes do Código)
- **Item escolhido:** **H02 — Motor de proxy reverso para mascaramento de rotas (`/apihub/{slug}`)**.
- **Justificativa:** O mapeamento de paths, reescrita do header `Host` e remoção do prefixo `/apihub` possuem entradas e saídas determinísticas fáceis de simular com `httptest.NewServer`. Escreveremos os testes de falha primeiro (vermelho), implementaremos o proxy (verde) e refatoraremos.
- *Discussão (Pior candidato a TDD):* O pior candidato a TDD no backlog é o **H04 (Interface Web HTML)**, pois layout visual e usabilidade sofrem ajustes subjetivos rápidos onde testes automatizados antes da interface geram retrabalho alto sem ganho lógico.

## 4. Registro dos Resultados
Os resultados dos testes automatizados são registrados a cada commit/PR pela saída do runner `go test -v -race ./...` no GitHub Actions. Para testes manuais e de aceitação/desempenho, será mantida uma tabela de registro no repositório contendo: **ID do Requisito/História | Resultado Esperado | Status (Aprovado/Reprovado) | Data | Responsável**.

## 5. Cronograma Previsto
- **Semana 1 (Sprint 1):** Escrita dos testes unitários de `service` (`H01`, `H03`) e ciclo TDD completo do `proxy` (`H02`).
- **Semana 2 (Sprint 1):** Testes de integração com simulação de timeout/falha (`H05`) e teste de benchmark de latência (`RNF-01`).
- **Semanas 3–4 (Sprint 2):** Testes de regressão contínua no CI, testes das histórias `H04`, `H06`, `H07`, `H08` e adaptação à mudança surpresa.

---

## Tabela: Histórias Prioritárias × Critérios de Aceitação

| História | Critérios de Aceitação Verificáveis | Tipo | Automatizável? |
| :--- | :--- | :---: | :---: |
| **H01 — Cadastro de Rota Mascarada** | 1. Enviar `POST /api/rotas` com nome, descrição, slug único e URL HTTP/HTTPS válida retorna `201 Created` e persiste o item. | Funcional | **Sim** |
| | 2. Tentar cadastrar um slug já existente ou URL sem esquema (`http://` ou `https://`) retorna `400 Bad Request` ou `409 Conflict`. | Funcional | **Sim** |
| | 3. A operação de escrita não deve gerar data race sob 50 requisições concorrentes (`go test -race`). | Não Funcional | **Sim** |
| **H02 — Proxy Reverso Mascarado** | 1. Chamar `GET /apihub/{slug}` encaminha query params e headers para a URL original e retorna o mesmo status code e body da API externa. | Funcional | **Sim** |
| | 2. O header `Host` enviado à API externa corresponde ao domínio de destino, e a resposta ao cliente omite a URL original. | Funcional | **Sim** |
| | 3. O overhead interno de roteamento do proxy deve ser ≤ 15 ms no p95 em execução local (`RNF-01`). | Não Funcional | **Sim** |
| **H03 — Catálogo de APIs Disponíveis** | 1. Chamar `GET /api/catalogo` retorna `200 OK` com array JSON contendo nome, descrição, URL original e URL mascarada completa de cada API. | Funcional | **Sim** |
| | 2. Quando não houver APIs cadastradas, retorna `200 OK` com array vazio `[]` em menos de 10 ms. | Não Funcional | **Sim** |