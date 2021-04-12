package gobuster

import (
	"bufio"
	"fmt"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
	"os/exec"
	"strings"
	"time"
)

func RunDNS(domain string) {
	if global.GVA_CONFIG.Search.GobusterFilePath == "" {
		fmt.Println("Please specify the file path of gobuster!")
		return
	}
	if global.GVA_CONFIG.Search.SubdomainWordList == "" {
		fmt.Println("Please specify the file path of subdomain wordlist file path!")
		return
	}
	cmdLines := global.GVA_CONFIG.Search.GobusterFilePath + " dns -d " + domain + " -w " +
		global.GVA_CONFIG.Search.SubdomainWordList
	cmdArgs := strings.Fields(cmdLines)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	stdout, err := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		global.GVA_LOG.Error("gobuster start error", zap.Any("error", err))
	}
	oneByte := make([]byte, 100)
	var foundDomain string
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			global.GVA_LOG.Error("gobuster output error", zap.Any("error", err))
			break
		}
		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		if strings.Contains(string(line), "Found") {
			foundDomain = strings.Replace(string(line), "Found: ", "", 1)
			foundDomain = strings.ToLower(foundDomain)
			foundDomain = strings.Replace(foundDomain, "\r", "", -1)
			foundDomain = strings.TrimSpace(foundDomain)
			println(foundDomain)
			subdomain := model.Subdomain{
				Domain:    domain,
				Subdomain: foundDomain,
				Status:    0,
			}
			err := service.CreateSubdomain(subdomain)
			if err != nil {
				global.GVA_LOG.Error("create subdomain failed!", zap.Any("error", err))
			}
		}
	}
	defer cmd.Wait()
}

func RunTask(duration time.Duration) {
	err, rules := service.GetValidRulesByType("domain")
	if err != nil {
		global.GVA_LOG.Error("get subdomain rules error", zap.Any("error", err))
	}
	if len(rules) == 0 {
		global.GVA_LOG.Warn("There's no rule for subdomain scan, please specify a subdomain rule at least")
	}
	for _, rule := range rules {
		start := time.Now()
		domain := rule.Content
		RunDNS(domain)
		fmt.Printf("Complete the scan of domain %s, cost %v, start to sleep %v seconds",
			domain, time.Since(start), duration*time.Second)
		time.Sleep(duration * time.Second)
	}
}
