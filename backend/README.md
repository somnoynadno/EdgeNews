# Backend

Основное API системы.

## Переменные окружения

Для запуска проекта необходимо передать следующие значения в env:

- ENV=PRODUCTION
- db_user=postgres
- db_pass=postgres
- db_name=edge
- db_host=postgres
- db_port=5432
- news_api_token=<your_token> // при наличии
- newscatcher_api_token=<your_token> // при наличии
- selenium_path=<path_to_standalone_selenium>
- gecko_driver_path=<path_to_gecko_driver>

Для прокидывания в докер контейнер эти значения должны быть записаны в ```.env``` 
файл данной директории (существует по умолчанию).

## Конфигурация

Модификация ```config/Config.go``` ведёт к изменению дефолтной конфигурации.

## Запуск

Сборка исполняемого файла: ``` $ go build -o main . ```
