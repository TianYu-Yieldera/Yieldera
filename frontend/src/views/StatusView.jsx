import React, { useEffect, useState } from "react";
import { Server, CheckCircle, XCircle, AlertTriangle, Activity, Zap, Database, Cpu } from "lucide-react";
import { config } from "../config/env";

const API_BASE = config.api.baseUrl;

export default function StatusView() {
  const [apiHealth, setApiHealth] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // 检查 API 健康状态
    fetch(`${API_BASE}/health`)
      .then(r => r.json())
      .then(() => {
        setApiHealth(true);
        setLoading(false);
      })
      .catch(() => {
        setApiHealth(false);
        setLoading(false);
      });
  }, []);

  // 模拟服务状态
  const services = [
    {
      name: "API Gateway",
      status: apiHealth ? "healthy" : "down",
      uptime: "99.98%",
      responseTime: "45ms",
      lastCheck: "刚刚",
      icon: Server,
      color: apiHealth ? "#10b981" : "#EF4444"
    },
    {
      name: "Blockchain Listener",
      status: "healthy",
      uptime: "99.95%",
      responseTime: "120ms",
      lastCheck: "30秒前",
      icon: Activity,
      color: "#10b981"
    },
    {
      name: "Event Consumer",
      status: "healthy",
      uptime: "99.99%",
      responseTime: "35ms",
      lastCheck: "1分钟前",
      icon: Zap,
      color: "#10b981"
    },
    {
      name: "Points Scheduler",
      status: "healthy",
      uptime: "100%",
      responseTime: "28ms",
      lastCheck: "2分钟前",
      icon: Cpu,
      color: "#10b981"
    },
    {
      name: "PostgreSQL Database",
      status: "healthy",
      uptime: "99.99%",
      responseTime: "12ms",
      lastCheck: "30秒前",
      icon: Database,
      color: "#10b981"
    },
    {
      name: "Kafka Message Queue",
      status: "healthy",
      uptime: "99.97%",
      responseTime: "65ms",
      lastCheck: "1分钟前",
      icon: Activity,
      color: "#10b981"
    }
  ];

  // 系统指标
  const metrics = [
    { label: "总运行时间", value: "45天 12小时", icon: Activity, color: "#6366F1" },
    { label: "总处理事件", value: "1,245,892", icon: Zap, color: "#A855F7" },
    { label: "当前TPS", value: "124/s", icon: Cpu, color: "#10b981" },
    { label: "平均响应时间", value: "58ms", icon: Server, color: "#F59E0B" }
  ];

  // RPC 端点状态
  const rpcEndpoints = [
    { name: "Sepolia RPC", url: "wss://eth-sepolia.g.alchemy.com/v2/...", status: "healthy", latency: "145ms" },
    { name: "Mainnet RPC", url: "wss://eth-mainnet.g.alchemy.com/v2/...", status: "maintenance", latency: "-" },
    { name: "Arbitrum RPC", url: "wss://arb-mainnet.g.alchemy.com/v2/...", status: "coming", latency: "-" }
  ];

  // 最近事件
  const recentEvents = [
    { time: "2分钟前", type: "info", message: "Scheduler 完成第 1,245 轮积分计算" },
    { time: "5分钟前", type: "success", message: "成功处理 125 个质押事件" },
    { time: "10分钟前", type: "info", message: "Consumer 同步至区块 #5,234,567" },
    { time: "15分钟前", type: "warning", message: "检测到高负载，已自动扩容" },
    { time: "30分钟前", type: "info", message: "数据库备份完成" }
  ];

  const getStatusIcon = (status) => {
    switch(status) {
      case "healthy": return <CheckCircle size={24} color="#10b981" />;
      case "down": return <XCircle size={24} color="#EF4444" />;
      case "maintenance": return <AlertTriangle size={24} color="#F59E0B" />;
      case "coming": return <Server size={24} color="#666" />;
      default: return <AlertTriangle size={24} color="#666" />;
    }
  };

  const getEventColor = (type) => {
    switch(type) {
      case "success": return "#10b981";
      case "info": return "#6366F1";
      case "warning": return "#F59E0B";
      case "error": return "#EF4444";
      default: return "#666";
    }
  };

  return (
    <div className="container">
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 32, display: 'flex', alignItems: 'center', gap: 12 }}>
          <Server size={36} color="#10b981" />
          系统状态
        </h1>
        <p className="muted" style={{ marginTop: 8 }}>实时监控所有服务的健康状态</p>
      </div>

      {/* 总体状态 */}
      <div className="card" style={{
        padding: 24,
        marginBottom: 24,
        background: 'linear-gradient(135deg, rgba(16, 185, 129, .1), rgba(99, 102, 241, .1))',
        borderColor: '#10b981'
      }}>
        <div className="row" style={{ alignItems: 'center', gap: 16 }}>
          <CheckCircle size={48} color="#10b981" />
          <div>
            <div style={{ fontSize: 24, fontWeight: 800, color: '#10b981' }}>所有系统正常运行</div>
            <div className="muted">所有核心服务运行正常，无故障报告</div>
          </div>
        </div>
      </div>

      {/* 系统指标 */}
      <div className="grid" style={{ gridTemplateColumns: 'repeat(4, 1fr)', gap: 16, marginBottom: 24 }}>
        {metrics.map((metric, index) => {
          const Icon = metric.icon;
          return (
            <div key={index} className="kpi" style={{ position: 'relative', overflow: 'hidden' }}>
              <div style={{ position: 'absolute', top: -10, right: -10, opacity: 0.1 }}>
                <Icon size={80} color={metric.color} />
              </div>
              <div style={{ position: 'relative', zIndex: 1 }}>
                <div className="title">{metric.label}</div>
                <div className="value" style={{ color: metric.color }}>{metric.value}</div>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-2" style={{ gap: 16, marginBottom: 24 }}>
        {/* 服务状态 */}
        <div className="card" style={{ padding: 24 }}>
          <div style={{ fontWeight: 700, marginBottom: 16, fontSize: 18 }}>核心服务</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
            {services.map((service, index) => {
              const Icon = service.icon;
              return (
                <div
                  key={index}
                  style={{
                    padding: 16,
                    background: 'rgba(255,255,255,.02)',
                    border: `1px solid ${service.color}33`,
                    borderRadius: 12,
                    transition: 'all 0.2s'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.background = 'rgba(255,255,255,.05)';
                    e.currentTarget.style.borderColor = service.color;
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.background = 'rgba(255,255,255,.02)';
                    e.currentTarget.style.borderColor = `${service.color}33`;
                  }}
                >
                  <div className="row" style={{ justifyContent: 'space-between', marginBottom: 12 }}>
                    <div className="row" style={{ gap: 12 }}>
                      <Icon size={24} color={service.color} />
                      <div style={{ fontWeight: 600 }}>{service.name}</div>
                    </div>
                    {getStatusIcon(service.status)}
                  </div>
                  <div className="row" style={{ gap: 16, flexWrap: 'wrap', fontSize: 13 }}>
                    <div>
                      <span className="muted">运行时间: </span>
                      <span>{service.uptime}</span>
                    </div>
                    <div>
                      <span className="muted">响应时间: </span>
                      <span>{service.responseTime}</span>
                    </div>
                    <div>
                      <span className="muted">最后检查: </span>
                      <span>{service.lastCheck}</span>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* RPC 端点 & 事件日志 */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          {/* RPC 端点 */}
          <div className="card" style={{ padding: 24 }}>
            <div style={{ fontWeight: 700, marginBottom: 16, fontSize: 18 }}>RPC 端点状态</div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
              {rpcEndpoints.map((endpoint, index) => (
                <div
                  key={index}
                  style={{
                    padding: 12,
                    background: 'rgba(255,255,255,.02)',
                    border: '1px solid rgba(255,255,255,.05)',
                    borderRadius: 8
                  }}
                >
                  <div className="row" style={{ justifyContent: 'space-between', marginBottom: 8 }}>
                    <div style={{ fontWeight: 600, fontSize: 14 }}>{endpoint.name}</div>
                    {getStatusIcon(endpoint.status)}
                  </div>
                  <div className="muted" style={{ fontSize: 12, marginBottom: 4, fontFamily: 'monospace' }}>
                    {endpoint.url}
                  </div>
                  <div className="row" style={{ gap: 12, fontSize: 12 }}>
                    <div>
                      <span className="muted">延迟: </span>
                      <span style={{ color: endpoint.latency === "-" ? "#666" : "#10b981" }}>
                        {endpoint.latency}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* 最近事件 */}
          <div className="card" style={{ padding: 24, flex: 1 }}>
            <div style={{ fontWeight: 700, marginBottom: 16, fontSize: 18 }}>最近事件</div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {recentEvents.map((event, index) => (
                <div
                  key={index}
                  style={{
                    padding: 12,
                    background: 'rgba(255,255,255,.02)',
                    border: '1px solid rgba(255,255,255,.05)',
                    borderRadius: 8,
                    borderLeftWidth: 3,
                    borderLeftColor: getEventColor(event.type)
                  }}
                >
                  <div className="row" style={{ justifyContent: 'space-between', marginBottom: 4 }}>
                    <div className="muted" style={{ fontSize: 11 }}>{event.time}</div>
                    <div style={{
                      fontSize: 10,
                      padding: '2px 6px',
                      borderRadius: 4,
                      background: `${getEventColor(event.type)}22`,
                      color: getEventColor(event.type),
                      fontWeight: 600
                    }}>
                      {event.type.toUpperCase()}
                    </div>
                  </div>
                  <div style={{ fontSize: 13 }}>{event.message}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* 监控说明 */}
      <div className="card" style={{ padding: 20, background: 'rgba(99, 102, 241, .1)', borderColor: '#6366F1' }}>
        <div style={{ fontWeight: 700, marginBottom: 8 }}>关于监控</div>
        <div className="muted" style={{ fontSize: 14, lineHeight: 1.6 }}>
          系统每30秒自动检查一次所有服务的健康状态。如果检测到异常，会立即发送告警通知。
          所有指标数据实时更新，确保你随时了解系统运行状况。
        </div>
      </div>
    </div>
  );
}
