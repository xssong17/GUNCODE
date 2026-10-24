<script lang="ts" setup>
import { ref } from 'vue'
import ToastHost from './components/ToastHost.vue'
import HomeView from './views/HomeView.vue'
import { useGunCode } from './composables/useGunCode'

/** 无边框窗口的运行时控制（wails runtime 注入在 window.runtime） */
const rt = () => (window as any).runtime
function winMin() { rt()?.WindowMinimise?.() }
function winMax() { rt()?.WindowToggleMaximise?.() }
function winClose() { rt()?.Quit?.() }

/** 全局数据（含版本信息）；导出/导入逻辑保留在 HomeView，标题栏按钮经 ref 调用 */
const { appInfo } = useGunCode()
const homeViewRef = ref<InstanceType<typeof HomeView> | null>(null)
async function onExport() { await homeViewRef.value?.onExport() }
async function onImport() { await homeViewRef.value?.onImport() }
</script>

<template>
  <div class="flex h-full flex-col bg-surface-sunken">
    <!-- 顶栏：无边框下自绘标题栏（整行可拖动，控制按钮区域不拖动） -->
    <header
      class="flex h-12 shrink-0 items-center border-b border-surface-border bg-surface"
      style="--wails-draggable: drag"
    >
      <div class="flex min-w-0 items-center gap-3 px-4">
        <span class="text-accent text-base leading-none">▲</span>
        <h1 class="truncate text-sm font-bold uppercase tracking-widest text-slate-100">
          GunCode <span class="text-accent">// 改枪码管理器</span>
        </h1>
      </div>

      <div class="ml-auto flex h-full items-center" style="--wails-draggable: no-drag">
        <span class="mr-3 text-[11px] text-slate-500">作者 书包呀 · 数据 v{{ appInfo?.datasetVersion ?? '—' }}</span>
        <span class="mr-2 rounded border border-accent/40 px-2 py-0.5 text-[10px] tracking-wider text-accent">
          三角洲行动
        </span>
        <div class="mr-1 flex h-full items-center gap-1">
          <button
            class="flex h-8 items-center rounded border border-surface-border px-2.5 text-xs text-slate-300 transition hover:bg-surface-raised hover:text-slate-100"
            title="把所有改枪方案导出为 JSON 备份"
            @click="onExport"
          >导出</button>
          <button
            class="flex h-8 items-center rounded border border-surface-border px-2.5 text-xs text-slate-300 transition hover:bg-surface-raised hover:text-slate-100"
            title="从 JSON 备份导入改枪方案"
            @click="onImport"
          >导入</button>
        </div>
        <div class="flex h-full items-stretch">
          <button
            class="flex w-10 items-center justify-center text-sm text-slate-400 transition hover:bg-surface-raised hover:text-slate-100"
            title="最小化"
            @click="winMin"
          >─</button>
          <button
            class="flex w-10 items-center justify-center text-xs text-slate-400 transition hover:bg-surface-raised hover:text-slate-100"
            title="最大化 / 还原"
            @click="winMax"
          >□</button>
          <button
            class="flex w-11 items-center justify-center text-sm text-slate-400 transition hover:bg-rose-600 hover:text-white"
            title="关闭"
            @click="winClose"
          >✕</button>
        </div>
      </div>
    </header>

    <main class="min-h-0 flex-1">
      <HomeView ref="homeViewRef" />
    </main>

    <ToastHost />
  </div>
</template>
