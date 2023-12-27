Name: {{.serviceName}}
Mode: dev
CpuThreshold: 900
Host: {{.host}}
Port: {{.port}}
MaxBytes: 33554432 # 32MB
MaxConns: 10000
Timeout: 60000 # 60s
Log:
  ServiceName: {{.serviceName}}
  Mode: file
  Path: logs
  Level: "info"
  Compress: false
  KeepDays: 7
  StackCooldownMillis: 100
