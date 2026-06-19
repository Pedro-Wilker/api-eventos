# 🚀 API Eventos

Uma API RESTful de alta performance desenvolvida em **Go (Golang)** para o gerenciamento de listas de convidados e presentes.

A aplicação utiliza uma arquitetura *multitenancy* simples, onde cada usuário (cliente) gerencia apenas os seus próprios dados, garantindo segurança e isolamento total através de autenticação via tokens JWT.

---

## 📑 Sumário

* [Tecnologias Utilizadas](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-tecnologias-utilizadas)
* [Regras de Negócio e Segurança](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-regras-de-neg%C3%B3cio-e-seguran%C3%A7a-cors-e-roles)
* [Observações de Implementação](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-observa%C3%A7%C3%B5es-de-implementa%C3%A7%C3%A3o)
* [Como Executar o Projeto Localmente](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-como-executar-o-projeto-localmente)
* [Documentação das Rotas](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-documenta%C3%A7%C3%A3o-das-rotas-endpoints)
  * [Rotas Públicas](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-rotas-p%C3%BAblicas)
  * [Rotas Protegidas](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-rotas-protegidas-exigem-token-jwt)
* [Códigos de Status HTTP](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-c%C3%B3digos-de-status-http)
* [Licença](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-licen%C3%A7a)

---

## 🛠️ Tecnologias Utilizadas

| Categoria       | Tecnologia                                                   |
| --------------- | ------------------------------------------------------------ |
| Linguagem       | [Go](https://go.dev/)                                           |
| Framework Web   | [Gin Gonic](https://gin-gonic.com/)                             |
| ORM             | [GORM](https://gorm.io/)                                        |
| Banco de Dados  | PostgreSQL (com uso de `JSONB`para otimização de arrays) |
| Segurança      | JWT (JSON Web Tokens) e Bcrypt para hash de senhas           |
| Identificadores | UUID (`google/uuid`) para QR Code de convidados            |
| Testes          | Node.js (Fetch API nativa)                                   |

---

## 🔒 Regras de Negócio e Segurança (CORS e Roles)

* **CORS Dinâmico:** a API está configurada para aceitar requisições de origens específicas em produção, incluindo domínios principais e qualquer subdomínio derivado deles:
  * `*.casamentopedroelara.com.br`
  * `*.artonbyte.com.br`
  * `localhost` (para desenvolvimento)
* **Lógica de Administrador (Role: `admin`):** usuários cadastrados com o cargo `admin` podem listar, editar e deletar convidados e presentes de **qualquer** usuário do sistema. Usuários com cargo `client` estão isolados em seu próprio *tenant* (só acessam seus próprios dados). Essa mesma regra de isolamento se aplica ao **check-in** (validação de entrada) de convidados.
* **QR Code automático:** todo convidado criado recebe automaticamente um `qr_code` (UUID v4) gerado pelo backend no momento da criação (hook `BeforeCreate` do GORM). O frontend **não precisa — e não deve —** enviar esse campo.

---

## ⚠️ Observações de Implementação

Esta seção documenta comportamentos reais do backend que fogem do "óbvio" e que quem for consumir a API precisa saber.

### 1. Campos de resposta em `PascalCase` para Convidados e Presentes

Os models `Guest` e `Gift` embutem `gorm.Model` **sem sobrescrever as tags JSON** dos campos herdados (`ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`). Isso significa que, embora o **payload de entrada** use `snake_case` (`nome_presente`, `quantidade_acompanhante`...), os campos herdados do GORM voltam na resposta em `PascalCase`:

```json
{
  "ID": 3,
  "CreatedAt": "2026-06-20T01:31:03.082571432+02:00",
  "UpdatedAt": "2026-06-20T01:31:03.082571432+02:00",
  "DeletedAt": null,
  "Name": "Liquidificador",
  "Description": "Potência 1000W, Cor Preta",
  "ReservedQty": 0,
  "TotalReservations": 1,
  "UserID": 5,
  "PhotoURL": "https://link-da-imagem.com/foto.jpg"
}
```

> 🔧 **Sugestão de melhoria:** adicionar tags `json:"id"`, `json:"created_at"` etc. nos models para manter o padrão `snake_case` em toda a API. Até lá, o frontend precisa tratar os dois formatos (campos próprios em `snake_case`, campos herdados do GORM em `PascalCase`).

### 2. `PUT` em Presentes substitui o objeto inteiro (sem merge parcial)

`UpdateGift` atribui **todos** os campos do corpo da requisição diretamente no registro, sem checar se vieram vazios. Ou seja, **qualquer campo omitido no `PUT` é apagado** (vira string vazia ou zero), e não apenas ignorado. Sempre envie o objeto completo, igual a um `PUT` REST tradicional — não trate como `PATCH`.

### 3. Criação de Convidados pode falhar com `500` se a tabela não estiver migrada

O campo `relacoes_acompanhante` (`CompanionRelations`) foi adicionado ao model `Guest`, mas o GORM tenta gravar **todas** as colunas do struct no `INSERT`, mesmo que o campo não seja enviado no JSON. Se a coluna correspondente (`companion_relations`) ainda não existir na tabela em produção — porque o `AutoMigrate` não rodou de novo após esse campo ser adicionado ao código — **toda** criação de convidado falha com `500 Falha ao salvar convidado`, mesmo sem enviar o campo.

> ✅ **Correção necessária no backend:** reiniciar/redeployar a aplicação em produção para que o `AutoMigrate` crie a coluna que falta (ou rodar a migration manualmente no banco).

### 4. `Criar Presente` aceita um campo extra não documentado anteriormente

Além de `quantas_reservas_por_presente` (capacidade total), o endpoint aceita `quantidade_reservas`, que permite já criar o presente com reservas iniciais (`ReservedQty`).

### 5. `ReserveGift` não usa transação/lock

A reserva é feita com um `Find` seguido de `if` e `Save`, sem transação ou lock no banco. Em cenários de alta concorrência (vários convidados reservando o mesmo presente ao mesmo tempo), existe risco de *race condition* e o limite de reservas pode ser ultrapassado. Não é um problema de documentação, mas vale registrar como melhoria futura (ex.: usar `SELECT ... FOR UPDATE` ou uma transação com verificação atômica).

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

Instale a dependência de geração de UUID (usada para o QR Code dos convidados):

```bash
go get github.com/google/uuid
go mod tidy
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

> ⚠️ Sempre que um novo campo for adicionado a um model (`Guest`, `Gift`, `User`), reinicie a aplicação em produção para que o `AutoMigrate` crie a coluna correspondente — veja [Observações de Implementação, item 3](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-observa%C3%A7%C3%B5es-de-implementa%C3%A7%C3%A3o).

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

| Campo        | Tipo   | Obrigatório | Descrição                                                                                               |
| ------------ | ------ | :----------: | --------------------------------------------------------------------------------------------------------- |
| `name`     | string |      ✅      | Nome completo do usuário                                                                                 |
| `email`    | string |      ✅      | E-mail único utilizado para login                                                                        |
| `password` | string |      ✅      | Senha em texto puro (mínimo 6 caracteres, será hasheada)                                                |
| `role`     | string |      ❌      | `client`ou `admin`. Qualquer valor diferente de `admin`(inclusive vazio) é tratado como `client` |

**Resposta de sucesso (201):**

```json
{
  "message": "Usuário criado com sucesso!"
}
```

> ⚠️ O endpoint **não** retorna os dados do usuário criado (nem `id`, nem `email`) — apenas a mensagem de confirmação. Se precisar do `user_id`, faça login em seguida e consulte `GET /api/me`.

**Resposta de erro (400) — e-mail já cadastrado ou erro de validação:**

```json
{
  "error": "Email já cadastrado ou erro no banco"
}
```

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

**Resposta de erro (401) — e-mail não encontrado ou senha incorreta:**

```json
{
  "error": "Email ou senha incorretos"
}
```

---

#### Listar Presentes do Cliente

Retorna a lista de presentes de um cliente específico. Ideal para renderizar a página pública do convite.

* **Método:** `GET`
* **Endpoint:** `/api/public/presentes/:user_id`

**Resposta de sucesso (200):**

```json
{
  "data": [
    {
      "ID": 3,
      "CreatedAt": "2026-06-20T01:31:03.082571+02:00",
      "UpdatedAt": "2026-06-20T01:31:03.082571+02:00",
      "DeletedAt": null,
      "Name": "Liquidificador",
      "Description": "Potência 1000W, Cor Preta",
      "ReservedQty": 0,
      "TotalReservations": 1,
      "UserID": 5,
      "PhotoURL": "https://link-da-imagem.com/foto.jpg"
    }
  ]
}
```

> ℹ️ Os campos vêm em `PascalCase` por herdarem de `gorm.Model` sem tags JSON customizadas — veja [Observações de Implementação](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-observa%C3%A7%C3%B5es-de-implementa%C3%A7%C3%A3o).

---

#### Reservar Presente

Incrementa a quantidade reservada (`ReservedQty`) de um presente, se houver disponibilidade (`ReservedQty < TotalReservations`).

* **Método:** `POST`
* **Endpoint:** `/api/public/presentes/:id/reservar`

**Resposta de sucesso (200):**

```json
{
  "message": "Presente reservado com sucesso!",
  "data": {
    "ID": 3,
    "Name": "Liquidificador",
    "ReservedQty": 1,
    "TotalReservations": 1,
    "UserID": 5,
    "PhotoURL": "https://link-da-imagem.com/foto.jpg"
  }
}
```

**Resposta de erro (404) — presente não encontrado:**

```json
{
  "error": "Presente não encontrado"
}
```

**Resposta de conflito (409) — presente esgotado:**

```json
{
  "error": "Este presente já foi totalmente reservado!"
}
```

> ⚠️ Não há transação/lock nessa operação — em cenários de alta concorrência existe risco teórico de duas reservas simultâneas ultrapassarem o limite. Veja [Observações de Implementação, item 5](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-observa%C3%A7%C3%B5es-de-implementa%C3%A7%C3%A3o).

---

### 🛡️ Rotas Protegidas (Exigem Token JWT)

#### Perfil do Usuário

Valida o token e retorna os dados de autorização do usuário logado.

* **Método:** `GET`
* **Endpoint:** `/api/me`

**Resposta de sucesso (200):**

```json
{
  "message": "Você está autenticado!",
  "user_id": 1,
  "role": "client"
}
```

---

#### Criar Convidado

Adiciona um novo convidado vinculado ao usuário autenticado. O `qr_code` é gerado automaticamente pelo backend —  **não envie esse campo** .

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
  "numeros_acompanhantes": ["11888888888"],
  "relacoes_acompanhante": ["Esposa", "Filho"]
}
```

| Campo                       | Tipo     | Obrigatório | Descrição                                                   |
| --------------------------- | -------- | :----------: | ------------------------------------------------------------- |
| `nome`                    | string   |      ✅      | Nome do convidado principal                                   |
| `quantidade_acompanhante` | int      |      ❌      | Número de acompanhantes                                      |
| `nome_acompanhante`       | string[] |      ❌      | Lista com os nomes dos acompanhantes                          |
| `email_convidado`         | string   |      ❌      | E-mail do convidado principal                                 |
| `numero_convidado`        | string   |      ❌      | Telefone do convidado principal                               |
| `emails_acompanhantes`    | string[] |      ❌      | Lista de e-mails dos acompanhantes                            |
| `numeros_acompanhantes`   | string[] |      ❌      | Lista de telefones dos acompanhantes                          |
| `relacoes_acompanhante`   | string[] |      ❌      | Relação de cada acompanhante (ex.:`"Esposa"`,`"Filho"`) |

> ℹ️ Apenas o campo `nome` é obrigatório; os demais são opcionais.
>
> ⚠️ Não envie o campo `qr_code` — ele é gerado automaticamente pelo backend.
>
> 🐛 **Bug conhecido em produção:** esse endpoint pode retornar `500 Falha ao salvar convidado` para **qualquer** convidado (mesmo sem enviar `relacoes_acompanhante`) se a tabela do banco ainda não tiver a coluna correspondente migrada. Veja [Observações de Implementação, item 3](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-observa%C3%A7%C3%B5es-de-implementa%C3%A7%C3%A3o).

**Resposta de sucesso (201):**

```json
{
  "message": "Convidado adicionado com sucesso!",
  "data": {
    "ID": 5,
    "CreatedAt": "2026-06-20T01:31:03.082571+02:00",
    "UpdatedAt": "2026-06-20T01:31:03.082571+02:00",
    "DeletedAt": null,
    "nome": "Carlos Almeida",
    "user_id": 1,
    "quantidade_acompanhante": 2,
    "nome_acompanhante": ["Ana Almeida", "Pedro Almeida"],
    "emails_acompanhantes": ["ana@email.com"],
    "numeros_acompanhantes": ["11888888888"],
    "relacoes_acompanhante": ["Esposa", "Filho"],
    "email_convidado": "carlos@email.com",
    "numero_convidado": "11999999999",
    "qr_code": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
    "entrada_registrada": false,
    "data_entrada": null,
    "usuario_validador": null
  }
}
```

**Resposta de erro (500) — falha ao salvar (ver bug conhecido acima):**

```json
{
  "error": "Falha ao salvar convidado"
}
```

---

#### Listar Convidados

Retorna todos os convidados que pertencem ao usuário autenticado. Se o usuário autenticado for `admin`, retorna os convidados de **todos** os usuários do sistema.

* **Método:** `GET`
* **Endpoint:** `/api/convidados`

**Resposta de sucesso (200):**

```json
{
  "data": [
    {
      "ID": 5,
      "nome": "Carlos Almeida",
      "user_id": 1,
      "quantidade_acompanhante": 2,
      "nome_acompanhante": ["Ana Almeida", "Pedro Almeida"],
      "email_convidado": "carlos@email.com",
      "numero_convidado": "11999999999",
      "qr_code": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
      "entrada_registrada": false,
      "data_entrada": null,
      "usuario_validador": null
    }
  ]
}
```

---

#### Atualizar Convidado

Edita as informações de um convidado existente. **Substitui o objeto inteiro** — campos omitidos voltam ao valor zero (mesma observação do `PUT` de presentes).

* **Método:** `PUT`
* **Endpoint:** `/api/convidados/:id`

**Corpo da requisição (JSON):** *(mesmo formato da criação)*

```json
{
  "nome": "Nome Atualizado",
  "quantidade_acompanhante": 3
}
```

**Resposta de sucesso (200):**

```json
{
  "message": "Convidado atualizado!",
  "data": { ... }
}
```

**Resposta de erro (404) — convidado não encontrado ou não pertence ao usuário:**

```json
{
  "error": "Convidado não encontrado ou sem permissão"
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

#### Buscar Convidado por Código (QR Code ou ID)

Retorna os dados de um convidado a partir do seu ID numérico ou do seu `qr_code` (UUID), sem alterar nada. Útil para pré-visualizar o convidado antes de confirmar o check-in (ex.: ao ler o QR Code na entrada do evento).

* **Método:** `GET`
* **Endpoint:** `/api/convidados/buscar/:codigo`

**Resposta de sucesso (200):**

```json
{
  "data": {
    "ID": 5,
    "nome": "Carlos Almeida",
    "qr_code": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
    "entrada_registrada": false,
    "data_entrada": null
  }
}
```

**Resposta de erro (404):**

```json
{
  "error": "Convidado não encontrado no sistema."
}
```

---

#### Check-in (Validar Entrada do Convidado)

Registra a entrada de um convidado no evento a partir do seu QR Code (ou ID). **Idempotente:** se o convidado já tiver entrada registrada, retorna conflito ao invés de duplicar o check-in.

* **Método:** `POST`
* **Endpoint:** `/api/convidados/checkin`

**Corpo da requisição (JSON):**

```json
{
  "codigo": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"
}
```

> `codigo` aceita tanto o `qr_code` (UUID) quanto o `id` numérico do convidado. O campo é  **obrigatório** .

**Resposta de erro (400) — campo `codigo` ausente:**

```json
{
  "error": "Informe o código do QR Code ou ID do convidado."
}
```

**Resposta de sucesso (200):**

```json
{
  "status": "valido",
  "mensagem": "Entrada autorizada com sucesso!",
  "data": {
    "ID": 5,
    "nome": "Carlos Almeida",
    "entrada_registrada": true,
    "data_entrada": "2026-06-19T20:15:00Z",
    "usuario_validador": 1
  }
}
```

**Resposta de conflito (409) — entrada já registrada anteriormente:**

```json
{
  "status": "duplicado",
  "mensagem": "Entrada já registrada em 19/06/2026 20:15",
  "data": { ... }
}
```

**Resposta de erro (404) — convidado não encontrado:**

```json
{
  "status": "invalido",
  "error": "Convidado não encontrado no sistema."
}
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
  "quantidade_reservas": 0,
  "quantas_reservas_por_presente": 1,
  "foto_presente": "https://link-da-imagem.com/foto.jpg"
}
```

| Campo                             | Tipo   | Obrigatório | Descrição                                                                                                     |
| --------------------------------- | ------ | :----------: | --------------------------------------------------------------------------------------------------------------- |
| `nome_presente`                 | string |      ✅      | Nome do presente                                                                                                |
| `descricao_presente`            | string |      ❌      | Detalhes adicionais (cor, potência, etc.)                                                                      |
| `quantidade_reservas`           | int    |      ❌      | Quantidade de reservas**já existentes**ao criar o presente (raramente usado fora de migração de dados) |
| `quantas_reservas_por_presente` | int    |      ❌      | Capacidade total de reservas. Se omitido ou ≤ 0, assume `1`                                                  |
| `foto_presente`                 | string |      ❌      | URL da imagem do presente                                                                                       |

**Resposta de sucesso (201):**

```json
{
  "message": "Presente adicionado com sucesso!",
  "data": {
    "ID": 3,
    "CreatedAt": "2026-06-20T01:31:03.082571432+02:00",
    "UpdatedAt": "2026-06-20T01:31:03.082571432+02:00",
    "DeletedAt": null,
    "Name": "Liquidificador",
    "Description": "Potência 1000W, Cor Preta",
    "ReservedQty": 0,
    "TotalReservations": 1,
    "UserID": 1,
    "PhotoURL": "https://link-da-imagem.com/foto.jpg"
  }
}
```

---

#### Listar Presentes

Retorna todos os presentes cadastrados pelo usuário autenticado. Se o usuário for `admin`, retorna os presentes de **todos** os usuários.

* **Método:** `GET`
* **Endpoint:** `/api/presentes`

**Resposta de sucesso (200):**

```json
{
  "data": [
    {
      "ID": 3,
      "Name": "Liquidificador",
      "Description": "Potência 1000W, Cor Preta",
      "ReservedQty": 0,
      "TotalReservations": 1,
      "UserID": 1,
      "PhotoURL": "https://link-da-imagem.com/foto.jpg"
    }
  ]
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
  "descricao_presente": "Potência 1000W, Cor Preta",
  "quantas_reservas_por_presente": 2,
  "foto_presente": "https://link-da-imagem.com/foto.jpg"
}
```

> ⚠️ **Atenção:** este `PUT` substitui o registro inteiro. Qualquer campo omitido (ex.: `descricao_presente`, `foto_presente`) é apagado, não mantido. Sempre envie o objeto completo. Veja [Observações de Implementação, item 2](https://claude.ai/chat/9e04e472-6a21-46c2-8370-1f06e5d2b0d4#-observa%C3%A7%C3%B5es-de-implementa%C3%A7%C3%A3o).

**Resposta de sucesso (200):**

```json
{
  "message": "Presente atualizado!",
  "data": {
    "ID": 3,
    "Name": "Liquidificador Atualizado",
    "Description": "Potência 1000W, Cor Preta",
    "ReservedQty": 0,
    "TotalReservations": 2,
    "UserID": 1,
    "PhotoURL": "https://link-da-imagem.com/foto.jpg"
  }
}
```

**Resposta de erro (404) — presente não encontrado ou não pertence ao usuário:**

```json
{
  "error": "Presente não encontrado ou sem permissão"
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

| Código | Significado                                                                     |
| ------- | ------------------------------------------------------------------------------- |
| `200` | Requisição bem-sucedida                                                       |
| `201` | Recurso criado com sucesso                                                      |
| `400` | Requisição inválida (campos obrigatórios ausentes/inválidos)               |
| `401` | Não autenticado (token ausente, inválido, expirado ou credenciais incorretas) |
| `403` | Acesso negado ao recurso                                                        |
| `404` | Recurso não encontrado                                                         |
| `409` | Conflito (ex.: presente já esgotado/totalmente reservado, check-in duplicado)  |
| `500` | Erro interno do servidor (ex.: falha ao salvar/deletar registro no banco)       |

---

## 📄 Licença

Este projeto está sob a licença MIT — veja o arquivo [LICENSE](https://claude.ai/chat/LICENSE) para mais detalhes. Sinta-se livre para usar, modificar e distribuir.

---

*Desenvolvido com dedicação por Pedro Wilker.*
