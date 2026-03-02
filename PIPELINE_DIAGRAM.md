# GitLab CI/CD Pipeline Diagram

## Full Pipeline Visualization

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              STAGE: PREPARE                                  │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐    │
│  │ download-dependencies                                               │    │
│  │ • go mod download                                                   │    │
│  │ • go mod verify                                                     │    │
│  │ ✓ Artifacts: go.mod, go.sum (1h)                                   │    │
│  └────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                               STAGE: TEST                                    │
│                         (All jobs run in parallel)                           │
│                                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │ test:unit    │  │  test:vet    │  │  test:fmt    │  │  test:lint   │   │
│  │              │  │              │  │              │  │              │   │
│  │ • Race det.  │  │ • go vet     │  │ • gofmt -l   │  │ • golangci   │   │
│  │ • Coverage   │  │              │  │              │  │              │   │
│  │ REQUIRED ✓   │  │ REQUIRED ✓   │  │ REQUIRED ✓   │  │ Optional ⚠   │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘   │
│                                                                              │
│  ┌──────────────┐                                                           │
│  │test:security │                                                           │
│  │              │                                                           │
│  │ • gosec scan │                                                           │
│  │ Optional ⚠   │                                                           │
│  └──────────────┘                                                           │
└─────────────────────────────────────────────────────────────────────────────┘
                  │
                  │ (Needs: test:unit, test:vet)
                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              STAGE: BUILD                                    │
│                         (All jobs run in parallel)                           │
│                                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │build:linux   │  │build:windows │  │build:macos   │  │build:macos   │   │
│  │    -x64      │  │    -x64      │  │   -arm64     │  │    -x64      │   │
│  │              │  │              │  │              │  │              │   │
│  │ • Native     │  │ • MinGW      │  │ • Clang      │  │ • Clang      │   │
│  │ • libpcap    │  │ • Cross-comp │  │ • Cross-comp │  │ • Cross-comp │   │
│  │              │  │              │  │              │  │              │   │
│  │ REQUIRED ✓   │  │ Optional ⚠   │  │ Optional ⚠   │  │ MANUAL 🔵    │   │
│  │ ~9 MB        │  │ ~9 MB        │  │ ~9 MB        │  │ ~9 MB        │   │
│  │ 30 days      │  │ 30 days      │  │ 30 days      │  │ 30 days      │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘   │
│                                                                              │
│  Each binary includes:                                                       │
│  • -X main.Version=${CI_COMMIT_TAG:-dev}                                   │
│  • -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)                        │
│  • -X main.GitCommit=${CI_COMMIT_SHA:0:8}                                  │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      │ (Only on: tags, master, main)
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                             STAGE: PACKAGE                                   │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐    │
│  │ package:all                                                         │    │
│  │ • Collect all binaries                                              │    │
│  │ • Generate SHA256SUMS.txt                                           │    │
│  │ • Create README.txt                                                 │    │
│  │ • Create release/ directory                                         │    │
│  │ ✓ Artifacts: release/ (90 days)                                    │    │
│  └────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐    │
│  │ release (Only on git tags)                                          │    │
│  │ • Create GitLab Release                                             │    │
│  │ • Attach all artifacts                                              │    │
│  │ • Generate release notes                                            │    │
│  └────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Pipeline Flow by Trigger

### On Every Push (Feature Branch)

```
download-dependencies
         │
         ▼
    ┌────┴────┐
    │  TESTS  │ (5 parallel jobs)
    └────┬────┘
         ▼
build:linux-x64 (Required)
build:windows-x64 (Attempt)
build:macos-arm64 (Attempt)

Status: ✅ if Linux builds + tests pass
```

### On Push to master/main

```
download-dependencies
         │
         ▼
    ┌────┴────┐
    │  TESTS  │ (5 parallel jobs)
    └────┬────┘
         ▼
    ┌────┴────┐
    │ BUILDS  │ (3 parallel jobs)
    └────┬────┘
         ▼
   package:all
   (Creates release package)

Status: ✅ Creates release/ with all successful builds
```

### On Git Tag Push (v1.0.0)

```
download-dependencies
         │
         ▼
    ┌────┴────┐
    │  TESTS  │ (5 parallel jobs)
    └────┬────┘
         ▼
    ┌────┴────┐
    │ BUILDS  │ (3 parallel jobs)
    └────┬────┘
         ▼
   package:all
         │
         ▼
    release
  (Creates GitLab Release)

Status: ✅ Full release with all artifacts
```

## Job Dependencies Graph

```
download-dependencies
         │
         ├─────────────┬─────────────┬─────────────┬─────────────┐
         ▼             ▼             ▼             ▼             ▼
    test:unit     test:vet     test:fmt    test:lint  test:security
         │             │             │             │             │
         └──────┬──────┘             │             │             │
                │                    │             │             │
         ┌──────┴──────┬─────────────┴─────────────┴─────────────┘
         ▼             ▼             ▼             ▼
   build:linux  build:windows build:macos  build:macos
     -x64         -x64        -arm64        -x64 (manual)
         │             │             │             │
         └──────┬──────┴──────┬──────┴─────────────┘
                │             │
                ▼             ▼
           package:all    release (tags only)
```

## Cache Flow

```
First Pipeline Run:
┌──────────┐
│ No Cache │
└────┬─────┘
     │
     ▼
┌──────────────────────────┐
│ Download all modules     │ (~2 min)
│ Build from scratch       │ (~3 min)
└────┬─────────────────────┘
     │
     ▼
┌──────────────────────────┐
│ Save to cache:           │
│ • .go/pkg/mod/          │
│ • .cache/go-build/      │
└──────────────────────────┘

Subsequent Runs (Same Branch):
┌──────────┐
│  Cache   │
│ HIT! ✓   │
└────┬─────┘
     │
     ▼
┌──────────────────────────┐
│ Use cached modules       │ (~30 sec)
│ Incremental build        │ (~1 min)
└──────────────────────────┘

Speed improvement: ~50%
```

## Artifact Flow

```
Test Stage:
  coverage.txt ──────────────┐
  gosec-report.json ─────────┤
                              │
Build Stage:                  │
  pcap_stripper-linux-amd64 ─┤
  pcap_stripper-windows.exe ─┤──▶ Available for
  pcap_stripper-macos-arm64 ─┤    download (30 days)
  pcap_stripper-macos-amd64 ─┘

Package Stage:
  release/
    ├─ pcap_stripper-linux-amd64
    ├─ pcap_stripper-windows-amd64.exe
    ├─ pcap_stripper-macos-arm64
    ├─ pcap_stripper-macos-amd64
    ├─ SHA256SUMS.txt
    └─ README.txt

  ▼
  GitLab Release (90 days)
```

## Timing Diagram

```
Time: 0s ──────────────────────────────────────────────▶ 10m

PREPARE  [██]                                   ~30s
           │
TEST     ──[████████]                          ~2-4m
           │ (parallel)
           │
BUILD    ──────────[██████]                    ~3-5m
                   │ (parallel)
                   │
PACKAGE  ──────────────────[██]                ~1m

Total: ~6-10 minutes
```

## Parallel Execution Visualization

```
Sequential (OLD):                      Parallel (NEW):
────────────────────                  ────────────────────
test:unit      [████]                 test:unit    [████]
test:vet          [██]                test:vet     [██]
test:fmt            [█]               test:fmt     [█]
test:lint              [████]         test:lint    [████]
test:security             [███]       test:security[███]
                                      ─────────────────────
Total: ~10 min                        Total: ~4 min ✓

build:linux    [███]                  build:linux  [███]
build:windows     [███]               build:windows[███]
build:macos          [███]            build:macos  [███]
                                      ─────────────────────
Total: ~9 min                         Total: ~3 min ✓

Overall Savings: ~12 minutes per run!
```

## Success/Failure Decision Tree

```
                    Pipeline Start
                         │
                         ▼
                  Dependencies OK?
                    │          │
                   Yes        No ──▶ FAIL ❌
                    │
                    ▼
              ┌────────────┐
              │   TESTS    │
              └────────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
   Unit Tests   go vet     gofmt
        │           │           │
       OK?         OK?         OK?
        │           │           │
       Yes         Yes         Yes
        └───────────┴───────────┘
                    │
               All Passed?
                │        │
               Yes       No ──▶ FAIL ❌
                │
                ▼
          Build Linux?
                │        │
               Yes       No ──▶ FAIL ❌
                │
                ▼
          ┌──────────┐
          │ OPTIONAL │
          │  BUILDS  │
          └──────────┘
                │
        Windows/macOS builds
        (failures allowed)
                │
                ▼
           SUCCESS ✅
                │
                ▼
      On master/tag? ───No──▶ END
                │
               Yes
                ▼
         Create Package
                │
                ▼
           Is Tag? ───No──▶ END
                │
               Yes
                ▼
        Create Release
                │
                ▼
              DONE
```

## Resource Usage Over Time

```
CPU Usage:
100% ├────────┐         ┌────────┐
     │        │         │        │
 75% │     ┌──┴─┐    ┌──┴─┐      │
     │     │    │    │    │      │
 50% │  ┌──┤    ├────┤    ├──────┤
     │  │  │    │    │    │      │
 25% ├──┤  └────┘    └────┘      └──
     │  │
  0% └──┴──────────────────────────▶
     Prep Test  Idle Build Idle Pkg

Memory Usage:
4GB  ├────────┐         ┌────────┐
     │        │         │        │
3GB  │     ┌──┴─┐    ┌──┴─┐      │
     │     │    │    │    │      │
2GB  │  ┌──┤    ├────┤    ├──────┤
     │  │  │    │    │    │      │
1GB  ├──┤  └────┘    └────┘      └──
     │  │
 0GB └──┴──────────────────────────▶
     Prep Test  Idle Build Idle Pkg
```

## Error Propagation

```
┌─────────────┐
│ Test Failed │
└──────┬──────┘
       │
       ▼
  Stop Pipeline ─────▶ No Builds Run
       │
       ▼
  Email Notification
       │
       ▼
  Pipeline Status: FAILED ❌


┌──────────────┐
│ Build Failed │
│  (Windows)   │
└──────┬───────┘
       │
       ▼
  Continue ─────▶ Linux build OK
       │
       ▼
  Warning Logged
       │
       ▼
  Pipeline Status: SUCCESS ✅ (with warnings)
```
