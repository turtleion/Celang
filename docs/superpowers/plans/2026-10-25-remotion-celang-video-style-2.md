# Remotion Celang Style 2 Video Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a complete, production-grade 60-second (1800 frames @ 30fps) Remotion motion graphics video for Celang (Style 2: Amethyst Cyber-Editorial & Futuristic Purple Glow) in `celang-video-style-2`, matching `PROMPT_VIDEO_MOTION_GRAPHIC_CELANG_STYLE_2.md` specifications with zero subtitles and 100% pure kinetic typography.

**Architecture:** A modular Remotion 4.x TypeScript project inside `celang-video-style-2`. The root composition sequences 9 bespoke scene components with synchronized visual transitions, particle canvas systems, 3D perspective transforms, holographic HUDs, laser wavefront morphing, and an integrated audio engine (neural Indonesian VO + ambient synthwave soundtrack + SFX cues).

**Tech Stack:** React 19 / Remotion 4.x, TypeScript, TailwindCSS, Lucide React (or SVG HUD icons), Python3 / edge-tts + ffmpeg for broadcast audio generation.

**Spec:** `/home/yosatsu/Coding/Celang/PROMPT_VIDEO_MOTION_GRAPHIC_CELANG_STYLE_2.md`

## Global Constraints
- Target Duration: Strict 60.0 seconds (1,800 frames at 30 fps).
- Resolution: 1920 × 1080 (Full HD, 16:9 widescreen).
- ZERO SUBTITLES / NO CAPTIONS: Absolutely no bottom subtitle/caption bar. All text is pure kinetic motion typography in center composition.
- Docs Purple Theme: Deep obsidian canvas (`#0b091a` – `#131127`), Brand Hero Purple (`#a8b1ff`), Electric Amethyst (`#7c3aed`), Luminous Lavender (`#c084fc`), Neon Indigo (`#6366f1`), Specular Chrome/White (`#FFFFFF`).
- Visual Elements: Prismatic crystal glass, neural algorithmic node graphs, holographic HUD typography, ultraviolet laser wavefronts. No cartoon toy puzzle blocks.
- Content: Focus on Scratch transition, breaking English barrier, natural Indonesian syntax, `celang format --style=pythonic`, polite error guidance, web playground. No complex features (no time-travel, concurrency channels, DB transactions).

## Review Focus
- Exact frame alignment across all 9 scenes ensuring no black frames or overlaps between Scene 1..9.
- Total frame count verification: Composition must be exactly 1800 frames.
- Zero subtitle compliance: Verify that no element resembles a bottom subtitle bar.
- Color palette compliance: Ensure all scenes follow the Docs Purple and Obsidian theme.
- Video preview & bundle sanity check: Ensure `npm run build` or Remotion render command compiles without error.

---

### Task 1: Project Scaffolding & Audio Generation
**Files:**
- Create: `celang-video-style-2/package.json`
- Create: `celang-video-style-2/remotion.config.ts`
- Create: `celang-video-style-2/tsconfig.json`
- Create: `celang-video-style-2/tailwind.config.js`
- Create: `celang-video-style-2/scripts/generate_audio.py`
- Output: `celang-video-style-2/public/audio/vo_scene_*.mp3`, `celang-video-style-2/public/audio/soundtrack.mp3`

**Interfaces:**
- Consumes: `edge-tts`, `ffmpeg`
- Produces: Complete Remotion workspace ready to render, with all VO voice tracks and ambient synthwave background audio in `public/audio/`.

- [ ] **Step 1: Scaffold Remotion project in `celang-video-style-2`**
Scaffold Remotion project using `npx create-video --yes --blank celang-video-style-2 --tailwind` or standard configuration.
- [ ] **Step 2: Install required dependencies**
Install `remotion`, `@remotion/cli`, `@remotion/media-utils`, `react`, `react-dom`, `lucide-react`, `tailwindcss`.
- [ ] **Step 3: Create audio generator script using edge-tts and synthesize VO and soundtrack**
Generate neural Indonesian voiceover for each of the 9 scenes with `id-ID-ArdiNeural` with pacing tuned to the scene lengths, plus generate a synchronized future synthwave bed and SFX via ffmpeg.
- [ ] **Step 4: Verify project setup and audio asset existence**
Run `ls -lh public/audio` and verify all audio files exist.

---

### Task 2: Core Design System, Constants, and Background Ambience
**Files:**
- Create: `celang-video-style-2/src/constants.ts`
- Create: `celang-video-style-2/src/components/DeepObsidianCanvas.tsx`
- Create: `celang-video-style-2/src/components/PrismaticGlow.tsx`
- Create: `celang-video-style-2/src/components/HolographicHud.tsx`
- Create: `celang-video-style-2/src/components/ParticlesBackground.tsx`

**Interfaces:**
- Consumes: Palette tokens and duration configs
- Produces: Reusable background environment with radial gradients, floating violet particles, ambient lighting, and HUD corner brackets.

- [ ] **Step 1: Define design tokens, color palette, and scene timings in `constants.ts`**
Define scene frame ranges (Scene 1: 0-180, Scene 2: 180-360, Scene 3: 360-570, Scene 4: 570-780, Scene 5: 780-990, Scene 6: 990-1200, Scene 7: 1200-1440, Scene 8: 1440-1620, Scene 9: 1620-1800).
- [ ] **Step 2: Implement `DeepObsidianCanvas` and `ParticlesBackground`**
Render layered gradient `#0b091a` to `#1e1b4b`, dynamic floating amethyst light particles, and faint isometric grid.
- [ ] **Step 3: Implement `HolographicHud` and `PrismaticGlow`**
Provide holographic technical frame borders, coordinates, specular light sweeps, and violet glow filters.

---

### Task 3: Scene 1 to Scene 3 (The Problem & Void)
**Files:**
- Create: `celang-video-style-2/src/scenes/Scene1AmethystIgnition.tsx`
- Create: `celang-video-style-2/src/scenes/Scene2LogicConstellation.tsx`
- Create: `celang-video-style-2/src/scenes/Scene3ChasmDoubleShock.tsx`

**Interfaces:**
- Consumes: Remotion `spring`, `interpolate`, `useCurrentFrame`
- Produces:
  - Scene 1: UV laser flare, 3D Z-axis extrusion of "INDONESIA", monumental "KODING", "SEMUA INGIN BISA KODING", "TETAPI... BAGAIMANA MEMULAINYA?".
  - Scene 2: Algorithmic node constellation mandala with `(jika)` and `(ulang)` nodes, neon purple data pulses, orbital typography: "VISUAL PROGRAMMING", "SCRATCH: LOGIKA INTUITIF", "BERBAHASA INDONESIA".
  - Scene 3: Crystal shattering break transition, rapid tilt-down into deep void, monolithic grey text pillars with cold code walls, diagonal sliced glitch typography with chromatic aberration: "TETAPI SAAT HARUS NAIK KELAS...", "JURANG PEMISAH YANG CURAM!", "THE DOUBLE SHOCK".

- [ ] **Step 1: Implement `Scene1AmethystIgnition` (0-180f)**
Laser ignite sweep, 3D typography extrusion, metallic chrome flare.
- [ ] **Step 2: Implement `Scene2LogicConstellation` (180-360f)**
Node graph network with SVG connections, pulsing data packets, orbital kinetic text tags.
- [ ] **Step 3: Implement `Scene3ChasmDoubleShock` (360-570f)**
Shattering prism geometry, vertical camera drop, diagonal sliced text with RGB split aberration.

---

### Task 4: Scene 4 to Scene 6 (The Barrier, Dawn of Celang & Native Clarity)
**Files:**
- Create: `celang-video-style-2/src/scenes/Scene4EnglishBarrier.tsx`
- Create: `celang-video-style-2/src/scenes/Scene5AmethystDawn.tsx`
- Create: `celang-video-style-2/src/scenes/Scene6NativeClarity.tsx`

**Interfaces:**
- Consumes: `constants.ts`, Remotion animation primitives
- Produces:
  - Scene 4: 3D rotating kinetic cylinder of keywords (`def`, `while`, `elif`, `SyntaxError`), monolithic cracked obsidian typography "TEMBOK BAHASA INGGRIS" with red-purple inner glow.
  - Scene 5: White flash shockwave, monumental 3D CELANG title explosion in amethyst crystal (`#a8b1ff`) with silver glint and luminous laser bridge spanning across the chasm.
  - Scene 6: Floating holographic HUD code card in 3D perspective, laser scanner typing Celang syntax (`buat target = "Masa Depan"`, `jika siap { cetak(...) }`), lavender glowing keyword pulses.

- [ ] **Step 1: Implement `Scene4EnglishBarrier` (570-780f)**
Kinetic rotating cylinder of syntax obstacles, imposing "TEMBOK BAHASA INGGRIS" typography.
- [ ] **Step 2: Implement `Scene5AmethystDawn` (780-990f)**
Shockwave explosion, monumental 3D Celang logo reveal with chromatic flares and glowing bridge rays.
- [ ] **Step 3: Implement `Scene6NativeClarity` (990-1200f)**
Holographic code editor card with typed syntax, glowing tokens, and clean typography badges.

---

### Task 5: Scene 7 to Scene 9 (The Quantum Morph, Guidance & Monolith CTA)
**Files:**
- Create: `celang-video-style-2/src/scenes/Scene7QuantumMorph.tsx`
- Create: `celang-video-style-2/src/scenes/Scene8HumanGuidance.tsx`
- Create: `celang-video-style-2/src/scenes/Scene9AmethystMonolith.tsx`

**Interfaces:**
- Consumes: Remotion animation primitives
- Produces:
  - Scene 7: Hologram terminal running `celang format --style=pythonic`, ultraviolet laser wavefront scanning left-to-right, quantum morphing Celang code into Python, interlocking luminous orbital rings.
  - Scene 8: Floating amethyst glass pane showing warm error guidance ("💡 Petunjuk Ramah..."), smooth transition into 3D interactive web playground badge sphere ("Playground Web • Coba Langsung di Browser").
  - Scene 9: Monolithic 3D CELANG brand monument with inner bioluminescence, Saturn-ring orbital kinetic text, specular light wipe revealing slogan and call-to-action ("celang.dev • github.com/celang").

- [ ] **Step 1: Implement `Scene7QuantumMorph` (1200-1440f)**
Laser wavefront scanline morph effect switching Indonesian Celang tokens to Python tokens with specular glints.
- [ ] **Step 2: Implement `Scene8HumanGuidance` (1440-1620f)**
Amethyst glass card with polite error message, transitioning to holographic cloud playground sphere.
- [ ] **Step 3: Implement `Scene9AmethystMonolith` (1620-1800f)**
Grand finale monolithic 3D identity, Saturn rings, glowing call-to-action and URL badges.

---

### Task 6: Master Composition Integration & Audio Sync
**Files:**
- Create: `celang-video-style-2/src/Root.tsx`
- Create: `celang-video-style-2/src/index.ts`
- Create: `celang-video-style-2/src/MainVideo.tsx`

**Interfaces:**
- Consumes: All 9 scenes, audio assets, `<Composition>` from Remotion
- Produces: Main composition `CelangStyle2Video` (1920x1080, 30fps, 1800 frames) with all scenes sequenced and audio tracks layered.

- [ ] **Step 1: Assemble `MainVideo.tsx` with `<Sequence>` for each scene and `<Audio>` tracks**
Sequence scenes precisely at frame offsets: 0, 180, 360, 570, 780, 990, 1200, 1440, 1620.
- [ ] **Step 2: Add dynamic audio ducking and SFX layers**
Mix the voice-over sequences with the synthwave background track and synchronized sound effect sweeps.
- [ ] **Step 3: Register composition in `Root.tsx`**
Define `Composition` with `id="CelangStyle2Video"`, `durationInFrames={1800}`, `fps={30}`, `width={1920}`, `height={1080}`.

---

### Task 7: Verification & Build Testing
**Files:**
- Test command: Remotion bundle test and still-frame verification
- Output: Render test or still image inspection

- [ ] **Step 1: Run TypeScript check (`npx tsc --noEmit`)**
Verify there are no type errors across all components.
- [ ] **Step 2: Verify Remotion bundle (`npx remotion bundle`)**
Confirm Remotion compiles the video without bundler errors.
- [ ] **Step 3: Render test frames from multiple scenes**
Render stills at frames 90 (Scene 1), 270 (Scene 2), 460 (Scene 3), 670 (Scene 4), 880 (Scene 5), 1100 (Scene 6), 1320 (Scene 7), 1530 (Scene 8), and 1710 (Scene 9) to verify visual fidelity and strict zero-subtitle compliance.
