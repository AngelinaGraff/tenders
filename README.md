# API для управления тендерами и предложениями.
### Стек
- Go
- Postgres
- Docker

### Настройка приложения производится через переменные окружения

- `SERVER_ADDRESS` — адрес и порт, который будет слушать HTTP сервер при запуске. Пример: 0.0.0.0:8080.
- `POSTGRES_CONN` — URL-строка для подключения к PostgreSQL в формате postgres://{username}:{password}@{host}:{5432}/{dbname}.
- `POSTGRES_JDBC_URL` — JDBC-строка для подключения к PostgreSQL в формате jdbc:postgresql://{host}:{port}/{dbname}.
- `POSTGRES_USERNAME` — имя пользователя для подключения к PostgreSQL.
- `POSTGRES_PASSWORD` — пароль для подключения к PostgreSQL.
- `POSTGRES_HOST` — хост для подключения к PostgreSQL (например, localhost).
- `POSTGRES_PORT` — порт для подключения к PostgreSQL (например, 5432).
- `POSTGRES_DATABASE` — имя базы данных PostgreSQL, которую будет использовать приложение.

### Бизнес-логика
#### Тендер

Тендеры могут создавать только пользователи от имени своей организации.

Доступные действия с тендером:

- **Создание**:

  - Тендер будет создан.

  - Доступен только ответственным за организацию.

  - Статус: `CREATED`.

- **Публикация**:

  - Тендер становится доступен всем пользователям.

  - Статус: `PUBLISHED`.

- **Закрытие**:

  - Тендер больше не доступен пользователям, кроме ответственных за организацию.

  - Статус: `CLOSED`.

- **Редактирование**:

  - Изменяются характеристики тендера.

  - Увеличивается версия.

#### Предложение

Предложения могут создавать пользователи от имени своей организации.

Предложение связано только с одним тендером. Один пользователь может быть ответственным в одной организации.

Доступные действия с предложениями:

- **Создание**:

  - Предложение будет создано.

  - Доступно только автору и ответственным за организацию.

  - Статус: `CREATED`.

- **Публикация**:

  - Предложение становится доступно ответственным за организацию и автору.

  - Статус: `PUBLISHED`.

- **Отмена**:

  - Виден только автору и ответственным за организацию.

  - Статус: `CANCELED`.

- **Редактирование**:

  - Изменяются характеристики предложения.

  - Увеличивается версия.

- **Согласование/отклонение**:

  - Доступно только ответственным за организацию, связанной с тендером.

  - Решение может быть принято любым ответственным.

  - При согласовании одного предложения, тендер автоматически закрывается.