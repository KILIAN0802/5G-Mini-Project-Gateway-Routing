# HƯỚNG DẪN DỰ ÁN & TÀI LIỆU THIẾT KẾ KIẾN TRÚC
## GATEWAY ROUTING & PDU SESSION SIMULATOR

Hệ thống mô phỏng quy trình định tuyến gói tin (Gateway Routing) và xử lý phiên kết nối (PDU Session) trong kiến trúc 5G Core đơn giản hóa. Hệ thống được thiết kế tối ưu hóa hiệu năng, chịu tải cao, tích hợp cân bằng tải (Load Balancing) và tự động phát hiện dịch vụ (Service Discovery).

---

## 1. SƠ ĐỒ KIẾN TRÚC & LUỒNG DỮ LIỆU (FLOW)

Dưới đây là mô hình tương tác giữa Client bắn tải (`auto-request`), Gateway định tuyến, và cụm các instance xử lý dịch vụ (`pdu-session`):

```mermaid
sequenceDiagram
    autonumber
    actor Client as Client (auto-request)
    participant GW as Gateway (Load Balancer)
    participant Registry as Registry (Service Discovery)
    participant PDU1 as PDU Session 1 (IP: 172.20.0.2)
    participant PDU2 as PDU Session 2 (IP: 172.20.0.3)

    Note over GW, Registry: Quét DNS Docker định kỳ<br/>net.LookupIP("pdu-session")
    Registry->>GW: Cập nhật danh sách IP khả dụng

    Client->>GW: Gửi yêu cầu khởi tạo session (HTTP POST /nsmf-pdusession/v1/sm-contexts)
    
    Note over GW: Chọn backend tối ưu dựa trên:<br/>RoundRobin / WeightedRR / LeastConnection
    
    GW->>PDU1: Chuyển tiếp gói tin (HTTP POST /create-session)
    Note over PDU1: Xử lý tạo session &<br/>Lưu thông tin vào Redis
    
    PDU1-->>GW: Trả về kết quả JSON (Active Session)
    GW-->>Client: Trả về kết quả cuối cùng cho Client
```

---

## 2. KIẾN TRÚC CÁC THÀNH PHẦN CHI TIẾT

### 2.1 Cổng Kết Nối Định Tuyến (Gateway)
Gateway đóng vai trò là điểm tiếp nhận duy nhất (Single Point of Entry) cho mọi yêu cầu khởi tạo session từ phía Client.
* **Cơ chế Service Discovery (Tự phát hiện dịch vụ)**:
  * Gateway không sử dụng địa chỉ IP cứng. Thay vào đó, một Goroutine nền chạy hàm `registry.ServiceDiscovery()`.
  * Nó gọi hàm `net.LookupIP("pdu-session")` để truy vấn DNS nội bộ của Docker. Docker Engine sẽ tự động trả về danh sách IP của tất cả các container thuộc dịch vụ `pdu-session` đang chạy.
  * Danh sách này liên tục được đồng bộ hóa với Registry nội bộ để thêm mới các instance vừa scale-up hoặc loại bỏ các instance bị lỗi (scale-down/crashed).
* **Cơ chế Health Check (Kiểm tra sức khỏe)**:
  * Một Goroutine định kỳ gửi request `/health` tới từng instance. Nếu instance không phản hồi trong khoảng thời gian cấu hình, cờ `Healthy` sẽ được chuyển thành `false` và instance đó tạm thời bị loại khỏi danh sách định tuyến.
* **Thuật toán Cân bằng tải (Load Balancing)**:
  * **Round Robin (RR)**: Định tuyến luân phiên đều giữa các instance khỏe mạnh.
  * **Weighted Round Robin (WRR)**: Phân phối yêu cầu dựa theo trọng số sức mạnh cấu hình trước của từng instance (ví dụ: máy mạnh gánh nhiều tải hơn).
  * **Least Connections (Load Balancer)**: Định tuyến gói tin mới vào instance có số lượng request đang xử lý (`ActiveRequests`) thấp nhất tại thời điểm đó (được thu thập qua API `/metrics`).


### 2.2 Công Cụ Bắn Tải Tự Động (Auto-Request)
Công cụ kiểm thử hiệu năng (Benchmark Tool) viết bằng Go, tối ưu hóa concurrency qua Goroutines và Channels.
* **Cơ chế Concurrency Control**: Sử dụng mô hình Worker Pool (giới hạn tối đa 500 luồng chạy song song) để đẩy tải cực hạn lên Gateway mà không làm sập bộ nhớ máy client.
---

## 3. HƯỚNG DẪN CẤU HÌNH HỆ THỐNG

### 3.1 Cấu hình tài nguyên & Biến môi trường (Docker Compose)
Tệp cấu hình `docker-compose.yml` khai báo các dịch vụ gồm Gateway (gán CPU core 0, GOGC=500), Redis, Redis Commander, 4 instance PDU Session và Auto-Request Client:

```yaml
services:
  # Container của gateway - gán cứng vào CPU Core 0
  gateway:
    build:
      context: ./gateway
    ports:
      - "8080:8080"
    cpuset: "0"
    environment:
      - GOGC=500
    depends_on:
      - pdu-session-1
      - pdu-session-2
      - pdu-session-3
      - pdu-session-4

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"

  redis-commander:
    image: rediscommander/redis-commander
    environment:
      - REDIS_HOSTS=redis:redis:6379
    ports:
      - "8081:8081"
    depends_on:
      - redis

  # Cụm các instance PDU Session (pdu-session-1 đến pdu-session-4)
  pdu-session-1:
    build:
      context: ./pdu-session
    environment:
      - INSTANCE_ID=pdu-session-1
      - PORT=9001
      - REDIS_ADDR=redis:6379
    networks:
      default:
        aliases:
          - pdu-session
    depends_on:
      - redis

  # (pdu-session-2, pdu-session-3, pdu-session-4 được định nghĩa tương tự)

  auto-request:
    build:
      context: ./auto-request
    environment:
      - TARGET_URL=http://gateway:8080/nsmf-pdusession/v1/sm-contexts
    depends_on:
      - gateway
```

### 3.2 Cấu hình đồng bộ Timeout trong Code
Để đảm bảo hệ thống không gặp lỗi `Gateway Timeout` hoặc `Client Timeout` khi xử lý dưới tải nặng, các mốc Timeout trong mã nguồn Go được đồng bộ lên **90 giây**:

* **Trong Gateway** (`gateway/main.go`):
  Khởi tạo `Timeout: 90 * time.Second` cho HTTP/2 Client pool dùng để gửi request chuyển tiếp (forward) tới các instance PDU Session:
  ```go
  pduClientPool[i] = &http.Client{
      Timeout: 90 * time.Second, // Đảm bảo đủ thời gian chờ PDU Session xử lý và phản hồi
      Transport: &http2.Transport{ ... },
  }
  ```

* **Trong Auto-Request Client** (`auto-request/main.go`):
  Cấu hình `ClientTimeout: 90 * time.Second` cho HTTP Client bắn tải tới Gateway:
  ```go
  var config = Config{
      targetURL:     "http://127.0.0.1:8080/nsmf-pdusession/v1/sm-contexts",
      ClientTimeout: 90 * time.Second, // Tránh Client ngắt kết nối sớm khi hệ thống chịu tải cao
  }
  ```

---

## 4. QUY TRÌNH VẬN HÀNH & KIỂM THỬ

### Bước 1: Khởi động và Cấu hình Tỷ lệ (Scaling) hệ thống
* Khởi động mặc định (4 replicas):
  ```bash
  docker compose up --build
  ```
* **Scale-up / Scale-down số lượng PDU Session**:
  Mở một terminal khác và chạy lệnh sau để thay đổi số lượng bản sao chạy song song (ví dụ: nâng lên 6 hoặc hạ xuống 2):
  ```bash
  docker compose up -d --scale pdu-session=6
  ```
  *(Gateway sẽ tự động phát hiện số lượng instance thay đổi và cập nhật lại bộ cân bằng tải trong vòng 10 giây)*

### Bước 2: Thực hiện kiểm thử chạy tải (Benchmark)
Khởi chạy công cụ bắn tải `auto-request` (bằng Docker Compose hoặc chạy trực tiếp bằng Go):

* **Cách 1: Chạy qua Docker Compose (Khuyên dùng)**
  ```bash
  docker compose run --rm auto-request
  ```

* **Cách 2: Chạy trực tiếp trên máy vật lý**
  ```bash
  cd auto-request
  go run .
  ```

#### Các bước thiết lập tham số kiểm thử:
Chương trình tương tác qua dòng lệnh và sẽ lần lượt yêu cầu nhập các tham số:
1. **Nhập số lượng request**: Tổng số gói tin cần gửi đi (ví dụ: `5000`).
2. **Nhập số lượng connection**: Số lượng kết nối HTTP/2 multiplexing song song tới Gateway (mặc định `1`).
3. **Nhập số lượng Worker song song**: Số lượng luồng Worker xử lý đồng thời (mặc định `100`).

#### Thao tác điều khiển:
* Nhập `exit` hoặc `quit` tại màn hình nhắc lệnh để thoát chương trình.
* Nhấn phím **ESC** bất cứ lúc nào để dừng khẩn cấp vòng lặp hoặc tiến trình bắn tải một cách an toàn.
