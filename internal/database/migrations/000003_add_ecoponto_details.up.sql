-- Adiciona as novas colunas para detalhes do ecoponto
-- 000003_add_ecoponto_details.up.sql
ALTER TABLE ecopontos
ADD COLUMN horario_funcionamento VARCHAR(255) NULL,
ADD COLUMN foto_url VARCHAR(255) NULL;