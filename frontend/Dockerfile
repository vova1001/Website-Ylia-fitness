FROM nginx:alpine

# Копируем фронтенд
COPY . /usr/share/nginx/html/

# Копируем nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]