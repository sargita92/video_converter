package ffmpeg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConverter(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("Erro ao criar converter: %v", err)
	}

	if converter.ffmpegPath == "" {
		t.Error("Caminho do FFmpeg não foi definido")
	}
}

func TestChangeExtension(t *testing.T) {
	tests := []struct {
		input    string
		newExt   string
		expected string
	}{
		{"video.h265", "mp4", "video.mp4"},
		{"path/to/video.h265", "avi", "path/to/video.avi"},
		{"video", "mkv", "video.mkv"},
	}

	for _, test := range tests {
		result := changeExtension(test.input, test.newExt)
		if result != test.expected {
			t.Errorf("changeExtension(%q, %q) = %q, esperado %q",
				test.input, test.newExt, result, test.expected)
		}
	}
}

func TestConvert(t *testing.T) {
	// Criar um arquivo de vídeo de teste
	testFile := filepath.Join(t.TempDir(), "test.h265")
	if err := os.WriteFile(testFile, []byte("test data"), 0644); err != nil {
		t.Fatalf("Erro ao criar arquivo de teste: %v", err)
	}

	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("Erro ao criar converter: %v", err)
	}

	progressChan := make(chan Progress)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for range progressChan {
			// Consumir progresso
		}
	}()

	outputPath, err := converter.Convert(testFile, FormatMP4, 23, progressChan)
	if err != nil {
		t.Errorf("Erro na conversão: %v", err)
	}

	// Verificar se o arquivo de saída foi criado
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Arquivo de saída não foi criado: %v", err)
	}

	// Limpar
	os.Remove(outputPath)
	<-done
} 