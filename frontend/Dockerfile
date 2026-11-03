# 源代码 → npm run build → dist 静态文件 → Nginx 提供访问

# 1.构建dist文件

FROM docker.m.daocloud.io/library/node:22-alpine AS builder

WORKDIR /app
ENV npm_config_registry=https://registry.npmmirror.com

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM docker.m.daocloud.io/library/alpine:3.21

# Build artifact only. The separate Nginx container mounts /dist from the host.
COPY --from=builder /app/dist /dist
