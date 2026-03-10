# crd-study

Projeto de estudo para compreender como criar e utilizar **Custom Resource Definitions (CRDs)** no **Kubernetes** utilizando **Go**, **client-go** e **Cobra CLI**.

Este projeto demonstra na prática como:

- criar um **CRD no Kubernetes**
- interagir com a API do Kubernetes usando **client-go**
- utilizar o **dynamic client** para recursos customizados
- observar mudanças em recursos com **watch**
- separar corretamente **spec** e **status**
- implementar uma **CLI em Go**
- aguardar o estado final de um recurso com um comando `wait`

## Pastas principais

### `cmd/`

Contém os comandos da CLI:

- `listpods`
- `list clusterrestores`
- `wait`

Esses comandos permitem interagir com o cluster Kubernetes diretamente da CLI.

---

### `internal/`

Implementação interna da aplicação:

- criação de clientes Kubernetes
- dynamic client
- lógica de watch
- leitura e validação de YAML

---

### `pkg/apis/gdch/v1/`

Define o **modelo do CRD em Go**, incluindo:

- `ClusterRestore`
- `ClusterRestoreSpec`
- `ClusterRestoreStatus`
- `ClusterRestoreList`

Essas estruturas representam o recurso customizado dentro do Kubernetes.

---

### `config/crd/bases/`

Contém o **manifesto do CRD** que é aplicado no cluster.

Esse arquivo registra o recurso `ClusterRestore` na API do Kubernetes.

---

### `examples/`

Contém exemplos de objetos `ClusterRestore` para testes.

---

### `hack/`

Scripts auxiliares do projeto (ex.: geração de código).

---

# Tecnologias Utilizadas

### Go

Linguagem principal do projeto.

---

### Kubernetes

Plataforma que está sendo estendida através de CRDs.

---

### client-go

Cliente oficial do Kubernetes para aplicações em Go.

Permite interagir diretamente com o **Kubernetes API Server**.

---

### Cobra CLI

Framework utilizado para construir a interface de linha de comando.

---

### Kind

Ferramenta que cria clusters Kubernetes locais usando Docker.

Ideal para desenvolvimento e testes.

---

### Bazel + Gazelle

Ferramentas de build que permitem:

- builds reproduzíveis
- gerenciamento avançado de dependências
- geração automática de arquivos Bazel para projetos Go

---

# Pré-requisitos

Antes de rodar o projeto, instale:

- Go
- Docker
- Kind
- kubectl

Verifique se estão instalados:

```bash
go version
docker info
kind version
kubectl version --client

