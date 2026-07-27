# Báo Cáo Tổng Hợp & Log Kết Quả Test Tải Đồng Thời (Benchmark Results)

---

## I. Báo Cáo Tổng Hợp & Đánh Giá

### 📊 Bảng So Sánh Các Đợt Test

| Lần Test | Connections | Workers | Thời gian | Tổng Request | Thành công | Thất bại | TPS Thành công | Tỉ lệ lỗi |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| **Lần 1** | **10** | **100** | **1m 19.72s** | **1,000,000** | **1,000,000** | **0** | **12,543.01 req/s** | **0%** |
| **Lần 2** | 20 | 500 | 1m 26.27s | 1,000,000 | 1,000,000 | 0 | **11,591.02** req/s | 0% |
| **Lần 3** | 128 | 2564 | 1m 42.57s | 1,000,000 | 1,000,000 | 0 | **9,749.44** req/s | 0% |
| **Lần 4** | 50 | 500 | 1m 30.53s | 1,000,000 | 1,000,000 | 0 | **11,045.49** req/s | 0% |
| **Lần 5** | 20 | 600 | 1m 24.61s | 1,000,000 | 1,000,000 | 0 | **11,818.84** req/s | 0% |

---

### 🔍 Chi Tiết Từng Lần Test

#### 1. Lần Test 1 (Hiệu năng cao nhất ⭐)
- **Cấu hình**: `Connections: 10`, `Workers: 100`
- **Thời gian chạy**: 1m19.725s (79.73s)
- **Tổng Request**: 1,000,000 | **Thành công**: 1,000,000 (100%) | **Thất bại**: 0 (0%)
- **TPS gửi đi**: 12,543.01 req/s | **TPS thành công**: 12,543.01 req/s
- **Phân phối tải (3 Instances)**:
  - `pdu-session-1`: 333,334 requests (33.34%)
  - `pdu-session-2`: 333,333 requests (33.33%)
  - `pdu-session-4`: 333,333 requests (33.33%)

#### 2. Lần Test 2
- **Cấu hình**: `Connections: 20`, `Workers: 500`
- **Thời gian chạy**: 1m26.273s (86.27s)
- **Tổng Request**: 1,000,000 | **Thành công**: 1,000,000 (100%) | **Thất bại**: 0 (0%)
- **TPS gửi đi**: 11,591.02 req/s | **TPS thành công**: 11,591.02 req/s
- **Phân phối tải (3 Instances)**:
  - `pdu-session-1`: 333,333 requests (33.33%)
  - `pdu-session-2`: 333,34 requests (33.34%)
  - `pdu-session-4`: 333,333 requests (33.33%)

#### 3. Lần Test 3
- **Cấu hình**: `Connections: 128`, `Workers: 2564`
- **Thời gian chạy**: 1m42.569s (102.57s)
- **Tổng Request**: 1,000,000 | **Thành công**: 1,000,000 (100%) | **Thất bại**: 0 (0%)
- **TPS gửi đi**: 9,749.44 req/s | **TPS thành công**: 9,749.44 req/s
- **Phân phối tải (3 Instances)**:
  - `pdu-session-1`: 333,333 requests (33.33%)
  - `pdu-session-2`: 333,333 requests (33.33%)
  - `pdu-session-4`: 333,334 requests (33.34%)

#### 4. Lần Test 4
- **Cấu hình**: `Connections: 50`, `Workers: 500`
- **Thời gian chạy**: 1m30.534s (90.53s)
- **Tổng Request**: 1,000,000 | **Thành công**: 1,000,000 (100%) | **Thất bại**: 0 (0%)
- **TPS gửi đi**: 11,045.49 req/s | **TPS thành công**: 11,045.49 req/s
- **Phân phối tải (3 Instances)**:
  - `pdu-session-1`: 333,34 requests (33.34%)
  - `pdu-session-2`: 333,333 requests (33.33%)
  - `pdu-session-4`: 333,333 requests (33.33%)

#### 5. Lần Test 5
- **Cấu hình**: `Connections: 20`, `Workers: 600`
- **Thời gian chạy**: 1m24.610s (84.61s)
- **Tổng Request**: 1,000,000 | **Thành công**: 1,000,000 (100%) | **Thất bại**: 0 (0%)
- **TPS gửi đi**: 11,818.84 req/s | **TPS thành công**: 11,818.84 req/s
- **Phân phối tải (3 Instances)**:
  - `pdu-session-1`: 333,333 requests (33.33%)
  - `pdu-session-2`: 333,334 requests (33.34%)
  - `pdu-session-4`: 333,333 requests (33.33%)

---

### 📌 Đánh Giá & Nhận Xét
1. **Khả năng Load Balancing**: Phân phối tải hoàn hảo giữa 3 instance active (~33.33% cho mỗi node).
2. **Điểm ngọt tối ưu (Optimal Point)**: Cấu hình `10 connections / 100 workers` đạt thông lượng cao nhất (**~12,543 TPS**).

---

## II. Log Nguyên Bản Từ Terminal (Raw Terminal Output)

```text

Đường dẫn đích (Target): http://gateway:8080/nsmf-pdusession/v1/sm-contexts
HTTP Client Timeout: 1m30s

Nhập số lượng request muốn bắn: 500000
Nhập số lượng connection (mặc định 1): 50
Nhập số lượng Worker song song (mặc định 100): 1000
2026/07/27 03:52:12 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 03:52:12 Thời gian chạy: 49.375298548s
2026/07/27 03:52:12 Tổng số Request gửi đi: 500000
2026/07/27 03:52:12 Thành công: 500000
2026/07/27 03:52:12 Thất bại: 0
2026/07/27 03:52:12 TPS gửi đi: 10126.52 req/s
2026/07/27 03:52:12 TPS thành công: 10126.52 req/s
2026/07/27 03:52:12 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 03:52:12 Instance pdu-session-1 xử lý: 125000 requests (25.00%)
2026/07/27 03:52:12 Instance pdu-session-4 xử lý: 125000 requests (25.00%)
2026/07/27 03:52:12 Instance pdu-session-2 xử lý: 125000 requests (25.00%)
2026/07/27 03:52:12 Instance pdu-session-3 xử lý: 125000 requests (25.00%)

Nhập số lượng request muốn bắn: 1000000
Nhập số lượng connection (mặc định 1): 500
Nhập số lượng Worker song song (mặc định 100): 100
2026/07/27 03:54:14 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 03:54:14 Thời gian chạy: 1m29.462054825s
2026/07/27 03:54:14 Tổng số Request gửi đi: 1000000
2026/07/27 03:54:14 Thành công: 1000000
2026/07/27 03:54:14 Thất bại: 0
2026/07/27 03:54:14 TPS gửi đi: 11177.92 req/s
2026/07/27 03:54:14 TPS thành công: 11177.92 req/s
2026/07/27 03:54:14 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 03:54:14 Instance pdu-session-1 xử lý: 250000 requests (25.00%)
2026/07/27 03:54:14 Instance pdu-session-2 xử lý: 250000 requests (25.00%)
2026/07/27 03:54:14 Instance pdu-session-4 xử lý: 250000 requests (25.00%)
2026/07/27 03:54:14 Instance pdu-session-3 xử lý: 250000 requests (25.00%)

Nhập số lượng request muốn bắn: 1000000
Nhập số lượng connection (mặc định 1): 10
Nhập số lượng Worker song song (mặc định 100): 100
2026/07/27 03:55:50 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 03:55:50 Thời gian chạy: 1m21.17783548s
2026/07/27 03:55:50 Tổng số Request gửi đi: 1000000
2026/07/27 03:55:50 Thành công: 1000000
2026/07/27 03:55:50 Thất bại: 0
2026/07/27 03:55:50 TPS gửi đi: 12318.63 req/s
2026/07/27 03:55:50 TPS thành công: 12318.63 req/s
2026/07/27 03:55:50 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 03:55:50 Instance pdu-session-4 xử lý: 250000 requests (25.00%)
2026/07/27 03:55:50 Instance pdu-session-1 xử lý: 250000 requests (25.00%)
2026/07/27 03:55:50 Instance pdu-session-3 xử lý: 250000 requests (25.00%)
2026/07/27 03:55:50 Instance pdu-session-2 xử lý: 250000 requests (25.00%)

Nhập số lượng request muốn bắn: 1000000
Nhập số lượng connection (mặc định 1): 10
Nhập số lượng Worker song song (mặc định 100): 200
2026/07/27 03:58:07 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 03:58:07 Thời gian chạy: 1m19.012660118s
2026/07/27 03:58:07 Tổng số Request gửi đi: 1000000
2026/07/27 03:58:07 Thành công: 1000000
2026/07/27 03:58:07 Thất bại: 0
2026/07/27 03:58:07 TPS gửi đi: 12656.20 req/s
2026/07/27 03:58:07 TPS thành công: 12656.20 req/s
2026/07/27 03:58:07 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 03:58:07 Instance pdu-session-2 xử lý: 250000 requests (25.00%)
2026/07/27 03:58:07 Instance pdu-session-4 xử lý: 250000 requests (25.00%)
2026/07/27 03:58:07 Instance pdu-session-3 xử lý: 250000 requests (25.00%)
2026/07/27 03:58:07 Instance pdu-session-1 xử lý: 250000 requests (25.00%)

Nhập số lượng request muốn bắn: 10000000
Nhập số lượng connection (mặc định 1): 10
Nhập số lượng Worker song song (mặc định 100): 200
2026/07/27 04:14:24 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 04:14:24 Thời gian chạy: 15m56.434457231s
2026/07/27 04:14:24 Tổng số Request gửi đi: 10000000
2026/07/27 04:14:24 Thành công: 10000000
2026/07/27 04:14:24 Thất bại: 0
2026/07/27 04:14:24 TPS gửi đi: 10455.50 req/s
2026/07/27 04:14:24 TPS thành công: 10455.50 req/s
2026/07/27 04:14:24 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 04:14:24 Instance pdu-session-4 xử lý: 2500000 requests (25.00%)
2026/07/27 04:14:24 Instance pdu-session-3 xử lý: 2500000 requests (25.00%)
2026/07/27 04:14:24 Instance pdu-session-2 xử lý: 2500000 requests (25.00%)
2026/07/27 04:14:24 Instance pdu-session-1 xử lý: 2500000 requests (25.00%)

Nhập số lượng request muốn bắn: 1000000
Nhập số lượng connection (mặc định 1): 10
Nhập số lượng Worker song song (mặc định 100): 100 
2026/07/27 04:19:43 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 04:19:43 Thời gian chạy: 1m19.725664306s
2026/07/27 04:19:43 Tổng số Request gửi đi: 1000000
2026/07/27 04:19:43 Thành công: 1000000
2026/07/27 04:19:43 Thất bại: 0
2026/07/27 04:19:43 TPS gửi đi: 12543.01 req/s
2026/07/27 04:19:43 TPS thành công: 12543.01 req/s
2026/07/27 04:19:43 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 04:19:43 Instance pdu-session-2 xử lý: 333333 requests (33.33%)
2026/07/27 04:19:43 Instance pdu-session-1 xử lý: 333334 requests (33.33%)
2026/07/27 04:19:43 Instance pdu-session-4 xử lý: 333333 requests (33.33%)

Nhập số lượng request muốn bắn: 1000000
Nhập số lượng connection (mặc định 1): 20
Nhập số lượng Worker song song (mặc định 100): 500
2026/07/27 04:22:15 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 04:22:15 Thời gian chạy: 1m26.273687953s
2026/07/27 04:22:15 Tổng số Request gửi đi: 1000000
2026/07/27 04:22:15 Thành công: 1000000
2026/07/27 04:22:15 Thất bại: 0
2026/07/27 04:22:15 TPS gửi đi: 11591.02 req/s
2026/07/27 04:22:15 TPS thành công: 11591.02 req/s
2026/07/27 04:22:15 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 04:22:15 Instance pdu-session-1 xử lý: 333333 requests (33.33%)
2026/07/27 04:22:15 Instance pdu-session-4 xử lý: 333333 requests (33.33%)
2026/07/27 04:22:15 Instance pdu-session-2 xử lý: 333334 requests (33.33%)

Nhập số lượng request muốn bắn: 1000000
Nhập số lượng connection (mặc định 1): 128
Nhập số lượng Worker song song (mặc định 100): 2564
2026/07/27 04:24:29 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 04:24:29 Thời gian chạy: 1m42.569978882s
2026/07/27 04:24:29 Tổng số Request gửi đi: 1000000
2026/07/27 04:24:29 Thành công: 1000000
2026/07/27 04:24:29 Thất bại: 0
2026/07/27 04:24:29 TPS gửi đi: 9749.44 req/s
2026/07/27 04:24:29 TPS thành công: 9749.44 req/s
2026/07/27 04:24:29 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 04:24:29 Instance pdu-session-4 xử lý: 333334 requests (33.33%)
2026/07/27 04:24:29 Instance pdu-session-2 xử lý: 333333 requests (33.33%)
2026/07/27 04:24:29 Instance pdu-session-1 xử lý: 333333 requests (33.33%)

Nhập số lượng request muốn bắn: 1000000
Nhập số lượng connection (mặc định 1): 50
Nhập số lượng Worker song song (mặc định 100): 500
2026/07/27 04:26:30 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 04:26:30 Thời gian chạy: 1m30.534712639s
2026/07/27 04:26:30 Tổng số Request gửi đi: 1000000
2026/07/27 04:26:30 Thành công: 1000000
2026/07/27 04:26:30 Thất bại: 0
2026/07/27 04:26:30 TPS gửi đi: 11045.49 req/s
2026/07/27 04:26:30 TPS thành công: 11045.49 req/s
2026/07/27 04:26:30 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 04:26:30 Instance pdu-session-2 xử lý: 333333 requests (33.33%)
2026/07/27 04:26:30 Instance pdu-session-1 xử lý: 333334 requests (33.33%)
2026/07/27 04:26:30 Instance pdu-session-4 xử lý: 333333 requests (33.33%)

Nhập số lượng request muốn bắn: 1000000
Nhập số lượng connection (mặc định 1): 20
Nhập số lượng Worker song song (mặc định 100): 600
2026/07/27 04:28:42 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 04:28:42 Thời gian chạy: 1m24.61070084s
2026/07/27 04:28:42 Tổng số Request gửi đi: 1000000
2026/07/27 04:28:42 Thành công: 1000000
2026/07/27 04:28:42 Thất bại: 0
2026/07/27 04:28:42 TPS gửi đi: 11818.84 req/s
2026/07/27 04:28:42 TPS thành công: 11818.84 req/s
2026/07/27 04:28:42 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 04:28:42 Instance pdu-session-2 xử lý: 333334 requests (33.33%)
2026/07/27 04:28:42 Instance pdu-session-1 xử lý: 333333 requests (33.33%)
2026/07/27 04:28:42 Instance pdu-session-4 xử lý: 333333 requests (33.33%)

```
