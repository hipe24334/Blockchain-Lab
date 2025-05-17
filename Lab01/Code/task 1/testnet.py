import os
import struct
import ecdsa
import hashlib
import base58
from bit import *
from bitcoin import *
from bit import PrivateKeyTestnet
from bitcoin.wallet import CBitcoinSecret, P2PKHBitcoinAddress


# Generate a random private key
while True:
    # Create random private key
    p_key = CBitcoinSecret.from_secret_bytes(os.urandom(32))
    
    # Take secret scalar behalfs bytes
    secret_bytes = p_key.to_bytes()
    
    # Turn secret scalar into hexa
    private_key = secret_bytes.hex()
    
    #Limit the length of key smaller 64 bytes
    if 0 < int.from_bytes(secret_bytes, 'big') < 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364140:
        break

# Derive the public key and Bitcoin address
public_key = p_key.pub

address = struct.pack('=B',111) +  P2PKHBitcoinAddress.from_pubkey(public_key)
checksum = hashlib.sha256(hashlib.sha256(address).digest()).digest()[:4]
binary_addr = address + checksum
bitcoin_addr = base58.b58encode(binary_addr)


print("Private Key:", private_key)
print("Public Key:", public_key.hex())
print("Bitcoin Address:", bitcoin_addr)

# Enter the txid and output_index
#After using website Bitcoin Testnet Faucet (fill BItcoin Address). We enter the transaction id as txid
txid = 'a7b8f9dcff1ccd78cc78bcf160b0d5c48747ef114566e8ff21470173dd313ea2'
output_index = 0

# Create a transaction input (UTXO)
txIn = {'txid': txid, 'vout': output_index}

# Create the transaction with input and output in one step
# The private_key_hex must be related to the Bitcoin address used for taking the txid 
private_key_hex = '006ae44a44aaea21b135f99843d6c8469a52cc2dfe7587b3db55e236b300180701'
pri_key = PrivateKeyTestnet.from_bytes(bytes.fromhex(private_key_hex))

# Create a transaction output to the desired destination
destination_address = bitcoin_addr
amount_to_send = 0.000001 
txOut = [(destination_address, amount_to_send,'btc')]
tx = pri_key.create_transaction(outputs = txOut)

# Create the transaction
pri_key.sign_transaction(tx)

# Print raw transaction
print("Successful unlock funds!")
print("The raw transaction is: ")
print(tx)












