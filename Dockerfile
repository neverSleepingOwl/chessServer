FROM nginx 
COPY main /
COPY client/ /usr/share/nginx/html/
EXPOSE 80
EXPOSE 8080
CMD ["nginx -s start"]
CMD ["./main"]
