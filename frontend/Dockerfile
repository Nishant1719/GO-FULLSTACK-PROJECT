# AWS Public ECR mirror of Docker Hub official images (fewer 504s / rate limits in CI than registry-1.docker.io)
FROM public.ecr.aws/docker/library/node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM public.ecr.aws/docker/library/nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
