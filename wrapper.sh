#!/bin/bash

echo "Nginx is running..."
exec ./main&
exec nginx -g "daemon off;"
