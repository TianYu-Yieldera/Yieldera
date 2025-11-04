// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IRWAMarketplace
 * @notice Interface for RWA primary and secondary market trading
 * @dev Supports both whole asset and fractional token trading
 */
interface IRWAMarketplace {
    /**
     * @notice Order types in the marketplace
     */
    enum OrderType {
        FixedPrice,      // Simple fixed price sale
        Auction,         // Time-based auction
        DutchAuction,    // Descending price auction
        Offer            // Bid/offer system
    }

    /**
     * @notice Order status
     */
    enum OrderStatus {
        Active,          // Open for trading
        Filled,          // Completely filled
        PartiallyFilled, // Partially filled
        Cancelled,       // Cancelled by creator
        Expired          // Expired due to time
    }

    /**
     * @notice Marketplace listing
     */
    struct Listing {
        uint256 listingId;
        uint256 assetId;
        address seller;
        OrderType orderType;
        OrderStatus status;
        uint256 amount;          // Tokens for sale
        uint256 price;           // Price per token (USD, 18 decimals)
        uint256 minPurchase;     // Minimum purchase amount
        uint256 filled;          // Amount already sold
        uint256 createdAt;
        uint256 expiresAt;
        bool isPrimaryMarket;    // True if initial issuance
    }

    /**
     * @notice Auction-specific data
     */
    struct Auction {
        uint256 startPrice;
        uint256 reservePrice;    // Minimum acceptable price
        uint256 currentBid;
        address currentBidder;
        uint256 bidIncrement;    // Minimum bid increase
        uint256 auctionEndTime;
    }

    /**
     * @notice Trade execution record
     */
    struct Trade {
        uint256 tradeId;
        uint256 listingId;
        uint256 assetId;
        address buyer;
        address seller;
        uint256 amount;
        uint256 price;
        uint256 totalValue;
        uint256 fee;
        uint256 timestamp;
    }

    // Events
    event ListingCreated(
        uint256 indexed listingId,
        uint256 indexed assetId,
        address indexed seller,
        uint256 amount,
        uint256 price,
        OrderType orderType
    );

    event ListingCancelled(
        uint256 indexed listingId,
        address indexed seller
    );

    event TradExecuted(
        uint256 indexed tradeId,
        uint256 indexed listingId,
        address indexed buyer,
        address seller,
        uint256 amount,
        uint256 price,
        uint256 totalValue
    );

    event BidPlaced(
        uint256 indexed listingId,
        address indexed bidder,
        uint256 bidAmount
    );

    event AuctionFinalized(
        uint256 indexed listingId,
        address indexed winner,
        uint256 finalPrice
    );

    // View functions
    function getListing(uint256 listingId) external view returns (Listing memory);
    function getActiveListings(uint256 assetId) external view returns (uint256[] memory);
    function getTradeHistory(uint256 assetId, uint256 count) external view returns (Trade[] memory);
    function calculateFee(uint256 amount, uint256 price) external view returns (uint256);
    function canCreateListing(address seller, uint256 assetId, uint256 amount) external view returns (bool);

    // State-changing functions
    function createListing(
        uint256 assetId,
        uint256 amount,
        uint256 price,
        OrderType orderType,
        uint256 duration,
        uint256 minPurchase
    ) external returns (uint256 listingId);

    function cancelListing(uint256 listingId) external;

    function buyTokens(
        uint256 listingId,
        uint256 amount
    ) external payable returns (uint256 tradeId);

    function placeBid(uint256 listingId, uint256 bidAmount) external payable;

    function finalizeAuction(uint256 listingId) external;
}
