//go:build windows

package platform

import (
	"Q-Solver/pkg/winapi"
	"os"
)

func init() {
	Current = &windowsPlatform{}
}

type windowsPlatform struct{}

func (p *windowsPlatform) GetWindowHandle() (WindowHandle, error) {
	hwnd, err := winapi.GetHwndByPid(uint32(os.Getpid()))
	return WindowHandle(hwnd), err
}

func (p *windowsPlatform) ApplyGhostMode(hwnd WindowHandle) error {
	winapi.ApplyGhostMode(uintptr(hwnd))
	return nil
}

func (p *windowsPlatform) SetClickThrough(hwnd WindowHandle, enabled bool) error {
	return winapi.SetWindowClickThrough(uintptr(hwnd), enabled)
}

func (p *windowsPlatform) SetDisplayAffinity(hwnd WindowHandle, hidden bool) error {
	affinity := winapi.WDA_NONE
	if hidden {
		affinity = winapi.WDA_EXCLUDEFROMCAPTURE
	}
	return winapi.SetWindowDisplayAffinity(uintptr(hwnd), uint32(affinity))
}

func (p *windowsPlatform) RestoreFocus(hwnd WindowHandle) error {
	winapi.RestoreFocus(uintptr(hwnd))
	return nil
}

func (p *windowsPlatform) RemoveFocus(hwnd WindowHandle) error {
	winapi.RemoveFocus(uintptr(hwnd))
	return nil
}

func (p *windowsPlatform) CheckScreenCaptureAccess() bool {
	return true // Windows 不需要截图权限
}

func (p *windowsPlatform) RequestScreenCaptureAccess() bool {
	return true // Windows 不需要截图权限
}

func (p *windowsPlatform) OpenScreenCaptureSettings() {
	// Windows 不需要截图权限，无操作
}

func (p *windowsPlatform) SetWindowLevel(hwnd WindowHandle, level int) error {
	// Windows 使用不同的置顶机制，在 ApplyGhostMode 中处理
	return nil
}
