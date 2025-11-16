import React, { useEffect, useRef } from 'react';

/**
 * Clean Protocol Network Visualization
 * Simple, professional network diagram showing Yieldera's integrations
 */
export default function ProtocolNetwork3D() {
  const canvasRef = useRef(null);

  // Protocol nodes data - updated to reflect current architecture
  // 增大节点尺寸，建立清晰的视觉层级
  const protocols = [
    { name: 'Yieldera', color: '#22D3EE', size: 90, isCenter: true },  // 中心节点最大
    { name: 'Arbitrum\nDeFi', color: '#3B82F6', size: 65, subtitle: 'High Yield' },  // 次要节点
    { name: 'Base\nTreasury', color: '#10B981', size: 65, subtitle: 'Stable RWA' },
    { name: 'AI Risk\nEngine', color: '#F59E0B', size: 65, subtitle: 'FastAPI' },
    { name: 'Aave V3', color: '#B650A0', size: 52 },  // 协议节点
    { name: 'Compound V3', color: '#00D395', size: 52 },
    { name: 'Uniswap V3', color: '#FF007A', size: 52 },
    { name: 'GMX V2', color: '#1E40AF', size: 52 },
    { name: 'Aerodrome', color: '#A855F7', size: 48, subtitle: 'Base DEX' }
  ];

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d', { alpha: true });
    canvas.width = canvas.offsetWidth;
    canvas.height = canvas.offsetHeight;

    const centerX = canvas.width / 2;
    const centerY = canvas.height / 2;

    // Generate starfield for cosmic background
    const generateStarfield = (count, sizeRange, brightnessRange) => {
      return Array.from({ length: count }, () => ({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        size: Math.random() * (sizeRange[1] - sizeRange[0]) + sizeRange[0],
        brightness: Math.random() * (brightnessRange[1] - brightnessRange[0]) + brightnessRange[0],
        twinkleSpeed: Math.random() * 2 + 1,
        twinklePhase: Math.random() * Math.PI * 2
      }));
    };

    // Multi-layer starfield for depth
    const distantStars = generateStarfield(200, [0.5, 1.2], [0.2, 0.4]);
    const midStars = generateStarfield(100, [1, 2], [0.4, 0.7]);
    const nearStars = generateStarfield(50, [1.5, 3], [0.6, 1]);

    // Calculate protocol node positions (star cluster)
    // 增加半径让节点充分利用整个宇宙空间
    const nodes = protocols.map((protocol, i) => {
      if (protocol.isCenter) {
        return { ...protocol, x: centerX, y: centerY };
      }

      const angle = ((i - 1) / (protocols.length - 1)) * Math.PI * 2;
      const radius = Math.min(centerX, centerY) * 0.82;  // 从0.68增加到0.82，接近边缘
      return {
        ...protocol,
        x: centerX + Math.cos(angle) * radius,
        y: centerY + Math.sin(angle) * radius,
        angle: angle
      };
    });

    let time = 0;
    let animationId;

    const animate = () => {
      time += 0.008;

      // Clear canvas
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      // Draw distant starfield
      distantStars.forEach(star => {
        const twinkle = Math.sin(time * star.twinkleSpeed + star.twinklePhase) * 0.3 + 0.7;
        ctx.fillStyle = `rgba(255, 255, 255, ${star.brightness * twinkle})`;
        ctx.beginPath();
        ctx.arc(star.x, star.y, star.size, 0, Math.PI * 2);
        ctx.fill();
      });

      // Draw mid-distance starfield with slight glow
      midStars.forEach(star => {
        const twinkle = Math.sin(time * star.twinkleSpeed + star.twinklePhase) * 0.4 + 0.6;
        ctx.shadowBlur = star.size * 2;
        ctx.shadowColor = `rgba(173, 216, 255, ${star.brightness * twinkle})`;
        ctx.fillStyle = `rgba(173, 216, 255, ${star.brightness * twinkle})`;
        ctx.beginPath();
        ctx.arc(star.x, star.y, star.size, 0, Math.PI * 2);
        ctx.fill();
        ctx.shadowBlur = 0;
      });

      // Draw near starfield with prominent glow
      nearStars.forEach(star => {
        const twinkle = Math.sin(time * star.twinkleSpeed + star.twinklePhase) * 0.5 + 0.5;
        ctx.shadowBlur = star.size * 3;
        ctx.shadowColor = `rgba(34, 211, 238, ${star.brightness * twinkle})`;
        ctx.fillStyle = `rgba(34, 211, 238, ${star.brightness * twinkle})`;
        ctx.beginPath();
        ctx.arc(star.x, star.y, star.size, 0, Math.PI * 2);
        ctx.fill();
        ctx.shadowBlur = 0;
      });

      // Draw cosmic nebula effect around center
      const nebulaGradient = ctx.createRadialGradient(centerX, centerY, 0, centerX, centerY, Math.min(centerX, centerY) * 0.8);
      nebulaGradient.addColorStop(0, 'rgba(34, 211, 238, 0.12)');
      nebulaGradient.addColorStop(0.4, 'rgba(59, 130, 246, 0.06)');
      nebulaGradient.addColorStop(0.7, 'rgba(147, 51, 234, 0.03)');
      nebulaGradient.addColorStop(1, 'transparent');
      ctx.fillStyle = nebulaGradient;
      ctx.globalAlpha = 0.6 + Math.sin(time) * 0.15;
      ctx.beginPath();
      ctx.arc(centerX, centerY, Math.min(centerX, centerY) * 0.8, 0, Math.PI * 2);
      ctx.fill();
      ctx.globalAlpha = 1;

      // Draw connections (stellar connections)
      nodes.forEach((node, i) => {
        if (!node.isCenter) {
          const centerNode = nodes[0];

          // Connection line with energy flow effect
          const pulseOpacity = 0.2 + Math.sin(time * 2 + i) * 0.15;
          const gradient = ctx.createLinearGradient(node.x, node.y, centerNode.x, centerNode.y);
          gradient.addColorStop(0, `rgba(34, 211, 238, ${pulseOpacity * 0.5})`);
          gradient.addColorStop(0.5, `rgba(34, 211, 238, ${pulseOpacity})`);
          gradient.addColorStop(1, `rgba(34, 211, 238, ${pulseOpacity * 0.5})`);

          ctx.strokeStyle = gradient;
          ctx.lineWidth = 2;
          ctx.shadowBlur = 8;
          ctx.shadowColor = 'rgba(34, 211, 238, 0.6)';
          ctx.beginPath();
          ctx.moveTo(node.x, node.y);
          ctx.lineTo(centerNode.x, centerNode.y);
          ctx.stroke();
          ctx.shadowBlur = 0;

          // Energy particle flowing along connection
          const progress = (time * 0.4 + i * 0.3) % 1;
          const particleX = node.x + (centerNode.x - node.x) * progress;
          const particleY = node.y + (centerNode.y - node.y) * progress;

          ctx.shadowBlur = 15;
          ctx.shadowColor = 'rgba(34, 211, 238, 0.9)';
          ctx.fillStyle = 'rgba(34, 211, 238, 0.95)';
          ctx.beginPath();
          ctx.arc(particleX, particleY, 4, 0, Math.PI * 2);
          ctx.fill();
          ctx.shadowBlur = 0;
        }
      });

      // Draw protocol nodes as celestial bodies
      nodes.forEach(node => {
        // Outer glow (stellar aura)
        const glowGradient = ctx.createRadialGradient(node.x, node.y, 0, node.x, node.y, node.size * 1.8);
        glowGradient.addColorStop(0, node.color + '60');
        glowGradient.addColorStop(0.3, node.color + '30');
        glowGradient.addColorStop(0.6, node.color + '15');
        glowGradient.addColorStop(1, 'transparent');

        ctx.fillStyle = glowGradient;
        ctx.beginPath();
        ctx.arc(node.x, node.y, node.size * 1.8, 0, Math.PI * 2);
        ctx.fill();

        // Core sphere
        const coreGradient = ctx.createRadialGradient(
          node.x - node.size * 0.15,
          node.y - node.size * 0.15,
          0,
          node.x,
          node.y,
          node.size * 0.65
        );
        coreGradient.addColorStop(0, node.isCenter ? 'rgba(255, 255, 255, 0.9)' : 'rgba(255, 255, 255, 0.3)');
        coreGradient.addColorStop(0.5, node.color + 'CC');
        coreGradient.addColorStop(1, node.color + '88');

        ctx.fillStyle = coreGradient;
        ctx.shadowBlur = 20;
        ctx.shadowColor = node.color;
        ctx.beginPath();
        ctx.arc(node.x, node.y, node.size * 0.65, 0, Math.PI * 2);
        ctx.fill();
        ctx.shadowBlur = 0;

        // Outer ring
        ctx.strokeStyle = node.color + 'DD';
        ctx.lineWidth = 2.5;
        ctx.beginPath();
        ctx.arc(node.x, node.y, node.size * 0.65, 0, Math.PI * 2);
        ctx.stroke();

        // Node label with glow - 根据节点大小调整文字
        ctx.shadowBlur = 10;
        ctx.shadowColor = 'rgba(0, 0, 0, 0.9)';
        ctx.fillStyle = 'white';

        // 根据节点类型使用不同字体大小
        const fontSize = node.isCenter ? 18 : (node.size > 60 ? 15 : 13);
        const fontWeight = node.isCenter ? 800 : 700;
        const lineHeight = node.isCenter ? 20 : (node.size > 60 ? 17 : 15);

        ctx.font = `${fontWeight} ${fontSize}px Inter, system-ui, sans-serif`;
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';

        // Handle multiline text
        const lines = node.name.split('\n');
        lines.forEach((line, idx) => {
          ctx.fillText(line, node.x, node.y + (idx - (lines.length - 1) / 2) * lineHeight);
        });
        ctx.shadowBlur = 0;

        // Rotating ring for center node (like Saturn rings)
        if (node.isCenter) {
          const ringSize = node.size * 0.8;
          const rotationAngle = time * 0.5;

          ctx.save();
          ctx.translate(node.x, node.y);
          ctx.rotate(rotationAngle);

          // Ring ellipse
          ctx.beginPath();
          ctx.ellipse(0, 0, ringSize, ringSize * 0.3, 0, 0, Math.PI * 2);
          ctx.strokeStyle = 'rgba(34, 211, 238, 0.4)';
          ctx.lineWidth = 2;
          ctx.stroke();

          ctx.restore();

          // Pulsing ring effect
          const pulseSize = node.size * 0.8 + Math.sin(time * 2) * 8;
          const pulseOpacity = 0.5 - Math.sin(time * 2) * 0.25;
          ctx.strokeStyle = `rgba(34, 211, 238, ${pulseOpacity})`;
          ctx.lineWidth = 3;
          ctx.shadowBlur = 15;
          ctx.shadowColor = 'rgba(34, 211, 238, 0.8)';
          ctx.beginPath();
          ctx.arc(node.x, node.y, pulseSize, 0, Math.PI * 2);
          ctx.stroke();
          ctx.shadowBlur = 0;
        }
      });

      animationId = requestAnimationFrame(animate);
    };

    animate();

    // Handle resize
    const handleResize = () => {
      canvas.width = canvas.offsetWidth;
      canvas.height = canvas.offsetHeight;
    };

    window.addEventListener('resize', handleResize);

    return () => {
      if (animationId) cancelAnimationFrame(animationId);
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  return (
    <div style={{
      position: 'relative',
      width: '100%',
      height: '550px',
      overflow: 'hidden',
      background: 'transparent'
    }}>
      <canvas
        ref={canvasRef}
        style={{
          width: '100%',
          height: '100%',
          position: 'relative',
          zIndex: 1
        }}
      />
    </div>
  );
}
