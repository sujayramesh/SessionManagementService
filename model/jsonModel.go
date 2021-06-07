package model

// Json request header struct to fill keepalive or TTL value
type SessionRequest struct {
	Expiry *int `json:"expiry"`
}

// Json response struct to fill the response code of the operation performed
type SessionResponse struct {
	ResponseCode int `json:"responseCode"`
}

// Json  used for HTTP GET Session management response and HTTP PUT Session management request
type SessionDetails struct {
	UUID     string `json:"uuid"`
	Duration int    `json:"duration"`
}

// Json request for HTTP DELETE Session management
type SessionUUID struct {
	UUID string `json:"uuid"`
}
