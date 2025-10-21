FROM nginx:alpine

# Копируем фронтенд из папки frontend/
COPY ./frontend/ /usr/share/nginx/html/

# Копируем nginx config из папки frontend/
COPY ./frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]