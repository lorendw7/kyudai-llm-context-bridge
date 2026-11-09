# 研究UI設計 / Research UI design

更新 / Updated: 2026-11-08。小型の自作React UIを先行し、第三者UIは後続の任意拡張とします。以下は提案で、UIコード・package manifest・依存はまだ作成／導入していません。 / Build a small owned React UI first; third-party clients are optional later. This is a proposal; UI code, package manifests and dependencies have not been created/installed.

## 選択：React + TypeScript + Vite + Mantine / Selection

Reactで画面を分割し、TypeScriptでAPI契約を表現し、Viteで開発・静的buildを行います。Mantineの部品とlibrary stylesを利用し、初期版では独自CSS、Tailwind、CSS Modulesを追加しません。default theme、layout props、component propsで実現できる範囲に画面を絞ります。

Use React for views, TypeScript for API contracts and Vite for development/static builds. Use Mantine components and library styles; add no handwritten CSS, Tailwind or CSS Modules initially. Keep the design within default themes, layout props and component props.

これは保守負担を抑える設計判断です。React公式資料はViteによる構築を説明していますが、routing・data取得を自分で管理する負担も示しています。小型local UIではその範囲を限定します。[React](https://react.dev/learn/build-a-react-app-from-scratch)、[Mantine/Vite](https://mantine.dev/guides/vite/)、[スタイル / Styles](https://mantine.dev/styles/styles-overview/)（2026-11-08参照）。

This is a maintenance decision. React's documentation describes Vite-based setup and also notes the responsibility for routing/data fetching. Limit that responsibility for this small local UI. Sources: [React](https://react.dev/learn/build-a-react-app-from-scratch), [Mantine/Vite](https://mantine.dev/guides/vite/), [styles](https://mantine.dev/styles/styles-overview/) (accessed 2026-11-08).

| 候補 / Candidate | 適合性 / Fit | 判断 / Decision |
| --- | --- | --- |
| Mantine | AppShell、layout、入力、dialog、feedbackを組み合わせやすい / Composable shell, layout, inputs, dialogs and feedback | 2026-11-08に選択。必要部品だけ使う / Selected 2026-11-08; use only required parts |
| Ant Design | Form/Table中心の管理画面向け / Fits form/table-heavy administration | 実験管理台が主用途になれば再検討 / Revisit if experiment administration dominates |
| Material UI | 一貫したMaterial Designの部品群 / Consistent Material Design components | 好みが合えば代替。追加製品のlicenseを個別確認 / Alternative if preferred; audit additional products separately |

Mantineだけを導入対象とし、他候補は同時導入しません。Mantine core LICENSEはMITを確認しましたが、実際に使うpackage、icons、test dependenciesと推移依存は未監査です。部品はdefault themeのままでも利用できますが、libraryを使うだけでaccessibilityが保証されるわけではありません。[資料台帳](sources.md)を参照してください。

Only Mantine is in installation scope; do not install the alternatives alongside it. Mantine core is verified as MIT; actual packages, icons, test dependencies and transitive dependencies remain unaudited. Components can use the default theme, but a library alone does not guarantee accessibility. See the [source register](sources.md).

## 3つの画面 / Three views

| 画面 / View | 操作と表示 / Interaction and display | Mantine候補 / Candidate components |
| --- | --- | --- |
| Chat / 会話 | session一覧、policyは作成時選択、modelはturnごとに選択、送信、取消、状態、再読込 / Sessions, policy at creation, model per turn, send/cancel/status/reload | AppShell, NavLink, Select, ScrollArea, Paper, Textarea, Button, Badge, Alert |
| Memory / 記憶 | 記憶text、source message、revision、訂正、検索、export/delete / Memory text, source message, revision, corrections, search, export/delete | Table, Accordion, TextInput, Drawer, Modal |
| Results / 結果 | 既存runのpolicy、model、品質、token、時間、失敗・欠測 / Existing runs: policy, model, quality, tokens, time, failures/missing data | Table, Select, Badge, Alert |

最初のUIはChatとMemoryを実装し、Resultsはrunnerのmetrics schemaが固定してから追加します。実験開始buttonは初期版に含めません。UI操作で大量推論や遠隔jobを暗黙起動させないためです。図表は必要になってからOSSのchart libraryを個別監査します。

Implement Chat and Memory first; add Results after the runner's metrics schema is fixed. The initial UI has no experiment-launch button, preventing accidental bulk inference or remote jobs. Audit an OSS chart library separately only when plots are needed.

## 画面とモジュール / View and module structure

```text
web/src/
  app/                 全体構成・日英navigation / shell, providers, JA/EN navigation
  features/chat/       会話・入力・turn状態 / session, transcript, composer, turn status
  features/memory/     記憶・出典・出力・削除 / memories, sources, export/delete
  features/results/    結果閲覧 / read-only runs and metrics
  lib/api/             API通信・検証・エラー変換 / API client, validation, errors
  lib/contracts/       版付きAPI型 / versioned API types
  locales/             日英文字列 / ja/en strings
  theme/               Mantine設定 / Mantine provider configuration
  test/                assistant担当のfixtureとtest設定 / assistant-owned fixtures/setup
```

上記は予定の配置です。各featureがAPI clientを利用し、network処理をcomponentに散らしません。状態はまずReactのlocal state／reducerで管理し、session ID・server revisionをkeyにします。複雑なrouting、Redux、global query cacheは初期追加しません。共通wrapperはdomain上の意味がある部品だけにし、Mantine全体を包み直しません。

This is a target layout. Features use the API client instead of scattering network logic across components. Start with React local state/reducers, keyed by session ID and server revision. Do not initially add complex routing, Redux or a global query cache. Create shared wrappers only for domain concepts, not for every Mantine component.

Goがmodel allowlist、session所有権、policy、記憶、計測を決定します。UIはそれを表示します。UIで履歴要約やtoken数の推定を行ってserverの実測値として表示しません。model変更時も同じsession IDを使います。

Go determines the model allowlist, session ownership, policy, memory and metrics. The UI displays them. Do not summarize history or estimate tokens in the UI and present them as server measurements. Switching models keeps the same session ID.

## 状態・安全な操作 / States and safe interaction

accepted、generating、memory_pending、complete、failed、recovery_required、deletingを区別します。answerが表示されてもmemoryが未確定なら次の送信を止め、保存済みanswerと回復状態を同時に表示します。pollingを中止してもserverの処理取消とはみなしません。

Distinguish accepted, generating, memory_pending, complete, failed, recovery_required and deleting. A visible answer with uncommitted memory blocks the next send and shows both the saved answer and recovery state. Stopping polling does not imply server cancellation.

入力は日本語IMEの変換確定を送信と混同せず、送信中の二重操作を抑止します。失敗時に入力を保持します。削除は対象session名と範囲をdialogで示し、完了後にmemory、message、選択状態をcacheから除去します。別sessionの遅い応答を現在画面へ混入させません。

Do not confuse Japanese IME composition confirmation with submission. Prevent duplicate sends and retain the draft on failure. A deletion dialog identifies the session and scope; completion clears memory, messages and selection from caches. Late responses from another session must not enter the current view.

モデルと記憶のtextは最初はplain textとして表示します。任意HTMLを実行しません。認証tokenはbrowser memory内だけに置き、再読込後は再認証します。`VITE_*`、localStorage、bundle、logに秘密を入れません。UI originを限定し、許可のないcross-origin要求をserverが拒否します。

Initially render model/memory text as plain text, without arbitrary HTML. Keep authentication tokens only in browser memory and reauthenticate after reload. Put no secrets in VITE_* variables, localStorage, bundles or logs. Restrict UI origins and reject unauthorized cross-origin requests at the server.

## 配信・保守 / Serving and maintenance

開発中はViteのlocal proxyからGoへ接続し、配布時は静的buildをGoが配信します。Nodeは本番推論runtimeに不要です。公開 `website/` とprivate研究 `web/` を分離し、GitHub Pagesからlocal APIを呼びません。依存導入やbuild設定は本人の実装範囲です。

During development, a local Vite proxy reaches Go; for distribution, Go serves the static build. Node is not needed in the inference runtime. Separate public website/ from private research web/; GitHub Pages does not call the local API. Dependency installation and build configuration belong to the user.

libraryのversionとlockfileを固定し、更新時はtypecheck、component tests、主要E2Eを実行します。日本語・英語のlabel、keyboard、focus復帰、狭幅を検証します。表示文言とerror codeの辞書を分離し、library更新がAPI仕様を変えないようにします。

Pin library versions and lockfiles. Upgrades run type checks, component tests and core E2E flows. Verify JA/EN labels, keyboard behavior, focus restoration and narrow layouts. Separate UI strings/error-code dictionaries so library upgrades do not change the API contract.

## 第三者UIを後で接続する場合 / Later third-party clients

NextChatはMITの候補で、custom endpointの設定とchat payloadを確認しました。ただし確認したpayloadにBridgeのsession IDはありません。base URLの変更だけでsession継続が完成するとは扱いません。Open WebUI現行版には追加branding条件があり、default OSSルートには選びません。LibreChatはより広いapplication stackを持つため初期採用しません。

NextChat is an MIT candidate; custom-endpoint configuration and chat payload code were inspected. The inspected payload does not include a Bridge session ID, so changing the base URL alone does not establish session continuity. Current Open WebUI has an additional branding condition and is not selected for the default OSS route. LibreChat has a broader application stack and is not selected initially.

互換層は `/compat/openai/v1/models` と `/compat/openai/v1/chat/completions` を別境界にする案です。認証済みownerと明示session binding、request ID、stream/error/cancelの対応が必要です。browser履歴のhashからsessionを推測しません。UI側の要約・自動title・retryは無効化または記録し、正本履歴を上書きさせません。

The proposed compatibility boundary uses /compat/openai/v1/models and /compat/openai/v1/chat/completions. It requires authenticated ownership, explicit session binding, request IDs and stream/error/cancel mapping. Never infer sessions from browser-history hashes. Disable or account for UI summaries, automatic titles and retries; do not allow them to overwrite authoritative history.

接続adapterの実装も本人が担当します。native APIで再現できるまで研究評価には使いません。直接Ollamaに接続するUIはBridgeの記憶・権限検証を迂回するため、本研究の連携demoとは扱いません。

The user also implements compatibility adapters. Do not use them for research evaluation until native-API reproduction is established. A UI connected directly to Ollama bypasses Bridge memory/authorization and is not a demonstration of this integration.
