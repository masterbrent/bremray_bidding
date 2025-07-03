// Mock photo service for R2 integration
// In production, this would handle actual uploads to Cloudflare R2

export interface UploadedPhoto {
  id: string;
  url: string;
  thumbnail: string;
  name: string;
}

export class PhotoService {
  // Mock upload to R2
  static async uploadPhotos(jobId: string, files: File[]): Promise<UploadedPhoto[]> {
    // In production, this would:
    // 1. Get presigned URLs from backend
    // 2. Upload directly to R2
    // 3. Return the final URLs
    
    const uploadedPhotos: UploadedPhoto[] = [];
    
    for (const file of files) {
      // Simulate upload delay
      await new Promise(resolve => setTimeout(resolve, 500));
      
      // Create object URL for preview (in production, this would be the R2 URL)
      const url = URL.createObjectURL(file);
      
      uploadedPhotos.push({
        id: `photo_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
        url,
        thumbnail: url, // In production, we'd generate actual thumbnails
        name: file.name
      });
    }
    
    return uploadedPhotos;
  }
  
  // Mock delete from R2
  static async deletePhotos(photoIds: string[]): Promise<void> {
    // In production, this would delete from R2
    await new Promise(resolve => setTimeout(resolve, 300));
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