package ffmpeg

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type OutputFormat string

const (
	FormatMP4  OutputFormat = "mp4"
	FormatAVI  OutputFormat = "avi"
	FormatMKV  OutputFormat = "mkv"
	FormatWebM OutputFormat = "webm"
)

type Progress struct {
	Frame     int
	FPS       float64
	Bitrate   string
	Time      string
	Progress  float64 // 0.0 to 1.0
}

type Converter struct {
	ffmpegPath string
}

func NewConverter() (*Converter, error) {
	// Verificar se o FFmpeg está instalado
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg não encontrado no PATH: %v", err)
	}

	return &Converter{
		ffmpegPath: path,
	}, nil
}

func (c *Converter) Convert(inputPath string, format OutputFormat, quality int, progressChan chan<- Progress) (string, error) {
	// Gerar nome do arquivo de saída
	outputPath := changeExtension(inputPath, string(format))

	// Configurar parâmetros de conversão
	args := []string{
		"-i", inputPath,
		"-c:v", "libx264", // Codec de vídeo
		"-preset", "medium",
		"-crf", fmt.Sprintf("%d", quality),
		"-progress", "pipe:1", // Envia progresso para stdout
	}

	// Adicionar parâmetros específicos do formato
	switch format {
	case FormatMP4:
		args = append(args, "-c:a", "aac")
	case FormatAVI:
		args = append(args, "-c:a", "mp3")
	case FormatMKV:
		args = append(args, "-c:a", "aac")
	case FormatWebM:
		args = append(args, "-c:v", "libvpx-vp9", "-c:a", "libopus")
	}

	args = append(args, outputPath)

	cmd := exec.Command(c.ffmpegPath, args...)
	
	// Capturar stdout para análise do progresso
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("erro ao criar pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("erro ao iniciar ffmpeg: %v", err)
	}

	// Analisar progresso
	go func() {
		defer close(progressChan)
		buf := make([]byte, 1024)
		progress := Progress{}
		timeRegex := regexp.MustCompile(`time=(\d+:\d+:\d+\.\d+)`)
		frameRegex := regexp.MustCompile(`frame=(\d+)`)
		fpsRegex := regexp.MustCompile(`fps=(\d+\.?\d*)`)
		bitrateRegex := regexp.MustCompile(`bitrate=(\d+\.?\d*kb/s)`)
		progressRegex := regexp.MustCompile(`out_time_ms=(\d+)`)

		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}

			output := string(buf[:n])
			lines := strings.Split(output, "\n")

			for _, line := range lines {
				if matches := timeRegex.FindStringSubmatch(line); len(matches) > 1 {
					progress.Time = matches[1]
				}
				if matches := frameRegex.FindStringSubmatch(line); len(matches) > 1 {
					progress.Frame, _ = strconv.Atoi(matches[1])
				}
				if matches := fpsRegex.FindStringSubmatch(line); len(matches) > 1 {
					progress.FPS, _ = strconv.ParseFloat(matches[1], 64)
				}
				if matches := bitrateRegex.FindStringSubmatch(line); len(matches) > 1 {
					progress.Bitrate = matches[1]
				}
				if matches := progressRegex.FindStringSubmatch(line); len(matches) > 1 {
					timeMs, _ := strconv.ParseInt(matches[1], 10, 64)
					// Simples estimativa de progresso baseada no tempo
					progress.Progress = float64(timeMs) / 1000000.0
					progressChan <- progress
				}
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("erro na conversão: %v", err)
	}

	return outputPath, nil
}

func changeExtension(path, newExt string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + "." + newExt
} 