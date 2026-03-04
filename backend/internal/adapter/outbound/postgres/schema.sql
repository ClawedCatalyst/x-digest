create extension if not exists pgcrypto;

create table if not exists app_users (
  id uuid primary key default gen_random_uuid(),
  x_user_id text not null unique,
  username text not null,
  created_at timestamptz not null default now()
);

create table if not exists tokens (
  user_id uuid primary key references app_users(id) on delete cascade,
  access_token_enc text not null,
  refresh_token_enc text,
  expires_at timestamptz not null,
  updated_at timestamptz not null default now()
);

-- Rename legacy daily_digest table to digest if it exists.
DO $$
BEGIN
  IF to_regclass('public.daily_digest') IS NOT NULL
     AND to_regclass('public.digest') IS NULL THEN
    EXECUTE 'ALTER TABLE daily_digest RENAME TO digest';
  END IF;
END$$;

create table if not exists digest (
  user_id uuid references app_users(id) on delete cascade,
  day date not null,
  data jsonb not null,
  created_at timestamptz not null default now(),
  primary key (user_id, day)
);
