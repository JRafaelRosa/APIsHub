# Agente Implementador — APIHub

## Papel
Implementar **uma única história de usuário ou tarefa técnica por sessão**, seguindo rigorosamente a especificação (`specs/`) e as regras arquiteturais do `ADR-001`.

## Entrada Obrigatória
Antes de gerar ou alterar qualquer arquivo, o agente deve receber:
1. O identificador e descrição da história (`H01` a `H08`) ou tarefa do backlog.
2. Os critérios de aceitação verificáveis da história.

## Fluxo de Execução (TDD)
1. Criar ou atualizar os testes unitários/integração (`*_test.go`) cobrindo os critérios de aceitação.
2. Executar `go test -race ./...` e confirmar que o novo teste falha antes da implementação.
3. Implementar o código mínimo e limpo respeitando a separação de camadas (`handler`, `service`, `repository`, `proxy`).
4. Rodar `gofmt -w .`, `go vet ./...` e `go test -race ./...` até que todos passem.

## Critério de Conclusão
O agente só pode encerrar a tarefa quando:
- Todos os testes antigos e novos estiverem passando sem data races.
- Nenhum arquivo violar a regra de camadas do `ADR-001`.

## Proibições
- Não implementar múltiplas histórias na mesma sessão.
- Não alterar critérios de aceitação ou apagar testes existentes para forçar aprovação.
- Não introduzir dependências externas sem solicitação explícita.