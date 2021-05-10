Name: "{{.serviceName}}"
Mode: "dev"
CpuThreshold: 900
Host: "{{.host}}"
Port: {{.port}}
MaxBytes: 33554432 # 32MB
MaxConns: 10000
Timeout: 60000 # 60s
Log:
  ServiceName: "{{.serviceName}}"
  Mode: "file"
  Path: "logs"
  Level: "info"
  Compress: false
  KeepDays: 7
  StackCooldownMillis: 100
Apollo:
  IsEnable: false
  AppId: "${APOLLO_APP_ID}"
  Cluster: "default"
  NameSpaceNames:
    - "application"
    - "{{.serviceName}}.yaml"
  CacheDir: "etc"
  MetaAddr: "${APOLLO_META_ADDR}"
  AccessKeySecret: "${APOLLO_ACCESS_KEY_SECRET}"
