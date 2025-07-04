# Cloudflare R2 Public Access Setup

## Getting Your R2 Public URL

1. **Log into Cloudflare Dashboard**
   - Go to R2 > Your Bucket (bremrayskyview)

2. **Enable Public Access**
   - Go to Settings tab
   - Under "Public Access", click "Allow Access"
   - You'll get a public URL like: `https://pub-[hash].r2.dev`

3. **Custom Domain (Optional)**
   - You can also connect a custom domain
   - Go to Settings > Custom Domains
   - Add your domain (e.g., `images.yourdomain.com`)

4. **Update .env**
   - Set `R2_PUBLIC_URL` to your public bucket URL
   - Do NOT include the bucket name in the URL
   - Example: `R2_PUBLIC_URL=https://pub-d1b8db0efacf4b3bb18b9d35baa24b86.r2.dev`

## Important Notes

- The R2 endpoint URL (`R2_ENDPOINT`) is for API access (uploading)
- The R2 public URL (`R2_PUBLIC_URL`) is for public access (viewing)
- These are different URLs!

## Testing

After updating the public URL:
1. Restart the backend
2. Upload a new photo
3. The photo URL should now be accessible in browser

## CORS Configuration

If you're using a custom domain, you may need to configure CORS in Cloudflare:
1. Go to R2 > Your Bucket > Settings
2. Add CORS rules to allow your frontend domain
