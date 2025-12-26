package main

import (
    "fmt"
    "net"
    "sort"
    "sync"
    "time"

    "github.com/miekg/dns"
)

var (
    dnsList = []string{"1.1.1.1", "8.8.8.8", "8.26.56.26", "9.9.9.9", "64.6.65.6", "91.239.100.100", "77.88.8.7", "156.154.70.1", "198.101.242.72", "176.103.130.130"}
    testFQDN = []string{"google.se", "fb.com", "amazon.de", "meta.ua", "mail.ru", "web.de", "svd.se", "ica.se"}
)

func getDNSTime(dnsServer string, repeat int, fqdnList []string, queryType uint16) float64 {
    client := new(dns.Client)
    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn(fqdnList[0]), queryType)

    start := time.Now()
    for i := 0; i < repeat; i++ {
        for _, fqdn := range fqdnList {
            m.SetQuestion(dns.Fqdn(fqdn), queryType)
            _, _, err := client.Exchange(m, net.JoinHostPort(dnsServer, "53"))
            if err != nil {
                continue
            }
        }
    }
    elapsed := time.Since(start)
    return elapsed.Seconds()
}

func testDNS(dnsList []string, repeat int, queryType uint16) map[string]float64 {
    results := make(map[string]float64)
    var wg sync.WaitGroup
    var mu sync.Mutex

    for _, dnsServer := range dnsList {
        wg.Add(1)
        go func(dnsServer string) {
            defer wg.Done()
            timeTaken := getDNSTime(dnsServer, repeat, testFQDN, queryType)
            mu.Lock()
            results[dnsServer] = timeTaken
            mu.Unlock()
        }(dnsServer)
    }

    wg.Wait()
    return results
}

func printResults(results map[string]float64) {
    var sortedResults []struct {
        dnsServer string
        timeTaken float64
    }
    for dnsServer, timeTaken := range results {
        sortedResults = append(sortedResults, struct {
            dnsServer string
            timeTaken float64
        }{dnsServer, timeTaken})
    }
    sort.Slice(sortedResults, func(i, j int) bool {
        return sortedResults[i].timeTaken < sortedResults[j].timeTaken
    })
    for _, result := range sortedResults {
        fmt.Printf("%s: %.4f\n", result.dnsServer, result.timeTaken)
    }
}

func main() {
    repeat := 3
    queryType := dns.TypeA
    results := testDNS(dnsList, repeat, queryType)
    fmt.Printf("DNS Query Test Results (Repeat: %d, Query Type: %d)\n", repeat, queryType)
    printResults(results)
}
