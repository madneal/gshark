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

func RunDNS(domain string) error {
	cmdLines := vars.GOBUSTER + " dns -d " + domain + " -w " + vars.SUBDOMAIN_WORDLIST
	cmdArgs := strings.Fields(cmdLines)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	stdout, err := cmd.StdoutPipe()
	cmd.Start()
	oneByte := make([]byte, 100)
	var foundDomian string
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			fmt.Println(err)
			break
		}
		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		if strings.Contains(string(line), "Found") {
			foundDomian = strings.Replace(string(line), "Found: ", "", 1)
			foundDomian = strings.ToLower(foundDomian)
			println(foundDomian)
			subdomain := models.Subdomain{
				Domain:    &domain,
				Subdomain: &foundDomian,
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
	return err
}

func RunTask(duration time.Duration) {
	rules, err := models.GetValidRulesByType("domain")
	if err != nil {
		logger.Log.Error(err)
	}
	for _, rule := range rules {
		domain := rule.Pattern
		err = RunDNS(domain)
		if err != nil {
			logger.Log.Error(err)
		}
		logger.Log.Infof("Complete the scan of domain %s, start to sleep %v seconds", domain, duration*time.Second)
		time.Sleep(duration * time.Second)
	}
}
