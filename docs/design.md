# 視覚設計と素材 / Visual design and assets

更新 / Updated: 2026-11-12 (Asia/Tokyo)。公開websiteと説明素材の設計資料です。研究UIの実装仕様は[ui.md](ui.md)、作業と検証記録は[PLAN.md](../PLAN.md)を参照します。図は設計案・合成例であり、実装画面や研究結果ではありません。

This reference covers the public website and explanatory assets. Research UI behavior belongs in [ui.md](ui.md); actions and validation records belong in [PLAN.md](../PLAN.md). Diagrams are proposals and synthetic examples, not product screenshots or research results.

## 素材が必要な理由 / Why use visual assets

本projectに必要なのは、モデル切替、正本と派生記憶、画面内の情報階層を説明する少数の図です。既存の橋形markと説明用SVGを保持し、公開websiteには暖白のcanvasと白い重点面を使います。写真・人物・吉祥物・大学logoは必要ありません。背景画像・外部font・icon packageを追加せず、SVGを編集可能な原本として保持します。

The useful assets explain model switching, authoritative messages versus derived memory, and interface hierarchy. Retain the bridge mark and explanatory SVGs while using a warm-white canvas and white focal surfaces on the public website. Photos, characters, mascots and university logos are unnecessary. Keep editable SVG sources without adding background images, remote fonts or icon packages.

## 配置と責務 / Placement and ownership

| 場所 / Location | 内容と利用 / Content and use |
| --- | --- |
| [website/assets/mark.svg](../website/assets/mark.svg) | 既存project mark。websiteで再利用し、大学公式markとはしない / Existing project mark, reused on the website without presenting it as an official university mark |
| [website/assets/context-continuity.svg](../website/assets/context-continuity.svg) | 1200 × 680。日英共通のA → B → A・訂正・出典図。websiteでは必要時に開き、READMEと同じ原本を参照 / Shared JA/EN continuity, correction and provenance diagram, disclosed on demand on the website and shared with README |
| [docs/assets/research-ui-layout.svg](assets/research-ui-layout.svg) | 1400 × 900。広幅／狭幅の研究UI配置案。ui.mdから参照し、製品screenshotとして使わない / Wide/narrow research UI layout proposal, referenced by ui.md; never label it a product screenshot |
| `website/index.html`, `website/ja.html` | 日英の公開説明。固定合成例のみ。実研究APIへ接続しない / English/Japanese public information pages with fixed synthetic examples and no research API connection |

2026-11-12の利用者指示により、公開websiteと美術素材はassistantが担当します。SVGは新しい実行依存を持たず、script・外部参照・埋込画像を含みません。研究用Go・Python・React `web/` は本人担当です。公開siteのCSSと、研究UIのMantine／独自CSSなしの方針を混同しません。

The user's 2026-11-12 instruction assigns the public website and artwork to the assistant. These SVGs add no runtime dependencies, scripts, external references or embedded images. The user retains research Go, Python and React `web/`. Public-site CSS is separate from the research UI's Mantine/no-handwritten-CSS decision.

## 参照と適用範囲 / Reference and scope

2026-11-12に[九州大学公式website](https://www.kyushu-u.ac.jp/ja/)の実際の表示とbrowserの算出styleを確認しました。白いbody／header、無彩色の文字、紫紅accentの観測は[資料台帳](sources.md)に履歴として保持します。これらは現在のproject paletteを拘束する値でも、大学brand guideline全体の調査結果でもありません。

The [official Kyushu University website](https://www.kyushu-u.ac.jp/ja/) was inspected visually and through browser-computed styles on 2026-11-12. Observations of its white body/header, neutral text and burgundy accents remain historical records in the [source register](sources.md). They neither prescribe the current project palette nor constitute a complete university brand-guideline audit.

大学siteは一般的な情報階層の参考に留め、配色はproject独自に調整します。大学logo、写真、文章、公式brandingは転用せず、独立projectの表示を維持します。参照siteのremote fontは導入せず、日本語には端末のMeiryoを先頭とするsystem fontを使います。

Use the university site as a general hierarchy reference while developing the project's own palette. No university logos, photos, copy or official branding are reused; retain the independent-project identification. Do not import the reference site's remote fonts; use system fonts with Meiryo first for Japanese.

## 視覚言語 / Visual language

| 要素 / Element | 用途 / Use |
| --- | --- |
| Warm white `#faf9f6` | 利用者が選択したpage canvas / User-selected page canvas |
| White `#ffffff` | masthead、合成demo、資料領域の重点面 / Focal surfaces for the masthead, synthetic demo and resources |
| Wine `#783b51` / hover `#5d2c3e` | navigation、link、選択状態の小さなaccent。成功・失敗の判定色にしない / Small navigation, link and selection accents; not a success/failure code |
| Heading `#292725` / body `#3d3935` / muted `#716b64` | 見出し、本文、補足を暖灰の明度で区分 / Warm-gray luminance levels for headings, body and supporting text |
| Selected surface `#f1ece6` / line `#e2ddd5` | 選択面と細い区切り。選択状態は文字・aria属性も併用 / Selected surfaces and fine dividers; also express selection through labels and ARIA |
| System sans-serif | 日本語はMeiryo優先。見出しと本文に端末のfontを使い、外部fontを取得しない / Meiryo-first Japanese; system fonts for headings and body without remote font downloads |

2026-11-12の最新の利用者選択は暖白canvas `#FAF9F6` です。白い重点面、暖灰の文字・区切り、低彩度の酒紅accentは、この選択に合わせたproject独自の配色です。light color schemeと既存layoutを維持し、今回の変更はpaletteに限定します。

The latest user choice on 2026-11-12 is the warm-white canvas `#FAF9F6`. White focal surfaces, warm-gray text/dividers and muted wine accents are complementary colors developed for this project. Keep the light color scheme and existing layout; this change is limited to the palette.

headerはproject名と言語切替の行、その下の白いnavigation行に分けます。heroは研究の説明と短い受入状況のasideを並べ、別のstatus帯で同じ情報を繰り返しません。本文は問い→合成例→構成→調査方法→正本文書の順で、sectionとfooterの余白を詰めます。資料は見出しの下に2列×2行の簡潔なlinkを置き、狭幅では1列にします。

The header has a project-name/language row followed by a white navigation row. The hero pairs the research introduction with a compact evidence aside, avoiding a separate strip that repeats the status. The body follows question → synthetic example → architecture → method → reference documents, with tighter sections and footer spacing. Place compact resource links in a two-by-two grid below their heading, switching to one column on narrow screens.

構成説明の継続図はnative HTMLの `details` / `summary` で必要時に開き、閉じた状態でも要点を本文で読めるようにします。訂正前後の違いと証拠区分、keyboard focusを維持し、heroに装飾図を戻しません。

The continuity diagram in the architecture section opens on demand through native HTML `details` / `summary`; the prose explains the essentials while it is closed. Preserve visible corrections, evidence labels and keyboard focus, without restoring decorative hero artwork.

既存SVGの梅紫 `#7b1746`、紙色 `#f5f3ef`、線色 `#ddd7da` は原本として保持し、今回の再設計では変更しません。既存project markも保持します。図内の色はwebsiteのsection背景ではありません。研究UIはdefault Mantine themeとcomponent propsを優先し、公開siteのpaletteをそのまま研究UI themeとして導入する要件にはしません。

Retain the existing SVGs' plum `#7b1746`, paper `#f5f3ef` and divider `#ddd7da` colors in their unchanged sources, along with the existing project mark. Colors within artwork are not website-section backgrounds. The research UI prioritizes the default Mantine theme and component props; the public-site palette does not require a matching research UI theme.

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

今回のpalette調整では既存SVGとlayoutを変更せず、新しい画像・依存・font・modelを追加しません。公開siteの表示変更と検証結果はPLAN.mdに記録し、設計方針を実測結果や公開完了の証拠とはしません。

This palette refinement leaves SVG sources and layout unchanged and adds no images, dependencies, fonts or models. Record site changes and validation results in PLAN.md; design decisions are not evidence of measurements or completed publication.

画像の値や語句を更新するときは関連本文と日英を同時に確認します。SVGの構文、外部参照の不在、実際の描画、代替text、local linkを確認します。画像の描画確認は、研究UIの動作検証やモデル評価の証拠にはなりません。

When changing labels or values, review the paired prose and both languages together. Check SVG syntax, absence of external references, actual rendering, alt text and local links. Rendering checks do not validate research UI behavior or model quality.
