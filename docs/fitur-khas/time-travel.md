# Variabel Time-Travel (Penjelajah Waktu)

Salah satu keunikan paling menonjol dari Celang adalah fitur **Time-Travel Variables**. Variabel di Celang tidak hanya menyimpan nilai saat ini, melainkan secara alami mengingat masa lalunya.

---

## ⏳ Properti Penjelajah Waktu

Setiap variabel yang dideklarasikan dengan `buat` secara otomatis memiliki 3 properti waktu:

| Properti | Padanan Inggris | Deskripsi |
| :--- | :--- | :--- |
| `.sebelumnya` | `.prev`, `.previous` | Mengembalikan nilai tepat satu langkah sebelum nilai saat ini. Jika belum pernah dimutasi, mengembalikan `null`. |
| `.awal` | `.initial`, `.first` | Mengembalikan nilai saat variabel pertama kali dideklarasikan (*genesis state*). |
| `.histori` | `.history` | Mengembalikan sebuah `Array` berisi linimasa rekaman mutasi variabel dari waktu ke waktu. |

---

## 📝 Contoh Penggunaan

```celang
// Deklarasi awal
buat harga = 50000

// Perubahan harga
harga = 75000
harga = 90000
harga = 120000

cetak("Harga sekarang   :", harga)             // 120000
cetak("Harga sebelumnya :", harga.sebelumnya)  // 90000
cetak("Harga awal mula  :", harga.awal)        // 50000
cetak("Histori kenaikan :", harga.histori)     // [50000, 75000, 90000, 120000]
```

---

## 🛡️ Jaminan Keamanan Memori: Arsitektur Ring-Buffer $O(1)$

Pertanyaan umum dari para insinyur sistem:
> *"Apakah menyimpan riwayat variabel akan memicu kebocoran memori (memory leak) jika variabel diubah jutaan kali dalam loop?"*

**Jawabannya: TIDAK.** 

Celang mengimplementasikan arsitektur **Anchor Snapshot + Bounded Sliding Ring-Buffer**:
1. **Anchor Snapshot (`InitialVal`)**: Nilai awal pertama kali disimpan permanen dalam sebuah pointer tetap ($O(1)$).
2. **Snapshot Cepat (`PreviousVal`)**: Nilai mutasi tepat sebelum nilai aktif selalu tersimpan di pointer terdedikasi.
3. **Bounded Window History**: Array riwayat dibatasi pada kapasitas tetap (default: 50 mutasi terakhir) menggunakan mekanisme geser in-place zero-allocation (`copy` buffer).

Artinya:
Bahkan jika sebuah variabel diubah **10.000.000 kali** di dalam loop, penggunaan memorinya **tetap konstan $O(1)$**, tidak pernah membengkak, dan memori lama dibersihkan secara aman oleh Garbage Collector.
