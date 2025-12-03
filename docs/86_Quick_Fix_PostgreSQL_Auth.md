# PostgreSQL認証エラーの即座に解決方法

**作成日**: 2025-01  
**問題**: `FATAL: password authentication failed for user "wantan"`

---

## 🔍 エラー内容

```
FATAL: password authentication failed for user "wantan"
```

macOSユーザー名 "wantan" でPostgreSQLに接続しようとしていますが、パスワード認証に失敗しています。

---

## ✅ 即座に解決する方法

### 方法1: 認証設定を変更（推奨・最も簡単）

認証方式を "trust" に変更して、パスワードなしで接続できるようにします。

```bash
# 自動修正スクリプトを実行
./scripts/fix_postgresql_auth.sh
```

このスクリプトは:
1. pg_hba.confのバックアップを作成
2. 認証方式を "trust" に変更
3. PostgreSQLを再起動
4. 接続テストを実行

---

### 方法2: 手動で認証設定を変更

```bash
# PostgreSQL PATH設定
source ./scripts/fix_postgresql_path.sh

# pg_hba.confファイルを編集
nano /opt/homebrew/var/postgresql@15/pg_hba.conf
# または
nano /usr/local/var/postgresql@15/pg_hba.conf

# ファイルの最後に以下を追加:
# local   all             all                                     trust
# host    all             all             127.0.0.1/32            trust
# host    all             all             ::1/128                 trust

# PostgreSQLを再起動
brew services restart postgresql@15
```

---

### 方法3: postgresユーザーで接続

PostgreSQLのデフォルトスーパーユーザー "postgres" を使用します。

```bash
# postgresユーザーで接続を試す
psql -U postgres -d postgres

# または、macOSユーザー名がpostgresの場合
psql -d postgres
```

---

### 方法4: ユーザーにパスワードを設定

```bash
# まず、何らかの方法でPostgreSQLに接続（方法1-3のいずれか）
psql -d postgres

# ユーザーにパスワードを設定
ALTER USER wantan WITH PASSWORD 'your_password';

# または、新しいユーザーを作成
CREATE USER tailorcloud WITH PASSWORD 'your_password';
CREATE DATABASE tailorcloud OWNER tailorcloud;
```

---

## 🎯 最も簡単な解決方法（推奨）

```bash
# 1. PostgreSQL PATH設定
source ./scripts/fix_postgresql_path.sh

# 2. 認証設定を自動修正
./scripts/fix_postgresql_auth.sh

# 3. データベースを作成
createdb tailorcloud

# 4. 接続確認
psql tailorcloud -c "SELECT 1;"
```

---

## 📝 注意事項

- **開発環境のみ**: "trust" 認証方式は開発環境でのみ使用してください
- **本番環境**: 本番環境では必ずパスワード認証を使用してください
- **セキュリティ**: ローカル接続のみで "trust" を使用しているため、リモート接続は保護されています

---

## 🔧 トラブルシューティング

### 問題: スクリプトが実行できない

```bash
# 実行権限を付与
chmod +x ./scripts/fix_postgresql_auth.sh
```

### 問題: PostgreSQLが再起動しない

```bash
# 手動で停止・起動
brew services stop postgresql@15
brew services start postgresql@15
```

### 問題: pg_hba.confが見つからない

```bash
# PostgreSQLデータディレクトリを確認
ls -la /opt/homebrew/var/postgresql@15/
# または
ls -la /usr/local/var/postgresql@15/
```

---

**最終更新日**: 2025-01

