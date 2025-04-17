package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sargita/video_converter/pkg/ffmpeg"
)

type App struct {
	window     fyne.Window
	converter  *ffmpeg.Converter
	filePath   string
	statusText *widget.Label
	progress   *widget.ProgressBar
	format     *widget.Select
	quality    *widget.Slider
	details    *widget.Label
}

func NewApp() *App {
	myApp := app.New()
	window := myApp.NewWindow("Conversor de Vídeo H.265")
	
	converter, err := ffmpeg.NewConverter()
	if err != nil {
		dialog.ShowError(err, window)
	}

	// Criar widgets de seleção
	format := widget.NewSelect([]string{
		string(ffmpeg.FormatMP4),
		string(ffmpeg.FormatAVI),
		string(ffmpeg.FormatMKV),
		string(ffmpeg.FormatWebM),
	}, nil)
	format.SetSelected(string(ffmpeg.FormatMP4))

	quality := widget.NewSlider(0, 51)
	quality.Value = 23
	quality.Step = 1

	return &App{
		window:    window,
		converter: converter,
		statusText: widget.NewLabel("Selecione um arquivo para converter"),
		progress:   widget.NewProgressBar(),
		format:     format,
		quality:    quality,
		details:    widget.NewLabel(""),
	}
}

func (a *App) Run() {
	// Criar widgets
	selectButton := widget.NewButton("Selecionar Arquivo", a.selectFile)
	convertButton := widget.NewButton("Converter", a.convertFile)
	convertButton.Disable()

	// Layout
	content := container.NewVBox(
		widget.NewLabel("Selecione um arquivo H.265 para converter:"),
		selectButton,
		widget.NewLabel("Formato de saída:"),
		a.format,
		widget.NewLabel("Qualidade (0-51, menor = melhor):"),
		a.quality,
		convertButton,
		a.progress,
		a.details,
		a.statusText,
	)

	a.window.SetContent(content)
	a.window.Resize(fyne.NewSize(400, 400))
	a.window.ShowAndRun()
}

func (a *App) selectFile() {
	dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if uri == nil {
			return
		}
		
		a.filePath = uri.URI().Path()
		a.statusText.SetText("Arquivo selecionado: " + a.filePath)
	}, a.window)
}

func (a *App) convertFile() {
	if a.filePath == "" {
		a.statusText.SetText("Nenhum arquivo selecionado")
		return
	}

	format := ffmpeg.OutputFormat(a.format.Selected)
	quality := int(a.quality.Value)
	
	a.statusText.SetText("Convertendo...")
	a.progress.SetValue(0)
	a.details.SetText("")
	
	progressChan := make(chan ffmpeg.Progress)
	
	go func() {
		outputPath, err := a.converter.Convert(a.filePath, format, quality, progressChan)
		if err != nil {
			dialog.ShowError(err, a.window)
			a.statusText.SetText("Erro na conversão")
			return
		}
		
		a.progress.SetValue(1)
		a.statusText.SetText("Conversão concluída! Arquivo salvo em: " + outputPath)
	}()

	// Atualizar interface com o progresso
	go func() {
		for progress := range progressChan {
			a.progress.SetValue(progress.Progress)
			a.details.SetText(fmt.Sprintf(
				"Frame: %d | FPS: %.2f | Bitrate: %s | Time: %s",
				progress.Frame,
				progress.FPS,
				progress.Bitrate,
				progress.Time,
			))
		}
	}()
} 