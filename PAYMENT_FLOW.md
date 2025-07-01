# Fluxo de Pagamentos com Stripe - Documentação

## Visão Geral

Este documento descreve o fluxo completo de pagamentos implementado no sistema de robôs, garantindo que:

1. **Robôs só são criados após confirmação de pagamento**
2. **Robôs funcionam apenas enquanto o plano estiver ativo**
3. **Falhas de pagamento refletem imediatamente no status do robô**

## Arquitetura

### Modelos Principais

1. **Payment**: Gerencia todos os pagamentos
2. **Subscription**: Gerencia assinaturas recorrentes
3. **Robot**: Robôs com controle de status e validade de plano
4. **Plan**: Planos associados a robôs (depreciado em favor de Subscription)

### Status de Robô

- `pending`: Aguardando pagamento
- `active`: Pagamento confirmado e funcionando
- `suspense`: Pagamento falhou ou plano expirado

## Fluxo de Criação de Robô

### 1. Usuário Solicita Criação (Frontend → Backend)

```http
POST /api/payments/robot
Authorization: Bearer {user_jwt_token}
Content-Type: application/json

{
  "robot_name": "MeuRobot",
  "plan_type": "basic" // ou "premium", "enterprise"
}
```

**Resposta:**
```json
{
  "session_id": "cs_test_...",
  "checkout_url": "https://checkout.stripe.com/...",
  "message": "Payment session created. Complete the payment to create your robot."
}
```

### 2. Usuário Completa Pagamento (Stripe Checkout)

O usuário é redirecionado para o Stripe Checkout e completa o pagamento.

### 3. Webhook do Stripe Confirma Pagamento

```http
POST /api/stripe/webhook
Content-Type: application/json
Stripe-Signature: t=...,v1=...

{
  "type": "checkout.session.completed",
  "data": {
    "object": {
      "id": "cs_test_...",
      "metadata": {
        "user_id": "uuid",
        "robot_name": "MeuRobot",
        "plan_type": "basic"
      }
    }
  }
}
```

### 4. Sistema Cria Robô Automaticamente

Quando o webhook é recebido:
1. Localiza o pagamento pela `session_id`
2. Atualiza status do pagamento para `completed`
3. Cria o robô com status `active`
4. Cria assinatura no banco de dados
5. Define `plan_valid_until` baseado na assinatura

## Validação Contínua

### Middleware de Autenticação do Robô

Toda requisição para `/api/conversa` passa por validações:

1. **Token JWT válido** com `robo_id`
2. **Robô existe** no banco de dados
3. **Status do robô** é `active`
4. **Assinatura ativa** ou `plan_valid_until` ainda válido

```go
// Exemplo de verificação no middleware
if robot.Status != "active" {
    return "Robot is not active. Please check your subscription"
}

if subscription == nil || !subscription.IsActive() {
    if robot.PlanValidUntil == nil || robot.PlanValidUntil.Before(time.Now()) {
        return "Subscription expired. Please renew your plan"
    }
}
```

## Endpoints Principais

### Criação de Pagamento
- **Endpoint**: `POST /api/payments/robot`
- **Auth**: JWT do usuário
- **Função**: Cria sessão de checkout no Stripe

### Webhook do Stripe
- **Endpoint**: `POST /api/stripe/webhook`
- **Auth**: Assinatura do Stripe
- **Função**: Processa eventos do Stripe

### Geração de Token do Robô
- **Endpoint**: `POST /api/robots/{id}/token`
- **Auth**: JWT do usuário
- **Função**: Gera token para robô se plano estiver ativo

### Conversa com Robô
- **Endpoint**: `POST /api/conversa`
- **Auth**: JWT do robô + validação de assinatura
- **Função**: Processa mensagens do robô

## Eventos do Stripe Suportados

### `checkout.session.completed`
- **Ação**: Cria robô e ativa assinatura
- **Resultado**: Robô fica disponível para uso

### `checkout.session.async_payment_failed`
- **Ação**: Marca pagamento como falhou
- **Resultado**: Robô não é criado

### `invoice.payment_succeeded`
- **Ação**: Renovação bem-sucedida
- **Resultado**: Estende validade da assinatura

### `invoice.payment_failed`
- **Ação**: Falha na renovação
- **Resultado**: Robô é suspenso após período de graça

### `customer.subscription.deleted`
- **Ação**: Cancelamento de assinatura
- **Resultado**: Robô é desativado

## Configuração Necessária

### Variáveis de Ambiente

```env
# Stripe
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_BASIC_PRICE_ID=price_...
STRIPE_PREMIUM_PRICE_ID=price_...
STRIPE_ENTERPRISE_PRICE_ID=price_...
STRIPE_SUCCESS_URL=https://seusite.com/success
STRIPE_CANCEL_URL=https://seusite.com/cancel

# JWT
JWT_SECRET_KEY=sua-chave-super-secreta
```

### Produtos no Stripe

1. Criar produtos no dashboard do Stripe
2. Criar preços recorrentes (mensal/anual)
3. Configurar webhook endpoint
4. Definir eventos do webhook

## Segurança

### Boas Práticas Implementadas

1. **Validação de assinatura** nos webhooks do Stripe
2. **Verificação contínua** de plano ativo
3. **Tokens JWT** com expiração
4. **Logs de auditoria** para pagamentos
5. **Metadata encriptada** em pagamentos sensíveis

### Prevenção de Fraudes

1. **Verificação de propriedade**: Usuário só pode usar robôs próprios
2. **Rate limiting**: Implementar nas rotas de conversa
3. **Monitoramento**: Logs de uso excessivo
4. **Validação de entrada**: Sanitização de dados

## Monitoramento

### Métricas Importantes

1. **Taxa de conversão** de pagamentos
2. **Robôs ativos** vs inativos
3. **Renovações bem-sucedidas** vs falhadas
4. **Uso por robô** (mensagens/dia)

### Alertas Recomendados

1. **Falhas de webhook** do Stripe
2. **Pagamentos pendentes** por muito tempo
3. **Robôs com plano expirado** ainda ativos
4. **Uso anômalo** de robôs

## Troubleshooting

### Robô não foi criado após pagamento
1. Verificar logs do webhook
2. Confirmar que evento foi recebido
3. Verificar se metadata estava correta
4. Reprocessar evento manualmente se necessário

### Robô parou de funcionar
1. Verificar status no banco de dados
2. Confirmar validade da assinatura
3. Verificar últimos eventos do Stripe
4. Reativar manualmente se apropriado

### Webhook não está chegando
1. Verificar endpoint configurado no Stripe
2. Confirmar que URL está acessível
3. Verificar assinatura do webhook
4. Testar com webhook test do Stripe

## Próximos Passos

### Melhorias Futuras

1. **Dashboard de administração** para gerenciar assinaturas
2. **API de relatórios** de uso e pagamentos
3. **Integração com outros provedores** de pagamento
4. **Sistema de créditos** como alternativa a assinaturas
5. **Notificações automáticas** para renovações
