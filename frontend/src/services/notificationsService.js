/**
 * Notifications Service Client
 * Manages user notifications and preferences
 */

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class NotificationsService {
  constructor() {
    this.baseURL = `${API_BASE}/api/v1/notifications`;
  }

  /**
   * Get user notifications
   * @param {string} userId - User wallet address
   * @param {number} limit - Number of notifications to fetch
   * @returns {Promise<Object>} User notifications
   */
  async getUserNotifications(userId, limit = 20) {
    try {
      const response = await fetch(`${this.baseURL}/${userId}?limit=${limit}`);

      if (!response.ok) {
        throw new Error(`Failed to fetch notifications: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Notifications fetch error:', error);
      throw error;
    }
  }

  /**
   * Mark notification as read
   * @param {string} userId - User wallet address
   * @param {string} timestamp - Notification timestamp
   * @returns {Promise<Object>} Update result
   */
  async markAsRead(userId, timestamp) {
    try {
      const response = await fetch(`${this.baseURL}/${userId}/${timestamp}/read`, {
        method: 'POST',
      });

      if (!response.ok) {
        throw new Error(`Failed to mark as read: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Mark as read error:', error);
      throw error;
    }
  }

  /**
   * Get user notification preferences
   * @param {string} userId - User wallet address
   * @returns {Promise<Object>} User preferences
   */
  async getPreferences(userId) {
    try {
      const response = await fetch(`${this.baseURL}/${userId}/preferences`);

      if (!response.ok) {
        throw new Error(`Failed to fetch preferences: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Preferences fetch error:', error);
      throw error;
    }
  }

  /**
   * Update user notification preferences
   * @param {string} userId - User wallet address
   * @param {Object} preferences - New preferences
   * @returns {Promise<Object>} Update result
   */
  async updatePreferences(userId, preferences) {
    try {
      const response = await fetch(`${this.baseURL}/${userId}/preferences`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(preferences),
      });

      if (!response.ok) {
        throw new Error(`Failed to update preferences: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Preferences update error:', error);
      throw error;
    }
  }

  /**
   * Get unread notification count
   * @param {Array} notifications - Array of notifications
   * @returns {number} Count of unread notifications
   */
  getUnreadCount(notifications) {
    if (!notifications || !Array.isArray(notifications)) return 0;
    return notifications.filter(n => !n.read).length;
  }

  /**
   * Format notification timestamp
   * @param {string} timestamp - ISO timestamp
   * @returns {string} Formatted time ago string
   */
  formatTimeAgo(timestamp) {
    const now = new Date();
    const notificationTime = new Date(timestamp);
    const diffMs = now - notificationTime;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;

    return notificationTime.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    });
  }

  /**
   * Get notification icon based on type
   * @param {string} type - Notification type
   * @returns {string} Icon name for lucide-react
   */
  getNotificationIcon(type) {
    const iconMap = {
      'yield': 'TrendingUp',
      'risk': 'AlertTriangle',
      'trade': 'ShoppingCart',
      'system': 'Settings',
      'alert': 'AlertCircle',
      'success': 'CheckCircle',
      'info': 'Info',
    };

    return iconMap[type] || 'Bell';
  }

  /**
   * Get notification color based on severity
   * @param {string} severity - Notification severity (low, medium, high, critical)
   * @returns {Object} Color scheme
   */
  getNotificationColor(severity) {
    const colorMap = {
      'low': { bg: 'rgb(239, 246, 255)', border: 'rgb(191, 219, 254)', text: 'rgb(29, 78, 216)' },
      'medium': { bg: 'rgb(254, 249, 195)', border: 'rgb(253, 224, 71)', text: 'rgb(234, 179, 8)' },
      'high': { bg: 'rgb(254, 243, 199)', border: 'rgb(251, 191, 36)', text: 'rgb(249, 115, 22)' },
      'critical': { bg: 'rgb(254, 226, 226)', border: 'rgb(252, 165, 165)', text: 'rgb(220, 38, 38)' },
    };

    return colorMap[severity] || colorMap['low'];
  }
}

// Export singleton instance
const notificationsService = new NotificationsService();
export default notificationsService;
