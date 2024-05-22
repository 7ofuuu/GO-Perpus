package main

import (
	"fmt"
)

const NMAX int = 100

type Buku struct {
	Judul   string
	Penulis string
	ISBN    int
	Salinan int // Untuk melacak siapa yang meminjam buku dan tanggal jatuh temponya
}
type tabBuku [NMAX]Buku

func main() {
	var data tabBuku
	var nData, pilihan int
	fmt.Println("masukan banyak buku")
	fmt.Scan(&nData)
	tambahData(&data, nData)

	for pilihan != 3 {
		fmt.Println()
		fmt.Println("----------------------------------------")
		fmt.Println("                 MENU                   ")
		fmt.Println("----------------------------------------")
		fmt.Println("1. pencarian buku berdasarkan judul")
		fmt.Println("2. pencarian buku berdasarkan penulis")
		fmt.Println("3. exit")
		fmt.Println()

		fmt.Println("masukkan pilihan")
		fmt.Scan(&pilihan)

		if pilihan == 1 {
			pencarianjudul(data, nData)

		} else if pilihan == 2 {
			pencarianpengarang(data, nData)

		}
	}

}

func tambahData(a *tabBuku, n int) {
	fmt.Println("masukkan data buku(Judul, Penulis, ISBN, Salinan):")
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i].Judul, &a[i].Penulis, &a[i].ISBN, &a[i].Salinan)
	}
}

func pencarianjudul(a tabBuku, n int) {
	var nama string

	fmt.Println("masukkan judul buku:")
	fmt.Scan(&nama)

	for i := 0; i < n; i++ {
		if a[i].Judul == nama {
			fmt.Println(a[i])
			return
		}

	}
}

func pencarianpengarang(a tabBuku, n int) {
	var nama string

	fmt.Println("masukkan penulis buku:")
	fmt.Scan(&nama)

	for i := 0; i < n; i++ {
		if a[i].Penulis == nama {
			fmt.Println(a[i])
			return
		}

	}
}
