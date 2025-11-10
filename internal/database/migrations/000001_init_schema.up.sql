-- Habilita a extensão PostGIS, se ainda não estiver habilitada
CREATE EXTENSION IF NOT EXISTS postgis;

-- 1. Cria a tabela principal de ecopontos
CREATE TABLE IF NOT EXISTS ecopontos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  nome VARCHAR(255) NOT NULL,
  tipo_residuo VARCHAR(100) NOT NULL,
  logradouro TEXT,
  bairro VARCHAR(100),
  coordenadas GEOMETRY(Point, 4326) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Cria o índice para buscas geográficas rápidas
CREATE INDEX IF NOT EXISTS idx_ecopontos_coordenadas ON ecopontos USING GIST (coordenadas);