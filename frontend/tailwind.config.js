/** Tailwind CSS v3 配置 —— 暗色主题的卡片式界面基础
 *  注意：使用 v3，勿升级到 v4（配置格式不同）
 */
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // 面板底色（深军绿暖调，贴合三角洲"军事情报板"质感，非冷黑科技风）
        surface: {
          DEFAULT: '#1b201a',
          raised: '#242a23',
          sunken: '#141815',
          border: '#394033',
        },
        // 强调色：三角洲琥珀橙（更暖）
        accent: {
          DEFAULT: '#f5a81a',
          soft: '#2a2010',
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
