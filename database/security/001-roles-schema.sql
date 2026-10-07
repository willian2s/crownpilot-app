-- Security baseline only. EF Core owns application tables and constraints.
BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crownpilot_runtime') THEN
        CREATE ROLE crownpilot_runtime
            NOLOGIN
            NOSUPERUSER
            NOCREATEDB
            NOCREATEROLE
            NOINHERIT
            NOREPLICATION
            NOBYPASSRLS;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crownpilot_migrator') THEN
        CREATE ROLE crownpilot_migrator
            NOLOGIN
            NOSUPERUSER
            NOCREATEDB
            NOCREATEROLE
            NOINHERIT
            NOREPLICATION
            NOBYPASSRLS;
    END IF;
END
$$;

ALTER ROLE crownpilot_runtime
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    NOINHERIT
    NOREPLICATION
    NOBYPASSRLS;

ALTER ROLE crownpilot_migrator
    NOLOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    NOINHERIT
    NOREPLICATION
    NOBYPASSRLS;

ALTER SCHEMA crownpilot OWNER TO crownpilot_migrator;

REVOKE ALL ON SCHEMA crownpilot FROM PUBLIC;
GRANT USAGE ON SCHEMA crownpilot TO crownpilot_runtime;
GRANT USAGE, CREATE ON SCHEMA crownpilot TO crownpilot_migrator;

-- The private schema is not part of the Supabase Data API exposed-schema set.
DO $$
DECLARE
    role_name text;
BEGIN
    FOREACH role_name IN ARRAY ARRAY['anon', 'authenticated', 'service_role'] LOOP
        IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = role_name) THEN
            EXECUTE format('REVOKE ALL ON SCHEMA crownpilot FROM %I', role_name);
        END IF;
    END LOOP;
END
$$;

-- Future EF-created objects receive least-privilege runtime grants. This does
-- not own domain DDL, policies, or domain data.
ALTER DEFAULT PRIVILEGES IN SCHEMA crownpilot
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crownpilot_runtime;
ALTER DEFAULT PRIVILEGES IN SCHEMA crownpilot
    GRANT USAGE, SELECT ON SEQUENCES TO crownpilot_runtime;
ALTER DEFAULT PRIVILEGES FOR ROLE crownpilot_migrator IN SCHEMA crownpilot
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crownpilot_runtime;
ALTER DEFAULT PRIVILEGES FOR ROLE crownpilot_migrator IN SCHEMA crownpilot
    GRANT USAGE, SELECT ON SEQUENCES TO crownpilot_runtime;

COMMIT;
