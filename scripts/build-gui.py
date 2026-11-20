#!/usr/bin/env python3
"""Build and check the independent unsigned Windows GUI without altering CLI resources."""

import argparse
import hashlib
import os
from pathlib import Path
import shutil
import struct
import subprocess
import tempfile
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parent.parent
WINRES = "github.com/tc-hib/go-winres@v0.3.3"


def verify_gui(path):
    data = path.read_bytes()

    def u16(i):
        return struct.unpack_from("<H", data, i)[0]

    def u32(i):
        return struct.unpack_from("<I", data, i)[0]

    def require(condition, message):
        if not condition:
            raise ValueError(message)

    require(data[:2] == b"MZ", "not a PE file")
    pe = u32(60)
    require(data[pe : pe + 4] == b"PE\0\0", "missing PE signature")
    opt = pe + 24
    require(u16(pe + 4) == 0x8664 and u16(opt) == 0x20B, "not PE32+ amd64")
    require(u16(opt + 68) == 2, "subsystem is not IMAGE_SUBSYSTEM_WINDOWS_GUI")
    require(u16(pe + 20) >= 136 and u32(opt + 108) >= 3, "no PE resource directory")
    res_rva, res_size = u32(opt + 128), u32(opt + 132)
    require(res_rva and res_size, "no PE resources")
    sections = []
    for i in range(u16(pe + 6)):
        section = opt + u16(pe + 20) + i * 40
        sections.append((u32(section + 12), u32(section + 16), u32(section + 20)))

    def pos(rva, size):
        for va, raw_size, raw in sections:
            if va <= rva and size <= raw_size - (rva - va):
                offset = raw + (rva - va)
                if offset + size <= len(data):
                    return offset
        raise ValueError(f"invalid resource RVA {rva:#x}")

    base = pos(res_rva, res_size)

    def bounded(relative, size):
        require(0 <= relative <= res_size - size, "out-of-bounds resource item")
        return base + relative

    dir_bit = 1 << 31

    def entries(relative):
        header = bounded(relative, 16)
        count = u16(header + 12) + u16(header + 14)
        body = bounded(relative + 16, count * 8)
        return [(u32(body + n * 8), u32(body + n * 8 + 4)) for n in range(count)]

    roots = dict(entries(0))
    require(all(roots.get(t, 0) & dir_bit for t in (3, 14, 24)), "icon or manifest resource absent")
    manifest_names = dict(entries(roots[24] & ~dir_bit))
    require(set(manifest_names) == {1} and manifest_names[1] & dir_bit, "expected one application manifest")
    manifests = entries(manifest_names[1] & ~dir_bit)
    require(len(manifests) == 1, "expected one manifest language")
    language, leaf = manifests[0]
    require(not language & dir_bit and not leaf & dir_bit, "invalid manifest resource leaf")
    offset = bounded(leaf, 16)
    size = u32(offset + 4)
    require(size > 0, "empty application manifest")
    start = pos(u32(offset), size)
    manifest = ET.fromstring(data[start : start + size])
    levels = manifest.findall(".//{urn:schemas-microsoft-com:asm.v3}requestedExecutionLevel")
    require(len(levels) == 1 and levels[0].attrib == {"level": "asInvoker", "uiAccess": "false"},
            "manifest must be asInvoker with uiAccess=false")
    dpi = manifest.find(".//{http://schemas.microsoft.com/SMI/2016/WindowsSettings}dpiAwareness")
    require(dpi is not None and dpi.text and dpi.text.split(",")[0].strip() == "PerMonitorV2",
            "missing PerMonitorV2 manifest setting")
    print(f"Verified: Windows amd64 GUI / asInvoker / icon / PerMonitorV2 / {len(data):,} bytes", flush=True)
    print(f"SHA-256: {hashlib.sha256(data).hexdigest()}", flush=True)


def build(output):
    if output.name.lower() != "speeddisc-gui.exe":
        raise ValueError("output filename must be SpeedDisc-GUI.exe")
    output.parent.mkdir(parents=True, exist_ok=True)
    # Fresh tree: the untaggable winres .syso for the conventional CLI can
    # already be present in the checkout. Never include or mutate that object.
    with tempfile.TemporaryDirectory(prefix="speeddisc-gui-") as tmp:
        stage = Path(tmp)
        for source in (*ROOT.glob("*.go"), ROOT / "go.mod", ROOT / "go.sum"):
            shutil.copy2(source, stage)
        if (ROOT / "internal").exists():
            shutil.copytree(ROOT / "internal", stage / "internal",
                            ignore=shutil.ignore_patterns("*.syso", "*.exe", "*.test", "__pycache__"))
        (stage / "winres").mkdir()
        # Copy only source assets, including any native UI go:embed icon.
        for name in ("speeddisc.ico", "icon.png", "speeddisc-gui.manifest", "gui-winres.json"):
            shutil.copy2(ROOT / "winres" / name, stage / "winres")
        host_env = os.environ.copy()
        for name in ("GOOS", "GOARCH", "GOFLAGS"):
            host_env.pop(name, None)
        subprocess.run(
            ["go", "run", WINRES, "make", "--arch", "amd64", "--in", "winres/gui-winres.json", "--out", "gui_rsrc"],
            check=True, cwd=stage, env=host_env,
        )
        exe = stage / "SpeedDisc-GUI.exe"
        cross_env = host_env | {"CGO_ENABLED": "0", "GOOS": "windows", "GOARCH": "amd64"}
        subprocess.run(
            ["go", "build", "-tags", "speeddiscgui", "-trimpath", "-buildvcs=false", "-ldflags", "-H windowsgui -s -w",
             "-o", str(exe), "."],
            check=True, cwd=stage, env=cross_env,
        )
        verify_gui(exe)
        shutil.copy2(exe, output)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    action = parser.add_mutually_exclusive_group(required=True)
    action.add_argument("--output", type=Path, help="place a new build here; filename must be SpeedDisc-GUI.exe")
    action.add_argument("--verify", type=Path, help="check an existing GUI PE without starting it")
    args = parser.parse_args()
    try:
        if args.output:
            build(args.output.resolve())
        else:
            verify_gui(args.verify)
    except (OSError, ValueError, struct.error, ET.ParseError, subprocess.CalledProcessError) as e:
        parser.exit(1, f"GUI build check failed: {e}\n")


if __name__ == "__main__":
    main()
