#!/bin/sh

read -p "Enter an email which to use for TLS: " EMAIL

echo "{
    email ${EMAIL}
}

import /etc/caddy/*.Caddyfile
" > /etc/caddy/Caddyfile