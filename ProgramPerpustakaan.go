package main

import (
	"fmt"
	"strings"
	"time"
)

const NMAX int = 100
const maxHariPinjam int = 7
const dendaPerHari int = 5000

type Buku struct {
	Judul     string
	Penulis   string
	ISBN      int
	Salinan   int
	Available bool // Untuk melacak ketersediaan buku
}

type Peminjaman struct {
	Judul         string
	Penulis       string
	ISBN          int
	TanggalPinjam time.Time // Menyimpan tanggal pinjam sebagai time.Time
}

type tabBuku [NMAX]Buku
type tabPeminjaman [NMAX]Peminjaman

func main() {
	var data tabBuku
	var peminjaman tabPeminjaman
	var nData, nPeminjaman, pilihan int
	var batas int

	for pilihan != 9 {
		fmt.Println()
		fmt.Println("----------------------------------------")
		fmt.Println("                 MENU                   ")
		fmt.Println("----------------------------------------")
		fmt.Println("1. tambah data buku")
		fmt.Println("2. tampilkan data buku")
		fmt.Println("3. pencarian buku berdasarkan judul")
		fmt.Println("4. pencarian buku berdasarkan penulis")
		fmt.Println("5. edit data")
		fmt.Println("6. hapus data")
		fmt.Println("7. peminjaman buku")
		fmt.Println("8. pengembalian buku")
		fmt.Println("9. exit")
		fmt.Println()

		fmt.Println("masukkan pilihan")
		fmt.Scan(&pilihan)

		if pilihan == 1 {
			fmt.Println("Masukan jumlah buku yang ingin ditambahkan")
			fmt.Scan(&batas)
			tambahData(&data, &nData, batas)
			// tambahData(&data, &nData)
			insertionSort(&data, nData)
		} else if pilihan == 2 {
			cetakData(&data, nData)
		} else if pilihan == 3 {
			pencarianjudul(data, nData)
		} else if pilihan == 4 {
			pencarianpengarang(data, nData)
		} else if pilihan == 5 {
			editBuku(&data, nData)
		} else if pilihan == 6 {
			hapusBuku(&data, &nData)
		} else if pilihan == 7 {
			pinjamBuku(&data, &peminjaman, &nData, &nPeminjaman)
		} else if pilihan == 8 {
			kembalikanBuku(&data, &peminjaman, &nData, &nPeminjaman)
		} else if pilihan == 9 {
			return
		}
	}
}

func tambahData(a *tabBuku, n *int, lim int) {
	var judul, penulis string
	var isbn, salinan int
	if *n > NMAX {
		*n = NMAX
	}

	fmt.Println("Masukkan data buku (Judul, Penulis, ISBN, Salinan):")
	for i := *n; i < *n+lim; i++ {
		fmt.Scan(&judul, &penulis, &isbn, &salinan)

		// Cek keunikan ISBN sebelum menambahkan
		a[i].Judul = judul
		a[i].Penulis = penulis
		a[i].ISBN = isbn
		a[i].Salinan = salinan
		a[i].Available = true // Inisialisasi semua buku sebagai tersedia
	}

	*n += lim
	klosama(&*a, &*n, lim, isbn)
}
func klosama(a *tabBuku, n *int, batas int, isbn int) {
	var i int = 0

	for i < *n+batas && a[i].ISBN == isbn || a[i].ISBN == a[i+1].ISBN {
		// a[1] = a[i+1]
		i++

		*n--
		fmt.Println("ISBN sudah terdaftar")
		if i > *n {
			return
		}
	}
}

func cetakData(a *tabBuku, n int) {
	fmt.Printf("%-20s %-25s %-20s %-20s %-20s\n", "Judul", "Penulis", "ISBN", "Salinan", "Status")
	fmt.Println("--------------------------------------------------------------------------------------------------")
	for i := 0; i < n; i++ {
		status := "Available"
		if !a[i].Available {
			status = "Borrowed"
		}
		fmt.Printf("%-20s %-25s %-20d %-20d %-20s\n", a[i].Judul, a[i].Penulis, a[i].ISBN, a[i].Salinan, status)
	}
}

func pencarianjudul(a tabBuku, n int) {
	var nama string

	fmt.Println("masukkan judul buku yang ingin dicari:")
	fmt.Scan(&nama)

	fmt.Printf("%-20s %-25s %-20s %-20s %-20s\n", "Judul", "Penulis", "ISBN", "Salinan", "Status")
	fmt.Println("--------------------------------------------------------------------------------------------------")
	for i := 0; i < n; i++ {
		status := "Available"
		if !a[i].Available {
			status = "Borrowed"
		}
		if strings.ToLower(a[i].Judul) == strings.ToLower(nama) {
			fmt.Printf("%-20s %-25s %-20d %-20d %-20s\n", a[i].Judul, a[i].Penulis, a[i].ISBN, a[i].Salinan, status)
			return
		}
	}
	fmt.Println("Buku Tidak Ditemukan")
}

func pencarianpengarang(a tabBuku, n int) {
	var nama string

	fmt.Println("masukkan penulis buku:")
	fmt.Scan(&nama)

	fmt.Printf("%-20s %-25s %-20s %-20s %-20s\n", "Judul", "Penulis", "ISBN", "Salinan", "Status")
	fmt.Println("--------------------------------------------------------------------------------------------------")
	for i := 0; i < n; i++ {
		status := "Available"
		if !a[i].Available {
			status = "Borrowed"
		}
		if strings.ToLower(a[i].Penulis) == strings.ToLower(nama) {
			fmt.Printf("%-20s %-25s %-20d %-20d %-20s\n", a[i].Judul, a[i].Penulis, a[i].ISBN, a[i].Salinan, status)

		}

		if i > n {
			fmt.Println("Data Tidak Ditemukan")
		}
	}
}

func editBuku(A *tabBuku, n int) {
	var pilihan int
	var noIs int
	var index int
	var noIsBaru int
	var isbnTaken bool
	var i int

	fmt.Println("Masukkan no ISBN buku yang ingin diedit: ")
	fmt.Scan(&noIs)

	index = cariISBN(*A, n, noIs)
	if index != -1 {
		fmt.Println("----------------------------------------")
		fmt.Println("     Pilih data yang ingin di edit      ")
		fmt.Println("----------------------------------------")
		fmt.Println("1. Edit Judul Buku.")
		fmt.Println("2. Edit Penulis.")
		fmt.Println("3. Edit ISBN.")
		fmt.Println("4. Edit Salinan.")
		fmt.Println("5. Edit Semua Data.")

		fmt.Println()
		fmt.Println("Masukan pilihan")
		fmt.Scan(&pilihan)
		if pilihan == 1 {
			fmt.Println("Masukkan Judul Baru:")
			fmt.Scan(&A[index].Judul)
		} else if pilihan == 2 {
			fmt.Println("Masukkan Penulis Baru:")
			fmt.Scan(&A[index].Penulis)
		} else if pilihan == 3 {
			isbnTaken = true
			for isbnTaken {
				fmt.Println("Masukkan ISBN :")
				fmt.Scan(&noIsBaru)
				isbnTaken = false
				for i = 0; i < index; i++ {
					if A[i].ISBN == noIsBaru {
						fmt.Println("ISBN sudah terdaftar. Silahkan masukkan ISBN lain")
						isbnTaken = true
					}
				}
			}
			A[index].ISBN = noIsBaru
		} else if pilihan == 4 {
			fmt.Println("Masukkan jumlah salinan:")
			fmt.Scan(&A[index].Salinan)
		} else if pilihan == 5 {
			isbnTaken = true
			for isbnTaken {
				fmt.Println("Masukkan ISBN baru:")
				fmt.Scan(&noIsBaru)
				isbnTaken = false
				for i = 0; i < index; i++ {
					if A[i].ISBN == noIsBaru {
						fmt.Println("ISBN sudah terdaftar. Silahkan masukkan ISBN lain")
						isbnTaken = true
					}
				}
			}
			A[index].ISBN = noIsBaru
			fmt.Println("Masukkan data buku:")
			fmt.Scan(&A[index].Judul, &A[index].Penulis, &A[index].ISBN, &A[index].Salinan)
		} else {
			fmt.Println("Pilihan tidak valid.")
		}
		fmt.Println("Data berhasil diedit.")
	} else {
		fmt.Println("Buku dengan ISBN tersebut tidak ditemukan.")
	}
}

func hapusBuku(A *tabBuku, n *int) {
	var noIs int
	var index int
	var i int

	if *n == 0 {
		fmt.Println("Belum ada buku yang tersedia.")
		return
	} else {
		cetakData(A, *n)
		fmt.Println()
		fmt.Println("Masukkan ISBN buku yang ingin dihapus: ")
		fmt.Scan(&noIs)

		index = cariISBN(*A, *n, noIs)
		if index != -1 {
			for i = index; i < *n-1; i++ {
				A[i] = A[i+1]
			}
			*n--
			fmt.Println("Data buku berhasil dihapus.")
		} else {
			fmt.Println("Buku dengan no ISBN tersebut tidak ditemukan.")
		}
	}
}

func cariISBN(A tabBuku, n int, noIs int) int {
	var left int
	var right int
	var mid int

	right = n - 1

	for left <= right {
		mid = (left + right) / 2
		if A[mid].ISBN == noIs {
			return mid
		} else if A[mid].ISBN < noIs {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

func insertionSort(A *tabBuku, n int) {
	var i, pass int
	var temp Buku

	for pass = 1; pass < n; pass++ {
		temp = A[pass]
		i = pass - 1
		for i >= 0 && A[i].ISBN > temp.ISBN {
			A[i+1] = A[i]
			i--
		}
		A[i+1] = temp
	}
}

func pinjamBuku(data *tabBuku, peminjaman *tabPeminjaman, nData *int, nPeminjaman *int) {
	var judul string
	var hari, bulan, tahun int
	fmt.Println("Masukkan judul buku yang ingin dipinjam:")
	fmt.Scan(&judul)

	index := -1
	for i := 0; i < *nData; i++ {
		if strings.ToLower(data[i].Judul) == strings.ToLower(judul) && data[i].Available {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Buku tidak tersedia untuk dipinjam.")
		return
	}

	fmt.Println("Masukkan tanggal peminjaman (format: hari bulan tahun):")
	fmt.Scan(&hari, &bulan, &tahun)

	tanggalPinjam := time.Date(tahun, time.Month(bulan), hari, 0, 0, 0, 0, time.UTC)
	data[index].Available = false
	peminjaman[*nPeminjaman] = Peminjaman{
		Judul:         data[index].Judul,
		Penulis:       data[index].Penulis,
		ISBN:          data[index].ISBN,
		TanggalPinjam: tanggalPinjam,
	}
	*nPeminjaman++
	fmt.Println("Buku berhasil dipinjam.")
}

func kembalikanBuku(data *tabBuku, peminjaman *tabPeminjaman, nData *int, nPeminjaman *int) {
	var judul string
	var hari, bulan, tahun int
	fmt.Println("Masukkan judul buku yang ingin dikembalikan:")
	fmt.Scan(&judul)

	index := -1
	for i := 0; i < *nPeminjaman; i++ {
		if strings.ToLower(peminjaman[i].Judul) == strings.ToLower(judul) {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Data peminjaman tidak ditemukan.")
		return
	}

	fmt.Println("Masukkan tanggal pengembalian (format: hari bulan tahun):")
	fmt.Scan(&hari, &bulan, &tahun)

	tanggalKembali := time.Date(tahun, time.Month(bulan), hari, 0, 0, 0, 0, time.UTC)
	tanggalPinjam := peminjaman[index].TanggalPinjam
	selisihHari := int(tanggalKembali.Sub(tanggalPinjam).Hours() / 24)

	if selisihHari > maxHariPinjam {
		denda := (selisihHari - maxHariPinjam) * dendaPerHari
		fmt.Printf("Anda terlambat mengembalikan buku. Denda: %d\n", denda)
	} else {
		fmt.Println("Buku dikembalikan tepat waktu.")
	}

	for i := 0; i < *nData; i++ {
		if data[i].ISBN == peminjaman[index].ISBN {
			data[i].Available = true
			break
		}
	}

	for i := index; i < *nPeminjaman-1; i++ {
		peminjaman[i] = peminjaman[i+1]
	}
	*nPeminjaman--

	fmt.Println("Buku berhasil dikembalikan.")
}
W