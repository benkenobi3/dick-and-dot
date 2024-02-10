# Для локального запуска приложения
run:
	docker compose -f compose.run.yaml up

# Для запуска окружения кроме запуска самого приложения
debug:
	docker compose -f compose.debug.yaml up
