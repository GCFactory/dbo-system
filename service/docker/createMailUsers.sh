#!/bin/bash

# Проверяем, установлен ли openssl
if ! command -v openssl &> /dev/null; then
    echo "Ошибка: openssl не установлен. Установите его сначала."
    exit 1
fi

# Указываем путь к файлу (можно изменить)
AUTH_FILE="./mail/auth/auth.txt"

# Проверяем, существует ли файл. Если нет — создаем.
if [ ! -f "$AUTH_FILE" ]; then
    echo "Файл $AUTH_FILE не найден. Создаю новый."
    touch "$AUTH_FILE"
    chmod 777 "$AUTH_FILE"
fi

# Запрашиваем имя пользователя
read -p "Введите имя пользователя: " username

# Запрашиваем пароль (без отображения ввода)
read -s -p "Введите пароль: " password
echo ""

# Генерируем хеш пароля (SHA-512)
password_hash=$(openssl passwd -6 "$password")

# Записываем в файл
echo "$username:$password_hash" >> "$AUTH_FILE"

# Выводим результат
echo "Пользователь '$username' добавлен в $AUTH_FILE"