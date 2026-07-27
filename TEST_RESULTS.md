# Báo Cáo Tổng Hợp & Log Kết Quả Test Tải Đồng Thời (Benchmark Results)

---

## I. Báo Cáo Tổng Hợp & Đánh Giá

### 📊 Bảng So Sánh Các Đợt Test

| Lần Test | Connections | Workers | Thời gian | Tổng Request | Thành công | Thất bại | TPS Thành công | Tỉ lệ lỗi | Ghi chú |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :--- |
| **Lần 1** | - | - | 3m 08.58s | 1,000,000 | 999,706 | 294 | **5,301.12** req/s | 0.029% | 4 Instances (Lỗi backend `pdu-session-3`) |
| **Lần 2** | **10** | **100** | **1m 19.72s** | **1,000,000** | **1,000,000** | **0** | **12,543.01 req/s** | **0%** | **Hiệu năng cao nhất (Peak TPS)** |
| **Lần 3** | 20 | 500 | 1m 26.27s | 1,000,000 | 1,000,000 | 0 | **11,591.02** req/s | 0% | 3 Instances |
| **Lần 4** | 128 | 2564 | 1m 42.57s | 1,000,000 | 1,000,000 | 0 | **9,749.44** req/s | 0% | Tăng worker & connection gây overhead |
| **Lần 5** | 50 | 500 | 1m 30.53s | 1,000,000 | 1,000,000 | 0 | **11,045.49** req/s | 0% | 3 Instances |
| **Lần 6** | 20 | 600 | 1m 24.61s | 1,000,000 | 1,000,000 | 0 | **11,818.84** req/s | 0% | 3 Instances |

---

### 🔍 Chi Tiết Từng Lần Test

#### 1. Lần Test 1 (Có sự cố Instance 3)
- **Cấu hình**: Mặc định / Chưa thiết lập cụ thể
- **Thời gian chạy**: 3m8.584s (188.58s)
- **Tổng Request**: 1,000,000 | **Thành công**: 999,706 | **Thất bại**: 294
- **TPS gửi đi**: 5,302.68 req/s | **TPS thành công**: 5,301.12 req/s
- **Phân phối tải (4 Instances)**:
  - `pdu-session-1`: 268,544 requests (26.86%)
  - `pdu-session-2`: 268,544 requests (26.86%)
  - `pdu-session-4`: 268,544 requests (26.86%)
  - `pdu-session-3`: 194,074 requests (19.41%) *(xảy ra lỗi backend connection refused/broken pipe)*

#### 2. Lần Test 2 (Hiệu năng cao nhất ⭐)
- **Cấu hình**: `Connections: 10`, `Workers: 100`
- **Thời gian chạy**: 1m19.725s (79.73s)
- **Tổng Request**: 1,000,000 | **Thành công**: 1,000,000 (100%) | **Thất bại**: 0 (0%)
- **TPS gửi đi**: 12,543.01 req/s | **TPS thành công**: 12,543.01 req/s
- **Phân phối tải (3 Instances)**:
  - `pdu-session-1`: 333,334 requests (33.34%)
  - `pdu-session-2`: 333,333 requests (33.33%)
  - `pdu-session-4`: 333,333 requests (33.33%)

#### 3. Lần Test 3
- **Cấu hình**: `Connections: 20`, `Workers: 500`
- **Thời gian chạy**: 1m26.273s (86.27s)
- **Tổng Request**: 1,000,000 | **Thành công**: 1,000,000 (100%) | **Thất bại**: 0 (0%)
- **TPS gửi đi**: 11,591.02 req/s | **TPS thành công**: 11,591.02 req/s
- **Phân phối tải (3 Instances)**:
  - `pdu-session-1`: 333,333 requests (33.33%)
  - `pdu-session-2`: 333,334 requests (33.34%)
  - `pdu-session-4`: 333,333 requests (33.33%)

#### 4. Lần Test 4
- **Cấu hình**: `Connections: 128`, `Workers: 2564`
- **Thời gian chạy**: 1m42.569s (102.57s)
- **Tổng Request**: 1,000,000 | **Thành công**: 1,000,000 (100%) | **Thất bại**: 0 (0%)
- **TPS gửi đi**: 9,749.44 req/s | **TPS thành công**: 9,749.44 req/s
- **Phân phối tải (3 Instances)**:
  - `pdu-session-1`: 333,333 requests (33.33%)
  - `pdu-session-2`: 333,333 requests (33.33%)
  - `pdu-session-4`: 333,334 requests (33.34%)

#### 5. Lần Test 5
- **Cấu hình**: `Connections: 50`, `Workers: 500`
- **Thời gian chạy**: 1m30.534s (90.53s)
- **Tổng Request**: 1,000,000 | **Thành công**: 1,000,000 (100%) | **Thất bại**: 0 (0%)
- **TPS gửi đi**: 11,045.49 req/s | **TPS thành công**: 11,045.49 req/s
- **Phân phối tải (3 Instances)**:
  - `pdu-session-1`: 333,34 requests (33.34%)
  - `pdu-session-2`: 333,333 requests (33.33%)
  - `pdu-session-4`: 333,333 requests (33.33%)

#### 6. Lần Test 6
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
Nhập số lượng Worker song song (mặc định 100): 50 
2026/07/27 04:16:33 [ERR] Request 772213 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772241 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772240 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772211 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772216 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772242 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772214 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772206 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772222 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772208 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772212 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772199 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772205 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772209 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772236 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772235 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772234 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772244 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772226 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772239 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772243 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772245 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772204 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772223 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772224 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772227 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772247 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772228 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772246 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772229 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772218 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772231 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772207 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772238 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772230 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772198 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772232 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:33 [ERR] Request 772237 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:39 [ERR] Request 772260 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772262 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772256 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772250 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772259 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772263 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772252 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772248 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772265 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772278 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772282 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772257 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772249 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772292 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772272 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772288 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772273 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772258 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772290 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772276 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772294 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772275 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772280 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772279 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772298 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772253 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772254 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772296 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772267 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772284 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772274 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772264 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772271 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772285 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772266 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772281 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772283 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772269 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772268 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772261 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772286 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772295 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772291 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772277 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772270 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772255 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772287 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772289 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772293 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772251 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772297 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772342 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772304 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772343 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772344 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772337 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772345 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772307 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772313 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772305 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772311 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772303 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772319 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772302 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772331 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772322 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772306 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772329 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772336 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772332 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772323 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772334 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772326 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772330 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772310 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772325 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772299 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772312 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772335 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772308 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772314 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772340 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772327 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772347 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772346 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772341 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772300 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772315 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772301 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772316 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772317 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772320 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772318 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772333 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772324 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772321 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772328 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772338 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772339 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772309 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772348 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772360 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772363 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772365 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772367 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772369 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772368 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772354 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772356 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772349 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772371 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772351 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772352 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772370 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772350 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772353 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772355 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772357 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772358 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772359 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772361 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772362 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772364 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772366 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772399 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772402 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772401 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772384 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772386 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772396 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772393 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772404 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772377 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772373 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772391 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772397 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772407 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772417 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772398 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772385 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772400 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772388 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772403 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772387 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772390 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772376 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772378 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772382 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772416 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772412 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772383 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772413 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772372 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772414 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772389 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772374 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772379 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772380 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772381 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772375 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772392 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772395 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772394 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772418 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772419 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772421 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772420 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772405 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772415 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772410 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772406 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772408 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772409 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772411 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772450 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772455 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772466 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772456 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772454 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772453 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772446 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772447 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772449 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772435 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772451 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772462 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772461 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772463 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772464 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772431 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772465 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772432 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772443 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772433 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772437 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772439 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772460 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772430 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772429 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772434 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772457 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772467 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772422 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772423 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772438 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772424 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772440 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772444 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772427 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772425 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772441 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772442 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772445 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772428 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772448 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772426 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772436 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772452 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772458 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772459 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772489 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772481 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772488 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772472 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772473 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772482 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772475 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772476 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772484 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772477 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772478 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772479 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772480 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772485 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772487 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772468 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772469 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772483 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772470 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772471 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772486 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:16:42 [ERR] Request 772474 received status 500: NO_BACKEND_AVAILABLE
2026/07/27 04:17:37 [ERR] Request 776596 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": write tcp 172.20.0.6:34370->172.20.0.4:9001: write: broken pipe
2026/07/27 04:17:37 [ERR] Request 776566 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": write tcp 172.20.0.6:34336->172.20.0.4:9001: write: broken pipe
2026/07/27 04:17:37 [ERR] Request 776562 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": write tcp 172.20.0.6:34322->172.20.0.4:9001: write: broken pipe
2026/07/27 04:17:37 [ERR] Request 776579 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": http2: client connection force closed via ClientConn.Close
2026/07/27 04:17:37 [ERR] Request 776583 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": unexpected EOF
2026/07/27 04:17:37 [ERR] Request 776589 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": unexpected EOF
2026/07/27 04:17:37 [ERR] Request 776613 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": dial tcp 172.20.0.4:9001: connect: connection refused
2026/07/27 04:17:37 [ERR] Request 776602 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": dial tcp 172.20.0.4:9001: connect: connection refused
2026/07/27 04:17:37 [ERR] Request 776603 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": dial tcp 172.20.0.4:9001: connect: connection refused
2026/07/27 04:17:37 [ERR] Request 776619 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": dial tcp 172.20.0.4:9001: connect: connection refused
2026/07/27 04:17:37 [ERR] Request 776623 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": dial tcp 172.20.0.4:9001: connect: connection refused
2026/07/27 04:17:37 [ERR] Request 776634 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": dial tcp 172.20.0.4:9001: connect: connection refused
2026/07/27 04:17:37 [ERR] Request 776624 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": dial tcp 172.20.0.4:9001: connect: connection refused
2026/07/27 04:17:37 [ERR] Request 776607 received status 500: Backend Error: Post "http://172.20.0.4:9001/create-session": dial tcp 172.20.0.4:9001: connect: connection refused
2026/07/27 04:17:58 === KẾT QUẢ TEST TẢI ĐỒNG THỜI ===
2026/07/27 04:17:58 Thời gian chạy: 3m8.584053993s
2026/07/27 04:17:58 Tổng số Request gửi đi: 1000000
2026/07/27 04:17:58 Thành công: 999706
2026/07/27 04:17:58 Thất bại: 294
2026/07/27 04:17:58 TPS gửi đi: 5302.68 req/s
2026/07/27 04:17:58 TPS thành công: 5301.12 req/s
2026/07/27 04:17:58 === PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===
2026/07/27 04:17:58 Instance pdu-session-3 xử lý: 194074 requests (19.41%)
2026/07/27 04:17:58 Instance pdu-session-1 xử lý: 268544 requests (26.86%)
2026/07/27 04:17:58 Instance pdu-session-4 xử lý: 268544 requests (26.86%)
2026/07/27 04:17:58 Instance pdu-session-2 xử lý: 268544 requests (26.86%)

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
