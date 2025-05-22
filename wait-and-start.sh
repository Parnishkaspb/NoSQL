#!/bin/bash
set -e

echo "⏳ Ждем доступности MongoDB ReplicaSet..."
until nc -z mongo1 27017 && nc -z mongo2 27017; do
  echo "🔄 Ожидание TCP-портов mongo1/mongo2..."
  sleep 2
done

echo "✅ Порты MongoDB открыты — запускаем приложение"
exec ./server