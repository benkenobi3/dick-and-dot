begin;

create unique index if not exists dick_user_id_chat_id_idx ON public.dick (user_id,chat_id);

commit;
