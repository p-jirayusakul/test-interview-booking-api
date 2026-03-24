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

## Notes
- โปรเจกต์นี้แยก layer ตามแนวคิด Clean Architecture
- ใช้ PostgreSQL เป็นฐานข้อมูลหลัก
- ใช้ Docker Compose สำหรับรันระบบทั้งหมด
