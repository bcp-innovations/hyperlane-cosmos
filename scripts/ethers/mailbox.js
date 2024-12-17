import { ethers } from "ethers";
import {config} from "dotenv";

config()

const provider = new ethers.JsonRpcProvider('https://rpc.api.moonbase.moonbeam.network', {
    chainId: 1287,
    name: 'moonbase-alpha',
});
let wallet = new ethers.Wallet(process.env.HYP_KEY, provider);

console.log("Current block-height: ", await provider.getBlockNumber())



// Sample send function
const send = async () => {

    const addressTo = '0x0c60e7eCd06429052223C78452F791AAb5C5CAc6';
    console.log(`Attempting to send transaction from ${wallet.address} to ${addressTo}`);

    // 5. Create tx object
    const tx = {
        to: addressTo,
        value: ethers.parseEther('0.001'),
    };

    // 6. Sign and send tx - wait for receipt
    const createReceipt = await wallet.sendTransaction(tx);
    await createReceipt.wait();
    console.log(`Transaction successful with hash: ${createReceipt.hash}`);
};

const abi = [
    // Read-Only Functions
    "function process(bytes calldata metadata, bytes calldata message) external payable",
    "function localDomain() external view returns (uint32)",
    "function latestDispatchedId() external view returns (bytes32)"
];

// This can be an address or an ENS name
const address = "0x3fc14e89a3073682730aAa82c0F29Dca147bF414";

// Read-Write; By connecting to a Signer, allows:
// - Everything from Read-Only (except as Signer, not anonymous)
// - Sending transactions for non-constant functions
const erc20_rw = new ethers.Contract(address, abi, wallet);
await erc20_rw.connect(wallet)

console.log(await erc20_rw.latestDispatchedId())
console.log(await erc20_rw.process("0x1234", "0x030000000300000001b32677d8121a50c7b960b8561ead86278a7d75ec786807983e1eebfcbc2d9cfc00000507000000000000000000000000f254e1ce6b468e5c118214d13faa2630110467150000000000000000000000000c60e7ecd06429052223c78452f791aab5c5cac6000000000000000000000000000000000000000000000000000000000000000b"))
