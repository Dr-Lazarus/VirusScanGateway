package handlers

type AnalysisResponse struct {
	Data struct {
		ID         string            `json:"id"`
		Type       string            `json:"type"`
		Links      map[string]string `json:"links"`
		Attributes struct {
			Stats struct {
				Malicious        int `json:"malicious"`
				Suspicious       int `json:"suspicious"`
				Undetected       int `json:"undetected"`
				Harmless         int `json:"harmless"`
				Timeout          int `json:"timeout"`
				ConfirmedTimeout int `json:"confirmed-timeout"`
				Failure          int `json:"failure"`
				TypeUnsupported  int `json:"type-unsupported"`
			} `json:"stats"`
			Results map[string]struct {
				Method        string      `json:"method"`
				EngineName    string      `json:"engine_name"`
				EngineVersion string      `json:"engine_version"`
				EngineUpdate  string      `json:"engine_update"`
				Category      string      `json:"category"`
				Result        interface{} `json:"result"`
			} `json:"results"`
			Date   int    `json:"date"`
			Status string `json:"status"`
		} `json:"attributes"`
	} `json:"data"`
	Meta struct {
		FileInfo struct {
			SHA256 string `json:"sha256"`
			MD5    string `json:"md5"`
			SHA1   string `json:"sha1"`
			Size   int    `json:"size"`
		} `json:"file_info"`
	} `json:"meta"`
}
