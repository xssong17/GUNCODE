<script lang="ts" setup>
import { ref } from 'vue'
import ToastHost from './components/ToastHost.vue'
import HomeView from './views/HomeView.vue'
import ChangelogModal from './components/ChangelogModal.vue'
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

/** 标题栏下拉菜单与「关于」弹窗 */
const menuOpen = ref(false)
const aboutOpen = ref(false)
const changelogOpen = ref(false)

/** 应用版本号（手动维护，随发布递增；涉及大的功能/UI改动时自增） */
const APP_VERSION = '1.0.0'
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
        <span class="mr-2 rounded border border-accent/40 px-2 py-0.5 text-[10px] tracking-wider text-accent">
          三角洲行动
        </span>
        <div class="relative mr-1 flex h-full items-center">
          <button
            class="flex h-full w-10 items-center justify-center text-lg text-slate-400 transition hover:bg-surface-raised hover:text-slate-100"
            title="菜单"
            @click.stop="menuOpen = !menuOpen"
          >⋮</button>
          <div
            v-if="menuOpen"
            class="absolute right-0 top-full z-50 mt-1 w-44 overflow-hidden rounded-md border border-surface-border bg-surface-raised py-1 shadow-xl"
            @click.stop
          >
            <button
              class="flex w-full items-center px-3 py-2 text-left text-sm text-slate-200 transition hover:bg-surface-sunken"
              @click="menuOpen = false; onExport()"
            >导出备份</button>
            <button
              class="flex w-full items-center px-3 py-2 text-left text-sm text-slate-200 transition hover:bg-surface-sunken"
              @click="menuOpen = false; onImport()"
            >导入备份</button>
            <div class="my-1 h-px bg-surface-border"></div>
            <button
              class="flex w-full items-center px-3 py-2 text-left text-sm text-slate-200 transition hover:bg-surface-sunken"
              @click="menuOpen = false; changelogOpen = true"
            >版本说明</button>
            <button
              class="flex w-full items-center px-3 py-2 text-left text-sm text-slate-200 transition hover:bg-surface-sunken"
              @click="menuOpen = false; aboutOpen = true"
            >关于</button>
          </div>
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

    <!-- 标题栏菜单外部点击关闭 -->
    <div v-if="menuOpen" class="fixed inset-0 z-40" @click="menuOpen = false"></div>

    <!-- 关于弹窗 -->
    <div
      v-if="aboutOpen"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      @click="aboutOpen = false"
    >
      <div class="w-80 rounded-lg border border-surface-border bg-surface p-5 shadow-xl" @click.stop>
        <div class="mb-3 flex items-center gap-2">
          <span class="text-accent text-xl leading-none">▲</span>
          <h2 class="text-base font-bold text-slate-100">GunCode</h2>
        </div>
        <p class="text-sm text-slate-400">三角洲行动 · 改枪码管理器</p>
        <div class="my-3 h-px bg-surface-border"></div>
        <dl class="space-y-1.5 text-sm">
          <div class="flex justify-between">
            <dt class="text-slate-400">作者</dt>
            <dd class="text-slate-100">书包呀</dd>
          </div>
          <div class="flex justify-between">
            <dt class="text-slate-400">版本</dt>
            <dd class="text-slate-100">v{{ APP_VERSION }}</dd>
          </div>
          <div class="flex justify-between">
            <dt class="text-slate-400">数据版本</dt>
            <dd class="text-slate-100">{{ appInfo?.datasetVersion ?? '—' }}</dd>
          </div>
        </dl>
        <button
          class="mt-4 w-full rounded-md bg-accent py-2 text-sm font-medium text-slate-900 transition hover:brightness-110"
          @click="aboutOpen = false"
        >关闭</button>
      </div>
    </div>

    <!-- 版本说明（更新日志） -->
    <ChangelogModal :open="changelogOpen" @close="changelogOpen = false" />
  </div>
</template>
