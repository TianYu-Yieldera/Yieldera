// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * @title TreasuryToken
 * @notice ERC20 token representing fractional ownership of US Treasury securities
 * @dev Each token represents a proportional share of the underlying treasury asset
 */
contract TreasuryToken is ERC20, ERC20Burnable, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");

    /// @notice Asset metadata
    uint256 public immutable assetId;
    string public cusip;
    string public maturityTerm;
    uint256 public maturityDate;
    uint256 public couponRate;

    /**
     * @notice Constructor
     * @param name_ Token name (e.g., "US Treasury 10Y Bond")
     * @param symbol_ Token symbol (e.g., "UST-10Y")
     * @param assetId_ Unique asset identifier
     * @param cusip_ CUSIP identifier
     * @param maturityTerm_ Maturity term (e.g., "10Y")
     * @param maturityDate_ Maturity timestamp
     * @param couponRate_ Annual coupon rate in basis points
     * @param factory Factory contract address
     */
    constructor(
        string memory name_,
        string memory symbol_,
        uint256 assetId_,
        string memory cusip_,
        string memory maturityTerm_,
        uint256 maturityDate_,
        uint256 couponRate_,
        address factory
    ) ERC20(name_, symbol_) {
        require(factory != address(0), "Invalid factory");

        assetId = assetId_;
        cusip = cusip_;
        maturityTerm = maturityTerm_;
        maturityDate = maturityDate_;
        couponRate = couponRate_;

        _grantRole(DEFAULT_ADMIN_ROLE, factory);
        _grantRole(MINTER_ROLE, factory);
        _grantRole(BURNER_ROLE, factory);
    }

    /**
     * @notice Mint tokens
     * @param to Recipient address
     * @param amount Amount to mint
     */
    function mint(address to, uint256 amount) external onlyRole(MINTER_ROLE) {
        _mint(to, amount);
    }

    /**
     * @notice Burn tokens from holder
     * @param from Holder address
     * @param amount Amount to burn
     */
    function burnFrom(address from, uint256 amount) public override onlyRole(BURNER_ROLE) {
        _burn(from, amount);
    }

    /**
     * @notice Check if treasury has matured
     * @return True if matured
     */
    function hasMatured() external view returns (bool) {
        return block.timestamp >= maturityDate;
    }

    /**
     * @notice Get asset metadata
     * @return id Asset ID
     * @return cusipId CUSIP identifier
     * @return term Maturity term
     * @return maturity Maturity date timestamp
     * @return coupon Coupon rate in basis points
     * @return supply Total token supply
     */
    function getAssetInfo() external view returns (
        uint256 id,
        string memory cusipId,
        string memory term,
        uint256 maturity,
        uint256 coupon,
        uint256 supply
    ) {
        return (
            assetId,
            cusip,
            maturityTerm,
            maturityDate,
            couponRate,
            totalSupply()
        );
    }
}
