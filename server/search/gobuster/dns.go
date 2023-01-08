package gobuster

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const timeout = 30 * time.Second

func preRun(ctx context.Context, domain string) error {
	guid := uuid.New()
	_, err := dnsLookup(ctx, fmt.Sprintf("%s.%s", guid, domain))
	if err == nil {
		err = errors.New("There may be a wildcard resolve for domain: " + domain)
		fmt.Println(err)
	}
	return nil
}

func getWordlist(wordlistFile string) (*bufio.Scanner, error) {
	if wordlistFile == "" {
		return nil, errors.New("has not provide wordlistFile")
	}
	wordlist, err := os.Open(wordlistFile)
	if err != nil {
		return nil, err
	}
	return bufio.NewScanner(wordlist), nil
}

func dnsLookup(ctx context.Context, domain string) ([]string, error) {
	var resolver net.Resolver
	return resolver.LookupHost(ctx, domain)
}

func RunDNS(ctx context.Context, domain string) {
	if err := preRun(ctx, domain); err != nil {
		global.GVA_LOG.Error("preRun error", zap.Error(err))
		return
	}
	var workerGroup sync.WaitGroup
	workerGroup.Add(10)
	wordChan := make(chan string, 10)
	for i := 0; i < 10; i++ {
		go worker(ctx, domain, wordChan, &workerGroup)
	}

	scanner, err := getWordlist(global.GVA_CONFIG.Search.SubdomainWordList)
	if err != nil {
		fmt.Println(err)
		return
	}
Scan:
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			break Scan
		default:
			word := scanner.Text()
			wordChan <- word
		}
	}
	close(wordChan)
	workerGroup.Wait()
}

func worker(ctx context.Context, domain string, wordChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case word, ok := <-wordChan:
			if !ok {
				return
			}
			wordCleaned := strings.TrimSpace(word)
			err := run(ctx, wordCleaned, domain)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func run(ctx context.Context, word, domain string) error {
	subdomain := fmt.Sprintf("%s.%s", word, domain)

	_, err := dnsLookup(ctx, subdomain)
	if err == nil {
		result := model.Subdomain{
			Domain:    domain,
			Subdomain: subdomain,
			Status:    0,
		}
		err = service.CreateSubdomain(result)
		if err != nil {
			return err
		}
	}
	return err
}

func RunTask(duration time.Duration) {
	err, rules := service.GetValidRulesByType("domain")
	if err != nil {
		global.GVA_LOG.Error("get subdomain rules error", zap.Any("error", err))
	}
	for _, rule := range rules {
		start := time.Now()
		domain := rule.Content
		RunDNS(context.Background(), domain)
		fmt.Printf("Complete the scan of domain %s, cost %v, start to sleep %v seconds",
			domain, time.Since(start), duration*time.Second)
		time.Sleep(duration * time.Second)
	}
}
