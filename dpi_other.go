//go:build !windows

package main

// displayScale 非 Windows 平台不做缩放换算。
func displayScale() float64 {
	return 1.0
}

// scaleSize 按缩放比例换算窗口尺寸。
func scaleSize(css int) int {
	return css
}
