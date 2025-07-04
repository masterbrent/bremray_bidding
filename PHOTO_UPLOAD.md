services/photoService.ts`)
   - Handles photo upload to backend API
   - Manages photo deletion
   - Provides download functionality

3. **Integration**
   - Job detail page shows photo gallery
   - Upload progress indication
   - Error handling with user feedback

## Configuration

### Environment Variables

Add these to your backend `.env` file:

```env
# Cloudflare R2 Configuration
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_access_key_id
R2_SECRET_ACCESS_KEY=your_secret_access_key
R2_BUCKET_NAME=your_bucket_name
R2_ENDPOINT=https://your_account_id.r2.cloudflarestorage.com
R2_PUBLIC_URL=https://your_public_url
```

### Cloudflare R2 Setup

1. **Create a bucket** in your Cloudflare dashboard
2. **Configure public access** for the bucket
3. **Generate API tokens** with read/write permissions
4. **Set up a custom domain** (optional) for public access

## Usage

### Taking Photos

1. Navigate to a job detail page
2. Click "Take Photos" button
3. Choose between:
   - **Use Camera**: Take photos directly
   - **Upload Photos**: Select from device

4. Preview photos before saving
5. Click "Add Photos" to upload

### Managing Photos

- **View**: Photos appear in the job's photo gallery
- **Download**: Click download icon on any photo
- **Delete**: Click delete icon (admin only)

## Security

- Photos are uploaded through the backend API
- Backend validates job existence before upload
- File size limits enforced (32MB max per request)
- Only image files accepted (JPEG, PNG, GIF, WebP)

## Future Enhancements

1. **Thumbnail Generation**: Create smaller versions for gallery
2. **Image Compression**: Optimize images before upload
3. **Bulk Operations**: Select multiple photos for download/delete
4. **Photo Annotations**: Add notes or markings to photos
5. **Photo Categories**: Organize photos by type (before/after, electrical panel, etc.)

## Troubleshooting

### Upload Fails

1. Check backend logs for R2 connection errors
2. Verify environment variables are set correctly
3. Ensure Cloudflare R2 bucket has proper permissions

### Photos Not Displaying

1. Check if R2 public URL is correct
2. Verify CORS settings allow frontend domain
3. Check browser console for loading errors

### Camera Not Working

1. Ensure HTTPS is used (required for camera access)
2. Check browser permissions for camera
3. Try different browser if issues persist
