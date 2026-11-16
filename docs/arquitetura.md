# Arquitetura do Sistema — APIHub

## Diagrama de Blocos

```mermaid
flowchart LR
    Cliente[Cliente HTTP / Browser]

    subgraph APIHub [Aplicação Go - APIHub]
        direction TB
        Handler[Camada Handler\nRotas REST & /apihub/*]
        Service[Camada Service\nValidação & Resolução de Slug]
        Repo[(Camada Repository\nCatálogo Thread-Safe + JSON)]
        Proxy[Camada Proxy\nReverseProxy + Sanitização]

        Handler --> Service
        Service --> Repo
        Handler --> Proxy
        Proxy --> Service
    end

    Ext1[API Externa 1\nEx: Pagamento PIX]
    Ext2[API Externa 2\nEx: Clima Região]

    Cliente -- "1. GET /apihub/clima/regiao" --> Handler
    Proxy -- "2. GET [https://api.clima.ext/v1](https://api.clima.ext/v1)" --> Ext2
    Ext2 -- "3. Resposta JSON" --> Proxy
    Handler -- "4. Resposta Mascarada" --> Cliente