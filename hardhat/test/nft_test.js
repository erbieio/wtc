const { expect } = require('chai');
const { deploy } = require('./helper');

describe('NFT', () => {
  it('Should deploy Chains Contract', async () => {
    const nft = await deploy('NFT');
    const user1 = '0x0000000000000000000000000000000000000001';
    const user2 = '0x0000000000000000000000000000000000000002';
    const id = 1;
    const uri = 'https://abc';
    describe('superMint', function () {
      it('Should mint one nft', async function () {
        const tx = await nft.superMint(user1, uri);
        await tx.wait();
        expect(await nft.ownerOf(id)).to.equal(user1);
        expect(await nft.tokenURI(id)).to.equal(uri);
        describe('superTransfer', function () {
          it('Should transfer one nft', async function () {
            const tx = await nft.superTransfer(user1, user2, id);
            await tx.wait();
            expect(await nft.ownerOf(id)).to.equal(user2);
          });
        });
      });
    });
  });
});
