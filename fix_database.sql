-- Script para corrigir o banco de dados

-- 1. Primeiro, verificar os índices existentes
.indexes plans

-- 2. Remover o índice único problemático (se existir)
DROP INDEX IF EXISTS idx_robot_active;

-- 3. Desativar planos duplicados
UPDATE plans 
SET active = false 
WHERE id NOT IN (
    SELECT MIN(id) 
    FROM plans 
    WHERE active = true 
    GROUP BY robot_id
);

-- 4. Verificar se ainda há duplicatas
SELECT robot_id, COUNT(*) as count 
FROM plans 
WHERE active = true 
GROUP BY robot_id 
HAVING COUNT(*) > 1;
