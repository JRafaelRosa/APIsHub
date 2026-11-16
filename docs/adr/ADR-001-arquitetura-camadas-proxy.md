# ADR-001: Adoção de Arquitetura em Camadas com Padrão Proxy para o APIHub

## Status
Aceito

## Contexto
O APIHub precisa atuar como ponto central de entrada para múltiplas APIs externas, mantendo um catálogo de rotas mascaradas e realizando o encaminhamento HTTP de baixa latência. O time possui 4 integrantes com perfis complementares trabalhando em sprints curtas de 2 semanas em Go, exigindo baixo acoplamento para permitir testes unitários isolados sem depender da internet ou de APIs reais durante o desenvolvimento.

## Decisão
Adotaremos uma **Arquitetura em Camadas** (inspirada em Sommerville, Fig. 6.8) organizada em quatro pacotes Go distintos com fluxo unidirecional de dependências:
1. **`handler` (Apresentação/HTTP):** Recebe requisições REST e de proxy, valida entradas e serializa respostas.
2. **`service` (Domínio/Regras de Negócio):** Valida unicidade de slugs, formata URLs mascaradas e coordena regras do catálogo.
3. **`repository` (Persistência):** Armazena e recupera os metadados das APIs (inicialmente em memória com proteção `sync.RWMutex` e persistência em arquivo JSON).
4. **`proxy` (Integração Externa):** Encapsula o encaminhamento HTTP reverso.

### Requisitos Não Funcionais Atendidos
- **RNF-01 (Baixa Latência ≤ 15 ms):** A separação clara permite manter o catálogo em estrutura indexada em memória (`map[string]APIRoute` com `RWMutex`) na camada de repositório, garantindo resolução de rota em tempo $O(1)$ sem I/O de disco a cada chamada de proxy.
- **RNF-03 (Portabilidade em Binário Único):** Ao isolar a camada de repositório por interface, evitamos dependências externas pesadas de banco de dados, mantendo o sistema compilável em um único binário Go.

### Desvantagem do Estilo Arquitetural e Mitigação
- **Desvantagem (Sommerville, p. 157):** Na arquitetura em camadas, a passagem de dados por múltiplas abstrações pode gerar custo extra de desempenho (overhead) e código repetitivo de mapeamento.
- **Mitigação:** Passaremos ponteiros de structs de domínio imutáveis (`*domain.APIRoute`) entre `repository`, `service` e `proxy` sem serializações intermediárias, além de reutilizar um único `http.Transport` com pool de conexões keep-alive na camada de proxy.

### Padrão de Projeto Pretendido
- **Padrão:** **Proxy** (Estrutural).
- **Problema:** Os clientes precisam consumir APIs externas heterogêneas sem conhecer suas URLs reais, cabeçalhos internos ou detalhes de conexão, além de exigir controle de timeout e ocultação de origem.
- **Consequência:** O componente `ReverseProxyAdapter` atua como substituto (proxy) das APIs externas, interceptando toda chamada para `/apihub/{slug}`, aplicando políticas de timeout e sanitização de headers antes de repassar ao servidor real.

## Consequências
- **Positivas:** Facilidade para aplicar TDD na camada `service` e `proxy` usando `httptest.NewServer`; divisão clara de tarefas entre os membros da equipe.
- **Negativas:** Exige disciplina para não pular camadas (ex.: handlers chamando repositórios diretamente).