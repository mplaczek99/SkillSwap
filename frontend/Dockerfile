# Build stage
FROM node:18-alpine as build-stage
WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Build the application with API URL from build arg
ARG VUE_APP_API_URL=http://localhost:8080
ENV VUE_APP_API_URL=${VUE_APP_API_URL}

RUN echo "Building with API URL: $VUE_APP_API_URL"
RUN npm run build

# Production stage
FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Create nginx config for SPA routing
RUN echo 'server { \
    listen 80; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html; \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

# Create the entrypoint script - properly with multiple RUN commands to avoid newline issues
RUN echo '#!/bin/sh' > /docker-entrypoint.sh
RUN echo 'if [ ! -z "$RUNTIME_API_URL" ]; then' >> /docker-entrypoint.sh
RUN echo '    echo "Updating API URL to: $RUNTIME_API_URL"' >> /docker-entrypoint.sh
RUN echo '    find /usr/share/nginx/html -type f -name "*.js" -exec sed -i "s|http://localhost:8080|$RUNTIME_API_URL|g" {} \\;' >> /docker-entrypoint.sh
RUN echo 'fi' >> /docker-entrypoint.sh
RUN echo 'exec nginx -g "daemon off;"' >> /docker-entrypoint.sh

RUN chmod +x /docker-entrypoint.sh

EXPOSE 80
CMD ["/docker-entrypoint.sh"]
