# Kontrol Alur (Percabangan & Perulangan)

Kontrol alur di Celang menggunakan kosakata yang intuitif dalam bahasa Indonesia, memudahkan pemula maupun pengembang berpengalaman membaca alur logika program seperti membaca teks sehari-hari.

---

## 1. Percabangan (`jika`, `selain jika`, `selain`)

Pernyataan kondisi di Celang mengevaluasi nilai kebenaran (*boolean condition*). Kurung tanda kurung `(...)` bersifat opsional.

### Sintaks Dasar

```celang
buat nilaiUjian = 85

jika nilaiUjian >= 90 {
    cetak("Predikat: A (Istimewa)")
} selain jika nilaiUjian >= 80 {
    cetak("Predikat: B (Sangat Baik)")
} selain jika nilaiUjian >= 70 {
    cetak("Predikat: C (Cukup)")
} selain {
    cetak("Predikat: D (Perlu Perbaikan)")
}
```

### Operator Pembanding & Logika

- **Pembanding**: `==` (sama dengan), `!=` (tidak sama dengan), `>`, `<`, `>=`, `<=`
- **Logika**: `&&` atau `dan`, `||` atau `atau`, `!` atau `bukan`

```celang
buat usia = 20
buat memilikiSim = benar

jika usia >= 17 dan memilikiSim {
    cetak("Boleh mengemudikan kendaraan bermotor.")
}
```

---

## 2. Perulangan Kondisional (`selama`)

Gunakan kata kunci `selama` (*while loop*) untuk menjalankan blok kode berulang kali selama suatu kondisi bernilai `benar`.

```celang
buat hitungan = 1

selama hitungan <= 5 {
    cetak("Perulangan ke-", hitungan)
    hitungan += 1
}

cetak("Perulangan selesai!")
```

---

## 3. Perulangan Koleksi & Rentang (`untuk`)

Kata kunci `untuk` (*for-in loop*) digunakan untuk mengiterasi deretan elemen array atau rentang bilangan.

### Iterasi Rentang Angka (`rentang`)

Fungsi bawaan `rentang(awal, akhir)` menghasilkan deretan angka berurutan:

```celang
// Mencetak angka 1 sampai 5
untuk i dalam rentang(1, 5) {
    cetak("Iterasi angka:", i)
}
```

### Iterasi Elemen Array

```celang
buat bahasaPemrograman = ["Celang", "Go", "Rust", "TypeScript"]

untuk bahasa dalam bahasaPemrograman {
    cetak("Bahasa favorit:", bahasa)
}
```

---

## 4. Pengendali Putaran (`berhenti` & `lanjut`)

- **`berhenti`** (*break*): Menghentikan seluruh proses perulangan seketika.
- **`lanjut`** (*continue*): Melompati sisa blok perulangan saat ini dan langsung berpindah ke iterasi berikutnya.

```celang
untuk angka dalam rentang(1, 10) {
    // Lewati angka genap
    jika angka % 2 == 0 {
        lanjut
    }

    // Berhenti jika sudah mencapai angka 7
    jika angka > 7 {
        berhenti
    }

    cetak("Angka ganjil terpilih:", angka)
}
```

### Output

```text
Angka ganjil terpilih: 1
Angka ganjil terpilih: 3
Angka ganjil terpilih: 5
Angka ganjil terpilih: 7
```
