<script lang="ts" setup>
import ToastHost from './components/ToastHost.vue'
import HomeView from './views/HomeView.vue'
</script>

<template>
  <div class="flex h-full flex-col bg-surface-sunken">
    <!-- 顶栏 -->
    <header class="flex h-12 shrink-0 items-center gap-3 border-b border-surface-border bg-surface px-4">
      <span class="text-accent text-base leading-none">▲</span>
      <h1 class="text-sm font-bold uppercase tracking-widest text-slate-100">
        GunCode <span class="text-accent">// 改枪码管理器</span>
      </h1>
      <span class="ml-auto rounded border border-accent/40 px-2 py-0.5 text-[10px] tracking-wider text-accent">
        三角洲行动
      </span>
    </header>

    <main class="min-h-0 flex-1">
      <HomeView />
    </main>

    <ToastHost />
  </div>
</template>
