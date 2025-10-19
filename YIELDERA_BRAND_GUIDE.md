# Yieldera 品牌视觉指南

**最后更新：2025-10-18**

---

## 🎨 品牌定位

**Yieldera** 是一个 DeFi 收益聚合 × RWA 资产投资平台，致力于让用户的资产自动增值。

### 核心理念
- **Enter the Yieldera** - 进入收益新时代
- **智能** - 自动优化收益策略
- **无感** - 用户无需操作，系统全自动
- **真实** - 收益可转换为真实世界资产

---

## 🎯 品牌名称

### 英文名称
**Yieldera**
- 发音：/jiːldˈerə/ (yeel-DER-uh)
- 词源：Yield（收益）+ Era（时代）
- 含义：收益的新时代

### 中文名称（备选）
- **盈时代** - 主推
- **益达** - 简洁易记

### Tagline（核心口号）
**Enter the Yieldera**
- 中文：进入收益新时代
- 备选：Your Gateway to Effortless Wealth（通往轻松财富的大门）

---

## 🎨 配色系统

### 主色调（Primary Colors）

#### 品牌紫色（Brand Purple）
```
渐变：#667eea → #764ba2
用途：品牌标识、主要 CTA、重要标题
```

**色值详情：**
- 起始色：`#667eea` (RGB: 102, 126, 234)
- 结束色：`#764ba2` (RGB: 118, 75, 162)
- CSS: `background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);`

#### 辅助色（Secondary Colors）

**粉红渐变**（RWA 商城）
```
#f093fb → #f5576c
用途：RWA 相关功能、次要 CTA
```

**青蓝渐变**（收益展示）
```
#4facfe → #00f2fe
用途：收益数据、正向指标
```

**绿色渐变**（成功状态）
```
#43e97b → #38f9d7
用途：成功消息、正向数据
```

### 中性色（Neutral Colors）

**背景色**
- 深色背景：`#111827`
- 卡片背景：`#1f2937`
- 悬停背景：`#374151`

**文字色**
- 主文字：`#ffffff`
- 次要文字：`#9ca3af`
- 禁用文字：`#6b7280`

**边框色**
- 默认边框：`#374151`
- 聚焦边框：`#667eea`

---

## 📝 排版系统

### 字体家族

**英文字体**
```css
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
             'Helvetica Neue', Arial, sans-serif;
```

**中文字体**
```css
font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC',
             'Microsoft YaHei', 'Hiragino Sans GB', sans-serif;
```

### 字号规范

| 用途 | 字号 | 粗细 | 行高 |
|------|-----|------|------|
| 超大标题（Hero） | 48px | 900 | 1.2 |
| 一级标题（H1） | 32px | 700 | 1.3 |
| 二级标题（H2） | 24px | 600 | 1.4 |
| 三级标题（H3） | 20px | 600 | 1.4 |
| 正文（Body） | 16px | 400 | 1.6 |
| 小字（Small） | 14px | 400 | 1.5 |
| 备注（Caption） | 13px | 400 | 1.4 |
| 按钮（Button） | 16px | 700 | 1 |

### 字重（Font Weight）

- **Regular**: 400（正文）
- **Medium**: 500（强调）
- **Semi-Bold**: 600（标题）
- **Bold**: 700（按钮、重要标题）
- **Black**: 900（超大标题）

---

## 🎭 Logo 设计指南

### Logo 元素

**图标（Mark）**
- 形状：向上箭头 + 圆形背景
- 颜色：品牌紫色渐变
- 含义：收益增长、向上突破

**文字（Wordmark）**
- 字体：粗体无衬线字体
- 颜色：渐变文字效果（仅在深色背景）
- 字母间距：-1px

### Logo 用法

**标准版（深色背景）**
```jsx
<div style={{
  fontWeight: 700,
  fontSize: '24px',
  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  WebkitBackgroundClip: 'text',
  WebkitTextFillColor: 'transparent'
}}>
  Yieldera
</div>
```

**简化版（小尺寸）**
```jsx
<div style={{
  fontWeight: 700,
  fontSize: '16px',
  color: '#667eea'
}}>
  Yieldera
</div>
```

### Logo 最小尺寸
- 数字平台：高度不小于 24px
- 印刷品：高度不小于 12mm

### Logo 留白区域
- 周围留白至少为 Logo 高度的 25%

---

## 🧩 UI 组件规范

### 按钮（Buttons）

#### 主要按钮（Primary）
```css
padding: 14px 32px;
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
border-radius: 8px;
color: #ffffff;
font-weight: 700;
font-size: 16px;
```

#### 次要按钮（Secondary）
```css
padding: 14px 32px;
background: rgba(102, 126, 234, 0.1);
border: 2px solid #667eea;
border-radius: 8px;
color: #667eea;
font-weight: 700;
```

#### 幽灵按钮（Ghost）
```css
padding: 14px 32px;
background: transparent;
border: 2px solid rgba(255,255,255,0.2);
border-radius: 8px;
color: #ffffff;
font-weight: 700;
```

### 卡片（Cards）

#### 标准卡片
```css
background: #1f2937;
border: 1px solid #374151;
border-radius: 12px;
padding: 24px;
```

#### 高亮卡片
```css
background: linear-gradient(135deg, rgba(102,126,234,.15), rgba(118,75,162,.15));
border: 2px solid #667eea;
border-radius: 12px;
padding: 24px;
```

#### 渐变卡片（Hero）
```css
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
border-radius: 16px;
padding: 48px 32px;
```

### 输入框（Inputs）

```css
padding: 12px 16px;
background: #374151;
border: 1px solid #4b5563;
border-radius: 8px;
color: #ffffff;
font-size: 16px;

/* Focus state */
border-color: #667eea;
outline: none;
box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
```

### 徽章（Badges）

#### 状态徽章
```css
padding: 4px 12px;
background: #667eea;
border-radius: 12px;
font-size: 12px;
font-weight: 700;
color: #ffffff;
```

#### 标签徽章
```css
padding: 6px 16px;
background: rgba(255,255,255,0.1);
border-radius: 20px;
font-size: 13px;
font-weight: 600;
color: #ffffff;
backdrop-filter: blur(10px);
```

---

## 🎬 动效规范

### 过渡时间（Transition Duration）

- **快速**：150ms - 按钮悬停、颜色变化
- **标准**：300ms - 卡片悬停、展开收起
- **慢速**：500ms - 页面切换、模态框

### 缓动函数（Easing）

```css
/* 标准缓动 */
transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);

/* 弹性缓动 */
transition-timing-function: cubic-bezier(0.68, -0.55, 0.265, 1.55);
```

### 悬停效果

**卡片悬停**
```css
transform: translateY(-5px);
box-shadow: 0 10px 20px rgba(0,0,0,0.3);
border-color: #667eea;
```

**按钮悬停**
```css
transform: scale(1.02);
box-shadow: 0 4px 14px rgba(102, 126, 234, 0.4);
```

---

## 📐 间距系统

使用 8px 基础单位的倍数：

| 名称 | 值 | 用途 |
|------|---|------|
| xs | 4px | 紧密间距 |
| sm | 8px | 小间距 |
| md | 16px | 标准间距 |
| lg | 24px | 大间距 |
| xl | 32px | 区块间距 |
| 2xl | 48px | 大区块间距 |
| 3xl | 64px | 超大间距 |

---

## 📱 响应式断点

```css
/* 手机 */
@media (max-width: 640px) { ... }

/* 平板 */
@media (min-width: 641px) and (max-width: 1024px) { ... }

/* 桌面 */
@media (min-width: 1025px) { ... }
```

---

## 🖼️ 图标规范

### 图标库
- **主要使用**：Lucide React
- **风格**：线性图标，2px 描边
- **尺寸**：16px、20px、24px

### 图标颜色
- 默认：`#ffffff`
- 次要：`#9ca3af`
- 强调：`#667eea`

### 常用图标

| 功能 | 图标 | 组件名 |
|------|------|--------|
| 理财金库 | 🏦 | Vault |
| RWA 资产 | 💎 | Gem |
| 收益趋势 | 📈 | TrendingUp |
| 闪电快速 | ⚡ | Zap |
| 钱包 | 👛 | Wallet |
| 排行榜 | 🏆 | Trophy |

---

## 🎯 品牌语言

### 语气与风格

**核心原则：**
- 专业但不呆板
- 自信但不傲慢
- 创新但不浮夸
- 友好但不随意

**推荐用词：**
- ✅ 智能、自动、优化、收益、资产、策略
- ✅ 无感知、轻松、高效、透明、安全
- ❌ 避免：暴富、赌博、保证、零风险

### 口号示例

**主口号**
- Enter the Yieldera
- 进入收益新时代

**功能口号**
- One Vault, All Yields（一个金库，所有收益）
- Smart Money, Smarter Yields（聪明的钱，更聪明的收益）
- From DeFi to Real Assets（从 DeFi 到真实资产）

**营销口号**
- Your Assets, Auto-Piloted（你的资产，自动驾驶）
- Sleep Well, Earn Well（睡得好，赚得好）
- The Future of Wealth（财富的未来）

---

## 📊 使用场景

### 1. 网站首页
- **Hero 区域**：大标题 Yieldera + Tagline + 渐变背景
- **CTA**：白色按钮（主要）+ 幽灵按钮（次要）
- **数据卡片**：使用不同渐变色区分类型

### 2. 理财金库页面
- **主色调**：品牌紫色渐变
- **统计卡片**：4 种不同渐变（紫、粉、蓝、绿）
- **投资策略**：使用 4 种协议对应的品牌色

### 3. RWA 商城页面
- **主色调**：粉红渐变
- **资产卡片**：深色背景，悬停高亮
- **分类标签**：不同颜色区分资产类型

### 4. 社交媒体
- **头像**：品牌 Logo 带渐变
- **封面**：Yieldera 大标题 + Enter the Yieldera
- **配图**：使用品牌渐变色作为背景

---

## ✅ 使用规范检查清单

设计新功能时，确保：

- [ ] 使用品牌紫色渐变作为主色调
- [ ] 所有标题使用正确字重（H1: 700, H2: 600）
- [ ] 按钮使用标准内边距（14px 32px）
- [ ] 卡片圆角统一使用 12px
- [ ] 间距使用 8px 的倍数
- [ ] 动画过渡时间为 300ms
- [ ] 悬停效果包含 transform + box-shadow
- [ ] 文案符合品牌语气
- [ ] 图标来自 Lucide React
- [ ] 响应式设计考虑移动端

---

## 🚀 未来扩展

### Logo 动画（计划中）
- 加载动画：向上箭头逐渐填充
- 悬停动画：渐变色流动效果

### 品牌插画（计划中）
- 风格：扁平化、几何感、渐变色
- 主题：收益增长、资产管理、智能自动化

### 品牌音效（计划中）
- 成功音：清脆愉悦
- 错误音：柔和提示
- 背景音：科技感氛围音乐

---

**品牌所有权：Yieldera Protocol**
**文档维护：产品设计团队**
**版本：v2.0.0**
