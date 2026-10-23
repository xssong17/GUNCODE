<script lang="ts" setup>
import ToastHost from './components/ToastHost.vue'
import HomeView from './views/HomeView.vue'
</script>

<template>
  <div class="flex h-full flex-col bg-surface-sunken">
    <!-- 顶栏 -->
    <header class="flex h-12 shrink-0 items-center gap-3 border-b border-surface-border bg-surface px-4">
      <h1 class="text-sm font-semibold tracking-wide text-slate-100">GunCode 改枪码管理器</h1>
      <span class="rounded bg-surface-raised px-1.5 py-0.5 text-[10px] text-slate-500">三角洲行动</span>
    </header>

    <main class="min-h-0 flex-1">
      <HomeView />
    </main>

    <ToastHost />
  </div>
</template>
