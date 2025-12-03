# ディスク容量不足エラー - 対処方法

**作成日**: 2025-01  
**問題**: Flutterインストール時に "No space left on device" エラー

---

## 🔴 問題の状況

### 現在のディスク容量

- **使用率**: 82% (185GB / 228GB)
- **空き容量**: 約3.1GB
- **必要な容量**: Flutterインストールに約2GB必要

### エラーメッセージ

```
symlink error: No space left on device
write error (disk full?)
```

---

## 💡 解決方法

### 方法1: ディスク容量を確保する（推奨）

#### 1. 大きなファイルを削除

```bash
# 大きなファイルを検索（10MB以上）
find ~ -type f -size +10M -exec ls -lh {} \; 2>/dev/null | sort -k5 -hr | head -20

# 一時ファイルを削除
rm -rf ~/Library/Caches/*
rm -rf /tmp/*

# Homebrewのキャッシュを削除
brew cleanup --prune=all
```

#### 2. 不要なアプリケーションを削除

```bash
# 大きなアプリケーションを検索
du -sh /Applications/* 2>/dev/null | sort -hr | head -20

# 不要なアプリを削除（手動で実施）
```

#### 3. Docker（使用している場合）のクリーンアップ

```bash
# Dockerイメージ・コンテナを削除
docker system prune -a --volumes
```

#### 4. Xcodeの不要なデータを削除

```bash
# Xcodeの派生データを削除
rm -rf ~/Library/Developer/Xcode/DerivedData/*

# Xcodeのアーカイブを削除
rm -rf ~/Library/Developer/Xcode/Archives/*
```

#### 5. 容量を確認

```bash
# 現在のディスク使用状況を確認
df -h /
```

**目標**: 空き容量を5GB以上確保する

---

### 方法2: 外部ストレージにFlutterをインストール

容量が確保できない場合、外部ディスクや別の場所にFlutterをインストール：

```bash
# 外部ディスクをマウント（例）
# 外部ディスクにFlutterをインストール
cd /Volumes/ExternalDisk
git clone https://github.com/flutter/flutter.git -b stable
export PATH="$PATH:/Volumes/ExternalDisk/flutter/bin"
```

---

### 方法3: 一時的に他のファイルを削除

```bash
# ダウンロードフォルダの大きなファイルを確認
du -sh ~/Downloads/* 2>/dev/null | sort -hr | head -10

# 不要なファイルを削除
# （注意: 削除前に内容を確認してください）
```

---

## 🔍 ディスク容量の確認コマンド

### ディスク全体の使用状況

```bash
df -h /
```

### 大きなディレクトリを検索

```bash
# ホームディレクトリ配下で大きなディレクトリを検索
du -sh ~/* 2>/dev/null | sort -hr | head -20
```

### 特定のディレクトリの容量

```bash
# Libraryフォルダの容量
du -sh ~/Library

# Downloadsフォルダの容量
du -sh ~/Downloads

# Applicationsフォルダの容量
du -sh /Applications
```

---

## ✅ 推奨アクション

### Step 1: すぐにできるクリーンアップ

```bash
# 1. Homebrewのキャッシュ削除
brew cleanup --prune=all

# 2. 一時ファイル削除
rm -rf /tmp/*
rm -rf ~/Library/Caches/*

# 3. ディスク容量を再確認
df -h /
```

### Step 2: 容量が確保できたら

```bash
# Flutterを再インストール
brew install --cask flutter
```

### Step 3: または、手動でインストール

容量を確保した後、公式サイトからダウンロード：

```bash
# 1. 公式サイトからダウンロード
# https://docs.flutter.dev/get-started/install/macos

# 2. 解凍
cd ~
unzip ~/Downloads/flutter_macos_*.zip

# 3. PATHに追加
echo 'export PATH="$PATH:$HOME/flutter/bin"' >> ~/.zshrc
source ~/.zshrc
```

---

## 📊 目標

- **空き容量**: 5GB以上を確保
- **使用率**: 80%以下にする

---

## ⚠️ 注意事項

- ファイルを削除する前に、**バックアップを取る**ことを推奨します
- 重要なファイルを誤って削除しないよう注意してください
- `rm -rf` コマンドは慎重に使用してください

---

**最終更新日**: 2025-01

