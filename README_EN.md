<div align="center">
  <img src="assets/banner.jpg" alt="Q-Solver Banner" width="100%" style="border-radius: 12px; box-shadow: 0 8px 30px rgba(0,0,0,0.12);">

  <h1>🧠 Q-Solver</h1>
  
  <h3>Real-time AI Assistant with Screen Analysis & Voice Intelligence</h3>
  
  <p><i>🎯 Screenshot → Think → Answer — Your invisible AI co-pilot</i></p>

  <p>
    <a href="https://github.com/jym66/Q-Solver/releases"><img src="https://img.shields.io/github/v/release/jym66/Q-Solver?color=blueviolet&label=Latest&style=for-the-badge" alt="Release"></a>
    <a href="https://github.com/jym66/Q-Solver/stargazers"><img src="https://img.shields.io/github/stars/jym66/Q-Solver?color=yellow&style=for-the-badge" alt="Stars"></a>
    <a href="https://github.com/jym66/Q-Solver/releases"><img src="https://img.shields.io/github/downloads/jym66/Q-Solver/total?color=green&style=for-the-badge" alt="Downloads"></a>
  </p>
  
  <p>
    <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white" alt="Vue">
    <img src="https://img.shields.io/badge/Wails-v2-E30613?logo=wails&logoColor=white" alt="Wails">
    <img src="https://img.shields.io/badge/Platform-Windows%20|%20macOS-0078D6?logo=windows&logoColor=white" alt="Platform">
  </p>

  <p>
    <a href="#-features">Features</a> •
    <a href="#-quick-start">Quick Start</a> •
    <a href="#-demo">Demo</a> •
    <a href="#-shortcuts">Shortcuts</a> •
    <a href="#-configuration">Configuration</a> •
    <a href="README.md">中文</a>
  </p>
  
  <br>
  
  <img src="assets/demo.gif" alt="Demo" width="90%" style="border-radius: 12px; box-shadow: 0 8px 30px rgba(0,0,0,0.3);">

</div>

> [!WARNING]
> **Development Stage Notice**: This project is currently in **early development stage**. Features may be unstable and bugs may occur. If you encounter any issues, please submit an Issue!

<br>

## 🔥 Why Q-Solver?

<table>
<tr>
<td width="50%">

### 🖼️ **Screenshot to Answer**
One hotkey captures your screen and gets instant AI analysis. Perfect for:
- 📝 Complex problem solving
- 💻 Code review & debugging  
- 📊 Data analysis
- 🎓 Learning assistance

</td>
<td width="50%">

### 🎙️ **Real-time Voice AI**
Live audio capture with instant AI responses:
- 🗣️ Real-time speech transcription
- 🤖 Instant AI answer suggestions
- 🗺️ Auto-generated mind maps
- ⚡ Ultra-low latency interaction

</td>
</tr>
</table>

<br>

## ✨ Features

### 🛡️ Stealth Mode — "Ghost Window"

| Feature | Description |
|:---:|:---|
| 🚫 **Invisible to Recording** | Most screenshot/screen recording software cannot capture this window |
| 👆 **Click-through** | Can enable click-through to interact with apps behind the window |
| 📌 **Always on Top** | Can be set to float above other windows |
| 🎯 **No Focus Stealing** | Tries to avoid interrupting your current work |

> ⚠️ These features may behave differently depending on your system/software environment. **Please test thoroughly before actual use.**

---

### 🎙️ Gemini Live API — Real-time Voice Interaction

> 💡 **Use Case**: Capture the other party's voice in real-time, AI generates answer suggestions simultaneously

| Feature | Description |
|:---:|:---|
| 🗣️ **Voice Transcription** | Real-time system audio capture and transcription |
| 🤖 **Instant Answers** | AI generates answer suggestions based on transcribed content |
| 🗺️ **Mind Map** | Automatically organize conversations into visual mind maps |
| 📤 **Export Notes** | One-click export to Markdown format |

---

### 🔌 Multi-Model Support — Choose Your AI

| Provider | Example Models | Highlights |
|:---:|:---|:---|
| **OpenAI** | GPT-4o, o1-preview | Strong general capabilities |
| **Gemini** | gemini-2.0-flash-exp | Supports Live API real-time voice |
| **Claude** | Claude 3.5 Sonnet | Excellent long-text understanding |
| **Custom** | Any OpenAI-compatible API | Supports OneAPI and similar services |

---

### 📄 Context Enhancement — Personalized Answers

- **📑 Resume Import**: PDF / Markdown format, AI gives targeted answers based on your background
- **🧠 Multi-turn Memory**: Maintains conversation context, no need to repeat background info
- **✨ Smart Parsing**: One-click parse PDF resume into structured Markdown

<br>

## 📸 Demo

| Screenshot Analysis | Real-time Voice | Mind Map |
|:---:|:---:|:---:|
| <img src="assets/img1.png" width="100%"/> | <img src="assets/img6.png" width="100%"/> | <img src="assets/img7.png" width="100%"/> |

<details>
<summary>📷 More Screenshots</summary>

| | |
|:---:|:---:|
| <img src="assets/img2.png" width="100%"/> | <img src="assets/img3.png" width="100%"/> |
| <img src="assets/img4.png" width="100%"/> | <img src="assets/img5.png" width="100%"/> |

</details>

<br>

## 🚀 Quick Start

### Option 1: Download Release (Recommended)

<a href="https://github.com/jym66/Q-Solver/releases">
  <img src="https://img.shields.io/badge/Download-Latest%20Release-blue?style=for-the-badge&logo=github" alt="Download">
</a>

> **macOS users**: You may need to add execute permission after downloading:
> ```bash
> chmod +x Q-Solver.app/Contents/MacOS/Q-Solver
> ```

### Option 2: Build from Source

```bash
# Prerequisites: Go 1.25+, Node.js 22+, Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clone & Run
git clone https://github.com/jym66/Q-Solver.git
cd Q-Solver
wails dev

# Build production
wails build -ldflags "-s -w" -tags prod
```

<br>

## ⌨️ Shortcuts

<table>
<tr>
<th>Action</th>
<th>Windows</th>
<th>macOS</th>
</tr>
<tr>
<td><b>📸 Screenshot & Solve</b></td>
<td><code>F8</code></td>
<td><code>⌘1</code></td>
</tr>
<tr>
<td>👁️ Show/Hide</td>
<td><code>F9</code></td>
<td><code>⌘2</code></td>
</tr>
<tr>
<td>🖱️ Click-through</td>
<td><code>F10</code></td>
<td><code>⌘3</code></td>
</tr>
<tr>
<td>↕️ Move Window</td>
<td><code>Alt + Arrow</code></td>
<td><code>⌘⌥ + Arrow</code></td>
</tr>
<tr>
<td>📜 Scroll</td>
<td><code>Alt + PgUp/PgDn</code></td>
<td><code>⌘⌥⇧ + ↑/↓</code></td>
</tr>
</table>

> **Note**: Windows supports custom hotkeys. macOS uses fixed shortcuts.

<br>

## ⚙️ Configuration

1. Click **Settings** icon (top-right corner)
2. Choose your **Provider**: OpenAI / Gemini / Claude / Custom
3. Enter your **API Key**
4. Select a **Model**
5. (Optional) Import **Resume/CV** for personalized answers

### Supported Providers

| Provider | Models | Live API |
|----------|--------|----------|
| OpenAI | GPT-4o, o1, etc. | ❌ |
| Gemini | gemini-2.0-flash-exp | ✅ |
| Claude | Claude 3.5+ | ❌ |
| Custom | Any OpenAI-compatible | ❌ |

<br>

## 🍎 macOS Setup

<details>
<summary><b>📸 Screenshot Permission</b></summary>

1. Go to **Settings** → **Screenshot** tab
2. Click **Grant Screenshot Permission**
3. Allow in System Preferences
4. Click **Refresh Permission Status**

</details>

<details>
<summary><b>🎙️ System Audio Capture (for Live API)</b></summary>

macOS requires a virtual audio driver for system audio capture:

```bash
# Install BlackHole
brew install blackhole-2ch
```

Then configure in **Audio MIDI Setup**:
1. Create **Multi-Output Device**
2. Add **BlackHole 2ch** + your speakers
3. Set as system output
4. Set BlackHole input volume to **1.0**

<img src="assets/img5.png" width="500"/>

</details>

<br>

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.25+, Wails v2 |
| Frontend | Vue 3, Vue Flow |
| AI | OpenAI SDK, Google GenAI, Anthropic SDK |
| Audio | malgo (miniaudio), WASAPI/BlackHole |
| UI | Native window APIs, CGO |

<br>

## ⚠️ Disclaimer

> **This project is for technical research and personal learning purposes only. Do not use it for any illegal or unethical purposes.**
> 
> The user assumes all responsibility for any consequences arising from the use of this software. The developer is not liable for any damages.

<br>

## ⭐ Star History

<div align="center">
  <a href="https://star-history.com/#jym66/Q-Solver&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=jym66/Q-Solver&type=Date&theme=dark" />
      <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=jym66/Q-Solver&type=Date" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=jym66/Q-Solver&type=Date" />
    </picture>
  </a>
</div>

<br>

## 📜 License

<p>
This project is licensed under <b>CC BY-NC 4.0</b> — for personal and educational use only.
</p>

<br>

---

<div align="center">
  <p>
    <b>If you find Q-Solver useful, please give it a ⭐ Star!</b>
  </p>
  <p>
    Made with ❤️ by <a href="https://github.com/jym66">jym66</a>
  </p>
</div>
