## 📊 Bảng So Sánh Các Đợt Test

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



