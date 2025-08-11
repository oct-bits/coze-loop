# golang极简镜像，编译拿到服务端产物
FROM golang:1.24-alpine AS backend_builder

# 1. 安装git(用于下载go mod依赖)
RUN apk add --no-cache git

# 2. 安装dlv(用于调试)
RUN go install "github.com/go-delve/delve/cmd/dlv@v1.25.1"

WORKDIR /coze-loop

# 2. 下载并缓存go mod依赖
COPY ./backend/go.mod ./backend/go.sum /coze-loop/src/backend/
RUN go mod download -C ./src/backend -x

# 3. 编译服务端
COPY ./backend/ /coze-loop/src/backend/
RUN mkdir -p ./bin && \
    go -C /coze-loop/src/backend build -gcflags="all=-N -l" -buildvcs=false -o /coze-loop/bin/main "./cmd"

# 最终镜像(coze-loop)，极简镜像
FROM compose-cn-beijing.cr.volces.com/coze/coze-loop:latest

WORKDIR /coze-loop

# 抽产物
COPY --from=backend_builder /coze-loop/bin/main /coze-loop/bin/main
COPY --from=backend_builder /go/bin/dlv /usr/local/bin/dlv