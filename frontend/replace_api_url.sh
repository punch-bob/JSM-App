#!/usr/bin/env sh

find '/usr/share/nginx/html' -name '*.js' -exec sed -i -e 's,_API_BASE_URL,'"$API_BASE_URL"',g' {} \;

find '/usr/share/nginx/html' -name '*.js' -exec sed -i -e 's,_API_AUTH_URL,'"$API_AUTH_URL"',g' {} \;
nginx -g "daemon off;"