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

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
