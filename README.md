# Conversor de Vídeo H.265

Este é um aplicativo desenvolvido em Go para conversão de vídeos no formato H.265 (HEVC) para outros formatos de vídeo comuns. O aplicativo utiliza FFmpeg para realizar as conversões e oferece uma interface gráfica simples e intuitiva.

## Requisitos

- Go 1.20 ou superior
- FFmpeg instalado no sistema
- Windows 10/11 (sistema operacional alvo)

## Funcionalidades

- Conversão de vídeos H.265 para:
  - MP4 (H.264)
  - AVI
  - MKV
  - WebM
- Interface gráfica simples e intuitiva
- Suporte a múltiplos arquivos
- Progresso da conversão em tempo real
- Configurações personalizáveis de qualidade

## Estrutura do Projeto

```
video_converter/
├── cmd/
│   └── main.go           # Ponto de entrada da aplicação
├── internal/
│   ├── converter/        # Lógica de conversão de vídeo
│   ├── gui/             # Interface gráfica
│   └── utils/           # Utilitários
├── pkg/
│   └── ffmpeg/          # Wrapper para FFmpeg
├── go.mod
├── go.sum
└── README.md
```

## Dependências Principais

- [FFmpeg](https://ffmpeg.org/) - Para processamento de vídeo
- [fyne](https://fyne.io/) - Para interface gráfica
- [go-ffmpeg](https://github.com/xfrr/goffmpeg) - Wrapper Go para FFmpeg

## Como Desenvolver

1. Clone o repositório:
```bash
git clone [URL_DO_REPOSITORIO]
cd video_converter
```

2. Instale as dependências:
```bash
go mod download
```

3. Certifique-se de ter o FFmpeg instalado e disponível no PATH do sistema

4. Execute o projeto:
```bash
go run cmd/main.go
```

## Como Usar

1. Inicie o aplicativo
2. Clique em "Selecionar Arquivo" para escolher um vídeo H.265
3. Selecione o formato de saída desejado
4. Ajuste a qualidade usando o slider (0-51, onde menor é melhor)
5. Clique em "Converter"
6. Acompanhe o progresso da conversão na barra de progresso

## Qualidade de Vídeo

O aplicativo utiliza o parâmetro CRF (Constant Rate Factor) para controlar a qualidade do vídeo:
- 0-18: Qualidade muito alta (quase sem perdas)
- 19-23: Qualidade alta (padrão)
- 24-28: Qualidade média
- 29-51: Qualidade baixa

## Testes

Para executar os testes unitários:
```bash
go test ./pkg/ffmpeg
```

## Licença

MIT License 