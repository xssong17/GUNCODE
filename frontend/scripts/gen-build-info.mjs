// 构建时生成 git 版本信息：读当前 HEAD 的 short hash，写入 src/generated/buildInfo.ts
// 由 frontend package.json 的 build 钩子调用，确保「构建版本」与本地 git 一致。
import { execSync } from 'node:child_process'
import { writeFileSync, mkdirSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))
const root = resolve(here, '..') // 项目根（frontend 上一级）

let v = 'dev'
try {
  v = execSync('git rev-parse --short HEAD', { cwd: root, encoding: 'utf8' }).trim()
} catch {
  // 非 git 环境，降级为 dev
}

const dir = resolve(here, '../src/generated')
mkdirSync(dir, { recursive: true })
writeFileSync(
  resolve(dir, 'buildInfo.ts'),
  `// 本文件由 scripts/gen-build-info.mjs 在构建时自动生成，请勿手改。\n/** 与本地 git 版本一致的构建号 */\nexport const GIT_VERSION = '${v}'\n`,
  'utf8',
)
console.log('[buildInfo] GIT_VERSION =', v)
