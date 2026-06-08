package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Api struct {
		Host     string `yaml:"host"`
		Key      string `yaml:"key"`
		NodeID   uint   `yaml:"node_id"`
		Interval int    `yaml:"interval"`
	} `yaml:"api"`
	Xray struct {
		Path string `yaml:"path"`
		Dir  string `yaml:"dir"`
	} `yaml:"xray"`
	Dns struct {
		Servers []string `yaml:"servers"`
	} `yaml:"dns"`
}

type ApiNodeConfigResponse struct {
	Data struct {
		ID          uint                   `json:"id"`
		Port        uint16                 `json:"port"`
		Address     string                 `json:"address"`
		Type        string                 `json:"type"`
		TrafficRate float32                `json:"traffic_rate"`
		Settings    map[string]interface{} `json:"settings"`
	} `json:"data"`
}

type ApiNodeUser struct {
	ID          uint   `json:"id"`
	UUID        string `json:"uuid"`
	Password    string `json:"password"`
	SpeedLimit  uint32 `json:"speed_limit"`
	DeviceLimit uint32 `json:"device_limit"`
}

type ApiNodeUsersResponse struct {
	Data []ApiNodeUser `json:"data"`
}

type TrafficLogPayload struct {
	UserID uint   `json:"user_id"`
	Up     uint64 `json:"up"`
	Down   uint64 `json:"down"`
}

type XrayStatsResponse struct {
	Stat []struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	} `json:"stat"`
}

func main() {
	// 1. Define command line parameters
	configPathFlag := flag.String("config", "config.yaml", "Path to config.yaml file")
	nodeIDFlag := flag.Uint("id", 0, "Node ID to manage (overrides yaml if set)")
	apiHostFlag := flag.String("api-host", "", "Panel API host URL (overrides yaml if set)")
	apiKeyFlag := flag.String("api-key", "", "Panel UniProxy Token (overrides yaml if set)")
	xrayPathFlag := flag.String("xray-path", "", "Path to the xray executable (overrides yaml if set)")
	xrayDirFlag := flag.String("xray-dir", "", "Directory containing Xray assets (overrides yaml if set)")
	intervalFlag := flag.Int("interval", 0, "Interval in seconds to check updates (overrides yaml if set)")
	flag.Parse()

	// Default config values
	var cfg Config
	cfg.Api.Host = "http://localhost:8080"
	cfg.Api.Key = "secret-uniproxy-token"
	cfg.Api.NodeID = 1
	cfg.Api.Interval = 15

	// Load config.yaml if it exists
	if _, err := os.Stat(*configPathFlag); err == nil {
		log.Printf("[Daemon] Loading configuration from %s", *configPathFlag)
		data, err := os.ReadFile(*configPathFlag)
		if err != nil {
			log.Printf("[Daemon] Warning: Failed to read %s: %v", *configPathFlag, err)
		} else {
			err = yaml.Unmarshal(data, &cfg)
			if err != nil {
				log.Printf("[Daemon] Warning: Failed to parse %s: %v", *configPathFlag, err)
			}
		}
	} else {
		log.Printf("[Daemon] Configuration file %s not found, relying on command-line flags or defaults", *configPathFlag)
	}

	// CLI flags override YAML config values
	if *nodeIDFlag != 0 {
		cfg.Api.NodeID = *nodeIDFlag
	}
	if *apiHostFlag != "" {
		cfg.Api.Host = *apiHostFlag
	}
	if *apiKeyFlag != "" {
		cfg.Api.Key = *apiKeyFlag
	}
	if *xrayPathFlag != "" {
		cfg.Xray.Path = *xrayPathFlag
	}
	if *xrayDirFlag != "" {
		cfg.Xray.Dir = *xrayDirFlag
	}
	if *intervalFlag != 0 {
		cfg.Api.Interval = *intervalFlag
	}

	log.Printf("[Daemon] Starting standalone Xray node supervisor")
	log.Printf("[Daemon] Node ID:  %d", cfg.Api.NodeID)
	log.Printf("[Daemon] API Host: %s", cfg.Api.Host)

	// 2. Resolve Xray paths
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[Daemon] Failed to get working directory: %v", err)
	}

	xrayPath := cfg.Xray.Path
	xrayDir := cfg.Xray.Dir

	// Auto-detect project paths if not specified
	if xrayPath == "" || xrayDir == "" {
		baseDir := wd
		if filepath.Base(baseDir) == "daemon" || filepath.Base(baseDir) == "backend" {
			baseDir = filepath.Dir(baseDir)
		}

		detectedXrayDir := filepath.Join(baseDir, "tools", "xray")
		detectedXrayPath := filepath.Join(detectedXrayDir, "xray")

		// Check if file exists in detected path
		if _, err := os.Stat(detectedXrayPath); err == nil {
			if xrayPath == "" {
				xrayPath = detectedXrayPath
			}
			if xrayDir == "" {
				xrayDir = detectedXrayDir
			}
		} else {
			// Fallback to local directory if auto-detect fails
			if xrayPath == "" {
				xrayPath = "./xray"
			}
			if xrayDir == "" {
				xrayDir = "./"
			}
		}
	}

	configPath := filepath.Join(xrayDir, "config_daemon.json")
	log.Printf("[Daemon] Xray path: %s", xrayPath)
	log.Printf("[Daemon] Xray dir:  %s", xrayDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("[Daemon] Shutdown signal received. Stopping Xray process...")
		cancel()
	}()

	// Start the supervisor loop
	runSupervisor(ctx, cfg.Api.Host, cfg.Api.Key, cfg.Api.NodeID, xrayPath, xrayDir, configPath, cfg.Api.Interval, cfg.Dns.Servers)
}

func fetchConfigFromApi(apiHost, apiKey string, nodeID uint, dnsServers []string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// 1. Fetch Node Config
	configUrl := fmt.Sprintf("%s/api/v1/server/UniProxy/config?node_id=%d", apiHost, nodeID)
	reqConfig, err := http.NewRequest("GET", configUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create config request: %w", err)
	}
	reqConfig.Header.Set("Token", apiKey)

	respConfig, err := client.Do(reqConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch node config: %w", err)
	}
	defer respConfig.Body.Close()

	if respConfig.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respConfig.Body)
		return nil, fmt.Errorf("node config request failed with status %d: %s", respConfig.StatusCode, string(body))
	}

	var configRes ApiNodeConfigResponse
	if err := json.NewDecoder(respConfig.Body).Decode(&configRes); err != nil {
		return nil, fmt.Errorf("failed to decode node config: %w", err)
	}

	// 2. Fetch Users
	usersUrl := fmt.Sprintf("%s/api/v1/server/UniProxy/user?node_id=%d", apiHost, nodeID)
	reqUsers, err := http.NewRequest("GET", usersUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create users request: %w", err)
	}
	reqUsers.Header.Set("Token", apiKey)

	respUsers, err := client.Do(reqUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer respUsers.Body.Close()

	if respUsers.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respUsers.Body)
		return nil, fmt.Errorf("users request failed with status %d: %s", respUsers.StatusCode, string(body))
	}

	var usersRes ApiNodeUsersResponse
	if err := json.NewDecoder(respUsers.Body).Decode(&usersRes); err != nil {
		return nil, fmt.Errorf("failed to decode users response: %w", err)
	}

	// 3. Convert API structures into Xray JSON config format
	clients := []map[string]interface{}{}
	for _, u := range usersRes.Data {
		clientObj := map[string]interface{}{
			"email": fmt.Sprintf("user_%d@hy-board.com", u.ID),
		}

		if configRes.Data.Type == "Trojan" || configRes.Data.Type == "Shadowsocks" {
			clientObj["password"] = u.Password
		} else {
			clientObj["id"] = u.UUID
		}
		clients = append(clients, clientObj)
	}

	// Build StreamSettings
	streamSettings := map[string]interface{}{
		"network": "tcp",
	}

	// Map settings from API directly
	for k, v := range configRes.Data.Settings {
		if k == "streamSettings" {
			if streamObj, ok := v.(map[string]interface{}); ok {
				for sk, sv := range streamObj {
					streamSettings[sk] = sv
				}
			}
		} else {
			streamSettings[k] = v
		}
	}

	protocol := "vless"
	var settings map[string]interface{}

	switch configRes.Data.Type {
	case "V2ray":
		protocol = "vmess"
		settings = map[string]interface{}{
			"clients": clients,
		}
	case "Trojan":
		protocol = "trojan"
		settings = map[string]interface{}{
			"clients": clients,
			"fallback": map[string]interface{}{
				"dest": 80,
			},
		}
	case "Shadowsocks":
		protocol = "shadowsocks"
		method := "aes-256-gcm"
		if m, ok := configRes.Data.Settings["method"].(string); ok && m != "" {
			method = m
		}

		// Map clients into shadowsocks users
		ssUsers := []map[string]interface{}{}
		for _, c := range clients {
			ssUser := map[string]interface{}{
				"email":    c["email"],
				"password": c["password"],
			}
			ssUsers = append(ssUsers, ssUser)
		}

		settings = map[string]interface{}{
			"method":   method,
			"password": "default-fallback-ss-multiuser-key",
			"users":    ssUsers,
			"network":  "tcp,udp",
		}
	default: // Vless
		protocol = "vless"
		settings = map[string]interface{}{
			"clients":    clients,
			"decryption": "none",
		}
	}

	inbound := map[string]interface{}{
		"port":           configRes.Data.Port,
		"protocol":       protocol,
		"settings":       settings,
		"streamSettings": streamSettings,
		"tag":            "proxy-in",
	}

	xrayConfig := map[string]interface{}{
		"log": map[string]interface{}{
			"access":   "",
			"error":    "",
			"loglevel": "info",
		},
		"stats": map[string]interface{}{},
		"api": map[string]interface{}{
			"tag":    "api",
			"listen": "127.0.0.1:10085",
			"services": []string{
				"StatsService",
			},
		},
		"policy": map[string]interface{}{
			"levels": map[string]interface{}{
				"0": map[string]interface{}{
					"statsUserUplink":   true,
					"statsUserDownlink": true,
				},
			},
			"system": map[string]interface{}{
				"statsInboundUplink":   true,
				"statsInboundDownlink": true,
			},
		},
		"inbounds": []interface{}{inbound},
		"outbounds": []interface{}{
			map[string]interface{}{
				"protocol": "freedom",
				"settings": map[string]interface{}{},
				"tag":      "direct",
			},
			map[string]interface{}{
				"protocol": "blackhole",
				"settings": map[string]interface{}{},
				"tag":      "blocked",
			},
		},
	}

	if len(dnsServers) > 0 {
		xrayConfig["dns"] = map[string]interface{}{
			"servers": dnsServers,
		}
	}

	return json.MarshalIndent(xrayConfig, "", "  ")
}

func runSupervisor(ctx context.Context, apiHost, apiKey string, nodeID uint, xrayPath, xrayDir, configPath string, interval int, dnsServers []string) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	var currentConfigHash string
	var cmd *exec.Cmd
	var xrayStarted bool

	// Set required assets env
	os.Setenv("XRAY_LOCATION_ASSET", xrayDir)

	// Function to start Xray process
	startProcess := func() {
		if cmd != nil && cmd.Process != nil {
			log.Println("[Daemon] Killing old Xray process before restarting...")
			_ = cmd.Process.Signal(syscall.SIGKILL)
			_, _ = cmd.Process.Wait()
		}

		log.Printf("[Daemon] Launching Xray binary: %s", xrayPath)
		cmd = exec.Command(xrayPath, "-c", configPath)
		cmd.Dir = xrayDir

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("[Daemon] Failed to get stdout pipe: %v", err)
			return
		}
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			log.Printf("[Daemon] Failed to get stderr pipe: %v", err)
			return
		}

		if err := cmd.Start(); err != nil {
			log.Printf("[Daemon] Failed to start Xray: %v", err)
			xrayStarted = false
			return
		}

		xrayStarted = true
		log.Printf("[Daemon] Xray process started successfully with PID: %d", cmd.Process.Pid)

		// Read stdout in background and print logs
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line) // Print to console
			}
		}()

		// Read stderr in background
		go func() {
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				fmt.Fprintln(os.Stderr, scanner.Text())
			}
		}()

		// Start a goroutine to wait for process exit
		go func(pCmd *exec.Cmd) {
			err := pCmd.Wait()
			log.Printf("[Daemon] Xray process exited: %v", err)
			xrayStarted = false
		}(cmd)
	}

	// Initial check and start
	configBytes, err := fetchConfigFromApi(apiHost, apiKey, nodeID, dnsServers)
	if err == nil {
		hash := fmt.Sprintf("%x", configBytes)
		currentConfigHash = hash

		// Write config
		if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
			log.Printf("[Daemon] Failed to write initial config: %v", err)
		} else {
			startProcess()
		}
	} else {
		log.Printf("[Daemon] Failed to generate initial config: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			// Clean up process on daemon stop
			if cmd != nil && cmd.Process != nil {
				log.Println("[Daemon] Cleaning up Xray process...")
				_ = cmd.Process.Signal(syscall.SIGTERM)
				time.Sleep(1 * time.Second)
				_ = cmd.Process.Signal(syscall.SIGKILL)
			}
			return

		case <-ticker.C:
			// 1. Check if process is still running
			if !xrayStarted {
				log.Println("[Daemon] Xray is not running. Attempting restart...")
				startProcess()
				continue
			}

			// 2. Fetch config and check for changes
			newConfigBytes, err := fetchConfigFromApi(apiHost, apiKey, nodeID, dnsServers)
			if err != nil {
				log.Printf("[Daemon] Error polling configuration from API: %v", err)
				continue
			}

			newHash := fmt.Sprintf("%x", newConfigBytes)
			if newHash != currentConfigHash {
				log.Println("[Daemon] Configuration change detected! Updating config and restarting Xray...")
				currentConfigHash = newHash
				if err := os.WriteFile(configPath, newConfigBytes, 0644); err != nil {
					log.Printf("[Daemon] Failed to write updated config: %v", err)
				} else {
					startProcess()
				}
			}

			// 3. Query precise traffic from Xray API via gRPC StatsService
			payloadList, err := queryTrafficFromXray(xrayPath)
			if err != nil {
				log.Printf("[Daemon] Error querying traffic from Xray API: %v", err)
			} else if len(payloadList) > 0 {
				go reportTrafficToApi(apiHost, apiKey, nodeID, payloadList)
			}
		}
	}
}

func queryTrafficFromXray(xrayPath string) ([]TrafficLogPayload, error) {
	cmd := exec.Command(xrayPath, "api", "statsquery", "-s", "127.0.0.1:10085", "-pattern", "user>>>", "-reset")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var res XrayStatsResponse
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		return nil, err
	}

	userStats := make(map[uint]*TrafficLogPayload)
	for _, item := range res.Stat {
		// Expecting format: user>>>user_2@hy-board.com>>>traffic>>>uplink
		parts := strings.Split(item.Name, ">>>")
		if len(parts) < 4 || parts[0] != "user" || parts[2] != "traffic" {
			continue
		}

		email := parts[1]
		direction := parts[3] // uplink or downlink

		if !strings.HasPrefix(email, "user_") {
			continue
		}

		// Parse user ID
		userIDStr := strings.TrimPrefix(email, "user_")
		atIndex := strings.Index(userIDStr, "@")
		if atIndex <= 0 {
			continue
		}
		userIDStr = userIDStr[:atIndex]

		var uID uint
		if _, err := fmt.Sscanf(userIDStr, "%d", &uID); err != nil {
			continue
		}

		bytesVal := uint64(item.Value)

		if bytesVal == 0 {
			continue
		}

		if _, ok := userStats[uID]; !ok {
			userStats[uID] = &TrafficLogPayload{UserID: uID}
		}

		if direction == "uplink" {
			userStats[uID].Up = bytesVal
		} else if direction == "downlink" {
			userStats[uID].Down = bytesVal
		}
	}

	payloadList := make([]TrafficLogPayload, 0, len(userStats))
	for _, p := range userStats {
		payloadList = append(payloadList, *p)
	}

	return payloadList, nil
}

func reportTrafficToApi(apiHost, apiKey string, nodeID uint, traffic []TrafficLogPayload) {
	client := &http.Client{Timeout: 10 * time.Second}
	pushUrl := fmt.Sprintf("%s/api/v1/server/UniProxy/push", apiHost)

	payload := map[string]interface{}{
		"node_id": nodeID,
		"traffic": traffic,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Daemon] Failed to marshal traffic payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", pushUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("[Daemon] Failed to create traffic report request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Daemon] Failed to report traffic to API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[Daemon] Traffic report request failed with status %d: %s", resp.StatusCode, string(body))
	} else {
		log.Printf("[Daemon] Reported traffic for %d users successfully", len(traffic))
	}
}
