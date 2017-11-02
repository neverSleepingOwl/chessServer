FROM nginx 
USER root
COPY main /
COPY html/* /data/www/
COPY html/js/* /data/www/js/
COPY html/images/ /data/www/images/
COPY nginx.conf /etc/nginx/nginx.conf
COPY wrapper.sh /
EXPOSE 80
EXPOSE 8080
CMD ["./wrapper.sh"]
