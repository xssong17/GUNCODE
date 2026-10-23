<script lang="ts" setup>
import HomeView from './views/HomeView.vue'
</script>

<template>
  <div class="flex h-full flex-col bg-surface-sunken">
    <!-- 顶栏骨架（S4 阶段接入搜索框与操作按钮） -->
    <header class="flex h-14 shrink-0 items-center gap-3 border-b border-surface-border bg-surface px-4">
      <h1 class="text-base font-semibold text-slate-100">GunCode 改枪码管理器</h1>
      <span class="rounded bg-surface-raised px-2 py-0.5 text-xs text-slate-400">三角洲行动</span>
    </header>

    <main class="min-h-0 flex-1">
      <HomeView />
    </main>
  </div>
</template>
