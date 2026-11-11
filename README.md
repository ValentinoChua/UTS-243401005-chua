# UTS-243401005-chua

# 1. Tema Dan Masalah
Tema yang saya dapatkan adalah terdapat keluhan pada fasilitas proyektor kelas yang mati 5 menit sebelum presentasi, lalu lima mahasiswa mengirimkan 5 pesan di grup whatsapp yang membuat petugas menerima 5 pesan yang sama, dan tidak ada yang tahu apakah kerusakan sudah ditangani. Masalah yang diselesaikan disini adalah keluhan yang dikirimkan tidak akan duplikat sehingga jika ada sebuah keluhan yang sama otomatis tidak akan terkirim lagi dan mahasiswa dapat melihat status dari keluhan.

# 2. Cara Menjalankan server, klien, dan bukti-tcp
cara menjalan kan server : go run .\server\
cara menjalan kan klien : go run ./klien (alamat) lihat (id) Atau go run ./klien (alamat) buat (ruang) (jenis)
cara menjalan kan bukti-tcp : go run ./bukti-tcp (alamat)

# 3. Endpoint
| Metode | Path | Kode Status |
| --- | --- | --- |
| GET | /keluhan | 200 |
| GET | /healthz | 200 |
| GET | /keluhan?status=terbuka | 200 / 400 |
| GET | /keluhan/{id} | 200 / 400 / 404 |
| POST | /keluhan | 201 / 400 / 409 |
| POST | /keluhan{id}/selesaikan | 200 / 400 / 404 / 409 |
| DELETE | /keluhan/{id} | 200 / 400 / 404 / 409 |

Contoh curl untuk tiap endpoint
A. Menampilkan seluruh keluhan
curl -i http://localhost:8080/keluhan

B. Memeriksa kondisi server
curl -i http://localhost:8080/healthz

C. Menampilkan keluhan berdasarkan status
curl -i "http://localhost:8080/keluhan?status=terbuka Atau curl -i "http://localhost:8080/keluhan?status=selesai

D. Menampilkan keluhan berdasarkan id
curl -i http://localhost:8080/keluhan/1

E. Menambahkan keluhan baru
curl -i -X POST http://localhost:8080/keluhan \
-H "Content-Type: application/json" \
-d '{"ruang":"B201","jenis":"proyektor"}'

F. Menyelesaikan keluhan 
curl -i -X POST http://localhost:8080/keluhan/1/selesaikan

G. Menghapus Keluhan
curl -i -X DELETE http://localhost:8080/keluhan/1

# 4. Kebijakan Retry
Jumlah percobaan : minimal 3 kali dan maksimal 5 kali
Jeda waktu : dimulai dari 200ms dan di kali 2 hingga terakhir ke 5 kali di 1.6 detik
Tidak menggunakan retry : POST tidak menggunakan retry, karena ada kemungkinan jika data telah di proses namun server membalas gagal dapat terjadi data terkirim lebih dari satu kali

# 5. Pemakaian AI
Project ini dibuat dengan panduan dari ChatGPT AI. ChatGPT AI membantu pengerjaan serta memberikan pertanyaan dan penjelasan terkait code serta pengetahuan umum terkait endpoint.
