// Service health check utilities
import { API_BASE_URL } from '../config';

export interface ServiceStatus {
  wave: 'connected' | 'disconnected' | 'checking';
  cloudflare: 'connected' | 'disconnected' | 'checking';
}

class ServiceHealthService {
  private status: ServiceStatus = {
    wave: 'checking',
    cloudflare: 'checking'
  };

  private listeners: ((status: ServiceStatus) => void)[] = [];

  // Subscribe to status changes
  subscribe(callback: (status: ServiceStatus) => void) {
    this.listeners.push(callback);
    // Immediately send current status
    callback(this.status);
    
    // Return unsubscribe function
    return () => {
      this.listeners = this.listeners.filter(cb => cb !== callback);
    };
  }

  // Notify all listeners of status change
  private notifyListeners() {
    this.listeners.forEach(callback => callback(this.status));
  }

  // Check Wave connection
  async checkWaveConnection(): Promise<boolean> {
    this.status.wave = 'checking';
    this.notifyListeners();

    try {
      const response = await fetch(`${API_BASE_URL}/health/wave`, {
        method: 'GET',
      });
      
      const isConnected = response.ok;
      this.status.wave = isConnected ? 'connected' : 'disconnected';
      this.notifyListeners();
      return isConnected;
    } catch (error) {
      console.error('Wave connection check failed:', error);
      this.status.wave = 'disconnected';
      this.notifyListeners();
      return false;
    }
  }

  // Check Cloudflare R2 connection
  async checkCloudflareConnection(): Promise<boolean> {
    this.status.cloudflare = 'checking';
    this.notifyListeners();

    try {
      const response = await fetch(`${API_BASE_URL}/health/cloudflare`, {
        method: 'GET',
      });
      
      const isConnected = response.ok;
      this.status.cloudflare = isConnected ? 'connected' : 'disconnected';
      this.notifyListeners();
      return isConnected;
    } catch (error) {
      console.error('Cloudflare connection check failed:', error);
      this.status.cloudflare = 'disconnected';
      this.notifyListeners();
      return false;
    }
  }

  // Check all services
  async checkAllServices() {
    await Promise.all([
      this.checkWaveConnection(),
      this.checkCloudflareConnection()
    ]);
  }

  // Get current status
  getStatus(): ServiceStatus {
    return { ...this.status };
  }
}

export const serviceHealth = new ServiceHealthService();

// Check services on startup
if (typeof window !== 'undefined') {
  serviceHealth.checkAllServices();
  
  // Recheck every 5 minutes
  setInterval(() => {
    serviceHealth.checkAllServices();
  }, 5 * 60 * 1000);
}
