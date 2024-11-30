// Import web3.js
import Web3 from 'web3';

// Set up Alchemy RPC URL (replace YOUR_ALCHEMY_KEY with your actual key)
const alchemyURL = 'https://eth-sepolia.g.alchemy.com/v2/frEh7TZEZEBIwyMsAx58Aeypq34caZ2s';

// Instantiate a web3 instance
const web3 = new Web3(new Web3.providers.HttpProvider(alchemyURL));

// Function to create a new Ethereum account
const createAccount = () => {
    const newAccount = web3.eth.accounts.create();
    console.log('New Account Address:', newAccount.address);
    console.log('Private Key (Keep it secret!):', newAccount.privateKey);
    return newAccount;
};

// Function to send ETH from one account to another
const sendETH = async (senderPrivateKey, receiverAddress, amountInEther) => {
    try {
        // Get the account object from private key
        const senderAccount = web3.eth.accounts.privateKeyToAccount(senderPrivateKey);

        // Add sender account to web3 wallet
        web3.eth.accounts.wallet.add(senderAccount);

        // Get current gas price from the network
        const gasPrice = await web3.eth.getGasPrice();

        // Prepare the transaction object
        const tx = {
            from: senderAccount.address,
            to: receiverAddress,
            value: web3.utils.toWei(amountInEther, 'ether'),
            gas: 21000, // standard for ETH transfer
            gasPrice: gasPrice,
        };

        // Send the transaction
        const receipt = await web3.eth.sendTransaction(tx);
        console.log('Transaction successful with hash:', receipt.transactionHash);
    } catch (error) {
        console.error('Error sending ETH:', error);
    }
};

// Function to get the balance of an account
const getBalance = async (address) => {
    try {
        const balance = await web3.eth.getBalance(address);
        console.log(`Balance of ${address}:`, web3.utils.fromWei(balance, 'ether'), 'ETH');
        return balance;
    } catch (error) {
        console.error('Error getting balance:', error);
    }
};

// Example usage
(async () => {
    // Create a new account
    const account = createAccount();

    // Example: Get balance of an address (replace with an actual address)
    await getBalance(account.address);

    // Example: Send ETH (replace with your sender private key and receiver address)
    // const senderPrivateKey = 'YOUR_PRIVATE_KEY';
    // const receiverAddress = 'RECEIVER_ADDRESS';
    // await sendETH(senderPrivateKey, receiverAddress, '0.01');
})();
