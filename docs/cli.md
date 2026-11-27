# CLI設計 / CLI design

更新 / Updated: 2026-11-27。利用者の報告によると、教員の方針は「まずcommandlineとAPIで検証し、その後にUIを加える」です。CLIはその検証のための入口で、他の人が書いたprogramも同じAPIを呼びます。研究UIは検証の後に加えます（[UI設計](ui.md)）。以下は提案で、CLIのcodeはまだありません。未決の点は[構成設計](architecture.md#open-questions)に載せ、利用者が決めます。

According to the user, the supervisor's direction is to verify through the command line and API first and add a UI afterwards. The CLI is the entry point for that verification, and programs written by other people call the same API. The research UI is added after verification (see the [UI design](ui.md)). Everything below is a proposal; no CLI code exists yet. Open points are listed in the [architecture](architecture.md#open-questions) for the user to decide.

## serviceとCLIの関係 / How the service and CLI relate

後でUIを加える方針なので、HTTP serviceは残します。CLI・runner・他の人のprogram・後の研究UIは、すべて同じ `/v1` を呼びます。比べた構成は次の三つで、A案を使います。

Because a UI is added later, the HTTP service stays. The CLI, the runner, other people's programs and the later research UI all call the same `/v1`. Three layouts were compared, and option A is used.

| 選択肢 / Option | 内容 / What it means | 利点と負担 / Benefit and cost |
| --- | --- | --- |
| **A · HTTP service + CLI（採用） / HTTP service + CLI (used)** | 既存の `/v1` serviceを保ち、検証の間はCLIを人の入口にする。runner、他のprogram、後のUIも同じAPIを使う / Keep the existing `/v1` service with the CLI as the human entry point during verification. The runner, other programs and the later UI use the same API | 受入済みのL1–L3とL4の設計をそのまま使える。serverとCLIの二つを起動する / Accepted L1–L3 code and the L4 design carry over unchanged. Two things run: the server and the CLI |
| B · HTTP APIのみ / HTTP API only | CLIを作らず、`curl` などで直接呼ぶ / No CLI; call the API with `curl` or similar | codeは最少。request IDの生成、状態照会、errorの読取りを毎回手で行う / Least code. Request IDs, status lookups and error reading are done by hand every time |
| C · serviceなしのCLI / CLI without a service | 各commandがDBとmodelを直接操作する / Each command opens the database and calls the model itself | processが一つ減る。認証、所有者の隔離、session内の排他を作り直すことになり、後のUIと他のprogramが呼ぶAPIがなくなる / One process fewer. Authentication, owner isolation and per-session exclusion must be rebuilt, and the later UI and other programs would have no API to call |

## Command一覧 / Commands

binary名と構成は未決です。下表は一つの `bridge` binaryにsubcommandを置く案で書いています。各commandは[native API](contracts.md)の一つのrouteに対応し、CLIだけの機能は足しません。

The binary name and layout are undecided. The table assumes one `bridge` binary with subcommands. Each command maps to one route of the [native API](contracts.md); the CLI adds no behavior of its own.

| Command | API | 使える段階 / Available from |
| --- | --- | --- |
| `bridge serve` | serviceを起動。現在の引数なし起動に相当 / Start the service, as the bare binary does today | S1 |
| `bridge models` | `GET /v1/models` | S1 |
| `bridge session create --policy ID` | `POST /v1/sessions` | S1 |
| `bridge session list` | `GET /v1/sessions` | S1 |
| `bridge session show SESSION [--conversation]` | `GET /v1/sessions/{id}`、`?view=conversation` | S1 |
| `bridge send SESSION [--model ID] [TEXT]` | `POST /v1/sessions/{id}/messages` の後、turnを照会 / then turn lookups | S1 |
| `bridge turn show SESSION REQUEST_ID` | `GET /v1/sessions/{id}/turns/{request_id}` | S1 |
| `bridge turn cancel SESSION REQUEST_ID` | `POST …/turns/{request_id}/cancel` | S1 |
| `bridge turn reconcile SESSION REQUEST_ID` | `POST …/turns/{request_id}/reconcile` | S2 |
| `bridge memory list SESSION` | `GET /v1/sessions/{id}/memories` | S3 |
| `bridge session export SESSION --out FILE` | `GET /v1/sessions/{id}/export` | S3 |
| `bridge session delete SESSION [--yes]` | `DELETE /v1/sessions/{id}` の後、完了を照会 / then completion lookups | S3 |

合成例（提案の書式で、まだ実行できません） / Synthetic example in the proposed syntax, not runnable yet:

```sh
bridge session create --policy full-history-v1
bridge send "$S" --model local-model-a "The deadline is November 12."
bridge send "$S" --model local-model-b "Correction: the deadline is November 15."
bridge send "$S" --model local-model-a "What is the deadline?"
```

対話型の `chat` 画面（REPL）は最初は作りません。一回のcommandが一つのturnに対応する方が、scriptとrunnerからも同じ形で使えます。必要が出た時に検討します。

No interactive `chat` loop (REPL) at first. One command per turn keeps the same shape for scripts and the runner. Reconsider when a concrete need appears.

<a id="next-model"></a>

## 次のモデル / Which model comes next

2026-11-18に利用者は「同じ会話では直前に使ったモデルを引き継ぐ」案を選びました。CLIでも後の[研究UI](ui.md#next-model)でも同じ決定を使います。`--model` を省いた `send` は、そのsessionの最新の完了回答で使われたモデルを使います。完了回答がない場合と、そのモデルが許可一覧から外れた場合は送信せず、`--model` の指定を求めます。既定のモデルを黙って選ぶことはしません。

On 2026-11-18 the user chose to keep the last model used in the same conversation. The CLI and the later [research UI](ui.md#next-model) both follow that decision. A `send` without `--model` uses the model of the most recent completed reply in that session. With no completed reply, or when that model is no longer allowlisted, the CLI sends nothing and asks for `--model`. It never picks a default silently.

この値はserverの会話記録から毎回求め、手元のfileには保存しません。使ったモデルは回答と一緒に表示します。Bで訂正した直後に `--model` なしで送ればB、Aに戻すには `--model` でAを一度指定します。runnerと評価scriptは常に `--model` を明示します。

The value is derived from the server's conversation record each time and is not stored in a local file. The model used is printed with the reply. After a correction on B, a send without `--model` goes to B; switch back by naming A once with `--model`. The runner and evaluation scripts always pass `--model` explicitly.

実行モデルが要求と違う証拠がある場合は「要求 B・実行 A」、確認できない場合は「要求 B・実行未確認」と表示し、次の送信では `--model` の指定を求めます。

When evidence shows a different executed model, print "requested B · executed A"; when it cannot be verified, print "requested B · execution unverified". In both cases the next send requires `--model`.

## 送信と待機 / Send and wait

`send` はrequest IDを生成し、送信前にstderrへ表示します。202を受けたらturnを照会し、終端状態で回答を表示して終了します。`--no-wait` は受理の確認だけで戻ります。`--request-id` で同じIDを指定すると再送になります。

`send` generates a request ID and prints it to stderr before sending. After the 202 it looks the turn up until a terminal state, prints the reply and exits. `--no-wait` returns once admission is confirmed. Passing the same ID with `--request-id` makes the call a replay.

| 場面 / Situation | CLIの動作 / CLI behavior |
| --- | --- |
| 待機中のCtrl+C、待機の期限切れ / Ctrl+C or wait timeout | 待機だけをやめる。turnはserverで続き、照会するcommandを表示する。取消は `turn cancel` で明示する / Stop waiting only. The turn continues on the server, and the CLI prints the lookup command. Cancellation is the explicit `turn cancel` |
| 送信の応答が届かない / Submission response lost | 自動で再送しない。同じ `--request-id` と同じ本文で再送すれば同じturnが返る / No automatic resend. Replaying with the same `--request-id` and text returns the same turn |
| request IDも失った / Request ID lost too | `session show --conversation` の一覧から受理済みturnを見つける / Find the admitted turn in the `session show --conversation` list |
| 409 busy | 進行中のturnのIDを示して終了。待ち行列も自動再試行も作らない / Report the active turn's ID and exit. No queue, no automatic retry |
| 失敗・取消の後の再試行 / Retry after failure or cancellation | 新しいIDでの明示的な送信。失敗した発言は次のモデルの入力に入らないと表示する / An explicit send under a fresh ID. State that the failed message is excluded from the next model's input |

## 状態と表示 / States and display

状態は[契約](contracts.md#turn-states)の名前をそのまま表示し、serverの応答から読みます。待ち時間から推測しません。未知の状態や非互換のschemaを受け取ったら、説明を出して送信を止めます。

Print the state names from the [contract](contracts.md#turn-states) unchanged and read them from the server. Never infer state from elapsed time. On an unknown state or incompatible schema, explain and refuse to send.

| Turn state | 表示 / Display | `send` の続き / What `send` does next |
| --- | --- | --- |
| `accepted`, `generating` | 状態名を表示 / Print the state name | 照会を続ける / Keep looking up |
| `answer_saved` | 保存済みの回答と「完了確認中」 / Saved reply plus "checking completion" | 確定状態まで照会 / Look up until a confirmed terminal state |
| `memory_pending` | 回答と「記憶更新未完了」 / Reply plus "memory update pending" | 回答を表示し、未解決として終了。次の送信はbusy / Print the reply and exit as unresolved. The next send is busy |
| `complete` | 回答、実行モデル、request ID / Reply, executed model, request ID | 成功で終了 / Exit successfully |
| `failed`, `cancelled` | 状態と理由code / State and reason code | 失敗で終了。自動再送しない / Exit with failure. No automatic resend |
| `recovery_required` | 「結果が未確定」と `turn reconcile` の案内 / "Outcome uncertain" plus a pointer to `turn reconcile` | 未解決として終了 / Exit as unresolved |

`deleting` と `deleted` はsessionの状態です。`session delete` は開始後に完了を照会し、`deleted` を確認してから成功で終了します。

`deleting` and `deleted` are session states. `session delete` looks up completion after starting and exits successfully only once `deleted` is confirmed.

## 出力と終了code / Output and exit codes

stdoutには結果だけを出します。進行状況、request ID、errorはstderrです。`--json` は一回のcommandにつき一つのJSON objectを出し、scriptはこちらを読みます。人向けの表示は色だけに意味を持たせず、端末でない出力先には進行表示を出しません。

stdout carries only the result. Progress, request IDs and errors go to stderr. `--json` prints one JSON object per command, and scripts read that form. Human output never relies on color alone, and progress is suppressed when the output is not a terminal.

| 終了code案 / Proposed exit code | 意味 / Meaning |
| --- | --- |
| 0 | 成功。`send` ではturnが `complete` / Success. For `send`, the turn is `complete` |
| 1 | turnが `failed` または `cancelled` / The turn is `failed` or `cancelled` |
| 2 | 使い方の誤り、不正な入力（422） / Usage error or invalid input (422) |
| 3 | 認証失敗、対象なし（401、404） / Authentication failure or no such object (401, 404) |
| 4 | 競合、busy（409） / Conflict or busy (409) |
| 5 | serviceに届かない、依存未準備（接続失敗、503） / Service unreachable or dependency not ready (connection failure, 503) |
| 6 | 未解決。待機の中断・期限、`memory_pending`、`recovery_required` / Unresolved: interrupted or timed-out wait, `memory_pending`, `recovery_required` |
| 7 | 非互換の応答、未知の状態 / Incompatible response or unknown state |

errorは[error envelope](contracts.md#turn-states)の `code` をそのまま表示し、短い英語の説明を添える案です。`--json` は言語に依存しません。日英の説明は本書に置きます。

Errors print the `code` from the [error envelope](contracts.md#turn-states) unchanged with a short English explanation. `--json` is language-neutral. The Japanese and English explanations live in this document.

モデルの回答と記憶のtextは信頼できない入力です。端末へ出す時は制御文字をescapeし、端末の操作列として解釈させません。`--json` とfile出力は元のtextを保ちます。

Model replies and memory text are untrusted input. When writing to a terminal, escape control characters so they are not interpreted as terminal sequences. `--json` and file output keep the original text.

## 設定と秘密 / Configuration and secrets

接続先は `BRIDGE_URL`（既定は `http://127.0.0.1:8080`）、tokenは `BRIDGE_TOKEN` から読みます。tokenをflagで受け取りません。引数はprocess一覧とshell履歴に残るためです。tokenを含む設定fileも最初は作りません。接続先は最初はloopbackに限ります。

Read the address from `BRIDGE_URL` (default `http://127.0.0.1:8080`) and the token from `BRIDGE_TOKEN`. The token is never a flag, because arguments appear in process listings and shell history. No token-bearing config file at first. The address is limited to loopback initially.

本文は引数、`--file`、stdinのいずれかで渡します。引数はshell履歴に残るため、合成data以外はfileかstdinを使います。日本語の本文が引数・stdin・pipeで壊れないことをWindowsで確かめます。

Message text comes from an argument, `--file` or stdin. Arguments stay in shell history, so anything other than synthetic data goes through a file or stdin. Verify on Windows that Japanese text survives arguments, stdin and pipes.

`session export` は `--out` のfileへ書き、既存fileを上書きしません。exportは会話本文を含むのでcommitしません。`session delete` は端末では対象sessionと範囲を示して確認を求め、端末でない場合は `--yes` がなければ実行しません。

`session export` writes to the `--out` file and refuses to overwrite an existing one. Exports contain conversation text and are never committed. `session delete` asks for confirmation on a terminal, naming the session and scope, and refuses to run non-interactively without `--yes`.

## 配置と依存 / Layout and dependencies

```text
cmd/bridge/         serve とCLI subcommandの入口 / entry point for serve and the CLI subcommands
internal/client/    /v1 のHTTP client、DTO、応答の検証、error envelope / HTTP client for /v1, DTOs, response validation and the error envelope
cmd/evaluate/       runner（S5）。internal/client を使う / runner (S5), built on internal/client
```

予定の配置で、空directoryを先に作りません。CLIとrunnerは同じclientを使い、DBやworkerへ直接触れません。clientはserver側のpackageをimportせず、serverと同じ契約fixtureで検証します。subcommandは標準の `flag` で作り、CLI frameworkの依存を足しません。

This is a target layout; do not create empty directories in advance. The CLI and runner share one client and never touch the database or worker directly. The client imports no server-side package and is verified against the same contract fixtures as the server. Subcommands use the standard `flag` package, adding no CLI framework dependency.

各commandは `run(args, stdin, stdout, stderr)` が終了codeを返す形にします。testはbinaryをbuildせずに、`httptest` 上の実handlerと合成 `ModelClient` に対して同じ関数を呼べます。

Each command is a `run(args, stdin, stdout, stderr)` function returning an exit code. Tests can call it against the real handler and a synthetic `ModelClient` under `httptest`, without building a binary.

## 検証 / Verification

将来の受入方針で、実行済みのtestではありません。 / These are future acceptance principles, not tests already executed.

| 接点 / Seam | 確認すること / What it checks |
| --- | --- |
| client | 契約fixtureでversion・欠落field・未知の状態。`httptest` でrequest ID、errorの対応、自動再送がないこと / Versions, missing fields and unknown states with contract fixtures. Request IDs, error mapping and the absence of automatic resends under `httptest` |
| command | 終了code、stdoutとstderrの分離、`--json` の形、tokenが出力に現れないこと / Exit codes, stdout/stderr separation, the `--json` shape and no token in any output |
| 流れ / Flow | 同じsessionのA → B → A、二重送信、応答喪失後の再送、Ctrl+C後の続行、busy、取消 / A → B → A in one session, double submit, replay after a lost response, continuation after Ctrl+C, busy and cancellation |
| binary（S2以降 / from S2） | buildしたbinaryを合成のmodel serverに対して起動し、同じ流れを通す / Launch the built binary against a synthetic model server and run the same flow |

## 研究UIとの関係 / Relation to the research UI

CLIで確かめた流れ（切替、二重送信、応答喪失後の回復、取消、出典、export、削除）は、後の研究UIがそのまま使うAPIの動作です。UIは同じ契約fixtureを使い、CLIの検証で固まった状態名・error code・回復の手順を変えません。UIで新しく必要になるのは、画面の状態管理と[browserでの安全対策](ui.md#browser-security)です。

The flows verified through the CLI (switching, double submit, recovery after a lost response, cancellation, sources, export and deletion) are the API behavior the later research UI relies on. The UI uses the same contract fixtures and changes none of the state names, error codes or recovery steps settled during CLI verification. What the UI adds is view-state management and [browser security](ui.md#browser-security).
