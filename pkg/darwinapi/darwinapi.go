//go:build darwin

package darwinapi

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework AppKit -framework CoreGraphics

#import <Cocoa/Cocoa.h>
#import <CoreGraphics/CoreGraphics.h>

// 检查截图权限 (macOS 10.15+)
bool CheckScreenCaptureAccess() {
    if (@available(macOS 10.15, *)) {
        return CGPreflightScreenCaptureAccess();
    }
    return true; // macOS 10.15 以下不需要权限
}

// 请求截图权限 (macOS 10.15+)
bool RequestScreenCaptureAccess() {
    if (@available(macOS 10.15, *)) {
        return CGRequestScreenCaptureAccess();
    }
    return true; // macOS 10.15 以下不需要权限
}

// 获取应用主窗口
void* GetMainWindow() {
    NSApplication* app = [NSApplication sharedApplication];
    NSWindow* window = [app mainWindow];
    if (window == nil && app.windows.count > 0) {
        window = app.windows[0];
    }
    return (__bridge void*)window;
}

// 设置窗口忽略鼠标事件（鼠标穿透）
void SetWindowIgnoresMouseEvents(void* nsWindow, bool ignores) {
    if (nsWindow == NULL) return;
    NSWindow* window = (__bridge NSWindow*)nsWindow;
    dispatch_async(dispatch_get_main_queue(), ^{
        [window setIgnoresMouseEvents:ignores];
    });
}

// 设置窗口级别（置顶）
void SetWindowLevel(void* nsWindow, int level) {
    if (nsWindow == NULL) return;
    NSWindow* window = (__bridge NSWindow*)nsWindow;
    dispatch_async(dispatch_get_main_queue(), ^{
        [window setLevel:level];
    });
}

// 设置窗口样式（无边框 + 圆角）
void SetWindowStyleMaskBorderless(void* nsWindow) {
    if (nsWindow == NULL) return;
    NSWindow* window = (__bridge NSWindow*)nsWindow;
    dispatch_async(dispatch_get_main_queue(), ^{
        // 无边框样式
        [window setStyleMask:NSWindowStyleMaskBorderless];
        // 透明背景
        [window setBackgroundColor:[NSColor clearColor]];
        [window setOpaque:NO];
        // 允许透明
        [window setHasShadow:YES];
    });
}

// 设置窗口 sharingType（防录屏，macOS 14+）
void SetWindowSharingType(void* nsWindow, int sharingType) {
    if (nsWindow == NULL) return;
    NSWindow* window = (__bridge NSWindow*)nsWindow;
    dispatch_async(dispatch_get_main_queue(), ^{
        if (@available(macOS 14.0, *)) {
            window.sharingType = sharingType;
        }
    });
}

// 设置窗口不激活
void SetWindowNotActivating(void* nsWindow, bool noActivate) {
    if (nsWindow == NULL) return;
    NSWindow* window = (__bridge NSWindow*)nsWindow;
    dispatch_async(dispatch_get_main_queue(), ^{
        if (noActivate) {
            [window setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces |
                                           NSWindowCollectionBehaviorStationary |
                                           NSWindowCollectionBehaviorIgnoresCycle];
        } else {
            [window setCollectionBehavior:NSWindowCollectionBehaviorDefault];
        }
    });
}

// 让窗口可以接收键盘焦点
void SetWindowCanBecomeKey(void* nsWindow) {
    if (nsWindow == NULL) return;
    NSWindow* window = (__bridge NSWindow*)nsWindow;
    dispatch_async(dispatch_get_main_queue(), ^{
        [window makeKeyAndOrderFront:nil];
    });
}

// 设置 contentView 圆角
void SetWindowCornerRadius(void* nsWindow, float radius) {
    if (nsWindow == NULL) return;
    NSWindow* window = (__bridge NSWindow*)nsWindow;
    dispatch_async(dispatch_get_main_queue(), ^{
        // 设置 contentView 层的圆角
        NSView* contentView = [window contentView];
        if (contentView) {
            [contentView setWantsLayer:YES];
            contentView.layer.cornerRadius = radius;
            contentView.layer.masksToBounds = YES;
        }
    });
}

// 设置窗口透明度 (0.0 - 1.0)
void SetWindowAlpha(void* nsWindow, float alpha) {
    if (nsWindow == NULL) return;
    NSWindow* window = (__bridge NSWindow*)nsWindow;
    dispatch_async(dispatch_get_main_queue(), ^{
        [window setAlphaValue:alpha];
    });
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	NSWindowSharingNone     = 0
	NSWindowSharingReadOnly = 1

	// NSFloatingWindowLevel 相当于 CGWindowLevelForKey(kCGFloatingWindowLevelKey)
	NSFloatingWindowLevel = 3
	// NSScreenSaverWindowLevel - 更高级别
	NSScreenSaverWindowLevel = 1000
)

// GetMainWindow 获取应用主窗口
func GetMainWindow() (unsafe.Pointer, error) {
	window := C.GetMainWindow()
	if window == nil {
		return nil, errors.New("无法获取主窗口")
	}
	return window, nil
}

// SetWindowIgnoresMouseEvents 设置鼠标穿透
func SetWindowIgnoresMouseEvents(window unsafe.Pointer, ignores bool) {
	C.SetWindowIgnoresMouseEvents(window, C.bool(ignores))
}

// SetWindowLevel 设置窗口层级
func SetWindowLevel(window unsafe.Pointer, level int) {
	C.SetWindowLevel(window, C.int(level))
}

// SetWindowStyleMaskBorderless 设置无边框样式
func SetWindowStyleMaskBorderless(window unsafe.Pointer) {
	C.SetWindowStyleMaskBorderless(window)
}

// SetWindowSharingType 设置防录屏
func SetWindowSharingType(window unsafe.Pointer, sharingType int) {
	C.SetWindowSharingType(window, C.int(sharingType))
}

// SetWindowNotActivating 设置不抢焦点
func SetWindowNotActivating(window unsafe.Pointer, noActivate bool) {
	C.SetWindowNotActivating(window, C.bool(noActivate))
}

// SetWindowCanBecomeKey 允许窗口获取键盘焦点
func SetWindowCanBecomeKey(window unsafe.Pointer) {
	C.SetWindowCanBecomeKey(window)
}

// SetWindowCornerRadius 设置窗口圆角
func SetWindowCornerRadius(window unsafe.Pointer, radius float32) {
	C.SetWindowCornerRadius(window, C.float(radius))
}

// SetWindowAlpha 设置窗口透明度
func SetWindowAlpha(window unsafe.Pointer, alpha float32) {
	C.SetWindowAlpha(window, C.float(alpha))
}

// CheckScreenCaptureAccess 检查是否有截图权限 (macOS 10.15+)
func CheckScreenCaptureAccess() bool {
	return bool(C.CheckScreenCaptureAccess())
}

// RequestScreenCaptureAccess 请求截图权限 (macOS 10.15+)
// 会弹出系统权限对话框，用户需要在系统偏好设置中授权
func RequestScreenCaptureAccess() bool {
	return bool(C.RequestScreenCaptureAccess())
}
