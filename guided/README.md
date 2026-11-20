# SpeedDisc Guided • preview / pré-visualização

**Not released. Do not distribute the preview executable. Standard-user Windows acceptance is still required.**

An isolated Go CLI/TUI companion, not a replacement for the v0.4 maintenance executable. Keep the legacy maintenance release and its permissions separate. This module never calls or imports its repair code. No website changes are included.

## One portable entry point

Run `SpeedDisc-Guided-preview.exe` in a normal Windows account. Choose PT/EN with L; choose a problem; review human-readable ranked evidence and next steps; optionally save JSON to a new path. Q exits. B returns from folder selection. EOF exits cleanly. Ctrl+C cancels a scan. No installer, elevation request, execution-policy change or admin bootstrap exists.

- **Disco cheio / Disk full:** select one local folder. Metadata-only scan, top ten files, logical bytes, visited file count, denied/skipped count. Excludes CrossDevice (case-insensitive), symbolic links, Windows reparse points, UNC/device paths and non-fixed drives. Uses Go `os.Root` to confine resolution. Stops after 30 seconds or 100,000 files. No file contents, cleanup, recycle-bin access or cache deletion. Counts are partial when paths are excluded/unavailable; logical sizes are not reclaimable physical storage. Do not concurrently rearrange the directory while scanning; results are snapshots, not a security audit.
- **Sem internet / No internet:** local adapter names, link flags and MTU only. No IP/MAC/SSID, Wi-Fi secrets, DNS requests, ping, external connections or network scanning. Link-up is not proof of working internet.
- **Arranque lento / Slow startup:** fixed read-only Windows CIM query through the system PowerShell executable, 10-second timeout, no profiles, hidden child window. Captures startup names only, never commands or arguments. Does not disable entries, inspect passwords, change registry or measure boot time. Policy/provider failure is reported, not bypassed. Inventory is explicitly partial.

Only an explicitly requested JSON report is written. Windows report destinations must be on fixed local drives with no reparse parent; UNC/device paths and CrossDevice are rejected. Reports use a confined parent-directory handle. Existing reports are never overwritten. Reports can include private filenames/application names: review before sharing. Errors are recorded; an unavailable collector is not reported as successful. No telemetry or automatic uploads.

## Build and test

Requires official Go 1.25+; no external module dependencies.

```
go test ./...
go vet ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o SpeedDisc-Guided-preview.exe .
```

This independent module deliberately does not reuse the old admin manifest or legacy build script. Default Go Windows console binary does not request elevation. An explicit `asInvoker` application manifest and signed release packaging should be reviewed during release engineering; do not copy `requireAdministrator` from the old app. Do not ask people to bypass Windows protection for an unsigned preview.

## Required Windows acceptance before any release

Disposable standard account on supported Windows 10/11; no Daniel-PC installation:

1. Start without UAC/elevation; exercise PT/EN, invalid selection, back, repeated runs, EOF/Ctrl+C, read-only/unwritable report destinations and existing report path.
2. Fixture tree: regular/nested files, denied directories, directory junction, symlink, mounted volume, CrossDevice (root and child), Unicode names, long paths. Verify partial counts and zero outside-root visits, no content changes, no network traversal. Include concurrent reparse replacement tests; do not claim TOCTOU hardening from Linux tests alone.
3. Offline/airplane mode, wired and Wi-Fi adapters, restricted CIM/PowerShell policies and timeout. Confirm no bypass and no unexpected traffic.
4. Startup names compared with Task Manager and known fixtures. Document omissions; never infer measured startup impact from presence alone.
5. Verify report privacy and PE manifest with standard-user execution. Cross-compilation is not a Windows runtime pass. No public download until these checks pass.

## Estado / Próximos passos

Implemented locally only. Fixture tests and cross-compilation are cloud checks. Windows runtime, junction/ACL and UI/console encoding acceptance remain pending. Parent task owns integration into repository `guided/`, review, publication and release decisions.
