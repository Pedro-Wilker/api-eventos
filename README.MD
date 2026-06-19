# 🚀 API Eventos

Uma API RESTful de alta performance desenvolvida em **Go (Golang)** para o gerenciamento de listas de convidados e presentes.

A aplicação utiliza uma arquitetura *multitenancy* simples, onde cada usuário (cliente) gerencia apenas os seus próprios dados, garantindo segurança e isolamento total através de autenticação via tokens JWT.

---

## 📑 Sumário

* [Tecnologias Utilizadas](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-tecnologias-utilizadas)
* [Regras de Negócio e Segurança](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-regras-de-neg%C3%B3cio-e-seguran%C3%A7a-cors-e-roles)
* [Como Executar o Projeto Localmente](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-como-executar-o-projeto-localmente)
* [Documentação das Rotas](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-documenta%C3%A7%C3%A3o-das-rotas-endpoints)
  * [Rotas Públicas](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-rotas-p%C3%BAblicas)
  * [Rotas Protegidas](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-rotas-protegidas-exigem-token-jwt)
* [Códigos de Status HTTP](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-c%C3%B3digos-de-status-http)
* [Licença](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-licen%C3%A7a)

---

## 🛠️ Tecnologias Utilizadas

| Categoria      | Tecnologia                                                   |
| -------------- | ------------------------------------------------------------ |
| Linguagem      | [Go](https://go.dev/)                                           |
| Framework Web  | [Gin Gonic](https://gin-gonic.com/)                             |
| ORM            | [GORM](https://gorm.io/)                                        |
| Banco de Dados | PostgreSQL (com uso de `JSONB`para otimização de arrays) |
| Segurança     | JWT (JSON Web Tokens) e Bcrypt para hash de senhas           |
| Testes         | Node.js (Fetch API nativa)                                   |

---

## 🔒 Regras de Negócio e Segurança (CORS e Roles)

* **CORS Dinâmico:** a API está configurada para aceitar requisições de origens específicas em produção, incluindo domínios principais e qualquer subdomínio derivado deles:
  * `*.casamentopedroelara.com.br`
  * `*.artonbyte.com.br`
  * `localhost` (para desenvolvimento)
* **Lógica de Administrador (Role: `admin`):** usuários cadastrados com o cargo `admin` podem listar, editar e deletar convidados e presentes de **qualquer** usuário do sistema. Usuários com cargo `client` estão isolados em seu próprio *tenant* (só acessam seus próprios dados).

---

## ⚙️ Como Executar o Projeto Localmente

### 1. Pré-requisitos

* [Go](https://go.dev/dl/) instalado (versão 1.21+ recomendada)
* [PostgreSQL](https://www.postgresql.org/) rodando localmente
* [Node.js](https://nodejs.org/) versão 18+ (apenas para rodar o script de testes)

### 2. Configuração do Ambiente

Clone o repositório:

```bash
git clone https://github.com/seu-usuario/api-eventos.git
cd api-eventos
```

Crie um arquivo `.env` na raiz do projeto com as suas credenciais do banco de dados e uma chave secreta para o JWT:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=sua_senha
DB_NAME=eventos_db
DB_PORT=5432
API_SECRET=sua_chave_secreta_super_segura
```

> ⚠️ **Importante:** nunca cometa o arquivo `.env` no repositório. Adicione-o ao `.gitignore`.

### 3. Rodando a Aplicação

O banco de dados e as tabelas serão criados/sincronizados automaticamente na primeira execução através do `AutoMigrate` do GORM.

```bash
# Baixe as dependências
go mod tidy

# Inicie o servidor
go run main.go
```

O servidor estará rodando em `http://localhost:8080`.

### 4. Rodando os Testes Automatizados

Com o servidor Go rodando, abra um novo terminal e execute o script de testes:

```bash
node test-api.js
```

---

## 📖 Documentação das Rotas (Endpoints)

Abaixo estão todas as rotas disponíveis na API.

> 🔐 **Atenção:** todas as rotas na seção "Protegidas" exigem o envio do cabeçalho `Authorization` no formato:
>
> ```
> Authorization: Bearer <SEU_TOKEN_JWT>
> ```

### 🔓 Rotas Públicas

#### Health Check

Verifica se a API está online.

* **Método:** `GET`
* **Endpoint:** `/ping`

**Resposta de sucesso (200):**

```json
{
  "message": "API online e pronta para o trabalho!"
}
```

---

#### Registrar Usuário

Cria um novo usuário (cliente ou admin) no sistema.

* **Método:** `POST`
* **Endpoint:** `/api/register`

**Corpo da requisição (JSON):**

```json
{
  "name": "João Silva",
  "email": "joao@email.com",
  "password": "senha_segura",
  "role": "client"
}
```

**Resposta de sucesso (201):**

```json
{
  "id": 1,
  "name": "João Silva",
  "email": "joao@email.com",
  "role": "client"
}
```

| Campo        | Tipo   | Obrigatório | Descrição                          |
| ------------ | ------ | :----------: | ------------------------------------ |
| `name`     | string |      ✅      | Nome completo do usuário            |
| `email`    | string |      ✅      | E-mail único utilizado para login   |
| `password` | string |      ✅      | Senha em texto puro (será hasheada) |
| `role`     | string |      ✅      | `client`ou `admin`               |

---

#### Login

Autentica o usuário e devolve o token JWT.

* **Método:** `POST`
* **Endpoint:** `/api/login`

**Corpo da requisição (JSON):**

```json
{
  "email": "joao@email.com",
  "password": "senha_segura"
}
```

**Resposta de sucesso (200):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI..."
}
```

**Resposta de erro (401):**

```json
{
  "error": "credenciais inválidas"
}
```

---

#### Listar Presentes do Cliente

Retorna a lista de presentes de um cliente específico. Ideal para renderizar a página pública do convite.

* **Método:** `GET`
* **Endpoint:** `/api/public/presentes/:user_id`

**Resposta de sucesso (200):**

```json
[
  {
    "id": 3,
    "nome_presente": "Liquidificador",
    "descricao_presente": "Potência 1000W, Cor Preta",
    "quantas_reservas_por_presente": 1,
    "foto_presente": "https://link-da-imagem.com/foto.jpg"
  }
]
```

---

#### Reservar Presente

Incrementa a quantidade reservada de um presente, se houver disponibilidade.

* **Método:** `POST`
* **Endpoint:** `/api/public/presentes/:id/reservar`

**Resposta de sucesso (200):**

```json
{
  "message": "Presente reservado com sucesso!",
  "data": { ... }
}
```

**Resposta de conflito (409) — presente esgotado:**

```json
{
  "error": "Este presente já foi totalmente reservado!"
}
```

---

### 🛡️ Rotas Protegidas (Exigem Token JWT)

#### Perfil do Usuário

Valida o token e retorna os dados de autorização do usuário logado.

* **Método:** `GET`
* **Endpoint:** `/api/me`

**Resposta de sucesso (200):**

```json
{
  "id": 1,
  "name": "João Silva",
  "email": "joao@email.com",
  "role": "client"
}
```

---

#### Criar Convidado

Adiciona um novo convidado vinculado ao usuário autenticado.

* **Método:** `POST`
* **Endpoint:** `/api/convidados`

**Corpo da requisição (JSON):**

```json
{
  "nome": "Carlos Almeida",
  "quantidade_acompanhante": 2,
  "nome_acompanhante": ["Ana Almeida", "Pedro Almeida"],
  "email_convidado": "carlos@email.com",
  "numero_convidado": "11999999999",
  "emails_acompanhantes": ["ana@email.com"],
  "numeros_acompanhantes": ["11888888888"]
}
```

| Campo                       | Tipo     | Obrigatório | Descrição                          |
| --------------------------- | -------- | :----------: | ------------------------------------ |
| `nome`                    | string   |      ✅      | Nome do convidado principal          |
| `quantidade_acompanhante` | int      |      ❌      | Número de acompanhantes             |
| `nome_acompanhante`       | string[] |      ❌      | Lista com os nomes dos acompanhantes |
| `email_convidado`         | string   |      ❌      | E-mail do convidado principal        |
| `numero_convidado`        | string   |      ❌      | Telefone do convidado principal      |
| `emails_acompanhantes`    | string[] |      ❌      | Lista de e-mails dos acompanhantes   |
| `numeros_acompanhantes`   | string[] |      ❌      | Lista de telefones dos acompanhantes |

> ℹ️ Apenas o campo `nome` é obrigatório; os demais são opcionais.

**Resposta de sucesso (201):**

```json
{
  "id": 5,
  "nome": "Carlos Almeida",
  "quantidade_acompanhante": 2,
  "nome_acompanhante": ["Ana Almeida", "Pedro Almeida"],
  "email_convidado": "carlos@email.com",
  "numero_convidado": "11999999999",
  "emails_acompanhantes": ["ana@email.com"],
  "numeros_acompanhantes": ["11888888888"],
  "user_id": 1
}
```

---

#### Listar Convidados

Retorna todos os convidados que pertencem ao usuário autenticado.

* **Método:** `GET`
* **Endpoint:** `/api/convidados`

**Resposta de sucesso (200):**

```json
[
  {
    "id": 5,
    "nome": "Carlos Almeida",
    "quantidade_acompanhante": 2,
    "nome_acompanhante": ["Ana Almeida", "Pedro Almeida"],
    "email_convidado": "carlos@email.com",
    "numero_convidado": "11999999999"
  }
]
```

---

#### Criar Presente

Adiciona um presente à lista do usuário autenticado.

* **Método:** `POST`
* **Endpoint:** `/api/presentes`

**Corpo da requisição (JSON):**

```json
{
  "nome_presente": "Liquidificador",
  "descricao_presente": "Potência 1000W, Cor Preta",
  "quantas_reservas_por_presente": 1,
  "foto_presente": "https://link-da-imagem.com/foto.jpg"
}
```

| Campo                             | Tipo   | Obrigatório | Descrição                                |
| --------------------------------- | ------ | :----------: | ------------------------------------------ |
| `nome_presente`                 | string |      ✅      | Nome do presente                           |
| `descricao_presente`            | string |      ❌      | Detalhes adicionais (cor, potência, etc.) |
| `quantas_reservas_por_presente` | int    |      ❌      | Quantidade disponível para reserva        |
| `foto_presente`                 | string |      ❌      | URL da imagem do presente                  |

**Resposta de sucesso (201):**

```json
{
  "id": 3,
  "nome_presente": "Liquidificador",
  "descricao_presente": "Potência 1000W, Cor Preta",
  "quantas_reservas_por_presente": 1,
  "foto_presente": "https://link-da-imagem.com/foto.jpg",
  "user_id": 1
}
```

---

#### Listar Presentes

Retorna todos os presentes cadastrados pelo usuário autenticado.

* **Método:** `GET`
* **Endpoint:** `/api/presentes`

**Resposta de sucesso (200):**

```json
[
  {
    "id": 3,
    "nome_presente": "Liquidificador",
    "descricao_presente": "Potência 1000W, Cor Preta",
    "quantas_reservas_por_presente": 1,
    "foto_presente": "https://link-da-imagem.com/foto.jpg"
  }
]
```

---

#### Atualizar Convidado

Edita as informações de um convidado existente.

* **Método:** `PUT`
* **Endpoint:** `/api/convidados/:id`

**Corpo da requisição (JSON):** *(mesmo formato da criação)*

```json
{
  "nome": "Nome Atualizado",
  "quantidade_acompanhante": 3
}
```

---

#### Deletar Convidado

Remove (soft delete) um convidado da lista.

* **Método:** `DELETE`
* **Endpoint:** `/api/convidados/:id`

**Resposta de sucesso (200):**

```json
{
  "message": "Convidado deletado com sucesso!"
}
```

---

#### Atualizar Presente

Edita as informações ou a disponibilidade de um presente existente.

* **Método:** `PUT`
* **Endpoint:** `/api/presentes/:id`

**Corpo da requisição (JSON):** *(mesmo formato da criação)*

```json
{
  "nome_presente": "Liquidificador Atualizado",
  "quantas_reservas_por_presente": 2
}
```

---

#### Deletar Presente

Remove (soft delete) um presente da lista.

* **Método:** `DELETE`
* **Endpoint:** `/api/presentes/:id`

**Resposta de sucesso (200):**

```json
{
  "message": "Presente deletado com sucesso!"
}
```

---

## 📋 Códigos de Status HTTP

| Código | Significado                                                       |
| ------- | ----------------------------------------------------------------- |
| `200` | Requisição bem-sucedida                                         |
| `201` | Recurso criado com sucesso                                        |
| `400` | Requisição inválida (campos obrigatórios ausentes/inválidos) |
| `401` | Não autenticado (token ausente, inválido ou expirado)           |
| `403` | Acesso negado ao recurso                                          |
| `404` | Recurso não encontrado                                           |
| `409` | Conflito (ex: presente já esgotado/totalmente reservado)         |
| `500` | Erro interno do servidor                                          |

---

## 📄 Licença

Este projeto está sob a licença MIT — veja o arquivo [LICENSE](https://claude.ai/chat/LICENSE) para mais detalhes. Sinta-se livre para usar, modificar e distribuir.

---

*Desenvolvido com dedicação por Pedro Wilker.*
