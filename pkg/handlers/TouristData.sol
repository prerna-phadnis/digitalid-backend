// TouristRegistry.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TouristRegistry {
    struct TouristData {
        string touristId;   // UUID from Go
        string dataHash;    // SHA-256 of JSON
        uint256 timestamp;
    }

    mapping(string => TouristData) private registry;
    string[] public allIds;

    event TouristRegistered(string touristId, string dataHash, uint256 timestamp);

    function registerTourist(string memory touristId, string memory dataHash) public {
        TouristData memory td = TouristData(touristId, dataHash, block.timestamp);
        registry[touristId] = td;
        allIds.push(touristId);

        emit TouristRegistered(touristId, dataHash, block.timestamp);
    }

    function getTourist(string memory touristId) public view returns (string memory, string memory, uint256) {
        TouristData memory td = registry[touristId];
        return (td.touristId, td.dataHash, td.timestamp);
    }

    function getAllIds() public view returns (string[] memory) {
        return allIds;
    }
}
