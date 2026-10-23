/** Tailwind CSS v3 配置 —— 暗色主题的卡片式界面基础
 *  注意：使用 v3，勿升级到 v4（配置格式不同）
 */
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // 面板底色（深炭黑，贴合《三角洲行动》黑底橙强调的界面）
        surface: {
          DEFAULT: '#17181b',
          raised: '#202226',
          sunken: '#0e0f11',
          border: '#2c2f35',
        },
        // 强调色：三角洲琥珀橙
        accent: {
          DEFAULT: '#ff9f1c',
          soft: '#2b1f0e',
        },
      },
      fontFamily: {
        // 改枪码必须等宽显示
        mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Consolas', 'monospace'],
      },
    },
  },
  plugins: [],
}
