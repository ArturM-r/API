CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE public.users (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              email TEXT NOT NULL UNIQUE,
                              password_hash TEXT NOT NULL,
                              created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);