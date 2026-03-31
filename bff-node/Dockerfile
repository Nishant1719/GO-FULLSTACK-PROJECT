# AWS Public ECR mirror of Docker Hub official images (more reliable pulls in CI)
FROM public.ecr.aws/docker/library/node:22-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev
COPY src ./src
EXPOSE 3000
CMD ["node", "src/index.js"]
