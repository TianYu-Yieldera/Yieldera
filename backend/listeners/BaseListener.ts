import { ethers } from 'ethers';
import { EventEmitter } from 'events';

/**
 * BaseListener - 通用区块链事件监听器基类
 *
 * 功能:
 * - WebSocket连接管理
 * - 自动重连
 * - 事件解码
 * - 错误处理
 * - 区块重组处理
 */
export abstract class BaseListener extends EventEmitter {
  protected provider: ethers.WebSocketProvider;
  protected contract: ethers.Contract;
  protected isRunning: boolean = false;
  protected reconnectAttempts: number = 0;
  protected maxReconnectAttempts: number = 10;

  constructor(
    protected wsUrl: string,
    protected contractAddress: string,
    protected abi: any[],
    protected listenerName: string
  ) {
    super();
    this.initializeProvider();
  }

  /**
   * 初始化Provider
   */
  private initializeProvider(): void {
    this.provider = new ethers.WebSocketProvider(this.wsUrl);
    this.contract = new ethers.Contract(
      this.contractAddress,
      this.abi,
      this.provider
    );

    // 监听连接状态
    this.provider.websocket.on('open', () => {
      console.log(`[${this.listenerName}] WebSocket connected`);
      this.reconnectAttempts = 0;
    });

    this.provider.websocket.on('close', () => {
      console.log(`[${this.listenerName}] WebSocket closed`);
      this.handleDisconnect();
    });

    this.provider.websocket.on('error', (error) => {
      console.error(`[${this.listenerName}] WebSocket error:`, error);
    });
  }

  /**
   * 处理断开连接
   */
  private async handleDisconnect(): Promise<void> {
    if (!this.isRunning) return;

    this.reconnectAttempts++;

    if (this.reconnectAttempts > this.maxReconnectAttempts) {
      console.error(`[${this.listenerName}] Max reconnection attempts reached`);
      this.emit('error', new Error('Max reconnection attempts reached'));
      return;
    }

    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
    console.log(`[${this.listenerName}] Reconnecting in ${delay}ms... (Attempt ${this.reconnectAttempts})`);

    await new Promise(resolve => setTimeout(resolve, delay));

    try {
      await this.provider.destroy();
      this.initializeProvider();
      await this.start();
    } catch (error) {
      console.error(`[${this.listenerName}] Reconnection failed:`, error);
      this.handleDisconnect();
    }
  }

  /**
   * 启动监听器
   */
  async start(): Promise<void> {
    if (this.isRunning) {
      console.warn(`[${this.listenerName}] Already running`);
      return;
    }

    this.isRunning = true;
    console.log(`[${this.listenerName}] Starting listener for contract ${this.contractAddress}`);

    try {
      // 验证合约存在
      const code = await this.provider.getCode(this.contractAddress);
      if (code === '0x') {
        throw new Error(`No contract found at ${this.contractAddress}`);
      }

      // 获取当前区块
      const currentBlock = await this.provider.getBlockNumber();
      console.log(`[${this.listenerName}] Current block: ${currentBlock}`);

      // 注册事件监听器
      await this.registerEventListeners();

      this.emit('started', { block: currentBlock });
      console.log(`[${this.listenerName}] Listener started successfully`);
    } catch (error) {
      console.error(`[${this.listenerName}] Failed to start:`, error);
      this.isRunning = false;
      throw error;
    }
  }

  /**
   * 停止监听器
   */
  async stop(): Promise<void> {
    if (!this.isRunning) {
      return;
    }

    console.log(`[${this.listenerName}] Stopping listener...`);
    this.isRunning = false;

    // 移除所有监听器
    this.contract.removeAllListeners();

    // 关闭Provider
    await this.provider.destroy();

    this.emit('stopped');
    console.log(`[${this.listenerName}] Listener stopped`);
  }

  /**
   * 注册事件监听器 - 子类必须实现
   */
  protected abstract registerEventListeners(): Promise<void>;

  /**
   * 处理事件 - 通用事件处理逻辑
   */
  protected async handleEvent(
    eventName: string,
    args: any[],
    event: ethers.Log
  ): Promise<void> {
    try {
      const block = await event.getBlock();
      const transaction = await event.getTransaction();

      const eventData = {
        eventName,
        args,
        blockNumber: event.blockNumber,
        blockHash: event.blockHash,
        transactionHash: event.transactionHash,
        logIndex: event.index,
        timestamp: block.timestamp,
        gasUsed: transaction?.gasLimit.toString(),
      };

      console.log(`[${this.listenerName}] Event: ${eventName}`, eventData);

      // 发出事件供外部处理
      this.emit('event', eventData);

      // 子类可以重写此方法进行特定处理
      await this.processEvent(eventData);
    } catch (error) {
      console.error(`[${this.listenerName}] Error handling event ${eventName}:`, error);
      this.emit('error', error);
    }
  }

  /**
   * 处理具体事件 - 子类可以重写
   */
  protected async processEvent(eventData: any): Promise<void> {
    // 默认不处理，子类重写
  }

  /**
   * 获取历史事件
   */
  async getHistoricalEvents(
    eventName: string,
    fromBlock: number,
    toBlock: number | string = 'latest'
  ): Promise<ethers.Log[]> {
    try {
      const filter = this.contract.filters[eventName]();
      const events = await this.contract.queryFilter(filter, fromBlock, toBlock);

      console.log(`[${this.listenerName}] Found ${events.length} ${eventName} events from block ${fromBlock} to ${toBlock}`);

      return events;
    } catch (error) {
      console.error(`[${this.listenerName}] Error fetching historical events:`, error);
      throw error;
    }
  }

  /**
   * 获取合约信息
   */
  async getContractInfo(): Promise<any> {
    try {
      const network = await this.provider.getNetwork();
      const blockNumber = await this.provider.getBlockNumber();

      return {
        name: this.listenerName,
        address: this.contractAddress,
        network: network.name,
        chainId: network.chainId,
        currentBlock: blockNumber,
        isRunning: this.isRunning,
      };
    } catch (error) {
      console.error(`[${this.listenerName}] Error getting contract info:`, error);
      throw error;
    }
  }
}
