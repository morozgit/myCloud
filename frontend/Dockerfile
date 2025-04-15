FROM node:18-alpine AS build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только файлы для установки зависимостей
COPY package.json package-lock.json ./

# Устанавливаем зависимости с кэшированием
RUN npm install --legacy-peer-deps

# Копируем остальные файлы проекта
COPY . .

# Запускаем сборку проекта
RUN npm run build --verbose

# Производственный образ Nginx
FROM nginx:stable-alpine

# Копируем собранное приложение из первого этапа
COPY --from=build /app/dist /usr/share/nginx/html

# Копируем Nginx-конфигурацию
COPY --from=build /app/nginx.conf /etc/nginx/conf.d/default.conf

# Экспонируем порт 3000
EXPOSE 3000

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]
