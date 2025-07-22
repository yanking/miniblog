# ================== 构建阶段 ==================
FROM golang:1.24 AS  builder

WORKDIR /opt/miniblog

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 复制项目源代码
COPY . .

RUN make build BINS=mb-apiserver

# ================== 运行阶段 ==================

# 使用更小的镜像作为最终镜像
FROM alpine:latest

# 安装运行程序所需的依赖
RUN apk --no-cache add ca-certificates tzdata

## 设置时区
ENV TZ=Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /opt/miniblog

# 创建必要的目录结构
RUN mkdir -p /opt/miniblog/configs

# 复制编译后的二进制文件到容器
COPY --from=builder /opt/miniblog/_output/mb-apiserver /opt/miniblog/bin/
# 设置二进制文件的执行权限
RUN chmod +x /opt/miniblog/bin/mb-apiserver

# 复制配置文件到容器
COPY ./configs/ /opt/miniblog/configs/

# 设置容器启动时执行的命令
ENTRYPOINT ["/opt/miniblog/bin/mb-apiserver","-c","/opt/miniblog/configs/mb-apiserver.yaml"]


# docker build -t localhost:5000/miniblog-apiserver:latest . -f ./build/docker/mb-apiserver/Dockerfile