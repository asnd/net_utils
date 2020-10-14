import timeit
import dnspython as dns
import dns.resolver

dns_list = ["1.1.1.1", "8.8.8.8"]
test_fqdn = ['google.se', 'fb.com', 'amazon.de']

t = timeit.timeit(dns.resolver.query('tutorialspoint.com', 'A'), 10)

print("Time is  - {}".format(t))


def run_dns_query(ip):
    try:
        dns.resolver.query(ip)
    except expression as identifier:
        pass
