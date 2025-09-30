# Common SDK - 멀티캐스트 통신 라이브러리

## 개요
Common SDK는 Go 언어로 작성된 멀티캐스트 통신을 위한 공통 라이브러리입니다. 네트워크 상의 여러 호스트 간에 효율적인 메시지 전송과 수신을 제공합니다.

## 주요 기능
- **멀티캐스트 송신/수신**: UDP 멀티캐스트를 통한 메시지 전송
- **메시지 분할**: 큰 메시지를 여러 프래그먼트로 분할하여 전송
- **자동 재조립**: 수신된 프래그먼트를 자동으로 재조립
- **핸들러 시스템**: 메시지 타입별 커스텀 핸들러 등록
- **호스트 정보 관리**: 네트워크 상의 호스트 정보 수집 및 관리

## 아키텍처

### 핵심 컴포넌트
- **Receiver**: 멀티캐스트 메시지 수신 및 처리
- **Sender**: 멀티캐스트 메시지 전송
- **Type**: 공통 데이터 구조 정의

### 메시지 처리 플로우
1. 메시지 타입별 핸들러 등록
2. 멀티캐스트 주소에서 메시지 수신
3. 프래그먼트 재조립
4. 등록된 핸들러로 메시지 전달

## API 사용법

### 초기화
```go
import "common-sdk/multicast"

func main() {
    multicast.Init()
    multicast.RunReceivers("224.0.0.1:9999")
}
```

### 커스텀 핸들러 등록
```go
multicast.RegisterHandler("customType", func(payload json.RawMessage, addr string) error {
    // 커스텀 메시지 처리 로직
    return nil
})
```

### 호스트 정보 조회
```go
hostData := multicast.GetHostData()
for hostname, info := range hostData {
    fmt.Printf("Host: %s, IPs: %v\n", hostname, info.IPs)
}
```

## 메시지 타입

### 지원되는 기본 메시지 타입
- **hostinfoSend**: 호스트 정보 요청 트리거
- **hostinfo**: 호스트 정보 전송

### 메시지 구조
```go
type GenericMessage struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}

type HostInfoReceiver struct {
    Version      string   `json:"version"`
    BuildDate    string   `json:"buildDate"`
    Revision     string   `json:"revision"`
    Hostname     string   `json:"hostname"`
    IPs          []string `json:"ips"`
    Endpoint     string   `json:"endpoint"`
    EndpointPort int      `json:"endpointPort"`
}
```

## 네트워크 설정

### 멀티캐스트 주소
- 기본 주소: `224.0.0.1:9999`
- 포트: 9999 (설정 가능)

### 네트워크 요구사항
- 멀티캐스트 지원 네트워크 인터페이스
- 방화벽에서 멀티캐스트 트래픽 허용

## 성능 특성
- **프래그먼트 크기**: 최대 1500 바이트
- **타임아웃**: 15초 (미완성 메시지 정리)
- **버퍼 크기**: 2048 바이트

## 사용 사례
- **서비스 디스커버리**: 네트워크 상의 서비스 자동 발견
- **클러스터 통신**: 분산 시스템 간 상태 동기화
- **모니터링**: 다중 호스트 상태 수집

## 개발자 정보
- **언어**: Go
- **의존성**: 표준 라이브러리만 사용
- **라이선스**: MIT
