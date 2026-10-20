<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useData, useRoute } from 'vitepress'

interface HeaderItem {
  level: number
  title: string
  slug: string
  link: string
  children?: HeaderItem[]
}

const { page, frontmatter } = useData()
const route = useRoute()

const isOpen = ref(true)
const searchQuery = ref('')
const activeLink = ref('')

const rawHeaders = computed<HeaderItem[]>(() => {
  return (page.value.headers || []) as HeaderItem[]
})

// Clean title of backticks or markdown tags
function cleanTitle(title: string): string {
  if (!title) return ''
  return title.replace(/`([^`]+)`/g, '$1').trim()
}

// Filter headers based on search query
const filteredHeaders = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return rawHeaders.value

  return rawHeaders.value
    .map(h => {
      const titleMatch = cleanTitle(h.title).toLowerCase().includes(query)
      const matchingChildren = (h.children || []).filter(c =>
        cleanTitle(c.title).toLowerCase().includes(query)
      )
      if (titleMatch || matchingChildren.length > 0) {
        return {
          ...h,
          children: titleMatch ? h.children : matchingChildren
        }
      }
      return null
    })
    .filter(Boolean) as HeaderItem[]
})

function scrollTo(link: string) {
  if (typeof window === 'undefined') return
  const id = link.startsWith('#') ? link.slice(1) : link
  const decodedId = decodeURIComponent(id)
  const el = document.getElementById(decodedId)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    activeLink.value = link
    if (window.history && window.history.pushState) {
      window.history.pushState(null, '', link)
    }
  }
}

let observer: IntersectionObserver | null = null

function setupObserver() {
  if (typeof window === 'undefined' || typeof IntersectionObserver === 'undefined') return
  if (observer) {
    observer.disconnect()
    observer = null
  }

  const allSlugs: string[] = []
  rawHeaders.value.forEach(h => {
    allSlugs.push(h.slug)
    if (h.children) {
      h.children.forEach(c => allSlugs.push(c.slug))
    }
  })

  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          activeLink.value = '#' + entry.target.id
          break
        }
      }
    },
    {
      rootMargin: '-80px 0px -70% 0px',
      threshold: 0
    }
  )

  allSlugs.forEach(slug => {
    const el = document.getElementById(decodeURIComponent(slug))
    if (el && observer) observer.observe(el)
  })
}

onMounted(() => {
  setupObserver()
})

onUnmounted(() => {
  if (observer) observer.disconnect()
})

watch(() => route.path, () => {
  searchQuery.value = ''
  if (typeof window !== 'undefined') {
    setTimeout(setupObserver, 250)
  }
})
</script>

<template>
  <nav
    v-if="rawHeaders.length >= 2 && frontmatter.pageToc !== false"
    class="page-toc-card"
    aria-label="Daftar Isi Halaman Ini"
  >
    <div class="page-toc-header">
      <div class="page-toc-title-area">
        <span class="page-toc-icon">📑</span>
        <span class="page-toc-heading">Daftar Isi Halaman</span>
        <span class="page-toc-badge">{{ rawHeaders.length }} Bagian</span>
      </div>

      <div class="page-toc-controls">
        <div v-if="rawHeaders.length >= 4" class="page-toc-search-wrap">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari bagian..."
            class="page-toc-search-input"
            aria-label="Cari bagian di halaman ini"
          />
          <span
            v-if="searchQuery"
            class="page-toc-search-clear"
            @click="searchQuery = ''"
            title="Hapus pencarian"
          >✕</span>
        </div>

        <button
          type="button"
          class="page-toc-toggle-btn"
          @click="isOpen = !isOpen"
          :aria-expanded="isOpen"
        >
          <span>{{ isOpen ? 'Sembunyikan' : 'Tampilkan' }}</span>
          <span class="page-toc-toggle-arrow" :class="{ rotated: !isOpen }">▾</span>
        </button>
      </div>
    </div>

    <!-- Collapsible Content -->
    <div v-show="isOpen" class="page-toc-body">
      <div v-if="filteredHeaders.length === 0" class="page-toc-empty">
        Tidak ditemukan topik yang cocok dengan "<em>{{ searchQuery }}</em>"
      </div>

      <div v-else class="page-toc-grid">
        <div
          v-for="(header, index) in filteredHeaders"
          :key="header.link"
          class="page-toc-item-card"
          :class="{ 'is-active': activeLink === header.link }"
        >
          <div class="page-toc-item-header">
            <span class="page-toc-num">{{ String(index + 1).padStart(2, '0') }}</span>
            <a
              :href="header.link"
              class="page-toc-item-link"
              @click.prevent="scrollTo(header.link)"
            >
              {{ cleanTitle(header.title) }}
            </a>
          </div>

          <!-- Subsections (H3) -->
          <ul v-if="header.children && header.children.length > 0" class="page-toc-sublist">
            <li
              v-for="sub in header.children"
              :key="sub.link"
              class="page-toc-subitem"
              :class="{ 'is-active': activeLink === sub.link }"
            >
              <a
                :href="sub.link"
                class="page-toc-sublink"
                @click.prevent="scrollTo(sub.link)"
              >
                <span class="page-toc-subdot">›</span>
                <span class="page-toc-subtitle">{{ cleanTitle(sub.title) }}</span>
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.page-toc-card {
  margin: 24px 0 32px 0;
  padding: 16px 20px;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-left: 4px solid var(--vp-c-brand-1);
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
  transition: border-color 0.25s, box-shadow 0.25s;
}

.page-toc-card:hover {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.07);
}

.page-toc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.page-toc-title-area {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-toc-icon {
  font-size: 1.15rem;
  line-height: 1;
}

.page-toc-heading {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--vp-c-text-1);
  letter-spacing: -0.01em;
}

.page-toc-badge {
  font-size: 0.72rem;
  font-weight: 600;
  color: var(--vp-c-brand-1);
  background: var(--vp-c-brand-soft);
  padding: 2px 8px;
  border-radius: 9999px;
  letter-spacing: 0.02em;
}

.page-toc-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.page-toc-search-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.page-toc-search-input {
  width: 140px;
  height: 28px;
  padding: 0 22px 0 8px;
  font-size: 0.8rem;
  color: var(--vp-c-text-1);
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
  outline: none;
  transition: all 0.2s ease;
}

.page-toc-search-input:focus {
  width: 170px;
  border-color: var(--vp-c-brand-1);
  box-shadow: 0 0 0 2px var(--vp-c-brand-soft);
}

.page-toc-search-clear {
  position: absolute;
  right: 6px;
  font-size: 0.7rem;
  color: var(--vp-c-text-3);
  cursor: pointer;
}

.page-toc-toggle-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 28px;
  padding: 0 10px;
  font-size: 0.8rem;
  font-weight: 500;
  color: var(--vp-c-text-2);
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.page-toc-toggle-btn:hover {
  color: var(--vp-c-brand-1);
  border-color: var(--vp-c-brand-1);
}

.page-toc-toggle-arrow {
  display: inline-block;
  font-size: 0.85rem;
  transition: transform 0.2s ease;
}

.page-toc-toggle-arrow.rotated {
  transform: rotate(-90deg);
}

.page-toc-body {
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px dashed var(--vp-c-divider);
}

.page-toc-empty {
  padding: 12px 0;
  font-size: 0.85rem;
  color: var(--vp-c-text-3);
  text-align: center;
}

.page-toc-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 10px;
}

.page-toc-item-card {
  padding: 10px 12px;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  transition: all 0.2s ease;
}

.page-toc-item-card:hover {
  border-color: var(--vp-c-brand-2);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.page-toc-item-card.is-active {
  border-color: var(--vp-c-brand-1);
  background: var(--vp-c-brand-soft);
}

.page-toc-item-header {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.page-toc-num {
  flex-shrink: 0;
  font-family: var(--vp-font-family-mono);
  font-size: 0.72rem;
  font-weight: 700;
  color: var(--vp-c-brand-1);
  background: var(--vp-c-brand-soft);
  padding: 1px 6px;
  border-radius: 4px;
}

.page-toc-item-link {
  font-size: 0.88rem;
  font-weight: 600;
  color: var(--vp-c-text-1);
  text-decoration: none;
  line-height: 1.4;
  transition: color 0.2s ease;
}

.page-toc-item-link:hover {
  color: var(--vp-c-brand-1);
}

.page-toc-sublist {
  list-style: none;
  margin: 8px 0 0 0;
  padding: 4px 0 0 16px;
  border-left: 2px solid var(--vp-c-divider);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-toc-subitem {
  line-height: 1.35;
}

.page-toc-sublink {
  display: inline-flex;
  align-items: baseline;
  gap: 6px;
  font-size: 0.8rem;
  color: var(--vp-c-text-2);
  text-decoration: none;
  transition: color 0.15s ease;
}

.page-toc-sublink:hover {
  color: var(--vp-c-brand-1);
}

.page-toc-subdot {
  color: var(--vp-c-brand-1);
  font-weight: bold;
}

.page-toc-subitem.is-active .page-toc-sublink {
  color: var(--vp-c-brand-1);
  font-weight: 600;
}

@media (max-width: 640px) {
  .page-toc-grid {
    grid-template-columns: 1fr;
  }
  .page-toc-search-input {
    width: 110px;
  }
  .page-toc-search-input:focus {
    width: 130px;
  }
}
</style>
