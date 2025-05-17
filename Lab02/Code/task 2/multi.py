import os
import struct
import hashlib
import base58
import string
from bitcoin.wallet import CBitcoinSecret, P2SHBitcoinAddress
from bitcoin.core.script import CScript, OP_2, OP_CHECKMULTISIG

# Generate two random private keys
private_key1 = CBitcoinSecret.from_secret_bytes(os.urandom(32))
private_key2 = CBitcoinSecret.from_secret_bytes(os.urandom(32))

# Derive the public keys
public_key1 = private_key1.pub
public_key2 = private_key2.pub

# Create a 2-of-2 multisig redeem script
redeem_script = CScript([OP_2, public_key1, public_key2, OP_2, OP_CHECKMULTISIG])

address = struct.pack('=B',111) +  P2SHBitcoinAddress.from_redeemScript(redeem_script)
checksum = hashlib.sha256(hashlib.sha256(address).digest()).digest()[:4]
binary_addr = address + checksum
bitcoin_addr = base58.b58encode(binary_addr)
bitcoin_address = bitcoin_addr.decode('utf-8')

print("Private Key 1:", private_key1)
print("Private Key 2:", private_key2)

print("Redeem Script:", redeem_script.hex())
print("Multisig Address:", bitcoin_address)
