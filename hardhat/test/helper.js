const { ethers } = require('hardhat');
const chalk = require('chalk');

const ContractProxy = (name, instance) => {
  const funcs = Object.keys(instance.interface);
  return new Proxy(
    { instance },
    {
      get: (target, property) => {
        if (funcs.includes(property)) {
          return async (...args) => {
            // Capture and print transaction gas consumption
            const tx = await target.instance[property](...args);
            if (tx.hash) {
              const { gasUsed } = await tx.wait();
              console.log(`${name}.${property} hash:${tx.hash} gasUsed:`, chalk.magenta(gasUsed.toString()));
            }
            return tx;
          };
        }
        if (property === 'connect') {
          return (...args) => ContractProxy(name, target.instance.connect(...args));
        }
        return target.instance[property];
      },
    },
  );
};

exports.deploy = async (name, ...args) => {
  const instance = await ethers.deployContract(name, args);
  console.log(`deploy contract ${name} address:`, chalk.magenta(instance.target));
  return ContractProxy(name, instance);
};
