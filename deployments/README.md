# TiRelease 部署

-   `secrets/`: Secrets，使用 sealed-secret 加密
    -   `sealedsecretconfig.yaml`：敏感配置，包括 DSN, github token, feishu token
-   `tirelease`:
    -   `config.yaml`：配置文件，用于指定运行参数
    -   `deployment.yaml`：TiRelease 服务,包含服务发现配置
    -   `canary-deployment.yaml`：TiRelease 伪 canary 服务，部署后可通过 http://tirelease.pingcap.net:30751/ 访问 nightly 环境

## 构建镜像

```sh
make image && make image-push
```

## 升级服务

```sh
make upgrade-service
```
