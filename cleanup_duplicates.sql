-- Script para limpar dados duplicados de planos
-- Este script deve ser executado antes de reiniciar a aplicação

-- 1. Primeiro, vamos desativar os planos duplicados (manter apenas o primeiro criado)
UPDATE plans 
SET active = false 
WHERE id NOT IN (
    SELECT MIN(id) 
    FROM plans 
    WHERE active = true 
    GROUP BY robot_id
);

-- 2. Opcional: Remover completamente os planos duplicados (cuidado!)
-- DELETE FROM plans 
-- WHERE id NOT IN (
--     SELECT MIN(id) 
--     FROM plans 
--     GROUP BY robot_id
-- );

-- 3. Verificar se ainda há duplicatas
SELECT robot_id, COUNT(*) as count 
FROM plans 
WHERE active = true 
GROUP BY robot_id 
HAVING COUNT(*) > 1;
