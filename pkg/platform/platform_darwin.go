//go:build darwin

package platform

import (
	"Q-Solver/pkg/darwinapi"
	"Q-Solver/pkg/logger"
	"unsafe"
)

func init() {
	Current = &darwinPlatform{}
}

type darwinPlatform struct {
	mainWindow unsafe.Pointer
}

func (p *darwinPlatform) GetWindowHandle() (WindowHandle, error) {
	window, err := darwinapi.GetMainWindow()
	if err != nil {
		return 0, err
	}
	p.mainWindow = window
	return WindowHandle(uintptr(window)), nil
}

func (p *darwinPlatform) ApplyGhostMode(hwnd WindowHandle) error {
	window := unsafe.Pointer(uintptr(hwnd))

	// 无边框 + 透明背景
	darwinapi.SetWindowStyleMaskBorderless(window)

	// 设置圆角
	darwinapi.SetWindowCornerRadius(window, 12.0)

	// 置顶
	darwinapi.SetWindowLevel(window, darwinapi.NSFloatingWindowLevel)

	// 防录屏macos14+才可以
	darwinapi.SetWindowSharingType(window, darwinapi.NSWindowSharingNone)

	// 不抢焦点
	darwinapi.SetWindowNotActivating(window, true)

	logger.Println("[macOS] 幽灵模式已激活")
	return nil
}

func (p *darwinPlatform) SetClickThrough(hwnd WindowHandle, enabled bool) error {
	window := unsafe.Pointer(uintptr(hwnd))
	darwinapi.SetWindowIgnoresMouseEvents(window, enabled)
	return nil
}

func (p *darwinPlatform) SetDisplayAffinity(hwnd WindowHandle, hidden bool) error {
	window := unsafe.Pointer(uintptr(hwnd))
	if hidden {
		darwinapi.SetWindowSharingType(window, darwinapi.NSWindowSharingNone)
	} else {
		darwinapi.SetWindowSharingType(window, darwinapi.NSWindowSharingReadOnly)
	}
	return nil
}

func (p *darwinPlatform) RestoreFocus(hwnd WindowHandle) error {
	window := unsafe.Pointer(uintptr(hwnd))
	darwinapi.SetWindowCanBecomeKey(window)
	return nil
}

func (p *darwinPlatform) RemoveFocus(hwnd WindowHandle) error {
	window := unsafe.Pointer(uintptr(hwnd))
	darwinapi.SetWindowNotActivating(window, true)
	return nil
}

func (p *darwinPlatform) CheckScreenCaptureAccess() bool {
	return darwinapi.CheckScreenCaptureAccess()
}

func (p *darwinPlatform) RequestScreenCaptureAccess() bool {
	return darwinapi.RequestScreenCaptureAccess()
}

func (p *darwinPlatform) OpenScreenCaptureSettings() {
	darwinapi.OpenScreenCaptureSettings()
}

func (p *darwinPlatform) SetWindowLevel(hwnd WindowHandle, level int) error {
	window := unsafe.Pointer(uintptr(hwnd))
	darwinapi.SetWindowLevel(window, level)
	return nil
}
