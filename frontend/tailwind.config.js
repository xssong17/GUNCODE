/** Tailwind CSS v3 配置 —— 暗色主题的卡片式界面基础
 *  注意：使用 v3，勿升级到 v4（配置格式不同）
 */
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // 面板底色：中等深度的暖军绿灰（非死黑、非军绿深、非浅白），贴合三角洲军事信息板
        surface: {
          DEFAULT: '#282d29',
          raised: '#303630',
          sunken: '#212522',
          border: '#424a42',
        },
        // 强调色：三角洲琥珀橙，只用于选中/主按钮/重要数据
        accent: {
          DEFAULT: '#f0a31a',
          soft: '#332713',
        },
      },
      fontFamily: {
        // 正文/中文：微软雅黑，清晰耐读
        sans: ['Microsoft YaHei UI', 'Microsoft YaHei', 'PingFang SC', 'system-ui', 'sans-serif'],
        // 标题/数字/英文：BarlowCondensed 战术窄体（只做强调）
        display: ['BarlowCondensed', 'Microsoft YaHei UI', 'Microsoft YaHei', 'sans-serif'],
        // 改枪码必须等宽显示
        mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Consolas', 'monospace'],
      },
    },
  },
  plugins: [],
}
