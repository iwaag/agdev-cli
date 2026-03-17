# agdev-cli Plan

## Goal

このプロジェクトを、AIエージェントからサブプロセス実行されることを前提にした Go 製 CLI の土台まで整える。
実処理はまだ入れず、以下が明確になっている状態を最初の到達点とする。

- コマンドツリーの構造
- 出力方針
- 設定と認証の流し方
- API / Socket 連携を置く場所
- Docker 前提のビルドと実行方法
- 今後の機能追加で破綻しないディレクトリ構成

Gemini の提案は妥当で、この用途では最初から Cobra を入れて始める方がよい。
理由は、`module command action [args...]` 型の階層コマンド、必須引数の検証、`--json` のような共通フラグ、ヘルプ生成を早い段階で整理できるため。

## Project Direction

### 想定ユースケース

- CLI の主利用者は人間ではなく AI エージェント
- CLI は長時間 AI 処理を持つバックエンドへのブリッジ
- 実行環境はコンテナ前提
- JWT 等の認証情報はログイン処理ではなく環境変数から読み込む
- 結果は stdout で機械可読に返し、エラーは stderr に分離する

### 非目標

- 初期段階では業務ロジックを入れない
- 初期段階では完全な Socket.io 実装までは行わない
- 初期段階では永続設定ファイル管理は必須にしない

## Recommended Initial Dependencies

最初に入れる依存は絞る。

### Required

- `github.com/spf13/cobra`
  - コマンドツリー、引数検証、ヘルプ生成
- `github.com/spf13/viper`
  - 環境変数ベース設定の読み込み

### Likely Needed Soon

- `github.com/go-resty/resty/v2`
  - HTTP クライアントを少ないボイラープレートで扱うため
- `github.com/stretchr/testify`
  - テスト記述の簡素化

### Defer Until Actually Needed

- Socket.io クライアント実装ライブラリ
  - 実際のバックエンド仕様が固まってから選定する
- 構造化ロガー
  - まずは標準ライブラリ + 出力ルールで十分

## Recommended Project Layout

初期段階では大きくしすぎず、CLI 層と内部ロジック層を分ける。

```text
.
├── cmd/
│   ├── root.go
│   ├── version.go
│   ├── image.go
│   ├── image_generate.go
│   ├── video.go
│   └── video_generate.go
├── internal/
│   ├── app/
│   │   └── exitcode.go
│   ├── config/
│   │   └── config.go
│   ├── output/
│   │   └── output.go
│   ├── api/
│   │   ├── client.go
│   │   └── types.go
│   └── socket/
│       └── client.go
├── .dockerignore
├── Dockerfile
├── main.go
├── go.mod
└── README.md
```

### Layout Rationale

- `cmd/`
  - Cobra のコマンド定義だけを置く
  - 引数解釈、フラグ定義、最低限のバリデーションに留める
- `internal/config/`
  - 環境変数から API URL、JWT、タイムアウトなどを読む
- `internal/api/`
  - REST 呼び出しの共通クライアント
- `internal/socket/`
  - 将来の進捗受信や双方向イベント用
- `internal/output/`
  - `--json` / テキスト出力の切り替え
- `internal/app/exitcode.go`
  - エージェントにとって重要な終了コード定義

## Command Model

最終イメージに近い形で、最初から階層コマンドを切る。

```text
agdev-cli
├── image
│   └── generate
├── video
│   └── generate
└── version
```

### Initial Command Policy

- `image generate <input-image> <prompt>`
- `video generate <first-frame> <last-frame>`
- まずは実処理なしで、引数検証とダミーのレスポンスだけ返せる形まで作る
- 各 `generate` コマンドは Cobra の `Args: cobra.ExactArgs(...)` を使う

## Output Contract

このプロジェクトでは人間向け CLI より、サブプロセスとしての安定性を優先する。

### stdout

- 正常結果のみ出す
- `--json` 指定時は必ず JSON を返す
- 将来的な標準レスポンス形の例

```json
{
  "status": "accepted",
  "job_id": "video-123",
  "message": "request queued"
}
```

### stderr

- エラー詳細
- デバッグメッセージ
- バックエンド接続失敗や認証不足の説明

### Exit Codes

初期段階から定義だけ固定しておく。

- `0`: success
- `1`: usage / validation error
- `2`: auth error
- `3`: network error
- `4`: timeout
- `5`: backend error
- `10`: unexpected internal error

## Config Policy

設定はまず環境変数のみで扱う。
設定ファイルやプロファイル概念は後回しでよい。

### Proposed Environment Variables

- `AGDEV_API_BASE_URL`
- `AGDEV_AUTH_TOKEN`
- `AGDEV_REQUEST_TIMEOUT`
- `AGDEV_OUTPUT_JSON`
- `AGDEV_LOG_LEVEL`

### Config Rules

- Viper で環境変数を自動読込
- 必須設定は起動時に検証
- JWT は CLI 内で保存しない
- トークン値はログ出力しない

## Docker Strategy

Gemini の提案どおり、マルチステージビルドを採用する。

### Build Policy

- ビルド用に `golang:1.22` 系イメージを使用
- 実行用は軽量イメージを使う
- 配布物は単一バイナリ

### Runtime Policy

- コンテナ起動時に環境変数を注入
- `ENTRYPOINT` で CLI を直接実行可能にする
- エージェントが `docker run ... image generate ...` のように扱える構成にする

## Implementation Phases

### Phase 1: Skeleton

目的:
CLI の骨組みとビルド経路だけを作る

作業:

- `main.go` を追加
- Cobra root command を追加
- `image`, `video`, `version` コマンドを追加
- `image generate`, `video generate` を追加
- `--json` のグローバルフラグを追加
- Dockerfile を追加
- README に最小実行例を追加

完了条件:

- ローカルで `go build ./...` が通る
- `agdev-cli image generate a.png "prompt"` がダミー成功を返す
- `agdev-cli video generate first.png last.png --json` が JSON を返す
- Docker build が通る

### Phase 2: Shared Infrastructure

目的:
本処理を入れる前に、共通基盤を固める

作業:

- `internal/config` を追加
- `internal/output` を追加
- `internal/app/exitcode.go` を追加
- `internal/api/client.go` を追加
- 共通のエラー変換ルールを追加

完了条件:

- 環境変数不足時に明示的エラーになる
- stdout / stderr / exit code の責務分離ができている
- ダミー API クライアント差し込みができる

### Phase 3: Backend Integration Stub

目的:
本番 API 実装前の結線ポイントを作る

作業:

- REST リクエスト生成部を追加
- JWT ヘッダ注入を追加
- タイムアウトとコンテキスト制御を追加
- 将来の Socket クライアント差し込み用 interface を用意

完了条件:

- モック API に対してコマンドから呼び出せる
- キャンセル時に context が伝播する
- エラー種別が exit code に落ちる

## Concrete Initial File Plan

最初の実装では以下を作るのが適切。

- `main.go`
- `cmd/root.go`
- `cmd/version.go`
- `cmd/image.go`
- `cmd/image_generate.go`
- `cmd/video.go`
- `cmd/video_generate.go`
- `internal/config/config.go`
- `internal/output/output.go`
- `internal/app/exitcode.go`
- `Dockerfile`
- `.dockerignore`
- `README.md`

## Notes For Future Implementation

### Why Cobra From Day One

- この CLI は単発コマンドではなく、明確に階層コマンド化される
- 位置引数とフラグの混在が前提
- `--help` と `--json` を早期に揃えたい
- 後から標準ライブラリ実装を Cobra に載せ替えるより、最初から寄せた方が安い

### Why Keep Business Logic Out of `cmd/`

- CLI 層に API 呼び出しや整形ロジックを入れるとテストが壊れやすい
- 今後 `image upscale`, `video status`, `chat start` などが増えても整理しやすい

### Why JSON Mode Early

- エージェント実行では人間向け文章より構造化出力が重要
- 早い段階で出力契約を決めておくと後からの破壊的変更を減らせる

## Proposed Immediate Next Step

この計画に基づいて、次は Phase 1 の Skeleton を実際に作る。
つまり、Cobra ベースのコマンド雛形、`--json` の共通フラグ、Dockerfile、ダミー応答までをこのリポジトリに実装する。
