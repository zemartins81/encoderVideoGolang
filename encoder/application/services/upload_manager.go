package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

// VideoUpload é uma estrutura que representa um upload de vídeo. Ela contém uma lista de caminhos
// para arquivos de vídeo, o caminho do vídeo, o bucket de saída e uma lista de erros.
type VideoUpload struct {
	// Paths é uma lista de caminhos para arquivos de vídeo.
	Paths []string
	// VideoPath é o caminho do vídeo.
	VideoPath string
	// OutputBucket é o bucket de saída para o upload.
	OutputBucket string
	// Errors é uma lista de erros que ocorreram durante o upload.
	Errors []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

// UploadObject envia um objeto para o bucket de saída, retornando um erro em caso de falha.
//
// Parâmetros:
//   - objectpath: o caminho do objeto a ser enviado.
//   - client: o cliente de armazenamento de objetos.
//   - ctx: o contexto de execução.
//
// Retorno:
//   - um erro se ocorrer um erro durante o envio.
func (vu *VideoUpload) UploadObject(objectpath string, client *storage.Client, ctx context.Context) error {
	// Divide o caminho do objeto em duas partes, a parte antes do caminho localStoragePath
	// e a parte após ela.
	path := strings.Split(objectpath, os.Getenv("localStoragePath")+"/")

	// Abre o arquivo a ser enviado.
	f, err := os.Open(objectpath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Cria um novo escritor de objetos no bucket de saída, definindo permissões de acesso.
	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	// Copia o conteúdo do arquivo para o escritor de objetos.
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	// Fecha o escritor de objetos, retornando um erro em caso de falha.
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

// loadPaths é uma função que carrega os caminhos dos arquivos de vídeo no diretório
// especificado em VideoPath, armazenando-os na lista Paths. Ele usa a função filepath.Walk
// para percorrer o diretório e suas subpastas, filtrando apenas os arquivos, não as pastas.
//
// Retorno:
//   - um erro, caso ocorra algum durante a execução.
func (vu *VideoUpload) loadPaths() error {
	// Percorre o diretório de vídeo, adicionando os caminhos de arquivos encontrados à lista Paths.
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		// Verifica se o arquivo atual é uma pasta, ignorando-a.
		if !info.IsDir() {
			// Adiciona o caminho do arquivo à lista Paths.
			vu.Paths = append(vu.Paths, path)
		}
		// Retorna um erro, caso ocorra durante a iteração.
		return nil
	})
	// Retorna um erro, caso ocorra durante a execução de filepath.Walk.
	if err != nil {
		return err
	}
	// Retorna nulo, indicando que a função foi concluída com sucesso.
	return nil
}

// getClientUpload retorna um cliente de armazenamento, um contexto e um erro.
// Ele cria um contexto em Background, conecta com a API do GCP e retorna um cliente
// para acessar o serviço de armazenamento de objetos.
//
// Retorno:
//   - cliente de armazenamento de objetos
//   - contexto de execução
//   - erro, em caso de falha na conexão com a API do GCP
func getClientUpload() (*storage.Client, context.Context, error) {
	// Cria um contexto em Background.
	ctx := context.Background()

	// Cria um novo cliente de armazenamento de objetos, conectando com a API do GCP.
	client, err := storage.NewClient(ctx)
	if err != nil {
		// Retorna o erro ocorrido durante a conexão com a API do GCP.
		return nil, nil, err
	}

	// Retorna o cliente de armazenamento de objetos, o contexto de execução e nulo, indicando que a função foi concluída com sucesso.
	return client, ctx, nil
}

// ProcessUpload é uma função que inicia o processo de upload dos vídeos.
// Ela carrega os caminhos dos arquivos, obtém um cliente de armazenamento e
// inicia uma série de workers que realizam o upload dos arquivos.
//
// Parâmetros:
//   - concurrency: número de workers que serão iniciados simultaneamente.
//   - doneUpload: canal de saída para informar se o upload foi concluído.
//
// Retorno:
//   - erro, caso ocorra durante a execução do processo.
func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	// Cria um canal para receber os índices dos arquivos a serem enviados.
	in := make(chan int, runtime.NumCPU())
	// Cria um canal para receber as respostas dos workers após o upload.
	returnChannel := make(chan string)

	// Carrega os caminhos dos arquivos a serem enviados.
	err := vu.loadPaths()
	if err != nil {
		return err
	}

	// Obtém um cliente de armazenamento para acessar o serviço de armazenamento de objetos.
	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	// Inicia os workers para realizar o upload dos arquivos.
	for process := 0; process < concurrency; process++ {
		go vu.upLoadWorker(in, returnChannel, uploadClient, ctx)
	}

	// Envia os índices dos arquivos para o canal de entrada dos workers.
	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	// Espera pela resposta dos workers após o upload.
	for r := range returnChannel {
		// Se uma resposta for diferente de vazia, significa que ocorreu um erro durante o upload.
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

// upLoadWorker é uma função que executa o upload de um arquivo. Ela recebe um canal de entrada para
// receber os índices dos arquivos, um canal de retorno para enviar as respostas após o upload,
// um cliente de armazenamento e um contexto. Ela inicia um loop onde espera um índice do
// canal de entrada, faz o upload do arquivo correspondente e envia uma resposta ao canal de retorno.
// Quando o canal de entrada for fechado, a função envia "upload completed" ao canal de retorno.
func (vu *VideoUpload) upLoadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {

	// Loop que espera os índices dos arquivos a serem enviados.
	for x := range in {
		// Faz o upload do arquivo correspondente ao índice recebido.
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)
		// Se ocorrer um erro durante o upload, o índice do arquivo é adicionado à lista de erros.
		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("Error on upload: %v. Error: %v", vu.Paths[x], err)
			returnChan <- err.Error()
		}
		// Envia uma resposta vazia para indicar que o upload foi concluído com sucesso.
		returnChan <- ""
	}

	// Envia "upload completed" quando o canal de entrada for fechado.
	returnChan <- "upload completed"
}
