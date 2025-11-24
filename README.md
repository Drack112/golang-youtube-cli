# golang-youtube-cli 🎬

> Uma interface de linha de comando (CLI) para buscar, visualizar e interagir com vídeos do YouTube diretamente pelo terminal, desenvolvida em Go.

---

## Sumário
- [golang-youtube-cli 🎬](#golang-youtube-cli-)
  - [Sumário](#sumário)
  - [Visão Geral](#visão-geral)
  - [Arquitetura 🏗️](#arquitetura-️)
  - [Fluxo de Execução 🔄](#fluxo-de-execução-)
  - [Instalação e Execução 🚀](#instalação-e-execução-)
  - [Dependências 📦](#dependências-)
  - [Exemplos de Uso 🖥️](#exemplos-de-uso-️)
  - [Dicas de Uso 💡](#dicas-de-uso-)
  - [Contribuição 🤝](#contribuição-)
  - [Licença 📄](#licença-)

---

## Visão Geral
O `golang-youtube-cli` permite realizar buscas no YouTube, visualizar resultados em uma interface textual interativa, acessar detalhes dos vídeos e reproduzi-los via player externo (ex: mpv). 

---

## Arquitetura 🏗️
O projeto segue uma estrutura modular, separando responsabilidades:

- **cmd/go-youtube/main.go**: Ponto de entrada. Inicializa o parser de flags, configura opções e inicia o programa TUI.
- **internal/**: Lógica principal dividida em submódulos:
  - **api/**: Realiza buscas e interações com a API do YouTube, incluindo paginação e parsing dos resultados.
  - **flags/**: Parser dos argumentos e opções da CLI, validação de entrada e modo interativo.
  - **handlers/**: Orquestra ações como busca, tratamento de erros e integração entre módulos.
  - **models/**: Estruturas de dados para vídeos, resultados de busca, canais, formatos, etc.
  - **player/**: Detecta e integra com players externos (mpv, yt-dlp), gerencia reprodução e streaming.
  - **tui/**: Implementa a interface textual interativa (Bubble Tea), views, navegação e estados.
  - **ui/**: Componentes visuais, estilos, renderização dos resultados e mensagens de erro.
- **pkg/**: Utilitários diversos (logger, http, manipulação de strings, versionamento).

---

## Fluxo de Execução 🔄
1. O usuário executa o binário ou `go run` passando argumentos ou inicia modo interativo.
2. O parser de flags valida e interpreta a entrada (termo de busca ou URL).
3. O módulo `api` realiza a busca, processa os resultados e retorna para o handler.
4. O handler prepara os dados para exibição e aciona a interface TUI.
5. O usuário navega pelos resultados, acessa detalhes ou inicia a reprodução do vídeo.
6. O módulo `player` integra com o player externo para streaming.

---

## Instalação e Execução 🚀
1. Instale o Go (>=1.18).
2. Clone o repositório:
  ```sh
  git clone https://github.com/Drack112/golang-youtube-cli.git
  cd golang-youtube-cli
  ```
3. Instale o player externo (recomendado: mpv) e yt-dlp/youtube-dl para streaming.
4. Execute:
  ```sh
  go run cmd/go-youtube/main.go
  ```
  Ou compile:
  ```sh
  go build -o go-youtube cmd/go-youtube/main.go
  ./go-youtube
  ```

---

## Dependências 📦
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) (TUI)
- [Lipgloss](https://github.com/charmbracelet/lipgloss) (estilos)
- [mpv](https://mpv.io/) (player externo)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) ou [youtube-dl](https://github.com/ytdl-org/youtube-dl) (streaming)

---

## Exemplos de Uso 🖥️

Busca por vídeos:
```sh
go run cmd/go-youtube/main.go search "golang tutorial"
```

Busca interativa:
```sh
go run cmd/go-youtube/main.go
```

Reprodução de vídeo:
Selecione o vídeo desejado na interface e pressione a tecla indicada para iniciar o player externo.

---

## Dicas de Uso 💡
- Use o modo interativo para explorar resultados rapidamente.
- Ative o modo debug para logs detalhados: `go run cmd/go-youtube/main.go -debug`
- Experimente diferentes termos de busca para resultados variados.
- Configure o player externo e yt-dlp para melhor experiência de streaming.

---

## Contribuição 🤝
Contribuições são bem-vindas! Para reportar bugs, sugerir melhorias ou enviar pull requests:
- Abra uma issue no repositório.
- Siga o padrão de código e documentação do projeto.
- Consulte os arquivos em `internal/` e `pkg/` para entender a estrutura.

---

## Licença 📄
Este projeto está sob a licença MIT. Consulte o arquivo LICENSE para mais detalhes.

---

Para dúvidas, sugestões ou contribuições, utilize as issues do repositório ou entre em contato diretamente.
