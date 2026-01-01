package api

import (
	"bytes"
	"encoding/xml"
	"log"
	"net/http"
)

type irccRequest struct {
	XMLName  xml.Name `xml:"s:Envelope"`
	SOAPNS   string   `xml:"xmlns:s,attr"`
	Encoding string   `xml:"s:encodingStyle,attr"`
	Body     struct {
		XSendIRCC struct {
			XMLName  xml.Name `xml:"u:X_SendIRCC"`
			Schemas  string   `xml:"xmlns:u,attr"`
			IRCCCode string   `xml:"IRCCCode"`
		}
	} `xml:"s:Body"`
}

func (server *Server) sendKeyRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	query_key := r.PathValue("id")
	key, ok := server.IrrCodes[query_key]
	if !ok {
		http.Error(w, "Unknown key", http.StatusBadRequest)
		return
	}

	irccreq := irccRequest{
		SOAPNS:   "http://schemas.xmlsoap.org/soap/envelope/",
		Encoding: "http://schemas.xmlsoap.org/soap/encoding/",
	}
	irccreq.Body.XSendIRCC.Schemas = "urn:schemas-sony-com:service:IRCC:1"
	irccreq.Body.XSendIRCC.IRCCCode = key

	bodyBytes, err := xml.MarshalIndent(irccreq, "", "  ")
	if err != nil {
		log.Printf("XML error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	targetURL := "http://" + server.TVIP + "/sony/IRCC"
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Set("SOAPACTION", `"urn:schemas-sony-com:service:IRCC:1#X_SendIRCC"`)
	req.Header.Set("X-Auth-PSK", server.TVPSK)

	resp, err := server.Client.Do(req)
	if err != nil {
		log.Printf("Query to TV error: %v", err)
		http.Error(w, "Can't reach the TV", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("TV answer %d", resp.StatusCode)
		http.Error(w, "TV error", http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
