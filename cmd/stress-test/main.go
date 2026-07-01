package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 1, "Número total de requisições")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "erro: --url é obrigatório")
		os.Exit(1)
	}
	if *requests <= 0 || *concurrency <= 0 {
		fmt.Fprintln(os.Stderr, "erro: --requests e --concurrency devem ser > 0")
		os.Exit(1)
	}

	rCh := createFullChannel(*requests)
	close(rCh)
	resultCh := make(chan int, *requests)

	start := time.Now()
	createAndRunJobs(rCh, resultCh, *concurrency, *url)
	elapsed := time.Since(start)

	close(resultCh)

	createReportAndPrint(resultCh, elapsed)

}

func createReportAndPrint(resultCh chan int, elapsed time.Duration) {

	scMap := make(map[int]int)

	total, ok200, failed := 0, 0, 0
	for sc := range resultCh {
		total++
		switch sc {
		case 200:
			ok200++
		case -1:
			failed++
		default:
			scMap[sc]++
		}
	}
	fmt.Printf("Total de requests: %d\n", total)
	fmt.Printf("Requests com status 200: %d\n", ok200)
	fmt.Printf("Falhas (erro de conexão/timeout): %d\n", failed)
	fmt.Printf("Tempo total: %.2f s\n", elapsed.Seconds())
	fmt.Printf("Throughput: %.2f req/s\n", float64(total)/elapsed.Seconds())

	fmt.Println("Status codes:")
	if ok200 > 0 {
		fmt.Printf("200: %d\n", ok200)
	}
	for key, value := range scMap {

		fmt.Printf("%d: %d\n", key, value)
	}

}

func createAndRunJobs(rCh chan byte, resultCh chan int, concurrency int, url string) {

	var wg sync.WaitGroup
	client := &http.Client{Timeout: 10 * time.Second}

	for range concurrency {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range rCh {
				resultCh <- doRequest(client, url)
			}
		}()
	}
	wg.Wait()
}

func doRequest(client *http.Client, url string) int {
	resp, err := client.Get(url)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}

func createFullChannel(qtd int) chan byte {

	c := make(chan byte, qtd)

	for range qtd {
		c <- 0
	}

	return c
}
