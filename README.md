# go-gen - Генратор базового проекта на Golang

## Установка и запуск

1. скачиваем генератор "go get github.com/0LuigiCode0/go-gen"
2. устанавливаем его "go install github.com/0LuigiCode0/go-gen"
3. заходим в папку где хотим создать проект
4. закидываем туда конфиг (можно и поменять его)
5. вызываем из папки "go-gen --file config.json"
6. запускаем проект (скорее всего будет ругаться на отсутствие пакетов, так что ставим их ручками)

## Разбор config.json
```json
{
   "module_name": "root",
   "go_version": 1.16,
   "work_dir": "root",
   "dbs": {
      "postgres": "postgres",
      "mongo": "mongodb",
   },
   "handlers": {
      "roots": "tcp",
      "mqtts": "mqtt",
      "wss": "ws"
   }
}
```
* `"module_name": "root"` - как будет называтся проект
* `"go_version": 1.16` - версия golang
* `"work_dir": "root"` - рабочая дирректория проекта, если указана то преокт создастся в ней
* `"dbs"` - массив бд, ключ это название базы внутри проекта (такжк будут названы и пакеты связанные с данноый бд), значение это выбор драйвера
    * `mongodb` - MongoDB
    * `postgres` - PostgreSQL
* ` "handlers"` - массив web интерфейсов, ключ это название интерфейса внутри проекта (такжк будут названы и пакеты связанные с данноым интерфейсом), значение это выбор web интрефейса
    * `tcp` - Обычное http соединение
    * `ws` - WebSocket
    * `mqtt` - MQTT