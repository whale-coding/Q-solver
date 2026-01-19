<div align="center">
  <img src="assets/banner.jpg" alt="Q-Solver Banner" width="100%" style="border-radius: 16px; box-shadow: 0 10px 40px rgba(0,0,0,0.15);">

  <br>

  <h1>🧠 Q-Solver</h1>
  
  <h3>AI 驱动的实时桌面助手 · 截图解题 · 语音对话</h3>
  
  <p><i>🎯 一键截图开启深度思考，实时语音连接智能未来</i></p>

  <p>
    <a href="https://github.com/jym66/Q-solver/stargazers"><img src="https://img.shields.io/github/stars/jym66/Q-solver?color=ffcb6b&style=for-the-badge&labelColor=30363d" alt="Stars"></a>
    <a href="https://github.com/jym66/Q-solver/releases"><img src="https://img.shields.io/github/v/release/jym66/Q-solver?color=89d185&style=for-the-badge&labelColor=30363d" alt="Release"></a>
    <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white&labelColor=30363d" alt="Go">
    <img src="https://img.shields.io/badge/Vue-3.x-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white&labelColor=30363d" alt="Vue">
    <img src="https://img.shields.io/badge/Wails-v2-E30613?style=for-the-badge&logo=wails&logoColor=white&labelColor=30363d" alt="Wails">
  </p>
  
  <p>
    <img src="https://img.shields.io/badge/macOS-000000?style=flat-square&logo=apple&logoColor=white" alt="macOS">
    <img src="https://img.shields.io/badge/Windows-0078D6?style=flat-square&logo=windows&logoColor=white" alt="Windows">
  </p>

  <br>

  <p>
    <a href="#-核心特性">特性</a> •
    <a href="#-快速开始">安装</a> •
    <a href="#-功能演示">演示</a> •
    <a href="#%EF%B8%8F-配置指南">配置</a> •
    <a href="README_EN.md">English Documentation</a>
  </p>
  
  <br>
  
  <img src="assets/demo.gif" alt="Demo" width="92%" style="border-radius: 12px; box-shadow: 0 12px 40px rgba(0,0,0,0.25); border: 1px solid rgba(255,255,255,0.1);">

</div>

<br>
<br>

> [!CAUTION]
> **🚧 开发阶段警告**：本项目目前处于**早期开发预览阶段 (Pre-Alpha)**。功能可能会随版本更新发生重大变化，建议仅用于测试和尝鲜。

<br>

<div align="center">

## 🌟 为什么选择 Q-Solver？

</div>

<table>
<tr>
<td width="50%" valign="top">

### 🖼️ 极速截图求解
只需一个快捷键，即刻捕获屏幕内容并进行 AI 分析。
- **📸 智能识别**：精准识别文字、公式、代码。
- **🧠 深度思考**：支持 o1, Claude 3.5 等强推理模型。
- **⚡️ 零干扰**：悬浮窗设计，不打断当前工作流。

</td>
<td width="50%" valign="top">

### 🎙️ 沉浸式语音交互
集成了 Google Gemini Live API，体验丝滑的实时对话。
- **🗣️ 双向通话**：毫秒级响应，如同真人交谈。
- **🗺️ 思维导图**：对话内容自动整理为可视化导图。
- **📝 智能笔记**：自动转录并总结重点，支持导出。

</td>
</tr>
</table>

<br>

<div align="center">

## ✨ 核心特性

</div>

### 🛡️ 隐身模式 (Stealth Mode)

专为隐私与多任务设计，打造“幽灵”般的窗口体验。

| 特性 | 描述 |
|:---|:---|
| **🚫 防录屏检测** | 窗口对大多数录屏/截屏软件不可见（macOS 14+ 支持更好） |
| **👻 鼠标穿透** | 开启后可透过窗口点击后方内容，互不影响 |
| **📌 全局置顶** | 始终悬浮在其他窗口之上，重要信息一眼即达 |
| **🔕 沉浸免打扰** | 精心设计的焦点管理，输入时不抢占主窗口焦点 |

---

### 🧠 多模型生态

不局限于单一模型，根据需求灵活切换。

| 模型系列 | 推荐场景 | 实时语音 |
|:---|:---|:---:|
| **Google Gemini 2.0** | ⚡️ 极速响应，Live API 语音对话首选 | ✅ |
| **OpenAI GPT-4o / o1** | 🎓 复杂逻辑推理，数学/编程难题求解 | ❌ |
| **Claude 3.5 Sonnet** | 📝 长文本分析，代码编写与审查 | ❌ |
| **Custom OpenAI兼容** | 🛠️ 对接 OneAPI 或本地 LLM 服务 | ❌ |

---

### 🎨 个性化与上下文

- **📄 简历/背景植入**：上传 PDF 简历，AI 将基于你的背景提供定制化建议（面试辅导、代码解释）。
- **💾 长时记忆**：智能维护对话上下文，拒绝“金鱼记忆”。

<br>

<div align="center">

## 📸 界面展示

</div>

| | | |
|:---:|:---:|:---:|
| <img src="assets/img1.png" width="100%" style="border-radius: 8px;"/> | <img src="assets/img6.png" width="100%" style="border-radius: 8px;"/> | <img src="assets/img7.png" width="100%" style="border-radius: 8px;"/> |

<br>
<br>

## 🚀 快速开始

### 📥 方式一：直接下载 (如果你想直接使用)

前往 [Releases 页面](https://github.com/jym66/Q-solver/releases) 下载对应系统的最新安装包。

> [!NOTE]
> **macOS 用户提示**：首次运行时如果提示“已损坏”或无法打开，请执行以下命令：
> ```bash
> xattr -cr /Applications/Q-Solver.app
> chmod +x /Applications/Q-Solver.app/Contents/MacOS/Q-Solver
> ```

### 🛠️ 方式二：源码构建 (如果你是开发者)

**环境要求**：Go 1.25+, Node.js 22+, Wails CLI

```bash
# 1. 安装 Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 2. 克隆仓库
git clone https://github.com/jym66/Q-solver.git
cd Q-Solver

# 3. 开发模式运行 (支持热重载)
wails dev

# 4. 编译发布版本
wails build -ldflags "-s -w" -tags prod
```

<br>

## ⌨️ 快捷键指南

> 💡 **提示**：目前 macOS 快捷键固定，Windows 后续支持自定义。

| 动作 | Windows | macOS |
|:---|:---:|:---:|
| **截图并提问** 📸 | `F8` | `⌘ + 1` |
| **显示/隐藏窗口** 👁️ | `F9` | `⌘ + 2` |
| **切换鼠标穿透** 👻 | `F10` | `⌘ + 3` |
| **微调窗口位置** ↕️ | `Alt + 方向键` | `⌘⌥ + 方向键` |
| **快速翻页** 📜 | `Alt + PgUp/Dn` | `⌘⌥⇧ + ↑/↓` |

<br>

## ⚙️ 配置与使用

1. 点击窗口右上角的 **设置 (Settings)** 图标。
2. 在 **提供商 (Provider)** 中选择你已有的 API 服务 (如 Gemini)。
3. 填入你的 **API Key**。
4. (可选) 开启 **Live API** 体验实时语音。

### 🍎 macOS 特别配置

macOS 需要额外权限以发挥完整功能：

<details>
<summary><b>🔐 屏幕录制权限 (必选)</b></summary>

为了实现截图功能，首次使用时：
1. 系统会弹窗提示请求 **屏幕录制** 权限。
2. 若未弹窗，请前往 **系统设置** -> **隐私与安全性** -> **屏幕录制**。
3. 勾选 **Q-Solver**。
4. **重启应用** 生效。

</details>

<details>
<summary><b>🎙️ 系统音频内录 (Live API 必选)</b></summary>

若想让 AI 听到电脑播放的声音（如会议内容），需要安装虚拟声卡：

1. 安装 [BlackHole](https://github.com/ExistentialAudio/BlackHole):
   ```bash
   brew install blackhole-2ch
   ```
2. 打开 **音频 MIDI 设置 (Audio MIDI Setup)**。
3. 创建 **多输出设备 (Multi-Output Device)**，同时勾选 **扬声器** 和 **BlackHole 2ch**。
4. 将该多输出设备设为系统默认输出。
5. 在 Q-Solver 设置中，确保音频输入包含 BlackHole。

<img src="assets/img5.png" width="90%" style="border-radius: 8px;"/>

</details>

<br>

## 🛠️ 技术栈概览

- **Core**: [Go](https://go.dev/) (Logic) + [Wails](https://wails.io/) (Binding)
- **UI**: [Vue 3](https://vuejs.org/) + [Vue Flow](https://vueflow.dev/) (Mind Map)
- **AI**: Gemini Protocol, OpenAI SDK
- **Audio**: Miniaudio (via malgo), BlackHole

<br>

<br>

## 📈 Star 趋势

<div align="center">
  <a href="https://star-history.com/#jym66/Q-solver&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=jym66/Q-solver&type=Date&theme=dark" />
      <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=jym66/Q-solver&type=Date" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=jym66/Q-solver&type=Date" />
    </picture>
  </a>
</div>

<br>

## 📄 许可证

本项目基于 **CC BY-NC 4.0** 协议开源，仅供 **非商业个人学习与研究** 使用。

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/jym66">jym66</a></p>
  <p>
    如果你觉得这个项目有趣，欢迎点个 <b>⭐ Star</b> 支持一下！
  </p>
</div>
