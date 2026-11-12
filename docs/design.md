# 視覚設計と素材 / Visual design and assets

更新 / Updated: 2026-11-12 (Asia/Tokyo)。公開websiteと説明素材の設計資料です。研究UIの実装仕様は[ui.md](ui.md)、作業と検証記録は[PLAN.md](../PLAN.md)を参照します。図は設計案・合成例であり、実装画面や研究結果ではありません。

This reference covers the public website and explanatory assets. Research UI behavior belongs in [ui.md](ui.md); actions and validation records belong in [PLAN.md](../PLAN.md). Diagrams are proposals and synthetic examples, not product screenshots or research results.

## 素材が必要な理由 / Why use visual assets

本projectに必要なのは、モデル切替、正本と派生記憶、画面内の情報階層を説明する少数の図です。既存の橋形markと梅紫・紙色を再利用します。写真・人物・吉祥物・大学logoは必要ありません。背景画像・外部font・icon packageを追加せず、単純なSVGを編集可能な原本として保持します。

The useful assets explain model switching, authoritative messages versus derived memory, and interface hierarchy. Reuse the existing bridge mark and plum/paper palette. Photos, characters, mascots and university logos are unnecessary. Keep simple editable SVG sources without adding background images, remote fonts or icon packages.

## 配置と責務 / Placement and ownership

| 場所 / Location | 内容と利用 / Content and use |
| --- | --- |
| [website/assets/mark.svg](../website/assets/mark.svg) | 既存project mark。websiteで再利用し、大学公式markとはしない / Existing project mark, reused on the website without presenting it as an official university mark |
| [website/assets/context-continuity.svg](../website/assets/context-continuity.svg) | 1200 × 680。日英共通のA → B → A・訂正・出典図。websiteとREADMEが同じ原本を参照 / Shared JA/EN continuity, correction and provenance diagram; the website and README reference one source |
| [docs/assets/research-ui-layout.svg](assets/research-ui-layout.svg) | 1400 × 900。広幅／狭幅の研究UI配置案。ui.mdから参照し、製品screenshotとして使わない / Wide/narrow research UI layout proposal, referenced by ui.md; never label it a product screenshot |
| `website/index.html`, `website/ja.html` | 日英の公開説明。固定合成例のみ。実研究APIへ接続しない / English/Japanese public information pages with fixed synthetic examples and no research API connection |

2026-11-12の利用者指示により、公開websiteと美術素材はassistantが担当します。SVGは新しい実行依存を持たず、script・外部参照・埋込画像を含みません。研究用Go・Python・React `web/` は本人担当です。公開siteのCSSと、研究UIのMantine／独自CSSなしの方針を混同しません。

The user's 2026-11-12 instruction assigns the public website and artwork to the assistant. These SVGs add no runtime dependencies, scripts, external references or embedded images. The user retains research Go, Python and React `web/`. Public-site CSS is separate from the research UI's Mantine/no-handwritten-CSS decision.

## 視覚言語 / Visual language

| 要素 / Element | 用途 / Use |
| --- | --- |
| Plum `#7b1746` | projectの識別、主要link、注目する訂正。成功・失敗の判定色にしない / Identity, primary links and the focal correction; not a success/failure code |
| Ink `#242125` / muted `#68636a` | 見出し・本文と補足の階層 / Heading/body and supporting text hierarchy |
| Paper `#f5f3ef` / white `#ffffff` | 説明領域と会話・記録を区分 / Separate explanatory areas from conversation/record surfaces |
| Line `#ddd7da` | 補助的な区切り。操作可能性や状態を線だけで伝えない / Supporting dividers; never the sole indicator of an action or state |

公開siteは文章中心の研究資料として、問い→合成例→構成→調査方法→正本文書の順に案内します。装飾より、訂正前後の違いと証拠区分を読み取りやすくします。研究UIはdefault Mantine themeとcomponent propsを優先し、このpaletteをそのままthemeとして導入する要件にはしません。

The public site reads as a research publication: question → synthetic example → architecture → method → reference documents. Prioritize visible corrections and evidence labels. The research UI prioritizes the default Mantine theme and component props; this palette is not a requirement to implement a custom theme.

## 画像と同等の文章 / Text equivalents

継続図：同じsessionと所有者・policyのままA → B → Aを切り替えます。合成例は11月12日の締切を15日に訂正し、最後のAにも15日を引き継ぐ目標です。正本message `m02` が訂正の出典で、memoryは出典とrevisionを持つ派生情報です。成功を実測した図ではありません。

Continuity diagram: switch A → B → A within one session with the same owner and policy. The synthetic example corrects a November 12 deadline to November 15 and aims to carry that correction back to A. Authoritative message `m02` supplies the correction; memory is derived and carries its source and revision. The diagram does not report measured success.

UI図：広幅はsession navigationと会話を並べ、出典は必要時に開きます。狭幅は一列にします。次turn用のmodelと過去回答のmodel labelを分離します。保存済み回答と `memory_pending` を同時に表示し、回復まで次の送信を止めます。図中の入力やbuttonは説明用で、操作可能な実装ではありません。

UI diagram: wide layouts place sessions beside the conversation and open sources on demand; narrow layouts use one column. Separate the next-turn model from the model labels on past answers. Show a saved answer together with `memory_pending` and block the next send until recovery. Drawn fields and buttons are explanatory, not interactive implementation.

説明画像には日英のtitle・descと意味のある代替textを付け、本文でも同じ情報を説明します。markが隣接するproject名を重複する場合は空のaltにします。狭幅では図を縮小するだけに依存せず、読みやすい本文と原寸画像linkを提供します。色に加えて文字を使い、keyboard focus・日本語の折返し・動きを減らす設定を維持します。

Give explanatory images JA/EN titles/descriptions and meaningful alt text, and repeat their substance in prose. Use an empty alt when the mark duplicates the neighboring project name. At narrow widths provide readable prose and a full-size image link instead of relying on scaled text alone. Combine color with labels and preserve keyboard focus, Japanese wrapping and reduced-motion behavior.

## 制作・更新 / Creation and maintenance

新図は2026-11-12にassistantが既存project paletteからSVGとして制作しました。第三者画像の転用、画像生成model、商用API、font downloadは使用していません。既存markの作者・制作過程を新たに推定しません。原著作物の配布条件はrepositoryの[MIT](../LICENSE)、資料の区分と参照日は[sources.md](sources.md)に記録します。

The assistant authored the new SVGs on 2026-11-12 using the existing project palette. No third-party imagery, image-generation model, commercial API or font download was used. Do not infer the original mark's authorship or production history. Original project material follows the repository's [MIT license](../LICENSE); evidence classes and inspection dates are recorded in [sources.md](sources.md).

画像の値や語句を更新するときは関連本文と日英を同時に確認します。SVGの構文、外部参照の不在、実際の描画、代替text、local linkを確認します。画像の描画確認は、研究UIの動作検証やモデル評価の証拠にはなりません。

When changing labels or values, review the paired prose and both languages together. Check SVG syntax, absence of external references, actual rendering, alt text and local links. Rendering checks do not validate research UI behavior or model quality.
