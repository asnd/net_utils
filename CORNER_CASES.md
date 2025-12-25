# PCAP Stripper - Corner Cases Analysis

## Executive Summary

This document provides a comprehensive analysis of corner cases in the PCAP Stripper tool, their potential impact, and how they are handled in the refactored code.

## Corner Case Categories

### 1. Input Validation Corner Cases

| Case | Risk | Detection | Handling | Code Location |
|------|------|-----------|----------|---------------|
| Empty input file path | High | Missing string check | Return error before processing | pcap_stripper.go:142 |
| Whitespace-only path | Medium | TrimSpace check | Cleaned before validation | pcap_stripper.go:147 |
| Non-existent file | High | os.Stat() check | Clear error message | pcap_stripper.go:150-152 |
| Directory instead of file | High | info.IsDir() check | Specific error message | pcap_stripper.go:159-161 |
| Unreadable file (permissions) | High | os.Open() attempt | Permission denied error | pcap_stripper.go:164-167 |
| Wrong file extension | Medium | Extension validation | Rejects non-.pcap/.cap files | pcap_stripper.go:171-174 |
| Symlink to non-existent file | Medium | os.Stat() follows links | Same as non-existent file | pcap_stripper.go:150 |
| Corrupted PCAP header | High | pcap.OpenOffline() error | Library error propagated | pcap_stripper.go:252 |
| Zero-byte file | Low | Handled gracefully | Returns 0 packets processed | pcap_stripper.go:298 |

### 2. Parameter Validation Corner Cases

| Case | Risk | Detection | Handling | Code Location |
|------|------|-----------|----------|---------------|
| Zero bytes to strip | High | Value check | Error: must be > 0 | pcap_stripper.go:181-183 |
| Negative bytes to strip | High | Value check | Error: must be > 0 | pcap_stripper.go:181-183 |
| Extremely large bytes (>65536) | Medium | Range check | Warning message | pcap_stripper.go:185-187 |
| Invalid position string | Medium | String matching | Error with valid options | pcap_stripper.go:193-203 |
| Case variations (BEGIN, End) | Low | Normalize to lowercase | Accepts all case variants | pcap_stripper.go:194 |
| Position aliases | Low | Switch statement | begin/start/head/front/end/tail/back | pcap_stripper.go:196-202 |

### 3. Output File Corner Cases

| Case | Risk | Detection | Handling | Code Location |
|------|------|-----------|----------|---------------|
| Output file already exists | Medium | os.Stat() check | Requires --force flag | pcap_stripper.go:238-244 |
| Input = Output (same file) | Critical | Absolute path comparison | Error prevents corruption | pcap_stripper.go:219-234 |
| Input = Output (via symlink) | Critical | filepath.Abs() resolution | Detected as same file | pcap_stripper.go:220-228 |
| No write permission | High | os.Create() error | Permission error message | pcap_stripper.go:265-268 |
| Disk full during write | High | Write/Flush errors | Error + cleanup partial file | pcap_stripper.go:352-355 |
| Output directory doesn't exist | High | os.Create() error | Clear error from OS | pcap_stripper.go:265 |
| Empty output path | Low | Auto-generation | Creates _stripped.pcap name | pcap_stripper.go:207-216 |

### 4. Packet Processing Corner Cases

| Case | Risk | Detection | Handling | Code Location |
|------|------|-----------|----------|---------------|
| Empty packet (0 bytes) | Medium | Length check | Skip, count as SkippedTooSmall | pcap_stripper.go:304-310 |
| Packet smaller than strip size | High | Length comparison | Skip, count in stats | pcap_stripper.go:315-322, 325-332 |
| Packet exactly strip size | Medium | >= comparison | Skipped (would result in 0-byte packet) | pcap_stripper.go:315, 325 |
| Very large packet (>65536) | Low | Handled by library | Processed normally | pcap_stripper.go:298 |
| Malformed packet metadata | Medium | gopacket parsing | Library handles or errors | pcap_stripper.go:337 |
| Write packet failure | High | Error check on write | Error + cleanup + abort | pcap_stripper.go:337-340 |

### 5. Resource and Performance Corner Cases

| Case | Risk | Detection | Handling | Code Location |
|------|------|-----------|----------|---------------|
| Very large files (>1GB) | Medium | Streaming processing | Uses iterator, not loading all | pcap_stripper.go:296-349 |
| Memory exhaustion | Low | Buffered I/O | 1MB buffer limits memory | pcap_stripper.go:271 |
| Slow I/O devices | Low | Progress reporting | Verbose mode shows progress | pcap_stripper.go:346-348 |
| Network file systems | Medium | Standard file I/O | Works but may be slow | - |
| Concurrent access | Medium | Not explicitly locked | OS-level file locking applies | - |
| Interrupted processing | Medium | Defer cleanup | Partial file removed on error | pcap_stripper.go:274-286 |

### 6. Edge Cases in Position Handling

| Case | Risk | Detection | Handling | Code Location |
|------|------|-----------|----------|---------------|
| Beginning: strip=packet_size | Medium | >= check | Packet skipped | pcap_stripper.go:315 |
| End: strip=packet_size | Medium | >= check | Packet skipped | pcap_stripper.go:325 |
| Beginning: strip>packet_size | Medium | >= check | Packet skipped | pcap_stripper.go:315 |
| End: strip>packet_size | Medium | >= check | Packet skipped | pcap_stripper.go:325 |
| Position with leading/trailing spaces | Low | TrimSpace + Lowercase | Normalized correctly | pcap_stripper.go:194 |

## Risk Assessment Matrix

```
Impact vs Probability:

Critical |                    | Input=Output (via symlink)
         |                    |
High     | Zero bytes,        | Disk full, No permissions,
         | Negative bytes,    | Non-existent file,
         | Small packets      | Write errors
         |                    |
Medium   | Wrong extension,   | Already exists, Large files,
         | Invalid position   | Malformed packets
         |                    |
Low      | Empty file,        | Case variations,
         | Aliases,           | Memory limits
         | Buffering          |
         |─────────────────────────────────────────
           Unlikely            Common
```

## Testing Matrix

### Unit Test Cases (Recommended)

```go
// Input validation
TestEmptyInputPath()
TestNonExistentFile()
TestDirectoryAsInput()
TestUnreadableFile()
TestWrongExtension()

// Parameter validation
TestZeroBytes()
TestNegativeBytes()
TestLargeBytes()
TestInvalidPosition()
TestPositionAliases()

// Output validation
TestOutputExists()
TestSameInputOutput()
TestNoWritePermission()

// Packet processing
TestEmptyPackets()
TestSmallPackets()
TestExactSizePackets()
TestNormalPackets()

// Edge cases
TestEmptyPcapFile()
TestSinglePacketFile()
TestAllPacketsTooSmall()
```

### Integration Test Cases (Recommended)

```bash
# Create test PCAP files with:
1. Normal packets (various sizes: 64, 128, 1500, 9000 bytes)
2. Small packets (1-10 bytes)
3. Empty file
4. Corrupted header
5. Mixed packet sizes

# Test scenarios:
- Strip from beginning with various byte counts
- Strip from end with various byte counts
- Process file with all packets too small
- Process with --force overwrite
- Process with auto-generated output
- Process with verbose mode
```

## Defensive Programming Patterns Used

### 1. Early Validation
All inputs validated before processing begins (pcap_stripper.go:74-105)

### 2. Fail Fast
Errors returned immediately rather than attempting recovery

### 3. Resource Cleanup
Defer statements ensure files closed and partial outputs removed (pcap_stripper.go:274-286)

### 4. Clear Error Messages
All errors include context about what failed and why

### 5. Safe Defaults
Position defaults to "beginning", output auto-generated if not specified

### 6. Defensive Checks
Even validated data checked again during processing (e.g., packet length)

### 7. Idempotency
Can safely re-run command (with --force) without side effects

## Known Limitations (Acceptable Risks)

### 1. Disk Space
**Issue**: Cannot predict if enough disk space for output
**Mitigation**: Error caught during write, partial file cleaned up
**Recommendation**: User should ensure adequate space

### 2. File System Limits
**Issue**: Max file size depends on file system (FAT32 = 4GB limit)
**Mitigation**: None (OS limitation)
**Recommendation**: Use modern file systems (ext4, NTFS, APFS)

### 3. Concurrent Access
**Issue**: Multiple processes accessing same file
**Mitigation**: None (relies on OS file locking)
**Recommendation**: User coordination

### 4. Network Latency
**Issue**: Network file systems may be slow
**Mitigation**: Buffered I/O helps, verbose mode shows progress
**Recommendation**: Copy to local disk for large files

### 5. Exotic PCAP Formats
**Issue**: Some custom PCAP variants might not parse
**Mitigation**: Library handles standard formats
**Recommendation**: Use standard tcpdump/Wireshark compatible files

## Recommendations for Users

### Best Practices:

1. **Always use verbose mode** for large files: `-v`
2. **Test with small file first** to verify parameters
3. **Check output file size** to ensure it's reasonable
4. **Keep backups** of important capture files
5. **Use force flag carefully** as it overwrites without confirmation
6. **Verify results** in Wireshark after processing

### Common Mistakes to Avoid:

1. ❌ Using same file for input and output
2. ❌ Stripping more bytes than packet size (will skip packets)
3. ❌ Forgetting to check available disk space
4. ❌ Processing files over slow network connections
5. ❌ Using tool on non-PCAP files despite correct extension

## Conclusion

The refactored PCAP Stripper handles **all critical corner cases** through comprehensive validation, defensive programming, and clear error reporting. The few acceptable limitations are documented and have recommended workarounds.

**Overall Risk Level**: LOW - Production ready with proper user guidance.
