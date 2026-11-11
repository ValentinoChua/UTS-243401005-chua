package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Keluhan struct {
	ID     int    `json:"id"`
	Ruang  string `json:"ruang"`
	Jenis  string `json:"jenis"`
	Status string `json:"status"`
}

type ResponseKeluhan struct {
	Jumlah int       `json:"jumlah"`
	Data   []Keluhan `json:"data"`
}

func kirimRequest(client *http.Client, request *http.Request) (*http.Response, error) {

	if request.Method == http.MethodPost {
		return client.Do(request)
	}

	var response *http.Response
	var err error

	for percobaan := 1; percobaan <= 5; percobaan++ {
		fmt.Println("Percobaan ke-", percobaan)

		response, err = client.Do(request)

		if err == nil {

			if response.StatusCode != http.StatusTooManyRequests &&
				response.StatusCode < 500 {
				return response, nil
			}

			if percobaan == 5 {
				return response, nil
			}

			fmt.Println("Server memberikan status:", response.Status)
			response.Body.Close()

		} else {
			fmt.Println("Gagal menghubungi server:", err)
		}

		if percobaan < 5 {
			jeda := 200 * time.Millisecond

			if response != nil {
				retryAfter := response.Header.Get("Retry-After")

				if retryAfter != "" {
					detik, err := strconv.Atoi(retryAfter)

					if err == nil && detik >= 0 {
						jeda = time.Duration(detik) * time.Second
					} else {
						for i := 1; i < percobaan; i++ {
							jeda *= 2
						}
					}
				} else {
					for i := 1; i < percobaan; i++ {
						jeda *= 2
					}
				}
			} else {
				for i := 1; i < percobaan; i++ {
					jeda *= 2
				}
			}

			fmt.Println("Menunggu:", jeda)
			time.Sleep(jeda)
		}
	}

	return nil, err
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Penggunaan:")
		fmt.Println("go run ./klien <alamat> lihat <id>")
		fmt.Println("go run ./klien <alamat> buat <ruang> <jenis>")
		os.Exit(1)
	}

	alamat := os.Args[1]
	perintah := os.Args[2]

	var url string
	var method string
	var requestBody []byte
	var err error

	switch perintah {
	case "lihat":
		if len(os.Args) != 4 {
			fmt.Println("Format: go run ./klien <alamat> lihat <id>")
			os.Exit(1)
		}

		id := os.Args[3]
		url = alamat + "/keluhan/" + id
		method = http.MethodGet

	case "buat":
		if len(os.Args) != 5 {
			fmt.Println("Format: go run ./klien <alamat> buat <ruang> <jenis>")
			os.Exit(1)
		}

		ruang := os.Args[3]
		jenis := os.Args[4]

		url = alamat + "/keluhan"
		method = http.MethodPost

		data := map[string]string{
			"ruang": ruang,
			"jenis": jenis,
		}

		requestBody, err = json.Marshal(data)
		if err != nil {
			fmt.Println("Gagal membuat JSON:", err)
			os.Exit(1)
		}

		fmt.Println("Mengirim keluhan:", string(requestBody))

	default:
		fmt.Println("Perintah tidak dikenal. Gunakan lihat atau buat.")
		os.Exit(1)
	}

	fmt.Println("Alamat server:", alamat)
	fmt.Println("Perintah:", perintah)

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	request, err := http.NewRequest(method, url, bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println("Gagal membuat request:", err)
		os.Exit(1)
	}

	if method == http.MethodPost {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := kirimRequest(&client, request)
	if err != nil {
		fmt.Println("Gagal menghubungi server:", err)
		os.Exit(1)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Println("Request gagal:", response.Status)
		os.Exit(1)
	}

	fmt.Println("Respons dari server:")

	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		fmt.Println("Gagal membaca respons:", err)
		os.Exit(1)
	}
}
