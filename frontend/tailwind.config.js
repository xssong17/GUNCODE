/** Tailwind CSS v3 配置 —— 暗色主题的卡片式界面基础
 *  注意：使用 v3，勿升级到 v4（配置格式不同）
 */
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // 面板底色（浅沙色/米白暖调，贴合三角洲主菜单"军事情报站"的明亮风格）
        surface: {
          DEFAULT: '#e8e4da',
          raised: '#f4f1ea',
          sunken: '#ddd8cb',
          border: '#c6bfad',
        },
        // 强调色：三角洲琥珀橙（浅底上用更深的橙更醒目）
        accent: {
          DEFAULT: '#e07f10',
          soft: '#f6ead2',
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
