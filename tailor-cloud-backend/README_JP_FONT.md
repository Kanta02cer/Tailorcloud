# 日本語フォント設定ガイド

TailorCloudの下請法PDF生成機能では、日本語フォント（Noto Sans JP）を使用してPDFを生成します。

## フォントファイルの配置

### 方法1: 環境変数でフォントディレクトリを指定

```bash
export FONT_DIR="/path/to/fonts"
```

### 方法2: デフォルトディレクトリに配置

プロジェクトルートの `assets/fonts/` ディレクトリに配置：

```
tailor-cloud-backend/
  assets/
    fonts/
      NotoSansJP-Regular.ttf
      NotoSansJP-Bold.ttf
```

## フォントファイルのダウンロード

1. [Google Fonts - Noto Sans JP](https://fonts.google.com/noto/specimen/Noto+Sans+JP) からダウンロード
2. または [Noto CJK](https://github.com/googlefonts/noto-cjk) からダウンロード

必要なファイル：
- `NotoSansJP-Regular.ttf`（標準フォント）
- `NotoSansJP-Bold.ttf`（太字フォント、オプション）

## フォントがない場合

フォントファイルが存在しない場合、システムは自動的にArialフォントにフォールバックします。警告メッセージが表示されますが、PDF生成は正常に動作します（日本語は正しく表示されません）。

## Docker環境での設定

Dockerコンテナ内でフォントを使用する場合：

```dockerfile
# Dockerfileに追加
COPY assets/fonts/ /app/assets/fonts/
ENV FONT_DIR=/app/assets/fonts
```

## 注意事項

- フォントファイルのライセンスに注意してください（Noto Sans JPはApache License 2.0）
- 本番環境ではフォントファイルをバージョン管理に含めるか、適切なストレージから取得するように設定してください

