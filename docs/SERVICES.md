# Сервисы инфраструктуры

## Redis
- **Порт**: 6379
- **Пользователь**: нет (используется пароль)
- **Пароль**: `redis` (или `REDIS_PASSWORD` из .env)
- **Connection String**: `redis://:<password>@localhost:6379`

### RedisInsight (UI)
- **URL**: http://localhost:5540

---

## MinIO
- **S3 API**: http://localhost:9000
- **Console UI**: http://localhost:9001
- **Access Key**: `minio` (или `MINIO_ROOT_USER` из .env)
- **Secret Key**: `minio12345` (или `MINIO_ROOT_PASSWORD` из .env)

---

## RabbitMQ
- **AMQP Port**: 5672
- **Management UI**: http://localhost:15672
- **Пользователь**: `rabbit` (или `RABBITMQ_USER` из .env)
- **Пароль**: `rabbit123` (или `RABBITMQ_PASSWORD` из .env)
- **Connection String**: `amqp://rabbit:rabbit123@localhost:5672/`

### Management UI
- **URL**: http://localhost:15672
- **Логин**: `rabbit`
- **Пароль**: `rabbit123`

---

## Запуск

```bash
# Скопировать пример .env файла
cp .env.example .env

# Запустить все сервисы
docker-compose up -d

# Остановить
docker-compose down

# Остановить и удалить данные
docker-compose down -v
```
