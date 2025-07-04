# Cloudflare R2 Public Access Setup

## Quick Setup Guide

### Step 1: Enable Public Access on Your R2 Bucket

1. **Log into Cloudflare Dashboard**
   - Go to R2 > Overview
   - Click on your bucket: `bremrayskyview`

2. **Enable Public Access**
   - Click on the "Settings" tab
   - Find "Public access" section
   - Click "Allow public access"
   - Confirm the action

3. **Get Your Public R2.dev URL**
   - After enabling public access, you'll see a URL like:
     - `https://pub-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.r2.dev`
   - Copy this URL (WITHOUT the bucket name at the end)

### Step 2: Update Your Backend Configuration

1. **Edit `/backend/.env`**:
   ```env
   R2_PUBLIC_URL=https://pub-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.r2.dev
   ```
   Replace the xxx with your actual public URL hash

2. **Restart the backend**:
   ```bash
   cd /path/to/electrical-bidding-app
   ./stop-all.sh
   ./start-all.sh
   ```

### Step 3: Test Photo Upload

1. Upload a new photo through the app
2. The photo should now display properly in the gallery

## Alternative: Custom Domain Setup (Optional)

If you want to use your own domain (e.g., `images.yourdomain.com`):

1. **Add Custom Domain in R2**:
   - Go to your bucket > Settings > Custom Domains
   - Click "Connect domain"
   - Enter your subdomain
   - Follow DNS configuration instructions

2. **Update `.env`**:
   ```env
   R2_PUBLIC_URL=https://images.yourdomain.com
   ```

## Troubleshooting

### Photos Still Not Showing?

1. **Check bucket permissions**:
   - Ensure "Allow public access" is enabled
   - Check if there are any bucket policies blocking access

2. **Verify the URL format**:
   - The public URL should NOT include the bucket name
   - Correct: `https://pub-xxx.r2.dev`
   - Wrong: `https://pub-xxx.r2.dev/bremrayskyview`

3. **Clear browser cache**:
   - Hard refresh: Cmd+Shift+R (Mac) or Ctrl+Shift+R (Windows)

4. **Check CORS settings** (if using custom domain):
   - Go to bucket > Settings > CORS
   - Add your frontend domain to allowed origins

### Testing R2 Access

You can test if your R2 bucket is publicly accessible:

```bash
# Test with curl (replace with your actual photo URL)
curl -I "https://pub-xxx.r2.dev/jobs/[job-id]/[timestamp]_[uuid].jpg"
```

Should return `HTTP/1.1 200 OK` if working correctly.

## Security Note

With public access enabled:
- All files in the bucket can be accessed if someone knows the URL
- The UUID in filenames provides some security through obscurity
- Consider implementing signed URLs for sensitive content

## Current Workaround

The app currently generates presigned URLs that expire after 7 days. This works without public access but requires photos to be re-uploaded after expiration.
