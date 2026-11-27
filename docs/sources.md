# 資料・ライセンス台帳 / Source and license register

参照日 / Access date: **2026-11-08 (Asia/Tokyo)**。特記のない資料はこの日に確認しました。上流資料の調査であり、導入済みdistributionの完全監査ではありません。新しい依存・modelを導入／取得していません。

Unless noted, sources were inspected on 2026-11-08. This is upstream-source research, not a complete audit of an installed distribution. No new dependencies or models were installed/downloaded.

## 証拠の区分 / Evidence classes

| 区分 / Class | このprojectでの例 / Project example |
| --- | --- |
| 確認済み要件 / Confirmed requirement | OSS、本人による実装、assistantのテスト、日英文書 / OSS, user implementation, assistant tests, bilingual docs |
| 利用者報告 / User-reported | 教員承認、玄界account・割当、「commandlineの形でAPIを提供し、UI画面は作らない」という教員の要件（2026-11-27） / Supervisor approval, Genkai account/allocation, and the supervisor's requirement of an API used from the command line with no UI screen (2026-11-27) |
| 確認済み設計選択 / Confirmed design choice | 公開siteの暖白canvas `#FAF9F6`（利用者選択、2026-11-12） / Public-site warm-white canvas `#FAF9F6` (user choice, 2026-11-12) |
| 設計提案 / Design proposal | SQLite、module境界、性能目標、[CLI](cli.md)のcommandと終了code / SQLite, module boundaries, performance targets, and [CLI](cli.md) commands and exit codes |
| 上流資料・vendor claim / Upstream documentation or vendor claim | 対応provider、公開benchmark、製品説明 / Supported providers, published benchmarks and product descriptions |
| 実測 / Measured result | 2026-11-08の初期scaffold：Go 1.26.0 Windowsで当時の2テスト成功。その後の受入は[PLAN.md](../PLAN.md)。モデル品質は未測定 / Initial scaffold snapshot on 2026-11-08: the two tests present then passed on Go 1.26.0 Windows. Subsequent acceptance is recorded in PLAN.md; model quality remains unmeasured |

## 公開デザイン素材 / Public design assets

2026-11-12参照。既存projectの `website/assets/mark.svg` と当初の梅紫・紙色を基に、[継続図](../website/assets/context-continuity.svg)をassistantが編集可能なSVGとして制作しました。同時に制作した研究UI配置図は、2026-11-27の研究UI取り下げに伴って削除しました。既存markは保持し、追加の第三者画像・font・依存・画像生成modelは使っていません。websiteのpalette調整でもSVG原本は変更せず、図内の色とpage paletteを区別します。原著作物の配布条件は[MIT](../LICENSE)。用途、寸法、代替説明は[視覚設計](design.md)を参照します。

Inspected 2026-11-12. The assistant authored the editable [continuity diagram](../website/assets/context-continuity.svg) using the existing `website/assets/mark.svg` and original plum/paper palette as references. The existing mark is retained; no additional third-party images, fonts, dependencies or image-generation models were used. The research UI layout figure made at the same time was deleted when the research UI was withdrawn on 2026-11-27. Website palette changes leave SVG sources unchanged, distinguishing artwork colors from the page palette. Original project material follows [MIT](../LICENSE). See [visual design](design.md) for uses, dimensions and text equivalents.

これらは設計提案・合成例です。model versionは該当なしで、実modelの応答・製品のscreenshot・benchmarkとして扱いません。外部資料の従来の参照日を今回の描画確認日で上書きしません。

These are design proposals and synthetic examples. Model version is not applicable; they are neither real-model responses, product screenshots nor benchmarks. The rendering review date does not replace earlier access dates for external sources.

## 公開siteの版面参照 / Public-site layout reference

| 資料 / Source | 参照日・確認内容 / Inspection date and observation | 採用範囲 / Application |
| --- | --- | --- |
| [九州大学公式website / Official Kyushu University website](https://www.kyushu-u.ac.jp/ja/) | 2026-11-12。実際の表示とbrowserの算出styleを確認。下表は過去の参照観測として保持するsnapshot / Visual and browser-computed-style inspection on 2026-11-12; preserve the sampled values below as a historical reference snapshot | 一般的な情報階層の参考。観測色は現在のproject paletteを拘束しない。大学logo・写真・文章・公式brandingは転用しない / General hierarchy reference; observed colors do not prescribe the current project palette. No reuse of university logos, photos, copy or official branding |

| 公式siteで確認した対象 / Sampled official-site element | 算出値・指定 / Computed value or declaration |
| --- | --- |
| body・header背景 / Body and header backgrounds | `rgb(255, 255, 255)` / `#ffffff` |
| body文字 / Body text | `rgb(51, 51, 51)` / `#333333` |
| 見出し・navigation文字 / Heading and navigation text | `rgb(0, 0, 0)` / `#000000` |
| footer背景 / Footer background | `rgb(0, 0, 0)` / `#000000` |
| 危機管理button背景・ALL NEWS枠 / Crisis-management button background and ALL NEWS border | `rgb(133, 2, 62)` / `#85023e` |
| Font指定 / Font declarations | bodyはMeiryoを含み、headerはNoto Sans JPを指定 / Body includes Meiryo; header specifies Noto Sans JP |

上記は当日の要素ごとの観測で、大学brand guideline全体や各端末の実使用fontを検証した結果ではありません。利用者は2026-11-12に暖白canvas `#FAF9F6` を確認し、少量の黒の追加を依頼しました。白い重点面・暖灰本文・酒紅 `#783b51` を保ち、見出し・navigation・brand、headerの細線と全幅footerにproject独自の墨黒 `#171717` を使います。大学siteの黒と同一色の採用ではありません。現在値は[視覚設計](design.md)に集約し、light color scheme、system font、本文layoutとSVG原本を維持します。

These are element-level observations on the inspection date, not verification of the complete university brand guidelines or the font actually used on every device. On 2026-11-12 the user confirmed the warm-white canvas `#FAF9F6` and requested modest black accents. Preserve white focal surfaces, warm-gray body text and wine `#783b51`, using project-selected near-black `#171717` for headings/navigation/brand, a thin header rule and the full-width footer. This does not copy the university site's exact black. Current values are maintained in [visual design](design.md); retain the light color scheme, system fonts, main-content layout and SVG sources.

大学siteは参考であり、現在の配色を大学による要件・公認や研究効果の証拠とはしません。今回のpalette調整で新しい画像・依存・font・modelは追加しません。

The university site is a reference; the current palette is not a university requirement, endorsement or research evidence. This palette refinement adds no images, dependencies, fonts or models.

## Coreの調査snapshot / Core survey snapshots

| 対象 / Component | 調査対象と根拠 / Inspected snapshot and evidence | 確認範囲・未確認 / Scope and remaining uncertainty |
| --- | --- | --- |
| Mem0 OSS | `2.2.1`, commit `753496cf1435d76c3a711f3af51303f7fdcb6420`: [manifest](https://github.com/mem0ai/mem0/blob/753496cf1435d76c3a711f3af51303f7fdcb6420/pyproject.toml), [LICENSE](https://github.com/mem0ai/mem0/blob/753496cf1435d76c3a711f3af51303f7fdcb6420/LICENSE) | Apache-2.0確認。Python `>=3.10,<4.0`。wheelと全推移依存は未監査 / Apache-2.0 verified; Python range from manifest; wheel/transitive dependencies unaudited |
| Ollama | commit `c1db5f9ef203b2e946161bf4c6bd4f2ec7f123b3`: [LICENSE](https://github.com/ollama/ollama/blob/c1db5f9ef203b2e946161bf4c6bd4f2ec7f123b3/LICENSE), [chat API](https://docs.ollama.com/api/chat), [FAQ](https://docs.ollama.com/faq) | MIT確認。binary、同梱library、cloud設定、token対応は固定版で要検証 / MIT verified; binary, bundled libraries, cloud settings and token behavior require pinned-version checks |
| Qdrant client | `v1.20.0`: [LICENSE](https://github.com/qdrant/qdrant-client/blob/v1.20.0/LICENSE), [local mode](https://github.com/qdrant/qdrant-client/blob/v1.20.0/README.md) | Apache-2.0、server不要のlocal modeを確認。Mem0との組合せは未実行 / Apache-2.0 and serverless local mode verified; Mem0 combination not executed |
| SQLite | [用途 / Uses](https://sqlite.org/whentouse.html), [copyright](https://sqlite.org/copyright.html), [WAL](https://sqlite.org/wal.html) | SQLite本体はpublic domain。driver、build tool、配布物は別確認 / Core SQLite is public domain; driver/build tools/distribution need separate review |
| Go SQLite driver候補 / Candidate | [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | 調査時表示 `v1.62.0`, BSD-3-Clause。採用version未固定、libc等の整合確認が必要 / Displayed v1.62.0, BSD-3-Clause; adoption unpinned, libc/dependency compatibility required |
| Model候補 / Candidates | [Qwen3-4B](https://huggingface.co/Qwen/Qwen3-4B), [Qwen3-8B](https://huggingface.co/Qwen/Qwen3-8B), [Qwen3-Embedding-0.6B](https://huggingface.co/Qwen/Qwen3-Embedding-0.6B) | publisher cardはApache-2.0。weights/tokenizer/template、変換・量子化artifactのrevision/hashは未固定 / Publisher cards identify Apache-2.0; exact artifacts/revisions/hashes unpinned |

Mem0 manifestにはqdrant-client、pydantic、openai、httpx、posthog、pytz、sqlalchemy、protobufが必須依存として記載されています。SDK packageの存在と商用API呼出しは別問題です。Ollama clientは必要分だけ別途解決し、広いextrasを無条件導入しません。

The Mem0 manifest lists qdrant-client, pydantic, openai, httpx, posthog, pytz, sqlalchemy and protobuf as required dependencies. Having an SDK package is distinct from calling a commercial API. Resolve the needed Ollama client separately; do not install broad extras blindly.

Mem0の[OSS概要](https://docs.mem0.ai/open-source/overview)はdefaultにcloud providerを含みます。[local LLM例](https://docs.mem0.ai/components/llms/models/ollama)と[local embedding例](https://docs.mem0.ai/components/embedders/models/ollama)はそれぞれ片側だけの設定なので、単独copyではlocal-onlyを証明しません。generation、embedding、vector/history、telemetryをすべて固定し、外向き通信を遮断して確認します。

Mem0's [OSS overview](https://docs.mem0.ai/open-source/overview) includes cloud providers in its defaults. The [local LLM](https://docs.mem0.ai/components/llms/models/ollama) and [embedding](https://docs.mem0.ai/components/embedders/models/ollama) examples configure one side at a time; copying either alone does not establish a local-only route. Pin generation, embedding, vector/history and telemetry settings, then verify with outbound traffic blocked.

## 取り下げたfrontend候補 / Withdrawn frontend candidates

2026-11-27に研究UIを取り下げたため、React、TypeScript、Vite、Mantine、Ant Design、Material UI、NextChat、Open WebUI、LibreChatと、Vitest、Testing Library、Playwrightの調査行をこの台帳から外しました。いずれも導入していません。調査記録はGit履歴（`a40978e` まで）に残ります。CLIはGo標準libraryだけで作る案で、新しい依存はありません。

Because the research UI was withdrawn on 2026-11-27, the survey rows for React, TypeScript, Vite, Mantine, Ant Design, Material UI, NextChat, Open WebUI and LibreChat, and for Vitest, Testing Library and Playwright, were removed from this register. None was ever installed. The survey remains in Git history up to `a40978e`. The CLI is proposed on the Go standard library alone, adding no dependency.

## 配布manifestの必須情報 / Required distribution-manifest information

採用時の台帳にはcomponent名、役割、直接／推移依存、version/commit、取得元、artifact SHA-256、LICENSE/NOTICE、runtime/build/test区分、OS/architecture、検証日、検証結果を記録します。Go、Python、Nodeのlockfileとruntime版、model/tokenizer/template、量子化recipeを同じrun manifestから参照できるようにします。

At adoption, record component, role, direct/transitive dependency, version/commit, source URL, artifact SHA-256, LICENSE/NOTICE, runtime/build/test category, OS/architecture, verification date and outcome. A run manifest should reference Go/Python/Node lockfiles and runtimes, model/tokenizer/template and quantization recipe.

上流repoのlicenseだけではwheel、container、binary、fonts/icons、GPU libraryの許諾まで確認したことになりません。必須の非OSS依存が残る構成は採用しません。NVIDIA等の学内GPU基盤をOSS要件のどこまでに含めるかは未確定で、review済みCPU routeを維持します。

An upstream repository license does not audit wheels, containers, binaries, fonts/icons or GPU libraries. Do not adopt configurations with unresolved mandatory non-OSS dependencies. The boundary for institution-provided GPU platforms such as NVIDIA remains unresolved; retain a reviewed CPU route.

## 玄界の資料 / Genkai sources

[PJM jobs](https://www.cc.kyushu-u.ac.jp/scp/usage/job/)、[Singularity](https://www.cc.kyushu-u.ac.jp/scp/system/genkai/software/Singularity/)、[storage](https://www.cc.kyushu-u.ac.jp/scp/usage/storage/)を参照しました。account、残枠、実際のmodule、利用可能GPU、隔離方式は未確認です。資料閲覧は遠隔実行を許可しません。

[PJM jobs](https://www.cc.kyushu-u.ac.jp/scp/usage/job/), [Singularity](https://www.cc.kyushu-u.ac.jp/scp/system/genkai/software/Singularity/) and [storage](https://www.cc.kyushu-u.ac.jp/scp/usage/storage/) were consulted. Entitlement, remaining allocation, actual modules/GPUs and isolation are unverified. Reading documentation does not authorize remote execution.

## 本projectのMIT / Project MIT license

[LICENSE](../LICENSE)の標準英語原文を保持します。説明：著作権表示と許諾文を保持する条件で、使用、複製、変更、配布等を許可し、無保証で提供します。この日本語説明はlicense原文を置換せず、第三者のコード・modelの条件を変更しません。

Keep the standard English text in [LICENSE](../LICENSE). Explanation: use, copying, modification and distribution are permitted subject to retaining the copyright/permission notice; the material is provided without warranty. This explanation does not replace the license or alter third-party code/model terms.
