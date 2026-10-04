#!/usr/bin/env python3
"""Check that SpeedDisc.exe is a 64-bit console program with requireAdministrator."""
import struct
import sys

def main() -> None:
    if len(sys.argv) != 2:
        sys.exit("usage: verify_pe.py SpeedDisc.exe")
    data = open(sys.argv[1], "rb").read()
    if data[:2] != b"MZ":
        sys.exit("not a PE file")
    e = struct.unpack_from("<I", data, 0x3C)[0]
    if data[e:e + 4] != b"PE\0\0":
        sys.exit("missing PE signature")
    machine = struct.unpack_from("<H", data, e + 4)[0]
    magic = struct.unpack_from("<H", data, e + 24)[0]
    if magic != 0x20B:
        sys.exit(f"expected PE32+ magic, got {magic:#x}")
    subsystem = struct.unpack_from("<H", data, e + 24 + 68)[0]
    print(f"machine={machine:#x} subsystem={subsystem} bytes={len(data)}")
    if machine != 0x8664:
        sys.exit("not amd64")
    if subsystem != 3:
        sys.exit("subsystem is not IMAGE_SUBSYSTEM_WINDOWS_CUI (console)")
    needle = b"requireAdministrator"
    if needle not in data and needle.decode().encode("utf-16le") not in data:
        sys.exit("requireAdministrator was not found in the executable")
    if b"requestedExecutionLevel" not in data and "requestedExecutionLevel".encode("utf-16le") not in data:
        sys.exit("requestedExecutionLevel was not found in the executable")
    print("manifest: requireAdministrator")
    print("subsystem: console")
    if not has_icon(data, e):
        sys.exit("Windows icon resource (RT_GROUP_ICON / RT_ICON) was not found")
    print("icon: RT_GROUP_ICON")


def has_icon(data: bytes, e: int) -> bool:
    """True when the PE resource tree contains an icon group and icon images."""
    nsec = struct.unpack_from("<H", data, e + 6)[0]
    optsz = struct.unpack_from("<H", data, e + 20)[0]
    opt = e + 24
    if struct.unpack_from("<H", data, opt)[0] != 0x20B:
        return False
    ndd = struct.unpack_from("<I", data, opt + 108)[0]
    if ndd < 3:
        return False
    res_rva, res_size = struct.unpack_from("<II", data, opt + 112 + 16)
    if res_rva == 0 or res_size == 0:
        return False
    sections = []
    sec_off = e + 24 + optsz
    for i in range(nsec):
        o = sec_off + i * 40
        vsz, va, rsz, raw = struct.unpack_from("<IIII", data, o + 8)
        sections.append((va, max(vsz, rsz), raw))

    def off_of(rva: int) -> int:
        for va, span, raw in sections:
            if va <= rva < va + max(span, 1):
                return raw + (rva - va)
        raise ValueError(f"rva {rva:#x}")

    try:
        root = off_of(res_rva)
    except ValueError:
        return False
    types = []

    def walk(off: int, level: int) -> None:
        named, ids = struct.unpack_from("<HH", data, off + 12)
        ent = off + 16
        for i in range(named + ids):
            name, rel = struct.unpack_from("<II", data, ent + i * 8)
            is_dir = rel & 0x80000000
            child = root + (rel & 0x7FFFFFFF)
            if level == 0 and (name & 0x80000000) == 0:
                types.append(name)
            if is_dir:
                walk(child, level + 1)

    try:
        walk(root, 0)
    except struct.error:
        return False
    # 3 = RT_ICON, 14 = RT_GROUP_ICON
    return 3 in types and 14 in types

if __name__ == "__main__":
    main()
