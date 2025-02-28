#!/bin/sh

set -e  # Остановить выполнение при ошибке

echo "Running migrations..."
./migration

echo "Starting application..."
exec ./todo-app  # Запускаем приложение
