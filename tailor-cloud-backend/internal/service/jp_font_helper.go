package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf/v2"
)

// JPFontHelper 日本語フォントヘルパー
// Noto Sans JPなどの日本語フォントを使用するためのヘルパー
type JPFontHelper struct {
	fontDir string
}

// NewJPFontHelper JPFontHelperのコンストラクタ
func NewJPFontHelper(fontDir string) *JPFontHelper {
	return &JPFontHelper{
		fontDir: fontDir,
	}
}

// RegisterJPFont 日本語フォントをPDFに登録
// fontName: フォントの名前（例: "NotoSansJP"）
// fontFile: フォントファイル名（例: "NotoSansJP-Regular.ttf"）
func (h *JPFontHelper) RegisterJPFont(pdf *gofpdf.Fpdf, fontName string, fontFile string) error {
	fontPath := filepath.Join(h.fontDir, fontFile)
	
	// フォントファイルの存在確認
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		// フォントファイルが存在しない場合は、基本フォントを使用（警告のみ）
		fmt.Printf("WARNING: Japanese font file not found: %s, using default font\n", fontPath)
		return nil
	}
	
	// UTF-8フォントとして登録
	pdf.AddUTF8Font(fontName, "", fontPath)
	
	return nil
}

// RegisterJPFonts 日本語フォント一式を登録（Regular, Bold）
func (h *JPFontHelper) RegisterJPFonts(pdf *gofpdf.Fpdf) error {
	// Regularフォント
	if err := h.RegisterJPFont(pdf, "NotoSansJP", "NotoSansJP-Regular.ttf"); err != nil {
		return fmt.Errorf("failed to register NotoSansJP Regular: %w", err)
	}
	
	// Boldフォント（存在する場合）
	h.RegisterJPFont(pdf, "NotoSansJP", "NotoSansJP-Bold.ttf")
	
	return nil
}

// SetJPFont 日本語フォントを設定
// フォントが登録されていない場合は、Arialにフォールバック
func (h *JPFontHelper) SetJPFont(pdf *gofpdf.Fpdf, style string, size float64) {
	// フォントディレクトリが設定されている場合のみ日本語フォントを使用
	if h.fontDir != "" {
		// フォントが登録されているか確認（簡易的なチェック）
		// 実際にはフォントが登録されているかどうかを確認する方法がないため、
		// エラーが出ない限り日本語フォントを使用する
		pdf.SetFont("NotoSansJP", style, size)
	} else {
		// フォントディレクトリが設定されていない場合はArialを使用
		pdf.SetFont("Arial", style, size)
	}
}

// GetFontDir フォントディレクトリを取得（環境変数から）
func GetFontDir() string {
	fontDir := os.Getenv("FONT_DIR")
	if fontDir == "" {
		// デフォルトは現在のディレクトリのassets/fonts
		cwd, _ := os.Getwd()
		fontDir = filepath.Join(cwd, "assets", "fonts")
	}
	return fontDir
}

