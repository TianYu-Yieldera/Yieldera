import { BigInt, Bytes } from "@graphprotocol/graph-ts"
import {
  CampaignCreated,
  Claimed,
  CampaignStatusUpdated
} from "../generated/YielderaAirdrop/YielderaAirdrop"
import {
  AirdropCampaign,
  User,
  Claim,
  DailySnapshot,
  GlobalStats,
  CampaignStatusUpdate
} from "../generated/schema"

// ==================== Helper Functions ====================

/**
 * 获取或创建全局统计
 */
function getOrCreateGlobalStats(): GlobalStats {
  let stats = GlobalStats.load("global")
  if (stats == null) {
    stats = new GlobalStats("global")
    stats.totalCampaigns = 0
    stats.activeCampaigns = 0
    stats.totalDistributed = BigInt.fromI32(0)
    stats.totalClaims = 0
    stats.totalUsers = 0
    stats.updatedAt = BigInt.fromI32(0)
  }
  return stats
}

/**
 * 获取或创建用户
 */
function getOrCreateUser(address: Bytes): User {
  let userId = address.toHex()
  let user = User.load(userId)

  if (user == null) {
    user = new User(userId)
    user.totalClaimed = BigInt.fromI32(0)
    user.claimCount = 0
    user.campaignCount = 0
    user.firstClaimAt = null
    user.lastClaimAt = null
    user.save()

    // 更新全局用户数
    let stats = getOrCreateGlobalStats()
    stats.totalUsers = stats.totalUsers + 1
    stats.save()
  }

  return user
}

/**
 * 获取日期字符串（YYYY-MM-DD）
 */
function getDateString(timestamp: BigInt): string {
  let date = new Date(timestamp.toI64() * 1000)
  let year = date.getUTCFullYear()
  let month = date.getUTCMonth() + 1
  let day = date.getUTCDate()

  let monthStr = month < 10 ? "0" + month.toString() : month.toString()
  let dayStr = day < 10 ? "0" + day.toString() : day.toString()

  return year.toString() + "-" + monthStr + "-" + dayStr
}

/**
 * 获取或创建每日快照
 */
function getOrCreateDailySnapshot(
  campaignId: string,
  dateStr: string,
  timestamp: BigInt
): DailySnapshot {
  let snapshotId = campaignId + "-" + dateStr
  let snapshot = DailySnapshot.load(snapshotId)

  if (snapshot == null) {
    snapshot = new DailySnapshot(snapshotId)
    snapshot.campaign = campaignId
    snapshot.date = dateStr
    snapshot.claimCount = 0
    snapshot.totalAmount = BigInt.fromI32(0)
    snapshot.uniqueClaimers = 0
    snapshot.cumulativeClaimCount = 0
    snapshot.cumulativeAmount = BigInt.fromI32(0)
  }

  return snapshot
}

// ==================== Event Handlers ====================

/**
 * 处理活动创建事件
 */
export function handleCampaignCreated(event: CampaignCreated): void {
  let campaignId = event.params.campaignId.toString()
  let campaign = new AirdropCampaign(campaignId)

  campaign.name = event.params.name
  campaign.description = ""  // 需要从合约读取
  campaign.merkleRoot = Bytes.empty()  // 需要从合约读取
  campaign.totalBudget = event.params.totalBudget
  campaign.claimedAmount = BigInt.fromI32(0)
  campaign.remainingBudget = event.params.totalBudget
  campaign.startTime = event.params.startTime
  campaign.endTime = event.params.endTime
  campaign.isActive = true
  campaign.participantCount = 0
  campaign.claimCount = 0
  campaign.creator = event.params.creator
  campaign.createdAt = event.block.timestamp
  campaign.updatedAt = event.block.timestamp

  campaign.save()

  // 更新全局统计
  let stats = getOrCreateGlobalStats()
  stats.totalCampaigns = stats.totalCampaigns + 1
  stats.activeCampaigns = stats.activeCampaigns + 1
  stats.updatedAt = event.block.timestamp
  stats.save()
}

/**
 * 处理领取事件
 */
export function handleClaimed(event: Claimed): void {
  let campaignId = event.params.campaignId.toString()
  let userAddress = event.params.user
  let amount = event.params.amount

  // 1. 更新活动
  let campaign = AirdropCampaign.load(campaignId)
  if (campaign != null) {
    campaign.claimedAmount = campaign.claimedAmount.plus(amount)
    campaign.remainingBudget = campaign.totalBudget.minus(campaign.claimedAmount)
    campaign.claimCount = campaign.claimCount + 1
    campaign.participantCount = campaign.participantCount + 1
    campaign.updatedAt = event.block.timestamp
    campaign.save()
  }

  // 2. 更新或创建用户
  let user = getOrCreateUser(userAddress)
  let isFirstClaim = user.claimCount == 0

  user.totalClaimed = user.totalClaimed.plus(amount)
  user.claimCount = user.claimCount + 1

  if (isFirstClaim) {
    user.firstClaimAt = event.block.timestamp
    user.campaignCount = 1
  } else {
    // 检查是否是新活动
    let claims = user.claims.load()
    let existingCampaigns = new Set<string>()
    for (let i = 0; i < claims.length; i++) {
      existingCampaigns.add(claims[i].campaign)
    }
    if (!existingCampaigns.has(campaignId)) {
      user.campaignCount = user.campaignCount + 1
    }
  }

  user.lastClaimAt = event.block.timestamp
  user.save()

  // 3. 创建领取记录
  let claimId = event.transaction.hash.toHex() + "-" + event.logIndex.toString()
  let claim = new Claim(claimId)
  claim.campaign = campaignId
  claim.user = user.id
  claim.amount = amount
  claim.timestamp = event.block.timestamp
  claim.blockNumber = event.block.number
  claim.transactionHash = event.transaction.hash
  claim.save()

  // 4. 更新每日快照
  let dateStr = getDateString(event.block.timestamp)
  let snapshot = getOrCreateDailySnapshot(campaignId, dateStr, event.block.timestamp)
  snapshot.claimCount = snapshot.claimCount + 1
  snapshot.totalAmount = snapshot.totalAmount.plus(amount)
  snapshot.uniqueClaimers = snapshot.uniqueClaimers + 1  // TODO: 需要去重

  if (campaign != null) {
    snapshot.cumulativeClaimCount = campaign.claimCount
    snapshot.cumulativeAmount = campaign.claimedAmount
  }

  snapshot.save()

  // 5. 更新全局统计
  let stats = getOrCreateGlobalStats()
  stats.totalDistributed = stats.totalDistributed.plus(amount)
  stats.totalClaims = stats.totalClaims + 1
  stats.updatedAt = event.block.timestamp
  stats.save()
}

/**
 * 处理活动状态更新事件
 */
export function handleCampaignStatusUpdated(event: CampaignStatusUpdated): void {
  let campaignId = event.params.campaignId.toString()
  let isActive = event.params.isActive

  // 1. 更新活动状态
  let campaign = AirdropCampaign.load(campaignId)
  if (campaign != null) {
    let wasActive = campaign.isActive
    campaign.isActive = isActive
    campaign.updatedAt = event.block.timestamp
    campaign.save()

    // 2. 更新全局激活活动数
    let stats = getOrCreateGlobalStats()
    if (wasActive && !isActive) {
      stats.activeCampaigns = stats.activeCampaigns - 1
    } else if (!wasActive && isActive) {
      stats.activeCampaigns = stats.activeCampaigns + 1
    }
    stats.updatedAt = event.block.timestamp
    stats.save()
  }

  // 3. 创建状态变更日志
  let updateId = event.transaction.hash.toHex() + "-" + event.logIndex.toString()
  let update = new CampaignStatusUpdate(updateId)
  update.campaign = campaignId
  update.isActive = isActive
  update.timestamp = event.block.timestamp
  update.blockNumber = event.block.number
  update.transactionHash = event.transaction.hash
  update.save()
}
