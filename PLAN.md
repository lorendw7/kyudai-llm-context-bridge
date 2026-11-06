# 作業計画 / Plan

更新：2026-11-06（Asia/Tokyo）。本書は唯一の作業一覧です。

## 現在の方針 / Current decisions

- A案：Goブリッジ + Ollama + Mem0 OSS。原型作成の指導教員承認、玄界アカウント・研究室計算枠の保有は利用者の報告による。
- 原型コードは利用者本人がすべて実装する。assistantは会話内で中国語・英語の説明、課題、レビューを提供し、明示的な変更依頼なしに代筆しない。文書とサイトは日英。
- 複数研究室が共有するのは推論資源。会話・記憶・権限は分離する。初回は合成データとローカル2モデルを使用し、商用APIは別枠。
- ソフトウェア・必須依存はOSS。正確な配布物、モデル、依存ライセンスとhashを確認する。CUDA等の玄界GPU基盤の適用範囲は未確認。導入、大型取得、遠隔実行は別途許可された範囲で行う。
- 公開対象は現在の設計、コード、検証、日英サイト。旧演示稿・生成source・作業履歴はローカル専用とし、GitHubに含めない。mainは今回の依頼により履歴を置き換え、現在のsnapshotから再開する。

## 現在の状態 / Status

- [x] localhost限定の`GET /health`とGoテストを用意。モデル、記憶、認証、永続化の実装と測定は未着手。
- [x] A案の責務、会話・記憶の分離、API境界、失敗処理、玄界配置、評価案をREADMEに整理。
- [x] 本人実装・会話内指導の役割をAGENTSとCONTRIBUTINGに反映。
- [x] 公開用snapshotを20ファイルに整理。原型コードを維持し、演示資料を追跡対象から除外。文書・機密情報検査、Goテスト、JavaScript構文、日英サイトの各段階・keyboard・広幅／狭幅表示を確認。

公開先は既存GitHub Pagesです。実際の配信状態は[GitHub Actions](https://github.com/lorendw7/kyudai-llm-context-bridge/actions)で確認できます。旧演示資料はローカル専用です。

## 実装と受入基準 / Implementation and acceptance

- [ ] 契約・合成ケース：`Session`、`Message`、model/memory interfaceとversioned API schemaを本人が作成。所有権、訂正、出典、削除、失敗の期待入出力を説明できること。
- [ ] 依存不要のGo実装：fake adapterでA → B → A、他利用者拒否、request ID、timeoutをテスト。モデル品質の測定とは扱わない。
- [ ] 依存・artifact固定：必要なpackagesのみ解決し、推移依存、LICENSE/NOTICE、model/tokenizer/template/量子化、hashを記録。未解決の非OSS必須依存がないこと。GPU基盤が要件外ならCPU案を利用。
- [ ] ローカル連携：許可範囲内でOllama、Mem0、local embedding、Qdrant local modeを接続。telemetry/cloud設定を確認し、外向き通信遮断下で訂正保持と他利用者隔離を確認。
- [ ] 保存と回復：再起動、重複要求、回答後の記憶更新失敗、削除と更新の競合、出典、exportをテスト。全派生storeとcheckpointに削除範囲が及ぶこと。
- [ ] 比較対照：全履歴・逐次要約・Mem0を同じモデル・ケースで実行。回答入力tokenと抽出・更新・要約・embedding・検索・retry・model loadを分けて記録。

## 玄界と評価 / Genkai and evaluation

- [ ] 環境確認：接続範囲の確認後、GPU権限、残枠と計上方式、module、OS、driver、Singularity、private storage、サービス隔離を記録。実際の環境に合う固定artifactを選ぶ。
- [ ] 包装：同一ジョブ内にbridge、memory、Ollama、runnerを配置。起動確認、timeout、途中保存、終了処理を用意。runtime install/downloadが不要で、他ジョブに干渉しないこと。
- [ ] Smoke test：B-group 1 GPU、`b-inter`最大30分を提案。起動、GPU認識、応答、記憶保存、正常終了を確認し失敗も記録。
- [ ] 開発pilot：1 GPU、`b-batch`最大60分を提案。開発12ケースの範囲で品質、token、VRAM、処理時間と消費枠を測定。上限で停止し欠測を残す。本評価の枠は実測から算出。
- [ ] 評価固定：開発12／最終36の独立した日英ケース、3方式×3反復、24 turns・6 probes、切替時点、rubric、overflow、token集計、seed、hashを固定。初期context案はREADMEを参照。最終ケースを見て調整しない。
- [ ] 最終評価：324 conversation runs（提案値）と独立contract tests。品質90%以上・全履歴比5 percentage points以内の低下・回答入力token30%以上削減を提案し、評価前に確定。隔離／削除は独立必須基準とし、言語・family・model別に失敗と不確実性を報告。
- [ ] 報告：確認済み要件、設計、vendor claim、実測を区別した日英報告。raw output、重み、DB、私的情報は公開しない。未達も研究結果として保持。

同時利用、学内SSO、保存期間の実運用規則、商用API、遠隔ブラウザdemoは必要が生じた段階で別途確定します。llama.cpp／vLLMやHy-Memoryへの変更は、A案の実測で必要性が判明した場合の別実験です。

## 工数見積り / Effort estimate

本人が学習しながら実装する場合の計画値。週20時間を前提に、最初の1週間の実績で更新します。

| 範囲 / Scope | 追加工数 / Incremental effort |
| --- | --- |
| 契約・合成テスト / Contracts and synthetic tests | 16–24 h |
| ローカル切替・記憶 / Local handoff and memory | 20–32 h |
| 保存・回復・比較 / Persistence, recovery and controls | 20–36 h |
| 玄界・評価・報告 / Genkai, evaluation and report | 24–48 h |

合計80–140時間、約4–7週間。最初のローカルdemoは累積36–56時間、約2–3週間。監査、互換性、手戻り、申請・queue待ちで延びる可能性があり、納期保証ではありません。常時運用サービス、SSO、商用API、fine-tuningは含みません。
