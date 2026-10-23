//go:build windows

package main

import "syscall"

var procGetDpiForSystem = syscall.NewLazyDLL("user32.dll").NewProc("GetDpiForSystem")

// displayScale 返回系统显示缩放比例（150% 缩放返回 1.5）。
//
// Wails v2 建窗时用的是物理像素，而 WebView2 按 CSS 像素布局，
// 在缩放显示器上会导致内容被裁掉右侧与底部，因此需要按缩放比例换算窗口尺寸。
func displayScale() float64 {
	if err := procGetDpiForSystem.Find(); err == nil {
		dpi, _, _ := procGetDpiForSystem.Call()
		if dpi > 0 {
			return float64(dpi) / 96.0
		}
	}
	return 1.0
}

// scaleSize 按缩放比例换算窗口尺寸。
func scaleSize(css int) int {
	v := int(float64(css)*displayScale() + 0.5)
	if v < 1 {
		return css
	}
	return v
}
