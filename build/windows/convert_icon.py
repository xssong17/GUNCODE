# -*- coding: utf-8 -*-
"""把最终选定的「鼠鼠看看配枪」图标转换为多尺寸 .ico，覆盖 Wails exe 图标。"""
from PIL import Image

src = r'E:\ai_yard\gunCode\build\windows\candidates\composite11_transparent.png'
dst = r'E:\ai_yard\gunCode\build\windows\icon.ico'

img = Image.open(src).convert('RGBA')
print('orig', img.size, img.mode)
# 圆角方形已是图标构图，缩到 256 以内
img.thumbnail((256, 256))
img.save(
    dst,
    format='ICO',
    sizes=[(256, 256), (128, 128), (64, 64), (48, 48), (32, 32), (24, 24), (16, 16)],
)
print('saved ico', img.size)
