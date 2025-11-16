import React, { useEffect, useRef } from 'react';

/**
 * Clean Protocol Network Visualization
 * Simple, professional network diagram showing Yieldera's integrations
 */
export default function ProtocolNetwork3D() {
  const canvasRef = useRef(null);

  // Protocol nodes data - simplified
  const protocols = [
    { name: 'Yieldera', color: '#22D3EE', size: 70, isCenter: true },
    { name: 'Aave', color: '#B650A0', size: 45 },
    { name: 'Compound', color: '#00D395', size: 45 },
    { name: 'Uniswap', color: '#FF007A', size: 45 },
    { name: 'GMX', color: '#3B82F6', size: 45 },
    { name: 'Treasury', color: '#10B981', size: 45 },
    { name: 'AI Risk', color: '#F59E0B', size: 45 }
  ];

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d', { alpha: true });
    canvas.width = canvas.offsetWidth;
    canvas.height = canvas.offsetHeight;

    const centerX = canvas.width / 2;
    const centerY = canvas.height / 2;

    // Calculate fixed positions
    const nodes = protocols.map((protocol, i) => {
      if (protocol.isCenter) {
        return { ...protocol, x: centerX, y: centerY };
      }

      const angle = ((i - 1) / (protocols.length - 1)) * Math.PI * 2;
      const radius = Math.min(centerX, centerY) * 0.6;
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
      time += 0.01;

      // Clear with transparency
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      // Draw connections first (behind nodes)
      nodes.forEach((node, i) => {
        if (!node.isCenter) {
          const centerNode = nodes[0];

          // Static connection line with subtle pulse
          const pulseOpacity = 0.15 + Math.sin(time + i) * 0.1;
          ctx.strokeStyle = `rgba(34, 211, 238, ${pulseOpacity})`;
          ctx.lineWidth = 2;
          ctx.beginPath();
          ctx.moveTo(node.x, node.y);
          ctx.lineTo(centerNode.x, centerNode.y);
          ctx.stroke();

          // Moving data particle on line
          const progress = (time * 0.5 + i * 0.3) % 1;
          const particleX = node.x + (centerNode.x - node.x) * progress;
          const particleY = node.y + (centerNode.y - node.y) * progress;

          ctx.fillStyle = 'rgba(34, 211, 238, 0.8)';
          ctx.beginPath();
          ctx.arc(particleX, particleY, 3, 0, Math.PI * 2);
          ctx.fill();
        }
      });

      // Draw nodes
      nodes.forEach(node => {
        // Subtle glow
        const gradient = ctx.createRadialGradient(node.x, node.y, 0, node.x, node.y, node.size);
        gradient.addColorStop(0, node.color + '40');
        gradient.addColorStop(0.5, node.color + '20');
        gradient.addColorStop(1, 'transparent');

        ctx.fillStyle = gradient;
        ctx.beginPath();
        ctx.arc(node.x, node.y, node.size * 1.2, 0, Math.PI * 2);
        ctx.fill();

        // Node circle
        ctx.fillStyle = node.isCenter ? node.color + '30' : '#0F172A';
        ctx.strokeStyle = node.color;
        ctx.lineWidth = 3;
        ctx.beginPath();
        ctx.arc(node.x, node.y, node.size * 0.5, 0, Math.PI * 2);
        ctx.fill();
        ctx.stroke();

        // Node label
        ctx.fillStyle = 'white';
        ctx.font = '600 12px Inter, system-ui, sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(node.name, node.x, node.y);

        // Subtle pulse ring for center node only
        if (node.isCenter) {
          const pulseSize = node.size * 0.5 + Math.sin(time * 2) * 5;
          const pulseOpacity = 0.4 - Math.sin(time * 2) * 0.2;
          ctx.strokeStyle = `rgba(34, 211, 238, ${pulseOpacity})`;
          ctx.lineWidth = 2;
          ctx.beginPath();
          ctx.arc(node.x, node.y, pulseSize, 0, Math.PI * 2);
          ctx.stroke();
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
      height: '400px',
      borderRadius: '16px',
      overflow: 'hidden',
      background: 'linear-gradient(135deg, rgba(10, 25, 47, 0.3) 0%, rgba(5, 10, 20, 0.5) 50%, rgba(0, 0, 0, 0.4) 100%)',
      border: '1px solid rgba(34, 211, 238, 0.15)',
      backdropFilter: 'blur(12px)',
      boxShadow: `
        inset 0 1px 0 rgba(255, 255, 255, 0.05),
        0 8px 32px rgba(0, 0, 0, 0.3),
        0 0 40px rgba(34, 211, 238, 0.08)
      `
    }}>
      {/* Glass reflection effect on top edge */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: '10%',
        right: '10%',
        height: '1px',
        background: 'linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent)',
        zIndex: 2
      }} />

      {/* Subtle corner accents */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: '60px',
        height: '60px',
        borderTop: '2px solid rgba(34, 211, 238, 0.2)',
        borderLeft: '2px solid rgba(34, 211, 238, 0.2)',
        borderTopLeftRadius: '16px',
        zIndex: 2
      }} />
      <div style={{
        position: 'absolute',
        top: 0,
        right: 0,
        width: '60px',
        height: '60px',
        borderTop: '2px solid rgba(34, 211, 238, 0.2)',
        borderRight: '2px solid rgba(34, 211, 238, 0.2)',
        borderTopRightRadius: '16px',
        zIndex: 2
      }} />

      <canvas
        ref={canvasRef}
        style={{
          width: '100%',
          height: '100%',
          position: 'relative',
          zIndex: 1
        }}
      />

      {/* Simplified info with enhanced styling */}
      <div style={{
        position: 'absolute',
        bottom: '16px',
        left: '50%',
        transform: 'translateX(-50%)',
        color: 'rgba(255, 255, 255, 0.7)',
        fontSize: '12px',
        fontWeight: '600',
        letterSpacing: '1px',
        textTransform: 'uppercase',
        padding: '6px 16px',
        borderRadius: '20px',
        background: 'rgba(0, 0, 0, 0.3)',
        border: '1px solid rgba(34, 211, 238, 0.2)',
        backdropFilter: 'blur(8px)',
        zIndex: 2,
        textShadow: '0 0 10px rgba(34, 211, 238, 0.5)'
      }}>
        Multi-Protocol Integration
      </div>
    </div>
  );
}
