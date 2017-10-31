FROM NGINX
COPY main /
COPY client/ /usr/share/nginx/html
EXPOSE 80
EXPOSE 443
EXPOSE 8080
