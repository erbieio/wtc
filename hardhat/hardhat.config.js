require('@nomicfoundation/hardhat-toolbox');
require('hardhat-deploy');

const INFURA_ID = '460f40a260564ac4a4f4b3fffb032dad';
const mnemonic = 'nothing across double orient color soap stomach hotel skirt bless true jelly';

// account list
// 0x394586580ff4170c8a0244837202cbabe9070f66 7bbfec284ee43e328438d46ec803863c8e1367ab46072f7864c07e0a03ba61fd
// 0x60f78431882884a68c94fe96ef45200749176d23 35a4c43ac13eab4c589b40828d757c37e52c00c2b9c4f51bc33cb690fff4d84a
// 0x981d431d2ac3aa32314021cec7b3b40cf7748049 43e7a174a002faf58dcc78ee6d6ced0f68094cfac6d38b2aa3537aa7dda7d105

module.exports = {
  defaultNetwork: 'localhost',

  networks: {
    hardhat: {
      chainId: 1337,
      accounts: {
        mnemonic,
      },
    },
    localhost: {
      url: 'http://127.0.0.1:8545',
      accounts: {
        mnemonic,
      },
    },
    matic: {
      url: 'https://rpc-mainnet.maticvigil.com',
      accounts: {
        mnemonic,
      },
    },
    mumbai: {
      url: 'https://rpc-mumbai.matic.today',
      accounts: {
        mnemonic,
      },
    },
    kovan: {
      url: `https://kovan.infura.io/v3/${INFURA_ID}`,
      accounts: {
        mnemonic,
      },
    },
    mainnet: {
      url: `https://mainnet.infura.io/v3/${INFURA_ID}`,
      accounts: {
        mnemonic,
      },
    },
    ropsten: {
      url: `https://ropsten.infura.io/v3/${INFURA_ID}`,
      accounts: {
        mnemonic,
      },
    },
    goerli: {
      url: `https://goerli.infura.io/v3/${INFURA_ID}`,
      accounts: {
        mnemonic,
      },
    },
  },
  solidity: {
    compilers: [
      {
        version: '0.8.21',
        settings: {
          optimizer: {
            enabled: true,
            runs: 20000,
          },
        },
      },
    ],
  },
  namedAccounts: {
    deployer: {
      default: 0,
    },
  },
};

// eslint-disable-next-line
task('wallet', 'Create a random account', async (_, { ethers }) => {
  const randomWallet = ethers.Wallet.createRandom();
  const privateKey = randomWallet._signingKey().privateKey;
  console.log('addr：' + randomWallet.address + '');
  console.log('key：' + privateKey);
});

// eslint-disable-next-line
task('accounts', 'Print a list of accounts on the node', async (_, { ethers }) => {
  const accounts = await ethers.provider.listAccounts();
  accounts.forEach((account) => console.log(account));
});

async function addr(ethers, addr, utils) {
  const { isAddress, getAddress } = utils;
  if (isAddress(addr)) {
    return getAddress(addr);
  }
  const accounts = await ethers.provider.listAccounts();
  if (accounts[addr] !== undefined) {
    return accounts[addr];
  }
  throw `not a canonical address: ${addr}`;
}

// eslint-disable-next-line
task('send', ' send ETH')
  .addParam('from', 'Sender address or account number on the node')
  .addOptionalParam('to', 'Receiver address or account number on the node')
  .addOptionalParam('amount', 'Amount of ETH to send')
  .setAction(async (taskArgs, { ethers }) => {
    const { parseUnits } = ethers.utils;
    const from = await addr(ethers, taskArgs.from, ethers.utils);
    const fromSigner = await ethers.provider.getSigner(from);
    const to = await addr(ethers, taskArgs.to, ethers.utils);
    const value = parseUnits(taskArgs.amount, 'ether').toHexString();
    return fromSigner.sendTransaction({ to, value });
  });
