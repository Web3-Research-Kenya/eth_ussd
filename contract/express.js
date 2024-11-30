// Import web3.js
import Web3 from 'web3';
import express from 'express';
import bodyParser from 'body-parser';

// Set up Alchemy RPC URL
const alchemyURL = 'https://eth-sepolia.g.alchemy.com/v2/frEh7TZEZEBIwyMsAx58Aeypq34caZ2s';

// Instantiate a web3 instance
const web3 = new Web3(new Web3.providers.HttpProvider(alchemyURL));

// Initialize Express app
const app = express();
app.use(bodyParser.json());

// Function to create a new Ethereum account
const createAccount = () => {
    const newAccount = web3.eth.accounts.create();
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
        return receipt.transactionHash;
    } catch (error) {
        throw new Error('Error sending ETH: ' + error.message);
    }
};

// Function to get the balance of an account
const getBalance = async (address) => {
    try {
        const balance = await web3.eth.getBalance(address);
        return web3.utils.fromWei(balance, 'ether');
    } catch (error) {
        throw new Error('Error getting balance: ' + error.message);
    }
};

// API Endpoint to create a new account
app.post('/create-account', (req, res) => {
    try {
        const account = createAccount();
        res.json({ address: account.address, privateKey: account.privateKey });
    } catch (error) {
        res.status(500).send('Error creating account');
    }
});

// API Endpoint to send ETH
app.post('/send-eth', async (req, res) => {
    const { senderPrivateKey, receiverAddress, amountInEther } = req.body;
    try {
        const txHash = await sendETH(senderPrivateKey, receiverAddress, amountInEther);
        res.json({ transactionHash: txHash });
    } catch (error) {
        res.status(500).send(error.message);
    }
});

// API Endpoint to get account balance
app.get('/get-balance/:address', async (req, res) => {
    const address = req.params.address;
    try {
        const balance = await getBalance(address);
        res.json({ address: address, balance: balance + ' ETH' });
    } catch (error) {
        res.status(500).send(error.message);
    }
});

// Start the server
const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});