FROM nginx:alpine

# Копируем готовые файлы фронтенда
COPY . /usr/share/nginx/html/

# Копируем nginx конфиг
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]