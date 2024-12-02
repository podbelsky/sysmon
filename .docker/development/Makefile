# Указываем конфигурацию для docker-compose
# Можно использовать наследование через символ ":", последний файл важнее
export COMPOSE_FILE=./docker-compose.yml

##
# Запустить/пересоздать контейнеров из образов
#
# Пример: make up
up:
	docker-compose up -d --force-recreate --remove-orphans

##
# Остановить и удалить контейнеры проекта
#
# Пример: make down
##
down:
	docker-compose down -v

##
# Запустить bash внутри контейнера
#
# Пример: make bash-golang
##
bash-%:
	docker-compose exec $* bash

##
# Посмотреть логи сервиса в фоллоу режиме
#
# Пример: make log-php или make log-nginx
##
logs-%:
	docker-compose logs -f $*

##
# Перезапустить сервис
#
# Пример: make restart-php
##
restart-%:
	docker-compose restart $*

##
# Остановить контейнеры
#
# Пример: make stop
##
stop:
	docker-compose stop