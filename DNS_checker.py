import time
import dns.resolver
import concurrent.futures

dns_list = ["1.1.1.1", "8.8.8.8", "8.26.56.26", "9.9.9.9", "64.6.65.6", "91.239.100.100", "77.88.8.7", "156.154.70.1", "198.101.242.72", "176.103.130.130"]
test_fqdn = ['google.se', 'fb.com', 'amazon.de', 'meta.ua', 'mail.ru', 'web.de', "svd.se", "ica.se"]

mr = dns.resolver.Resolver()
for s in mr.nameservers:
    dns_list.append(s)

def get_dns_time(dns_server, repeat=3, fqdn_list=test_fqdn, query_type='A'):
    mr = dns.resolver.Resolver()
    mr.nameservers = [dns_server]
    start = time.time()
    for _ in range(repeat):
        for fqdn in fqdn_list:
            try:
                mr.query(fqdn, query_type)
            except dns.exception.DNSException:
                pass
        mr.cache = ''
    end = time.time()
    return end - start

def test_dns(dns_list, repeat=1000, query_type='A'):
    results = {}
    with concurrent.futures.ThreadPoolExecutor() as executor:
        futures = [executor.submit(get_dns_time, dns_server, repeat, test_fqdn, query_type) for dns_server in dns_list]
        for dns_server, future in zip(dns_list, futures):
            results[dns_server] = future.result()
    return results

def print_results(results):
    sorted_results = sorted(results.items(), key=lambda item: item[1])
    for dns_server, time_taken in sorted_results:
        print(f"{dns_server}: {round(time_taken, 4)}")

def main():
    repeat = 1000
    query_type = 'A'
    results = test_dns(dns_list, repeat, query_type)
    print(f"DNS Query Test Results (Repeat: {repeat}, Query Type: {query_type})")
    print_results(results)

if __name__ == "__main__":
    main()