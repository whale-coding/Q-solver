package platform

// WindowHandle 窗口句柄类型（Windows 为 HWND，macOS 为 NSWindow 指针）
type WindowHandle uintptr

// Platform 平台相关操作接口
type Platform interface {
	// 获取当前进程主窗口句柄
	GetWindowHandle() (WindowHandle, error)

	// 应用幽灵模式（无边框、置顶、防录屏、不抢焦点）
	ApplyGhostMode(hwnd WindowHandle) error

	// 设置鼠标穿透
	SetClickThrough(hwnd WindowHandle, enabled bool) error

	// 设置防录屏状态
	SetDisplayAffinity(hwnd WindowHandle, hidden bool) error

	// 恢复焦点
	RestoreFocus(hwnd WindowHandle) error

	// 移除焦点
	RemoveFocus(hwnd WindowHandle) error

	// 检查截图权限 (macOS 10.15+，Windows 直接返回 true)
	CheckScreenCaptureAccess() bool

	// 请求截图权限 (macOS 10.15+，Windows 直接返回 true)
	RequestScreenCaptureAccess() bool
}

// Current 当前平台实现（由条件编译决定）
var Current Platform
