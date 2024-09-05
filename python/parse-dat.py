import datetime
import glob

def parse_blocks(blks: list[bytes], block_height_start: int=0, block_height_end: int=-1, input_remainder: bytes='') -> list[dict, int, str]:
    blocks = {}
    block_remainder = ''
    # block_height = 0
    blks[0] = input_remainder + blks[0].hex()
    i = 0

    for block in blks:
        data_part = None
        if i == 0:
            data_part = block
            i += 1
        else:
            data_part = block.hex()
        data = block_remainder + data_part
        # split on magic number
        raw_blks = data.split('f9beb4d9')
        block_remainder = 'f9beb4d9' + raw_blks[-1]
        # loop through blocks from split to parse them individually
        for raw_blk in raw_blks[1:-1]:
            tmp_block = {}
            # print(f"i: {i}, b: {b}")
            block_size_raw = raw_blk[:8]
            block_size_swap = bytearray.fromhex(block_size_raw)
            block_size_swap.reverse()
            block_size = (int(block_size_swap.hex(), 16))*2+8
            if len(raw_blk) != block_size:
                raise ValueError(f"error, expected len(raw_blk) == block_size but got {len(raw_blk)} == {block_size}")

            block_header = raw_blk[16:168]
            block_header_swap = bytearray.fromhex(block_header)
            block_header_swap.reverse()
            # print(f"block header hex: {block_header_swap.hex()}")
            # print(f"length of block header: {len(block_header)}")

            version_dat = raw_blk[8:16]
            version_dat_swap = bytearray.fromhex(version_dat)
            version_dat_swap.reverse()
            tmp_block['version'] = int(version_dat_swap.hex(), 16)
            # print(f"version: {int(version_dat_swap.hex(), 16)}")

            prev_block = block_header[:64]
            prev_block_swap = bytearray.fromhex(prev_block)
            prev_block_swap.reverse()
            tmp_block['prev block'] = prev_block_swap.hex()
            # print(f"prev block: {int(prev_block_swap.hex(), 16)}")

            merkle_root = block_header[64:128]
            merkle_root_swap = bytearray.fromhex(merkle_root)
            merkle_root_swap.reverse()
            tmp_block['merkle root'] = merkle_root_swap.hex()
            # print(f"merkle root: {merkle_root_swap.hex()}")

            block_time = block_header[128:136]
            block_time_swap = bytearray.fromhex(block_time)
            block_time_swap.reverse()
            block_time_utc = datetime.datetime.fromtimestamp(int(block_time_swap.hex(), 16))
            tmp_block['timestamp'] = block_time_utc
            # print(f"block time: {block_time_utc}")

            block_bits = block_header[136:144]
            block_bits_swap = bytearray.fromhex(block_bits)
            block_bits_swap.reverse()
            tmp_block['bits'] = block_bits_swap.hex()
            # print(f"block bits swap: {block_bits_swap.hex()}")

            block_nonce = block_header[144:152]
            block_nonce_swap = bytearray.fromhex(block_nonce)
            block_nonce_swap.reverse()
            tmp_block['nonce'] = block_nonce_swap.hex()
            # print(f"block nonce swap: {int(block_nonce_swap.hex(), 16)}")

            # tx_ids_raw = raw_blk[168:block_size]
            tmp_block['tx ids raw '] = raw_blk[168:block_size]
            # print(f"tx ids raw: {tx_ids_raw}")
            blocks[block_height_start] = tmp_block
            block_height_start += 1
        if block_height_start >= block_height_end and block_height_end >= 0:
            break
    return blocks, block_height_start-1, block_remainder

glob_path  = "C:\\Users\\david\\OneDrive\\Documents\\code\\python\\Blockchain\\Bitcoin\\data\\bitcoin_data\\*.dat"
glob_paths = glob.glob(glob_path)
parsed_blocks, block_height, block_remainder = None, None, None

for i, dat_file in enumerate(glob_paths):
    print(f"parsing file: blk0000{i}.dat")
    blk = None
    with open(dat_file, 'rb') as f:
        blk = f.readlines()
        f.close()
    if dat_file == 'C:\\Users\\david\\OneDrive\\Documents\\code\\python\\Blockchain\\Bitcoin\\data\\bitcoin_data\\blk00000.dat':
        parsed_blocks, block_height, block_remainder = parse_blocks(blks=blk)
    else:
        parsed_blocks, block_height, block_remainder = parse_blocks(blks=blk, block_height_start=block_height,  input_remainder=block_remainder)
    print(f"{(block_height):,d}")