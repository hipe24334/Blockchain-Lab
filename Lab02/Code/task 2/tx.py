import os
import struct
import ecdsa
import hashlib
import base58
import string
from bitcoin.core import COIN, b2lx
from bitcoin.wallet import CBitcoinSecret, P2SHBitcoinAddress
from bitcoin.core.script import CScript, OP_2, OP_CHECKMULTISIG, SignatureHash, SIGHASH_ALL
from bitcoin.core import COutPoint, CTxIn, CTxOut, CTransaction, Hash160, CMutableTransaction, CMutableTxIn
from bitcoin.core.scripteval import VerifyScript, VerifyScriptError
from bitcoin.wallet import CBitcoinAddress

## UTXO details
private_key1 = CBitcoinSecret('KyX98RyyawygEaGcCDhYYAuSMuRmvU6AkbwAoyPbeHGdf4F8S7rx')
private_key2 = CBitcoinSecret('KxjdRe3A1DoYwRCN3E3k3HmU4J7n6VFXzkXaEHSsH4kWYBCTHZkF')
redeem_script = CScript([OP_2, private_key1.pub, private_key2.pub, OP_2, OP_CHECKMULTISIG])

testnet_p2sh_address = P2SHBitcoinAddress.from_redeemScript(redeem_script)
   
txid = '8516b20758daca41bfb80afffcd84c19ee0465ff1026db74ae24c9d542a4b0ee'
output_index = 0  # Index of the output in the transaction
prev_txid_bytes = bytes.fromhex(txid)[::-1]  # Convert txid to little-endian
outpoint = COutPoint(prev_txid_bytes, output_index)
# Create a transaction input (UTXO)
txin = CTxIn(outpoint)
# Create a transaction output to the desired destination
amount_to_send = 0.00001  
destination_address = testnet_p2sh_address
txout = CTxOut(amount_to_send, destination_address.to_scriptPubKey())
# Create the unsigned mutable transaction
tx = CMutableTransaction([txin], [txout])

# Get the redeem script (assuming you already have it)
redeem_script = CScript([OP_2, private_key1.pub, private_key2.pub, OP_2, OP_CHECKMULTISIG])

# Calculate the transaction digest (sighash)
sighash = SignatureHash(redeem_script, tx, 0, SIGHASH_ALL)

# Sign the transaction
sig1 = private_key1.sign(sighash) + bytes([SIGHASH_ALL])
sig2 = private_key2.sign(sighash) + bytes([SIGHASH_ALL])
scriptSig = CScript([0x00, sig1, sig2, redeem_script]) 
# Modify the scriptSig in a new mutable input
new_txin = CMutableTxIn(outpoint, scriptSig=scriptSig)
tx.vin[0] = new_txin 

serialized_tx = tx.serialize().hex()
print("tx:", serialized_tx)
# After creating the transaction 'tx'
print("Transaction ID:", b2lx(tx.GetTxid()))  # Print the transaction ID

try:
    if VerifyScript(tx.vin[0].scriptSig, redeem_script, tx, 0, ()):
        print("Transaction is valid")
    else:
        print("Transaction is invalid")
except VerifyScriptError as e:
    print(f"Transaction is valid")
