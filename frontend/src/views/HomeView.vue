<script lang="ts" setup>
// S0 骨架占位：S4 阶段在此搭建「左侧武器栏 + 筛选条 + 卡片网格」三栏布局
</script>

<template>
  <div class="flex h-full">
    <!-- 左侧：枪械分类与武器列表 -->
    <aside class="w-60 shrink-0 overflow-y-auto border-r border-surface-border bg-surface p-3">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-500">枪械分类</div>
      <div class="mt-3 space-y-1">
        <div
          v-for="c in ['突击步枪', '射手步枪', '狙击步枪', '轻机枪', '冲锋枪', '霰弹枪', '手枪', '特殊武器']"
          :key="c"
          class="rounded-md px-2 py-1.5 text-sm text-slate-400 transition-colors hover:bg-surface-raised hover:text-slate-200"
        >
          {{ c }}
        </div>
      </div>
    </aside>

    <!-- 右侧：卡片网格 -->
    <section class="min-w-0 flex-1 overflow-y-auto p-4">
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
        <div class="card">
          <div class="flex items-center justify-between">
            <span class="text-sm font-semibold text-slate-200">M4A1</span>
            <span class="text-xs text-accent">★</span>
          </div>
          <p class="mt-1 text-xs text-slate-500">中距离架枪（示例卡片）</p>
          <div class="code-text mt-3">AKM突击步枪-烽火地带-6KJ5IFS0A021USGQ8AC2S</div>
          <div class="mt-3 flex gap-2">
            <button class="rounded-md bg-accent px-3 py-1 text-xs font-medium text-slate-900 hover:brightness-110">
              复制
            </button>
            <button class="rounded-md border border-surface-border px-3 py-1 text-xs text-slate-300 hover:bg-surface-sunken">
              编辑
            </button>
          </div>
        </div>
      </div>

      <p class="mt-6 text-center text-xs text-slate-600">
        S0 骨架：Tailwind v3 + Vue3 主题已就绪，真实数据与交互在 S2~S6 接入
      </p>
    </section>
  </div>
</template>
