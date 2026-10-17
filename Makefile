# ==============================================================================
# 🌟 Celang Makefile
# Sistem Otomasi Build Celang & Dokumentasi VitePress
# ==============================================================================

BINARY_NAME := celang
CMD_PATH    := ./cmd/celang
GO          := go
NPM         := npm

.PHONY: all build build-all docs-build docs-dev docs-preview preview test run clean help

# Target default: Mengompilasi binary celang
all: build

## 🔨 Build Binary Celang
build:
	@echo "📦 Mengompilasi binary $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "✅ Binary '$(BINARY_NAME)' siap digunakan! (./$(BINARY_NAME))"

## 🌐 Dokumentasi: Build aset statis VitePress
docs-build:
	@echo "📚 Membangun dokumentasi VitePress..."
	$(NPM) run docs:build

## 🚀 Dokumentasi: Auto-Build & Run Preview Server (Sesuai Permintaan)
docs-preview: docs-build
	@echo "🌐 Menjalankan server preview dokumentasi VitePress..."
	$(NPM) run docs:preview

## ⚡ Shortcut cepat untuk docs-preview
preview: docs-preview

## 💻 Dokumentasi: Menjalankan server pengembangan lokal (Hot-Reload)
docs-dev:
	@echo "🔥 Menjalankan live dev server dokumentasi..."
	$(NPM) run docs:dev

## 📦 Build Celang sekaligus Dokumentasi
build-all: build docs-build
	@echo "🎉 Seluruh build (Celang & Dokumentasi) selesai!"

## 🧪 Menjalankan seluruh pengujian unit Celang
test:
	@echo "🧪 Menjalankan pengujian (go test)..."
	$(GO) test -count=1 ./...

## 🏃 Menjalankan skrip Celang (Contoh: make run ARGS="examples/universal_mode.ce")
run: build
	@if [ -z "$(ARGS)" ]; then \
		echo "ℹ️  Menjalankan REPL interaktif..."; \
		./$(BINARY_NAME); \
	else \
		echo "▶️  Menjalankan: ./$(BINARY_NAME) run $(ARGS)"; \
		./$(BINARY_NAME) run $(ARGS); \
	fi

## 🧹 Membersihkan artefak build binary dan cache dokumentasi
clean:
	@echo "🧹 Membersihkan artefak build..."
	rm -f $(BINARY_NAME)
	rm -rf docs/.vitepress/dist docs/.vitepress/cache
	@echo "✨ Bersih!"

## ❓ Bantuan / Panduan Perintah
help:
	@echo "🌟 Panduan Penggunaan Makefile Celang:"
	@echo "  make build         - Kompilasi binary compiler & VM Celang"
	@echo "  make docs-preview  - Otomatis build dokumentasi lalu jalankan server preview"
	@echo "  make preview       - Alias cepat untuk 'make docs-preview'"
	@echo "  make docs-build    - Bangun aset HTML statis dokumentasi VitePress"
	@echo "  make docs-dev      - Jalankan server dokumentasi mode dev (hot-reload)"
	@echo "  make build-all     - Kompilasi binary Celang sekaligus build dokumentasi"
	@echo "  make test          - Jalankan seluruh unit test Celang"
	@echo "  make run [ARGS=..] - Jalankan program Celang (default: REPL)"
	@echo "  make clean         - Hapus file binary dan cache dist"
	@echo "  make help          - Menampilkan bantuan ini"
