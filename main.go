package main

import (
	"fmt"
	"net"
	"time"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckPrinter ตรวจสอบว่าเครื่องพิมพ์ออนไลน์หรือไม่
func CheckPrinter(ip string, port string) bool {
	address := fmt.Sprintf("%s:%s", ip, port)
	timeout := 2 * time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// PrintTestPage ส่งคำสั่งพิมพ์ทดสอบไปยังเครื่องพิมพ์
func PrintTestPage(ip string, port string) error {
	address := fmt.Sprintf("%s:%s", ip, port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// คำสั่ง ESC/POS หรือ PCL สำหรับทดสอบพิมพ์ (ขึ้นอยู่กับเครื่องพิมพ์)
	testPrintCmd := []byte("\x1B\x40Hello Printer!\n\n\x1D\x56\x42\x00") // Reset & Print "Hello Printer!"

	_, err = conn.Write(testPrintCmd)
	if err != nil {
		return err
	}

	return nil
}

// สร้าง handler สำหรับ Gin
func handler(w http.ResponseWriter, r *http.Request) {
	r := gin.Default()

	printerIP := "192.168.1.155"
	printerPort := "9100"

	r.GET("/check-printer", func(c *gin.Context) {
		if CheckPrinter(printerIP, printerPort) {
			c.JSON(200, gin.H{"status": "online", "message": "Printer is reachable"})
		} else {
			c.JSON(200, gin.H{"status": "offline", "message": "Cannot connect to printer"})
		}
	})

	r.GET("/print-test", func(c *gin.Context) {
		err := PrintTestPage(printerIP, printerPort)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "success", "message": "Test page sent to printer"})
	})

	// รัน Gin Server บน Port 8080
	r.Run(":8080")
}

func main() {
	http.HandleFunc("/", handler) // ระบุ handler
	http.ListenAndServe(":8080", nil)
}
