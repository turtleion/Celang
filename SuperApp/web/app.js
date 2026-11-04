const kunciToken = "superapp.token";
const $ = (id) => document.getElementById(id);
const pesan = (teks) => { $("pesan").textContent = teks; };
const tampil = (el, data) => { el.textContent = JSON.stringify(data, null, 2); };

function token() {
  return sessionStorage.getItem(kunciToken) || "";
}

function headers() {
  return {
    "Content-Type": "application/json",
    "Authorization": "Bearer " + token()
  };
}

async function panggil(jalur, opsi) {
  const res = await fetch(jalur, opsi);
  const data = await res.json().catch(() => ({ ok: false, galat: "respons bukan JSON" }));
  if (!res.ok) {
    throw new Error(data.galat || ("HTTP " + res.status));
  }
  return data;
}

function bukaKonsol(nama, peran) {
  $("login").hidden = true;
  $("konsol").hidden = false;
  $("siapa").textContent = nama + " (" + peran + ")";
  muatCatatan();
}

$("form-masuk").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  pesan("");
  try {
    const data = await panggil("/api/masuk", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ nama: $("nama").value, sandi: $("sandi").value })
    });
    sessionStorage.setItem(kunciToken, data.token);
    bukaKonsol($("nama").value, data.peran);
  } catch (e) {
    pesan(e.message);
  }
});

$("keluar").addEventListener("click", async () => {
  try { await panggil("/api/keluar", { headers: headers() }); } catch (e) {}
  sessionStorage.removeItem(kunciToken);
  $("konsol").hidden = true;
  $("login").hidden = false;
  $("keluaran").textContent = "";
});

$("kirim").addEventListener("click", async () => {
  $("keluaran").textContent = "memproses...";
  try {
    const data = await panggil("/api/konsol", {
      method: "POST",
      headers: headers(),
      body: JSON.stringify({ perintah: $("perintah").value, argumen: $("argumen").value })
    });
    tampil($("keluaran"), data);
  } catch (e) {
    $("keluaran").textContent = "galat: " + e.message;
  }
});

async function muatCatatan() {
  try {
    const data = await panggil("/api/catatan", { headers: headers() });
    const ul = $("daftar-catatan");
    ul.innerHTML = "";
    for (const item of data.catatan) {
      const li = document.createElement("li");
      const tombol = document.createElement("button");
      tombol.className = "kecil";
      tombol.textContent = item.id.slice(0, 8);
      tombol.addEventListener("click", () => bukaCatatan(item.id));
      li.appendChild(tombol);
      li.appendChild(document.createTextNode(" " + new Date(item.dibuat).toLocaleString()));
      ul.appendChild(li);
    }
  } catch (e) {
    $("daftar-catatan").textContent = "galat: " + e.message;
  }
}

async function bukaCatatan(id) {
  try {
    const data = await panggil("/api/catatan/" + encodeURIComponent(id), { headers: headers() });
    tampil($("isi-catatan"), data);
  } catch (e) {
    $("isi-catatan").textContent = "galat: " + e.message;
  }
}

$("simpan").addEventListener("click", async () => {
  try {
    const data = await panggil("/api/catatan", {
      method: "POST",
      headers: headers(),
      body: JSON.stringify({ isi: $("catatan").value })
    });
    $("catatan").value = "";
    tampil($("isi-catatan"), data);
    muatCatatan();
  } catch (e) {
    $("isi-catatan").textContent = "galat: " + e.message;
  }
});

$("muat-statistik").addEventListener("click", async () => {
  try {
    const data = await panggil("/api/statistik", { headers: headers() });
    tampil($("statistik"), data);
  } catch (e) {
    $("statistik").textContent = "galat: " + e.message;
  }
});

if (token()) {
  muatCatatan();
  $("login").hidden = true;
  $("konsol").hidden = false;
  $("siapa").textContent = "sesi tersimpan";
}
