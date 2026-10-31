# Pemrograman Konkuren & Asinkron (`sambil` & `bersamaan`)

Dalam kehidupan sehari-hari, kita sering melakukan banyak hal sekaligus: mendengarkan musik **sambil** menyapu lantai, atau memasak nasi dan memanggang ikan secara **bersamaan** agar matang pada waktu yang sama.

Celang menghadirkan kemampuan multitasking alami ini menggunakan dua kata yang sangat akrab: **`sambil`** dan **`bersamaan`**.

---

## 1. Mengerjakan di Latar Belakang: `sambil`

Gunakan kata kunci `sambil` ketika Anda ingin menjalankan suatu pekerjaan di latar belakang tanpa membuat program utama berhenti menunggu.

### Analogi
Seperti menyalakan mesin cuci sambil Anda melanjutkan mengetik di komputer. Anda tidak perlu berdiri di depan mesin cuci; Anda langsung lanjut bekerja.

```celang
cetak("Memulai program...")

sambil {
    cetak("[Latar] Mengunduh berkas pembaruan...")
    waktu.tidur(100)
    cetak("[Latar] Berkas berhasil diunduh!")
}

cetak("Program utama langsung melayani pengguna tanpa menunggu...")
waktu.tidur(150)
cetak("Selesai!")
```

### Hasil yang Terlihat di Layar:
```text
Memulai program...
Program utama langsung melayani pengguna tanpa menunggu...
[Latar] Mengunduh berkas pembaruan...
[Latar] Berkas berhasil diunduh!
Selesai!
```

---

## 2. Menyelesaikan Tugas Bersama-sama: `bersamaan`

Jika `sambil` digunakan untuk tugas latar belakang yang tidak perlu ditunggu, kata kunci `bersamaan` digunakan ketika Anda memiliki beberapa tugas yang ingin dikerjakan **secara serentak**, dan Anda **harus menunggu semua tugas itu selesai** sebelum melangkah ke tahap berikutnya.

### Analogi
Seperti tiga orang sahabat yang membagi tugas membersihkan rumah: satu orang menyapu, satu mengepel, dan satu mengelap jendela. Mereka bekerja serentak bersama-sama, dan baru akan pergi makan siang setelah ketiga tugas tersebut selesai.

```celang
fungsi olahLaporan(namaBerkas) {
    cetak("Mulai memproses:", namaBerkas)
    waktu.tidur(50)
    cetak("Selesai:", namaBerkas)
}

cetak("--- Memulai Pemrosesan Serentak ---")

bersamaan {
    olahLaporan("laporan_januari.csv")
    olahLaporan("laporan_februari.csv")
    olahLaporan("laporan_maret.csv")
}

// Baris di bawah ini HANYA akan berjalan setelah ketiga laporan selesai diproses!
cetak("--- Semua laporan selesai! Melanjutkan ke tahap akhir. ---")
```

---

## 3. Pilihan Gaya Bahasa

Tindakan multitasking ini juga dapat ditulis dalam gaya bahasa lain sesuai preferensi Anda:

::: code-group

```celang [🇮🇩 Bahasa Indonesia]
sambil {
    unduhData()
}

bersamaan {
    prosesA()
    prosesB()
}
```

```celang [🌾 Basa Jawa]
sinambi {
    unduhData()
}

bebarengan {
    prosesA()
    prosesB()
}
```

```celang [⚡ Gaya JavaScript]
async {
    unduhData()
}

parallel {
    prosesA()
    prosesB()
}
```

:::
