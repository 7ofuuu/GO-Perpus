package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Buku struct mewakili detail sebuah buku
type Buku struct {
	Judul    string
	Penulis  string
	ISBN     string
	Salinan  int
	Peminjam []Peminjam // Untuk melacak siapa yang meminjam buku dan tanggal jatuh temponya
}

// Peminjam struct mewakili seorang peminjam dengan nama dan tanggal jatuh tempo
type Peminjam struct {
	Nama     string
	TglJatuh time.Time
}

// Perpustakaan struct mewakili perpustakaan dengan koleksi bukunya
type Perpustakaan struct {
	Buku []Buku
}

// TambahBuku menambahkan buku baru ke perpustakaan
func TambahBuku(perpus *Perpustakaan, buku Buku) {
	perpus.Buku = append(perpus.Buku, buku)
}

// CariBuku mencari buku berdasarkan judul atau penulis
func CariBuku(perpus *Perpustakaan, kataKunci string) []Buku {
	hasil := []Buku{}
	for _, buku := range perpus.Buku {
		if strings.Contains(strings.ToLower(buku.Judul), strings.ToLower(kataKunci)) || strings.Contains(strings.ToLower(buku.Penulis), strings.ToLower(kataKunci)) {
			hasil = append(hasil, buku)
		}
	}
	return hasil
}

// EditBuku mengedit detail sebuah buku di perpustakaan
func EditBuku(perpus *Perpustakaan, isbn string, detailBaru Buku) {
	for i := range perpus.Buku {
		if perpus.Buku[i].ISBN == isbn {
			perpus.Buku[i] = detailBaru
			break
		}
	}
}

// HapusBuku menghapus sebuah buku dari perpustakaan
func HapusBuku(perpus *Perpustakaan, isbn string) {
	for i := range perpus.Buku {
		if perpus.Buku[i].ISBN == isbn {
			perpus.Buku = append(perpus.Buku[:i], perpus.Buku[i+1:]...)
			break
		}
	}
}

// PinjamBuku meminjam sebuah buku dari perpustakaan
func PinjamBuku(perpus *Perpustakaan, isbn string, peminjam string) {
	indeksBuku := cariIndeksBuku(perpus, isbn)
	if indeksBuku == -1 {
		fmt.Println("Buku tidak ditemukan.")
		return
	}

	if !salinanTersedia(perpus, indeksBuku) {
		fmt.Println("Tidak ada salinan yang tersedia untuk dipinjam.")
		return
	}

	if sudahMeminjamBukuYangSama(perpus, indeksBuku, peminjam) {
		fmt.Println("Anda sudah meminjam buku ini.")
		return
	}

	if !dapatMeminjamLebihBanyakBuku(perpus, peminjam) {
		fmt.Println("Anda tidak dapat meminjam lebih dari 3 buku.")
		return
	}

	meminjamBukuDariIndeks(perpus, indeksBuku, peminjam)
}

// cariIndeksBuku mencari indeks buku berdasarkan ISBN
func cariIndeksBuku(perpus *Perpustakaan, isbn string) int {
	for i, buku := range perpus.Buku {
		if buku.ISBN == isbn {
			return i
		}
	}
	return -1
}

// salinanTersedia memeriksa apakah ada salinan yang tersedia untuk dipinjam
func salinanTersedia(perpus *Perpustakaan, indeksBuku int) bool {
	return perpus.Buku[indeksBuku].Salinan > 0
}

// sudahMeminjamBukuYangSama memeriksa apakah pengguna sudah meminjam buku yang sama
func sudahMeminjamBukuYangSama(perpus *Perpustakaan, indeksBuku int, peminjam string) bool {
	for _, p := range perpus.Buku[indeksBuku].Peminjam {
		if p.Nama == peminjam {
			return true
		}
	}
	return false
}

// dapatMeminjamLebihBanyakBuku memeriksa apakah pengguna telah meminjam kurang dari 3 buku
func dapatMeminjamLebihBanyakBuku(perpus *Perpustakaan, peminjam string) bool {
	hitung := 0
	for _, buku := range perpus.Buku {
		for _, p := range buku.Peminjam {
			if p.Nama == peminjam {
				hitung++
				if hitung >= 3 {
					return false
				}
			}
		}
	}
	return true
}

// meminjamBukuDariIndeks menangani proses peminjaman sebenarnya
func meminjamBukuDariIndeks(perpus *Perpustakaan, indeksBuku int, peminjam string) {
	tglJatuh := time.Now().AddDate(0, 0, 7)
	perpus.Buku[indeksBuku].Salinan--
	perpus.Buku[indeksBuku].Peminjam = append(perpus.Buku[indeksBuku].Peminjam, Peminjam{Nama: peminjam, TglJatuh: tglJatuh})
	fmt.Printf("Buku berhasil dipinjam oleh %s. Tanggal jatuh tempo: %s\n", peminjam, tglJatuh.Format("02-Jan-2006"))
}

// KembalikanBuku mengembalikan sebuah buku ke perpustakaan
func KembalikanBuku(perpus *Perpustakaan, isbn string, peminjam string) {
	for i := range perpus.Buku {
		if perpus.Buku[i].ISBN == isbn {
			for j := range perpus.Buku[i].Peminjam {
				if perpus.Buku[i].Peminjam[j].Nama == peminjam {
					perpus.Buku[i].Salinan++
					perpus.Buku[i].Peminjam = append(perpus.Buku[i].Peminjam[:j], perpus.Buku[i].Peminjam[j+1:]...)
					fmt.Printf("Buku berhasil dikembalikan oleh %s\n", peminjam)
					return
				}
			}
		}
	}
	fmt.Println("Buku tidak dipinjam oleh pengguna ini.")
}

// SimpanPerpustakaan menyimpan data perpustakaan ke sebuah file
func SimpanPerpustakaan(perpus *Perpustakaan, namaFile string) {
	file, err := os.Create(namaFile)
	if err != nil {
		fmt.Println("Kesalahan membuat file:", err)
		return
	}
	defer file.Close()

	for _, buku := range perpus.Buku {
		_, err := fmt.Fprintf(file, "%s,%s,%s,%d\n", buku.Judul, buku.Penulis, buku.ISBN, buku.Salinan)
		if err != nil {
			fmt.Println("Kesalahan menulis ke file:", err)
			return
		}
	}
}

// MuatPerpustakaan memuat data perpustakaan dari sebuah file
func MuatPerpustakaan(perpus *Perpustakaan, namaFile string) {
	file, err := os.Open(namaFile)
	if err != nil {
		fmt.Println("Kesalahan membuka file:", err)
		return
	}
	defer file.Close()

	var judul, penulis, isbn string
	var salinan int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		baris := scanner.Text()
		fields := strings.Split(baris, ",")
		if len(fields) != 4 {
			continue
		}
		judul = fields[0]
		penulis = fields[1]
		isbn = fields[2]
		fmt.Sscanf(fields[3], "%d", &salinan)
		perpus.Buku = append(perpus.Buku, Buku{Judul: judul, Penulis: penulis, ISBN: isbn, Salinan: salinan})
	}
}

// ambilInput adalah fungsi pembantu untuk mendapatkan input dari pengguna
func ambilInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + " ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	// Inisialisasi perpustakaan
	perpustakaan := Perpustakaan{}

	// Muat data perpustakaan dari file
	MuatPerpustakaan(&perpustakaan, "perpustakaan.txt")

	for {
		fmt.Println("\nPilihan:")
		fmt.Println("1. Tambah Buku")
		fmt.Println("2. Hapus Buku")
		fmt.Println("3. Edit Buku")
		fmt.Println("4. Cari Buku")
		fmt.Println("5. Tampilkan Semua Buku")
		fmt.Println("6. Pinjam Buku")
		fmt.Println("7. Kembalikan Buku")
		fmt.Println("8. Keluar")

		pilihan := ambilInput("Masukkan pilihan Anda:")

		switch pilihan {
		case "1":
			fmt.Println("Menambahkan buku baru:")
			judul := ambilInput("Judul:")
			penulis := ambilInput("Penulis:")
			isbn := ambilInput("ISBN:")
			salinan := ambilInput("Jumlah salinan:")
			var salinanInt int
			fmt.Sscanf(salinan, "%d", &salinanInt)
			TambahBuku(&perpustakaan, Buku{Judul: judul, Penulis: penulis, ISBN: isbn, Salinan: salinanInt})
			fmt.Println("Buku berhasil ditambahkan.")
			// Simpan data perpustakaan ke file setelah menambahkan buku baru
			SimpanPerpustakaan(&perpustakaan, "perpustakaan.txt")
		case "2":
			isbn := ambilInput("Masukkan ISBN buku yang akan dihapus:")
			HapusBuku(&perpustakaan, isbn)
			fmt.Println("Buku berhasil dihapus.")
			// Simpan data perpustakaan ke file setelah menghapus buku
			SimpanPerpustakaan(&perpustakaan, "perpustakaan.txt")
		case "3":
			isbn := ambilInput("Masukkan ISBN buku yang akan diedit:")
			judulBaru := ambilInput("Judul baru:")
			penulisBaru := ambilInput("Penulis baru:")
			isbnBaru := ambilInput("ISBN baru:")
			salinanBaru := ambilInput("Jumlah salinan baru:")
			var salinanBaruInt int
			fmt.Sscanf(salinanBaru, "%d", &salinanBaruInt)
			EditBuku(&perpustakaan, isbn, Buku{Judul: judulBaru, Penulis: penulisBaru, ISBN: isbnBaru, Salinan: salinanBaruInt})
			fmt.Println("Buku berhasil diedit.")
			// Simpan data perpustakaan ke file setelah mengedit buku
			SimpanPerpustakaan(&perpustakaan, "perpustakaan.txt")
		case "4":
			kataKunci := ambilInput("Masukkan judul atau penulis untuk mencari:")
			hasilCari := CariBuku(&perpustakaan, kataKunci)
			fmt.Println("Hasil pencarian:")
			tampilkanBuku(hasilCari)
		case "5":
			fmt.Println("Semua Buku:")
			tampilkanBuku(perpustakaan.Buku)
		case "6":
			isbn := ambilInput("Masukkan ISBN buku yang akan dipinjam:")
			peminjam := ambilInput("Masukkan nama Anda:")
			PinjamBuku(&perpustakaan, isbn, peminjam)
			// Simpan data perpustakaan ke file setelah meminjam buku
			SimpanPerpustakaan(&perpustakaan, "perpustakaan.txt")
		case "7":
			isbn := ambilInput("Masukkan ISBN buku yang akan dikembalikan:")
			peminjam := ambilInput("Masukkan nama Anda:")
			KembalikanBuku(&perpustakaan, isbn, peminjam)
			// Simpan data perpustakaan ke file setelah mengembalikan buku
			SimpanPerpustakaan(&perpustakaan, "perpustakaan.txt")
		case "8":
			fmt.Println("Keluar...")
			return
		default:
			fmt.Println("Pilihan tidak valid. Silakan masukkan angka antara 1 dan 8.")
		}
	}
}

// tampilkanBuku menampilkan daftar buku dalam format tabel
func tampilkanBuku(buku []Buku) {
	fmt.Printf("%-30s %-20s %-15s %-7s %-50s\n", "Judul", "Penulis", "ISBN", "Salinan", "Peminjam")
	for _, b := range buku {
		var detailPeminjam []string
		for _, p := range b.Peminjam {
			detailPeminjam = append(detailPeminjam, fmt.Sprintf("%s (Jatuh tempo: %s)", p.Nama, p.TglJatuh.Format("02-Jan-2006")))
		}
		peminjam := strings.Join(detailPeminjam, ", ")
		fmt.Printf("%-30s %-20s %-15s %-7d %-50s\n", b.Judul, b.Penulis, b.ISBN, b.Salinan, peminjam)
	}
}
