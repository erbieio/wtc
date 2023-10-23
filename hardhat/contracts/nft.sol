// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract NFT is ERC721, ERC721URIStorage, Ownable {
    uint256 private _maxId;

    constructor() ERC721("Erbie NFT", "ET") Ownable(msg.sender) {}

    function superMint(address to, string memory uri) public onlyOwner {
        _maxId++;
        _safeMint(to, _maxId);
        _setTokenURI(_maxId, uri);
    }

    function superTransfer(
        address from,
        address to,
        uint256 tokenId
    ) public onlyOwner {
        _safeTransfer(from, to, tokenId);
    }

    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
