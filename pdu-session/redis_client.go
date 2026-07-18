package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type redisConn struct {
	net.Conn
	rd *bufio.Reader
}

var (
	redisAddr string
	pool      chan *redisConn
)

// initRedis lấy địa chỉ từ biến môi trường và khởi tạo Connection Pool
func initRedis() {
	redisAddr = GetEnv("REDIS_ADDR", "redis:6379")
	limitStr := GetEnv("REDIS_POOL_SIZE", "150")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 150
	}
	pool = make(chan *redisConn, limit)
}

// getRedisConn lấy kết nối từ pool hoặc tạo mới nếu pool trống
func getRedisConn() (*redisConn, error) {
	select {
	case conn := <-pool:
		return conn, nil
	default:
		c, err := net.DialTimeout("tcp", redisAddr, 5*time.Second)
		if err != nil {
			return nil, err
		}
		if tcpConn, ok := c.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(30 * time.Second)
		}
		return &redisConn{
			Conn: c,
			rd:   bufio.NewReader(c),
		}, nil
	}
}

// releaseRedisConn trả lại kết nối vào pool hoặc đóng nếu có lỗi
func releaseRedisConn(conn *redisConn, err error) {
	if err != nil {
		conn.Close()
		return
	}
	select {
	case pool <- conn:
		// Trả về pool thành công
	default:
		// Pool đầy, đóng kết nối
		conn.Close()
	}
}

// SaveSessionInRedis lưu thông tin session dạng JSON vào Redis
func SaveSessionInRedis(supi string, data string) error {
	conn, err := getRedisConn()
	if err != nil {
		return err
	}
	var opErr error
	defer func() {
		releaseRedisConn(conn, opErr)
	}()

	key := "session:" + supi
	// Gửi lệnh SET key value theo định dạng RESP
	var buf bytes.Buffer
	buf.WriteString("*3\r\n")
	buf.WriteString("$3\r\nSET\r\n")
	buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(key), key))
	buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(data), data))
	_, opErr = conn.Write(buf.Bytes())
	if opErr != nil {
		return opErr
	}
	// Đọc phản hồi từ connection reader
	line, opErr := conn.rd.ReadString('\n')
	if opErr != nil {
		return opErr
	}
	if line[0] == '-' {
		opErr = fmt.Errorf("redis error: %s", strings.TrimSpace(line[1:]))
		return opErr
	}
	return nil
}

// GetAllSessionsFromRedis lấy danh sách toàn bộ session
func GetAllSessionsFromRedis() (map[string]string, error) {
	conn, err := getRedisConn()
	if err != nil {
		return nil, err
	}
	var opErr error
	defer func() {
		releaseRedisConn(conn, opErr)
	}()

	// Gửi lệnh KEYS session:*
	var buf bytes.Buffer
	buf.WriteString("*2\r\n")
	buf.WriteString("$4\r\nKEYS\r\n")
	pattern := "session:*"
	buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(pattern), pattern))
	_, opErr = conn.Write(buf.Bytes())
	if opErr != nil {
		return nil, opErr
	}
	keys, opErr := readRESPArray(conn.rd)
	if opErr != nil {
		return nil, opErr
	}
	sessions := make(map[string]string)
	if len(keys) == 0 {
		return sessions, nil
	}

	// Gửi lệnh MGET để lấy toàn bộ value trong 1 roundtrip
	buf.Reset()
	buf.WriteString(fmt.Sprintf("*%d\r\n", 1+len(keys)))
	buf.WriteString("$4\r\nMGET\r\n")
	for _, key := range keys {
		buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(key), key))
	}
	_, opErr = conn.Write(buf.Bytes())
	if opErr != nil {
		return nil, opErr
	}
	values, opErr := readRESPArray(conn.rd)
	if opErr != nil {
		return nil, opErr
	}

	for i, key := range keys {
		if i < len(values) && values[i] != "" {
			supi := strings.TrimPrefix(key, "session:")
			sessions[supi] = values[i]
		}
	}
	return sessions, nil
}
func readRESPBulkString(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if line[0] == '-' {
		return "", fmt.Errorf("redis error: %s", strings.TrimSpace(line[1:]))
	}
	if line[0] != '$' {
		return "", fmt.Errorf("invalid resp type: %c", line[0])
	}
	lengthStr := strings.TrimSpace(line[1:])
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}
	if length == -1 {
		return "", nil
	}
	buf := make([]byte, length+2)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return "", err
	}
	return string(buf[:length]), nil
}
func readRESPArray(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if line[0] == '-' {
		return nil, fmt.Errorf("redis error: %s", strings.TrimSpace(line[1:]))
	}
	if line[0] != '*' {
		return nil, fmt.Errorf("invalid resp type: %c", line[0])
	}
	countStr := strings.TrimSpace(line[1:])
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, err
	}
	if count == -1 {
		return nil, nil
	}
	results := make([]string, count)
	for i := 0; i < count; i++ {
		val, err := readRESPBulkString(reader)
		if err != nil {
			return nil, err
		}
		results[i] = val
	}
	return results, nil
}
