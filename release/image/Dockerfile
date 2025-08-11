# golang极简镜像，编译拿到服务端产物
FROM golang:1.24-alpine AS backend_builder

# 1. 安装git(用于下载go mod依赖)
RUN apk add --no-cache git

WORKDIR /coze-loop

# 2. 下载并缓存go mod依赖
COPY ./backend/go.mod ./backend/go.sum /coze-loop/src/backend/
RUN go mod download -C ./src/backend -x

# 3. 编译服务端
COPY ./backend/ /coze-loop/src/backend/
RUN mkdir -p ./bin && \
    go -C /coze-loop/src/backend build -buildvcs=false -o /coze-loop/bin/main "./cmd"

# node极简镜像(含Node.js+npm)额外安装rush后编译拿到前端产物
FROM node:20.13.1-alpine AS frontend_builder

# 1. 安装 curl 和 bash 之类的基本工具（alpine）
RUN apk add --no-cache bash

# 2. 安装 pnpm 和 rush
RUN corepack enable && \
    corepack prepare pnpm@8.15.8 --activate && \
    npm install -g @microsoft/rush@5.147.1

WORKDIR /coze-loop

# 3. 编译前端
COPY ./frontend/ /coze-loop/src/frontend/
COPY ./common/ /coze-loop/src/common/
COPY ./rush.json /coze-loop/src/rush.json
RUN mkdir -p /coze-loop/resources && \
    sh /coze-loop/src/frontend/apps/cozeloop/build-artifact.sh /coze-loop/resources

# 最终镜像(coze-loop)，极简镜像
FROM alpine:3.22.0

WORKDIR /coze-loop

# 1. 抽产物
COPY --from=backend_builder /coze-loop/bin/main /coze-loop/bin/main
COPY --from=frontend_builder /coze-loop/resources/ /coze-loop/resources/

# 2. 抽配置
COPY ./conf/ /coze-loop/conf/