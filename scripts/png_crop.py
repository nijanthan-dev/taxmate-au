#!/usr/bin/env python3
"""Crop 8-bit RGB/RGBA PNG files without third-party dependencies."""

import argparse
import binascii
import struct
import zlib
from pathlib import Path
from typing import Iterable, List, Tuple

PNG_SIGNATURE = b"\x89PNG\r\n\x1a\n"
COLOR_CHANNELS = {2: 3, 6: 4}


def read_chunks(raw: bytes) -> Iterable[Tuple[bytes, bytes]]:
    if not raw.startswith(PNG_SIGNATURE):
        raise ValueError("not a PNG file")
    offset = len(PNG_SIGNATURE)
    while offset < len(raw):
        if offset + 8 > len(raw):
            raise ValueError("truncated PNG chunk")
        length = struct.unpack(">I", raw[offset : offset + 4])[0]
        chunk_type = raw[offset + 4 : offset + 8]
        data_start = offset + 8
        data_end = data_start + length
        crc_end = data_end + 4
        if crc_end > len(raw):
            raise ValueError("truncated PNG chunk data")
        yield chunk_type, raw[data_start:data_end]
        offset = crc_end
        if chunk_type == b"IEND":
            break


def bytes_per_pixel(color_type: int, bit_depth: int) -> int:
    if bit_depth != 8 or color_type not in COLOR_CHANNELS:
        raise ValueError("only 8-bit RGB/RGBA PNG files are supported")
    return COLOR_CHANNELS[color_type]


def paeth(a: int, b: int, c: int) -> int:
    p = a + b - c
    pa = abs(p - a)
    pb = abs(p - b)
    pc = abs(p - c)
    if pa <= pb and pa <= pc:
        return a
    if pb <= pc:
        return b
    return c


def unfilter_scanlines(data: bytes, width: int, height: int, bpp: int) -> List[bytearray]:
    row_len = width * bpp
    rows: List[bytearray] = []
    offset = 0
    prior = bytearray(row_len)
    for _ in range(height):
        if offset + row_len + 1 > len(data):
            raise ValueError("truncated PNG image data")
        filter_type = data[offset]
        offset += 1
        row = bytearray(data[offset : offset + row_len])
        offset += row_len
        for index in range(row_len):
            left = row[index - bpp] if index >= bpp else 0
            up = prior[index]
            upper_left = prior[index - bpp] if index >= bpp else 0
            if filter_type == 1:
                row[index] = (row[index] + left) & 0xFF
            elif filter_type == 2:
                row[index] = (row[index] + up) & 0xFF
            elif filter_type == 3:
                row[index] = (row[index] + ((left + up) // 2)) & 0xFF
            elif filter_type == 4:
                row[index] = (row[index] + paeth(left, up, upper_left)) & 0xFF
            elif filter_type != 0:
                raise ValueError(f"unsupported PNG filter: {filter_type}")
        rows.append(row)
        prior = row
    return rows


def write_chunk(chunk_type: bytes, data: bytes) -> bytes:
    crc = binascii.crc32(chunk_type)
    crc = binascii.crc32(data, crc) & 0xFFFFFFFF
    return struct.pack(">I", len(data)) + chunk_type + data + struct.pack(">I", crc)


def crop_png(input_path: Path, output_path: Path, left: int, top: int, width: int, height: int) -> None:
    ihdr = None
    compressed = bytearray()
    for chunk_type, data in read_chunks(input_path.read_bytes()):
        if chunk_type == b"IHDR":
            ihdr = data
        elif chunk_type == b"IDAT":
            compressed.extend(data)
    if ihdr is None:
        raise ValueError("missing PNG IHDR")
    source_width, source_height, bit_depth, color_type, compression, filter_method, interlace = struct.unpack(">IIBBBBB", ihdr)
    if compression != 0 or filter_method != 0 or interlace != 0:
        raise ValueError("interlaced or nonstandard PNG files are not supported")
    if left < 0 or top < 0 or width <= 0 or height <= 0 or left + width > source_width or top + height > source_height:
        raise ValueError("crop rectangle outside image bounds")
    bpp = bytes_per_pixel(color_type, bit_depth)
    rows = unfilter_scanlines(zlib.decompress(bytes(compressed)), source_width, source_height, bpp)
    raw = bytearray()
    for row in rows[top : top + height]:
        start = left * bpp
        end = start + width * bpp
        raw.append(0)
        raw.extend(row[start:end])
    output_ihdr = struct.pack(">IIBBBBB", width, height, bit_depth, color_type, compression, filter_method, interlace)
    output = PNG_SIGNATURE + write_chunk(b"IHDR", output_ihdr) + write_chunk(b"IDAT", zlib.compress(bytes(raw), 9)) + write_chunk(b"IEND", b"")
    output_path.write_bytes(output)


def main() -> None:
    parser = argparse.ArgumentParser(description="Crop an 8-bit RGB/RGBA PNG.")
    parser.add_argument("input", type=Path)
    parser.add_argument("output", type=Path)
    parser.add_argument("left", type=int)
    parser.add_argument("top", type=int)
    parser.add_argument("width", type=int)
    parser.add_argument("height", type=int)
    args = parser.parse_args()
    crop_png(args.input, args.output, args.left, args.top, args.width, args.height)


if __name__ == "__main__":
    main()
