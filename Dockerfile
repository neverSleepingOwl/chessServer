FROM nginx 
COPY main /
COPY client/ /usr/share/nginx/html/
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
EXPOSE 8080
CMD ["nginx -g daemon off;"]
CMD ["./main"]
