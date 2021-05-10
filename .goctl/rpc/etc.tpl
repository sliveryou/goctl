Name: "{{.serviceName}}.rpc"
Mode: "dev"
ListenOn: "0.0.0.0:46600"
Timeout: 60000 # 60s
Log:
  ServiceName: "{{.serviceName}}-service"
  Mode: "file"
  Path: "logs/{{.serviceName}}-service"
  Level: "info"
  Compress: false
  KeepDays: 7
  StackCooldownMillis: 100
Etcd:
  Hosts:
    - "localhost:2379"
  Key: "{{.serviceName}}.rpc"
Apollo:
  IsEnable: false
  AppId: "${APOLLO_APP_ID}"
  Cluster: "default"
  NameSpaceNames:
    - "application"
    - "{{.serviceName}}-service.yaml"
  CacheDir: "etc/{{.serviceName}}-service"
  MetaAddr: "${APOLLO_META_ADDR}"
  AccessKeySecret: "${APOLLO_ACCESS_KEY_SECRET}"
DB:
  User: "${DB_USER}"
  Password: "${DB_PASSWORD}"
  Host: "${DB_HOST}"
  Port: ${DB_PORT}
  Database: "${DB_DATABASE}"
  MaxIdleConns: 10
  MaxOpenConns: 20
  LogLevel: "info"
