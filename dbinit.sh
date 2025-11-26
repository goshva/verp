#!/bin/bash

# 1. Настройка PostgreSQL для внешних подключений
echo "[1/6] Настройка PostgreSQL для внешних подключений..."

# 1.1. Разрешаем прослушивание всех IP-адресов
sudo sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/g" /etc/postgresql/17/main/postgresql.conf

# 1.2. Устанавливаем порт 5431
sudo sed -i "s/port = 5432/port = 5431/g" /etc/postgresql/15/main/postgresql.conf

# 2. Настройка pg_hba.conf для доступа из Windows
echo "[2/6] Настройка pg_hba.conf для доступа из Windows..."

# 2.1. Получаем IP-адрес Windows (предполагается, что это основной шлюз WSL)
WINDOWS_IP=$(ip route | awk '/default/ {print $3}')
echo "Обнаружен IP-адрес Windows: $WINDOWS_IP"

# 2.2. Добавляем правило для подключения из Windows
sudo bash -c "echo \"host    all             all             $WINDOWS_IP/32            md5\" >> /etc/postgresql/17/main/pg_hba.conf"

# 3. Перезапуск PostgreSQL
echo "[3/6] Перезапуск PostgreSQL..."
sudo service postgresql restart

# 4. Проверка, что PostgreSQL слушает порт 5431
echo "[4/6] Проверка, что PostgreSQL слушает порт 5431..."
sudo netstat -tulnp | grep 5431

# 5. Создание пользователя и базы данных (если ещё не созданы)
echo "[5/6] Создание пользователя и базы данных..."
sudo -u postgres psql <<-EOSQL
    -- Меняем пароль для пользователя postgres
    ALTER USER postgres WITH PASSWORD 'password';

    -- Создаём пользователя venderp
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'venderp') THEN
            CREATE USER venderp WITH PASSWORD 'password';
        END IF;
    END
    \$\$;

    -- Создаём базу данных venderp
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'venderp') THEN
            CREATE DATABASE venderp OWNER venderp;
        END IF;
    END
    \$\$;

    -- Даём все права пользователю venderp на базу данных venderp
    GRANT ALL PRIVILEGES ON DATABASE venderp TO venderp;
EOSQL

# 6. Вывод информации о подключении
echo "[6/6] Настройка завершена!"
echo "Вы можете подключаться из Windows с параметрами:"
echo "DB_HOST=$WINDOWS_IP (или IP-адрес WSL: $(hostname -I))"
echo "DB_PORT=5431"
echo "DB_USER=venderp"
echo "DB_PASSWORD=password"
echo "DB_NAME=venderp"
echo "SSL_MODE=disable"

