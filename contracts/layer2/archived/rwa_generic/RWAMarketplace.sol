// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/rwa/IRWAMarketplace.sol";
import "../interfaces/rwa/IRWAAsset.sol";
import "../interfaces/rwa/IRWACompliance.sol";
import "../interfaces/rwa/IRWAValuation.sol";

/**
 * @title RWAMarketplace
 * @notice Primary and secondary market for RWA fractional tokens
 * @dev Supports multiple order types with compliance checks
 *
 * Key Features:
 * - Primary market: Initial token issuance from asset issuer
 * - Secondary market: Peer-to-peer trading between investors
 * - Multiple order types: Fixed price, Auction, Dutch auction, Offers
 * - Compliance-checked trades (KYC, jurisdiction, accreditation)
 * - Fee collection (configurable fee rate)
 * - Trade history and analytics
 *
 * Trading Flow:
 * 1. Seller creates listing via createListing()
 * 2. Buyer checks compliance via canInvestInAsset()
 * 3. Buyer purchases via buyTokens() or placeBid()
 * 4. Compliance checked, tokens transferred, fees collected
 * 5. Trade recorded in history
 */
contract RWAMarketplace is IRWAMarketplace, AccessControl, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    bytes32 public constant MARKETPLACE_ADMIN_ROLE = keccak256("MARKETPLACE_ADMIN_ROLE");
    bytes32 public constant FEE_COLLECTOR_ROLE = keccak256("FEE_COLLECTOR_ROLE");

    // Listing counter
    uint256 private listingIdCounter;
    uint256 private tradeIdCounter;

    // Listings storage
    mapping(uint256 => Listing) private listings;
    mapping(uint256 => Auction) private auctions;

    // Trade history
    Trade[] private tradeHistory;
    mapping(uint256 => uint256[]) private assetTrades; // assetId => tradeIds

    // Integration contracts
    IRWAAsset public immutable assetFactory;
    IRWACompliance public immutable complianceContract;
    IRWAValuation public immutable valuationContract;

    // Fee configuration
    uint256 public feeRate = 250; // 2.5% (basis points)
    uint256 public constant MAX_FEE_RATE = 1000; // 10% maximum
    uint256 public constant FEE_PRECISION = 10000;
    address public feeCollector;

    // Accumulated fees
    mapping(address => uint256) public collectedFees; // token => amount

    // Statistics
    uint256 public totalListings;
    uint256 public totalTrades;
    uint256 public totalVolume; // Total USD volume

    /**
     * @notice Constructor
     * @param admin Admin address
     * @param factory RWAAssetFactory address
     * @param compliance RWACompliance address
     * @param valuation RWAValuation address
     * @param feeCollector_ Fee collector address
     */
    constructor(
        address admin,
        address factory,
        address compliance,
        address valuation,
        address feeCollector_
    ) {
        require(admin != address(0), "Invalid admin");
        require(factory != address(0), "Invalid factory");
        require(compliance != address(0), "Invalid compliance");
        require(valuation != address(0), "Invalid valuation");
        require(feeCollector_ != address(0), "Invalid fee collector");

        assetFactory = IRWAAsset(factory);
        complianceContract = IRWACompliance(compliance);
        valuationContract = IRWAValuation(valuation);
        feeCollector = feeCollector_;

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(MARKETPLACE_ADMIN_ROLE, admin);
        _grantRole(FEE_COLLECTOR_ROLE, admin);
    }

    // =============================================================
    //                     LISTING MANAGEMENT
    // =============================================================

    /**
     * @notice Create marketplace listing
     * @param assetId Asset identifier
     * @param amount Number of tokens to sell
     * @param price Price per token in USD (18 decimals)
     * @param orderType Type of order (FixedPrice, Auction, etc.)
     * @param duration Listing duration in seconds
     * @param minPurchase Minimum purchase amount
     * @return listingId New listing ID
     */
    function createListing(
        uint256 assetId,
        uint256 amount,
        uint256 price,
        OrderType orderType,
        uint256 duration,
        uint256 minPurchase
    ) external override whenNotPaused nonReentrant returns (uint256 listingId) {
        require(assetFactory.isAssetActive(assetId), "Asset not active");
        require(amount > 0, "Invalid amount");
        require(price > 0, "Invalid price");
        require(duration > 0 && duration <= 365 days, "Invalid duration");
        require(minPurchase <= amount, "Min purchase > amount");

        // Get fractional token
        address tokenAddress = assetFactory.getFractionalToken(assetId);
        require(tokenAddress != address(0), "Asset not fractionalized");

        IERC20 token = IERC20(tokenAddress);

        // Check seller compliance
        require(
            complianceContract.canInvestInAsset(msg.sender, assetId),
            "Seller not compliant"
        );

        // Check seller has enough tokens
        require(token.balanceOf(msg.sender) >= amount, "Insufficient balance");

        // Transfer tokens to marketplace (escrow)
        token.safeTransferFrom(msg.sender, address(this), amount);

        listingIdCounter++;
        listingId = listingIdCounter;

        Listing storage listing = listings[listingId];
        listing.listingId = listingId;
        listing.assetId = assetId;
        listing.seller = msg.sender;
        listing.orderType = orderType;
        listing.status = OrderStatus.Active;
        listing.amount = amount;
        listing.price = price;
        listing.minPurchase = minPurchase;
        listing.filled = 0;
        listing.createdAt = block.timestamp;
        listing.expiresAt = block.timestamp + duration;
        listing.isPrimaryMarket = false; // Can add logic to detect primary market

        // Set up auction if needed
        if (orderType == OrderType.Auction || orderType == OrderType.DutchAuction) {
            Auction storage auction = auctions[listingId];
            auction.startPrice = price;
            auction.reservePrice = price * 80 / 100; // Default 80% reserve
            auction.currentBid = 0;
            auction.currentBidder = address(0);
            auction.bidIncrement = price / 100; // 1% minimum increment
            auction.auctionEndTime = block.timestamp + duration;
        }

        totalListings++;

        emit ListingCreated(listingId, assetId, msg.sender, amount, price, orderType);
    }

    /**
     * @notice Cancel active listing
     * @param listingId Listing to cancel
     */
    function cancelListing(uint256 listingId) external override nonReentrant {
        Listing storage listing = listings[listingId];

        require(listing.seller == msg.sender, "Not seller");
        require(listing.status == OrderStatus.Active, "Not active");
        require(block.timestamp < listing.expiresAt, "Listing expired");

        // Calculate remaining tokens
        uint256 remaining = listing.amount - listing.filled;
        require(remaining > 0, "Listing fully filled");

        // Return escrowed tokens
        address tokenAddress = assetFactory.getFractionalToken(listing.assetId);
        IERC20(tokenAddress).safeTransfer(msg.sender, remaining);

        listing.status = OrderStatus.Cancelled;

        emit ListingCancelled(listingId, msg.sender);
    }

    // =============================================================
    //                        TRADING
    // =============================================================

    /**
     * @notice Buy tokens from listing
     * @param listingId Listing to buy from
     * @param amount Amount of tokens to buy
     * @return tradeId Trade identifier
     */
    function buyTokens(
        uint256 listingId,
        uint256 amount
    ) external payable override whenNotPaused nonReentrant returns (uint256 tradeId) {
        Listing storage listing = listings[listingId];

        require(listing.status == OrderStatus.Active, "Listing not active");
        require(block.timestamp < listing.expiresAt, "Listing expired");
        require(amount >= listing.minPurchase, "Below min purchase");
        require(amount <= listing.amount - listing.filled, "Insufficient tokens");

        // Only fixed price for now (auction uses placeBid)
        require(listing.orderType == OrderType.FixedPrice, "Use placeBid for auctions");

        // Check buyer compliance
        require(
            complianceContract.canInvestInAsset(msg.sender, listing.assetId),
            "Buyer not compliant"
        );

        // Calculate cost and fee
        uint256 totalCost = (amount * listing.price) / 1e18;
        uint256 fee = calculateFee(amount, listing.price);
        uint256 sellerReceives = totalCost - fee;

        require(msg.value >= totalCost, "Insufficient payment");

        // Transfer tokens to buyer
        address tokenAddress = assetFactory.getFractionalToken(listing.assetId);
        IERC20(tokenAddress).safeTransfer(msg.sender, amount);

        // Transfer payment to seller
        payable(listing.seller).transfer(sellerReceives);

        // Collect fee
        collectedFees[address(0)] += fee; // ETH fees (address(0) represents ETH)

        // Update listing
        listing.filled += amount;
        if (listing.filled == listing.amount) {
            listing.status = OrderStatus.Filled;
        } else {
            listing.status = OrderStatus.PartiallyFilled;
        }

        // Record trade
        tradeIdCounter++;
        tradeId = tradeIdCounter;

        Trade memory trade = Trade({
            tradeId: tradeId,
            listingId: listingId,
            assetId: listing.assetId,
            buyer: msg.sender,
            seller: listing.seller,
            amount: amount,
            price: listing.price,
            totalValue: totalCost,
            fee: fee,
            timestamp: block.timestamp
        });

        tradeHistory.push(trade);
        assetTrades[listing.assetId].push(tradeId);

        totalTrades++;
        totalVolume += totalCost;

        // Refund excess payment
        if (msg.value > totalCost) {
            payable(msg.sender).transfer(msg.value - totalCost);
        }

        emit TradExecuted(
            tradeId,
            listingId,
            msg.sender,
            listing.seller,
            amount,
            listing.price,
            totalCost
        );
    }

    // =============================================================
    //                       AUCTION LOGIC
    // =============================================================

    /**
     * @notice Place bid on auction
     * @param listingId Auction listing
     * @param bidAmount Bid amount in USD (18 decimals)
     */
    function placeBid(
        uint256 listingId,
        uint256 bidAmount
    ) external payable override nonReentrant {
        Listing storage listing = listings[listingId];
        Auction storage auction = auctions[listingId];

        require(listing.status == OrderStatus.Active, "Auction not active");
        require(listing.orderType == OrderType.Auction, "Not an auction");
        require(block.timestamp < auction.auctionEndTime, "Auction ended");

        // Check buyer compliance
        require(
            complianceContract.canInvestInAsset(msg.sender, listing.assetId),
            "Bidder not compliant"
        );

        // Validate bid
        if (auction.currentBid == 0) {
            require(bidAmount >= auction.startPrice, "Bid below start price");
        } else {
            require(
                bidAmount >= auction.currentBid + auction.bidIncrement,
                "Bid increment too small"
            );
        }

        require(bidAmount >= auction.reservePrice, "Bid below reserve");
        require(msg.value >= bidAmount, "Insufficient payment");

        // Refund previous bidder
        if (auction.currentBidder != address(0)) {
            payable(auction.currentBidder).transfer(auction.currentBid);
        }

        // Update auction
        auction.currentBid = bidAmount;
        auction.currentBidder = msg.sender;

        // Refund excess
        if (msg.value > bidAmount) {
            payable(msg.sender).transfer(msg.value - bidAmount);
        }

        emit BidPlaced(listingId, msg.sender, bidAmount);
    }

    /**
     * @notice Finalize auction after end time
     * @param listingId Auction to finalize
     */
    function finalizeAuction(uint256 listingId) external override nonReentrant {
        Listing storage listing = listings[listingId];
        Auction storage auction = auctions[listingId];

        require(listing.orderType == OrderType.Auction, "Not an auction");
        require(block.timestamp >= auction.auctionEndTime, "Auction not ended");
        require(listing.status == OrderStatus.Active, "Already finalized");

        // Get token address
        address tokenAddress = assetFactory.getFractionalToken(listing.assetId);

        // Check if reserve met
        if (auction.currentBid < auction.reservePrice) {
            // Auction failed - return tokens to seller
            IERC20(tokenAddress).safeTransfer(listing.seller, listing.amount);

            listing.status = OrderStatus.Cancelled;
            return;
        }

        // Calculate fee
        uint256 fee = (auction.currentBid * feeRate) / FEE_PRECISION;
        uint256 sellerReceives = auction.currentBid - fee;

        // Transfer tokens to winner
        IERC20(tokenAddress).safeTransfer(auction.currentBidder, listing.amount);

        // Transfer payment to seller
        payable(listing.seller).transfer(sellerReceives);

        // Collect fee
        collectedFees[address(0)] += fee;

        // Update listing
        listing.status = OrderStatus.Filled;
        listing.filled = listing.amount;

        // Record trade
        tradeIdCounter++;
        uint256 tradeId = tradeIdCounter;

        Trade memory trade = Trade({
            tradeId: tradeId,
            listingId: listingId,
            assetId: listing.assetId,
            buyer: auction.currentBidder,
            seller: listing.seller,
            amount: listing.amount,
            price: auction.currentBid,
            totalValue: auction.currentBid,
            fee: fee,
            timestamp: block.timestamp
        });

        tradeHistory.push(trade);
        assetTrades[listing.assetId].push(tradeId);

        totalTrades++;
        totalVolume += auction.currentBid;

        emit AuctionFinalized(listingId, auction.currentBidder, auction.currentBid);
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get listing details
     * @param listingId Listing identifier
     * @return Listing struct
     */
    function getListing(uint256 listingId) external view override returns (Listing memory) {
        return listings[listingId];
    }

    /**
     * @notice Get active listings for asset
     * @param assetId Asset identifier
     * @return Array of active listing IDs
     */
    function getActiveListings(uint256 assetId) external view override returns (uint256[] memory) {
        uint256 count = 0;

        // Count active listings
        for (uint256 i = 1; i <= listingIdCounter; i++) {
            if (listings[i].assetId == assetId &&
                listings[i].status == OrderStatus.Active &&
                block.timestamp < listings[i].expiresAt) {
                count++;
            }
        }

        // Build array
        uint256[] memory activeListings = new uint256[](count);
        uint256 index = 0;

        for (uint256 i = 1; i <= listingIdCounter; i++) {
            if (listings[i].assetId == assetId &&
                listings[i].status == OrderStatus.Active &&
                block.timestamp < listings[i].expiresAt) {
                activeListings[index] = i;
                index++;
            }
        }

        return activeListings;
    }

    /**
     * @notice Get trade history for asset
     * @param assetId Asset identifier
     * @param count Number of recent trades to return
     * @return Array of Trade structs
     */
    function getTradeHistory(
        uint256 assetId,
        uint256 count
    ) external view override returns (Trade[] memory) {
        uint256[] storage tradeIds = assetTrades[assetId];
        uint256 returnCount = count > tradeIds.length ? tradeIds.length : count;

        Trade[] memory trades = new Trade[](returnCount);

        for (uint256 i = 0; i < returnCount; i++) {
            uint256 tradeIndex = tradeIds[tradeIds.length - returnCount + i] - 1;
            trades[i] = tradeHistory[tradeIndex];
        }

        return trades;
    }

    /**
     * @notice Calculate trading fee
     * @param amount Token amount
     * @param price Price per token
     * @return Fee amount
     */
    function calculateFee(uint256 amount, uint256 price) public view override returns (uint256) {
        uint256 totalValue = (amount * price) / 1e18;
        return (totalValue * feeRate) / FEE_PRECISION;
    }

    /**
     * @notice Check if seller can create listing
     * @param seller Seller address
     * @param assetId Asset identifier
     * @param amount Amount to sell
     * @return True if can create listing
     */
    function canCreateListing(
        address seller,
        uint256 assetId,
        uint256 amount
    ) external view override returns (bool) {
        // Check asset active
        if (!assetFactory.isAssetActive(assetId)) {
            return false;
        }

        // Check compliance
        if (!complianceContract.canInvestInAsset(seller, assetId)) {
            return false;
        }

        // Check balance
        address tokenAddress = assetFactory.getFractionalToken(assetId);
        if (tokenAddress == address(0)) {
            return false;
        }

        IERC20 token = IERC20(tokenAddress);
        if (token.balanceOf(seller) < amount) {
            return false;
        }

        return true;
    }

    /**
     * @notice Get auction details
     * @param listingId Listing identifier
     * @return Auction struct
     */
    function getAuction(uint256 listingId) external view returns (Auction memory) {
        return auctions[listingId];
    }

    // =============================================================
    //                     ADMIN FUNCTIONS
    // =============================================================

    /**
     * @notice Update fee rate
     * @param newFeeRate New fee rate in basis points
     */
    function setFeeRate(uint256 newFeeRate) external onlyRole(MARKETPLACE_ADMIN_ROLE) {
        require(newFeeRate <= MAX_FEE_RATE, "Fee too high");
        feeRate = newFeeRate;
    }

    /**
     * @notice Update fee collector address
     * @param newFeeCollector New fee collector
     */
    function setFeeCollector(address newFeeCollector) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newFeeCollector != address(0), "Invalid address");
        feeCollector = newFeeCollector;
    }

    /**
     * @notice Withdraw collected fees
     * @param token Token address (address(0) for ETH)
     */
    function withdrawFees(address token) external onlyRole(FEE_COLLECTOR_ROLE) {
        uint256 amount = collectedFees[token];
        require(amount > 0, "No fees to withdraw");

        collectedFees[token] = 0;

        if (token == address(0)) {
            payable(feeCollector).transfer(amount);
        } else {
            IERC20(token).safeTransfer(feeCollector, amount);
        }
    }

    /**
     * @notice Pause marketplace
     */
    function pause() external onlyRole(MARKETPLACE_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause marketplace
     */
    function unpause() external onlyRole(MARKETPLACE_ADMIN_ROLE) {
        _unpause();
    }

    /**
     * @notice Receive ETH
     */
    receive() external payable {}
}
