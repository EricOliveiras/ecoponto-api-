-- Remove as colunas
-- 000003_add_ecoponto_details.down.sql
ALTER TABLE ecopontos
DROP COLUMN horario_funcionamento,
DROP COLUMN foto_url;