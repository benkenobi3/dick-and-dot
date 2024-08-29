begin;

drop index concurrently if exists dick_user_id_chat_id_idx;

create index concurrently if not exists chat_id_user_id_index ON public.dick (chat_id, user_id);

commit;