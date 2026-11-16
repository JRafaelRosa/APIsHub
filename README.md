# APIHub — Centralizador e Mascarador de Rotas de API

## Tema e Escopo
O **APIHub** é um gateway centralizador destinado a desenvolvedores e equipes de software que consomem múltiplas APIs externas dispersas. O sistema resolve a fragmentação de endpoints ao catalogar serviços externos (com descrição, finalidade e URL original) e expô-los sob rotas internas unificadas e mascaradas (ex.: `/apihub/pagamento/pix` e `/apihub/clima/regiao`). Não entram no escopo desta versão: cobrança/monetização de chamadas, balanceamento de carga multi-região e emissão de certificados SSL customizados.

## Linguagem e Stack
- **Linguagem:** Go (Golang 1.22+) — escolhida pela baixa latência de I/O, concorrência nativa via goroutines e suporte nativo robusto a proxy reverso HTTP (`net/http/httputil`).

## Equipe

| Nome Completo | GitHub | O que cobre bem | Papel | Justificativa da Escolha |
| :--- | :--- | :--- | :--- | :--- |
| Pessoa 1 | @usuario1 | Programar, negociar | Arquiteto Técnico | Possui experiência prévia com desenvolvimento backend e definição de arquiteturas de serviços de rede. |
| Pessoa 2 | @usuario2 | Negociar, documentar | Product Owner | Tem facilidade em traduzir necessidades de uso em requisitos claros e priorizar entregas de valor. |
| Pessoa 3 | @usuario3 | Programar, testar | Scrum Master | Organiza o fluxo de tarefas do time, acompanha impedimentos e garante a disciplina de testes e prazos. |
| Pessoa 4 | @usuario4 | Testar, documentar | Desenvolvedor / QA | Perfil analítico focado na escrita de casos de teste automatizados e manutenção da documentação técnica. |

*(Nota: Todas as 4 competências — negociar, programar, testar e documentar — estão cobertas de forma redundante pelos integrantes do grupo).*

## Definição de Pronto
Para que qualquer história ou tarefa seja considerada concluída e integrada à branch principal (`main`), ela deve cumprir todos os itens abaixo:
1. O código compila sem erros ou warnings em ambiente Go padrão (`go build ./...`).
2. Todos os testes automatizados (unidade e integração) passam (`go test -race ./...`).
3. A análise estática e formatação de código passam sem alertas no CI (`gofmt` e `go vet ./...`).
4. Todos os critérios de aceitação definidos na história foram verificados e atendidos.
5. O Pull Request foi revisado e aprovado manualmente por pelo menos um integrante humano diferente do autor.
6. A estrutura do código respeita rigorosamente a separação de camadas definida no `ADR-001`.
7. Qualquer uso de agente de IA na tarefa foi registrado no arquivo `docs/registro-ia.md`.