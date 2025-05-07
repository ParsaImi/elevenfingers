#!/bin/bash

if ! [ -x "$(command -v docker compose)" ]; then
  echo 'Error: docker-compose is not installed.' >&2
  exit 1
fi

domains=(parsaimi.xyz www.parsaimi.xyz)
newdomains=(game.parsaimi.xyz)
serverdomains=(cloud.parsaimi.xyz)
rsa_key_size=4096
data_path="./data/certbot"
email="parsaemani17@gmail.com" # Adding a valid address is strongly recommended
staging=0 # Set to 1 if you're testing your setup to avoid hitting request limits

if [ -d "$data_path" ]; then
  read -p "Existing data found for $domains. Continue and replace existing certificate? (y/N) " decision
  if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
    exit
  fi
fi


if [ ! -e "$data_path/conf/options-ssl-nginx.conf" ] || [ ! -e "$data_path/conf/ssl-dhparams.pem" ]; then
  echo "### Downloading recommended TLS parameters ..."
  mkdir -p "$data_path/conf"
  curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot-nginx/certbot_nginx/_internal/tls_configs/options-ssl-nginx.conf > "$data_path/conf/options-ssl-nginx.conf"
  curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot/certbot/ssl-dhparams.pem > "$data_path/conf/ssl-dhparams.pem"
  echo
fi

echo "### Creating dummy certificate for $domains ..."
path="/etc/letsencrypt/live/$domains"
mkdir -p "$data_path/conf/live/$domains"
docker compose run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:$rsa_key_size -days 1\
    -keyout '$path/privkey.pem' \
    -out '$path/fullchain.pem' \
    -subj '/CN=localhost'" certbot
echo

echo "### Creating dummy certificate for $newdomains ..."
path="/etc/letsencrypt/live/$newdomains"
mkdir -p "$data_path/conf/live/$newdomains"
docker compose run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:$rsa_key_size -days  1\
    -keyout '$path/privkey.pem' \
    -out '$path/fullchain.pem' \
    -subj '/CN=localhost'" certbot
echo

echo "### Creating dummy certificate for $serverdomains..."
path="/etc/letsencrypt/live/$serverdomains"
mkdir -p "$data_path/conf/live/$serverdomains"
docker compose run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:$rsa_key_size -days  1\
    -keyout '$path/privkey.pem' \
    -out '$path/fullchain.pem' \
    -subj '/CN=localhost'" certbot
echo

echo "### Starting nginx ..."
docker compose up --force-recreate -d nginx
echo

echo "### Deleting dummy certificate for $domains ..."
docker compose run --rm --entrypoint "\
  rm -Rf /etc/letsencrypt/live/$domains && \
  rm -Rf /etc/letsencrypt/archive/$domains && \
  rm -Rf /etc/letsencrypt/renewal/$domains.conf" certbot
echo

echo "### Deleting dummy certificate for $newdomains ..."
docker compose run --rm --entrypoint "\
  rm -Rf /etc/letsencrypt/live/$newdomains && \
  rm -Rf /etc/letsencrypt/archive/$newdomains && \
  rm -Rf /etc/letsencrypt/renewal/$newdomains.conf" certbot
echo


echo "### Deleting dummy certificate for $serverdomains ..."
docker compose run --rm --entrypoint "\
  rm -Rf /etc/letsencrypt/live/$serverdomains && \
  rm -Rf /etc/letsencrypt/archive/$serverdomains && \
  rm -Rf /etc/letsencrypt/renewal/$serverdomains.conf" certbot
echo



echo "### Requesting Let's Encrypt certificate for $domains ..."
#Join $domains to -d args
domain_args=""
for domain in "${domains[@]}"; do
  domain_args="$domain_args -d $domain"
done


echo "### Requesting Let's Encrypt certificate for $newdomains ..."
#Join $newdomains to -d args
newdomain_args=""
for domain in "${newdomains[@]}"; do
  newdomain_args="$newdomain_args -d $domain"
done

echo "### Requesting Let's Encrypt certificate for $serverdomains ..."
#Join $serverdomains to -d args
serverdomain_args=""
for domain in "${serverdomains[@]}"; do
  serverdomain_args="$serverdomain_args -d $domain"
done

# Select appropriate email arg
case "$email" in
  "") email_arg="--register-unsafely-without-email" ;;
  *) email_arg="--email $email" ;;
esac

# Enable staging mode if needed
if [ $staging != "0" ]; then staging_arg="--staging"; fi

docker compose run --rm --entrypoint "\
  certbot certonly --webroot -w /var/www/certbot \
    $staging_arg \
    $email_arg \
    $serverdomain_args \
    --rsa-key-size $rsa_key_size \
    --agree-tos \
    --force-renewal" certbot
echo

docker compose run --rm --entrypoint "\
  certbot certonly --webroot -w /var/www/certbot \
    $staging_arg \
    $email_arg \
    $domain_args \
    --rsa-key-size $rsa_key_size \
    --agree-tos \
    --force-renewal" certbot
echo

docker compose run --rm --entrypoint "\
  certbot certonly --webroot -w /var/www/certbot \
    $staging_arg \
    $email_arg \
    $newdomain_args \
    --rsa-key-size $rsa_key_size \
    --agree-tos \
    --force-renewal" certbot
echo




echo "### Reloading nginx ..."
docker compose exec nginx nginx -s reload
