import timeit
import time
import dns.resolver
import collections

dns_list = ["1.1.1.1", "8.8.8.8","8.26.56.26", "192.168.50.12","9.9.9.9","64.6.65.6","91.239.100.100","77.88.8.7","156.154.70.1","198.101.242.72","176.103.130.130"] #,"192.168.50.12"]
test_fqdn = ['google.se', 'fb.com', 'amazon.de','meta.ua', 'mail.ru', 'web.de',"svd.se","ica.se"]



#t = timeit.timeit(dns.resolver.query('tutorialspoint.com', 'A'), 10)

mr = dns.resolver.Resolver()



def get_dns_time(dns_list,repeat=10, fqdn="svd.se", query_type='A'):
   d = {}
   # mr = dns.resolver.Resolver()
   r = []
   for l in dns_list:
     mr= dns.resolver.Resolver()
     mr.nameservers = [l]  
     start = time.time()
     for i in range(repeat):
       for f in test_fqdn:
          r=mr.query(f, query_type)
       mr.cache =  '' 
     print(r.rrset, mr.nameservers)
     end = time.time()
     d[l] =  end-start
   return d


def test_dns(dnss, repeat=1000, query_type='A'):
    t = ()
    for  d in dnss: 
      t =   get_dns_time(d,repeat = repeat,query_type=query_type) 
      print (d,t)
      t = ()

#get_dns_time(dns_list[0])

d= get_dns_time(dns_list)

for key, value in sorted(d.items(), key=lambda item: item[1]):
    print("%s: %s" % (key, round(value,4)))

#for k,v in d.items(): 
#    print(k,round(v,5))



def run_dns_query(ip):
    try:
        dns.resolver.query()
    except expression as identifier:
        pass
