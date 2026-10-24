<script setup lang="ts">
/** 版本说明（更新日志）弹窗，风格参考 shushu.fan/changelog */
defineProps<{ open: boolean }>()
const emit = defineEmits<{ (e: 'close'): void }>()

interface Entry {
  date: string
  type: '新增' | '优化' | '修复'
  title: string
  desc: string
}

const changelog: Entry[] = [
  { date: '2026-10-24', type: '新增', title: '改枪码粘贴识别', desc: '粘贴分享码自动识别枪械、模式与标题，无需手动从分类里选择。' },
  { date: '2026-10-24', type: '新增', title: '导入 / 导出备份', desc: '一键把全部改枪方案导出为 JSON 备份，可整体导入恢复；入口收进标题栏 ⋮ 菜单。' },
  { date: '2026-10-24', type: '新增', title: '价格标签与快捷筛选', desc: '价格标签单选、标签 / 类型 / 模式快捷筛选、排序与仅看收藏。' },
  { date: '2026-10-24', type: '新增', title: '标题栏菜单与「关于」', desc: '无边框标题栏加 ⋮ 菜单，集中导出、导入、关于与版本说明；关于展示作者书包呀、应用与数据版本。' },
  { date: '2026-10-24', type: '优化', title: '卡片交互', desc: '点击卡片即复制改枪码；收藏角标、编辑 / 删除图标收进卡片右上角；hover 高亮提示可复制。' },
  { date: '2026-10-24', type: '优化', title: '界面风格', desc: '采用三角洲风格的暖军绿灰配色、琥珀橙强调与白底装备图，正文微软雅黑 + Barlow Condensed 数字。' },
  { date: '2026-10-24', type: '修复', title: '烽火地带双前缀', desc: '修复粘贴改枪码被重复补「武器名-烽火地带」前缀的问题。' },
  { date: '2026-10-24', type: '修复', title: '武器数据核对', desc: '68 把武器与官方全表对齐：补录缺失枪械、更正分类、补齐官方图样。' },
]

const typeStyle: Record<Entry['type'], string> = {
  新增: 'border-emerald-500/50 bg-emerald-500/10 text-emerald-400',
  优化: 'border-amber-500/50 bg-amber-500/10 text-amber-400',
  修复: 'border-rose-500/50 bg-rose-500/10 text-rose-400',
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    @click="emit('close')"
  >
    <div
      class="flex max-h-[80vh] w-[26rem] flex-col overflow-hidden rounded-lg border border-surface-border bg-surface shadow-xl"
      @click.stop
    >
      <div class="flex items-start justify-between border-b border-surface-border px-5 py-3.5">
        <div>
          <h2 class="text-base font-bold text-slate-100">版本说明</h2>
          <p class="mt-0.5 text-xs text-slate-500">公示各版本的更新、改动，以及各种缝缝补补。</p>
        </div>
        <button
          class="-mt-1 text-lg text-slate-500 transition hover:text-slate-200"
          title="关闭"
          @click="emit('close')"
        >✕</button>
      </div>

      <div class="flex-1 overflow-y-auto px-5 py-3">
        <div
          v-for="(e, i) in changelog"
          :key="i"
          class="flex gap-3 py-2.5"
        >
          <div class="w-14 shrink-0 pt-0.5 text-right text-[11px] text-slate-500">{{ e.date }}</div>
          <div class="flex-1 border-l border-surface-border pl-3">
            <div class="flex flex-wrap items-center gap-2">
              <span
                class="rounded border px-1.5 py-0.5 text-[10px]"
                :class="typeStyle[e.type]"
              >{{ e.type }}</span>
              <span class="text-sm font-medium text-slate-100">{{ e.title }}</span>
            </div>
            <p class="mt-1 text-xs leading-relaxed text-slate-400">{{ e.desc }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
