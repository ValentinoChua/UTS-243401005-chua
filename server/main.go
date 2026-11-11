package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Keluhan struct {
	ID     int    `json:"id"`
	Ruang  string `json:"ruang"`
	Jenis  string `json:"jenis"`
	Status string `json:"status"`
}

var daftarKeluhan = []Keluhan{}
var nextID = 1
var mu sync.Mutex

func kirimError(w http.ResponseWriter, pesan string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]string{
		"galat": pesan,
	}

	json.NewEncoder(w).Encode(response)
}

func createKeluhanHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		kirimError(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	var keluhan Keluhan

	err := json.NewDecoder(r.Body).Decode(&keluhan)
	if err != nil {
		kirimError(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	keluhan.Ruang = strings.TrimSpace(keluhan.Ruang)
	keluhan.Jenis = strings.TrimSpace(keluhan.Jenis)

	if keluhan.Ruang == "" || keluhan.Jenis == "" {
		kirimError(w, "Ruang dan jenis keluhan wajib diisi", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, data := range daftarKeluhan {
		if strings.EqualFold(data.Ruang, keluhan.Ruang) &&
			strings.EqualFold(data.Jenis, keluhan.Jenis) &&
			data.Status == "terbuka" {

			kirimError(w, "Keluhan serupa masih terbuka", http.StatusConflict)
			return
		}
	}

	keluhan.ID = nextID
	nextID++
	keluhan.Status = "terbuka"

	daftarKeluhan = append(daftarKeluhan, keluhan)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(keluhan)
}

func getKeluhanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		kirimError(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")

	if status != "" && status != "terbuka" && status != "selesai" {
		kirimError(w, "Status tidak dikenal", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	hasil := []Keluhan{}

	for _, data := range daftarKeluhan {
		if status == "" || data.Status == status {
			hasil = append(hasil, data)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Jumlah int       `json:"jumlah"`
		Data   []Keluhan `json:"data"`
	}{
		Jumlah: len(hasil),
		Data:   hasil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getKeluhanByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		kirimError(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/keluhan/")
	id, err := strconv.Atoi(idText)

	if err != nil || id <= 0 {
		kirimError(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, data := range daftarKeluhan {
		if data.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
			return
		}
	}

	kirimError(w, "Keluhan tidak ditemukan", http.StatusNotFound)
}

func selesaikanKeluhanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		kirimError(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/keluhan/")
	idText = strings.TrimSuffix(idText, "/selesaikan")

	id, err := strconv.Atoi(idText)

	if err != nil || id <= 0 {
		kirimError(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, data := range daftarKeluhan {
		if data.ID == id {
			if data.Status == "selesai" {
				kirimError(w, "data sudah selesai", http.StatusConflict)
				return
			}

			daftarKeluhan[i].Status = "selesai"

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(daftarKeluhan[i])
			return
		}
	}

	kirimError(w, "Keluhan tidak ditemukan", http.StatusNotFound)
}

func deleteKeluhanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		kirimError(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/keluhan/")
	id, err := strconv.Atoi(idText)

	if err != nil || id <= 0 {
		kirimError(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, data := range daftarKeluhan {
		if data.ID == id {
			if data.Status == "selesai" {
				kirimError(w, "Keluhan yang sudah selesai tidak dapat dihapus", http.StatusConflict)
				return
			}

			daftarKeluhan = append(daftarKeluhan[:i], daftarKeluhan[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]string{
				"status": "dihapus",
			}

			json.NewEncoder(w).Encode(response)
			return
		}
	}

	kirimError(w, "Keluhan tidak ditemukan", http.StatusNotFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server Keluhan Fasilitas Kelas Berjalan!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		kirimError(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/healthz", healthHandler)

	http.HandleFunc("/keluhan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createKeluhanHandler(w, r)
		} else if r.Method == http.MethodGet {
			getKeluhanHandler(w, r)
		} else {
			kirimError(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/keluhan/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/selesaikan") {
			selesaikanKeluhanHandler(w, r)
		} else if r.Method == http.MethodDelete {
			deleteKeluhanHandler(w, r)
		} else {
			getKeluhanByIDHandler(w, r)
		}
	})

	fmt.Println("Server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
