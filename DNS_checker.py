import time
import dns.resolver
import concurrent.futures
from typing import List, Dict

DNS_SERVERS = ["1.1.1.1", "8.8.8.8", "8.26.56.26", "9.9.9.9", "64.6.65.6", "91.239.100.100", "77.88.8.7", "156.154.70.1", "198.101.242.72", "176.103.130.130"]
TEST_FQDNS = ['google.se', 'fb.com', 'amazon.de', 'meta.ua', 'mail.ru', 'web.de', "svd.se", "ica.se"]

def get_dns_time(dns_server: str, repeat: int = 3, fqdn_list: List[str] = TEST_FQDNS, query_type: str = 'A') -> float:
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

def test_dns(dns_servers: List[str], repeat: int = 1000, query_type: str = 'A') -> Dict[str, float]:
    results = {}
    with concurrent.futures.ThreadPoolExecutor() as executor:
        futures = {executor.submit(get_dns_time, dns_server, repeat, TEST_FQDNS, query_type): dns_server for dns_server in dns_servers}
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
    repeat = 1000
    query_type = 'A'
    results = test_dns(DNS_SERVERS, repeat, query_type)
    print(f"DNS Query Test Results (Repeat: {repeat}, Query Type: {query_type})")
    print_results(results)

if __name__ == "__main__":
    main()