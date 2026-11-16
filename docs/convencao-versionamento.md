# Convenção de Versionamento — APIHub

## Regras Principais
1. **Repositório:** Hospedado no GitHub sob controle da equipe, com acesso público (ou concedido explicitamente ao professor `petroskilp`). A branch `main` é protegida contra push direto.
2. **Branches:** Toda tarefa ou história deve ser desenvolvida em uma branch isolada criada a partir da `main`, seguindo o padrão `tipo/id-descricao-curta` (ex.: `feat/h01-cadastro-rotas`, `fix/proxy-host-header`, `docs/adr-001`).
3. **Fusão (Merge):** Nenhuma branch pode ser integrada à `main` sem abrir um Pull Request (PR) e obter revisão e aprovação manual de pelo menos **1 integrante humano** diferente do autor do código.
4. **Teste e Lint antes de integrar:** Antes de qualquer merge, o pipeline de Integração Contínua (CI) deve executar com sucesso:
    - Verificação de formatação (`test -z "$(gofmt -l .)"`)
    - Análise estática (`go vet ./...`)
    - Suíte completa de testes com detector de concorrência (`go test -race ./...`)

## Tratamento de Exceções
- **Conflitos de Merge (duas pessoas alteram o mesmo trecho):** O autor do PR mais recente deve atualizar sua branch local com `git fetch origin && git rebase origin/main` (ou `git merge origin/main`), resolver o conflito manualmente garantindo que a lógica de ambas as alterações seja preservada, rodar `go test -race ./...` localmente e atualizar o PR para nova revisão.
- **Build Quebrado na `main` ("Nunca quebre a construção"):** Caso um merge quebre o build ou os testes na `main`, o desenvolvimento de novas features é imediatamente pausado. O autor do último merge tem prioridade máxima para reverter o commit (`git revert`) ou subir um hotfix corrigindo o problema antes que qualquer outra branch seja integrada.

## Responsabilidade pelo Código Integrado
- **Regra do Grupo:** Todo código gerado por agentes de IA (OpenCode/Copilot) é tratado apenas como rascunho; **o integrante humano autor do Pull Request e o revisor humano que o aprovou respondem integralmente pela qualidade, segurança e funcionamento do código integrado na `main`.**