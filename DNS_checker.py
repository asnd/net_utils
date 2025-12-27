import time
import dns.resolver
import concurrent.futures
import argparse
from typing import List, Dict

DEFAULT_DNS_SERVERS = ["1.1.1.1", "8.8.8.8", "8.26.56.26", "9.9.9.9", "64.6.65.6", "91.239.100.100", "77.88.8.7", "156.154.70.1", "198.101.242.72", "176.103.130.130"]
DEFAULT_TEST_FQDNS = ['google.se', 'fb.com', 'amazon.de', 'meta.ua', 'mail.ru', 'web.de', "svd.se", "ica.se"]

def get_dns_time(dns_server: str, repeat: int, fqdn_list: List[str], query_type: str = 'A') -> float:
    resolver = dns.resolver.Resolver()
    resolver.nameservers = [dns_server]
    start_time = time.time()
    for _ in range(repeat):
        for fqdn in fqdn_list:
            try:
                resolver.resolve(fqdn, query_type)
            except dns.exception.DNSException:
                pass
        resolver.cache = {}
    end_time = time.time()
    return end_time - start_time

def test_dns(dns_servers: List[str], fqdn_list: List[str], repeat: int = 1000, query_type: str = 'A') -> Dict[str, float]:
    results = {}
    with concurrent.futures.ThreadPoolExecutor() as executor:
        futures = {executor.submit(get_dns_time, dns_server, repeat, fqdn_list, query_type): dns_server for dns_server in dns_servers}
        for future in concurrent.futures.as_completed(futures):
            dns_server = futures[future]
            try:
                results[dns_server] = future.result()
            except Exception as e:
                results[dns_server] = float('inf')
                print(f"Error testing {dns_server}: {e}")
    return results

def print_results(results: Dict[str, float]) -> None:
    sorted_results = sorted(results.items(), key=lambda item: item[1])
    for dns_server, time_taken in sorted_results:
        print(f"{dns_server}: {round(time_taken, 4)}")

def main() -> None:
    parser = argparse.ArgumentParser(description="DNS Query Performance Checker")
    parser.add_argument("--servers", help="Comma-separated list of DNS servers to test", type=str)
    parser.add_argument("--domains", help="Comma-separated list of domains to query", type=str)
    parser.add_argument("--repeat", help="Number of times to repeat the test", type=int, default=3)
    
    args = parser.parse_args()

    if args.servers:
        dns_servers = [s.strip() for s in args.servers.split(',')]
    else:
        dns_servers = DEFAULT_DNS_SERVERS

    if args.domains:
        fqdn_list = [d.strip() for d in args.domains.split(',')]
    else:
        fqdn_list = DEFAULT_TEST_FQDNS

    repeat = args.repeat
    query_type = 'A'
    
    print(f"DNS Query Test Results (Repeat: {repeat}, Query Type: {query_type})")
    results = test_dns(dns_servers, fqdn_list, repeat, query_type)
    print_results(results)

if __name__ == "__main__":
    main()
