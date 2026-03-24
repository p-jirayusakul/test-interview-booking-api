# test-interview-booking-api

ออกแบบ Microservice จอง Event ในระบบ booking API
- จำกัดจำนวนคนที่จองได้
- มี Waitlist สำรองคนที่จองไม่ทัน
- concurrency สูง (หลาย user จองพร้อมกัน)
- ป้องกัน double-booking
- ป้องกัน race condition


## Tech Stack
- Go 1.26
- Echo
- PostgreSQL
- Docker / Docker Compose

## Clean Architecture
โปรเจกต์นี้แยกโค้ดตามแนวคิด Clean Architecture:
- `domain` เก็บ entity และ validation
- `usecase` เก็บ business flow
- `delivery` รับ request และส่ง response
- `infrastructure` คุยกับฐานข้อมูล
- `bootstrap` ทำหน้าที่ประกอบ dependency ทั้งระบบ

## How to run with Docker

### 1. Start services
```bash
docker compose up --build
```

### 2. API endpoint
โดยค่าเริ่มต้น API จะรันที่:
```text
http://localhost:8080/api/v1
```

## API Documentation

Wiki:
- [Home](https://github.com/p-jirayusakul/test-interview-booking-api/wiki)
- [Search Event](https://github.com/p-jirayusakul/test-interview-booking-api/wiki/API-Spec:-Search-Event)
- [จอง Event](https://github.com/p-jirayusakul/test-interview-booking-api/wiki/API-Spec:-%E0%B8%88%E0%B8%AD%E0%B8%87-Event)


Swagger UI:
```text
http://localhost:8080/swagger/index.html
```

## Testing

### Run unit tests
```bash
go test ./internal/usecase -run ^TestEventsUseCase
```

### Run Concurrent tests
```bash
go test ./internal/usecase -run TestEventsConcurrent
```

### K6 Load test
- Total Requests: 100
- Concurrent Users: 50

Results:
- Confirmed: 50
- Waitlist: 5
- Rejected: 45

Performance:
- Avg latency: ~109ms
- P95 latency: ~231ms

Observations:
- No overbooking occurred
- System maintained data consistency under concurrent load
- Increased latency is expected due to row-level locking (SELECT FOR UPDATE)

```bash
TOTAL RESULTS 

    checks_total.......: 100     390.016818/s
    checks_succeeded...: 100.00% 100 out of 100
    checks_failed......: 0.00%   0 out of 100

    ✓ status is 201 or 409

    HTTP
    http_req_duration..............: avg=109.17ms min=4.86ms med=95.38ms max=246.03ms p(90)=222.84ms p(95)=231.93ms
      { expected_response:true }...: avg=64.87ms  min=4.86ms med=52.2ms  max=195.06ms p(90)=145.34ms p(95)=166.48ms
    http_req_failed................: 45.00% 45 out of 100
    http_reqs......................: 100    390.016818/s

    EXECUTION
    iteration_duration.............: avg=113.98ms min=5.56ms med=99.19ms max=255.89ms p(90)=232.14ms p(95)=241.15ms
    iterations.....................: 100    390.016818/s

    NETWORK
    data_received..................: 27 kB  105 kB/s
    data_sent......................: 19 kB  76 kB/s
```

## Notes
- โปรเจกต์นี้แยก layer ตามแนวคิด Clean Architecture
- ใช้ PostgreSQL เป็นฐานข้อมูลหลัก
- ใช้ Docker Compose สำหรับรันระบบทั้งหมด
