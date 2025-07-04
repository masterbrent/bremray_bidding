// Photo service for handling photo uploads and downloads
import { API_BASE_URL } from '../config';

export interface UploadedPhoto {
  id: string;
  url: string;
  thumbnail: string;
  name: string;
}

export class PhotoService {
  // Upload photos to backend API which handles R2 upload
  static async uploadPhotos(jobId: string, files: File[]): Promise<UploadedPhoto[]> {
    const formData = new FormData();
    
    // Add all files to the form data
    for (const file of files) {
      formData.append('photos', file);
    }
    
    try {
      const response = await fetch(`${API_BASE_URL}/jobs/${jobId}/photos`, {
        method: 'POST',
        body: formData,
      });
      
      if (!response.ok) {
        throw new Error(`Upload failed: ${response.statusText}`);
      }
      
      const uploadedPhotos = await response.json();
      
      // Transform the response to match our interface
      return uploadedPhotos.map((photo: any) => ({
        id: photo.id,
        url: photo.url,
        thumbnail: photo.url, // Backend doesn't generate thumbnails yet
        name: photo.caption || `Photo ${photo.id.slice(0, 8)}`
      }));
    } catch (error) {
      console.error('Error uploading photos:', error);
      throw error;
    }
  }
  
  // Delete photos via API
  static async deletePhotos(jobId: string, photoIds: string[]): Promise<void> {
    try {
      // Delete each photo
      for (const photoId of photoIds) {
        const response = await fetch(`${API_BASE_URL}/jobs/${jobId}/photos/${photoId}`, {
          method: 'DELETE',
        });
        
        if (!response.ok) {
          console.error(`Failed to delete photo ${photoId}`);
        }
      }
    } catch (error) {
      console.error('Error deleting photos:', error);
      throw error;
    }
  }
  
  // Download photo
  static async downloadPhoto(photo: UploadedPhoto): Promise<void> {
    try {
      const response = await fetch(photo.url);
      const blob = await response.blob();
      
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = photo.name || 'photo.jpg';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(link.href);
    } catch (error) {
      console.error('Error downloading photo:', error);
    }
  }
  
  // Download multiple photos as zip (would use JSZip in production)
  static async downloadPhotos(photos: UploadedPhoto[]): Promise<void> {
    // For now, download individually
    for (const photo of photos) {
      await this.downloadPhoto(photo);
      // Add small delay between downloads
      await new Promise(resolve => setTimeout(resolve, 100));
    }
  }
}