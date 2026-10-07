# syntax=docker/dockerfile:1
FROM node:24-alpine AS builder
WORKDIR /app
COPY frontend/vue-ui/package*.json ./
RUN --mount=type=cache,target=/root/.npm npm ci
COPY frontend/vue-ui/ ./
RUN npm run build

FROM nginxinc/nginx-unprivileged:1.30-alpine
ENV BACKEND_PORT=8080 NGINX_ENVSUBST_FILTER=^BACKEND_PORT$
COPY docker/nginx/nginx.conf /etc/nginx/templates/default.conf.template
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
