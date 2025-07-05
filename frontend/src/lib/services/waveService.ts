import { API_BASE_URL } from '../config';

export class WaveService {
  static async sendToWave(jobId: string): Promise<{
    invoiceNumber: string;
    invoiceUrl: string;
    message: string;
  }> {
    const response = await fetch(`${API_BASE_URL}/jobs/${jobId}/send-to-wave`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || 'Failed to send to Wave');
    }

    return response.json();
  }
}
