FROM nginxinc/nginx-unprivileged:1.30-alpine
ENV BACKEND_PORT=8080 NGINX_ENVSUBST_FILTER=^BACKEND_PORT$
COPY nginx.conf /etc/nginx/templates/default.conf.template
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
