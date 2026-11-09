# セキュリティと機密情報 / Security and sensitive data

現在は研究原型の準備段階です。hosted LLM service、認証system、model API、React研究UIはまだありません。既存static siteは固定の合成textだけを使用し、modelを呼び出しません。

The project is preparing a research prototype. It does not yet provide a hosted LLM service, authentication system, model API or React research UI. The existing static site uses fixed synthetic text and makes no model calls.

## 報告 / Reporting

credential、私的会話、個人情報、攻撃可能な詳細をpublic issueに書かないでください。GitHub private vulnerability reportingが有効ならSecurity → Report a vulnerabilityを使います。有効でなければ機密を含めずprivate連絡方法を求めるissueを作成します。応答時間は保証しません。

Do not post credentials, private conversations, personal information or exploitable details in public issues. If GitHub private vulnerability reporting is enabled, use Security → Report a vulnerability. Otherwise, request a private reporting channel without including sensitive details. No response-time guarantee is offered.

## 原型の境界 / Prototype boundaries

認証、所有者確認、local-only provider設定、削除検証は[設計](docs/architecture.md)と[契約](docs/contracts.md)の要件であり、実装済みの保護とは扱いません。loopbackは同一hostの他processを隔離しません。raw model API、worker、DB fileも保護対象です。

Authentication, ownership checks, local-only provider configuration and deletion verification are design/contract requirements, not implemented protections. Loopback does not isolate other processes on a host. Raw model APIs, workers and DB files also require protection.

UIのbuild変数、localStorage、console logに秘密を保存しません。model出力と検索memoryは信頼できないtextです。export、backup、WAL、test trace、browser screenshotも本文を含み得るため、初期検証は合成dataに限定します。

Do not store secrets in UI build variables, localStorage or console logs. Model output and retrieved memory are untrusted text. Exports, backups, WAL, test traces and browser screenshots can contain content; initial checks use synthetic data only.

credentialを誤公開した場合、発行元で失効させます。最新fileから消してもGit履歴からは消えません。削除後のbackup restoreでdataを復活させないよう、削除記録と保持範囲を管理します。ignore規則はstageと履歴のreviewを代替しません。

If a credential is exposed, revoke it with its issuer. Removing it from the latest file does not erase Git history. Track deletion records and retention so backup restores do not resurrect deleted data. Ignore rules do not replace staged-content/history review.
