create table if not exists public.dick (
    id bigserial not null,
    user_id bigint not null,
    chat_id bigint not null,
    length bigint not null,
    updated_at timestamptz not null default now()
);