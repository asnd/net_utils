# GitLab CI/CD Pipeline Documentation

## Overview

This GitLab CI/CD pipeline automatically builds, tests, and packages the PCAP Stripper tool for multiple platforms.

## Pipeline Stages

```
┌──────────────┐
│   PREPARE    │  Download dependencies, verify modules
└──────┬───────┘
       │
       ▼
┌──────────────┐
│     TEST     │  Unit tests, linting, security scanning
└──────┬───────┘
       │
       ▼
┌──────────────┐
│    BUILD     │  Multi-platform compilation
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   PACKAGE    │  Create release artifacts
└──────────────┘
```

## Stage Details

### 1. PREPARE Stage

**Job: `download-dependencies`**
- Downloads all Go module dependencies
- Verifies module checksums
- Caches modules for subsequent stages
- **Duration**: ~30 seconds
- **Artifacts**: go.mod, go.sum (1 hour)

### 2. TEST Stage

**Job: `test:unit`**
- Runs unit tests with race detection
- Generates code coverage report
- **Coverage Format**: Cobertura
- **Required for**: All build jobs
- **Duration**: ~1-2 minutes
- **Artifacts**: coverage.txt (1 week)

**Job: `test:vet`**
- Static analysis with `go vet`
- Checks for common Go mistakes
- **Required for**: All build jobs
- **Duration**: ~30 seconds
- **Failure**: Blocks pipeline

**Job: `test:fmt`**
- Validates code formatting
- Ensures `gofmt` compliance
- **Duration**: ~10 seconds
- **Failure**: Blocks pipeline
- **Fix**: Run `gofmt -w .` locally

**Job: `test:lint`** (Optional)
- Comprehensive linting with golangci-lint
- Checks: gofmt, govet, errcheck, staticcheck, unused, etc.
- **Duration**: ~2-3 minutes
- **Failure**: Warning only (allow_failure: true)

**Job: `test:security`** (Optional)
- Security scanning with gosec
- Generates JSON security report
- **Duration**: ~1-2 minutes
- **Failure**: Warning only (allow_failure: true)
- **Artifacts**: gosec-report.json (1 week)

### 3. BUILD Stage

All builds include:
- Version information from git tags/commits
- Binary stripping (-s -w flags)
- Build metadata (version, commit, timestamp)

**Job: `build:linux-x64`**
- **Platform**: Linux AMD64
- **CGO**: Enabled (required for libpcap)
- **Dependencies**: libpcap-dev
- **Output**: `pcap_stripper-linux-amd64`
- **Size**: ~8-10 MB
- **Duration**: ~1 minute
- **Failure**: Blocks pipeline
- **Artifacts**: 30 days

**Job: `build:windows-x64`**
- **Platform**: Windows AMD64
- **CGO**: Enabled with MinGW cross-compiler
- **Cross-Compiler**: gcc-mingw-w64-x86-64
- **Output**: `pcap_stripper-windows-amd64.exe`
- **Size**: ~8-10 MB
- **Duration**: ~2-3 minutes
- **Failure**: Warning only (allow_failure: true)
- **Note**: Requires Npcap/WinPcap runtime on target system
- **Artifacts**: 30 days

**Job: `build:macos-arm64`**
- **Platform**: macOS ARM64 (Apple Silicon)
- **CGO**: Enabled with clang
- **Compiler**: clang
- **Output**: `pcap_stripper-macos-arm64`
- **Size**: ~8-10 MB
- **Duration**: ~1-2 minutes
- **Failure**: Warning only (allow_failure: true)
- **Note**: May require code signing for Gatekeeper
- **Artifacts**: 30 days

**Job: `build:macos-x64`** (Manual)
- **Platform**: macOS AMD64 (Intel)
- **CGO**: Enabled with clang
- **Trigger**: Manual only
- **Output**: `pcap_stripper-macos-amd64`
- **Size**: ~8-10 MB
- **Duration**: ~1-2 minutes
- **Failure**: Warning only (allow_failure: true)
- **Artifacts**: 30 days

### 4. PACKAGE Stage

**Job: `package:all`**
- Collects all successful builds
- Generates SHA256 checksums
- Creates release README
- **Triggers on**: tags, master, main branches only
- **Dependencies**: All build jobs
- **Output**: release/ directory
- **Artifacts**: 90 days
- **Contents**:
  - All platform binaries
  - SHA256SUMS.txt
  - README.txt

**Job: `release`** (Tags only)
- Creates GitLab Release
- Attaches all artifacts
- **Triggers on**: Git tags only
- **Image**: registry.gitlab.com/gitlab-org/release-cli:latest

## Build Flags and Metadata

Each binary is built with the following ldflags:
```bash
-ldflags="-s -w \
  -X 'main.Version=${CI_COMMIT_TAG:-dev}' \
  -X 'main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)' \
  -X 'main.GitCommit=${CI_COMMIT_SHA:0:8}'"
```

**Flags explained:**
- `-s`: Strip symbol table
- `-w`: Strip DWARF debugging information
- `-X`: Set string variable values at build time

**Version Information:**
Users can check build info with:
```bash
./pcap_stripper --version
# Output: dev (commit: a1b2c3d4, built: 2024-01-15T10:30:00Z)
```

## Caching Strategy

**Cache Key**: `${CI_COMMIT_REF_SLUG}`
- Caches Go modules: `.go/pkg/mod/`
- Caches build artifacts: `.cache/go-build/`
- **Benefit**: Speeds up subsequent builds by ~50%

## Artifacts Summary

| Artifact | Stage | Expiry | Size | Use Case |
|----------|-------|--------|------|----------|
| go.mod/go.sum | prepare | 1 hour | <1 KB | Dependency tracking |
| coverage.txt | test | 1 week | ~10 KB | Coverage analysis |
| gosec-report.json | test | 1 week | ~50 KB | Security audit |
| *-linux-amd64 | build | 30 days | ~9 MB | Linux binary |
| *-windows-amd64.exe | build | 30 days | ~9 MB | Windows binary |
| *-macos-arm64 | build | 30 days | ~9 MB | macOS ARM binary |
| *-macos-amd64 | build | 30 days | ~9 MB | macOS Intel binary |
| release/ | package | 90 days | ~40 MB | All-in-one package |

## Platform-Specific Notes

### Linux x64
- **✅ Fully supported**
- Primary development/testing platform
- Requires: libpcap installed (`apt install libpcap0.8`)
- Works on: Ubuntu, Debian, RHEL, Arch, etc.

### Windows x64
- **⚠️ Cross-compiled** (limitations possible)
- Requires: Npcap or WinPcap installed
- Download Npcap: https://npcap.com/
- Install in "WinPcap compatibility mode"
- May need admin rights to capture packets

### macOS ARM64 (Apple Silicon)
- **⚠️ Cross-compiled** (limitations possible)
- M1, M2, M3 chips
- Requires: libpcap (included in macOS)
- May need: `xattr -d com.apple.quarantine pcap_stripper-macos-arm64`
- May need: Code signing for distribution

### macOS x64 (Intel)
- **⚠️ Manual build only**
- Intel-based Macs
- Same requirements as ARM64
- Built on-demand to save CI time

## Triggering Builds

### Automatic Triggers
```bash
# Any push to master/main
git push origin master

# Push a tag (triggers release)
git tag v1.0.0
git push origin v1.0.0
```

### Manual Triggers
- Navigate to CI/CD → Pipelines
- Click "Run Pipeline"
- Select branch
- Optionally trigger manual jobs (macOS x64)

## Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `PACKAGE_NAME` | pcap_stripper | Output binary name |
| `CGO_ENABLED` | 1 | Enable CGO for libpcap |
| `GOPATH` | `$CI_PROJECT_DIR/.go` | Go workspace |
| `GOCACHE` | `$CI_PROJECT_DIR/.cache/go-build` | Build cache |
| `CI_COMMIT_TAG` | - | Git tag (for releases) |
| `CI_COMMIT_SHA` | - | Git commit hash |
| `CI_COMMIT_REF_SLUG` | - | Branch/tag slug |

## Troubleshooting

### Build Failures

**Linux build fails:**
```bash
# Missing libpcap-dev
Solution: Already installed in before_script

# CGO compilation errors
Check: libpcap-dev version
```

**Windows build fails:**
```bash
# MinGW not found
Solution: Already installed in before_script

# Link errors
Issue: Missing WinPcap SDK (expected with cross-compilation)
Status: allow_failure: true
```

**macOS build fails:**
```bash
# Clang errors
Solution: Already installed in before_script

# Architecture mismatch
Issue: Cross-compilation without macOS SDK
Status: allow_failure: true
```

### Test Failures

**Unit tests fail:**
```bash
# Race conditions
Fix: Address race detector warnings

# Test coverage too low
Action: Add more tests or adjust coverage threshold
```

**Lint warnings:**
```bash
# Formatting issues
Fix: gofmt -w .

# Static analysis warnings
Fix: Address go vet warnings
```

**Security warnings:**
```bash
# gosec findings
Review: Check gosec-report.json
Action: Fix or suppress false positives
```

## Pipeline Optimization

### Speed Improvements
1. **Parallel testing**: All test jobs run in parallel
2. **Caching**: Go modules and build cache
3. **Minimal dependencies**: Only install what's needed
4. **Artifact reuse**: Tests run once, builds reuse results

### Cost Optimization
1. **Manual builds**: macOS x64 is manual-only
2. **Allow failures**: Cross-platform builds don't block
3. **Artifact expiry**: Aggressive cleanup (30-90 days)
4. **Smart caching**: Per-branch cache keys

## Security Considerations

1. **Dependency verification**: `go mod verify` in prepare stage
2. **Security scanning**: gosec runs on every commit
3. **Signed commits**: Recommended (not enforced)
4. **Artifact checksums**: SHA256SUMS.txt in releases
5. **Supply chain**: All tools from official sources

## Creating a Release

### Step 1: Update version
```bash
# Ensure code is ready
go test ./...
go vet ./...
gofmt -w .
```

### Step 2: Create tag
```bash
# Semantic versioning
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### Step 3: Monitor pipeline
- CI/CD → Pipelines
- Watch for test/build success
- Check artifact generation

### Step 4: Verify release
- Navigate to Repository → Releases
- Verify all binaries attached
- Test download and checksums

## Best Practices

1. **Always tag releases**: Use semantic versioning (v1.2.3)
2. **Keep CHANGELOG**: Document changes between versions
3. **Test locally first**: Run tests before pushing
4. **Monitor pipelines**: Check for warnings/failures
5. **Update dependencies**: Regular `go get -u` and test
6. **Review security reports**: Check gosec findings
7. **Verify artifacts**: Download and test binaries

## CI/CD Metrics

**Expected Pipeline Duration:**
- Prepare: 30s
- Test: 2-4 minutes (parallel)
- Build: 3-5 minutes (parallel)
- Package: 1 minute
- **Total**: ~6-10 minutes

**Resource Usage:**
- CPU: 2-4 cores during builds
- Memory: 2-4 GB peak
- Disk: ~500 MB cache
- Bandwidth: ~100 MB artifacts

## Future Improvements

### Potential Enhancements
1. ✅ Add ARM64 Linux builds
2. ✅ Add integration tests
3. ✅ Add benchmark tests
4. ✅ Add Docker builds
5. ✅ Add dependency scanning (Dependabot)
6. ✅ Add SBOM generation
7. ✅ Add performance regression tests

### Advanced Features
- Code signing for macOS/Windows
- Notarization for macOS
- Windows Authenticode signing
- Docker multi-arch images
- Automated changelog generation
- Automated GitHub release mirroring

## Support

For CI/CD issues:
1. Check pipeline logs
2. Review this documentation
3. Check GitLab CI/CD docs
4. Open an issue in the repository

## References

- [GitLab CI/CD Documentation](https://docs.gitlab.com/ee/ci/)
- [Go Cross Compilation](https://golang.org/doc/install/source#environment)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [golangci-lint](https://golangci-lint.run/)
- [gosec](https://github.com/securego/gosec)
