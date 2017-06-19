package packet

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// InformData captures the unpacked representation of an inform packet from the device.
type InformData struct {
	BootVersion string `json:"bootrom_version,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	State       int    `json:"state,omitempty"`

	Interfaces []Interface `json:"if_table,omitempty"`
	RadioInfo  []RadioInfo `json:"radio_table,omitempty"`

	IsDiscovery bool `json:"discovery_response,omitempty"`
}

// Interface captures the interface information reported in a Inform packet.
type Interface struct {
	Drops      int    `json:"drops"`
	FullDuplex bool   `json:"full_duplex"`
	IP         string `json:"ip"`
	Latency    int    `json:"latency"`
	MAC        string `json:"mac"`
	Name       string `json:"name"`
	Netmask    string `json:"netmask"`
	NumPorts   int    `json:"num_port"`
	Speed      int    `json:"speed"`
	Up         bool   `json:"up"`
	Uptime     int    `json:"uptime"`
}

// RadioInfo describes a wireless interface on the device.
type RadioInfo struct {
	BuiltinAntennaGain int    `json:"builtin_ant_gain"`
	HasInternalAntenna bool   `json:"builtin_antenna"`
	MaxTx              int    `json:"max_txpower"`
	Name               string `json:"name"`
	Radio              string `json:"radio"`
}

// FormatDiscoveryResponse decodes a JSON inform payload representing a discoveryResponse packet
func FormatDiscoveryResponse(d []byte) (*InformData, error) {
	var out InformData
	return &out, json.Unmarshal(d, &out)
}

// CommandData encapsulates the data sent in an instruction to a AP.
type CommandData struct {
	Type            string `json:"_type"`
	ServerTimestamp string `json:"server_time_in_utc,omitempty"`

	Interval int `json:"interval,omitempty"`

	ConfigVersion    string `json:"cfgversion,omitempty"`
	ManagementConfig string `json:"mgmt_cfg,omitempty"`
	SystemConfig     string `json:"system_cfg,omitempty"`
	BlockedStations  string `json:"blocked_sta"`
}

// MakeNoop creates the payload section of a noop response.
func MakeNoop(pollDelay int) ([]byte, error) {
	return json.Marshal(CommandData{Type: "noop", Interval: pollDelay})
}

// MakeMgmtConfigUpdate creates the payload section of a response which sets all configuration.
func MakeMgmtConfigUpdate(mgmtCfg, configVersion string) ([]byte, error) {
	return json.Marshal(CommandData{
		Type:             "setparam",
		ServerTimestamp:  unixMicroPSTString(),
		ManagementConfig: strings.Replace(mgmtCfg, "\n", "\\n", -1),
		ConfigVersion:    configVersion,
	})
}

// MakeConfigUpdate creates the payload section of a response which sets all configuration.
func MakeConfigUpdate(sysCfg, mgmtCfg, configVersion string) ([]byte, error) {
	return json.Marshal(CommandData{
		Type:             "setparam",
		SystemConfig:     sysCfg,
		ServerTimestamp:  unixMicroPSTString(),
		ManagementConfig: mgmtCfg,
		ConfigVersion:    configVersion,
	})
}

//Credit: mcrute - https://github.com/mcrute/go-inform/blob/master/inform/tx_messages.go
func unixMicroPST() int64 {
	tnano := time.Now().UnixNano()
	return tnano / int64(time.Millisecond)
}

func unixMicroPSTString() string {
	return strconv.FormatInt(unixMicroPST(), 10)
}
