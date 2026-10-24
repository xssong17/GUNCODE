# -*- coding: utf-8 -*-
"""优化版图标：大尺寸清晰缩放；小尺寸(<=48px)做高对比+边缘锐化，提升任务栏辨识度。"""
from PIL import Image, ImageEnhance, ImageFilter

src = r'E:\ai_yard\gunCode\build\windows\candidates\composite10.png'
dst = r'E:\ai_yard\gunCode\build\windows\icon.ico'

SIZES = [256, 128, 64, 48, 32, 24, 16]

base = Image.open(src).convert('RGBA')

def build_frame(s):
    img = base.resize((s, s), Image.LANCZOS)
    if s <= 48:
        # 小尺寸：增强对比 + 边缘锐化，让轮廓更醒目
        img = ImageEnhance.Contrast(img).enhance(1.55)
        img = img.filter(ImageFilter.UnsharpMask(radius=1, percent=220, threshold=0))
        img = ImageEnhance.Sharpness(img).enhance(1.4)
        # 轻微提高饱和度，让枪/猫颜色更分明
        img = ImageEnhance.Color(img).enhance(1.2)
    return img

frames = [build_frame(s) for s in SIZES]
# 第一帧带 sizes 保存，其余作为 append_images 追加
frames[0].save(dst, format='ICO', sizes=[(s, s) for s in SIZES],
               append_images=frames[1:])
print('saved optimized ico, sizes=', SIZES)
