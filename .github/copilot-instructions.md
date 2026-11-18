# Instruções Globais do Projeto — APIHub (Go)

## Objetivo do Sistema
O APIHub é um gateway centralizador em Go que cataloga APIs externas (nome, descrição, URL original) e expõe rotas mascaradas no formato `/apihub/{slug}` via proxy reverso HTTP.

## Regra Obrigatória de Arquitetura (ADR-001)
- **Arquitetura em Camadas:** O código deve respeitar estritamente a divisão `handler` -> `service` -> `repository` e `proxy` -> `service`.
- **Proibição de Acoplamento:** Nenhum handler HTTP ou função de transporte de proxy pode acessar diretamente mapas em memória, arquivos JSON ou estruturas do pacote `repository`. Toda consulta deve passar pela interface definida na camada `service`.
- **Padrão de Projeto:** Toda comunicação externa deve ser encapsulada no adaptador `Proxy` (`net/http/httputil.ReverseProxy`), sobrescrevendo o header `req.Host` para o host de destino e aplicando timeout máximo de 5 segundos (`RNF-02`).

## Comandos de Build, Lint e Teste
- Formatar: `gofmt -w .`
- Análise estática: `go vet ./...`
- Testes com detector de race condition: `go test -v -race ./...`
- Build: `go build -o apihub ./cmd/api`

## Convenções de Código
- Idioma de domínio, rotas e mensagens de erro de negócio: **Português (PT-BR)**.
- Nomes de pacotes Go curtos, em minúsculo e no singular (`handler`, `service`, `repository`, `proxy`, `domain`).
- Sempre propagar `context.Context` nas chamadas de rede e repositório.

## Proibições
- Nunca gravar chaves de API, senhas ou tokens reais no repositório.
- Nunca concluir uma implementação sem testes automatizados cobrindo o fluxo principal e pelo menos um caso de erro.
- Nunca adicionar bibliotecas externas de terceiros sem justificativa arquitetural (preferir a biblioteca padrão do Go).

## Responsabilidade
O agente propõe alterações; o integrante humano que abre e aprova o Pull Request responde integralmente pelo código integrado.