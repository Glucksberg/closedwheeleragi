# 📱 Telegram Chat Integration - Complete Guide

**Date**: 2026-02-08
**Status**: ✅ **FULLY FUNCTIONAL**

---

## 🎯 Overview

Você agora pode **conversar diretamente** com o ClosedWheelerAGI via Telegram, como se estivesse usando a TUI local!

### Funcionalidades
- ✅ Chat completo com o AGI via Telegram
- ✅ Execução de ferramentas remotamente
- ✅ Aprovações de operações sensíveis
- ✅ Comandos administrativos (/status, /logs, /diff)
- ✅ Respostas longas divididas automaticamente
- ✅ Notificações de progresso

---

## 🚀 Setup Rápido

### 1. Obter Bot Token

```bash
# 1. Abra o Telegram e encontre @BotFather
# 2. Envie: /newbot
# 3. Escolha um nome: "MeuAGI Bot"
# 4. Escolha um username: "meu_agi_bot"
# 5. Copie o token: 1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
```

### 2. Configurar no .env

```bash
# Adicione ao .env:
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
```

### 3. Habilitar no config.json

```json
{
  "telegram": {
    "enabled": true,
    "bot_token": "",  // Lido do .env
    "chat_id": 0,     // Será configurado no próximo passo
    "notify_on_tool_start": true
  }
}
```

### 4. Obter seu Chat ID

```bash
# 1. Inicie o ClosedWheeler
./ClosedWheeler

# 2. Abra o Telegram e encontre seu bot
# 3. Envie: /start

# O bot responderá:
# "👋 Olá! Seu Chat ID é: 123456789"
```

### 5. Configurar Chat ID

```json
{
  "telegram": {
    "enabled": true,
    "bot_token": "",
    "chat_id": 123456789,  // ← Cole aqui
    "notify_on_tool_start": true
  }
}
```

### 6. Reinicie o Agent

```bash
# Pare e reinicie
Ctrl+C
./ClosedWheeler

# Agora você está conectado! 🎉
```

---

## 💬 Como Usar

### Conversação Normal

Simplesmente envie mensagens para o bot como se estivesse conversando com o AGI localmente:

**Você:** `Analise o arquivo main.go e me diga o que ele faz`

**AGI:**
```
💭 Pensando...

📝 Analisando o arquivo main.go...

O arquivo main.go é o ponto de entrada da aplicação ClosedWheelerAGI.

Principais funcionalidades:
1. Parse de flags de linha de comando
2. Carregamento de configuração
3. Inicialização do agente
4. Setup do Telegram (se habilitado)
5. Execução da TUI

O fluxo principal:
- Verifica se há API key configurada
- Se não houver, executa setup interativo
- Cria instância do agent
- Inicia polling do Telegram
- Roda a interface TUI
```

**Você:** `Crie uma função para validar email em Go`

**AGI:**
```
💭 Pensando...

🔧 Criando função de validação de email...

Criei a função ValidateEmail no arquivo utils/validator.go:

```go
func ValidateEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}
```

A função usa regex para validar o formato do email.

Deseja que eu adicione testes unitários também?
```

### Comandos Administrativos

#### `/start`
Mostra informações iniciais e seu Chat ID.

```
/start

👋 Olá! Bem-vindo ao ClosedWheelerAGI

Seu Chat ID: 123456789

Configure este ID no config.json para ativar o controle remoto.
```

#### `/help`
Lista todos os comandos disponíveis.

```
/help

🤖 ClosedWheelerAGI - Comandos Telegram

Comandos Disponíveis:
/start - Informações iniciais
/help - Esta mensagem
/status - Status da memória
/logs - Últimos logs
/diff - Git diff

Conversação:
Envie qualquer mensagem sem "/" para conversar!
```

#### `/status`
Mostra status da memória e projeto.

```
/status

📊 AGI Status

Memory: STM: 5 │ WM: 12 │ LTM: 45
Project: ClosedWheelerAGI (27 files, Go)
```

#### `/logs`
Mostra últimos logs do sistema.

```
/logs

📜 Últimos Logs:
[2026-02-08 10:30:45] INFO: Agent started
[2026-02-08 10:30:46] INFO: Telegram connected
[2026-02-08 10:31:00] INFO: Tool call: read_file
[2026-02-08 10:31:02] INFO: Response generated
```

#### `/diff`
Mostra diferenças no repositório Git.

```
/diff

🔍 Git Diff:
diff --git a/main.go b/main.go
index 123..456
+++ b/main.go
@@ -10,3 +10,5 @@
+func newFunction() {
+    // code
+}
```

---

## 🔐 Aprovações de Ferramentas Sensíveis

Quando o AGI precisa executar ferramentas sensíveis (configuradas em `permissions.sensitive_tools`), você receberá uma solicitação de aprovação:

**AGI:**
```
⚠️ Solicitação de Aprovação

Ferramenta: git_commit
Argumentos: {"message": "Add new feature"}

[✅ Aprovar] [❌ Negar]
```

**Você:** Clica em "✅ Aprovar"

**AGI:**
```
✅ Aprovado!
Executando git_commit...
Commit criado com sucesso: abc123
```

### Timeout de Aprovação

Se você não responder em 5 minutos (configurável em `permissions.telegram_approval_timeout`):

```
⏰ Timeout: Operação negada automaticamente após 5 minutos.
```

---

## 📊 Respostas Longas

Respostas longas são automaticamente divididas em partes:

**Você:** `Explique toda a arquitetura do projeto`

**AGI:**
```
📝 Resposta (parte 1/3):

A arquitetura do ClosedWheelerAGI segue uma estrutura modular...
[conteúdo da parte 1]

(Continuação 2/3)
[conteúdo da parte 2]

(Continuação 3/3)
[conteúdo da parte 3]
```

---

## ⚙️ Configurações Avançadas

### Permissões

Configure quais comandos e ferramentas são permitidos via Telegram:

```json
{
  "permissions": {
    "allowed_commands": ["*"],  // Todos os comandos
    // ou
    "allowed_commands": ["/status", "/logs", "/help"],  // Apenas específicos

    "allowed_tools": ["*"],  // Todas as ferramentas
    // ou
    "allowed_tools": ["read_file", "list_files", "git_status"],  // Apenas leitura

    "sensitive_tools": [
      "git_commit",
      "git_push",
      "exec_command",
      "write_file",
      "delete_file"
    ],

    "auto_approve_non_sensitive": false,  // Requer aprovação manual
    "require_approval_for_all": false,     // Ou requer para tudo
    "telegram_approval_timeout": 300       // 5 minutos
  }
}
```

### Notificações

Habilite notificações quando ferramentas são executadas:

```json
{
  "telegram": {
    "enabled": true,
    "notify_on_tool_start": true  // Notifica ao iniciar ferramentas
  }
}
```

Com isso habilitado, você receberá:

```
🔧 Executando: read_file
Arquivo: /path/to/file.go

[resultado da ferramenta]
```

---

## 🔒 Segurança

### Chat ID Único

Apenas o Chat ID configurado pode:
- Executar comandos
- Conversar com o AGI
- Aprovar/negar operações

**Outros usuários recebem**:
```
🔒 Acesso negado.
Seu Chat ID (987654321) não está autorizado.
```

### Audit Log

Todas as interações via Telegram são logadas em `.agi/audit.log`:

```json
{"timestamp":"2026-02-08T10:30:45Z","action":"command","name":"/status","allowed":true,"user_id":123456789}
{"timestamp":"2026-02-08T10:31:00Z","action":"tool","name":"read_file","allowed":true,"user_id":123456789}
{"timestamp":"2026-02-08T10:31:30Z","action":"approval","name":"git_commit","allowed":true,"reason":"approved by user","user_id":123456789}
```

### Tokens Seguros

**✅ Faça**:
- Guarde bot token no `.env` (gitignored)
- Use chat IDs específicos
- Revise o audit log regularmente
- Habilite aprovações para ferramentas sensíveis

**❌ Não Faça**:
- Compartilhe o bot token
- Use chat_id = 0 (qualquer um)
- Desabilite audit log
- Auto-aprove tudo

---

## 🧪 Testes

### Teste Básico de Conexão

```bash
# 1. Configure bot token e chat ID
# 2. Inicie o agent
./ClosedWheeler

# 3. Envie /start no Telegram
# Deve responder com seu Chat ID
```

### Teste de Chat

```bash
# Envie qualquer mensagem
"Olá, você está funcionando?"

# Deve responder
"💭 Pensando..."
"Sim, estou funcionando perfeitamente! Como posso ajudar?"
```

### Teste de Aprovação

```bash
# Configure git_commit como sensível
# Peça ao AGI: "Faça um commit das alterações"

# Deve aparecer:
"⚠️ Solicitação de Aprovação
Ferramenta: git_commit
[✅ Aprovar] [❌ Negar]"
```

---

## 🐛 Troubleshooting

### Bot não responde

**Problema**: Mensagens não chegam

**Soluções**:
```bash
# 1. Verifique se Telegram está habilitado
cat .agi/config.json | grep -A5 telegram

# 2. Verifique bot token
echo $TELEGRAM_BOT_TOKEN

# 3. Verifique chat ID
# Deve estar configurado no config.json

# 4. Verifique logs
tail -f .agi/agent.log | grep -i telegram
```

### "Acesso negado"

**Problema**: Bot diz que não tenho permissão

**Solução**:
```bash
# 1. Envie /start para obter seu Chat ID
# 2. Configure no config.json:
{
  "telegram": {
    "chat_id": SEU_CHAT_ID_AQUI
  }
}

# 3. Reinicie o agent
```

### Respostas cortadas

**Problema**: Resposta parece incompleta

**Explicação**: Telegram tem limite de 4096 caracteres por mensagem. Respostas longas são automaticamente divididas em partes.

**Solução**: Aguarde todas as partes serem enviadas (marcadas como "parte 1/3", "parte 2/3", etc.)

### Timeout de aprovação

**Problema**: Operação negada mesmo que eu queira aprovar

**Solução**:
```json
// Aumente o timeout (em segundos)
{
  "permissions": {
    "telegram_approval_timeout": 600  // 10 minutos
  }
}
```

---

## 📱 Exemplos de Conversação

### Análise de Código

**Você:**
```
Analise o arquivo pkg/agent/agent.go e me diga se há algum problema
```

**AGI:**
```
💭 Pensando...

🔍 Analisando pkg/agent/agent.go...

O arquivo está bem estruturado, mas identifiquei algumas melhorias:

1. ✅ Boa separação de responsabilidades
2. ✅ Tratamento de erros adequado
3. ⚠️ Falta documentação em alguns métodos privados
4. ⚠️ Poderia ter mais testes unitários

Sugestões:
- Adicionar godoc comments nos métodos handleTelegramChat e splitMessage
- Criar testes para a função de split de mensagens

Deseja que eu implemente essas melhorias?
```

### Geração de Código

**Você:**
```
Crie um middleware de autenticação JWT para uma API REST em Go
```

**AGI:**
```
💭 Pensando...

🔧 Criando middleware de autenticação JWT...

Criei o arquivo middleware/auth.go com o middleware completo:

```go
package middleware

import (
    "net/http"
    "github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secretKey string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Implementation here...
        })
    }
}
```

Também criei testes unitários em middleware/auth_test.go.

Deseja que eu crie exemplos de uso?
```

### Refatoração

**Você:**
```
O método getUserData() está muito grande, refatore-o
```

**AGI:**
```
💭 Pensando...

🔧 Refatorando getUserData()...

Dividi o método em 3 funções menores:

1. validateUserInput() - Validação de entrada
2. fetchUserFromDB() - Busca no banco
3. enrichUserData() - Enriquecimento com dados extras

O código ficou mais testável e legível.

Deseja que eu também crie testes para essas novas funções?
```

---

## 🎯 Best Practices

### Para Desenvolvedores

1. **Use comandos para ações rápidas**
   ```
   /status    # Verificação rápida
   /logs      # Debug
   /diff      # Review de código
   ```

2. **Use chat para tarefas complexas**
   ```
   "Analise todo o pacote auth e sugira melhorias de segurança"
   "Crie testes para todos os métodos da classe User"
   "Refatore o código seguindo os princípios SOLID"
   ```

3. **Aprove apenas operações que você entende**
   - Leia os argumentos antes de aprovar
   - Se não entender, negue e pergunte ao AGI

4. **Monitore o audit log**
   ```bash
   tail -f .agi/audit.log | jq
   ```

### Para Times

1. **Um bot por projeto**
   - Cada projeto deve ter seu próprio bot
   - Facilita controle de acesso e auditoria

2. **Chat IDs diferentes para dev/prod**
   - Desenvolvedores têm acesso total em dev
   - Apenas leads têm acesso em prod

3. **Revise configurações de permissões**
   - Minimize ferramentas permitidas
   - Require approval para operações críticas

---

## 📊 Estatísticas de Uso

Monitore o uso via audit log:

```bash
# Total de mensagens enviadas
cat .agi/audit.log | grep "\"action\":\"command\"" | wc -l

# Ferramentas mais usadas
cat .agi/audit.log | jq -r 'select(.action=="tool") | .name' | sort | uniq -c | sort -rn

# Taxa de aprovação
cat .agi/audit.log | jq -r 'select(.action=="approval") | .allowed' | grep true | wc -l
```

---

## 🚀 Próximas Funcionalidades

### Planejado
- [ ] Suporte para múltiplos usuários (RBAC)
- [ ] Webhooks (mais rápido que polling)
- [ ] Comandos customizados via plugins
- [ ] Respostas com formatação rica (imagens, arquivos)
- [ ] Notificações proativas (build failed, deploy complete)

---

**Status**: ✅ **PRODUCTION READY**
**Build**: ✅ **11MB**
**Chat via Telegram**: ✅ **100% FUNCIONAL**

*Converse com seu AGI de qualquer lugar! 🌎📱*
