# Setting Up Public Photo Bucket

## Steps to Create a New Public Bucket for Photos

### 1. Create New R2 Bucket
1. Log into Cloudflare Dashboard
2. Go to R2 > Overview
3. Click "Create bucket"
4. Name it: `bremray-photos-public` (or your preferred name)
5. Keep default settings and create

### 2. Enable Public Access
1. Click on the new bucket
2. Go to "Settings" tab
3. Find "Public access" section
4. Click "Allow public access"
5. Confirm the action
6. Copy the public URL (looks like: `https://pub-xxxxxxxxxxxxxx.r2.dev`)

### 3. Update Backend Configuration
1. Edit `/backend/.env`:
```env
# Keep existing private bucket for backups
R2_BUCKET_NAME=bremrayskyview

# Add new public bucket for photos
R2_PUBLIC_BUCKET_NAME=bremray-photos-public
R2_PUBLIC_BUCKET_URL=https://pub-xxxxxxxxxxxxxx.r2.dev
```

### 4. Migrate Existing Photos (Optional)
If you want to move existing photos to the new bucket:

```bash
# You can use rclone or the Cloudflare dashboard to copy files
# from bremrayskyview/jobs/* to bremray-photos-public/jobs/*
```

### 5. Update Database URLs
After migrating photos, update the database:

```sql
-- Update photo URLs to use the new public bucket URL
UPDATE job_photos 
SET url = REPLACE(url, 'old-bucket-url', 'https://pub-xxxxxxxxxxxxxx.r2.dev');
```

### 6. Restart Backend
```bash
cd /path/to/electrical-bidding-app/backend
./stop-all.sh
./start-all.sh
```

## Benefits of This Approach
- ✅ Database backups remain private in the original bucket
- ✅ Photos get permanent public URLs that never expire
- ✅ No need for presigned URLs
- ✅ Better performance (direct access, no API calls needed)
- ✅ Clear separation of concerns (backups vs public assets)

## Security Note
The public bucket will make all photos accessible to anyone with the URL. The random UUIDs in filenames provide security through obscurity, but this is suitable for job photos that aren't highly sensitive.
