FROM nginx:alpine

# Install envsubst
RUN apk add --no-cache gettext

# Copy frontend files
COPY frontend/ /usr/share/nginx/html/

# Copy nginx config template
COPY nginx/default.conf.template /etc/nginx/templates/default.conf.template

# Copy config template and create entrypoint
COPY frontend/js/config.template.js /tmp/config.template.js

RUN echo '#!/bin/sh' > /docker-entrypoint.d/30-config.sh && \
    echo 'envsubst "\$API_URL \$ADMIN_SECRET_PATH" < /tmp/config.template.js > /usr/share/nginx/html/js/config.js' >> /docker-entrypoint.d/30-config.sh && \
    chmod +x /docker-entrypoint.d/30-config.sh

EXPOSE 4200
