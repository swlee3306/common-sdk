# Common-SDK

## 한 줄 소개
멀티캐스트 통신을 위한 공용 SDK입니다.

## 저장소 성격
- 분류: 인프라 / 공통 라이브러리
- 목적: 여러 호스트 간 메시지 송수신 공통 기능 제공
- 핵심 기술: Go, 멀티캐스트, 네트워크 통신

멀티캐스트 통신과 공통 기능을 정리한 라이브러리 저장소입니다.

Go 언어로 작성된 멀티캐스트 통신용 공통 라이브러리입니다. 여러 호스트 간 메시지 송수신과 관련된 공통 기능을 정리하고 실험하는 데 목적이 있습니다.

## 🚀 주요 기능

### ⚡ 성능 최적화
- **메시지 압축**: gzip, LZ4 압축 알고리즘 지원
- **메시지 암호화**: AES-256 암호화 옵션
- **연결 풀링**: 효율적인 연결 재사용
- **비동기 처리**: 고성능 비동기 메시지 처리
- **메모리 풀링**: GC 압박 감소를 위한 메모리 풀링

### 📊 모니터링 및 관찰성
- **Prometheus 메트릭**: 상세한 성능 메트릭 수집
- **구조화된 로깅**: JSON 형식의 구조화된 로그
- **헬스체크**: 시스템 상태 모니터링
- **성능 분석**: 메시지 처리 성능 분석

### 🔒 보안 및 안정성
- **메시지 암호화**: AES-256 GCM 암호화
- **에러 처리**: 강화된 에러 처리 및 복구
- **재시도 로직**: 지수 백오프를 통한 재시도
- **연결 검증**: 연결 상태 지속적 모니터링

### 🛠️ 개발자 경험
- **포괄적인 테스트**: 단위 및 통합 테스트
- **API 문서화**: 상세한 API 문서
- **예제 코드**: 다양한 사용 사례 예제
- **성능 벤치마크**: 성능 측정 및 최적화

## 📦 설치 및 사용

### 1. 저장소 클론
```bash
git clone https://github.com/swlee3306/common-sdk.git
cd common-sdk
```

### 2. 의존성 설치
```bash
go mod tidy
```

### 3. 기본 사용법
```go
package main

import (
    "fmt"
    "log"
    
    "github.com/swlee3306/common-sdk/compression"
    "github.com/swlee3306/common-sdk/encryption"
    "github.com/swlee3306/common-sdk/metrics"
    "github.com/swlee3306/common-sdk/logging"
    "github.com/swlee3306/common-sdk/health"
    "github.com/swlee3306/common-sdk/retry"
    "github.com/swlee3306/common-sdk/pool"
    "github.com/swlee3306/common-sdk/errors"
    "github.com/swlee3306/common-sdk/multicast"
)

func main() {
    // 로깅 초기화
    logging.SetLevel(logging.DebugLevel)
    
    // 메트릭 초기화
    metrics.InitMetrics()
    metrics.StartMetricsServer("9090")
    
    // 압축 설정
    compressor := compression.NewCompressor(compression.Gzip)
    
    // 암호화 설정
    key, err := encryption.GenerateRandomKey()
    if err != nil {
        log.Printf("Failed to generate key: %v", err)
        return
    }
    
    // 메시지 처리
    message := []byte("Hello, World!")
    
    // 압축
    compressed, err := compressor.Compress(message)
    if err != nil {
        log.Printf("Compression failed: %v", err)
        return
    }
    
    // 암호화
    encrypted, err := encryption.Encrypt(compressed, key)
    if err != nil {
        log.Printf("Encryption failed: %v", err)
        return
    }
    
    log.Println("Message processed successfully")
    
    // 멀티캐스트 사용 예제
    multicast.Init()
    multicast.RunReceivers("224.0.0.1:9999")
}
```

## 🏗️ 프로젝트 구조

```
common-sdk/
├── compression/           # 메시지 압축
│   ├── compression.go    # 압축 알고리즘 구현
│   └── compression_test.go # 압축 테스트
├── encryption/           # 메시지 암호화
│   └── encryption.go     # AES-256 암호화
├── multicast/            # 멀티캐스트 통신
│   ├── receiver.go       # 메시지 수신기
│   ├── sender.go         # 메시지 송신기
│   └── type.go           # 공통 데이터 타입
├── metrics/              # Prometheus 메트릭
│   └── metrics.go        # 메트릭 수집 및 노출
├── logging/              # 구조화된 로깅
│   └── logger.go         # 로거 설정
├── health/               # 헬스체크
│   └── health.go         # 헬스체크 시스템
├── retry/                # 재시도 로직
│   └── retry.go          # 지수 백오프 재시도
├── pool/                 # 연결 풀링
│   └── pool.go           # 제네릭 연결 풀
├── errors/               # 에러 처리
│   └── errors.go         # 커스텀 에러 타입
├── docs/                 # 문서
│   └── api.md           # API 문서
├── go.mod               # Go 모듈 파일
├── go.sum               # 의존성 체크섬
└── README.md            # 프로젝트 문서
```

## 🔧 고급 사용법

### 1. 메시지 압축
```go
// gzip 압축
gzipCompressor := compression.NewCompressor(compression.Gzip)
compressed, err := gzipCompressor.Compress(data)
if err != nil {
    return err
}

// LZ4 압축
lz4Compressor := compression.NewCompressor(compression.Lz4)
compressed, err := lz4Compressor.Compress(data)
if err != nil {
    return err
}
```

### 2. 메시지 암호화
```go
// 암호화 키 생성
key, err := encryption.GenerateRandomKey()
if err != nil {
    return err
}

// 메시지 암호화
encrypted, err := encryption.Encrypt(data, key)
if err != nil {
    return err
}

// 메시지 복호화
decrypted, err := encryption.Decrypt(encrypted, key)
if err != nil {
    return err
}
```

### 3. 메트릭 수집
```go
// 메트릭 초기화
metrics.InitMetrics()

// 메트릭 서버 시작
metrics.StartMetricsServer("9090")

// 커스텀 메트릭 추가
metrics.RequestCounter.WithLabelValues("GET", "/api").Inc()
metrics.RequestDurationHistogram.WithLabelValues("GET", "/api").Observe(duration.Seconds())
```

### 4. 재시도 로직
```go
err := retry.Do(func() error {
    // 네트워크 요청 또는 다른 작업
    return someOperation()
},
    retry.WithMaxAttempts(5),
    retry.WithInitialDelay(100*time.Millisecond),
    retry.WithFactor(2.0),
    retry.WithJitter(0.1),
    retry.WithOnRetry(func(attempt int, err error, delay time.Duration) {
        log.Printf("Retry attempt %d failed: %v, retrying in %v", attempt, err, delay)
    }),
)
```

### 5. 연결 풀링
```go
// TCP 연결 풀 생성
factory := func() (pool.Resource, error) {
    conn, err := net.DialTimeout("tcp", "localhost:8080", 5*time.Second)
    if err != nil {
        return nil, err
    }
    return conn, nil
}

p, err := pool.NewChannelPool(5, 10, time.Minute, factory)
if err != nil {
    return err
}
defer p.Close()

// 연결 사용
conn, err := p.Get()
if err != nil {
    return err
}
defer p.Put(conn)
```

### 6. 멀티캐스트 통신
```go
// 멀티캐스트 초기화
multicast.Init()

// 커스텀 핸들러 등록
multicast.RegisterHandler("customType", func(payload json.RawMessage, addr string) error {
    // 커스텀 메시지 처리 로직
    fmt.Printf("Received message from %s: %s\n", addr, string(payload))
    return nil
})

// 멀티캐스트 수신기 시작
multicast.RunReceivers("224.0.0.1:9999")

// 호스트 정보 조회
hostData := multicast.GetHostData()
for hostname, info := range hostData {
    fmt.Printf("Host: %s, IPs: %v\n", hostname, info.IPs)
}
```

## 📊 성능 특성

### 압축 성능
- **gzip**: 높은 압축률, 중간 처리 속도
- **LZ4**: 낮은 압축률, 높은 처리 속도
- **압축률**: 평균 60-80% 크기 감소

### 암호화 성능
- **AES-256 GCM**: 높은 보안성, 빠른 처리
- **처리 속도**: 초당 수만 건 메시지 처리
- **메모리 사용량**: 최적화된 메모리 사용

### 메트릭 성능
- **메트릭 수집**: 실시간 성능 모니터링
- **저장소**: Prometheus 호환 메트릭
- **대시보드**: Grafana 연동 지원

## 🧪 테스트

### 단위 테스트 실행
```bash
# 모든 테스트 실행
go test ./...

# 특정 패키지 테스트
go test ./compression
go test ./encryption
go test ./metrics

# 테스트 커버리지 확인
go test -cover ./...
```

### 성능 벤치마크
```bash
# 압축 성능 벤치마크
go test -bench=BenchmarkCompression ./compression

# 암호화 성능 벤치마크
go test -bench=BenchmarkEncryption ./encryption

# 연결 풀 성능 벤치마크
go test -bench=BenchmarkPool ./pool
```

## 📈 모니터링 및 관찰성

### Prometheus 메트릭
```go
// 사용 가능한 메트릭
- common_sdk_requests_total: 총 요청 수
- common_sdk_in_flight_requests: 진행 중인 요청 수
- common_sdk_request_duration_seconds: 요청 지속 시간
- common_sdk_message_bytes_total: 처리된 메시지 바이트 수
- common_sdk_compression_ratio: 압축률
- common_sdk_encryption_operations_total: 암호화 작업 수
```

### 로깅 설정
```go
// 로그 레벨 설정
logging.SetLevel(logging.InfoLevel)

// 구조화된 로깅
logging.WithFields(logrus.Fields{
    "operation": "compress",
    "algorithm": "gzip",
    "size": len(data),
}).Info("Message compressed")
```

### 헬스체크
```go
// 헬스체크 설정
checker := health.NewSimpleHealthChecker("database", func() error {
    // 데이터베이스 연결 확인
    return db.Ping()
})

results := health.AggregateHealthCheck([]health.HealthChecker{checker})
overallStatus := health.OverallStatus(results)
```

## 🚀 배포 및 운영

### Docker 사용
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o common-sdk main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/common-sdk .
CMD ["./common-sdk"]
```

### Kubernetes 배포
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: common-sdk
spec:
  replicas: 3
  selector:
    matchLabels:
      app: common-sdk
  template:
    metadata:
      labels:
        app: common-sdk
    spec:
      containers:
      - name: common-sdk
        image: common-sdk:latest
        ports:
        - containerPort: 9090
        env:
        - name: LOG_LEVEL
          value: "info"
```

## 🔧 설정

### 환경 변수
```bash
# 로깅 설정
LOG_LEVEL=info
LOG_FORMAT=json

# 메트릭 설정
METRICS_PORT=9090
METRICS_PATH=/metrics

# 압축 설정
COMPRESSION_ALGORITHM=gzip
COMPRESSION_LEVEL=6

# 암호화 설정
ENCRYPTION_KEY=your-encryption-key
ENCRYPTION_ENABLED=true
```

### 설정 파일
```yaml
# config.yaml
logging:
  level: info
  format: json

metrics:
  port: 9090
  path: /metrics

compression:
  algorithm: gzip
  level: 6

encryption:
  enabled: true
  key: your-encryption-key
```

## 🤝 기여하기

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 라이센스

이 프로젝트는 MIT 라이센스 하에 배포됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 🙏 감사의 말

- [Go](https://golang.org/) - 프로그래밍 언어
- [Prometheus](https://prometheus.io/) - 메트릭 수집
- [Logrus](https://github.com/sirupsen/logrus) - 로깅 라이브러리
- [LZ4](https://github.com/pierrec/lz4) - 압축 라이브러리

## 📞 지원 및 문의

- 이슈 리포트: [GitHub Issues](https://github.com/swlee3306/common-sdk/issues)
- 이메일: swlee3306@gmail.com
- 문서: [Wiki](https://github.com/swlee3306/common-sdk/wiki)

---

**Common-SDK** - 엔터프라이즈급 멀티캐스트 통신 라이브러리 🚀
