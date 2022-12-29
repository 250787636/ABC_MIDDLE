### 使用方法

`docker构建方式`

>build目录下需要保护的文件

1. run.sh：启动脚本
2. Dockerfile：docker构建文件
3. deploy.ini：配置文件
4. resources：资源工具文件
5. abc：项目编译后文件

> 编译命令: go build -o ./build/abc
> 
> docker构建时指定image和tag
> 构建命令：docker build -t harbor.bangcle.net:8029/custom_development/tp-abc-front:ver1.0.1 .
> 
> docker其他资源文件
> 
1. /data/docker/tp-abc-backend: 后端docker映射卷目录
2. /data/docker/tp-abc-backend/logs:映射日志目录
3. /data/docker/tp-abc-backend/media:临时文件存放
4. /data/docker/tp-abc-backend/deploy.ini:配置文件映射

> docker-compose目录结构：
1. /data/compose/abc/.env:存放需要启动的镜像tag，前后端
2. /data/compose/docker-compose.yml:docker配置文件

> 启动命令:docker-compose up -d