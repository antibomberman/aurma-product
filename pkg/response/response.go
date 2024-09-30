package response

import (
	"bytes"
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(w http.ResponseWriter, message string, data interface{}) {
	switch data.(type) {
	case proto.Message:
		successProto(w, message, data.(proto.Message))
		return
	default:
		response := Response{
			Success: true,
			Message: message,
			Data:    data,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

}
func Fail(w http.ResponseWriter, message string) {
	response := Response{
		Success: false,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

func successProto(w http.ResponseWriter, message string, data proto.Message) {
	m := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}

	jsonData, err := m.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal proto to JSON", http.StatusInternalServerError)
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString(`{"success": true, "message": "` + message + `", "data": `)
	buffer.Write(jsonData)
	buffer.WriteString(`}`)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buffer.Bytes())
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
