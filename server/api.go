package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/afrizal423/go-NIK-parse/ParseNIK"
	"github.com/otiai10/gosseract/v2"
)

type Server struct{}

type Response struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    KTPData
}

type KTPData struct {
	NIK              string `json:"nik"`
	Nama             string `json:"nama"`
	TempatLahir      string `json:"tempat_lahir"`
	TanggalLahir     string `json:"tanggal_lahir"`
	JenisKelamin     string `json:"jenis_kelamin"`
	Alamat           string `json:"alamat"`
	RT_RW            string `json:"rt_rw"`
	KelDesa          string `json:"kel_desa"`
	Kecamatan        string `json:"kecamatan"`
	Agama            string `json:"agama"`
	StatusPerkawinan string `json:"status_perkawinan"`
	Pekerjaan        string `json:"pekerjaan"`
	Kewarganegaraan  string `json:"kewarganegaraan"`
	BerlakuHingga    string `json:"berlaku_hingga"`
}

type OCRResponse struct {
	Text string `json:"text"`
}

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Error:   false,
		Code:    200,
		Message: "Success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) ParseNiks(w http.ResponseWriter, r *http.Request) {

	tesnik := "3503081905980002"
	coba := ParseNIK.GetdataNIK(tesnik)
	data, _ := json.Marshal(coba)
	fmt.Println(string(data))

	response := Response{
		Error:   false,
		Code:    200,
		Message: "Success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

	err := r.ParseForm()
	if err != nil {
		fmt.Printf("[main][func: handleBinaryFileUpload] Failed when Parsing Form: %s", errors.Unwrap(err).Error())
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusInternalServerError)
		fmt.Println("[main][func: handleBinaryFileUpload] Error Retrieving the File:", err)
		return
	}
	defer file.Close()

	fileExt := filepath.Ext(handler.Filename)
	if !(fileExt == ".jpg" || fileExt == ".jpeg" || fileExt == ".png") {
		http.Error(w, "Invalid File Type", http.StatusBadRequest)
		fmt.Printf("[main][func: handleBinaryFileUpload] Invalid File Type. Uploaded: %s\n", fileExt)
		return
	}

	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, file); err != nil {
		http.Error(w, "Failed when processing file", http.StatusInternalServerError)
		fmt.Println("[main][func: handleBinaryFileUpload] Failed when processing file:", err)
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	err = client.SetImageFromBytes(fileBuffer.Bytes())
	if err != nil {
		http.Error(w, "Failed to set image from buffer", http.StatusInternalServerError)
		log.Fatalf("Failed to set image from buffer: %v", err)
		return
	}

	text, err := client.Text()
	if err != nil {
		http.Error(w, "Failed to read text from image", http.StatusInternalServerError)
		log.Fatalf("Failed to read text from image: %v", err)
		return
	}
	log.Println("Extracted text from image: ", text)

	ktpData := extractKTPData(text)

	jsonResponse, err := json.Marshal(ktpData)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		log.Fatalf("Failed to marshal JSON: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}
