#!/bin/bash

# Replace with your actual domain and email
domain="parsaimi.xyz"
email="parsaemani17@gmail.com"

# Create required directories
mkdir -p ./certbot/conf
mkdir -p ./certbot/www

# Create dummy certificate for $domain
mkdir -p ./certbot/conf/live/$domain
openssl req -x509 -nodes -newkey rsa:4096 -days 1 \
  -keyout ./certbot/conf/live/$domain/privkey.pem \
  -out ./certbot/conf/live/$domain/fullchain.pem \
  -subj "/CN=localhost"

# Create required directories for SSL
mkdir -p ./certbot/conf/live/$domain
touch ./certbot/conf/options-ssl-nginx.conf
touch ./certbot/conf/ssl-dhparams.pem


echo "Starting containers..."
docker compose up -d

echo "Sleeping for 5 seconds to ensure nginx is up and running..."
sleep 5

# Request Let's Encrypt certificate
echo "Requesting Let's Encrypt certificate for $domain..."
docker compose run --rm certbot certonly --webroot -w /var/www/certbot \
  --email $email --agree-tos --no-eff-email \
  -d $domain

echo "Restarting Nginx to apply new certificate..."
docker compose exec nginx nginx -s reload

echo "SSL setup complete!"
