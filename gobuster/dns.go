package gobuster

import (
	"bufio"
	"fmt"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/models"
	"github.com/madneal/gshark/vars"
	"os/exec"
	"strings"
	"time"
)

func RunDNS(domain string) {
	if vars.GOBUSTER == "" {
		fmt.Println("Please specify the file path of gobuster!")
		return
	}
	if vars.SUBDOMAIN_WORDLIST == "" {
		fmt.Println("Please specify the file path of subdomain wordlist file path!")
		return
	}
	cmdLines := vars.GOBUSTER + " dns -d " + domain + " -w " + vars.SUBDOMAIN_WORDLIST
	cmdArgs := strings.Fields(cmdLines)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	stdout, err := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		logger.Log.Error(err)
	}
	oneByte := make([]byte, 100)
	var foundDomain string
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			logger.Log.Error(err)
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
			subdomain := models.Subdomain{
				Domain:    &domain,
				Subdomain: &foundDomain,
				Status:    0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			_, err := subdomain.Insert()
			if err != nil {
				logger.Log.Error(err)
			}
		}
	}
	defer cmd.Wait()
}

func RunTask(duration time.Duration) {
	rules, err := models.GetValidRulesByType("domain")
	if err != nil {
		logger.Log.Error(err)
	}
	for _, rule := range rules {
		start := time.Now()
		domain := rule.Pattern
		RunDNS(domain)
		if err != nil {
			logger.Log.Error(err)
		}
		logger.Log.Infof("Complete the scan of domain %s, cost %v, start to sleep %v seconds",
			domain, time.Since(start), duration*time.Second)
		time.Sleep(duration * time.Second)
	}
}
