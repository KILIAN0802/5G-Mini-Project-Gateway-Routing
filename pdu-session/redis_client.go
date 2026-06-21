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

var redisAddr string

// initRedis lấy địa chỉ từ biến môi trường

func initRedis() {
	redisAddr = GetEnv("REDIS_ADDR", "redis:6379")
}

// getRedisConn tạo kết nối TCP tới Redis

func getRedisConn() (net.Conn, error){
	return net.DialTimeout("tcp", redisAddr, 5*time.Second)
}

// SaveSessionInRedis lưu thông tin session dạng JSON vào Redis
func SaveSessionInRedis(supi string, data string) error {
	conn, err := getRedisConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	key := "session:" + supi
	// Gửi lệnh SET key value theo định dạng RESP
	var buf bytes.Buffer
	buf.WriteString("*3\r\n")
	buf.WriteString("$3\r\nSET\r\n")
	buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(key), key))
	buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(data), data))
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	// Đọc phản hồi
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if line[0] == '-' {
		return fmt.Errorf("redis error: %s", strings.TrimSpace(line[1:]))
	}
	return nil
}
// GetAllSessionsFromRedis lấy danh sách toàn bộ session
func GetAllSessionsFromRedis() (map[string]string, error) {
	conn, err := getRedisConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	// Gửi lệnh KEYS session:*
	var buf bytes.Buffer
	buf.WriteString("*2\r\n")
	buf.WriteString("$4\r\nKEYS\r\n")
	pattern := "session:*"
	buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(pattern), pattern))
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)
	keys, err := readRESPArray(reader)
	if err != nil {
		return nil, err
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
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}
	values, err := readRESPArray(reader)
	if err != nil {
		return nil, err
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
