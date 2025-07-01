# Script de Teste - Fluxo de Pagamentos

## Pré-requisitos

1. Configure as variáveis de ambiente no arquivo `.env`
2. Execute o servidor: `go run cmd/app/main.go`
3. Tenha uma conta Stripe configurada com produtos e preços

## Teste 1: Registro e Login do Usuário

### 1.1 Registrar usuário
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 1.2 Fazer login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Salve o token JWT retornado para usar nos próximos testes.**

## Teste 2: Criar Pagamento para Robô

### 2.1 Tentar criar robô diretamente (deve falhar)
```bash
curl -X POST http://localhost:8080/api/robots \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -d '{
    "name": "TestRobot"
  }'
```

**Resultado esperado:** Erro informando que deve usar o endpoint de pagamento.

### 2.2 Criar sessão de pagamento
```bash
curl -X POST http://localhost:8080/api/payments/robot \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -d '{
    "robot_name": "TestRobot",
    "plan_type": "basic"
  }'
```

**Resultado esperado:** 
```json
{
  "session_id": "cs_test_...",
  "checkout_url": "https://checkout.stripe.com/...",
  "message": "Payment session created. Complete the payment to create your robot."
}
```

## Teste 3: Verificar Estado Antes do Pagamento

### 3.1 Listar robôs (deve estar vazio)
```bash
curl -X GET http://localhost:8080/api/robots \
  -H "Authorization: Bearer {JWT_TOKEN}"
```

**Resultado esperado:** Lista vazia ou sem robôs ativos.

## Teste 4: Simular Webhook do Stripe

### 4.1 Simular evento de pagamento confirmado
```bash
curl -X POST http://localhost:8080/api/stripe/webhook \
  -H "Content-Type: application/json" \
  -H "Stripe-Signature: t=1234567890,v1=fake_signature" \
  -d '{
    "type": "checkout.session.completed",
    "data": {
      "object": {
        "id": "{SESSION_ID_DO_PASSO_2.2}",
        "customer": {
          "id": "cus_test_customer"
        },
        "subscription": {
          "id": "sub_test_subscription"
        },
        "metadata": {
          "user_id": "{USER_ID}",
          "robot_name": "TestRobot",
          "plan_type": "basic"
        }
      }
    }
  }'
```

**Nota:** Este teste pode falhar devido à validação de assinatura do Stripe. Para teste completo, use o webhook test do Stripe.

## Teste 5: Verificar Robô Criado

### 5.1 Listar robôs (deve mostrar o robô criado)
```bash
curl -X GET http://localhost:8080/api/robots \
  -H "Authorization: Bearer {JWT_TOKEN}"
```

**Resultado esperado:** Lista com o TestRobot com status "active".

### 5.2 Gerar token para o robô
```bash
curl -X POST http://localhost:8080/api/robots/{ROBOT_ID}/token \
  -H "Authorization: Bearer {JWT_TOKEN}"
```

**Resultado esperado:** Token JWT do robô.

## Teste 6: Testar Conversa com Robô

### 6.1 Enviar mensagem usando token do robô
```bash
curl -X POST http://localhost:8080/api/conversa \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {ROBOT_JWT_TOKEN}" \
  -d '{
    "mensagem": "Olá, robô!"
  }'
```

**Resultado esperado:** Resposta da IA configurada.

## Teste 7: Verificar Validação de Plano

### 7.1 Simular plano expirado (modificar diretamente no banco)
```sql
UPDATE robots 
SET plan_valid_until = '2023-01-01 00:00:00' 
WHERE name = 'TestRobot';
```

### 7.2 Tentar conversar novamente
```bash
curl -X POST http://localhost:8080/api/conversa \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {ROBOT_JWT_TOKEN}" \
  -d '{
    "mensagem": "Esta mensagem deve falhar"
  }'
```

**Resultado esperado:** Erro 402 (Payment Required) informando que a assinatura expirou.

## Teste 8: Verificar Webhook de Falha de Pagamento

### 8.1 Simular falha de pagamento
```bash
curl -X POST http://localhost:8080/api/stripe/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "type": "checkout.session.async_payment_failed",
    "data": {
      "object": {
        "id": "{SESSION_ID}",
        "metadata": {
          "user_id": "{USER_ID}",
          "robot_name": "FailedRobot",
          "plan_type": "basic"
        }
      }
    }
  }'
```

**Resultado esperado:** Pagamento marcado como falhou, robô não é criado.

## Verificação do Banco de Dados

### Verificar tabelas criadas
```sql
.tables
-- Deve mostrar: users, robots, plans, payments, subscriptions, conversa_logs

SELECT * FROM payments;
SELECT * FROM subscriptions;
SELECT * FROM robots WHERE name = 'TestRobot';
```

## Logs Importantes

Durante os testes, monitore os logs do servidor para:

1. **Criação de pagamento**: Confirmação de que a sessão foi criada
2. **Recebimento de webhook**: Logs de processamento de eventos
3. **Criação de robô**: Confirmação de que o robô foi criado automaticamente
4. **Validação de middleware**: Logs de verificação de plano ativo

## Troubleshooting Comum

### Erro: "invalid plan type"
- Verifique se está usando "basic", "premium" ou "enterprise"

### Erro: "failed to create payment session"
- Verifique as variáveis de ambiente do Stripe
- Confirme que os PRICE_IDs estão corretos

### Webhook não funciona
- Para teste local, use ngrok ou similar
- Configure o endpoint no dashboard do Stripe
- Verifique o WEBHOOK_SECRET

### Robô não foi criado após webhook
- Verifique logs do servidor
- Confirme que o USER_ID no metadata está correto
- Verifique se o pagamento foi encontrado no banco

## Próximos Testes

1. **Teste de renovação**: Simular renovação automática
2. **Teste de cancelamento**: Cancelar assinatura e verificar desativação
3. **Teste de múltiplos robôs**: Criar vários robôs para o mesmo usuário
4. **Teste de diferentes planos**: Testar premium e enterprise
5. **Teste de concorrência**: Múltiplos usuários criando robôs simultaneamente
