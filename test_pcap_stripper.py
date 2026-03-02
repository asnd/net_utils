#!/usr/bin/env python3

import os
import sys
import tempfile
import unittest
from pathlib import Path

try:
    from scapy.all import Ether, IP, TCP, wrpcap, rdpcap
except ImportError:
    print("Error: scapy is required for testing. Install with: pip install scapy")
    sys.exit(1)


class TestPcapStripper(unittest.TestCase):
    def setUp(self):
        self.test_dir = tempfile.mkdtemp()
        self.input_pcap = os.path.join(self.test_dir, "input.pcap")
        self.output_pcap = os.path.join(self.test_dir, "output.pcap")

        packets = []
        for i in range(5):
            pkt = Ether() / IP(src=f"10.0.0.{i+1}", dst=f"10.0.0.{i+10}") / TCP(sport=1000+i, dport=80)
            packets.append(pkt)

        wrpcap(self.input_pcap, packets)

    def tearDown(self):
        for f in [self.input_pcap, self.output_pcap]:
            if os.path.exists(f):
                os.remove(f)
        os.rmdir(self.test_dir)

    def test_strip_bytes_basic(self):
        from pcap_stripper import strip_bytes_from_packets

        result = strip_bytes_from_packets(self.input_pcap, self.output_pcap, 20, verbose=False)

        self.assertTrue(result)
        self.assertTrue(os.path.exists(self.output_pcap))

        output_packets = rdpcap(self.output_pcap)
        self.assertEqual(len(output_packets), 5)

        for pkt in output_packets:
            self.assertGreaterEqual(len(pkt), 14)

    def test_strip_bytes_skips_small_packets(self):
        from pcap_stripper import strip_bytes_from_packets

        small_packets = [Ether() / IP(src="10.0.0.1", dst="10.0.0.2")]
        wrpcap(self.input_pcap, small_packets)

        result = strip_bytes_from_packets(self.input_pcap, self.output_pcap, 100, verbose=False)

        self.assertTrue(result)
        output_packets = rdpcap(self.output_pcap)
        self.assertEqual(len(output_packets), 0)

    def test_strip_eth_basic(self):
        from pcap_stripper import strip_until_second_ethernet

        inner_pkt = Ether() / IP(src="10.0.0.1", dst="10.0.0.2")
        outer_pkt = Ether() / inner_pkt
        wrpcap(self.input_pcap, [outer_pkt])

        result = strip_until_second_ethernet(self.input_pcap, self.output_pcap, verbose=False)

        self.assertTrue(result)
        self.assertTrue(os.path.exists(self.output_pcap))

    def test_filter_bpf(self):
        from pcap_stripper import filter_packets_bpf

        packets = [
            Ether() / IP(src="10.0.0.1", dst="10.0.0.2") / TCP(sport=80, dport=443),
            Ether() / IP(src="10.0.0.3", dst="10.0.0.4") / UDP(sport=53, dport=53),
        ]
        wrpcap(self.input_pcap, packets)

        result = filter_packets_bpf(self.input_pcap, self.output_pcap, "tcp", verbose=False)

        self.assertTrue(result)
        self.assertTrue(os.path.exists(self.output_pcap))


class UDP:
    def __init__(self, sport=0, dport=0):
        self.sport = sport
        self.dport = dport


if __name__ == "__main__":
    unittest.main()
