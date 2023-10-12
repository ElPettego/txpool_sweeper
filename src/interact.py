#!/usr/bin/env python3.8

from web3 import Web3
from web3.middleware import geth_poa
from w3_utils import get_keys, handle_tx
import argparse
import dis

def main():
    actions = ['front', 'back', 'withdraw']
    networks = {
        'polygon': 'https://polygon-bor.publicnode.com',
        'binance': 'https://bsc.publicnode.com'
    }

    parser = argparse.ArgumentParser(description='interact with smart contract')

    # parser.add_argument('action',           help=f'action to exec on smart contract: {actions}', choices=actions)
    parser.add_argument('network',          help='network to exec the tx')
    parser.add_argument('keys',             help='keys couple to exec transaction')
    parser.add_argument('contract_address', help='address of the deployed contract')
    parser.add_argument('abi_path',         help='path to abi of the deployed contract')
    parser.add_argument('token_address',    help='address of the token to do tha shit')
    parser.add_argument('--nonce',    '-n',  type=int,   default=-1,     help='nonce for the tx')
    parser.add_argument('--address1', '-a1', type=str,   default='0x',   help='address of the token1 to buy/sell')
    parser.add_argument('--address2', '-a2', type=str,   default='0x',   help='address of the token2 to buy/sell')
    parser.add_argument('--amount',   '-q',  type=float, default=0.0001, help='amount of the token to buy in eth')
    parser.add_argument('--gasprice', '-g',  type=int,   default=-1,     help='gas price for the tx')
    
    args = parser.parse_args()
    
    provider = networks[args.network]

    PRIVATE, PUBLIC = get_keys(args.keys)

    W3 = Web3(Web3.HTTPProvider(provider))
    W3.middleware_onion.inject(geth_poa.geth_poa_middleware, layer=0)
    
    contract_address = args.contract_address
    with open(args.abi_path, 'r') as f:
        abi = f.read()

    # print(abi)

    contract = W3.eth.contract(W3.to_checksum_address(contract_address), abi=abi)

    args.nonce = W3.eth.get_transaction_count(PUBLIC)
    args.gasprice = W3.eth.gas_price

    # print(contract)

    tx = contract.functions.doThaShit(
        W3.to_wei(args.amount, 'ether'),
        W3.to_checksum_address(args.token_address)
        ).build_transaction({
            'from':     PUBLIC,
            'gasPrice': args.gasprice,
            'nonce':    args.nonce,
        }
    )

    handle_tx(W3, tx, PRIVATE)


if __name__ == "__main__":
    dis.dis(main)
    # main()
