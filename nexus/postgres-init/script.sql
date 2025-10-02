CREATE SCHEMA IF NOT EXISTS nexus;

ALTER USER nexus_admin SET search_path TO nexus, public;
ALTER DATABASE nexus SET search_path TO nexus, public;

CREATE EXTENSION IF NOT EXISTS pg_trgm;