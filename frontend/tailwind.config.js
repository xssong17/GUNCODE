/** Tailwind CSS v3 配置 —— 暗色主题的卡片式界面基础
 *  注意：使用 v3，勿升级到 v4（配置格式不同）
 */
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // 面板底色（深色主题）
        surface: {
          DEFAULT: '#1b1e24',
          raised: '#242830',
          sunken: '#15171c',
          border: '#333844',
        },
        accent: {
          DEFAULT: '#f0a020',
          soft: '#3a2c14',
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
