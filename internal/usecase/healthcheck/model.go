package healthcheck

type Response struct {
	Message    string `json:"message"`
	ServerTime string `json:"server_time"`
	Version    string `json:"version"`
}
