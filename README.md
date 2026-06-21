# Hướng Dẫn Chạy Hệ Thống (Gateway Routing & PDU Session)

## 1. Cấu hình chế độ trễ (Delay Mode)
Mở [docker-compose.yml](docker-compose.yml), tại mục `pdu-session` > `environment`, cấu hình biến `DELAY_MODE`:
* **Trễ cố định 15s**: `DELAY_MODE=fixed`
* **Trễ ngẫu nhiên (mod 20s)**: `DELAY_MODE=random`

---

## 2. Khởi chạy hệ thống (Docker Compose)
Mở terminal tại thư mục gốc và chạy lệnh:
```bash
docker compose up --build
```
* **Giám sát tài nguyên (CPU < 65%)**: Mở terminal mới chạy lệnh `docker stats` để theo dõi.

---

## 3. Khởi chạy Client bắn tải (auto-request)
Mở terminal mới tại thư mục `auto-request/` và chạy lệnh:
```bash
go run .
```

Chọn chế độ tương tác:
* **Chế độ 1**: Tự động bắn 500 requests theo chu kỳ ngẫu nhiên 5-7s. (Nhấn **ESC** để dừng).
* **Chế độ 2**: Nhập số lượng request bất kỳ từ bàn phím $\rightarrow$ nhấn **Enter** để bắn. (Nhấn **ESC** hoặc gõ `exit` để dừng).



