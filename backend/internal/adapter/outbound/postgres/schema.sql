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

create table if not exists daily_digest (
  user_id uuid references app_users(id) on delete cascade,
  day date not null,
  data jsonb not null,
  created_at timestamptz not null default now(),
  primary key (user_id, day)
);
