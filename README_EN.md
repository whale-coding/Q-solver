<div align="center">
  <img src="assets/banner.jpg" alt="Q-Solver Banner" width="100%" style="border-radius: 12px; box-shadow: 0 8px 30px rgba(0,0,0,0.12);">

  <h1 style="font-size: 3rem; margin: 20px 0;">Q-Solver</h1>
  <p style="font-size: 1.2rem; color: #666;">📝 Screenshot, Think, Solve — Your Distraction-Free AI Assistant</p>

  <p>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white" alt="Go"></a>
    <a href="https://vuejs.org"><img src="https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white" alt="Vue"></a>
    <a href="https://wails.io"><img src="https://img.shields.io/badge/Wails-v2-E30613?logo=wails&logoColor=white" alt="Wails"></a>
    <img src="https://img.shields.io/badge/Platform-Windows%20|%20macOS-0078D6?logo=windows&logoColor=white" alt="Platform">
    <img src="https://img.shields.io/badge/License-CC_BY--NC_4.0-lightgrey" alt="License">
  </p>

  <p>
    <a href="#✨-Features">Features</a> •
    <a href="#🚀-Quick-Start">Install</a> •
    <a href="#⚙️-Configuration">Config</a> •
    <a href="#⌨️-Shortcuts">Shortcuts</a> •
    <a href="README.md">中文文档</a>
  </p>
</div>

> [!WARNING]
> **Development Phase Note**: This project is currently in early active development (Alpha/Beta) and may contain minor bugs or unstable features. If you encounter any issues, please submit an Issue, and we will fix it ASAP!

<br>

## 📖 Introduction

**Q-Solver** is a desktop AI assistant designed for **high-pressure, high-focus scenarios**. It deeply integrates top-tier LLMs like **OpenAI**, **Google Gemini**, and **Anthropic Claude**. With minimal screenshot interactions, it provides real-time code analysis, Q&A, and content creation assistance.

When you need to **stay focused while getting instant help**, Q-Solver is your best choice—floating window that doesn't break your flow, stealth mode that leaves no trace, and context awareness that truly understands your needs.

---

## ✨ Features

### 🛡️ Stealth Mode
Designed for scenarios where you **need to focus without interruption**:
- **Borderless / Semi-Transparent**: Always on top, never blocking your workspace
- **Anti-Focus Stealing**: Window stays visible when switching apps, answers always accessible
- **Stealth / Anti-Detection**: Special window attributes to evade screen recording and screenshot tools
- **Click-Through**: Click through the window to interact with apps behind it, zero interference

### 🔌 Comprehensive Model Support
- **Native SDK Integration**: Built-in native SDKs for Google Gemini, Anthropic Claude, and OpenAI (not just HTTP forwarding) to ensure optimal streaming response and stability.
- **Custom Access**: Supports aggregation services like OneAPI for flexible network adaptation.

### 📄 Context Awareness
- **Background Import**: Import PDF/Markdown documents, and the AI generates personalized answers based on your background.
- **Smart Memory**: Supports multi-turn conversation context, automatically clearing when a new topic starts.

### 🎙️ Gemini Live API (Experimental)
- **Real-time Voice Conversation**: Supports bidirectional real-time voice calls with Gemini models, providing an ultra-low latency interaction experience.
- **Voice Transcription Mode**: Captures and transcribes the other party's voice in real-time, with AI providing instant answer suggestions to help you handle any conversation scenario.
- **Note**: This feature is currently in an **experimental stage**. Connection instability may occur due to network conditions.
- **Supported Models**: Please select `gemini-2.0-flash-exp`.

### ⚡ Extreme Performance
- **Go + Wails**: Modern tech stack with native-level performance and extremely low memory footprint.
- **Vue 3 Frontend**: Responsive, modern UI that is smooth and fluid.
- **Global Shortcuts**: Wake up with `F8` in a split second without switching windows.

---

## 🖼️ Preview

<div align="center">
  <img src="assets/demo.gif" style="border-radius: 8px; width: 100%; box-shadow: 0 4px 20px rgba(0,0,0,0.1);" />
  <p><i>👆 Live Demo: Screenshot -> Thinking -> Answer</i></p>
</div>

| | |
|:---:|:---:|
| <img src="assets/img1.png" style="border-radius: 8px; width: 100%;" /> | <img src="assets/img2.png" style="border-radius: 8px; width: 100%;" /> |
| <img src="assets/img3.png" style="border-radius: 8px; width: 100%;" /> | <img src="assets/img4.png" style="border-radius: 8px; width: 100%;" /> |

---

## 🚀 Quick Start

### Method 1: Download Installer (Recommended)
Go to the [Releases](https://github.com/jym66/Q-Solver/releases) page to download the latest Windows installer (`.exe`).

### Method 2: Build from Source
If you are a developer, you can clone the source code for secondary development:

```bash
# Prerequisites: Go 1.25+, Node.js 22+, Wails CLI

# 1. Clone repository
git clone https://github.com/jym66/Q-Solver.git
cd Q-Solver

# 2. Run in Dev Mode (Hot Reload)
wails dev

# 3. Build Production Version
wails build -ldflags "-s -w" -tags prod
```

---

## 🍎 macOS Instructions

> ⚠️ **Compatibility Note**: The macOS version is currently in **compatibility adaptation stage** and may have bugs or instability. Please submit an Issue if you encounter any problems!

### Shortcuts
The macOS version supports the following fixed shortcuts. **Customization is not supported yet**:

| Function | Shortcut |
|----------|----------|
| Screenshot & Solve | `⌘1` |
| Show/Hide | `⌘2` |
| Click-Through | `⌘3` |
| Move Window Up | `⌘⌥↑` |
| Move Window Down | `⌘⌥↓` |
| Move Window Left | `⌘⌥←` |
| Move Window Right | `⌘⌥→` |
| Scroll Up | `⌘⌥⇧↑` |
| Scroll Down | `⌘⌥⇧↓` |

> **Note**: Custom shortcuts are currently only available on Windows.

### Screenshot Permission
On first launch, you need to grant screenshot permission:
1. Open **Settings** -> **Screenshot** tab
2. Click **Grant Screenshot Permission** button
3. Enable the app in System Preferences
4. Return to the app and click **Refresh Permission Status**

### System Audio Capture (Gemini Live API)
macOS does not natively support recording system audio. You need to install a virtual audio driver:

1. **Install BlackHole** (2-channel version recommended):
   ```bash
   brew install blackhole-2ch
   ```

2. **Configure Multi-Output Device**:
   - Open **Audio MIDI Setup** (Applications > Utilities)
   - Click the **+** button at the bottom left, select **Create Multi-Output Device**
   - Check both **BlackHole 2ch** and your speakers/headphones
   - Right-click the new multi-output device and select **Use This Device For Sound Output**

3. **Adjust BlackHole Volume** (Important!):
   - In **Audio MIDI Setup**, select **BlackHole 2ch** on the left panel
   - In the **Input** tab on the right, drag the **Main** volume slider to **1.0** (far right)
   - If the volume is not set to 1.0, the captured system audio will be very quiet

   <img src="assets/img5.png" width="600" style="border-radius: 8px; margin: 10px 0;" />

4. **Grant Microphone Permission**: When launching the Live API feature for the first time, please allow microphone access.

> **Note**: If BlackHole is not installed, the "Interviewer Mode" in Live API will not be able to capture system audio.

---

## ⚙️ Configuration

Launch the software and click the **Settings** icon in the top right corner:

1.  **Select Provider**: Supports OpenAI, Gemini, Claude, or Custom.
2.  **API Key**: Enter the API Key for the corresponding platform.
3.  **Model Parameters**:
    *   **Temperature**: Control answer randomness (0.0 - 2.0).
    *   **Thinking Budget**: Set Token budget for o1/Claude-3.5 thinking process.
4.  **Custom Mode**:
    *   Select `Custom` provider.
    *   Enter API URL (e.g., OneAPI).
    *   Supports automatic protocol switching based on model prefix (Gemini/Claude).

---

## ⌨️ Shortcuts

### Windows Shortcuts

Windows version supports **customizable** shortcuts. Default configuration:

| Key | Function |
| :--- | :--- |
| **F8** | **Screenshot & Ask** (Core Function) |
| **F9** | Show / Hide Main Window |
| **F10** | Toggle Mouse Click-Through |
| **Alt + ⬆️/⬇️/⬅️/➡️** | Fine-tune Window Position |
| **Alt + PgUp/PgDn** | Quick History Navigation |

### macOS Shortcuts

macOS version uses **fixed** shortcuts. Customization is not supported yet:

| Key | Function |
| :--- | :--- |
| **⌘1** | **Screenshot & Solve** (Core Function) |
| **⌘2** | Show / Hide Main Window |
| **⌘3** | Toggle Mouse Click-Through |
| **⌘⌥↑/↓/←/→** | Fine-tune Window Position |
| **⌘⌥⇧↑/↓** | Quick History Navigation |

## Support

If you like this project, please give it a ⭐️ **Star**, it helps us a lot!

**Disclaimer**: This project is for technical research and personal learning only. Please do not use it for any illegal purposes.

<br>

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/jym66">jym66</a></p>
</div>
