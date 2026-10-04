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

if __name__ == "__main__":
    main()
