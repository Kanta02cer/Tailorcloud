#!/bin/bash

# TailorCloud Flutter App セットアップスクリプト

set -e

echo "🚀 TailorCloud Flutter App セットアップを開始します..."
echo ""

# Flutterのバージョンを確認
if ! command -v flutter &> /dev/null; then
    echo "❌ Flutterがインストールされていません"
    echo "   インストール方法: brew install --cask flutter"
    echo "   または: https://docs.flutter.dev/get-started/install/macos"
    exit 1
fi

echo "✅ Flutterが見つかりました"
flutter --version
echo ""

# Flutter Doctorで環境を確認
echo "📋 Flutter環境を確認中..."
flutter doctor
echo ""

# 依存パッケージをインストール
echo "📦 依存パッケージをインストール中..."
flutter pub get
echo ""

# コード生成を実行
echo "🔨 コード生成を実行中..."
echo "   (Freezed, Riverpod Generator, JSON Serializable)"
flutter pub run build_runner build --delete-conflicting-outputs
echo ""

echo "✅ セットアップが完了しました！"
echo ""
echo "次のステップ:"
echo "  1. iOSシミュレーターを起動: open -a Simulator"
echo "  2. アプリを実行: flutter run"
echo "  3. または、利用可能なデバイスを確認: flutter devices"

